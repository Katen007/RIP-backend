package handler

import (
	"lab1_rip/internal/app/repository"

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

func (h *Handler) RegisterHandler(r *gin.Engine) {
	r.GET("/texts", h.GetTexts)
	r.GET("/texts/:id", h.GetText)
	r.GET("/readIndxs/:id", h.GetReadIndxs)
	r.POST("/texts", h.AddTextToReadIndxs)
	r.POST("/readIndxs", h.DeleteReadIndexs)
}

// RegisterStatic То же самое, что и с маршрутами, регистрируем статику
func (h *Handler) RegisterStatic(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
