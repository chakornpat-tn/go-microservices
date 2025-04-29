package server

import (
	"github.com/chakornpat-tn/go-microservices/modules/auth/authHandler"
	"github.com/chakornpat-tn/go-microservices/modules/auth/authRepository"
	"github.com/chakornpat-tn/go-microservices/modules/auth/authUsecase"
)

func (s *server) authService() {
	repo := authRepository.NewAuthRepository(s.db)
	usecase := authUsecase.NewAuthUsecase(repo)
	httpHandler := authHandler.NewAuthHttpHandler(s.cfg, usecase)
	grpcHandler := authHandler.NewAuthGrpcHandler(usecase)

	_ = httpHandler
	_ = grpcHandler

	auth := s.app.Group("/auth_v1")

	// Health Check
	auth.GET("", s.healthCheckService)

}
