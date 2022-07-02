package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	conf "gitlab.com/mlcprojects/wms/config"
	"gitlab.com/mlcprojects/wms/controllers"
)

const version = "/api/v1"

var cf = conf.Cf

func SetUpRoutes(e *echo.Echo) {

	config := middleware.JWTConfig{
		TokenLookup:    "header:accessToken",
		ParseTokenFunc: controllers.ValidateAccessToken,
	}

	apiGroup := e.Group(version)
	apiGroup.Use(middleware.JWTWithConfig(config))

	// entry point
	apiGroup.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})
	// histories
	apiGroup.GET("/histories", controllers.GetHistories)
	apiGroup.POST("/histories", controllers.CreateHistory)
	apiGroup.GET("/history", controllers.GetItemHistory)
	// items
	apiGroup.GET("/items", controllers.GetItems)
	apiGroup.GET("/item", controllers.GetItem)
	apiGroup.POST("/items", controllers.CreateItem)
	apiGroup.PATCH("/items", controllers.UpdateItem)
	apiGroup.PATCH("/allocate", controllers.AllocateItem)
	// skus
	apiGroup.GET("/skus", controllers.GetSkus)
	apiGroup.POST("/skus", controllers.CreateSku)
	// users
	apiGroup.GET("/users", controllers.GetUsers)
	apiGroup.GET("/users/:id", controllers.GetUser)
	apiGroup.POST("/users", controllers.CreateUser)
	apiGroup.PATCH("/users", controllers.EditUser)
	apiGroup.DELETE("/users", controllers.DeleteUser)
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
	apiGroup.PATCH("/locations", controllers.EditLocation)
	apiGroup.GET("/locbystatus", controllers.GetLocationByStatus)
	apiGroup.GET("/location", controllers.GetLocation)
	// batches
	apiGroup.GET("/batches", controllers.GetBatches)
	apiGroup.POST("/batches", controllers.CreateBatch)
	// login
	e.POST("/login", controllers.Login)
	// refresh
	e.POST("/refresh", controllers.Refresh)
	// logout
	e.POST("/logout", controllers.Logout)
}
