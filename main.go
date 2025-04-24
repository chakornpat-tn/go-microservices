package main

import (
	"context"
	"log"
	"os"

	"github.com/chakornpat-tn/go-microservices/config"
)

func main() {
	ctx := context.Background()
	_ = ctx

	//Init Env Config
	cfg := config.LoadConfig(
		func() string {
			if len(os.Args) < 2 {
				log.Fatal("Please provide the config file path")
			}
			return os.Args[1]
		}(),
	)

	log.Println(cfg)
}
