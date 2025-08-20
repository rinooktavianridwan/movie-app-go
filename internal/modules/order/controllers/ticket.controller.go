package controllers

import (
	"net/http"
	"strconv"

	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/order/requests"
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

    result, err := c.TicketService.GetTicketsByUser(userID.(uint), page, perPage)
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
        userIDUint := userID.(uint)
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

func (c *TicketController) UpdateStatus(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var req requests.UpdateTicketRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
        ticket, err = c.TicketService.UpdateTicketStatus(uint(id), &req, nil)
    } else {
        userIDUint := userID.(uint)
        ticket, err = c.TicketService.UpdateTicketStatus(uint(id), &req, &userIDUint)
    }

    if err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
        } else {
            ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "ticket status updated",
        "ticket_id": ticket.ID,
        "status": ticket.Status,
    })
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