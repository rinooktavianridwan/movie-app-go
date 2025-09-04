package services

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/studio/repositories"
	"movie-app-go/internal/modules/studio/requests"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"

	"gorm.io/gorm"
)

type StudioService struct {
	StudioRepo   *repositories.StudioRepository
	FacilityRepo *repositories.FacilityRepository
}

func NewStudioService(studioRepo *repositories.StudioRepository, facilityRepo *repositories.FacilityRepository) *StudioService {
	return &StudioService{StudioRepo: studioRepo, FacilityRepo: facilityRepo}
}

func (s *StudioService) CreateStudio(req *requests.CreateStudioRequest) error {
	count, err := s.FacilityRepo.CountFacilitiesByIDs(req.FacilityIDs)
	if err != nil {
		return err
	}
	if count != int64(len(req.FacilityIDs)) {
		return utils.ErrInvalidFacilityIDs
	}

	studio := models.Studio{
		Name:         req.Name,
		SeatCapacity: req.SeatCapacity,
	}

	return s.StudioRepo.WithTransaction(func(tx *gorm.DB) error {
		if err := s.StudioRepo.CreateWithTx(tx, &studio); err != nil {
			return err
		}

		facilityStudios := make([]models.FacilityStudio, 0, len(req.FacilityIDs))
		for _, fid := range req.FacilityIDs {
			facilityStudios = append(facilityStudios, models.FacilityStudio{
				StudioID:   studio.ID,
				FacilityID: fid,
			})
		}

		return s.StudioRepo.CreateFacilityStudiosWithTx(tx, facilityStudios)
	})
}

func (s *StudioService) GetAllStudiosPaginated(page, perPage int) (repository.PaginationResult[models.Studio], error) {
	return s.StudioRepo.GetAllPaginated(page, perPage)
}

func (s *StudioService) GetStudioByID(id uint) (*models.Studio, error) {
	studio, err := s.StudioRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrStudioNotFound
		}
		return nil, err
	}
	return studio, nil
}

func (s *StudioService) UpdateStudio(id uint, req *requests.CreateStudioRequest) error {
	studio, err := s.GetStudioByID(id)
	if err != nil {
		return err
	}

	count, err := s.FacilityRepo.CountFacilitiesByIDs(req.FacilityIDs)
	if err != nil {
		return err
	}
	if count != int64(len(req.FacilityIDs)) {
		return utils.ErrInvalidFacilityIDs
	}

	studio.Name = req.Name
	studio.SeatCapacity = req.SeatCapacity

	return s.StudioRepo.WithTransaction(func(tx *gorm.DB) error {
		if err := s.StudioRepo.UpdateWithTx(tx, studio); err != nil {
			return err
		}

		if err := s.StudioRepo.DeleteFacilityStudiosWithTx(tx, id); err != nil {
			return err
		}

		facilityStudios := make([]models.FacilityStudio, 0, len(req.FacilityIDs))
		for _, fid := range req.FacilityIDs {
			facilityStudios = append(facilityStudios, models.FacilityStudio{
				StudioID:   studio.ID,
				FacilityID: fid,
			})
		}

		return s.StudioRepo.CreateFacilityStudiosWithTx(tx, facilityStudios)
	})
}

func (s *StudioService) DeleteStudio(id uint) error {
	_, err := s.GetStudioByID(id)
	if err != nil {
		return err
	}

	scheduleCount, err := s.StudioRepo.CountSchedulesByStudioID(id)
	if err != nil {
		return err
	}
	if scheduleCount > 0 {
		return utils.ErrStudioHasSchedules
	}

	return s.StudioRepo.WithTransaction(func(tx *gorm.DB) error {
        if err := s.StudioRepo.DeleteFacilityStudiosWithTx(tx, id); err != nil {
            return err
        }

        if err := s.StudioRepo.DeleteWithTx(tx, id); err != nil {
            return err
        }

        return nil
    })
}
