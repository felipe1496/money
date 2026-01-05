package tests

import (
	"errors"
	"testing"

	"github.com/felipe1496/open-wallet/internal/resources/users"
	"github.com/felipe1496/open-wallet/internal/resources/users/mocks"

	"github.com/stretchr/testify/assert"
)

func TestUsersUseCase_List(t *testing.T) {
	t.Run("should list users successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockUsersRepo)
		uc := users.NewUsersUseCase(mockRepo)

		filter := users.UserFilter{Username: "alice"}

		expectedUsers := []users.User{
			{ID: "1", Username: "alice", Name: "Alice", Email: "alice@gmail.com"},
			{ID: "2", Username: "alice2", Name: "Alice2", Email: "alice2@gmail.com"},
		}

		// Repo returns users successfully
		mockRepo.
			On("ListUsers", filter).
			Return(expectedUsers, nil)

		result, err := uc.List(filter)

		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := new(mocks.MockUsersRepo)
		uc := users.NewUsersUseCase(mockRepo)

		filter := users.UserFilter{Email: "john@gmail.com"}

		mockRepo.
			On("ListUsers", filter).
			Return([]users.User{}, errors.New("db exploded"))

		result, err := uc.List(filter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, users.FailedToFetchUsersError)
		mockRepo.AssertExpectations(t)
	})
}

func TestUsersUseCase_Create(t *testing.T) {
	t.Run("should return error if username is already taken", func(t *testing.T) {
		mockRepo := new(mocks.MockUsersRepo)
		uc := users.NewUsersUseCase(mockRepo)

		input := users.CreateUserInput{Username: "johndoethegreat", Name: "John", Email: "john@gmail.com"}

		mockRepo.
			On("ListUsers", users.UserFilter{Username: input.Username}).
			Return([]users.User{
				{
					ID:       "1",
					Username: "johndoethegreat",
					Name:     "Urek",
					Email:    "urek@gmail.com",
				},
			}, nil)

		result, err := uc.Create(input)

		assert.Equal(t, users.User{}, result)
		assert.ErrorIs(t, err, users.UsernameAlreadyExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error if email already exists", func(t *testing.T) {
		mockRepo := new(mocks.MockUsersRepo)
		uc := users.NewUsersUseCase(mockRepo)

		input := users.CreateUserInput{
			Username: "johndoethegreat",
			Name:     "John",
			Email:    "john@gmail.com",
		}

		mockRepo.
			On("ListUsers", users.UserFilter{Username: input.Username}).
			Return([]users.User{}, nil)

		mockRepo.
			On("ListUsers", users.UserFilter{Email: input.Email}).
			Return([]users.User{
				{
					ID:       "1",
					Username: "rolling_stone",
					Name:     "Urek",
					Email:    "john@gmail.com",
				},
			}, nil)

		result, err := uc.Create(input)

		assert.Equal(t, users.User{}, result)
		assert.ErrorIs(t, err, users.EmailAlreadyExists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should create user successfully", func(t *testing.T) {
		mockRepo := new(mocks.MockUsersRepo)
		uc := users.NewUsersUseCase(mockRepo)

		input := users.CreateUserInput{
			Username: "johndoethegreat",
			Name:     "John",
			Email:    "john@gmail.com",
		}

		expectedUser := users.User{
			ID:       "123",
			Username: input.Username,
			Name:     input.Name,
			Email:    input.Email,
		}

		mockRepo.
			On("ListUsers", users.UserFilter{Username: input.Username}).
			Return([]users.User{}, nil)

		mockRepo.
			On("ListUsers", users.UserFilter{Email: input.Email}).
			Return([]users.User{}, nil)

		mockRepo.
			On("CreateUser", input).
			Return(expectedUser, nil)

		result, err := uc.Create(input)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

}
