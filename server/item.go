package server

import (
	"log"

	"github.com/chakornpat-tn/go-microservices/modules/item/itemHandler"
	itemPb "github.com/chakornpat-tn/go-microservices/modules/item/itemPb"
	"github.com/chakornpat-tn/go-microservices/modules/item/itemRepository"
	"github.com/chakornpat-tn/go-microservices/modules/item/itemUsecase"
	"github.com/chakornpat-tn/go-microservices/pkg/grpccon"
)

func (s *server) itemService() {
	repo := itemRepository.NewItemRepository(s.db)
	usecase := itemUsecase.NewItemUsecase(repo)
	httpHandler := itemHandler.NewItemHttpHandler(s.cfg, usecase)
	grpcHandler := itemHandler.NewItemGrpcHandler(usecase)

	_ = httpHandler

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.ItemUrl)
		itemPb.RegisterItemGrpcServiceServer(grpcServer, grpcHandler)
		log.Printf("Item gRPC server listening on %s", s.cfg.Grpc.ItemUrl)
		grpcServer.Serve(lis)
	}()

	item := s.app.Group("/item_v1")

	// Health Check
	item.GET("", s.healthCheckService)

}
