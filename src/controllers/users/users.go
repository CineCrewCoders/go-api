package users

import (
	"encoding/json"
	"main/models/database"
	"main/models/users"
	"log"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func CreateUser(username string) int {
	user := users.User{
		Username: username,
		Watched: []primitive.ObjectID{},
		PlanToWatch: []primitive.ObjectID{},
		Rated: []users.Rated{},
	}
	collection := database.Db.Collection("users")
	res, err := collection.InsertOne(database.Ctx, user)
	if err != nil {
		log.Fatal(err)
		return http.StatusBadRequest
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return http.StatusOK
}

func GetUserById(c *gin.Context) string {
	id, _ := primitive.ObjectIDFromHex(c.Query("id"))
	myUser := users.User{}
	collection := database.Db.Collection("users")
	filter := bson.M{"_id": id}
	err := collection.FindOne(database.Ctx, filter).Decode(&myUser)
	if err != nil {
		log.Fatal(err)
	}
	userJSON, _ := json.Marshal(myUser)
	return string(userJSON)
}

func GetUserByUsername(c *gin.Context) string {
	username := c.Query("username")
	myUser := users.User{}
	collection := database.Db.Collection("users")
	filter := bson.M{"username": username}
	err := collection.FindOne(database.Ctx, filter).Decode(&myUser)
	if err != nil {
		log.Fatal(err)
	}
	userJSON, _ := json.Marshal(myUser)
	return string(userJSON)
}

func AddPlanToWatch(username string, movieId primitive.ObjectID) int {
	collection := database.Db.Collection("users")
	filter := bson.M{"username": username}
	update := bson.M{"$push": bson.M{"planToWatch": movieId}}
	_, err := collection.UpdateOne(database.Ctx, filter, update)
	if err != nil {
		log.Fatal(err)
		return http.StatusBadRequest
	}
	return http.StatusOK
}

