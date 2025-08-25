package report

import (
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/report/controllers"
	"movie-app-go/internal/modules/report/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportModule struct {
	ReportController *controllers.ReportController
}

func NewReportModule(db *gorm.DB) *ReportModule {
	reportService := services.NewReportService(db)

	return &ReportModule{
		ReportController: controllers.NewReportController(reportService),
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *ReportModule) {
	rg.GET("/reports/daily/all", middleware.AdminOnly(), module.ReportController.GetDailyReportAll)
	rg.GET("/reports/daily/movie", middleware.AdminOnly(), module.ReportController.GetDailyReportByMovie)
	rg.GET("/reports/daily/studio", middleware.AdminOnly(), module.ReportController.GetDailyReportByStudio)
	rg.GET("/reports/monthly/all", middleware.AdminOnly(), module.ReportController.GetMonthlyReportAll)
	rg.GET("/reports/monthly/movie", middleware.AdminOnly(), module.ReportController.GetMonthlyReportByMovie)
	rg.GET("/reports/monthly/studio", middleware.AdminOnly(), module.ReportController.GetMonthlyReportByStudio)
}
