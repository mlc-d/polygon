package models

import (
	"context"
	"fmt"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type History struct {
	Id         uint      `bun:",pk,autoincrement" json:"id"`
	ItemID     uint      `json:"item"`
	SkuID      uint      `json:"sku"`
	LocationID uint      `json:"location"`
	StatusID   uint      `json:"status"`
	UserID     uint      `json:"user"`
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
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
