package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/migrations"
	"gitlab.com/mlcprojects/wms/routes"
	"net/http"
)

func main() {
	ctx := context.Background()
	database.InitDB()
	migrations.MigrateModels(ctx)
	migrations.InsertDefaults(ctx)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"localhost"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	routes.SetUpRoutes(e)
	e.Logger.Fatal(e.Start(":1998"))
}
