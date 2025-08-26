package controllers

import (
	"net/http"
	"strconv"
	"time"

	"movie-app-go/internal/modules/report/services"

	"github.com/gin-gonic/gin"
)

type ReportController struct {
	ReportService *services.ReportService
}

func NewReportController(reportService *services.ReportService) *ReportController {
	return &ReportController{ReportService: reportService}
}

func (c *ReportController) GetDailyReportAll(ctx *gin.Context) {
    startDateStr := ctx.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
    endDateStr := ctx.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

    startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
        return
    }

    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
        return
    }

    result, err := c.ReportService.GetDailySalesReport(startDate, endDate, page, perPage)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, result)
}

func (c *ReportController) GetMonthlyReportAll(ctx *gin.Context) {
    yearStr := ctx.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
    year, err := strconv.Atoi(yearStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
        return
    }

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

    result, err := c.ReportService.GetMonthlySalesReport(year, page, perPage)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, result)
}

func (c *ReportController) GetDailyReportByMovie(ctx *gin.Context) {
	startDateStr := ctx.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := ctx.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
        return
    }

    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
        return
    }

    result, err := c.ReportService.GetDailyReportByMovie(startDate, endDate, page, perPage)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, result)
}

func (c *ReportController) GetMonthlyReportByMovie(ctx *gin.Context) {
	yearStr := ctx.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
    year, err := strconv.Atoi(yearStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
        return
    }

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

    result, err := c.ReportService.GetMonthlyReportByMovie(year, page, perPage)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, result)
}

func (c *ReportController) GetDailyReportByStudio(ctx *gin.Context) {
	startDateStr := ctx.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := ctx.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
        return
    }

    endDate, err := time.Parse("2006-01-02", endDateStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
        return
    }

    result, err := c.ReportService.GetDailyReportByStudio(startDate, endDate, page, perPage)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, result)
}

func (c *ReportController) GetMonthlyReportByStudio(ctx *gin.Context) {
	yearStr := ctx.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
    year, err := strconv.Atoi(yearStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
        return
    }

    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

    result, err := c.ReportService.GetMonthlyReportByStudio(year, page, perPage)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, result)
}
