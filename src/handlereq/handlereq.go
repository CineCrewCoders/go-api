package handlereq

import (
	"encoding/json"
	"log"
	"main/controllers/movies"
	"main/controllers/users"

	"github.com/gin-gonic/gin"
	"net/http"
)

type SignUpRequest struct {
    Username string `json:"username"`
    // Add other fields as needed
}

func HandleRequests() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	router.GET("/movies", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(movies.GetMovies()))
	})

	router.GET("/movie/:id", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(movies.GetMovieById(c)))
	})

	router.GET("/search", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(movies.SearchMovies(c)))
	})

	router.POST("/signup", func(c *gin.Context) {
		if c.Request.Header.Get("Content-Type") != "" && c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type header must be application/json"})
			return
		}
	
		var signUpRequest SignUpRequest
		if err := c.ShouldBindJSON(&signUpRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}
	
		if signUpRequest.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
			return
		}

		res := users.CreateUser(signUpRequest.Username)
		if res == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
		} else if res == http.StatusBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	})

	router.GET("/user", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(users.GetUserById(c)))
	})

	router.Run("0.0.0.0:6000")
	log.Println("Server started on: http://localhost:6000")
}