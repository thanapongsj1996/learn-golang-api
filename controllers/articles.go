package controllers

import (
	"learn-golang-api/cloudbucket"
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
	Title      string                `form:"title" binding:"required"`
	Body       string                `form:"body" binding:"required"`
	Excerpt    string                `form:"excerpt" binding:"required"`
	CategoryID uint                  `form:"categoryId" binding:"required"`
	Image      *multipart.FileHeader `form:"image"`
}

type updateArticleForm struct {
	Title      string                `form:"title"`
	Body       string                `form:"body"`
	Excerpt    string                `form:"excerpt"`
	CategoryID uint                  `form:"categoryId"`
	Image      *multipart.FileHeader `form:"image"`
}

type articleResponse struct {
	ID         uint   `json:"id"`
	Title      string `json:"title"`
	Excerpt    string `json:"excerpt"`
	Body       string `json:"body"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryID"`
	Category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	User struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	} `json:"user"`
}

type articlesPaging struct {
	Items  []articleResponse `json:"items"`
	Paging *pagingResult     `json:"paging"`
}

func (a *Articles) FindAll(ctx *gin.Context) {
	var articles []models.Article
	query := a.DB.Preload("User").Preload("Category").Order("id")

	// default limit => 12
	// /articles => limit => 12, page => 1
	// /articles?limit=10 => limit => 10, page => 1
	// /articles?page=10 => limit => 12, page => 10
	// /articles?limit=10&page=2 => limit => 10, page => 2

	pagination := pagination{ctx: ctx, query: query, records: &articles}
	paging := pagination.paginate()

	serializedArticles := []articleResponse{}
	copier.Copy(&serializedArticles, &articles)

	ctx.JSON(http.StatusOK, gin.H{"articles": articlesPaging{Items: serializedArticles, Paging: paging}})
}

func (a *Articles) FindOne(ctx *gin.Context) {
	article, err := a.findArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	serializedArticle := articleResponse{}
	copier.Copy(&serializedArticle, &article)

	ctx.JSON(http.StatusOK, gin.H{"article": serializedArticle})
}

func (a *Articles) Create(ctx *gin.Context) {
	form := creatArticleForm{}
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// form => article
	article := models.Article{}
	sub, _ := ctx.Get("sub")
	user := sub.(*models.User)

	copier.Copy(&article, &form)
	article.User = *user

	// article => DB
	if err := a.DB.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// a.setArticleImage(ctx, &article)
	cloudbucket.HandleFileUploadToBucket(ctx, "image", "profile")
	serializedArticle := articleResponse{}
	copier.Copy(&serializedArticle, &article)

	ctx.JSON(http.StatusCreated, gin.H{"article": serializedArticle})
}

func (a *Articles) Update(ctx *gin.Context) {
	var form updateArticleForm

	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	article, err := a.findArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := a.DB.Model(&article).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	a.setArticleImage(ctx, &article)

	var serializedArticle articleResponse
	copier.Copy(&serializedArticle, &article)
	ctx.JSON(http.StatusOK, gin.H{"article": serializedArticle})
}

func (a *Articles) Delete(ctx *gin.Context) {
	article, err := a.findArticleByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// ใช้ลบแบบ soft delete
	// a.DB.Delete(&article)

	// ใช้ลบแบบ hard delete
	// a.DB.Unscoped().Delete(&article)

	if err := a.DB.Delete(&article).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
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
	dst := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		return err
	}

	article.Image = os.Getenv("HOST") + "/" + dst
	a.DB.Save(article)

	return nil
}

func (a *Articles) findArticleByID(ctx *gin.Context) (models.Article, error) {
	var article models.Article
	id := ctx.Param("id")

	if err := a.DB.Preload("User").Preload("Category").First(&article, id).Error; err != nil {
		return article, err
	}

	return article, nil
}
