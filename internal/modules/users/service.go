package users

import (
	"context"

	"github.com/rafli2460/culinary-blog-api/internal/domain"
)

type UserService interface {
	Register(ctx context.Context, user domain.User) error
	UpdateRole(ctx context.Context, id int64, role string) error
	DeleteUser(ctx context.Context, id int64) error
	Login(ctx context.Context, username, password string) (string, error)
}

