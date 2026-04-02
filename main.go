package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"dynamic-analytics/internal/api"
	"dynamic-analytics/internal/llm"
)

func main() {
	// Load .env file (ignore error if not present)
	godotenv.Load()

	// Configuration from env vars
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	port := getEnv("PORT", "8080")
	apiKey := os.Getenv("ANTHROPIC_API_KEY")

	if apiKey == "" {
		log.Fatal("ANTHROPIC_API_KEY is required")
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	slog.Info("connected to MongoDB", "uri", mongoURI)

	db := client.Database("recruitment")

	// Initialize LLM client
	llmClient := llm.NewClient(apiKey)

	// Create router and start server
	router := api.NewRouter(db, llmClient)

	slog.Info("starting server", "port", port)
	fmt.Printf("Server running on http://localhost:%s\n", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
