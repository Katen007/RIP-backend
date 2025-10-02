package repository

import (
	"lab1_rip/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) ReadIndxsTextsUpdate(readIndxsID, textID int, fields map[string]any) error {
	res := r.db.Model(&ds.ReadIndxsToText{}).Where("read_indxs_id = ? AND text_id = ?", readIndxsID, textID).Updates(fields)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}
func (r *Repository) ReadIndxsTextsDelete(readIndxsID, textID int) error {
	res := r.db.Where("read_indxs_id = ? AND text_id = ?", readIndxsID, textID).Delete(&ds.ReadIndxsToText{})
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}
