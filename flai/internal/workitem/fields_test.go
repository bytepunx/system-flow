package workitem_test

import (
	"os"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/threads"
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
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("%s is not the code's fields; write these lines:\n%s", release.FieldsFile, strings.Join(want, "\n"))
	}
	if !strings.HasSuffix(release.FieldsFile, "internal/workitem/front-matter-fields.txt") {
		t.Errorf("release.FieldsFile is %s", release.FieldsFile)
	}
}
