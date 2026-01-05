package transactions

import (
	"log"
	"os"

	"github.com/felipe1496/open-wallet/db"

	"github.com/felipe1496/open-wallet/internal/middlewares"
	"github.com/felipe1496/open-wallet/internal/services"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	db, err := db.Conn(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	jwtService := services.NewJWTService()
	handler := NewHandler(db)
	transactionsGroup := router.Group("/api/v1/transactions")
	{
		transactionsGroup.POST("/simple-expense",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.CreateSimpleExpense)
		transactionsGroup.GET("/entries/:period",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryOptsMiddleware(),
			handler.ListViewEntries)
		transactionsGroup.DELETE("/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.DeleteTransaction)
		transactionsGroup.POST("/income",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.CreateIncome)
		transactionsGroup.POST("/installment",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.CreateInstallment)
		transactionsGroup.PATCH("/simple-expense/:transaction_id",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.UpdateSimpleExpense)
	}
}
