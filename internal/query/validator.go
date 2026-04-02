package query

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

var forbiddenStages = []string{"$out", "$merge", "$delete", "$update"}
const maxStages = 8

// Validate checks a MongoDB aggregation pipeline for safety.
func Validate(pipeline []bson.D) error {
	if len(pipeline) == 0 {
		return fmt.Errorf("empty pipeline")
	}
	if len(pipeline) > maxStages {
		return fmt.Errorf("pipeline has %d stages, max allowed is %d", len(pipeline), maxStages)
	}

	for i, stage := range pipeline {
		if len(stage) == 0 {
			return fmt.Errorf("stage %d is empty", i)
		}
		stageName := stage[0].Key
		for _, forbidden := range forbiddenStages {
			if strings.EqualFold(stageName, forbidden) {
				return fmt.Errorf("forbidden stage %q at position %d", stageName, i)
			}
		}
		// Block $lookup to external collections
		if stageName == "$lookup" {
			return fmt.Errorf("$lookup is not allowed")
		}
	}
	return nil
}
