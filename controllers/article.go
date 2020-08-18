package controllers

import (
	"learn-golang-api/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Articles struct{}

type creatArticleForm struct {
	Title string                `form:"title" binding:"required"`
	Body  string                `form:"body" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}

var articles []models.Article = []models.Article{
	{ID: 1, Title: "Title1", Body: "body1"},
	{ID: 2, Title: "Title2", Body: "body2"},
	{ID: 3, Title: "Title3", Body: "body3"},
	{ID: 4, Title: "Title4", Body: "body4"},
	{ID: 5, Title: "Title5", Body: "body5"},
	{ID: 6, Title: "Title6", Body: "body6"},
}

func (a *Articles) FindAll(ctx *gin.Context) {
	result := articles

	if limit := ctx.Query("limit"); limit != "" {
		n, _ := strconv.Atoi(limit)
		result = result[:n]
	}

	ctx.JSON(http.StatusOK, gin.H{"articles": result})
}

func (a *Articles) FindOne(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	for _, item := range articles {
		if item.ID == uint(id) {
			ctx.JSON(http.StatusOK, gin.H{"article": item})
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
}

func (a *Articles) Create(ctx *gin.Context) {
	var form creatArticleForm

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addedArticle := models.Article{
		ID:    uint(len(articles) + 1),
		Title: form.Title,
		Body:  form.Body,
	}

	// Get file
	file, _ := ctx.FormFile("image")

	// Create path
	// ID => 8, uploads/articles/8/image.png
	path := "uploads/articles/" + strconv.Itoa(int(addedArticle.ID))
	os.MkdirAll(path, 0755)

	// Upload file
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		// ...
	}

	// Attach file to article
	addedArticle.Image = os.Getenv("HOST") + "/" + filename

	articles = append(articles, addedArticle)
	ctx.JSON(http.StatusCreated, gin.H{"article": addedArticle})
}
