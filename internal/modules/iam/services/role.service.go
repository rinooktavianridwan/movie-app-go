package services

import (
    "movie-app-go/internal/models"
    "movie-app-go/internal/modules/iam/repositories"
    "movie-app-go/internal/repository"
)

type RoleService struct {
    RoleRepo *repositories.RoleRepository
}

func NewRoleService(roleRepo *repositories.RoleRepository) *RoleService {
    return &RoleService{RoleRepo: roleRepo}
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {
    return s.RoleRepo.GetAll()
}

func (s *RoleService) GetAllRolesPaginated(page, perPage int) (repository.PaginationResult[models.Role], error) {
    return s.RoleRepo.GetAllPaginated(page, perPage)
}

func (s *RoleService) GetRoleByID(id uint) (*models.Role, error) {
    return s.RoleRepo.GetByID(id)
}