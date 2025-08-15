package services

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/studio/requests"
    "movie-app-go/internal/repository"
	"gorm.io/gorm"
)

type FacilityService struct {
	DB *gorm.DB
}

func NewFacilityService(db *gorm.DB) *FacilityService {
	return &FacilityService{DB: db}
}

func (s *FacilityService) CreateFacility(req *requests.CreateFacilityRequest) (*models.Facility, error){
	facility := models.Facility{Name: req.Name}
	if err := s.DB.Create(&facility).Error; err != nil {
		return nil, err
	}
	return &facility, nil
}

func (s *FacilityService) GetAllFacilitiesPaginated(page, perPage int) (repository.PaginationResult[models.Facility], error) {
    return repository.Paginate[models.Facility](s.DB, page, perPage)
}

func (s *FacilityService) GetFacilityByID(id uint) (*models.Facility, error) {
    var facility models.Facility
    if err := s.DB.First(&facility, id).Error; err != nil {
        return nil, err
    }
    return &facility, nil
}

func (s *FacilityService) UpdateFacility(id uint, req *requests.UpdateFacilityRequest) (*models.Facility, error) {
    var facility models.Facility
    if err := s.DB.First(&facility, id).Error; err != nil {
        return nil, err
    }
    facility.Name = req.Name
    if err := s.DB.Save(&facility).Error; err != nil {
        return nil, err
    }
    return &facility, nil
}

func (s *FacilityService) DeleteFacility(id uint) error {
    if err := s.DB.Delete(&models.Facility{}, id).Error; err != nil {
        return err
    }
    return nil
}