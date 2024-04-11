package movies

import (
	"encoding/json"	
	"log"
	"main/models/database"
	"main/models/movies"
	"main/models/people"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetMovie(c *gin.Context) string {
	collection := database.Db.Collection("movies")
	
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	cursor, err := collection.Find(database.Ctx, bson.D{{"_id", objID}})

	if err != nil {
		log.Fatal(err)
	}

	var movie movies.Movie
	for cursor.Next(database.Ctx) {
		var movieDb movies.MovieDb
		cursor.Decode(&movieDb)
		var director people.People
		director = GetPerson(movieDb.Director)
		var cast []people.People
		for _, actor := range movieDb.Cast {
			cast = append(cast, GetPerson(actor))
		}
		movie = movies.Movie{
			ID: movieDb.ID,
			Name: movieDb.Name,
			Year: movieDb.Year,
			Genres: movieDb.Genres,
			Cast: cast,
			Director: director,
			Description: movieDb.Description,
			Rating: movieDb.Rating,
		}
	}

	movieJson, err := json.Marshal(movie)
	if err != nil {
		log.Fatal(err)
	}

	return string(movieJson)
}

func GetPerson(name string) people.People {
	collection := database.Db.Collection("people")
	
	cursor, err := collection.Find(database.Ctx, bson.D{{"name", name}})
	if err != nil {
		log.Fatal(err)
	}

	var person people.People
	for cursor.Next(database.Ctx) {
		cursor.Decode(&person)
	}

	return person
}

func GetMovies() string {
	collection := database.Db.Collection("movies")
	
	cursor, err := collection.Find(database.Ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	var moviesArr []movies.Movie
	for cursor.Next(database.Ctx) {
		var movie movies.MovieDb
		cursor.Decode(&movie)
		var director people.People
		director = GetPerson(movie.Director)
		var cast []people.People
		for _, actor := range movie.Cast {
			cast = append(cast, GetPerson(actor))
		}
		moviesArr = append(moviesArr, movies.Movie{
			ID: movie.ID,
			Name: movie.Name,
			Year: movie.Year,
			Genres: movie.Genres,
			Cast: cast,
			Director: director,
			Description: movie.Description,
			Rating: movie.Rating,
		})
	}
	
	moviesJson, err := json.Marshal(moviesArr)
	if err != nil {
		log.Fatal(err)
	}

	return string(moviesJson)
}
