package handler

import (
	"fmt"
	"lab1_rip/internal/app/repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) API_ReadIndxsCartIcon(c *gin.Context) {
	uid := getUserID(c)
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
	data, err := h.Repository.ReadIndxsList(repository.ReadIndxsFilter{Status: q.Status, DateFrom: df, DateTo: dt})
	if err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.JSON(200, data)
}

func (h *Handler) API_ReadIndxsGet(c *gin.Context) {
	id := mustIntParam(c, "id")
	data, err := h.Repository.ReadIndxsGet(id)
	if err != nil {
		h.errorHandler(c, 404, err)
		return
	}
	c.JSON(200, data)
}

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

func (h *Handler) API_ReadIndxsForm(c *gin.Context) {
	id := mustIntParam(c, "id")
	if err := h.Repository.ReadIndxsForm(id); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	c.Status(204)
}

func (h *Handler) API_ReadIndxsModerate(c *gin.Context) {
	id := mustIntParam(c, "id")
	uid := getUserID(c)
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

func (h *Handler) API_ReadIndxsDelete(c *gin.Context) {
	id := mustIntParam(c, "id")
	if err := h.Repository.ReadIndxsSoftDelete(id); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.Status(204)
}
