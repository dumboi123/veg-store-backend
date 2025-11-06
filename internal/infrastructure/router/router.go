package router

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"veg-store-backend/docs"
	"veg-store-backend/injection/core"
	"veg-store-backend/internal/restful/handler"
	"veg-store-backend/internal/restful/middleware"
	"veg-store-backend/util"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/fx"
)

type Router struct {
	ApiPath string
	Engine  *gin.Engine
}

func NewRouter() *Router {
	router := &Router{
		ApiPath: core.Configs.Server.ApiPrefix + core.Configs.Server.ApiVersion,
		Engine:  initGinEngine(),
	}

	if core.Configs.Mode != "prod" && core.Configs.Mode != "production" {
		// Register Swagger UI in non-production modes
		router.registerSwaggerUI()
	}

	return router
}

func (router *Router) HttpRun() error {
	// Run HTTP Server
	core.Logger.Info(fmt.Sprintf("Starting HTTP server on port %s...", core.Configs.Server.Port))
	httpServer := &http.Server{
		Addr:           ":" + core.Configs.Server.Port,
		Handler:        router.Engine,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server startup failed %w", err)
	}

	return nil
}

func initGinEngine() *gin.Engine {
	engine := gin.New()

	// Custom log for Gin per request
	core.UseGinRequestLogging(engine)

	// Register Custom recovery handler for Gin
	engine.Use(gin.CustomRecovery(func(ginContext *gin.Context, recovered interface{}) {
		handler.CustomRecoveryHandler(core.GetHttpContext(ginContext), recovered)
	}))

	err := engine.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic("Failed to set trusted proxies" + err.Error())
	}

	// Configure CORS
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     core.Configs.Cors.AllowOrigins,
		AllowMethods:     core.Configs.Cors.AllowMethods,
		AllowHeaders:     core.Configs.Cors.AllowHeaders,
		AllowCredentials: core.Configs.Cors.AllowCredentials,
	}))

	// Register all middlewares
	engine.Use(
		middleware.Locale(util.DefaultLocale),
		middleware.HttpContext(),
		middleware.TraceID(),
		middleware.ErrorHandler(),
	)

	return engine
}

func (router *Router) registerSwaggerUI() {
	docs.SwaggerInfo.Host = core.Configs.Swagger.Host
	docs.SwaggerInfo.BasePath = router.ApiPath

	swaggerUiPrefix := docs.SwaggerInfo.BasePath + "/swagger-ui/*any"
	router.Engine.GET(swaggerUiPrefix, ginSwagger.WrapHandler(swaggerFiles.Handler)) /*,
	ginSwagger.URL("http://localhost:8080"+router.ApiPath+"/swagger-ui/doc.json"),
	ginSwagger.DefaultModelsExpandDepth(1)*/
}

var RouterModule = fx.Options(fx.Provide(NewRouter))
