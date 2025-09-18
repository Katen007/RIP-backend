package handler

import (
	"lab1_rip/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//	func addToCartID(id int) {
//		cartMu.Lock()
//		cartIDs = append(cartIDs, id)
//		cartMu.Unlock()
//	}
//
//	func removeOneFromCart(id int) {
//		cartMu.Lock()
//		for i, v := range cartIDs {
//			if v == id {
//				cartIDs = append(cartIDs[:i], cartIDs[i+1:]...)
//				break
//			}
//		}
//		cartMu.Unlock()
//	}

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetTexts(ctx *gin.Context) {
	var texts []repository.Text
	var err error
	readIndxsId := 1
	items, _ := h.Repository.GetReadIndxsComponents(1)
	components := items.Texts
	count := len(components)

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

// func (h *Handler) AddToCart(ctx *gin.Context) {
// 	id, err := strconv.Atoi(ctx.Param("id"))
// 	if err == nil {
// 		addToCartID(id)
// 	}
// 	ref := ctx.Request.Referer()
// 	if ref == "" {
// 		ref = "/hello"
// 	}
// 	ctx.Redirect(http.StatusFound, ref)
// }

//	func (h *Handler) RemoveFromCart(ctx *gin.Context) {
//		id, err := strconv.Atoi(ctx.Param("id"))
//		if err == nil {
//			removeOneFromCart(id)
//		}
//		ctx.Redirect(http.StatusFound, "/cart")
//	}

func (h *Handler) GetReadIndxs(ctx *gin.Context) {
	//ids := getCartCopy()
	// items := make([]repository.Order, 0, len(ids))
	total := 0
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	// for _, id := range ids {
	// 	o, err := h.Repository.GetOrder(id)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	// items = append(items, o)
	// 	total += o.Price
	// }S
	items, _ := h.Repository.GetReadIndxsComponents(id)
	itemsCalc := items.Texts

	ctx.HTML(http.StatusOK, "readIndxs.html", gin.H{
		"readIndxs": items,
		"total":     total,
		"items":     itemsCalc,
	})
}
