package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"movie-app-go/internal/modules/report/responses"
	"movie-app-go/internal/modules/report/services"
	"movie-app-go/internal/utils"

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
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid start_date format. Use YYYY-MM-DD"))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid end_date format. Use YYYY-MM-DD"))
		return
	}

	if endDate.Before(startDate) {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("End date must be after start date"))
		return
	}

	result, err := c.ReportService.GetDailySalesReport(startDate, endDate, page, perPage)
	if err != nil {
		if errors.Is(err, utils.ErrNoReportData) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	response := responses.PaginatedDailySalesResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      result.Data,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Daily sales report retrieved successfully",
		response,
	))
}

func (c *ReportController) GetMonthlyReportAll(ctx *gin.Context) {
	yearStr := ctx.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid year format"))
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.ReportService.GetMonthlySalesReport(year, page, perPage)
	if err != nil {
		if errors.Is(err, utils.ErrNoReportData) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	response := responses.PaginatedMonthlySalesResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      result.Data,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Monthly sales report retrieved successfully",
		response,
	))
}

func (c *ReportController) GetDailyReportByMovie(ctx *gin.Context) {
	startDateStr := ctx.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := ctx.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid start_date format. Use YYYY-MM-DD"))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid end_date format. Use YYYY-MM-DD"))
		return
	}

	if endDate.Before(startDate) {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("End date must be after start date"))
		return
	}

	result, err := c.ReportService.GetDailyReportByMovie(startDate, endDate, page, perPage)
	if err != nil {
		if errors.Is(err, utils.ErrNoReportData) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	response := responses.PaginatedDailyMovieResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      result.Data,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Daily movie report retrieved successfully",
		response,
	))
}

func (c *ReportController) GetMonthlyReportByMovie(ctx *gin.Context) {
	yearStr := ctx.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid year format"))
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.ReportService.GetMonthlyReportByMovie(year, page, perPage)
	if err != nil {
		if errors.Is(err, utils.ErrNoReportData) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	response := responses.PaginatedMonthlyMovieResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      result.Data,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Monthly movie report retrieved successfully",
		response,
	))
}

func (c *ReportController) GetDailyReportByStudio(ctx *gin.Context) {
	startDateStr := ctx.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDateStr := ctx.DefaultQuery("end_date", time.Now().Format("2006-01-02"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid start_date format. Use YYYY-MM-DD"))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid end_date format. Use YYYY-MM-DD"))
		return
	}

	if endDate.Before(startDate) {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("End date must be after start date"))
		return
	}

	result, err := c.ReportService.GetDailyReportByStudio(startDate, endDate, page, perPage)
	if err != nil {
		if errors.Is(err, utils.ErrNoReportData) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	response := responses.PaginatedDailyStudioResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      result.Data,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Daily studio report retrieved successfully",
		response,
	))
}

func (c *ReportController) GetMonthlyReportByStudio(ctx *gin.Context) {
	yearStr := ctx.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid year format"))
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.ReportService.GetMonthlyReportByStudio(year, page, perPage)
	if err != nil {
		if errors.Is(err, utils.ErrNoReportData) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	response := responses.PaginatedMonthlyStudioResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      result.Data,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Monthly studio report retrieved successfully",
		response,
	))
}
