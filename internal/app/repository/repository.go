package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Order struct { // вот наша новая структура
	ID          int // поля структур, которые передаются в шаблон
	Title       string
	ImageURL    string
	Description string
	Price       int
}

var ApplicationComponents = map[int][]Order{
	1: {
		{ID: 1, Title: "Научная статья", ImageURL: "http://localhost:9000/img/img/sci.jpg", Description: `Мы анализируем научные тексты: проверяем структуру,
терминологическую насыщенность и вычисляем индексы читабельности (Флеша, Фога,
SMOG). Отчёт покажет, насколько материал удобен для чтения и как его упростить.`, Price: 125},
		{ID: 2, Title: "Исторический документ", ImageURL: "http://localhost:9000/img/img/his.jpg", Description: `Подходит для архивных и исторических материалов. Оцениваем
длину предложений, редкие конструкции и устаревшие формы, строим метрики
читабельности и даём рекомендации по адаптации.`, Price: 400},
	},
}

func (r *Repository) GetOrders() ([]Order, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД
	orders := []Order{
		{ID: 1, Title: "Научная статья", ImageURL: "http://localhost:9000/img/img/sci.jpg", Description: `Мы анализируем научные тексты: проверяем структуру,
терминологическую насыщенность и вычисляем индексы читабельности (Флеша, Фога,
SMOG). Отчёт покажет, насколько материал удобен для чтения и как его упростить.`, Price: 125},
		{ID: 2, Title: "Исторический документ", ImageURL: "http://localhost:9000/img/img/his.jpg", Description: `Подходит для архивных и исторических материалов. Оцениваем
длину предложений, редкие конструкции и устаревшие формы, строим метрики
читабельности и даём рекомендации по адаптации.`, Price: 400},
		{ID: 3, Title: "Художественный текст", ImageURL: "http://localhost:9000/img/img/lit.jpg", Description: `Анализ художественных фрагментов: ритм и длины фраз, доля
диалогов, простота восприятия. Сравниваем индексы с популярной литературой
и подсказываем, где можно повысить «легкость» чтения.`, Price: 500},
	}

	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
	if len(orders) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return orders, nil
}
func (r *Repository) GetOrder(id int) (Order, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Order{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}
func (r *Repository) GetOrdersByTitle(title string) ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	}

	var result []Order
	for _, order := range orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}

	return result, nil
}

func (r *Repository) GetApplicationComponents(id int) ([]Order, error) {
	return ApplicationComponents[id], nil
}
