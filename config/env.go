package config

import (
	"os"
	"strings"
)

type Variables struct {
	GoogleClientID     string
	GoogleClientSecret string
	JWTSecret          string
	DatabaseURL        string
	Port               string
	Origins            []string
}

func GetEnv() Variables {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	jwtSecret := os.Getenv("JWT_SECRET")
	databaseURL := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	origins := os.Getenv("ORIGINS")

	var originsList []string
	if origins != "" {
		originsList = strings.Split(origins, ",")
	}

	return Variables{
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		JWTSecret:          jwtSecret,
		DatabaseURL:        databaseURL,
		Port:               port,
		Origins:            originsList,
	}
}
