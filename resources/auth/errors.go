package auth

import "errors"

var (
	GoogleAuthFailedErr       = errors.New("authentication with Google failed")
	GoogleDintProvideEmailErr = errors.New("google did not provide an email")
	JwtGenErr                 = errors.New("failed to generate JWT token")
	GoogleEmailNotVerifiedErr = errors.New("google email is not verified")
)
