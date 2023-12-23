package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Address string `json:"address"`
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

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
    router.GET("/albums/:id/:name", getSingleAlbums)
	router.Run("localhost:8085")
}