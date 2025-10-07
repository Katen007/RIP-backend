package ds

type ReadIndxsToText struct {
	TextID      int `gorm:"primaryKey;autoIncrement:false" json:"-"`
	ReadIndxsID int `gorm:"primaryKey;autoIncrement:false" json:"-"`
	// Formula      string `gorm:"varchar(30)"` // УДАЛЯЕМ

	Calculation    int `gorm:"default:1" json:"calculation"`
	CountWords     int `gorm:"not null;default:0" json:"count_words"`
	CountSentences int `gorm:"not null;default:0" json:"count_sentences"`
	CountSyllables int `gorm:"not null;default:0" json:"count_syllables"`

	Text      Text      `gorm:"foreignKey:TextID;references:ID" json:"data"`
	ReadIndxs ReadIndxs `gorm:"foreignKey:ReadIndxsID;references:ID" json:"-"`
}
