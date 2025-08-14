package iam

import (
    "movie-app-go/internal/modules/iam/services"
    "movie-app-go/internal/modules/iam/controllers"
    "movie-app-go/internal/middleware"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type IAMModule struct {
    UserController *controllers.UserController
    AuthController  *controllers.AuthController
}

func NewIAMModule(db *gorm.DB) *IAMModule {
    userService := services.NewUserService(db)
    authService := services.NewAuthService(db)
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
    r.DELETE("/users/:id", middleware.AdminOnly(), iam.UserController.Delete)

    r.POST("/auth/register", iam.AuthController.Register)
    r.POST("/auth/login", iam.AuthController.Login)
}