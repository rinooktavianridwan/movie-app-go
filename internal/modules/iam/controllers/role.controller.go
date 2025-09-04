package controllers

import (
	"movie-app-go/internal/modules/iam/responses"
	"movie-app-go/internal/modules/iam/services"
	"movie-app-go/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	RoleService *services.RoleService
}

func NewRoleController(roleService *services.RoleService) *RoleController {
	return &RoleController{RoleService: roleService}
}

func (c *RoleController) GetAllRoles(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.RoleService.GetAllRolesPaginated(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse("Failed to retrieve roles"))
		return
	}

	roleResponses := responses.ToRoleResponses(result.Data)
	response := responses.PaginatedRoleResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      roleResponses,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Roles retrieved successfully",
		response,
	))
}

func (c *RoleController) GetRoleByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid role ID"))
		return
	}

	role, err := c.RoleService.GetRoleByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, utils.NotFoundResponse("Role not found"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Role retrieved successfully",
		responses.ToRoleResponse(role),
	))
}
