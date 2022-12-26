package controller

import (
	"api-pos/db"
	"api-pos/dto"
	"api-pos/models"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Category struct{}

func (c Category) FindAll(ctx *gin.Context){
	var categories []models.Category
	db.Conn.Find(&categories)

	var result []dto.CategoryResponse
	for _, category := range categories{
		result = append(result, dto.CategoryResponse{
			ID: category.ID,
			Name: category.Name,
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func (c Category) FindOne(ctx *gin.Context){
	id := ctx.Param("id")
	var category models.Category
	if err := db.Conn.First(&category, id).Error; errors.Is(err, gorm.ErrRecordNotFound){
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.CategoryResponse{
		ID: category.ID,
		Name: category.Name,
	})
}

func (c Category) Create(ctx *gin.Context){
	var form dto.CatRequest
	//map body to struct
	if err := ctx.ShouldBindJSON(&form); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"eeror": err.Error()})
		return
	}

	category := models.Category{
		Name: form.Name,
	}

	if err := db.Conn.Create(&category).Error; err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.CategoryResponse{
		ID: category.ID,
		Name: category.Name,
	})
}

func (c Category) Update(ctx *gin.Context){
	id := ctx.Param("id")
	var form dto.CatRequest
	if err := ctx.ShouldBindJSON(&form); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := db.Conn.First(&category, id).Error; errors.Is(err, gorm.ErrRecordNotFound){
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	category.Name = form.Name
	db.Conn.Save(&category)
	ctx.JSON(http.StatusOK, dto.CategoryResponse{
		ID: category.ID,
		Name: category.Name,
	})
}

func (c Category) Delete(ctx *gin.Context){
	id := ctx.Param("id")
	//soft delete
	//db.Conn.Delete(&models.Category{}, id)
	//Hard delete
	db.Conn.Unscoped().Delete(&models.Category{}, id)

	ctx.JSON(http.StatusOK, gin.H{"delete at ": time.Now()})
}