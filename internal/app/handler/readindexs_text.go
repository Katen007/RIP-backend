package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type mmBody struct {
	ReadIndxsID    int  `json:"read_indxs_id" binding:"required"`
	TextID         int  `json:"text_id"      binding:"required"`
	CountWords     *int `json:"count_words"`
	CountSentences *int `json:"count_sentences"`
	CountSyllables *int `json:"count_syllables"`
}

func (h *Handler) API_ReadIndxsTextsUpdate(c *gin.Context) {
	var b mmBody
	if err := c.ShouldBindJSON(&b); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	f := map[string]any{}
	if b.CountWords != nil {
		f["count_words"] = *b.CountWords
	}
	if b.CountSentences != nil {
		f["count_sentences"] = *b.CountSentences
	}
	if b.CountSyllables != nil {
		f["count_syllables"] = *b.CountSyllables
	}
	if err := h.Repository.ReadIndxsTextsUpdate(b.ReadIndxsID, b.TextID, f); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.Status(204)
}

type mmBodyDelete struct {
	ReadIndxsID int `json:"read_indxs_id" binding:"required"`
	TextID      int `json:"text_id"      binding:"required"`
}

func (h *Handler) API_ReadIndxsTextsDelete(c *gin.Context) {
	var b mmBodyDelete
	if err := c.ShouldBindJSON(&b); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	if err := h.Repository.ReadIndxsTextsDelete(b.ReadIndxsID, b.TextID); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.Status(http.StatusNoContent)
}
