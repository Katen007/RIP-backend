package repository

import (
	"lab1_rip/internal/app/ds"
)

func (r *Repository) GetTexts() ([]ds.Text, error) {
	// тут мы пользуемся ORM
	var texts []ds.Text
	err := r.db.Find(&texts).Error
	if err != nil {
		return nil, err
	}
	return texts, nil
}

func (r *Repository) GetTextsByTitle(title string) ([]ds.Text, error) {
	var texts []ds.Text
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&texts).Error
	if err != nil {
		return nil, err
	}
	return texts, nil
}

func (r *Repository) GetText(id int) (ds.Text, error) {
	order := ds.Text{}
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return ds.Text{}, err
	}
	return order, nil
}
