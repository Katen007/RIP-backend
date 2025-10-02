package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type mmBody struct {
	Calculation *int    `json:"calculation"`
	Formula     *string `json:"formula"`
}

func (h *Handler) API_ReadIndxsTextsUpdate(c *gin.Context) {
	var read_indx_id, text_id int
	read_indx_id = mustIntParam(c, "id")
	if c.IsAborted() {
		return
	}
	text_id = mustIntParam(c, "text_id")
	if c.IsAborted() {
		return
	}
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
	if err := h.Repository.ReadIndxsTextsUpdate(read_indx_id, text_id, m); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) API_ReadIndxsTextsDelete(c *gin.Context) {
	var read_indx_id, text_id int
	read_indx_id = mustIntParam(c, "id")
	if c.IsAborted() {
		return
	}
	text_id = mustIntParam(c, "text_id")
	if c.IsAborted() {
		return
	}
	if err := h.Repository.ReadIndxsTextsDelete(read_indx_id, text_id); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.Status(http.StatusNoContent)
}
