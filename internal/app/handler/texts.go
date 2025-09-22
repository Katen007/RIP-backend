package handler

import (
	"lab1_rip/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetTexts(ctx *gin.Context) {
	var texts []ds.Text
	var err error
	var readIndxsId int = -1
	readIndxs, err := h.Repository.GetReadIndxs(1)
	if err != nil {
		logrus.Error(err)
	} else {
		readIndxsId = readIndxs.ID
	}
	count, err := h.Repository.GetCountTexts(1)
	if err != nil {
		logrus.Error(err)
	}

	searchQuery := ctx.Query("searchTexts")
	// addItem := ctx.Query("addItem")
	// if addItem != "" {
	// 	count += 1
	// }
	if searchQuery == "" { // если поле поиска пусто, то просто получаем из репозитория все записи
		texts, err = h.Repository.GetTexts()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		texts, err = h.Repository.GetTextsByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "texts.html", gin.H{

		"texts":          texts,
		"searchTexts":    searchQuery,
		"readIndxsCount": count,
		"readIndxsId":    readIndxsId,
	})
}
func (h *Handler) GetText(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	text, err := h.Repository.GetText(id)

	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "text.html", gin.H{
		"text": text,
	})
}
