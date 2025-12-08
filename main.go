package main

import (
	"PROJECTUAS_BE/app/repository"
	"PROJECTUAS_BE/app/service"
	"PROJECTUAS_BE/config"

	// "PROJECTUAS_BE/middleware"
	"PROJECTUAS_BE/routes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// ===============================
	// Memproses set env pada file env
	// ===============================
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// ===============================
	// 🟨 Connect to PostgreSQL
	// ===============================
	pgDB := config.ConnectPG()
	defer pgDB.Close()
	fmt.Println("PostgreSQL Connected via config.ConnectPG()")

	// ===============================
	// 🟨 Init Auth + Generate Sample Token
	// ===============================
	app := fiber.New()
	userRepo := repository.NewUserRepository(pgDB)
	authService := service.NewAuthService(userRepo)

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

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Cannot connect to MongoDB:", err)
	}
	fmt.Println("✅ Connected to MongoDB!")

	// ===============================
	// 🟨 Setup Routes
	// ===============================
	routes.SetupRoutes(app, authService)

	// ===============================
	// 🟨 Run Server
	// ===============================
	if port == "" {
		port = "3000"
	}
	fmt.Printf("🚀 Server running on port %s\n", port)

	log.Fatal(app.Listen(":" + port))
}
