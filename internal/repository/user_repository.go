package repository

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rafli2460/culinary-blog-api/internal/config"
	"github.com/rafli2460/culinary-blog-api/internal/models"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	GetByUsername(ctx context.Context, username string) (models.User, error)
	Create(ctx context.Context, user *models.User) error

	GetAllUsers(ctx context.Context, user *models.User) error
	GetStats(ctx context.Context) (models.UserStats, error)
	UpdateRole(ctx context.Context, userID int, newRole string) error
	DeleteUser(ctx context.Context, userID int) error
}

type userRepository struct {
	db *config.Database
}

func NewUserRepository(db *config.Database) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	query := `SELECT id, username, password, created_at FROM users WHERE username = ?`
	err := r.db.Read.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, errors.New("user not found")
		}
		log.Error().Err(err).Str("username", username).Msg("Error Database: User not found")
		return user, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users(username, password, created_at) VALUES (:username, :password, NOW())`
	_, err := r.db.Write.NamedExecContext(ctx, query, user)
	if err != nil {
		log.Error().Err(err).Str("username", user.Username).Msg("Error Database: User already exist")
		return err
	}
	return nil
}

func (r *userRepository) GetAllUsers(ctx context.Context, user *models.User) error {
	panic("not implemented") // TODO: Implement
}

func (r *userRepository) GetStats(ctx context.Context) (models.UserStats, error) {
	panic("not implemented") // TODO: Implement
}

func (r *userRepository) UpdateRole(ctx context.Context, userID int, newRole string) error {
	panic("not implemented") // TODO: Implement
}

func (r *userRepository) DeleteUser(ctx context.Context, userID int) error {
	panic("not implemented") // TODO: Implement
}
