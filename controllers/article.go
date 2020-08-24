package controllers

import (
	"learn-golang-api/models"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Articles struct {
	DB *gorm.DB
}

type creatArticleForm struct {
	Title   string                `form:"title" binding:"required"`
	Body    string                `form:"body" binding:"required"`
	Excerpt string                `form:"excerpt" binding:"required"`
	Image   *multipart.FileHeader `form:"image" binding:"required"`
}

type createdArticleResponse struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Excerpt string `json:"excerpt"`
	Body    string `json:"body"`
	Image   string `json:"image"`
}

var articles []models.Article = []models.Article{}

func (a *Articles) FindAll(ctx *gin.Context) {

}

func (a *Articles) FindOne(ctx *gin.Context) {

}

func (a *Articles) Create(ctx *gin.Context) {
	form := creatArticleForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// form => article
	article := models.Article{}
	copier.Copy(&article, &form)

	// article => DB
	if err := a.DB.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(ctx, &article)
	serializedArticle := createdArticleResponse{}
	copier.Copy(&serializedArticle, &article)

	ctx.JSON(http.StatusCreated, gin.H{"article": serializedArticle})
}

func (a *Articles) setArticleImage(ctx *gin.Context, article *models.Article) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return err
	}

	if article.Image != "" {
		// Path รูปเดิม http://127.0.0.1:5000/upload/articles/<ID>/image.png
		// 1. ดึง path /upload/articles/<ID>/image.png
		article.Image = strings.Replace(article.Image, os.Getenv("HOST"), "", 1)
		// 2. ตำแหน่งรูป <WD>/upload/articles/<ID>/image.png
		pwd, _ := os.Getwd()
		// 3. ลบรูปออก Remove <WD>/upload/articles/<ID>/image.png
		os.Remove(pwd + article.Image)
	}

	path := "uploads/articles/" + strconv.Itoa(int(article.ID))
	os.MkdirAll(path, 0755)
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	article.Image = os.Getenv("HOST") + "/" + filename
	a.DB.Save(article)

	return nil
}
