package main

import (
    // "fmt"
    "main/models/database"
    "main/handlereq"

    // "go.mongodb.org/mongo-driver/bson"
)

func main() {
    database.Db, database.Ctx = database.Start()
    
    defer database.Db.Client().Disconnect(database.Ctx)

    // collection := database.Db.Collection("people")

    // count, err := collection.CountDocuments(database.Ctx, bson.D{})
    // if err != nil {
    //     panic(err)
    // }

    // fmt.Printf("Number of documents in collection: %d\n", count)
    handlereq.HandleRequests()
}