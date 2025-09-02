package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/order/requests"
	"movie-app-go/internal/modules/order/responses"
	"movie-app-go/internal/modules/order/services"
	"movie-app-go/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionController struct {
	TransactionService *services.TransactionService
}

func NewTransactionController(s *services.TransactionService) *TransactionController {
	return &TransactionController{TransactionService: s}
}

func (c *TransactionController) Create(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")

	var req requests.CreateTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	userIDUint := uint(userID.(float64))
	err := c.TransactionService.CreateTransaction(userIDUint, &req)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrScheduleNotFound):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		case errors.Is(err, utils.ErrInsufficientSeats):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		case errors.Is(err, utils.ErrSeatAlreadyBooked):
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, err.Error()))
		case errors.Is(err, utils.ErrPromoNotFound):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(
		http.StatusCreated,
		"Transaction created successfully",
		nil,
	))
}

func (c *TransactionController) GetMyTransactions(ctx *gin.Context) {
	userID, _ := ctx.Get("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	userIDUint := uint(userID.(float64))
	result, err := c.TransactionService.GetTransactionsByUser(userIDUint, page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
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

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"User transactions retrieved successfully",
		response,
	))
}

func (c *TransactionController) GetAll(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "10"))

	result, err := c.TransactionService.GetAllTransactions(page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
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

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Transactions retrieved successfully",
		response,
	))
}

func (c *TransactionController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("invalid id"))
		return
	}

	userID, exists := ctx.Get("user_id")
	isAdmin, adminExists := ctx.Get("is_admin")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.UnauthorizedResponse("unauthenticated"))
		return
	}

	var transaction *models.Transaction
	if adminExists && isAdmin.(bool) {
		transaction, err = c.TransactionService.GetTransactionByID(uint(id), nil)
	} else {
		userIDUint := uint(userID.(float64))
		transaction, err = c.TransactionService.GetTransactionByID(uint(id), &userIDUint)
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse("transaction not found"))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Transaction retrieved successfully",
		responses.ToTransactionResponse(transaction),
	))
}

func (c *TransactionController) ProcessPayment(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("invalid transaction id"))
		return
	}

	var req requests.ProcessPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("invalid request payload"))
		return
	}

	err = c.TransactionService.ProcessPayment(uint(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrTransactionNotFound):
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		case errors.Is(err, utils.ErrTransactionAlreadyPaid):
			ctx.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, err.Error()))
		case errors.Is(err, utils.ErrTransactionExpired):
			ctx.JSON(http.StatusGone, utils.ErrorResponse(http.StatusGone, err.Error()))
		case errors.Is(err, utils.ErrPaymentProcessingFailed):
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Payment processed successfully",
		nil,
	))
}
