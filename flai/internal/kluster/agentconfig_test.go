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

	"github.com/goccy/go-yaml"
)

// readAgentKubeconfig reads story's namespace's kubeconfig as kubectl would.
func (f deployFixture) readAgentKubeconfig(t *testing.T, story string) kubeconfig {
	t.Helper()
	data, err := os.ReadFile(AgentKubeconfigPath(f.root, story))
	if err != nil {
		t.Fatalf("no kubeconfig of %s's namespace: %v", story, err)
	}
	var kc kubeconfig
	if err := yaml.Unmarshal(data, &kc); err != nil {
		t.Fatalf("the namespace's kubeconfig does not read: %v\n%s", err, data)
	}
	return kc
}

func TestAgentKubeconfigPath(t *testing.T) {
	if got, want := AgentKubeconfigPath("/p", "S-0578"), "/p/.flai-cache/kluster/agents/s-0578.kubeconfig"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

// S-0763 criterion 1: a passing deploy makes the ServiceAccount flai-agent,
// binds it to edit, and writes the namespace's kubeconfig with its 24-hour
// token, the cluster's server and certificate authority, and the namespace.
func TestAPassingDeployWritesTheNamespacesKubeconfig(t *testing.T) {
	f := deploySetup(t)
	var out bytes.Buffer
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &out); err != nil {
		t.Fatal(err)
	}
	calls := f.calls(t)
	if want := []string{"deploy " + f.worktree + " ", applyS0578 + KubeconfigPath(f.root), tokenS0578 + KubeconfigPath(f.root)}; len(calls) < 3 || !slices.Equal(calls[len(calls)-3:], want) {
		t.Errorf("after deploy, the account applied and a 24-hour token made with flai's kubeconfig:\n got %q\nwant %q", calls, want)
	}
	applied, err := os.ReadFile(filepath.Join(f.env, "kubectl.apply"))
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range []string{"kind: ServiceAccount\nmetadata:\n  name: flai-agent\n  namespace: s-0578", "kind: RoleBinding\nmetadata:\n  name: flai-agent\n  namespace: s-0578", "kind: ClusterRole\n  name: edit", "- kind: ServiceAccount\n  name: flai-agent\n  namespace: s-0578"} {
		if !strings.Contains(string(applied), w) {
			t.Errorf("the manifest applied has %q:\n%s", w, applied)
		}
	}

	path := AgentKubeconfigPath(f.root, "S-0578")
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Errorf("the kubeconfig is its owner's alone: %v", fi.Mode().Perm())
	}
	kc := f.readAgentKubeconfig(t, "S-0578")
	if kc.CurrentContext != "k3d-authstar-integration" || len(kc.Contexts) != 1 || kc.Contexts[0].Name != "k3d-authstar-integration" {
		t.Errorf("its one context is the cluster's, current: %+v", kc)
	}
	if x := kc.Contexts[0].Context; x.Namespace != "s-0578" || x.User != AgentAccount || x.Cluster != "k3d-authstar-integration" {
		t.Errorf("its context is the namespace's, as flai-agent: %+v", x)
	}
	if len(kc.Clusters) != 1 || kc.Clusters[0].Cluster != (kubeCluster{Server: "https://0.0.0.0:6550", CertificateAuthorityData: "RkFLRS1DQQ=="}) {
		t.Errorf("its cluster is flai's kubeconfig's for the context: %+v", kc.Clusters)
	}
	if len(kc.Users) != 1 || kc.Users[0].Name != AgentAccount || kc.Users[0].User.Token != fakeToken {
		t.Errorf("its user is the token: %+v", kc.Users)
	}
	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "client-key-data") || strings.Contains(string(data), "k3d-other") {
		t.Errorf("nothing else of flai's kubeconfig is copied:\n%s", data)
	}

	log, err := os.ReadFile(LogPath(f.root, "s-0578", StepDeploy))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(log), "serviceaccount/flai-agent configured") || !strings.Contains(string(log), "wrote the namespace s-0578's kubeconfig "+path) {
		t.Errorf("the deploy step's log tells it: %q", log)
	}
	if strings.Contains(string(log), fakeToken) || strings.Contains(out.String(), fakeToken) {
		t.Errorf("the token is in no log:\n%s\n%s", log, out.String())
	}
}

// S-0763 criterion 1: a deploy again makes a new token and writes it.
func TestADeployAgainWritesANewKubeconfig(t *testing.T) {
	f := deploySetup(t)
	f.deployed(t, "S-0578")
	t.Setenv("FAKE_TOKEN", "fake-token-2")
	f.reset(t)
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitB, &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	if calls := f.calls(t); !slices.Contains(calls, tokenS0578+KubeconfigPath(f.root)) {
		t.Errorf("a new token is made: %q", calls)
	}
	if kc := f.readAgentKubeconfig(t, "S-0578"); kc.Users[0].User.Token != "fake-token-2" {
		t.Errorf("the kubeconfig has the new token: %+v", kc.Users)
	}
}

func TestADeployStepThatFailsMakesNoKubeconfig(t *testing.T) {
	f := deploySetup(t)
	t.Setenv("FAKE_FAIL", "deploy")
	if _, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{}); err == nil {
		t.Fatal("the deploy was to fail")
	}
	for _, c := range f.calls(t) {
		if strings.HasPrefix(c, "kubectl ") {
			t.Errorf("no kubectl runs after the deploy step fails: %q", c)
		}
	}
	if _, err := os.Stat(AgentKubeconfigPath(f.root, "S-0578")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("no kubeconfig is written: %v", err)
	}
}

// A kubeconfig that cannot be made fails the deploy at the deploy step, with
// kubectl's exit code and output in its log, and the record's steps as they
// are.
func TestAKubeconfigThatCannotBeMadeFailsTheDeployStep(t *testing.T) {
	f := deploySetup(t)
	t.Setenv("FAKE_KUBECTL_FAIL", "yes")
	r, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{})
	var se *StepError
	if !errors.As(err, &se) || se.Step != StepDeploy || se.ExitCode != 5 || !strings.Contains(err.Error(), "kubectl apply") {
		t.Fatalf("the failure is the deploy step's, with kubectl's exit code: %v", err)
	}
	if got, want := stepNames(r.Steps), []string{"build=passed", "load=passed", "deploy=failed"}; !slices.Equal(got, want) {
		t.Errorf("the record keeps its steps:\n got %q\nwant %q", got, want)
	}
	if kept, _, _ := ReadRecord(f.root, "s-0578"); kept.Outcome != OutcomeFailed || kept.FailedStep != StepDeploy || kept.ExitCode != 5 {
		t.Errorf("the kept record says the deploy step failed: %+v", kept)
	}
	if log, _ := os.ReadFile(LogPath(f.root, "s-0578", StepDeploy)); !strings.Contains(string(log), "deploy ran") || !strings.Contains(string(log), "kubectl broke") {
		t.Errorf("the deploy step's log has kubectl's output: %q", log)
	}
	if _, err := os.Stat(AgentKubeconfigPath(f.root, "S-0578")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("no kubeconfig is written: %v", err)
	}
}

func TestAKubeconfigWithoutTheClustersContextFailsTheDeployStep(t *testing.T) {
	f := deploySetup(t)
	if err := os.WriteFile(KubeconfigPath(f.root), []byte("apiVersion: v1\nkind: Config\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	_, err := f.c.Deploy(context.Background(), f.m, "S-0578", "review", f.worktree, commitA, &bytes.Buffer{})
	var se *StepError
	if !errors.As(err, &se) || se.Step != StepDeploy || !strings.Contains(err.Error(), "no context k3d-authstar-integration") || !strings.Contains(err.Error(), "flai kluster up") {
		t.Fatalf("a kubeconfig without the cluster's context fails the deploy step, saying what to do: %v", err)
	}
	for _, c := range f.calls(t) {
		if strings.HasPrefix(c, "kubectl ") {
			t.Errorf("no kubectl runs: %q", c)
		}
	}
}

func TestClusterEntry(t *testing.T) {
	cfg := func(cluster string) []byte {
		return []byte("clusters:\n- name: k3d-a\n  cluster:\n" + cluster + "contexts:\n- name: k3d-a\n  context:\n    cluster: k3d-a\n    user: admin\n- name: k3d-b\n  context:\n    cluster: k3d-gone\n")
	}
	e, err := clusterEntry(cfg("    server: https://a:6443\n    certificate-authority: ca.crt\n"), "/p/.flai-cache/kluster/kubeconfig", "k3d-a")
	if err != nil || e.Server != "https://a:6443" || e.CertificateAuthority != "/p/.flai-cache/kluster/ca.crt" {
		t.Errorf("a certificate authority file is named from the kubeconfig's folder: %+v %v", e, err)
	}
	if e, err := clusterEntry(cfg("    server: https://a:6443\n    certificate-authority: /etc/ca.crt\n"), "/p/kubeconfig", "k3d-a"); err != nil || e.CertificateAuthority != "/etc/ca.crt" {
		t.Errorf("an absolute one is kept: %+v %v", e, err)
	}
	for _, tc := range []struct{ name, data, context, want string }{
		{"no context", string(cfg("    server: https://a:6443\n")), "k3d-c", "no context k3d-c"},
		{"no cluster", string(cfg("    server: https://a:6443\n")), "k3d-b", "names the cluster k3d-gone, which it does not list"},
		{"no server", string(cfg("    insecure-skip-tls-verify: true\n")), "k3d-a", "names no server"},
		{"not YAML", "clusters: [", "k3d-a", "read the cluster's kubeconfig"},
	} {
		if _, err := clusterEntry([]byte(tc.data), "/p/kubeconfig", tc.context); err == nil || !strings.Contains(err.Error(), tc.want) || !strings.Contains(err.Error(), "flai kluster up") {
			t.Errorf("%s: %v", tc.name, err)
		}
	}
}

// S-0763: the agent's variables are the deploy's, at the recorded commit,
// with KUBECONFIG the namespace's kubeconfig alone and no images file; on
// the host, an API host changes nothing (ADR-0319 §3).
func TestAgentEnv(t *testing.T) {
	f := deploySetup(t)
	if _, err := f.c.AgentEnv("S-0578", ""); !errors.Is(err, ErrNotDeployed) {
		t.Errorf("a namespace never deployed has no variables: %v", err)
	}
	if _, err := f.c.AgentEnv("main", ""); err == nil {
		t.Error("a name that is no story ID has variables")
	}
	f.deployed(t, "S-0578")
	env, err := f.c.AgentEnv("s-0578", "host.docker.internal")
	if err != nil {
		t.Fatal(err)
	}
	var kube []string
	for _, v := range env {
		if strings.HasPrefix(v, "KUBECONFIG=") {
			kube = append(kube, v)
		}
		if strings.HasPrefix(v, "FLAI_KLUSTER_IMAGES=") {
			t.Errorf("the agent gets no images file: %s", v)
		}
	}
	if want := []string{"KUBECONFIG=" + AgentKubeconfigPath(f.root, "S-0578")}; !slices.Equal(kube, want) {
		t.Errorf("one KUBECONFIG, the namespace's: got %q, want %q", kube, want)
	}
	tg, err := StoryTarget("S-0578", "", commitA)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range f.c.Env(tg) {
		if !strings.HasPrefix(v, "KUBECONFIG=") && !slices.Contains(env, v) {
			t.Errorf("the agent gets the deploy's %s: %q", v, env)
		}
	}
}

func TestTokenOfPassesOverAWarning(t *testing.T) {
	for out, want := range map[string]string{
		"fake-token\n": "fake-token",
		"Warning: a deprecated API\nfake-token\n":  "fake-token",
		"fake-token\nWarning: something went by\n": "fake-token",
		"\n  \n": "",
	} {
		if got := tokenOf(out); got != want {
			t.Errorf("tokenOf(%q) = %q, want %q", out, got, want)
		}
	}
}

// copyFixture is a cluster run from a main checkout with no
// .flai-cache/kluster, as a story's container mounts it, and S-0578's copy
// of its namespace at commitA, with a kubeconfig naming the API server at
// 127.0.0.1:6443 (ADR-0319).
func copyFixture(t *testing.T) (Cluster, Folder) {
	t.Helper()
	f := setup(t, authstarConfig)
	cp := Folder(AgentCopyDir(f.root, "S-0578"))
	kc := "apiVersion: v1\nkind: Config\nclusters:\n- name: k3d-authstar-integration\n  cluster:\n    server: https://127.0.0.1:6443\n    certificate-authority-data: Y2E=\ncontexts:\n- name: k3d-authstar-integration\n  context:\n    cluster: k3d-authstar-integration\n    user: flai-agent\n    namespace: s-0578\nusers:\n- name: flai-agent\n  user:\n    token: agent-token\ncurrent-context: k3d-authstar-integration\n"
	for path, data := range map[string]string{
		cp.Record("s-0578"):          `{"story":"S-0578","namespace":"s-0578","commit":"` + commitA + `","outcome":"passed"}`,
		cp.AgentKubeconfig("s-0578"): kc,
	} {
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(Dir(f.root)); err == nil {
		t.Fatalf("the fixture has %s", Dir(f.root))
	}
	return f.c, cp
}

// ADR-0319 §2: the readers read the main checkout's .flai-cache/kluster when
// it is there, and the story's copy otherwise.
func TestAgentFolder(t *testing.T) {
	root := t.TempDir()
	if f, copied := AgentFolder(root, "s-0578"); !copied || f != Folder(AgentCopyDir(root, "S-0578")) {
		t.Errorf("with no .flai-cache/kluster, the story's copy: %s %v", f, copied)
	}
	if err := os.MkdirAll(Dir(root), 0o755); err != nil {
		t.Fatal(err)
	}
	if f, copied := AgentFolder(root, "S-0578"); copied || f != Folder(Dir(root)) {
		t.Errorf("with .flai-cache/kluster, the main checkout's: %s %v", f, copied)
	}
}

// ADR-0319 §2, §3: read from the copy, KUBECONFIG and FLAI_KLUSTER_STATE
// name the copy's kubeconfig and state folder; with an API host, KUBECONFIG
// names the container's kubeconfig, written from the copy's with the
// server's host replaced, 0600, and removed once the copy's is gone.
func TestAgentEnvFromTheCopy(t *testing.T) {
	c, cp := copyFixture(t)
	env, err := c.AgentEnv("S-0578", "")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"KUBECONFIG=" + cp.AgentKubeconfig("s-0578"), "FLAI_KLUSTER_STATE=" + cp.State("s-0578"), "FLAI_KLUSTER_COMMIT=" + commitA} {
		if !slices.Contains(env, want) {
			t.Errorf("the copy's variables have %q: %q", want, env)
		}
	}
	ckc := cp.ContainerKubeconfig("s-0578")
	if _, err := os.Stat(ckc); err == nil {
		t.Errorf("with no API host, no container kubeconfig is written")
	}
	env, err = c.AgentEnv("S-0578", "host.docker.internal")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(env, "KUBECONFIG="+ckc) {
		t.Errorf("KUBECONFIG names the container's kubeconfig %s: %q", ckc, env)
	}
	data, err := os.ReadFile(ckc)
	if err != nil {
		t.Fatal(err)
	}
	var kc kubeconfig
	if err := yaml.Unmarshal(data, &kc); err != nil || len(kc.Clusters) != 1 || len(kc.Users) != 1 {
		t.Fatalf("the container's kubeconfig reads: %v\n%s", err, data)
	}
	if e := kc.Clusters[0].Cluster; e.Server != "https://host.docker.internal:6443" || e.TLSServerName != "127.0.0.1" || e.CertificateAuthorityData != "Y2E=" || kc.Users[0].User.Token != "agent-token" {
		t.Errorf("the container's kubeconfig names the API host and keeps the rest: %+v %+v", e, kc.Users)
	}
	if fi, err := os.Stat(ckc); err != nil || fi.Mode().Perm() != 0o600 {
		t.Errorf("the container's kubeconfig is readable by its owner alone: %v %v", fi, err)
	}
	if err := os.Remove(cp.AgentKubeconfig("s-0578")); err != nil {
		t.Fatal(err)
	}
	if env, err = c.AgentEnv("S-0578", "host.docker.internal"); err != nil || !slices.Contains(env, "KUBECONFIG="+ckc) {
		t.Errorf("with no kubeconfig in the copy, KUBECONFIG names the container's all the same: %v %q", err, env)
	}
	if _, err := os.Stat(ckc); err == nil {
		t.Errorf("a container kubeconfig whose source is gone is removed")
	}
}

// ADR-0319 §3: each cluster's server host is replaced, its scheme and port
// kept, and tls-server-name is the host it replaced unless it names one.
func TestWithAPIHost(t *testing.T) {
	for _, tc := range []struct{ name, cluster, host, server, tls string }{
		{"an IPv4 host", "    server: https://127.0.0.1:6443\n", "host.docker.internal", "https://host.docker.internal:6443", "127.0.0.1"},
		{"a name", "    server: https://kluster.local:443\n", "host.docker.internal", "https://host.docker.internal:443", "kluster.local"},
		{"an IPv6 host", "    server: https://[::1]:6443\n", "host.docker.internal", "https://host.docker.internal:6443", "::1"},
		{"an IPv6 API host", "    server: https://127.0.0.1:6443\n", "fd00::1", "https://[fd00::1]:6443", "127.0.0.1"},
		{"no port", "    server: https://0.0.0.0\n", "host.docker.internal", "https://host.docker.internal", "0.0.0.0"},
		{"a server name of its own", "    server: https://0.0.0.0:6443\n    tls-server-name: k3d-demo-serverlb\n", "host.docker.internal", "https://host.docker.internal:6443", "k3d-demo-serverlb"},
	} {
		out, err := withAPIHost([]byte("clusters:\n- name: k3d-a\n  cluster:\n"+tc.cluster), "/p/s-0578.kubeconfig", tc.host)
		if err != nil {
			t.Errorf("%s: %v", tc.name, err)
			continue
		}
		var kc kubeconfig
		if err := yaml.Unmarshal(out, &kc); err != nil || len(kc.Clusters) != 1 {
			t.Errorf("%s: %v\n%s", tc.name, err, out)
			continue
		}
		if e := kc.Clusters[0].Cluster; e.Server != tc.server || e.TLSServerName != tc.tls {
			t.Errorf("%s: server %q tls-server-name %q, want %q %q", tc.name, e.Server, e.TLSServerName, tc.server, tc.tls)
		}
	}
	for _, data := range []string{"clusters: [", "clusters:\n- name: k3d-a\n  cluster:\n    server: 127.0.0.1:6443\n"} {
		if _, err := withAPIHost([]byte(data), "/p/s-0578.kubeconfig", "host.docker.internal"); err == nil || !strings.Contains(err.Error(), "/p/s-0578.kubeconfig") {
			t.Errorf("%q is refused, naming the file: %v", data, err)
		}
	}
}
