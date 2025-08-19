package main

import (
	"fmt"
	"movie-app-go/database"
	// "movie-app-go/database/seed"
	"movie-app-go/internal/modules/iam"
	"movie-app-go/internal/modules/genre"
	"movie-app-go/internal/modules/studio"
	"movie-app-go/internal/modules/movie"
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

	// Setup Gin
	r := gin.Default()
	iam.RegisterRoutes(r.Group("/api"), iamModule)
	studio.RegisterRoutes(r.Group("/api"), studioModule)
	movie.RegisterRoutes(r.Group("/api"), movieModule)
	genre.RegisterRoutes(r.Group("/api"), genreModule)

	// Run server
	r.Run(":" + port)
}
