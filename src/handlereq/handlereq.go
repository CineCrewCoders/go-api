package handlereq

import (
	"encoding/json"
	"log"
	"main/controllers/movies"
	"main/controllers/users"
	"main/controllers/helpers"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type SignUpRequest struct {
	UserId string `json:"userId"`
    Username string `json:"username"`
}

type MovieIDRequest struct {
	List string `json:"list"`
	MovieID string `json:"movieId"`
}

type Rated struct {
	MovieID string `json:"movieId"`
	Score float64 `json:"score"`
}

func HandleRequests() {
	router := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.Use(cors.Default())

	router.GET("/movies", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(movies.GetMovies()))
	})

	router.GET("/movies/:id", func(c *gin.Context) {
		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			return
		}

		userID := c.Request.Header.Get("UserId")
		
		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			return
		}

		c.JSON(200, json.RawMessage(movies.GetMovieById(userID, c)))
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

		if signUpRequest.UserId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required"})
			return
		}

		// userId is string
		userID := signUpRequest.UserId
	
		if signUpRequest.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
			return
		}

		res := users.CreateUser(userID , signUpRequest.Username)
		if res == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
		} else if res == http.StatusBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate id or username"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	})

	// router.GET("/user/:id", func(c *gin.Context) {
	// 	c.JSON(200, json.RawMessage(users.GetUserById(c)))
	// })

	router.GET("/user/:username", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(users.GetUserByUsername(c)))
	})

	router.POST("/movies/list", func(c *gin.Context) {
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type header must be application/json"})
			return
		}
	
		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			return
		}

		var movieIDRequest MovieIDRequest
		if err := c.ShouldBindJSON(&movieIDRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}
	
		if movieIDRequest.List == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "List is required"})
			return
		}
		log.Println(movieIDRequest.List)
		movieID, err := primitive.ObjectIDFromHex(movieIDRequest.MovieID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movieId format"})
			return
		}
	
		res := users.AddMovieToList(userID, movieID, movieIDRequest.List)
		if res == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{"message": "Movie added to plan to watch list successfully"})
		} else if res == http.StatusBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		} else if res == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	})

	router.GET("movies/plan_to_watch", func(c *gin.Context) {
		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			return
		}

		c.JSON(200, json.RawMessage(users.GetUserMovieList(userID, "PlanToWatch")))
	})

	router.GET("movies/watched", func(c *gin.Context) {
		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			return
		}

		c.JSON(200, json.RawMessage(users.GetUserMovieList(userID, "Watched")))
	})

	router.POST("movies/rate", func(c *gin.Context) {
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type header must be application/json"})
			return
		}

		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			return
		}

		var rated Rated
		if err := c.ShouldBindJSON(&rated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		if rated.MovieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "MovieId is required"})
			return
		}

		movieID, err := primitive.ObjectIDFromHex(rated.MovieID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movieId format"})
			return
		}

		if rated.Score < 1 || rated.Score > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Score must be between 0 and 10"})
			return
		}

		res := users.RateMovie(userID, movieID, rated.Score)
		if res == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{"message": "Movie rated successfully"})
		} else if res == http.StatusBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		} else if res == http.StatusConflict {
			c.JSON(http.StatusConflict, gin.H{"error": "Movie already rated"})
		} else if res == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	})

	router.PUT("movies/rate", func(c *gin.Context) {
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type header must be application/json"})
			return
		}

		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			return
		}

		var rated Rated
		if err := c.ShouldBindJSON(&rated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		if rated.MovieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "MovieId is required"})
			return
		}

		movieID, err := primitive.ObjectIDFromHex(rated.MovieID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movieId format"})
			return
		}

		if rated.Score < 1 || rated.Score > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Score must be between 0 and 10"})
			return
		}

		res := users.UpdateMovieRating(userID, movieID, rated.Score)
		if res == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{"message": "Movie rating updated successfully"})
		} else if res == http.StatusBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		} else if res == http.StatusConflict {
			c.JSON(http.StatusConflict, gin.H{"error": "Movie not rated yet or the score is the same"})
		} else if res == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	})

	router.DELETE("movies/watched", func(c *gin.Context) {
		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			return
		}

		if c.Query("movieId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "MovieId is required"})
			return
		}

		movieID, err := primitive.ObjectIDFromHex(c.Query("movieId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movieId format"})
			return
		}

		res := users.RemoveMovieFromList(userID, movieID, "Watched")
		if res == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{"message": "Movie removed from watched list successfully"})
		} else if res == http.StatusBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		} else if res == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found in the list"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	})

	router.DELETE("movies/plan_to_watch", func(c *gin.Context) {
		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			return
		}

		if c.Query("movieId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "MovieId is required"})
			return
		}

		movieID, err := primitive.ObjectIDFromHex(c.Query("movieId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movieId format"})
			return
		}

		res := users.RemoveMovieFromList(userID, movieID, "PlanToWatch")
		if res == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{"message": "Movie removed from plan to watch list successfully"})
		} else if res == http.StatusBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		} else if res == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found in the list"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	})


	router.Run("0.0.0.0:5678")
	log.Println("Server started on: http://localhost:5678")
}