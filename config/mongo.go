package config

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "github.com/joho/godotenv"
)

var Client *mongo.Client

func ConnectMongo() *mongo.Database {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found or couldn't load it — relying on env vars")
    }

    uri := os.Getenv("MONGO_URI")
    dbName := os.Getenv("MONGO_DB")
    if uri == "" || dbName == "" {
        log.Fatal("MONGO_URI or MONGO_DB not set in environment")
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOpts := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        log.Fatalf("mongo connect error: %v", err)
    }

    // ping
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatalf("mongo ping error: %v", err)
    }

    Client = client
    fmt.Println("✅ Connected to MongoDB")
    return client.Database(dbName)
}
