package handler

import (
	"lab1_rip/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterAPI(r *gin.Engine) {
	api := r.Group("/api", LabFixedUser())

	// Text (услуги)
	api.GET("/texts", h.API_TextsList)
	api.GET("/texts/:id", h.API_TextGet)
	api.POST("/texts", h.API_TextCreate)
	api.PUT("/texts/:id", h.API_TextUpdate)
	api.DELETE("/texts/:id", h.API_TextDelete)
	api.POST("/texts/:id/image", h.API_TextUploadImage)
	api.POST("/texts/:id/add-to-draft", h.API_TextAddToDraft)

	// ReadIndxs (заявки)
	api.GET("/readindxs/cart-icon", h.API_ReadIndxsCartIcon)
	api.GET("/readindxs", h.API_ReadIndxsList)
	api.GET("/readindxs/:id", h.API_ReadIndxsGet)
	api.PUT("/readindxs/:id", h.API_ReadIndxsUpdate)
	api.PUT("/readindxs/:id/form", h.API_ReadIndxsForm)
	api.PUT("/readindxs/:id/moderate", h.API_ReadIndxsModerate)
	api.DELETE("/readindxs/:id", h.API_ReadIndxsDelete)

	// m-m
	api.PUT("/readindxs-texts", h.API_ReadIndxsTextsUpdate)
	api.DELETE("/readindxs-texts", h.API_ReadIndxsTextsDelete)
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
func mustIntParam(c *gin.Context, name string) int {
	raw := c.Param(name)
	id, err := strconv.Atoi(raw)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":      "error",
			"description": "invalid '" + name + "' path param: expected integer",
		})
		return 0
	}
	return id
}
