package inventoryHandler

import (
	"context"

	inventoryPb "github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryPb"
	"github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryUsecase"
)

type (
	inventoryGrpcHandler struct {
		inventoryPb.UnimplementedInventoryGrpcServiceServer
		inventoryUsecase inventoryUsecase.InventoryUsecaseService
	}
)

func NewInventoryGrpcHandler(inventoryUsecase inventoryUsecase.InventoryUsecaseService) *inventoryGrpcHandler {
	return &inventoryGrpcHandler{
		inventoryUsecase: inventoryUsecase,
	}
}

func (g *inventoryGrpcHandler) IsAvailableToSell(ctx context.Context, req *inventoryPb.IsAvailableToSellReq) (*inventoryPb.IsAvailableToSellRes, error) {
	return nil, nil
}
