package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"movie-app-go/internal/modules/genre/requests"
	"movie-app-go/internal/modules/genre/responses"
	"movie-app-go/internal/modules/genre/services"
	"movie-app-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type GenreController struct {
	GenreService *services.GenreService
}

func NewGenreController(s *services.GenreService) *GenreController {
	return &GenreController{GenreService: s}
}

func (c *GenreController) Create(ctx *gin.Context) {
	var req requests.CreateGenreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	err := c.GenreService.CreateGenre(req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(
		http.StatusCreated,
		"Genre created successfully",
		nil,
	))
}

func (c *GenreController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.GenreService.GetAllGenresPaginated(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	genreResponses := responses.ToGenreResponses(result.Data)
	response := responses.PaginatedGenreResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      genreResponses,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Genres retrieved successfully",
		response,
	))
}

func (c *GenreController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	genre, err := c.GenreService.GetGenreByID(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrGenreNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Genre retrieved successfully",
		responses.ToGenreResponse(genre),
	))
}

func (c *GenreController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req requests.UpdateGenreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	err := c.GenreService.UpdateGenre(uint(id), req.Name)
	if err != nil {
		if errors.Is(err, utils.ErrGenreNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Genre updated successfully",
		nil,
	))
}

func (c *GenreController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.GenreService.DeleteGenre(uint(id)); err != nil {
		switch {
		case errors.Is(err, utils.ErrGenreNotFound):
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		case errors.Is(err, utils.ErrGenreInUse):
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Genre deleted successfully",
		nil,
	))
}
