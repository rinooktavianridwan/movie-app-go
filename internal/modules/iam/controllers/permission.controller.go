package controllers

import (
	"movie-app-go/internal/modules/iam/services"
	"movie-app-go/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	PermissionService *services.PermissionService
}

func NewPermissionController(permissionService *services.PermissionService) *PermissionController {
	return &PermissionController{PermissionService: permissionService}
}

func (c *PermissionController) GetAll(ctx *gin.Context) {
    pageStr := ctx.DefaultQuery("page", "1")
    perPageStr := ctx.DefaultQuery("per_page", "10")

    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1
    }

    perPage, err := strconv.Atoi(perPageStr)
    if err != nil || perPage < 1 {
        perPage = 10
    }

    result, err := c.PermissionService.GetAllPaginated(page, perPage)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(
        http.StatusOK,
        "Permissions retrieved successfully",
        result,
    ))
}

func (c *PermissionController) GetByID(ctx *gin.Context) {
    idParam := ctx.Param("id")
    id, err := strconv.ParseUint(idParam, 10, 32)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid permission ID"))
        return
    }

    permission, err := c.PermissionService.GetByID(uint(id))
    if err != nil {
        if err.Error() == "permission not found" {
            ctx.JSON(http.StatusNotFound, utils.NotFoundResponse("Permission not found"))
            return
        }
        ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(
        http.StatusOK,
        "Permission retrieved successfully",
        permission,
    ))
}

func (c *PermissionController) GetByResource(ctx *gin.Context) {
    resource := ctx.Param("resource")
    if resource == "" {
        ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Resource parameter is required"))
        return
    }

    permissions, err := c.PermissionService.GetByResource(resource)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(
        http.StatusOK,
        "Permissions retrieved successfully",
        permissions,
    ))
}

func (c *PermissionController) GetAllGroupedByResource(ctx *gin.Context) {
    result, err := c.PermissionService.GetAllGroupedByResource()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
        return
    }

    ctx.JSON(http.StatusOK, utils.SuccessResponse(
        http.StatusOK,
        "Permissions grouped by resource retrieved successfully",
        result,
    ))
}
