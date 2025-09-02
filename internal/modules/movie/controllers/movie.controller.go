package controllers

import (
	"errors"
	"mime/multipart"
	"net/http"
	"strconv"

	"movie-app-go/internal/modules/movie/requests"
	"movie-app-go/internal/modules/movie/responses"
	"movie-app-go/internal/modules/movie/services"
	"movie-app-go/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MovieController struct {
	MovieService *services.MovieService
}

func NewMovieController(s *services.MovieService) *MovieController {
	return &MovieController{MovieService: s}
}

func (c *MovieController) Create(ctx *gin.Context) {
	var req requests.CreateMovieRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	var posterFile *multipart.FileHeader
	if file, err := ctx.FormFile("poster"); err == nil {
		posterFile = file
	}

	err := c.MovieService.CreateMovie(&req, posterFile)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidGenreIDs):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(
		http.StatusCreated,
		"Movie created successfully",
		nil,
	))
}

func (c *MovieController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.MovieService.GetAllMoviesPaginated(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	movieResponses := responses.ToMovieResponses(result.Data)
	response := responses.PaginatedMovieResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      movieResponses,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Movies retrieved successfully",
		response,
	))
}

func (c *MovieController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	movie, err := c.MovieService.GetMovieByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse("Movie not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Movie retrieved successfully",
		responses.ToMovieResponse(movie),
	))
}

func (c *MovieController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req requests.UpdateMovieRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	var posterFile *multipart.FileHeader
	if file, err := ctx.FormFile("poster"); err == nil {
		posterFile = file
	}

	err := c.MovieService.UpdateMovie(uint(id), &req, posterFile)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrMovieNotFound):
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		case errors.Is(err, utils.ErrInvalidGenreIDs):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Movie updated successfully",
		nil,
	))
}

func (c *MovieController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	err := c.MovieService.DeleteMovie(uint(id))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrMovieNotFound):
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		case errors.Is(err, utils.ErrMovieHasSchedules):
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Movie deleted successfully",
		nil,
	))
}
