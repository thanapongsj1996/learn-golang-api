package controllers

import (
	"learn-golang-api/config"
	"learn-golang-api/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

type Users struct {
	DB *gorm.DB
}

type creatUsersForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type updateUsersForm struct {
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"omitempty,min=8"`
	Name     string `json:"name"`
}

type userResponse struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Role   string `json:"role"`
}

type usersPaging struct {
	Items  []userResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}

func (u *Users) FindAll(ctx *gin.Context) {
	var users []models.User
	query := u.DB.Order("id").Find(&users)

	pagination := pagination{ctx: ctx, query: query, records: &users}
	paging := pagination.paginate()

	var serializedUsers []userResponse
	copier.Copy(&serializedUsers, &users)
	ctx.JSON(http.StatusOK, gin.H{"users": usersPaging{Items: serializedUsers, Paging: paging}})
}

func (u *Users) FindOne(ctx *gin.Context) {
	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedUsers models.User
	copier.Copy(&serializedUsers, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUsers})
}

func (u *Users) findUserByID(ctx *gin.Context) (*models.User, error) {
	id := ctx.Param("id")
	var user models.User

	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func setUserImage(ctx *gin.Context, user *models.User) error {
	file, err := ctx.FormFile("avatar")
	if err != nil || file == nil {
		return err
	}

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}

	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.MkdirAll(path, os.ModePerm)
	dst := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		return err
	}

	db := config.GetDB()
	user.Avatar = os.Getenv("HOST") + "/" + dst
	db.Save(user)

	return nil
}
