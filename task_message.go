package manusai

import (
	"encoding/json"
	"fmt"
)

// UnmarshalJSON accepts the timestamp formats returned by Manus API v2 while
// preserving Timestamp as Unix milliseconds for SDK users.
func (m *TaskMessage) UnmarshalJSON(data []byte) error {
	type taskMessageWire struct {
		ID               string                 `json:"id"`
		Type             string                 `json:"type"`
		Timestamp        json.RawMessage        `json:"timestamp"`
		UserMessage      map[string]interface{} `json:"user_message,omitempty"`
		AssistantMessage map[string]interface{} `json:"assistant_message,omitempty"`
		ErrorMessage     map[string]interface{} `json:"error_message,omitempty"`
		StatusUpdate     map[string]interface{} `json:"status_update,omitempty"`
	}

	var wire taskMessageWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}

	timestamp, err := parseTimestamp(wire.Timestamp, "task message timestamp")
	if err != nil {
		return fmt.Errorf("decode task message: %w", err)
	}

	*m = TaskMessage{
		ID:               wire.ID,
		Type:             wire.Type,
		Timestamp:        timestamp,
		UserMessage:      wire.UserMessage,
		AssistantMessage: wire.AssistantMessage,
		ErrorMessage:     wire.ErrorMessage,
		StatusUpdate:     wire.StatusUpdate,
	}
	return nil
}
