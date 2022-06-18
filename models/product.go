package models

import (
	"context"
	"time"

	"gitlab.com/mlcprojects/wms/database"
)

type Product struct {
	Id          uint      `bun:",pk,autoincrement" json:"id,omitempty"`
	Name        string    `bun:",notnull,unique" json:"name,omitempty"`
	Ref         string    `bun:",notnull,unique" json:"ref,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at,omitempty"`
	UpdatedAt   time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at,omitempty"`
	DeletedAt   time.Time `bun:",soft_delete,nullzero" json:"deleted_at,omitempty"`
}

func CreateProduct(ctx context.Context, p *Product) (err error) {
	db := database.DB
	_, err = db.NewInsert().
		Model(p).
		Exec(ctx)
	return
}

func GetProducts(ctx context.Context) (products []Product) {
	db := database.DB
	err := db.NewSelect().
		Model(&products).
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
