package auth

import (
	"log"
	"rango-backend/db"
	"rango-backend/services"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	db, err := db.Conn("postgres://docker:docker@localhost:5432/docker?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	handler := NewHandler(db, services.NewGoogleService(), services.NewJWTService())
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/login/google", handler.LoginGoogle)
	}
}
