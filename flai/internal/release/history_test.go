package release

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// pendingByProcess is Pending as it was before S-0157, one git process per
// question: the reference the kept history must answer exactly as.
func pendingByProcess(r execx.Runner, root string, m manifest.Manifest, repo *workitem.Repo) ([]*PendingPlan, error) {
	plans := map[string]*Plan{}
	itemPlan := func(id string) (*Plan, error) {
		if p, ok := plans[id]; ok {
			return p, nil
		}
		it, err := repo.Get(id)
		if err != nil {
			return nil, err
		}
		var parent *workitem.Item
		if it.Parent != "" {
			parent, _ = repo.Get(it.Parent)
		}
		p, err := Compute(r, root, m, it, parent, "")
		if err != nil {
			return nil, err
		}
		plans[id] = p
		return p, nil
	}
	var out []*PendingPlan
	for _, proj := range m.Projects {
		cur, err := CurrentVersion(r, root, proj)
		if err != nil {
			return nil, err
		}
		boundary, err := componentBoundary(r, root, proj, cur)
		if err != nil {
			return nil, err
		}
		revRange := "HEAD"
		if boundary != "" {
			revRange = boundary + "..HEAD"
		}
		ids, err := acceptedIDsIn(r, root, revRange)
		if err != nil {
			return nil, err
		}
		pp := &PendingPlan{Component: proj, From: cur}
		seenFile := map[string]bool{}
		for _, id := range ids {
			ip, err := itemPlan(id)
			if err != nil {
				continue
			}
			for _, st := range ip.Steps {
				if st.Component.Name != proj.Name {
					continue
				}
				pp.Level = higher(pp.Level, st.Level)
				pp.Items = append(pp.Items, PendingItem{ID: id, Title: ip.Title, Level: st.Level})
				for _, f := range st.Files {
					if !seenFile[f] {
						seenFile[f] = true
						pp.Files = append(pp.Files, f)
					}
				}
			}
		}
		if len(pp.Items) == 0 || pp.Level == "" || pp.Level == None {
			continue
		}
		sort.Strings(pp.Files)
		pp.To = cur.Bumped(pp.Level)
		if proj.Kind == "template" {
			pp.Version = filepath.ToSlash(filepath.Join(proj.Path, "template.yaml"))
		} else {
			pp.Tag = proj.Name + "/v" + pp.To.String()
		}
		out = append(out, pp)
	}
	return out, nil
}

// acceptedIDsIn lists the items accepted by a commit in revRange, each once,
// oldest first: the reference's own git log.
func acceptedIDsIn(r execx.Runner, root, revRange string) ([]string, error) {
	out, err := r.Run(root, "git", "log", "--reverse", "--format=%s", revRange)
	if err != nil {
		return nil, err
	}
	var ids []string
	seen := map[string]bool{}
	for _, l := range strings.Split(out, "\n") {
		if m := acceptedSubject.FindStringSubmatch(l); m != nil && !seen[m[1]] {
			seen[m[1]] = true
			ids = append(ids, m[1])
		}
	}
	return ids, nil
}

// componentBoundary is the ref marking a component's last publish, or ""
// when it has never been released. A code component's tag is that ref. A
// template has no local tag; the commit that set its current version stands
// in.
func componentBoundary(r execx.Runner, root string, p manifest.Project, cur Version) (string, error) {
	if cur == (Version{}) {
		return "", nil
	}
	if p.Kind != "template" {
		return p.Name + "/v" + cur.String(), nil
	}
	out, err := r.Run(root, "git", "log", "-n1", "--format=%H", "-S", "version: "+cur.String(), "--", filepath.Join(p.Path, "template.yaml"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// counting records the git commands a runner starts.
type counting struct {
	execx.Runner
	mu  sync.Mutex
	ran []string
}

func (c *counting) Run(dir, name string, args ...string) (string, error) {
	return c.RunInput(dir, name, "", args...)
}

func (c *counting) RunInput(dir, name, input string, args ...string) (string, error) {
	c.mu.Lock()
	c.ran = append(c.ran, name+" "+args[0])
	c.mu.Unlock()
	return c.Runner.RunInput(dir, name, input, args...)
}

func (c *counting) take() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	ran := c.ran
	c.ran = nil
	return ran
}

// S-0157: Pending from the kept history answers exactly what one git
// process per question answered, as the repository changes under it: items
// pending in several components, none, a release just cut, HEAD moved back,
// and a tag off HEAD's history.
func TestPendingFromTheKeptHistoryAnswersAsTheProcessesDid(t *testing.T) {
	root, r := gitRepo(t) // cli tagged 0.10.0, tpl at 1.0.0, web never released
	repo := &workitem.Repo{Root: root, Manifest: m}
	git := func(args ...string) {
		t.Helper()
		if out, err := r.Run(root, "git", args...); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	commitFiles := func(msg string, files map[string]string) {
		t.Helper()
		for rel, body := range files {
			p := filepath.Join(root, rel)
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		git("add", "-A")
		git("commit", "-q", "--allow-empty", "-m", msg)
	}
	same := func(step string, wantIDs ...string) {
		t.Helper()
		want, err := pendingByProcess(r, root, m, repo)
		if err != nil {
			t.Fatalf("%s: the reference: %v", step, err)
		}
		got, err := Pending(r, root, m, repo)
		if err != nil {
			t.Fatalf("%s: %v", step, err)
		}
		wj, _ := json.Marshal(want)
		gj, _ := json.Marshal(got)
		if string(wj) != string(gj) {
			t.Errorf("%s: plans differ\nprocesses: %s\nhistory:   %s", step, wj, gj)
		}
		ids := PendingIDs(r, root, m, repo)
		var names []string
		for id := range ids {
			names = append(names, id)
		}
		sort.Strings(names)
		sort.Strings(wantIDs)
		if strings.Join(names, " ") != strings.Join(wantIDs, " ") {
			t.Errorf("%s: PendingIDs %v, want %v", step, names, wantIDs)
		}
	}

	same("nothing accepted")

	// Story commits before the acceptance, one naming two items, one empty.
	// An item touching two components with no tag to say which it delivers
	// to cannot be planned, and is left out of both's plans; PendingIDs still
	// names it, as unpublished (I-0024).
	commitFiles("feat: [S-101] a command", map[string]string{"cli/a.go": "package main\n"})
	commitFiles("fix: [S-102] [S-101] both", map[string]string{"cli/b.go": "package main\n"})
	commitFiles("chore: [S-102] nothing", nil)
	commitFiles("feat: [S-100] two components", map[string]string{"cli/t.go": "package main\n", "web/t.js": "// t\n"})
	writeAcceptedItem(t, root, r, "S-101", workitem.Story, "feature", "A command", nil)
	writeAcceptedItem(t, root, r, "S-102", workitem.Story, "remediation", "Both", nil)
	writeAcceptedItem(t, root, r, "S-100", workitem.Story, "feature", "Two components", nil)
	same("two items, one left out", "S-100", "S-101", "S-102")

	writeAcceptedItem(t, root, r, "S-103", workitem.Story, "improvement", "Template text", map[string]string{"tpl/root/y.md": "y\n"})
	writeAcceptedItem(t, root, r, "S-104", workitem.Story, "research", "A finding", map[string]string{"cli/r.go": "package main\n"})
	writeAcceptedItem(t, root, r, "S-107", workitem.Story, "feature", "The web page", map[string]string{"web/e.js": "// e\n"})
	same("three components and research", "S-100", "S-101", "S-102", "S-103", "S-107")

	// Publish as flai release --pending does: bump the template's version
	// file, commit, and tag the code components.
	plans, err := Pending(r, root, m, repo)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range plans {
		if err := ApplyPending(p, root, time.Date(2026, 9, 29, 0, 0, 0, 0, time.UTC)); err != nil {
			t.Fatal(err)
		}
	}
	commitFiles("chore: release", nil)
	for _, p := range plans {
		if _, err := TagPending(r, root, p); err != nil {
			t.Fatal(err)
		}
	}
	same("a release just cut")

	writeAcceptedItem(t, root, r, "S-105", workitem.Story, "remediation", "After the release", map[string]string{"tpl/root/z.md": "z\n"})
	writeAcceptedItem(t, root, r, "S-108", workitem.Story, "remediation", "Also after", map[string]string{"cli/c.go": "package main\n"})
	same("accepted after the release", "S-105", "S-108")

	ahead, err := r.Run(root, "git", "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}
	git("reset", "-q", "--hard", "HEAD~2")
	same("HEAD moved back to the release")
	git("reset", "-q", "--hard", ahead)
	same("HEAD forward again to commits known but not reachable last time", "S-105", "S-108")
	// an acceptance naming no item there is named too, so nothing is lost
	git("commit", "-q", "--amend", "-m", "chore: [S-109] accept and archive")
	same("the last commit amended", "S-105", "S-109")

	// A web tag on a commit HEAD never reaches: git decides what follows it.
	git("checkout", "-q", "-b", "side")
	commitFiles("feat: [S-106] off the line", map[string]string{"web/side.js": "// side\n"})
	git("tag", "web/v9.0.0")
	git("checkout", "-q", "main")
	writeAcceptedItem(t, root, r, "S-107", workitem.Story, "remediation", "Web again", map[string]string{"web/f.js": "// f\n"})
	same("a tag off HEAD's history", "S-105", "S-107", "S-109")
}

// S-0157: while HEAD and the tags are unchanged, Pending starts one git
// process; an acceptance or a publish costs at most three.
func TestPendingStartsOneGitProcessWhileNothingChanged(t *testing.T) {
	root, sys := gitRepo(t)
	repo := &workitem.Repo{Root: root, Manifest: m}
	r := &counting{Runner: sys}
	writeAcceptedItem(t, root, sys, "S-101", workitem.Story, "remediation", "Fix one", map[string]string{"cli/a.go": "package main\n"})

	if _, err := Pending(r, root, m, repo); err != nil {
		t.Fatal(err)
	}
	if ran := r.take(); len(ran) > 4 {
		t.Errorf("the first look reads the history whole, in at most four processes: %v", ran)
	}
	for i := 0; i < 3; i++ {
		PendingIDs(r, root, m, repo)
		if ran := r.take(); len(ran) != 1 || ran[0] != "git show-ref" {
			t.Fatalf("look %d with nothing changed: %v", i, ran)
		}
	}

	writeAcceptedItem(t, root, sys, "S-102", workitem.Story, "feature", "Add", map[string]string{"cli/b.go": "package main\n"})
	ids := PendingIDs(r, root, m, repo)
	if ran := r.take(); len(ran) > 3 || !ids["S-101"] || !ids["S-102"] {
		t.Errorf("after an acceptance: %v, pending %v", ran, ids)
	}

	if _, err := sys.Run(root, "git", "tag", "-a", "cli/v0.11.0", "-m", "published"); err != nil {
		t.Fatal(err)
	}
	ids = PendingIDs(r, root, m, repo)
	if ran := r.take(); len(ran) > 3 || len(ids) != 0 {
		t.Errorf("after a publish: %v, pending %v", ran, ids)
	}
}

func TestBracketedIsWhatGrepFinds(t *testing.T) {
	got := bracketed("feat: [S-12] and [[T-0003]] [S-4 ] [x] [E-7]\n\nsee [S-99")
	var ids []string
	for id := range got {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	if strings.Join(ids, " ") != "E-7 S-12 T-0003" {
		t.Errorf("bracketed: %v", ids)
	}
}

// S-0157 on this repository: the kept history answers what the processes
// did, with its real tags, template, and archive.
func TestPendingOnThisRepositoryAnswersAsTheProcessesDid(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: reads the monorepo")
	}
	root, err := filepath.Abs(filepath.Join("..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "system-flow.yaml")); err != nil {
		t.Skip("monorepo not present")
	}
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	r := execx.System{}
	want, err := pendingByProcess(r, root, repo.Manifest, repo)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range want {
		t.Logf("%s: %d items pending, %s", p.Component.Name, len(p.Items), p.Level)
	}
	for _, look := range []string{"first", "kept"} {
		got, err := Pending(r, root, repo.Manifest, repo)
		if err != nil {
			t.Fatal(err)
		}
		wj, _ := json.Marshal(want)
		gj, _ := json.Marshal(got)
		if string(wj) != string(gj) {
			t.Errorf("%s look: plans differ\nprocesses: %s\nhistory:   %s", look, wj, gj)
		}
	}
}
