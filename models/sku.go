package models

import (
	"context"
	"github.com/uptrace/bun"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type Sku struct {
	Id        uint      `bun:",pk,autoincrement" json:"id"`
	Sku       string    `bun:",notnull,unique" json:"sku"`
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
		Column("sku.id").
		Column("sku.sku").
		Column("sku.product_id").
		/*Relation("Product", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ColumnExpr("ref as product_ref")
		}).*/
		Column("sku.created_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}

func GetSku(ctx context.Context) (s *Sku) {
	db := database.DB
	err := db.NewSelect().
		Model(&s).
		Column("sku.id").
		Column("sku.sku").
		Column("sku.product_id").
		Relation("Product", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ColumnExpr("ref as product_ref")
		}).
		Column("sku.created_at").
		Where("sku_id = ?").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
