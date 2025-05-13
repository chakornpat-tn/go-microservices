package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	Credential struct {
		ID           bson.ObjectID `bson:"_id,omitempty"`
		PlayerID     string        `bson:"player_id"`
		RoleCode     int           `bson:"role_id"`
		AccessToken  string        `bson:"access_token"`
		RefreshToken string        `bson:"refresh_token"`
		CreatedAt    time.Time     `bson:"created_at"`
		UpdatedAt    time.Time     `bson:"updated_at"`
	}

	Role struct {
		ID    bson.ObjectID `json:"_id" bson:"_id,omitempty"`
		Title string        `json:"title" bson:"title"`
		Code  int           `json:"code" bson:"code"`
	}
)
