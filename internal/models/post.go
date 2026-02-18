package models

import "time"

type Post struct {
	ID       int       `db:"id" json:"id"`
	UserID   int       `db:"user_id" json:"user_id"`
	Title    string    `db:"title" json:"title"`
	Content  string    `db:"content" json:"content"`
	Image    *string   `db:"image" json:"image"`
	CreateAt time.Time `db:"created_at" json:"created_at"`
}

type PostDetail struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	Image     *string   `db:"image" json:"image"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	Username  string    `db:"username" json:"author"`
}
