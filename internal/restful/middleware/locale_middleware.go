package middleware

import (
	"fmt"
	"strings"
	"veg-store-backend/util"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Locale auto read "Accept-Language" header
func Locale(defaultLocale string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		lang := ginCtx.GetHeader("Accept-Language")

		zap.L().Info(fmt.Sprintf("Accept-Language: %s", lang))

		if lang == "" {
			lang = defaultLocale

		} else {
			// Just get the first language tag, e.g. "en-US,en;q=0.9" -> "en"
			lang = strings.Split(lang, ",")[0]
			lang = strings.Split(lang, "-")[0]
			lang = strings.TrimSpace(lang)
		}

		// Save locale to context
		ginCtx.Set(util.LocaleContextKey, lang)
		ginCtx.Next()
	}
}
