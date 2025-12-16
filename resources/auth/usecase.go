package auth

import (
	"rango-backend/resources/users"
	"rango-backend/services"
)

type AuthUseCase interface {
	LoginWithGoogle(accessToken string) (users.User, error)
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

func (uc *AuthUseCaseImpl) LoginWithGoogle(accessToken string) (users.User, error) {
	userInfo, err := uc.googleService.GetUserInfo(accessToken)

	if err != nil {
		return users.User{}, GoogleAuthFailedErr
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
