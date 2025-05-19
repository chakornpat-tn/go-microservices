package middlewareRepository

import (
	"context"
	"errors"
	"log"
	"time"

	authPb "github.com/chakornpat-tn/go-microservices/modules/auth/authPb"
	"github.com/chakornpat-tn/go-microservices/pkg/grpccon"
)

type (
	MiddlewareRepositoryService interface {
		AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error
	}

	middlewareRepository struct{}
)

func NewMiddlewareRepository() MiddlewareRepositoryService {
	return &middlewareRepository{}
}

func (r *middlewareRepository) AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error {
	ctx, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpccClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gPRC connection failed: %s", err.Error())
		return errors.New("gPRC connection failed")
	}

	result, err := conn.Auth().AccessTokenSearch(ctx, &authPb.AccessTokenSearchReq{
		AccessToken: accessToken,
	})
	if err != nil {
		log.Printf("Error: AccessTokenSearch failed: %s", err.Error())
		return errors.New("error: access token search failed")
	}

	if result == nil {
		log.Printf("Error: AccessTokenSearch invalid response")
		return errors.New("error: access token search invalid")
	}

	if !result.IsValid {
		log.Printf("Error: AccessTokenSearch invalid response")
		return errors.New("error: access token search invalid")
	}

	return nil
}
