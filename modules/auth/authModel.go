package auth

import (
	"time"

	"github.com/chakornpat-tn/go-microservices/modules/player"
)

type (
	PlayerLoginReq struct {
		Email    string `json:"email" form:"email" validate:"required,email,max=255"`
		Password string `json:"password" form:"password" validate:"required,max=255"`
	}

	RefreshTokenReq struct {
		RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required,max=500"`
	}

	InsertPlayerRole struct {
		PlayerID string `json:"player_id" validate:"required"`
		RoleCode int    `json:"role_id" validate:"required"`
	}

	ProfileIntercepter struct {
		*player.PlayerProfile
		Credential *CredentialRes `json:"credential"`
	}

	CredentialRes struct {
		ID           string    `json:"_id" bson:"_id,omitempty"`
		PlayerID     string    `json:"player_id" bson:"player_id"`
		RoleCode     int       `json:"role_id" bson:"role_id"`
		AccessToken  string    `json:"access_token" bson:"access_token"`
		RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}
)
