package models

import (
	"context"
	"fmt"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type Item struct {
	Id         uint      `bun:",pk,autoincrement" json:"id"`
	UIC        string    `bun:",notnull,unique,type:varchar(6)" json:"uic"`
	SkuID      uint      `json:"sku_id"`
	Sku        *Sku      `bun:"rel:belongs-to,join:sku_id=id"`
	LoteID     uint      `json:"lote_id"`
	Lote       *Lote     `bun:"rel:belongs-to,join:lote_id=id"`
	LocationID uint      `json:"location_id"`
	Location   *Location `bun:"rel:belongs-to,join:location_id=id"`
	StatusID   uint      `json:"status_id"`
	Status     *Status   `bun:"rel:belongs-to,join:status_id=id"`
	UserID     uint      `json:"user_id"`
	User       *Location `bun:"rel:belongs-to,join:user_id=id"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamp" json:"created_at"`
	UpdatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamp" json:"updated_at"`
	DeletedAt  time.Time `bun:",soft_delete,nullzero,type:timestamp" json:"deleted_at"`
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
	i.UpdatedAt = time.Now()
	db := database.Pgdb

	cmd, err := db.Prepare(`
		UPDATE items
		SET location_id = l.id,
    		status_id = l.status_id,
    		updated_at = $1
		FROM (SELECT id, status_id FROM locations WHERE id = $2) as l
		WHERE items.id = $3;
	`)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer cmd.Close()
	cmd.Exec(i.UpdatedAt, i.LocationID, i.Id)

	return
}
