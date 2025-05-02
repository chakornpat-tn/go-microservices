package playerUsecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chakornpat-tn/go-microservices/modules/player"
	"github.com/chakornpat-tn/go-microservices/modules/player/playerRepository"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	PlayerUsecaseService interface {
		CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error)
		FindOnePlayer(pctx context.Context, playerID string) (*player.PlayerProfile, error)
	}

	playerUsecase struct {
		playerRepo playerRepository.PlayerRepositoryService
	}
)

func NewPlayerUsecase(playerRepo playerRepository.PlayerRepositoryService) PlayerUsecaseService {
	return &playerUsecase{
		playerRepo: playerRepo,
	}
}

func (u *playerUsecase) CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error) {

	if !u.playerRepo.IsUniquePlayer(pctx, req.Email, req.Username) {
		return nil, errors.New("erros: email or username already exits")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error: hash password failed")
	}

	playerID, err := u.playerRepo.InsertOnePlayer(pctx, &player.Player{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
		PlayerRoles: []player.PlayerRole{
			{
				RoleTitle: "player",
				RoleCode:  0,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return u.FindOnePlayer(pctx, playerID.Hex())
}

func (u *playerUsecase) FindOnePlayer(pctx context.Context, playerID string) (*player.PlayerProfile, error) {
	result, err := u.playerRepo.FindOnePlayer(pctx, playerID)
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Print("Error loading location")
		return nil, err
	}

	return &player.PlayerProfile{
		ID:        result.ID.Hex(),
		Username:  result.Username,
		Email:     result.Email,
		CreatedAt: result.CreatedAt.In(loc),
		UpdatedAt: result.UpdatedAt.In(loc),
	}, nil
}
