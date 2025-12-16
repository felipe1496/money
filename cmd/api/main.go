package main

import (
	"log"
	"rango-backend/resources/auth"
	"rango-backend/resources/transactions"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Println("Nenhum arquivo .env encontrado")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth.Router(r)
	transactions.Router(r)

	r.Run(":8080")
}
