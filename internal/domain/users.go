package domain

import "time"

type User struct {
	ID        int64     `db:"id" json:"id,omitempty"`
	Username  string    `db:"username" json:"username,omitempty"`
	Password  string    `db:"password" json:"password,omitempty"`
	Role      string    `db:"role" json:"role,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at,omitempty"`
}

func (User) TableName() string {
	return "users"
}