package manusai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const unixMillisecondsThreshold int64 = 100000000000

// parseTimestamp decodes the timestamp formats used by Manus API v2 and
// normalizes them to Unix milliseconds.
func parseTimestamp(raw json.RawMessage, field string) (int64, error) {
	value := strings.TrimSpace(string(raw))
	if value == "" || value == "null" {
		return 0, nil
	}

	if strings.HasPrefix(value, "\"") {
		var text string
		if err := json.Unmarshal(raw, &text); err != nil {
			return 0, fmt.Errorf("%s: invalid string value", field)
		}
		if timestamp, err := strconv.ParseInt(text, 10, 64); err == nil {
			return normalizeUnixMilliseconds(timestamp), nil
		}
		if timestamp, err := time.Parse(time.RFC3339Nano, text); err == nil {
			return timestamp.UnixMilli(), nil
		}
		return 0, fmt.Errorf("%s: %q is not a Unix timestamp or RFC3339 timestamp", field, text)
	}

	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var number json.Number
	if err := decoder.Decode(&number); err != nil {
		return 0, fmt.Errorf("%s: expected a JSON number or string", field)
	}
	timestamp, err := strconv.ParseInt(number.String(), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %q is not an integer Unix timestamp", field, number.String())
	}
	return normalizeUnixMilliseconds(timestamp), nil
}

func normalizeUnixMilliseconds(timestamp int64) int64 {
	if timestamp > -unixMillisecondsThreshold && timestamp < unixMillisecondsThreshold {
		return timestamp * 1000
	}
	return timestamp
}

func parseInteger(raw json.RawMessage, field string) (int64, error) {
	value := strings.TrimSpace(string(raw))
	if value == "" || value == "null" {
		return 0, nil
	}
	if strings.HasPrefix(value, "\"") {
		var text string
		if err := json.Unmarshal(raw, &text); err != nil {
			return 0, fmt.Errorf("%s: invalid string value", field)
		}
		value = text
	}
	integer, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %q is not an integer", field, value)
	}
	return integer, nil
}
