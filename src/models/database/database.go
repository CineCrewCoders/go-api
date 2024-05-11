package database

import (
    "fmt"
    "context"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
)

var Db *mongo.Database
var Ctx context.Context

func Start() (*mongo.Database, context.Context) {
    if Db != nil {
        return Db, Ctx
    }
    uri := "mongodb://mongodb:27017"
    appCtx := context.Background()

    client, connectErr := mongo.Connect(appCtx, options.Client().ApplyURI(uri))
    if connectErr != nil {
        panic(connectErr)
        return nil, nil
    }

    pingErr := client.Ping(appCtx, readpref.Primary())
    if pingErr != nil {
        panic(pingErr)
        return nil, nil
    }

    fmt.Println("Connected to MongoDB successfully!")

    databaseName := "cinecrew"

    Db = client.Database(databaseName)
    return Db, appCtx
}
