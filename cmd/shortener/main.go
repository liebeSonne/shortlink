package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/config"
	"github.com/liebeSonne/shortlink/internal/handler"
	handlerauth "github.com/liebeSonne/shortlink/internal/handler/auth"
	"github.com/liebeSonne/shortlink/internal/handler/compress"
	"github.com/liebeSonne/shortlink/internal/handler/cookie"
	internalio "github.com/liebeSonne/shortlink/internal/io"
	applogger "github.com/liebeSonne/shortlink/internal/logger"
	"github.com/liebeSonne/shortlink/internal/repository"
	"github.com/liebeSonne/shortlink/internal/repository/database"
	"github.com/liebeSonne/shortlink/internal/repository/filestorage"
	"github.com/liebeSonne/shortlink/internal/repository/memory"
	"github.com/liebeSonne/shortlink/internal/service"
)

const appID = "shortlink"
const envPrefix = ""

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	closer := internalio.MultiCloser{}
	defer func() {
		closeErr := closer.Close()
		if closeErr != nil {
			log.Fatalf("error closing closer: %v", closeErr)
		}
	}()

	cfg := initConfig()
	logger := initLogger(cfg, &closer)

	err := runMigrator(cfg)
	if err != nil {
		logger.Fatalw("error run migrator", "error", err)
	}

	err = runApp(ctx, cfg, logger, &closer)
	if err != nil {
		logger.Fatalw("error starting server", "error", err)
	}
}

func runApp(
	ctx context.Context,
	cfg config.Config,
	logger applogger.Logger,
	closer *internalio.MultiCloser,
) error {
	router, err := initRouter(ctx, cfg, logger, closer)
	if err != nil {
		return err
	}

	logger.Infow("starting server",
		"addr", cfg.ServerAddress,
		"baseURL", cfg.BaseURL,
		"logLevel", cfg.LogLevel,
		"logFile", cfg.LogFile,
		"storage", cfg.FileStoragePath,
	)

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	serverErrors := make(chan error, 1)

	go func() {
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
		logger.Infow("starting server shutdown")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			srv.Close()
			return err
		}
		logger.Infow("server shutdown complete")
	}

	return nil
}

func initRouter(
	ctx context.Context,
	cfg config.Config,
	logger applogger.Logger,
	closer *internalio.MultiCloser,
) (http.Handler, error) {
	dbClient, err := initDatabaseClient(ctx, cfg, closer)
	if err != nil {
		return nil, fmt.Errorf("error initializing database client: %w", err)
	}

	tokenService := auth.NewTokenService(cfg.AuthSecretKey, cfg.AuthTokenExpires)
	cookieService := cookie.NewService(cfg.AuthCookieTokenKey)
	userService := service.NewUserService()

	shortLinkRepository, err := initShortLinkRepository(cfg, closer, dbClient)
	if err != nil {
		return nil, fmt.Errorf("error initializing short link repository: %w", err)
	}
	shortIDGenerator := service.NewShortIDGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortIDGenerator, service.DefaultMaxAttemptsToGenerateUniqueID)
	shortLinkDeleter := service.NewShortLinkDeleter(ctx, logger, func(input service.InputDelete) error {
		return shortLinkService.DeleteIDs(ctx, input.IDs, input.UserID)
	})
	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, shortLinkRepository, cfg.BaseURL, shortLinkDeleter, logger)
	db := createDatabase(cfg)

	databaseHandler := handler.NewDatabaseHandler(db, logger)
	rootRouter := handler.NewRootRouter(shortLinkHandler, databaseHandler, cfg.EnableLogs)

	router := rootRouter.Router().(http.Handler)

	router = handlerauth.NewAuthMiddleware(router, tokenService, cookieService, logger)
	router = cookie.NewAuthCookieMiddleware(router, tokenService, cookieService, userService, logger)

	router, err = compress.NewCompressorMiddleware(router, compress.CompressorConfig{
		Encodings:    []compress.Encoding{compress.GzipEncoding},
		ContentTypes: &[]string{"application/json", "text/html"},
	})
	if err != nil {
		return nil, err
	}

	router = handler.LoggingMiddleware(router, logger)

	return router, nil
}

var configToLoggerLogLevelMap = map[string]applogger.LogLevel{
	config.LogLevelDebug: applogger.DebugLevel,
	config.LogLevelInfo:  applogger.InfoLevel,
	config.LogLevelWarn:  applogger.WarnLevel,
	config.LogLevelError: applogger.ErrorLevel,
	config.LogLevelFatal: applogger.FatalLevel,
	config.LogLevelPanic: applogger.PanicLevel,
}

func initConfig() config.Config {
	cfg, err := config.LoadConfig(appID, envPrefix)
	if err != nil {
		log.Fatalf("error get config: %s", err.Error())
	}
	return cfg
}

func initLogger(cfg config.Config, closer *internalio.MultiCloser) applogger.Logger {
	loggerLevel, ok := configToLoggerLogLevelMap[cfg.LogLevel]
	if !ok {
		log.Fatalf("unknown log level: %s", cfg.LogLevel)
	}

	logWriter := initLogWriter(cfg, closer)

	logger, err := applogger.NewZapLogger(loggerLevel, logWriter)
	if err != nil {
		log.Fatalf("error init logger: %s", err.Error())
	}
	return logger
}

func initLogWriter(cfg config.Config, closer *internalio.MultiCloser) io.Writer {
	if cfg.LogFile != nil && *cfg.LogFile != "" {
		file, err := os.OpenFile(*cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		if closer != nil {
			closer.AddCloser(internalio.CloserFunc(
				func() error {
					return file.Close()
				},
			))
		}

		return file
	}

	return os.Stderr
}

func initShortLinkRepository(
	cfg config.Config,
	closer *internalio.MultiCloser,
	dbClient *database.Client,
) (repository.ShortLinkRepository, error) {
	if dbClient != nil {
		repo := database.NewShortLinkRepository((*dbClient).Pool())
		return repo, nil
	}

	if cfg.FileStoragePath != nil && *cfg.FileStoragePath != "" {
		repo, err := crateFileShortLinkRepository(*cfg.FileStoragePath, closer)
		if err != nil {
			return nil, err
		}
		return repo, nil
	}

	return memory.NewMemoryShortLinkRepository(), nil
}

func crateFileShortLinkRepository(
	fileStoragePath string,
	closer *internalio.MultiCloser,
) (repository.ShortLinkRepository, error) {
	repo, err := filestorage.NewFileShortLinkRepository(fileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("error on init file short link repository: %w", err)
	}

	if closer != nil {
		closer.AddCloser(internalio.CloserFunc(
			func() error {
				return repo.Close()
			},
		))
	}

	return repo, nil
}
