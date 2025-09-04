package seed

import (
    "log"
    "movie-app-go/internal/models"

    "gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {
    log.Println("Seeding roles...")

    roles := []models.Role{
        {
            Name:        "admin",
            Description: "Administrator with full access",
        },
        {
            Name:        "manager", 
            Description: "Manager with limited administrative access",
        },
        {
            Name:        "staff",
            Description: "Staff with basic operational access",
        },
        {
            Name:        "customer",
            Description: "Regular customer",
        },
    }

    for _, role := range roles {
        var existingRole models.Role
        if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
            if err == gorm.ErrRecordNotFound {
                if err := db.Create(&role).Error; err != nil {
                    log.Printf("Error seeding role %s: %v", role.Name, err)
                    return err
                }
                log.Printf("Role %s created successfully", role.Name)
            } else {
                return err
            }
        } else {
            log.Printf("Role %s already exists, skipping", role.Name)
        }
    }

    log.Println("Roles seeding completed")
    return nil
}