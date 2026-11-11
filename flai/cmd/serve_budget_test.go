package cmd

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// budgetStates are three budget states at budgetAt: claude-code meclod held
// until the reset the API reported, codex past its backoff with no probe, and
// anthropic past its until with a probe claimed a minute before and not yet
// spawned, which holds.
const budgetStates = `{
"claude-code meclod":{"provider":"claude-code meclod","limit":"usage_limit_reached","status":429,"since":"2026-10-28T08:00:00Z","until":"2026-10-28T09:00:00Z","reset":"2026-10-28T09:00:00Z","backoff":"5m0s"},
"codex":{"provider":"codex","since":"2026-10-28T07:55:00Z","until":"2026-10-28T08:00:00Z","backoff":"5m0s"},
"anthropic":{"provider":"anthropic","status":429,"since":"2026-10-28T07:50:00Z","until":"2026-10-28T08:00:00Z","backoff":"10m0s",
"probe":{"root":"/r","work":"S-0001","started":"2026-10-28T08:29:00Z","claimed":"2026-10-28T08:29:00Z"}}}`

var budgetAt = time.Date(2026, 10, 28, 8, 30, 0, 0, time.UTC)

// budgetFolder gives t a config of its own and returns its serve folder, with
// budget.json holding state; "" writes none.
func budgetFolder(t *testing.T, state string) serve.Dir {
	t.Helper()
	cfgPath := filepath.Join(t.TempDir(), "cfg.json")
	useConfig(t, cfgPath)
	dir := serve.DirFor(cfgPath)
	if err := os.MkdirAll(string(dir), 0o700); err != nil {
		t.Fatal(err)
	}
	if state != "" {
		if err := os.WriteFile(filepath.Join(string(dir), "budget.json"), []byte(state), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

// budgetKeys are the keys of dir's budget states, in order.
func budgetKeys(t *testing.T, dir serve.Dir) []string {
	t.Helper()
	states, err := dir.Budgets()
	if err != nil {
		t.Fatal(err)
	}
	keys := []string{}
	for _, b := range states {
		keys = append(keys, b.Provider)
	}
	return keys
}

// budgetJournal are the serve.budget entries of dir's journal.
func budgetJournal(t *testing.T, dir serve.Dir) []hostapi.Entry {
	t.Helper()
	f, err := os.Open(dir.Journal())
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	var out []hostapi.Entry
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var e hostapi.Entry
		if err := json.Unmarshal(sc.Bytes(), &e); err != nil {
			t.Fatalf("journal line %q: %v", sc.Text(), err)
		}
		if e.Method == "serve.budget" {
			out = append(out, e)
		}
	}
	return out
}

// S-0720: flai serve budget lists each state, whether it holds a start now,
// its limit, its probe, and the command that clears it, as text and JSON.
func TestServeBudgetLists(t *testing.T) {
	t.Parallel()
	dir := budgetFolder(t, budgetStates)
	cwd := t.TempDir()
	out, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget")
	if code != 0 {
		t.Fatalf("budget: %d %s", code, errOut)
	}
	for _, want := range []string{
		"anthropic: holds every start on it while its probe runs, its hold having ended at 2026-10-28T08:00:00Z\n  limit: HTTP 429, set 2026-10-28T07:50:00Z\n  probe: S-0001 in /r, started 2026-10-28T08:29:00Z\n  clear it with: flai serve budget clear anthropic\n",
		"claude-code meclod: holds every start on it until 2026-10-28T09:00:00Z, the reset the API reported\n  limit: usage_limit_reached (429), set 2026-10-28T08:00:00Z\n  clear it with: flai serve budget clear 'claude-code meclod'\n",
		"codex: holds no start now: its hold ended at 2026-10-28T08:00:00Z, and the next start on it goes as the probe\n  limit: the harness's usage limit, set 2026-10-28T07:55:00Z\n",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("budget lacks %q:\n%s", want, out)
		}
	}
	js, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget", "--json")
	var got []struct {
		Provider string             `json:"provider"`
		Held     bool               `json:"held"`
		Until    time.Time          `json:"until"`
		Probe    *serve.BudgetProbe `json:"probe"`
	}
	if err := json.Unmarshal([]byte(js), &got); code != 0 || err != nil {
		t.Fatalf("budget --json: %d %v %s\n%s", code, err, errOut, js)
	}
	if len(got) != 3 || got[0].Provider != "anthropic" || !got[0].Held || got[0].Probe == nil || got[0].Probe.Work != "S-0001" ||
		got[1].Provider != "claude-code meclod" || !got[1].Held || !got[1].Until.Equal(time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)) ||
		got[2].Provider != "codex" || got[2].Held {
		t.Errorf("budget --json: %+v", got)
	}

	if err := os.Remove(filepath.Join(string(dir), "budget.json")); err != nil {
		t.Fatal(err)
	}
	if out, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget"); code != 0 || !strings.HasPrefix(out, "no budget state is set") {
		t.Errorf("budget with none: %d %q %s", code, out, errOut)
	}
	if js, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget", "--json"); code != 0 || strings.TrimSpace(js) != "[]" {
		t.Errorf("budget --json with none: %d %q %s", code, js, errOut)
	}
}

// S-0720: flai serve budget clear <key> removes that state alone, journals
// it as serve.budget with why, and sends its notice through the channel.
func TestServeBudgetClearByKey(t *testing.T) {
	t.Parallel()
	dir := budgetFolder(t, budgetStates)
	cwd := t.TempDir()
	notices := filepath.Join(t.TempDir(), "notices.jsonl")
	if runtime.GOOS != "windows" {
		if _, errOut, code := runInAt(t, cwd, budgetAt, "serve", "notify", "set", "--command", "cat >> '"+notices+"'; echo >> '"+notices+"'"); code != 0 {
			t.Fatalf("notify set: %d %s", code, errOut)
		}
	}
	out, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget", "clear", "claude-code meclod")
	if want := "cleared the budget state of claude-code meclod, on usage_limit_reached (429) until 2026-10-28T09:00:00Z; starts on it go again\n"; code != 0 || out != want {
		t.Fatalf("clear:\n got %d %q %s\nwant %q", code, out, errOut, want)
	}
	if keys := budgetKeys(t, dir); strings.Join(keys, ",") != "anthropic,codex" {
		t.Errorf("clear by key left %v", keys)
	}
	entries := budgetJournal(t, dir)
	if len(entries) != 1 || entries[0].By == "" || entries[0].Outcome != "done" || entries[0].At != "2026-10-28T08:30:00Z" ||
		!strings.Contains(entries[0].Detail, "cleared by hand with flai serve budget clear") || !strings.Contains(entries[0].Detail, "claude-code meclod") {
		t.Errorf("journal: %+v", entries)
	}
	if runtime.GOOS != "windows" {
		data, err := os.ReadFile(notices)
		if err != nil {
			t.Fatalf("no notice sent: %v (%s)", err, errOut)
		}
		var n map[string]any
		if err := json.Unmarshal([]byte(strings.TrimSpace(string(data))), &n); err != nil {
			t.Fatalf("notice %q: %v", data, err)
		}
		if n["kind"] != "budget" || n["state"] != "cleared" || n["provider"] != "claude-code meclod" || n["why"] != "cleared by hand with flai serve budget clear" {
			t.Errorf("notice: %v", n)
		}
	}
	js, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget", "clear", "codex", "--json")
	var got struct {
		Removed []serve.Budget `json:"removed"`
		Unread  string         `json:"unread"`
	}
	if err := json.Unmarshal([]byte(js), &got); code != 0 || err != nil || len(got.Removed) != 1 || got.Removed[0].Provider != "codex" || got.Unread != "" {
		t.Errorf("clear --json: %d %v %s\n%s", code, err, errOut, js)
	}
	if keys := budgetKeys(t, dir); strings.Join(keys, ",") != "anthropic" {
		t.Errorf("clear --json by key left %v", keys)
	}
}

// S-0720: flai serve budget clear with no key removes every state and
// journals each; a key that names none is an error naming the keys there are.
func TestServeBudgetClearAllAndUnknown(t *testing.T) {
	t.Parallel()
	dir := budgetFolder(t, budgetStates)
	cwd := t.TempDir()
	out, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget", "clear", "claude-code")
	if code == 0 || out != "" || !strings.Contains(errOut, "there is no budget state of") || !strings.Contains(errOut, "flai serve budget lists them") {
		t.Errorf("clear of no state: %d %q %s", code, out, errOut)
	}
	want := `there is no budget state of "claude-code"; there are: "anthropic", "claude-code meclod", "codex" (flai serve budget lists them; quote a key of more than one word)`
	if err := noBudgetState(dir, "claude-code"); err == nil || err.Error() != want {
		t.Errorf("clear of no state:\n got %v\nwant %s", err, want)
	}
	if keys := budgetKeys(t, dir); len(keys) != 3 || len(budgetJournal(t, dir)) != 0 {
		t.Errorf("clear of no state changed something: %v %v", keys, budgetJournal(t, dir))
	}
	js, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget", "clear", "--json")
	var got struct {
		Removed []serve.Budget `json:"removed"`
	}
	if err := json.Unmarshal([]byte(js), &got); code != 0 || err != nil || len(got.Removed) != 3 {
		t.Fatalf("clear --json of every state: %d %v %s\n%s", code, err, errOut, js)
	}
	if keys := budgetKeys(t, dir); len(keys) != 0 {
		t.Errorf("clear of every state left %v", keys)
	}
	if entries := budgetJournal(t, dir); len(entries) != 3 {
		t.Errorf("clear of every state journalled %d entries, not 3: %+v", len(entries), entries)
	}
	if out, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget", "clear"); code != 0 || out != "no budget state is set; nothing to clear\n" {
		t.Errorf("clear with none set: %d %q %s", code, out, errOut)
	}
	if _, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget", "clear", "codex"); code == 0 || !strings.Contains(errOut, "none is set, so flai serve holds no start") {
		t.Errorf("clear by key with none set: %d %s", code, errOut)
	}
}

// S-0720: a budget.json that does not parse is replaced by a clear with no
// key, which says so and why, and refused by a clear of one key.
func TestServeBudgetClearUnparsed(t *testing.T) {
	t.Parallel()
	dir := budgetFolder(t, "not json")
	cwd := t.TempDir()
	if _, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget"); code == 0 || !strings.Contains(errOut, "flai serve budget clear") {
		t.Errorf("budget of a file that does not parse: %d %s", code, errOut)
	}
	if _, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget", "clear", "codex"); code == 0 || !strings.Contains(errOut, "does not parse") {
		t.Errorf("clear by key of a file that does not parse: %d %s", code, errOut)
	}
	out, errOut, code := runInAt(t, cwd, budgetAt, "serve", "budget", "clear")
	if code != 0 || !strings.HasPrefix(out, "replaced "+filepath.Join(string(dir), "budget.json")+", which did not parse (") || !strings.Contains(out, "with no budget state") {
		t.Errorf("clear of a file that does not parse: %d %q %s", code, out, errOut)
	}
	if keys := budgetKeys(t, dir); len(keys) != 0 {
		t.Errorf("replaced file holds %v", keys)
	}
}
