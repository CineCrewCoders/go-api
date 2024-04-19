package movies

import (
	"encoding/json"	
	"log"
	"main/models/database"
	"main/models/movies"
	"main/controllers/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
	"strconv"
)

func GetMovies() string {
	moviesSlice := []movies.Movie{}
	collection := database.Db.Collection("movies")
	cursor, err := collection.Find(database.Ctx, bson.M{})
	if err != nil {
		log.Println(err)
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

func GetMovieById(c *gin.Context) string {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	myMovie := movies.MovieDb{}
	collection := database.Db.Collection("movies")
	filter := bson.M{"_id": id}
	err := collection.FindOne(database.Ctx, filter).Decode(&myMovie)
	if err != nil {
		log.Println(err)
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

func SearchMovies(c *gin.Context) string {
	title := c.Query("title")
	genres := c.QueryArray("genres")
	minScore := c.Query("min_score")
	allMovies := []movies.Movie{}
	collection := database.Db.Collection("movies")
	cursor, err := collection.Find(database.Ctx, bson.M{})
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(database.Ctx)
	for cursor.Next(database.Ctx) {
		var myMovie movies.MovieDb
		cursor.Decode(&myMovie)
		cast := strings.Split(myMovie.Actors, ", ")
		allMovies = append(allMovies, movies.Movie{
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
	filteredMovies := allMovies
	if title != "" {
		filteredMovies = SearchMoviesByTitle(title, filteredMovies)
	}
	if len(genres) > 0 {
		filteredMovies = SearchMoviesByGenre(genres, filteredMovies)
	}
	if minScore != "" {
		minScoreFloat, err := strconv.ParseFloat(minScore, 64)
		if err != nil {
			log.Println(err)
		}
		filteredMovies = SearchMoviesByMinScore(minScoreFloat, filteredMovies)
	}
	filteredMoviesJSON, _ := json.Marshal(filteredMovies)
	return string(filteredMoviesJSON)
}


func SearchMoviesByTitle(title string, allMovies []movies.Movie) []movies.Movie {
	moviesSlice := []movies.Movie{}
	collection := database.Db.Collection("movies")
	cursor, err := collection.Find(database.Ctx, bson.M{"title": primitive.Regex{Pattern: title, Options: "i"}})
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(database.Ctx)
	for cursor.Next(database.Ctx) {
		var myMovie movies.MovieDb
		cursor.Decode(&myMovie)
		cast := strings.Split(myMovie.Actors, ", ")
		if helpers.ContainsMovieID(allMovies, myMovie.ID) {
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
	}
	return moviesSlice
}

func SearchMoviesByGenre(genre []string, allMovies []movies.Movie) []movies.Movie {
	genreFilter := bson.M{"genres": bson.M{"$in": genre}}
	moviesSlice := []movies.Movie{}
	cursor, err := database.Db.Collection("movies").Find(database.Ctx, genreFilter)
	if err != nil {
		log.Println(err)
	}
	defer cursor.Close(database.Ctx)
	for cursor.Next(database.Ctx) {
		var myMovie movies.MovieDb
		cursor.Decode(&myMovie)
		cast := strings.Split(myMovie.Actors, ", ")
		if helpers.ContainsMovieID(allMovies, myMovie.ID) {
			newMovie := movies.Movie{
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
			moviesSlice = append(moviesSlice, newMovie)
		}
	}

	return moviesSlice
}

func SearchMoviesByMinScore(minScore float64, allMovies []movies.Movie) []movies.Movie {
    moviesSlice := []movies.Movie{}
    cursor, err := database.Db.Collection("movies").Find(database.Ctx, bson.M{"rating.average": bson.M{"$gte": minScore}})
    if err != nil {
        log.Println(err)
    }
    defer cursor.Close(database.Ctx)

    for cursor.Next(database.Ctx) {
        var myMovie movies.MovieDb
        if err := cursor.Decode(&myMovie); err != nil {
            log.Println(err)
        }

        if helpers.ContainsMovieID(allMovies, myMovie.ID) {
            cast := strings.Split(myMovie.Actors, ", ")
            newMovie := movies.Movie{
                ID:        myMovie.ID,
                Title:     myMovie.Title,
                Year:      myMovie.Year,
                Runtime:   myMovie.Runtime,
                Genres:    myMovie.Genres,
                Actors:    cast,
                Director:  myMovie.Director,
                Plot:      myMovie.Plot,
                PosterURL: myMovie.PosterURL,
                Rating:    myMovie.Rating,
            }
            moviesSlice = append(moviesSlice, newMovie)
        }
    }

    return moviesSlice
}
