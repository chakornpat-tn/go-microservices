package paymentRepository

import (
	"context"

	"github.com/chakornpat-tn/go-microservices/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	PaymentRepositoryService interface{}

	paymentRepository struct {
		db *mongo.Client
	}
)

func NewPaymentRepository(db *mongo.Client) PaymentRepositoryService {
	return &paymentRepository{
		db: db,
	}
}

func (r *paymentRepository) paymentDbConn(pctx context.Context, cfg *config.Config) *mongo.Database {
	return r.db.Database("payment_db")
}
