package itemRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/chakornpat-tn/go-microservices/modules/item"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	ItemRepositoryService interface {
		IsUniqueItem(pctx context.Context, title string) bool
		InsertOneItem(pctx context.Context, req *item.Item) (bson.ObjectID, error)
		FindOneItem(pctx context.Context, itemID string) (*item.Item, error)
	}

	itemRepository struct {
		db *mongo.Client
	}
)

func NewItemRepository(db *mongo.Client) ItemRepositoryService {
	return &itemRepository{
		db: db,
	}
}

func (r *itemRepository) itemDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("item_db")
}

func (r *itemRepository) IsUniqueItem(pctx context.Context, title string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result := new(item.Item)
	if err := col.FindOne(
		ctx,
		bson.M{"title": title},
	).Decode(result); err != nil {
		log.Printf("Error: IsUniqueItem failed: %s", err.Error())
		return true
	}

	return false
}

func (r *itemRepository) InsertOneItem(pctx context.Context, req *item.Item) (bson.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	itemID, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOneItem failed: %s", err.Error())
		return bson.NilObjectID, errors.New("error: insert one item failed")
	}

	return itemID.InsertedID.(bson.ObjectID), nil

}

func (r *itemRepository) FindOneItem(pctx context.Context, itemID string) (*item.Item, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result := new(item.Item)
	if err := col.FindOne(ctx, bson.M{
		"_id": utils.ConvToObjID(itemID),
	}).Decode(result); err != nil {
		log.Printf("Error: FindOneItem failed: %s", err.Error())
		return nil, errors.New("error: find one item failed")
	}

	return result, nil
}
