package llm

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// LLMResponse represents the parsed response from the LLM.
type LLMResponse struct {
	Pipeline    []bson.D     `json:"-"`
	RawPipeline any          `json:"pipeline"`
	UIType      string       `json:"ui_type"`
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	ChartConfig *ChartConfig `json:"chart_config,omitempty"`
	Columns     []Column     `json:"columns,omitempty"`
	StatConfig  *StatConfig  `json:"stat_config,omitempty"`
}

// ChartConfig holds chart rendering configuration.
type ChartConfig struct {
	XField string `json:"x_field"`
	YField string `json:"y_field"`
	XLabel string `json:"x_label"`
	YLabel string `json:"y_label"`
}

// Column defines a table column.
type Column struct {
	Field  string `json:"field"`
	Header string `json:"header"`
	Format string `json:"format,omitempty"`
}

// StatConfig holds stat card configuration.
type StatConfig struct {
	ValueField string `json:"value_field"`
	Format     string `json:"format"`
	Label      string `json:"label"`
}

// ParseResponse parses the raw LLM text response into an LLMResponse struct.
func ParseResponse(raw string) (*LLMResponse, error) {
	// Strip markdown fences if present
	cleaned := strings.TrimSpace(raw)
	if strings.HasPrefix(cleaned, "```json") {
		cleaned = strings.TrimPrefix(cleaned, "```json")
		cleaned = strings.TrimSuffix(cleaned, "```")
		cleaned = strings.TrimSpace(cleaned)
	} else if strings.HasPrefix(cleaned, "```") {
		cleaned = strings.TrimPrefix(cleaned, "```")
		cleaned = strings.TrimSuffix(cleaned, "```")
		cleaned = strings.TrimSpace(cleaned)
	}

	// Parse the outer JSON structure
	var rawResp struct {
		Pipeline    json.RawMessage `json:"pipeline"`
		UIType      string          `json:"ui_type"`
		Title       string          `json:"title"`
		Summary     string          `json:"summary"`
		ChartConfig *ChartConfig    `json:"chart_config,omitempty"`
		Columns     []Column        `json:"columns,omitempty"`
		StatConfig  *StatConfig     `json:"stat_config,omitempty"`
	}

	if err := json.Unmarshal([]byte(cleaned), &rawResp); err != nil {
		return nil, fmt.Errorf("JSON parse error: %w", err)
	}

	// Parse pipeline stages from JSON to bson.D
	pipeline, err := parsePipeline(rawResp.Pipeline)
	if err != nil {
		return nil, fmt.Errorf("pipeline parse error: %w", err)
	}

	// Keep raw pipeline for API response
	var rawPipeline any
	json.Unmarshal(rawResp.Pipeline, &rawPipeline)

	return &LLMResponse{
		Pipeline:    pipeline,
		RawPipeline: rawPipeline,
		UIType:      rawResp.UIType,
		Title:       rawResp.Title,
		Summary:     rawResp.Summary,
		ChartConfig: rawResp.ChartConfig,
		Columns:     rawResp.Columns,
		StatConfig:  rawResp.StatConfig,
	}, nil
}

// parsePipeline converts JSON pipeline stages into []bson.D.
func parsePipeline(raw json.RawMessage) ([]bson.D, error) {
	// First unmarshal as a generic structure to get the JSON
	var stages []map[string]any
	if err := json.Unmarshal(raw, &stages); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pipeline: %w", err)
	}

	// Convert each stage through JSON -> BSON extended JSON
	var pipeline []bson.D
	for i, stage := range stages {
		stageBytes, err := json.Marshal(stage)
		if err != nil {
			return nil, fmt.Errorf("stage %d marshal error: %w", i, err)
		}

		var bsonStage bson.D
		if err := bson.UnmarshalExtJSON(stageBytes, false, &bsonStage); err != nil {
			return nil, fmt.Errorf("stage %d BSON parse error: %w", i, err)
		}
		pipeline = append(pipeline, bsonStage)
	}

	return pipeline, nil
}
