package api

import (
	"cdn/common/constants"
	"cdn/config"
	"cdn/middlewares"
	"cdn/services/document"
	"cdn/services/image"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"

	"github.com/gin-gonic/gin"
)

type APIControllers struct {
	ImageController    *image.Controller
	DocumentController *document.Controller
}

var Controllers = &APIControllers{}

// Register all API endpoints
func registerEndpoints(humaApi *huma.API) {
	image.RegisterEndpoints(humaApi, Controllers.ImageController)
	document.RegisterEndpoints(humaApi, Controllers.DocumentController)
}

// Set up and start the API: set up API documentation,
// configure middlewares, and security measures.
func Start() {
	// Set up gin for your API
	gin.SetMode(config.Env.GinMode)
	gin.ForceConsoleColor()
	engine := gin.Default()
	engine.HandleMethodNotAllowed = true
	engine.ForwardedByClientIP = true
	err := engine.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		panic(err)
	}
	// Register gin middlewares
	engine.Use(middlewares.GinContextRegister())

	// Group API
	ginGroup := engine.Group(config.Env.ApiGroup)

	// OpenAPI documentation based on huma
	humaConfig := huma.DefaultConfig(constants.OPEN_API_TITLE, constants.OPEN_API_VERSION)
	// Custom CreateHooks to remove $schema links
	humaConfig.CreateHooks = []func(huma.Config) huma.Config{
		func(c huma.Config) huma.Config {
			return c
		},
	}
	humaConfig.DocsPath = ""
	humaConfig.Servers = []*huma.Server{
		{URL: config.Env.ApiGroup},
	}
	humaConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		constants.SECURITY_AUTH_NAME: {
			Type:        "apiKey",
			Name:        "X-API-Key",
			In:          "header",
			Scheme:      "token",
			Description: "Api key used to perform create, update and delete actions.",
		},
	}
	humaConfig.Info.Description = constants.OPEN_API_DESCRIPTION
	humaApi := humagin.NewWithGroup(engine, ginGroup, humaConfig)
	// Register huma middlewares
	humaApi.UseMiddleware(
		middlewares.HeadersMiddleware(humaApi),
		middlewares.CorsMiddleware(humaApi),
		middlewares.AuthMiddleware(humaApi),
	)

	// Serve static files as favicon
	engine.StaticFS("/assets", http.Dir(constants.ASSET_APP_PATH))
	engine.StaticFS("/openapi", http.Dir(constants.ASSET_OPENAPI_PATH))

	// Register endpoints
	ginGroup.GET("/docs", func(ctx *gin.Context) {
		ctx.Data(200, "text/html", config.OpenAPITemplates.Docs)
	})
	ginGroup.GET("/docs/scalar", func(ctx *gin.Context) {
		ctx.Data(200, "text/html", config.OpenAPITemplates.Scalar)
	})
	ginGroup.GET("/docs/swagger", func(ctx *gin.Context) {
		ctx.Data(200, "text/html", config.OpenAPITemplates.Swagger)
	})
	ginGroup.GET("/docs/redocly", func(ctx *gin.Context) {
		ctx.Data(200, "text/html", config.OpenAPITemplates.Redocly)
	})
	registerEndpoints(&humaApi)

	// Start to listen
	formattedPort := fmt.Sprintf(":%d", config.Env.AppPort)
	err = engine.Run(formattedPort)
	if err != nil {
		panic(err)
	}
}
