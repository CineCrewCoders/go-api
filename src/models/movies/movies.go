package movies

import (
	// "main/models/genres"
	"main/models/people"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rating struct {
	NumVotes int     `json:"num_votes"`
	Average  float64 `json:"average"`
}

type Movie struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string `json:"name"`
	Year int `json:"year"`
	// Genres []genres.Genre `json:"genres"`
	Genres []string `json:"genres"`
	Cast []people.People `json:"cast"`
	Director people.People `json:"director"`
	Description string `json:"description"`
	// Poster string `json:"poster"`
	Rating Rating `json:"rating"`
}

type MovieDb struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string `json:"name"`
	Year int `json:"year"`
	Genres []string `json:"genres"`
	Cast []string `json:"cast"`
	Director string `json:"director"`
	Description string `json:"description"`
	// Poster string `json:"poster"`
	Rating Rating `json:"rating"`
}
