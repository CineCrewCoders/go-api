package people

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type People struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name string `json:"name"`
	Birthday time.Time `json:"birthday"`
	Profession string `json:"profession"`
}
