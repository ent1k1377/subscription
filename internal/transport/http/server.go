package http

import (
	"context"
	"net/http"
	"subscriptions/internal/transport/http/handler/subscription"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer          *http.Server
	engine              *gin.Engine
	subscriptionHandler *subscription.Handler
}

func NewServer(subscriptionHandler *subscription.Handler) *Server {
	engine := gin.Default()
	httpServer := &http.Server{
		Addr:    ":8080", // TODO Добавить Config
		Handler: engine,
	}

	return &Server{
		httpServer:          httpServer,
		engine:              engine,
		subscriptionHandler: subscriptionHandler,
	}
}

func (s *Server) Start() error {
	s.SetRoutes()

	return s.httpServer.ListenAndServe()
}

func (s *Server) Close(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) SetRoutes() {
	api := s.engine.Group("/api/subscriptions")
	{
		api.POST("/create", s.subscriptionHandler.Create)
		api.GET("/:uuid", s.subscriptionHandler.GetSubscription)
		api.PUT("/:uuid", s.subscriptionHandler.UpdateSubscription)
		api.DELETE("/:uuid", s.subscriptionHandler.DeleteSubscription)
		api.GET("/list", s.subscriptionHandler.ListSubscriptions)
	}
}
