package models

import (
	"context"
	"errors"
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
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp,type:timestamptz" json:"updated_at"`
	DeletedAt time.Time `bun:",soft_delete,nullzero" json:"deleted_at"`
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
		Exec(ctx); err != nil {
		return
	}
	return
}

func GetUsers(ctx context.Context) (users []User) {
	db := database.DB
	err := db.NewSelect().
		Model(&users).
		ExcludeColumn("password").
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
		ExcludeColumn("password").
		ExcludeColumn("updated_at").
		ExcludeColumn("deleted_at").
		Where("id = ?", u.Id).
		Scan(ctx)
	if err != nil {
		err = errors.New("not in database")
	}
	return user, err
}

func UpdateUser(ctx context.Context, u *User) (err error) {
	db := database.DB
	newPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 7)
	if err != nil {
		return
	}
	_, err = db.NewUpdate().
		Model(u).
		Value("password", string(newPassword)).
		Exec(ctx)
	return
}

func ValidateUser(ctx context.Context, u *User) (r uint, err error) {
	db := database.DB
	user := User{}
	err = db.NewSelect().
		Model(&user).
		Column("password").
		Column("role_id").
		Where("name = ?", u.Name).
		Scan(ctx)
	if err != nil {
		return
	}
	return user.RoleID, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
}
