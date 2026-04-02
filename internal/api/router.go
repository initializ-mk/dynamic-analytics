package api

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"

	"dynamic-analytics/internal/llm"
)

// NewRouter creates and configures the Chi router.
func NewRouter(db *mongo.Database, llmClient *llm.Client) http.Handler {
	r := chi.NewRouter()

	// Middleware
	allowedOrigins := []string{"http://localhost:3000", "http://localhost:5173"}
	if extra := os.Getenv("CORS_ORIGINS"); extra != "" {
		for _, o := range strings.Split(extra, ",") {
			allowedOrigins = append(allowedOrigins, strings.TrimSpace(o))
		}
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	h := &Handler{
		db:        db,
		llmClient: llmClient,
	}

	r.Get("/api/health", h.HealthHandler)
	r.Get("/api/schema", h.SchemaHandler)
	r.Post("/api/query", h.QueryHandler)

	return r
}
