package schedule

import (
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/schedule/controllers"
	"movie-app-go/internal/modules/schedule/repositories"
	movierepos "movie-app-go/internal/modules/movie/repositories"
	studiorepos "movie-app-go/internal/modules/studio/repositories"
	"movie-app-go/internal/modules/schedule/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScheduleModule struct {
	ScheduleController *controllers.ScheduleController
}

func NewScheduleModule(db *gorm.DB) *ScheduleModule {
	studioRepo := studiorepos.NewStudioRepository(db)
	scheduleRepo := repositories.NewScheduleRepository(db)
    movieRepo := movierepos.NewMovieRepository(db)
	scheduleService := services.NewScheduleService(scheduleRepo, movieRepo, studioRepo)

	return &ScheduleModule{
		ScheduleController: controllers.NewScheduleController(scheduleService),
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *ScheduleModule, mf *middleware.Factory) {
    rg.POST("/schedules", mf.Auth(), mf.RequirePermission("schedules.create"), module.ScheduleController.Create)
    rg.GET("/schedules", module.ScheduleController.GetAll)
    rg.GET("/schedules/:id", module.ScheduleController.GetByID)
    rg.PUT("/schedules/:id", mf.Auth(), mf.RequirePermission("schedules.update"), module.ScheduleController.Update)
    rg.DELETE("/schedules/:id", mf.Auth(), mf.RequirePermission("schedules.delete"), module.ScheduleController.Delete)
}
