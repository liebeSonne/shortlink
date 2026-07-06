package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/config"
	"github.com/liebeSonne/shortlink/internal/handler"
	"github.com/liebeSonne/shortlink/internal/handler/audit"
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

//	@title						Shortener API
//	@description				Сервис сокращенных ссылок.
//	@version					1.0
//	@host						localhost:8080
//	@SecurityDefinitions.apikey	cookieAuth
//	@In							cookie
//	@Name						session_token

const appID = "shortlink"
const envPrefix = ""
const gracefulShutdownTimeout = 10 * time.Second

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	closer := internalio.MultiCloser{}
	defer func() {
		closeErr := closer.Close()
		if closeErr != nil {
			log.Fatalf("error closing closer: %v", closeErr)
		}
	}()

	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("error initializing config: %v", err)
	}

	logger, err := initLogger(cfg, &closer)
	if err != nil {
		log.Fatalf("error initializing logger: %v", err)
	}

	err = runMigrator(cfg)
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
		var err error
		if cfg.EnableHTTPS {
			err = srv.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
		gracefulShutdown(srv, logger)
	}

	return nil
}

func gracefulShutdown(
	srv *http.Server,
	logger applogger.Logger,
) {
	logger.Infow("starting server shutdown")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		logger.Errorw("shutdown server error", "error", err)

		// дополнительное время при ошибке
		forceCtx, forceCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer forceCancel()

		<-forceCtx.Done()

		err = srv.Close()
		if err != nil {
			logger.Errorw("close server error", "error", err)
		}
	}

	logger.Infow("server shutdown complete")
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

	shortLinkRepository, err := initShortLinkRepository(cfg, closer, dbClient, logger)
	if err != nil {
		return nil, fmt.Errorf("error initializing short link repository: %w", err)
	}
	shortIDGenerator := service.NewShortIDGenerator()
	shortLinkService := service.NewShortLinkService(shortLinkRepository, shortIDGenerator, service.DefaultMaxAttemptsToGenerateUniqueID)
	shortLinkDeleter := service.NewShortLinkDeleter(ctx, logger, func(input service.InputDelete) error {
		return shortLinkService.DeleteIDs(ctx, input.IDs, input.UserID)
	})

	auditPublisher, err := initAuditPublisher(ctx, cfg, logger, closer)
	if err != nil {
		return nil, fmt.Errorf("error initializing audit publisher: %w", err)
	}

	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, shortLinkRepository, cfg.BaseURL, shortLinkDeleter, logger, auditPublisher)
	db := createDatabase(cfg)

	databaseHandler := handler.NewDatabaseHandler(db, logger)
	rootRouter := handler.NewRootRouter(shortLinkHandler, databaseHandler, cfg.EnableLogs)

	router := rootRouter.Router().(http.Handler)

	router = handlerauth.NewAuthMiddleware(router, tokenService, cookieService, logger)
	router = cookie.NewAuthCookieMiddleware(router, tokenService, cookieService, userService, logger)

	router, err = compress.NewCompressorMiddleware(router, compress.CompressorConfig{
		Encodings:    []compress.Encoding{compress.GzipEncoding},
		ContentTypes: &[]string{"application/json", "text/html"},
	}, logger)
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

func initConfig() (config.Config, error) {
	cfg, err := config.LoadConfig(appID, envPrefix)
	if err != nil {
		return config.Config{}, fmt.Errorf("error get config: %s", err.Error())
	}
	return cfg, nil
}

func initLogger(cfg config.Config, closer *internalio.MultiCloser) (applogger.Logger, error) {
	loggerLevel, ok := configToLoggerLogLevelMap[cfg.LogLevel]
	if !ok {
		return nil, fmt.Errorf("unknown log level: %s", cfg.LogLevel)
	}

	logWriter, err := initLogWriter(cfg, closer)
	if err != nil {
		return nil, fmt.Errorf("error initializing log writer: %w", err)
	}

	logger, err := applogger.NewZapLogger(loggerLevel, logWriter)
	if err != nil {
		return nil, fmt.Errorf("error init logger: %w", err)
	}
	return logger, nil
}

func initLogWriter(cfg config.Config, closer *internalio.MultiCloser) (io.Writer, error) {
	if cfg.LogFile != nil && *cfg.LogFile != "" {
		file, err := os.OpenFile(*cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("error opening log file: %w", err)
		}

		if closer != nil {
			closer.AddCloser(internalio.CloserFunc(
				func() error {
					return file.Close()
				},
			))
		}

		return file, nil
	}

	return os.Stderr, nil
}

func initShortLinkRepository(
	cfg config.Config,
	closer *internalio.MultiCloser,
	dbClient *database.Client,
	logger applogger.Logger,
) (repository.ShortLinkRepository, error) {
	if dbClient != nil {
		repo := database.NewShortLinkRepository((*dbClient).Pool(), logger)
		return repo, nil
	}

	if cfg.FileStoragePath != nil && *cfg.FileStoragePath != "" {
		repo, err := crateFileShortLinkRepository(*cfg.FileStoragePath, closer, logger)
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
	logger applogger.Logger,
) (repository.ShortLinkRepository, error) {
	repo, err := filestorage.NewFileShortLinkRepository(fileStoragePath, logger)
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

func initAuditPublisher(
	ctx context.Context,
	cfg config.Config,
	logger applogger.Logger,
	closer *internalio.MultiCloser,
) (audit.Publisher, error) {
	auditPublisher := audit.NewPublisher(ctx, 5, 5, logger)

	if cfg.AuditFile != nil && *cfg.AuditFile != "" {
		auditFileObserver, err := audit.NewFileObserver(*cfg.AuditFile, logger)
		if err != nil {
			return nil, fmt.Errorf("error initializing audit file observer: %w", err)
		}
		if closer != nil {
			closer.AddCloser(internalio.CloserFunc(
				func() error {
					return auditFileObserver.Close()
				},
			))
		}
		auditPublisher.Subscribe(auditFileObserver)
	}

	if cfg.AuditURL != nil && *cfg.AuditURL != "" {
		auditURLObserver := audit.NewURLObserver(*cfg.AuditURL, 3, time.Second*3, logger)
		auditPublisher.Subscribe(auditURLObserver)
	}

	return auditPublisher, nil
}
