package repository

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rafli2460/culinary-blog-api/internal/config"
	"github.com/rafli2460/culinary-blog-api/internal/models"
	"github.com/rafli2460/culinary-blog-api/pkg/logger"
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
			return user, logger.ValidationError("user not found")
		}
		return user, logger.LogErrorWithFields(err, "Error Database: User not found", map[string]interface{}{
			"username": username,
		})
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users(username, password, created_at) VALUES (:username, :password, NOW())`
	_, err := r.db.Write.NamedExecContext(ctx, query, user)
	if err != nil {
		return logger.LogErrorWithFields(err, "Error Database: User already exists", map[string]interface{}{
			"username": user.Username,
		})
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
		return nil, logger.LogErrorWithFields(err, "Failed to gather user data", map[string]interface{}{
			"search": search,
		})
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
		return stats, logger.LogError(err, "error gathering user statistics")
	}

	return stats, nil
}

func (r *userRepository) UpdateRole(ctx context.Context, userID int, newRole string) error {
	query := `UPDATE users SET role = ? WHERE id = ?`

	_, err := r.db.Write.ExecContext(ctx, query, newRole, userID)
	if err != nil {
		return logger.LogErrorWithFields(err, "error changing role", map[string]interface{}{
			"user_id":  userID,
			"new_role": newRole,
		})
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := r.db.Write.ExecContext(ctx, query, userID)
	if err != nil {
		return logger.LogErrorWithFields(err, "error deleting user", map[string]interface{}{
			"user_id": userID,
		})
	}

	return nil
}
