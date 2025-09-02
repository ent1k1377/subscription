package http

import (
	"context"
	"net/http"

	"github.com/ent1k1377/subscriptions/internal/config"
	"github.com/ent1k1377/subscriptions/internal/transport/http/handler/subscription"

	"github.com/ent1k1377/subscriptions/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	httpServer          *http.Server
	engine              *gin.Engine
	subscriptionHandler *subscription.Handler
}

func NewServer(cfg config.ServerConfig, subscriptionHandler *subscription.Handler) *Server {
	engine := gin.Default()
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
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
	docs.SwaggerInfo.BasePath = "/api/"
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := s.engine.Group("/api/subscriptions")
	{
		api.POST("/create", s.subscriptionHandler.Create)
		api.GET("/:uuid", s.subscriptionHandler.GetSubscription)
		api.PUT("/:uuid", s.subscriptionHandler.UpdateSubscription)
		api.DELETE("/:uuid", s.subscriptionHandler.DeleteSubscription)
		api.GET("/list", s.subscriptionHandler.ListSubscriptions)
		api.GET("/sum", s.subscriptionHandler.TotalCostSubscriptions)
	}
}
