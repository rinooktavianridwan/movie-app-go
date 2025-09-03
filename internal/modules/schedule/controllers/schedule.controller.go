package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"movie-app-go/internal/modules/schedule/options"
	"movie-app-go/internal/modules/schedule/requests"
	"movie-app-go/internal/modules/schedule/responses"
	"movie-app-go/internal/modules/schedule/services"
	"movie-app-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type ScheduleController struct {
	ScheduleService *services.ScheduleService
}

func NewScheduleController(s *services.ScheduleService) *ScheduleController {
	return &ScheduleController{ScheduleService: s}
}

func (c *ScheduleController) Create(ctx *gin.Context) {
	var req requests.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	err := c.ScheduleService.CreateSchedule(&req)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidMovieIDs):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		case errors.Is(err, utils.ErrStudioNotFound):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		case errors.Is(err, utils.ErrScheduleConflict):
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, err.Error()))
		case errors.Is(err, utils.ErrInvalidTimeRange):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		case errors.Is(err, utils.ErrPastDate):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(
		http.StatusCreated,
		"Schedule created successfully",
		nil,
	))
}

func (c *ScheduleController) GetAll(ctx *gin.Context) {
	opts, err := options.ParseScheduleOptions(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	result, err := c.ScheduleService.GetAllSchedulesPaginated(opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	scheduleResponses := responses.ToScheduleResponses(result.Data)
	response := responses.PaginatedScheduleResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      scheduleResponses,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Schedules retrieved successfully",
		response,
	))
}

func (c *ScheduleController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	schedule, err := c.ScheduleService.GetScheduleByID(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrScheduleNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Schedule retrieved successfully",
		responses.ToScheduleResponse(schedule),
	))
}

func (c *ScheduleController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req requests.UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	err := c.ScheduleService.UpdateSchedule(uint(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrScheduleNotFound):
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		case errors.Is(err, utils.ErrInvalidMovieIDs):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		case errors.Is(err, utils.ErrStudioNotFound):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		case errors.Is(err, utils.ErrScheduleConflict):
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, err.Error()))
		case errors.Is(err, utils.ErrInvalidTimeRange):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		case errors.Is(err, utils.ErrPastDate):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Schedule updated successfully",
		nil,
	))
}

func (c *ScheduleController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	err := c.ScheduleService.DeleteSchedule(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrScheduleNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Schedule deleted successfully",
		nil,
	))
}
