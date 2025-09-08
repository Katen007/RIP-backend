package api

import (
	"lab1_rip/internal/app/handler"
	"lab1_rip/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")
	// слева название папки, в которую выгрузится наша статика
	// справа путь к папке, в которой лежит статика

	r.GET("/texts", handler.GetOrders)
	r.GET("/order/:id", handler.GetOrder)
	r.GET("/cart/:id", handler.ShowCart) // страница заявки
	// r.GET("/cart/add/:id", handler.AddToCart)         // добавить
	// r.GET("/cart/remove/:id", handler.RemoveFromCart) // удалить одну штуку

	r.Run()
	log.Println("Server down")
}
