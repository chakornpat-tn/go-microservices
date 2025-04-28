package payment

type (
	ItemServiceReq struct {
		Items []*ItemServiceReqDatum `json:"items" validate:"required"`
	}

	ItemServiceReqDatum struct {
		ItemID string  `json:"item_id" validate:"required, max=64"`
		Price  float64 `json:"price"`
	}
)
