package movie

import (
	"movie-app-go/internal/modules/movie/controllers"
	"movie-app-go/internal/modules/movie/services"
	"movie-app-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MovieModule struct {
	MovieController *controllers.MovieController
	GenreController *controllers.GenreController
}

func NewMovieModule(db *gorm.DB) *MovieModule {
	movieService := services.NewMovieService(db)
	genreService := services.NewGenreService(db)

	return &MovieModule{
		MovieController: controllers.NewMovieController(movieService),
		GenreController: controllers.NewGenreController(genreService),
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *MovieModule) {
	// Movie
	rg.POST("/movies", middleware.AdminOnly(), module.MovieController.Create)
	rg.GET("/movies", module.MovieController.GetAll)
	rg.GET("/movies/:id", module.MovieController.GetByID)
	rg.PUT("/movies/:id", middleware.AdminOnly(), module.MovieController.Update)
	rg.DELETE("/movies/:id", middleware.AdminOnly(), module.MovieController.Delete)

	// Genre
	rg.POST("/genres", middleware.AdminOnly(), module.GenreController.Create)
	rg.GET("/genres", module.GenreController.GetAll)
	rg.GET("/genres/:id", module.GenreController.GetByID)
	rg.PUT("/genres/:id", middleware.AdminOnly(), module.GenreController.Update)
	rg.DELETE("/genres/:id", middleware.AdminOnly(), module.GenreController.Delete)
}
