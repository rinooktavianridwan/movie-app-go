package controllers

import (
	"net/http"
	"strconv"

	"movie-app-go/internal/modules/genre/requests"
	"movie-app-go/internal/modules/genre/responses"
	"movie-app-go/internal/modules/genre/services"

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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	genre, err := c.GenreService.CreateGenre(req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, responses.ToGenreResponse(genre))
}

func (c *GenreController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.GenreService.GetAllGenresPaginated(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := responses.ToGenreResponses(result.Data)
	response := responses.PaginatedGenreResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      resp,
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *GenreController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	genre, err := c.GenreService.GetGenreByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
		return
	}
	ctx.JSON(http.StatusOK, responses.ToGenreResponse(genre))
}

func (c *GenreController) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req requests.UpdateGenreRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	genre, err := c.GenreService.UpdateGenre(uint(id), req.Name)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
		return
	}
	ctx.JSON(http.StatusOK, responses.ToGenreResponse(genre))
}

func (c *GenreController) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := c.GenreService.DeleteGenre(uint(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "genre deleted"})
}
