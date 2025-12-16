package services

import (
	"errors"
)

var (
	ErrGoogleNetwork    = errors.New("network error")
	ErrGoogleAuthFailed = errors.New("authentication with Google failed")
	ErrGoogleResponse   = errors.New("error reading Google response")
	ErrGoogleDecode     = errors.New("error decoding Google user JSON")
)
