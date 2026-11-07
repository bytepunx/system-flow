package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
	"github.com/bytepunx/system-flow/flai/internal/prompt"
)

// fakeKlusterCLI is a kluster on PATH that logs its working directory and
// arguments, one call a line, to $FAKE_LOG; prints $FAKE_STATUS's file for
// status; writes a kubeconfig where kubeconfig --output names; and exits 3
// for the subcommand $FAKE_FAIL.
const fakeKlusterCLI = `#!/bin/sh
echo "$(pwd) $*" >> "$FAKE_LOG"
sub=""
out=""
prev=""
for a in "$@"; do
  case "$prev" in --config|--provider|--name) prev=""; continue ;; --output) out="$a"; prev=""; continue ;; esac
  case "$a" in --config|--provider|--name|--output) prev="$a" ;; --*) ;; *) [ -z "$sub" ] && sub="$a" ;; esac
done
if [ "$sub" = "$FAKE_FAIL" ]; then echo "kluster: $sub broke" >&2; exit 3; fi
case "$sub" in
  up) echo "creating cluster"; echo "cluster ready" ;;
  down) echo "deleted" ;;
  status) cat "$FAKE_STATUS" ;;
  kubeconfig) echo "apiVersion: v1" > "$out" ;;
esac
`

// klusterProject is a project whose system-flow.yaml has a kluster: block
// naming integration/kluster.yaml, with the fake kluster first on PATH and
// nothing else of the host's; it returns the root and the fake's log.
func klusterProject(t *testing.T, status string) (root, log string) {
	t.Helper()
	root = tempProject(t)
	m := filepath.Join(root, "system-flow.yaml")
	b, _ := os.ReadFile(m)
	b = append(b, []byte("kluster:\n  config: integration/kluster.yaml\n  deploy: [scripts/kluster-deploy.sh]\n")...)
	writeFiles(t, root, map[string]string{
		"system-flow.yaml":         string(b),
		".gitignore":               ".flai-cache/\n",
		"integration/kluster.yaml": "name: demo-integration\nprofile: authstar\nprovider: k3d\n",
	})
	bin := t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "kluster"), []byte(fakeKlusterCLI), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+"/usr/bin"+string(os.PathListSeparator)+"/bin")
	log = filepath.Join(t.TempDir(), "log")
	statusFile := filepath.Join(t.TempDir(), "status")
	if err := os.WriteFile(statusFile, []byte(status), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FAKE_LOG", log)
	t.Setenv("FAKE_STATUS", statusFile)
	t.Setenv("FAKE_FAIL", "")
	return root, log
}

// runKluster runs the CLI in dir as runIn does, with stdin a terminal or
// not, and answer given to any prompt.
func runKluster(t *testing.T, dir string, terminal bool, answer string, args ...string) (string, string, int) {
	t.Helper()
	var out, errOut bytes.Buffer
	a := &app{configPath: configOf(t), out: &out, errOut: &errOut, cwd: dir, clock: func() time.Time { return time.Date(2026, 11, 7, 1, 0, 0, 0, time.UTC) }}
	a.stdinIsTerminal = &terminal
	a.prompter = prompt.New(strings.NewReader(answer), &out)
	root := testRoot(t, a)
	root.SetArgs(args)
	code := 0
	if err := root.Execute(); err != nil {
		a.fail(err)
		code = 1
	}
	return out.String(), errOut.String(), code
}

func readLog(t *testing.T, log string) string {
	t.Helper()
	b, _ := os.ReadFile(log)
	return string(b)
}

// S-0577: flai kluster up runs kluster up with the block's config in the
// main checkout, streams its output, and keeps the kubeconfig under
// .flai-cache/kluster.
func TestKlusterUp(t *testing.T) {
	root, log := klusterProject(t, "")
	out, errOut, code := runKluster(t, root, false, "", "kluster", "up")
	if code != 0 {
		t.Fatalf("up: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(out, "creating cluster") || !strings.Contains(out, "demo-integration is up") || !strings.Contains(out, "k3d-demo-integration") {
		t.Errorf("up shows kluster's output and what it made: %s", out)
	}
	calls := readLog(t, log)
	if !strings.Contains(calls, root+" --config integration/kluster.yaml --provider k3d up") {
		t.Errorf("kluster up ran with the block's config in the root: %s", calls)
	}
	if _, err := os.Stat(kluster.KubeconfigPath(root)); err != nil {
		t.Errorf("the kubeconfig is kept: %v", err)
	}
}

// S-0577: flai kluster down runs kluster down with --yes, so kluster does
// not prompt; flai asks first unless --yes, and refuses on a non-terminal.
func TestKlusterDown(t *testing.T) {
	root, log := klusterProject(t, "")
	out, errOut, code := runKluster(t, root, false, "", "kluster", "down")
	if code == 0 || !strings.Contains(errOut, "pass --yes") || readLog(t, log) != "" {
		t.Fatalf("down without --yes on a non-terminal is refused, running nothing: %d %s %s", code, out, errOut)
	}
	out, errOut, code = runKluster(t, root, true, "n\n", "kluster", "down")
	if code == 0 || !strings.Contains(errOut, "left up") || readLog(t, log) != "" {
		t.Fatalf("down answered no runs nothing: %d %s %s", code, out, errOut)
	}
	out, errOut, code = runKluster(t, root, true, "y\n", "kluster", "down")
	if code != 0 || !strings.Contains(out, "demo-integration is down") {
		t.Fatalf("down answered yes: %d %s %s", code, out, errOut)
	}
	if calls := readLog(t, log); !strings.Contains(calls, "down --name demo-integration --yes") {
		t.Errorf("kluster down ran without its own prompt: %s", calls)
	}
	out, errOut, code = runKluster(t, root, false, "", "kluster", "down", "--yes")
	if code != 0 || !strings.Contains(out, "demo-integration is down") {
		t.Fatalf("down --yes: %d %s %s", code, out, errOut)
	}
}

// S-0577: flai kluster status reports the cluster's name, whether it runs,
// and its age, as text and with --json as one object; an absent cluster is
// a report, not a failure.
func TestKlusterStatus(t *testing.T) {
	root, _ := klusterProject(t, "NAME              RUNNING  AGE\nother             yes      9h\ndemo-integration  yes      2h50m39s\n")
	out, errOut, code := runKluster(t, root, false, "", "kluster", "status")
	if code != 0 {
		t.Fatalf("status: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{"demo-integration: running", "profile   authstar", "provider  k3d", "context   k3d-demo-integration", "age       2h50m39s"} {
		if !strings.Contains(out, want) {
			t.Errorf("status says %q: %s", want, out)
		}
	}
	out, errOut, code = runKluster(t, root, false, "", "kluster", "status", "--json")
	if code != 0 {
		t.Fatalf("status --json: %d %s %s", code, out, errOut)
	}
	var s map[string]any
	if err := json.Unmarshal([]byte(out), &s); err != nil {
		t.Fatalf("status --json is one object: %v %s", err, out)
	}
	if s["name"] != "demo-integration" || s["running"] != true || s["exists"] != true || s["age"] != "2h50m39s" || s["age_seconds"] != float64(2*3600+50*60+39) {
		t.Errorf("status --json: %v", s)
	}

	absent, _ := klusterProject(t, "NAME   RUNNING  AGE\nother  yes      9h\n")
	out, errOut, code = runKluster(t, absent, false, "", "kluster", "status")
	if code != 0 || !strings.Contains(out, "demo-integration: not created") {
		t.Errorf("an absent cluster is reported: %d %s %s", code, out, errOut)
	}
}

// S-0577: a project with no kluster: block, and a host with no kluster on
// PATH, are refused with what to add or how to install kluster.
func TestKlusterRefusals(t *testing.T) {
	bare := tempProject(t)
	for _, sub := range []string{"up", "down", "status"} {
		_, errOut, code := runKluster(t, bare, false, "", "kluster", sub, "--yes")
		if code == 0 || !strings.Contains(errOut, "has no kluster: block") || !strings.Contains(errOut, "config: integration/kluster.yaml") {
			t.Errorf("%s without a block: %d %s", sub, code, errOut)
		}
	}
	root, _ := klusterProject(t, "")
	t.Setenv("PATH", t.TempDir())
	for _, sub := range []string{"up", "down", "status"} {
		_, errOut, code := runKluster(t, root, false, "", "kluster", sub, "--yes")
		if code == 0 || !strings.Contains(errOut, "kluster is not installed") || !strings.Contains(errOut, "install.sh") || !strings.Contains(errOut, "go install github.com/bytepunx/kluster/kluster@latest") {
			t.Errorf("%s without kluster: %d %s", sub, code, errOut)
		}
	}
}

// S-0577: a failing kluster fails the command with kluster's own words.
func TestKlusterUpFails(t *testing.T) {
	root, _ := klusterProject(t, "")
	t.Setenv("FAKE_FAIL", "up")
	out, errOut, code := runKluster(t, root, false, "", "kluster", "up")
	if code == 0 || !strings.Contains(out, "kluster: up broke") || !strings.Contains(errOut, "kluster up of the cluster demo-integration") {
		t.Errorf("a failed up: %d %s %s", code, out, errOut)
	}
	if _, err := os.Stat(kluster.KubeconfigPath(root)); err == nil {
		t.Error("a failed up keeps no kubeconfig")
	}
}
