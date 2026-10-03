package cmd

import (
	"context"
	"encoding/json"
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
	// E-0001 follows its stories into review, where it cannot be cancelled,
	// and is accepted with the last of them (S-0200): the preview is of an
	// epic with nothing under it.
	if _, errOut, code := runIn(t, root, "epic", "new", "Spare"); code != 0 {
		t.Fatal(errOut)
	}
	all := func() (publish string) {
		t.Helper()
		says("item.show", answered("item.show", `{"id":"S-0002"}`, "edit", "S-0002", "--show"), `"hash"`)
		says("item.move.preview", answered("item.move.preview", `{"id":"E-0002"}`, "move", "E-0002", "cancelled", "--by=designer", "--reason=preview", "--dry-run"), `"dry_run":true`)
		answered("stats.get", `{}`, "stats")
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
