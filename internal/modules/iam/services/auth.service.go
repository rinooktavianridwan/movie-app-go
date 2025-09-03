package services

import (
	"fmt"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/iam/requests"
	"movie-app-go/internal/utils"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Register(req *requests.RegisterRequest) error {
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
		return err
	}
	user.Password = string(hashed)

	if err := s.DB.Unscoped().Where("email = ?", user.Email).First(&existing).Error; err == nil {
		if existing.DeletedAt.Valid {
			existing.DeletedAt = gorm.DeletedAt{}
			existing.Name = user.Name
			existing.Password = user.Password
			existing.IsAdmin = user.IsAdmin
			if err := s.DB.Unscoped().Save(&existing).Error; err != nil {
				return err
			}
			return nil
		}
		return utils.ErrEmailAlreadyExists
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *AuthService) Login(req *requests.LoginRequest) (*models.User, string, error) {
	var user models.User
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
            return nil, "", utils.ErrInvalidCredentials
        }
        return nil, "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", utils.ErrInvalidCredentials
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

func (s *AuthService) Logout(tokenString string) error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret"
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["exp"] == nil {
		return fmt.Errorf("invalid token claims")
	}
	exp := int64(claims["exp"].(float64)) - time.Now().Unix()
	if exp < 0 {
		exp = 0
	}
	return BlacklistToken(tokenString, time.Duration(exp)*time.Second)
}
