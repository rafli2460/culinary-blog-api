package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rafli2460/culinary-blog-api/internal/models"
	"github.com/rafli2460/culinary-blog-api/internal/repository"
	"github.com/rafli2460/culinary-blog-api/pkg/logger"
	"github.com/rs/zerolog/log"
)

type PostService interface {
	CreatePost(ctx context.Context, userID int, title, content string, file *multipart.FileHeader) error
	DeletePost(ctx context.Context, postID, currentUserID int, currentUserRole string) error
	UpdatePost(ctx context.Context, postID, currentUserID int, currentUserRole string, title, content string, file *multipart.FileHeader) error
	GetPost(ctx context.Context, id int) (*models.PostDetail, error)
	GetAllPosts(ctx context.Context, page int, limit int) ([]models.PostDetail, error)
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{postRepo: repo}
}

func (s *postService) CreatePost(ctx context.Context, userID int, title string, content string, file *multipart.FileHeader) error {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" {
		return logger.ValidationError("title is required")
	}

	if content == "" {
		return logger.ValidationError("content is required")
	}

	var imageName *string

	if file != nil {
		const maxFileSize = 5 * 1024 * 1024
		if file.Size > maxFileSize {
			return logger.ValidationError("file size cannot exceed 5MB")
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{
			".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
		}
		if !allowedExts[ext] {
			return logger.ValidationError("file format is invalid (only JPG, PNG, GIF, WEBP)")
		}

		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return logger.LogError(err, "error creating upload directory")
		}

		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		dstPath := filepath.Join(uploadDir, newFileName)

		src, err := file.Open()
		if err != nil {
			return logger.LogError(err, "error reading uploaded file")
		}
		defer src.Close()

		dst, err := os.Create(dstPath)
		if err != nil {
			return logger.LogError(err, "failed to create file")
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return logger.LogError(err, "failed to save image file")
		}

		imageName = &newFileName
	}

	post := &models.Post{
		UserID:  userID,
		Title:   title,
		Content: content,
		Image:   imageName,
	}

	return s.postRepo.Create(ctx, post)
}

func (s *postService) DeletePost(ctx context.Context, postID int, currentUserID int, currentUserRole string) error {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return logger.ValidationError("post not found")
	}

	if post.UserID != currentUserID && currentUserRole != "admin" {
		return logger.ValidationError("access denied: you do not have permission to delete this post")
	}

	if post.Image != nil && *post.Image != "" {
		imagePath := filepath.Join("uploads", *post.Image)

		if err := os.Remove(imagePath); err != nil {
			log.Warn().Err(err).Str("file", imagePath).Msg("Failed to delete physical image file, it might not exist")
		} else {
			log.Info().Str("file", imagePath).Msg("Physical image file successfully deleted")
		}
	}

	return s.postRepo.Delete(ctx, postID)
}

func (s *postService) UpdatePost(ctx context.Context, postID int, currentUserID int, currentUserRole string, title, content string, file *multipart.FileHeader) error {
	existingPost, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return logger.ValidationError("post not found")
	}

	if existingPost.UserID != currentUserID && currentUserRole != "admin" {
		return logger.ValidationError("access denied: you do not have permission to edit this post")
	}

	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)
	if title == "" || content == "" {
		return logger.ValidationError("title and content cannot be empty")
	}

	finalImageName := existingPost.Image

	if file != nil {
		if file.Size > 5*1024*1024 {
			return logger.ValidationError("file size exceeds maximum limit of 5MB")
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
		if !allowedExts[ext] {
			return logger.ValidationError("invalid file format")
		}

		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		uploadDir := "./uploads"
		dstPath := filepath.Join(uploadDir, newFileName)

		src, err := file.Open()
		if err != nil {
			return logger.LogError(err, "failed to read new image file")
		}
		defer src.Close()

		dst, err := os.Create(dstPath)
		if err != nil {
			return logger.LogError(err, "failed to create file on server")
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return logger.LogError(err, "failed to save new image file")
		}

		if existingPost.Image != nil && *existingPost.Image != "" {
			oldImagePath := filepath.Join(uploadDir, *existingPost.Image)
			if err := os.Remove(oldImagePath); err != nil {
				log.Warn().Err(err).Str("file", oldImagePath).Msg("Failed to delete old image")
			}
		}

		finalImageName = &newFileName
	}

	existingPost.Title = title
	existingPost.Content = content
	existingPost.Image = finalImageName

	return s.postRepo.Update(ctx, existingPost)
}

func (s *postService) GetPost(ctx context.Context, id int) (*models.PostDetail, error) {
	return s.postRepo.GetPostDetailByID(ctx, id)
}

func (s *postService) GetAllPosts(ctx context.Context, page int, limit int) ([]models.PostDetail, error) {
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	return s.postRepo.GetAll(ctx, limit, offset)
}
