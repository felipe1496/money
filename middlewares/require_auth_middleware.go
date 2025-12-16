package middlewares

import (
	"net/http"
	"rango-backend/services"
	"rango-backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuthMiddleware(JWTService services.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "missing token"))
			ctx.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := JWTService.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, "invalid token"))
			ctx.Abort()
			return
		}

		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
