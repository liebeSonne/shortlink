package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"

	"github.com/liebeSonne/shortlink/internal/auth"
	"github.com/liebeSonne/shortlink/internal/config"
	"github.com/liebeSonne/shortlink/internal/handler"
	"github.com/liebeSonne/shortlink/internal/handler/audit"
	handlerauth "github.com/liebeSonne/shortlink/internal/handler/auth"
	"github.com/liebeSonne/shortlink/internal/handler/compress"
	"github.com/liebeSonne/shortlink/internal/handler/cookie"
	grpchandler "github.com/liebeSonne/shortlink/internal/handler/grpc"
	"github.com/liebeSonne/shortlink/internal/handler/subnet"
	internalio "github.com/liebeSonne/shortlink/internal/io"
	applogger "github.com/liebeSonne/shortlink/internal/logger"
	"github.com/liebeSonne/shortlink/internal/repository"
	"github.com/liebeSonne/shortlink/internal/repository/database"
	"github.com/liebeSonne/shortlink/internal/repository/filestorage"
	"github.com/liebeSonne/shortlink/internal/repository/memory"
	"github.com/liebeSonne/shortlink/internal/service"
)

const grpcHandlerAuthTokenKey = "authorization"

type dependencyContainer struct {
	HTTPHandler                  http.Handler
	ShortenerGRPCServer          *grpchandler.ShortenerGRPCServer
	ShortenerGRPCAuthInterceptor grpc.UnaryServerInterceptor
}

func newDependencyContainer(
	ctx context.Context,
	cfg config.Config,
	logger applogger.Logger,
	closer *internalio.MultiCloser,
) (*dependencyContainer, error) {
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
	closer.AddCloser(internalio.CloserFunc(
		func() error {
			return shortLinkDeleter.Stop()
		},
	))

	auditPublisher, err := initAuditPublisher(ctx, cfg, logger, closer)
	if err != nil {
		return nil, fmt.Errorf("error initializing audit publisher: %w", err)
	}

	shortLinkHandler := handler.NewShortLinkHandler(shortLinkService, shortLinkRepository, cfg.BaseURL, shortLinkDeleter, logger, auditPublisher)
	db := createDatabase(cfg)

	statsHandler := subnet.NewTrustedSubnetMiddleware(shortLinkHandler.HandleInternalStats, cfg.TrustedSubnet, logger)

	databaseHandler := handler.NewDatabaseHandler(db, logger)
	rootRouter := handler.NewRootRouter(shortLinkHandler, databaseHandler, cfg.EnableLogs, statsHandler)

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

	shortenerServer := grpchandler.NewShortenerGRPCServer(shortLinkService, shortLinkRepository, cfg.BaseURL, logger, auditPublisher)
	shortenerGRPCAuthInterceptor := grpchandler.AuthUnaryInterceptor(grpcHandlerAuthTokenKey, tokenService, logger)

	return &dependencyContainer{
		HTTPHandler:                  router,
		ShortenerGRPCServer:          shortenerServer,
		ShortenerGRPCAuthInterceptor: shortenerGRPCAuthInterceptor,
	}, nil
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
