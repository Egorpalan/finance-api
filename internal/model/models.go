package model

import "time"

type Transaction struct {
	ID              int64     `json:"id"`
	SenderID        int64     `json:"sender_id"`
	ReceiverID      int64     `json:"receiver_id"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"`
	Timestamp       time.Time `json:"timestamp"`
}

type User struct {
	ID      int64   `json:"id"`
	Balance float64 `json:"balance"`
}
