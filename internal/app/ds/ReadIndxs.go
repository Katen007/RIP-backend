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
