package service

import (
	"context"
	"errors"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rafli2460/culinary-blog-api/internal/models"
	"github.com/rafli2460/culinary-blog-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req models.RegisterRequest) error
	Login(ctx context.Context, req models.LoginRequest) (string, error)

	GetAllUsers(ctx context.Context, search string) ([]models.User, error)
	GetStats(ctx context.Context) (models.UserStats, error)
	UpdateRole(ctx context.Context, targetUserID int, currentAdminID int, newRole string) error
	DeleteUser(ctx context.Context, targetUserID int, currentAdminID int) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) Register(ctx context.Context, req models.RegisterRequest) error {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return errors.New("Please enter username")
	}
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(req.Username)
	if !validUsername {
		return errors.New("username can only contain letters, numbers, and underscores")
	}

	_, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return errors.New("username is already taken")
	} else if err.Error() != "user not found" {
		return errors.New("internal server error")
	}

	req.Password = strings.TrimSpace(req.Password)
	if req.Password == "" {
		return errors.New("Please enter password")
	}
	if len(req.Password) < 6 {
		return errors.New("Password must be 6 characters long")
	}
	if req.Password != strings.TrimSpace(req.ConfirmPassword) {
		return errors.New("passwords do not match")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("error generating password")
	}

	newUser := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	return s.userRepo.Create(ctx, newUser)
}

func (s *userService) Login(ctx context.Context, req models.LoginRequest) (string, error) {
	req.Username = strings.TrimSpace(req.Username)

	if req.Username == "" || req.Password == "" {
		return "", errors.New("username and password cannot be empty")
	}

	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.New("error creating authentication token")
	}

	return tokenString, nil
}

func (s *userService) GetAllUsers(ctx context.Context, search string) ([]models.User, error) {
	search = strings.TrimSpace(search)
	return s.userRepo.GetAllUsers(ctx, search)
}

func (s *userService) GetStats(ctx context.Context) (models.UserStats, error) {
	stats, err := s.userRepo.GetStats(ctx)
	if err != nil {
		return stats, err
	}

	stats.UserCount = stats.TotalUsers - stats.AdminCount

	return stats, nil
}

func (s *userService) UpdateRole(ctx context.Context, targetUserID int, currentAdminID int, newRole string) error {
	if targetUserID == currentAdminID {
		return errors.New("action denied: you can't change your own role")
	}

	newRole = strings.ToLower(strings.TrimSpace(newRole))
	if newRole != "admin" && newRole != "user" {
		return errors.New("invalid role, must be 'admin' or 'user'")
	}

	return s.userRepo.UpdateRole(ctx, targetUserID, newRole)
}

func (s *userService) DeleteUser(ctx context.Context, targetUserID int, currentAdminID int) error {
	if targetUserID == currentAdminID {
		return errors.New("action denied: you can't delete your own account")
	}

	return s.userRepo.Delete(ctx, targetUserID)
}
