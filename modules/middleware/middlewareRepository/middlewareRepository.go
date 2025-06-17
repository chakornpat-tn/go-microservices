package middlewareRepository

import (
	"context"
	"errors"
	"log"
	"time"

	authPb "github.com/chakornpat-tn/go-microservices/modules/auth/authPb"
	"github.com/chakornpat-tn/go-microservices/pkg/grpccon"
	"github.com/chakornpat-tn/go-microservices/pkg/jwtauth"
)

type (
	MiddlewareRepositoryService interface {
		AccessTokenSearch(pctx context.Context, grpcUrl, accessToken string) error
		RolesCount(pctx context.Context, grpcUrl string) (int64, error)
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

	jwtauth.SetApiKeyInContext(&ctx)

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

func (r *middlewareRepository) RolesCount(pctx context.Context, grpcUrl string) (int64, error) {
	ctx, cancel := context.WithTimeout(pctx, 5*time.Second)
	defer cancel()

	conn, err := grpccon.NewGrpccClient(grpcUrl)
	if err != nil {
		log.Printf("Error: gPRC connection failed: %s", err.Error())
		return -1, errors.New("gPRC connection failed")
	}

	jwtauth.SetApiKeyInContext(&ctx)

	result, err := conn.Auth().RolesCount(ctx, &authPb.RolesCountReq{})

	if err != nil {
		log.Printf("Error: RolesCount failed: %s", err.Error())
		return -1, errors.New("error: roles count failed")
	}

	if result == nil {
		log.Printf("Error: RolesCount failed")
		return -1, errors.New("error: roles count failed")
	}

	return result.Count, nil
}
