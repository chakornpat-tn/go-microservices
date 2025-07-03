package inventoryUsecase

import (
	"context"
	"fmt"

	"github.com/chakornpat-tn/go-microservices/modules/inventory"
	"github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryRepository"
	"github.com/chakornpat-tn/go-microservices/modules/item"
	"github.com/chakornpat-tn/go-microservices/modules/models"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type (
	InventoryUsecaseService interface {
		FindPlayerItems(pctx context.Context, basePaginateUrl, playerID string, req *inventory.InventorySearchReq) (*models.PaginateRes, error)
	}

	inventoryUsecase struct {
		inventoryRepo inventoryRepository.InventoryRepositoryService
	}
)

func NewInventoryUsecase(inventoryRepo inventoryRepository.InventoryRepositoryService) InventoryUsecaseService {
	return &inventoryUsecase{
		inventoryRepo: inventoryRepo,
	}
}

func (u *inventoryUsecase) FindPlayerItems(pctx context.Context, basePaginateUrl, playerID string, req *inventory.InventorySearchReq) (*models.PaginateRes, error) {
	filter := bson.D{}
	// Find many items filter
	if req.Start != "" {
		filter = append(filter, bson.E{Key: "_id", Value: bson.D{{Key: "$gt", Value: utils.ConvToObjID(req.Start)}}})
	}
	filter = append(filter, bson.E{Key: "player_id", Value: playerID})

	opts := options.Find().SetSort(bson.D{{
		Key:   "_id",
		Value: 1,
	}}).SetLimit(int64(req.Limit))

	// Find
	InventoryData, err := u.inventoryRepo.FindPlayerItems(pctx, filter, opts)
	if err != nil {
		return nil, err
	}

	results := make([]*inventory.ItemInInventory, 0)
	for _, v := range InventoryData {
		results = append(results, &inventory.ItemInInventory{
			InventoryID: v.ID.Hex(),
			PlayerID:    v.PlayerID,
			ItemShowCase: &item.ItemShowCase{
				ItemID: v.ItemID,
			},
		})
	}

	total, err := u.inventoryRepo.CountPlayerItems(pctx, playerID)
	if err != nil {
		return nil, err
	}

	return &models.PaginateRes{
		Data:  results,
		Total: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s/%s?limit=%d", basePaginateUrl, playerID, req.Limit),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].InventoryID,
			Href:  fmt.Sprintf("%s/%s?limit=%d&start=%s", basePaginateUrl, playerID, req.Limit, results[len(results)-1].InventoryID),
		},
	}, nil
}
