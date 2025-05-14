package authRepository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/chakornpat-tn/go-microservices/modules/auth"
	playerPb "github.com/chakornpat-tn/go-microservices/modules/player/playerPb"
	"github.com/chakornpat-tn/go-microservices/pkg/grpccon"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type (
	AuthRepositoryService interface {
		CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error)
		InsertOnePlayerCredential(pctx context.Context, req *auth.Credential) (bson.ObjectID, error)
		FindOnePlayerCredential(pctx context.Context, credentialId string) (*auth.Credential, error)
	}

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

func (r *authRepository) InsertOnePlayerCredential(pctx context.Context, req *auth.Credential) (bson.ObjectID, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	result, err := col.InsertOne(ctx, req)
	if err != nil {
		log.Printf("Error: InsertOnePlayerCredential failed: %s", err.Error())
		return bson.NilObjectID, errors.New("error: InsertOnePlayerCredential failed")
	}

	return result.InsertedID.(bson.ObjectID), nil
}

func (r *authRepository) FindOnePlayerCredential(pctx context.Context, credentialId string) (*auth.Credential, error) {
	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
	defer cancel()

	db := r.authDbConn(ctx)
	col := db.Collection("auth")

	var result auth.Credential
	err := col.FindOne(ctx, bson.M{
		"_id": utils.ConvToObjID(credentialId),
	}).Decode(&result)
	if err != nil {
		log.Printf("Error: FindOnePlayerCredential failed: %s", err.Error())
		return nil, errors.New("error: FindOnePlayerCredential failed")
	}

	return &result, nil
}
