package services

import (
	"net/http"
	"rango-backend/utils"
)

var (
	FailedGoogleAuthenticationErr = utils.NewHTTPError(http.StatusUnauthorized, "google authentication failed")
)
