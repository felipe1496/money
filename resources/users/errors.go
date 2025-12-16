package users

import "errors"

var (
	FailedToFetchUsersError = errors.New("failed to fetch users")
	FailerToCreateUserError = errors.New("failed to create user")
	UsernameAlreadyExists   = errors.New("username already taken")
	EmailAlreadyExists      = errors.New("user with this email already exists")
)
