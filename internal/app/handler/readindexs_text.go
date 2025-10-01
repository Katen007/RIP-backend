package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type mmBody struct {
	ReadIndxsID int     `json:"readIndxsId" binding:"required"`
	TextID      int     `json:"textId"      binding:"required"`
	Calculation *int    `json:"calculation"`
	Formula     *string `json:"formula"`
}

func (h *Handler) API_ReadIndxsTextsUpdate(c *gin.Context) {
	var b mmBody
	if err := c.BindJSON(&b); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	m := map[string]any{}
	if b.Calculation != nil {
		m["calculation"] = *b.Calculation
	}
	if b.Formula != nil {
		m["formula"] = *b.Formula
	}
	if err := h.Repository.ReadIndxsTextsUpdate(b.ReadIndxsID, b.TextID, m); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) API_ReadIndxsTextsDelete(c *gin.Context) {
	var b mmBody
	if err := c.BindJSON(&b); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	if err := h.Repository.ReadIndxsTextsDelete(b.ReadIndxsID, b.TextID); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.Status(http.StatusNoContent)
}
