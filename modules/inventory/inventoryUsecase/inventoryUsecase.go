package inventoryUsecase

import (
	"context"
	"fmt"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/inventory"
	"github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryRepository"
	"github.com/chakornpat-tn/go-microservices/modules/item"
	itemPb "github.com/chakornpat-tn/go-microservices/modules/item/itemPb"
	"github.com/chakornpat-tn/go-microservices/modules/models"
	"github.com/chakornpat-tn/go-microservices/modules/payment"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type (
	InventoryUsecaseService interface {
		FindPlayerItems(pctx context.Context, cfg *config.Config, playerID string, req *inventory.InventorySearchReq) (*models.PaginateRes, error)
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
		AddPlayerItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq)
		RemovePlayerItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq)
		RollbackAddPlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventory)
		RollbackRemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventory)
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

func (u *inventoryUsecase) GetOffset(pctx context.Context) (int64, error) {
	return u.inventoryRepo.GetOffset(pctx)
}

func (u *inventoryUsecase) UpserOffset(pctx context.Context, offset int64) error {
	return u.inventoryRepo.UpserOffset(pctx, offset)
}

func (u *inventoryUsecase) FindPlayerItems(pctx context.Context, cfg *config.Config, playerID string, req *inventory.InventorySearchReq) (*models.PaginateRes, error) {
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
	inventoryData, err := u.inventoryRepo.FindPlayerItems(pctx, filter, opts)
	if err != nil {
		return nil, err
	}
	if len(inventoryData) == 0 {
		return &models.PaginateRes{
			Data:  make([]*inventory.ItemInInventory, 0),
			Total: 0,
			Limit: req.Limit,
			First: models.FirstPaginate{
				Href: fmt.Sprintf("%s/%s?limit=%d", cfg.Paginate.InventoryNextPageBasedUrl, playerID, req.Limit),
			},
			Next: models.NextPaginate{
				Start: "",
				Href:  "",
			},
		}, nil
	}

	itemData, err := u.inventoryRepo.FindItemsInIDs(pctx, cfg.Grpc.ItemUrl, &itemPb.FindItemsInIdsReq{
		Ids: func() []string {
			itemIds := make([]string, 0)
			for _, v := range inventoryData {
				itemIds = append(itemIds, v.ItemID)
			}
			return itemIds
		}(),
	})

	itemMaps := make(map[string]*item.ItemShowCase)
	for _, v := range itemData.Items {
		itemMaps[v.Id] = &item.ItemShowCase{
			ItemID:   v.Id,
			Title:    v.Title,
			Price:    v.Price,
			ImageUrl: v.ImgUrl,
			Damage:   int(v.Damage),
		}
	}

	results := make([]*inventory.ItemInInventory, 0)
	for _, v := range inventoryData {
		results = append(results, &inventory.ItemInInventory{
			InventoryID: v.ID.Hex(),
			PlayerID:    v.PlayerID,
			ItemShowCase: func() *item.ItemShowCase {
				return &item.ItemShowCase{
					ItemID:   v.ItemID,
					Title:    itemMaps[v.ItemID].Title,
					Price:    itemMaps[v.ItemID].Price,
					Damage:   itemMaps[v.ItemID].Damage,
					ImageUrl: itemMaps[v.ItemID].ImageUrl,
				}
			}(),
		})
	}

	// Count
	total, err := u.inventoryRepo.CountPlayerItems(pctx, playerID)
	if err != nil {
		return nil, err
	}

	return &models.PaginateRes{
		Data:  results,
		Total: total,
		Limit: req.Limit,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s/%s?limit=%d", cfg.Paginate.InventoryNextPageBasedUrl, playerID, req.Limit),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].InventoryID,
			Href:  fmt.Sprintf("%s/%s?limit=%d&start=%s", cfg.Paginate.InventoryNextPageBasedUrl, playerID, req.Limit, results[len(results)-1].InventoryID),
		},
	}, nil
}

func (u *inventoryUsecase) AddPlayerItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) {
	inventoryID, err := u.inventoryRepo.InsertOnePlayerItem(pctx, &inventory.Inventory{
		PlayerID: req.PlayerID,
		ItemID:   req.ItemID,
	})
	if err != nil {
		u.inventoryRepo.AddPlayerItemRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryID:   inventoryID.Hex(),
			TransactionID: "",
			PlayerID:      req.PlayerID,
			ItemID:        req.ItemID,
			Amount:        0,
			Error:         err.Error(),
		})
		return
	}

	u.inventoryRepo.AddPlayerItemRes(pctx, cfg, &payment.PaymentTransferRes{
		InventoryID:   inventoryID.Hex(),
		TransactionID: "",
		PlayerID:      req.PlayerID,
		ItemID:        req.ItemID,
		Amount:        0,
		Error:         "",
	})

}

func (u *inventoryUsecase) RemovePlayerItemRes(pctx context.Context, cfg *config.Config, req *inventory.UpdateInventoryReq) {
	if !u.inventoryRepo.FindOnePlayerItem(pctx, req.PlayerID, req.ItemID) {
		u.inventoryRepo.RemovePlayerItemRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryID:   "",
			TransactionID: "",
			PlayerID:      req.PlayerID,
			ItemID:        req.ItemID,
			Amount:        0,
			Error:         "error: item not found",
		})
		return

	}
	if err := u.inventoryRepo.DeleteOnePlayerItem(pctx, req.PlayerID, req.ItemID); err != nil {
		u.inventoryRepo.RemovePlayerItemRes(pctx, cfg, &payment.PaymentTransferRes{
			InventoryID:   "",
			TransactionID: "",
			PlayerID:      req.PlayerID,
			ItemID:        req.ItemID,
			Amount:        0,
			Error:         err.Error(),
		})
		return
	}

	u.inventoryRepo.RemovePlayerItemRes(pctx, cfg, &payment.PaymentTransferRes{
		InventoryID:   "",
		TransactionID: "",
		PlayerID:      req.PlayerID,
		ItemID:        req.ItemID,
		Amount:        0,
		Error:         "",
	})
}

func (u *inventoryUsecase) RollbackAddPlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventory) {
	u.inventoryRepo.DeleteOneInventory(pctx, req.InventoryID)
}

func (u *inventoryUsecase) RollbackRemovePlayerItem(pctx context.Context, cfg *config.Config, req *inventory.RollbackPlayerInventory) {
	u.inventoryRepo.InsertOnePlayerItem(pctx, &inventory.Inventory{
		PlayerID: req.PlayerID,
		ItemID:   req.ItemID,
	})
}
