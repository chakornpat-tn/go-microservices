package main

import (
	"context"
	"log"
	"os"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/pkg/database"
)

func main() {
	ctx := context.Background()

	//Init Env Config
	cfg := config.LoadConfig(
		func() string {
			if len(os.Args) < 2 {
				log.Fatal("Please provide the config file path")
			}
			return os.Args[1]
		}(),
	)

	//Db connection
	db := database.DbConn(ctx, &cfg)
	defer db.Disconnect(ctx)

}
