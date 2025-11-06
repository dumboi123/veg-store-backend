package route

import (
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/restful/handler"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	*Route[*handler.UserHandler]
}

func NewUserRoutes(userHandler *handler.UserHandler, router *router.Router) *UserRoutes {
	return &UserRoutes{
		Route: &Route[*handler.UserHandler]{
			Handler: userHandler,
			Router:  router,
		},
	}
}

func (routes *UserRoutes) Setup() {
	api := routes.Router.Engine.Group(routes.Router.ApiPath + "/user")
	{
		api.GET("/hello", func(ginContext *gin.Context) {
			routes.Handler.Hello(core.GetHttpContext(ginContext))
		})
		api.GET("/details/:id", func(ginContext *gin.Context) {
			routes.Handler.Details(core.GetHttpContext(ginContext))
		})
		api.GET("/ping", func(ginContext *gin.Context) {
			routes.Handler.HealthCheck(core.GetHttpContext(ginContext))
		})
		api.GET("/", func(ginContext *gin.Context) {
			routes.Handler.GetAllUsers(core.GetHttpContext(ginContext))
		})
	}
}
