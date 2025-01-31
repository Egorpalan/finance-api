package server

import (
	"context"
	"github.com/Egorpalan/finance-api/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	Router     *gin.Engine
	Config     *config.Config
}

func NewServer(cfg *config.Config) *Server {
	router := gin.Default()
	return &Server{
		Router: router,
		Config: cfg,
		httpServer: &http.Server{
			Addr:    ":" + cfg.ServerPort,
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	logrus.Printf("Server is running on port %s", s.Config.ServerPort)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	logrus.Info("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
