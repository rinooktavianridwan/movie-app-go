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

func RegisterRoutes(rg *gin.RouterGroup, module *StudioModule, mf *middleware.Factory) {
    // Facility routes
    rg.POST("/facilities", mf.Auth(), mf.RequirePermission("studios.create"), module.FacilityController.Create)
    rg.GET("/facilities", mf.Auth(), mf.RequirePermission("studios.read"), module.FacilityController.GetAll)
    rg.GET("/facilities/:id", mf.Auth(), mf.RequirePermission("studios.read"), module.FacilityController.GetByID)
    rg.PUT("/facilities/:id", mf.Auth(), mf.RequirePermission("studios.update"), module.FacilityController.Update)
    rg.DELETE("/facilities/:id", mf.Auth(), mf.RequirePermission("studios.delete"), module.FacilityController.Delete)

    // Studio routes
    rg.POST("/studios", mf.Auth(), mf.RequirePermission("studios.create"), module.StudioController.Create)
    rg.GET("/studios", module.StudioController.GetAll)
    rg.GET("/studios/:id", module.StudioController.GetByID)
    rg.PUT("/studios/:id", mf.Auth(), mf.RequirePermission("studios.update"), module.StudioController.Update)
    rg.DELETE("/studios/:id", mf.Auth(), mf.RequirePermission("studios.delete"), module.StudioController.Delete)
}
