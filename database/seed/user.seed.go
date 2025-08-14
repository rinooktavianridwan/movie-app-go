package seed

import (
    "movie-app-go/internal/models"
    "gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
    users := []models.User{
        {Name: "Admin", Email: "admin@bioskop.com", Password: "admin123", IsAdmin: true},
        {Name: "User", Email: "user@bioskop.com", Password: "user123", IsAdmin: false},
    }
    for _, u := range users {
        var existing models.User
        if err := db.Where("email = ?", u.Email).First(&existing).Error; err == gorm.ErrRecordNotFound {
            if err := db.Create(&u).Error; err != nil {
                return err
            }
        }
    }
    return nil
}