package itemUsecase

import (
	"context"
	"errors"

	"github.com/chakornpat-tn/go-microservices/modules/item"
	"github.com/chakornpat-tn/go-microservices/modules/item/itemRepository"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
)

type (
	ItemUsecaseService interface {
		CreateItem(pctx context.Context, req *item.CreateItemReq) (any, error)
	}

	itemUsecase struct {
		itemRepo itemRepository.ItemRepositoryService
	}
)

func NewItemUsecase(itemRepo itemRepository.ItemRepositoryService) ItemUsecaseService {
	return &itemUsecase{
		itemRepo: itemRepo,
	}
}

func (u *itemUsecase) CreateItem(pctx context.Context, req *item.CreateItemReq) (any, error) {
	if !u.itemRepo.IsUniqueItem(pctx, req.Title) {
		return "", errors.New("error: item title already exists")
	}

	itemID, err := u.itemRepo.InsertOneItem(pctx, &item.Item{
		Title:       req.Title,
		Price:       req.Price,
		Damage:      req.Damage,
		UsageStatus: true,
		ImageUrl:    req.ImageUrl,
		CreatedAt:   utils.LocalTime(),
		UpdatedAt:   utils.LocalTime(),
	})
	if err != nil {
		return "", errors.New("error: insert one item failed")
	}
	return itemID.Hex(), nil
}
