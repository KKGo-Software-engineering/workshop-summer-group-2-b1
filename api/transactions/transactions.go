package transactions

import "time"

type Transaction struct {
	ID               int64     `json:"id"`
	Date             time.Time `json:"date"`
	Amount           float64   `json:"amount"`
	Category         string    `json:"category"`
	Transaction_type string    `json:"transaction_type"`
	Note             string    `json:"note"`
	Image_url        string    `json:"image_url"`
}
