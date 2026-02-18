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
		return errors.New("username only contain character, number, and underscore")
	}

	_, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return errors.New("username ini sudah digunakan")
	} else if err.Error() != "user not found" {
		return errors.New("terjadi kesalahan pada server")
	}

	req.Password = strings.TrimSpace(req.Password)
	if req.Password == "" {
		return errors.New("Please enter password")
	}
	if len(req.Password) < 6 {
		return errors.New("Password must be 6 characters long")
	}
	if req.Password != strings.TrimSpace(req.ConfirmPassword) {
		return errors.New("password is not match")
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
		return "", errors.New("wrong username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("wrong username or password")
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
