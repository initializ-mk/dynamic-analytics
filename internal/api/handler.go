package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"dynamic-analytics/internal/llm"
	"dynamic-analytics/internal/query"
	"dynamic-analytics/internal/schema"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	db        *mongo.Database
	llmClient *llm.Client
}

// QueryRequest is the incoming request body for POST /api/query.
type QueryRequest struct {
	Query string `json:"query"`
}

// QueryResponse is the API response for query results.
type QueryResponse struct {
	Success       bool             `json:"success"`
	UIType        string           `json:"ui_type,omitempty"`
	ChartConfig   *llm.ChartConfig `json:"chart_config,omitempty"`
	Columns       []llm.Column     `json:"columns,omitempty"`
	StatConfig    *llm.StatConfig  `json:"stat_config,omitempty"`
	Data          []map[string]any `json:"data,omitempty"`
	Title         string           `json:"title,omitempty"`
	Summary       string           `json:"summary,omitempty"`
	GenPipeline   any              `json:"generated_pipeline,omitempty"`
	Meta          *Meta            `json:"meta,omitempty"`
	Error         string           `json:"error,omitempty"`
}

// Meta holds execution metadata.
type Meta struct {
	ExecutionTimeMs int64 `json:"execution_time_ms"`
	ResultCount     int   `json:"result_count"`
}

// HealthHandler returns a simple health check response.
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// SchemaHandler returns the database schema definitions.
func (h *Handler) SchemaHandler(w http.ResponseWriter, r *http.Request) {
	schemas := []schema.CollectionSchema{
		schema.GetCandidatesSchema(),
		schema.GetRecruitersSchema(),
	}
	writeJSON(w, http.StatusOK, schemas)
}

// QueryHandler processes natural language queries via the LLM.
func (h *Handler) QueryHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, QueryResponse{
			Success: false,
			Error:   "invalid request body",
		})
		return
	}

	if req.Query == "" {
		writeJSON(w, http.StatusBadRequest, QueryResponse{
			Success: false,
			Error:   "query is required",
		})
		return
	}

	slog.Info("query received", "query", req.Query)

	// Call LLM to generate pipeline
	llmResp, err := h.llmClient.GenerateQuery(r.Context(), req.Query)
	if err != nil {
		slog.Error("LLM generation failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, QueryResponse{
			Success: false,
			Error:   "Failed to generate query: " + err.Error(),
		})
		return
	}

	// Validate pipeline
	if err := query.Validate(llmResp.Pipeline); err != nil {
		slog.Error("pipeline validation failed", "error", err)
		writeJSON(w, http.StatusBadRequest, QueryResponse{
			Success: false,
			Error:   "Generated pipeline is not safe: " + err.Error(),
		})
		return
	}

	// Execute pipeline
	collection := h.db.Collection("candidates")
	results, err := query.Execute(r.Context(), collection, llmResp.Pipeline)
	if err != nil {
		slog.Error("pipeline execution failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, QueryResponse{
			Success: false,
			Error:   "Query execution failed: " + err.Error(),
		})
		return
	}

	elapsed := time.Since(start).Milliseconds()
	slog.Info("query completed", "results", len(results), "time_ms", elapsed)

	writeJSON(w, http.StatusOK, QueryResponse{
		Success:     true,
		UIType:      llmResp.UIType,
		ChartConfig: llmResp.ChartConfig,
		Columns:     llmResp.Columns,
		StatConfig:  llmResp.StatConfig,
		Data:        results,
		Title:       llmResp.Title,
		Summary:     llmResp.Summary,
		GenPipeline: llmResp.RawPipeline,
		Meta: &Meta{
			ExecutionTimeMs: elapsed,
			ResultCount:     len(results),
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
