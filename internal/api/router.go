package api

import (
	"net/http"
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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
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
