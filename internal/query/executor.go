package query

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Execute runs a MongoDB aggregation pipeline and returns the results.
func Execute(ctx context.Context, collection *mongo.Collection, pipeline []bson.D) ([]map[string]any, error) {
	execCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	mongoPipeline := make(mongo.Pipeline, len(pipeline))
	for i, stage := range pipeline {
		mongoPipeline[i] = stage
	}

	cursor, err := collection.Aggregate(execCtx, mongoPipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(execCtx)

	var results []bson.M
	if err := cursor.All(execCtx, &results); err != nil {
		return nil, err
	}

	// Convert bson.M to map[string]any
	out := make([]map[string]any, len(results))
	for i, r := range results {
		out[i] = map[string]any(r)
	}
	return out, nil
}
