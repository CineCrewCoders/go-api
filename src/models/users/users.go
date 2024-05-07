package users

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Watched []primitive.ObjectID `json:"watched"`
	PlanToWatch []primitive.ObjectID `json:"plan_to_watch"`
	Rated []Rated `json:"rated"`
}

type Rated struct {
	MovieID primitive.ObjectID `bson:"movie_id"`
	Score float64 `json:"score"`
}