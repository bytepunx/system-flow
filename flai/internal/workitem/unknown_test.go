package workitem

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// captureWarnings sends the default logger to a buffer for the test.
func captureWarnings(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prev := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	t.Cleanup(func() { slog.SetDefault(prev) })
	return &buf
}

func TestUnknownFieldsKeepTheirLines(t *testing.T) {
	fm := "id: S-0001\nhold: until S-0002\nreviewers:\n- alex\n- sam\ntitle: T\nlater:\n  at: 2026-10-01T00:00:00Z\n"
	got := UnknownFields(fm, Item{})
	want := []Field{
		{Name: "hold", Raw: "hold: until S-0002\n"},
		{Name: "reviewers", Raw: "reviewers:\n- alex\n- sam\n"},
		{Name: "later", Raw: "later:\n  at: 2026-10-01T00:00:00Z\n"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %+v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("field %d: got %+v want %+v", i, got[i], want[i])
		}
	}
}

// S-0181: an item a newer flai wrote, with a field this flai does not know,
// is listed with a warning naming the field and the file, and a write keeps
// the field as it was.
func TestAnItemWithAnUnknownFieldIsListedWarnedAndKept(t *testing.T) {
	logs := captureWarnings(t)
	r := newProject(t)
	e := mustCreate(t, r, Epic, "E", "")
	s := mustCreate(t, r, Story, "S", e.ID)
	data, _ := os.ReadFile(s.Path)
	extra := "hold:\n  - S-0009\n  - S-0010\nreviewer: sam\n"
	doc := strings.Replace(string(data), "\n---\n", "\n"+extra+"---\n", 1)
	if err := os.WriteFile(s.Path, []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}

	items, err := r.List(false)
	if err != nil {
		t.Fatalf("listing refused the item: %v", err)
	}
	var got *Item
	for _, it := range items {
		if it.ID == s.ID {
			got = it
		}
	}
	if got == nil || strings.Join(FieldNames(got.Unknown), ",") != "hold,reviewer" {
		t.Fatalf("unknown fields not kept: %+v", got)
	}
	if !strings.Contains(logs.String(), "level=WARN") || !strings.Contains(logs.String(), "fields=\"hold, reviewer\"") || !strings.Contains(logs.String(), s.Path) {
		t.Errorf("no warning naming the fields and the file:\n%s", logs)
	}
	if _, err := r.List(false); err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(logs.String(), "level=WARN"); n != 1 {
		t.Errorf("the warning is given once per file and fields, got %d", n)
	}

	got.Owner = "sam"
	if err := r.Save(got); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(s.Path)
	if !strings.Contains(string(after), "owner: sam\n") || !strings.Contains(string(after), extra+"---\n") {
		t.Errorf("a write dropped or changed the unknown fields:\n%s", after)
	}
}

// S-0181: a flai below the manifest's minimum stops at Open, before any item
// is read, naming the version needed, rather than reading past fields it
// does not know.
func TestOpenRefusesAFlaiBelowTheMinimum(t *testing.T) {
	logs := captureWarnings(t)
	r := newProject(t)
	e := mustCreate(t, r, Epic, "E", "")
	data, _ := os.ReadFile(e.Path)
	_ = os.WriteFile(e.Path, []byte(strings.Replace(string(data), "\n---\n", "\nhold: x\n---\n", 1)), 0o644)
	man := filepath.Join(r.Root, manifest.File)
	m, _ := os.ReadFile(man)
	_ = os.WriteFile(man, append(m, []byte("flai:\n  minimum: 1.27.0\n")...), 0o644)
	was := buildinfo.Version
	t.Cleanup(func() { buildinfo.Version = was })
	buildinfo.Version = "1.26.4"

	_, err := Open(r.Root)
	var old *manifest.TooOldError
	if !errors.As(err, &old) || !strings.Contains(err.Error(), "needs flai 1.27.0 or newer") {
		t.Fatalf("Open: %v", err)
	}
	if logs.Len() != 0 {
		t.Errorf("items were read before the refusal:\n%s", logs)
	}
}
