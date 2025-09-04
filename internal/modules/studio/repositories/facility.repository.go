package repositories

import (
    "movie-app-go/internal/models"
    "movie-app-go/internal/repository"
    "gorm.io/gorm"
)

type FacilityRepository struct {
    DB *gorm.DB
}

func NewFacilityRepository(db *gorm.DB) *FacilityRepository {
    return &FacilityRepository{DB: db}
}

func (r *FacilityRepository) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.Facility], error) {
    return repository.Paginate[models.Facility](r.DB, page, perPage)
}

func (r *FacilityRepository) GetByID(id uint) (*models.Facility, error) {
    var facility models.Facility
    if err := r.DB.First(&facility, id).Error; err != nil {
        return nil, err
    }
    return &facility, nil
}

func (r *FacilityRepository) Create(facility *models.Facility) error {
    return r.DB.Create(facility).Error
}

func (r *FacilityRepository) Update(facility *models.Facility) error {
    return r.DB.Save(facility).Error
}

func (r *FacilityRepository) Delete(id uint) error {
    return r.DB.Delete(&models.Facility{}, id).Error
}

func (r *FacilityRepository) ExistsByName(name string) (bool, error) {
    var count int64
    err := r.DB.Model(&models.Facility{}).Where("name = ?", name).Count(&count).Error
    return count > 0, err
}

func (r *FacilityRepository) ExistsByNameExceptID(name string, facilityID uint) (bool, error) {
    var count int64
    err := r.DB.Model(&models.Facility{}).Where("name = ? AND id != ?", name, facilityID).Count(&count).Error
    return count > 0, err
}

func (r *FacilityRepository) CountFacilitiesByIDs(facilityIDs []uint) (int64, error) {
    var count int64
    err := r.DB.Model(&models.Facility{}).Where("id IN ?", facilityIDs).Count(&count).Error
    return count, err
}

func (r *FacilityRepository) CountStudiosByFacilityID(facilityID uint) (int64, error) {
    var count int64
    err := r.DB.Model(&models.FacilityStudio{}).Where("facility_id = ?", facilityID).Count(&count).Error
    return count, err
}
