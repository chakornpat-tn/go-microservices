package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/middleware/middlewareHandler"
	"github.com/chakornpat-tn/go-microservices/modules/middleware/middlewareRepository"
	"github.com/chakornpat-tn/go-microservices/modules/middleware/middlewareUsecase"
	"github.com/chakornpat-tn/go-microservices/pkg/jwtauth"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	server struct {
		app        *echo.Echo
		db         *mongo.Client
		cfg        *config.Config
		middleware middlewareHandler.MiddlewareHandlerService
	}
)

func newMiddleware(cfg *config.Config) middlewareHandler.MiddlewareHandlerService {
	repo := middlewareRepository.NewMiddlewareRepository()
	usecase := middlewareUsecase.NewMiddlewareUsecase(repo)
	return middlewareHandler.NewMiddlewareHandler(cfg, usecase)

}

func (s *server) graceFulShutdown(pctx context.Context, quit <-chan os.Signal) {
	log.Println("Start Service:", s.cfg.App.Name)
	<-quit
	log.Println("Shutting down Service:", s.cfg.App.Name)

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	if err := s.app.Shutdown(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}

}

func (s *server) httpListening() {
	if err := s.app.Start(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error: %v", err)
	}
}

func Start(pctx context.Context, cfg *config.Config, db *mongo.Client) {
	s := &server{
		app:        echo.New(),
		db:         db,
		cfg:        cfg,
		middleware: newMiddleware(cfg),
	}

	jwtauth.SetApiKey(cfg.Jwt.ApiSecretKey)

	// Basic Middleware
	// Request Timeout
	s.app.Use(middleware.TimeoutWithConfig(
		middleware.TimeoutConfig{
			Skipper:      middleware.DefaultSkipper,
			Timeout:      30 * time.Second,
			ErrorMessage: "Error Request Timeout.",
		},
	))

	// CORS
	s.app.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			Skipper:      middleware.DefaultSkipper,
			AllowOrigins: []string{"*"},
			AllowMethods: []string{echo.GET, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		},
	))

	// Body Limit
	s.app.Use(middleware.BodyLimit("10M"))

	// Custom Middleware
	switch s.cfg.App.Name {
	case "auth":
		s.authService()
	case "player":
		s.playerService()
	case "item":
		s.itemService()
	case "inventory":
		s.inventoryService()
	case "payment":
		s.paymentService()
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	s.app.Use(middleware.Logger())

	go s.graceFulShutdown(pctx, quit)

	// Listening
	s.httpListening()

}
