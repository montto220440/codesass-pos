package main

import (
	"api-pos/db"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") == "production"{
		gin.SetMode(gin.ReleaseMode)
	}else{
		if err := godotenv.Load(); err != nil{
			log.Fatal("Error loading .env file")
		}
	}
	if err := godotenv.Load(); err != nil{
		log.Fatal("Error loading env")
	}
	//connect db
	db.Connectdb()
	db.Migrate()

	cors_donfig := cors.DefaultConfig()
	cors_donfig.AllowAllOrigins = true

	//create folder
	os.MkdirAll("uploads/products", 0755)
	r := gin.Default()
	r.Use(cors.New(cors_donfig))
	r.Static("/uploads", "./uploads")
	serveRoutes(r)
	port := os.Getenv("POST")
	if port != ""{
		port = "5000"
	}
	r.Run(":" + port)
}

