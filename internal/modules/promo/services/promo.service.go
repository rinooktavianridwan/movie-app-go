package services

import (
	"fmt"
	"log"
	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"
	movierepos "movie-app-go/internal/modules/movie/repositories"
	notificationServices "movie-app-go/internal/modules/notification/services"
	"movie-app-go/internal/modules/promo/options"
	"movie-app-go/internal/modules/promo/repositories"
	"movie-app-go/internal/modules/promo/requests"
	"movie-app-go/internal/modules/promo/responses"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"
	"time"

	"gorm.io/gorm"
)

type PromoService struct {
	PromoRepo           *repositories.PromoRepository
	MovieRepo           *movierepos.MovieRepository
	NotificationService *notificationServices.NotificationService
}

func NewPromoService(
	promoRepo *repositories.PromoRepository,
	movieRepo *movierepos.MovieRepository,
	notificationService *notificationServices.NotificationService,
) *PromoService {
	return &PromoService{
		PromoRepo:           promoRepo,
		MovieRepo:           movieRepo,
		NotificationService: notificationService,
	}
}

func (s *PromoService) CreatePromo(req *requests.CreatePromoRequest) error {
	exists, err := s.PromoRepo.ExistsByCode(req.Code)
	if err != nil {
		return err
	}
	if exists {
		return utils.ErrPromoCodeExists
	}

	if len(req.MovieIDs) > 0 {
		count, err := s.PromoRepo.CountMoviesByIDs(req.MovieIDs)
		if err != nil {
			return err
		}
		if count != int64(len(req.MovieIDs)) {
			return utils.ErrInvalidMovieIDs
		}
	}

	return s.PromoRepo.WithTransaction(func(tx *gorm.DB) error {
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

		if err := s.PromoRepo.CreateWithTx(tx, &promo); err != nil {
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
			if err := s.PromoRepo.CreatePromoMoviesWithTx(tx, promoMovies); err != nil {
				return err
			}
		}

		if promo.IsActive {
			s.sendPromoNotificationAsync(promo.ID, promo.Name, promo.Code)
		}

		return nil
	})
}

func (s *PromoService) GetAllPromosPaginated(opts options.PromoOptions) (repository.PaginationResult[models.Promo], error) {
	return s.PromoRepo.GetAllWithOptions(opts)
}

func (s *PromoService) GetPromoByID(id uint) (*models.Promo, error) {
	promo, err := s.PromoRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrPromoNotFound
		}
		return nil, err
	}
	return promo, nil
}

func (s *PromoService) UpdatePromo(id uint, req *requests.UpdatePromoRequest) error {
	promo, err := s.GetPromoByID(id)
	if err != nil {
		return err
	}

	var wasInactive bool = !promo.IsActive

	if req.MovieIDs != nil && len(req.MovieIDs) > 0 {
		count, err := s.PromoRepo.CountMoviesByIDs(req.MovieIDs)
		if err != nil {
			return err
		}
		if count != int64(len(req.MovieIDs)) {
			return utils.ErrInvalidMovieIDs
		}
	}

	return s.PromoRepo.WithTransaction(func(tx *gorm.DB) error {
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

		if err := s.PromoRepo.UpdateWithTx(tx, promo); err != nil {
			return err
		}

		if req.MovieIDs != nil {
			if err := s.PromoRepo.DeletePromoMoviesWithTx(tx, id); err != nil {
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
				if err := s.PromoRepo.CreatePromoMoviesWithTx(tx, promoMovies); err != nil {
					return err
				}
			}
		}

		if wasInactive && promo.IsActive {
			s.sendPromoNotificationAsync(promo.ID, promo.Name, promo.Code)
		}

		return nil
	})
}

func (s *PromoService) TogglePromoStatus(id uint) error {
	promo, err := s.GetPromoByID(id)
	if err != nil {
		return err
	}

	wasInactive := !promo.IsActive
	promo.IsActive = !promo.IsActive

	err = s.PromoRepo.WithTransaction(func(tx *gorm.DB) error {
		return s.PromoRepo.UpdateWithTx(tx, promo)
	})

	if err != nil {
		return err
	}

	if wasInactive && promo.IsActive {
		s.sendPromoNotificationAsync(promo.ID, promo.Name, promo.Code)
	}

	return nil
}

func (s *PromoService) DeletePromo(id uint) error {
	_, err := s.GetPromoByID(id)
	if err != nil {
		return err
	}
	usageCount, err := s.PromoRepo.CountTransactionsByPromoID(id)
	if err != nil {
		return err
	}
	if usageCount > 0 {
		return utils.ErrPromoInUse
	}

	return s.PromoRepo.WithTransaction(func(tx *gorm.DB) error {
		return s.PromoRepo.DeleteWithTx(tx, id)
	})
}

func (s *PromoService) ValidatePromo(userID uint, req *requests.ValidatePromoRequest) (*responses.PromoValidationResponse, error) {
	promo, err := s.PromoRepo.GetByCode(req.PromoCode)
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
		movieIDMap := make(map[uint]bool)
		for _, id := range req.MovieIDs {
			movieIDMap[id] = true
		}

		isValidForMovie := false
		for _, pm := range promo.PromoMovies {
			if movieIDMap[pm.MovieID] {
				isValidForMovie = true
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

func (s *PromoService) GetPromoByCode(code string) (*models.Promo, error) {
	promo, err := s.PromoRepo.GetByCode(code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrPromoNotFound
		}
		return nil, err
	}
	return promo, nil
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
