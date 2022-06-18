package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/mlcprojects/wms/database"
)

type Item struct {
	Id         uint      `bun:",pk,autoincrement" json:"id"`
	UIC        string    `bun:",notnull,unique,type:varchar(6)" json:"uic"`
	SkuID      uint      `json:"sku_id"`
	Sku        *Sku      `bun:"rel:belongs-to,join:sku_id=id"`
	BatchID    uint      `json:"batch_id"`
	Batch      *Batch    `bun:"rel:belongs-to,join:batch_id=id"`
	LocationID uint      `json:"location_id"`
	Location   *Location `bun:"rel:belongs-to,join:location_id=id"`
	StatusID   uint      `json:"status_id"`
	Status     *Status   `bun:"rel:belongs-to,join:status_id=id"`
	UserID     uint      `json:"user_id"`
	User       *User     `bun:"rel:belongs-to,join:user_id=id"`
	CreatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz"`
	UpdatedAt  time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt  time.Time `bun:",soft_delete,nullzero,type:timestamptz"`
}

type PublicItem struct {
	Id         uint      `json:"id"`
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
	UpdatedAt  time.Time `json:"updated_at"`
}

func (i *Item) CreateItem(ctx context.Context) (err error) {
	//just in case
	i.StatusID = 1

	db := database.DB
	_, err = db.NewInsert().
		Model(i).
		Exec(ctx)

	if err != nil {
		err = errors.New(i.UIC)
		return err
	}

	var itemId uint

	err = db.NewSelect().
		Model(&Item{}).
		Column("id").
		Where("uic = ?", i.UIC).
		Limit(1).
		Scan(ctx, &itemId)

	if err == nil {
		fmt.Println("crear historia")
		err = CreateHistory(ctx, &History{
			ItemID:     itemId,
			SkuID:      i.SkuID,
			LocationID: i.LocationID,
			StatusID:   i.StatusID,
			UserID:     i.UserID,
		})
	}
	return
}

func GetItems(ctx context.Context) (items []PublicItem) {
	db := database.DB
	err := db.NewSelect().
		Model(&items).
		ModelTableExpr("items").
		Column("items.id", "items.uic", "items.sku_id", "items.batch_id", "items.location_id", "items.status_id", "items.user_id", "items.updated_at").
		ColumnExpr("sku.sku as sku").
		ColumnExpr("b.batch as batch").
		ColumnExpr("l.location as location").
		ColumnExpr("s.status as status").
		ColumnExpr("u.name as user").
		Join("left join skus as sku").JoinOn("sku.id = items.sku_id").
		Join("left join batches as b").JoinOn("b.id = items.batch_id").
		Join("left join locations as l").JoinOn("l.id = items.location_id").
		Join("left join statuses as s").JoinOn("s.id = items.status_id").
		Join("left join users as u").JoinOn("u.id = items.user_id").
		Limit(500).
		OrderExpr("items.id DESC").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}

func GetItem(ctx context.Context, uic string) (i PublicItem, err error) {
	db := database.DB
	err = db.NewSelect().
		Model(&i).
		ModelTableExpr("items").
		Column("items.id", "items.uic", "items.sku_id", "items.batch_id", "items.location_id", "items.status_id", "items.user_id", "items.updated_at").
		ColumnExpr("sku.sku as sku").
		ColumnExpr("b.batch as batch").
		ColumnExpr("l.location as location").
		ColumnExpr("s.status as status").
		ColumnExpr("u.name as user").
		Join("left join skus as sku").JoinOn("sku.id = items.sku_id").
		Join("left join batches as b").JoinOn("b.id = items.batch_id").
		Join("left join locations as l").JoinOn("l.id = items.location_id").
		Join("left join statuses as s").JoinOn("s.id = items.status_id").
		Join("left join users as u").JoinOn("u.id = items.user_id").
		Where("items.uic = ?", uic).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return PublicItem{}, err
	}
	return i, err
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

func (i *Item) AllocateItem(ctx context.Context) (err error) {
	i.UpdatedAt = time.Now()
	db := database.Pgdb

	// get complete current information about the item
	qry, err := db.Query(`
		SELECT
		    id,
		    sku_id,
		    status_id,
		    location_id
		FROM items
		WHERE uic = $1
	`, i.UIC)

	if err != nil {
		return
	}

	var currentLocation uint

	for qry.Next() {
		err = qry.Scan(&i.Id, &i.SkuID, &i.StatusID, &currentLocation)
		if err != nil || i.Id == 0 {
			err = errors.New(i.UIC)
			return
		}
	}

	// checks for redundante requests, e.g. moving the item to the same location where it already is at
	if i.StatusID == 4 || i.LocationID == currentLocation || i.Id == 0 {
		err = errors.New(i.UIC)
		return
	}

	// actual update
	cmd, err := db.Prepare(`
		UPDATE items
		SET location_id = l.id,
    		status_id = l.status_id,
    		user_id = $1,
    		updated_at = $2
		FROM (SELECT id, status_id FROM locations WHERE id = $3) as l
		WHERE items.uic = $4;
	`)
	if err != nil {
		return
	}
	defer cmd.Close()
	cmd.Exec(i.UserID, i.UpdatedAt, i.LocationID, i.UIC)

	// controller doesn't know the default status for the received location
	// here we get that value in order to create a record in the history table
	qry, err = db.Query(`
		SELECT
			status_id
		FROM locations
		WHERE id = $1;
	`, i.LocationID)

	if err != nil {
		return
	}

	var newStatusID uint

	for qry.Next() {
		err = qry.Scan(&newStatusID)
		if err != nil || i.Id == 0 {
			err = errors.New(i.UIC)
			return
		}
	}

	// create a history for this movement
	err = CreateHistory(ctx, &History{
		ItemID:     i.Id,
		SkuID:      i.SkuID,
		LocationID: i.LocationID,
		StatusID:   newStatusID,
		UserID:     i.UserID,
	})
	return
}
