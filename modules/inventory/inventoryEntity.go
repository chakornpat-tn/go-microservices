package inventory

type (
	Inventory struct {
		ID       string `json:"_id" bson:"_id,omitempty" `
		PlayerID string `json:"player_id" bson:"player_id" `
		ItemID   string `json:"item_id" bson:"item_id" `
	}
)
