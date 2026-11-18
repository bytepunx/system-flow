package kluster

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
)

// fakeCA is the certificate authority file a cluster's kubeconfig names in
// copySetup; not a real certificate.
const fakeCA = "-----BEGIN CERTIFICATE-----\nFAKE-CA-FILE\n-----END CERTIFICATE-----\n"

// copySetup is teardownSetup with flai's kubeconfig naming its certificate
// authority as a file beside it, ca.crt, and a file in the namespace
// s-0578's state folder, as a deploy's commands would leave one.
func copySetup(t *testing.T) deployFixture {
	t.Helper()
	f := teardownSetup(t)
	kc := strings.Replace(clusterKubeconfig, "certificate-authority-data: RkFLRS1DQQ==", "certificate-authority: ca.crt", 1)
	if err := os.WriteFile(KubeconfigPath(f.root), []byte(kc), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(filepath.Dir(KubeconfigPath(f.root)), "ca.crt"), []byte(fakeCA), 0o600); err != nil {
		t.Fatal(err)
	}
	f.stateFile(t, "s-0578", "db/url", "postgres://s0578")
	return f
}

// stateFile writes data at rel in namespace ns's state folder.
func (f deployFixture) stateFile(t *testing.T, ns, rel, data string) {
	t.Helper()
	path := filepath.Join(StatePath(f.root, ns), rel)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

// copied lists the files under dir by their path relative to it, each with
// its content, and fails on a file or folder the owner's alone does not
// hold.
func copied(t *testing.T, dir string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fi, err := d.Info()
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		if d.IsDir() {
			if fi.Mode().Perm() != 0o700 {
				t.Errorf("the copy's folder %s is its owner's alone: %v", rel, fi.Mode().Perm())
			}
			return nil
		}
		if fi.Mode().Perm() != 0o600 {
			t.Errorf("the copy's file %s is its owner's alone: %v", rel, fi.Mode().Perm())
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		out[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatalf("read the copy %s: %v", dir, err)
	}
	return out
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

func TestAgentCopyDirAndItsFolder(t *testing.T) {
	if got, want := AgentCopyDir("/p", "s-0578"), "/p/.flai-cache/agents/S-0578/kluster"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	c := Folder(AgentCopyDir("/p", "S-0578"))
	for got, want := range map[string]string{
		c.Record("s-0578"):          "/p/.flai-cache/agents/S-0578/kluster/deployments/s-0578.json",
		c.LogDir("s-0578"):          "/p/.flai-cache/agents/S-0578/kluster/logs/s-0578",
		c.Log("s-0578", StepDeploy): "/p/.flai-cache/agents/S-0578/kluster/logs/s-0578/deploy.log",
		c.State("s-0578"):           "/p/.flai-cache/agents/S-0578/kluster/state/s-0578",
		c.AgentKubeconfig("s-0578"): "/p/.flai-cache/agents/S-0578/kluster/agents/s-0578.kubeconfig",
	} {
		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}
	main := Folder(Dir("/p"))
	if RecordPath("/p", "s-0578") != main.Record("s-0578") || LogPath("/p", "s-0578", StepBuild) != main.Log("s-0578", StepBuild) || StatePath("/p", "s-0578") != main.State("s-0578") || AgentKubeconfigPath("/p", "S-0578") != main.AgentKubeconfig("s-0578") {
		t.Error("the main checkout's paths are its Folder's")
	}
}

// S-0834 criterion 1 (ADR-0319 §1): a deploy of a story's namespace leaves
// in the story's agent cache folder its record, its step logs, its state
// folder, and its kubeconfig with the certificate authority inlined, and
// nothing of another namespace or of the cluster's kubeconfig.
func TestADeployWritesTheNamespacesCopyInTheStorysAgentCacheFolder(t *testing.T) {
	f := copySetup(t)
	f.deployed(t, "S-0578")
	f.stateFile(t, "s-0579", "db/url", "postgres://s0579")
	f.deployed(t, "S-0579")

	dir := AgentCopyDir(f.root, "S-0578")
	got := copied(t, dir)
	want := []string{"agents/s-0578.kubeconfig", "deployments/s-0578.json", "logs/s-0578/build.log", "logs/s-0578/deploy.log", "logs/s-0578/images", "logs/s-0578/load.log", "state/s-0578/db/url"}
	if !slices.Equal(sortedKeys(got), want) {
		t.Fatalf("the copy holds the namespace's record, logs, state, and kubeconfig alone:\n got %q\nwant %q", sortedKeys(got), want)
	}
	for _, rel := range []string{"deployments/s-0578.json", "logs/s-0578/build.log", "logs/s-0578/load.log", "logs/s-0578/deploy.log", "state/s-0578/db/url"} {
		orig, err := os.ReadFile(filepath.Join(Dir(f.root), rel))
		if err != nil {
			t.Fatal(err)
		}
		if got[rel] != string(orig) {
			t.Errorf("the copy's %s is the main checkout's:\n got %q\nwant %q", rel, got[rel], orig)
		}
	}
	r, ok, err := Folder(dir).ReadRecord("s-0578")
	if err != nil || !ok || r.Outcome != OutcomePassed || r.Commit != commitA || len(r.Steps) != 3 {
		t.Errorf("the copied record reads as the deploy's last write, passed: %+v %v %v", r, ok, err)
	}
	if !strings.Contains(got["logs/s-0578/deploy.log"], "wrote the namespace s-0578's kubeconfig") {
		t.Errorf("the deploy step's log is copied whole, after it ended: %q", got["logs/s-0578/deploy.log"])
	}

	var kc kubeconfig
	if err := yaml.Unmarshal([]byte(got["agents/s-0578.kubeconfig"]), &kc); err != nil {
		t.Fatal(err)
	}
	if len(kc.Clusters) != 1 || kc.Clusters[0].Cluster != (kubeCluster{Server: "https://0.0.0.0:6550", CertificateAuthorityData: base64.StdEncoding.EncodeToString([]byte(fakeCA))}) {
		t.Errorf("the copied kubeconfig carries the certificate authority itself and names no file: %+v", kc.Clusters)
	}
	if len(kc.Users) != 1 || kc.Users[0].User.Token != fakeToken || len(kc.Contexts) != 1 || kc.Contexts[0].Context.Namespace != "s-0578" {
		t.Errorf("the copied kubeconfig is the namespace's, with its token: %+v", kc)
	}
	orig := f.readAgentKubeconfig(t, "S-0578")
	if ca := orig.Clusters[0].Cluster.CertificateAuthority; ca != filepath.Join(filepath.Dir(KubeconfigPath(f.root)), "ca.crt") {
		t.Errorf("the main checkout's kubeconfig still names the file: %q", ca)
	}
	for rel, data := range got {
		if strings.Contains(data, "client-key-data") || strings.Contains(data, "s-0579") || strings.Contains(data, "s0579") {
			t.Errorf("the copy's %s holds nothing of the cluster's kubeconfig or of another namespace:\n%s", rel, data)
		}
	}
	if _, err := os.Stat(filepath.Join(AgentCopyDir(f.root, "S-0579"), "state", "s-0579", "db", "url")); err != nil {
		t.Errorf("the other story's namespace has its own copy: %v", err)
	}
}

// ADR-0319 §1: the copy follows each step's end, so a step's log is there
// once the step ends, while the job still runs.
func TestTheCopyFollowsEachStepsEnd(t *testing.T) {
	f := copySetup(t)
	t.Setenv("FAKE_RECORD", filepath.Join(AgentCopyDir(f.root, "S-0578"), "deployments", "s-0578.json"))
	f.deployed(t, "S-0578")
	data, err := os.ReadFile(filepath.Join(f.env, "deploy.record"))
	if err != nil {
		t.Fatalf("the copied record is there while the deploy step runs: %v", err)
	}
	var r Record
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatal(err)
	}
	if names := stepNames(r.Steps); !slices.Equal(names, []string{"build=passed", "load=passed", "deploy=running"}) {
		t.Errorf("the copy the deploy step sees has the steps before it ended and it running: %q", names)
	}
	if _, err := os.Stat(Folder(AgentCopyDir(f.root, "S-0578")).Log("s-0578", StepLoad)); err != nil {
		t.Errorf("the load step's log was copied when it ended: %v", err)
	}
}

// ADR-0319 §1: what the main checkout no longer holds leaves the copy, so a
// redeploy that cleared the logs, or a command that removed a state file, is
// reflected.
func TestARedeployRemovesFromTheCopyWhatTheMainCheckoutNoLongerHolds(t *testing.T) {
	f := copySetup(t)
	f.deployed(t, "S-0578")
	if err := os.RemoveAll(filepath.Join(StatePath(f.root, "s-0578"), "db")); err != nil {
		t.Fatal(err)
	}
	f.stateFile(t, "s-0578", "token", "t2")
	f.m.Kluster.Build = nil
	f.deployed(t, "S-0578")
	got := sortedKeys(copied(t, AgentCopyDir(f.root, "S-0578")))
	if want := []string{"agents/s-0578.kubeconfig", "deployments/s-0578.json", "logs/s-0578/deploy.log", "state/s-0578/token"}; !slices.Equal(got, want) {
		t.Errorf("the copy mirrors the main checkout after the redeploy:\n got %q\nwant %q", got, want)
	}
	if _, err := os.Stat(filepath.Join(Folder(AgentCopyDir(f.root, "S-0578")).State("s-0578"), "db")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a state folder the main checkout no longer holds leaves the copy: %v", err)
	}
}

// ADR-0319 §1: a checks run kept in the record is kept in the copy's. Its
// log lies under flai serve's folder, not the namespace's logs, so the copy
// names it and does not hold it.
func TestAChecksRunKeptInTheRecordIsKeptInTheCopy(t *testing.T) {
	f := copySetup(t)
	f.deployed(t, "S-0578")
	checksLog := filepath.Join(t.TempDir(), "S-0578-checks.log")
	if err := RecordChecks(f.root, "s-0578", commitA, OutcomeFailed, checksLog, "agent-S-0578"); err != nil {
		t.Fatal(err)
	}
	r, ok, err := Folder(AgentCopyDir(f.root, "S-0578")).ReadRecord("s-0578")
	if err != nil || !ok || r.ChecksCommit != commitA || r.ChecksOutcome != OutcomeFailed || r.ChecksLog != checksLog || r.ChecksAskedBy != "agent-S-0578" || r.Outcome != OutcomePassed {
		t.Errorf("the copied record keeps the checks' outcome: %+v %v %v", r, ok, err)
	}
}

// ADR-0319 §1: a teardown that passes removes the copy whole, and leaves the
// rest of the story's agent cache folder.
func TestATeardownThatPassesRemovesTheCopy(t *testing.T) {
	f := copySetup(t)
	f.deployed(t, "S-0578")
	settings := filepath.Join(filepath.Dir(AgentCopyDir(f.root, "S-0578")), "settings.json")
	if err := os.WriteFile(settings, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := f.c.Teardown(context.Background(), f.m, "S-0578", f.worktree, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(AgentCopyDir(f.root, "S-0578")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the copy is removed once the namespace is torn down: %v", err)
	}
	if _, err := os.Stat(settings); err != nil {
		t.Errorf("the rest of the story's agent cache folder stays: %v", err)
	}
}

// ADR-0319 §1: a teardown that fails keeps the copy, its record failed at
// the teardown step and the teardown's log beside the deploy's.
func TestATeardownThatFailsKeepsTheCopy(t *testing.T) {
	f := copySetup(t)
	f.deployed(t, "S-0578")
	f.m.Kluster.Teardown = []string{"scripts/teardown.sh"}
	t.Setenv("FAKE_FAIL", "teardown")
	if _, err := f.c.Teardown(context.Background(), f.m, "S-0578", f.worktree, &bytes.Buffer{}); err == nil {
		t.Fatal("the teardown fails")
	}
	c := Folder(AgentCopyDir(f.root, "S-0578"))
	r, ok, err := c.ReadRecord("s-0578")
	if err != nil || !ok || r.Outcome != OutcomeFailed || r.FailedStep != StepTeardown {
		t.Errorf("the copied record is the failed teardown's: %+v %v %v", r, ok, err)
	}
	if data, err := os.ReadFile(c.Log("s-0578", StepTeardown)); err != nil || !strings.Contains(string(data), "teardown broke") {
		t.Errorf("the teardown's log is copied: %q %v", data, err)
	}
	for _, p := range []string{c.Log("s-0578", StepDeploy), c.AgentKubeconfig("s-0578"), filepath.Join(c.State("s-0578"), "db", "url")} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("%s stays in the copy: %v", p, err)
		}
	}
}

// ADR-0319 §1: the default namespace, which no story's agent works in, gets
// no copy.
func TestAPromotionWritesNoCopy(t *testing.T) {
	f := copySetup(t)
	if _, err := f.c.Promote(context.Background(), f.m, "main", &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(f.root, ".flai-cache", "agents")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("no agent cache folder is written for the default namespace: %v", err)
	}
}

// ADR-0319 §1: a copy that cannot be written fails nothing; the step's log
// says so, and the main checkout's record and logs stay the evidence.
func TestACopyThatFailsFailsNoDeploy(t *testing.T) {
	f := copySetup(t)
	agent := filepath.Dir(AgentCopyDir(f.root, "S-0578"))
	if err := os.MkdirAll(filepath.Dir(agent), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(agent, []byte("not a folder"), 0o600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	r, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &out)
	if err != nil || r.Outcome != OutcomePassed {
		t.Fatalf("the deploy passes without its copy: %+v %v", r, err)
	}
	for _, step := range []string{StepBuild, StepLoad, StepDeploy} {
		data, err := os.ReadFile(LogPath(f.root, "s-0578", step))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "the copy of the namespace s-0578 for S-0578's agent, in "+AgentCopyDir(f.root, "S-0578")+", is not up to date") {
			t.Errorf("the %s step's log says the copy was not written:\n%s", step, data)
		}
	}
	if data, _ := os.ReadFile(LogPath(f.root, "s-0578", StepBuild)); !strings.Contains(string(data), "build ran") {
		t.Errorf("the note leaves the step's own output in its log:\n%s", data)
	}
	if strings.Contains(out.String(), "not up to date") {
		t.Errorf("the progress lines say nothing of the copy: %q", out.String())
	}
	if kept, ok, _ := ReadRecord(f.root, "s-0578"); !ok || kept.Outcome != OutcomePassed {
		t.Errorf("the main checkout's record is kept: %+v", kept)
	}
}

func TestInlineAuthority(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "ca.crt"), []byte(fakeCA), 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "s-0578.kubeconfig")
	for _, ca := range []string{"ca.crt", filepath.Join(dir, "ca.crt")} {
		in, err := agentKubeconfig(kubeCluster{Server: "https://127.0.0.1:6550", CertificateAuthority: ca}, "k3d-x", "s-0578", "tok")
		if err != nil {
			t.Fatal(err)
		}
		out, err := inlineAuthority(in, path)
		if err != nil {
			t.Fatal(err)
		}
		var kc kubeconfig
		if err := yaml.Unmarshal(out, &kc); err != nil {
			t.Fatal(err)
		}
		if kc.Clusters[0].Cluster != (kubeCluster{Server: "https://127.0.0.1:6550", CertificateAuthorityData: base64.StdEncoding.EncodeToString([]byte(fakeCA))}) || kc.Users[0].User.Token != "tok" {
			t.Errorf("the file %s is inlined: %+v", ca, kc)
		}
	}
	in, err := agentKubeconfig(kubeCluster{Server: "https://127.0.0.1:6550", CertificateAuthorityData: "RkFLRS1DQQ=="}, "k3d-x", "s-0578", "tok")
	if err != nil {
		t.Fatal(err)
	}
	if out, err := inlineAuthority(in, path); err != nil || !bytes.Equal(out, in) {
		t.Errorf("a kubeconfig naming no file is copied as it is: %v\n%s", err, out)
	}
	missing, err := agentKubeconfig(kubeCluster{Server: "https://127.0.0.1:6550", CertificateAuthority: "gone.crt"}, "k3d-x", "s-0578", "tok")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := inlineAuthority(missing, path); err == nil || !strings.Contains(err.Error(), filepath.Join(dir, "gone.crt")) {
		t.Errorf("a file that is not there is named: %v", err)
	}
}

// TestTheCopyWritesAgainOnlyTheFilesThatChanged shows that mirrorDir leaves
// a copied file whose size and modification time match its source's, and
// writes again one whose source changed, with the source's time.
func TestTheCopyWritesAgainOnlyTheFilesThatChanged(t *testing.T) {
	src, dst := t.TempDir(), filepath.Join(t.TempDir(), "copy")
	for name, data := range map[string]string{"build.log": "built\n", "deploy.log": "deploying\n"} {
		if err := os.WriteFile(filepath.Join(src, name), []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := mirrorDir(src, dst); err != nil {
		t.Fatal(err)
	}
	// A copy that kept its source's size and time is taken as current, so a
	// marker of the same length written over it with that time stays.
	from, err := os.Stat(filepath.Join(src, "build.log"))
	if err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(dst, "build.log")
	if err := os.WriteFile(marker, []byte("KEPT!\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(marker, from.ModTime(), from.ModTime()); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "deploy.log"), []byte("deploying\ndeployed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := mirrorDir(src, dst); err != nil {
		t.Fatal(err)
	}
	if got := copied(t, dst); got["build.log"] != "KEPT!\n" || got["deploy.log"] != "deploying\ndeployed\n" {
		t.Fatalf("copy = %q, want build.log left as it was and deploy.log written again", got)
	}
	d, err := os.Stat(filepath.Join(dst, "deploy.log"))
	if err != nil {
		t.Fatal(err)
	}
	s, err := os.Stat(filepath.Join(src, "deploy.log"))
	if err != nil {
		t.Fatal(err)
	}
	if !d.ModTime().Equal(s.ModTime()) {
		t.Fatalf("deploy.log's copy has time %v, want its source's %v", d.ModTime(), s.ModTime())
	}
}
