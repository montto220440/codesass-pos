package models

import "gorm.io/gorm"

type Order_item struct {
	gorm.Model
	SKU string `gorm:"not null"`
	Name string `gorm:"not null"`
	Image string `gorm:"not null"`
	Price float64 `gorm:"not null"`
	Quantity uint `gorm:"not null"`
	OrderID uint `gorm:"not null"`
}