package handlereq

import (
	"encoding/json"
	"log"
	"main/controllers/people"

	"github.com/gin-gonic/gin"
)

func HandleRequests() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.GET("/actors", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(people.GetActors()))
	})

	router.Run("0.0.0.0:6000")
	log.Println("Server started on: http://localhost:6000")
}