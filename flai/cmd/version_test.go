package cmd

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

// flai version --json prints the task-done rules number beside the version,
// which an outdated MCP server reads from the installed flai (ADR-0180).
func TestVersionJSONPrintsTheTaskDoneRules(t *testing.T) {
	var out bytes.Buffer
	a := &app{out: &out, jsonOut: true}
	if err := a.printVersion(versionReport{Info: buildinfo.Get(), Release: buildinfo.OwnStamp()}); err != nil {
		t.Fatal(err)
	}
	var v struct {
		Version       string `json:"version"`
		TaskDoneRules *int   `json:"task_done_rules"`
	}
	if err := json.Unmarshal(out.Bytes(), &v); err != nil {
		t.Fatalf("not json: %v\n%s", err, out.String())
	}
	if v.TaskDoneRules == nil || *v.TaskDoneRules != buildinfo.TaskDoneRules || v.Version == "" {
		t.Errorf("flai version --json = %s, want task_done_rules %d", out.String(), buildinfo.TaskDoneRules)
	}
}
