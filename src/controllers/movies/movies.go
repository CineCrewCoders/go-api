package movies

import (
	"encoding/json"	
	"log"
	"main/models/database"
	"main/models/movies"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

func GetMovies() string {
	moviesSlice := []movies.Movie{}
	collection := database.Db.Collection("movies")
	cursor, err := collection.Find(database.Ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(database.Ctx)
	for cursor.Next(database.Ctx) {
		var myMovie movies.MovieDb
		cursor.Decode(&myMovie)
		cast := strings.Split(myMovie.Actors, ", ")
		moviesSlice = append(moviesSlice, movies.Movie{
			ID: myMovie.ID,
			Title: myMovie.Title,
			Year: myMovie.Year,
			Runtime: myMovie.Runtime,
			Genres: myMovie.Genres,
			Actors: cast,
			Director: myMovie.Director,
			Plot: myMovie.Plot,
			PosterURL: myMovie.PosterURL,
			Rating: myMovie.Rating,
		})
	}
	moviesJSON, _ := json.Marshal(moviesSlice)
	return string(moviesJSON)
}

func GetMovie(c *gin.Context) string {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	myMovie := movies.MovieDb{}
	collection := database.Db.Collection("movies")
	filter := bson.M{"_id": id}
	err := collection.FindOne(database.Ctx, filter).Decode(&myMovie)
	if err != nil {
		log.Fatal(err)
	}

	cast := strings.Split(myMovie.Actors, ", ")
	movie := movies.Movie{
		ID: myMovie.ID,
		Title: myMovie.Title,
		Year: myMovie.Year,
		Runtime: myMovie.Runtime,
		Genres: myMovie.Genres,
		Actors: cast,
		Director: myMovie.Director,
		Plot: myMovie.Plot,
		PosterURL: myMovie.PosterURL,
		Rating: myMovie.Rating,
	}
	movieJSON, _ := json.Marshal(movie)
	return string(movieJSON)
}
