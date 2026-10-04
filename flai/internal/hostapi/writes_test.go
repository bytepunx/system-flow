package hostapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
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
	"item.order":        {`{"id":"S-0002","before":"S-0001",` + rid + `}`, "order S-0002 --before=S-0001 --json", ""},
	"item.block":        {`{"id":"T-0001","reason":"tide",` + rid + `}`, "block T-0001 --reason=tide --json", ""},
	"board.limit":       {`{"column":"in-progress","limit":3,` + rid + `}`, "board limit --json -- in-progress 3", ""},
	"item.unblock":      {`{"id":"T-0001",` + rid + `}`, "unblock T-0001 --json", ""},
	"item.new":          {`{"type":"story","title":" --json  is my title ","parent":"E-0001","tags":["cli"],"touches":["flai/cmd"],"topics":["logging"," release"],"body":"## Goal\nx\n",` + rid + `}`, "story new --nature=feature --owner=olive --epic=E-0001 --tag=cli --touches=flai/cmd --topics=logging --topics=release --body-stdin --autocommit --trailer=" + Trailer + " --json -- --json is my title", "## Goal\nx\n"},
	"item.template":     {`{"type":"epic"}`, "epic new --print-body --json", ""},
	"item.finalize":     {`{"id":"S-0001",` + rid + `}`, "edit S-0001 --no-draft --by=olive --autocommit --trailer=" + Trailer + " --json", ""},
	"issue.list":        {`{"story":"S-0198"}`, "issue list --story=S-0198 --json", ""},
	"issue.story":       {`{"id":"I-0007","epic":"E-0002",` + rid + `}`, "issue story I-0007 --owner=olive --autocommit --trailer=" + Trailer + " --epic=E-0002 --json", ""},
	"item.edit":         {`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","title":" --json  is my title ","nature":"remediation","tags":["cli","dashboard"],"touches":[],"topics":["logging"],"after":["S-0128","S-129"],"parent":"E-0002","body":"## Goal\nx\n",` + rid + `}`, "edit S-0001 --hash=" + strings.Repeat("a", 64) + " --by=olive --autocommit --trailer=" + Trailer + " --title=--json is my title --nature=remediation --parent=E-0002 --tag=cli --tag=dashboard --clear-touches --topics=logging --after=S-0128 --after=S-129 --body-stdin --json", "## Goal\nx\n"},
	"accept.run":        {`{"id":"S-0001","include_uncommitted":true,` + rid + `}`, "accept S-0001 --by=olive --yes --json", ""},
	"stream.log":        {`{"id":"S-0001","entry":"--not a flag",` + rid + `}`, "stream log S-0001 --json -- --not a flag", ""},
	"stream.answer":     {`{"id":"S-0001","question":"--not a flag?","answer":"--also not",` + rid + `}`, "stream answer S-0001 --by=olive --json -- --not a flag? --also not", ""},
	"thread.new":        {`{"on":"S-0001","heading":"Goal","title":"How deep?","text":"Eight metres?",` + rid + `}`, "thread new --on=S-0001 --by=olive --heading=Goal --json -- How deep? Eight metres?", ""},
	"thread.reply":      {`{"id":"TH-0001","text":"Nine.",` + rid + `}`, "thread reply TH-0001 --by=olive --json -- Nine.", ""},
	"thread.resolve":    {`{"id":"TH-0001","reason":"answered",` + rid + `}`, "thread resolve TH-0001 --by=olive --reason=answered --json", ""},
	"doc.show":          {`{"path":"design/system/overview.md"}`, "doc show --json -- design/system/overview.md", ""},
	"doc.save":          {`{"path":"design/system/overview.md","content":"# new\n","hash":"abcdef0123456789","message":"tidy",` + rid + `}`, "doc save --hash=abcdef0123456789 --trailer=" + Trailer + " --message=tidy --json -- design/system/overview.md", "# new\n"},
	"adr.new":           {`{"title":"Decide it","status":"accepted","refines":["ADR-0016"],"supersedes":["7"],"body":"## Context\nx\n",` + rid + `}`, "adr new --status=accepted --supersedes=7 --refines=16 --body-stdin --autocommit --trailer=" + Trailer + " --json -- Decide it", "## Context\nx\n"},
	"adr.template":      {`{}`, "adr new --print-body --json", ""},
	"adr.accept":        {`{"id":"ADR-0028",` + rid + `}`, "adr accept 28 --autocommit --trailer=" + Trailer + " --json", ""},
	"publish.run":       {`{` + rid + `}`, "release --pending --json", ""},
	"dashboard.status":  {`{}`, "dashboard status --json", ""},
	"dashboard.check":   {`{}`, "dashboard check --json", ""},
	"dashboard.restart": {`{` + rid + `}`, "dashboard restart --json", ""},
	"dashboard.upgrade": {`{` + rid + `}`, "dashboard upgrade --json", ""},
	"dashboard.stop":    {`{` + rid + `}`, "dashboard stop --json", ""},
	"checks.status":     {`{"id":"S-0001"}`, "checks status S-0001 --json", ""},
	// S-0106: flai host, through flai serve
	"host.status":   {`{}`, "host status --json", ""},
	"host.check":    {`{}`, "host check --json", ""},
	"host.process":  {`{"process":"serve","action":"restart",` + rid + `}`, "host restart serve --json", ""},
	"host.upgrade":  {`{` + rid + `}`, "host upgrade --json", ""},
	"checks.tail":   {`{"id":"S-0001","from":128}`, "checks tail S-0001 --from=128 --wait=20 --json", ""},
	"checks.run":    {`{"id":"S-0001",` + rid + `}`, "checks run S-0001 --json", ""},
	"checks.cancel": {`{"id":"S-0001",` + rid + `}`, "checks cancel S-0001 --json", ""},
	"agent.restart": {`{"id":"S-0001",` + rid + `}`, "serve agent restart S-0001 --json", ""},
	"agent.start":   {`{"id":"S-0001",` + rid + `}`, "serve agent start S-0001 --json", ""},
	"agent.commit":  {`{"id":"S-0001",` + rid + `}`, "serve agent commit S-0001 --json", ""},
	"agent.stop":    {`{"id":"S-0001",` + rid + `}`, "serve agent stop S-0001 --json", ""},
	"plan.run":      {`{"id":"E-0001",` + rid + `}`, "plan E-0001 --json", ""},
	// S-0105: the host's settings, each a flai command gated by the settings action
	"settings.action": {`{"action":"push","on":true,` + rid + `}`, "serve enable push --json", ""},
	"settings.default_agent": {`{"agent":{"harness":"claude-code","model":"claude-opus-5-5","config":{"effort":"high"},"roles":{"verify":{"model":"sonnet"}}},` + rid + `}`,
		"agent set --replace --autocommit --trailer=Co-Authored-By: flaiover <flaiover@localhost> --harness=claude-code --model=claude-opus-5-5 --config=effort=high --role-model=verify=sonnet --json", ""},
	"settings.agent":           {`{"name":"builder","attended_minutes":10,"command":["claude","-p","work on {story}; echo $HOME"],` + rid + `}`, "serve agent set --name=builder --json -- claude -p work on {story}; echo $HOME", ""},
	"settings.harness":         {`{"name":"claude-code","program":"/opt/claude","args":["--permission-mode","acceptEdits"],` + rid + `}`, "serve agent harness claude-code --program=/opt/claude --json -- --permission-mode acceptEdits", ""},
	"settings.check":           {`{"name":"flai","command":["scripts/flai-test.sh","--short"],` + rid + `}`, "serve checks set --name=flai --json -- scripts/flai-test.sh --short", ""},
	"settings.checks_timeout":  {`{"minutes":20,` + rid + `}`, "serve checks timeout 20 --json", ""},
	"settings.import":          {`{"folder":"/home/me/git","add":true,` + rid + `}`, "serve import add --json -- /home/me/git", ""},
	"settings.serve":           {`{"root":"/home/me/git/blog","key":"blog",` + rid + `}`, "serve project add --json -- /home/me/git/blog", ""},
	"settings.unserve":         {`{"root":"/home/me/git/blog",` + rid + `}`, "serve project remove --json -- /home/me/git/blog", ""},
	"settings.mcp_token":       {`{` + rid + `}`, "mcp token --rotate --json", ""},
	"settings.dashboard_token": {`{` + rid + `}`, "dashboard token --rotate --no-restart --json", ""},
}

// refused is, per method, params that must never reach a command line.
var refused = map[string][]string{
	"item.move":     {`{"id":"--help","to":"ready",` + rid + `}`, `{"id":"S-0001","to":"--yes",` + rid + `}`, `{"id":"S-0001; rm -rf /","to":"ready",` + rid + `}`, `{"id":"S-0001","to":"ready"}`, `{"id":"S-0001","to":"ready","request_id":"x"}`},
	"item.order":    {`{"id":"S-0002","before":"--top",` + rid + `}`, `{"id":"S-0002","top":true,"bottom":true,` + rid + `}`, `{"id":"S-0002",` + rid + `}`},
	"item.block":    {`{"id":"T-0001","reason":"  ",` + rid + `}`, `{"id":"T","reason":"x",` + rid + `}`},
	"board.limit":   {`{"column":"backlog","limit":3,` + rid + `}`, `{"column":"--json","limit":3,` + rid + `}`, `{"column":"ready","limit":-1,` + rid + `}`, `{"column":"ready","limit":100,` + rid + `}`, `{"column":"ready",` + rid + `}`},
	"item.unblock":  {`{"id":"--json",` + rid + `}`},
	"item.finalize": {`{"id":"--help",` + rid + `}`, `{"id":"E-0001",` + rid + `}`, `{"id":"T-0001",` + rid + `}`, `{"id":"S-0001 --draft",` + rid + `}`, `{"id":"S-0001"}`},
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
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","after":["--clear-after"],` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","topics":["--clear-topics"],` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","topics":["flai/cmd"],` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","body":"  ",` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","title":"t"}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","clear_cost_of_delay":true,"cost_of_delay":{"value":"5"},` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","clear_forecast":true,"forecast":{"duration":"6h"},` + rid + `}`,
		`{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `","cost_of_delay":{"value":5},` + rid + `}`,
	},
	"item.new":          {`{"type":"task","title":"t","body":"b",` + rid + `}`, `{"type":"story","title":"t","body":"b","parent":"S-0001",` + rid + `}`, `{"type":"epic","title":"t","body":"b","parent":"E-0001",` + rid + `}`, `{"type":"epic","title":"t","body":"b","tags":["a,b"],` + rid + `}`, `{"type":"epic","title":"t","body":"b","tags":["--owner=eve"],` + rid + `}`, `{"type":"epic","title":"t","body":"b","topics":["two words"],` + rid + `}`, `{"type":"epic","title":"t","body":"b","topics":["--owner=eve"],` + rid + `}`, `{"type":"epic","title":"t","body":"b","nature":"urgent",` + rid + `}`, `{"type":"epic","title":"","body":"b",` + rid + `}`, `{"type":"epic","title":"t","body":"  ",` + rid + `}`, `{"type":"epic","title":"t","body":"b","draft":true,` + rid + `}`},
	"item.template":     {`{"type":"task"}`},
	"issue.list":        {`{"story":"--all"}`, `{"story":"T-0001"}`, `{"story":"S-1 --all"}`, `[]`},
	"issue.story":       {`{"id":"--help",` + rid + `}`, `{"id":"S-0001",` + rid + `}`, `{"id":"I-0001 --story=S-1",` + rid + `}`, `{"id":"I-0001","epic":"S-0001",` + rid + `}`, `{"id":"I-0001","epic":"--json",` + rid + `}`, `{"id":"I-0001"}`},
	"accept.run":        {`{"id":"S-1 --no-push",` + rid + `}`},
	"stream.log":        {`{"id":"S-0001","entry":"  ",` + rid + `}`},
	"stream.answer":     {`{"id":"../../etc","question":"q","answer":"a",` + rid + `}`, `{"id":"S-0001","question":"  ","answer":"a",` + rid + `}`, `{"id":"S-0001","question":"q","answer":"  ",` + rid + `}`},
	"thread.new":        {`{"on":"../secret.md","title":"t","text":"x",` + rid + `}`, `{"on":"--by=eve","title":"t","text":"x",` + rid + `}`, `{"on":"S-0001","title":"","text":"x",` + rid + `}`, `{"on":"src/main.go","title":"t","text":"x",` + rid + `}`},
	"thread.reply":      {`{"id":"--by=eve","text":"x",` + rid + `}`, `{"id":"TH-0001","text":"",` + rid + `}`},
	"thread.resolve":    {`{"id":"TH-1; ls",` + rid + `}`},
	"doc.show":          {`{"path":"../../etc/passwd"}`, `{"path":"system-flow.yaml"}`, `{"path":"-rf.md"}`, `{"path":"src/notes.md"}`},
	"doc.save":          {`{"path":"design/x.md","content":"c","hash":"--force",` + rid + `}`, `{"path":".git/hooks/pre-push.md","content":"c","hash":"abcdef0123",` + rid + `}`, `{"path":"design/x.md","hash":"abcdef0123",` + rid + `}`, `{"path":"design/x.sh","content":"c","hash":"abcdef0123",` + rid + `}`},
	"adr.new":           {`{"title":"t","body":"b","status":"final",` + rid + `}`, `{"title":"t","body":"b","refines":["--autocommit"],` + rid + `}`, `{"title":"t","body":"b","supersedes":["0"],` + rid + `}`},
	"adr.template":      {`[]`},
	"adr.accept":        {`{"id":"--trailer=x",` + rid + `}`},
	"publish.run":       {`"--force"`, `{}`},
	"dashboard.status":  {`"--force"`},
	"dashboard.check":   {`"--force"`},
	"dashboard.restart": {`"--force"`, `{}`},
	"dashboard.upgrade": {`"--force"`, `{}`},
	"dashboard.stop":    {`"--force"`, `{}`},
	"checks.status":     {`{"id":"--help"}`, `{"id":"../S-0001"}`},
	"host.status":       {`"--force"`},
	"host.check":        {`"--force"`},
	"host.process": {`{"process":"--config=/tmp/x","action":"stop",` + rid + `}`, `{"process":"serve","action":"--help",` + rid + `}`,
		`{"process":"dashboard","action":"stop",` + rid + `}`, `{"process":"serve","action":"restart"}`},
	"host.upgrade":  {`"--force"`, `{}`},
	"checks.tail":   {`{"id":"S-0001","from":-1}`, `{"id":"--help","from":0}`},
	"checks.run":    {`{"id":"--help",` + rid + `}`, `{"id":"S-0001"}`},
	"checks.cancel": {`{"id":"--help",` + rid + `}`, `{"id":"S-0001"}`},
	"agent.restart": {`{"id":"--help",` + rid + `}`, `{"id":"T-0001",` + rid + `}`, `{"id":"S-0001"}`},
	"agent.start":   {`{"id":"--help",` + rid + `}`, `{"id":"T-0001",` + rid + `}`, `{"id":"S-0001"}`},
	"agent.commit":  {`{"id":"--help",` + rid + `}`, `{"id":"T-0001",` + rid + `}`, `{"id":"S-0001"}`},
	"agent.stop":    {`{"id":"--help",` + rid + `}`, `{"id":"T-0001",` + rid + `}`, `{"id":"S-0001"}`},
	"plan.run":      {`{"id":"--help",` + rid + `}`, `{"id":"T-0001",` + rid + `}`, `{"id":"E-0001 --json",` + rid + `}`, `{"id":"E-0001"}`},
	"settings.action": {`{"action":"settings","on":true,` + rid + `}`, `{"action":"settings","on":false,` + rid + `}`, `{"action":"--all-projects","on":true,` + rid + `}`,
		`{"action":"push",` + rid + `}`, `{"action":"push","on":true}`},
	"settings.default_agent": {`{"agent":{"harness":"--dangerously-skip-permissions"},` + rid + `}`, `{"agent":{"model":"m","config":{"Bad Key":"v"}},` + rid + `}`,
		`{"agent":{"model":"m","config":{"k":"a\nb"}},` + rid + `}`},
	"settings.agent": {`{"name":"--x",` + rid + `}`, `{"attended_minutes":10,` + rid + `}`, `{"command":[],` + rid + `}`, `{"command":["","x"],` + rid + `}`,
		`{"command":"claude -p",` + rid + `}`, `{"command":["claude","a\nb"],` + rid + `}`, `{` + rid + `}`},
	"settings.harness": {`{"name":"command","program":"x",` + rid + `}`, `{"name":"--help","program":"x",` + rid + `}`, `{"name":"claude-code",` + rid + `}`,
		`{"name":"claude-code","args":["a\nb"],` + rid + `}`, `{"name":"claude-code","args":"--x",` + rid + `}`},
	"settings.check":           {`{"name":"--x","command":["a"],` + rid + `}`, `{"name":"a","command":[],` + rid + `}`, `{"name":"a","command":[" "],` + rid + `}`},
	"settings.checks_timeout":  {`{"minutes":0,` + rid + `}`, `{"minutes":"20",` + rid + `}`},
	"settings.import":          {`{"folder":"relative/git","add":true,` + rid + `}`, `{"folder":"/a/../b","add":true,` + rid + `}`, `{"folder":"/a",` + rid + `}`},
	"settings.serve":           {`{"root":"blog",` + rid + `}`, `{"root":"/a/../b",` + rid + `}`, `{"root":"--help",` + rid + `}`, `{"root":"/a","key":"--x",` + rid + `}`, `{"root":"/a"}`},
	"settings.unserve":         {`{"key":"blog",` + rid + `}`, `{"root":"/a\nb",` + rid + `}`, `{"root":"/a","key":"a b",` + rid + `}`},
	"settings.mcp_token":       {`"--force"`, `{}`},
	"settings.dashboard_token": {`"--force"`, `{}`},
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
	action := specs()[name].action
	if action == "" {
		return Host{}
	}
	return Host{Enabled: func(a, _ string) bool { return a == action }}
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

// A story need not belong to an epic (S-0092): item.new leaves out --epic
// entirely rather than refusing for a missing parent.
func TestItemNewStoryWithNoEpic(t *testing.T) {
	p := withDocs(t)
	rec := &recorder{ran: Ran{Stdout: []byte(`{"ok":true}`)}}
	params := `{"type":"story","title":"Standalone","body":"## Goal\nx\n",` + rid + `}`
	_, rerr := writeMethods(rec.run, time.Now, hostFor("item.new"))["item.new"](context.Background(), p, json.RawMessage(params))
	if rerr != nil {
		t.Fatalf("story with no epic: %+v", rerr)
	}
	if len(rec.runs) != 1 {
		t.Fatalf("ran %d commands", len(rec.runs))
	}
	args := strings.Join(rec.runs[0].Args, " ")
	if strings.Contains(args, "--epic") {
		t.Errorf("a story with no parent should not carry --epic: %s", args)
	}
	if !strings.HasPrefix(args, "story new --nature=feature --owner=olive --body-stdin") {
		t.Errorf("ran %s", args)
	}
}

// S-0198: issue.list without a story lists every issue, issue.story without
// an epic makes a story under none, and a story flai check refuses is
// Refused, as item.new's is.
func TestIssueMethodsWithTheirOptionsLeftOut(t *testing.T) {
	p := withDocs(t)
	for name, c := range map[string]struct{ params, args string }{
		"issue.list":  {`{}`, "issue list --json"},
		"issue.story": {`{"id":"I-12",` + rid + `}`, "issue story I-12 --owner=olive --autocommit --trailer=" + Trailer + " --json"},
	} {
		rec := &recorder{ran: Ran{Stdout: []byte(`[]`)}}
		if _, e := writeMethods(rec.run, time.Now, Host{})[name](context.Background(), p, json.RawMessage(c.params)); e != nil {
			t.Fatalf("%s: %+v", name, e)
		}
		if got := strings.Join(rec.runs[0].Args, " "); got != c.args {
			t.Errorf("%s ran %q, want %q", name, got, c.args)
		}
	}
	rec := &recorder{ran: Ran{Exit: 4, Stdout: []byte(`{"refused":{"findings":[{"rule":"item.heading"}]}}`), Events: []map[string]any{{"level": "FATAL", "err": "refused: flai check has 1 finding(s)"}}}}
	_, e := writeMethods(rec.run, time.Now, Host{})["issue.story"](context.Background(), p, json.RawMessage(good["issue.story"].params))
	if e == nil || e.Code != Refused || e.Data.(map[string]any)["findings"] == nil {
		t.Errorf("a refusal: %+v", e)
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
		if specs()[name].reads {
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
	// action governs publish.run, not accept.run
	ran, e := call("accept.run", Ran{Stdout: []byte(`{"id":"S-0001"}`)})
	if e != nil || len(ran) != 1 {
		t.Fatalf("an acceptance with the action off: %v %+v", ran, e)
	}
	if len(journal) != 0 {
		t.Errorf("acceptance is no host action: %+v", journal)
	}
	ran, e = call("publish.run", Ran{})
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
	ran, e = call("publish.run", Ran{Stdout: []byte(`{"pushed":true,"tags":["cli/v1.1.0"],"published":["origin v1.1.0 (abc1234)"]}`)})
	if e != nil || len(ran) != 1 || ran[0] != "release --pending --json" {
		t.Fatalf("publish.run: %v %+v", ran, e)
	}
	_, _ = call("publish.run", Ran{Stdout: []byte(`{"pushed":false,"reason":"nothing pending"}`)})
	_, _ = call("publish.run", Ran{Stdout: []byte(`{"pushed":false,"tags":["cli/v1.1.0"],"push_error":"rejected"}`)})
	_, e = call("publish.run", Ran{Exit: 3, Events: []map[string]any{{"level": "FATAL", "err": "conflict: main and origin/main have diverged"}}})
	if e == nil || e.Code != Conflict || e.Message != "main and origin/main have diverged" {
		t.Errorf("a remote that moved is a conflict with its reason: %+v", e)
	}
	// accept.run is never journalled as a host action (S-0087): it never
	// pushes, so it never uses the push action.
	want := []struct{ method, outcome, detail string }{
		{"publish.run", "done", "pushed with tags cli/v1.1.0; published origin v1.1.0 (abc1234)"},
		{"publish.run", "done", "nothing pushed: nothing pending"},
		{"publish.run", "failed", "not pushed: rejected"},
		{"publish.run", "failed", "main and origin/main have diverged"},
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

// Nothing a dashboard can ask for changes the host's settings unless the
// operator turned on the settings action in a shell (ADR-0029, S-0105), and
// nothing it can ask for turns that action on or off, or reads the journal
// or the rest of the configuration.
func TestOnlySettingsTouchesTheHostConfiguration(t *testing.T) {
	for name := range Methods("test", nil) {
		for _, word := range []string{"config", "journal", "enable", "disable"} {
			if strings.Contains(name, word) {
				t.Errorf("%s: the host's configuration and journal are the operator's, on the host", name)
			}
		}
	}
	for name, sp := range specs() {
		args, _, e := sp.build(channel.Project{Root: t.TempDir()}, json.RawMessage(good[name].params))
		if e != nil || len(args) == 0 {
			continue
		}
		if args[0] == "config" {
			t.Errorf("%s runs flai config", name)
		}
		host := args[0] == "serve" || args[0] == "agent" || (len(args) > 1 && args[1] == "token" && (args[0] == "dashboard" || args[0] == "mcp"))
		// S-0116, S-0115: starting a story's agent is the agent action's, and changes no setting; so is stopping it (S-0170)
		now := (name == "agent.restart" || name == "agent.start" || name == "agent.commit" || name == "agent.stop") && sp.action == ActionAgent && strings.Join(args[:3], " ") == "serve "+strings.Replace(name, ".", " ", 1)
		if host && !now && sp.action != ActionSettings {
			t.Errorf("%s runs flai %s without the settings action", name, strings.Join(args[:2], " "))
		}
		if strings.HasPrefix(name, "settings.") && sp.action != ActionSettings {
			t.Errorf("%s is not gated by the settings action", name)
		}
	}
	if info := enabledActions(Host{}, "/x"); len(info) != len(Actions)-len(shellOnly) || info[ActionPush] {
		t.Errorf("project.info names every action a dashboard sees and says none is on: %+v", info)
	}
}

// ADR-0067: publishing is the one way accepted work reaches the remote. No
// method pushes accepted work or says what a push would send; publish.run
// stays under the push action; and auto-publish, the operator's shell tool,
// is neither shown to a dashboard nor turned on or off by one.
func TestTheDashboardPublishesAndNeverPushes(t *testing.T) {
	p := withDocs(t)
	host := Host{Enabled: func(string, string) bool { return true }}
	table := MethodsFor("test", nil, host)
	for _, name := range []string{"push.run", "push.pending"} {
		if _, ok := table[name]; ok {
			t.Errorf("%s is still a method: the dashboard publishes, it does not push", name)
		}
	}
	if _, ok := table["publish.run"]; !ok || specs()["publish.run"].action != ActionPush {
		t.Errorf("publish.run is a method under the push action: %+v", specs()["publish.run"])
	}
	if info := enabledActions(host, p.Root); !info[ActionPush] {
		t.Errorf("project.info says the push action is on: %+v", info)
	} else if _, ok := info[ActionAutoPublish]; ok {
		t.Errorf("project.info names auto-publish: %+v", info)
	}
	if DashboardSees(ActionAutoPublish) || !DashboardSees(ActionPush) || DashboardSees("pull") {
		t.Error("a dashboard sees push, and neither auto-publish nor what is no host action")
	}
	rec := &recorder{ran: Ran{Stdout: []byte(`{}`)}}
	m := writeMethods(rec.run, time.Now, host)
	for i, on := range []string{"true", "false"} {
		params := fmt.Sprintf(`{"action":"auto-publish","on":%s,"request_id":"req-0000000%d"}`, on, i+2)
		_, e := m["settings.action"](context.Background(), p, json.RawMessage(params))
		if e == nil || e.Code != channel.CodeInvalidParams || !strings.Contains(e.Message, "flai serve enable auto-publish") || !strings.Contains(e.Message, "flai serve disable auto-publish") || !strings.Contains(e.Message, "ADR-0067") {
			t.Errorf("settings.action %s: %+v, want a refusal naming what to run in a shell", params, e)
		}
	}
	if len(rec.runs) != 0 {
		t.Errorf("a refused setting ran %+v", rec.runs)
	}
}

// S-0105: a setting kept for every project needs the settings action on for
// every project; one of this project's needs it here; either refusal says
// what the operator runs, and is journalled.
func TestSettingsNeedTheShellsConsent(t *testing.T) {
	p := withDocs(t)
	var journal []Entry
	here := Host{Enabled: func(a, root string) bool { return a == ActionSettings && root == p.Root }, Record: func(e Entry) { journal = append(journal, e) }}
	rec := &recorder{ran: Ran{Stdout: []byte(`{"ok":true}`)}}
	m := writeMethods(rec.run, time.Now, here)
	if _, e := m["settings.action"](context.Background(), p, json.RawMessage(good["settings.action"].params)); e != nil {
		t.Fatalf("a project's own setting, enabled here: %+v", e)
	}
	_, e := m["settings.agent"](context.Background(), p, json.RawMessage(good["settings.agent"].params))
	if e == nil || e.Code != Disabled || !strings.Contains(e.Message, "flai serve enable settings --all-projects") || e.Data.(map[string]any)["hostwide"] != true {
		t.Fatalf("a host-wide setting, enabled here only: %+v", e)
	}
	_, e = writeMethods(rec.run, time.Now, Host{})["settings.action"](context.Background(), p, json.RawMessage(`{"action":"agent","on":true,"request_id":"req-00000002"}`))
	if e == nil || e.Code != Disabled || !strings.Contains(e.Message, "flai serve enable settings") {
		t.Fatalf("not enabled: %+v", e)
	}
	if len(rec.runs) != 1 {
		t.Errorf("ran %d commands, want the one allowed", len(rec.runs))
	}
	if len(journal) != 2 || journal[0].Outcome != "done" || journal[0].Action != ActionSettings || journal[0].Detail != "flai serve enable push" || journal[1].Outcome != "disabled" {
		t.Errorf("journal: %+v", journal)
	}
}

// S-0122: serving or removing a project asks the settings action of that
// project, not of the one the request arrives through, and the journal names
// the project served or removed.
func TestServingAProjectNeedsItsOwnConsent(t *testing.T) {
	p := withDocs(t)
	const blog = "/home/me/git/blog"
	var journal []Entry
	host := Host{Enabled: func(a, root string) bool { return a == ActionSettings && root == blog }, Record: func(e Entry) { journal = append(journal, e) }}
	rec := &recorder{ran: Ran{Stdout: []byte(`{"removed":{"key":"blog"}}`)}}
	m := writeMethods(rec.run, time.Now, host)
	if _, e := m["settings.unserve"](context.Background(), p, json.RawMessage(`{"root":"`+blog+`","key":"blog",`+rid+`}`)); e != nil {
		t.Fatalf("enabled for the project removed, not for the one asking: %+v", e)
	}
	_, e := m["settings.serve"](context.Background(), p, json.RawMessage(`{"root":"/home/me/git/notes","key":"notes","request_id":"req-00000002"}`))
	if e == nil || e.Code != Disabled || !strings.Contains(e.Message, "not enabled for /home/me/git/notes") || !strings.Contains(e.Message, "flai serve enable settings") {
		t.Fatalf("enabled for the project asking only: %+v", e)
	}
	if len(rec.runs) != 1 || rec.runs[0].Dir != p.Root || strings.Join(rec.runs[0].Args, " ") != "serve project remove --json -- "+blog {
		t.Errorf("runs: %+v", rec.runs)
	}
	if len(journal) != 2 {
		t.Fatalf("journal: %+v", journal)
	}
	if j := journal[0]; j.Outcome != "done" || j.Action != ActionSettings || j.Project != "blog" || j.Root != blog || j.Detail != "flai serve project remove -- "+blog {
		t.Errorf("removed: %+v", j)
	}
	if j := journal[1]; j.Outcome != "disabled" || j.Project != "notes" || j.Root != "/home/me/git/notes" {
		t.Errorf("refused: %+v", j)
	}
}

// S-0105: the settings page reads what the host keeps, and whether it may
// change it here and everywhere, and never a token.
func TestSettingsGetSaysWhatMayChange(t *testing.T) {
	p := withDocs(t)
	host := Host{Enabled: func(a, root string) bool { return a == ActionSettings && root == p.Root },
		Settings: func(root string) any { return map[string]any{"root": root} }}
	res, e := MethodsFor("test", nil, host)["settings.get"](context.Background(), p, json.RawMessage(`{}`))
	if e != nil {
		t.Fatal(e)
	}
	got := res.(map[string]any)
	if got["here"] != true || got["everywhere"] != false || got["enable_everywhere"] != "flai serve enable settings --all-projects" || got["host"].(map[string]any)["root"] != p.Root {
		t.Errorf("settings.get: %+v", got)
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
	// Nothing the dashboard asks for configures an agent; the three that
	// start one, a ready story's now (S-0115), a new one for a story whose
	// agent dropped or failed (S-0116), and one to commit what a story in
	// review left uncommitted (S-0140), and the one that stops a story's
	// agent (S-0170), need the agent action.
	starts := map[string]bool{"agent.restart": true, "agent.start": true, "agent.commit": true, "agent.stop": true}
	for name := range Methods("test", nil) {
		if strings.HasPrefix(name, "agent.") && name != "agent.status" && name != "agent.stream" && !starts[name] {
			t.Errorf("%s: nothing the dashboard can ask for configures an agent", name)
		}
	}
	for name := range starts {
		if sp := specs()[name]; sp.action != ActionAgent {
			t.Errorf("%s is gated by %q, not the agent action", name, sp.action)
		}
	}
}

// S-0208: plan.run starts the planner for an epic or a story only while the
// operator has enabled the plan action for the project, says what enables it
// when not, and journals both, the start as what it started.
func TestPlanRunNeedsThePlanAction(t *testing.T) {
	p := withDocs(t)
	var journal []Entry
	on := false
	host := Host{Enabled: func(action, root string) bool { return on && action == ActionPlan && root == p.Root }, Record: func(e Entry) { journal = append(journal, e) }}
	rec := &recorder{ran: Ran{Stdout: []byte(`{"item":"E-0001","agent":"planner-E-0001","harness":"command","command":"plan-it","pid":42,"log":"/l","session":"s","started":"2026-10-03T10:00:00Z"}`)}}
	m := writeMethods(rec.run, time.Now, host)
	_, e := m["plan.run"](context.Background(), p, json.RawMessage(`{"id":"E-0001","request_id":"req-00000001"}`))
	if e == nil || e.Code != Disabled || !strings.Contains(e.Message, `the host action "plan" is not enabled`) || !strings.Contains(e.Message, "flai serve enable plan") {
		t.Fatalf("off: %+v", e)
	}
	if len(rec.runs) != 0 {
		t.Fatalf("a refused plan ran %+v", rec.runs)
	}
	on = true
	if _, e := m["plan.run"](context.Background(), p, json.RawMessage(`{"id":"E-0001","request_id":"req-00000002"}`)); e != nil {
		t.Fatalf("on: %+v", e)
	}
	if len(rec.runs) != 1 || rec.runs[0].Dir != p.Root || strings.Join(rec.runs[0].Args, " ") != "plan E-0001 --json" {
		t.Errorf("runs: %+v", rec.runs)
	}
	if len(journal) != 2 || journal[0].Outcome != "disabled" || journal[0].Action != ActionPlan || journal[0].Method != "plan.run" {
		t.Fatalf("journal: %+v", journal)
	}
	if j := journal[1]; j.Outcome != "done" || j.Action != ActionPlan || j.Detail != "started plan-it to plan E-0001 as planner-E-0001 (pid 42)" {
		t.Errorf("started: %+v", j)
	}
	if !DashboardSees(ActionPlan) {
		t.Error("the dashboard is not told of the plan action, so it cannot offer Plan")
	}
}

// S-0142: agent.stream reads a story's agent's stream through the host, from
// an offset or the tail, and refuses what is not a story or has no agent.
func TestAgentStreamAsksTheHost(t *testing.T) {
	p := withDocs(t)
	type asked struct {
		root, story string
		after       int64
	}
	var got []asked
	host := Host{AgentStream: func(root, story string, after int64) (any, error) {
		got = append(got, asked{root, story, after})
		if story == "S-0002" {
			return nil, fmt.Errorf("%s: %w", story, ErrNoAgent)
		}
		if story == "S-0003" {
			return nil, errors.New("disk on fire")
		}
		return map[string]any{"story": story, "next": 7}, nil
	}}
	m := MethodsFor("test", nil, host)["agent.stream"]
	res, e := m(context.Background(), p, json.RawMessage(`{"story":"S-0001","after":5}`))
	if e != nil {
		t.Fatal(e)
	}
	if b, _ := json.Marshal(res); string(b) != `{"next":7,"story":"S-0001"}` {
		t.Errorf("answer: %s", b)
	}
	if _, e := m(context.Background(), p, json.RawMessage(`{"story":"S-0001"}`)); e != nil {
		t.Fatal(e)
	}
	if want := []asked{{p.Root, "S-0001", 5}, {p.Root, "S-0001", -1}}; len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("asked %+v, want %+v", got, want)
	}
	for raw, code := range map[string]int{
		`{"story":"T-0001"}`:            channel.CodeInvalidParams,
		`{"story":"S-0001","after":-1}`: channel.CodeInvalidParams,
		`{"story":"S-0002"}`:            NotFound,
		`{"story":"S-0003"}`:            channel.CodeInternal,
	} {
		if _, e := m(context.Background(), p, json.RawMessage(raw)); e == nil || e.Code != code {
			t.Errorf("%s: %+v, want code %d", raw, e, code)
		}
	}
	if _, e := Methods("test", nil)["agent.stream"](context.Background(), p, json.RawMessage(`{"story":"S-0001"}`)); e == nil || e.Code != NotFound {
		t.Errorf("on a host that starts no agents: %+v", e)
	}
}

// S-0082: checks.run and checks.cancel are off until enabled, refused and
// journalled the same generic way every other host action is (already
// proven for push above); this checks describeChecksRun's own journal line.
func TestChecksHostActionJournalEntry(t *testing.T) {
	p := withDocs(t) // owner: olive
	var journal []Entry
	on := false
	host := Host{
		Enabled: func(action, root string) bool { return on && action == ActionChecks && root == p.Root },
		Record:  func(e Entry) { journal = append(journal, e) },
	}
	call := func(name string, ran Ran) *channel.Error {
		rec := &recorder{ran: ran}
		_, e := writeMethods(rec.run, time.Now, host)[name](context.Background(), p, json.RawMessage(good[name].params))
		return e
	}
	if e := call("checks.run", Ran{}); e == nil || e.Code != Disabled {
		t.Fatalf("off: %+v", e)
	}
	on = true
	if e := call("checks.run", Ran{Stdout: []byte(`{"story":"S-0001","outcome":"passed"}`)}); e != nil {
		t.Fatalf("on: %+v", e)
	}
	if e := call("checks.cancel", Ran{Stdout: []byte(`{"story":"S-0001","outcome":"cancelled"}`)}); e != nil {
		t.Fatalf("cancel: %+v", e)
	}
	if len(journal) != 3 {
		t.Fatalf("journal: %+v", journal)
	}
	if journal[0].Outcome != "disabled" || journal[0].Action != ActionChecks {
		t.Errorf("entry 0: %+v", journal[0])
	}
	if journal[1].Outcome != "done" || journal[1].Detail != "S-0001: passed" || journal[1].Action != ActionChecks {
		t.Errorf("entry 1: %+v", journal[1])
	}
	if journal[2].Outcome != "done" || journal[2].Detail != "S-0001: cancelled" {
		t.Errorf("entry 2: %+v", journal[2])
	}
}

// S-0081: dashboard.restart, dashboard.upgrade, and dashboard.stop can each
// end the very WebSocket connection their own request arrived on (a restart
// or a successful upgrade stops the container answering it; a stop does
// when it is the last project). Found live: the connection dying cancelled
// the request's context, which exec.CommandContext turned into a SIGINT to
// the flai subprocess mid-restart, killing it between "stop the old
// container" and "start the new one" and leaving nothing running at all.
// These specs, and checks.run since S-0082 (a browser closing the review
// page mid-run must not kill it either), run with their own background
// context instead.
func TestADetachedWriteSurvivesItsOwnConnectionDying(t *testing.T) {
	p := withDocs(t)
	for _, name := range []string{"dashboard.restart", "dashboard.upgrade", "dashboard.stop", "checks.run", "host.process", "host.upgrade"} {
		t.Run(name, func(t *testing.T) {
			parentCtx, cancelParent := context.WithCancel(context.Background())
			t.Cleanup(cancelParent)
			var sawCtx context.Context
			run := func(ctx context.Context, r Run) (Ran, error) {
				sawCtx = ctx
				cancelParent() // the request's own connection dies, as a real restart's does
				if ctx.Err() != nil {
					t.Fatal("the command's context was already cancelled by the connection dying mid-run")
				}
				return Ran{Stdout: []byte(`{"container":"flaiover"}`)}, nil
			}
			if _, e := writeMethods(run, time.Now, hostFor(name))[name](parentCtx, p, json.RawMessage(good[name].params)); e != nil {
				t.Fatalf("%s: %+v", name, e)
			}
			// The real assertion already happened inside run, live, the moment
			// the connection died: ctx.Err() was nil there, or the test failed
			// with t.Fatal above. sawCtx is checked here only for identity —
			// checking Err() this late would just see the method's own,
			// entirely normal, deferred cancel of its own context on return,
			// which happens whether or not the connection ever died.
			if sawCtx == nil {
				t.Fatal("run was never called")
			}
			if sawCtx == parentCtx {
				t.Errorf("%s: ran with the request's own context, not a detached one", name)
			}
		})
	}
}

// detachedWrites are the methods whose record outlives flai serve (S-0109).
func detachedWrites() []string {
	var out []string
	for name, sp := range specs() {
		if sp.detachTimeout > 0 && !sp.reads {
			out = append(out, name)
		}
	}
	return out
}

// withRecords is the host for a method, keeping its records where the test
// can read them, as a file beside flai serve's state would.
func withRecords(name string, kept map[string]Request) Host {
	h := hostFor(name)
	h.Requests = func(change func(map[string]Request)) error { change(kept); return nil }
	return h
}

// S-0109: a detached write is recorded before it acts, and its outcome when
// it ends.
func TestADetachedWriteIsRecordedBeforeItActs(t *testing.T) {
	p := withDocs(t)
	if len(detachedWrites()) < 6 {
		t.Fatalf("detached writes: %v", detachedWrites())
	}
	for _, name := range detachedWrites() {
		kept := map[string]Request{}
		key := p.Root + "|" + name + "|req-00000001"
		run := func(context.Context, Run) (Ran, error) {
			if r, ok := kept[key]; !ok || r.Done || r.PID != os.Getpid() {
				t.Errorf("%s acted before it was recorded as started: %+v", name, kept)
			}
			return Ran{Stdout: []byte(`{"n":1}`)}, nil
		}
		if _, e := writeMethods(run, time.Now, withRecords(name, kept))[name](context.Background(), p, json.RawMessage(good[name].params)); e != nil {
			t.Fatalf("%s: %+v", name, e)
		}
		if r := kept[key]; !r.Done || r.Result == nil || string(r.Result.Data) != `{"n":1}` {
			t.Errorf("%s: the outcome is not recorded: %+v", name, r)
		}
	}
}

// S-0109, found live in S-0107: a serve restart asked for from the dashboard
// ran twice, because the serve that took it ended before recording it and
// the next serve acted on the dashboard's repeat.
func TestARepeatAfterServeWentAwayIsAnsweredNotDone(t *testing.T) {
	p := withDocs(t)
	now := time.Date(2026, 9, 24, 7, 0, 0, 0, time.UTC)
	for _, name := range detachedWrites() {
		kept := map[string]Request{}
		ran := 0
		run := func(context.Context, Run) (Ran, error) { ran++; return Ran{Stdout: []byte(`{"n":1}`)}, nil }
		// The serve that took it recorded it as started, and ended.
		gone := os.Getpid() + 1
		kept[p.Root+"|"+name+"|req-00000001"] = Request{Started: now, PID: gone, Until: now.Add(time.Hour)}
		// The next serve is asked again.
		res, e := writeMethods(run, func() time.Time { return now.Add(time.Second) }, withRecords(name, kept))[name](context.Background(), p, json.RawMessage(good[name].params))
		if e != nil || ran != 0 {
			t.Fatalf("%s: ran %d times; %+v", name, ran, e)
		}
		w := res.(Written)
		if string(w.Data) != `{"repeat":{"request_id":"req-00000001","started":"2026-09-24T07:00:00Z","under_way":false}}` || len(w.Warnings) != 1 || !strings.Contains(w.Warnings[0], "not done again") {
			t.Errorf("%s answered %s %v", name, w.Data, w.Warnings)
		}
	}
}

// A repeat that arrives while the write is still running here is told it is
// under way, and a repeat after it ended gets its answer, from the next
// serve too.
func TestARepeatOfADetachedWriteUnderWayOrDone(t *testing.T) {
	p := withDocs(t)
	const name = "host.upgrade"
	kept := map[string]Request{}
	var mu sync.Mutex
	h := hostFor(name)
	h.Requests = func(change func(map[string]Request)) error { mu.Lock(); defer mu.Unlock(); change(kept); return nil }
	started, release := make(chan struct{}), make(chan struct{})
	ran := 0
	run := func(context.Context, Run) (Ran, error) {
		ran++
		close(started)
		<-release
		return Ran{Stdout: []byte(`{"restarting":true}`)}, nil
	}
	m := writeMethods(run, time.Now, h)[name]
	done := make(chan any)
	go func() {
		res, _ := m(context.Background(), p, json.RawMessage(good[name].params))
		done <- res
	}()
	<-started
	res, e := m(context.Background(), p, json.RawMessage(good[name].params))
	if w, ok := res.(Written); e != nil || !ok || !strings.Contains(string(w.Data), `"under_way":true`) {
		t.Errorf("while under way: %#v %+v", res, e)
	}
	close(release)
	first := <-done
	again, e := writeMethods(run, time.Now, h)[name](context.Background(), p, json.RawMessage(good[name].params))
	if e != nil || ran != 1 || string(again.(Written).Data) != string(first.(Written).Data) {
		t.Errorf("after it ended, from the next serve: ran %d, %#v %+v", ran, again, e)
	}
}

// A detached write whose record cannot be written is refused, not done
// unrecorded.
func TestADetachedWriteThatCannotBeRecordedIsNotDone(t *testing.T) {
	p := withDocs(t)
	for _, name := range detachedWrites() {
		h := hostFor(name)
		h.Requests = func(func(map[string]Request)) error { return errors.New("disk full") }
		rec := &recorder{}
		_, e := writeMethods(rec.run, time.Now, h)[name](context.Background(), p, json.RawMessage(good[name].params))
		if e == nil || e.Code != channel.CodeInternal || !strings.Contains(e.Message, "disk full") || len(rec.runs) != 0 {
			t.Errorf("%s: %+v, ran %d", name, e, len(rec.runs))
		}
	}
}

// A method with no detachTimeout keeps the old behavior: its exec context is
// the request's own, so a connection dying mid-run does cancel it. This
// guards against detaching everything by accident.
func TestAnOrdinaryWriteIsStillCancelledWithItsConnection(t *testing.T) {
	p := withDocs(t)
	parentCtx, cancelParent := context.WithCancel(context.Background())
	defer cancelParent()
	var sawCtx context.Context
	run := func(ctx context.Context, r Run) (Ran, error) {
		sawCtx = ctx
		cancelParent()
		return Ran{Stdout: []byte(`{"id":"S-0001"}`)}, nil
	}
	if _, e := writeMethods(run, time.Now, Host{})["item.move"](parentCtx, p, json.RawMessage(good["item.move"].params)); e != nil {
		t.Fatalf("item.move: %+v", e)
	}
	if sawCtx.Err() == nil {
		t.Error("an ordinary write's exec context should be the request's own and so get cancelled with it")
	}
}

// S-0103: the dashboard's agent reaches flai as flags, checked first, since
// each value becomes an argument.
func TestTheAgentReachesFlaiAsFlags(t *testing.T) {
	p := channel.Project{Key: "harbour", Root: "/p"}
	rec := &recorder{ran: Ran{Stdout: []byte(`{"ok":true}`)}}
	m := writeMethods(rec.run, time.Now, Host{})
	run := func(name, params string) ([]string, *channel.Error) {
		n := len(rec.runs)
		_, e := m[name](context.Background(), p, json.RawMessage(params))
		if len(rec.runs) == n {
			return nil, e
		}
		return rec.runs[len(rec.runs)-1].Args, e
	}
	args, e := run("item.new", `{"type":"story","title":"T","body":"b","request_id":"req-00000001","agent":{"harness":"claude-code","model":"claude-opus-5-5","config":{"effort":"high"}}}`)
	joined := strings.Join(args, " ")
	if e != nil || !strings.Contains(joined, "--harness=claude-code --model=claude-opus-5-5 --agent-config=effort=high") {
		t.Errorf("item.new: %v %s", e, joined)
	}
	if _, e := run("item.new", `{"type":"epic","title":"T","body":"b","request_id":"req-00000002","agent":{"model":"x"}}`); e == nil {
		t.Error("an epic was given an agent")
	}
	if _, e := run("item.new", `{"type":"story","title":"T","body":"b","request_id":"req-00000003","agent":{"harness":"; rm -rf /"}}`); e == nil {
		t.Error("an invalid harness reached a command line")
	}
	hash := strings.Repeat("a", 64)
	args, e = run("item.edit", `{"id":"S-0001","hash":"`+hash+`","request_id":"req-00000004","agent":{"model":"claude-sonnet-5"}}`)
	if e != nil || !strings.Contains(strings.Join(args, " "), "--clear-agent --model=claude-sonnet-5") {
		t.Errorf("item.edit replaces: %v %v", e, args)
	}
	args, e = run("item.edit", `{"id":"S-0001","hash":"`+hash+`","request_id":"req-00000005","agent":null}`)
	if e != nil || !strings.HasSuffix(strings.Join(args, " "), "--clear-agent --json") {
		t.Errorf("item.edit clears: %v %v", e, args)
	}
	args, _ = run("item.edit", `{"id":"S-0001","hash":"`+hash+`","request_id":"req-00000006","title":"New"}`)
	if strings.Contains(strings.Join(args, " "), "agent") {
		t.Errorf("an edit without agent touched it: %v", args)
	}
	// S-0189: roles reach flai as role flags, so a save from the dashboard
	// keeps them
	args, e = run("item.edit", `{"id":"S-0001","hash":"`+hash+`","request_id":"req-00000007","agent":{"model":"m","roles":{"verify":{"model":"sonnet","config":{"effort":"medium"}},"explore":{"harness":"claude-code"}}}}`)
	if e != nil || !strings.Contains(strings.Join(args, " "), "--clear-agent --model=m --role-harness=explore=claude-code --role-model=verify=sonnet --role-config=verify.effort=medium") {
		t.Errorf("item.edit with roles: %v %v", e, args)
	}
	for _, bad := range []string{`{"verify":{}}`, `{"Verify":{"model":"m"}}`, `{"verify":{"config":{"k":" "}}}`} {
		if _, e := run("item.new", `{"type":"story","title":"T","body":"b","request_id":"req-00000009","agent":{"roles":`+bad+`}}`); e == nil {
			t.Errorf("an invalid role reached a command line: %s", bad)
		}
	}
}

// S-0118: a restart that flai queued for want of room is journalled as queued,
// not as a start with no process.
func TestAnAgentRestartThatWasQueuedIsJournalledAsQueued(t *testing.T) {
	describe := describeAgentNow("restarted", "")
	if _, detail := describe(Written{Data: json.RawMessage(`{"story":"S-0001","agent":"agent-S-0001","queued":"2026-09-26T03:00:00Z"}`)}, nil); detail != "queued another agent for S-0001 until there is room" {
		t.Errorf("queued: %q", detail)
	}
	if _, detail := describe(Written{Data: json.RawMessage(`{"story":"S-0001","agent":"agent-S-0001","command":"claude","pid":42}`)}, nil); detail != "restarted claude for S-0001 as agent-S-0001 (pid 42)" {
		t.Errorf("started: %q", detail)
	}
	commit := describeAgentNow("started", " to commit what its worktree holds")
	if _, detail := commit(Written{Data: json.RawMessage(`{"story":"S-0001","agent":"agent-S-0001","command":"claude","pid":42}`)}, nil); detail != "started claude for S-0001 as agent-S-0001 to commit what its worktree holds (pid 42)" {
		t.Errorf("commit: %q", detail)
	}
}

// S-0199: the planning data reaches flai as flags: draft on a new story,
// each cost of delay and forecast key on an edit (an empty value removes
// it), the clears, and finalize on a move to ready as --yes.
func TestPlanningDataReachesFlaiAsFlags(t *testing.T) {
	p := channel.Project{Key: "harbour", Root: "/p"}
	rec := &recorder{ran: Ran{Stdout: []byte(`{"ok":true}`)}}
	m := writeMethods(rec.run, time.Now, Host{})
	n := 0
	run := func(name, params string) (string, *channel.Error) {
		n++
		before := len(rec.runs)
		_, e := m[name](context.Background(), p, json.RawMessage(strings.Replace(params, "{", fmt.Sprintf(`{"request_id":"req-%08d",`, n), 1)))
		if len(rec.runs) == before {
			return "", e
		}
		return strings.Join(rec.runs[len(rec.runs)-1].Args, " "), e
	}
	if args, e := run("item.new", `{"type":"story","title":"T","body":"b","draft":true}`); e != nil || !strings.Contains(args, " --draft --body-stdin") {
		t.Errorf("a new draft: %v %s", e, args)
	}
	if args, _ := run("item.new", `{"type":"story","title":"T","body":"b"}`); strings.Contains(args, "draft") {
		t.Errorf("a story not a draft: %s", args)
	}
	edit := `{"id":"S-0001","hash":"` + strings.Repeat("a", 64) + `",`
	for params, want := range map[string]string{
		`"draft":true}`:  " --draft --json",
		`"draft":false}`: " --no-draft --json",
		`"cost_of_delay":{"revenue_per_week":"1200.5","penalty_per_week":"300","time_lost_per_cycle":" 4h ","value":"2000"}}`: " --revenue-per-week=1200.5 --penalty-per-week=300 --time-lost-per-cycle=4h --cost-of-delay-value=2000 --json",
		`"cost_of_delay":{"penalty_per_week":""}}`: " --penalty-per-week= --json",
		`"clear_cost_of_delay":true}`:              " --clear-cost-of-delay --json",
		`"forecast":{"duration":"6h","delivery":"2026-10-09T17:00:00Z","basis":"three  tasks\n like S-0185"}}`: " --forecast-duration=6h --forecast-delivery=2026-10-09T17:00:00Z --forecast-basis=three tasks like S-0185 --json",
		`"forecast":{"basis":""}}`: " --forecast-basis= --json",
		`"clear_forecast":true}`:   " --clear-forecast --json",
		`"draft":true,"clear_forecast":true,"cost_of_delay":{"value":"5"}}`: " --draft --clear-forecast --cost-of-delay-value=5 --json",
	} {
		if args, e := run("item.edit", edit+params); e != nil || !strings.HasSuffix(args, "--trailer="+Trailer+want) {
			t.Errorf("item.edit %s: %v %s", params, e, args)
		}
	}
	for _, params := range []string{`"cost_of_delay":{}}`, `"forecast":{}}`} {
		if _, e := run("item.edit", edit+params); e == nil || e.Message != "nothing to change" {
			t.Errorf("an empty block %s: %+v", params, e)
		}
	}
	if args, e := run("item.move", `{"id":"S-0001","to":"ready","finalize":true}`); e != nil || args != "move S-0001 ready --by=designer --yes --json" {
		t.Errorf("finalize: %v %s", e, args)
	}
	if args, _ := run("item.move", `{"id":"S-0001","to":"ready"}`); strings.Contains(args, "--yes") {
		t.Errorf("a move without finalize: %s", args)
	}
	if args, _ := run("item.move", `{"id":"S-0001","to":"in-progress","finalize":true}`); strings.Contains(args, "--yes") {
		t.Errorf("finalize on a move not to ready: %s", args)
	}
}

// S-0204: item.new passes the cost of delay inputs as flags, trimmed, with
// an absent or empty key left out, beside the owner the inputs are set by;
// the value is the planner's and is refused.
func TestItemNewPassesCostOfDelayInputs(t *testing.T) {
	p := channel.Project{Key: "harbour", Root: "/p"}
	rec := &recorder{ran: Ran{Stdout: []byte(`{"ok":true}`)}}
	m := writeMethods(rec.run, time.Now, Host{})
	n := 0
	run := func(params string) (string, *channel.Error) {
		n++
		before := len(rec.runs)
		_, e := m["item.new"](context.Background(), p, json.RawMessage(strings.Replace(params, "{", fmt.Sprintf(`{"request_id":"req-%08d",`, n), 1)))
		if len(rec.runs) == before {
			return "", e
		}
		return strings.Join(rec.runs[len(rec.runs)-1].Args, " "), e
	}
	for params, want := range map[string]string{
		`{"type":"story","title":"T","body":"b","cost_of_delay":{"revenue_per_week":" 1200.5 ","penalty_per_week":"300","time_lost_per_cycle":"4h"}}`: "story new --nature=feature --owner=designer --revenue-per-week=1200.5 --penalty-per-week=300 --time-lost-per-cycle=4h --body-stdin",
		`{"type":"epic","title":"T","body":"b","cost_of_delay":{"revenue_per_week":"","penalty_per_week":"  ","time_lost_per_cycle":"2h"}}`:           "epic new --nature=feature --owner=designer --time-lost-per-cycle=2h --body-stdin",
		`{"type":"story","title":"T","body":"b","cost_of_delay":{}}`:                                                                                  "story new --nature=feature --owner=designer --body-stdin",
	} {
		if args, e := run(params); e != nil || !strings.HasPrefix(args, want) {
			t.Errorf("%s: %v %s", params, e, args)
		}
	}
	if args, e := run(`{"type":"story","title":"T","body":"b","cost_of_delay":{"revenue_per_week":"1","value":"9"}}`); e == nil || e.Code != channel.CodeInvalidParams || !strings.Contains(e.Message, "the value is the planner's") || args != "" {
		t.Errorf("a value on create: %+v %s", e, args)
	}
	if args, e := run(`{"type":"task","title":"T","body":"b","cost_of_delay":{"revenue_per_week":"1"}}`); e == nil || args != "" {
		t.Errorf("a task: %+v %s", e, args)
	}
}

// S-0201: item.finalize runs flai edit --no-draft as the owner, committed
// as the dashboard's other edits are. What the edit leaves unchanged was not
// a draft, and is refused as a rule in so many words.
func TestItemFinalizeFinalizesADraftAndRefusesWhatIsNot(t *testing.T) {
	p := channel.Project{Key: "harbour", Root: "/p"}
	finalize := func(ran Ran, params string) (any, *channel.Error, string) {
		rec := &recorder{ran: ran}
		res, e := writeMethods(rec.run, time.Now, Host{})["item.finalize"](context.Background(), p, json.RawMessage(params))
		if len(rec.runs) == 0 {
			return res, e, ""
		}
		return res, e, strings.Join(rec.runs[0].Args, " ")
	}
	res, e, args := finalize(Ran{Stdout: []byte(`{"id":"S-0007","changed":["draft"],"committed":true}`)}, `{"id":"S-0007",`+rid+`}`)
	if e != nil || args != "edit S-0007 --no-draft --by=designer --autocommit --trailer="+Trailer+" --json" {
		t.Errorf("finalize: %+v %s", e, args)
	}
	if w, ok := res.(Written); !ok || !strings.Contains(string(w.Data), `"changed":["draft"]`) {
		t.Errorf("the answer: %#v", res)
	}
	_, e, _ = finalize(Ran{Stdout: []byte(`{"id":"S-0007","changed":[],"unchanged":true}`)}, `{"id":"S-0007",`+rid+`}`)
	if e == nil || e.Code != Rule || !strings.HasPrefix(e.Message, "S-0007 is not a draft: there is nothing to finalize") {
		t.Errorf("a story that is not a draft: %+v", e)
	}
	if _, e, args := finalize(Ran{}, `{"id":"E-0001",`+rid+`}`); e == nil || e.Code != channel.CodeInvalidParams || e.Message != "E-0001 is not a story" || args != "" {
		t.Errorf("an epic: %+v %s", e, args)
	}
	// what flai check refuses reaches the dashboard as flai said it
	refusal := Ran{Exit: 4, Stdout: []byte(`{"refused":{"findings":[{"rule":"item.front-matter"}]}}`), Events: []map[string]any{{"level": "FATAL", "err": "refused: flai check has 1 finding(s)"}}}
	if _, e, _ := finalize(refusal, `{"id":"S-0007",`+rid+`}`); e == nil || e.Code != Refused {
		t.Errorf("a refusal of flai check: %+v", e)
	}
}

// S-0199: flai's refusal of a draft story to ready reaches the dashboard as
// a rule, in workitem's own words.
func TestADraftRefusedToReadyIsARule(t *testing.T) {
	p := harbour(t)
	repo, err := workitem.Open(p.Root)
	if err != nil {
		t.Fatal(err)
	}
	draft, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Drafted", Owner: "olive", Draft: true, Body: "## Acceptance criteria\n- [ ] works\n", Now: t0})
	if err != nil {
		t.Fatal(err)
	}
	_, refusal := repo.Transition(draft, workitem.Ready, "olive", "", t0)
	if refusal == nil || !strings.HasPrefix(refusal.Error(), "rule: ") {
		t.Fatalf("workitem should refuse a draft to ready as a rule: %v", refusal)
	}
	rec := &recorder{ran: Ran{Exit: 1, Events: []map[string]any{{"level": "FATAL", "err": refusal.Error()}}}}
	_, e := writeMethods(rec.run, time.Now, Host{})["item.move"](context.Background(), p, json.RawMessage(`{"id":"`+draft.ID+`","to":"ready",`+rid+`}`))
	if e == nil || e.Code != Rule || !strings.HasPrefix(e.Message, draft.ID+" is a draft: finalize it first") {
		t.Errorf("a draft to ready: %+v", e)
	}
}
