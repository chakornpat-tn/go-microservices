package itemUsecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/chakornpat-tn/go-microservices/modules/item"
	"github.com/chakornpat-tn/go-microservices/modules/item/itemRepository"
	"github.com/chakornpat-tn/go-microservices/modules/models"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type (
	ItemUsecaseService interface {
		CreateItem(pctx context.Context, req *item.CreateItemReq) (*item.ItemShowCase, error)
		FindOneItem(pcxt context.Context, itemID string) (*item.ItemShowCase, error)
		FindManyItems(pctx context.Context, basePaginateUrl string, req *item.ItemSearchReq) (*models.PaginateRes, error)
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
func (u *itemUsecase) FindManyItems(pctx context.Context, basePaginateUrl string, req *item.ItemSearchReq) (*models.PaginateRes, error) {
	findItemsFilter := bson.D{}
	findItemsOpts := options.Find().SetSort(bson.D{{
		Key:   "_id",
		Value: 1,
	}}).SetLimit(int64(req.Limit))

	countItemsFilter := bson.D{}

	// Find many items filter
	if req.Start != "" {
		req.Start = strings.TrimPrefix(req.Start, "item:")
		findItemsFilter = append(findItemsFilter, bson.E{Key: "_id", Value: bson.D{{Key: "$gt", Value: utils.ConvToObjID(req.Start)}}})
	}

	if req.Title != "" {
		findItemsFilter = append(findItemsFilter, bson.E{Key: "title", Value: bson.Regex{Pattern: fmt.Sprintf(".*%s.*", req.Title), Options: "i"}})
		countItemsFilter = append(countItemsFilter, bson.E{Key: "title", Value: bson.Regex{Pattern: fmt.Sprintf(".*%s.*", req.Title), Options: "i"}})
	}

	findItemsFilter = append(findItemsFilter, bson.E{Key: "usage_status", Value: true})
	countItemsFilter = append(countItemsFilter, bson.E{Key: "usage_status", Value: true})

	// Find
	results, err := u.itemRepo.FindManyItem(pctx, findItemsFilter, findItemsOpts)
	if err != nil {
		return nil, err
	}

	total, err := u.itemRepo.CountItems(pctx, countItemsFilter)
	if err != nil {
		return nil, err
	}
	return &models.PaginateRes{
		Data:  results,
		Limit: req.Limit,
		Total: total,
		First: models.FirstPaginate{
			Href: fmt.Sprintf("%s?limit=%d&title=%s", basePaginateUrl, req.Limit, req.Title),
		},
		Next: models.NextPaginate{
			Start: results[len(results)-1].ItemID,
			Href:  fmt.Sprintf("%s?limit=%d&title=%s&start=%s", basePaginateUrl, req.Limit, req.Title, results[len(results)-1].ItemID),
		},
	}, nil
}
