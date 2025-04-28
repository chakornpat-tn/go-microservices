package grpccon

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/chakornpat-tn/go-microservices/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authPb "github.com/chakornpat-tn/go-microservices/modules/auth/authPb"
	inventoryPb "github.com/chakornpat-tn/go-microservices/modules/inventory/inventoryPb"
	itemPb "github.com/chakornpat-tn/go-microservices/modules/item/itemPb"
	playerPb "github.com/chakornpat-tn/go-microservices/modules/player/playerPb"
)

type (
	GrpcClientFactoryHandler interface {
	}

	GrpcClientFactory struct {
		client *grpc.ClientConn
	}

	grpcAuth struct{}
)

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

func NewGrpccClient(ctx context.Context, host string) (GrpcClientFactoryHandler, error) {
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
	grpcServer := grpc.NewServer(opts...)

	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("Error: Failed to listen: %v", err)
	}

	return grpcServer, lis
}
