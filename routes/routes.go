package routes

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	conf "gitlab.com/mlcprojects/wms/config"
	"gitlab.com/mlcprojects/wms/controllers"
)

var (
	cf       = conf.Cf
	endpoint = cf.Api.Endpoint
	version  = cf.Api.Version
)

func SetUpRoutes(e *echo.Echo) {

	config := middleware.JWTConfig{
		TokenLookup:    "header:accessToken",
		ParseTokenFunc: controllers.ValidateAccessToken,
	}
	fmt.Println(endpoint + version)
	apiGroup := e.Group(endpoint + version)
	apiGroup.Use(middleware.JWTWithConfig(config))

	// histories
	apiGroup.GET("/histories", controllers.GetHistories)
	apiGroup.POST("/histories", controllers.CreateHistory)
	// items
	apiGroup.GET("/items", controllers.GetItems)
	apiGroup.POST("/items", controllers.CreateItem)
	apiGroup.PATCH("/items", controllers.UpdateItem)
	// skus
	apiGroup.GET("/skus", controllers.GetSkus)
	apiGroup.POST("/skus", controllers.CreateSku)
	// users
	apiGroup.GET("/users", controllers.GetUsers)
	apiGroup.GET("/users/:id", controllers.GetUser)
	apiGroup.POST("/users", controllers.CreateUser)
	// products
	apiGroup.GET("/products", controllers.GetProducts)
	apiGroup.POST("/products", controllers.CreateProduct)
	// roles
	apiGroup.GET("/roles", controllers.GetRoles)
	apiGroup.POST("/roles", controllers.CreateRole)
	// statuses
	apiGroup.GET("/statuses", controllers.GetStatuses)
	apiGroup.POST("/statuses", controllers.CreateStatus)
	// locations
	apiGroup.GET("/locations", controllers.GetLocations)
	apiGroup.POST("/locations", controllers.CreateLocation)
	// lotes
	apiGroup.GET("/lotes", controllers.GetLotes)
	apiGroup.POST("/lotes", controllers.CreateLote)

	// test
	// apiGroup.POST("/testing", controllers.SampleHandler)

	//login
	e.POST("/login", controllers.Login)
	//refresh
	e.POST("/refresh", controllers.Refresh)
}
