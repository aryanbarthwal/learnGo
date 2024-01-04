package main

import (
	"Go-Proj-01/handlers"
	"Go-Proj-01/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func main() {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host=localhost user=DBusername password=DBpassword dbname=DBname sslmode=disable TimeZone=Asia/Kolkata",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})


	if err != nil {
		fmt.Println(err, "Failed to connect to the database")
		panic("Failed to connect to the database")
	}
	
	db.AutoMigrate(&models.Model{})

	router := gin.Default()

	router.Use(func(c *gin.Context) {
 	c.Set("db", db)
	c.Next()
	})

	router.GET("/albums", handlers.GetAlbums)
    router.GET("/albums/:id/:name", handlers.GetSingleAlbums)
    router.POST("/register", handlers.RegisterUser)
    router.POST("/login", handlers.LoginUser)
	
	router.GET("/getRandomData", handlers.GoAndCallJsonFormatterApi)
	router.Run("localhost:8085")
}