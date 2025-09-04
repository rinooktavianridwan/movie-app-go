package seed

import (
	"movie-app-go/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) error {
	var adminRole, customerRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "customer").First(&customerRole).Error; err != nil {
		return err
	}

	users := []models.User{
		{Name: "Admin", Email: "admin@bioskop.com", Password: "admin123", RoleID: &adminRole.ID},
		{Name: "User", Email: "user@bioskop.com", Password: "user123", RoleID: &customerRole.ID},
		{Name: "Alice", Email: "alice@mail.com", Password: "alice123", RoleID: &customerRole.ID},
		{Name: "Bob", Email: "bob@mail.com", Password: "bob123", RoleID: &customerRole.ID},
		{Name: "Charlie", Email: "charlie@mail.com", Password: "charlie123", RoleID: &customerRole.ID},
		{Name: "Diana", Email: "diana@mail.com", Password: "diana123", RoleID: &customerRole.ID},
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
