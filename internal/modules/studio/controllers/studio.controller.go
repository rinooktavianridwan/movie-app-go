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

type StudioController struct {
	StudioService *services.StudioService
}

func NewStudioController(s *services.StudioService) *StudioController {
	return &StudioController{StudioService: s}
}

func (c *StudioController) Create(ctx *gin.Context) {
	var req requests.CreateStudioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}
	err := c.StudioService.CreateStudio(&req)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidFacilityIDs) {
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(
		http.StatusCreated,
		"Studio created successfully",
		nil,
	))
}

func (c *StudioController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.StudioService.GetAllStudiosPaginated(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
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

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Studios retrieved successfully",
		response,
	))
}

func (c *StudioController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	studio, err := c.StudioService.GetStudioByID(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrStudioNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Studio retrieved successfully",
		responses.ToStudioResponse(studio),
	))
}

func (c *StudioController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req requests.CreateStudioRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}
	err := c.StudioService.UpdateStudio(uint(id), &req)
	if err != nil {
		if errors.Is(err, utils.ErrStudioNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else if errors.Is(err, utils.ErrInvalidFacilityIDs) {
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Studio updated successfully",
		nil,
	))
}

func (c *StudioController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	err := c.StudioService.DeleteStudio(uint(id))
	if err != nil {
        switch {
        case errors.Is(err, utils.ErrStudioNotFound):
            ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
        case errors.Is(err, utils.ErrStudioHasSchedules):
            ctx.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, err.Error()))
        default:
            ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
        }
        return
    }

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Studio deleted successfully",
		nil,
	))
}
