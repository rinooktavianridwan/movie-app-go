package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"movie-app-go/internal/constants"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/notification/options"
	"movie-app-go/internal/modules/notification/requests"
	"movie-app-go/internal/modules/notification/responses"
	"movie-app-go/internal/modules/notification/services"
	"movie-app-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationService *services.NotificationService
}

func NewNotificationController(notificationService *services.NotificationService) *NotificationController {
	return &NotificationController{
		NotificationService: notificationService,
	}
}

func (c *NotificationController) getUserID(ctx *gin.Context) (uint, error) {
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		return 0, gin.Error{Err: fmt.Errorf("user not authenticated"), Type: gin.ErrorTypePublic}
	}

	switch v := userIDInterface.(type) {
	case uint:
		return v, nil
	case float64:
		return uint(v), nil
	case int:
		return uint(v), nil
	default:
		return 0, gin.Error{Err: fmt.Errorf("invalid user ID type"), Type: gin.ErrorTypePublic}
	}
}

func (c *NotificationController) GetUserNotifications(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	opts, err := options.ParseNotificationOptions(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	result, err := c.NotificationService.GetUserNotifications(userID, opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	notificationResponses := responses.ToNotificationResponses(result.Data)
	response := responses.PaginatedNotificationResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      notificationResponses,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"User notifications retrieved successfully",
		response,
	))
}

func (c *NotificationController) GetNotificationByID(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid notification ID"))
		return
	}

	notification, err := c.NotificationService.GetNotificationByID(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrNotificationNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	if notification.UserID != userID {
		ctx.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, "Access denied to this notification"))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Notification retrieved successfully",
		responses.ToNotificationResponse(*notification),
	))
}

func (c *NotificationController) MarkAsRead(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid notification ID"))
		return
	}

	notification, err := c.NotificationService.MarkAsRead(uint(id), userID)
	if err != nil {
		if errors.Is(err, utils.ErrNotificationNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Notification marked as read successfully",
		responses.ToNotificationResponse(*notification),
	))
}

func (c *NotificationController) MarkAsUnread(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid notification ID"))
		return
	}

	notification, err := c.NotificationService.MarkAsUnread(uint(id), userID)
	if err != nil {
		if errors.Is(err, utils.ErrNotificationNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Notification marked as unread successfully",
		responses.ToNotificationResponse(*notification),
	))
}

func (c *NotificationController) BulkMarkAsRead(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	var req requests.BulkMarkReadRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	if err := c.NotificationService.BulkMarkAsRead(&req, userID); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Notifications marked as read successfully",
		nil,
	))
}

func (c *NotificationController) MarkAllAsRead(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	if err := c.NotificationService.MarkAllAsRead(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"All notifications marked as read successfully",
		nil,
	))
}

func (c *NotificationController) GetNotificationStats(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	stats, err := c.NotificationService.GetNotificationStats(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Notification statistics retrieved successfully",
		stats,
	))
}

func (c *NotificationController) DeleteNotification(ctx *gin.Context) {
	userID, err := c.getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, utils.ErrorResponse(http.StatusUnauthorized, err.Error()))
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid notification ID"))
		return
	}

	if err := c.NotificationService.DeleteNotification(uint(id), userID); err != nil {
		if errors.Is(err, utils.ErrNotificationNotFound) {
            ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
        } else {
            ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
        }
        return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
        http.StatusOK,
        "Notification deleted successfully",
        nil,
    ))
}

func (c *NotificationController) CreateNotification(ctx *gin.Context) {
	var req requests.CreateNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	err := c.NotificationService.CreateNotification(&req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(
        http.StatusCreated,
        "Notification created successfully",
        nil,
    ))
}

func (c *NotificationController) CreateSystemNotification(ctx *gin.Context) {
	var req requests.CreateSystemNotificationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	var userIDs []uint
	err := c.NotificationService.DB.Model(&models.User{}).
        Where("is_admin = ?", false).
        Pluck("id", &userIDs).Error
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse("Failed to get users"))
        return
    }

	if len(userIDs) == 0 {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("No users found"))
		return
	}

	err = c.NotificationService.CreateBulkNotifications(
		userIDs,
		req.Title,
		req.Message,
		constants.NotificationTypeSystem,
		req.Data,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(
        http.StatusCreated,
        "System notification sent successfully",
        map[string]interface{}{
            "user_count": len(userIDs),
            "message":    fmt.Sprintf("System notification sent to %d users", len(userIDs)),
        },
    ))
}
