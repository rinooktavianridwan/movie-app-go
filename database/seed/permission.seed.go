package seed

import (
	"log"
	"movie-app-go/internal/models"

	"gorm.io/gorm"
)

func SeedPermissions(db *gorm.DB) error {
	log.Println("Seeding permissions...")

	permissions := []models.Permission{
		// User management
		{Name: "users.create", Resource: "users", Action: "create", Description: "Create new users"},
		{Name: "users.read", Resource: "users", Action: "read", Description: "View users"},
		{Name: "users.update", Resource: "users", Action: "update", Description: "Update users"},
		{Name: "users.delete", Resource: "users", Action: "delete", Description: "Delete users"},

		// Role management
		{Name: "roles.create", Resource: "roles", Action: "create", Description: "Create new roles"},
		{Name: "roles.read", Resource: "roles", Action: "read", Description: "View roles"},
		{Name: "roles.update", Resource: "roles", Action: "update", Description: "Update roles"},
		{Name: "roles.delete", Resource: "roles", Action: "delete", Description: "Delete roles"},

		// Movie management
		{Name: "movies.create", Resource: "movies", Action: "create", Description: "Create new movies"},
		{Name: "movies.read", Resource: "movies", Action: "read", Description: "View movies"},
		{Name: "movies.update", Resource: "movies", Action: "update", Description: "Update movies"},
		{Name: "movies.delete", Resource: "movies", Action: "delete", Description: "Delete movies"},

		// Studio management
		{Name: "studios.create", Resource: "studios", Action: "create", Description: "Create new studios"},
		{Name: "studios.read", Resource: "studios", Action: "read", Description: "View studios"},
		{Name: "studios.update", Resource: "studios", Action: "update", Description: "Update studios"},
		{Name: "studios.delete", Resource: "studios", Action: "delete", Description: "Delete studios"},

		// Schedule management
		{Name: "schedules.create", Resource: "schedules", Action: "create", Description: "Create new schedules"},
		{Name: "schedules.read", Resource: "schedules", Action: "read", Description: "View schedules"},
		{Name: "schedules.update", Resource: "schedules", Action: "update", Description: "Update schedules"},
		{Name: "schedules.delete", Resource: "schedules", Action: "delete", Description: "Delete schedules"},

		// Order management
		{Name: "orders.create", Resource: "orders", Action: "create", Description: "Create new orders"},
		{Name: "orders.read", Resource: "orders", Action: "read", Description: "View orders"},
		{Name: "orders.update", Resource: "orders", Action: "update", Description: "Update orders"},
		{Name: "orders.cancel", Resource: "orders", Action: "cancel", Description: "Cancel orders"},

		// Promo management
		{Name: "promos.create", Resource: "promos", Action: "create", Description: "Create new promos"},
		{Name: "promos.read", Resource: "promos", Action: "read", Description: "View promos"},
		{Name: "promos.update", Resource: "promos", Action: "update", Description: "Update promos"},
		{Name: "promos.delete", Resource: "promos", Action: "delete", Description: "Delete promos"},

		// Notification management
		{Name: "notifications.create", Resource: "notifications", Action: "create", Description: "Create notifications"},
		{Name: "notifications.read", Resource: "notifications", Action: "read", Description: "View notifications"},
		{Name: "notifications.update", Resource: "notifications", Action: "update", Description: "Update notifications"},
		{Name: "notifications.delete", Resource: "notifications", Action: "delete", Description: "Delete notifications"},

		// Report access
		{Name: "reports.view", Resource: "reports", Action: "read", Description: "View reports"},
	}

	for _, permission := range permissions {
		var existingPermission models.Permission
		if err := db.Where("name = ?", permission.Name).First(&existingPermission).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&permission).Error; err != nil {
					log.Printf("Error seeding permission %s: %v", permission.Name, err)
					return err
				}
				log.Printf("Permission %s created successfully", permission.Name)
			} else {
				return err
			}
		} else {
			log.Printf("Permission %s already exists, skipping", permission.Name)
		}
	}

	log.Println("Permissions seeding completed")
	return nil
}
