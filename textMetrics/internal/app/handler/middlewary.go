package handler

import (
	"bytes"
	"fmt"
	"io"
	"lab1_rip/internal/app/ds"
	jwtutils "lab1_rip/internal/pkg/jwtUtils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const CtxUserID = "user_id"

func (h *Handler) AuthoMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtTokenStr, err := GetTokenInHeader(ctx.Request.Header)

		if err != nil {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
			return
		}
		logrus.Printf("token : %s", jwtTokenStr)
		if err := h.Redis.GetBlackListJWT(ctx, jwtTokenStr); err == nil {
			h.errorHandler(ctx, http.StatusUnauthorized, fmt.Errorf("unauthorized: token in blacklist"))
			return
		}
		user, err := h.GetJWT(jwtTokenStr)
		if err != nil {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
			return
		}
		logrus.Printf("user1 : %v", user.IsModerator)
		ctx.Set("user", user)
		ctx.Set("token", jwtTokenStr)
		ctx.Set(CtxUserID, user.ID)
		ctx.Next()
	}
}

func (h *Handler) ModeratorValidateMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtTokenStr, err := GetTokenInHeader(ctx.Request.Header)
		if err != nil {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
			return
		}
		user, err := h.GetJWT(jwtTokenStr)
		if err != nil {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
			return
		}
		if !user.IsModerator {
			h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("not allowed role"))
			return
		}
		ctx.Set("user", user)
		ctx.Set("token", jwtTokenStr)
		ctx.Next()
	}
}
func (h *Handler) GetUserID(ctx *gin.Context) int {
	dto, err := h.GetUserDTO(ctx)
	if err != nil {
		h.errorHandler(ctx, 404, err)
		return 0
	}
	return int(dto.ID)
}

func (h *Handler) GetUserDTO(ctx *gin.Context) (ds.User, error) {
	userRaw, ok := ctx.Get("user")
	if !ok {
		return ds.User{}, fmt.Errorf("user not found")
	}
	user, ok := userRaw.(ds.User)
	if !ok {
		return ds.User{}, fmt.Errorf("invalid user id type: expected int, got %T", user)
	}
	return user, nil
}

func (h *Handler) GetJWT(jwtTokenStr string) (ds.User, error) {
	claims, err := jwtutils.ValidateJwtToken(jwtTokenStr, h.Config.SecretKey)
	if err != nil {
		return ds.User{}, err
	}
	userId, ok := claims["sub"].(float64)
	logrus.Info(userId)
	if !ok {
		return ds.User{}, fmt.Errorf("bad jwt credentials")
	}
	is_moderator, ok := claims["is_moderator"].(bool)
	logrus.Info(is_moderator)
	if !ok {
		return ds.User{}, fmt.Errorf("bad jwt credentials")
	}
	login, ok := claims["login"].(string)
	if !ok {
		return ds.User{}, fmt.Errorf("bad jwt credentials")
	}
	user := ds.User{ID: int(userId), IsModerator: is_moderator, Login: login}
	return user, nil

}

func GetTokenInHeader(header http.Header) (string, error) {
	jwtPrefix := "Bearer "
	jwtStr := header.Get("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) {
		return "", fmt.Errorf("access token not found")
	}
	jwtTokenStr := jwtStr[len(jwtPrefix):]
	return jwtTokenStr, nil
}
func (h *Handler) ReadIndxsAccessMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := h.GetUserDTO(ctx)
		if err != nil {
			h.errorHandler(ctx, http.StatusNotFound, err)
			return
		}
		if user.IsModerator {
			ctx.Next()
			return
		}
		id := mustIntParam(ctx, "id")
		if ctx.IsAborted() {
			return
		}
		readIndxs, err := h.Repository.GetReadIndxsById(id)
		if err != nil {
			h.errorHandler(ctx, http.StatusNotFound, err)
			return
		}
		if readIndxs.CreatorId != user.ID {
			h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("access denied"))
			return
		}
		ctx.Next()
	}
}

func (h *Handler) ReadIndxsToTextsAccessMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := h.GetUserDTO(ctx)
		if err != nil {
			h.errorHandler(ctx, http.StatusNotFound, err)
			return
		}

		// прочитаем raw body и восстановим позже
		data, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		// восстановим поток для следующих хендлеров
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		var b mmBodyDelete
		if err := ctx.ShouldBindJSON(&b); err != nil {
			h.errorHandler(ctx, 400, err)
			return
		}

		// ещё раз восстановим body (на всякий случай) чтобы downstream мог читать
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		ctx.Set("body", b)

		readIndxs, err := h.Repository.ReadIndxsTextsGet(b.ReadIndxsID, b.TextID)
		if err != nil {
			h.errorHandler(ctx, 500, err)
			return
		}
		if readIndxs.ReadIndxs.CreatorId != user.ID {
			h.errorHandler(ctx, http.StatusForbidden, fmt.Errorf("access denied"))
			return
		}
		ctx.Next()
	}
}

// ...existing code...
