package ds

type ReadIndxsToText struct {
	TextID      int       `gorm:"primaryKey; autoIncrement:false"`
	ReadIndxsID int       `gorm:"primaryKey; autoIncrement:false"`
	Formula     string    `gorm:"varchar(30)"`
	Calculation int       `gorm:"default:1"`
	Text        Text      `gorm:"foreignKey:TextID; references:ID"`
	ReadIndxs   ReadIndxs `gorm:"foreignKey:ReadIndxsID; references:ID" json:"-"`
}
