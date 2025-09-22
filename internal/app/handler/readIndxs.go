package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetReadIndxs(ctx *gin.Context) {
	//ids := getCartCopy()
	// items := make([]repository.Order, 0, len(ids))
	total := 0
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}
	readIndex, err := h.Repository.GetReadIndxsById(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	items, _ := h.Repository.GetTextsInReadIndxs(id)

	ctx.HTML(http.StatusOK, "readIndxs.html", gin.H{
		"readIndx": readIndex,
		"total":    total,
		"items":    items,
	})
}

func (h *Handler) AddTextToReadIndxs(ctx *gin.Context) {
	tId := ctx.PostForm("text_id")
	textId, err := strconv.Atoi(tId)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info(textId)

	if err := h.Repository.AddTextToReadIndxs(textId, 1); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/texts")
}

func (h *Handler) DeleteReadIndexs(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	if err := h.Repository.DeleteReadIndxs(id); err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/texts")
}
