package auth

import (
	"net/http"
	"rango-backend/utils"
)

var (
	GoogleAuthFailedErr       = utils.NewHTTPError(http.StatusUnauthorized, "authentication with Google failed")
	GoogleDintProvideEmailErr = utils.NewHTTPError(http.StatusUnauthorized, "google did not provide an email")
	JwtGenErr                 = utils.NewHTTPError(http.StatusUnauthorized, "failed to generate JWT token")
	GoogleEmailNotVerifiedErr = utils.NewHTTPError(http.StatusUnauthorized, "Google email not verified")
)
