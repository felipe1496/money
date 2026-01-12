package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/felipe1496/open-wallet/config"
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
		log.Println("Nenhum arquivo .env encontrado")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetEnv().Origins,
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth.Router(r)
	transactions.Router(r)
	categories.Router(r)

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	if port == "" {
		port = ":8080"
	}

	r.Run(port)
}
