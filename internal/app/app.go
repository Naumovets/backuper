package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Naumovets/Backuper/internal/api"
	"github.com/Naumovets/Backuper/internal/config"
	clog "github.com/Naumovets/Backuper/internal/logger"
	"github.com/Naumovets/Backuper/internal/scheduler"
	"github.com/Naumovets/Backuper/internal/services"
	"go.uber.org/zap"
)

type app struct {
	config *config.Config
}

func NewApp() *app {
	config, err := config.Load(os.Getenv("config_path"))
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	return &app{
		config: config,
	}
}

func (a *app) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger, err := clog.InitLogger(a.config.Backuper.Logging)
	if err != nil {
		log.Fatalf("cannot init logger: %s\n", err)
	}

	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	ctxWithLogger := context.WithValue(ctx, clog.CtxKeyLogger, logger)

	servicer := services.NewService(a.config)

	srv := api.InitRest(servicer)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	sch := scheduler.NewScheduler(a.config)

	sch.Run(ctxWithLogger)

	<-ctx.Done()
	stop()

	logger.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: ", zap.Error(err))
	}

	logger.Info("Server exiting")
}
