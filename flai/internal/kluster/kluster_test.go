package kluster

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// fakeKluster is a kluster on PATH that logs its working directory and
// arguments, one call a line, to $FAKE_LOG; prints $FAKE_STATUS's file for
// status; writes a kubeconfig where kubeconfig --output names; and exits
// $FAKE_EXIT for the subcommand $FAKE_FAIL.
const fakeKluster = `#!/bin/sh
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

// resolvedTempDir is t.TempDir with its symlinks resolved, for a folder a
// fake script runs in: $(pwd) prints the resolved path, which on macOS,
// where /var is a link to /private/var, is not the one t.TempDir names.
func resolvedTempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

type fixture struct {
	root, log, status string
	c                 Cluster
}

// setup makes a project with a kluster config and puts the fake kluster
// first on PATH, with nothing else on it.
func setup(t *testing.T, config string) fixture {
	t.Helper()
	root, bin := resolvedTempDir(t), t.TempDir()
	if err := os.WriteFile(filepath.Join(bin, "kluster"), []byte(fakeKluster), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "integration"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "integration", "kluster.yaml"), []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", bin+string(os.PathListSeparator)+"/usr/bin"+string(os.PathListSeparator)+"/bin")
	f := fixture{root: root, log: filepath.Join(t.TempDir(), "log"), status: filepath.Join(t.TempDir(), "status")}
	t.Setenv("FAKE_LOG", f.log)
	t.Setenv("FAKE_STATUS", f.status)
	t.Setenv("FAKE_FAIL", "")
	m := manifest.Manifest{Name: "authstar", Kluster: &manifest.Kluster{Config: "integration/kluster.yaml", Deploy: []string{"scripts/kluster-deploy.sh"}}}
	c, err := New(root, m, System{})
	if err != nil {
		t.Fatal(err)
	}
	f.c = c
	return f
}

const authstarConfig = "name: authstar-integration\nprofile: authstar\nprovider: k3d\ntrust-domain: authstar.cluster.local\n"

func (f fixture) calls(t *testing.T) []string {
	t.Helper()
	data, err := os.ReadFile(f.log)
	if err != nil {
		return nil
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n")
}

func TestNewReadsTheConfig(t *testing.T) {
	f := setup(t, "name: shop\n")
	if f.c.Name != "shop" || f.c.Provider != "k3d" || f.c.Profile != "" || f.c.Context() != "k3d-shop" || f.c.Config != "integration/kluster.yaml" {
		t.Errorf("the cluster is what its config names, k3d by default: %+v", f.c)
	}
	if _, err := New(f.root, manifest.Manifest{Name: "x"}, System{}); !errors.Is(err, ErrNoBlock) || !strings.Contains(err.Error(), "kluster:") || !strings.Contains(err.Error(), "turn the kluster host action on") || strings.Contains(err.Error(), "add one") {
		t.Errorf("no block says to turn the kluster host action on, not to edit the file: %v", err)
	}
	m := manifest.Manifest{Kluster: &manifest.Kluster{Config: "nowhere.yaml"}}
	if _, err := New(f.root, m, System{}); err == nil || !strings.Contains(err.Error(), "nowhere.yaml") {
		t.Errorf("a config that is not there is said: %v", err)
	}
}

func TestArgs(t *testing.T) {
	c := Cluster{Root: "/r", Config: "integration/kluster.yaml", Name: "a", Provider: "k3d"}
	for _, tc := range []struct {
		got  []string
		want string
	}{
		{c.UpArgs(), "--config integration/kluster.yaml --provider k3d up"},
		{c.DownArgs(), "--config integration/kluster.yaml --provider k3d down --name a --yes"},
		{c.StatusArgs(), "--config integration/kluster.yaml --provider k3d status"},
		{c.KubeconfigArgs(), "--config integration/kluster.yaml --provider k3d kubeconfig --name a --output " + filepath.Join("/r", ".flai-cache", "kluster", "kubeconfig")},
	} {
		if got := strings.Join(tc.got, " "); got != tc.want {
			t.Errorf("args %q, want %q", got, tc.want)
		}
	}
}

func TestUpStreamsAndWritesTheKubeconfig(t *testing.T) {
	f := setup(t, authstarConfig)
	var out bytes.Buffer
	if err := f.c.Up(context.Background(), &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "creating cluster") || !strings.Contains(out.String(), "cluster ready") {
		t.Errorf("kluster's output is streamed: %q", out.String())
	}
	calls := f.calls(t)
	want := []string{
		f.root + " --config integration/kluster.yaml --provider k3d up",
		f.root + " --config integration/kluster.yaml --provider k3d kubeconfig --name authstar-integration --output " + KubeconfigPath(f.root),
	}
	if !slices.Equal(calls, want) {
		t.Errorf("up runs kluster up then kubeconfig, in the root:\n got %q\nwant %q", calls, want)
	}
	if data, err := os.ReadFile(KubeconfigPath(f.root)); err != nil || !strings.Contains(string(data), "apiVersion") {
		t.Errorf("the kubeconfig is kept under .flai-cache/kluster: %v %q", err, data)
	}
}

func TestUpThatFailsWritesNoKubeconfig(t *testing.T) {
	f := setup(t, authstarConfig)
	t.Setenv("FAKE_FAIL", "up")
	var out bytes.Buffer
	err := f.c.Up(context.Background(), &out)
	if err == nil || !strings.Contains(err.Error(), "kluster up of the cluster authstar-integration") || !strings.Contains(out.String(), "up broke") {
		t.Errorf("a failed up says what failed, and kluster's own words reach out: %v %q", err, out.String())
	}
	if len(f.calls(t)) != 1 {
		t.Errorf("nothing runs after a failed up: %q", f.calls(t))
	}
	if _, err := os.Stat(KubeconfigPath(f.root)); err == nil {
		t.Error("a failed up writes no kubeconfig")
	}
}

func TestDownSkipsKlustersPrompt(t *testing.T) {
	f := setup(t, authstarConfig)
	var out bytes.Buffer
	if err := f.c.Down(context.Background(), &out); err != nil {
		t.Fatal(err)
	}
	if calls := f.calls(t); len(calls) != 1 || calls[0] != f.root+" --config integration/kluster.yaml --provider k3d down --name authstar-integration --yes" {
		t.Errorf("down runs kluster down --yes in the root: %q", calls)
	}
	t.Setenv("FAKE_FAIL", "down")
	if err := f.c.Down(context.Background(), &out); err == nil || !strings.Contains(err.Error(), "kluster down") {
		t.Errorf("a failed down is said: %v", err)
	}
}

func TestStatus(t *testing.T) {
	f := setup(t, authstarConfig)
	for _, tc := range []struct {
		name, table string
		want        Status
	}{
		{"running", "NAME                  RUNNING  AGE\nauthstar-integration  yes      2h50m39s\nmarketeer             yes      1582h7m49s\n",
			Status{Exists: true, Running: true, Age: "2h50m39s", AgeSeconds: 2*3600 + 50*60 + 39, Uptime: 2*time.Hour + 50*time.Minute + 39*time.Second}},
		{"stopped", "NAME                  RUNNING  AGE\nauthstar-integration  no       5m\n",
			Status{Exists: true, Running: false, Age: "5m", AgeSeconds: 300, Uptime: 5 * time.Minute}},
		{"absent", "NAME       RUNNING  AGE\nmarketeer  yes      1h\n", Status{}},
		{"no clusters at all", "No clusters found.\n", Status{}},
		{"a wider table", "NAME                  PROVIDER  RUNNING  NODES  AGE\nauthstar-integration  k3d       yes      1      12m\n",
			Status{Exists: true, Running: true, Age: "12m", AgeSeconds: 720, Uptime: 12 * time.Minute}},
		{"an age that is not a duration", "NAME                  RUNNING  AGE\nauthstar-integration  yes      3d\n",
			Status{Exists: true, Running: true, Age: "3d"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(f.status, []byte(tc.table), 0o644); err != nil {
				t.Fatal(err)
			}
			got, err := f.c.Status(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			tc.want.Name, tc.want.Profile, tc.want.Provider, tc.want.Context = "authstar-integration", "authstar", "k3d", "k3d-authstar-integration"
			if got != tc.want {
				t.Errorf("status\n got %+v\nwant %+v", got, tc.want)
			}
		})
	}
}

func TestStatusThatFails(t *testing.T) {
	f := setup(t, authstarConfig)
	t.Setenv("FAKE_FAIL", "status")
	if _, err := f.c.Status(context.Background()); err == nil || !strings.Contains(err.Error(), "kluster status") || !strings.Contains(err.Error(), "status broke") {
		t.Errorf("a failed status says so, with kluster's words: %v", err)
	}
	if err := os.WriteFile(f.status, []byte("NAME  AGE\nauthstar-integration  1h\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FAKE_FAIL", "")
	if _, err := f.c.Status(context.Background()); err == nil || !strings.Contains(err.Error(), "RUNNING") {
		t.Errorf("a table without RUNNING is an error: %v", err)
	}
}

func TestMissingBinary(t *testing.T) {
	f := setup(t, authstarConfig)
	t.Setenv("PATH", t.TempDir())
	var out bytes.Buffer
	for name, err := range map[string]error{
		"up":   f.c.Up(context.Background(), &out),
		"down": f.c.Down(context.Background(), &out),
	} {
		if !errors.Is(err, ErrNotInstalled) {
			t.Errorf("%s without kluster: %v", name, err)
		}
	}
	_, err := f.c.Status(context.Background())
	if !errors.Is(err, ErrNotInstalled) || !strings.Contains(err.Error(), "install.sh") || !strings.Contains(err.Error(), "go install github.com/bytepunx/kluster/kluster@latest") {
		t.Errorf("the refusal says how to install kluster: %v", err)
	}
}
