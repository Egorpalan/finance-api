package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Egorpalan/finance-api/internal/model"
	"github.com/Egorpalan/finance-api/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	tx, ok := args.Get(0).(pgx.Tx)
	if !ok {
		return nil, errors.New("failed to cast to pgx.Tx")
	}
	return tx, args.Error(1)
}

func (m *MockRepository) GetUserByIDForUpdate(ctx context.Context, tx pgx.Tx, userID int64) (*model.User, error) {
	args := m.Called(ctx, tx, userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) UpdateUserBalance(ctx context.Context, tx pgx.Tx, userID int64, amount float64) error {
	args := m.Called(ctx, tx, userID, amount)
	return args.Error(0)
}

func (m *MockRepository) AddTransaction(ctx context.Context, tx pgx.Tx, transaction *model.Transaction) error {
	args := m.Called(ctx, tx, transaction)
	return args.Error(0)
}

func (m *MockRepository) GetUserTransactions(ctx context.Context, userID int64) ([]model.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *MockRepository) CreateUser(ctx context.Context, balance float64) (*model.User, error) {
	args := m.Called(ctx, balance)
	return args.Get(0).(*model.User), args.Error(1)
}

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestGetTransactions(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := service.NewService(mockRepo)

	userID := int64(1)
	transactions := []model.Transaction{
		{SenderID: 1, ReceiverID: 2, Amount: 50, TransactionType: "transfer", Timestamp: time.Now()},
		{SenderID: 1, ReceiverID: 3, Amount: 30, TransactionType: "transfer", Timestamp: time.Now()},
	}

	mockRepo.On("GetUserTransactions", mock.Anything, userID).Return(transactions, nil)

	result, err := svc.GetTransactions(userID)

	assert.NoError(t, err)

	assert.Equal(t, transactions, result)

	mockRepo.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	svc := service.NewService(mockRepo)

	balance := 100.0
	newUser := &model.User{ID: 1, Balance: balance}

	mockRepo.On("CreateUser", mock.Anything, balance).Return(newUser, nil)

	result, err := svc.CreateUser(balance)

	assert.NoError(t, err)

	assert.Equal(t, newUser, result)

	mockRepo.AssertExpectations(t)
}
