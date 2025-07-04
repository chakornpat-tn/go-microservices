package inventory

import "go.mongodb.org/mongo-driver/v2/bson"

type (
	Inventory struct {
		ID       bson.ObjectID `json:"_id" bson:"_id,omitempty" `
		PlayerID string        `json:"player_id" bson:"player_id" `
		ItemID   string        `json:"item_id" bson:"item_id" `
	}
)
