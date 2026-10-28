package pilot

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
)

// fakeClaude is a claude for pilot to start: it records its arguments, one
// after another each ended by a NUL, and the variables pilot gives it, but
// of the token only whether it is set; writes a line to its standard error;
// and then, by PILOT_TEST_MODE, prints the file out and exits with
// PILOT_TEST_EXIT, traps SIGTERM and exits 143 once it comes, or is killed by
// the signal pilot forwards. It writes ready once it waits for a signal.
const fakeClaude = `#!/bin/sh
d="$PILOT_TEST_DIR"
: > "$d/argv"
for a in "$@"; do printf '%s\000' "$a" >> "$d/argv"; done
{
  printf 'CLAUDE_CONFIG_DIR=%s\n' "$CLAUDE_CONFIG_DIR"
  printf 'FLAI_AGENT=%s\n' "$FLAI_AGENT"
  printf 'FLAI_SESSION=%s\n' "$FLAI_SESSION"
  printf 'FLAI_STORY=%s\n' "$FLAI_STORY"
  if [ -n "$FLAI_MCP_TOKEN" ]; then echo 'token=set'; else echo 'token=unset'; fi
} > "$d/env"
echo 'the harness on its standard error' >&2
case "$PILOT_TEST_MODE" in
trap)
  trap 'kill "$pid" 2>/dev/null; echo trapped > "$d/trapped"; exit 143' TERM
  sleep 30 > /dev/null 2>&1 &
  pid=$!
  echo ready > "$d/ready"
  wait "$pid"
  exit 1 ;;
die)
  echo ready > "$d/ready"
  exec sleep 30 > /dev/null ;;
esac
cat "$d/out"
exit "${PILOT_TEST_EXIT:-0}"
`

// printed is what the fake claude prints: a line, a line longer than any
// line buffer's default, and a last line without a newline.
var printed = "{\"type\":\"system\"}\n" + strings.Repeat("x", 200<<10) + "\n" + `{"type":"result"} without a newline`

// lockedBuffer is a bytes.Buffer that pilot and the harness's standard
// error copier may write to at once.
type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

// fake puts the fake claude first on PATH in a new folder, with printed as
// its output, and gives the folder.
func fake(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fake harness is a POSIX shell script")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "claude"), []byte(fakeClaude), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "out"), []byte(printed), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(filepath.ListSeparator)+os.Getenv("PATH"))
	return dir
}

// environ is pilot's environment for the fake claude in dir: this test's,
// with the MCP server's variables and extra.
func environ(dir string, extra ...string) []string {
	return append(os.Environ(), append([]string{
		"FLAI_MCP_URL=" + urlValue,
		"FLAI_MCP_TOKEN=" + tokenValue,
		"FLAI_MCP_HOST=" + hostValue,
		"PILOT_TEST_DIR=" + dir,
	}, extra...)...)
}

// result is what one run of pilot gave.
type result struct {
	code           int
	stdout, stderr string
}

// run runs Main on req, given as JSON, failing t when it has not returned
// within the time limit, or when it printed the token.
func run(t *testing.T, ctx context.Context, req any, env []string, signals <-chan os.Signal) result {
	t.Helper()
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	var stdout, stderr lockedBuffer
	codes := make(chan int, 1)
	go func() {
		codes <- Main(ctx, bytes.NewReader(data), &stdout, &stderr, env, signals)
	}()
	select {
	case code := <-codes:
		r := result{code: code, stdout: stdout.String(), stderr: stderr.String()}
		if strings.Contains(r.stdout, tokenValue) || strings.Contains(r.stderr, tokenValue) {
			t.Errorf("pilot printed the token's value; stderr: %s", r.stderr)
		}
		return r
	case <-time.After(30 * time.Second):
		t.Fatalf("pilot has not returned after 30s; stderr: %s", stderr.String())
		return result{}
	}
}

// pilotLines are stderr's lines that are pilot's own.
func pilotLines(stderr string) []string {
	var out []string
	for _, l := range strings.Split(strings.TrimSuffix(stderr, "\n"), "\n") {
		if strings.HasPrefix(l, "pilot: ") {
			out = append(out, l)
		}
	}
	return out
}

// argv is what the fake claude in dir was last started with.
func argv(t *testing.T, dir string) []string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, "argv"))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Split(strings.TrimSuffix(string(data), "\x00"), "\x00")
}

// follows is whether args has flag followed by value.
func follows(args []string, flag, value string) bool {
	i := slices.Index(args, flag)
	return i >= 0 && i+1 < len(args) && args[i+1] == value
}

// waitFor waits for path to exist, failing t, without stopping it, after
// the time limit: it runs beside the test, not in it.
func waitFor(t *testing.T, path string) {
	deadline := time.Now().Add(10 * time.Second)
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			t.Errorf("%s was not written within 10s", path)
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestRunRelaysOutputAndExitCode(t *testing.T) {
	for _, code := range []int{0, 7} {
		dir := fake(t)
		r := run(t, t.Context(), valid(t.TempDir()), environ(dir, "PILOT_TEST_EXIT="+strconv.Itoa(code)), nil)
		if r.code != code {
			t.Errorf("pilot exited %d, want the harness's %d; stderr: %s", r.code, code, r.stderr)
		}
		if r.stdout != printed {
			t.Errorf("pilot relayed %d bytes, want the %d the harness printed, unchanged", len(r.stdout), len(printed))
		}
		if !strings.Contains(r.stderr, "the harness on its standard error\n") {
			t.Errorf("the harness's standard error did not reach pilot's: %s", r.stderr)
		}
		lines := pilotLines(r.stderr)
		if len(lines) != 2 || !strings.Contains(lines[0], "starting claude-code") || !strings.HasSuffix(lines[1], "exited with code "+strconv.Itoa(code)) {
			t.Errorf("pilot's own lines are %q, want a start line and an end line with the code", lines)
		}
	}
}

func TestRunStartsNewSessionThenResumesIt(t *testing.T) {
	dir := fake(t)
	cache := t.TempDir()
	req := valid(cache)
	mcp := filepath.Join(cache, "mcp.json")

	first := run(t, t.Context(), req, environ(dir), nil)
	args := argv(t, dir)
	if first.code != 0 || !follows(args, "--session-id", req.Session) || slices.Contains(args, "--resume") {
		t.Errorf("the first run exited %d with --session-id %s and no --resume? argv %q", first.code, req.Session, args)
	}
	if !follows(args, "--mcp-config", mcp) || !slices.Contains(args, "--strict-mcp-config") {
		t.Errorf("argv does not name %s in --mcp-config with --strict-mcp-config: %q", mcp, args)
	}
	if !strings.Contains(first.stderr, "in a new session") {
		t.Errorf("the first run's start line does not say the session is new: %s", first.stderr)
	}
	if strings.Contains(first.stderr, req.Session) {
		t.Errorf("pilot's lines name the session: %s", first.stderr)
	}

	second := run(t, t.Context(), req, environ(dir), nil)
	args = argv(t, dir)
	if second.code != 0 || !follows(args, "--resume", req.Session) || slices.Contains(args, "--session-id") {
		t.Errorf("the second run exited %d with --resume %s and no --session-id? argv %q", second.code, req.Session, args)
	}
	if !strings.Contains(second.stderr, "resuming its session") {
		t.Errorf("the second run's start line does not say it resumes: %s", second.stderr)
	}
}

func TestRunWritesTheMCPFileWithoutTheToken(t *testing.T) {
	dir := fake(t)
	cache := t.TempDir()
	mcp := filepath.Join(cache, "mcp.json")
	// A file left from before, open to all, is rewritten and closed.
	if err := os.WriteFile(mcp, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if r := run(t, t.Context(), valid(cache), environ(dir), nil); r.code != 0 {
		t.Fatalf("pilot exited %d; stderr: %s", r.code, r.stderr)
	}
	info, err := os.Stat(mcp)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("the MCP file's mode is %v, want 0600", info.Mode().Perm())
	}
	data, err := os.ReadFile(mcp)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), tokenValue) || !strings.Contains(string(data), "${FLAI_MCP_TOKEN}") {
		t.Errorf("the MCP file holds the token's value, or does not name its variable: %s", data)
	}
}

func TestRunGivesTheHarnessItsEnvironment(t *testing.T) {
	dir := fake(t)
	cache := t.TempDir()
	req := valid(cache)
	if r := run(t, t.Context(), req, environ(dir), nil); r.code != 0 {
		t.Fatalf("pilot exited %d; stderr: %s", r.code, r.stderr)
	}
	data, err := os.ReadFile(filepath.Join(dir, "env"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"CLAUDE_CONFIG_DIR=" + filepath.Join(cache, "claude"),
		"FLAI_AGENT=" + req.Name,
		"FLAI_SESSION=" + req.Session,
		"FLAI_STORY=" + req.Story,
		"token=set",
	} {
		if !strings.Contains(string(data), want+"\n") {
			t.Errorf("the harness's environment lacks %s:\n%s", want, data)
		}
	}
}

func TestRunStartsAStrategicAgentAsItsDefinition(t *testing.T) {
	dir := fake(t)
	root := t.TempDir()
	def := filepath.Join(root, ".claude", "agents", "planner.md")
	if err := os.MkdirAll(filepath.Dir(def), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(def, []byte("---\nname: planner\ndescription: Plans one epic or one story.\n---\nPlan it.\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	req := valid(t.TempDir())
	req.Story, req.Root, req.Role, req.Item = "", root, conventions.RolePlan, "E-0029"
	if r := run(t, t.Context(), req, environ(dir), nil); r.code != 0 {
		t.Fatalf("pilot exited %d; stderr: %s", r.code, r.stderr)
	}
	if args := argv(t, dir); !follows(args, "--agent", "planner") {
		t.Errorf("the planner's argv has no --agent planner: %q", args)
	}
}

func TestRunForwardsASignalToTheHarness(t *testing.T) {
	dir := fake(t)
	signals := make(chan os.Signal, 1)
	go func() {
		waitFor(t, filepath.Join(dir, "ready"))
		signals <- syscall.SIGTERM
	}()
	r := run(t, t.Context(), valid(t.TempDir()), environ(dir, "PILOT_TEST_MODE=trap"), signals)
	if _, err := os.Stat(filepath.Join(dir, "trapped")); err != nil {
		t.Errorf("the harness did not get the SIGTERM sent to pilot: %v", err)
	}
	if r.code != 143 {
		t.Errorf("pilot exited %d, want the harness's 143; stderr: %s", r.code, r.stderr)
	}
}

func TestRunEndsTheHarnessWhenItsContextIsDone(t *testing.T) {
	dir := fake(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	go func() {
		waitFor(t, filepath.Join(dir, "ready"))
		cancel()
	}()
	r := run(t, ctx, valid(t.TempDir()), environ(dir, "PILOT_TEST_MODE=die"), nil)
	if want := 128 + int(syscall.SIGTERM); r.code != want {
		t.Errorf("pilot exited %d for a harness SIGTERM killed, want %d; stderr: %s", r.code, want, r.stderr)
	}
}

func TestMainRefusesInOneLine(t *testing.T) {
	cases := []struct {
		name  string
		req   func(cache string) any
		empty bool // PATH holds no claude
	}{
		{"empty object", func(string) any { return map[string]any{} }, false},
		{"unsafe session", func(cache string) any { r := valid(cache); r.Session = "../S-0507"; return r }, false},
		{"no program", func(cache string) any { return valid(cache) }, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := fake(t)
			if c.empty {
				t.Setenv("PATH", t.TempDir())
			}
			r := run(t, t.Context(), c.req(t.TempDir()), environ(dir), nil)
			if r.code != ExitUsage {
				t.Errorf("pilot exited %d, want %d", r.code, ExitUsage)
			}
			if lines := strings.Split(strings.TrimSuffix(r.stderr, "\n"), "\n"); len(lines) != 1 || !strings.HasPrefix(lines[0], "pilot: ") {
				t.Errorf("pilot's standard error is not one pilot: line: %q", r.stderr)
			}
			if r.stdout != "" {
				t.Errorf("pilot wrote %q on standard output", r.stdout)
			}
			if _, err := os.Stat(filepath.Join(dir, "argv")); err == nil {
				t.Error("the harness was started")
			}
		})
	}
}
