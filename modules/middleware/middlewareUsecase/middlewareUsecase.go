package middlewareUsecase

import (
	"errors"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/middleware/middlewareRepository"
	"github.com/chakornpat-tn/go-microservices/pkg/jwtauth"
	"github.com/chakornpat-tn/go-microservices/pkg/rbac"
	"github.com/labstack/echo/v4"
)

type (
	MiddlewareUsecaseService interface {
		JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error)
		RbacAuthorization(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error)
	}

	middlewareUsecase struct {
		middlewareRepo middlewareRepository.MiddlewareRepositoryService
	}
)

func NewMiddlewareUsecase(middlewareRepo middlewareRepository.MiddlewareRepositoryService) MiddlewareUsecaseService {
	return &middlewareUsecase{
		middlewareRepo: middlewareRepo,
	}

}

func (u *middlewareUsecase) JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error) {
	ctx := c.Request().Context()

	claims, err := jwtauth.ParseToken(cfg.Jwt.AccessTokenSecretKey, accessToken)
	if err != nil {
		return nil, err
	}

	if err := u.middlewareRepo.AccessTokenSearch(ctx, cfg.Grpc.AuthUrl, accessToken); err != nil {
		return nil, err
	}

	c.Set("player_id", claims.PlayerID)
	c.Set("role_code", claims.RoleCode)

	return c, nil
}

func (u *middlewareUsecase) RbacAuthorization(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error) {
	ctx := c.Request().Context()

	playerRoleCode := c.Get("role_code").(int)

	rolesCount, err := u.middlewareRepo.RolesCount(ctx, cfg.Grpc.AuthUrl)
	if err != nil {
		return nil, err
	}

	playerRoleBinary := rbac.IntToBinary(playerRoleCode, int(rolesCount))

	for i := 0; i < int(rolesCount); i++ {
		if playerRoleBinary[i]&expected[i] == 1 {
			return c, nil
		}
	}

	return nil, errors.New("error: permission denied")
}
