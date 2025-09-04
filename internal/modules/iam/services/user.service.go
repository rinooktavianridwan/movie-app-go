package services

import (
	"mime/multipart"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/iam/repositories"
	"movie-app-go/internal/modules/iam/requests"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.User], error) {
	return s.UserRepo.GetAllPaginated(page, perPage)
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	user, err := s.UserRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(id uint, req *requests.UserUpdateRequest) error {
	user, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if req.Email != user.Email {
		exists, err := s.UserRepo.ExistsByEmailExceptID(req.Email, id)
		if err != nil {
			return err
		}
		if exists {
			return utils.ErrEmailAlreadyExists
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashed)
	}
	if req.IsAdmin != nil {
		user.IsAdmin = *req.IsAdmin
	}

	return s.UserRepo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	_, err := s.GetByID(id)
	if err != nil {
        return err
    }

	return s.UserRepo.Delete(id)
}

func (s *UserService) UpdateAvatar(userID uint, file *multipart.FileHeader) error {
	user, err := s.GetByID(userID)
    if err != nil {
        return err
    }

	if user.Avatar != nil && *user.Avatar != "" {
        utils.DeleteFile(*user.Avatar)
    }

	avatarPath, err := utils.SaveFile(file, "uploads/avatars", "image", 5)
    if err != nil {
        return err
    }

	relativePath := strings.TrimPrefix(avatarPath, "./")
    user.Avatar = &relativePath

	if err := s.UserRepo.Update(user); err != nil {
        utils.DeleteFile(avatarPath)
        return err
    }

    return nil
}
