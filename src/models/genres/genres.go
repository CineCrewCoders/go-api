package genres

import "go.mongodb.org/mongo-driver/bson/primitive"

type Genre struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string `json:"name"`
}