package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/iam/requests"
	"os"
	"time"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Register(req *requests.RegisterRequest) (*models.User, error) {
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}

	var existing models.User
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashed)

	if err := s.DB.Unscoped().Where("email = ?", user.Email).First(&existing).Error; err == nil {
		if existing.DeletedAt.Valid {
			existing.DeletedAt = gorm.DeletedAt{}
			existing.Name = user.Name
			existing.Password = user.Password
			existing.IsAdmin = user.IsAdmin
			if err := s.DB.Unscoped().Save(&existing).Error; err != nil {
				return nil, err
			}
			return &existing, nil
		}
		return nil, fmt.Errorf("email sudah terdaftar")
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) Login(req *requests.LoginRequest) (*models.User, string, error) {
	var user models.User
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, "", fmt.Errorf("email tidak ditemukan")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", fmt.Errorf("email atau password salah")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"email":    user.Email,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, "", err
	}
	return &user, tokenString, nil
}
