package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type pagination struct {
	ctx     *gin.Context
	query   *gorm.DB
	records interface{}
}

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

func (p *pagination) paginate() *pagingResult {
	// 1. Get limit, page
	page, _ := strconv.Atoi(p.ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(p.ctx.DefaultQuery("limit", "12"))

	// 2. Count records
	ch := make(chan int)
	go p.countRecords(ch)

	// 3. Find records
	// limit , offset
	// EX. limit => 10
	// page => 1 , 1 -10 , offser = 0
	// page => 2 , 11 -20 , offser = 10
	// page => 3 , 21 -30 , offser = 20
	offset := (page - 1) * limit
	p.query.Limit(limit).Offset(offset).Find(p.records)

	// 4. Total page
	count := <-ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	// 5. Find nextPage
	var nextPage int
	if page == totalPage {
		nextPage = page
	} else {
		nextPage = page + 1
	}

	// 6. Create pagingResult
	result := pagingResult{
		Page:      page,
		Limit:     limit,
		Count:     count,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		TotalPage: totalPage,
	}

	return &result
}

func (p *pagination) countRecords(ch chan int) {
	var count int
	p.query.Model(p.records).Count(&count)

	ch <- count
}
