package http

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"subscriptions/internal/transport/http/handler/subscription"
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
		api.GET("/:uuid", nil)
		api.PUT("/:uuid", nil)
		api.DELETE("/:uuid", nil)
		api.GET("/list", nil)
	}
}
