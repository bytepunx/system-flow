package hostapi

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

const rid = `"request_id":"req-00000001"`

// good is one valid call per write method and the command line it must become.
// TestEveryWriteIsCovered fails when a method is added without an entry here
// and in refused below.
var good = map[string]struct {
	params string
	args   string
	stdin  string
}{
	"item.move":         {`{"id":"S-0001","to":"cancelled","reason":"  not\n needed ",` + rid + `}`, "move S-0001 cancelled --by=olive --reason=not needed --json", ""},
	"item.move.preview": {`{"id":"E-0001"}`, "move E-0001 cancelled --by=olive --reason=preview --dry-run --json", ""},
	"item.order":        {`{"id":"S-0002","before":"S-0001",` + rid + `}`, "order S-0002 --before=S-0001 --json", ""},
	"item.block":        {`{"id":"T-0001","reason":"tide",` + rid + `}`, "block T-0001 --reason=tide --json", ""},
	"item.unblock":      {`{"id":"T-0001",` + rid + `}`, "unblock T-0001 --json", ""},
	"item.new":          {`{"type":"story","title":" --json  is my title ","parent":"E-0001","tags":["cli"],"touches":["flai/cmd"],"body":"## Goal\nx\n",` + rid + `}`, "story new --nature=feature --owner=olive --epic=E-0001 --tag=cli --touches=flai/cmd --body-stdin --autocommit --trailer=" + Trailer + " --json -- --json is my title", "## Goal\nx\n"},
	"item.template":     {`{"type":"epic"}`, "epic new --print-body --json", ""},
	"item.show":         {`{"id":"S-0001"}`, "edit S-0001 --show --json", ""},
	"item.edit":         {`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","title":" --json  is my title ","nature":"remediation","tags":["cli","dashboard"],"touches":[],"parent":"E-0002","body":"## Goal\nx\n",` + rid + `}`, "edit S-0001 --hash=" + strings.Repeat("a", 64) + " --by=olive --autocommit --trailer=" + Trailer + " --title=--json is my title --nature=remediation --parent=E-0002 --tag=cli --tag=dashboard --clear-touches --body-stdin --json", "## Goal\nx\n"},
	"accept.preview":    {`{"id":"S-0001"}`, "accept S-0001 --dry-run --json", ""},
	"accept.run":        {`{"id":"S-0001","include_uncommitted":true,` + rid + `}`, "accept S-0001 --by=olive --yes --json", ""},
	"stream.diff":       {`{"id":"S-0001"}`, "stream diff S-0001 --json", ""},
	"stream.log":        {`{"id":"S-0001","entry":"--not a flag",` + rid + `}`, "stream log S-0001 --json -- --not a flag", ""},
	"thread.new":        {`{"on":"S-0001","heading":"Goal","title":"How deep?","text":"Eight metres?",` + rid + `}`, "thread new --on=S-0001 --by=olive --heading=Goal --json -- How deep? Eight metres?", ""},
	"thread.reply":      {`{"id":"TH-0001","text":"Nine.",` + rid + `}`, "thread reply TH-0001 --by=olive --json -- Nine.", ""},
	"thread.resolve":    {`{"id":"TH-0001","reason":"answered",` + rid + `}`, "thread resolve TH-0001 --by=olive --reason=answered --json", ""},
	"doc.show":          {`{"path":"design/system/overview.md"}`, "doc show --json -- design/system/overview.md", ""},
	"doc.save":          {`{"path":"design/system/overview.md","content":"# new\n","hash":"abcdef0123456789","message":"tidy",` + rid + `}`, "doc save --hash=abcdef0123456789 --trailer=" + Trailer + " --message=tidy --json -- design/system/overview.md", "# new\n"},
	"adr.new":           {`{"title":"Decide it","status":"accepted","refines":["ADR-0016"],"supersedes":["7"],"body":"## Context\nx\n",` + rid + `}`, "adr new --status=accepted --supersedes=7 --refines=16 --body-stdin --autocommit --trailer=" + Trailer + " --json -- Decide it", "## Context\nx\n"},
	"adr.template":      {`{}`, "adr new --print-body --json", ""},
	"adr.accept":        {`{"id":"ADR-0028",` + rid + `}`, "adr accept 28 --autocommit --trailer=" + Trailer + " --json", ""},
	"stats.get":         {`{"since":"30d","type":"story","by":"nature"}`, "stats --since=30d --type=story --by=nature --json", ""},
	"push.pending":      {`{}`, "push --pending --dry-run --json", ""},
	"push.run":          {`{` + rid + `}`, "push --pending --publish --json", ""},
}

// refused is, per method, params that must never reach a command line.
var refused = map[string][]string{
	"item.move":         {`{"id":"--help","to":"ready",` + rid + `}`, `{"id":"S-0001","to":"--yes",` + rid + `}`, `{"id":"S-0001; rm -rf /","to":"ready",` + rid + `}`, `{"id":"S-0001","to":"ready"}`, `{"id":"S-0001","to":"ready","request_id":"x"}`},
	"item.move.preview": {`{"id":"../x"}`},
	"item.order":        {`{"id":"S-0002","before":"--top",` + rid + `}`, `{"id":"S-0002","top":true,"bottom":true,` + rid + `}`, `{"id":"S-0002",` + rid + `}`},
	"item.block":        {`{"id":"T-0001","reason":"  ",` + rid + `}`, `{"id":"T","reason":"x",` + rid + `}`},
	"item.unblock":      {`{"id":"--json",` + rid + `}`},
	"item.show":         {`{"id":"--show"}`, `{"id":"../S-0001"}`},
	"item.edit": {
		`{"id":"S-0001","title":"no hash",` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `",` + rid + `}`,
		`{"id":"--help","hash":"` + strings.Repeat("a", 64) + `","title":"t",` + rid + `}`,
		`{"id":"S-0001","hash":"--autocommit","title":"t",` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","title":"  ",` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","nature":"--by=eve",` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","parent":"--parent=E-1",` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","tags":["--by=eve"],` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","touches":["a,b"],` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","body":"  ",` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","title":"t"}`,
	},
	"item.new":       {`{"type":"task","title":"t","body":"b",` + rid + `}`, `{"type":"story","title":"t","body":"b","parent":"S-0001",` + rid + `}`, `{"type":"epic","title":"t","body":"b","parent":"E-0001",` + rid + `}`, `{"type":"epic","title":"t","body":"b","tags":["a,b"],` + rid + `}`, `{"type":"epic","title":"t","body":"b","tags":["--owner=eve"],` + rid + `}`, `{"type":"epic","title":"t","body":"b","nature":"urgent",` + rid + `}`, `{"type":"epic","title":"","body":"b",` + rid + `}`, `{"type":"epic","title":"t","body":"  ",` + rid + `}`},
	"item.template":  {`{"type":"task"}`},
	"accept.preview": {`{"id":"--no-push"}`},
	"accept.run":     {`{"id":"S-1 --no-push",` + rid + `}`},
	"stream.diff":    {`{"id":"../../etc"}`},
	"stream.log":     {`{"id":"S-0001","entry":"  ",` + rid + `}`},
	"thread.new":     {`{"on":"../secret.md","title":"t","text":"x",` + rid + `}`, `{"on":"--by=eve","title":"t","text":"x",` + rid + `}`, `{"on":"S-0001","title":"","text":"x",` + rid + `}`, `{"on":"src/main.go","title":"t","text":"x",` + rid + `}`},
	"thread.reply":   {`{"id":"--by=eve","text":"x",` + rid + `}`, `{"id":"TH-0001","text":"",` + rid + `}`},
	"thread.resolve": {`{"id":"TH-1; ls",` + rid + `}`},
	"doc.show":       {`{"path":"../../etc/passwd"}`, `{"path":"system-flow.yaml"}`, `{"path":"-rf.md"}`, `{"path":"src/notes.md"}`},
	"doc.save":       {`{"path":"design/x.md","content":"c","hash":"--force",` + rid + `}`, `{"path":".git/hooks/pre-push.md","content":"c","hash":"abcdef0123",` + rid + `}`, `{"path":"design/x.md","hash":"abcdef0123",` + rid + `}`, `{"path":"design/x.sh","content":"c","hash":"abcdef0123",` + rid + `}`},
	"adr.new":        {`{"title":"t","body":"b","status":"final",` + rid + `}`, `{"title":"t","body":"b","refines":["--autocommit"],` + rid + `}`, `{"title":"t","body":"b","supersedes":["0"],` + rid + `}`},
	"adr.template":   {`[]`},
	"adr.accept":     {`{"id":"--trailer=x",` + rid + `}`},
	"stats.get":      {`{"since":"30d; ls"}`, `{"type":"folder"}`, `{"by":"owner"}`},
	"push.pending":   {`"--force"`},
	"push.run":       {`"--force"`, `{}`},
}

type recorder struct {
	runs []Run
	ran  Ran
}

func (r *recorder) run(_ context.Context, run Run) (Ran, error) {
	r.runs = append(r.runs, run)
	if run.OnEvent != nil {
		for _, ev := range r.ran.Events {
			run.OnEvent(ev)
		}
	}
	return r.ran, nil
}

// hostFor is a host nobody has touched, except for a method that is itself a
// host action, which can only show what it runs with that action enabled.
func hostFor(name string) Host {
	action := specs(Host{})[name].action
	if action == "" {
		return Host{}
	}
	return Host{Enabled: func(a, _ string) bool { return a == action }}
}

func TestEveryWriteIsCovered(t *testing.T) {
	for name := range specs(Host{}) {
		if _, ok := good[name]; !ok {
			t.Errorf("%s has no valid call in good: add one with the command line it must become", name)
		}
		if len(refused[name]) == 0 {
			t.Errorf("%s has no refusals in refused: say what must never reach its command line", name)
		}
	}
	if len(good) != len(specs(Host{})) {
		t.Errorf("good names %d methods, there are %d", len(good), len(specs(Host{})))
	}
}

func TestWritesBecomeTheCommandLinesTheDashboardUsedToRun(t *testing.T) {
	p := withDocs(t) // owner: olive
	for name, c := range good {
		rec := &recorder{ran: Ran{Stdout: []byte(`{"ok":true}`)}}
		res, rerr := writeMethods(rec.run, time.Now, hostFor(name))[name](context.Background(), p, json.RawMessage(c.params))
		if rerr != nil {
			t.Errorf("%s: %+v", name, rerr)
			continue
		}
		if len(rec.runs) != 1 || strings.Join(rec.runs[0].Args, " ") != c.args || rec.runs[0].Stdin != c.stdin || rec.runs[0].Dir != p.Root {
			t.Errorf("%s ran\n  %q\nwant\n  %q (stdin %q, dir %s)", name, strings.Join(rec.runs[0].Args, " "), c.args, rec.runs[0].Stdin, rec.runs[0].Dir)
		}
		if w, ok := res.(Written); !ok || string(w.Data) != `{"ok":true}` || w.Warnings == nil {
			t.Errorf("%s answered %#v", name, res)
		}
	}
}

func TestWhatIsNotDataNeverReachesACommandLine(t *testing.T) {
	p := withDocs(t)
	for name, list := range refused {
		for _, params := range list {
			rec := &recorder{}
			_, rerr := writeMethods(rec.run, time.Now, Host{})[name](context.Background(), p, json.RawMessage(params))
			if rerr == nil || rerr.Code != channel.CodeInvalidParams || len(rec.runs) != 0 {
				t.Errorf("%s %s: error %+v, ran %d command(s)", name, params, rerr, len(rec.runs))
			}
		}
	}
}

func TestARepeatedWriteIsAnsweredFromTheJournal(t *testing.T) {
	p := withDocs(t)
	now := time.Date(2026, 9, 20, 9, 0, 0, 0, time.UTC)
	for name, c := range good {
		if specs(Host{})[name].reads {
			continue
		}
		rec := &recorder{ran: Ran{Stdout: []byte(`{"n":1}`)}}
		m := writeMethods(rec.run, func() time.Time { return now }, hostFor(name))[name]
		first, e1 := m(context.Background(), p, json.RawMessage(c.params))
		rec.ran = Ran{Stdout: []byte(`{"n":2}`)}
		second, e2 := m(context.Background(), p, json.RawMessage(c.params))
		if e1 != nil || e2 != nil || len(rec.runs) != 1 || string(first.(Written).Data) != string(second.(Written).Data) {
			t.Errorf("%s: ran %d times; %v %v", name, len(rec.runs), e1, e2)
		}
		// another request ID is another write, and an old entry is forgotten
		other := strings.Replace(c.params, "req-00000001", "req-00000002", 1)
		if _, e := m(context.Background(), p, json.RawMessage(other)); e != nil || len(rec.runs) != 2 {
			t.Errorf("%s: a new request ID must run: %d %v", name, len(rec.runs), e)
		}
		now = now.Add(journalKeeps + time.Minute)
		if _, e := m(context.Background(), p, json.RawMessage(c.params)); e != nil || len(rec.runs) != 3 {
			t.Errorf("%s: a forgotten request ID must run: %d %v", name, len(rec.runs), e)
		}
	}
}

func TestOutcomes(t *testing.T) {
	p := withDocs(t)
	call := func(name string, ran Ran) (any, *channel.Error) {
		rec := &recorder{ran: ran}
		return writeMethods(rec.run, time.Now, Host{})[name](context.Background(), p, json.RawMessage(good[name].params))
	}
	fatal := func(msg string) []map[string]any { return []map[string]any{{"level": "FATAL", "err": msg}} }

	res, _ := call("item.move", Ran{Stdout: []byte(`{"id":"S-0001"}`), Events: []map[string]any{{"level": "WARN", "msg": "policy", "detail": "WIP limit for ready is 5, this makes 6"}, {"level": "INFO", "msg": "x"}}})
	if w := res.(Written); len(w.Warnings) != 1 || w.Warnings[0] != "WIP limit for ready is 5, this makes 6" {
		t.Errorf("warnings: %+v", w)
	}
	if _, e := call("item.move", Ran{Exit: 1, Events: fatal("rule: S-0001 cannot go from review to cancelled")}); e == nil || e.Code != Rule || e.Message != "S-0001 cannot go from review to cancelled" {
		t.Errorf("a rule: %+v", e)
	}
	if _, e := call("item.move", Ran{Exit: 1, Events: fatal("disk full")}); e == nil || e.Code != channel.CodeInternal || e.Message != "disk full" {
		t.Errorf("a failure: %+v", e)
	}
	// S-0077: a write aimed at an item that is not there is 404, as a read is,
	// not a failure of flai's. The words are workitem's own, taken from the
	// lookup itself, so a change of wording fails here and not in a dashboard.
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	_, missing := repo.Get("S-9999")
	if missing == nil {
		t.Fatal("S-9999 should not exist in the fixture")
	}
	for _, said := range []string{missing.Error(), "move: " + missing.Error()} {
		if _, e := call("item.move", Ran{Exit: 1, Events: fatal(said)}); e == nil || e.Code != NotFound || e.Message != said {
			t.Errorf("a missing item, said as %q: %+v", said, e)
		}
	}
	if _, e := call("item.move", Ran{Exit: 1, Events: fatal("the remote was not found")}); e == nil || e.Code != channel.CodeInternal {
		t.Errorf("only an item's absence is a 404: %+v", e)
	}
	_, e := call("doc.save", Ran{Exit: 3, Stdout: []byte(`{"conflict":{"hash":"h2","current":"# theirs\n"}}`), Events: fatal("conflict: the file changed")})
	if e == nil || e.Code != Conflict || e.Message != "the file changed" || e.Data.(map[string]any)["hash"] != "h2" {
		t.Errorf("a conflict: %+v", e)
	}
	_, e = call("item.new", Ran{Exit: 4, Stdout: []byte(`{"refused":{"findings":[{"rule":"item.heading"}]}}`), Events: fatal("refused: flai check has 1 finding(s)")})
	if e == nil || e.Code != Refused || e.Data.(map[string]any)["findings"] == nil {
		t.Errorf("a refusal: %+v", e)
	}
	if _, e := call("item.move", Ran{Stdout: []byte("not json")}); e == nil || !strings.Contains(e.Message, "not JSON") {
		t.Errorf("output that is not JSON: %+v", e)
	}
	if res, e := call("item.unblock", Ran{}); e != nil || string(res.(Written).Data) != "null" {
		t.Errorf("no output: %+v %+v", res, e)
	}
}

func TestAcceptanceSendsEachStepAsProgress(t *testing.T) {
	p := withDocs(t)
	rec := &recorder{ran: Ran{Stdout: []byte(`{"id":"S-0001"}`), Events: []map[string]any{{"level": "INFO", "msg": "story branch merged"}, {"level": "INFO", "msg": "archived"}}}}
	var steps []string
	ctx := channel.WithProgress(context.Background(), func(v any) { steps = append(steps, v.(map[string]any)["msg"].(string)) })
	if _, e := writeMethods(rec.run, time.Now, Host{})["accept.run"](ctx, p, json.RawMessage(good["accept.run"].params)); e != nil {
		t.Fatal(e)
	}
	if strings.Join(steps, "|") != "story branch merged|archived" {
		t.Errorf("steps: %v", steps)
	}
	// a move is not streamed
	steps = nil
	_, _ = writeMethods(rec.run, time.Now, Host{})["item.move"](ctx, p, json.RawMessage(good["item.move"].params))
	if len(steps) != 0 {
		t.Errorf("a move sent progress: %v", steps)
	}
}

// S-0078, I-0028: flai runs on the host as the operator. What a dashboard may
// make it do with their credentials is off until they enable it by name, and
// every such request is written down, whatever became of it.
func TestHostActionsAreOffUntilEnabledAndJournalled(t *testing.T) {
	p := withDocs(t) // owner: olive
	var journal []Entry
	on := false
	host := Host{
		Enabled: func(action, root string) bool { return on && action == ActionPush && root == p.Root },
		Record:  func(e Entry) { journal = append(journal, e) },
	}
	at := time.Date(2026, 9, 20, 13, 0, 0, 0, time.UTC)
	call := func(name string, ran Ran) ([]string, *channel.Error) {
		rec := &recorder{ran: ran}
		_, e := writeMethods(rec.run, func() time.Time { return at }, host)[name](context.Background(), p, json.RawMessage(good[name].params))
		var lines []string
		for _, r := range rec.runs {
			lines = append(lines, strings.Join(r.Args, " "))
		}
		return lines, e
	}

	// off: an acceptance never pushes or tags either way (S-0087); the push
	// action governs push.run, not accept.run
	ran, e := call("accept.run", Ran{Stdout: []byte(`{"id":"S-0001"}`)})
	if e != nil || len(ran) != 1 {
		t.Fatalf("an acceptance with the action off: %v %+v", ran, e)
	}
	if len(journal) != 0 {
		t.Errorf("acceptance is no host action: %+v", journal)
	}
	ran, e = call("push.run", Ran{})
	if e == nil || e.Code != Disabled || len(ran) != 0 || !strings.Contains(e.Message, "flai serve enable push") || e.Data.(map[string]any)["enable"] != "flai serve enable push" {
		t.Fatalf("a disabled action runs nothing and says what to run: %v %+v", ran, e)
	}
	if len(journal) != 1 || journal[0].Outcome != "disabled" || journal[0].Action != ActionPush || journal[0].By != "olive" || journal[0].Project != p.Key || journal[0].Root != p.Root || journal[0].At != "2026-09-20T13:00:00Z" || journal[0].RequestID == "" {
		t.Errorf("the refusal is journalled: %+v", journal)
	}

	// on: the same command line either way
	on, journal = true, nil
	ran, e = call("accept.run", Ran{Stdout: []byte(`{"id":"S-0001"}`)})
	if e != nil || len(ran) != 1 {
		t.Fatalf("an acceptance with the action on: %v %+v", ran, e)
	}
	ran, e = call("push.run", Ran{Stdout: []byte(`{"pushed":true,"unpushed":{"acceptances":["S-0001"],"tags":["cli/v1.1.0"]}}`)})
	if e != nil || len(ran) != 1 || ran[0] != "push --pending --publish --json" {
		t.Fatalf("push.run: %v %+v", ran, e)
	}
	_, _ = call("push.run", Ran{Stdout: []byte(`{"pushed":false,"reason":"nothing pending"}`)})
	_, e = call("push.run", Ran{Exit: 3, Events: []map[string]any{{"level": "FATAL", "err": "conflict: main and origin/main have diverged"}}})
	if e == nil || e.Code != Conflict || e.Message != "main and origin/main have diverged" {
		t.Errorf("a remote that moved is a conflict with its reason: %+v", e)
	}
	// accept.run is never journalled as a host action (S-0087): it never
	// pushes, so it never uses the push action.
	want := []struct{ method, outcome, detail string }{
		{"push.run", "done", "pushed with tags cli/v1.1.0"},
		{"push.run", "done", "nothing pushed: nothing pending"},
		{"push.run", "failed", "main and origin/main have diverged"},
	}
	if len(journal) != len(want) {
		t.Fatalf("journal: %+v", journal)
	}
	for i, w := range want {
		if got := journal[i]; got.Method != w.method || got.Outcome != w.outcome || got.Detail != w.detail || got.Action != ActionPush {
			t.Errorf("entry %d: %+v, want %+v", i, got, w)
		}
	}
}

// Nothing a dashboard can ask for reads or changes what is enabled.
func TestNoMethodTouchesTheHostConfiguration(t *testing.T) {
	for name := range Methods("test", nil) {
		for _, word := range []string{"config", "enable", "disable", "action", "journal"} {
			if strings.Contains(name, word) {
				t.Errorf("%s: the host's configuration and journal are the operator's, on the host", name)
			}
		}
	}
	for name, sp := range specs(Host{Enabled: func(string, string) bool { return true }}) {
		args, _, e := sp.build(channel.Project{Root: t.TempDir()}, json.RawMessage(good[name].params))
		if e != nil {
			continue
		}
		if len(args) > 0 && (args[0] == "config" || args[0] == "serve") {
			t.Errorf("%s runs flai %s", name, args[0])
		}
	}
	if info := enabledActions(Host{}, "/x"); len(info) != len(Actions) || info[ActionPush] {
		t.Errorf("project.info names every action and says none is on: %+v", info)
	}
}

// S-0079: the dashboard is told what flai serve did about starting an agent,
// and can do nothing about it.
func TestAgentStatusIsReadOnly(t *testing.T) {
	p := withDocs(t)
	state := map[string]any{"command": "claude", "waiting": "an agent is attending"}
	host := Host{
		Enabled: func(action, root string) bool { return action == ActionAgent && root == p.Root },
		Agent:   func(string) any { return state },
	}
	res, e := MethodsFor("test", nil, host)["agent.status"](context.Background(), p, json.RawMessage(`{}`))
	if e != nil {
		t.Fatal(e)
	}
	got, _ := json.Marshal(res)
	if string(got) != `{"enabled":true,"state":{"command":"claude","waiting":"an agent is attending"}}` {
		t.Errorf("status: %s", got)
	}
	res, _ = Methods("test", nil)["agent.status"](context.Background(), p, json.RawMessage(`{}`))
	if got, _ := json.Marshal(res); string(got) != `{"enabled":false}` {
		t.Errorf("on a host nobody touched: %s", got)
	}
	for name := range Methods("test", nil) {
		if strings.HasPrefix(name, "agent.") && name != "agent.status" {
			t.Errorf("%s: nothing the dashboard can ask for starts, stops, or configures an agent", name)
		}
	}
}
