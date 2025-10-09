package ds

import (
	"time"
)

type ReadIndxs struct {
	ID               int               `gorm:"primaryKey" json:"id"`
	Status           string            `gorm:"type:varchar(50);not null;default:'DRAFT'" json:"status"`
	DateCreate       time.Time         `gorm:"autoCreateTime;not null" json:"date_create"`
	CreatorId        int               `gorm:"not null;" json:"-"`
	DateForm         time.Time         `json:"date_form"`
	DateEnd          time.Time         `json:"date_end"`
	ModeratorId      *int              `json:"-"`
	Comments         string            `gorm:"type:varchar(200);" json:"commenst"`
	Contacts         string            `gorm:"type:varchar(50);" json:"contacts"`
	Creator          User              `gorm:"foreignKey:CreatorId;" json:"creator"`
	Moderator        User              `gorm:"foreignKey:ModeratorId;" json:"moderator"`
	ReadIndxsToTexts []ReadIndxsToText `gorm:"foreignKey:ReadIndxsID" json:"texts"`
}

type ReadIndxsDTO struct {
	ID           int        `json:"id"`
	Status       string     `json:"status"`
	DateCreate   time.Time  `json:"date_create"`
	DateForm     *time.Time `json:"date_form"`
	DateEnd      *time.Time `json:"date_end"`
	Comments     string     `json:"commenst"`
	Contacts     *string    `json:"contacts"`
	Creator      string     `json:"creator"`
	Moderator    *string    `json:"moderator"`
	Calculations []int      `json:"calculations"`
}
type ReadIndxsInfoDTO struct {
	ID         int               `json:"id"`
	Status     string            `json:"status"`
	DateCreate time.Time         `json:"date_create"`
	DateForm   *time.Time        `json:"date_form"`
	DateEnd    *time.Time        `json:"date_end"`
	Comments   string            `json:"commenst"`
	Contacts   *string           `json:"contacts"`
	Creator    string            `json:"creator"`
	Moderator  *string           `json:"moderator"`
	Texts      []ReadIndxsToText `json:"texts"`
}

func ToReadIndxsDTO(orm ReadIndxs) ReadIndxsDTO {
	calcs := make([]int, len(orm.ReadIndxsToTexts))
	for i, el := range orm.ReadIndxsToTexts {
		calcs[i] = int(el.Calculation)
	}
	dto := ReadIndxsDTO{
		ID:           orm.ID,
		Status:       orm.Status,
		DateCreate:   orm.DateCreate,
		DateForm:     &orm.DateForm,
		DateEnd:      &orm.DateEnd,
		Comments:     orm.Comments,
		Contacts:     &orm.Contacts,
		Creator:      orm.Creator.Login,
		Moderator:    &orm.Moderator.Login,
		Calculations: calcs,
	}
	return dto
}

func ToReadIndxsInfoDTO(orm ReadIndxs) ReadIndxsInfoDTO {
	dto := ReadIndxsInfoDTO{
		ID:         orm.ID,
		Status:     orm.Status,
		DateCreate: orm.DateCreate,
		DateForm:   &orm.DateForm,
		DateEnd:    &orm.DateEnd,
		Comments:   orm.Comments,
		Contacts:   &orm.Contacts,
		Creator:    orm.Creator.Login,
		Moderator:  &orm.Moderator.Login,
		Texts:      orm.ReadIndxsToTexts,
	}
	return dto
}

func ToReadIndxsInfoListDTO(orms []ReadIndxs) []ReadIndxsInfoDTO {
	var dtoList []ReadIndxsInfoDTO
	for _, el := range orms {
		dtoList = append(dtoList, ToReadIndxsInfoDTO(el))
	}
	return dtoList
}

func ToReadIndxsListDTO(orms []ReadIndxs) []ReadIndxsDTO {
	var dtoList []ReadIndxsDTO
	for _, el := range orms {
		dtoList = append(dtoList, ToReadIndxsDTO(el))
	}
	return dtoList
}
