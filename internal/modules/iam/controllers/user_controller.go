package controllers

import (
	"errors"
	"movie-app-go/internal/modules/iam/requests"
	"movie-app-go/internal/modules/iam/responses"
	"movie-app-go/internal/modules/iam/services"
	"movie-app-go/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{Service: service}
}

func (c *UserController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.Service.GetAllPaginated(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	resp := responses.ToUserResponses(result.Data)
	response := responses.PaginatedUserResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      resp,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Users retrieved successfully",
		response,
	))
}

func (c *UserController) GetByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	user, err := c.Service.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse("User not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"User retrieved successfully",
		responses.ToUserResponse(user),
	))
}

func (c *UserController) Update(ctx *gin.Context) {
	var req requests.UserUpdateRequest
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	err := c.Service.Update(uint(id), &req)
	if err != nil {
		if errors.Is(err, utils.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"User updated successfully",
		nil,
	))
}

func (c *UserController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.Service.Delete(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse("User not found"))
			return
		}
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"User deleted successfully",
		nil,
	))
}

func (c *UserController) UploadAvatar(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Unauthorized"))
		return
	}

	file, err := ctx.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("File avatar diperlukan"))
		return
	}

	userIDFloat := uint(userID.(float64))
	err = c.Service.UpdateAvatar(userIDFloat, file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Avatar berhasil diupload",
		nil,
	))
}
