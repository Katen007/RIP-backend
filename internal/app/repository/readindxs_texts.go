package repository

import (
	"lab1_rip/internal/app/ds"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) ReadIndxsTextsUpdate(readIndxsID, textID int, fields map[string]any) error {
	allow := map[string]bool{
		"count_words": true, "count_sentences": true, "count_syllables": true,
	}
	upd := map[string]any{}
	for k, v := range fields {
		if allow[k] {
			upd[k] = v
		}
	}
	if len(upd) == 0 {
		return nil
	}

	return r.db.Model(&ds.ReadIndxsToText{}).
		Where("read_indxs_id = ? AND text_id = ?", readIndxsID, textID).
		Updates(upd).Error
}

func (r *Repository) ReadIndxsTextsDelete(readIndxsID, textID int) error {
	var redindx ds.ReadIndxsToText
	res := r.db.Model(&ds.ReadIndxsToText{}).Where("read_indxs_id = ? AND text_id = ?", readIndxsID, textID).Delete(&redindx)
	logrus.Error(res.Error)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *Repository) ReadIndxsTextsGet(readIndxsID, textID int) (ds.ReadIndxsToText, error) {
	var exist ds.ReadIndxsToText
	res := r.db.Model(ds.ReadIndxsToText{}).Preload("ReadIndxs").Where("read_indxs_id = ? AND text_id = ?", readIndxsID, textID).First(&exist)
	if res.RowsAffected == 0 {
		return ds.ReadIndxsToText{}, gorm.ErrRecordNotFound
	}
	if res.Error != nil {
		return ds.ReadIndxsToText{}, res.Error
	}
	return exist, nil
}
