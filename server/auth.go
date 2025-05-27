package server

import (
	"log"

	"github.com/chakornpat-tn/go-microservices/modules/auth/authHandler"
	authPb "github.com/chakornpat-tn/go-microservices/modules/auth/authPb"
	"github.com/chakornpat-tn/go-microservices/modules/auth/authRepository"
	"github.com/chakornpat-tn/go-microservices/modules/auth/authUsecase"
	"github.com/chakornpat-tn/go-microservices/pkg/grpccon"
)

func (s *server) authService() {
	repo := authRepository.NewAuthRepository(s.db)
	usecase := authUsecase.NewAuthUsecase(repo)
	httpHandler := authHandler.NewAuthHttpHandler(s.cfg, usecase)
	grpcHandler := authHandler.NewAuthGrpcHandler(usecase)

	// Grpc
	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)
		authPb.RegisterAuthGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("Auth gRPC server listening on %s", s.cfg.Grpc.AuthUrl)
		grpcServer.Serve(lis)
	}()

	auth := s.app.Group("/auth_v1")

	// Health Check
	auth.GET("", s.healthCheckService)
	auth.POST("/login", httpHandler.Login)
	auth.POST("/logout", httpHandler.Logout)
	auth.POST("/refresh-token", httpHandler.RefreshToken)

}
