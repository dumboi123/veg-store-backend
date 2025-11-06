package injection_test

import (
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/restful/handler"

	"github.com/gin-gonic/gin"
)

func MockUserRoutes(handler *handler.UserHandler) *gin.Engine {
	mockRouter := router.NewRouter()
	api := mockRouter.Engine.Group(mockRouter.ApiPath + "/user")
	{
		api.GET("/hello", func(ginCtx *gin.Context) {
			handler.Hello(MockHttpContext(ginCtx))
		})
		api.GET("/details/:id", func(ginCtx *gin.Context) {
			handler.Details(MockHttpContext(ginCtx))
		})
		api.GET("/ping", func(ginCtx *gin.Context) {
			handler.HealthCheck(MockHttpContext(ginCtx))
		})
		api.GET("/", func(ginCtx *gin.Context) {
			handler.GetAllUsers(MockHttpContext(ginCtx))
		})
	}
	return mockRouter.Engine
}
