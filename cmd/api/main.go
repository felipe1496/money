package main

import (
	"log"
	"os"
	"time"

	docs "github.com/felipe1496/open-wallet/docs"

	"github.com/felipe1496/open-wallet/internal/resources/auth"
	"github.com/felipe1496/open-wallet/internal/resources/categories"
	"github.com/felipe1496/open-wallet/internal/resources/transactions"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Money API
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	// add swagger
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	/* origins := os.Getenv("ORIGINS")
	if origins == "" {
		log.Fatal("ORIGINS cannot be empty")
	}

	originsList := strings.Split(origins, ",") */

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://openwallet.vercel.app"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth.Router(r)
	transactions.Router(r)
	categories.Router(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
