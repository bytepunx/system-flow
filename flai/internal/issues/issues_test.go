package issues

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

var t0 = time.Date(2026, 9, 16, 10, 0, 0, 0, time.UTC)

func repo(t *testing.T) *workitem.Repo {
	t.Helper()
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(root, "design"), 0o755)
	r, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestLifecycle(t *testing.T) {
	r := repo(t)
	if id := NextID(r); id != "I-0001" {
		t.Fatalf("first id %s", id)
	}
	is, err := New(r, NewOptions{Title: "Lint: version mismatch", Class: "efficiency", Cost: "5m", Note: "found in S-004", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	if is.ID != "I-0001" || is.Count != 1 || is.Cost != "5m" || !strings.HasSuffix(is.Path, "I-0001-lint-version-mismatch.md") {
		t.Fatalf("new: %+v", is)
	}
	raw, _ := os.ReadFile(is.Path)
	if !strings.Contains(string(raw), `title: "Lint: version mismatch"`) || !strings.Contains(string(raw), "### 2026-09-16T10:00:00Z\nfound in S-004") {
		t.Errorf("file:\n%s", raw)
	}
	back, err := Read(is.Path)
	if err != nil || back.Marshal() != string(raw) {
		t.Fatalf("round trip: %v", err)
	}
	if err := Bump(back, "15m", "again", t0.Add(time.Hour)); err != nil {
		t.Fatal(err)
	}
	if back.Count != 2 || back.Cost != "10m" || back.LastReported != "2026-09-16T11:00:00Z" || back.FirstReported != "2026-09-16T10:00:00Z" {
		t.Errorf("bump: %+v", back)
	}
	if !strings.Contains(back.Body, "### 2026-09-16T11:00:00Z\nagain\n\n## Remediation") {
		t.Errorf("instance not inserted before remediation:\n%s", back.Body)
	}
	if err := Bump(back, "", "no cost this time", t0.Add(2*time.Hour)); err != nil || back.Count != 3 || back.Cost != "10m" {
		t.Errorf("bump without cost: %v %+v", err, back)
	}
	second, _ := New(r, NewOptions{Title: "Two", Class: "defect", Now: t0})
	if second.ID != "I-0002" || second.Cost != "" {
		t.Errorf("second: %+v", second)
	}
	list, err := WriteSummary(r, t0)
	if err != nil || len(list) != 2 {
		t.Fatal(err)
	}
	sum, _ := os.ReadFile(filepath.Join(Dir(r), SummaryFile))
	s := string(sum)
	if !strings.Contains(s, "| [I-0001](I-0001-lint-version-mismatch.md) | efficiency | Lint: version mismatch | 3 | 10m | 30m |") || !strings.Contains(s, "| [I-0002](I-0002-two.md) | defect | Two | 1 | - | - |") {
		t.Errorf("summary:\n%s", s)
	}
	if strings.Index(s, "[I-0001]") > strings.Index(s, "[I-0002]") {
		t.Error("most expensive first")
	}
	if err := Close(back, "fixed by S-010", t0.Add(3*time.Hour)); err != nil || back.Status != "closed" || !strings.Contains(back.Body, "Closed 2026-09-16T13:00:00Z: fixed by S-010") {
		t.Errorf("close: %v %+v", err, back)
	}
	if err := Bump(back, "", "", t0); err == nil {
		t.Error("bump on closed should fail")
	}
	list, _ = WriteSummary(r, t0)
	if table := SummaryTable(list); strings.Contains(table, "I-0001") || !strings.Contains(table, "I-0002") {
		t.Errorf("closed issue in table:\n%s", table)
	}
	if _, err := New(r, NewOptions{Title: "x", Class: "bug", Now: t0}); err == nil {
		t.Error("bad class should fail")
	}
	if _, err := New(r, NewOptions{Title: "x", Class: "defect", Cost: "soon", Now: t0}); err == nil {
		t.Error("bad cost should fail")
	}
}

func TestValidate(t *testing.T) {
	good := &Issue{ID: "I-0001", Title: "t", Class: "defect", Status: "open", Count: 1, Cost: "5m", FirstReported: "2026-09-16T10:00:00Z", LastReported: "2026-09-16T10:00:00Z", Updated: "2026-09-16T10:00:00Z"}
	if err := good.Validate(); err != nil {
		t.Fatal(err)
	}
	bad := *good
	bad.Class, bad.Status, bad.Count, bad.Cost = "meh", "maybe", 0, "lots"
	err := bad.Validate()
	for _, want := range []string{"class", "status", "count", "cost"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("expected %s in %v", want, err)
		}
	}
	if normalise("1h0m0s") != "1h" || normalise("30m0s") != "30m" || normalise("2h30m0s") != "2h30m" {
		t.Error("normalise")
	}
}
