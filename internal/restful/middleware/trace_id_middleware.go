package middleware

import (
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceID() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		traceID := ginCtx.GetHeader("X-Request-ID")
		if traceID == "" {
			traceID = uuid.NewString()
		}
		ginCtx.Set(util.TraceIDContextKey, traceID)
		ginCtx.Writer.Header().Set("X-Request-ID", traceID)
		ginCtx.Next()
	}
}
