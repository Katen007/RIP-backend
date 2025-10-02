package repository

import (
	"errors"
	"lab1_rip/internal/app/ds"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrLoginTaken = errors.New("login already taken")

func (r *Repository) UserByLogin(login string) (ds.User, error) {
	var u ds.User
	err := r.db.Where("login = ?", login).First(&u).Error
	return u, err
}

func (r *Repository) UserByID(id int) (ds.User, error) {
	var u ds.User
	err := r.db.First(&u, id).Error
	return u, err
}

func (r *Repository) UserCreate(login, hashed string, isModerator bool) (ds.User, error) {
	// уникальность логина гарантируется БД, но проверим красиво
	if _, err := r.UserByLogin(login); err == nil {
		return ds.User{}, ErrLoginTaken
	} else if err != gorm.ErrRecordNotFound {
		return ds.User{}, err
	}

	u := ds.User{Login: login, HashedPassword: hashed, IsModerator: isModerator}
	return u, r.db.Create(&u).Error
}

func (r *Repository) UserUpdateMe(id int, fields map[string]any) (ds.User, error) {
	// позволяем менять только login и hashed_password
	whitelist := map[string]bool{"login": true, "hashed_password": true}
	upd := map[string]any{}
	for k, v := range fields {
		if whitelist[k] {
			upd[k] = v
		}
	}
	if len(upd) == 0 {
		return r.UserByID(id)
	}

	var out ds.User
	tx := r.db.Model(&ds.User{}).
		Where("id = ?", id).
		Clauses(clause.Returning{}).
		Updates(upd).
		Scan(&out)
	if tx.Error != nil {
		return ds.User{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return ds.User{}, gorm.ErrRecordNotFound
	}
	return out, nil
}
