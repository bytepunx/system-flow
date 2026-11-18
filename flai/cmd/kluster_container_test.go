package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/kluster"
)

// containerKubeconfig is the copy's kubeconfig of S-0007's namespace, its
// certificate authority inlined, naming the API server as the host sees it.
const containerKubeconfig = `apiVersion: v1
kind: Config
clusters:
- name: k3d-demo-integration
  cluster:
    server: https://127.0.0.1:6443
    certificate-authority-data: Y2E=
contexts:
- name: k3d-demo-integration
  context:
    cluster: k3d-demo-integration
    user: flai-agent
    namespace: s-0007
users:
- name: flai-agent
  user:
    token: agent-token
current-context: k3d-demo-integration
`

// containerProject is the main checkout as S-0007's agent container mounts
// it (ADR-0253, ADR-0319): system-flow.yaml with a kluster: block, wip/, the
// story's worktree, whose .git file names the main checkout and which holds
// the kluster config, and the story's copy of its namespace under
// .flai-cache/agents/S-0007/kluster, with a record at klusterHead, a deploy
// log, a state folder, and a kubeconfig; and no .flai-cache/kluster and no
// kluster config at the main checkout's root. The session is S-0007's
// agent, with a KUBECONFIG of its own and no API host. It returns the main
// checkout, the worktree, and the copy.
func containerProject(t *testing.T) (main, wt string, cp kluster.Folder) {
	t.Helper()
	main, err := filepath.EvalSymlinks(tempProject(t))
	if err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(filepath.Join(main, "system-flow.yaml"))
	manifest := string(b) + "kluster:\n  config: integration/kluster.yaml\n  deploy: [scripts/kluster-deploy.sh]\n  namespace: demo\n"
	wt = filepath.Join(main, ".flai-cache", "worktrees", "S-0007")
	cp = kluster.Folder(kluster.AgentCopyDir(main, "S-0007"))
	host := kluster.Folder(kluster.Dir(main))
	r := kluster.Record{Story: "S-0007", Namespace: "s-0007", Commit: klusterHead, Outcome: kluster.OutcomePassed, StartedAt: queuedAt, Checkout: wt,
		Steps: []kluster.Step{{Name: kluster.StepDeploy, Outcome: kluster.OutcomePassed, Log: host.Log("s-0007", kluster.StepDeploy)}}}
	rec, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	rel := func(p string) string {
		s, err := filepath.Rel(main, p)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	files := map[string]string{}
	files["system-flow.yaml"] = manifest
	files[".git/worktrees/S-0007/HEAD"] = "ref: refs/heads/story/S-0007\n"
	files[rel(filepath.Join(wt, ".git"))] = "gitdir: " + filepath.Join(main, ".git", "worktrees", "S-0007") + "\n"
	files[rel(filepath.Join(wt, "system-flow.yaml"))] = manifest
	files[rel(filepath.Join(wt, "integration", "kluster.yaml"))] = "name: demo-integration\nprofile: authstar\nprovider: k3d\n"
	files[rel(cp.Record("s-0007"))] = string(rec)
	files[rel(cp.Log("s-0007", kluster.StepDeploy))] = "deployed from the copy\n"
	files[rel(filepath.Join(cp.State("s-0007"), "seed"))] = "state\n"
	files[rel(cp.AgentKubeconfig("s-0007"))] = containerKubeconfig
	writeFiles(t, main, files)
	for _, absent := range []string{kluster.Dir(main), filepath.Join(main, "integration", "kluster.yaml")} {
		if _, err := os.Stat(absent); err == nil {
			t.Fatalf("a container does not mount %s", absent)
		}
	}
	t.Setenv("FLAI_STORY", "")
	t.Setenv("FLAI_AGENT", "agent-S-0007")
	t.Setenv("KUBECONFIG", "/elsewhere/kubeconfig")
	t.Setenv("FLAI_KLUSTER_API_HOST", "")
	return main, wt, cp
}

// I-0361, ADR-0319 §2: in a story's container, with no .flai-cache/kluster
// and no kluster config at the main checkout's root, flai kluster exec runs
// in the story's worktree with the variables of the story's copy: KUBECONFIG
// and FLAI_KLUSTER_STATE name the copy's kubeconfig and state folder, and
// the commit is the copied record's. With no API host the kubeconfig is the
// copy's own, and nothing is written beside it.
func TestKlusterExecInAContainerReadsTheStorysCopy(t *testing.T) {
	main, wt, cp := containerProject(t)
	out, errOut, code := runKlusterExit(t, wt, "kluster", "exec", "S-0007", "--", "sh", "-c", "pwd; env")
	if code != 0 {
		t.Fatalf("exec in a container: %d %s %s", code, out, errOut)
	}
	lines := strings.Split(out, "\n")
	if lines[0] != wt {
		t.Errorf("the command runs in the worktree %s: %q", wt, lines[0])
	}
	for _, want := range []string{
		"FLAI_STORY=S-0007",
		"FLAI_KLUSTER_NAME=demo-integration",
		"FLAI_KLUSTER_NAMESPACE=s-0007",
		"FLAI_KLUSTER_COMMIT=" + klusterHead,
		"FLAI_KLUSTER_STATE=" + cp.State("s-0007"),
		"FLAI_KLUSTER_CONTEXT=k3d-demo-integration",
		"KUBECONFIG=" + cp.AgentKubeconfig("s-0007"),
	} {
		if !slices.Contains(lines, want) {
			t.Errorf("the command gets %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "/elsewhere") || strings.Contains(out, kluster.Dir(main)) {
		t.Errorf("the command gets nothing of the shell's KUBECONFIG nor of .flai-cache/kluster:\n%s", out)
	}
	if strings.Contains(errOut, "has no kubeconfig") {
		t.Errorf("the copy's kubeconfig is there: %s", errOut)
	}
	if _, err := os.Stat(filepath.Join(string(cp), "agents", "s-0007.container.kubeconfig")); err == nil {
		t.Errorf("with no API host, no container kubeconfig is written")
	}
}

// ADR-0319 §3: off the host's network, FLAI_KLUSTER_API_HOST names the API
// server's host, and exec writes the copy's kubeconfig with the server's
// host replaced and tls-server-name the host it replaced, 0600, and names it
// as KUBECONFIG; the copy's own is left as it was.
func TestKlusterExecInAContainerNamesTheAPIHost(t *testing.T) {
	_, wt, cp := containerProject(t)
	t.Setenv("FLAI_KLUSTER_API_HOST", "host.docker.internal")
	ckc := filepath.Join(string(cp), "agents", "s-0007.container.kubeconfig")
	out, errOut, code := runKlusterExit(t, wt, "kluster", "exec", "S-0007", "--", "sh", "-c", `echo "$KUBECONFIG"; cat "$KUBECONFIG"`)
	if code != 0 {
		t.Fatalf("exec with an API host: %d %s %s", code, out, errOut)
	}
	first, kc, _ := strings.Cut(out, "\n")
	if first != ckc {
		t.Errorf("KUBECONFIG names the container's kubeconfig %s: %q", ckc, first)
	}
	for _, want := range []string{"server: https://host.docker.internal:6443", "tls-server-name: 127.0.0.1", "token: agent-token", "namespace: s-0007", "certificate-authority-data: Y2E="} {
		if !strings.Contains(kc, want) {
			t.Errorf("the container's kubeconfig has %q:\n%s", want, kc)
		}
	}
	if fi, err := os.Stat(ckc); err != nil || fi.Mode().Perm() != 0o600 {
		t.Errorf("the container's kubeconfig is readable by its owner alone: %v %v", fi, err)
	}
	if b, _ := os.ReadFile(cp.AgentKubeconfig("s-0007")); string(b) != containerKubeconfig {
		t.Errorf("the copy's kubeconfig is left as it was:\n%s", b)
	}
	if strings.Contains(errOut, "has no kubeconfig") {
		t.Errorf("the container's kubeconfig is there: %s", errOut)
	}
}

// I-0361, ADR-0319 §2, §3: flai kluster env prints the copy's variables in a
// container, and with an API host KUBECONFIG names the container's
// kubeconfig.
func TestKlusterEnvInAContainerPrintsTheCopysVariables(t *testing.T) {
	_, wt, cp := containerProject(t)
	out, errOut, code := runKlusterExit(t, wt, "kluster", "env", "S-0007", "--json")
	var vars map[string]string
	if code != 0 || json.Unmarshal([]byte(out), &vars) != nil {
		t.Fatalf("env in a container: %d %s %s", code, out, errOut)
	}
	if vars["KUBECONFIG"] != cp.AgentKubeconfig("s-0007") || vars["FLAI_KLUSTER_STATE"] != cp.State("s-0007") || vars["FLAI_KLUSTER_COMMIT"] != klusterHead || vars["FLAI_KLUSTER_NAME"] != "demo-integration" || len(vars) != 9 {
		t.Errorf("env prints the copy's nine variables: %v", vars)
	}
	t.Setenv("FLAI_KLUSTER_API_HOST", "host.docker.internal")
	out, errOut, code = runKlusterExit(t, wt, "kluster", "env", "S-0007")
	if want := "export KUBECONFIG='" + filepath.Join(string(cp), "agents", "s-0007.container.kubeconfig") + "'"; code != 0 || !slices.Contains(strings.Split(out, "\n"), want) {
		t.Errorf("env with an API host prints %q: %d %s %s", want, code, out, errOut)
	}
}

// I-0361, ADR-0319 §2: flai kluster log prints the copied step log in a
// container, read from the copy whatever path the copied record's step
// names, with the step's outcome; --json gives the log's path from the main
// checkout.
func TestKlusterLogInAContainerPrintsTheCopiedLog(t *testing.T) {
	_, wt, _ := containerProject(t)
	out, errOut, code := runKluster(t, wt, false, "", "kluster", "log", "S-0007", "--step", "deploy")
	if code != 0 || out != "deployed from the copy\n" {
		t.Fatalf("log in a container: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runKluster(t, wt, false, "", "kluster", "log", "S-0007", "--step", "deploy", "--json")
	var got kluster.LogPart
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("log --json in a container: %d %s %s", code, out, errOut)
	}
	if got.Log != ".flai-cache/agents/S-0007/kluster/logs/s-0007/deploy.log" || !got.Exists || got.Outcome != kluster.OutcomePassed {
		t.Errorf("log --json reads the copy: %+v", got)
	}
}
