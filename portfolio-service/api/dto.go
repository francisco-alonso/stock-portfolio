package api

type PositionDto struct {
	UserID   string  `json:"user_id"`
	Asset    string  `json:"asset"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}