package repository

import (
	"errors"
	"fmt"
	"lab1_rip/internal/app/ds"
	"time"

	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) ReadIndxsCartIcon(userId int) (*int, int, error) {
	var ri ds.ReadIndxs
	if err := r.db.Where("creator_id = ? AND status = ?", userId, ds.StatusDraft).First(&ri).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	var c int64
	if err := r.db.Model(&ds.ReadIndxsToText{}).Where("read_indxs_id = ?", ri.ID).Count(&c).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	return &ri.ID, int(c), nil
}

type ReadIndxsFilter struct {
	Status   string
	DateFrom *time.Time
	DateTo   *time.Time
}

func (r *Repository) ReadIndxsList(user ds.User, f ReadIndxsFilter) ([]ds.ReadIndxs, error) {
	q := r.db.Model(ds.ReadIndxs{}).Preload("Creator").Preload("Moderator").Preload("ReadIndxsToTexts.Text").
		Where("status IN ?", []string{ds.StatusFormed, ds.StatusCompleted, ds.StatusRejected})
	if !user.IsModerator {
		q = q.Where("creator_id = ?", user.ID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.DateFrom != nil {
		q = q.Where("date_form::date >= ?", f.DateFrom.Format("2006-01-02"))
	}
	if f.DateTo != nil {
		q = q.Where("date_form::date <= ?", f.DateTo.Format("2006-01-02"))
	}
	var v []ds.ReadIndxs
	res := q.Limit(10).Find(&v)
	if res.RowsAffected == 0 {
		return v, gorm.ErrRecordNotFound
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return v, nil
}

func (r *Repository) ReadIndxsGet(id int) (ds.ReadIndxs, error) {
	var ri ds.ReadIndxs
	res := r.db.Preload("Creator").Preload("Moderator").
		Preload("ReadIndxsToTexts.Text").Where("id = ? AND status <> ?", id, ds.StatusDeleted).
		Find(&ri)
	if res.RowsAffected == 0 {
		return ri, gorm.ErrRecordNotFound
	}
	return ri, res.Error
}

func (r *Repository) ReadIndxsUpdateThematic(id int, fields map[string]any) error {
	for _, k := range []string{"id", "status", "creator_id", "moderator_id", "date_create", "date_form", "date_end"} {
		delete(fields, k)
	}
	return r.db.Model(&ds.ReadIndxs{}).Where("id = ?", id).Updates(fields).Error
}

func (r *Repository) ReadIndxsUpdateCalcs(field ds.UpdateCalculationReadIndxs) error {
	//for _, k := range fields {
	err := r.db.Model(&ds.ReadIndxsToText{}).Where("text_id = ? AND read_indxs_id = ?", field.TextID, field.ReadIndxsID).Updates(map[string]interface{}{"calculation": field.Calculation}).Error
	if err != nil {
		return err
	}
	//}
	return nil
}

func (r *Repository) ReadIndxsForm(id int) error {
	var ri ds.ReadIndxs
	if err := r.db.First(&ri, id).Error; err != nil {
		return err
	}
	if ri.Status != ds.StatusDraft {
		return errors.New("only DRAFT can be formed")
	}

	var cnt int64
	if err := r.db.Model(&ds.ReadIndxsToText{}).Where("read_indxs_id = ?", id).Count(&cnt).Error; err != nil {
		return err
	}
	if ri.Contacts == "" || cnt == 0 {
		return errors.New("contacts and at least one text are required")
	}

	return r.db.Model(&ds.ReadIndxs{}).Where("id = ?", id).
		Updates(map[string]any{"status": ds.StatusFormed, "date_form": time.Now()}).Error
}

// лучше вынести в конфиг
const asyncServiceURL = "http://localhost:8000/api/calculate_readability/"

func sendToAsyncService(riID int, it ds.ReadIndxsToText) {
	payload := map[string]interface{}{
		"read_indxs_id":   riID,
		"text_id":         it.TextID,
		"count_words":     it.CountWords,
		"count_sentences": it.CountSentences,
		"count_syllables": it.CountSyllables,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("async: marshal error: %v", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, asyncServiceURL, bytes.NewReader(body))
	if err != nil {
		log.Printf("async: new request error: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	go func() {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("async: request error: %v", err)
			return
		}
		defer resp.Body.Close()
		log.Printf("async: status from Django: %s", resp.Status)
	}()
}

func (r *Repository) ReadIndxsModerate(id, moderatorID int, action string) (ds.ReadIndxs, error) {
	var ri ds.ReadIndxs
	if err := r.db.Preload("ReadIndxsToTexts.Text").First(&ri, id).Error; err != nil {
		return ds.ReadIndxs{}, err
	}
	if ri.Status != ds.StatusFormed {
		return ds.ReadIndxs{}, errors.New("only FORMED can be moderated")
	}

	// вместо локального calcFlesch — отправляем каждую запись в асинхронный сервис
	for _, it := range ri.ReadIndxsToTexts {
		// важно: CountWords / CountSentences / CountSyllables уже должны быть заполнены
		sendToAsyncService(ri.ID, it)
	}

	newStatus := ds.StatusCompleted
	if action == "reject" {
		newStatus = ds.StatusRejected
	}

	if err := r.db.Model(&ds.ReadIndxs{}).Where("id = ?", id).
		Updates(map[string]any{
			"status":       newStatus,
			"moderator_id": moderatorID,
			"date_end":     time.Now(),
		}).Error; err != nil {
		return ds.ReadIndxs{}, err
	}

	return r.ReadIndxsGet(id)
}

func (r *Repository) ReadIndxsSoftDelete(id int) error {
	return r.db.Model(&ds.ReadIndxs{}).
		Where("id = ? AND status <> ?", id, ds.StatusDeleted).
		Update("status", ds.StatusDeleted).Error
}

func (r *Repository) GetCurrentReadIndxs(userId int) (ds.ReadIndxs, error) {
	var exist_calc ds.ReadIndxs
	findErr := r.db.Where("creator_id = ? AND status = ?", userId, ds.StatusDraft).First(&exist_calc).Error
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
		return fmt.Errorf("record already exist")
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
