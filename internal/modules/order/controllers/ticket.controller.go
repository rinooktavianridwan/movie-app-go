package controllers

import (
	"net/http"
	"strconv"
	"time"

	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/order/responses"
	"movie-app-go/internal/modules/order/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TicketController struct {
	TicketService *services.TicketService
}

func NewTicketController(s *services.TicketService) *TicketController {
	return &TicketController{TicketService: s}
}

func (c *TicketController) GetMyTickets(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	userIDUint := uint(userID.(float64))
	result, err := c.TicketService.GetTicketsByUser(userIDUint, page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	ctx.JSON(http.StatusOK, response)
}

func (c *TicketController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.TicketService.GetAllTickets(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	ctx.JSON(http.StatusOK, response)
}

func (c *TicketController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userID, exists := ctx.Get("user_id")
	isAdmin, adminExists := ctx.Get("is_admin")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var ticket *models.Ticket
	if adminExists && isAdmin.(bool) {
		ticket, err = c.TicketService.GetTicketByID(uint(id), nil)
	} else {
		userIDUint := uint(userID.(float64))
		ticket, err = c.TicketService.GetTicketByID(uint(id), &userIDUint)
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, responses.ToTicketResponse(ticket))
}

func (c *TicketController) GetBySchedule(ctx *gin.Context) {
	scheduleID, err := strconv.Atoi(ctx.Param("schedule_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule_id"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "50"))

	result, err := c.TicketService.GetTicketsBySchedule(uint(scheduleID), page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	ctx.JSON(http.StatusOK, response)
}

func (c *TicketController) ScanTicket(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	userID, exists := ctx.Get("user_id")
	isAdmin, adminExists := ctx.Get("is_admin")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var ticket *models.Ticket
	if adminExists && isAdmin.(bool) {
		// Admin can scan any ticket
		ticket, err = c.TicketService.ScanTicket(uint(id), nil)
	} else {
		// User can only scan their own tickets
		userIDUint := uint(userID.(float64))
		ticket, err = c.TicketService.ScanTicket(uint(id), &userIDUint)
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		} else if err.Error() == "ticket already used" {
			ctx.JSON(http.StatusConflict, gin.H{"error": "ticket already scanned"})
		} else if err.Error() == "cancelled ticket cannot be used" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "ticket is cancelled"})
		} else if err.Error() == "payment not confirmed" {
			ctx.JSON(http.StatusPaymentRequired, gin.H{"error": "payment must be completed first"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "ticket scanned successfully",
		"ticket_id":  ticket.ID,
		"status":     ticket.Status,
		"scanned_at": time.Now(),
		"movie":      ticket.Schedule.Movie.Title,
		"studio":     ticket.Schedule.Studio.Name,
		"seat":       ticket.SeatNumber,
	})
}
