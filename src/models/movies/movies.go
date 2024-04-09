import (
	"models/genres"
	"models/people"
)
package movies

type Movie struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Year int `json:"year"`
	Genres []Genre `json:"genres"`
	Cast []People `json:"cast"`
	Directors People `json:"director"`
	Description string `json:"description"`
	Poster string `json:"poster"`
	Rating (int, float64) `json:"rating"`
}
