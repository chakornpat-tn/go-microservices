package migration

import (
	"context"
	"log"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/player"
	"github.com/chakornpat-tn/go-microservices/pkg/database"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func playerDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("player_db")
}

func PlayerMigrate(pctx context.Context, cfg *config.Config) {
	db := playerDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("player_transactions")
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{bson.E{Key: "_id", Value: 1}}},
		{Keys: bson.D{bson.E{Key: "player_id", Value: 1}}},
	})
	log.Println(indexs)

	col = db.Collection("players")
	indexs, _ = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{bson.E{Key: "_id", Value: 1}}},
		{Keys: bson.D{bson.E{Key: "email", Value: 1}}},
	})
	log.Println(indexs)

	documents := func() []any {
		roles := []*player.Player{
			{
				Email:    "player001@email.com",
				Password: "123456",
				Username: "Player001",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email:    "player002@email.com",
				Password: "123456",
				Username: "Player002",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "player",
						RoleCode:  0,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
			{
				Email:    "admin001@email.com",
				Password: "123456",
				Username: "Admin001",
				PlayerRoles: []player.PlayerRole{
					{
						RoleTitle: "admin",
						RoleCode:  1,
					},
				},
				CreatedAt: utils.LocalTime(),
				UpdatedAt: utils.LocalTime(),
			},
		}
		docs := make([]any, 0)
		for _, r := range roles {
			docs = append(docs, r)
		}
		return docs
	}()

	results, err := col.InsertMany(pctx, documents)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate player completed: ", results)

	playerTransactions := make([]any, 0)
	for _, p := range results.InsertedIDs {
		playerTransactions = append(playerTransactions, &player.PlayerTransactions{
			PlayerID:  "player:" + p.(bson.ObjectID).Hex(),
			Amount:    1000,
			CreatedAt: utils.LocalTime(),
		})
	}
	col = db.Collection("player_transactions")
	results, err = col.InsertMany(pctx, playerTransactions)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate player_transactions completed: ", results)
	col = db.Collection("player_transactions_queue")
	result, err := col.InsertOne(pctx, bson.M{"offset": -1}, nil)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate player_transactions_queue completed: ", result)

}
