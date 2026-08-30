package manusai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// UnmarshalJSON accepts the timestamp formats returned by Manus API v2 while
// preserving Timestamp as a Unix timestamp for SDK users. Numeric timestamps
// are retained as returned; RFC3339 timestamps are normalized to Unix milliseconds.
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

	timestamp, err := parseTaskMessageTimestamp(wire.Timestamp)
	if err != nil {
		return fmt.Errorf("decode task message timestamp: %w", err)
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

func parseTaskMessageTimestamp(raw json.RawMessage) (int64, error) {
	value := strings.TrimSpace(string(raw))
	if value == "" || value == "null" {
		return 0, nil
	}

	if strings.HasPrefix(value, "\"") {
		var text string
		if err := json.Unmarshal(raw, &text); err != nil {
			return 0, fmt.Errorf("invalid string value %s", value)
		}
		if timestamp, err := strconv.ParseInt(text, 10, 64); err == nil {
			return timestamp, nil
		}
		if timestamp, err := time.Parse(time.RFC3339Nano, text); err == nil {
			return timestamp.UnixMilli(), nil
		}
		return 0, fmt.Errorf("%q is not a Unix timestamp or RFC3339 timestamp", text)
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var number json.Number
	if err := decoder.Decode(&number); err != nil {
		return 0, fmt.Errorf("%s is not a JSON number or string", value)
	}
	timestamp, err := strconv.ParseInt(number.String(), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%q is not an integer Unix timestamp", number.String())
	}
	return timestamp, nil
}
