package models

import (
	"context"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type Item struct {
	Id         uint      `bun:",pk,autoincrement" json:"id"`
	UIC        string    `bun:",notnull,unique,type:varchar(6)" json:"uic"`
	SkuID      uint      `json:"sku"`
	LocationID uint      `json:"location"`
	StatusID   uint      `json:"status"`
	UserID     uint      `json:"user"`
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

func AllocateItem(ctx context.Context, i *Item) (err error) {
	updatedAt := time.Now()
	db := database.DB
	_, err = db.NewUpdate().
		Model(&Item{}).
		Column("location_id", "updated_at").
		Set("location_id = ?", i.LocationID).
		Set("updated_at = ?", updatedAt).
		Where("id = ?", i.Id).
		Exec(ctx)
	if err == nil {
		_ = CreateHistory(ctx, &History{
			ItemID:     i.Id,
			LocationID: i.LocationID,
		})
	}
	return
}
