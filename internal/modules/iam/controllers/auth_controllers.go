package controllers

import (
	"errors"
	"movie-app-go/internal/modules/iam/requests"
	"movie-app-go/internal/modules/iam/responses"
	"movie-app-go/internal/modules/iam/services"
	"movie-app-go/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	err := a.AuthService.Register(&req)
	if err != nil {
		if errors.Is(err, utils.ErrEmailAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(
		http.StatusCreated,
		"User registered successfully",
		nil,
	))
}

func (a *AuthController) Login(ctx *gin.Context) {
	var req requests.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	user, token, err := a.AuthService.Login(&req)
	if err != nil {
		if errors.Is(err, utils.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Login successful",
		gin.H{
			"user":  responses.ToUserResponse(user),
			"token": token,
		},
	))
}

func (a *AuthController) Logout(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Authorization header missing"))
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid authorization header"))
		return
	}
	tokenString := parts[1]
	if err := a.AuthService.Logout(tokenString); err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse(err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Logout successful",
		nil,
	))
}
