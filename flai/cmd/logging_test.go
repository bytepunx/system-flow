package cmd

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

// events parses JSON log lines from stderr.
func events(t *testing.T, errOut string) []map[string]any {
	t.Helper()
	var out []map[string]any
	for _, line := range strings.Split(strings.TrimSpace(errOut), "\n") {
		if line == "" {
			continue
		}
		var ev map[string]any
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			t.Fatalf("stderr line is not a JSON event: %q", line)
		}
		out = append(out, ev)
	}
	return out
}

func find(evs []map[string]any, msg string) map[string]any {
	for _, e := range evs {
		if e["msg"] == msg {
			return e
		}
	}
	return nil
}

func TestLogEventsPinFieldNames(t *testing.T) {
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("LOG_FORMAT", "")
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))

	// first run creates the config: info event with component and path
	out, errOut, code := runCLI(t, "config", "get")
	if code != 0 || !strings.HasPrefix(out, "{") {
		t.Fatalf("stdout must stay the command's output: %d %q", code, out)
	}
	ev := find(events(t, errOut), "config created with defaults")
	if ev == nil || ev["level"] != "INFO" || ev["component"] != "config" || ev["path"] == nil || ev["ts"] == nil {
		t.Errorf("config created event: %v", ev)
	}

	// failure: exactly one fatal event at the boundary with err
	_, errOut, code = runCLI(t, "config", "get", "nope")
	if code != 1 {
		t.Fatalf("exit %d", code)
	}
	evs := events(t, errOut)
	fatal := find(evs, "command failed")
	if fatal == nil || fatal["level"] != "FATAL" || fatal["component"] != "cmd" || !strings.Contains(fatal["err"].(string), "unknown key") {
		t.Errorf("fatal event: %v", fatal)
	}
	n := 0
	for _, e := range evs {
		if e["level"] == "FATAL" || e["level"] == "ERROR" {
			n++
		}
	}
	if n != 1 {
		t.Errorf("errors must be logged once, got %d", n)
	}

	// --verbose shows debug events; LOG_FORMAT=text switches the handler
	_, errOut, _ = runCLI(t, "config", "get", "--verbose")
	if ev := find(events(t, errOut), "config loaded"); ev == nil || ev["level"] != "DEBUG" {
		t.Errorf("verbose debug event: %s", errOut)
	}
	_, errOut, _ = runCLI(t, "config", "get")
	if strings.Contains(errOut, "config loaded") {
		t.Error("debug event shown without --verbose")
	}
	t.Setenv("LOG_FORMAT", "text")
	_, errOut, _ = runCLI(t, "config", "get", "--verbose")
	if !strings.Contains(errOut, "level=DEBUG") || !strings.Contains(errOut, "msg=\"config loaded\"") || !strings.Contains(errOut, "ts=") {
		t.Errorf("text format: %s", errOut)
	}
	t.Setenv("LOG_FORMAT", "")
	t.Setenv("LOG_LEVEL", "error")
	_, errOut, _ = runCLI(t, "config", "get", "--verbose")
	if !strings.Contains(errOut, "config loaded") {
		t.Error("--verbose must win over LOG_LEVEL")
	}
}
