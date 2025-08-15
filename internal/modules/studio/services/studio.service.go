package services

import (
	"errors"
	"gorm.io/gorm"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/studio/requests"
	"movie-app-go/internal/repository"
)

type StudioService struct {
	DB *gorm.DB
}

func NewStudioService(db *gorm.DB) *StudioService {
	return &StudioService{DB: db}
}

func (s *StudioService) CreateStudio(req *requests.CreateStudioRequest) (*models.Studio, error) {
	var studio *models.Studio

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		studio = &models.Studio{
			Name:         req.Name,
			SeatCapacity: req.SeatCapacity,
		}
		if err := tx.Create(studio).Error; err != nil {
			return err
		}

		var count int64
		if err := tx.Model(&models.Facility{}).Where("id IN ?", req.FacilityIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(req.FacilityIDs)) {
			return errors.New("some facility_ids are invalid")
		}
		
		facilityStudios := make([]models.FacilityStudio, 0, len(req.FacilityIDs))
		for _, fid := range req.FacilityIDs {
			facilityStudios = append(facilityStudios, models.FacilityStudio{
				StudioID:   studio.ID,
				FacilityID: fid,
			})
		}
		if len(facilityStudios) > 0 {
			if err := tx.Create(&facilityStudios).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return studio, nil
}

func (s *StudioService) GetAllStudiosPaginated(page, perPage int) (repository.PaginationResult[models.Studio], error) {
	return repository.PaginateWithPreload[models.Studio](s.DB, page, perPage, "FacilityStudios.Facility")
}

func (s *StudioService) GetStudioByID(id uint) (*models.Studio, error) {
	var studio models.Studio
	if err := s.DB.Preload("FacilityStudios.Facility").First(&studio, id).Error; err != nil {
		return nil, err
	}
	return &studio, nil
}

func (s *StudioService) UpdateStudio(id uint, req *requests.CreateStudioRequest) (*models.Studio, error) {
	var studio models.Studio
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&studio, id).Error; err != nil {
			return err
		}
		studio.Name = req.Name
		studio.SeatCapacity = req.SeatCapacity
		if err := tx.Save(&studio).Error; err != nil {
			return err
		}
		if err := tx.Where("studio_id = ?", id).Delete(&models.FacilityStudio{}).Error; err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&models.Facility{}).Where("id IN ?", req.FacilityIDs).Count(&count).Error; err != nil {
			return err
		}
		if count != int64(len(req.FacilityIDs)) {
			return errors.New("some facility_ids are invalid")
		}
		
		facilityStudios := make([]models.FacilityStudio, 0, len(req.FacilityIDs))
		for _, fid := range req.FacilityIDs {
			facilityStudios = append(facilityStudios, models.FacilityStudio{
				StudioID:   studio.ID,
				FacilityID: fid,
			})
		}
		if len(facilityStudios) > 0 {
			if err := tx.Create(&facilityStudios).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &studio, nil
}

func (s *StudioService) DeleteStudio(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// Hapus relasi facility_studio dulu
		if err := tx.Where("studio_id = ?", id).Delete(&models.FacilityStudio{}).Error; err != nil {
			return err
		}
		// Hapus studio
		if err := tx.Delete(&models.Studio{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
