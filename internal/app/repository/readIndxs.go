package repository

import (
	"lab1_rip/internal/app/ds"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) GetCurrentReadIndxs(userId int) (ds.ReadIndxs, error) {
	exist_calc, findErr := r.GetReadIndxs(userId)
	if findErr != nil {
		if findErr == gorm.ErrRecordNotFound {
			return r.CreateReadIndxs(userId)
		} else {
			return ds.ReadIndxs{}, findErr
		}
	}
	return exist_calc, nil
}

func (r *Repository) GetReadIndxsById(id int) (ds.ReadIndxs, error) {
	var readIndxs ds.ReadIndxs
	findErr := r.db.Where("id = ? AND status != ?", id, "DELETED").First(&readIndxs).Error
	if findErr != nil {
		return ds.ReadIndxs{}, findErr
	}
	return readIndxs, nil
}
func (r *Repository) GetReadIndxs(userID int) (ds.ReadIndxs, error) {
	var readIndxs ds.ReadIndxs
	findErr := r.db.Where("creator_id = ? AND status != ?", userID, "DELETED").First(&readIndxs).Error
	if findErr != nil {
		return ds.ReadIndxs{}, findErr
	}
	return readIndxs, nil
}

func (r *Repository) CreateReadIndxs(userId int) (ds.ReadIndxs, error) {
	newReadIndxs := ds.ReadIndxs{
		CreatorId:   userId,
		ModeratorId: nil,
	}
	createErr := r.db.Create(&newReadIndxs).Error
	if createErr != nil {
		return ds.ReadIndxs{}, createErr
	}
	return newReadIndxs, nil
}
func (r *Repository) AddTextToReadIndxs(textID int, userId int) error {
	readIndxs, err := r.GetCurrentReadIndxs(userId)
	if err != nil {
		return err
	}

	var existing ds.ReadIndxsToText
	check := r.db.Where("text_id = ? AND read_indxs_id = ?", textID, readIndxs.ID).First(&existing)
	logrus.Error(textID, readIndxs.ID)
	if check.Error == nil {
		logrus.Error("not em")
		return nil
	}
	if check.Error != nil && check.Error != gorm.ErrRecordNotFound {
		return check.Error
	}

	textToReadIndxs := ds.ReadIndxsToText{
		TextID:      textID,
		ReadIndxsID: readIndxs.ID,
	}
	err = r.db.Create(&textToReadIndxs).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetTextsInReadIndxs(readIndxsID int) ([]ds.ReadIndxsToText, error) {
	var readIndxs ds.ReadIndxs
	err := r.db.Preload("ReadIndxsToTexts.Text").First(&readIndxs, readIndxsID).Error
	if err != nil {
		return nil, err
	}
	items := readIndxs.ReadIndxsToTexts
	return items, nil
}

func (r *Repository) GetCountTexts(userId int) (int, error) {
	readIndxs, err := r.GetReadIndxs(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}

	var count int64
	err = r.db.Model(&ds.ReadIndxsToText{}).
		Where("read_indxs_id = ?", readIndxs.ID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *Repository) DeleteReadIndxs(readIndxsId int) error {
	var err error
	var ids []int
	query := "UPDATE read_indxs SET status=$1 WHERE id = $2 AND status!=$1"
	res := r.db.Raw(query, "DELETED", readIndxsId).Scan(&ids)
	if err = res.Error; err != nil {
		return err
	}
	return nil
}
