package playerUsecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/chakornpat-tn/go-microservices/modules/player"
	playerPb "github.com/chakornpat-tn/go-microservices/modules/player/playerPb"
	"github.com/chakornpat-tn/go-microservices/modules/player/playerRepository"
	"github.com/chakornpat-tn/go-microservices/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	PlayerUsecaseService interface {
		CreatePlayer(pctx context.Context, req *player.CreatePlayerReq) (*player.PlayerProfile, error)
		FindOnePlayer(pctx context.Context, playerID string) (*player.PlayerProfile, error)
		AddPlayerMonney(pctx context.Context, req *player.CreatePlayerTransactionReq) (*player.PlayerSavingAccount, error)
		GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error)
		FindOnePlayerCredential(pctx context.Context, email, password string) (*playerPb.PlayerProfile, error)
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

func (u *playerUsecase) AddPlayerMonney(pctx context.Context, req *player.CreatePlayerTransactionReq) (*player.PlayerSavingAccount, error) {
	if err := u.playerRepo.InsertOnePlayerTranscation(pctx, &player.PlayerTransactions{
		PlayerID:  req.PlayerID,
		Amount:    req.Amount,
		CreatedAt: utils.LocalTime(),
	}); err != nil {
		return nil, err
	}

	return u.playerRepo.GetPlayerSavingAccount(pctx, req.PlayerID)
}

func (u *playerUsecase) GetPlayerSavingAccount(pctx context.Context, playerId string) (*player.PlayerSavingAccount, error) {
	return u.playerRepo.GetPlayerSavingAccount(pctx, playerId)
}

func (u *playerUsecase) FindOnePlayerCredential(pctx context.Context, email, password string) (*playerPb.PlayerProfile, error) {
	result, err := u.playerRepo.FindOnePlayerCredential(pctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		log.Printf("Error: FindOnePlayerCredential: %s", err.Error())
		return nil, errors.New("error: email or password is invalid")
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	roleCode := 0
	for _, v := range result.PlayerRoles {
		roleCode += v.RoleCode
	}

	return &playerPb.PlayerProfile{
		Id:        result.ID.Hex(),
		Username:  result.Username,
		Email:     result.Email,
		RoleCode:  int32(roleCode),
		CreatedAt: result.CreatedAt.In(loc).String(),
		UpdatedAt: result.UpdatedAt.In(loc).String(),
	}, nil
}
