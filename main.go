package main

import (
	"context"
	"flag"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"gitlab.com/mlcprojects/wms/routes"
	"log"
	"net/http"
)

func MigrateModels(ctx context.Context) {
	db := database.DB
	if _, err := db.NewCreateTable().
		Model(&models.History{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	if _, err := db.NewCreateTable().
		Model(&models.Item{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	if _, err := db.NewCreateTable().
		Model(&models.Location{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	if _, err := db.NewCreateTable().
		Model(&models.Lote{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	if _, err := db.NewCreateTable().
		Model(&models.Product{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	if _, err := db.NewCreateTable().
		Model(&models.Role{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	if _, err := db.NewCreateTable().
		Model(&models.Sku{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	if _, err := db.NewCreateTable().
		Model(&models.Status{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	if _, err := db.NewCreateTable().
		Model(&models.User{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	log.Print("tables successfully migrated")
}

func insertDefaults(ctx context.Context) {
	f := flag.Bool("d", false, "insert default values")
	flag.Parse()
	if *f {
		db := database.DB
		var rolesList = []models.Role{
			{Role: "admin", IsAdmin: true},
			{Role: "manager", IsAdmin: true},
			{Role: "supervisor", IsAdmin: true},
			{Role: "leader", IsAdmin: true},
			{Role: "operator", IsAdmin: false},
			{Role: "publisher", IsAdmin: false},
		}
		if _, err := db.NewInsert().Model(&rolesList).Exec(ctx); err != nil {
			return
		}
		var statusList = []models.Status{
			{Status: "created"},
			{Status: "held"},
			{Status: "ready"},
			{Status: "taken"},
			{Status: "processing"},
			{Status: "gone"},
			{Status: "lost"},
			{Status: "deleted"},
		}
		if _, err := db.NewInsert().Model(&statusList).Exec(ctx); err != nil {
			return
		}
	}
	log.Print("default values created")
}

func main() {
	ctx := context.Background()
	database.InitDB()
	MigrateModels(ctx)
	insertDefaults(ctx)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	routes.SetUpRoutes(e)
	e.Logger.Fatal(e.Start(":1998"))
}
