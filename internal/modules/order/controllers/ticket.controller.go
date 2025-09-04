package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/order/responses"
	"movie-app-go/internal/modules/order/services"
	"movie-app-go/internal/utils"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	TicketService *services.TicketService
}

func NewTicketController(s *services.TicketService) *TicketController {
	return &TicketController{TicketService: s}
}

func (c *TicketController) GetMyTickets(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	userIDUint, ok := userID.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid user ID"))
		return
	}
	result, err := c.TicketService.GetTicketsByUser(userIDUint, page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ticketResponses := responses.ToTicketResponses(result.Data)
	response := responses.PaginatedTicketResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      ticketResponses,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"User tickets retrieved successfully",
		response,
	))
}

func (c *TicketController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.TicketService.GetAllTickets(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	ticketResponses := responses.ToTicketResponses(result.Data)
	response := responses.PaginatedTicketResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      ticketResponses,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Tickets retrieved successfully",
		response,
	))
}

func (c *TicketController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid id"))
		return
	}
	userID, _ := ctx.Get("user_id")
	isAdmin, adminExists := ctx.Get("is_admin")

	var ticket *models.Ticket
	if adminExists && isAdmin.(bool) {
		ticket, err = c.TicketService.GetTicketByID(uint(id), nil)
	} else {
		userIDUint, ok := userID.(uint)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid user ID"))
			return
		}
		ticket, err = c.TicketService.GetTicketByID(uint(id), &userIDUint)
	}

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrTicketNotFound):
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		case errors.Is(err, utils.ErrUnauthorizedAccess):
			ctx.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Ticket retrieved successfully",
		responses.ToTicketResponse(ticket),
	))
}

func (c *TicketController) GetBySchedule(ctx *gin.Context) {
	scheduleID, _ := strconv.Atoi(ctx.Param("schedule_id"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "50"))

	result, err := c.TicketService.GetTicketsBySchedule(uint(scheduleID), page, perPage)
	if err != nil {
		if errors.Is(err, utils.ErrScheduleNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ticketResponses := responses.ToTicketResponses(result.Data)
	response := responses.PaginatedTicketResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      ticketResponses,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Schedule tickets retrieved successfully",
		response,
	))
}

func (c *TicketController) ScanTicket(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid id"))
		return
	}

	userID, exists := ctx.Get("user_id")
	isAdmin, adminExists := ctx.Get("is_admin")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("User not authenticated"))
		return
	}

	if adminExists && isAdmin.(bool) {
		err = c.TicketService.ScanTicket(uint(id), nil)
	} else {
		userIDUint, ok := userID.(uint)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("Invalid user ID"))
			return
		}
		err = c.TicketService.ScanTicket(uint(id), &userIDUint)
	}

	if err != nil {
		switch {
		case errors.Is(err, utils.ErrTicketNotFound):
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		case errors.Is(err, utils.ErrTicketAlreadyScanned):
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, err.Error()))
		case errors.Is(err, utils.ErrTicketNotPaid):
			ctx.JSON(http.StatusPaymentRequired, utils.ErrorResponse(http.StatusPaymentRequired, err.Error()))
		case errors.Is(err, utils.ErrUnauthorizedAccess):
			ctx.JSON(http.StatusForbidden, utils.ErrorResponse(http.StatusForbidden, err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Ticket scanned successfully",
		nil,
	))
}
