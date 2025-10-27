package handler

import (
	"fmt"
	"lab1_rip/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
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
	// 1. CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173", // твой React dev (vite)
			"http://localhost:3000", // если у тебя CRA/Next или fallback
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,           // нужно true если у тебя куки/сессия
		MaxAge:           12 * time.Hour, // кэш preflight
	}))

	// 2. Группа API как и было
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
	api.GET("/readindxs/my-text-cart", h.API_ReadIndxsCartIcon)
	api.GET("/readindxs", h.API_ReadIndxsList)
	api.GET("/readindxs/:id", h.API_ReadIndxsGet)
	api.PUT("/readindxs/:id", h.API_ReadIndxsUpdate)
	api.PUT("/readindxs/:id/form", h.API_ReadIndxsForm)
	api.PUT("/readindxs/:id/moderate", h.API_ReadIndxsModerate)
	api.DELETE("/readindxs/:id", h.API_ReadIndxsDelete)

	// m-m
	api.PUT("/readindxs-texts/", h.API_ReadIndxsTextsUpdate)
	api.DELETE("/readindxs-texts/", h.API_ReadIndxsTextsDelete)

	// auth / user
	api.POST("/users/register", h.API_UserRegister)
	api.POST("/auth/login", h.API_AuthLogin)
	api.POST("/auth/logout", h.API_AuthLogout)
	api.GET("/users/me", h.API_UserMe)
	api.PUT("/users/me", h.API_UserUpdateMe)
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
