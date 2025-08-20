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

type TransactionController struct {
	TransactionService *services.TransactionService
}

func NewTransactionController(s *services.TransactionService) *TransactionController {
	return &TransactionController{TransactionService: s}
}

// Create Transaction (User creates booking)
func (c *TransactionController) Create(ctx *gin.Context) {
	// Get user ID from JWT token
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req requests.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := c.TransactionService.CreateTransaction(userID.(uint), &req)
	if err != nil {
		if err.Error() == "schedule not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "some seats are already booked" ||
			err.Error() == "seat numbers exceed studio capacity" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":        "transaction created successfully",
		"transaction_id": transaction.ID,
		"total_amount":   transaction.TotalAmount,
		"payment_status": transaction.PaymentStatus,
	})
}

// Get My Transactions (User gets their own transactions)
func (c *TransactionController) GetMyTransactions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.TransactionService.GetTransactionsByUser(userID.(uint), page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transactionResponses := responses.ToTransactionResponses(result.Data)
	response := responses.PaginatedTransactionResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      transactionResponses,
	}

	ctx.JSON(http.StatusOK, response)
}

// Get All Transactions (Admin only)
func (c *TransactionController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.TransactionService.GetAllTransactions(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transactionResponses := responses.ToTransactionResponses(result.Data)
	response := responses.PaginatedTransactionResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      transactionResponses,
	}

	ctx.JSON(http.StatusOK, response)
}

// Get Transaction by ID
func (c *TransactionController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Check if user is admin or owner
	userID, exists := ctx.Get("user_id")
	isAdmin, adminExists := ctx.Get("is_admin")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var transaction *models.Transaction
	if adminExists && isAdmin.(bool) {
		// Admin can see all transactions
		transaction, err = c.TransactionService.GetTransactionByID(uint(id), nil)
	} else {
		// User can only see their own transactions
		userIDUint := userID.(uint)
		transaction, err = c.TransactionService.GetTransactionByID(uint(id), &userIDUint)
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, responses.ToTransactionResponse(transaction))
}

// Update Transaction Status (Admin only - for payment processing)
func (c *TransactionController) UpdateStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req requests.UpdateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := c.TransactionService.UpdateTransactionStatus(uint(id), &req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":        "transaction status updated",
		"transaction_id": transaction.ID,
		"payment_status": transaction.PaymentStatus,
	})
}
