package hostapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
)

// ActionImport is how the journal names an import from the board (S-0098).
// It is not among Actions and is not enabled with flai serve enable: the
// operator consents by naming, with flai serve import add, the folder a
// repository is in.
const ActionImport = "import"

// CandidatePrefix begins a candidate's key (serve.CandidatePrefix, which
// imports this package); what follows it is the key the project is imported with.
const CandidatePrefix = "import-"

// NotCommitted says an import was applied but not committed, because a test
// failed; Data is the import's answer, with the tests and their output.
const NotCommitted = -32013

func importSpecs() map[string]spec {
	return map[string]spec{
		// what an import would do, and which tests it would run, changing nothing
		"import.preview": {reads: true, build: func(channel.Project, json.RawMessage) ([]string, string, *channel.Error) {
			return []string{"import", ".", "--dry-run"}, "", nil
		}},
		// the import, its tests, and its commit; long tests outlive the page
		"import.run": {record: ActionImport, progress: true, detachTimeout: 2 * time.Hour, exits: map[int]int{5: NotCommitted}, describe: describeImport,
			build: func(p channel.Project, _ json.RawMessage) ([]string, string, *channel.Error) {
				// the key flai serve chose, unique among what it serves, not the
				// template's default of the name's initials, which collide easily
				key := strings.TrimPrefix(p.Key, CandidatePrefix)
				return []string{"import", ".", "--yes", "--commit", "--var", "project_key=" + key, "--var", "project_name=" + p.Name, "--trailer", Trailer}, "", nil
			}},
	}
}

// ImportMethods are all that a repository offered for import answers
// (S-0098): there is no project there yet, so nothing else of the table.
// imported is called with the root once an import has been applied,
// committed or not, so that flai serve serves it as a project from then on.
func ImportMethods(run Runner, now func() time.Time, host Host, imported func(root string)) map[string]channel.Method {
	if now == nil {
		now = time.Now
	}
	m := methodsFrom(importSpecs(), run, now, host)
	apply := m["import.run"]
	m["import.run"] = func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
		res, err := apply(ctx, p, raw)
		if imported != nil && (err == nil || err.Code == NotCommitted) {
			imported(p.Root)
		}
		return res, err
	}
	return m
}

func describeImport(res any, err *channel.Error) (outcome, detail string) {
	var said struct {
		Commit struct {
			Committed bool   `json:"committed"`
			Commit    string `json:"commit"`
			Reason    string `json:"reason"`
		} `json:"commit"`
	}
	switch {
	case err != nil && err.Code == NotCommitted:
		return "done", "imported, not committed: " + err.Message
	case err != nil:
		return "failed", err.Message
	}
	w, _ := res.(Written)
	_ = json.Unmarshal(w.Data, &said)
	return "done", fmt.Sprintf("imported and committed as %s", said.Commit.Commit)
}
