package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/realtime-chat-backend/utils"
)

// region "SessionMiddleware" checks if a user session exists and sets the session expiration.
func SessionMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)        // Get the default session
		sessionUserID := session.Get("id")      // Retrieve user ID from session
		sessionUserMail := session.Get("email") // Retrieve user email from session

		// Check if session ID or email is missing
		if sessionUserID == nil || sessionUserMail == nil {
			ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized", "Authorization failed")) // Respond with an unauthorized status
			ctx.Abort()                                                                                       // Abort the request
			return
		}

		// Extend the session expiration time to 24 hours
		session.Set("Expires", time.Now().Add(24*time.Hour))

		// Create a new context with user ID and email
		socketCtx := context.WithValue(ctx.Request.Context(), "id", sessionUserID.(string))
		socketCtx = context.WithValue(socketCtx, "email", sessionUserMail.(string))

		// Update the request context
		ctx.Request = ctx.Request.WithContext(socketCtx)

		session.Save() // Save the session
		return
	}
}

// endregion

// region "CombinedAuthMiddleware" checks for a valid session or JWT token for authorization.
func CombinedAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := sessions.Default(ctx)        // Get the default session
		sessionMail := session.Get("email")     // Retrieve user email from session
		token := ctx.GetHeader("Authorization") // Get the Authorization header

		// If a session exists, extend its expiration
		if sessionMail != nil {
			session.Set("Expires", time.Now().Add(24*time.Hour)) // Extend expiration
			session.Save()                                       // Save the session
			ctx.Set("email", sessionMail.(string))               // Set email in context from session
			ctx.Next()                                           // Continue to the next handler
			return
		}

		// If no session exists, check for a JWT token
		if token != "" {
			err := utils.VerifyToken(token) // Verify the JWT token
			if err == nil {
				claims, _ := utils.GetClaims(token)
				ctx.Set("email", claims["user_email"].(string)) // Set email in context from JWT claims
				ctx.Next()                                      // Continue to the next handler
				return
			}
		}

		// If neither session nor token is valid, respond with unauthorized status
		ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized", "Authorization failed"))
		ctx.Abort() // Abort the request
	}
}

// endregion
