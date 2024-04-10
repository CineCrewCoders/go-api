package people

import (
	"encoding/json"	
	"log"
	"main/models/database"
	"main/models/people"

	"go.mongodb.org/mongo-driver/bson"
)

func GetActors() string {
	collection := database.Db.Collection("people")
	
	cursor, err := collection.Find(database.Ctx, bson.D{{"profession", "actor"}})
	if err != nil {
		log.Fatal(err)
	}

	var actors []people.People
	for cursor.Next(database.Ctx) {
		var actor people.People
		cursor.Decode(&actor)
		actors = append(actors, actor)
	}

	actorsJson, err := json.Marshal(actors)
	if err != nil {
		log.Fatal(err)
	}

	return string(actorsJson)
}