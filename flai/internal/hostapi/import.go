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
// failed; Data is the import's answer, with the tests and their output. It is
// flai import --commit's exit 5, which the dashboard still maps; import.run
// runs flai import --formula, which commits before any test and never answers
// it (S-0424).
const NotCommitted = -32013

func importSpecs() map[string]spec {
	return map[string]spec{
		// what an import would do, and which tests it would run, changing nothing
		"import.preview": {reads: true, build: func(channel.Project, json.RawMessage) ([]string, string, *channel.Error) {
			return []string{"import", ".", "--dry-run"}, "", nil
		}},
		// the import, its commit, and the import formula's instance, which
		// checks it and runs its tests as a story (S-0424); exit 1 with an
		// answer is an import committed whose instance was not made
		"import.run": {record: ActionImport, progress: true, detachTimeout: 2 * time.Hour, answers: []int{1}, judge: importInstance, describe: describeImport,
			build: func(p channel.Project, _ json.RawMessage) ([]string, string, *channel.Error) {
				// the key flai serve chose, unique among what it serves, not the
				// template's default of the name's initials, which collide easily
				key := strings.TrimPrefix(p.Key, CandidatePrefix)
				return []string{"import", ".", "--yes", "--formula", "--var", "project_key=" + key, "--var", "project_name=" + p.Name, "--trailer", Trailer}, "", nil
			}},
	}
}

// ImportMethods are all that a repository offered for import answers
// (S-0098): there is no project there yet, so nothing else of the table.
// imported is called with the root once an import has been applied and
// committed, its formula's instance made or not, so that flai serve serves
// it as a project from then on.
func ImportMethods(run Runner, now func() time.Time, host Host, imported func(root string)) map[string]channel.Method {
	if now == nil {
		now = time.Now
	}
	m := methodsFrom(importSpecs(), run, now, host)
	apply := m["import.run"]
	m["import.run"] = func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
		res, err := apply(ctx, p, raw)
		if imported != nil && (err == nil || instanceNotMade(err)) {
			imported(p.Root)
		}
		return res, err
	}
	return m
}

// importInstance reads flai import --formula's answer and refuses one whose
// formula carries an error: the import is applied and committed, but the
// import formula's instance was not made, its exit 1, which answers counts as
// an answer so that the commit reaches the dashboard. It is an error in flai's
// words, saying how to make the instance, with the answer as data.
func importInstance(w Written) *channel.Error {
	var said struct {
		Formula struct {
			Error string `json:"error"`
		} `json:"formula"`
	}
	_ = json.Unmarshal(w.Data, &said) // an answer that is not the import's made its instance
	if said.Formula.Error == "" {
		return nil
	}
	var data map[string]any
	_ = json.Unmarshal(w.Data, &data)
	return &channel.Error{Code: channel.CodeInternal, Message: said.Formula.Error, Data: data}
}

// instanceNotMade reports whether err is importInstance's refusal: an import
// applied and committed whose formula's instance was not made.
func instanceNotMade(err *channel.Error) bool {
	data, _ := err.Data.(map[string]any)
	formula, _ := data["formula"].(map[string]any)
	_, said := formula["error"]
	return said
}

// describeImport says for the journal what an import did: its commit, and the
// import formula's instance that checks it and runs its tests.
func describeImport(res any, err *channel.Error) (outcome, detail string) {
	var said struct {
		Commit struct {
			Commit string `json:"commit"`
		} `json:"commit"`
		Formula struct {
			Story string `json:"story"`
		} `json:"formula"`
	}
	switch {
	case err != nil && instanceNotMade(err):
		return "done", err.Message
	case err != nil:
		return "failed", err.Message
	}
	w, _ := res.(Written)
	_ = json.Unmarshal(w.Data, &said)
	return "done", fmt.Sprintf("imported and committed as %s; %s checks it and runs its tests", said.Commit.Commit, said.Formula.Story)
}
