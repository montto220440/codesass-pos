package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Name string 
	Email string 
	Tel string 
	Product []Order_item `gorm:"foreignKey:OrderID"`
}