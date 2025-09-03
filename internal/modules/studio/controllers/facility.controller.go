package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"movie-app-go/internal/modules/studio/requests"
	"movie-app-go/internal/modules/studio/responses"
	"movie-app-go/internal/modules/studio/services"
	"movie-app-go/internal/utils"

	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	err := c.FacilityService.CreateFacility(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusCreated, utils.SuccessResponse(
		http.StatusCreated,
		"Facility created successfully",
		nil,
	))
}

func (c *FacilityController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.FacilityService.GetAllFacilitiesPaginated(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
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

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Facilities retrieved successfully",
		response,
	))
}

func (c *FacilityController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	facility, err := c.FacilityService.GetFacilityByID(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrFacilityNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Facility retrieved successfully",
		responses.ToFacilityResponse(facility),
	))
}

func (c *FacilityController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req requests.UpdateFacilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}
	err := c.FacilityService.UpdateFacility(uint(id), &req)
	if err != nil {
		if errors.Is(err, utils.ErrFacilityNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Facility updated successfully",
		nil,
	))
}

func (c *FacilityController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	err := c.FacilityService.DeleteFacility(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrFacilityNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Facility deleted successfully",
		nil,
	))
}
