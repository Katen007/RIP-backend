package handler

import (
	"fmt"
	"lab1_rip/internal/app/ds"

	"github.com/gin-gonic/gin"
)

type TextsFilters struct {
	Title  string `form:"title"`
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

func (h *Handler) API_TextsList(c *gin.Context) {
	var f = TextsFilters{
		Limit:  10,
		Offset: 0,
	}
	err := c.ShouldBindQuery(&f)
	if err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	data, err := h.Repository.TextsList(f.Title, f.Limit, f.Offset)
	if err != nil {
		h.errorHandler(c, 500, err)
		return
	}
	fmt.Println(data[0].Title)
	textDto := ds.ToTextsListDTO(data)
	c.JSON(200, textDto)
}
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
func (h *Handler) API_TextAddToDraft(c *gin.Context) {
	textID := mustIntParam(c, "id")
	userID := getUserID(c)
	if c.IsAborted() {
		return
	}
	if err := h.Repository.AddTextToReadIndxs(textID, userID); err != nil {
		h.errorHandler(c, 400, err)
		return
	}
	c.Status(204)
}
