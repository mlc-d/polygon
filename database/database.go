package database

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	conf "gitlab.com/mlcprojects/wms/config"
)

var (
	DB     *bun.DB
	Ctx    = context.Background()
	config = conf.Cf
)

func InitDB() {
	fmt.Println(config.Db.Dsn)
	sqldb, err := sql.Open("mysql", config.Db.Dsn)
	if err != nil {
		panic(err)
	}

	DB = bun.NewDB(sqldb, mysqldialect.New())
}
