package repository

import (
	"lab1_rip/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) TextsList(title string) ([]ds.Text, error) {
	q := r.db.Where("is_delete = false")
	if title != "" {
		q = q.Where("title ILIKE ?", "%"+title+"%")
	}
	var res []ds.Text
	return res, q.Limit(10).Order("id ASC").Find(&res).Error
}
func (r *Repository) TextByID(id int) (ds.Text, error) {
	var t ds.Text
	return t, r.db.Where("id = ? AND is_delete = false", id).First(&t).Error
}
func (r *Repository) TextCreate(t *ds.Text) error { return r.db.Create(t).Error }
func (r *Repository) TextUpdate(id int, fields map[string]any) error {
	delete(fields, "id")
	delete(fields, "is_delete")
	res := r.db.Model(&ds.Text{}).Where("id = ?", id).Updates(fields)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error

}
func (r *Repository) TextSoftDelete(id int) error {
	res := r.db.Model(&ds.Text{}).Where("id = ? AND is_delete = false", id).
		Update("is_delete", true)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return res.Error
}
func (r *Repository) TextUpdateImageKey(id int, key *string) error {
	res := r.db.Model(&ds.Text{}).Where("id = ? AND is_delete = false", id).Update("image_url", key)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return res.Error
}
