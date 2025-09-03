package repositories

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/promo/options"
	"movie-app-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type PromoRepository struct {
	DB *gorm.DB
}

func NewPromoRepository(db *gorm.DB) *PromoRepository {
	return &PromoRepository{DB: db}
}

func (r *PromoRepository) GetAllWithOptions(opts options.PromoOptions) (repository.PaginationResult[models.Promo], error) {
	query := r.DB.Model(&models.Promo{}).Preload("PromoMovies.Movie")
	
	if opts.Search != "" {
		query = query.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?",
			"%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}

	if opts.IsActive != nil {
		query = query.Where("is_active = ?", *opts.IsActive)
	}

	if opts.DiscountType != "" {
		query = query.Where("discount_type = ?", opts.DiscountType)
	}

	if opts.ValidOnly {
		now := time.Now()
		query = query.Where("valid_from <= ? AND valid_until >= ?", now, now)
	}

	if opts.MovieID != nil {
		query = query.Joins("JOIN promo_movies pm ON promos.id = pm.promo_id").
			Where("pm.movie_id = ?", *opts.MovieID)
	}

	return repository.Paginate[models.Promo](query, opts.Page, opts.PerPage)
}

func (r *PromoRepository) GetByID(id uint) (*models.Promo, error) {
	var promo models.Promo
	if err := r.DB.Preload("PromoMovies.Movie").First(&promo, id).Error; err != nil {
		return nil, err
	}
	return &promo, nil
}

func (r *PromoRepository) GetByCode(code string) (*models.Promo, error) {
	var promo models.Promo
	if err := r.DB.Preload("PromoMovies").Where("code = ?", code).First(&promo).Error; err != nil {
		return nil, err
	}
	return &promo, nil
}

func (r *PromoRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	err := r.DB.Model(&models.Promo{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}

func (r *PromoRepository) CountMoviesByIDs(movieIDs []uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Movie{}).Where("id IN ?", movieIDs).Count(&count).Error
	return count, err
}

func (r *PromoRepository) CountTransactionsByPromoID(promoID uint) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Transaction{}).Where("promo_id = ?", promoID).Count(&count).Error
	return count, err
}

// Transaction methods
func (r *PromoRepository) CreateWithTx(tx *gorm.DB, promo *models.Promo) error {
	return tx.Create(promo).Error
}

func (r *PromoRepository) UpdateWithTx(tx *gorm.DB, promo *models.Promo) error {
	return tx.Save(promo).Error
}

func (r *PromoRepository) DeleteWithTx(tx *gorm.DB, id uint) error {
	return tx.Delete(&models.Promo{}, id).Error
}

func (r *PromoRepository) CreatePromoMoviesWithTx(tx *gorm.DB, promoMovies []models.PromoMovie) error {
	if len(promoMovies) == 0 {
		return nil
	}
	return tx.Create(&promoMovies).Error
}

func (r *PromoRepository) DeletePromoMoviesWithTx(tx *gorm.DB, promoID uint) error {
	return tx.Where("promo_id = ?", promoID).Delete(&models.PromoMovie{}).Error
}

func (r *PromoRepository) WithTransaction(fn func(*gorm.DB) error) error {
	return r.DB.Transaction(fn)
}
