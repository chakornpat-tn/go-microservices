package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/item"
	"github.com/chakornpat-tn/go-microservices/pkg/database"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func itemDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return database.DbConn(pctx, cfg).Database("item_db")
}

func ItemMigrate(pctx context.Context, cfg *config.Config) {
	db := itemDbConn(pctx, cfg)
	defer db.Client().Disconnect(pctx)

	col := db.Collection("items")
	indexs, _ := col.Indexes().CreateMany(pctx, []mongo.IndexModel{
		{Keys: bson.D{bson.E{Key: "_id", Value: 1}}},
		{Keys: bson.D{bson.E{Key: "title", Value: 1}}},
	})
	fmt.Println(indexs)

	documents := func() []any {
		items := []*item.Item{
			{
				Title:       "Diamond Sword",
				Price:       1000,
				ImageUrl:    "https://i.imgur.com/1Y8tQZM.png",
				UsageStatus: true,
				Damage:      100,
				CreatedAt:   utils.LocalTime(),
				UpdatedAt:   utils.LocalTime(),
			},
			{
				Title:       "Iron Sword",
				Price:       500,
				ImageUrl:    "https://i.imgur.com/1Y8tQZM.png",
				UsageStatus: true,
				Damage:      50,
				CreatedAt:   utils.LocalTime(),
				UpdatedAt:   utils.LocalTime(),
			},
			{
				Title:       "Wooden Sword",
				Price:       100,
				ImageUrl:    "https://i.imgur.com/1Y8tQZM.png",
				UsageStatus: true,
				Damage:      20,
				CreatedAt:   utils.LocalTime(),
				UpdatedAt:   utils.LocalTime(),
			},
		}
		docs := make([]any, 0)
		for _, item := range items {
			docs = append(docs, item)
		}
		return docs
	}()

	result, err := col.InsertMany(pctx, documents)
	if err != nil {
		panic(err)
	}
	log.Println("Migrate items completed: ", result)
}
