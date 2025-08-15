package main

import (
    "github.com/joho/godotenv"
    "github.com/gin-gonic/gin"
    "os"
    "fmt"
    "movie-app-go/database"
    // "movie-app-go/database/seed"
    "movie-app-go/internal/modules/iam"
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
    // if err := seed.SeedUsers(db); err != nil {
    //     fmt.Println("Gagal seeding user:", err)
    //     return
    // }

    // Dependency injection
    iamModule := iam.NewIAMModule(db)

    // Setup Gin
    r := gin.Default()
    iam.RegisterRoutes(r.Group("/api"), iamModule)

    // Run server
    r.Run(":" + port)
}