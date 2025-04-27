package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Credential struct {
		ID           primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		PlayerID     string             `json:"player_id" bson:"player_id"`
		RoleCode     int                `json:"role_id" bson:"role_id"`
		AccessToken  string             `json:"access_token" bson:"access_token"`
		RefreshToken string             `json:"refresh_token" bson:"refresh_token"`
		CreatedAt    time.Time          `json:"created_at"`
		UpdatedAt    time.Time          `json:"updated_at"`
	}

	Role struct {
		ID    primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Title string             `json:"title" bson:"title"`
		Code  int                `json:"code" bson:"code"`
	}
)
