package categories

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
	group := router.Group("/api/v1/categories")
	{
		group.POST("",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.Create)
		group.GET("", middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryOptsMiddleware(),
			handler.List)
		group.DELETE("/:category_id",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.DeleteByID)
		group.GET("/:period",
			middlewares.RequireAuthMiddleware(jwtService),
			middlewares.QueryOptsMiddleware(),
			handler.ListCategoryAmountPerPeriod)
		group.PATCH("/:category_id",
			middlewares.RequireAuthMiddleware(jwtService),
			handler.Update)
	}
}
