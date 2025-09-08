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

var count int = 0

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}
func (h *Handler) GetOrders(ctx *gin.Context) {
	var orders []repository.Order
	var err error
	cartId := 1

	searchQuery := ctx.Query("query")
	addItem := ctx.Query("addItem")
	if addItem != "" {
		count += 1
	}
	if searchQuery == "" { // если поле поиска пусто, то просто получаем из репозитория все записи
		orders, err = h.Repository.GetOrders()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		orders, err = h.Repository.GetOrdersByTitle(searchQuery) // в ином случае ищем заказ по заголовку
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{

		"orders":    orders,
		"query":     searchQuery,
		"cartCount": count,
		"cartId":    cartId,
	})
}
func (h *Handler) GetOrder(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	order, err := h.Repository.GetOrder(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"order": order,
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

func (h *Handler) ShowCart(ctx *gin.Context) {
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
	items, _ := h.Repository.GetApplicationComponents(id)
	for _, i := range items {
		total += i.Price
	}

	ctx.HTML(http.StatusOK, "cart.html", gin.H{
		"items":     items,
		"total":     total,
		"cartCount": len(items),
	})
}
