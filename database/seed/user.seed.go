package seed

import (
    "movie-app-go/internal/models"
    "gorm.io/gorm"
    "golang.org/x/crypto/bcrypt"
)

func SeedUsers(db *gorm.DB) error {
    users := []models.User{
        {Name: "Admin", Email: "admin@bioskop.com", Password: "admin123", IsAdmin: true},
        {Name: "User", Email: "user@bioskop.com", Password: "user123", IsAdmin: false},
        {Name: "Alice", Email: "alice@mail.com", Password: "alice123", IsAdmin: false},
        {Name: "Bob", Email: "bob@mail.com", Password: "bob123", IsAdmin: false},
        {Name: "Charlie", Email: "charlie@mail.com", Password: "charlie123", IsAdmin: false},
        {Name: "Diana", Email: "diana@mail.com", Password: "diana123", IsAdmin: false},
    }
    
    var usersToInsert []models.User
    for _, u := range users {
        hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashed)
        usersToInsert = append(usersToInsert, u)
    }

    if len(usersToInsert) > 0 {
        if err := db.Create(&usersToInsert).Error; err != nil {
            return err
        }
    }
    return nil
}