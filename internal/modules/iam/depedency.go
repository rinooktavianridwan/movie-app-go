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
	UserController *controllers.UserController
	AuthController *controllers.AuthController
}

func NewIAMModule(db *gorm.DB) *IAMModule {
	services.InitRedis()
	userRepo := repositories.NewUserRepository(db)
	authRepo := repositories.NewAuthRepository(db)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(authRepo)
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)
	return &IAMModule{
		UserController: userController,
		AuthController: authController,
	}
}

func RegisterRoutes(r *gin.RouterGroup, iam *IAMModule) {
	r.GET("/users", middleware.AdminOnly(), iam.UserController.GetAll)
	r.GET("/users/:id", middleware.AdminOnly(), iam.UserController.GetByID)
	r.PUT("/users/:id", middleware.AdminOnly(), iam.UserController.Update)
	r.PUT("/users/upload-avatar", middleware.Auth(), iam.UserController.UploadAvatar)
	r.DELETE("/users/:id", middleware.AdminOnly(), iam.UserController.Delete)

	r.POST("/auth/register", iam.AuthController.Register)
	r.POST("/auth/login", iam.AuthController.Login)
	r.POST("/auth/logout", middleware.Auth(), iam.AuthController.Logout)
}
