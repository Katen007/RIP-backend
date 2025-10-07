package ds

type Text struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	IsDelete    bool   `gorm:"type:boolean not null;default:false" json:"-"`
	ImageURL    string `gorm:"type:varchar(256)" json:"image_url"`
	Title       string `gorm:"type:varchar(256);not null" json:"title"`
	Description string `gorm:"type:varchar(1024)" json:"-"`
	Price       int    `json:"-"`

	ReadIndxsToTexts []ReadIndxsToText `gorm:"foreignKey:TextID" json:"-"`
}

type TextDTO struct {
	ID          int    `json:"id"`
	ImageURL    string `json:"image_url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

func ToTextDTO(orm Text) TextDTO {
	dto := TextDTO{
		ID:          orm.ID,
		ImageURL:    orm.ImageURL,
		Title:       orm.Title,
		Description: orm.Description,
		Price:       orm.Price,
	}
	return dto
}

func ToTextsListDTO(orms []Text) []TextDTO {
	var dtoList []TextDTO
	for _, el := range orms {
		dtoList = append(dtoList, ToTextDTO(el))
	}
	return dtoList
}
