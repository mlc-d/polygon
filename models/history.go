package models

import (
	"context"
	"fmt"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type History struct {
	Id         uint      `bun:",pk,autoincrement" json:"id"`
	ItemID     uint      `json:"item_id"`
	Item       *Item     `bun:"rel:belongs-to,join:item_id=id"`
	SkuID      uint      `json:"sku_id"`
	Sku        *Sku      `bun:"rel:belongs-to,join:sku_id=id"`
	LocationID uint      `json:"location_id"`
	Location   *Location `bun:"rel:belongs-to,join:location_id=id"`
	StatusID   uint      `json:"status_id"`
	Status     *Status   `bun:"rel:belongs-to,join:status_id=id"`
	UserID     uint      `json:"user_id"`
	User       *Location `bun:"rel:belongs-to,join:user_id=id"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
}

func CreateHistory(ctx context.Context, h *History) (err error) {
	fmt.Println("llegué a histories")
	fmt.Println(h)
	db := database.DB
	_, err = db.NewInsert().
		Model(h).
		Exec(ctx)
	return
}

func GetHistories(ctx context.Context) (histories []History) {
	db := database.DB
	err := db.NewSelect().
		Model(&histories).
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
