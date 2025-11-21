package main

import (
	"PROJECTUAS_BE/config"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config.ConnectMongo()
	// ===============================
	// 🟨 Connect to PostgreSQL
	// ===============================
	pgDB := config.ConnectPG()
	defer pgDB.Close()
	fmt.Println("PostgreSQL Connected via config.ConnectPG()")

	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// ===============================
	// 🟨 Connect to MongoDB
	// ===============================

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")
	port := os.Getenv("SERVER_PORT")

	if mongoURI == "" || dbName == "" {
		log.Fatal("❌ MONGO_URI or MONGO_DB not found in .env")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("❌ MongoDB connection error:", err)
	}

	// 🔹 Tes koneksi
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Cannot connect to MongoDB:", err)
	}
	fmt.Println("✅ Connected to MongoDB!")

	if port == "" {
		port = "3000"
	}
	fmt.Printf("🚀 Server running on port %s\n", port)
	// router.Run(":" + port)
}
