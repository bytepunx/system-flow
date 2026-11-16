package check

import (
	"encoding/json"
	"strings"
)

// ResultFromOutput is the flai check --json result in out, a command's
// combined output, and whether it holds one: an object with a findings key.
func ResultFromOutput(out string) (*Result, bool) {
	raw, ok := JSONObject(out)
	if !ok {
		return nil, false
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(raw, &keys); err != nil {
		return nil, false
	}
	if _, ok := keys["findings"]; !ok {
		return nil, false
	}
	var res Result
	if err := json.Unmarshal(raw, &res); err != nil {
		return nil, false
	}
	return &res, true
}

// JSONObject is the first JSON object in out, a command's combined output,
// that flai's --json printed, indented, so beginning on a line that is "{"
// alone, and whether there is one: build and log lines may come before it
// or after it.
func JSONObject(out string) (json.RawMessage, bool) {
	for offset := 0; offset < len(out); {
		line, _, _ := strings.Cut(out[offset:], "\n")
		if strings.TrimSuffix(line, "\r") == "{" {
			var raw json.RawMessage
			if err := json.NewDecoder(strings.NewReader(out[offset:])).Decode(&raw); err == nil {
				return raw, true
			}
		}
		offset += len(line) + 1
	}
	return nil, false
}
