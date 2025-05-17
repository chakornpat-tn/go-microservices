package authUsecase

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/auth"
	"github.com/chakornpat-tn/go-microservices/modules/auth/authRepository"
	"github.com/chakornpat-tn/go-microservices/modules/player"
	playerPb "github.com/chakornpat-tn/go-microservices/modules/player/playerPb"
	"github.com/chakornpat-tn/go-microservices/pkg/jwtauth"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
)

type (
	AuthUsecaseService interface {
		Login(pctx context.Context, cfg *config.Config, req *auth.PlayerLoginReq) (*auth.ProfileIntercepter, error)
		RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileIntercepter, error)
	}

	authUsecase struct {
		authRepo authRepository.AuthRepositoryService
	}
)

func NewAuthUsecase(authRepo authRepository.AuthRepositoryService) AuthUsecaseService {
	return &authUsecase{
		authRepo: authRepo,
	}
}

func (u *authUsecase) Login(pctx context.Context, cfg *config.Config, req *auth.PlayerLoginReq) (*auth.ProfileIntercepter, error) {
	profile, err := u.authRepo.CredentialSearch(pctx, cfg.Grpc.PlayerUrl, &playerPb.CredentialSearchReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	profile.Id = "player:" + profile.Id
	accessToken := jwtauth.NewAccessToken(cfg.Jwt.AccessTokenSecretKey, cfg.Jwt.AccessDuration, &jwtauth.Claims{
		PlayerID: profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	refreshToken := jwtauth.NewAccessToken(cfg.Jwt.RefreshTokenSecretKey, cfg.Jwt.RefreshDuration, &jwtauth.Claims{
		PlayerID: profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	credentialID, err := u.authRepo.InsertOnePlayerCredential(pctx, &auth.Credential{
		PlayerID:     profile.Id,
		RoleCode:     int(profile.RoleCode),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		CreatedAt:    utils.LocalTime(),
		UpdatedAt:    utils.LocalTime(),
	})

	credential, err := u.authRepo.FindOnePlayerCredential(pctx, credentialID.Hex())
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileIntercepter{
		PlayerProfile: &player.PlayerProfile{
			ID:        profile.Id,
			Username:  profile.Username,
			Email:     profile.Email,
			CreatedAt: utils.ConvertStringTimeToTime(profile.CreatedAt).In(loc),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.UpdatedAt).In(loc),
		},
		Credential: &auth.CredentialRes{
			ID:           credential.ID.Hex(),
			PlayerID:     credential.PlayerID,
			RoleCode:     credential.RoleCode,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt.In(loc),
			UpdatedAt:    credential.UpdatedAt.In(loc),
		},
	}, nil
}

func (u *authUsecase) RefreshToken(pctx context.Context, cfg *config.Config, req *auth.RefreshTokenReq) (*auth.ProfileIntercepter, error) {
	claims, err := jwtauth.ParseToken(cfg.Jwt.RefreshTokenSecretKey, req.RefreshToken)
	if err != nil {
		log.Printf("Error: ParseToken failed: %s", err)
		return nil, errors.New(err.Error())
	}

	profile, err := u.authRepo.FindOnePlayerProfileToRefresh(pctx, cfg.Grpc.PlayerUrl, &playerPb.FindOnePlayerProfileToRefreshReq{
		PlayerId: strings.TrimPrefix(claims.PlayerID, "player:"),
	})
	if err != nil {
		return nil, err
	}

	accessToken := jwtauth.NewAccessToken(cfg.Jwt.AccessTokenSecretKey, cfg.Jwt.AccessDuration, &jwtauth.Claims{
		PlayerID: profile.Id,
		RoleCode: int(profile.RoleCode),
	}).SignToken()

	refreshToken := jwtauth.ReloadToken(cfg.Jwt.RefreshTokenSecretKey, claims.ExpiresAt.Unix(), &jwtauth.Claims{
		PlayerID: profile.Id,
		RoleCode: int(profile.RoleCode),
	})

	if err := u.authRepo.UpdateOnePlayerCredential(pctx, req.CredentialId, &auth.UpdateRefreshTokenReq{
		PlayerId:     profile.Id,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UpdatedAt:    utils.LocalTime(),
	}); err != nil {
		return nil, err
	}

	credential, err := u.authRepo.FindOnePlayerCredential(pctx, req.CredentialId)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &auth.ProfileIntercepter{
		PlayerProfile: &player.PlayerProfile{
			ID:        "player:" + profile.Id,
			Email:     profile.Email,
			Username:  profile.Username,
			CreatedAt: utils.ConvertStringTimeToTime(profile.CreatedAt).In(loc),
			UpdatedAt: utils.ConvertStringTimeToTime(profile.UpdatedAt).In(loc),
		},
		Credential: &auth.CredentialRes{
			ID:           credential.ID.Hex(),
			PlayerID:     credential.PlayerID,
			RoleCode:     credential.RoleCode,
			AccessToken:  credential.AccessToken,
			RefreshToken: credential.RefreshToken,
			CreatedAt:    credential.CreatedAt.In(loc),
			UpdatedAt:    credential.UpdatedAt.In(loc),
		},
	}, nil
}
