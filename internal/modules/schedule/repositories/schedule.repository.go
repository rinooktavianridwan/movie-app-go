package repositories

import (
    "movie-app-go/internal/models"
    "movie-app-go/internal/modules/schedule/options"
    "movie-app-go/internal/repository"
    "time"
    "gorm.io/gorm"
)

type ScheduleRepository struct {
    DB *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
    return &ScheduleRepository{DB: db}
}

func (r *ScheduleRepository) GetAllWithOptions(opts *options.GetAllScheduleOptions) (repository.PaginationResult[models.Schedule], error) {
    query := r.DB.Preload("Movie").Preload("Studio")

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

func (r *ScheduleRepository) GetByID(id uint) (*models.Schedule, error) {
    var schedule models.Schedule
    if err := r.DB.Preload("Movie").Preload("Studio").First(&schedule, id).Error; err != nil {
        return nil, err
    }
    return &schedule, nil
}

func (r *ScheduleRepository) CheckScheduleConflict(studioID uint, startTime, endTime time.Time, excludeID *uint) (bool, error) {
    query := r.DB.Model(&models.Schedule{}).
        Where("studio_id = ? AND ((start_time <= ? AND end_time > ?) OR (start_time < ? AND end_time >= ?) OR (start_time >= ? AND end_time <= ?))",
            studioID, startTime, endTime, startTime, endTime, startTime, endTime)

    if excludeID != nil {
        query = query.Where("id != ?", *excludeID)
    }

    var count int64
    err := query.Count(&count).Error
    return count > 0, err
}

func (r *ScheduleRepository) CountTicketsByScheduleID(scheduleID uint) (int64, error) {
    var count int64
    err := r.DB.Model(&models.Ticket{}).Where("schedule_id = ?", scheduleID).Count(&count).Error
    return count, err
}

// Transaction methods
func (r *ScheduleRepository) CreateWithTx(tx *gorm.DB, schedule *models.Schedule) error {
    return tx.Create(schedule).Error
}

func (r *ScheduleRepository) UpdateWithTx(tx *gorm.DB, schedule *models.Schedule) error {
    return tx.Save(schedule).Error
}

func (r *ScheduleRepository) DeleteWithTx(tx *gorm.DB, id uint) error {
    return tx.Delete(&models.Schedule{}, id).Error
}

func (r *ScheduleRepository) GetByIDWithTx(tx *gorm.DB, id uint) (*models.Schedule, error) {
    var schedule models.Schedule
    if err := tx.Preload("Movie").Preload("Studio").First(&schedule, id).Error; err != nil {
        return nil, err
    }
    return &schedule, nil
}

func (r *ScheduleRepository) WithTransaction(fn func(*gorm.DB) error) error {
    return r.DB.Transaction(fn)
}