package handler

// UserCredentials пример тела логина.
// swagger:model UserCredentials
type UserCredentials struct {
	// required: true
	Login string `json:"login"`
	// required: true
	Password string `json:"password"`
}

// UserRegisterRequest тело регистрации.
// swagger:model UserRegisterRequest
type UserRegisterRequest struct {
	// min length: 3
	// max length: 25
	// required: true
	Login string `json:"login" binding:"required,min=3,max=25"`
	// min length: 4
	// max length: 64
	// required: true
	Password    string `json:"password" binding:"required,min=4,max=64"`
	IsModerator bool   `json:"isModerator"`
}

// UserResponse ответ с данными пользователя.
// swagger:model UserResponse
type UserResponse struct {
	ID          int    `json:"id"`
	Login       string `json:"login"`
	IsModerator bool   `json:"isModerator"`
}

// GenericOK стандартный ок.
// swagger:model GenericOK
type GenericOK struct {
	Ok bool `json:"ok"`
}

// ErrorResponse стандартная ошибка.
// swagger:model ErrorResponse
type ErrorResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

// CartIconResponse ответ для корзины чтений.
// swagger:model CartIconResponse
type CartIconResponse struct {
	DraftReadIndxsID int `json:"draft_readIndxs_id"`
	TextsCount       int `json:"texts_count"`
}

// TextDTO пример ответа текста.
// swagger:model TextDTO
type TextDTO struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	ImageURL    string `json:"image_url"`
}

// TextCreateRequest тело создания текста.
// swagger:model TextCreateRequest
type TextCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

// TextUpdateRequest частичное обновление текста.
// swagger:model TextUpdateRequest
type TextUpdateRequest map[string]any

// ReadIndxsListItem элемент списка индексов чтения (пример).
// swagger:model ReadIndxsListItem
type ReadIndxsListItem struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	TextsCnt  int    `json:"texts_count"`
}

// ReadIndxsListResponse список индексов.
// swagger:model ReadIndxsListResponse
type ReadIndxsListResponse struct {
	Items []ReadIndxsListItem `json:"items"`
}

// ReadIndxsInfoResponse детальная инфа по индексу (пример).
// swagger:model ReadIndxsInfoResponse
type ReadIndxsInfoResponse struct {
	ID       int       `json:"id"`
	Status   string    `json:"status"`
	Texts    []TextDTO `json:"texts"`
	Thematic any       `json:"thematic"`
}

// ReadIndxsModerateRequest тело модерации.
// swagger:model ReadIndxsModerateRequest
type ReadIndxsModerateRequest struct {
	Action string `json:"action"`
}

// ReadIndxsModerateResponse ответ модерации (пример).
// swagger:model ReadIndxsModerateResponse
type ReadIndxsModerateResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

// mmBody тело обновления метрик текста в индексе.
// swagger:model ReadIndxsTextMetricsUpdate
type mmBody struct {
	// required: true
	ReadIndxsID int `json:"read_indxs_id" binding:"required"`
	// required: true
	TextID         int  `json:"text_id"      binding:"required"`
	CountWords     *int `json:"count_words"`
	CountSentences *int `json:"count_sentences"`
	CountSyllables *int `json:"count_syllables"`
}

// mmBodyDelete тело удаления текста из индекса.
// swagger:model ReadIndxsTextDelete
type mmBodyDelete struct {
	// required: true
	ReadIndxsID int `json:"read_indxs_id" binding:"required"`
	// required: true
	TextID int `json:"text_id"      binding:"required"`
}
