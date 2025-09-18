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

type Text struct { // вот наша новая структура
	ID          int // поля структур, которые передаются в шаблон
	Title       string
	ImageURL    string
	Description string
	Price       int
}
type CalculationText struct {
	Text      Text
	IndexRead float32
}
type Calculation struct {
	Texts []CalculationText
}

var texts = []Text{
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

var ReadIndxsComponents = map[int]Calculation{
	1: {
		Texts: []CalculationText{
			{Text: texts[1], IndexRead: 15},
		},
	},
}

func (r *Repository) GetTexts() ([]Text, error) {
	// имитируем работу с БД. Типа мы выполнили sql запрос и получили эти строки из БД

	// обязательно проверяем ошибки, и если они появились - передаем выше, то есть хендлеру
	// тут я снова искусственно обработаю "ошибку" чисто чтобы показать вам как их передавать выше
	if len(texts) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return texts, nil
}
func (r *Repository) GetText(id int) (Text, error) {
	// тут у вас будет логика получения нужной услуги, тоже наверное через цикл в первой лабе, и через запрос к БД начиная со второй
	texts, err := r.GetTexts()
	if err != nil {
		return Text{}, err // тут у нас уже есть кастомная ошибка из нашего метода, поэтому мы можем просто вернуть ее
	}

	for _, text := range texts {
		if text.ID == id {
			return text, nil // если нашли, то просто возвращаем найденный заказ (услугу) без ошибок
		}
	}
	return Text{}, fmt.Errorf("заказ не найден") // тут нужна кастомная ошибка, чтобы понимать на каком этапе возникла ошибка и что произошло
}
func (r *Repository) GetTextsByTitle(title string) ([]Text, error) {
	texts, err := r.GetTexts()
	if err != nil {
		return []Text{}, err
	}

	var result []Text
	for _, text := range texts {
		if strings.Contains(strings.ToLower(text.Title), strings.ToLower(title)) {
			result = append(result, text)
		}
	}

	return result, nil
}

func (r *Repository) GetReadIndxsComponents(id int) (Calculation, error) {
	return ReadIndxsComponents[id], nil
}
