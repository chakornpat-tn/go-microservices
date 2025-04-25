package database

import (
	"context"
	"log"
	"time"

	"github.com/chakornpat-tn/go-microservices/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func DbConn(pctx context.Context, cfg *config.Config) *mongo.Client {

	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.Db.Url))
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	return client
}
