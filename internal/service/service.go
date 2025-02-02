package service

import (
	"context"
	"errors"
	"github.com/Egorpalan/finance-api/internal/model"
	"github.com/Egorpalan/finance-api/internal/repository"
	"time"
)

type Service struct {
	repo repository.RepositoryInterface
}

func NewService(repo repository.RepositoryInterface) *Service {
	return &Service{repo: repo}
}

func (s *Service) TopUpBalance(userID int64, amount float64) error {
	ctx := context.Background()

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = s.repo.GetUserByIDForUpdate(ctx, tx, userID)
	if err != nil {
		return err
	}

	err = s.repo.UpdateUserBalance(ctx, tx, userID, amount)
	if err != nil {
		return err
	}

	transaction := &model.Transaction{
		SenderID:        1,
		ReceiverID:      userID,
		Amount:          amount,
		TransactionType: "top_up",
		Timestamp:       time.Now(),
	}
	err = s.repo.AddTransaction(ctx, tx, transaction)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) TransferMoney(senderID, receiverID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	ctx := context.Background()
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	sender, err := s.repo.GetUserByIDForUpdate(ctx, tx, senderID)
	if err != nil {
		return err
	}

	_, err = s.repo.GetUserByIDForUpdate(ctx, tx, receiverID)
	if err != nil {
		return err
	}

	if sender.Balance < amount {
		return errors.New("insufficient balance")
	}

	err = s.repo.UpdateUserBalance(ctx, tx, senderID, -amount)
	if err != nil {
		return err
	}

	err = s.repo.UpdateUserBalance(ctx, tx, receiverID, amount)
	if err != nil {
		return err
	}

	transactionSender := &model.Transaction{
		SenderID:        senderID,
		ReceiverID:      receiverID,
		Amount:          -amount,
		TransactionType: "transfer",
		Timestamp:       time.Now(),
	}
	transactionReceiver := &model.Transaction{
		SenderID:        senderID,
		ReceiverID:      receiverID,
		Amount:          amount,
		TransactionType: "transfer",
		Timestamp:       time.Now(),
	}

	if err := s.repo.AddTransaction(ctx, tx, transactionSender); err != nil {
		return err
	}
	if err := s.repo.AddTransaction(ctx, tx, transactionReceiver); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *Service) GetTransactions(userID int64) ([]model.Transaction, error) {
	ctx := context.Background()
	return s.repo.GetUserTransactions(ctx, userID)
}

func (s *Service) CreateUser(balance float64) (*model.User, error) {
	ctx := context.Background()
	return s.repo.CreateUser(ctx, balance)
}
