package server

import (
	"log"

	"github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryHandler"
	inventoryPb "github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryPb"
	"github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryRepository"
	"github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryUsecase"
	"github.com/chakornpat-tn/go-microservices/pkg/grpccon"
)

func (s *server) inventoryService() {
	repo := inventoryRepository.NewInventoryRepository(s.db)
	usecase := inventoryUsecase.NewInventoryUsecase(repo)
	httpHandler := inventoryHandler.NewInventoryHttpHandler(s.cfg, usecase)
	grpcHandler := inventoryHandler.NewInventoryGrpcHandler(usecase)
	queueHandler := inventoryHandler.NewInventoryQueueHandler(s.cfg, usecase)

	_ = httpHandler
	_ = queueHandler

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.InventoryUrl)
		inventoryPb.RegisterInventoryGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("Inventory gRPC server listening on %s", s.cfg.Grpc.InventoryUrl)
		grpcServer.Serve(lis)
	}()

	inventory := s.app.Group("/inventory_v1")

	// Health Check
	inventory.GET("", s.healthCheckService)

}
