package models

import (
	"context"
	"gitlab.com/mlcprojects/wms/database"
	"time"
)

type Role struct {
	Id        uint      `bun:",pk,autoincrement" json:"id"`
	Role      string    `bun:",notnull,unique" json:"role"`
	IsAdmin   bool      `bun:",notnull" json:"-"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"-"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"-"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"-"`
}

func CreateRole(ctx context.Context, r *Role) (err error) {
	db := database.DB
	_, err = db.NewInsert().
		Model(r).
		Exec(ctx)
	return
}

func GetRoles(ctx context.Context) (roles []Role) {
	db := database.DB
	err := db.NewSelect().
		Model(&roles).
		ExcludeColumn("created_at").
		ExcludeColumn("deleted_at").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}
