package controllers

import (
	"net/http"
	"strconv"

	"movie-app-go/internal/modules/schedule/options"
	"movie-app-go/internal/modules/schedule/requests"
	"movie-app-go/internal/modules/schedule/responses"
	"movie-app-go/internal/modules/schedule/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schedule, err := c.ScheduleService.CreateSchedule(&req)
	if err != nil {
		if err.Error() == "movie not found" || err.Error() == "studio not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "start_time must be before end_time" ||
			err.Error() == "date cannot be in the past" ||
			err.Error() == "schedule conflict: studio is already booked at this time" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusCreated, responses.ToScheduleResponse(schedule))
}

func (c *ScheduleController) GetAll(ctx *gin.Context) {
	// Parse options dari query parameters
	opts, err := options.ParseScheduleOptions(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Panggil service dengan options struct
	result, err := c.ScheduleService.GetAllSchedulesPaginated(opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	ctx.JSON(http.StatusOK, response)
}

func (c *ScheduleController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	schedule, err := c.ScheduleService.GetScheduleByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "schedule not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, responses.ToScheduleResponse(schedule))
}

func (c *ScheduleController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req requests.UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	schedule, err := c.ScheduleService.UpdateSchedule(uint(id), &req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "schedule not found"})
		} else if err.Error() == "movie not found" || err.Error() == "studio not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "start_time must be before end_time" ||
			err.Error() == "date cannot be in the past" ||
			err.Error() == "schedule conflict: studio is already booked at this time" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, responses.ToScheduleResponse(schedule))
}

func (c *ScheduleController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.ScheduleService.DeleteSchedule(uint(id)); err != nil {
		if err.Error() == "schedule not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "cannot delete schedule: tickets already exist" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "schedule deleted"})
}
