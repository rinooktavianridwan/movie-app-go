package main

import (
    "github.com/joho/godotenv"
	"os"
	"fmt"
	"movie-app-go/database"
)

func main() {
    godotenv.Load()

	port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

	_, err := database.ConnectDB()
    if err != nil {
        fmt.Println("Gagal koneksi ke database:", err)
        return
    }
}