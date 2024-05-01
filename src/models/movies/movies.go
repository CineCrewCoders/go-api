package movies

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieDb struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Title string `json:"title"`
	Year string `json:"year"`
	Runtime string `json:"runtime"`
	Genres []string `json:"genres"`
	Actors string `json:"actors"`
	Director string `json:"director"`
	Plot string `json:"plot"`
	PosterURL string `json:"poster_url"`
	Rating Rating `json:"rating"`
}

type Movie struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Title string `json:"title"`
	Year string `json:"year"`
	Runtime string `json:"runtime"`
	Genres []string `json:"genres"`
	Actors []string `json:"actors"`
	Director string `json:"director"`
	Plot string `json:"plot"`
	PosterURL string `json:"poster_url"`
	Rating Rating `json:"rating"`
	Score float64 `json:"score,omitempty"`
}

type Rating struct {
	NumVotes int     `json:"num_votes"`
	Average  float64 `json:"average"`
}
