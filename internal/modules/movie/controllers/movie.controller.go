package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "movie-app-go/internal/modules/movie/services"
    "movie-app-go/internal/modules/movie/requests"
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
    ctx.JSON(http.StatusCreated, movie)
}

func (c *MovieController) GetAll(ctx *gin.Context) {
    movies, err := c.MovieService.GetAllMovies()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, movies)
}

func (c *MovieController) GetByID(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    movie, genreIDs, err := c.MovieService.GetMovieByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{
        "id":        movie.ID,
        "title":     movie.Title,
        "overview":  movie.Overview,
        "duration":  movie.Duration,
        "genre_ids": genreIDs,
    })
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
        ctx.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
        return
    }
    ctx.JSON(http.StatusOK, movie)
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