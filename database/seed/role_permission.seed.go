package seed

import (
	"fmt"
	"log"
	"movie-app-go/internal/models"

	"gorm.io/gorm"
)

func SeedRolePermissions(db *gorm.DB) error {
    log.Println("Seeding role permissions...")

    var adminRole, managerRole, staffRole, customerRole models.Role
    if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
        return err
    }
    if err := db.Where("name = ?", "manager").First(&managerRole).Error; err != nil {
        return err
    }
    if err := db.Where("name = ?", "staff").First(&staffRole).Error; err != nil {
        return err
    }
    if err := db.Where("name = ?", "customer").First(&customerRole).Error; err != nil {
        return err
    }

    var permissions []models.Permission
    if err := db.Find(&permissions).Error; err != nil {
        return err
    }

    permissionMap := make(map[string]uint)
    for _, p := range permissions {
        permissionMap[p.Name] = p.ID
    }

    rolePermissions := []struct {
        RoleID      uint
        RoleName    string
        Permissions []string
    }{
        {
            RoleID:   adminRole.ID,
            RoleName: "admin",
            Permissions: []string{
                "users.create", "users.read", "users.update", "users.delete",
                "roles.create", "roles.read", "roles.update", "roles.delete",
                "movies.create", "movies.read", "movies.update", "movies.delete",
                "studios.create", "studios.read", "studios.update", "studios.delete",
                "schedules.create", "schedules.read", "schedules.update", "schedules.delete",
                "orders.create", "orders.read", "orders.update", "orders.cancel",
                "promos.create", "promos.read", "promos.update", "promos.delete",
                "notifications.create", "notifications.read", "notifications.update", "notifications.delete",
                "reports.view",
            },
        },
        {
            RoleID:   managerRole.ID,
            RoleName: "manager",
            Permissions: []string{
                "users.read", "users.update",
                "roles.read",
                "movies.create", "movies.read", "movies.update",
                "studios.create", "studios.read", "studios.update",
                "schedules.create", "schedules.read", "schedules.update",
                "orders.read", "orders.update", "orders.cancel",
                "promos.create", "promos.read", "promos.update",
                "notifications.create", "notifications.read",
                "reports.view",
            },
        },
        {
            RoleID:   staffRole.ID,
            RoleName: "staff",
            Permissions: []string{
                "movies.read",
                "studios.read",
                "schedules.read",
                "orders.read", "orders.update",
                "promos.read",
                "notifications.read",
            },
        },
        {
            RoleID:   customerRole.ID,
            RoleName: "customer",
            Permissions: []string{
                "movies.read",
                "studios.read",
                "schedules.read",
                "orders.create", "orders.read",
            },
        },
    }

    var allRolePermissions []models.RolePermission

    for _, rp := range rolePermissions {
        log.Printf("Preparing permissions for role: %s", rp.RoleName)
        
        for _, permName := range rp.Permissions {
            if permID, exists := permissionMap[permName]; exists {
                allRolePermissions = append(allRolePermissions, models.RolePermission{
                    RoleID:       rp.RoleID,
                    PermissionID: permID,
                })
            } else {
                log.Printf("Permission %s not found", permName)
            }
        }
    }

    var existingRolePermissions []models.RolePermission
    if err := db.Find(&existingRolePermissions).Error; err != nil {
        return err
    }

    existingMap := make(map[string]bool)
    for _, existing := range existingRolePermissions {
        key := fmt.Sprintf("%d_%d", existing.RoleID, existing.PermissionID)
        existingMap[key] = true
    }

    var newRolePermissions []models.RolePermission
    for _, rp := range allRolePermissions {
        key := fmt.Sprintf("%d_%d", rp.RoleID, rp.PermissionID)
        if !existingMap[key] {
            newRolePermissions = append(newRolePermissions, rp)
        }
    }

    if len(newRolePermissions) > 0 {
        if err := db.Create(&newRolePermissions).Error; err != nil {
            return err
        }
        log.Printf("Successfully created %d new role permissions", len(newRolePermissions))
    } else {
        log.Println("All role permissions already exist")
    }

    log.Println("Role permissions seeding completed")
    return nil
}
