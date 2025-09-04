package repositories

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/repository"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	DB *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{DB: db}
}

func (r *PermissionRepository) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.Permission], error) {
	return repository.Paginate[models.Permission](
		r.DB.Model(&models.Permission{}),
		page,
		perPage,
	)
}

func (r *PermissionRepository) GetByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := r.DB.First(&permission, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) GetByResource(resource string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.DB.Where("resource = ?", resource).Order("action ASC").Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) GetAll() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.DB.Order("resource ASC, action ASC").Find(&permissions).Error
	return permissions, err
}
