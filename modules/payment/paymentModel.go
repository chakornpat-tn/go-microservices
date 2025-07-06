package payment

type (
	ItemServiceReq struct {
		Items []*ItemServiceReqDatum `json:"items" validate:"required"`
	}

	ItemServiceReqDatum struct {
		ItemID string  `json:"item_id" validate:"required, max=64"`
		Price  float64 `json:"price"`
	}

	PaymentTransferReq struct {
		PlayerID string  `json:"player_id"`
		ItemID   string  `json:"item_id"`
		Amount   float64 `json:"amount"`
	}

	PaymentTransferRes struct {
		InventoryID   string `json:"inventory_id"`
		TransactionID string `json:"transaction_id"`
		PlayerID      string `json:"player_id"`
		ItemID        string `json:"item_id"`
		Amount        string `json:"amount"`
		Error         string `json:"error"`
	}
)
