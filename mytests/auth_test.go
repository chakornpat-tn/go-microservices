package mytests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/auth"
	"github.com/chakornpat-tn/go-microservices/modules/auth/authRepository"
	"github.com/chakornpat-tn/go-microservices/modules/auth/authUsecase"
	"github.com/chakornpat-tn/go-microservices/modules/player"
	playerPb "github.com/chakornpat-tn/go-microservices/modules/player/playerPb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	testLogin struct {
		name     string
		ctx      context.Context
		cfg      *config.Config
		req      *auth.PlayerLoginReq
		expected *auth.ProfileIntercepter
		err      error
		isErr    bool
		mock     func()
	}
)

func TestLogin(t *testing.T) {
	repoMock := new(authRepository.AuthRepositoryMock)
	usecase := authUsecase.NewAuthUsecase(repoMock)

	cfg := NewTestConfig()
	ctx := context.Background()

	credentialIDSuccess := bson.NewObjectID()

	tests := []testLogin{
		{
			name:  "Success",
			ctx:   ctx,
			cfg:   cfg,
			isErr: false,
			req: &auth.PlayerLoginReq{
				Email:    "success@gmail.com",
				Password: "123456",
			},
			expected: &auth.ProfileIntercepter{
				PlayerProfile: &player.PlayerProfile{
					ID:        "player:001",
					Email:     "success@gmail.com",
					Username:  "player001",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				Credential: &auth.CredentialRes{
					ID:           credentialIDSuccess.Hex(),
					PlayerID:     "player:001",
					RoleCode:     0,
					AccessToken:  "accessToken",
					RefreshToken: "refreshToken",
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
				},
			},
			mock: func() {
				repoMock.On("CredentialSearch", ctx, cfg.Grpc.PlayerUrl, &playerPb.CredentialSearchReq{
					Email:    "success@gmail.com",
					Password: "123456",
				}).Return(&playerPb.PlayerProfile{
					Id:        "001",
					Email:     "success@gmail.com",
					Username:  "player001",
					RoleCode:  0,
					CreatedAt: "0001-01-01 00:00:00 +0000 UTC",
					UpdatedAt: "0001-01-01 00:00:00 +0000 UTC",
				}, nil).Once()

				repoMock.On("AccessToken", cfg, mock.AnythingOfType("*jwtauth.Claims")).Return("accessToken").Once()
				repoMock.On("RefreshToken", cfg, mock.AnythingOfType("*jwtauth.Claims")).Return("refreshToken").Once()

				repoMock.On("InsertOnePlayerCredential", ctx, mock.AnythingOfType("*auth.Credential")).Return(credentialIDSuccess, nil).Once()
				repoMock.On("FindOnePlayerCredential", ctx, credentialIDSuccess.Hex()).Return(&auth.Credential{
					ID:           credentialIDSuccess,
					PlayerID:     "player:001",
					RoleCode:     0,
					AccessToken:  "accessToken",
					RefreshToken: "refreshToken",
					CreatedAt:    time.Time{},
					UpdatedAt:    time.Time{},
				}, nil).Once()
			},
		},
		{
			name:  "Error: CredentialSearch failed",
			ctx:   ctx,
			cfg:   cfg,
			isErr: true,
			req: &auth.PlayerLoginReq{
				Email:    "failed@gmail.com",
				Password: "123456",
			},
			expected: nil,
			err:      errors.New("error: email or password is incorrect"),
			mock: func() {
				repoMock.On("CredentialSearch", ctx, cfg.Grpc.PlayerUrl, &playerPb.CredentialSearchReq{
					Email:    "failed@gmail.com",
					Password: "123456",
				}).Return(nil, errors.New("error: email or password is incorrect")).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := usecase.Login(test.ctx, test.cfg, test.req)

			if test.isErr {
				assert.Error(t, err)
				assert.Equal(t, test.err, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				result.PlayerProfile.CreatedAt = time.Time{}
				result.PlayerProfile.UpdatedAt = time.Time{}
				result.Credential.CreatedAt = time.Time{}
				result.Credential.UpdatedAt = time.Time{}
				assert.Equal(t, test.expected, result)
			}
			repoMock.AssertExpectations(t)
		})
	}
}
