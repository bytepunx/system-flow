package kluster

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/url"
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
	return Folder(Dir(root)).AgentKubeconfig(strings.ToLower(story))
}

// AgentKubeconfig is the kubeconfig held to namespace ns kept in f.
func (f Folder) AgentKubeconfig(ns string) string {
	return filepath.Join(string(f), "agents", ns+".kubeconfig")
}

// inlineAuthority is the kubeconfig data, read from path, with the
// certificate authority file each of its clusters names, made absolute from
// path's folder, written into it as certificate-authority-data in place of
// the file, so that it names no file outside it (ADR-0319 §1).
func inlineAuthority(data []byte, path string) ([]byte, error) {
	var kc kubeconfig
	if err := yaml.Unmarshal(data, &kc); err != nil {
		return nil, fmt.Errorf("read the namespace's kubeconfig %s: %w; deploy the namespace again", path, err)
	}
	for i := range kc.Clusters {
		e := &kc.Clusters[i].Cluster
		if e.CertificateAuthority == "" {
			continue
		}
		ca := e.CertificateAuthority
		if !filepath.IsAbs(ca) {
			ca = filepath.Join(filepath.Dir(path), ca)
		}
		pem, err := os.ReadFile(ca)
		if err != nil {
			return nil, fmt.Errorf("read the certificate authority %s the kubeconfig %s names: %w; bring the cluster up again with flai kluster up, then deploy the namespace again", ca, path, err)
		}
		e.CertificateAuthorityData, e.CertificateAuthority = base64.StdEncoding.EncodeToString(pem), ""
	}
	out, err := yaml.Marshal(kc)
	if err != nil {
		return nil, fmt.Errorf("encode the kubeconfig %s with its certificate authority inlined: %w", path, err)
	}
	return out, nil
}

// APIHostEnv is the variable a story's container is given off the host's
// network, naming the host the container reaches the cluster's API server at
// (ADR-0319 §3).
const APIHostEnv = "FLAI_KLUSTER_API_HOST"

// AgentFolder is the folder flai kluster log, env, and exec read story's
// namespace from: the main checkout's .flai-cache/kluster under root when it
// is there, as on the host, else the story's copy, AgentCopyDir, as in the
// story's container (ADR-0319 §2); copied says which.
func AgentFolder(root, story string) (f Folder, copied bool) {
	if info, err := os.Stat(Dir(root)); err == nil && info.IsDir() {
		return Folder(Dir(root)), false
	}
	return Folder(AgentCopyDir(root, story)), true
}

// ContainerKubeconfig is the kubeconfig of namespace ns kept in f with the
// API server named by the host APIHostEnv gives (ADR-0319 §3).
func (f Folder) ContainerKubeconfig(ns string) string {
	return filepath.Join(string(f), "agents", ns+".container.kubeconfig")
}

// AgentEnv is the environment a story's agent gets for its namespace, at the
// commit its record names, read from AgentFolder: Env's variables, with
// KUBECONFIG the namespace's own kubeconfig in place of flai's and
// FLAI_KLUSTER_STATE the folder's state folder. Where the folder is the
// story's copy and apiHost, APIHostEnv's value, is not empty, KUBECONFIG is
// ContainerKubeconfig, written again from the copy's kubeconfig. A story with
// no record is ErrNotDeployed.
func (c Cluster) AgentEnv(story, apiHost string) ([]string, error) {
	t, err := StoryTarget(story, "", "")
	if err != nil {
		return nil, err
	}
	f, copied := AgentFolder(c.Root, t.Story)
	r, ok, err := f.ReadRecord(t.Namespace)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("%s: %w; deploy it first", t.Story, ErrNotDeployed)
	}
	kube := f.AgentKubeconfig(t.Namespace)
	if copied && apiHost != "" {
		if kube, err = writeContainerKubeconfig(f, t.Namespace, apiHost); err != nil {
			return nil, err
		}
	}
	t.Commit = r.Commit
	env := c.Env(t)
	for i, v := range env {
		switch {
		case strings.HasPrefix(v, "KUBECONFIG="):
			env[i] = "KUBECONFIG=" + kube
		case strings.HasPrefix(v, "FLAI_KLUSTER_STATE="):
			env[i] = "FLAI_KLUSTER_STATE=" + f.State(t.Namespace)
		}
	}
	return env, nil
}

// writeContainerKubeconfig writes f's kubeconfig of namespace ns at
// ContainerKubeconfig with apiHost as its server's host, whole and readable
// by its owner alone, and returns its path. With no kubeconfig in f, as
// before a deploy step passes, the one written before is removed, so that
// none outlives its source, and the path is returned all the same.
func writeContainerKubeconfig(f Folder, ns, apiHost string) (string, error) {
	src, dst := f.AgentKubeconfig(ns), f.ContainerKubeconfig(ns)
	data, err := os.ReadFile(src)
	if errors.Is(err, fs.ErrNotExist) {
		if err := os.Remove(dst); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return "", fmt.Errorf("remove the container's kubeconfig %s, whose source %s is gone: %w", dst, src, err)
		}
		return dst, nil
	}
	if err != nil {
		return "", fmt.Errorf("read the namespace %s's kubeconfig %s: %w", ns, src, err)
	}
	out, err := withAPIHost(data, src, apiHost)
	if err != nil {
		return "", err
	}
	if err := atomicfile.WriteFile(dst, out, 0o600); err != nil {
		return "", fmt.Errorf("write the namespace %s's kubeconfig for the container, naming the API server at %s, to %s: %w", ns, apiHost, dst, err)
	}
	return dst, nil
}

// withAPIHost is the kubeconfig data, read from path, with each cluster's
// server host replaced by host, its scheme and port kept, and its
// tls-server-name the host it replaced unless it names one already, so that
// the API server's certificate is still verified against the name it was
// issued for (ADR-0319 §3).
func withAPIHost(data []byte, path, host string) ([]byte, error) {
	var kc kubeconfig
	if err := yaml.Unmarshal(data, &kc); err != nil {
		return nil, fmt.Errorf("read the namespace's kubeconfig %s: %w; ask for the namespace's deploy again", path, err)
	}
	for i := range kc.Clusters {
		e := &kc.Clusters[i].Cluster
		u, err := url.Parse(e.Server)
		if err != nil || u.Host == "" {
			return nil, fmt.Errorf("the cluster %s in the kubeconfig %s names the server %q, which is not a URL with a host; ask for the namespace's deploy again", kc.Clusters[i].Name, path, e.Server)
		}
		was := u.Hostname()
		switch {
		case u.Port() != "":
			u.Host = net.JoinHostPort(host, u.Port())
		case strings.Contains(host, ":"):
			u.Host = "[" + host + "]"
		default:
			u.Host = host
		}
		e.Server = u.String()
		if e.TLSServerName == "" {
			e.TLSServerName = was
		}
	}
	out, err := yaml.Marshal(kc)
	if err != nil {
		return nil, fmt.Errorf("encode the kubeconfig %s with the API server at %s: %w", path, host, err)
	}
	return out, nil
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
