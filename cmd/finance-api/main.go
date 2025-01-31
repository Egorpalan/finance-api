package main

import (
	"context"
	"github.com/Egorpalan/finance-api/internal/handler"
	"github.com/Egorpalan/finance-api/internal/repository"
	"github.com/Egorpalan/finance-api/internal/service"
	"github.com/Egorpalan/finance-api/pkg/config"
	"github.com/Egorpalan/finance-api/pkg/database"
	"github.com/Egorpalan/finance-api/pkg/server"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	cfg, err := config.LoadConfig(".env.example")
	if err != nil {
		logrus.Fatalf("error initializing configs %s", err.Error())
	}

	db, err := database.NewPostgresConnection(cfg)
	if err != nil {

		logrus.Fatalf("Failed to connect to database: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := server.NewServer(cfg)

	srv.Router.POST("/balance/top-up", handlers.TopUpBalance)
	srv.Router.POST("/transfer", handlers.TransferMoney)
	srv.Router.GET("/transactions/:user_id", handlers.GetTransactions)
	srv.Router.POST("/user", handlers.CreateUserHandler)

	go func() {
		if err := srv.Run(); err != nil {
			logrus.Fatalf("Failed to start server: %s", err.Error())
		}
	}()
	logrus.Info("App started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit
	logrus.Info("App shutdown")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Failed to shutdown server: %s", err.Error())
	}

}
