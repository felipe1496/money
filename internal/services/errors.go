package services

import (
	"net/http"

	"github.com/felipe1496/open-wallet/internal/utils"
)

var (
	FailedGoogleAuthenticationErr = utils.NewHTTPError(http.StatusUnauthorized, "google authentication failed")
)
