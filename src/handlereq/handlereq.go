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
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
	
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_request_duration_seconds",
		Help: "Duration of HTTP requests in seconds.",
	},
	[]string{"method", "endpoint"},
)

// Register the metric with Prometheus
func init() {
	prometheus.MustRegister(httpDuration)
}

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
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"UserId", "Content-Type"},
	}))


	
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	router.GET("/movies", func(c *gin.Context) {
		start := time.Now()
		c.JSON(200, json.RawMessage(movies.GetMovies()))
		elapsed := time.Since(start).Seconds()
		log.Println(elapsed)
		httpDuration.WithLabelValues("GET", "/movies").Observe(elapsed)
	})

	router.GET("/movies/:id", func(c *gin.Context) {
		start := time.Now()
		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("GET", "/movies/:id").Observe(elapsed)
			return
		}

		userID := c.Request.Header.Get("UserId")
		
		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("GET", "/movies/:id").Observe(elapsed)
			return
		}

		c.JSON(200, json.RawMessage(movies.GetMovieById(userID, c)))
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("GET", "/movies/:id").Observe(elapsed)
	})

	router.GET("/search", func(c *gin.Context) {
		start := time.Now()
		c.JSON(200, json.RawMessage(movies.SearchMovies(c)))
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("GET", "/search").Observe(elapsed)
	})

	router.POST("/signup", func(c *gin.Context) {
		start := time.Now()
		if c.Request.Header.Get("Content-Type") != "" && c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type header must be application/json"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/signup").Observe(elapsed)
			return
		}
	
		var signUpRequest SignUpRequest
		if err := c.ShouldBindJSON(&signUpRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/signup").Observe(elapsed)
			return
		}

		if signUpRequest.UserId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/signup").Observe(elapsed)
			return
		}

		// userId is string
		userID := signUpRequest.UserId
	
		if signUpRequest.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/signup").Observe(elapsed)
			return
		}

		res := users.CreateUser(userID , signUpRequest.Username)
		if res == http.StatusOK {
			c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/signup").Observe(elapsed)
		} else if res == http.StatusBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate id or username"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/signup").Observe(elapsed)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/signup").Observe(elapsed)
		}
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("POST", "/signup").Observe(elapsed)
	})

	// router.GET("/user/:id", func(c *gin.Context) {
	// 	c.JSON(200, json.RawMessage(users.GetUserById(c)))
	// })

	router.GET("/user/:username", func(c *gin.Context) {
		start := time.Now()
		c.JSON(200, json.RawMessage(users.GetUserByUsername(c)))
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("GET", "/user/:username").Observe(elapsed)
	})

	router.POST("/movies/list", func(c *gin.Context) {
		start := time.Now()
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type header must be application/json"})
			return
		}
	
		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/list").Observe(elapsed)
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/list").Observe(elapsed)
			return
		}

		var movieIDRequest MovieIDRequest
		if err := c.ShouldBindJSON(&movieIDRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/list").Observe(elapsed)
			return
		}
	
		if movieIDRequest.List == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "List is required"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/list").Observe(elapsed)
			return
		}
		log.Println(movieIDRequest.List)
		movieID, err := primitive.ObjectIDFromHex(movieIDRequest.MovieID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movieId format"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/list").Observe(elapsed)
			return
		}
	
		res := users.AddMovieToList(userID, movieID, movieIDRequest.List)
		if res == http.StatusOK {
			if movieIDRequest.List == "Watched" {
				c.JSON(http.StatusOK, gin.H{"message": "Movie added to watched list successfully"})
			} else {
				c.JSON(http.StatusOK, gin.H{"message": "Movie added to plan to watch list successfully"})
			}
		} else if res == http.StatusBadRequest {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
		} else if res == http.StatusNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("POST", "/movies/list").Observe(elapsed)
	})

	router.GET("movies/plan_to_watch", func(c *gin.Context) {
		start := time.Now()
		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("GET", "/movies/plan_to_watch").Observe(elapsed)
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("GET", "/movies/plan_to_watch").Observe(elapsed)
			return
		}

		c.JSON(200, json.RawMessage(users.GetUserMovieList(userID, "PlanToWatch")))
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("GET", "/movies/plan_to_watch").Observe(elapsed)
	})

	router.GET("movies/watched", func(c *gin.Context) {
		start := time.Now()
		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("GET", "/movies/watched").Observe(elapsed)
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("GET", "/movies/watched").Observe(elapsed)
			return
		}

		c.JSON(200, json.RawMessage(users.GetUserMovieList(userID, "Watched")))
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("GET", "/movies/watched").Observe(elapsed)
	})

	router.POST("movies/rate", func(c *gin.Context) {
		start := time.Now()
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type header must be application/json"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/rate").Observe(elapsed)
			return
		}

		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/rate").Observe(elapsed)
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/rate").Observe(elapsed)
			return
		}

		var rated Rated
		if err := c.ShouldBindJSON(&rated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/rate").Observe(elapsed)
			return
		}

		if rated.MovieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "MovieId is required"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/rate").Observe(elapsed)
			return
		}

		movieID, err := primitive.ObjectIDFromHex(rated.MovieID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movieId format"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/rate").Observe(elapsed)
			return
		}

		if rated.Score < 1 || rated.Score > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Score must be between 1 and 10"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("POST", "/movies/rate").Observe(elapsed)
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
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("POST", "/movies/rate").Observe(elapsed)
	})

	router.PUT("movies/rate", func(c *gin.Context) {
		start := time.Now()
		if c.Request.Header.Get("Content-Type") != "application/json" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type header must be application/json"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("PUT", "/movies/rate").Observe(elapsed)
			return
		}

		if c.Request.Header.Get("UserId") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UserId is required in the header"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("PUT", "/movies/rate").Observe(elapsed)
			return
		}

		userID := c.Request.Header.Get("UserId")

		if !helpers.UserExists(userID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User does not exist"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("PUT", "/movies/rate").Observe(elapsed)
			return
		}

		var rated Rated
		if err := c.ShouldBindJSON(&rated); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("PUT", "/movies/rate").Observe(elapsed)
			return
		}

		if rated.MovieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "MovieId is required"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("PUT", "/movies/rate").Observe(elapsed)
			return
		}

		movieID, err := primitive.ObjectIDFromHex(rated.MovieID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movieId format"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("PUT", "/movies/rate").Observe(elapsed)
			return
		}

		if rated.Score < 1 || rated.Score > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Score must be between 0 and 10"})
			elapsed := time.Since(start).Seconds()
			httpDuration.WithLabelValues("PUT", "/movies/rate").Observe(elapsed)
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
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("PUT", "/movies/rate").Observe(elapsed)
	})

	router.DELETE("movies/watched", func(c *gin.Context) {
		start := time.Now()
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
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("DELETE", "/movies/watched").Observe(elapsed)
	})

	router.DELETE("movies/plan_to_watch", func(c *gin.Context) {
		start := time.Now()
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
		elapsed := time.Since(start).Seconds()
		httpDuration.WithLabelValues("DELETE", "/movies/plan_to_watch").Observe(elapsed)
	})


	router.Run("0.0.0.0:5678")
	log.Println("Server started on: http://0.0.0.0:5678")
}