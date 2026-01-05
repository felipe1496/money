package users

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
)

var (
	FailedToFetchUsersError = utils.NewHTTPError(http.StatusInternalServerError, "failed to fetch users")
	FailedToCreateUserError = utils.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	UsernameAlreadyExists   = utils.NewHTTPError(http.StatusConflict, "user with this username already exists")
	EmailAlreadyExists      = utils.NewHTTPError(http.StatusConflict, "user with this email already exists")
)
