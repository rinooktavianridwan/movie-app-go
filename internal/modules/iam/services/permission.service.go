package services

import (
	"movie-app-go/internal/modules/iam/repositories"
	"movie-app-go/internal/modules/iam/responses"
	"movie-app-go/internal/repository"
)

type PermissionService struct {
	PermissionRepo *repositories.PermissionRepository
}

func NewPermissionService(permissionRepo *repositories.PermissionRepository) *PermissionService {
	return &PermissionService{PermissionRepo: permissionRepo}
}

func (s *PermissionService) GetAllPaginated(page, perPage int) (repository.PaginationResult[responses.PermissionResponse], error) {
	result, err := s.PermissionRepo.GetAllPaginated(page, perPage)
	if err != nil {
		return repository.PaginationResult[responses.PermissionResponse]{}, err
	}

	var permissionResponses []responses.PermissionResponse
	for _, permission := range result.Data {
		permissionResponses = append(permissionResponses, responses.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Resource:    permission.Resource,
			Action:      permission.Action,
			Description: permission.Description,
		})
	}

	return repository.PaginationResult[responses.PermissionResponse]{
		Data:       permissionResponses,
		Total:      result.Total,
		Page:       result.Page,
		PerPage:    result.PerPage,
		TotalPages: result.TotalPages,
	}, nil
}

func (s *PermissionService) GetByID(id uint) (*responses.PermissionResponse, error) {
	permission, err := s.PermissionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return &responses.PermissionResponse{
		ID:          permission.ID,
		Name:        permission.Name,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Description: permission.Description,
	}, nil
}

func (s *PermissionService) GetByResource(resource string) ([]responses.PermissionResponse, error) {
	permissions, err := s.PermissionRepo.GetByResource(resource)
	if err != nil {
		return nil, err
	}

	var result []responses.PermissionResponse
	for _, permission := range permissions {
		result = append(result, responses.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Resource:    permission.Resource,
			Action:      permission.Action,
			Description: permission.Description,
		})
	}

	return result, nil
}

func (s *PermissionService) GetAllGroupedByResource() (map[string][]responses.PermissionResponse, error) {
	permissions, err := s.PermissionRepo.GetAll()
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]responses.PermissionResponse)
	for _, permission := range permissions {
		resource := permission.Resource

		permissionResponse := responses.PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Resource:    permission.Resource,
			Action:      permission.Action,
			Description: permission.Description,
		}

		grouped[resource] = append(grouped[resource], permissionResponse)
	}

	return grouped, nil
}
