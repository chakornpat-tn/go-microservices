package playerRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/chakornpat-tn/go-microservices/modules/player"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	PlayerRepositoryService interface {
		IsUniquePlayer(pctx context.Context, email, username string) bool
		InsertOnePlayer(pctx context.Context, req *player.Player) (bson.ObjectID, error)
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
