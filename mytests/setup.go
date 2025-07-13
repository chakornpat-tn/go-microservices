package mytests

import (
	"github.com/chakornpat-tn/go-microservices/config"
)

func NewTestConfig() *config.Config {
	cfg := config.LoadConfig("../env/test/.env")
	return &cfg
}
