package models

import (
	"context"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type Sku struct {
	Id        uint      `bun:",pk,autoincrement" json:"id"`
	Sku       string    `bun:",notnull,unique" json:"name"`
	ProductID uint      `json:"product_id"`
	Product   *Product  `bun:"rel:belongs-to,join:product_id=id"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"deleted_at"`
}

func CreateSku(ctx context.Context, s *Sku) (err error) {
	db := database.DB
	_, err = db.NewInsert().
		Model(s).
		Exec(ctx)
	return
}

func GetSkus(ctx context.Context) (skus []Sku) {
	db := database.DB
	err := db.NewSelect().
		Model(&skus).
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
