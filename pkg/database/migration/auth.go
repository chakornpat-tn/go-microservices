package migration

import (
	"context"
	"log"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/auth"
	"github.com/chakornpat-tn/go-microservices/pkg/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func authDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("auth_db")
}

func AuthMigrate(pctx context.Context, cfg *config.Config) {
	db := authDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("auth")
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{bson.E{Key: "_id", Value: 1}}},
		{Keys: bson.D{bson.E{Key: "player_id", Value: 1}}},
		{Keys: bson.D{bson.E{Key: "refresh_token", Value: 1}}},
	})

	for _, index := range indexs {
		log.Printf("Index: %s", index)
	}

	col = db.Collection("roles")
	indexs, _ = col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{bson.E{Key: "_id", Value: 1}}},
		{Keys: bson.D{bson.E{Key: "code", Value: 1}}},
	})
	for _, index := range indexs {
		log.Printf("Index: %s", index)
	}

	// roles data
	documents := func() []any {
		roles := []*auth.Role{
			{
				Title: "user",
				Code:  0,
			},
			{
				Title: "admin",
				Code:  1,
			},
		}
		docs := make([]any, 0)
		for _, r := range roles {
			docs = append(docs, r)
		}
		return docs
	}()

	result, err := col.InsertMany(pctx, documents)
	if err != nil {
		panic(err)
	}
	log.Printf("Migrate auth completed: ", result)
}
