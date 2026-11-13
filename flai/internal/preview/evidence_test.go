package preview

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
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
	got := unevidenced(criteria, ev, []string{"a.go", "b.go"}, nil)
	if len(got) != 2 {
		t.Fatalf("blockers = %+v, want two", got)
	}
	for i, says := range []string{"criterion 2 (two) names no file the branch changes (it names wip/x.md, gone.go)", "no item for acceptance criterion 3 (three)"} {
		if got[i].Code != BlockCriterionUnevidenced || !strings.Contains(got[i].Message, says) {
			t.Errorf("blocker %d = %+v, want %s saying %q", i, got[i], BlockCriterionUnevidenced, says)
		}
	}
}

// S-0662: an item may name the story's kluster record in place of files,
// "kluster deploy and checks passed at <commit>", or the deploy alone, with
// a note after it; anything else after "kluster" is read as files.
func TestParseEvidenceReadsAnItemNamingTheKlusterRecord(t *testing.T) {
	ev, err := ParseEvidence(strings.Join([]string{
		"Verdict: pass",
		"- 1: kluster deploy and checks passed at `" + kevC3[:12] + "` — the API answers on its namespace",
		"- 2: Kluster deploy passed at " + strings.ToUpper(kevC2),
		"- 3: kluster deploy passed at the head",
	}, "\n"))
	if err != nil {
		t.Fatal(err)
	}
	want := []CriterionEvidence{
		{N: 1, Files: []string{}, Kluster: kevC3[:12]},
		{N: 2, Files: []string{}, Kluster: strings.ToUpper(kevC2)},
		{N: 3, Files: []string{"kluster deploy passed at the head"}},
	}
	if !reflect.DeepEqual(ev.Criteria, want) {
		t.Errorf("criteria = %+v, want %+v", ev.Criteria, want)
	}
}

// S-0662, ADR-0292 §3, §6: an item naming the kluster record is evidence
// only while the kluster host action is on and the manifest names a deploy,
// the record's deploy and, while the checks host action is on, its checks
// run passed at the story branch's head (or at a commit followed only by
// changes under design, docs, and wip), that head is the commit verified,
// and the item names the record's commit or one at the head the same way.
// Otherwise the blocker names the criterion and what the record says.
func TestUnevidencedTakesTheKlusterRecordOnlyAtTheVerifiedHead(t *testing.T) {
	passed, failed := kluster.OutcomePassed, kluster.OutcomeFailed
	for name, c := range map[string]struct {
		deploy, checks    string // the record's commits; "" for no record
		outcome, checksIs string
		klusterOn         bool
		checksOn          bool
		noDeploy          bool
		verified, named   string
		says              string // "" when the item is evidence
	}{
		"at the verified head":                {deploy: kevC3, checks: kevC3, outcome: passed, checksIs: passed, klusterOn: true, checksOn: true, verified: kevC3, named: kevC3[:7]},
		"followed only by docs":               {deploy: kevC2, checks: kevC2, outcome: passed, checksIs: passed, klusterOn: true, checksOn: true, verified: kevC3, named: kevC2[:12]},
		"naming the head, at a docs-only gap": {deploy: kevC2, checks: kevC2, outcome: passed, checksIs: passed, klusterOn: true, checksOn: true, verified: kevC3, named: kevC3},
		"deploy alone with checks off":        {deploy: kevC3, outcome: passed, klusterOn: true, verified: kevC3, named: kevC3},
		"judged at the head for a landing":    {deploy: kevC3, checks: kevC3, outcome: passed, checksIs: passed, klusterOn: true, checksOn: true, named: kevC3},
		"naming another commit": {deploy: kevC3, checks: kevC3, outcome: passed, checksIs: passed, klusterOn: true, checksOn: true, verified: kevC3, named: kevC1[:12],
			says: kevC1[:12] + " is not the commit of the record's deploy, " + kevC3[:12] + ", nor at the head " + kevC3[:12] + ": 2 paths outside design, docs, and wip changed"},
		"naming no commit": {deploy: kevC3, checks: kevC3, outcome: passed, checksIs: passed, klusterOn: true, checksOn: true, verified: kevC3, named: "abcdef1",
			says: "abcdef1 names no commit in this repository"},
		"verified another commit": {deploy: kevC3, checks: kevC3, outcome: passed, checksIs: passed, klusterOn: true, checksOn: true, verified: kevC2, named: kevC3,
			says: "the record is judged at the head of story/S-0001, " + kevC3[:12] + ", which is not the commit the acceptance verified, " + kevC2[:12]},
		"after a failed checks run": {deploy: kevC3, checks: kevC3, outcome: passed, checksIs: failed, klusterOn: true, checksOn: true, verified: kevC3, named: kevC3,
			says: "the record says S-0001's namespace s-0001 is deployed at the head " + kevC3[:12] + ", and its checks run there ended failed"},
		"deployed before a code change": {deploy: kevC1, checks: kevC1, outcome: passed, checksIs: passed, klusterOn: true, checksOn: true, verified: kevC3, named: kevC1,
			says: "the record says S-0001's namespace s-0001 is not deployed at the head " + kevC3[:12]},
		"never deployed": {klusterOn: true, checksOn: true, verified: kevC3, named: kevC3,
			says: "no deploy of it is recorded"},
		"with the action off": {deploy: kevC3, checks: kevC3, outcome: passed, checksIs: passed, checksOn: true, verified: kevC3, named: kevC3,
			says: "no kluster evidence applies to S-0001: the kluster host action is off for this project, so name the changed files that meet it"},
		"with no deploy": {deploy: kevC3, checks: kevC3, outcome: passed, checksIs: passed, klusterOn: true, checksOn: true, noDeploy: true, verified: kevC3, named: kevC3,
			says: "the manifest's kluster: block names no deploy"},
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			if c.deploy != "" {
				kevRecord(t, root, kluster.Record{Commit: c.deploy, Outcome: c.outcome, ChecksCommit: c.checks, ChecksOutcome: c.checksIs})
			}
			m := manifest.Manifest{Layout: map[string]string{"design": "design", "docs": "docs", "wip": "wip"}, Kluster: &manifest.Kluster{Deploy: []string{"scripts/deploy.sh"}}}
			if c.noDeploy {
				m.Kluster.Deploy = nil
			}
			ev, err := ParseEvidence("Verdict: pass\n- 1: `a.go`\n- 2: kluster deploy and checks passed at " + c.named)
			if err != nil {
				t.Fatal(err)
			}
			kl := klusterJudgeOf(kevGit{}, root, m, "S-0001", c.verified, c.klusterOn, c.checksOn)
			got := unevidenced([]itemedit.Criterion{{N: 1, Text: "one"}, {N: 2, Text: "it runs"}}, ev, []string{"a.go"}, kl)
			if c.says == "" {
				if len(got) != 0 {
					t.Errorf("blockers = %+v, want none", got)
				}
				return
			}
			want := "the evidence for acceptance criterion 2 (it runs) names the kluster record at " + c.named + ", but "
			if len(got) != 1 || got[0].Code != BlockCriterionUnevidenced || !strings.HasPrefix(got[0].Message, want) || !strings.Contains(got[0].Message, c.says) {
				t.Errorf("blockers = %+v, want one %s starting %q and saying %q", got, BlockCriterionUnevidenced, want, c.says)
			}
		})
	}
}

// The commits of kevGit's history, oldest first: c1 and c2 change code, c3,
// the head of story/S-0001, only the docs and wip folders.
var (
	kevC1 = strings.Repeat("a1", 20)
	kevC2 = strings.Repeat("b2", 20)
	kevC3 = strings.Repeat("c3", 20)
)

// kevGit is an execx.Runner that answers git from kevC1..kevC3: rev-parse
// resolves story/S-0001 to kevC3 and a commit from a prefix of it,
// merge-base --is-ancestor walks the history, and diff --name-only lists the
// files each commit after the first up to the second changed.
type kevGit struct{}

var (
	kevParent = map[string]string{kevC2: kevC1, kevC3: kevC2}
	kevFiles  = map[string][]string{kevC1: {"cli/a.go"}, kevC2: {"cli/a.go", "cli/b.go"}, kevC3: {"docs/users/guide.md", "wip/agents/S-0001-x.md"}}
)

func (kevGit) Run(_, _ string, args ...string) (string, error) {
	switch args[0] {
	case "rev-parse":
		ref := strings.TrimSuffix(args[len(args)-1], "^{commit}")
		if ref == "refs/heads/story/S-0001" {
			return kevC3 + "\n", nil
		}
		for _, c := range []string{kevC1, kevC2, kevC3} {
			if len(ref) >= 7 && strings.HasPrefix(c, strings.ToLower(ref)) {
				return c + "\n", nil
			}
		}
		return "", errors.New("exit status 1")
	case "merge-base":
		for c := args[3]; c != ""; c = kevParent[c] {
			if c == args[2] {
				return "", nil
			}
		}
		return "", errors.New("exit status 1")
	case "diff":
		from, to := args[len(args)-2], args[len(args)-1]
		var out []string
		for c := to; c != "" && c != from; c = kevParent[c] {
			out = append(out, kevFiles[c]...)
		}
		return strings.Join(out, "\x00"), nil
	}
	return "", errors.New("unexpected git " + strings.Join(args, " "))
}

func (g kevGit) RunInput(dir, name, _ string, args ...string) (string, error) {
	return g.Run(dir, name, args...)
}

func (kevGit) LookPath(name string) (string, error) { return name, nil }

// kevRecord writes r as S-0001's namespace's deployment record under root.
func kevRecord(t *testing.T, root string, r kluster.Record) {
	t.Helper()
	r.Story, r.Namespace = "S-0001", "s-0001"
	path := kluster.RecordPath(root, "s-0001")
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}
