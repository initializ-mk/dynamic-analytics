package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Client wraps the Anthropic API client.
type Client struct {
	client anthropic.Client
	model  string
}

// NewClient creates a new LLM client. It reads ANTHROPIC_API_KEY from env by default.
func NewClient(apiKey string) *Client {
	opts := []option.RequestOption{}
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}
	return &Client{
		client: anthropic.NewClient(opts...),
		model:  "claude-sonnet-4-20250514",
	}
}

// GenerateQuery sends a user query to Claude and returns the parsed LLM response.
func (c *Client) GenerateQuery(ctx context.Context, userQuery string) (*LLMResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	systemPrompt := BuildSystemPrompt()

	slog.Info("sending query to LLM", "query", userQuery, "model", c.model)

	response, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:       anthropic.Model(c.model),
		MaxTokens:   2048,
		Temperature: anthropic.Float(0),
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userQuery)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("LLM request failed: %w", err)
	}

	if len(response.Content) == 0 {
		return nil, fmt.Errorf("empty response from LLM")
	}

	text := response.Content[0].Text
	slog.Info("received LLM response", "length", len(text))

	parsed, err := ParseResponse(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w (raw: %s)", err, truncate(text, 500))
	}

	// Log the generated pipeline
	pipelineJSON, _ := json.Marshal(parsed.Pipeline)
	slog.Info("generated pipeline", "pipeline", string(pipelineJSON), "ui_type", parsed.UIType)

	return parsed, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
