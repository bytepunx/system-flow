// Package kluster runs a project's local Kubernetes cluster with the kluster
// command of bytepunx/kluster, from the manifest's kluster: block: it brings
// the cluster up, takes it down, and reads its status as data
// (design/system/kluster.md, ADR-0254). flai shells out to kluster as it does
// to git and docker (ADR-0010); the program is reached through Exec, so tests
// put a fake kluster on PATH.
package kluster

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Binary is the program flai runs the cluster with.
const Binary = "kluster"

// ErrNotInstalled says kluster is not on PATH, and how to install it.
var ErrNotInstalled = errors.New("kluster is not installed: flai runs a project's kluster with the kluster command, and none is on PATH; install it with curl -fsSL https://raw.githubusercontent.com/bytepunx/kluster/main/install.sh | bash, or go install github.com/bytepunx/kluster/kluster@latest, then run this again")

// ErrNoBlock says the project's manifest names no cluster, and which block
// to add.
var ErrNoBlock = errors.New("this project's system-flow.yaml has no kluster: block, so flai runs no cluster for it; add one naming the kluster config and the deploy command, such as kluster: {config: integration/kluster.yaml, deploy: [scripts/kluster-deploy.sh]}, as design/system/kluster.md § The kluster: block describes, then run this again")

// Exec runs programs for the package: LookPath finds one on PATH, Run runs
// one in dir, writing its standard output and error to out as they come, so
// that a bring-up of many minutes shows its progress, and Command runs one as
// Cmd describes it, with the environment and standard input a deploy's steps
// need.
type Exec interface {
	LookPath(name string) (string, error)
	Run(ctx context.Context, dir string, out io.Writer, name string, args ...string) error
	Command(ctx context.Context, cmd Cmd) error
}

// System is the real Exec.
type System struct{}

// LookPath reports where name is on PATH.
func (System) LookPath(name string) (string, error) { return exec.LookPath(name) }

// Run runs name with args in dir, its output streamed to out.
func (System) Run(ctx context.Context, dir string, out io.Writer, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdout = out
	cmd.Stderr = out
	return cmd.Run()
}

// Cluster is a project's cluster as its manifest names it, run from Root.
type Cluster struct {
	// Root is the directory kluster runs in, the project's main checkout,
	// where Config is read from and the kubeconfig is kept.
	Root string
	// Config is the kluster config file, from Root.
	Config string
	// Name, Provider, and Profile are what the config file names.
	Name, Provider, Profile string
	// Exec runs kluster.
	Exec Exec
}

// New is the cluster m's kluster: block names, run in root with x: it reads
// the config file for the cluster's name, provider, and profile. A manifest
// with no block is ErrNoBlock.
func New(root string, m manifest.Manifest, x Exec) (Cluster, error) {
	if m.Kluster == nil {
		return Cluster{}, ErrNoBlock
	}
	cfg := strings.TrimSpace(m.Kluster.Config)
	if cfg == "" {
		return Cluster{}, fmt.Errorf("the kluster: block of %s names no config; write kluster.config, the kluster config file from the repository root, and run flai check", filepath.Join(root, "system-flow.yaml"))
	}
	kc, err := manifest.ReadKlusterConfig(os.DirFS(root), filepath.ToSlash(cfg))
	if err != nil {
		return Cluster{}, fmt.Errorf("read the kluster config %s in %s: %w; fix kluster.config in system-flow.yaml or the file, and run flai check", cfg, root, err)
	}
	if kc.Name == "" {
		return Cluster{}, fmt.Errorf("the kluster config %s in %s names no cluster; give it a name:, and run flai check", cfg, root)
	}
	return Cluster{Root: root, Config: cfg, Name: kc.Name, Provider: kc.Provider, Profile: kc.Profile, Exec: x}, nil
}

// KubeconfigPath is where flai keeps the kubeconfig of the cluster run from
// root; it is never merged into ~/.kube/config.
func KubeconfigPath(root string) string {
	return filepath.Join(root, ".flai-cache", "kluster", "kubeconfig")
}

// Context is the cluster's kube context, k3d- and its name.
func (c Cluster) Context() string { return "k3d-" + c.Name }

// global are the arguments every kluster call starts with.
func (c Cluster) global() []string {
	args := []string{"--config", c.Config}
	if c.Provider != "" {
		args = append(args, "--provider", c.Provider)
	}
	return args
}

// UpArgs are the arguments of kluster up.
func (c Cluster) UpArgs() []string { return append(c.global(), "up") }

// DownArgs are the arguments of kluster down, which skip kluster's own
// confirmation: flai asks its own.
func (c Cluster) DownArgs() []string {
	return append(c.global(), "down", "--name", c.Name, "--yes")
}

// StatusArgs are the arguments of kluster status.
func (c Cluster) StatusArgs() []string { return append(c.global(), "status") }

// KubeconfigArgs are the arguments of kluster kubeconfig that write the
// cluster's kubeconfig to KubeconfigPath.
func (c Cluster) KubeconfigArgs() []string {
	return append(c.global(), "kubeconfig", "--name", c.Name, "--output", KubeconfigPath(c.Root))
}

// installed returns ErrNotInstalled when kluster is not on PATH.
func (c Cluster) installed() error {
	if _, err := c.Exec.LookPath(Binary); err != nil {
		return ErrNotInstalled
	}
	return nil
}

// run runs kluster with args in Root, its output to out; a failure says what
// was run, on which cluster, and what to do.
func (c Cluster) run(ctx context.Context, out io.Writer, what string, args []string) error {
	if err := c.Exec.Run(ctx, c.Root, out, Binary, args...); err != nil {
		return fmt.Errorf("kluster %s of the cluster %s, from %s in %s, failed: %w; read kluster's output for the cause, fix it, and run it again", what, c.Name, c.Config, c.Root, err)
	}
	return nil
}

// Up brings the cluster up with kluster up, its output streamed to out, and
// then writes its kubeconfig to KubeconfigPath.
func (c Cluster) Up(ctx context.Context, out io.Writer) error {
	if err := c.installed(); err != nil {
		return err
	}
	if err := c.run(ctx, out, "up", c.UpArgs()); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(KubeconfigPath(c.Root)), 0o755); err != nil {
		return fmt.Errorf("make the folder for the kubeconfig of the cluster %s: %w", c.Name, err)
	}
	return c.run(ctx, out, "kubeconfig", c.KubeconfigArgs())
}

// Down takes the cluster down with kluster down, without kluster's own
// prompt, its output streamed to out.
func (c Cluster) Down(ctx context.Context, out io.Writer) error {
	if err := c.installed(); err != nil {
		return err
	}
	return c.run(ctx, out, "down", c.DownArgs())
}

// Status is what flai knows of the cluster: what its config names, and
// whether kluster status lists it, whether it runs, and its age.
type Status struct {
	Name     string `json:"name"`
	Profile  string `json:"profile,omitempty"`
	Provider string `json:"provider"`
	Context  string `json:"context"`
	// Exists is whether kluster status lists the cluster; false means it
	// was never brought up or was taken down.
	Exists bool `json:"exists"`
	// Running is whether kluster says it runs.
	Running bool `json:"running"`
	// Age is the cluster's age as kluster prints it, such as 2h50m39s.
	Age string `json:"age,omitempty"`
	// AgeSeconds is Age in seconds, when it reads as a duration.
	AgeSeconds int64 `json:"age_seconds,omitempty"`
	// Uptime is Age as a duration, zero when it does not read as one.
	Uptime time.Duration `json:"-"`
}

// Status reads the cluster's status from kluster status.
func (c Cluster) Status(ctx context.Context) (Status, error) {
	if err := c.installed(); err != nil {
		return Status{}, err
	}
	var buf bytes.Buffer
	if err := c.Exec.Run(ctx, c.Root, &buf, Binary, c.StatusArgs()...); err != nil {
		return Status{}, fmt.Errorf("kluster status, for the cluster %s from %s in %s, failed: %w; run kluster status there to see why\n%s", c.Name, c.Config, c.Root, err, strings.TrimSpace(buf.String()))
	}
	s := Status{Name: c.Name, Profile: c.Profile, Provider: c.Provider, Context: c.Context()}
	row, ok, err := ParseStatus(buf.String(), c.Name)
	if err != nil {
		return Status{}, fmt.Errorf("read kluster status for the cluster %s: %w", c.Name, err)
	}
	if !ok {
		return s, nil
	}
	s.Exists, s.Running, s.Age = true, row.Running, row.Age
	if d, err := time.ParseDuration(row.Age); err == nil {
		s.Uptime, s.AgeSeconds = d, int64(d/time.Second)
	}
	return s, nil
}

// Row is one cluster of kluster status's table.
type Row struct {
	Name    string
	Running bool
	Age     string
}

// ParseStatus finds the cluster name in the table kluster status prints, a
// header naming NAME, RUNNING, and AGE above one row per cluster, columns
// matched by their header, so a column a later kluster adds is passed over.
// ok is false when the table does not list name, or when kluster printed no
// table, as it may with no cluster at all. A header that lacks NAME or
// RUNNING is an error.
func ParseStatus(out, name string) (row Row, ok bool, err error) {
	lines := strings.Split(strings.ReplaceAll(out, "\r\n", "\n"), "\n")
	header := -1
	var col map[string]int
	for i, l := range lines {
		f := strings.Fields(l)
		if len(f) > 0 && f[0] == "NAME" {
			header, col = i, map[string]int{}
			for j, h := range f {
				col[h] = j
			}
			break
		}
	}
	if header < 0 {
		return Row{}, false, nil
	}
	ni, okName := col["NAME"]
	ri, okRun := col["RUNNING"]
	if !okName || !okRun {
		return Row{}, false, fmt.Errorf("kluster status printed a table without NAME and RUNNING columns: %q", strings.TrimSpace(lines[header]))
	}
	ai, okAge := col["AGE"]
	for _, l := range lines[header+1:] {
		f := strings.Fields(l)
		if ni >= len(f) || f[ni] != name {
			continue
		}
		r := Row{Name: name}
		if ri < len(f) {
			r.Running = strings.EqualFold(f[ri], "yes") || strings.EqualFold(f[ri], "true")
		}
		if okAge && ai < len(f) {
			r.Age = f[ai]
		}
		return r, true, nil
	}
	return Row{}, false, nil
}
