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

	GetAllUsers(ctx context.Context, search string) ([]models.User, error)
	GetStats(ctx context.Context) (models.UserStats, error)
	UpdateRole(ctx context.Context, userID int, newRole string) error
	Delete(ctx context.Context, userID int) error
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
		log.Error().Err(err).Str("username", user.Username).Msg("Error Database: User already exists")
		return err
	}
	return nil
}

func (r *userRepository) GetAllUsers(ctx context.Context, search string) ([]models.User, error) {
	var users []models.User
	query := `SELECT id, username, role, created_at FROM users`
	var args []interface{}

	if search != "" {
		query += ` WHERE username like ? OR role LIKE ?`
		likeSearch := "%" + search + "%"
		args = append(args, likeSearch, likeSearch)
	}

	query += ` ORDER BY created_at DESC`

	err := r.db.Read.SelectContext(ctx, &users, query, args...)
	if err != nil {
		log.Error().Err(err).Str("search", search).Msg("Failed to gather user data")
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetStats(ctx context.Context) (models.UserStats, error) {
	var stats models.UserStats

	query := `
		SELECT
			COUNT(*) as total_users,
			COALESCE(SUM(CASE WHEN role = 'admin' THEN 1 ELSE 0 END), 0) as admin_count
		FROM users
	`

	err := r.db.Read.GetContext(ctx, &stats, query)
	if err != nil {
		log.Error().Err(err).Msg("error gathering user statistics")
		return stats, err
	}

	return stats, nil
}

func (r *userRepository) UpdateRole(ctx context.Context, userID int, newRole string) error {
	query := `UPDATE users SET role = ? WHERE id = ?`

	_, err := r.db.Write.ExecContext(ctx, query, newRole, userID)
	if err != nil {
		log.Error().Err(err).Int("user_id", userID).Str("new_role", newRole).Msg("error changing role")
		return err
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Write.ExecContext(ctx, query, userID)
	if err != nil {
		log.Error().Err(err).Int("user_id", userID).Msg("error deleting user")
		return err
	}

	return nil
}
