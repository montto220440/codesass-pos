package controller

import (
	"api-pos/db"
	"api-pos/dto"
	"api-pos/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Order struct {}

func (o Order) FindAll(ctx *gin.Context) {
	var orders []models.Order
	db.Conn.Preload("Product").Find(&orders)

	var result []dto.OrderResponse
	for _, order := range orders{
		orderResult := dto.OrderResponse{
			ID: order.ID,
			Name: order.Name,
			Tel: order.Tel,
			Email: order.Email,
		}
		var product []dto.OrderProductResponse
		for _, pd := range order.Product{
			product = append(product, dto.OrderProductResponse{
				ID: pd.ID,
				SKU: pd.SKU,
				Name: pd.Name,
				Image: pd.Image,
				Price: pd.Price,
				Quantity: pd.Quantity,
			})
		}
		orderResult.Products = product
		result =append(result, orderResult)
	}
	ctx.JSON(http.StatusOK, result)
}

func (o Order) FindOne(ctx *gin.Context) {
	id := ctx.Param("id")
	var order models.Order
	
	query := db.Conn.Preload("Product").First(&order, id)
	if err := query.Error; errors.Is(err, gorm.ErrRecordNotFound){
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	result := dto.OrderResponse{
		ID: order.ID,
		Name: order.Name,
		Tel: order.Tel,
		Email: order.Email,
	}
	var product []dto.OrderProductResponse
	for _, pd := range order.Product{
		product = append(product, dto.OrderProductResponse{
			ID: pd.ID,
			SKU: pd.SKU,
			Name: pd.Name,
			Image: pd.Image,
			Price: pd.Price,
			Quantity: pd.Quantity,
		})
	}
	result.Products= product

	ctx.JSON(http.StatusOK, result)
}

func (o Order) Create(ctx *gin.Context) {
	var form dto.OrderRequest
	if err := ctx.ShouldBindJSON(&form); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var order models.Order
	var orderItem []models.Order_item
	for _, product := range form.Products{
		orderItem = append(orderItem, models.Order_item{
			SKU: product.SKU,
			Name: product.Name,
			Price: product.Price,
			Quantity: product.Quantity,
			Image: product.Image,
		})
	}

	order.Name = form.Name
	order.Tel = form.Tel
	order.Email = form.Email
	order.Product = orderItem
	db.Conn.Create(&order)

	result := dto.OrderResponse{
		ID: order.ID,
		Name: order.Name,
		Tel: order.Tel,
		Email: order.Email,
	}
	var product []dto.OrderProductResponse
	for _, pd := range order.Product{
		product = append(product, dto.OrderProductResponse{
			ID: pd.ID,
			SKU: pd.SKU,
			Name: pd.Name,
			Image: pd.Image,
			Price: pd.Price,
			Quantity: pd.Quantity,
		})
	}
	result.Products = product

	ctx.JSON(http.StatusCreated, result)
}