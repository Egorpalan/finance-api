package tests

import (
	"context"
	"testing"
	"time"

	"github.com/Egorpalan/finance-api/internal/model"
	"github.com/Egorpalan/finance-api/internal/service"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTx struct {
	mock.Mock
}

func (m *MockTx) Begin(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockTx) Commit(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Rollback(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTx) Conn() *pgx.Conn {
	args := m.Called()
	return args.Get(0).(*pgx.Conn)
}

func (m *MockTx) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	args := m.Called(ctx, sql, arguments)
	return args.Get(0).(pgconn.CommandTag), args.Error(1)
}

func (m *MockTx) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	panic("not implemented")
}

func (m *MockTx) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	panic("not implemented")
}

func (m *MockTx) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	panic("not implemented")
}

func (m *MockTx) SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults {
	panic("not implemented")
}

func (m *MockTx) LargeObjects() pgx.LargeObjects {
	panic("not implemented")
}

func (m *MockTx) Prepare(ctx context.Context, name, sql string) (*pgconn.StatementDescription, error) {
	panic("not implemented")
}

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, userID int64) (*model.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) UpdateUserBalance(ctx context.Context, tx pgx.Tx, user *model.User) error {
	args := m.Called(ctx, tx, user)
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

func TestTopUpBalance(t *testing.T) {
	mockRepo := new(MockRepository)
	mockTx := new(MockTx)
	service := service.NewService(mockRepo)

	userID := int64(1)
	amount := 100.0
	user := &model.User{ID: userID, Balance: 500.0}

	mockRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("GetUserByID", mock.Anything, userID).Return(user, nil)
	mockRepo.On("UpdateUserBalance", mock.Anything, mockTx, user).Return(nil)
	mockRepo.On("AddTransaction", mock.Anything, mockTx, mock.Anything).Return(nil)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	err := service.TopUpBalance(userID, amount)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestTransferMoney(t *testing.T) {
	mockRepo := new(MockRepository)
	mockTx := new(MockTx)
	service := service.NewService(mockRepo)

	senderID := int64(1)
	receiverID := int64(2)
	amount := 50.0
	sender := &model.User{ID: senderID, Balance: 100.0}
	receiver := &model.User{ID: receiverID, Balance: 200.0}

	mockRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)
	mockRepo.On("GetUserByID", mock.Anything, senderID).Return(sender, nil)
	mockRepo.On("GetUserByID", mock.Anything, receiverID).Return(receiver, nil)
	mockRepo.On("UpdateUserBalance", mock.Anything, mockTx, sender).Return(nil)
	mockRepo.On("UpdateUserBalance", mock.Anything, mockTx, receiver).Return(nil)
	mockRepo.On("AddTransaction", mock.Anything, mockTx, mock.Anything).Return(nil)
	mockTx.On("Commit", mock.Anything).Return(nil)
	mockTx.On("Rollback", mock.Anything).Return(nil)

	err := service.TransferMoney(senderID, receiverID, amount)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestGetTransactions(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo)

	userID := int64(1)
	transactions := []model.Transaction{
		{ID: 1, SenderID: 1, ReceiverID: 2, Amount: 50.0, TransactionType: "transfer", Timestamp: time.Now()},
	}

	mockRepo.On("GetUserTransactions", mock.Anything, userID).Return(transactions, nil)

	result, err := service.GetTransactions(userID)

	assert.NoError(t, err)
	assert.Equal(t, transactions, result)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := service.NewService(mockRepo)

	balance := 100.0
	user := &model.User{ID: 1, Balance: balance}

	mockRepo.On("CreateUser", mock.Anything, balance).Return(user, nil)

	result, err := service.CreateUser(balance)

	assert.NoError(t, err)
	assert.Equal(t, user, result)
	mockRepo.AssertExpectations(t)
}
