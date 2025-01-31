package repository

import (
	"context"
	"fmt"
	"github.com/Egorpalan/finance-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryInterface interface {
	GetUserByID(userID int64) (*model.User, error)
	UpdateUserBalance(user *model.User) error
	AddTransaction(transaction *model.Transaction) error
	GetUserTransactions(userID int64) ([]model.Transaction, error)
	CreateUser(balance float64) (*model.User, error)
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByID(userID int64) (*model.User, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID: %d", userID)
	}
	var user model.User
	query := "SELECT id, balance FROM users WHERE id=$1"
	err := r.db.QueryRow(context.Background(), query, userID).Scan(&user.ID, &user.Balance)
	if err != nil {
		return nil, fmt.Errorf("could not find user with id %d: %v", userID, err)
	}
	return &user, nil
}

func (r *Repository) UpdateUserBalance(user *model.User) error {
	query := "UPDATE users SET balance=$1 WHERE id=$2"
	_, err := r.db.Exec(context.Background(), query, user.Balance, user.ID)
	return err
}

func (r *Repository) AddTransaction(transaction *model.Transaction) error {
	query := `INSERT INTO transactions (sender_id, receiver_id, amount, transaction_type, created_at)
				VALUES ($1, $2, $3, $4, current_timestamp)`
	_, err := r.db.Exec(context.Background(), query, transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.TransactionType)
	return err
}

func (r *Repository) GetUserTransactions(userID int64) ([]model.Transaction, error) {
	var transactions []model.Transaction
	query := `SELECT id, sender_id, receiver_id, amount, transaction_type, created_at
              FROM transactions WHERE sender_id=$1 OR receiver_id=$1 ORDER BY created_at DESC LIMIT 10`
	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.SenderID, &t.ReceiverID, &t.Amount, &t.TransactionType, &t.Timestamp); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *Repository) CreateUser(balance float64) (*model.User, error) {
	var user model.User
	query := `INSERT INTO users (balance) VALUES ($1) RETURNING id, balance`
	err := r.db.QueryRow(context.Background(), query, balance).Scan(&user.ID, &user.Balance)
	if err != nil {
		return nil, fmt.Errorf("could not create user: %v", err)
	}
	return &user, nil
}
