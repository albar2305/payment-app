package model

import "time"

type Transaction struct {
	ID                 string    `json:"id"`
	SenderCustomerId   string    `json:"sender_customer_id"`
	ReceiverMerchantId string    `json:"receiver_merchant_id"`
	Amount             int64     `json:"amount"`
	CreatedAt          time.Time `json:"created_at"`
}

type CreateTransactionRequest struct {
	UserId             string `json:"user_id"`
	ReceiverMerchantId string `json:"receiver_merchant_id"`
	Amount             int64  `json:"amount"`
}

type TransactionResponse struct {
	ID        string           `json:"id"`
	Customer  CustomerResponse `json:"customer"`
	Merchant  Merchant         `json:"merchant"`
	Amount    int64            `json:"amount"`
	CreatedAt time.Time        `json:"created_at"`
}
