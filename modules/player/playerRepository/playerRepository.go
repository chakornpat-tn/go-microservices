package playerRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/chakornpat-tn/go-microservices/modules/models"
	"github.com/chakornpat-tn/go-microservices/modules/player"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type (
	PlayerRepositoryService interface {
		IsUniquePlayer(pctx context.Context, email, username string) bool
		InsertOnePlayer(pctx context.Context, req *player.Player) (bson.ObjectID, error)
		FindOnePlayer(pctx context.Context, playerID string) (*player.PlayerProfileBson, error)
		InsertOnePlayerTranscation(pctx context.Context, req *player.PlayerTransactions) error
		GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error)
		FindOnePlayerCredential(pctx context.Context, email string) (*player.Player, error)
		FindOnePlayerProfileToRefresh(pctx context.Context, playerID string) (*player.Player, error)
		GetOffset(pctx context.Context) (int64, error)
		UpserOffset(pctx context.Context, offset int64) error
	}

	playerRepository struct {
		db *mongo.Client
	}
)

func NewPlayerRepository(db *mongo.Client) PlayerRepositoryService {
	return &playerRepository{
		db: db,
	}
}

func (r *playerRepository) playerDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("player_db")
}

func (r *playerRepository) GetOffset(pctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions")

	result := new(models.KafkaOffset)
	if err := col.FindOne(ctx, bson.M{}).Decode(result); err != nil {
		log.Printf("\nError: get offset failed: %s\n", err.Error())
		return -1, errors.New("error:get offset failed")
	}

	return result.Offset, nil
}

func (r *playerRepository) UpserOffset(pctx context.Context, offset int64) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions")

	result, err := col.UpdateOne(ctx, bson.M{}, bson.M{"$set": bson.M{"offset": offset}}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		log.Printf("Error: upsert offset failed: %s", err.Error())
		return errors.New("error:upsert offset failed")
	}

	log.Printf("\n Upsert offset result: %v \n", result)

	return nil
}

func (r *playerRepository) IsUniquePlayer(pctx context.Context, email, username string) bool {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player")

	player := new(player.Player)
	if err := col.FindOne(ctx, bson.M{"$or": []bson.M{
		{"username": username},
		{"email": email},
	}},
	).Decode(player); err != nil {
		return true
	}

	return false
}
func (r *playerRepository) InsertOnePlayer(pctx context.Context, req *player.Player) (bson.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	playerId, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOnePlayer: %s", err.Error())
		return bson.NilObjectID, errors.New("error: insert one player failed")
	}

	return playerId.InsertedID.(bson.ObjectID), nil
}

func (r *playerRepository) FindOnePlayer(pctx context.Context, playerID string) (*player.PlayerProfileBson, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("players")

	result := new(player.PlayerProfileBson)
	if err := col.FindOne(
		ctx,
		bson.M{"_id": utils.ConvToObjID(playerID)},
		options.FindOne().SetProjection(bson.M{
			"_id":        1,
			"username":   1,
			"email":      1,
			"created_at": 1,
			"updated_at": 1,
		}),
	).Decode(result); err != nil {
		log.Printf("Error: findOnePlayer: %s", err.Error())
		return nil, errors.New("error: player not found")
	}

	return result, nil
}

func (r *playerRepository) InsertOnePlayerTranscation(pctx context.Context, req *player.PlayerTransactions) error {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOnePlayerTranscation: %s", err.Error())
		return errors.New("error: insert one player transaction failed")
	}

	log.Printf("Result: InsertOnePlayerTranscation: %s", result.InsertedID)

	return nil
}

func (r *playerRepository) GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(ctx)
	col := db.Collection("player_transactions")

	filter := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "player_id", Value: playerId}}}},
		bson.D{{
			Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: "$player_id"},
				{Key: "balance", Value: bson.D{{Key: "$sum", Value: "$amount"}}},
			},
		}},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "player_id", Value: "$_id"},
				{Key: "_id", Value: 0},
				{Key: "balance", Value: 1},
			}},
		},
	}
	cursors, err := col.Aggregate(ctx, filter)
	if err != nil {
		log.Printf("Error: GetPlayerSavingAccount: %s", err.Error())
		return nil, errors.New("error: failed to get player saving account")
	}

	result := new(player.PlayerSavingAccount)
	for cursors.Next(ctx) {
		if err := cursors.Decode(result); err != nil {
			log.Printf("Error: GetPlayerSavingAccount: %s", err.Error())
			return nil, errors.New("error: failed to decode player saving account")
		}
	}

	return result, nil
}

func (r *playerRepository) FindOnePlayerCredential(pctx context.Context, email string) (*player.Player, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(pctx)
	col := db.Collection("players")

	result := new(player.Player)
	if err := col.FindOne(ctx, bson.M{"email": email}).Decode(result); err != nil {
		log.Printf("Error: FindOnePlayerCredential: %s", err.Error())
		return nil, errors.New("error: email not found")
	}

	return result, nil
}

func (r *playerRepository) FindOnePlayerProfileToRefresh(pctx context.Context, playerID string) (*player.Player, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.playerDbConn(pctx)
	col := db.Collection("players")

	result := new(player.Player)
	if err := col.FindOne(ctx, bson.M{"_id": utils.ConvToObjID(playerID)}).Decode(result); err != nil {
		log.Printf("Error: FindOnePlayerProfileToRefresh: %s", err.Error())
		return nil, errors.New("error: player profile not found")
	}

	return result, nil
}
