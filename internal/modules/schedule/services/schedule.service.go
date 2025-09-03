package services

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/schedule/options"
	"movie-app-go/internal/modules/schedule/requests"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"
	"time"

	"gorm.io/gorm"
)

type ScheduleService struct {
	DB *gorm.DB
}

func NewScheduleService(db *gorm.DB) *ScheduleService {
	return &ScheduleService{DB: db}
}

func (s *ScheduleService) GetAllSchedulesPaginated(opts *options.GetAllScheduleOptions) (repository.PaginationResult[models.Schedule], error) {
	query := s.DB.Preload("Movie").Preload("Studio")

	if opts.MovieTitle != "" {
		query = query.Joins("JOIN movies ON schedules.movie_id = movies.id").
			Where("movies.title ILIKE ?", "%"+opts.MovieTitle+"%")
	}

	if opts.StudioID != nil {
		query = query.Where("schedules.studio_id = ?", *opts.StudioID)
	}

	if opts.DateFrom != nil {
		query = query.Where("schedules.date >= ?", *opts.DateFrom)
	}
	if opts.DateTo != nil {
		query = query.Where("schedules.date <= ?", *opts.DateTo)
	}

	query = query.Order("date ASC, start_time ASC")

	return repository.Paginate[models.Schedule](query, opts.Page, opts.PerPage)
}

func (s *ScheduleService) GetScheduleByID(id uint) (*models.Schedule, error) {
	var schedule models.Schedule
	if err := s.DB.Preload("Movie").Preload("Studio").First(&schedule, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrScheduleNotFound
		}
		return nil, err
	}
	return &schedule, nil
}

func (s *ScheduleService) CreateSchedule(req *requests.CreateScheduleRequest) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var movie models.Movie
		if err := tx.First(&movie, req.MovieID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.ErrInvalidMovieIDs
			}
			return err
		}

		var studio models.Studio
		if err := tx.First(&studio, req.StudioID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.ErrStudioNotFound
			}
			return err
		}

		bufferMinutes := 30
		totalMinutes := int(movie.Duration) + bufferMinutes
		endTime := req.StartTime.Add(time.Duration(totalMinutes) * time.Minute)

		if req.Date.Before(time.Now().Truncate(24 * time.Hour)) {
			return utils.ErrPastDate
		}

		var conflictCount int64
		if err := tx.Model(&models.Schedule{}).
			Where("studio_id = ? AND date = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
				req.StudioID, req.Date, req.StartTime, req.StartTime, endTime, endTime, req.StartTime, endTime).
			Count(&conflictCount).Error; err != nil {
			return err
		}
		if conflictCount > 0 {
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

		if err := tx.Create(&schedule).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *ScheduleService) UpdateSchedule(id uint, req *requests.UpdateScheduleRequest) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var schedule models.Schedule
		if err := tx.First(&schedule, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.ErrScheduleNotFound
			}
			return err
		}

		if req.MovieID != nil {
			var movie models.Movie
			if err := tx.First(&movie, req.MovieID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return utils.ErrInvalidMovieIDs
				}
				return err
			}
		}

		if req.StudioID != nil {
			var studio models.Studio
			if err := tx.First(&studio, req.StudioID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					return utils.ErrStudioNotFound
				}
				return err
			}
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
			var movie models.Movie
			if err := tx.First(&movie, schedule.MovieID).Error; err != nil {
				return err
			}

			bufferMinutes := 30
			totalMinutes := int(movie.Duration) + bufferMinutes
			schedule.EndTime = schedule.StartTime.Add(time.Duration(totalMinutes) * time.Minute)
		}

		var conflictCount int64
		if err := tx.Model(&models.Schedule{}).
			Where("id != ? AND studio_id = ? AND date = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
				id, schedule.StudioID, schedule.Date, schedule.StartTime, schedule.StartTime, schedule.EndTime, schedule.EndTime, schedule.StartTime, schedule.EndTime).
			Count(&conflictCount).Error; err != nil {
			return err
		}
		if conflictCount > 0 {
			return utils.ErrScheduleConflict
		}

		if err := tx.Save(&schedule).Error; err != nil {
			return err
		}
		return nil
	})
}

func (s *ScheduleService) DeleteSchedule(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var schedule models.Schedule
		if err := tx.First(&schedule, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return utils.ErrScheduleNotFound
			}
			return err
		}

		var ticketCount int64
        if err := tx.Model(&models.Ticket{}).Where("schedule_id = ?", id).Count(&ticketCount).Error; err != nil {
            return err
        }
        if ticketCount > 0 {
            return utils.ErrScheduleHasTickets
        }

		if err := tx.Delete(&models.Schedule{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
