package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `db:"id" json:"id,omitempty"`
	Username  string    `db:"username" json:"username,omitempty"`
	Password  string    `db:"password" json:"password,omitempty"`
	Role      string    `db:"role" json:"role,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
}

type UserService interface {
	Register(ctx context.Context, user User) (err error)
	UpdateRole(ctx context.Context, id int64, role string) (err error)
	DeleteUser(ctx context.Context, id int64) (err error)
	Login(ctx context.Context, username, password string) (user User, err error)
}

func (User) TableName() string {
	return "users"
}