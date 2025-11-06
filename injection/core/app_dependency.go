package core

import (
	"context"

	"go.uber.org/fx"
)

type AppDependency interface {
	Name() string
	Start() error
	Stop() error
}

// LifecycleModule helper: tạo fx.Module cho bất cứ AppDependency nào
func LifecycleModule(constructor interface{}) fx.Option {
	return fx.Options(
		fx.Provide(constructor),
		fx.Invoke(func(lifecycle fx.Lifecycle, d AppDependency) {
			lifecycle.Append(fx.Hook{
				OnStart: func(context context.Context) error {
					return d.Start()
				},
				OnStop: func(context context.Context) error {
					return d.Stop()
				},
			})
		}),
	)
}
