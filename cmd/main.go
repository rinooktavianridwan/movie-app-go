package main

import (
	"fmt"
	"movie-app-go/database"
	"movie-app-go/internal/jobs"


	// "movie-app-go/database/seed"
	"movie-app-go/internal/modules/genre"
	"movie-app-go/internal/modules/iam"
	"movie-app-go/internal/modules/movie"
	"movie-app-go/internal/modules/order"
	"movie-app-go/internal/modules/promo" 
	"movie-app-go/internal/modules/schedule"
	"movie-app-go/internal/modules/studio"
	"movie-app-go/internal/modules/report"
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

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	db, err := database.ConnectDB()
	if err != nil {
		fmt.Println("Gagal koneksi ke database:", err)
		return
	}

	queueService := jobs.NewQueueService(redisAddr)
	workerService := jobs.NewWorkerService(redisAddr, db)

	go func() {
		if err := workerService.Start(); err != nil {
			fmt.Printf("Could not start worker: %v\n", err)
		}
	}()

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
	promoModule := promo.NewPromoModule(db)
	orderModule := order.NewOrderModule(db, queueService, promoModule.PromoService)
	reportModule := report.NewReportModule(db)

	// Setup Gin
	r := gin.Default()
	iam.RegisterRoutes(r.Group("/api"), iamModule)
	studio.RegisterRoutes(r.Group("/api"), studioModule)
	movie.RegisterRoutes(r.Group("/api"), movieModule)
	genre.RegisterRoutes(r.Group("/api"), genreModule)
	schedule.RegisterRoutes(r.Group("/api"), scheduleModule)
	promo.RegisterRoutes(r.Group("/api"), promoModule)
	order.RegisterRoutes(r.Group("/api"), orderModule)
	report.RegisterRoutes(r.Group("/api"), reportModule)

	// Run server
	r.Run(":" + port)
}
