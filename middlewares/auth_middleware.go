package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"

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
	}
}

func SessionMiddleware() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		session:=sessions.Default(ctx)
		sessionUserID:=session.Get("id")
		if sessionUserID == nil {
			ctx.JSON(http.StatusUnauthorized, helpers.NewErrorResponse(http.StatusUnauthorized, "Unauthorized", "Authorization failed"))
			ctx.Abort()
			return
		}
		session.Set("Expires", time.Now().Add(24*time.Hour))
		session.Save()
		return
	}
}
