package models

import (
	"context"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type Lote struct {
	Id        uint      `bun:",pk,autoincrement" json:"id"`
	Lote      string    `bun:",notnull,unique" json:"lote"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"deleted_at"`
}

func CreateLote(ctx context.Context, l *Lote) (err error) {
	db := database.DB
	_, err = db.NewInsert().
		Model(l).
		Exec(ctx)
	return
}

func GetLotes(ctx context.Context) (lotes []Lote) {
	db := database.DB
	err := db.NewSelect().
		Model(&lotes).
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
