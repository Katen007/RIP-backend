package handler

import (
	"errors"
	"fmt"
	"lab1_rip/internal/app/models"
	"lab1_rip/internal/app/repository"

	jwtutils "lab1_rip/internal/pkg/jwtUtils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Response[T any] struct {
	Ok   bool `json:"ok"`
	Data T    `json:"data"`
}

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
// @Success      200   {object}  Response[models.AuthoResp]
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /auth/login [post]
func (h *Handler) API_AuthLogin(c *gin.Context) {
	var credentials struct {
		Login    string `json:"login" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&credentials); err != nil {
		h.errorHandler(c, http.StatusBadRequest, err)
		return
	}
	password := credentials.Password
	user, err := h.Repository.UserByLogin(credentials.Login)
	if err != nil {
		h.errorHandler(c, 401, err)
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		h.errorHandler(c, 401, err)
		return
	}
	claims := jwt.MapClaims{}
	duration := h.Config.ExpiresAtMinutes
	exp := time.Now().Add(duration)
	logrus.Info(time.Unix(int64(exp.Unix()), 0).Format(time.ANSIC))
	claims["sub"] = user.ID
	claims["exp"] = exp.Unix()
	claims["login"] = user.Login
	claims["is_moderator"] = user.IsModerator
	claims["token_type"] = "access"

	tokenStr, err := jwtutils.CreateJwtToken(claims, h.Config.SecretKey)
	if err != nil {
		h.errorHandler(c, 401, err)
		return
	}

	dto := models.AuthoResp{
		TokenType:   "access",
		ExpiresIn:   exp.Format(time.ANSIC),
		AccessToken: tokenStr,
	}
	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"data": dto,
	})
}

// API_AuthLogout godoc
// @Summary      Логаут
// @Tags         auth
// @Produce      json
// @Param        body  body      nil  false  "Логаут"
// @Success      200  {object}  map[string]bool  "ok=true"
// @Router       /auth/logout [post]
// @Security BearerAuth
func (h *Handler) API_AuthLogout(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		h.errorHandler(c, 401, nil)
		return
	}
	tokenStr, ok := token.(string)
	if !ok {
		h.errorHandler(c, 400, nil)
		return
	}
	err := h.Redis.SetBlackListJWT(c, tokenStr, h.Config.ExpiresAtMinutes)
	if err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// API_UserMe godoc
// @Summary      Текущий пользователь
// @Tags         users
// @Produce      json
// @Success      200  {object}  ds.User
// @Failure      401  {object}  map[string]interface{}
// @Security BearerAuth
// @Router       /users/me [get]
func (h *Handler) API_UserMe(c *gin.Context) {
	user, err := h.GetUserDTO(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "description": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": user.ID, "login": user.Login, "isModerator": user.IsModerator,
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
// @Security BearerAuth
// @Router       /users/me [put]
func (h *Handler) API_UserUpdateMe(c *gin.Context) {
	uid_, ok := c.Get(CtxUserID)
	//uid := uid_.(int)
	if !ok || uid_ == nil {
		h.errorHandler(c, 401, errors.New("unauthorized: uid not found in context"))
		//c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "description": "unauthorized"})
		return
	}
	fmt.Println("uid_", uid_)
	uid, ok := uid_.(int)
	if !ok {
		h.errorHandler(c, 500, errors.New("invalid uid type in context"))
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
