package auth

import (
	"fmt"
	"rango-backend/resources/users"
	"rango-backend/services"
	"strings"

	"github.com/oklog/ulid/v2"
)

type AuthUseCase interface {
	LoginWithGoogle(code string) (users.User, error)
}

type AuthUseCaseImpl struct {
	googleService services.GoogleService
	usersUseCase  users.UsersUseCase
}

func NewAuthUseCase(googleService services.GoogleService, usersUseCase users.UsersUseCase) AuthUseCase {
	return &AuthUseCaseImpl{
		googleService: googleService,
		usersUseCase:  usersUseCase,
	}
}

func (uc *AuthUseCaseImpl) LoginWithGoogle(code string) (users.User, error) {
	userAccessToken, err := uc.googleService.GetUserAccessToken(code)

	if err != nil {
		return users.User{}, err
	}

	userInfo, err := uc.googleService.GetUserInfo(*userAccessToken)

	if err != nil {
		return users.User{}, err
	}

	if !*userInfo.EmailVerified {
		return users.User{}, GoogleEmailNotVerifiedErr
	}

	if userInfo.Email == nil {
		return users.User{}, GoogleDintProvideEmailErr
	}

	userExists, err := uc.usersUseCase.List(users.UserFilter{Email: *userInfo.Email})

	if err != nil {
		return users.User{}, err
	}

	var userRes users.User

	if len(userExists) == 0 {
		createUserInput := users.CreateUserInput{
			Name: userInfo.Name,
		}

		createUserInput.Email = *userInfo.Email

		createUserInput.AvatarURL = userInfo.Picture

		createUserInput.Username = fmt.Sprintf("%s_%s", strings.ToLower(strings.ReplaceAll(userInfo.Name, " ", "_")), ulid.Make().String())

		createdUser, err := uc.usersUseCase.Create(createUserInput)

		if err != nil {
			return users.User{}, err
		}

		userRes = createdUser
	} else {
		userRes = userExists[0]
	}

	return userRes, nil
}
