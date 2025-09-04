package repositories

import (
    "movie-app-go/internal/models"
    "movie-app-go/internal/repository"
    "gorm.io/gorm"
)

type StudioRepository struct {
    DB *gorm.DB
}

func NewStudioRepository(db *gorm.DB) *StudioRepository {
    return &StudioRepository{DB: db}
}

func (r *StudioRepository) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.Studio], error) {
    return repository.Paginate[models.Studio](
        r.DB.Preload("FacilityStudios.Facility"),
        page,
        perPage,
    )
}

func (r *StudioRepository) GetByID(id uint) (*models.Studio, error) {
    var studio models.Studio
    if err := r.DB.Preload("FacilityStudios.Facility").First(&studio, id).Error; err != nil {
        return nil, err
    }
    return &studio, nil
}

func (r *StudioRepository) CountSchedulesByStudioID(studioID uint) (int64, error) {
    var count int64
    err := r.DB.Model(&models.Schedule{}).Where("studio_id = ?", studioID).Count(&count).Error
    return count, err
}

// Transaction methods
func (r *StudioRepository) CreateWithTx(tx *gorm.DB, studio *models.Studio) error {
    return tx.Create(studio).Error
}

func (r *StudioRepository) UpdateWithTx(tx *gorm.DB, studio *models.Studio) error {
    return tx.Save(studio).Error
}

func (r *StudioRepository) DeleteWithTx(tx *gorm.DB, id uint) error {
    return tx.Delete(&models.Studio{}, id).Error
}

func (r *StudioRepository) DeleteFacilityStudiosWithTx(tx *gorm.DB, studioID uint) error {
    return tx.Where("studio_id = ?", studioID).Delete(&models.FacilityStudio{}).Error
}

func (r *StudioRepository) CreateFacilityStudiosWithTx(tx *gorm.DB, facilityStudios []models.FacilityStudio) error {
    if len(facilityStudios) == 0 {
        return nil
    }
    return tx.Create(&facilityStudios).Error
}

func (r *StudioRepository) WithTransaction(fn func(*gorm.DB) error) error {
    return r.DB.Transaction(fn)
}