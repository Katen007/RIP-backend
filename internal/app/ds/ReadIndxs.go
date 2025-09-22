package ds

import (
	"time"
)

type ReadIndxs struct {
	ID               int       `gorm:"primaryKey"`
	Status           string    `gorm:"type:varchar(50);not null;default:'DRAFT'"`
	DateCreate       time.Time `gorm:"autoCreateTime;not null"`
	CreatorId        int       `gorm:"not null;"`
	DateForm         time.Time
	DateEnd          time.Time
	ModeratorId      *int
	Comments         string            `gorm:"type:varchar(200);"`
	Contacts         string            `gorm:"type:varchar(50);"`
	Creator          User              `gorm:"foreignKey:CreatorId;"`
	Moderator        User              `gorm:"foreignKey:ModeratorId;"`
	ReadIndxsToTexts []ReadIndxsToText `gorm:"foreignKey:ReadIndxsID"`
}
