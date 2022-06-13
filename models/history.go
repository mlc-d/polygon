package models

import (
	"context"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type History struct {
	Id         uint      `bun:",pk,autoincrement" json:"id"`
	ItemID     uint      `json:"item_id"`
	Item       *Item     `bun:"rel:belongs-to,join:item_id=id" json:"-"`
	SkuID      uint      `json:"sku_id"`
	Sku        *Sku      `bun:"rel:belongs-to,join:sku_id=id" json:"-"`
	LocationID uint      `json:"location_id"`
	Location   *Location `bun:"rel:belongs-to,join:location_id=id" json:"-"`
	StatusID   uint      `json:"status_id"`
	Status     *Status   `bun:"rel:belongs-to,join:status_id=id" json:"-"`
	UserID     uint      `json:"user_id"`
	User       *User     `bun:"rel:belongs-to,join:user_id=id" json:"-"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
}

type PublicHistory struct {
	Id         uint      `json:"id"`
	ItemID     uint      `json:"item_id"`
	Uic        string    `json:"uic"`
	SkuID      uint      `json:"sku_id"`
	Sku        string    `json:"sku"`
	BatchID    uint      `json:"batch_id"`
	Batch      string    `json:"batch"`
	LocationID uint      `json:"location_id"`
	Location   string    `json:"location"`
	StatusID   uint      `json:"status_id"`
	Status     string    `json:"status"`
	UserID     uint      `json:"user_id"`
	User       string    `json:"user"`
	CreatedAt  time.Time `json:"created_at"`
}

func CreateHistory(ctx context.Context, h *History) (err error) {
	h.CreatedAt = time.Now()
	db := database.DB
	_, err = db.NewInsert().
		Model(h).
		Exec(ctx)
	return
}

func GetItemHistory(ctx context.Context, uic string) (err error, history []PublicHistory) {
	db := database.DB
	var itemId uint

	err = db.NewSelect().
		Model(&Item{}).
		Column("id").
		Where("uic = ?", uic).
		Limit(1).
		Scan(ctx, &itemId)

	if err != nil {
		return err, nil
	}

	err = db.NewSelect().
		Model(&History{}).
		ModelTableExpr("histories").
		Column("histories.id", "histories.item_id", "histories.sku_id", "histories.location_id", "histories.status_id", "histories.user_id", "histories.created_at").
		ColumnExpr("sku.sku as sku").
		ColumnExpr("i.uic as uic").
		ColumnExpr("l.location as location").
		ColumnExpr("s.status as status").
		ColumnExpr("u.name as user").
		Join("left join skus as sku").JoinOn("sku.id = histories.sku_id").
		Join("left join items as i").JoinOn("i.id = histories.item_id").
		Join("left join locations as l").JoinOn("l.id = histories.location_id").
		Join("left join statuses as s").JoinOn("s.id = histories.status_id").
		Join("left join users as u").JoinOn("u.id = histories.user_id").
		Where("item_id = ?", itemId).
		Order("histories.created_at DESC").
		Scan(ctx, &history)

	if err != nil {
		return err, nil
	}
	return nil, history
}

func GetHistories(ctx context.Context) (histories []PublicHistory) {
	db := database.DB
	err := db.NewSelect().
		Model(&histories).
		ModelTableExpr("histories").
		Column("histories.id", "histories.uic", "histories.sku_id", "histories.batch_id", "histories.location_id", "histories.status_id", "histories.user_id", "histories.updated_at").
		ColumnExpr("sku.sku as sku").
		ColumnExpr("b.batch as batch").
		ColumnExpr("l.location as location").
		ColumnExpr("s.status as status").
		ColumnExpr("u.name as user").
		Join("items as i").JoinOn("i.id = histories.item_id").
		Join("skus as sku").JoinOn("sku.id = items.sku_id").
		Join("locations as l").JoinOn("l.id = items.location_id").
		Join("statuses as s").JoinOn("s.id = items.status_id").
		Join("users as u").JoinOn("u.id = items.user_id").
		Limit(500).
		OrderExpr("histories.id DESC").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
