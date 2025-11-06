package middleware

import (
	"veg-store-backend/injection/core"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
)

/*
This middleware injects the application context into each Gin request context.
It provides access to shared resources like configuration, logger, and localizer.
*/

func HttpContext() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		httpContext := &core.HttpContext{
			Translator: core.Translator,
			Gin:        ginContext,
		}

		ginContext.Set(util.AppContextKey, httpContext)
		ginContext.Next()
	}
}
