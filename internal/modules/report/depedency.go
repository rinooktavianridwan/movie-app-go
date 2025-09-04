package report

import (
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/report/controllers"
	"movie-app-go/internal/modules/report/repositories"
	"movie-app-go/internal/modules/report/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportModule struct {
	ReportController *controllers.ReportController
}

func NewReportModule(db *gorm.DB) *ReportModule {
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)

	return &ReportModule{
		ReportController: controllers.NewReportController(reportService),
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *ReportModule, mf *middleware.Factory) {
    reports := rg.Group("/reports")
    reports.Use(mf.Auth(), mf.RequirePermission("reports.view"))
    {
        // Daily reports
        reports.GET("/daily", module.ReportController.GetDailyReportAll)
        reports.GET("/daily/movies", module.ReportController.GetDailyReportByMovie)
        reports.GET("/daily/studios", module.ReportController.GetDailyReportByStudio)
        
        // Monthly reports
        reports.GET("/monthly", module.ReportController.GetMonthlyReportAll)
        reports.GET("/monthly/movies", module.ReportController.GetMonthlyReportByMovie)
        reports.GET("/monthly/studios", module.ReportController.GetMonthlyReportByStudio)
    }
}
