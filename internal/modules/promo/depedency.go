package promo

import (
	"movie-app-go/internal/middleware"
	movierepos "movie-app-go/internal/modules/movie/repositories"
	notificationRepositories "movie-app-go/internal/modules/notification/repositories"
	notificationServices "movie-app-go/internal/modules/notification/services"
	"movie-app-go/internal/modules/promo/controllers"
	"movie-app-go/internal/modules/promo/repositories"
	"movie-app-go/internal/modules/promo/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PromoModule struct {
	PromoController *controllers.PromoController
	PromoService    *services.PromoService
}

func NewPromoModule(db *gorm.DB) *PromoModule {
	promoRepo := repositories.NewPromoRepository(db)
	movieRepo := movierepos.NewMovieRepository(db)

	notificationRepo := notificationRepositories.NewNotificationRepository(db)
	notificationService := notificationServices.NewNotificationService(notificationRepo)

	promoService := services.NewPromoService(promoRepo, movieRepo, notificationService)
	promoController := controllers.NewPromoController(promoService)

	return &PromoModule{
		PromoController: promoController,
		PromoService:    promoService,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *PromoModule) {
	rg.POST("/promos", middleware.AdminOnly(), module.PromoController.Create)
	rg.GET("/promos", middleware.Auth(), module.PromoController.GetAllPromos)
	rg.GET("/promos/:id", middleware.Auth(), module.PromoController.GetPromoByID)
	rg.PUT("/promos/:id", middleware.AdminOnly(), module.PromoController.Update)
	rg.POST("/promos/:id/toggle", middleware.AdminOnly(), module.PromoController.TogglePromoStatus)
	rg.DELETE("/promos/:id", middleware.AdminOnly(), module.PromoController.Delete)
}
