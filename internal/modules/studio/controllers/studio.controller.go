package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"movie-app-go/internal/modules/studio/requests"
	"movie-app-go/internal/modules/studio/responses"
	"movie-app-go/internal/modules/studio/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StudioController struct {
	StudioService *services.StudioService
}

func NewStudioController(s *services.StudioService) *StudioController {
	return &StudioController{StudioService: s}
}

func (c *StudioController) Create(ctx *gin.Context) {
	var req requests.CreateStudioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	studio, err := c.StudioService.CreateStudio(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, responses.ToStudioResponse(studio))
}

func (c *StudioController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.StudioService.GetAllStudiosPaginated(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := responses.ToStudioResponses(result.Data)
	response := responses.PaginatedStudioResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      resp,
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *StudioController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	studio, err := c.StudioService.GetStudioByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "studio not found"})
		return
	}
	ctx.JSON(http.StatusOK, responses.ToStudioResponse(studio))
}

func (c *StudioController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req requests.CreateStudioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	studio, err := c.StudioService.UpdateStudio(uint(id), &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "studio not found"})
		} else if err.Error() == "some facility_ids are invalid" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, responses.ToStudioResponse(studio))
}

func (c *StudioController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.StudioService.DeleteStudio(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "studio not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "studio deleted"})
}
