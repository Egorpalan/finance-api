package service

import (
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
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	user.Balance += amount
	err = s.repo.UpdateUserBalance(user)
	if err != nil {
		return err
	}

	transaction := &model.Transaction{
		SenderID:        1,
		ReceiverID:      userID,
		Amount:          amount,
		TransactionType: "top_up",
	}
	err = s.repo.AddTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) TransferMoney(senderID, receiverID int64, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	sender, err := s.repo.GetUserByID(senderID)
	if err != nil {
		return err
	}

	receiver, err := s.repo.GetUserByID(receiverID)
	if err != nil {
		return err
	}

	if sender.Balance < amount {
		return errors.New("insufficient balance")
	}

	sender.Balance -= amount
	receiver.Balance += amount

	if err := s.repo.UpdateUserBalance(sender); err != nil {
		return err
	}
	if err := s.repo.UpdateUserBalance(receiver); err != nil {
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

	if err := s.repo.AddTransaction(transactionSender); err != nil {
		return err
	}
	if err := s.repo.AddTransaction(transactionReceiver); err != nil {
		return err
	}

	return nil
}

func (s *Service) GetTransactions(userID int64) ([]model.Transaction, error) {
	transactions, err := s.repo.GetUserTransactions(userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *Service) CreateUser(balance float64) (*model.User, error) {
	return s.repo.CreateUser(balance)
}
