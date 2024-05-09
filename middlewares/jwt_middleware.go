package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/helpers"
)

func JwtMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "Unauthorized", "Authorization token is required"))
			ctx.Abort()
			return
		}

		err := helpers.VerifyToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "Unauthorized", "Authorization failed"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}