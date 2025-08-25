package controllers

import (
	"net/http"
	"strconv"

	"movie-app-go/internal/modules/promo/options"
	"movie-app-go/internal/modules/promo/requests"
	"movie-app-go/internal/modules/promo/responses"
	"movie-app-go/internal/modules/promo/services"

	"github.com/gin-gonic/gin"
)

type PromoController struct {
	PromoService *services.PromoService
}

func NewPromoController(promoService *services.PromoService) *PromoController {
	return &PromoController{
		PromoService: promoService,
	}
}

func (c *PromoController) CreatePromo(ctx *gin.Context) {
	var req requests.CreatePromoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed"})
		return
	}

	promo, err := c.PromoService.CreatePromo(&req)
	if err != nil {
		if err.Error() == "promo code already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else if err.Error() == "some movie_ids are invalid" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	promoResponse := responses.ToPromoResponse(*promo)
	ctx.JSON(http.StatusCreated, promoResponse)
}

func (c *PromoController) GetAllPromos(ctx *gin.Context) {
	opts, err := options.ParsePromoOptions(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := c.PromoService.GetAllPromosPaginated(opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	promoResponses := responses.ToPromoResponses(result.Data)
	paginatedResponse := responses.PaginatedPromoResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      promoResponses,
	}

	ctx.JSON(http.StatusOK, paginatedResponse)
}

func (c *PromoController) GetPromoByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	promo, err := c.PromoService.GetPromoByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "promo not found"})
		return
	}

	promoResponse := responses.ToPromoResponse(*promo)
	ctx.JSON(http.StatusOK, promoResponse)
}

func (c *PromoController) UpdatePromo(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req requests.UpdatePromoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "validation failed"})
		return
	}

	promo, err := c.PromoService.UpdatePromo(uint(id), &req)
	if err != nil {
		if err.Error() == "promo not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "some movie_ids are invalid" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	promoResponse := responses.ToPromoResponse(*promo)
	ctx.JSON(http.StatusOK, promoResponse)
}

func (c *PromoController) TogglePromoStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	promo, err := c.PromoService.TogglePromoStatus(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "promo not found"})
		return
	}

	promoResponse := responses.ToPromoResponse(*promo)
	ctx.JSON(http.StatusOK, gin.H{
		"message": "promo status updated successfully",
		"promo":   promoResponse,
	})
}

func (c *PromoController) DeletePromo(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := c.PromoService.DeletePromo(uint(id)); err != nil {
		if err.Error() == "promo not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "promo deleted successfully"})
}
