package mocks

import (
	"github.com/felipe1496/open-wallet/internal/resources/transactions"
	"github.com/felipe1496/open-wallet/internal/utils"

	"github.com/stretchr/testify/mock"
)

type MockTransactionsRepo struct {
	mock.Mock
}

func (m *MockTransactionsRepo) CreateEntry(input transactions.CreateEntryDTO, db utils.Executer) (transactions.Entry, error) {
	args := m.Called(input, db)
	return args.Get(0).(transactions.Entry), args.Error(1)
}

func (m *MockTransactionsRepo) CreateTransaction(input transactions.CreateTransactionDTO, db utils.Executer) (transactions.Transaction, error) {
	args := m.Called(input, db)
	return args.Get(0).(transactions.Transaction), args.Error(1)
}

func (m *MockTransactionsRepo) ListTransactions(db utils.Executer) ([]transactions.Transaction, error) {
	args := m.Called(db)
	return args.Get(0).([]transactions.Transaction), args.Error(1)
}

func (m *MockTransactionsRepo) ListEntries(db utils.Executer) ([]transactions.Entry, error) {
	args := m.Called(db)
	return args.Get(0).([]transactions.Entry), args.Error(1)
}
