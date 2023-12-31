package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"crypto/sha256"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Model struct {
  ID        uint           `gorm:"primaryKey"`
  Name 	string
  Email 	string
  Password string
}

type Response struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Address string `json:"address"`
}

type RegisterUser struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

var albums = []Response {
	{Id: "1", Name: "Album 1", Address: "Address 1"},
	{Id: "2", Name: "Album 2", Address: "Address 2"},
	{Id: "3", Name: "Album 3", Address: "Address 3"},
}


func getSingleAlbums(c *gin.Context) {
	id := c.Param("id")
	var data Response
	for i := 0; i < len(albums); i++ {
		if albums[i].Id == id {
			data = albums[i]
		}
	}
	c.IndentedJSON(http.StatusOK, data)
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func ConvertToSha256(toConvert string) string {
 	conversion := sha256.Sum256([]byte(toConvert))
 	return fmt.Sprintf("%x\n", conversion)
}

func registerUser(c *gin.Context) {
		var user RegisterUser
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		db := c.MustGet("db").(*gorm.DB)

		if(user.Name == "" || user.Email == "" || user.Password == "") {
			c.JSON(400, gin.H{"error": "Please provide all the fields"})
			return
		}

	//  passwordHash := ConvertToSha256(user.Password)

		if(len(user.Password) < 8) {
			c.JSON(400, gin.H{"error": "Password should be at least 8 characters long"})
			return
		}

		userData := Model{Name: user.Name, Email: user.Email, Password: user.Password}

	    db.Create(&userData)

		fmt.Printf("Received User:", user)
     	c.IndentedJSON(http.StatusOK, gin.H{"message": "User Registered Successfully!!"})
}

func loginUser(c *gin.Context) {
		var userLogin LoginUser
		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		db := c.MustGet("db").(*gorm.DB)

		var user Model
		if err := db.First(&user, "Email = ?", userLogin.Email).Error; err != nil {
			c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
		}

		if user.Password != userLogin.Password {
			c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
		}

    	c.IndentedJSON(http.StatusOK, gin.H{"message": "User Logged In Successfully!!"})
}

type JsonFormatter struct {
    UserId int `json:"userId"`
    Id int `json:"id"`
    Title string `json:"title"`
    Body string `json:"body"`

}

func goAndCallJsonFormatterApi(c *gin.Context) {
	var jsonFormatter []JsonFormatter

	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from the external API"})
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&jsonFormatter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode JSON response"})
		return
	}

	c.IndentedJSON(http.StatusOK, jsonFormatter)
}

func main() {

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host=localhost user=DBusername password=DBpassword dbname=DBname sslmode=disable TimeZone=Asia/Kolkata",
		PreferSimpleProtocol: true,
	}), &gorm.Config{})


	if err != nil {
		fmt.Println(err, "Failed to connect to the database")
		panic("Failed to connect to the database")
	}
	
	db.AutoMigrate(&Model{})

	router := gin.Default()

	router.Use(func(c *gin.Context) {
 	c.Set("db", db)
	c.Next()
	})

	router.GET("/albums", getAlbums)
    router.GET("/albums/:id/:name", getSingleAlbums)
    router.POST("/register", registerUser)
    router.POST("/login", loginUser)
	
	router.GET("/getRandomData", goAndCallJsonFormatterApi)
	router.Run("localhost:8085")
}