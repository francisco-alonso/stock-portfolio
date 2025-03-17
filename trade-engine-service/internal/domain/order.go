package domain

type Order struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Asset    string  `json:"asset"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}
