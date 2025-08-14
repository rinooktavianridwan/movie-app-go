package controllers

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "movie-app-go/internal/modules/studio/services"
    "movie-app-go/internal/modules/studio/requests"
)

type StudioController struct {
    StudioService *services.StudioService
}

func NewStudioController(s *services.StudioService) *StudioController {
    return &StudioController{StudioService: s}
}

func (c *StudioController) Create(ctx *gin.Context) {
    var req requests.CreateStudioRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    studio, err := c.StudioService.CreateStudio(&req)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, studio)
}

func (c *StudioController) GetAll(ctx *gin.Context) {
    studios, err := c.StudioService.GetAllStudios()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, studios)
}

func (c *StudioController) GetByID(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    studio, facilityIDs, err := c.StudioService.GetStudioByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "studio not found"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{
        "id":           studio.ID,
        "name":         studio.Name,
        "seat_capacity": studio.SeatCapacity,
        "facility_ids": facilityIDs,
    })
}

func (c *StudioController) Update(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    var req requests.CreateStudioRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    studio, err := c.StudioService.UpdateStudio(uint(id), &req)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "studio not found"})
        return
    }
    ctx.JSON(http.StatusOK, studio)
}

func (c *StudioController) Delete(ctx *gin.Context) {
    id, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }
    if err := c.StudioService.DeleteStudio(uint(id)); err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "studio not found"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "studio deleted"})
}