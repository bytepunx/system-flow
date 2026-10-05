package cmd

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// sameAnswer runs the command a read of the dashboard's used to start, with
// --json, and the read as flai serve now answers it in its own process, on
// one fixture at one instant, and fails when they answer differently: the
// JSON, the warnings, or whether it failed and with what (S-0159). It
// returns the read's answer, or its error for the caller to check its code.
func sameAnswer(t *testing.T, root string, host hostapi.Host, method, params string, args ...string) (string, *channel.Error) {
	t.Helper()
	at := time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC)
	out, errOut, code := runInAt(t, root, at, append(args, "--json")...)
	ctx, rec := perf.Start(context.Background())
	res, rerr := hostapi.MethodsFor("test", func() time.Time { return at }, host)[method](ctx, channel.Project{Key: "t", Root: root}, json.RawMessage(params))
	for _, ph := range rec.Phases() {
		if strings.HasPrefix(ph.Name, "exec.flai") {
			t.Errorf("%s started flai: phase %s", method, ph.Name)
		}
	}
	warnings, fatal := []string{}, ""
	for _, line := range strings.Split(errOut, "\n") {
		var ev map[string]any
		if json.Unmarshal([]byte(line), &ev) != nil {
			continue
		}
		switch ev["level"] {
		case "WARN":
			w, _ := ev["detail"].(string)
			if w == "" {
				w, _ = ev["msg"].(string)
			}
			warnings = append(warnings, w)
		case "FATAL":
			fatal, _ = ev["err"].(string)
		}
	}
	if code != 0 {
		if rerr == nil || !strings.Contains(fatal, rerr.Message) {
			t.Errorf("%s: flai %v failed with %q, the read answered %+v", method, args, fatal, rerr)
		}
		return "", rerr
	}
	if rerr != nil {
		t.Errorf("%s: flai %v answered, the read failed: %+v", method, args, rerr)
		return "", rerr
	}
	w, ok := res.(hostapi.Written)
	if !ok {
		t.Fatalf("%s answered %T, not what flai's answer became", method, res)
	}
	var want, got any
	if err := json.Unmarshal([]byte(out), &want); err != nil {
		t.Fatalf("flai %v: %v\n%s", args, err, out)
	}
	if err := json.Unmarshal(w.Data, &got); err != nil {
		t.Fatalf("%s: %v\n%s", method, err, w.Data)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("%s: flai %v printed\n%s\nthe read answered\n%s", method, args, out, w.Data)
	}
	if !reflect.DeepEqual(warnings, w.Warnings) {
		t.Errorf("%s: flai warned %q, the read %q", method, warnings, w.Warnings)
	}
	return string(w.Data), nil
}

// S-0217: the orchestrator's reads answer as flai order --by, flai promote
// --candidates, and flai release --evaluate print them, and the order read
// never applies its order.
func TestTheOrchestrationReadsAnswerWhatTheCommandsPrint(t *testing.T) {
	root, _ := evaluateProject(t)
	t.Setenv("LOG_FORMAT", "json")
	run := func(args ...string) {
		t.Helper()
		if _, errOut, code := runIn(t, root, args...); code != 0 {
			t.Fatalf("flai %v: %s", args, errOut)
		}
	}
	fill := func(id string) {
		t.Helper()
		file, _ := filepath.Glob(filepath.Join(root, "wip/kanban/stories", id+"-*.md"))
		s, _ := os.ReadFile(file[0])
		body := strings.Replace(string(s), "## Goal\n", "## Goal\n\nDo it.\n", 1)
		body = strings.Replace(body, "## Acceptance criteria\n- [ ]\n", "## Acceptance criteria\n- [ ] ok\n", 1)
		if err := os.WriteFile(file[0], []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// ready: S-0004 valued and forecast, S-0005 with neither; backlog: S-0003
	// with neither, S-0006 a candidate
	for _, s := range [][]string{{"Valued", "docs"}, {"Bare", "template"}, {"Candidate", "scripts"}} {
		run("story", "new", s[0], "--epic", "E-0001", "--touches", s[1])
	}
	run("edit", "S-0004", "--cost-of-delay-value", "200", "--forecast-duration", "2h")
	run("edit", "S-0006", "--cost-of-delay-value", "50", "--forecast-duration", "1h")
	for _, id := range []string{"S-0004", "S-0005", "S-0006"} {
		fill(id)
	}
	run("move", "S-0005", "ready")
	run("move", "S-0004", "ready")
	setReleasePolicy(t, root, "    policy: threshold\n    value: 300\n")
	manifest := filepath.Join(root, "system-flow.yaml")
	data, _ := os.ReadFile(manifest)
	if err := os.WriteFile(manifest, append(data, "  policy: cod\n"...), 0o644); err != nil {
		t.Fatal(err)
	}
	board := filepath.Join(root, "wip/kanban/board.md")
	before, _ := os.ReadFile(board)

	answer := func(method, params string, args ...string) string {
		t.Helper()
		got, rerr := sameAnswer(t, root, hostapi.Host{}, method, params, args...)
		if rerr != nil {
			t.Fatalf("%s %s: %+v", method, params, rerr)
		}
		return got
	}
	for _, policy := range workitem.OrderPolicies {
		answer("order.by", `{"policy":"`+policy+`"}`, "order", "--by", policy)
	}
	if got := answer("order.by", `{}`, "order", "--by", "cod"); !strings.Contains(got, `"policy":"cod"`) || !strings.Contains(got, `"applied":false`) || !strings.Contains(got, `"id":"S-0004"`) {
		t.Errorf("order.by with no policy is the manifest's: %s", got)
	}
	if after, _ := os.ReadFile(board); string(after) != string(before) {
		t.Error("order.by wrote board.md")
	}
	if got := answer("promote.candidates", `{}`, "promote", "--candidates"); !strings.Contains(got, `"id":"S-0006"`) || !strings.Contains(got, `"id":"S-0003"`) {
		t.Errorf("promote.candidates lacks S-0006 as a candidate or S-0003 refused: %s", got)
	}
	answer("promote.candidates", `{"limit":1}`, "promote", "--candidates", "--limit", "1")
	if got := answer("release.evaluate", `{}`, "release", "--evaluate"); !strings.Contains(got, `"met":true`) {
		t.Errorf("release.evaluate: 300 a week pending meets a threshold of 300: %s", got)
	}
	setReleasePolicy(t, root, "    policy: theme\n    epic: E-0001\n")
	if got := answer("release.evaluate", `{}`, "release", "--evaluate"); !strings.Contains(got, `"met":false`) {
		t.Errorf("release.evaluate: E-0001 has stories not accepted: %s", got)
	}
}

// TestTheReadsAnswerWhatTheCommandsPrint holds the reads flai serve answers
// in its own process to the commands the dashboard used to run for them,
// through a project's life: stories in review, accepted and not published,
// a remote that moved on. The dashboard publishes and never pushes
// (ADR-0067): it has no read of what flai push --pending would send.
func TestTheReadsAnswerWhatTheCommandsPrint(t *testing.T) {
	root, remote := pendingProject(t)
	t.Setenv("LOG_FORMAT", "json")
	host := hostapi.Host{Enabled: func(action, root string) bool {
		cfg, _, err := config.Load(config.ResolvePath(""))
		return err == nil && cfg.ActionEnabled(action, root)
	}}
	answered := func(method, params string, args ...string) string {
		t.Helper()
		data, _ := sameAnswer(t, root, host, method, params, args...)
		return data
	}
	same := func(method, params string, args ...string) *channel.Error {
		t.Helper()
		_, rerr := sameAnswer(t, root, host, method, params, args...)
		return rerr
	}
	// says fails when an answer lacks what the fixture is there to show
	says := func(what, data string, parts ...string) {
		t.Helper()
		for _, p := range parts {
			if !strings.Contains(data, p) {
				t.Errorf("%s lacks %s: %s", what, p, data)
			}
		}
	}
	// E-0001 follows its first story into review, where it cannot be
	// cancelled, and stays there (S-0200): the preview is of an epic with
	// nothing under it.
	if _, errOut, code := runIn(t, root, "epic", "new", "Spare"); code != 0 {
		t.Fatal(errOut)
	}
	// a strategic agent's activity document, which stats reads (ADR-0079)
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := repo.AppendActivity(workitem.ActivityPlanner, workitem.ActivityEntry{At: time.Now(), Summary: "Planned E-0001.", Items: []string{"E-0001"}, Seconds: 60, Cost: 0.1}); err != nil {
		t.Fatal(err)
	}
	// a thread on S-0002, which stats reads for the time its agent waited
	// (S-0205)
	if _, err := threads.New(repo, threads.NewOptions{Title: "Which way", On: "S-0002", Author: "agent", Text: "Left?", Now: time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC)}); err != nil {
		t.Fatal(err)
	}
	all := func() (publish string) {
		t.Helper()
		says("item.show", answered("item.show", `{"id":"S-0002"}`, "edit", "S-0002", "--show"), `"hash"`)
		says("item.move.preview", answered("item.move.preview", `{"id":"E-0002"}`, "move", "E-0002", "cancelled", "--by=designer", "--reason=preview", "--dry-run"), `"dry_run":true`)
		// and the threads, the board's limit, and the stories' commits
		// (S-0205): S-0002's on its branch changed cli/another.go
		says("stats.get", answered("stats.get", `{}`, "stats"), `"strategic":[{"kind":"planner","cost":0.1,"seconds":60,"activities":1`,
			`"waiting":{"weeks":[`, `"claims":{"limit":2,"days":[`, `"drift":[`, `{"id":"S-0002","committed":["cli/another.go"]`, `"wait_threads_seconds":`)
		answered("stats.get", `{"since":"12w","type":"task","by":"parent"}`, "stats", "--since=12w", "--type=task", "--by=parent")
		says("stats.get by the hour", answered("stats.get", `{"since":"7d","bucket":"hour"}`, "stats", "--since=7d", "--bucket=hour"), `"bucket":"hour"`, `"spend":{"epic":`)
		return answered("publish.preview", `{}`, "release", "--pending", "--dry-run")
	}

	publish := all()
	says("publish.preview before any acceptance", publish, `"plans":null`)
	says("stream.diff", answered("stream.diff", `{"id":"S-0002"}`, "stream", "diff", "S-0002"), `"cli/another.go"`)
	says("accept.preview", answered("accept.preview", `{"id":"S-0002"}`, "accept", "S-0002", "--dry-run"), `"branch":"story/S-0002"`)
	for method, id := range map[string]string{"item.show": "S-0099", "stream.diff": "S-0099", "accept.preview": "S-0099", "item.move.preview": "S-0099"} {
		args := map[string][]string{
			"item.show": {"edit", id, "--show"}, "stream.diff": {"stream", "diff", id},
			"accept.preview": {"accept", id, "--dry-run"}, "item.move.preview": {"move", id, "cancelled", "--by=designer", "--reason=preview", "--dry-run"},
		}[method]
		if rerr := same(method, `{"id":"`+id+`"}`, args...); rerr == nil || rerr.Code != hostapi.NotFound {
			t.Errorf("%s of an item that is not there: %+v, want not found", method, rerr)
		}
	}
	if rerr := same("accept.preview", `{"id":"T-0002"}`, "accept", "T-0002", "--dry-run"); rerr == nil || rerr.Code != channel.CodeInternal {
		t.Errorf("accept.preview of a task: %+v", rerr)
	}

	for _, id := range []string{"S-0001", "S-0002"} {
		if _, errOut, code := runIn(t, root, "accept", id); code != 0 {
			t.Fatalf("accept %s: %s", id, errOut)
		}
	}
	publish = all()
	says("publish.preview once accepted", publish, `"cli"`, `"S-0002"`)
	// with auto-publish on in a shell, still no read or write says what a
	// push would send, or sends it: asked for, each is a method there is not
	if _, errOut, code := runIn(t, root, "serve", "enable", "auto-publish"); code != 0 {
		t.Fatal(errOut)
	}
	for _, method := range []string{"push.pending", "push.run"} {
		if _, ok := hostapi.MethodsFor("test", nil, host)[method]; ok {
			t.Errorf("%s is a method the dashboard can call", method)
		}
		if out, errOut, code := runIn(t, root, "hostapi", method, `{"request_id":"3f0c1a52-7d3b-4f0e-9a51-0c2d4e6f8a30"}`); code == 0 || !strings.Contains(out+errOut, "flai offers no method "+method) {
			t.Errorf("%s answered as a method: %d %s %s", method, code, out, errOut)
		}
	}

	// someone else pushed: the remote has a commit this clone lacks
	other := filepath.Join(t.TempDir(), "other")
	gitIn(t, filepath.Dir(other), "clone", "-q", remote, other)
	gitIn(t, other, "config", "user.email", "o@o")
	gitIn(t, other, "config", "user.name", "o")
	gitIn(t, other, "commit", "-q", "--allow-empty", "-m", "elsewhere")
	gitIn(t, other, "push", "-q", "origin", "main")
	gitIn(t, root, "fetch", "-q", "origin")

	// published from another clone: this one lacks the tag (S-0174), which
	// a remote's answer kept from the reads above would not show yet
	defer func(d time.Duration) { release.RemoteTTL = d }(release.RemoteTTL)
	release.RemoteTTL = 0
	gitIn(t, remote, "tag", "cli/v1.2.0", "main")
	publish = answered("publish.preview", `{}`, "release", "--pending", "--dry-run")
	says("publish.preview from a clone missing the remote's tags", publish, `"plans":null`, `"remote":"cli/v1.2.0"`, `"branch":{"upstream":"origin/main"`, `"fix":"git fetch --tags origin \u0026\u0026 git merge origin/main"`)
}
