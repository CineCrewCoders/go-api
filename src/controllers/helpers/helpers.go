package helpers

import (
	"main/models/movies"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ContainsMovieID(movies []movies.Movie, targetID primitive.ObjectID) bool {
    for _, movie := range movies {
        if movie.ID == targetID {
            return true
        }
    }
    return false
}