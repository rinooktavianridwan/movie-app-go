package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"movie-app-go/internal/modules/promo/options"
	"movie-app-go/internal/modules/promo/requests"
	"movie-app-go/internal/modules/promo/responses"
	"movie-app-go/internal/modules/promo/services"
	"movie-app-go/internal/utils"

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

func (c *PromoController) Create(ctx *gin.Context) {
	var req requests.CreatePromoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	err := c.PromoService.CreatePromo(&req)
	if err != nil {
		if errors.Is(err, utils.ErrPromoNotFound) {
			ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusCreated, utils.SuccessResponse(
		http.StatusCreated,
		"Promo created successfully",
		nil,
	))
}

func (c *PromoController) GetAllPromos(ctx *gin.Context) {
	opts, err := options.ParsePromoOptions(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	result, err := c.PromoService.GetAllPromosPaginated(opts)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		return
	}

	promoResponses := responses.ToPromoResponses(result.Data)
	response := responses.PaginatedPromoResponse{
		Page:      result.Page,
		PerPage:   result.PerPage,
		Total:     result.Total,
		TotalPage: result.TotalPages,
		Data:      promoResponses,
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Promos retrieved successfully",
		response,
	))
}

func (c *PromoController) GetPromoByID(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	promo, err := c.PromoService.GetPromoByID(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrPromoNotFound) {
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		} else {
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Promo retrieved successfully",
		responses.ToPromoResponse(*promo),
	))
}

func (c *PromoController) Update(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req requests.UpdatePromoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse(err.Error()))
		return
	}

	err := c.PromoService.UpdatePromo(uint(id), &req)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrPromoNotFound):
			ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
		http.StatusOK,
		"Promo updated successfully",
		nil,
	))
}

func (c *PromoController) TogglePromoStatus(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid ID"))
		return
	}

	promo, err := c.PromoService.TogglePromoStatus(uint(id))
	if err != nil {
		if errors.Is(err, utils.ErrPromoNotFound) {
            ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
        } else {
            ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
        }
        return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
        http.StatusOK,
        "Promo status updated successfully",
        responses.ToPromoResponse(*promo),
    ))
}

func (c *PromoController) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := c.PromoService.DeletePromo(uint(id)); err != nil {
		switch {
        case errors.Is(err, utils.ErrPromoNotFound):
            ctx.JSON(http.StatusNotFound, utils.NotFoundResponse(err.Error()))
        case errors.Is(err, utils.ErrPromoInUse):
            ctx.JSON(http.StatusConflict, utils.ErrorResponse(http.StatusConflict, err.Error()))
        default:
            ctx.JSON(http.StatusInternalServerError, utils.InternalServerErrorResponse(err.Error()))
        }
        return
	}

	ctx.JSON(http.StatusOK, utils.SuccessResponse(
        http.StatusOK,
        "Promo deleted successfully",
        nil,
    ))
}
