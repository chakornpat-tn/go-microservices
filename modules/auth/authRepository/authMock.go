package authRepository

import (
	"context"

	"github.com/chakornpat-tn/go-microservices/config"
	"github.com/chakornpat-tn/go-microservices/modules/auth"
	playerPb "github.com/chakornpat-tn/go-microservices/modules/player/playerPb"
	"github.com/chakornpat-tn/go-microservices/pkg/jwtauth"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AuthRepositoryMock struct {
	mock.Mock
}

func NewAuthRepositoryMock() AuthRepositoryService {
	return &AuthRepositoryMock{}
}

func (m *AuthRepositoryMock) CredentialSearch(pctx context.Context, grpcUrl string, req *playerPb.CredentialSearchReq) (*playerPb.PlayerProfile, error) {
	args := m.Called(pctx, grpcUrl, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*playerPb.PlayerProfile), args.Error(1)
}
func (m *AuthRepositoryMock) InsertOnePlayerCredential(pctx context.Context, req *auth.Credential) (bson.ObjectID, error) {
	args := m.Called(pctx, req)
	return args.Get(0).(bson.ObjectID), args.Error(1)
}
func (m *AuthRepositoryMock) FindOnePlayerCredential(pctx context.Context, credentialId string) (*auth.Credential, error) {
	args := m.Called(pctx, credentialId)
	return args.Get(0).(*auth.Credential), args.Error(1)
}
func (m *AuthRepositoryMock) FindOnePlayerProfileToRefresh(pctx context.Context, grpcUrl string, req *playerPb.FindOnePlayerProfileToRefreshReq) (*playerPb.PlayerProfile, error) {
	args := m.Called(pctx, grpcUrl, req)
	return args.Get(0).(*playerPb.PlayerProfile), args.Error(1)
}
func (m *AuthRepositoryMock) UpdateOnePlayerCredential(pctx context.Context, credentialId string, req *auth.UpdateRefreshTokenReq) error {
	args := m.Called(pctx, credentialId, req)
	return args.Error(0)
}
func (m *AuthRepositoryMock) DeleteOnePlayerCredential(pctx context.Context, credentialID string) (int64, error) {
	args := m.Called(pctx, credentialID)
	return int64(args.Int(0)), args.Error(1)
}
func (m *AuthRepositoryMock) FindOneAccessToken(pctx context.Context, accessToken string) (*auth.Credential, error) {
	args := m.Called(pctx, accessToken)
	return args.Get(0).(*auth.Credential), args.Error(1)
}
func (m *AuthRepositoryMock) RolesCount(pctx context.Context) (int64, error) {
	args := m.Called(pctx)
	return int64(args.Int(0)), args.Error(1)
}

func (m *AuthRepositoryMock) RefreshToken(cfg *config.Config, claims *jwtauth.Claims) string {
	return m.Called(cfg, claims).String(0)
}

func (m *AuthRepositoryMock) AccessToken(cfg *config.Config, claims *jwtauth.Claims) string {
	return m.Called(cfg, claims).String(0)
}
