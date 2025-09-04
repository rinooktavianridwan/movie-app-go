package services

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/studio/repositories"
	"movie-app-go/internal/modules/studio/requests"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"

	"gorm.io/gorm"
)

type FacilityService struct {
	FacilityRepo *repositories.FacilityRepository
}

func NewFacilityService(facilityRepo *repositories.FacilityRepository) *FacilityService {
	return &FacilityService{FacilityRepo: facilityRepo}
}

func (s *FacilityService) CreateFacility(req *requests.CreateFacilityRequest) error {
	exists, err := s.FacilityRepo.ExistsByName(req.Name)
	if err != nil {
		return err
	}
	if exists {
		return utils.ErrFacilityAlreadyExists
	}

	facility := models.Facility{
		Name: req.Name,
	}

	return s.FacilityRepo.Create(&facility)
}

func (s *FacilityService) GetAllFacilitiesPaginated(page, perPage int) (repository.PaginationResult[models.Facility], error) {
	return s.FacilityRepo.GetAllPaginated(page, perPage)
}

func (s *FacilityService) GetFacilityByID(id uint) (*models.Facility, error) {
	facility, err := s.FacilityRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrFacilityNotFound
		}
		return nil, err
	}
	return facility, nil
}

func (s *FacilityService) UpdateFacility(id uint, req *requests.UpdateFacilityRequest) error {
	facility, err := s.GetFacilityByID(id)
	if err != nil {
		return err
	}

	if req.Name != facility.Name {
		exists, err := s.FacilityRepo.ExistsByNameExceptID(req.Name, id)
		if err != nil {
			return err
		}
		if exists {
			return utils.ErrFacilityAlreadyExists
		}
	}

	facility.Name = req.Name

	return s.FacilityRepo.Update(facility)
}

func (s *FacilityService) DeleteFacility(id uint) error {
	_, err := s.GetFacilityByID(id)
	if err != nil {
		return err
	}

	studioCount, err := s.FacilityRepo.CountStudiosByFacilityID(id)
	if err != nil {
		return err
	}
	if studioCount > 0 {
		return utils.ErrFacilityInUse
	}

	return s.FacilityRepo.Delete(id)
}
