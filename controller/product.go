package controller

import (
	"api-pos/db"
	"api-pos/dto"
	"api-pos/models"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct{}

func (p Product) FindAll(ctx *gin.Context) {
	categoryid := ctx.Query("categoryid")
	search := ctx.Query("search")
	status := ctx.Query("status")

	//db.Conn.Where("category_id = ?", categoryid).Find(&products)

	var products []models.Product
	query := db.Conn.Preload("Category")
	if categoryid != "" {
		query = query.Where("category_id = ?", categoryid)
	}
	if search != "" {
		query = query.Where("sku = ? OR name LIKE ?", search, "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Order("created_at desc").Find(&products)

	var result []dto.ReadProductResponse
	for _, product := range products{
		result = append(result, dto.ReadProductResponse{
			ID: product.ID,
			SKU: product.SKU,
			Name: product.Name,
			Desc: product.Desc,
			Price: product.Price,
			Status: product.Status,
			Image: product.Image,
			Category: dto.CategoryResponse{
				ID: product.Category.ID,
				Name: product.Category.Name,
			},
		})
	}

	ctx.JSON(http.StatusOK, result)
}

func (p Product) FindOne(ctx *gin.Context) {
	//get param
	id := ctx.Param("id")
	var product models.Product
	query := db.Conn.Preload("Category").First(&product, id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound){
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.ReadProductResponse{
		ID: product.ID,
			SKU: product.SKU,
			Name: product.Name,
			Desc: product.Desc,
			Price: product.Price,
			Status: product.Status,
			Image: product.Image,
			Category: dto.CategoryResponse{
				ID: product.Category.ID,
				Name: product.Category.Name,
			},
	})
}

func (p Product) Create(ctx *gin.Context) {
	var form dto.ProductRequest
	if err := ctx.ShouldBind(&form); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//create path
	imagePath := "./uploads/products/" + uuid.New().String()
	//save img
	ctx.SaveUploadedFile(image, imagePath)

	product := models.Product{
		SKU: form.SKU,
		Name: form.Name,
		Desc: form.Desc,
		Price: form.Price,
		Status: form.Status,
		CategoryID: form.CategoryID,
		Image: imagePath,
	}

	if err := db.Conn.Create(&product).Error; err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.CreateOrUpdateProductResponse{
		ID: product.ID,
		SKU: product.SKU,
		Name: product.Name,
		Desc: product.Desc,
		Price: product.Price,
		Status: product.Status,
		CategoryID: product.CategoryID,
		Image: product.Image,
	})
}
func (p Product) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var form dto.ProductRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	if err := db.Conn.First(&product, id).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if image != nil {
		imagePath := "./uploads/products/" + uuid.New().String()
		ctx.SaveUploadedFile(image, imagePath)
		os.Remove(product.Image)
		product.Image = imagePath
	}
	product.SKU = form.SKU
	product.Name = form.Name
	product.Desc = form.Desc
	product.Price = form.Price
	product.Status = form.Status
	product.CategoryID = form.CategoryID
	db.Conn.Save(&product)

	ctx.JSON(http.StatusOK, dto.CreateOrUpdateProductResponse{
		ID:         product.ID,
		SKU:        product.SKU,
		Name:       product.Name,
		Desc:       product.Desc,
		Price:      product.Price,
		Status:     product.Status,
		CategoryID: product.CategoryID,
		Image:      product.Image,
	})

}
func (p Product) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	db.Conn.Unscoped().Delete(&models.Product{}, id)

	ctx.JSON(http.StatusOK, gin.H{"delete at ": time.Now()})
}