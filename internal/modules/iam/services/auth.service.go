package services

import (
	"fmt"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/iam/repositories"
	"movie-app-go/internal/modules/iam/requests"
	"movie-app-go/internal/utils"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	AuthRepo *repositories.AuthRepository
}

func NewAuthService(authRepo *repositories.AuthRepository) *AuthService {
	return &AuthService{AuthRepo: authRepo}
}

func (s *AuthService) Register(req *requests.RegisterRequest) error {
	exists, err := s.AuthRepo.ExistsByEmail(req.Email)
    if err != nil {
        return err
    }
    if exists {
        return utils.ErrEmailAlreadyExists
    }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

	user := models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: string(hashedPassword),
        IsAdmin:  false,
    }

	return s.AuthRepo.Create(&user)
}

func (s *AuthService) Login(req *requests.LoginRequest) (*models.User, string, error) {
	user, err := s.AuthRepo.GetUserByEmail(req.Email)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, "", utils.ErrInvalidCredentials
        }
        return nil, "", err
    }

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
        return nil, "", utils.ErrInvalidCredentials
    }

	token, err := s.generateJWT(user)
    if err != nil {
        return nil, "", err
    }

    return user, token, nil
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

func (s *AuthService) generateJWT(user *models.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id":  user.ID,
        "email":    user.Email,
        "is_admin": user.IsAdmin,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}