package models

import (
	"context"
	"errors"
	"github.com/uptrace/bun"
	"gitlab.com/mlcprojects/wms/database"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	Id        uint      `bun:",pk,autoincrement" json:"id"`
	Name      string    `bun:",notnull,unique" json:"name"`
	Password  string    `bun:",notnull,type:text" json:"password"`
	RoleID    uint      `json:"role_id"`
	Role      *Role     `bun:"rel:belongs-to,join:role_id=id"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"created_at"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"-"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"-"`
}

type PublicUser struct {
	Id        uint      `json:"id"`
	Name      string    `json:"name"`
	RoleID    uint      `json:"role_id"`
	Role      string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateUser(ctx context.Context, u *User) (err error) {
	db := database.DB
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 7)
	if err != nil {
		return
	}
	u.Password = string(password)
	if _, err = db.NewInsert().
		Model(u).
		Column("name", "password", "role_id", "created_at", "updated_at").
		Exec(ctx); err != nil {
		return
	}
	return
}

func GetUsers(ctx context.Context) (users []PublicUser) {
	db := database.DB
	err := db.NewSelect().
		Model(&users).
		ModelTableExpr("users").
		Column("users.id", "users.name", "r.role", "users.created_at").
		ColumnExpr("r.id as role_id").
		Join("left join roles as r").JoinOn("r.id = users.role_id").
		Where("users.deleted_at IS NULL").
		Order("users.id ASC").
		Scan(ctx)
	if err != nil {
		panic(err.Error())
	}
	return
}

func GetUser(ctx context.Context, u *User) (user User, err error) {
	db := database.DB
	err = db.NewSelect().
		Model(&user).
		Column("user.id").
		Column("user.name").
		Column("user.id").
		Column("user.role_id").
		Relation("Roles", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ColumnExpr("role as role_name")
		}).
		Column("user.created_at").
		Where("user.name = ? AND user.deleted_at IS NOT NULL", u.Name).
		Scan(ctx)
	if err != nil {
		err = errors.New("not in database")
	}
	return user, err
}

func UpdateUser(ctx context.Context, u *User) (err error) {
	u.UpdatedAt = time.Now()
	db := database.DB
	var newPassword []byte
	if u.Password != "" {
		newPassword, err = bcrypt.GenerateFromPassword([]byte(u.Password), 7)
		if err != nil {
			return
		}
		u.Password = string(newPassword)
	}
	_, err = db.NewUpdate().
		Model(u).
		OmitZero().
		WherePK().
		Exec(ctx)
	return
}

func DeleteUser(ctx context.Context, u *User) (err error) {
	db := database.DB
	_, err = db.NewDelete().
		Model(u).
		WherePK().
		Exec(ctx)
	return
}

func ValidateUser(ctx context.Context, u *User) (r, i uint, err error) {
	db := database.DB
	user := User{}
	err = db.NewSelect().
		Model(&user).
		Column("id").
		Column("password").
		Column("role_id").
		Where("name = ?", u.Name).
		Scan(ctx)
	if err != nil {
		return
	}
	return user.RoleID, user.Id, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
}
