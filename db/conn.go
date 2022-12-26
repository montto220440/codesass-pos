package db

import (
	"log"
	"api-pos/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Conn *gorm.DB

func Connectdb() {
	dsn := "root@tcp(127.0.0.1:3306)/pos_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(
		postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)},
	)
	if err!=nil{
		log.Fatal("connect db failed")
	}
	Conn = db
}

func Migrate(){
	Conn.AutoMigrate(
		&models.Category{},
		&models.Product{},
		&models.Order{},
		&models.Order_item{},
	)
}