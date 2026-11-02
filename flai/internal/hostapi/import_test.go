package hostapi

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
)

// S-0098: a repository offered for import answers the import methods and
// nothing else, runs flai import with the key flai serve chose, and says when
// an import has been applied and committed, its formula's instance made or
// not (S-0424).
func TestImportMethods(t *testing.T) {
	p := channel.Project{Key: "import-widget-2", Name: "widget", Root: "/repos/widget"}
	var imported []string
	var journal []Entry
	host := Host{Record: func(e Entry) { journal = append(journal, e) }}
	rec := &recorder{ran: Ran{Stdout: []byte(`{"key":"widget-2","commit":{"committed":true,"commit":"abc1234"},"formula":{"story":"S-0001","status":"ready","formula":"import"}}`)}}
	m := ImportMethods(rec.run, time.Now, host, func(root string) { imported = append(imported, root) })

	names := []string{}
	for n := range m {
		names = append(names, n)
	}
	slices.Sort(names)
	if strings.Join(names, " ") != "import.preview import.run" {
		t.Fatalf("methods: %v", names)
	}

	// a write: no request ID, no import
	if _, e := m["import.run"](context.Background(), p, json.RawMessage(`{}`)); e == nil || e.Code != channel.CodeInvalidParams {
		t.Errorf("without a request id: %+v", e)
	}
	res, e := m["import.run"](context.Background(), p, json.RawMessage(`{"request_id":"req-00000001"}`))
	if e != nil {
		t.Fatal(e)
	}
	args := rec.runs[len(rec.runs)-1].Args
	want := []string{"import", ".", "--yes", "--formula", "--var", "project_key=widget-2", "--var", "project_name=widget", "--trailer", Trailer, "--json"}
	if !slices.Equal(args, want) || rec.runs[len(rec.runs)-1].Dir != "/repos/widget" {
		t.Errorf("ran %v in %s", args, rec.runs[len(rec.runs)-1].Dir)
	}
	if w, _ := res.(Written); !strings.Contains(string(w.Data), "abc1234") {
		t.Errorf("answer: %+v", res)
	}
	if !slices.Equal(imported, []string{"/repos/widget"}) {
		t.Errorf("imported: %v", imported)
	}
	// the journal names the commit and the instance that runs the tests
	if len(journal) != 1 || journal[0].Action != ActionImport || journal[0].Outcome != "done" ||
		journal[0].Detail != "imported and committed as abc1234; S-0001 checks it and runs its tests" {
		t.Errorf("journal: %+v", journal)
	}

	// the instance not made: the import is committed all the same, so the
	// repository is served from now on, and the error says so with the
	// answer as data
	notMade := "the import is committed as abc1234, but the import formula's instance was not made: no import formula; make it with flai formula run import --input path=/repos/widget --ready --autocommit in /repos/widget"
	rec.ran = Ran{Exit: 1, Stdout: []byte(`{"key":"widget-2","commit":{"committed":true,"commit":"abc1234"},"formula":{"error":"` + notMade + `"}}`),
		Events: []map[string]any{{"level": "FATAL", "err": notMade}}}
	_, e = m["import.run"](context.Background(), p, json.RawMessage(`{"request_id":"req-00000002"}`))
	if e == nil || e.Code != channel.CodeInternal || e.Message != notMade {
		t.Fatalf("no instance: %+v", e)
	}
	if c, _ := e.Data.(map[string]any)["commit"].(map[string]any); c["commit"] != "abc1234" {
		t.Errorf("data: %+v", e.Data)
	}
	if len(imported) != 2 || journal[1].Outcome != "done" || journal[1].Detail != notMade {
		t.Errorf("after no instance: %v %+v", imported, journal[1])
	}

	// a failure of the import itself: nothing to serve
	rec.ran = Ran{Exit: 1, Events: []map[string]any{{"level": "FATAL", "err": "/repos/widget has uncommitted changes"}}}
	if _, e = m["import.run"](context.Background(), p, json.RawMessage(`{"request_id":"req-00000003"}`)); e == nil || !strings.Contains(e.Message, "uncommitted") {
		t.Errorf("refused: %+v", e)
	}
	if len(imported) != 2 || journal[2].Outcome != "failed" {
		t.Errorf("after a refusal: %v %+v", imported, journal[2])
	}

	// the preview changes nothing and needs no request id
	rec.ran = Ran{Stdout: []byte(`{"tests":[]}`)}
	if _, e := m["import.preview"](context.Background(), p, json.RawMessage(`{}`)); e != nil {
		t.Fatal(e)
	}
	if args := rec.runs[len(rec.runs)-1].Args; !slices.Equal(args, []string{"import", ".", "--dry-run", "--json"}) {
		t.Errorf("preview ran %v", args)
	}
}
