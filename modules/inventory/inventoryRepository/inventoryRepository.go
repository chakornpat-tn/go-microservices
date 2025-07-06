package inventoryRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/chakornpat-tn/go-microservices/modules/inventory"
	itemPb "github.com/chakornpat-tn/go-microservices/modules/item/itemPb"
	"github.com/chakornpat-tn/go-microservices/modules/models"
	"github.com/chakornpat-tn/go-microservices/pkg/grpccon"
	"github.com/chakornpat-tn/go-microservices/pkg/jwtauth"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type (
	InventoryRepositoryService interface {
		FindItemsInIDs(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
		FindPlayerItems(pctx context.Context, filter bson.D, opts ...options.Lister[options.FindOptions]) ([]*inventory.Inventory, error)
		CountPlayerItems(pctx context.Context, playerID string) (int64, error)
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
	}

	inventoryRepository struct {
		db *mongo.Client
	}
)

func NewInventoryRepository(db *mongo.Client) InventoryRepositoryService {
	return &inventoryRepository{
		db: db,
	}
}

func (r *inventoryRepository) GetOffset(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory_queue")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("\nError: get offset failed: %s\n", err.Error())
		return -1, errors.New("error:get offset failed")
	}

	return result.Offset, nil
}

func (r *inventoryRepository) UpserOffset(pctx context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("players_inventory_queue")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		log.Printf("Error: upsert offset failed: %s", err.Error())
		return errors.New("error:upsert offset failed")
	}

	log.Printf("\n Upsert offset result: %v \n", result)

	return nil
}

func (r *inventoryRepository) inventoryDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("inventory_db")
}

func (r *inventoryRepository) FindItemsInIDs(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpccClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc client connection failed: %s", err.Error())
		return nil, errors.New("error:grpc connection failed")
	}

	jwtauth.SetApiKeyInContext(&ctx)

	result, err := conn.Item().FindItemsInIds(ctx, req)
	if err != nil {
		log.Printf("Error: gRPC FindItemsInIds failed: %s", err.Error())
		return nil, errors.New("error:gRPC find items in ids failed")
	}

	if len(result.Items) == 0 {
		log.Printf("Error: gRPC FindItemsInIds failed")
		return nil, errors.New("error:gRPC find items in ids failed")
	}

	return result, nil
}

func (r *inventoryRepository) FindPlayerItems(pctx context.Context, filter bson.D, opts ...options.Lister[options.FindOptions]) ([]*inventory.Inventory, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("inventories")

	cursors, err := col.Find(ctx, filter, opts...)
	if err != nil {
		log.Printf("Error: FindPlayerItems failed: %s", err.Error())
		return nil, errors.New("error: find player items failed")
	}

	results := make([]*inventory.Inventory, 0)
	for cursors.Next(ctx) {
		result := new(inventory.Inventory)
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: FindPlayerItems failed: %s", err.Error())
			return nil, errors.New("error: find player items failed")
		}

		results = append(results, result)
	}

	return results, nil
}

func (r *inventoryRepository) CountPlayerItems(pctx context.Context, playerID string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.inventoryDbConn(ctx)
	col := db.Collection("inventories")

	count, err := col.CountDocuments(ctx, bson.M{"player_id": playerID})
	if err != nil {
		log.Printf("Error: PlayerItems Failed: %s", err.Error())
		return -1, errors.New("error: count player items failed")
	}

	return count, nil
}
