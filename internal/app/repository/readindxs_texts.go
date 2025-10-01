package repository

import "lab1_rip/internal/app/ds"

func (r *Repository) ReadIndxsTextsUpdate(readIndxsID, textID int, fields map[string]any) error {
	return r.db.Model(&ds.ReadIndxsToText{}).
		Where("read_indxs_id = ? AND text_id = ?", readIndxsID, textID).
		Updates(fields).Error
}
func (r *Repository) ReadIndxsTextsDelete(readIndxsID, textID int) error {
	return r.db.Where("read_indxs_id = ? AND text_id = ?", readIndxsID, textID).
		Delete(&ds.ReadIndxsToText{}).Error
}
