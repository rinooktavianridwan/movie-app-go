package depedency

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "movie-app-go/internal/modules/studio/controllers"
    "movie-app-go/internal/modules/studio/services"
)

type StudioModule struct {
    StudioController   *controllers.StudioController
    FacilityController *controllers.FacilityController
}

func InitStudioModule(db *gorm.DB) *StudioModule {
    studioService := services.NewStudioService(db)
    facilityService := services.NewFacilityService(db)

    return &StudioModule{
        StudioController:   controllers.NewStudioController(studioService),
        FacilityController: controllers.NewFacilityController(facilityService),
    }
}

func RegisterStudioRoutes(rg *gin.RouterGroup, module *StudioModule) {
    // Facility
    rg.POST("/facilities", module.FacilityController.Create)
    rg.GET("/facilities", module.FacilityController.GetAll)
    rg.GET("/facilities/:id", module.FacilityController.GetByID)
    rg.PUT("/facilities/:id", module.FacilityController.Update)
    rg.DELETE("/facilities/:id", module.FacilityController.Delete)

    // Studio
    rg.POST("/studios", module.StudioController.Create)
    rg.GET("/studios", module.StudioController.GetAll)
    rg.GET("/studios/:id", module.StudioController.GetByID)
    rg.PUT("/studios/:id", module.StudioController.Update)
    rg.DELETE("/studios/:id", module.StudioController.Delete)
}