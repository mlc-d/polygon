package models

import (
	"context"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type Batch struct {
	Id        uint      `bun:",pk,autoincrement" json:"id"`
	Batch     string    `bun:",notnull,unique" json:"batch"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"deleted_at"`
}

func CreateBatch(ctx context.Context, l *Batch) (err error) {
	db := database.DB
	_, err = db.NewInsert().
		Model(l).
		Exec(ctx)
	return
}

func GetBatches(ctx context.Context) (batches []Batch) {
	db := database.DB
	err := db.NewSelect().
		Model(&batches).
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
