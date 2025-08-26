package promo

import (
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/promo/controllers"
	"movie-app-go/internal/modules/promo/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PromoModule struct {
	PromoController *controllers.PromoController
	PromoService    *services.PromoService
}

func NewPromoModule(db *gorm.DB) *PromoModule {
	promoService := services.NewPromoService(db)
	promoController := controllers.NewPromoController(promoService)

	return &PromoModule{
		PromoController: promoController,
		PromoService:    promoService,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *PromoModule) {
	rg.POST("/promos", middleware.AdminOnly(), module.PromoController.CreatePromo)
	rg.GET("/promos", middleware.Auth(), module.PromoController.GetAllPromos)
	rg.GET("/promos/:id", middleware.Auth(), module.PromoController.GetPromoByID)
	rg.PUT("/promos/:id", middleware.AdminOnly(), module.PromoController.UpdatePromo)
	rg.POST("/promos/:id/toggle", middleware.AdminOnly(), module.PromoController.TogglePromoStatus)
	rg.DELETE("/promos/:id", middleware.AdminOnly(), module.PromoController.DeletePromo)
}
