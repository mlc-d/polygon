package migrations

import (
	"context"
	"flag"
	"fmt"
	"gitlab.com/mlcprojects/wms/database"
	"gitlab.com/mlcprojects/wms/models"
	"log"
)

func MigrateModels(ctx context.Context) {
	db := database.DB
	// create 'roles'
	if _, err := db.NewCreateTable().
		Model(&models.Role{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	// create 'statuses'
	if _, err := db.NewCreateTable().
		Model(&models.Status{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	// create 'lotes'
	if _, err := db.NewCreateTable().
		Model(&models.Lote{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	// create 'products'
	if _, err := db.NewCreateTable().
		Model(&models.Product{}).
		IfNotExists().
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	// create 'skus'
	if _, err := db.NewCreateTable().
		Model(&models.Sku{}).
		IfNotExists().
		WithForeignKeys().
		ForeignKey("(product_id) REFERENCES products (id) ON DELETE RESTRICT").
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	// 'locations'
	if _, err := db.NewCreateTable().
		Model(&models.Location{}).
		IfNotExists().
		WithForeignKeys().
		ForeignKey("(status_id) REFERENCES statuses (id) ON DELETE RESTRICT").
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	// 'users'
	if _, err := db.NewCreateTable().
		Model(&models.User{}).
		IfNotExists().
		WithForeignKeys().
		ForeignKey("(role_id) REFERENCES roles (id) ON DELETE RESTRICT").
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	// 'items'
	if _, err := db.NewCreateTable().
		Model(&models.Item{}).
		IfNotExists().
		WithForeignKeys().
		ForeignKey("(sku_id) REFERENCES skus (id) ON DELETE RESTRICT").
		ForeignKey("(lote_id) REFERENCES lotes (id) ON DELETE RESTRICT").
		ForeignKey("(location_id) REFERENCES locations (id) ON DELETE RESTRICT").
		ForeignKey("(status_id) REFERENCES statuses (id) ON DELETE RESTRICT").
		ForeignKey("(user_id) REFERENCES users (id) ON DELETE RESTRICT").
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	// 'histories'
	if _, err := db.NewCreateTable().
		Model(&models.History{}).
		IfNotExists().
		WithForeignKeys().
		ForeignKey("(item_id) REFERENCES items (id) ON DELETE RESTRICT").
		ForeignKey("(sku_id) REFERENCES skus (id) ON DELETE RESTRICT").
		ForeignKey("(location_id) REFERENCES locations (id) ON DELETE RESTRICT").
		ForeignKey("(status_id) REFERENCES statuses (id) ON DELETE RESTRICT").
		ForeignKey("(user_id) REFERENCES users (id) ON DELETE RESTRICT").
		Exec(ctx); err != nil {
		panic(err.Error())
	}
	log.Print("tables successfully migrated")
}

func InsertDefaults(ctx context.Context) {
	f := flag.Bool("d", false, "insert default values")
	flag.Parse()
	if *f {
		db := database.DB
		var rolesList = []models.Role{
			{Role: "dev", IsAdmin: true},
			{Role: "admin", IsAdmin: true},
			{Role: "manager", IsAdmin: true},
			{Role: "supervisor", IsAdmin: true},
			{Role: "leader", IsAdmin: true},
			{Role: "publisher", IsAdmin: false},
			{Role: "operator", IsAdmin: false},
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
		fmt.Println("Please enter password for 'dev' user:")
		var pass string
		if _, err := fmt.Scanln(&pass); err != nil {
			panic(err.Error())
		}
		if err := models.CreateUser(ctx, &models.User{
			Name:     "dev",
			Password: pass,
			RoleID:   1,
		}); err != nil {
			panic(err.Error())
		}

	}
	log.Print("default values created")
}
