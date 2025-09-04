package studio

import (
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/studio/controllers"
	"movie-app-go/internal/modules/studio/repositories"
	"movie-app-go/internal/modules/studio/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StudioModule struct {
	StudioController   *controllers.StudioController
	FacilityController *controllers.FacilityController
}

func NewStudioModule(db *gorm.DB) *StudioModule {
	studioRepo := repositories.NewStudioRepository(db)
	facilityRepo := repositories.NewFacilityRepository(db)
	studioService := services.NewStudioService(studioRepo, facilityRepo)
	facilityService := services.NewFacilityService(facilityRepo)

	return &StudioModule{
		StudioController:   controllers.NewStudioController(studioService),
		FacilityController: controllers.NewFacilityController(facilityService),
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *StudioModule) {
	// Facility
	rg.POST("/facilities", middleware.AdminOnly(), module.FacilityController.Create)
	rg.GET("/facilities", middleware.AdminOnly(), module.FacilityController.GetAll)
	rg.GET("/facilities/:id", middleware.AdminOnly(), module.FacilityController.GetByID)
	rg.PUT("/facilities/:id", middleware.AdminOnly(), module.FacilityController.Update)
	rg.DELETE("/facilities/:id", middleware.AdminOnly(), module.FacilityController.Delete)

	// Studio
	rg.POST("/studios", middleware.AdminOnly(), module.StudioController.Create)
	rg.GET("/studios", middleware.AdminOnly(), module.StudioController.GetAll)
	rg.GET("/studios/:id", middleware.AdminOnly(), module.StudioController.GetByID)
	rg.PUT("/studios/:id", middleware.AdminOnly(), module.StudioController.Update)
	rg.DELETE("/studios/:id", middleware.AdminOnly(), module.StudioController.Delete)
}
