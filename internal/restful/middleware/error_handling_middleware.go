package middleware

import (
	"net/http"
	"strings"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/dto"
	"veg-store-backend/internal/application/exception"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ginContext.Next()

		if len(ginContext.Errors) == 0 {
			return
		}

		err := ginContext.Errors.Last().Err

		var subError exception.SubError
		switch rootError := err.(type) {
		case exception.SubError:
			subError = rootError

		default:
			core.Logger.Error("Unhandled error", zap.Error(err))
			ginContext.JSON(http.StatusInternalServerError, dto.HttpResponse[any]{
				HttpStatus: http.StatusInternalServerError,
				Code:       "internal/server-error",
				Message:    "Internal Server Error",
			})
			return
		}

		// Map code prefix -> HTTP httpStatus
		httpStatus := mapErrorCodeToStatus(subError.Code)

		// Get trace_id
		traceID := util.GetTraceId(ginContext)

		core.Logger.Error("Request failed",
			zap.String("trace_id", traceID),
			zap.String("code", subError.Code),
			zap.String("message", subError.MessageKey),
			zap.String("path", ginContext.FullPath()),
			zap.String("method", ginContext.Request.Method),
		)

		ginContext.JSON(httpStatus, dto.HttpResponse[any]{
			HttpStatus: httpStatus,
			Code:       subError.Code,
			Message:    core.Translator.T(util.GetLocale(ginContext), subError.MessageKey),
		})
	}
}

func mapErrorCodeToStatus(code string) int {
	switch {
	case strings.HasPrefix(code, "invalid/"):
		return http.StatusBadRequest
	case strings.HasPrefix(code, "auth/unauthenticated"):
		return http.StatusUnauthorized
	case strings.HasPrefix(code, "auth/forbidden"):
		return http.StatusForbidden
	case strings.HasPrefix(code, "not_found/"):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
