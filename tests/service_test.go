package tests

import (
	"github.com/Egorpalan/finance-api/internal/model"
	"github.com/Egorpalan/finance-api/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByID(userID int64) (*model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) UpdateUserBalance(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockRepository) AddTransaction(transaction *model.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockRepository) GetUserTransactions(userID int64) ([]model.Transaction, error) {
	args := m.Called(userID)
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *MockRepository) CreateUser(balance float64) (*model.User, error) {
	args := m.Called(balance)
	return args.Get(0).(*model.User), args.Error(1)
}

func TestTopUpBalance(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo)

	mockRepo.On("GetUserByID", int64(1)).Return(&model.User{ID: 1, Balance: 100.00}, nil)
	mockRepo.On("UpdateUserBalance", mock.Anything).Return(nil)
	mockRepo.On("AddTransaction", mock.Anything).Return(nil)

	err := service.TopUpBalance(1, 50.00)

	mockRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestTransferMoney(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo)

	mockRepo.On("GetUserByID", int64(1)).Return(&model.User{ID: 1, Balance: 100.00}, nil)
	mockRepo.On("GetUserByID", int64(2)).Return(&model.User{ID: 2, Balance: 50.00}, nil)
	mockRepo.On("UpdateUserBalance", mock.Anything).Return(nil)
	mockRepo.On("AddTransaction", mock.Anything).Return(nil)

	err := service.TransferMoney(1, 2, 30.00)

	mockRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestTransferMoney_InsufficientBalance(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo)

	mockRepo.On("GetUserByID", int64(1)).Return(&model.User{ID: 1, Balance: 10.00}, nil)
	mockRepo.On("GetUserByID", int64(2)).Return(&model.User{ID: 2, Balance: 50.00}, nil)

	err := service.TransferMoney(1, 2, 30.00)

	mockRepo.AssertExpectations(t)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "insufficient balance")
}

func TestGetTransactions(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo)

	mockRepo.On("GetUserTransactions", int64(1)).Return([]model.Transaction{
		{ID: 1, SenderID: 1, ReceiverID: 2, Amount: -30.00, TransactionType: "transfer", Timestamp: time.Now()},
		{ID: 2, SenderID: 1, ReceiverID: 1, Amount: 50.00, TransactionType: "top_up", Timestamp: time.Now()},
	}, nil)

	transactions, err := service.GetTransactions(1)

	mockRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
	assert.Equal(t, transactions[0].Amount, -30.00)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo)

	mockRepo.On("CreateUser", 100.00).Return(&model.User{ID: 1, Balance: 100.00}, nil)

	user, err := service.CreateUser(100.00)

	mockRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, int64(1))
	assert.Equal(t, user.Balance, 100.00)
}
