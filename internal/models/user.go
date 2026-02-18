package models

import "time"

type User struct {
	ID        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"-"`
	Role      string    `db:"role" json:"role"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UserStats struct {
	TotalUsers int `db:"total_users" json:"total_users"`
	AdminCount int `db:"admin_count" json:"admin_count"`
	UserCount  int `json:"user_count"`
}

type UpdateRoleRequest struct {
	Role string `db:"role"`
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
