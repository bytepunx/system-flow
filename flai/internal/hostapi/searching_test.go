package hostapi

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
)

func TestSearchQuery(t *testing.T) {
	p := withDocs(t)
	var res SearchResult
	if err := call(t, p, "search.query", `{"q":"berths"}`, &res); err != nil {
		t.Fatal(err)
	}
	if res.Query != "berths" || res.Indexed < 8 || len(res.Hits) == 0 || res.Hits[0].ItemID != "S-0001" || res.Hits[0].Kind != "item" || res.Hits[0].Scope != "wip" {
		t.Errorf("%+v", res)
	}
	for _, h := range res.Hits {
		if strings.HasPrefix(h.Path, "src/") || h.Path == "secret.md" || strings.Contains(h.Path, ".hidden") {
			t.Errorf("indexed outside the project's documents: %s", h.Path)
		}
	}

	// docs are a scope of their own
	_ = call(t, p, "search.query", `{"q":"users"}`, &res)
	for _, h := range res.Hits {
		if h.Scope == "docs" {
			t.Errorf("docs without asking: %+v", h)
		}
	}
	_ = call(t, p, "search.query", `{"q":"users","docs":true}`, &res)
	if len(res.Hits) == 0 || res.Hits[0].Path != "docs/users/index.md" {
		t.Errorf("with docs: %+v", res.Hits)
	}

	// a file written after the index was built is found: the index follows the files
	before := res.Indexed
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(filepath.Join(p.Root, "design/system/lighthouse.md"), []byte("---\ntitle: Lighthouse\n---\n\n# Lighthouse\n\nA lamp on a headland.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = call(t, p, "search.query", `{"q":"headland"}`, &res)
	if res.Indexed != before+1 || len(res.Hits) != 1 || res.Hits[0].Title != "Lighthouse" || !strings.Contains(res.Hits[0].Snippet, "lamp on a headland") {
		t.Errorf("after a new file: indexed %d (was %d) %+v", res.Indexed, before, res.Hits)
	}

	for params, want := range map[string]string{`{"q":"x","limit":1000}`: "limit", `{"q":"` + strings.Repeat("a", 501) + `"}`: "500"} {
		if err := call(t, p, "search.query", params, &res); err == nil || err.Code != channel.CodeInvalidParams || !strings.Contains(err.Message, want) {
			t.Errorf("%.30s: %+v", params, err)
		}
	}
}
