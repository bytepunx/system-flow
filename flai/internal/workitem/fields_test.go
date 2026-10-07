package workitem_test

import (
	"os"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/usage"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// S-0181: release.FieldsFile lists the front-matter fields this flai reads.
// Publishing a flai release that changed it raises the project's
// flai.minimum, so the list must be the code's: when a field is added or
// removed, change the file with it.
func TestFieldsFileIsTheCode(t *testing.T) {
	data, err := os.ReadFile("front-matter-fields.txt")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, line := range strings.Split(string(data), "\n") {
		if line != "" && !strings.HasPrefix(line, "#") {
			got = append(got, line)
		}
	}
	want := []string{
		"item: " + strings.Join(workitem.KnownFields(workitem.Item{}), " "),
		"thread: " + strings.Join(workitem.KnownFields(threads.Thread{}), " "),
		"issue: " + strings.Join(workitem.KnownFields(issues.Issue{}), " "),
		// a story's agent block (S-0189): an older flai would drop a key of it
		// that it does not know when it rewrote the story
		"item.agent: " + strings.Join(workitem.KnownFields(manifest.Agent{}), " "),
		"item.agent.roles: " + strings.Join(workitem.KnownFields(manifest.Role{}), " "),
		// the planning blocks (S-0199), likewise
		"item.cost_of_delay: " + strings.Join(workitem.KnownFields(workitem.CostOfDelay{}), " "),
		"item.cost_of_delay.inputs: " + strings.Join(workitem.KnownFields(workitem.CostInputs{}), " "),
		"item.forecast: " + strings.Join(workitem.KnownFields(workitem.Forecast{}), " "),
		"item.finalized: " + strings.Join(workitem.KnownFields(workitem.Finalized{}), " "),
		// the usage block, its models, and its strategic entries (S-0225),
		// likewise
		"item.usage: " + strings.Join(workitem.KnownFields(usage.Usage{}), " "),
		"item.usage.models: " + strings.Join(workitem.KnownFields(usage.Model{}), " "),
		"item.usage.strategic: " + strings.Join(workitem.KnownFields(usage.Strategic{}), " "),
		// and the counts of each day of its turns (S-0293)
		"item.usage.turns: " + strings.Join(workitem.KnownFields(usage.TurnDay{}), " "),
		// the types each type-restricted field is valid on (S-0176): an older
		// flai refuses an item that carries one on a type it does not allow
		"item.types: " + workitem.FieldTypes(),
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("%s is not the code's fields; write these lines:\n%s", release.FieldsFile, strings.Join(want, "\n"))
	}
	if !strings.HasSuffix(release.FieldsFile, "internal/workitem/front-matter-fields.txt") {
		t.Errorf("release.FieldsFile is %s", release.FieldsFile)
	}
}
