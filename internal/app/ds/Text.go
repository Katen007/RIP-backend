package ds

type Text struct {
	ID          int    `gorm:"primaryKey"`
	IsDelete    bool   `gorm:"type:boolean not null;default:false"`
	ImageURL    string `gorm:"type:varchar(256)"`
	Title       string `gorm:"type:varchar(256);not null"`
	Description string `gorm:"type:varchar(1024)"`
	Price       int

	ReadIndxsToTexts []ReadIndxsToText `gorm:"foreignKey:TextID"`
}
