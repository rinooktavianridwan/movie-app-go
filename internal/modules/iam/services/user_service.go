package services

import (
	"fmt"
	"mime/multipart"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/iam/requests"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (s *UserService) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.User], error) {
	return repository.Paginate[models.User](s.DB, page, perPage)
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Update(id uint, req *requests.UserUpdateRequest) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	user.Name = req.Name
	user.Email = req.Email
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) Delete(id uint) error {
	if err := s.DB.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateAvatar(userID uint, file *multipart.FileHeader) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user tidak ditemukan")
	}

	if user.Avatar != nil && *user.Avatar != "" {
		utils.DeleteFile(*user.Avatar)
	}

	avatarPath, err := utils.SaveFile(file, "uploads/avatars", "image", 5)
	if err != nil {
		return nil, err
	}

	relativePath := strings.TrimPrefix(avatarPath, "./")
	user.Avatar = &relativePath

	if err := s.DB.Save(&user).Error; err != nil {
		utils.DeleteFile(avatarPath)
		return nil, err
	}

	return &user, nil
}
