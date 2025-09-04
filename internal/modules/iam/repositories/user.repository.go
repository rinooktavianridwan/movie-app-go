package repositories

import (
    "movie-app-go/internal/models"
    "movie-app-go/internal/repository"
    "gorm.io/gorm"
)

type UserRepository struct {
    DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{
        DB: db,
    }
}

func (r *UserRepository) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.User], error) {
    return repository.Paginate[models.User](r.DB, page, perPage)
}

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
    var user models.User
    if err := r.DB.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *UserRepository) ExistsByEmailExceptID(email string, userID uint) (bool, error) {
    var count int64
    err := r.DB.Model(&models.User{}).Where("email = ? AND id != ?", email, userID).Count(&count).Error
    return count > 0, err
}

func (r *UserRepository) Update(user *models.User) error {
    return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
    return r.DB.Delete(&models.User{}, id).Error
}
