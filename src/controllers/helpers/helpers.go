package helpers

import (
	"main/models/movies"
    "main/models/users"
    "main/models/database"
    "log"
    "go.mongodb.org/mongo-driver/bson"
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

func IsMovieInList(userIDs primitive.ObjectID, targetID primitive.ObjectID, list string) bool {
    var movieList []primitive.ObjectID
    var userData users.User
    collection := database.Db.Collection("users")
    filter := bson.M{"_id": userIDs}
    err := collection.FindOne(database.Ctx, filter).Decode(&userData)
    if err != nil {
        log.Println(err)
    }
    if list == "Watched" {
        movieList = userData.Watched
        for _, movieID := range movieList {
            if movieID == targetID {
                return true
            }
        }
    } else if list == "PlanToWatch" {
        movieList = userData.PlanToWatch
        for _, movieID := range movieList {
            if movieID == targetID {
                return true
            }
        }
    }

    return false

}