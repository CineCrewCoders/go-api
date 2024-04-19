package users

import (
	"encoding/json"
	"main/models/database"
	"main/models/movies"
	"main/models/users"
	"main/controllers/helpers"
	"log"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
)

func CreateUser(userId primitive.ObjectID, username string) int {
	user := users.User{
		ID: userId,
		Username: username,
		Watched: []primitive.ObjectID{},
		PlanToWatch: []primitive.ObjectID{},
		Rated: []users.Rated{},
	}
	collection := database.Db.Collection("users")
	res, err := collection.InsertOne(database.Ctx, user)
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return http.StatusOK
}

// func GetUserById(c *gin.Context) string {
// 	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
// 	myUser := users.User{}
// 	collection := database.Db.Collection("users")
// 	filter := bson.M{"_id": id}
// 	err := collection.FindOne(database.Ctx, filter).Decode(&myUser)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	userJSON, _ := json.Marshal(myUser)
// 	return string(userJSON)
// }

func GetUserByUsername(c *gin.Context) string {
	username := c.Param("username")
	myUser := users.User{}
	collection := database.Db.Collection("users")
	filter := bson.M{"username": username}
	err := collection.FindOne(database.Ctx, filter).Decode(&myUser)
	if err != nil {
		log.Println(err)
	}
	userJSON, _ := json.Marshal(myUser)
	return string(userJSON)
}

func AddMovieToList(userID primitive.ObjectID, movieId primitive.ObjectID, list string) int {
	collection := database.Db.Collection("users")
	filter := bson.M{"_id": userID}
	if list != "Watched" && list != "PlanToWatch" {
		return http.StatusBadRequest
	}

	collectionMovies := database.Db.Collection("movies")
	filterMovie := bson.M{"_id": movieId}
	err := collectionMovies.FindOne(database.Ctx, filterMovie).Err()
	if err != nil {
		log.Println(err)
		return http.StatusNotFound
	}

	update := bson.M{}
	if list == "Watched" {
		update = bson.M{"$push": bson.M{"Watched": movieId}}
	} else {
		update = bson.M{"$push": bson.M{"PlanToWatch": movieId}}
	}
	log.Println(update)
	log.Println(list)
	_, err = collection.UpdateOne(database.Ctx, filter, update)
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest
	}
	return http.StatusOK
}

func GetUserMovieList(userID primitive.ObjectID, list string) string {
	myUser := users.User{}
	collection := database.Db.Collection("users")
	filter := bson.M{"_id": userID}
	err := collection.FindOne(database.Ctx, filter).Decode(&myUser)
	if err != nil {
		log.Fatal(err)
	}
	var movieList []primitive.ObjectID
	if list == "Watched" {
		movieList = myUser.Watched
	} else {
		movieList = myUser.PlanToWatch
	}

	finalList := []movies.Movie{}
	movieCollection := database.Db.Collection("movies")
	for _, movieID := range movieList {
		myMovie := movies.MovieDb{}
		filter := bson.M{"_id": movieID}
		err := movieCollection.FindOne(database.Ctx, filter).Decode(&myMovie)
		if err != nil {
			log.Println(err)
		}
		cast := strings.Split(myMovie.Actors, ", ")
		finalList = append(finalList, movies.Movie{
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
	moviesJSON, _ := json.Marshal(finalList)
	return string(moviesJSON)
}

func RateMovie(userID primitive.ObjectID, movieID primitive.ObjectID, score float64) int {
	log.Println("score: ", score)
    collection := database.Db.Collection("users")
    filter := bson.M{"_id": userID, "Rated": bson.M{"$not": bson.M{"$elemMatch": bson.M{"movie_id": movieID}}}}
    update := bson.M{"$push": bson.M{"Rated": bson.M{"movie_id": movieID, "score": score}}}

	collectionMovies := database.Db.Collection("movies")
	filterMovie := bson.M{"_id": movieID}
	err := collectionMovies.FindOne(database.Ctx, filterMovie).Err()
	if err != nil {
		log.Println(err)
		return http.StatusNotFound
	}

	movie := movies.MovieDb{}
	err = collectionMovies.FindOne(database.Ctx, filterMovie).Decode(&movie)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}

	movieRating := movie.Rating
	if movieRating.Average == 0 {
		movieRating.NumVotes = 1
		movieRating.Average = score
	} else {
		movieRating.NumVotes += 1
		movieRating.Average = (movieRating.Average*float64(movieRating.NumVotes-1) + score) / float64(movieRating.NumVotes)
	}
	updateMovie := bson.M{"$set": bson.M{"rating": movieRating}}
	_, err = collectionMovies.UpdateOne(database.Ctx, filterMovie, updateMovie)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
    
    result, err := collection.UpdateOne(database.Ctx, filter, update)
    if err != nil {
        log.Println(err)
        return http.StatusBadRequest
    }
    
    if result.ModifiedCount == 0 {
        return http.StatusConflict 
    }

    return http.StatusOK
}

func UpdateMovieRating(userID primitive.ObjectID, movieID primitive.ObjectID, score float64) int {
	collection := database.Db.Collection("users")
	oldScore := helpers.GetMovieScore(userID, movieID)
	filter := bson.M{"_id": userID, "Rated": bson.M{"$elemMatch": bson.M{"movie_id": movieID}}}
	update := bson.M{"$set": bson.M{"Rated.$.score": score}}

	collectionMovies := database.Db.Collection("movies")
	filterMovie := bson.M{"_id": movieID}
	err := collectionMovies.FindOne(database.Ctx, filterMovie).Err()
	if err != nil {
		log.Println(err)
		return http.StatusNotFound
	}

	movie := movies.MovieDb{}
	err = collectionMovies.FindOne(database.Ctx, filterMovie).Decode(&movie)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}

	movieRating := movie.Rating
	movieRating.Average = (movieRating.Average*float64(movieRating.NumVotes) - oldScore + score) / float64(movieRating.NumVotes)
	updateMovie := bson.M{"$set": bson.M{"rating": movieRating}}
	_, err = collectionMovies.UpdateOne(database.Ctx, filterMovie, updateMovie)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError
	}
	
	result, err := collection.UpdateOne(database.Ctx, filter, update)
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest
	}
	
	if result.ModifiedCount == 0 {
		return http.StatusConflict 
	}

	return http.StatusOK
}

func RemoveMovieFromList(userID primitive.ObjectID, movieID primitive.ObjectID, list string) int {
	collection := database.Db.Collection("users")
	filter := bson.M{"_id": userID}
	update := bson.M{}
	if list == "Watched" {
		if !helpers.IsMovieInList(userID, movieID, "Watched") {
			log.Println("Movie not in list")
			return http.StatusNotFound
		}
		update = bson.M{"$pull": bson.M{"Watched": movieID}}
	} else {
		if !helpers.IsMovieInList(userID, movieID, "PlanToWatch") {
			log.Println("Movie not in list")
			return http.StatusNotFound
		}
		update = bson.M{"$pull": bson.M{"PlanToWatch": movieID}}
	}
	
	_, err := collection.UpdateOne(database.Ctx, filter, update)
	if err != nil {
		log.Println(err)
		return http.StatusBadRequest
	}
	return http.StatusOK
}