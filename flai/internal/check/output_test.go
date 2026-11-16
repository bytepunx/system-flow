package check

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// ResultFromOutput finds flai check --json's indented result among build
// and log lines, and holds no result in output without one.
func TestResultFromOutputFindsTheResultAmongOtherLines(t *testing.T) {
	note := Finding{Level: Warning, Rule: "doc.updated", Path: "design/system/overview.md", Line: 1, Message: "updated is behind the file's last change", Outside: true}
	data, err := json.MarshalIndent(Result{Findings: []Finding{note}, Items: 12}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	result := string(data) + "\n"
	for _, tc := range []struct {
		name, out string
		want      bool
	}{
		{"alone", strings.TrimSpace(result), true},
		{"after a build line and a log line", "building bin/flai\nlevel=INFO msg=\"{\"\n" + result, true},
		{"before a log line", result + "level=WARN msg=done", true},
		{"with CRLF line ends", strings.ReplaceAll(result, "\n", "\r\n"), true},
		{"no findings key", "{\n  \"version\": \"1.70.0\"\n}", false},
		{"cut short", "building bin/flai\n{\n  \"findings\": [", false},
		{"no object", "go: build failed", false},
		{"empty", "", false},
	} {
		res, ok := ResultFromOutput(tc.out)
		if ok != tc.want {
			t.Errorf("%s: ok %v, want %v", tc.name, ok, tc.want)
			continue
		}
		if ok && (len(res.Findings) != 1 || !reflect.DeepEqual(res.Findings[0], note) || res.Items != 12) {
			t.Errorf("%s: %+v", tc.name, res)
		}
	}
}

// JSONObject takes the first object that begins on a line of "{" alone,
// passing over braces inside other lines, and finds none without one.
func TestJSONObjectTakesTheFirstObjectOnALineOfItsOwn(t *testing.T) {
	for _, tc := range []struct {
		name, out, want string
		ok              bool
	}{
		{"version", "building bin/flai\n{\n  \"version\": \"1.2.3\"\n}\n", `{"version":"1.2.3"}`, true},
		{"a brace inside a log line", "msg={\"a\":1}\n{\n  \"b\": 2\n}\n{\n  \"c\": 3\n}\n", `{"b":2}`, true},
		{"not JSON after the brace", "{\nnot json\n", "", false},
		{"none", "flai 1.2.3\n", "", false},
	} {
		raw, ok := JSONObject(tc.out)
		if ok != tc.ok {
			t.Errorf("%s: ok %v, want %v", tc.name, ok, tc.ok)
			continue
		}
		if !ok {
			continue
		}
		var got, want any
		if err := json.Unmarshal(raw, &got); err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if err := json.Unmarshal([]byte(tc.want), &want); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: %s, want %s", tc.name, raw, tc.want)
		}
	}
}
