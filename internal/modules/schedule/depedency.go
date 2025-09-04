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

func RegisterRoutes(rg *gin.RouterGroup, module *ScheduleModule) {
	// Schedule routes
	rg.POST("/schedules", middleware.AdminOnly(), module.ScheduleController.Create)
	rg.GET("/schedules", module.ScheduleController.GetAll)      // Public untuk melihat jadwal
	rg.GET("/schedules/:id", module.ScheduleController.GetByID) // Public untuk melihat detail jadwal
	rg.PUT("/schedules/:id", middleware.AdminOnly(), module.ScheduleController.Update)
	rg.DELETE("/schedules/:id", middleware.AdminOnly(), module.ScheduleController.Delete)
}
