package controllers

import (
	"errors"
	"movie-app-go/internal/modules/iam/requests"
	"movie-app-go/internal/modules/iam/responses"
	"movie-app-go/internal/modules/iam/services"
	"movie-app-go/internal/utils"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

	userIDUint, ok := userID.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid user ID"))
		return
	}
	err = c.Service.UpdateAvatar(userIDUint, file)
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

func (c *UserController) DownloadImportTemplate(ctx *gin.Context) {
	templatePath := filepath.Join("internal", "modules", "iam", "templates", "User.xlsx")

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "template file tidak ditemukan"})
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename=users_import_template.xlsx")
	ctx.File(templatePath)
}

func (c *UserController) ImportUserExcelSingleSheet(ctx *gin.Context) {
	var req requests.ImportUserRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("file diperlukan"))
		return
	}

	var sheetName *string
	if sn := strings.TrimSpace(req.SheetName); sn != "" {
		req.SheetName = sn
		sheetName = &req.SheetName
	}

	err = c.Service.ImportExcelSingleSheet(file, sheetName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("gagal parsing file"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Import User Berhasil",
		nil,
	))
}

func (c *UserController) ImportUserExcelMultiSheet(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("file diperlukan"))
		return
	}

	err = c.Service.ImportExcelMultiSheet(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("gagal parsing file"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Import User Berhasil",
		nil,
	))
}

func (c *UserController) ExportUsers(ctx *gin.Context) {
	file, err := c.Service.ExportUsersExcel()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Disposition", "attachment; filename=users_export.xlsx")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file)
}
