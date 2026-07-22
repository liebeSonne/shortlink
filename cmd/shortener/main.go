package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/liebeSonne/shortlink/api/proto"
	"github.com/liebeSonne/shortlink/internal/config"
	internalio "github.com/liebeSonne/shortlink/internal/io"
	applogger "github.com/liebeSonne/shortlink/internal/logger"
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
	dependency, err := newDependencyContainer(ctx, cfg, logger, closer)
	if err != nil {
		return fmt.Errorf("error initializing dependency conteiter: %v", err)
	}

	logger.Infow("starting server",
		"addr", cfg.ServerAddress,
		"grpcAddr", cfg.GRPCServerAddress,
		"baseURL", cfg.BaseURL,
		"logLevel", cfg.LogLevel,
		"logFile", cfg.LogFile,
		"storage", cfg.FileStoragePath,
		"trustedSubnet", cfg.TrustedSubnet,
	)

	// GRPC
	grpcListen, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		return fmt.Errorf("could not listen grpc on address '%s': %w", cfg.GRPCServerAddress, err)
	}

	grpcOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(dependency.ShortenerGRPCAuthInterceptor),
	}

	if cfg.EnableHTTPS {
		cert, err := tls.LoadX509KeyPair(cfg.TLSCertFile, cfg.TLSKeyFile)
		if err != nil {
			return fmt.Errorf("error on loading tls cert, key: %w", err)
		}
		tlsCredentials := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		grpcOptions = append(grpcOptions, grpc.Creds(tlsCredentials))
	}

	grpcServer := grpc.NewServer(grpcOptions...)
	pb.RegisterShortenerServiceServer(grpcServer, dependency.ShortenerGRPCServer)

	// HTTP
	httpServer := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: dependency.HTTPHandler,
	}

	serverErrors := make(chan error, 1)

	go func() {
		var err error
		if cfg.EnableHTTPS {
			err = httpServer.ListenAndServeTLS(cfg.TLSCertFile, cfg.TLSKeyFile)
		} else {
			err = httpServer.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	go func() {
		err := grpcServer.Serve(grpcListen)
		if err != nil {
			serverErrors <- err
		}
	}()

	select {
	case err := <-serverErrors:
		return err
	case <-ctx.Done():
		gracefulShutdown(httpServer, grpcServer, logger)
	}

	return nil
}

func gracefulShutdown(
	srv *http.Server,
	grpcSrv *grpc.Server,
	logger applogger.Logger,
) {
	logger.Infow("starting server shutdown")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()

	grpcSrv.GracefulStop()

	err := srv.Shutdown(shutdownCtx)
	if err != nil {
		logger.Errorw("shutdown server error", "error", err)

		// дополнительное время при ошибке
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()

		done := make(chan error, 1)
		go func() {
			done <- srv.Close()
		}()

		select {
		case err := <-done:
			logger.Errorw("close server error", "error", err)
		case <-timer.C:
			logger.Errorw("force shutdown timeout")
		}
	}

	logger.Infow("server shutdown complete")
}
