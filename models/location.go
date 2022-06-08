package models

import (
	"context"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type Location struct {
	Id        uint      `bun:",pk,autoincrement" json:"id"`
	Location  string    `bun:",notnull,unique" json:"location"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamp" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"deleted_at"`
}

func CreateLocation(ctx context.Context, l *Location) (err error) {
	db := database.DB
	_, err = db.NewInsert().
		Model(l).
		Exec(ctx)
	return
}

func GetLocations(ctx context.Context) (locations []Location) {
	db := database.DB
	err := db.NewSelect().
		Model(&locations).
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
