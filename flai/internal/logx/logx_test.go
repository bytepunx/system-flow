package logx

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
)

func TestLevelsAndFormats(t *testing.T) {
	var buf bytes.Buffer
	l := New(&buf, Options{IsTerminal: false})
	l.Debug("hidden")
	l.Info("config created", "component", "config", "path", "/x")
	line := strings.TrimSpace(buf.String())
	var ev map[string]any
	if err := json.Unmarshal([]byte(line), &ev); err != nil {
		t.Fatalf("not json by default off-terminal: %s", line)
	}
	for _, k := range []string{"ts", "level", "msg", "component", "path"} {
		if _, ok := ev[k]; !ok {
			t.Errorf("missing %s in %v", k, ev)
		}
	}
	if ev["level"] != "INFO" || !strings.HasSuffix(ev["ts"].(string), "Z") {
		t.Errorf("level/ts: %v", ev)
	}
	buf.Reset()
	l = New(&buf, Options{IsTerminal: true, Verbose: true})
	l.Debug("shown", "component", "x")
	if !strings.Contains(buf.String(), "level=DEBUG") || !strings.Contains(buf.String(), `msg=shown`) || !strings.Contains(buf.String(), "ts=") {
		t.Errorf("text verbose: %s", buf.String())
	}
	buf.Reset()
	l = New(&buf, Options{Level: "error", Format: "json"})
	l.Warn("nope")
	Fatal(l, "command failed", "err", "boom")
	if strings.Contains(buf.String(), "nope") || !strings.Contains(buf.String(), `"level":"FATAL"`) || !strings.Contains(buf.String(), `"err":"boom"`) {
		t.Errorf("error level + fatal: %s", buf.String())
	}
	buf.Reset()
	_ = New(&buf, Options{Level: "loud", Format: "xml", IsTerminal: true})
	out := buf.String()
	if !strings.Contains(out, "unknown LOG_LEVEL") || !strings.Contains(out, "unknown LOG_FORMAT") {
		t.Errorf("bad values should warn: %s", out)
	}
	if String(LevelFatal) != "FATAL" || String(slog.LevelWarn) != "WARN" {
		t.Error("level names")
	}
}
