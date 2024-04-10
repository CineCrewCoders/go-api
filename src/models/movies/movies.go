import (
	"models/genres"
	"models/people"
	"go.mongodb.org/mongo-driver/bson/primitive"
)
package movies

type Movie struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Title string `json:"title"`
	Year int `json:"year"`
	Genres []Genre `json:"genres"`
	Cast []People `json:"cast"`
	Director People `json:"director"`
	Description string `json:"description"`
	// Poster string `json:"poster"`
	Rating (int, float64) `json:"rating"`
}
