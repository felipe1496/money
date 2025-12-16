package auth

type LoginGoogleRequest struct {
	Code string `json:"code"`
}
