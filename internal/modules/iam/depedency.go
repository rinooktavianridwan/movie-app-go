package iam

import (
	"movie-app-go/internal/middleware"
	"movie-app-go/internal/modules/iam/controllers"
	"movie-app-go/internal/modules/iam/repositories"
	"movie-app-go/internal/modules/iam/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAMModule struct {
	UserController       *controllers.UserController
	AuthController       *controllers.AuthController
	RoleController       *controllers.RoleController
	PermissionController *controllers.PermissionController
	PermissionRepo       *repositories.PermissionRepository
}

func NewIAMModule(db *gorm.DB) *IAMModule {
	services.InitRedis()

	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)

	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(authRepo, roleRepo)
	roleService := services.NewRoleService(roleRepo)
	permissionService := services.NewPermissionService(permissionRepo)

	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)
	roleController := controllers.NewRoleController(roleService)
	permissionController := controllers.NewPermissionController(permissionService)

	return &IAMModule{
		UserController:       userController,
		AuthController:       authController,
		RoleController:       roleController,
		PermissionController: permissionController,
		PermissionRepo:       permissionRepo,
	}
}

func RegisterRoutes(rg *gin.RouterGroup, module *IAMModule, mf *middleware.Factory) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", module.AuthController.Register)
		auth.POST("/login", module.AuthController.Login)
		auth.POST("/logout", mf.Auth(), module.AuthController.Logout)
	}

	users := rg.Group("/users", mf.Auth())
	{
		users.GET("", mf.RequirePermission("users.read"), module.UserController.GetAll)
		users.GET("/:id", mf.RequirePermission("users.read"), module.UserController.GetByID)
		users.PUT("/:id", mf.RequirePermission("users.update"), module.UserController.Update)
		users.DELETE("/:id", mf.RequirePermission("users.delete"), module.UserController.Delete)
		users.POST("/:id/avatar", module.UserController.UploadAvatar)
	}

	roles := rg.Group("/roles", mf.Auth())
	{
		roles.GET("", mf.RequirePermission("roles.read"), module.RoleController.GetAllRoles)
		roles.GET("/:id", mf.RequirePermission("roles.read"), module.RoleController.GetRoleByID)
	}

	permissions := rg.Group("/permissions", mf.Auth())
	{
		permissions.GET("", mf.RequirePermission("roles.read"), module.PermissionController.GetAll)
		permissions.GET("/grouped", mf.RequirePermission("roles.read"), module.PermissionController.GetAllGroupedByResource)
		permissions.GET("/resource/:resource", mf.RequirePermission("roles.read"), module.PermissionController.GetByResource)
		permissions.GET("/:id", mf.RequirePermission("roles.read"), module.PermissionController.GetByID)
	}
}
