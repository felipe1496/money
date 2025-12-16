package transactions

import (
	"log"
	"rango-backend/db"
	"rango-backend/middlewares"
	"rango-backend/services"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	db, err := db.Conn("postgres://docker:docker@localhost:5432/docker?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	handler := NewHandler(db)
	transactionsGroup := router.Group("/api/v1/transactions")
	{
		transactionsGroup.POST("/simple-expense",
			middlewares.RequireAuthMiddleware(services.NewJWTService()),
			handler.CreateSimpleExpense)
		transactionsGroup.GET("/entries/:period",
			middlewares.RequireAuthMiddleware(services.NewJWTService()),
			middlewares.QueryOptsMiddleware(),
			handler.ListViewEntries)
	}
}
