package main

import (
	"fmt"
	"movie-app-go/database"

	// "movie-app-go/database/seed"
	"movie-app-go/internal/modules/iam"
	"movie-app-go/internal/modules/genre"
	"movie-app-go/internal/modules/movie"
	"movie-app-go/internal/modules/schedule"
	"movie-app-go/internal/modules/studio"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println("Gagal koneksi ke database:", err)
		return
	}

	// Jalankan seeder user
	// if err := seed.RunAllSeeders(db); err != nil {
	// 	panic(err)
	// }

	// Dependency injection
	iamModule := iam.NewIAMModule(db)
	studioModule := studio.NewStudioModule(db)
	movieModule := movie.NewMovieModule(db)
	genreModule := genre.NewGenreModule(db)
	scheduleModule := schedule.NewScheduleModule(db)

	// Setup Gin
	r := gin.Default()
	iam.RegisterRoutes(r.Group("/api"), iamModule)
	studio.RegisterRoutes(r.Group("/api"), studioModule)
	movie.RegisterRoutes(r.Group("/api"), movieModule)
	genre.RegisterRoutes(r.Group("/api"), genreModule)
	schedule.RegisterRoutes(r.Group("/api"), scheduleModule)

	// Run server
	r.Run(":" + port)
}
