package authRepository

import (
	"context"
	"errors"
	"log"
	"time"

	playerPb "github.com/chakornpat-tn/go-microservices/modules/player/playerPb"
	"github.com/chakornpat-tn/go-microservices/pkg/grpccon"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	AuthRepositoryService interface{}

	authRepository struct {
		db *mongo.Client
	}
)

func NewAuthRepository(db *mongo.Client) AuthRepositoryService {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) authDbConn(pctx context.Context) *mongo.Database {
	return r.db.Database("auth_db")
}

func (r *authRepository) CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error) {
	ctx, cancel := context.WithTimeout(pctx, 30*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpccClient(grpcUrl)
	if err != nil {
		log.Printf("Error: grpc client connection failed: %s", err.Error())
		return nil, errors.New("error:grpc connection failed")
	}

	result, err := conn.Player().CredentialSearch(ctx, req)
	if err != nil {
		log.Printf("Error: gRPC CredentialSearch failed: %s", err.Error())
		return nil, errors.New(err.Error())
	}

	return result, nil
}
