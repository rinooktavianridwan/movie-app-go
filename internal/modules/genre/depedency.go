package genre

import (
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/genre/controllers"
	"movie-app-go/internal/modules/genre/repositories"
	"movie-app-go/internal/modules/genre/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GenreModule struct {
	GenreController *controllers.GenreController
}

func NewGenreModule(db *gorm.DB) *GenreModule {
	genreRepo := repositories.NewGenreRepository(db)
	genreService := services.NewGenreService(genreRepo)

	return &GenreModule{
		GenreController: controllers.NewGenreController(genreService),
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *GenreModule, mf *middleware.Factory) {
	rg.POST("/genres", mf.Auth(), mf.RequirePermission("movies.create"), module.GenreController.Create)
	rg.GET("/genres", module.GenreController.GetAll)
	rg.GET("/genres/:id", module.GenreController.GetByID)
	rg.PUT("/genres/:id", mf.Auth(), mf.RequirePermission("movies.update"), module.GenreController.Update)
	rg.DELETE("/genres/:id", mf.Auth(), mf.RequirePermission("movies.delete"), module.GenreController.Delete)
}
