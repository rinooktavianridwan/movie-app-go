package movie

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "movie-app-go/internal/modules/movie/controllers"
    "movie-app-go/internal/modules/movie/services"
)

type MovieModule struct {
    MovieController *controllers.MovieController
    GenreController *controllers.GenreController
}

func InitMovieModule(db *gorm.DB) *MovieModule {
    movieService := services.NewMovieService(db)
    genreService := services.NewGenreService(db)

    return &MovieModule{
        MovieController: controllers.NewMovieController(movieService),
        GenreController: controllers.NewGenreController(genreService),
    }
}

func RegisterMovieRoutes(rg *gin.RouterGroup, module *MovieModule) {
    // Movie
    rg.POST("/movies", module.MovieController.Create)
    rg.GET("/movies", module.MovieController.GetAll)
    rg.GET("/movies/:id", module.MovieController.GetByID)
    rg.PUT("/movies/:id", module.MovieController.Update)
    rg.DELETE("/movies/:id", module.MovieController.Delete)

    // Genre
    rg.POST("/genres", module.GenreController.Create)
    rg.GET("/genres", module.GenreController.GetAll)
    rg.GET("/genres/:id", module.GenreController.GetByID)
    rg.PUT("/genres/:id", module.GenreController.Update)
    rg.DELETE("/genres/:id", module.GenreController.Delete)
}