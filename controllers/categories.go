package controllers

import (
	"learn-golang-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Categories struct {
	DB *gorm.DB
}

type categoryResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Articles []struct {
		ID    uint   `json:"id"`
		Title string `json:"title"`
	} `json:"articles"`
}

type allCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type createCategoryForm struct {
	Name string `json:"name" binding:"required"`
	Desc string `json:"desc" binding:"required"`
}

type updateCategoryForm struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (c *Categories) FindAll(ctx *gin.Context) {
	var categories []models.Category

	if err := c.DB.Order("id").Find(&categories).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var serializedCategories []allCategoryResponse
	copier.Copy(&serializedCategories, &categories)
	ctx.JSON(http.StatusOK, gin.H{"categories": serializedCategories})
}

func (c *Categories) FindOne(ctx *gin.Context) {
	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedCategory categoryResponse
	copier.Copy(&serializedCategory, &category)
	ctx.JSON(http.StatusOK, gin.H{"category": serializedCategory})
}

func (c *Categories) Create(ctx *gin.Context) {
	var form createCategoryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	copier.Copy(&category, &form)

	if err := c.DB.Create(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedCategory categoryResponse
	copier.Copy(&serializedCategory, &category)
	ctx.JSON(http.StatusCreated, gin.H{"category": serializedCategory})
}

func (c *Categories) Update(ctx *gin.Context) {
	var form updateCategoryForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Model(&category).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedCategory categoryResponse
	copier.Copy(&serializedCategory, &category)
	ctx.JSON(http.StatusOK, serializedCategory)
}

func (c *Categories) Delete(ctx *gin.Context) {
	category, err := c.findCategoryByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := c.DB.Unscoped().Delete(&category).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Categories) findCategoryByID(ctx *gin.Context) (models.Category, error) {
	var category models.Category
	id := ctx.Param("id")

	if err := c.DB.Preload("Articles").Find(&category, id).Error; err != nil {
		return category, err
	}

	return category, nil
}
