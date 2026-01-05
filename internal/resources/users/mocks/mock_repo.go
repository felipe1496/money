package mocks

import (
	"github.com/felipe1496/open-wallet/internal/resources/users"

	"github.com/stretchr/testify/mock"
)

type MockUsersRepo struct {
	mock.Mock
}

func (m *MockUsersRepo) ListUsers(filter users.UserFilter) ([]users.User, error) {
	args := m.Called(filter)
	return args.Get(0).([]users.User), args.Error(1)
}

func (m *MockUsersRepo) CreateUser(input users.CreateUserInput) (users.User, error) {
	args := m.Called(input)
	return args.Get(0).(users.User), args.Error(1)
}
