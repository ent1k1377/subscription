package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"subscriptions/internal/database/postgres"
	"subscriptions/internal/database/postgres/repository"
	"subscriptions/internal/service"
	myhttp "subscriptions/internal/transport/http"
	"subscriptions/internal/transport/http/handler/subscription"
	"syscall"
)

type App struct {
	server *myhttp.Server
	db     *postgres.DB
}

func New() *App {
	// TODO Load Config
	pool, err := postgres.GetConnection()
	if err != nil {
		panic(err)
	}
	
	db := postgres.NewDB(pool)
	subscriptionRepo := repository.NewSubscription(pool)
	subscriptionService := service.NewSubscription(subscriptionRepo)
	subscriptionHandler := subscription.NewHandler(subscriptionService)
	
	server := myhttp.NewServer(subscriptionHandler)
	
	return &App{
		server: server,
		db:     db,
	}
}

func (a *App) Run() {
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
