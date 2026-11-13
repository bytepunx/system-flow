package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
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
// PATH, are refused with how to get a block or how to install kluster.
func TestKlusterRefusals(t *testing.T) {
	bare := tempProject(t)
	for _, sub := range []string{"up", "down", "status"} {
		_, errOut, code := runKluster(t, bare, false, "", "kluster", sub, "--yes")
		if code == 0 || !strings.Contains(errOut, "has no kluster: block") || !strings.Contains(errOut, "flai serve enable kluster") {
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

// The fakes a deploy runs, each logging one line a call to $FAKE_LOG: docker
// lists one server node and a load balancer, saves an image to the file -o
// names, and reads an import from its standard input; kubectl only logs; git
// answers the head of a checkout with $FAKE_HEAD, the main branch with main,
// and makes a worktree with the project's scripts in it, and fails anything
// else as outside a repository.
const (
	fakeDocker = `#!/bin/sh
echo "docker $*" >> "$FAKE_LOG"
case "$1" in
  ps) echo "k3d-demo-integration-server-0 server"; echo "k3d-demo-integration-serverlb loadbalancer" ;;
  save) echo saved > "$3" ;;
  exec) cat > /dev/null ;;
esac
`
	fakeKubectl = `#!/bin/sh
echo "kubectl $*" >> "$FAKE_LOG"
case "$*" in
  *" create token "*) echo fake-token ;;
  *" apply "*) cat > /dev/null ;;
esac
`
	// klusterKubeconfig is flai's kubeconfig of demo-integration, with fake
	// credentials, so a deploy can make the namespace's own kubeconfig.
	klusterKubeconfig = `apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: RkFLRS1DQQ==
    server: https://0.0.0.0:6550
  name: k3d-demo-integration
contexts:
- context:
    cluster: k3d-demo-integration
    user: admin@k3d-demo-integration
  name: k3d-demo-integration
current-context: k3d-demo-integration
users:
- name: admin@k3d-demo-integration
  user:
    client-certificate-data: RkFLRS1DRVJU
    client-key-data: RkFLRS1LRVk=
`
	fakeGit = `#!/bin/sh
case "$*" in
  "rev-parse HEAD"|"rev-parse --verify main^{commit}") echo "git $(pwd) $*" >> "$FAKE_LOG"; echo "$FAKE_HEAD"; exit 0 ;;
  "rev-parse --abbrev-ref HEAD") echo main; exit 0 ;;
  "worktree prune") exit 0 ;;
  "worktree add --detach --force "*) echo "git $*" >> "$FAKE_LOG"; mkdir -p "$5"; echo "gitdir: x" > "$5/.git"; cp -R scripts "$5/"; exit 0 ;;
  "-C "*) echo "git $*" >> "$FAKE_LOG"; exit 0 ;;
esac
echo "fatal: not a git repository" >&2
exit 128
`
	klusterBuildScript = `#!/bin/sh
echo "build $(pwd) story=$FLAI_STORY ns=$FLAI_KLUSTER_NAMESPACE tag=$FLAI_KLUSTER_TAG" >> "$FAKE_LOG"
echo "demo/app:$FLAI_KLUSTER_TAG" > "$FLAI_KLUSTER_IMAGES"
`
	klusterDeployScript = `#!/bin/sh
echo "deploy $(pwd) story=$FLAI_STORY ns=$FLAI_KLUSTER_NAMESPACE prefix=$FLAI_KLUSTER_PREFIX tag=$FLAI_KLUSTER_TAG context=$FLAI_KLUSTER_CONTEXT" >> "$FAKE_LOG"
if [ -n "$FAKE_DEPLOY_EXIT" ]; then echo "deploy broke"; exit "$FAKE_DEPLOY_EXIT"; fi
echo "rolled out"
`
	klusterTeardownScript = `#!/bin/sh
echo "teardown $(pwd) story=$FLAI_STORY ns=$FLAI_KLUSTER_NAMESPACE" >> "$FAKE_LOG"
`
)

const klusterHead = "0123456789abcdef0123456789abcdef01234567"

// deployProject is a klusterProject whose block builds, deploys, and tears
// down with scripts, its default namespace demo, with a kubeconfig kept, the
// fakes on PATH, and the story S-0007 in status with a worktree holding the
// scripts. It returns the root and the fakes' log.
func deployProject(t *testing.T, status string) (root, log string) {
	t.Helper()
	root, log = klusterProject(t, "")
	b, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	b = append(b, []byte("  namespace: demo\n  build: [scripts/kluster-build.sh]\n  teardown: [scripts/kluster-teardown.sh]\n")...)
	scripts := map[string]string{
		"scripts/kluster-build.sh":    klusterBuildScript,
		"scripts/kluster-deploy.sh":   klusterDeployScript,
		"scripts/kluster-teardown.sh": klusterTeardownScript,
	}
	files := map[string]string{
		"system-flow.yaml":               string(b),
		".flai-cache/kluster/kubeconfig": klusterKubeconfig,
		"wip/kanban/stories/S-0007-demo.md": "---\nid: S-0007\ntitle: Demo\ntype: story\nnature: feature\nstatus: " + status +
			"\ncreated: 2026-11-07T00:00:00Z\nupdated: 2026-11-07T00:00:00Z\ntransitions:\n  - {to: " + status + ", at: 2026-11-07T00:00:00Z, by: alex}\n---\n\n# S-0007 Demo\n",
	}
	for p, s := range scripts {
		files[p] = s
		files[".flai-cache/worktrees/S-0007/"+p] = s
	}
	writeFiles(t, root, files)
	for p := range scripts {
		for _, f := range []string{filepath.Join(root, p), filepath.Join(root, ".flai-cache/worktrees/S-0007", p)} {
			if err := os.Chmod(f, 0o755); err != nil {
				t.Fatal(err)
			}
		}
	}
	bin := strings.SplitN(os.Getenv("PATH"), string(os.PathListSeparator), 2)[0]
	for name, s := range map[string]string{"docker": fakeDocker, "kubectl": fakeKubectl, "git": fakeGit} {
		if err := os.WriteFile(filepath.Join(bin, name), []byte(s), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	t.Setenv("FAKE_HEAD", klusterHead)
	t.Setenv("FAKE_DEPLOY_EXIT", "")
	return root, log
}

// S-0578: flai kluster deploy builds in the story's worktree, loads the
// image into the cluster's node, deploys into the story's namespace with its
// variables, and prints the record, which --json prints as data.
func TestKlusterDeploy(t *testing.T) {
	root, log := deployProject(t, "in-progress")
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0007")
	out, errOut, code := runKluster(t, root, false, "", "kluster", "deploy", "S-7")
	if code != 0 {
		t.Fatalf("deploy: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{"s-0007: build, log", "S-0007 in s-0007: passed", "commit  " + klusterHead, "tag     s-0007-0123456789ab", "images  demo/app:s-0007-0123456789ab", "deploy   passed, log " + kluster.LogPath(root, "s-0007", "deploy")} {
		if !strings.Contains(out, want) {
			t.Errorf("deploy prints %q: %s", want, out)
		}
	}
	calls := readLog(t, log)
	for _, want := range []string{
		"build " + wt + " story=S-0007 ns=s-0007 tag=s-0007-0123456789ab",
		"docker save -o",
		"docker exec -i k3d-demo-integration-server-0 ctr -n k8s.io images import -",
		"deploy " + wt + " story=S-0007 ns=s-0007 prefix=s0007_ tag=s-0007-0123456789ab context=k3d-demo-integration",
	} {
		if !strings.Contains(calls, want) {
			t.Errorf("the deploy ran %q:\n%s", want, calls)
		}
	}
	if strings.Contains(calls, "serverlb") {
		t.Errorf("the load balancer is no node to load into:\n%s", calls)
	}

	out, errOut, code = runKluster(t, root, false, "", "kluster", "deploy", "S-0007", "--json")
	if code != 0 {
		t.Fatalf("deploy --json: %d %s %s", code, out, errOut)
	}
	var r map[string]any
	if err := json.Unmarshal([]byte(out), &r); err != nil {
		t.Fatalf("deploy --json prints only the record: %v %s", err, out)
	}
	if r["story"] != "S-0007" || r["namespace"] != "s-0007" || r["commit"] != klusterHead || r["outcome"] != "passed" || len(r["steps"].([]any)) != 3 {
		t.Errorf("deploy --json: %v", r)
	}
	if !strings.Contains(errOut, "s-0007: deploy, log") {
		t.Errorf("with --json the step lines go to standard error: %s", errOut)
	}
}

// S-0578: a step that fails stops the deploy, exits 1, and names the step,
// its exit code, and its log; the record says so.
func TestKlusterDeployFails(t *testing.T) {
	root, _ := deployProject(t, "review")
	t.Setenv("FAKE_DEPLOY_EXIT", "4")
	out, errOut, code := runKluster(t, root, false, "", "kluster", "deploy", "S-0007")
	log := kluster.LogPath(root, "s-0007", "deploy")
	if code == 0 || !strings.Contains(out, "S-0007 in s-0007: failed") || !strings.Contains(out, "deploy   failed, exit code 4, log "+log) {
		t.Fatalf("a failed deploy prints its record: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "the deploy step for the namespace s-0007 failed with exit code 4") || !strings.Contains(errOut, log) {
		t.Errorf("the error names the step, its exit code, and its log: %s", errOut)
	}
	if b, _ := os.ReadFile(log); !strings.Contains(string(b), "deploy broke") {
		t.Errorf("the step's output is in its log: %q", b)
	}
	out, _, code = runKluster(t, root, false, "", "kluster", "deploy", "S-0007", "--json")
	var r map[string]any
	if code == 0 || json.Unmarshal([]byte(out), &r) != nil || r["outcome"] != "failed" || r["failed_step"] != "deploy" || r["exit_code"] != float64(4) {
		t.Errorf("a failed deploy --json prints the record and exits 1: %d %s", code, out)
	}
}

// S-0578: a story not in progress or in review, and a story with no
// worktree, are refused, and nothing runs.
func TestKlusterDeployRefusals(t *testing.T) {
	root, log := deployProject(t, "ready")
	_, errOut, code := runKluster(t, root, false, "", "kluster", "deploy", "S-0007")
	if code == 0 || !strings.Contains(errOut, "S-0007 is ready, and flai deploys only a story in progress or in review") {
		t.Errorf("a ready story is refused: %d %s", code, errOut)
	}
	root2, log2 := deployProject(t, "in-progress")
	if err := os.RemoveAll(filepath.Join(root2, ".flai-cache", "worktrees", "S-0007")); err != nil {
		t.Fatal(err)
	}
	_, errOut, code = runKluster(t, root2, false, "", "kluster", "deploy", "S-0007")
	if code == 0 || !strings.Contains(errOut, "S-0007 has no worktree") || !strings.Contains(errOut, "flai stream open S-0007") {
		t.Errorf("a story with no worktree is refused: %d %s", code, errOut)
	}
	if calls := readLog(t, log) + readLog(t, log2); strings.Contains(calls, "build ") || strings.Contains(calls, "deploy ") {
		t.Errorf("a refusal runs nothing:\n%s", calls)
	}
}

// S-0578: flai kluster teardown asks first, or needs --yes without a
// terminal; it runs the teardown command and deletes the namespace, and a
// story never deployed has nothing to tear down.
func TestKlusterTeardown(t *testing.T) {
	root, log := deployProject(t, "in-progress")
	out, errOut, code := runKluster(t, root, false, "", "kluster", "teardown", "S-0007")
	if code != 0 || !strings.Contains(out, "S-0007 has no deployment recorded; nothing to tear down") {
		t.Fatalf("a story never deployed: %d %s %s", code, out, errOut)
	}
	out, _, code = runKluster(t, root, false, "", "kluster", "teardown", "S-0007", "--json")
	var nd map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &nd) != nil || nd["deployed"] != false || nd["story"] != "S-0007" {
		t.Errorf("a story never deployed, --json: %d %s", code, out)
	}
	if _, _, code := runKluster(t, root, false, "", "kluster", "deploy", "S-0007"); code != 0 {
		t.Fatal("deploy")
	}
	before := readLog(t, log)
	_, errOut, code = runKluster(t, root, false, "", "kluster", "teardown", "S-0007")
	if code == 0 || !strings.Contains(errOut, "pass --yes") || readLog(t, log) != before {
		t.Fatalf("teardown without --yes on a non-terminal is refused, running nothing: %d %s", code, errOut)
	}
	out, errOut, code = runKluster(t, root, false, "", "kluster", "teardown", "S-0007", "--yes")
	if code != 0 || !strings.Contains(out, "S-0007's namespace s-0007 is torn down") {
		t.Fatalf("teardown --yes: %d %s %s", code, out, errOut)
	}
	calls := readLog(t, log)
	wt := filepath.Join(root, ".flai-cache", "worktrees", "S-0007")
	if !strings.Contains(calls, "teardown "+wt+" story=S-0007 ns=s-0007") || !strings.Contains(calls, "kubectl --context k3d-demo-integration delete namespace s-0007 --wait --ignore-not-found") {
		t.Errorf("teardown runs the command in the worktree, then deletes the namespace:\n%s", calls)
	}
	if _, err := os.Stat(kluster.RecordPath(root, "s-0007")); !os.IsNotExist(err) {
		t.Errorf("the record is removed: %v", err)
	}
}

// S-0578: flai kluster promote deploys main's head from a worktree of it
// into the default namespace, and flai kluster deployments lists it first.
func TestKlusterPromoteAndDeployments(t *testing.T) {
	root, log := deployProject(t, "in-progress")
	out, errOut, code := runKluster(t, root, false, "", "kluster", "deployments")
	if code != 0 || strings.TrimSpace(out) != "no deployments" {
		t.Fatalf("no deployments: %d %s %s", code, out, errOut)
	}
	out, _, code = runKluster(t, root, false, "", "kluster", "deployments", "--json")
	if code != 0 || strings.TrimSpace(out) != "[]" {
		t.Errorf("no deployments, --json, is an empty array: %d %s", code, out)
	}
	if _, errOut, code := runKluster(t, root, false, "", "kluster", "deploy", "S-0007"); code != 0 {
		t.Fatalf("deploy: %s", errOut)
	}
	out, errOut, code = runKluster(t, root, false, "", "kluster", "promote")
	if code != 0 || !strings.Contains(out, "demo: passed") || !strings.Contains(out, "tag     main-0123456789ab") {
		t.Fatalf("promote: %d %s %s", code, out, errOut)
	}
	main := kluster.MainWorktree(root)
	calls := readLog(t, log)
	if !strings.Contains(calls, "git worktree add --detach --force "+main+" "+klusterHead) || !strings.Contains(calls, "deploy "+main+" story= ns=demo prefix= tag=main-0123456789ab") {
		t.Errorf("promote deploys main's head from its worktree into the default namespace:\n%s", calls)
	}
	out, errOut, code = runKluster(t, root, false, "", "kluster", "deployments")
	if code != 0 {
		t.Fatalf("deployments: %d %s %s", code, out, errOut)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "NAMESPACE") || !strings.HasPrefix(lines[1], "demo ") || !strings.Contains(lines[1], "(main)") || !strings.HasPrefix(lines[2], "s-0007 ") || !strings.Contains(lines[2], "0123456789ab") || !strings.Contains(lines[2], "passed") {
		t.Errorf("deployments lists the default namespace first:\n%s", out)
	}
	out, _, code = runKluster(t, root, false, "", "kluster", "deployments", "--json")
	var list []map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &list) != nil || len(list) != 2 || list[0]["namespace"] != "demo" || list[0]["default"] != true || list[1]["story"] != "S-0007" {
		t.Errorf("deployments --json: %d %s", code, out)
	}
}

// S-0578: flai kluster --help lists the four commands S-0578 adds, and each
// refuses a project with no kluster: block.
func TestKlusterDeployCommandsListedAndRefused(t *testing.T) {
	root, _ := klusterProject(t, "")
	out, _, _ := runKluster(t, root, false, "", "kluster", "--help")
	for _, sub := range []string{"deploy", "teardown", "promote", "deployments"} {
		if !strings.Contains(out, "\n  "+sub+" ") {
			t.Errorf("--help lists %s: %s", sub, out)
		}
	}
	bare := tempProject(t)
	for _, args := range [][]string{{"deploy", "S-0001"}, {"teardown", "S-0001", "--yes"}, {"promote"}, {"deployments"}} {
		_, errOut, code := runKluster(t, bare, false, "", append([]string{"kluster"}, args...)...)
		if code == 0 || !strings.Contains(errOut, "has no kluster: block") {
			t.Errorf("%s without a block: %d %s", args[0], code, errOut)
		}
	}
}

// jobProject is a project whose session is the agent of S-0007, by its name
// alone, with --job's waits shortened; it returns the root.
func jobProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_STORY", "")
	t.Setenv("FLAI_AGENT", "agent-S-0007")
	was := klusterJobWait
	klusterJobWait = kluster.Wait{Take: 3 * time.Second, Answer: 3 * time.Second, Poll: 5 * time.Millisecond}
	t.Cleanup(func() { klusterJobWait = was })
	return tempProject(t)
}

// serveKluster stands in for flai serve: it takes story's request of kind in
// root's requests folder, writes jobs[0] as the story's job state, answers
// with what answer makes of the request, and writes each later job state
// 30ms after the one before. The request arrives on the channel, which is
// closed when none came within five seconds.
func serveKluster(t *testing.T, root, story, kind string, answer func(kluster.Request) kluster.Answer, jobs ...kluster.AskedJob) <-chan kluster.Request {
	t.Helper()
	got := make(chan kluster.Request, 1)
	go func() {
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			req, ok, err := kluster.Take(root, story, kind)
			if err != nil {
				t.Errorf("take the request: %v", err)
				break
			}
			if !ok {
				time.Sleep(2 * time.Millisecond)
				continue
			}
			got <- req
			if len(jobs) > 0 {
				if err := kluster.WriteAskedJob(root, jobs[0]); err != nil {
					t.Errorf("the job state: %v", err)
				}
			}
			if err := kluster.WriteAnswer(root, answer(req)); err != nil {
				t.Errorf("the answer: %v", err)
			}
			for _, j := range jobs[min(1, len(jobs)):] {
				time.Sleep(30 * time.Millisecond)
				if err := kluster.WriteAskedJob(root, j); err != nil {
					t.Errorf("the job state: %v", err)
				}
			}
			return
		}
		close(got)
	}()
	return got
}

// queuedAt is when the jobs of the --job tests were queued.
var queuedAt = time.Date(2026, 11, 13, 7, 33, 53, 0, time.UTC)

// askedJob is S-0007's job of kind, queued, running once started, and ended
// with outcome once outcome is not empty.
func askedJob(kind string, started bool, outcome string) kluster.AskedJob {
	j := kluster.AskedJob{ID: kluster.JobID("S-0007", kind, queuedAt), Kind: kind, Story: "S-0007", By: "agent-S-0007", QueuedAt: queuedAt}
	if started {
		j.StartedAt = queuedAt.Add(time.Second)
	}
	if outcome != "" {
		j.EndedAt, j.Outcome, j.Commit = queuedAt.Add(time.Minute), outcome, klusterHead
	}
	return j
}

// S-0660 criterion 1, ADR-0292 §1: flai kluster deploy --job runs nothing
// here: it asks flai serve as the session's agent and prints the job's id,
// its kind, and the jobs ahead, and what to wait for; teardown --job asks
// without --yes, and a job of the story found running is said to be so.
func TestKlusterJobAsksFlaiServeForTheOwnStory(t *testing.T) {
	root := jobProject(t)
	deploy := askedJob(kluster.KindDeploy, false, "")
	log := filepath.Join(kluster.Dir(root), "logs", "s-0007")
	got := serveKluster(t, root, "S-0007", kluster.KindDeploy, func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Job: deploy.ID, Kind: r.Kind, Story: r.Story, Ahead: 2, Started: true, Log: log}
	}, deploy)
	out, errOut, code := runKluster(t, root, false, "", "kluster", "deploy", "S-7", "--job")
	if code != 0 {
		t.Fatalf("deploy --job: %d %s %s", code, out, errOut)
	}
	for _, want := range []string{"kluster deploy job " + deploy.ID + " queued for S-0007, 2 jobs ahead of it", "Wait for its kluster event (wait_for_events), or read it later with flai kluster deployments", "its output goes to " + log} {
		if !strings.Contains(out, want) {
			t.Errorf("deploy --job prints %q: %s", want, out)
		}
	}
	if req := <-got; req.Story != "S-0007" || req.Kind != kluster.KindDeploy || req.By != "agent-S-0007" || req.At.IsZero() {
		t.Errorf("the request flai serve read: %+v", req)
	}

	teardown := askedJob(kluster.KindTeardown, true, "")
	serveKluster(t, root, "S-0007", kluster.KindTeardown, func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Job: teardown.ID, Kind: r.Kind, Story: r.Story}
	}, teardown)
	out, errOut, code = runKluster(t, root, false, "", "kluster", "teardown", "S-0007", "--job")
	if code != 0 || !strings.Contains(out, "kluster teardown job "+teardown.ID+" of S-0007 was already running, so nothing more was queued; no job ahead of it") {
		t.Errorf("teardown --job without --yes, found running: %d %s %s", code, out, errOut)
	}
}

// S-0660 criterion 3, ADR-0292 §1: a session with a story of its own is
// refused another story's job before any request is written; a shell with no
// session story asks for the story given; --wait goes with --job alone.
func TestKlusterJobRefusesAnotherStory(t *testing.T) {
	root := jobProject(t)
	t.Setenv("FLAI_STORY", "S-0007")
	for _, sub := range []string{"deploy", "teardown"} {
		_, errOut, code := runKluster(t, root, false, "", "kluster", sub, "S-0008", "--job")
		if code == 0 || !strings.Contains(errOut, "this session works S-0007") || !strings.Contains(errOut, "nothing was asked for S-0008") {
			t.Errorf("%s --job for another story: %d %s", sub, code, errOut)
		}
	}
	if _, err := os.Stat(kluster.RequestsDir(root)); !os.IsNotExist(err) {
		t.Errorf("a refusal wrote a request: %v", err)
	}
	_, errOut, code := runKluster(t, root, false, "", "kluster", "deploy", "S-0007", "--wait", "1m")
	if code == 0 || !strings.Contains(errOut, "--wait waits for the job flai serve runs, so it goes with --job") {
		t.Errorf("--wait without --job: %d %s", code, errOut)
	}

	t.Setenv("FLAI_STORY", "")
	t.Setenv("FLAI_AGENT", "alex")
	got := serveKluster(t, root, "S-0008", kluster.KindDeploy, func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Job: kluster.JobID(r.Story, r.Kind, queuedAt), Kind: r.Kind, Story: r.Story, Started: true}
	})
	if out, errOut, code := runKluster(t, root, false, "", "kluster", "deploy", "S-8", "--job"); code != 0 || !strings.Contains(out, "queued for S-0008") {
		t.Errorf("the operator's shell asks for the story given: %d %s %s", code, out, errOut)
	}
	if req := <-got; req.Story != "S-0008" || req.By != "alex" {
		t.Errorf("the request: %+v", req)
	}
}

// ADR-0292 §1: a refusal from flai serve prints its reason and exits 1, its
// answer with --json; a request no flai serve takes is withdrawn, saying so.
func TestKlusterJobRefusedOrNotTaken(t *testing.T) {
	root := jobProject(t)
	reason := "the kluster host action is off for this project; the operator turns it on with flai serve enable kluster"
	refuse := func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Kind: r.Kind, Story: r.Story, Reason: reason}
	}
	serveKluster(t, root, "S-0007", kluster.KindDeploy, refuse)
	out, errOut, code := runKluster(t, root, false, "", "kluster", "deploy", "S-0007", "--job")
	if code == 0 || !strings.Contains(errOut, reason) || out != "" {
		t.Errorf("a refusal: %d %q %s", code, out, errOut)
	}
	serveKluster(t, root, "S-0007", kluster.KindTeardown, refuse)
	out, errOut, code = runKluster(t, root, false, "", "kluster", "teardown", "S-0007", "--job", "--json")
	var res map[string]map[string]any
	if code == 0 || json.Unmarshal([]byte(out), &res) != nil || res["answer"]["reason"] != reason || res["answer"]["kind"] != kluster.KindTeardown {
		t.Errorf("a refusal --json prints the answer and exits 1: %d %s %s", code, out, errOut)
	}

	klusterJobWait = kluster.Wait{Take: 50 * time.Millisecond, Poll: 5 * time.Millisecond}
	_, errOut, code = runKluster(t, root, false, "", "kluster", "deploy", "S-0007", "--job")
	if code == 0 || !strings.Contains(errOut, "flai serve is not running") {
		t.Errorf("no flai serve: %d %s", code, errOut)
	}
	if _, err := os.Stat(kluster.RequestFile(root, "S-0007", kluster.KindDeploy, kluster.RequestSuffix)); !os.IsNotExist(err) {
		t.Errorf("the request was not withdrawn: %v", err)
	}
}

// S-0660 criterion 2, ADR-0292 §5: --wait waits for the job's end and
// prints its outcome, commit, steps with their exit codes and logs, and
// error, exiting 1 on a failure; a job that outlasts the wait is said to be
// still running, with how to read it later, and exits 0.
func TestKlusterJobWaits(t *testing.T) {
	root := jobProject(t)
	running := askedJob(kluster.KindDeploy, true, "")
	failed := askedJob(kluster.KindDeploy, true, kluster.OutcomeFailed)
	code4 := 4
	buildLog, deployLog := kluster.LogPath(root, "s-0007", "build"), kluster.LogPath(root, "s-0007", "deploy")
	failed.FailedStep, failed.ExitCode, failed.Error, failed.Log = "deploy", &code4, "the deploy step failed with exit code 4", deployLog
	failed.Steps = []kluster.Step{{Name: "build", Outcome: kluster.OutcomePassed, Log: buildLog}, {Name: "deploy", Outcome: kluster.OutcomeFailed, ExitCode: 4, Log: deployLog}}
	answer := func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Job: running.ID, Kind: r.Kind, Story: r.Story, Started: true}
	}
	serveKluster(t, root, "S-0007", kluster.KindDeploy, answer, running, failed)
	out, errOut, code := runKluster(t, root, false, "", "kluster", "deploy", "S-0007", "--job", "--wait", "5s")
	for _, want := range []string{"queued for S-0007", "S-0007's deploy job " + running.ID + ": failed", "commit  0123456789ab\n", "build    passed, log " + buildLog, "deploy   failed, exit code 4, log " + deployLog, "failed at deploy, exit code 4", "error   the deploy step failed with exit code 4"} {
		if !strings.Contains(out, want) {
			t.Errorf("--wait prints %q: %s", want, out)
		}
	}
	if code == 0 {
		t.Errorf("a failed job exits 1: %s %s", out, errOut)
	}

	serveKluster(t, root, "S-0007", kluster.KindDeploy, answer, running)
	out, errOut, code = runKluster(t, root, false, "", "kluster", "deploy", "S-0007", "--job", "--wait", "100ms")
	if code != 0 || !strings.Contains(out, "It is still running after 100ms: wait for its kluster event (wait_for_events), or read it later with flai kluster deployments") {
		t.Errorf("a job that outlasts the wait: %d %s %s", code, out, errOut)
	}
}

// S-0660 criterion 2: --json prints the answer and, with --wait, the job's
// state and whether it ended, as one object.
func TestKlusterJobJSON(t *testing.T) {
	root := jobProject(t)
	running := askedJob(kluster.KindTeardown, true, "")
	passed := askedJob(kluster.KindTeardown, true, kluster.OutcomePassed)
	serveKluster(t, root, "S-0007", kluster.KindTeardown, func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Job: running.ID, Kind: r.Kind, Story: r.Story, Ahead: 1, Started: true}
	}, running, passed)
	out, errOut, code := runKluster(t, root, false, "", "kluster", "teardown", "S-0007", "--job", "--wait", "5s", "--json")
	var res struct {
		Answer kluster.Answer    `json:"answer"`
		Job    *kluster.AskedJob `json:"job"`
		Ended  *bool             `json:"ended"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &res) != nil {
		t.Fatalf("--job --wait --json: %d %s %s", code, out, errOut)
	}
	if res.Answer.Job != running.ID || res.Answer.Ahead != 1 || !res.Answer.Started || res.Job == nil || res.Job.Outcome != kluster.OutcomePassed || res.Job.Commit != klusterHead || res.Ended == nil || !*res.Ended {
		t.Errorf("--json: %s", out)
	}

	serveKluster(t, root, "S-0007", kluster.KindTeardown, func(r kluster.Request) kluster.Answer {
		return kluster.Answer{Job: running.ID, Kind: r.Kind, Story: r.Story, Started: true}
	})
	out, _, code = runKluster(t, root, false, "", "kluster", "teardown", "S-0007", "--job", "--json")
	var bare map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &bare) != nil || bare["job"] != nil || bare["ended"] != nil || bare["answer"] == nil {
		t.Errorf("--json without --wait prints the answer alone: %d %s", code, out)
	}
}

// logStory writes text as S-0007's step log under root, and a record of its
// namespace listing the step as passed.
func logStory(t *testing.T, root, step, text string) {
	t.Helper()
	r := kluster.Record{Story: "S-0007", Namespace: "s-0007", Outcome: kluster.OutcomePassed, StartedAt: queuedAt, Steps: []kluster.Step{{Name: step, Outcome: kluster.OutcomePassed, Log: kluster.LogPath(root, "s-0007", step)}}}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	writeFiles(t, root, map[string]string{
		".flai-cache/kluster/logs/s-0007/" + step + ".log": text,
		".flai-cache/kluster/deployments/s-0007.json":      string(b),
	})
	if _, ok, err := kluster.ReadRecord(root, "s-0007"); !ok || err != nil {
		t.Fatalf("the record is not where kluster reads it: %v", err)
	}
}

// S-0763 criterion 4, ADR-0292 §2: flai kluster log prints a step's whole
// log, read from 0 in kluster.log's 64 KB parts, with no kluster: block and
// no flai serve; --json prints kluster.log's fields with the whole text.
func TestKlusterLogPrintsAStepsLog(t *testing.T) {
	root := jobProject(t)
	text := "a" + strings.Repeat("é", kluster.LogChunk) + "deployed\n"
	logStory(t, root, kluster.StepDeploy, text)
	out, errOut, code := runKluster(t, root, false, "", "kluster", "log", "S-7", "--step", "deploy")
	if code != 0 || out != text {
		t.Fatalf("log: %d, %d bytes of %d, %s", code, len(out), len(text), errOut)
	}
	out, errOut, code = runKluster(t, root, false, "", "kluster", "log", "S-0007", "--step", "deploy", "--json")
	var got kluster.LogPart
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("log --json: %d %s %s", code, out, errOut)
	}
	want := kluster.LogPart{Namespace: "s-0007", Step: "deploy", Log: ".flai-cache/kluster/logs/s-0007/deploy.log", Exists: true, Offset: 0, Next: int64(len(text)), Size: int64(len(text)), Text: text, Outcome: kluster.OutcomePassed}
	if got != want {
		t.Errorf("log --json: got %+v", got)
	}
	t.Setenv("FLAI_AGENT", "alex")
	writeFiles(t, root, map[string]string{".flai-cache/kluster/logs/s-0008/build-api.log": "built\n"})
	if out, errOut, code := runKluster(t, root, false, "", "kluster", "log", "S-0008", "--step", "build-api"); code != 0 || out != "built\n" {
		t.Errorf("the operator's shell reads the story given: %d %q %s", code, out, errOut)
	}
}

// S-0763 criterion 4, ADR-0292 §2: a session with a story of its own is
// refused another story's log; a step kluster.log refuses is refused; a log
// that is not there is said to be absent, exiting 1.
func TestKlusterLogRefusals(t *testing.T) {
	root := jobProject(t)
	logStory(t, root, kluster.StepDeploy, "deployed\n")
	writeFiles(t, root, map[string]string{".flai-cache/kluster/logs/s-0008/deploy.log": "theirs\n"})
	out, errOut, code := runKluster(t, root, false, "", "kluster", "log", "S-0008", "--step", "deploy")
	if code == 0 || out != "" || !strings.Contains(errOut, "this session works S-0007") || !strings.Contains(errOut, "nothing was read for S-0008") {
		t.Errorf("another story's log: %d %q %s", code, out, errOut)
	}
	for _, args := range [][]string{{"--step", "images"}, {"--step", "build-"}, {"--step", "../../kubeconfig"}, {}} {
		out, errOut, code := runKluster(t, root, false, "", append([]string{"kluster", "log", "S-0007"}, args...)...)
		if code == 0 || out != "" || !strings.Contains(errOut, "is not a step") {
			t.Errorf("%v: %d %q %s", args, code, out, errOut)
		}
	}
	if _, errOut, code := runKluster(t, root, false, "", "kluster", "log", "E-0001", "--step", "deploy"); code == 0 || !strings.Contains(errOut, "is not a story ID") {
		t.Errorf("an epic: %d %s", code, errOut)
	}
	out, errOut, code = runKluster(t, root, false, "", "kluster", "log", "S-0007", "--step", "teardown")
	if code == 0 || out != "" || !strings.Contains(errOut, "S-0007's namespace s-0007 has no teardown log at .flai-cache/kluster/logs/s-0007/teardown.log") {
		t.Errorf("an absent log: %d %q %s", code, out, errOut)
	}
	out, _, code = runKluster(t, root, false, "", "kluster", "log", "S-0007", "--step", "teardown", "--json")
	var got kluster.LogPart
	if code == 0 || json.Unmarshal([]byte(out), &got) != nil || got.Exists || got.Text != "" {
		t.Errorf("an absent log --json prints exists false and exits 1: %d %s", code, out)
	}
}

// execProject is a deployProject whose session is the agent of S-0007, with
// no kluster on PATH, a KUBECONFIG and build variables of its own, and a
// record of S-0007's namespace at klusterHead with outcome, with the
// namespace's kubeconfig written when kubeconfig is true. It returns the
// root and S-0007's worktree.
func execProject(t *testing.T, outcome string, kubeconfig bool) (root, wt string) {
	t.Helper()
	root, _ = deployProject(t, "in-progress")
	bin := strings.SplitN(os.Getenv("PATH"), string(os.PathListSeparator), 2)[0]
	if err := os.Remove(filepath.Join(bin, "kluster")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FLAI_STORY", "")
	t.Setenv("FLAI_AGENT", "agent-S-0007")
	t.Setenv("KUBECONFIG", "/elsewhere/kubeconfig")
	t.Setenv("FLAI_KLUSTER_IMAGES", "/elsewhere/images")
	t.Setenv("FLAI_KLUSTER_COMPONENT", "api")
	r := kluster.Record{Story: "S-0007", Namespace: "s-0007", Commit: klusterHead, Outcome: outcome, StartedAt: queuedAt}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	files := map[string]string{".flai-cache/kluster/deployments/s-0007.json": string(b)}
	if kubeconfig {
		files[".flai-cache/kluster/agents/s-0007.kubeconfig"] = "apiVersion: v1\n"
	}
	writeFiles(t, root, files)
	return root, filepath.Join(root, ".flai-cache", "worktrees", "S-0007")
}

// runKlusterExit runs the CLI in dir as runKluster does, exiting with an
// exitError's code as run does.
func runKlusterExit(t *testing.T, dir string, args ...string) (string, string, int) {
	t.Helper()
	var out, errOut bytes.Buffer
	a := &app{configPath: configOf(t), out: &out, errOut: &errOut, cwd: dir, clock: func() time.Time { return time.Date(2026, 11, 7, 1, 0, 0, 0, time.UTC) }}
	terminal := false
	a.stdinIsTerminal = &terminal
	root := testRoot(t, a)
	root.SetArgs(args)
	root.SetIn(strings.NewReader("from stdin\n"))
	code := 0
	if err := root.Execute(); err != nil {
		var ee *exitError
		if code = 1; errors.As(err, &ee) {
			code = ee.code
		}
		if ee == nil || ee.msg != "" {
			a.fail(err)
		}
	}
	return out.String(), errOut.String(), code
}

// S-0763 criterion 3, ADR-0292 §2: flai kluster exec runs the command in the
// story's worktree with ADR-0254 §4's variables at the record's commit, but
// the build's, and KUBECONFIG the namespace's own in place of the shell's,
// with no kluster on PATH; standard input reaches it, and its exit code is
// flai's.
func TestKlusterExecRunsInTheWorktreeWithTheNamespacesVariables(t *testing.T) {
	root, wt := execProject(t, kluster.OutcomePassed, true)
	out, errOut, code := runKlusterExit(t, root, "kluster", "exec", "S-7", "--", "sh", "-c", "pwd; env; cat")
	if code != 0 {
		t.Fatalf("exec: %d %s %s", code, out, errOut)
	}
	lines := strings.Split(out, "\n")
	if lines[0] != wt {
		t.Errorf("the command runs in the worktree %s: %q", wt, lines[0])
	}
	for _, want := range []string{
		"FLAI_STORY=S-0007",
		"FLAI_KLUSTER_NAME=demo-integration",
		"FLAI_KLUSTER_NAMESPACE=s-0007",
		"FLAI_KLUSTER_PREFIX=s0007_",
		"FLAI_KLUSTER_COMMIT=" + klusterHead,
		"FLAI_KLUSTER_TAG=s-0007-0123456789ab",
		"FLAI_KLUSTER_STATE=" + kluster.StatePath(root, "s-0007"),
		"FLAI_KLUSTER_CONTEXT=k3d-demo-integration",
		"KUBECONFIG=" + kluster.AgentKubeconfigPath(root, "S-0007"),
		"PWD=" + wt,
		"FAKE_HEAD=" + klusterHead,
		"from stdin",
	} {
		if !slices.Contains(lines, want) {
			t.Errorf("the command gets %q:\n%s", want, out)
		}
	}
	for _, not := range []string{"FLAI_KLUSTER_IMAGES=", "FLAI_KLUSTER_COMPONENT=", "/elsewhere", "KUBECONFIG=" + kluster.KubeconfigPath(root)} {
		if strings.Contains(out, not) {
			t.Errorf("the command does not get %q:\n%s", not, out)
		}
	}
	if strings.Contains(errOut, "has no kubeconfig") {
		t.Errorf("a namespace with its kubeconfig says nothing of it: %s", errOut)
	}
	out, errOut, code = runKlusterExit(t, root, "kluster", "exec", "S-0007", "--", "sh", "-c", "echo out; echo err >&2; exit 7")
	if code != 7 || out != "out\n" || !strings.Contains(errOut, "err\n") {
		t.Errorf("the command's exit code and output are flai's: %d %q %q", code, out, errOut)
	}
	t.Setenv("FLAI_AGENT", "alex")
	if out, errOut, code := runKlusterExit(t, root, "kluster", "exec", "S-0007", "--", "sh", "-c", "echo $FLAI_KLUSTER_NAMESPACE"); code != 0 || out != "s-0007\n" {
		t.Errorf("the operator's shell runs for the story given: %d %q %s", code, out, errOut)
	}
}

// S-0763 criterion 3: a namespace whose record failed still runs the
// command, for diagnosis, and a namespace with no kubeconfig yet is said to
// have none.
func TestKlusterExecRunsForAFailedRecord(t *testing.T) {
	root, _ := execProject(t, kluster.OutcomeFailed, false)
	out, errOut, code := runKlusterExit(t, root, "kluster", "exec", "S-0007", "--", "sh", "-c", "echo $KUBECONFIG")
	if code != 0 || out != kluster.AgentKubeconfigPath(root, "S-0007")+"\n" {
		t.Fatalf("a failed record runs: %d %q %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "S-0007's namespace s-0007 has no kubeconfig at "+kluster.AgentKubeconfigPath(root, "S-0007")) {
		t.Errorf("a missing kubeconfig is said: %s", errOut)
	}
}

// S-0763 criterion 3, ADR-0292 §2: no record, another story than the
// session's, no worktree, and no command are refused, and nothing runs.
func TestKlusterExecRefusals(t *testing.T) {
	root, wt := execProject(t, kluster.OutcomePassed, true)
	marker := filepath.Join(t.TempDir(), "ran")
	touch := []string{"--", "touch", marker}
	cases := []struct {
		name string
		args []string
		want []string
	}{
		{"another story", append([]string{"S-0008"}, touch...), []string{"this session works S-0007", "nothing was run for S-0008"}},
		{"no command", []string{"S-0007", "--"}, []string{"no command to run: give it after --"}},
		{"no dash", []string{"S-0007", "touch"}, []string{"no command to run: give it after --"}},
		{"no story", []string{"--", "touch", marker}, []string{"give the story's ID, then -- and the command"}},
		{"nothing", nil, []string{"give the story's ID, then -- and the command"}},
		{"two stories", append([]string{"S-0007", "S-0008"}, touch...), []string{"one story at a time"}},
		{"an epic", append([]string{"E-0001"}, touch...), []string{"is not a story ID"}},
	}
	for _, c := range cases {
		out, errOut, code := runKlusterExit(t, root, append([]string{"kluster", "exec"}, c.args...)...)
		if code == 0 || out != "" {
			t.Errorf("%s: %d %q %s", c.name, code, out, errOut)
		}
		for _, w := range c.want {
			if !strings.Contains(errOut, w) {
				t.Errorf("%s says %q: %s", c.name, w, errOut)
			}
		}
	}
	if err := os.RemoveAll(wt); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runKlusterExit(t, root, append([]string{"kluster", "exec", "S-0007"}, touch...)...); code == 0 || !strings.Contains(errOut, "S-0007 has no worktree at "+wt) || !strings.Contains(errOut, "flai stream open S-0007") {
		t.Errorf("no worktree: %d %s", code, errOut)
	}
	if err := os.Remove(kluster.RecordPath(root, "s-0007")); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runKlusterExit(t, root, append([]string{"kluster", "exec", "S-0007"}, touch...)...); code == 0 || !strings.Contains(errOut, "S-0007's namespace s-0007 has no record") || !strings.Contains(errOut, "flai kluster deploy S-0007 --job") {
		t.Errorf("no record: %d %s", code, errOut)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Errorf("a refusal runs nothing")
	}
}

// envNames are the names of flai kluster env's export lines, in order.
func envNames(t *testing.T, out string) []string {
	t.Helper()
	var names []string
	for _, l := range strings.Split(strings.TrimSuffix(out, "\n"), "\n") {
		rest, ok := strings.CutPrefix(l, "export ")
		name, _, eq := strings.Cut(rest, "=")
		if !ok || !eq {
			t.Fatalf("not an export line: %q in\n%s", l, out)
		}
		names = append(names, name)
	}
	return names
}

// S-0661 criterion 3, ADR-0292 §2: flai kluster env prints exec's variables,
// with KUBECONFIG the namespace's own, as export lines sorted by name that
// eval sets in the shell, or as one JSON object; it needs no worktree, and
// the operator's shell prints them for the story given.
func TestKlusterEnvPrintsTheNamespacesVariables(t *testing.T) {
	root, wt := execProject(t, kluster.OutcomePassed, true)
	if err := os.RemoveAll(wt); err != nil {
		t.Fatal(err)
	}
	out, errOut, code := runKlusterExit(t, root, "kluster", "env", "S-7")
	if code != 0 {
		t.Fatalf("env: %d %s %s", code, out, errOut)
	}
	agentKC := kluster.AgentKubeconfigPath(root, "S-0007")
	lines := strings.Split(out, "\n")
	for _, want := range []string{
		"export FLAI_STORY='S-0007'",
		"export FLAI_KLUSTER_NAME='demo-integration'",
		"export FLAI_KLUSTER_NAMESPACE='s-0007'",
		"export FLAI_KLUSTER_PREFIX='s0007_'",
		"export FLAI_KLUSTER_COMMIT='" + klusterHead + "'",
		"export FLAI_KLUSTER_TAG='s-0007-0123456789ab'",
		"export FLAI_KLUSTER_STATE='" + kluster.StatePath(root, "s-0007") + "'",
		"export FLAI_KLUSTER_CONTEXT='k3d-demo-integration'",
		"export KUBECONFIG='" + agentKC + "'",
	} {
		if !slices.Contains(lines, want) {
			t.Errorf("env prints %q:\n%s", want, out)
		}
	}
	names := envNames(t, out)
	if len(names) != 9 || !slices.IsSorted(names) {
		t.Errorf("env prints the nine variables sorted by name: %v", names)
	}
	for _, not := range []string{"FLAI_KLUSTER_IMAGES", "FLAI_KLUSTER_COMPONENT", "/elsewhere", "'" + kluster.KubeconfigPath(root) + "'"} {
		if strings.Contains(out, not) {
			t.Errorf("env does not print %q:\n%s", not, out)
		}
	}
	if strings.Contains(errOut, "has no kubeconfig") {
		t.Errorf("a namespace with its kubeconfig says nothing of it: %s", errOut)
	}
	sh := exec.Command("sh", "-c", `eval "$1"; echo "$KUBECONFIG"; echo "$FLAI_KLUSTER_NAMESPACE"`, "sh", out)
	if got, err := sh.Output(); err != nil || string(got) != agentKC+"\ns-0007\n" {
		t.Errorf("eval sets the variables in place of the shell's: %v %q", err, got)
	}

	out, errOut, code = runKlusterExit(t, root, "kluster", "env", "S-0007", "--json")
	var vars map[string]string
	if code != 0 || json.Unmarshal([]byte(out), &vars) != nil {
		t.Fatalf("env --json: %d %s %s", code, out, errOut)
	}
	if len(vars) != 9 || vars["KUBECONFIG"] != agentKC || vars["FLAI_KLUSTER_NAMESPACE"] != "s-0007" || vars["FLAI_KLUSTER_COMMIT"] != klusterHead {
		t.Errorf("env --json prints the nine variables by name: %v", vars)
	}

	t.Setenv("FLAI_AGENT", "alex")
	if out, errOut, code := runKlusterExit(t, root, "kluster", "env", "S-0007"); code != 0 || !strings.Contains(out, "export FLAI_KLUSTER_NAMESPACE='s-0007'\n") {
		t.Errorf("the operator's shell prints the story given: %d %q %s", code, out, errOut)
	}
}

// S-0661 criterion 3: a value is single-quoted so that sh, and zsh where the
// host has it, read it back as it is, a ' and what a shell expands in it
// included.
func TestKlusterEnvQuotesForTheShell(t *testing.T) {
	values := []string{"", "plain", "it's", "''", `a "b" $HOME ` + "`id`" + ` \n *`, "two\nlines"}
	shells := []string{"sh"}
	if _, err := exec.LookPath("zsh"); err == nil {
		shells = append(shells, "zsh")
	}
	for _, v := range values {
		for _, s := range shells {
			got, err := exec.Command(s, "-c", `eval "$1"; printf %s "$V"`, s, "export V="+shellQuote(v)).Output()
			if err != nil || string(got) != v {
				t.Errorf("%s reads %s back as %q: %v %q", s, shellQuote(v), v, err, got)
			}
		}
	}
	if got := shellQuote("it's"); got != `'it'\''s'` {
		t.Errorf("a ' is written as '\\'': %s", got)
	}
}

// S-0661 criterion 3: a namespace whose record failed before its kubeconfig
// was written prints its variables all the same, and says on standard error
// that the kubeconfig is not there yet.
func TestKlusterEnvPrintsForAFailedRecord(t *testing.T) {
	root, _ := execProject(t, kluster.OutcomeFailed, false)
	kc := kluster.AgentKubeconfigPath(root, "S-0007")
	out, errOut, code := runKlusterExit(t, root, "kluster", "env", "S-0007")
	if code != 0 || !strings.Contains(out, "export KUBECONFIG='"+kc+"'\n") {
		t.Fatalf("a failed record prints: %d %q %s", code, out, errOut)
	}
	if !strings.Contains(errOut, "S-0007's namespace s-0007 has no kubeconfig at "+kc) {
		t.Errorf("a missing kubeconfig is said: %s", errOut)
	}
}

// S-0661 criterion 3, ADR-0292 §2: another story than the session's, not a
// story, and a namespace with no record are refused, and nothing is printed.
func TestKlusterEnvRefusals(t *testing.T) {
	root, _ := execProject(t, kluster.OutcomePassed, true)
	check := func(name string, want []string, args ...string) {
		t.Helper()
		out, errOut, code := runKlusterExit(t, root, append([]string{"kluster", "env"}, args...)...)
		if code == 0 || out != "" {
			t.Errorf("%s: %d %q %s", name, code, out, errOut)
		}
		for _, w := range want {
			if !strings.Contains(errOut, w) {
				t.Errorf("%s says %q: %s", name, w, errOut)
			}
		}
	}
	check("another story", []string{"this session works S-0007", "nothing was run for S-0008"}, "S-0008")
	check("an epic", []string{"is not a story ID"}, "E-0001")
	check("two stories", []string{"accepts 1 arg"}, "S-0007", "S-0008")
	if err := os.Remove(kluster.RecordPath(root, "s-0007")); err != nil {
		t.Fatal(err)
	}
	check("no record", []string{"S-0007's namespace s-0007 has no record", "flai kluster deploy S-0007 --job"}, "S-0007")
}
