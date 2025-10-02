package ds

type User struct {
	ID             int    `gorm:"primaryKey;autoIncrement" json:"-"`
	Login          string `gorm:"type:varchar(25);unique;not null"`
	HashedPassword string `gorm:"type:varchar(100);not null" json:"-"`
	IsModerator    bool   `gorm:"type:boolean;default:false" json:"-"`
}
