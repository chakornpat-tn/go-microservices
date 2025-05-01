package playerUsecase

import (
	"context"
	"errors"

	"github.com/chakornpat-tn/go-microservices/modules/player"
	"github.com/chakornpat-tn/go-microservices/modules/player/playerRepository"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	PlayerUsecaseService interface {
		CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (string, error)
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

func (u *playerUsecase) CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (string, error) {

	if !u.playerRepo.IsUniquePlayer(pctx, req.Email, req.Username) {
		return "", errors.New("erros: email or username already exits")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("error: hash password failed")
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
		return "", err
	}

	return playerID.Hex(), nil
}
