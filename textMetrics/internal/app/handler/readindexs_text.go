package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// API_ReadIndxsTextsUpdate обновляет метрики текста внутри индекса чтения.
// @Summary      Update text metrics in read index
// @Description  Partial update of word/sentence/syllable counters for specific text in a read index
// @Tags         read-indxs
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body      mmBody  true  "metrics payload"
// @Success      204   {string}  string  "No Content"
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Security BearerAuth
// @Router       /readindxs-texts [put]
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

// API_ReadIndxsTextsDelete удаляет текст из индекса чтения.
// @Summary      Remove text from read index
// @Tags         read-indxs
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body      mmBodyDelete  true  "text delete payload"
// @Success      204   {string}  string  "No Content"
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Security BearerAuth
// @Router       /readindxs-texts [delete]
func (h *Handler) API_ReadIndxsTextsDelete(c *gin.Context) {
	logrus.Info(1)
	b, ok := c.Get("body")
	if !ok {
		h.errorHandler(c, 400, nil)
		return
	}
	body, ok := b.(mmBodyDelete)
	if !ok {
		h.errorHandler(c, 400, nil)
		return
	}
	if err := h.Repository.ReadIndxsTextsDelete(body.ReadIndxsID, body.TextID); err != nil {
		h.errorHandler(c, 404, err)
		return
	}
	c.Status(http.StatusNoContent)
}
