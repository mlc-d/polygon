package models

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/mlcprojects/wms/database"
)

type Location struct {
	Id        uint      `bun:",pk,autoincrement" json:"id,omitempty"`
	Location  string    `bun:",notnull,unique" json:"location,omitempty"`
	StatusID  uint      `json:"status_id,omitempty"`
	Status    *Status   `bun:"rel:belongs-to,join:status_id=id" json:"status,omitempty"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at,omitempty"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at,omitempty"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"deleted_at,omitempty"`
}

func CreateLocation(ctx context.Context, l *Location) (err error) {
	db := database.DB
	_, err = db.NewInsert().
		Model(l).
		Exec(ctx)
	return
}

func GetLocations(ctx context.Context) (locations []Location) {
	db := database.DB
	err := db.NewSelect().
		Model(&locations).
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}

func GetLocation(ctx context.Context, loc string) (l Location, err error) {
	db := database.DB
	err = db.NewSelect().
		Model(&l).
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Where("location = ?", loc).
		Scan(ctx)
	fmt.Println("encontré: ", l)
	if err != nil || l.Id < 1 {
		return Location{}, err
	}
	return l, nil
}

func GetLocationByStatus(ctx context.Context, statusId uint) (locations []Location, err error) {
	db := database.DB

	err = db.NewSelect().
		Model(&locations).
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Where("status_id = ?", statusId).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return locations, nil
}

func EditLocation(ctx context.Context, l *Location) (err error) {
	l.UpdatedAt = time.Now()
	db := database.DB
	_, err = db.NewUpdate().
		Model(l).
		OmitZero().
		WherePK().
		Exec(ctx)
	return
}
