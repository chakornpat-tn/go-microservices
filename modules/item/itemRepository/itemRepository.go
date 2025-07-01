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
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type (
	ItemRepositoryService interface {
		IsUniqueItem(pctx context.Context, title string) bool
		InsertOneItem(pctx context.Context, req *item.Item) (bson.ObjectID, error)
		FindOneItem(pctx context.Context, itemID string) (*item.Item, error)
		FindManyItem(pctx context.Context, filter bson.D, opts ...options.Lister[options.FindOptions]) ([]*item.ItemShowCase, error)
		CountItems(pctx context.Context, filter bson.D) (int64, error)
		UpdateOneItem(pctx context.Context, ItemID string, req bson.M) error
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

func (r *itemRepository) FindManyItem(pctx context.Context, filter bson.D, opts ...options.Lister[options.FindOptions]) ([]*item.ItemShowCase, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	cursor, err := col.Find(ctx, filter, opts...)
	if err != nil {
		log.Printf("Error: FindManyItem failed: %s", err.Error())
		return make([]*item.ItemShowCase, 0), errors.New("error: find many item failed")
	}

	results := make([]*item.ItemShowCase, 0)
	for cursor.Next(ctx) {
		result := new(item.Item)
		if err := cursor.Decode(result); err != nil {
			log.Printf("Error: FindManyItem failed: %s", err.Error())
			return make([]*item.ItemShowCase, 0), errors.New("error: find many item failed")
		}
		results = append(results, &item.ItemShowCase{
			ItemID:   "item:" + result.ID.Hex(),
			Title:    result.Title,
			Price:    result.Price,
			Damage:   result.Damage,
			ImageUrl: result.ImageUrl,
		})
	}

	return results, nil
}

func (r *itemRepository) CountItems(pctx context.Context, filter bson.D) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	count, err := col.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("Error: CountItems failed: %s", err.Error())
		return -1, errors.New("error: count items failed")
	}

	return count, nil
}

func (r *itemRepository) UpdateOneItem(pctx context.Context, ItemID string, req bson.M) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.itemDbConn(ctx)
	col := db.Collection("items")

	result, err := col.UpdateOne(ctx, bson.M{"_id": utils.ConvToObjID(ItemID)}, bson.M{"$set": req})
	if err != nil {
		log.Printf("Error: UpdateOneItem failed: %s", err.Error())
		return errors.New("error: update one item failed")
	}
	log.Printf("UpdateOneItem result: %v", result)

	return nil
}
