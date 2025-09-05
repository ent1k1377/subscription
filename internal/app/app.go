package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/ent1k1377/subscriptions/internal/config"
	"github.com/ent1k1377/subscriptions/internal/database/postgres"
	"github.com/ent1k1377/subscriptions/internal/database/postgres/repository"
	"github.com/ent1k1377/subscriptions/internal/service"
	myhttp "github.com/ent1k1377/subscriptions/internal/transport/http"
	"github.com/ent1k1377/subscriptions/internal/transport/http/handler/subscription"
)

type App struct {
	server *myhttp.Server
	db     *postgres.DB
	logger *slog.Logger
}

func New() *App {
	cfg := config.MustLoadConfig()

	baseLogger := setupLogger(cfg.LoggerConfig.Level)
	slog.SetDefault(baseLogger)
	baseLogger.Info("Application configuration initialized", slog.Any("cfg", cfg))
	baseLogger.Info("Initialized application")

	pool, err := postgres.GetConnection(cfg.DatabaseConfig)
	if err != nil {
		panic(err)
	}
	baseLogger.Info("Successful connection to the database")

	db := postgres.NewDB(pool)
	subscriptionRepo := repository.NewSubscription(pool, baseLogger)
	subscriptionService := service.NewSubscription(baseLogger, subscriptionRepo)
	subscriptionHandler := subscription.NewHandler(baseLogger, subscriptionService)

	server := myhttp.NewServer(cfg.ServerConfig, baseLogger, subscriptionHandler)

	return &App{
		server: server,
		db:     db,
		logger: baseLogger,
	}
}

func (a *App) Run() {
	a.logger.Info("Starting server")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			cancel()
		}
	}()

	<-ctx.Done()

	_ = a.server.Close(context.Background())
	// TODO лог ошибки, да и вообще надо получше сделать shutdown
	a.db.Close()
}
