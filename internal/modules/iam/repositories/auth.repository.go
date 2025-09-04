package repositories

import (
	"movie-app-go/internal/models"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		DB: db,
	}
}

func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Preload("Role.Permissions").Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }

	return &user, nil
}

func (r *AuthRepository) GetUserByEmailWithRole(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Preload("Role.Permissions").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *AuthRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
