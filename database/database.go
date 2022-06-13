package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	conf "gitlab.com/mlcprojects/wms/config"
	"log"
)

var (
	DB     *bun.DB
	Pgdb   *sql.DB
	Ctx    = context.Background()
	config = conf.Cf
)

func InitDB() {
	fmt.Println(config.Db.Dsn)
	dsn := config.Db.Dsn
	Pgdb = sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	DB = bun.NewDB(Pgdb, pgdialect.New())
	DB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	log.Print("connected to database")
}
