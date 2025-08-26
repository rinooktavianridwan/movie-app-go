package services

import (
	"fmt"
	"log"
	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"
	notificationServices "movie-app-go/internal/modules/notification/services"
	"movie-app-go/internal/modules/promo/options"
	"movie-app-go/internal/modules/promo/requests"
	"movie-app-go/internal/modules/promo/responses"
	"movie-app-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type PromoService struct {
	DB                  *gorm.DB
	NotificationService *notificationServices.NotificationService
}

func NewPromoService(db *gorm.DB) *PromoService {
	return &PromoService{DB: db, NotificationService: notificationServices.NewNotificationService(db)}
}

func (s *PromoService) CreatePromo(req *requests.CreatePromoRequest) (*models.Promo, error) {
	var existingPromo models.Promo
	if err := s.DB.Where("code = ?", req.Code).First(&existingPromo).Error; err == nil {
		return nil, fmt.Errorf("promo code already exists")
	}

	var createdPromo *models.Promo
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		if len(req.MovieIDs) > 0 {
			var count int64
			if err := tx.Model(&models.Movie{}).Where("id IN ?", req.MovieIDs).Count(&count).Error; err != nil {
				return err
			}
			if count != int64(len(req.MovieIDs)) {
				return fmt.Errorf("some movie_ids are invalid")
			}
		}

		promo := models.Promo{
			Name:          req.Name,
			Code:          req.Code,
			Description:   req.Description,
			DiscountType:  req.DiscountType,
			DiscountValue: req.DiscountValue,
			MinTickets:    req.MinTickets,
			MaxDiscount:   req.MaxDiscount,
			UsageLimit:    req.UsageLimit,
			IsActive:      req.IsActive,
			ValidFrom:     req.ValidFrom,
			ValidUntil:    req.ValidUntil,
		}

		if err := tx.Create(&promo).Error; err != nil {
			return err
		}

		if len(req.MovieIDs) > 0 {
			promoMovies := make([]models.PromoMovie, len(req.MovieIDs))
			for i, movieID := range req.MovieIDs {
				promoMovies[i] = models.PromoMovie{
					PromoID: promo.ID,
					MovieID: movieID,
				}
			}
			if err := tx.Create(&promoMovies).Error; err != nil {
				return err
			}
		}

		createdPromo = &promo
		return nil
	})

	if err != nil {
		return nil, err
	}

	if createdPromo.IsActive {
		s.sendPromoNotificationAsync(createdPromo.ID, createdPromo.Name, createdPromo.Code)
	}

	return createdPromo, nil
}

func (s *PromoService) GetAllPromosPaginated(opts options.PromoOptions) (repository.PaginationResult[models.Promo], error) {
	query := s.DB.Model(&models.Promo{}).
		Preload("PromoMovies.Movie")

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

func (s *PromoService) GetPromoByID(id uint) (*models.Promo, error) {
	var promo models.Promo
	if err := s.DB.Preload("PromoMovies.Movie").First(&promo, id).Error; err != nil {
		return nil, fmt.Errorf("promo not found")
	}
	return &promo, nil
}

func (s *PromoService) GetPromoByCode(code string) (*models.Promo, error) {
	var promo models.Promo
	if err := s.DB.Preload("PromoMovies").Where("code = ?", code).First(&promo).Error; err != nil {
		return nil, fmt.Errorf("promo not found")
	}
	return &promo, nil
}

func (s *PromoService) UpdatePromo(id uint, req *requests.UpdatePromoRequest) (*models.Promo, error) {
	var updatedPromo *models.Promo
	var wasInactive bool

	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var promo models.Promo
		if err := tx.First(&promo, id).Error; err != nil {
			return fmt.Errorf("promo not found")
		}

		wasInactive = !promo.IsActive

		if len(req.MovieIDs) > 0 {
			var count int64
			if err := tx.Model(&models.Movie{}).Where("id IN ?", req.MovieIDs).Count(&count).Error; err != nil {
				return err
			}
			if count != int64(len(req.MovieIDs)) {
				return fmt.Errorf("some movie_ids are invalid")
			}
		}

		if req.Name != nil {
			promo.Name = *req.Name
		}
		if req.Description != nil {
			promo.Description = *req.Description
		}
		if req.DiscountType != nil {
			promo.DiscountType = *req.DiscountType
		}
		if req.DiscountValue != nil {
			promo.DiscountValue = *req.DiscountValue
		}
		if req.MinTickets != nil {
			promo.MinTickets = *req.MinTickets
		}
		if req.MaxDiscount != nil {
			promo.MaxDiscount = req.MaxDiscount
		}
		if req.UsageLimit != nil {
			promo.UsageLimit = req.UsageLimit
		}
		if req.IsActive != nil {
			promo.IsActive = *req.IsActive
		}
		if req.ValidFrom != nil {
			promo.ValidFrom = *req.ValidFrom
		}
		if req.ValidUntil != nil {
			promo.ValidUntil = *req.ValidUntil
		}

		if err := tx.Save(&promo).Error; err != nil {
			return err
		}

		if req.MovieIDs != nil {
			if err := tx.Where("promo_id = ?", id).Delete(&models.PromoMovie{}).Error; err != nil {
				return err
			}

			if len(req.MovieIDs) > 0 {
				promoMovies := make([]models.PromoMovie, len(req.MovieIDs))
				for i, movieID := range req.MovieIDs {
					promoMovies[i] = models.PromoMovie{
						PromoID: promo.ID,
						MovieID: movieID,
					}
				}
				if err := tx.Create(&promoMovies).Error; err != nil {
					return err
				}
			}
		}

		updatedPromo = &promo
		return nil
	})

	if err != nil {
		return nil, err
	}

	if wasInactive && updatedPromo.IsActive {
		s.sendPromoNotificationAsync(updatedPromo.ID, updatedPromo.Name, updatedPromo.Code)
	}

	return updatedPromo, nil
}

func (s *PromoService) TogglePromoStatus(id uint) (*models.Promo, error) {
	var promo models.Promo
	if err := s.DB.First(&promo, id).Error; err != nil {
		return nil, fmt.Errorf("promo not found")
	}

	wasInactive := !promo.IsActive
	promo.IsActive = !promo.IsActive
	if err := s.DB.Save(&promo).Error; err != nil {
		return nil, err
	}

	if wasInactive && promo.IsActive {
		s.sendPromoNotificationAsync(promo.ID, promo.Name, promo.Code)
	}

	return &promo, nil
}

func (s *PromoService) DeletePromo(id uint) error {
	var promo models.Promo
	if err := s.DB.First(&promo, id).Error; err != nil {
		return fmt.Errorf("promo not found")
	}

	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("promo_id = ?", id).Delete(&models.PromoMovie{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&promo).Error; err != nil {
			return err
		}

		return nil
	})
}

func (s *PromoService) ValidatePromo(userID uint, req *requests.ValidatePromoRequest) (*responses.PromoValidationResponse, error) {
	promo, err := s.GetPromoByCode(req.PromoCode)
	if err != nil {
		return &responses.PromoValidationResponse{
			IsValid: false,
			Message: "Promo code not found",
		}, nil
	}

	if !promo.IsActive {
		return &responses.PromoValidationResponse{
			IsValid: false,
			Message: "Promo is not active",
		}, nil
	}

	now := time.Now()
	if now.Before(promo.ValidFrom) || now.After(promo.ValidUntil) {
		return &responses.PromoValidationResponse{
			IsValid: false,
			Message: "Promo has expired or not yet active",
		}, nil
	}

	if len(req.SeatNumbers) < promo.MinTickets {
		return &responses.PromoValidationResponse{
			IsValid: false,
			Message: fmt.Sprintf("Minimum %d tickets required", promo.MinTickets),
		}, nil
	}

	if promo.UsageLimit != nil && promo.UsageCount >= *promo.UsageLimit {
		return &responses.PromoValidationResponse{
			IsValid: false,
			Message: "Promo usage limit exceeded",
		}, nil
	}

	if len(promo.PromoMovies) > 0 {
		isValidForMovie := false
		for _, pm := range promo.PromoMovies {
			for _, movieID := range req.MovieIDs {
				if pm.MovieID == movieID {
					isValidForMovie = true
					break
				}
			}
			if isValidForMovie {
				break
			}
		}
		if !isValidForMovie {
			return &responses.PromoValidationResponse{
				IsValid: false,
				Message: "Promo is not applicable for selected movies",
			}, nil
		}
	}

	var discountAmount float64
	if promo.DiscountType == constants.DiscountTypePercentage {
		discountAmount = req.TotalAmount * (promo.DiscountValue / 100)
		if promo.MaxDiscount != nil && discountAmount > *promo.MaxDiscount {
			discountAmount = *promo.MaxDiscount
		}
	} else {
		discountAmount = promo.DiscountValue
	}

	finalAmount := req.TotalAmount - discountAmount

	return &responses.PromoValidationResponse{
		IsValid:        true,
		DiscountAmount: discountAmount,
		FinalAmount:    finalAmount,
		Message:        "Promo applied successfully",
	}, nil
}

func (s *PromoService) sendPromoNotificationAsync(promoID uint, promoName, promoCode string) {
	go func() {
		if err := s.NotificationService.CreatePromoNotification(
			promoID,
			promoName,
			promoCode,
		); err != nil {
			log.Printf("Failed to create promo notifications: %v", err)
		} else {
			log.Printf("Promo notification sent for promo: %s", promoName)
		}
	}()
}
