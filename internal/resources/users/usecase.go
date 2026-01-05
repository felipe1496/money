package users

type UsersUseCase interface {
	List(filter UserFilter) ([]User, error)
	Create(input CreateUserInput) (User, error)
}

type UsersUseCaseImpl struct {
	repo UsersRepo
}

func NewUsersUseCase(repo UsersRepo) UsersUseCase {
	return &UsersUseCaseImpl{repo: repo}
}

func (uc *UsersUseCaseImpl) List(filter UserFilter) ([]User, error) {
	users, err := uc.repo.ListUsers(filter)

	if err != nil {
		return nil, FailedToFetchUsersError
	}

	return users, nil
}

func (uc *UsersUseCaseImpl) Create(input CreateUserInput) (User, error) {

	if userAlreadyExists, err := uc.List(UserFilter{Username: input.Username}); err == nil && len(userAlreadyExists) > 0 {
		return User{}, UsernameAlreadyExists
	}

	if userAlreadyExists, err := uc.List(UserFilter{Email: input.Email}); err == nil && len(userAlreadyExists) > 0 {
		return User{}, EmailAlreadyExists
	}

	user, err := uc.repo.CreateUser(input)
	if err != nil {
		return User{}, FailedToCreateUserError
	}

	return user, nil
}
