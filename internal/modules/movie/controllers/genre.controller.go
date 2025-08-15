package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "movie-app-go/internal/modules/movie/services"
)

type GenreController struct {
    GenreService *services.GenreService
}

func NewGenreController(s *services.GenreService) *GenreController {
    return &GenreController{GenreService: s}
}

func (c *GenreController) Create(ctx *gin.Context) {
    var req struct {
        Name string `json:"name" binding:"required"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    genre, err := c.GenreService.CreateGenre(req.Name)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, genre)
}

func (c *GenreController) GetAll(ctx *gin.Context) {
    genres, err := c.GenreService.GetAllGenres()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, genres)
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
    ctx.JSON(http.StatusOK, genre)
}

func (c *GenreController) Update(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    var req struct {
        Name string `json:"name" binding:"required"`
    }
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    genre, err := c.GenreService.UpdateGenre(uint(id), req.Name)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "genre not found"})
        return
    }
    ctx.JSON(http.StatusOK, genre)
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