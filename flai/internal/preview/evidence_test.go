package preview

import (
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
)

func TestParseEvidenceReadsTheVerdictAndTheFilesOfEachCriterion(t *testing.T) {
	text := strings.Join([]string{
		"# Evidence",
		"",
		"Verdict: pass, tests and lint clean",
		"",
		"- 1: flai/cmd/accept.go, ./flai/internal/preview/accept.go",
		"- 2: `flai/internal/guard/guard.go`, `docs/users/flai.md` — the refusal, as documented",
		"* 3: `design/adrs/`, with a note, and commas",
		"- 4:",
		"- not a criterion: ignored",
	}, "\r\n")
	ev, err := ParseEvidence(text)
	if err != nil {
		t.Fatal(err)
	}
	if ev.Verdict != "pass, tests and lint clean" {
		t.Errorf("verdict = %q", ev.Verdict)
	}
	if ev.Text != text {
		t.Error("the evidence does not keep its text")
	}
	want := []CriterionEvidence{
		{N: 1, Files: []string{"flai/cmd/accept.go", "flai/internal/preview/accept.go"}},
		{N: 2, Files: []string{"flai/internal/guard/guard.go", "docs/users/flai.md"}},
		{N: 3, Files: []string{"design/adrs"}},
		{N: 4, Files: []string{}},
	}
	if !reflect.DeepEqual(ev.Criteria, want) {
		t.Errorf("criteria = %+v, want %+v", ev.Criteria, want)
	}
}

func TestParseEvidenceRefusesWhatCannotBeChecked(t *testing.T) {
	for name, tc := range map[string]struct{ text, says string }{
		"no verdict":          {"- 1: a.go", "no Verdict line"},
		"empty verdict":       {"Verdict:  \n- 1: a.go", "Verdict line is empty"},
		"two verdicts":        {"Verdict: pass\nVerdict: fail\n- 1: a.go", "two Verdict lines"},
		"no criterion":        {"Verdict: pass\n\nAll good.", "lists no criterion"},
		"criterion twice":     {"Verdict: pass\n- 1: a.go\n- 1: b.go", "criterion 1 twice"},
		"number out of range": {"Verdict: pass\n- 99999999999999999999: a.go", "not a criterion's number"},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := ParseEvidence(tc.text)
			if err == nil || !strings.Contains(err.Error(), tc.says) {
				t.Errorf("error = %v, want one saying %q", err, tc.says)
			}
		})
	}
}

func TestUnevidencedNamesEachCriterionWithoutAChangedFile(t *testing.T) {
	ev, err := ParseEvidence("Verdict: pass\n- 1: `a.go` — meets it\n- 2: `wip/x.md`, `gone.go`")
	if err != nil {
		t.Fatal(err)
	}
	criteria := []itemedit.Criterion{{N: 1, Text: "one"}, {N: 2, Text: "two"}, {N: 3, Text: "three"}}
	got := unevidenced(criteria, ev, []string{"a.go", "b.go"})
	if len(got) != 2 {
		t.Fatalf("blockers = %+v, want two", got)
	}
	for i, says := range []string{"criterion 2 (two) names no file the branch changes (it names wip/x.md, gone.go)", "no item for acceptance criterion 3 (three)"} {
		if got[i].Code != BlockCriterionUnevidenced || !strings.Contains(got[i].Message, says) {
			t.Errorf("blocker %d = %+v, want %s saying %q", i, got[i], BlockCriterionUnevidenced, says)
		}
	}
}
