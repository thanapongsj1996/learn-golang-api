package models

import "github.com/jinzhu/gorm"

type Article struct {
	gorm.Model
	Title   string `gorm:"unique;not null" json:"title"`
	Excerpt string `gorm:"not null"`
	Body    string `gorm:"not null"`
	Image   string `gorm:"not null"`
}