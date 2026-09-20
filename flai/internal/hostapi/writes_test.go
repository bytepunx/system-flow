package hostapi

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
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
}

// refused is, per method, params that must never reach a command line.
var refused = map[string][]string{
	"item.move":         {`{"id":"--help","to":"ready",` + rid + `}`, `{"id":"S-0001","to":"--yes",` + rid + `}`, `{"id":"S-0001; rm -rf /","to":"ready",` + rid + `}`, `{"id":"S-0001","to":"ready"}`, `{"id":"S-0001","to":"ready","request_id":"x"}`},
	"item.move.preview": {`{"id":"../x"}`},
	"item.order":        {`{"id":"S-0002","before":"--top",` + rid + `}`, `{"id":"S-0002","top":true,"bottom":true,` + rid + `}`, `{"id":"S-0002",` + rid + `}`},
	"item.block":        {`{"id":"T-0001","reason":"  ",` + rid + `}`, `{"id":"T","reason":"x",` + rid + `}`},
	"item.unblock":      {`{"id":"--json",` + rid + `}`},
	"item.new":          {`{"type":"task","title":"t","body":"b",` + rid + `}`, `{"type":"story","title":"t","body":"b","parent":"S-0001",` + rid + `}`, `{"type":"epic","title":"t","body":"b","parent":"E-0001",` + rid + `}`, `{"type":"epic","title":"t","body":"b","tags":["a,b"],` + rid + `}`, `{"type":"epic","title":"t","body":"b","tags":["--owner=eve"],` + rid + `}`, `{"type":"epic","title":"t","body":"b","nature":"urgent",` + rid + `}`, `{"type":"epic","title":"","body":"b",` + rid + `}`, `{"type":"epic","title":"t","body":"  ",` + rid + `}`},
	"item.template":     {`{"type":"task"}`},
	"accept.preview":    {`{"id":"--no-push"}`},
	"accept.run":        {`{"id":"S-1 --no-push",` + rid + `}`},
	"stream.diff":       {`{"id":"../../etc"}`},
	"stream.log":        {`{"id":"S-0001","entry":"  ",` + rid + `}`},
	"thread.new":        {`{"on":"../secret.md","title":"t","text":"x",` + rid + `}`, `{"on":"--by=eve","title":"t","text":"x",` + rid + `}`, `{"on":"S-0001","title":"","text":"x",` + rid + `}`, `{"on":"src/main.go","title":"t","text":"x",` + rid + `}`},
	"thread.reply":      {`{"id":"--by=eve","text":"x",` + rid + `}`, `{"id":"TH-0001","text":"",` + rid + `}`},
	"thread.resolve":    {`{"id":"TH-1; ls",` + rid + `}`},
	"doc.show":          {`{"path":"../../etc/passwd"}`, `{"path":"system-flow.yaml"}`, `{"path":"-rf.md"}`, `{"path":"src/notes.md"}`},
	"doc.save":          {`{"path":"design/x.md","content":"c","hash":"--force",` + rid + `}`, `{"path":".git/hooks/pre-push.md","content":"c","hash":"abcdef0123",` + rid + `}`, `{"path":"design/x.md","hash":"abcdef0123",` + rid + `}`, `{"path":"design/x.sh","content":"c","hash":"abcdef0123",` + rid + `}`},
	"adr.new":           {`{"title":"t","body":"b","status":"final",` + rid + `}`, `{"title":"t","body":"b","refines":["--autocommit"],` + rid + `}`, `{"title":"t","body":"b","supersedes":["0"],` + rid + `}`},
	"adr.template":      {`[]`},
	"adr.accept":        {`{"id":"--trailer=x",` + rid + `}`},
	"stats.get":         {`{"since":"30d; ls"}`, `{"type":"folder"}`, `{"by":"owner"}`},
	"push.pending":      {`"--force"`},
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

func TestEveryWriteIsCovered(t *testing.T) {
	for name := range specs() {
		if _, ok := good[name]; !ok {
			t.Errorf("%s has no valid call in good: add one with the command line it must become", name)
		}
		if len(refused[name]) == 0 {
			t.Errorf("%s has no refusals in refused: say what must never reach its command line", name)
		}
	}
	if len(good) != len(specs()) {
		t.Errorf("good names %d methods, there are %d", len(good), len(specs()))
	}
}

func TestWritesBecomeTheCommandLinesTheDashboardUsedToRun(t *testing.T) {
	p := withDocs(t) // owner: olive
	for name, c := range good {
		rec := &recorder{ran: Ran{Stdout: []byte(`{"ok":true}`)}}
		res, rerr := writeMethods(rec.run, time.Now)[name](context.Background(), p, json.RawMessage(c.params))
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
			_, rerr := writeMethods(rec.run, time.Now)[name](context.Background(), p, json.RawMessage(params))
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
		if specs()[name].reads {
			continue
		}
		rec := &recorder{ran: Ran{Stdout: []byte(`{"n":1}`)}}
		m := writeMethods(rec.run, func() time.Time { return now })[name]
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
		return writeMethods(rec.run, time.Now)[name](context.Background(), p, json.RawMessage(good[name].params))
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
	if _, e := writeMethods(rec.run, time.Now)["accept.run"](ctx, p, json.RawMessage(good["accept.run"].params)); e != nil {
		t.Fatal(e)
	}
	if strings.Join(steps, "|") != "story branch merged|archived" {
		t.Errorf("steps: %v", steps)
	}
	// a move is not streamed
	steps = nil
	_, _ = writeMethods(rec.run, time.Now)["item.move"](ctx, p, json.RawMessage(good["item.move"].params))
	if len(steps) != 0 {
		t.Errorf("a move sent progress: %v", steps)
	}
}
