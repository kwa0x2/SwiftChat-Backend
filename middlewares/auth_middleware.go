package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized", "Authorization token is required"))
			ctx.Abort()
			return
		}

		err := utils.VerifyToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized", "Authorization failed"))
			ctx.Abort()
			return
		}
	}
}

func SessionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		sessionUserID := session.Get("id")
		sessionUserMail := session.Get("mail")

		if sessionUserID == nil || sessionUserMail == nil {
			ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized", "Authorization failed"))
			ctx.Abort()
			return
		}

		session.Set("Expires", time.Now().Add(24*time.Hour))

		socketCtx := context.WithValue(ctx.Request.Context(), "id", sessionUserID.(string))
		socketCtx = context.WithValue(socketCtx, "mail", sessionUserMail.(string))

		ctx.Request = ctx.Request.WithContext(socketCtx)
		session.Save()
		return
	}
}

func CombinedAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)
		sessionMail := session.Get("mail")
		token := ctx.GetHeader("Authorization")

		// ğer session varsa
		if sessionMail != nil {
			session.Set("Expires", time.Now().Add(24*time.Hour))
			session.Save()
			ctx.Set("mail", sessionMail.(string)) // sessiondaki maili contexte atar
			ctx.Next()
			return
		}

		// session yoksa jwt kontrolu
		if token != "" {
			err := utils.VerifyToken(token)
			if err == nil {
				claims, _ := utils.GetClaims(token)
				ctx.Set("mail", claims["user_email"].(string)) // jwtden gelen maili contexte atmak icin
				ctx.Next()
				return
			}
		}

		// ikiside yoksa
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized", "Authorization failed"))
		ctx.Abort()
	}
}
