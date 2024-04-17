package handlereq

import (
	"encoding/json"
	"log"
	"main/controllers/movies"

	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.GET("/movies", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(movies.GetMovies()))
	})

	router.GET("/movie/:id", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(movies.GetMovie(c)))
	})

	router.Run("0.0.0.0:6000")
	log.Println("Server started on: http://localhost:6000")
}