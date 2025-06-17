package grpccon

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/pkg/jwtauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	authPb "github.com/chakornpat-tn/go-microservices/modules/auth/authPb"
	inventoryPb "github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryPb"
	itemPb "github.com/chakornpat-tn/go-microservices/modules/item/itemPb"
	playerPb "github.com/chakornpat-tn/go-microservices/modules/player/playerPb"
)

type (
	GrpcClientFactoryHandler interface {
		Auth() authPb.AuthGrpcServiceClient
		Player() playerPb.PlayerGrpcServiceClient
		Item() itemPb.ItemGrpcServiceClient
	}

	GrpcClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct {
		secretKey string
	}
)

func (g *grpcAuth) unaryAuthorization(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("Error: metadata not found in context")
		return nil, errors.New("error: metadata not found in context")
	}

	authHeader, ok := md["auth"]
	if !ok || len(authHeader) == 0 {
		log.Printf("Error: auth header not found in metadata")
		return nil, errors.New("error: auth header not found in metadata")
	}

	_, err := jwtauth.ParseToken(g.secretKey, authHeader[0])
	if err != nil {
		log.Printf("Error: failed to parse token: %s", err.Error())
		return nil, errors.New("error: failed to parse token")
	}

	return handler(ctx, req)
}

func (g *GrpcClientFactory) Auth() authPb.AuthGrpcServiceClient {
	return authPb.NewAuthGrpcServiceClient(g.client)
}

func (g *GrpcClientFactory) Inventory() inventoryPb.InventoryGrpcServiceClient {
	return inventoryPb.NewInventoryGrpcServiceClient(g.client)
}

func (g *GrpcClientFactory) Item() itemPb.ItemGrpcServiceClient {
	return itemPb.NewItemGrpcServiceClient(g.client)
}

func (g *GrpcClientFactory) Player() playerPb.PlayerGrpcServiceClient {
	return playerPb.NewPlayerGrpcServiceClient(g.client)
}

func NewGrpccClient(host string) (GrpcClientFactoryHandler, error) {
	opts := make([]grpc.DialOption, 0)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	clientConn, err := grpc.NewClient(host, opts...)
	if err != nil {
		log.Printf("Error: grpc client connection failed: %v", err)
		return nil, errors.New("error: grpc client connection failed")
	}

	return &GrpcClientFactory{
		client: clientConn,
	}, nil
}

func NewGrpcServer(cfg *config.Jwt, host string) (*grpc.Server, net.Listener) {
	opts := make([]grpc.ServerOption, 0)

	grpcAuth := &grpcAuth{
		secretKey: cfg.ApiSecretKey,
	}

	opts = append(opts, grpc.UnaryInterceptor(grpcAuth.unaryAuthorization))

	grpcServer := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Error: Failed to listen: %v", err)
	}

	return grpcServer, lis
}
