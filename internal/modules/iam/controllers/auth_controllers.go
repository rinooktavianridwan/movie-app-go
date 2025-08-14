package controllers

import (
    "net/http"
    "movie-app-go/internal/modules/iam/services"
    "github.com/gin-gonic/gin"
    "movie-app-go/internal/modules/iam/responses"
    "movie-app-go/internal/modules/iam/requests"
)

type AuthController struct {
    AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
    return &AuthController{AuthService: authService}
}

func (a *AuthController) Register(ctx *gin.Context) {
    var req requests.RegisterRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, err := a.AuthService.Register(&req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, responses.ToUserResponse(user))
}

func (a *AuthController) Login(ctx *gin.Context) {
    var req requests.LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user, token, err := a.AuthService.Login(&req)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{
        "user": responses.ToUserResponse(user),
        "token": token,
    })
}