package handler

import (
	"lab1_rip/internal/app/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// API_UserRegister godoc
// @Summary      Регистрация пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      UserCredentials  true  "Регистрация"
// @Success      201   {object}  ds.User
// @Failure      400   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}  "login already taken"
// @Failure      500   {object}  map[string]interface{}
// @Router       /users/register [post]
func (h *Handler) API_UserRegister(c *gin.Context) {
	var in struct {
		Login       string `json:"login" binding:"required,min=3,max=25"`
		Password    string `json:"password" binding:"required,min=4,max=64"`
		IsModerator bool   `json:"isModerator"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		h.errorHandler(c, 500, err)
		return
	}

	u, err := h.Repository.UserCreate(in.Login, string(hash), in.IsModerator)
	if err != nil {
		if err == repository.ErrLoginTaken {
			c.JSON(http.StatusConflict, gin.H{"status": "error", "description": "login already taken"})
			return
		}
		h.errorHandler(c, 500, err)
		return
	}

	// опционально сразу логинить:
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "login": u.Login, "isModerator": u.IsModerator})
}

// API_AuthLogin godoc
// @Summary      Логин
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      UserCredentials  true  "Логин"
// @Success      200   {object}  map[string]bool  "ok=true"
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /auth/login [post]
func (h *Handler) API_AuthLogin(c *gin.Context) {
	var in struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}
	u, err := h.Repository.UserByLogin(in.Login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "description": "invalid credentials"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(in.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "description": "invalid credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// API_AuthLogout godoc
// @Summary      Логаут
// @Tags         auth
// @Produce      json
// @Param        body  body      nil  false  "Логаут"
// @Success      200  {object}  map[string]bool  "ok=true"
// @Router       /auth/logout [post]
func (h *Handler) API_AuthLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// API_UserMe godoc
// @Summary      Текущий пользователь
// @Tags         users
// @Produce      json
// @Success      200  {object}  ds.User
// @Failure      401  {object}  map[string]interface{}
// @Router       /users/me [get]
func (h *Handler) API_UserMe(c *gin.Context) {
	uid_, ok := c.Get("user_id")
	uid := uid_.(int)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "description": "unauthorized"})
		return
	}
	u, err := h.Repository.UserByID(uid)
	if err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": u.ID, "login": u.Login, "isModerator": u.IsModerator,
	})
}

// API_UserUpdateMe godoc
// @Summary      Обновить свои данные
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body      UserCredentials  true  "Новые поля"
// @Success      200   {object}  ds.User
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /users/me [patch]
func (h *Handler) API_UserUpdateMe(c *gin.Context) {
	uid_, ok := c.Get("user_id")
	uid := uid_.(int)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "description": "unauthorized"})
		return
	}
	var in struct {
		Login    *string `json:"login"`
		Password *string `json:"password"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}

	fields := map[string]any{}
	if in.Login != nil && *in.Login != "" {
		fields["login"] = *in.Login
	}
	if in.Password != nil && *in.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*in.Password), bcrypt.DefaultCost)
		if err != nil {
			h.errorHandler(c, 500, err)
			return
		}
		fields["hashed_password"] = string(hash)
	}
	u, err := h.Repository.UserUpdateMe(uid, fields)
	if err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": u.ID, "login": u.Login, "isModerator": u.IsModerator})
}
