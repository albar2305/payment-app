package model

import "time"

type Customer struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateCustomerRequest struct {
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
	Balance int64  `json:"balance"`
}

type AddCustomerBalanceParams struct {
	Amount int64  `json:"amount"`
	ID     string `json:"id"`
}

type CustomerResponse struct {
	ID        string       `json:"id"`
	User      UserResponse `json:"user_id"`
	Name      string       `json:"name"`
	Balance   int64        `json:"balance"`
	CreatedAt time.Time    `json:"created_at"`
}
