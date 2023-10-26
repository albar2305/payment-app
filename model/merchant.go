package model

import "time"

type Merchant struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	BusinesType string    `json:"busines_type"`
	Balance     int64     `json:"balance"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateMerchantRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	BusinesType string `json:"busines_type"`
	Balance     int64  `json:"balance"`
}
