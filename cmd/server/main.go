package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/application/exception"
	"veg-store-backend/internal/application/service"
	"veg-store-backend/internal/infrastructure/repository"
	"veg-store-backend/internal/infrastructure/router"
	"veg-store-backend/internal/restful/handler"
	"veg-store-backend/internal/restful/route"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// @title RESTFUL API VERSION 1.0
// @version 1.0
// @description This is an API document for veg-store-backend.
// @termsOfService http://example.com/terms/

// @contact.name Nhan Le
// @contact.url http://example.com/support
// @contact.email benlun1201@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.

// @schemes http https

func main() {
	injectGlobalComponents(determineMode())

	app := fx.New(
		repository.UserRepositoryModule,
		service.UserServiceModule,
		handler.UserHandlerModule,
		router.RouterModule,
		route.RoutesModule,

		fx.Invoke(func(lifecycle fx.Lifecycle, appRouter *router.Router, routesCollection route.RoutesCollection) {
			routesCollection.Setup()

			lifecycle.Append(fx.Hook{
				OnStart: func(context context.Context) error {
					go func() {
						if err := appRouter.HttpRun(); err != nil {
							core.Logger.Fatal("Server failed to start: ", zap.Error(err))
						}
					}()
					return nil
				},
				OnStop: func(context context.Context) error {
					core.Logger.Info("Shutting down server...")
					return nil
				},
			})
		}),
	)

	// Graceful shutdown
	startContext, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()
	if err := app.Start(startContext); err != nil {
		log.Fatal(err)
	}

	// Wait for OS signal to terminate
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	stopContext, cancel := context.WithTimeout(context.Background(), fx.DefaultTimeout)
	defer cancel()
	if err := app.Stop(stopContext); err != nil {
		log.Fatal(err)
	}
}

func injectGlobalComponents(mode string) {
	core.Configs.Mode = mode
	core.Logger = core.InitLogger()       // Initialize Logger
	core.Configs = core.Load()            // Load configuration
	core.Translator = core.InitI18n()     // Initialize i18n Translator
	core.Error = exception.InitAppError() // Initialize App Error
}

func determineMode() string {
	mode := os.Getenv("MODE")
	switch mode {
	case "production", "prod":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	if mode == "" {
		zap.NewExample().Warn("No 'MODE' is defined. Server will run in 'dev' mode by default.")
		return "dev"
	}
	return mode
}
