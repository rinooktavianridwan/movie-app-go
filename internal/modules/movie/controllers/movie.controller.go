package controllers

import (
	"net/http"
	"strconv"

	"movie-app-go/internal/modules/movie/requests"
	"movie-app-go/internal/modules/movie/responses"
	"movie-app-go/internal/modules/movie/services"

	"github.com/gin-gonic/gin"
)

type MovieController struct {
	MovieService *services.MovieService
}

func NewMovieController(s *services.MovieService) *MovieController {
	return &MovieController{MovieService: s}
}

func (c *MovieController) Create(ctx *gin.Context) {
	var req requests.CreateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movie, err := c.MovieService.CreateMovie(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, responses.ToMovieResponse(movie))
}

func (c *MovieController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.MovieService.GetAllMoviesPaginated(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mapping ke response
	movieResponses := responses.ToMovieResponses(result.Data)
	paginatedResponse := responses.PaginatedMovieResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      movieResponses,
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

func (c *MovieController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	movie, err := c.MovieService.GetMovieByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}
	ctx.JSON(http.StatusOK, responses.ToMovieResponse(movie))
}

func (c *MovieController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req requests.UpdateMovieRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movie, err := c.MovieService.UpdateMovie(uint(id), &req)
	if err != nil {
		if err.Error() == "some genre_ids are invalid" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "movie not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, responses.ToMovieResponse(movie))
}

func (c *MovieController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.MovieService.DeleteMovie(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "movie deleted"})
}
