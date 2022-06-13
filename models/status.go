package models

import (
	"context"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type Status struct {
	Id        uint      `bun:",pk,autoincrement" json:"id"`
	Status    string    `bun:",notnull,unique" json:"status"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"-"`
}

type PublicStatus struct {
}

func (s *Status) CreateStatus(ctx context.Context) (err error) {
	db := database.DB
	_, err = db.NewInsert().
		Model(s).
		Exec(ctx)
	return
}

func GetStatuses(ctx context.Context) (statuses []Status) {
	db := database.DB
	err := db.NewSelect().
		Model(&statuses).
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
