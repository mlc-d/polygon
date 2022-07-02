package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	conf "gitlab.com/mlcprojects/wms/config"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/migrations"
	"gitlab.com/mlcprojects/wms/routes"
	"net/http"
	"os"
)

var (
	config = conf.Cf
)

func main() {
	ctx := context.Background()
	database.InitDB()
	migrations.MigrateModels(ctx)
	migrations.InsertDefaults(ctx)

	e := echo.New()

	origins := config.Origins.Url

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{origins},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	routes.SetUpRoutes(e)

	port := os.Getenv("PORT")

	if port == "" {
		port = "1998"
	}

	e.Logger.Fatal(e.Start(":1998"))
}
