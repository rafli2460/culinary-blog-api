package users

import (
	"context"

	"github.com/rafli2460/culinary-blog-api/internal/domain"
	"github.com/rafli2460/culinary-blog-api/internal/server"
)

type Service struct {
	app  server.App
	repo UserRepository
}

func Init(ctx context.Context, a *server.App) (domain.UserService, *Handler) {
	r := initRepository(ctx, a)

	svc := &Service{
		app:  *a,
		repo: r,
	}

	h := NewHandler(svc)

	return svc, h
}

func (s *Service) Register(ctx context.Context, user domain.User) (err error) {
	err = s.repo.Insert(ctx, user)

	return
}

func (s *Service) UpdateRole(ctx context.Context, id int64, role string) (err error) {
	err = s.repo.UpdateRole(ctx, id, role)

	return
}

func (s *Service) DeleteUser(ctx context.Context, id int64) (err error) {
	err = s.repo.Delete(ctx, id)

	return
}

func (s *Service) Login(ctx context.Context, username string, password string) (user domain.User, err error) {
	_, err = s.repo.FindByUsername(ctx, username)

	return
}
