package handler

import (
	"lab1_rip/internal/app/ds"

	"github.com/gin-gonic/gin"
)

// API_TextsList список текстов с фильтром по заголовку.
// @Summary      List texts
// @Tags         texts
// @Produce      json
// @Param        title  query   string  false  "title contains"
// @Success      200    {array} TextDTO
// @Failure      500    {object} ErrorResponse
// @Router       /texts [get]
func (h *Handler) API_TextsList(c *gin.Context) {
	title := c.Query("title")
	data, err := h.Repository.TextsList(title)
	if err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	textDto := ds.ToTextsListDTO(data)
	c.JSON(200, textDto)
}

// API_TextGet получить текст по id.
// @Summary      Get text by id
// @Tags         texts
// @Produce      json
// @Param        id   path  int  true  "Text ID"
// @Success      200  {object}  TextDTO
// @Failure      404  {object}  ErrorResponse
// @Router       /texts/{id} [get]
func (h *Handler) API_TextGet(c *gin.Context) {
	id := mustIntParam(c, "id")
	data, err := h.Repository.TextByID(id)
	if err != nil {
		h.errorHandler(c, 404, err)
		return
	}
	textDto := ds.ToTextDTO(data)
	c.JSON(200, textDto)
}

// API_TextCreate создать текст.
// @Summary      Create text
// @Tags         texts
// @Accept       json
// @Produce      json
// @Param        body  body  TextCreateRequest  true  "new text payload"
// @Success      201   {object}  TextDTO
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /texts [post]
func (h *Handler) API_TextCreate(c *gin.Context) {
	var dto struct {
		Title, Description string
		Price              int
	}
	if err := c.BindJSON(&dto); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	t := ds.Text{Title: dto.Title, Description: dto.Description, Price: dto.Price}
	if err := h.Repository.TextCreate(&t); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.JSON(201, t)
}

// API_TextUpdate частичное обновление текста.
// @Summary      Update text (partial)
// @Tags         texts
// @Accept       json
// @Produce      json
// @Param        id    path  int                true  "Text ID"
// @Param        body  body  TextUpdateRequest  true  "fields to update"
// @Success      204   {string}  string "No Content"
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /texts/{id} [patch]
func (h *Handler) API_TextUpdate(c *gin.Context) {
	id := mustIntParam(c, "id")
	var m map[string]any
	if err := c.BindJSON(&m); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	if err := h.Repository.TextUpdate(id, m); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.Status(204)
}

// API_TextDelete мягкое удаление текста (и его изображения).
// @Summary      Soft delete text
// @Tags         texts
// @Produce      json
// @Param        id   path  int  true "Text ID"
// @Success      204  {string} string "No Content"
// @Failure      400  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /texts/{id} [delete]
func (h *Handler) API_TextDelete(c *gin.Context) {
	id := mustIntParam(c, "id")
	if c.IsAborted() {
		return
	}
	text, err := h.Repository.TextByID(id)
	if err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	imgUrl := text.ImageURL
	if err := h.Repository.TextSoftDelete(id); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	err = h.Repository.DeleteComponentImg(c, &imgUrl)
	if err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.Status(204)
}

// API_TextUploadImage загрузка изображения к тексту.
// В твоей реализации читается бинарное тело напрямую из Request.Body.
// @Summary      Upload text image
// @Tags         texts
// @Accept       octet-stream
// @Produce      json
// @Param        id    path   int     true  "Text ID"
// @Param        file  body   []byte  true  "raw binary image (Content-Type: image/*)"
// @Success      201   {object}  map[string]string  "image_key"
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /texts/{id}/image [post]
func (h *Handler) API_TextUploadImage(c *gin.Context) {
	id := mustIntParam(c, "id")
	if c.IsAborted() {
		return
	}
	fileSize, contentType := h.fileParams(c)
	if c.IsAborted() {
		return
	}
	text, err := h.Repository.TextByID(id)
	if err != nil {
		h.errorHandler(c, 400, err)
		return
	}

	file := c.Request.Body
	defer file.Close()

	url, err := h.Repository.UploadComponentImg(c, file, "img/", fileSize, contentType)
	if err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	if text.ImageURL != "" {
		h.Repository.DeleteComponentImg(c, &text.ImageURL)
	}
	if err := h.Repository.TextUpdateImageKey(id, &url); err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	c.JSON(201, gin.H{"image_key": url})
}

// API_TextAddToDraft godoc
// @Summary      Добавить текст в черновик индексов пользователя
// @Tags         texts
// @Param        id   path  int  true  "textID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]interface{}
// @Router       /texts/{id}/add-to-draft [post]
func (h *Handler) API_TextAddToDraft(c *gin.Context) {
	textID := mustIntParam(c, "id")
	userID := h.GetUserID(c)
	if c.IsAborted() {
		return
	}
	if err := h.Repository.AddTextToReadIndxs(textID, userID); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	c.Status(204)
}
