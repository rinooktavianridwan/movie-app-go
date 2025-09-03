package services

import (
	"movie-app-go/internal/models"
	movierepos "movie-app-go/internal/modules/movie/repositories"
	"movie-app-go/internal/modules/schedule/options"
	"movie-app-go/internal/modules/schedule/repositories"
	"movie-app-go/internal/modules/schedule/requests"
	studiorepos "movie-app-go/internal/modules/studio/repositories"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"
	"time"

	"gorm.io/gorm"
)

type ScheduleService struct {
	ScheduleRepo *repositories.ScheduleRepository
	MovieRepo    *movierepos.MovieRepository
	StudioRepo   *studiorepos.StudioRepository
}

func NewScheduleService(
	scheduleRepo *repositories.ScheduleRepository,
	movieRepo *movierepos.MovieRepository,
	studioRepo *studiorepos.StudioRepository,
) *ScheduleService {
	return &ScheduleService{
		ScheduleRepo: scheduleRepo,
		MovieRepo:    movieRepo,
		StudioRepo:   studioRepo,
	}
}

func (s *ScheduleService) GetAllSchedulesPaginated(opts *options.GetAllScheduleOptions) (repository.PaginationResult[models.Schedule], error) {
	return s.ScheduleRepo.GetAllWithOptions(opts)
}

func (s *ScheduleService) GetScheduleByID(id uint) (*models.Schedule, error) {
	schedule, err := s.ScheduleRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrScheduleNotFound
		}
		return nil, err
	}
	return schedule, nil
}

func (s *ScheduleService) CreateSchedule(req *requests.CreateScheduleRequest) error {
	movie, err := s.MovieRepo.GetByID(req.MovieID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrMovieNotFound
		}
		return err
	}

	_, err = s.StudioRepo.GetByID(req.StudioID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrStudioNotFound
		}
		return err
	}

	if req.Date.Before(time.Now().Truncate(24 * time.Hour)) {
		return utils.ErrPastDate
	}

	return s.ScheduleRepo.WithTransaction(func(tx *gorm.DB) error {
		bufferMinutes := 30
		totalMinutes := int(movie.Duration) + bufferMinutes
		endTime := req.StartTime.Add(time.Duration(totalMinutes) * time.Minute)

		hasConflict, err := s.ScheduleRepo.CheckScheduleConflict(
			req.StudioID, req.StartTime, endTime, nil)
		if err != nil {
			return err
		}
		if hasConflict {
			return utils.ErrScheduleConflict
		}

		schedule := models.Schedule{
			MovieID:   req.MovieID,
			StudioID:  req.StudioID,
			StartTime: req.StartTime,
			EndTime:   endTime,
			Date:      req.Date,
			Price:     req.Price,
		}

		return s.ScheduleRepo.CreateWithTx(tx, &schedule)
	})
}

func (s *ScheduleService) UpdateSchedule(id uint, req *requests.UpdateScheduleRequest) error {
	return s.ScheduleRepo.WithTransaction(func(tx *gorm.DB) error {
		schedule, err := s.ScheduleRepo.GetByIDWithTx(tx, id)
		if err != nil {
			return err
		}

		if req.MovieID != nil {
			_, err := s.MovieRepo.GetByID(*req.MovieID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return utils.ErrMovieNotFound
				}
				return err
			}
			schedule.MovieID = *req.MovieID
		}

		if req.StudioID != nil {
			_, err := s.StudioRepo.GetByID(*req.StudioID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return utils.ErrStudioNotFound
				}
				return err
			}
			schedule.StudioID = *req.StudioID
		}

		if req.StartTime != nil {
			schedule.StartTime = *req.StartTime
		}

		if req.Date != nil {
			if req.Date.Before(time.Now().Truncate(24 * time.Hour)) {
				return utils.ErrPastDate
			}
			schedule.Date = *req.Date
		}

		if req.Price != nil {
			schedule.Price = *req.Price
		}

		if req.MovieID != nil || req.StartTime != nil {
			movie, err := s.MovieRepo.GetByID(schedule.MovieID)
			if err != nil {
				return err
			}

			bufferMinutes := 30
			totalMinutes := int(movie.Duration) + bufferMinutes
			schedule.EndTime = schedule.StartTime.Add(time.Duration(totalMinutes) * time.Minute)
		}

		hasConflict, err := s.ScheduleRepo.CheckScheduleConflict(
			schedule.StudioID, schedule.StartTime, schedule.EndTime, &id)
		if err != nil {
			return err
		}
		if hasConflict {
			return utils.ErrScheduleConflict
		}

		return s.ScheduleRepo.UpdateWithTx(tx, schedule)
	})
}

func (s *ScheduleService) DeleteSchedule(id uint) error {
	_, err := s.GetScheduleByID(id)
	if err != nil {
		return err
	}

	ticketCount, err := s.ScheduleRepo.CountTicketsByScheduleID(id)
	if err != nil {
		return err
	}
	if ticketCount > 0 {
		return utils.ErrScheduleHasTickets
	}

	return s.ScheduleRepo.WithTransaction(func(tx *gorm.DB) error {
		return s.ScheduleRepo.DeleteWithTx(tx, id)
	})
}
