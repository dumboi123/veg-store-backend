package route

import (
	"veg-store-backend/internal/infrastructure/router"

	"go.uber.org/fx"
)

type RoutesCollection []Routes

type Route[THandler any] struct {
	Handler THandler
	Router  *router.Router
}

type Routes interface {
	Setup()
}

func NewRoutesCollection(userRoutes *UserRoutes) RoutesCollection {
	return RoutesCollection{
		userRoutes,
	}
}

func (routesCollection RoutesCollection) Setup() {
	for _, routes := range routesCollection {
		routes.Setup()
	}
}

var RoutesModule = fx.Options(
	fx.Provide(NewUserRoutes),
	fx.Provide(NewRoutesCollection),
)
