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

type PostRepository interface {
	Create(ctx context.Context, post *models.Post) error
	GetByID(ctx context.Context, id int) (*models.Post, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, post *models.Post) error
	GetAll(ctx context.Context, limit int, offset int) ([]models.PostDetail, error)
	GetPostDetailByID(ctx context.Context, id int) (*models.PostDetail, error)
}

type postRepository struct {
	db *config.Database
}

func NewPostRepository(db *config.Database) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *models.Post) error {
	query := `INSERT INTO posts(user_id, title, content, image, created_at)
			  VALUES(:user_id, :title, :content, :image, NOW())`
	_, err := r.db.Write.NamedExecContext(ctx, query, post)
	if err != nil {
		return logger.LogErrorWithFields(err, "failed to save post into database", map[string]interface{}{
			"title": post.Title,
		})
	}
	return nil
}

func (r *postRepository) GetByID(ctx context.Context, id int) (*models.Post, error) {
	var post models.Post
	query := `SELECT id, user_id, title, content, image, created_at FROM posts WHERE id = ?`

	err := r.db.Read.GetContext(ctx, &post, query, id)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *postRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = ?`

	_, err := r.db.Write.ExecContext(ctx, query, id)
	if err != nil {
		return logger.LogErrorWithFields(err, "Failed to delete post from database", map[string]interface{}{
			"post_id": id,
		})
	}
	return nil
}

func (r *postRepository) Update(ctx context.Context, post *models.Post) error {
	query := `UPDATE posts SET title = :title, content = :content, image = :image WHERE id = :id`
	_, err := r.db.Write.NamedExecContext(ctx, query, post)
	if err != nil {
		return logger.LogErrorWithFields(err, "failed to update post in database", map[string]interface{}{
			"post_id": post.ID,
		})
	}
	return nil
}

func (r *postRepository) GetPostDetailByID(ctx context.Context, id int) (*models.PostDetail, error) {
	var post models.PostDetail

	query := `
		SELECT posts.id, posts.title, posts.content, posts.image, posts.created_at, users.username 
		FROM posts 
		JOIN users ON posts.user_id = users.id 
		WHERE posts.id = ?`

	err := r.db.Read.GetContext(ctx, &post, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, logger.ValidationError("post not found")
		}
		return nil, logger.LogErrorWithFields(err, "Failed to retrieve post details", map[string]interface{}{
			"id": id,
		})
	}

	return &post, nil
}

func (r *postRepository) GetAll(ctx context.Context, limit int, offset int) ([]models.PostDetail, error) {
	posts := make([]models.PostDetail, 0)

	query := `
		SELECT posts.id, posts.title, posts.content, posts.image, posts.created_at, users.username 
		FROM posts 
		JOIN users ON posts.user_id = users.id 
		ORDER BY posts.created_at DESC
		LIMIT ? OFFSET ?`

	err := r.db.Read.SelectContext(ctx, &posts, query, limit, offset)
	if err != nil {
		return nil, logger.LogError(err, "Failed to retrieve post list")
	}

	return posts, nil
}
