package main

import (
    "fmt"
    "time"
    "context"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    "go.mongodb.org/mongo-driver/bson"
)

func main() {
    // fmt.Println("Hello, World!")
    uri := "mongodb://localhost:27017"
    appCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

    client, connectErr := mongo.Connect(appCtx, options.Client().ApplyURI(uri))
    if connectErr != nil {
        panic(connectErr)
    }

    pingErr := client.Ping(appCtx, readpref.Primary())
    if pingErr != nil {
        panic(pingErr)
    }

    fmt.Println("Connected to MongoDB successfully!")
    databaseName := "cinecrew"
    collectionName := "movies"

    // Access the specified database
    database := client.Database(databaseName)

    // Access the specified collection
    collection := database.Collection(collectionName)

    // Count the number of documents in the collection
    count, err := collection.CountDocuments(appCtx, bson.D{})
    if err != nil {
        panic(err)
    }

    fmt.Printf("Number of documents in '%s' collection: %d\n", collectionName, count)


    defer cancel()
    defer func() {
        disconnectErr := client.Disconnect(appCtx)
        if disconnectErr != nil {
            panic(disconnectErr)
        }
    } ()

}