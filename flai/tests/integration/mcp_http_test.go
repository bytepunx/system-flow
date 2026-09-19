// Package integration holds tests that need something running beside the
// code: here, a flaiover dashboard. See design/conventions/code-quality.md.
package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// bearer adds the dashboard token to every request.
type bearer struct {
	token string
	next  http.RoundTripper
}

func (b bearer) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.Header.Set("Authorization", "Bearer "+b.token)
	r.Header.Set("X-Flai-Agent", "integration-test")
	return b.next.RoundTrip(r)
}

// TestMCPOverHTTP connects to a running dashboard's /mcp with the official
// Go SDK's Streamable HTTP client (ADR-0024, S-0043). It needs a dashboard:
//
//	flai dashboard
//	FLAIOVER_MCP_URL=http://localhost:4242/mcp \
//	FLAIOVER_TOKEN="$(cat .flai-cache/dashboard.token)" \
//	go test ./tests/integration/ -run TestMCPOverHTTP -v
//
// Without those variables it skips, so the ordinary tiers are unaffected.
func TestMCPOverHTTP(t *testing.T) {
	url, token := os.Getenv("FLAIOVER_MCP_URL"), os.Getenv("FLAIOVER_TOKEN")
	if url == "" || token == "" {
		t.Skip("set FLAIOVER_MCP_URL and FLAIOVER_TOKEN to test against a running dashboard")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// without the token the endpoint answers 401, not a login page
	resp, err := http.Post(url, "application/json", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"ping"}`))
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("without a token: status %d, want 401", resp.StatusCode)
	}

	client := mcp.NewClient(&mcp.Implementation{Name: "flai-integration", Version: "0"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{
		Endpoint:   url,
		HTTPClient: &http.Client{Transport: bearer{token: token, next: http.DefaultTransport}},
	}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer func() { _ = session.Close() }()

	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("tools/list: %v", err)
	}
	var names []string
	for _, tool := range tools.Tools {
		names = append(names, tool.Name)
	}
	sort.Strings(names)
	for _, want := range []string{"board", "inbox", "item_get", "wait_for_events"} {
		if i := sort.SearchStrings(names, want); i >= len(names) || names[i] != want {
			t.Errorf("tool %s is missing from %v", want, names)
		}
	}

	for _, name := range []string{"inbox", "board"} {
		res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: name, Arguments: map[string]any{}})
		if err != nil || res.IsError {
			t.Fatalf("%s: %v %+v", name, err, res)
		}
		data, _ := json.Marshal(res.StructuredContent)
		var out map[string]any
		if err := json.Unmarshal(data, &out); err != nil {
			t.Fatalf("%s result: %v", name, err)
		}
		switch name {
		case "inbox":
			if out["agent"] != "integration-test" {
				t.Errorf("the session should run as the agent named in X-Flai-Agent, got %v", out["agent"])
			}
			if _, ok := out["ready"]; !ok {
				t.Errorf("inbox should list ready work: %v", out)
			}
		case "board":
			if _, ok := out["columns"]; !ok {
				t.Errorf("board should have columns: %v", out)
			}
		}
	}
}
