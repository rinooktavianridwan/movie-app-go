package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"movie-app-go/internal/modules/studio/requests"
	"movie-app-go/internal/modules/studio/responses"
	"movie-app-go/internal/modules/studio/services"
)

type FacilityController struct {
	FacilityService *services.FacilityService
}

func NewFacilityController(s *services.FacilityService) *FacilityController {
	return &FacilityController{FacilityService: s}
}

func (c *FacilityController) Create(ctx *gin.Context) {
	var req requests.CreateFacilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	facility, err := c.FacilityService.CreateFacility(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, responses.FacilityResponse{ID: facility.ID, Name: facility.Name})
}

func (c *FacilityController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.FacilityService.GetAllFacilitiesPaginated(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := responses.ToFacilityResponses(result.Data)
    response := responses.PaginatedFacilityResponse{
        Page:      result.Page,
        PerPage:   result.PerPage,
        Total:     result.Total,
        TotalPage: result.TotalPages,
        Data:      resp,
    }
    ctx.JSON(http.StatusOK, response)
}

func (c *FacilityController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	facility, err := c.FacilityService.GetFacilityByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "facility not found"})
		return
	}
	ctx.JSON(http.StatusOK, responses.FacilityResponse{ID: facility.ID, Name: facility.Name})
}

func (c *FacilityController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req requests.UpdateFacilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	facility, err := c.FacilityService.UpdateFacility(uint(id), &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "facility not found"})
		return
	}
	ctx.JSON(http.StatusOK, responses.FacilityResponse{ID: facility.ID, Name: facility.Name})
}

func (c *FacilityController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.FacilityService.DeleteFacility(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "facility not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "facility deleted"})
}
