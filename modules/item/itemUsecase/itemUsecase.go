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
		CreateItem(pctx context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error)
		FindOneItem(pcxt context.Context, itemID string) (*item.ItemShowCase, error)
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

func (u *itemUsecase) CreateItem(pctx context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error) {
	if !u.itemRepo.IsUniqueItem(pctx, req.Title) {
		return nil, errors.New("error: item title already exists")
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
		return nil, errors.New("error: insert one item failed")
	}
	return u.FindOneItem(pctx, itemID.Hex())
}

func (u *itemUsecase) FindOneItem(pcxt context.Context, itemID string) (*item.ItemShowCase, error) {
	result, err := u.itemRepo.FindOneItem(pcxt, itemID)
	if err != nil {
		return nil, errors.New("error: find one item failed")
	}

	return &item.ItemShowCase{
		ItemID:   "item:" + result.ID.Hex(),
		Title:    result.Title,
		Price:    result.Price,
		Damage:   result.Damage,
		ImageUrl: result.ImageUrl,
	}, nil
}
