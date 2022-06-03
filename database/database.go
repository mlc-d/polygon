package database

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	conf "gitlab.com/mlcprojects/wms/config"
	"log"
)

var (
	DB     *bun.DB
	Ctx    = context.Background()
	config = conf.Cf
)

func InitDB() {
	dsn := config.Db.Dsn
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	DB = bun.NewDB(pgdb, pgdialect.New())
	log.Print("connected to database")
}
