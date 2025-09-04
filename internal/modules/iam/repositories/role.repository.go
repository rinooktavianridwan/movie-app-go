package repositories

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/repository"

	"gorm.io/gorm"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

func (r *RoleRepository) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.Role], error) {
	query := r.DB.Preload("Permissions")
	return repository.Paginate[models.Role](query, page, perPage)
}

func (r *RoleRepository) GetByID(id uint) (*models.Role, error) {
	var role models.Role
	if err := r.DB.Preload("Permissions").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.DB.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetByName(name string) (*models.Role, error) {
    var role models.Role
    if err := r.DB.Where("name = ?", name).First(&role).Error; err != nil {
        return nil, err
    }
    return &role, nil
}
