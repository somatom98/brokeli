package event_store

import (
	"encoding/json"
	"fmt"
)

// DecodeEvent turns a stored event payload back into its typed value.
func DecodeEvent[T any](content any) (T, error) {
	var zero T

	if event, ok := content.(T); ok {
		return event, nil
	}

	data, err := json.Marshal(content)
	if err != nil {
		return zero, fmt.Errorf("failed to marshal event content: %w", err)
	}

	var decoded T
	if err := json.Unmarshal(data, &decoded); err != nil {
		return zero, fmt.Errorf("failed to unmarshal event content: %w", err)
	}

	return decoded, nil
}
