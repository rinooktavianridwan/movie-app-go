package services

import (
	"fmt"
	"movie-app-go/internal/modules/schedule/options"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/schedule/requests"
	"movie-app-go/internal/repository"
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
		return nil, err
	}
	return &schedule, nil
}

func (s *ScheduleService) CreateSchedule(req *requests.CreateScheduleRequest) (*models.Schedule, error) {
	var schedule models.Schedule
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var movie models.Movie
		if err := tx.First(&movie, req.MovieID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("movie not found")
			}
			return err
		}

		var studio models.Studio
		if err := tx.First(&studio, req.StudioID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("studio not found")
			}
			return err
		}

		bufferMinutes := 30
		totalMinutes := int(movie.Duration) + bufferMinutes
		endTime := req.StartTime.Add(time.Duration(totalMinutes) * time.Minute)

		if req.Date.Before(time.Now().Truncate(24 * time.Hour)) {
			return fmt.Errorf("date cannot be in the past")
		}

		var conflictCount int64
		if err := tx.Model(&models.Schedule{}).
			Where("studio_id = ? AND date = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
				req.StudioID, req.Date, req.StartTime, req.StartTime, endTime, endTime, req.StartTime, endTime).
			Count(&conflictCount).Error; err != nil {
			return err
		}
		if conflictCount > 0 {
			return fmt.Errorf("schedule conflict: studio is already booked at this time")
		}

		schedule = models.Schedule{
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

	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (s *ScheduleService) UpdateSchedule(id uint, req *requests.UpdateScheduleRequest) (*models.Schedule, error) {
	var schedule models.Schedule
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&schedule, id).Error; err != nil {
			return err
		}

		var movie models.Movie
		if err := tx.First(&movie, req.MovieID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("movie not found")
			}
			return err
		}

		var studio models.Studio
		if err := tx.First(&studio, req.StudioID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("studio not found")
			}
			return err
		}

		bufferMinutes := 30
		totalMinutes := int(movie.Duration) + bufferMinutes
		endTime := req.StartTime.Add(time.Duration(totalMinutes) * time.Minute)

		if req.Date.Before(time.Now().Truncate(24 * time.Hour)) {
			return fmt.Errorf("date cannot be in the past")
		}

		var conflictCount int64
		if err := tx.Model(&models.Schedule{}).
			Where("id != ? AND studio_id = ? AND date = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
				id, req.StudioID, req.Date, req.StartTime, req.StartTime, endTime, endTime, req.StartTime, endTime).
			Count(&conflictCount).Error; err != nil {
			return err
		}
		if conflictCount > 0 {
			return fmt.Errorf("schedule conflict: studio is already booked at this time")
		}

		schedule.MovieID = req.MovieID
		schedule.StudioID = req.StudioID
		schedule.StartTime = req.StartTime
		schedule.EndTime = endTime
		schedule.Date = req.Date
		schedule.Price = req.Price

		if err := tx.Save(&schedule).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (s *ScheduleService) DeleteSchedule(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		var schedule models.Schedule
		if err := tx.First(&schedule, id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("schedule not found")
			}
			return err
		}

		var ticketCount int64
		if err := tx.Model(&models.Ticket{}).Where("schedule_id = ?", id).Count(&ticketCount).Error; err != nil {
			return err
		}
		if ticketCount > 0 {
			return fmt.Errorf("cannot delete schedule: tickets already exist")
		}

		if err := tx.Delete(&models.Schedule{}, id).Error; err != nil {
			return err
		}
		return nil
	})
}
