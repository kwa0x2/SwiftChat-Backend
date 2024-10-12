package utils

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func HandleErrorWithSentry(ctx *gin.Context, err error, additionalData map[string]interface{}) {
	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			for key, value := range additionalData {
				scope.SetExtra(key, value)
			}
			hub.CaptureException(err)
		})
	}
}
