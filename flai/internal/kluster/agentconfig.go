package kluster

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
)

// The ServiceAccount, and its RoleBinding, a story's namespace's kubeconfig
// authenticates as; the ClusterRole it is bound to; and how long its token is
// valid (ADR-0292 §2).
const (
	AgentAccount       = "flai-agent"
	AgentClusterRole   = "edit"
	AgentTokenDuration = "24h"
)

// AgentKubeconfigPath is the kubeconfig held to story's namespace, which a
// deploy of it writes and its teardown removes.
func AgentKubeconfigPath(root, story string) string {
	return filepath.Join(Dir(root), "agents", strings.ToLower(story)+".kubeconfig")
}

// AgentEnv is the environment a story's agent gets for its namespace, at the
// commit its record names: Env's variables, with KUBECONFIG the namespace's
// own kubeconfig in place of flai's. A story with no record is
// ErrNotDeployed.
func (c Cluster) AgentEnv(story string) ([]string, error) {
	t, err := StoryTarget(story, "", "")
	if err != nil {
		return nil, err
	}
	r, ok, err := ReadRecord(c.Root, t.Namespace)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%s: %w; deploy it first", t.Story, ErrNotDeployed)
	}
	t.Commit = r.Commit
	env := c.Env(t)
	for i, v := range env {
		if strings.HasPrefix(v, "KUBECONFIG=") {
			env[i] = "KUBECONFIG=" + AgentKubeconfigPath(c.Root, t.Story)
		}
	}
	return env, nil
}

// kubeconfig is the part of a kubeconfig flai reads and writes.
type kubeconfig struct {
	APIVersion     string              `yaml:"apiVersion"`
	Kind           string              `yaml:"kind"`
	Clusters       []kubeconfigCluster `yaml:"clusters"`
	Contexts       []kubeconfigContext `yaml:"contexts"`
	Users          []kubeconfigUser    `yaml:"users"`
	CurrentContext string              `yaml:"current-context"`
}

type kubeconfigCluster struct {
	Name    string      `yaml:"name"`
	Cluster kubeCluster `yaml:"cluster"`
}

// kubeCluster is a cluster entry: the API server and how its certificate is
// trusted.
type kubeCluster struct {
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty"`
	CertificateAuthority     string `yaml:"certificate-authority,omitempty"`
	InsecureSkipTLSVerify    bool   `yaml:"insecure-skip-tls-verify,omitempty"`
	TLSServerName            string `yaml:"tls-server-name,omitempty"`
}

type kubeconfigContext struct {
	Name    string      `yaml:"name"`
	Context kubeContext `yaml:"context"`
}

type kubeContext struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace,omitempty"`
}

type kubeconfigUser struct {
	Name string   `yaml:"name"`
	User kubeUser `yaml:"user"`
}

type kubeUser struct {
	Token string `yaml:"token"`
}

// clusterEntry is the cluster entry that flai's kubeconfig, data read from
// path, names for context, its certificate authority file, when it names
// one, made absolute from the kubeconfig's folder.
func clusterEntry(data []byte, path, context string) (kubeCluster, error) {
	var kc kubeconfig
	if err := yaml.Unmarshal(data, &kc); err != nil {
		return kubeCluster{}, fmt.Errorf("read the cluster's kubeconfig %s: %w; write it again with flai kluster up", path, err)
	}
	name := ""
	for _, x := range kc.Contexts {
		if x.Name == context {
			name = x.Context.Cluster
		}
	}
	if name == "" {
		return kubeCluster{}, fmt.Errorf("the cluster's kubeconfig %s has no context %s naming a cluster; write it again with flai kluster up", path, context)
	}
	for _, cl := range kc.Clusters {
		if cl.Name != name {
			continue
		}
		if cl.Cluster.Server == "" {
			return kubeCluster{}, fmt.Errorf("the cluster %s in the kubeconfig %s names no server; write it again with flai kluster up", name, path)
		}
		e := cl.Cluster
		if e.CertificateAuthority != "" && !filepath.IsAbs(e.CertificateAuthority) {
			e.CertificateAuthority = filepath.Join(filepath.Dir(path), e.CertificateAuthority)
		}
		return e, nil
	}
	return kubeCluster{}, fmt.Errorf("the context %s in the cluster's kubeconfig %s names the cluster %s, which it does not list; write it again with flai kluster up", context, path, name)
}

// agentKubeconfig is the kubeconfig of namespace ns: the cluster entry e,
// the token as its user, and ns as its context's namespace, under context,
// the cluster's context, so that FLAI_KLUSTER_CONTEXT names it too.
func agentKubeconfig(e kubeCluster, context, ns, token string) ([]byte, error) {
	kc := kubeconfig{
		APIVersion:     "v1",
		Kind:           "Config",
		Clusters:       []kubeconfigCluster{{Name: context, Cluster: e}},
		Contexts:       []kubeconfigContext{{Name: context, Context: kubeContext{Cluster: context, User: AgentAccount, Namespace: ns}}},
		Users:          []kubeconfigUser{{Name: AgentAccount, User: kubeUser{Token: token}}},
		CurrentContext: context,
	}
	return yaml.Marshal(kc)
}

// agentManifest is the ServiceAccount flai-agent in namespace ns and the
// RoleBinding flai-agent of it to the ClusterRole edit, for kubectl apply.
func agentManifest(ns string) string {
	return `apiVersion: v1
kind: ServiceAccount
metadata:
  name: ` + AgentAccount + `
  namespace: ` + ns + `
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: ` + AgentAccount + `
  namespace: ` + ns + `
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: ` + AgentClusterRole + `
subjects:
- kind: ServiceAccount
  name: ` + AgentAccount + `
  namespace: ` + ns + `
`
}

// tokenOf is the token in kubectl create token's output: its last line that
// is one word. Cmd gives a command's standard error in the same writer as its
// output, so a warning kubectl prints is passed over rather than taken into
// the token.
func tokenOf(out string) string {
	lines := strings.Split(out, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if l := strings.TrimSpace(lines[i]); l != "" && !strings.ContainsAny(l, " \t") {
			return l
		}
	}
	return ""
}

// writeAgentKubeconfig makes the ServiceAccount flai-agent in t's namespace,
// bound to the ClusterRole edit, with flai's kubeconfig, and writes the
// namespace's kubeconfig at AgentKubeconfigPath with a new token of it, whole
// and readable by its owner alone. Its progress goes to log; the token does
// not.
func (c Cluster) writeAgentKubeconfig(ctx context.Context, t Target, log io.Writer) error {
	src := KubeconfigPath(c.Root)
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("read the cluster's kubeconfig %s for the kubeconfig of the namespace %s: %w; write it again with flai kluster up", src, t.Namespace, err)
	}
	entry, err := clusterEntry(data, src, c.Context())
	if err != nil {
		return err
	}
	if _, err := c.Exec.LookPath("kubectl"); err != nil {
		return fmt.Errorf("kubectl is not on PATH, and flai makes the namespace %s's kubeconfig with it; install kubectl, then deploy again", t.Namespace)
	}
	env := c.Env(t)
	fmt.Fprintf(log, "making the ServiceAccount %s in the namespace %s, bound to the ClusterRole %s\n", AgentAccount, t.Namespace, AgentClusterRole)
	if err := c.Exec.Command(ctx, Cmd{Dir: c.Root, Env: env, Stdin: strings.NewReader(agentManifest(t.Namespace)), Out: log, Name: "kubectl",
		Args: []string{"--context", c.Context(), "--namespace", t.Namespace, "apply", "-f", "-"}}); err != nil {
		return fmt.Errorf("make the ServiceAccount %s and its RoleBinding in the namespace %s with kubectl apply: %w", AgentAccount, t.Namespace, err)
	}
	var tok bytes.Buffer
	if err := c.Exec.Command(ctx, Cmd{Dir: c.Root, Env: env, Out: &tok, Name: "kubectl",
		Args: []string{"--context", c.Context(), "--namespace", t.Namespace, "create", "token", AgentAccount, "--duration", AgentTokenDuration}}); err != nil {
		_, _ = log.Write(tok.Bytes())
		return fmt.Errorf("make a token of the ServiceAccount %s in the namespace %s with kubectl create token: %w", AgentAccount, t.Namespace, err)
	}
	token := tokenOf(tok.String())
	if token == "" {
		return fmt.Errorf("kubectl create token gave no token for the ServiceAccount %s in the namespace %s; run the deploy again", AgentAccount, t.Namespace)
	}
	out, err := agentKubeconfig(entry, c.Context(), t.Namespace, token)
	if err != nil {
		return fmt.Errorf("encode the kubeconfig of the namespace %s: %w", t.Namespace, err)
	}
	path := AgentKubeconfigPath(c.Root, t.Story)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("make the folder of the namespaces' kubeconfigs %s: %w", filepath.Dir(path), err)
	}
	if err := atomicfile.WriteFile(path, out, 0o600); err != nil {
		return fmt.Errorf("write the kubeconfig of the namespace %s to %s: %w", t.Namespace, path, err)
	}
	fmt.Fprintf(log, "wrote the namespace %s's kubeconfig %s, its token valid for %s\n", t.Namespace, path, AgentTokenDuration)
	return nil
}
