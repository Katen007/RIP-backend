package handler

import (
	"fmt"
	"lab1_rip/internal/app/ds"
	"lab1_rip/internal/app/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// API_ReadIndxsCartIcon возвращает текущий draft и число текстов.
// @Summary      Read indices cart icon info
// @Tags         readindxs
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  CartIconResponse
// @Failure      400  {object}  ErrorResponse
// @Security BearerAuth
// @Router       /readindxs/my-text-cart [get]
func (h *Handler) API_ReadIndxsCartIcon(c *gin.Context) {
	uid := h.GetUserID(c)
	if c.IsAborted() {
		h.errorHandler(c, 400, fmt.Errorf("invalid user id"))
		return
	}
	id, cnt, err := h.Repository.ReadIndxsCartIcon(uid)
	if err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	c.JSON(200, gin.H{"draft_readIndxs_id": id, "texts_count": cnt})
}

// API_ReadIndxsList список индексов с фильтрами.
// @Summary      List read indices
// @Tags         readindxs
// @Produce      json
// @Security     ApiKeyAuth
// @Param        status     query     string false "status filter"
// @Param        date_from  query     string false "YYYY-MM-DD"
// @Param        date_to    query     string false "YYYY-MM-DD"
// @Success      200  {object}  ReadIndxsListResponse
// @Failure      500  {object}  ErrorResponse
// @Security BearerAuth
// @Router       /readindxs [get]
func (h *Handler) API_ReadIndxsList(c *gin.Context) {
	var q struct {
		Status   string `form:"status"`
		DateFrom string `form:"date_from"`
		DateTo   string `form:"date_to"`
	}
	_ = c.ShouldBindQuery(&q)
	var df, dt *time.Time = nil, nil
	if q.DateFrom != "" {
		t, _ := time.Parse("2006-01-02", q.DateFrom)
		df = &t
	}
	if q.DateTo != "" {
		t, _ := time.Parse("2006-01-02", q.DateTo)
		dt = &t
	}
	logrus.Info(df, dt, q.DateFrom)
	user, err := h.GetUserDTO(c)
	if err != nil {
		h.errorHandler(c, 401, err)
		return
	}
	data, err := h.Repository.ReadIndxsList(user, repository.ReadIndxsFilter{Status: q.Status, DateFrom: df, DateTo: dt})

	if err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	dto := ds.ToReadIndxsListDTO(data)
	c.JSON(200, dto)
}

// API_ReadIndxsGet получить индекс по id.
// @Summary      Get read index by id
// @Tags         readindxs
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path  int  true  "ReadIndxs ID"
// @Success      200  {object}  ReadIndxsInfoResponse
// @Failure      404  {object}  ErrorResponse
// @Security BearerAuth
// @Router       /readindxs/{id} [get]
func (h *Handler) API_ReadIndxsGet(c *gin.Context) {
	id := mustIntParam(c, "id")
	data, err := h.Repository.ReadIndxsGet(id)
	if err != nil {
		h.errorHandler(c, 404, err)
		return
	}
	dto := ds.ToReadIndxsInfoDTO(data)
	c.JSON(200, dto)
}

// API_ReadIndxsUpdate частичное обновление тематик/полей индекса.
// @Summary      Update read index (partial)
// @Tags         readindxs
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id    path  int   true  "ReadIndxs ID"
// @Param        body  body  map[string]any true "fields to update"
// @Success      204   {string}  string  "No Content"
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Security BearerAuth
// @Router       /readindxs/{id} [patch]
func (h *Handler) API_ReadIndxsUpdate(c *gin.Context) {
	id := mustIntParam(c, "id")
	var body map[string]any
	if err := c.BindJSON(&body); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	if err := h.Repository.ReadIndxsUpdateThematic(id, body); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.Status(204)
}

// API_ReadIndxsForm формирует индекс (server-side операция).
// @Summary      Build/form read index
// @Tags         readindxs
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path  int  true  "ReadIndxs ID"
// @Success      204  {string}  string "No Content"
// @Failure      400  {object}  ErrorResponse
// @Security BearerAuth
// @Router       /readindxs/{id}/form [post]
func (h *Handler) API_ReadIndxsForm(c *gin.Context) {
	id := mustIntParam(c, "id")
	if err := h.Repository.ReadIndxsForm(id); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	c.Status(204)
}

// API_ReadIndxsModerate модерация индекса.
// @Summary      Moderate read index
// @Tags         readindxs
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id    path  int  true  "ReadIndxs ID"
// @Param        body  body  ReadIndxsModerateRequest true "action payload"
// @Success      200  {object}  ReadIndxsModerateResponse
// @Failure      400  {object}  ErrorResponse
// @Security BearerAuth
// @Router       /readindxs/{id}/moderate [post]
func (h *Handler) API_ReadIndxsModerate(c *gin.Context) {
	id := mustIntParam(c, "id")
	uid := h.GetUserID(c)
	var body struct {
		Action string `json:"action"`
	}
	if err := c.BindJSON(&body); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	data, err := h.Repository.ReadIndxsModerate(id, uid, body.Action)
	if err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	c.JSON(200, data)
}

// API_ReadIndxsDelete мягкое удаление индекса.
// @Summary      Soft delete read index
// @Tags         readindxs
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path  int  true  "ReadIndxs ID"
// @Success      204  {string}  string "No Content"
// @Failure      500  {object}  ErrorResponse
// @Security BearerAuth
// @Router       /readindxs/{id} [delete]
func (h *Handler) API_ReadIndxsDelete(c *gin.Context) {
	id := mustIntParam(c, "id")
	if err := h.Repository.ReadIndxsSoftDelete(id); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.Status(204)
}
