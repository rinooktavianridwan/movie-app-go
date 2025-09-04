package movie

import (
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/movie/controllers"
	"movie-app-go/internal/modules/movie/repositories"
	"movie-app-go/internal/modules/movie/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MovieModule struct {
	MovieController *controllers.MovieController
}

func NewMovieModule(db *gorm.DB) *MovieModule {
	movieRepo := repositories.NewMovieRepository(db)
	movieService := services.NewMovieService(movieRepo)

	return &MovieModule{
		MovieController: controllers.NewMovieController(movieService),
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *MovieModule, mf *middleware.Factory) {
    rg.POST("/movies", mf.Auth(), mf.RequirePermission("movies.create"), module.MovieController.Create)
    rg.GET("/movies", module.MovieController.GetAll)
    rg.GET("/movies/:id", module.MovieController.GetByID)
    rg.PUT("/movies/:id", mf.Auth(), mf.RequirePermission("movies.update"), module.MovieController.Update)
    rg.DELETE("/movies/:id", mf.Auth(), mf.RequirePermission("movies.delete"), module.MovieController.Delete)
}
