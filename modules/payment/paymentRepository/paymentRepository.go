package paymentRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/chakornpat-tn/go-microservices/config"
	itemPb "github.com/chakornpat-tn/go-microservices/modules/item/itemPb"
	"github.com/chakornpat-tn/go-microservices/pkg/grpccon"
	"github.com/chakornpat-tn/go-microservices/pkg/jwtauth"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	PaymentRepositoryService interface {
		FindItemsInIDs(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error)
	}

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

func (r *paymentRepository) FindItemsInIDs(pctx context.Context, grpcUrl string, req *itemPb.FindItemsInIdsReq) (*itemPb.FindItemsInIdsRes, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpccClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc client connection failed: %s", err.Error())
		return nil, errors.New("error:grpc connection failed")
	}

	jwtauth.SetApiKeyInContext(&ctx)

	result, err := conn.Item().FindItemsInIds(ctx, req)
	if err != nil {
		log.Printf("Error: gRPC FindItemsInIds failed: %s", err.Error())
		return nil, errors.New("error:gRPC find items in ids failed")
	}

	if len(result.Items) == 0 {
		log.Printf("Error: gRPC FindItemsInIds failed")
		return nil, errors.New("error:gRPC find items in ids failed")
	}

	return result, nil
}
