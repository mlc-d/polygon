package models

import (
	"context"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type Item struct {
	Id         uint      `bun:",pk,autoincrement" json:"id"`
	UIC        string    `bun:",notnull,unique,type:varchar(6)" json:"uic"`
	SkuID      uint      `json:"sku_id"`
	Sku        *Sku      `bun:"rel:belongs-to,join:sku_id=id"`
	LocationID uint      `json:"location_id"`
	Location   *Location `bun:"rel:belongs-to,join:location_id=id"`
	StatusID   uint      `json:"status_id"`
	Status     *Status   `bun:"rel:belongs-to,join:status_id=id"`
	UserID     uint      `json:"user_id"`
	User       *Location `bun:"rel:belongs-to,join:user_id=id"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt  time.Time `bun:",soft_delete,nullzero,type:timestamptz" json:"deleted_at"`
}

func CreateItem(ctx context.Context, i *Item) (err error) {
	db := database.DB
	_, err = db.NewInsert().
		Model(i).
		Exec(ctx)
	if err == nil {
		_ = CreateHistory(ctx, &History{
			ItemID:     i.Id,
			SkuID:      i.SkuID,
			LocationID: i.LocationID,
			StatusID:   i.StatusID,
			UserID:     i.UserID,
		})
	}
	return
}

func GetItems(ctx context.Context) (items []Item) {
	db := database.DB
	err := db.NewSelect().
		Model(&items).
		//ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}

func UpdateItem(ctx context.Context, i *Item) (err error) {
	i.UpdatedAt = time.Now()
	db := database.DB
	_, err = db.NewUpdate().
		Model(i).
		ExcludeColumn("id").
		ExcludeColumn("uic").
		ExcludeColumn("created_at").
		ExcludeColumn("deleted_at").
		Where("id = ?", i.Id).
		Returning("NULL").
		Exec(ctx)
	if err == nil {
		_ = CreateHistory(ctx, &History{
			ItemID:     i.Id,
			SkuID:      i.SkuID,
			LocationID: i.LocationID,
			StatusID:   i.StatusID,
			UserID:     i.UserID,
		})
	}
	return
}
