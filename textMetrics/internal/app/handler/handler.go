package handler

import (
	"fmt"
	"lab1_rip/internal/app/config"
	"lab1_rip/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
	Config     *config.Config
	Redis      *repository.Redis
}

func NewHandler(r *repository.Repository, c *config.Config, redis *repository.Redis) *Handler {
	return &Handler{
		Repository: r,
		Config:     c,
		Redis:      redis,
	}
}

func (h *Handler) RegisterAPI(r *gin.Engine) {
	api := r.Group("/api")

	// Text (услуги)
	api.GET("/texts", h.API_TextsList)
	api.GET("/texts/:id", h.API_TextGet)
	api.POST("/texts", h.API_TextCreate)
	api.PUT("/texts/:id", h.ModeratorValidateMiddleware(), h.API_TextUpdate)
	api.DELETE("/texts/:id", h.ModeratorValidateMiddleware(), h.API_TextDelete)
	api.POST("/texts/:id/image", h.ModeratorValidateMiddleware(), h.API_TextUploadImage)
	api.POST("/texts/:id/add-to-draft", h.AuthoMiddleware(), h.API_TextAddToDraft)

	// ReadIndxs (заявки)
	api.GET("/readindxs/my-text-cart", h.AuthoMiddleware(), h.API_ReadIndxsCartIcon)
	api.GET("/readindxs", h.AuthoMiddleware(), h.API_ReadIndxsList)
	api.GET("/readindxs/:id", h.AuthoMiddleware(), h.ReadIndxsAccessMiddleware(), h.API_ReadIndxsGet)
	api.PUT("/readindxs/:id", h.AuthoMiddleware(), h.ReadIndxsAccessMiddleware(), h.API_ReadIndxsUpdate)
	api.PUT("/readindxs/update-calc", h.API_ReadIndxsUpdateCalcs)
	api.PUT("/readindxs/:id/form", h.AuthoMiddleware(), h.ReadIndxsAccessMiddleware(), h.API_ReadIndxsForm)
	api.PUT("/readindxs/:id/moderate", h.API_ReadIndxsModerate) //h.ModeratorValidateMiddleware(), h.API_ReadIndxsModerate)
	api.DELETE("/readindxs/:id", h.AuthoMiddleware(), h.ReadIndxsAccessMiddleware(), h.API_ReadIndxsDelete)

	// m-m
	api.PUT("/readindxs-texts", h.AuthoMiddleware(), h.ReadIndxsToTextsAccessMiddleware(), h.API_ReadIndxsTextsUpdate)
	api.DELETE("/readindxs-texts", h.AuthoMiddleware(), h.ReadIndxsToTextsAccessMiddleware(), h.API_ReadIndxsTextsDelete)

	api.POST("/users/register", h.API_UserRegister)
	api.POST("/auth/login", h.API_AuthLogin)
	api.POST("/auth/logout", h.AuthoMiddleware(), h.API_AuthLogout) // можно и без мидлвари
	api.GET("/users/me", h.AuthoMiddleware(), h.API_UserMe)
	api.PUT("/users/me", h.AuthoMiddleware(), h.API_UserUpdateMe)
}

// errorHandler для более удобного вывода ошибок
func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
	ctx.Abort()
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

func (h *Handler) fileParams(c *gin.Context) (int64, string) {
	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	contentLengthStr := c.Request.Header.Get("Content-Length")
	if contentLengthStr == "" {
		h.errorHandler(c, http.StatusBadRequest, fmt.Errorf("content-Length header is required"))
		return 0, ""
	}
	fileSize, err := strconv.ParseInt(contentLengthStr, 10, 64)
	if err != nil {
		h.errorHandler(c, http.StatusBadRequest, fmt.Errorf("invalid Content-Length header"))
		return 0, ""
	}
	return fileSize, contentType
}
