package repository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // подключаемся к БД
	if err != nil {
		return nil, err
	}

	// Возвращаем объект Repository с подключенной базой данных
	return &Repository{
		db: db,
	}, nil
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
