package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

func newSelfUpgradeCmd(a *app) *cobra.Command {
	var version, dir string
	var check, list bool
	var src releaseSource
	c := &cobra.Command{
		Use:   "self-upgrade",
		Short: "Install the latest flai release over this binary, or list the published ones",
		Long: `Resolves the newest flai release on GitHub (or --version), downloads the
archive for this platform and checksums.txt, verifies the SHA-256, and
replaces the running executable. While the repository is private the API
needs a token: GITHUB_TOKEN, GH_TOKEN, or a gh auth login session. The same
operation from a shell is install.sh at the repository root.

--version installs that release, an earlier one included, the way the newest
is installed; a version that is not a published release is refused, naming
the published ones, before anything is downloaded. A release below the
flai.minimum of the project in this folder, or of a project flai serve
serves, is warned about and installed: going back past it is your call.

--list prints the published releases, newest first, and installs nothing. It
marks the installed one (the flai self-upgrade would replace), the newest,
and any below such a project's flai.minimum; with --json it prints a list of
version, tag, published, installed, latest, and below_minimum (each project
and minimum the release is below, left out when none).`,
		Example: `  flai self-upgrade --check
  flai self-upgrade --list
  flai self-upgrade
  flai self-upgrade --version 1.0.3
  flai self-upgrade --dir /usr/local/bin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opt := src.options(a, version)
			if list {
				return a.listFlaiReleases(cmd, opt, dir)
			}
			rel, err := selfupgrade.Resolve(cmd.Context(), opt)
			if err != nil {
				return err
			}
			t, err := installTargetFor(dir)
			if err != nil {
				return err
			}
			current, running, dest, moved := t.current, t.running, t.dest, t.moved
			upToDate := version == "" && rel.Version == current
			if check {
				if a.jsonOut {
					return a.printJSON(map[string]any{"current": current, "latest": rel.Version, "tag": rel.Tag, "up_to_date": upToDate, "path": dest})
				}
				state := "upgrade available"
				if upToDate {
					state = "up to date"
				}
				fmt.Fprintf(a.out, "flai %s installed at %s; latest is %s (%s)\n", current, dest, rel.Version, state)
				return nil
			}
			if upToDate {
				if a.jsonOut {
					return a.printJSON(map[string]any{"current": current, "latest": rel.Version, "up_to_date": true, "path": dest})
				}
				fmt.Fprintf(a.out, "flai %s is already the latest release\n", current)
				return nil
			}
			// ADR-0117 §7: below a project's minimum is warned about, not refused
			for _, m := range belowMinimum(rel.Version, a.flaiMinimums()) {
				a.logger().Warn("chosen release is below a project's flai.minimum", "component", "selfupgrade", "version", rel.Version, "minimum", m.Minimum, "project", m.Project)
			}
			a.logger().Info("downloading release", "tag", rel.Tag, "os", runtime.GOOS, "arch", runtime.GOARCH)
			bin, err := selfupgrade.Download(cmd.Context(), opt, rel)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return err
			}
			if err := selfupgrade.Replace(dest, bin); err != nil {
				return fmt.Errorf("install to %s: %w (try --dir with a writable directory, or sudo)", dest, err)
			}
			if a.jsonOut {
				out := map[string]any{"previous": current, "installed": rel.Version, "tag": rel.Tag, "path": dest}
				if moved {
					out["running"] = running
				}
				return a.printJSON(out)
			}
			was := current
			if was == "" {
				was = "not installed there"
			}
			fmt.Fprintf(a.out, "installed flai %s to %s (was %s)\n", rel.Version, dest, was)
			if moved {
				fmt.Fprintf(a.out, "  the flai that ran this is inside a project (%s), so the release was installed under your home instead;\n  put %s on your PATH ahead of %s\n", running, filepath.Dir(dest), filepath.Dir(running))
			}
			return nil
		},
	}
	c.Flags().StringVar(&version, "version", "", "install this published release instead of the latest, e.g. 1.0.3")
	c.Flags().StringVar(&dir, "dir", "", "install into this directory instead of over the running binary")
	c.Flags().BoolVar(&check, "check", false, "report the installed and latest versions without installing")
	c.Flags().BoolVar(&list, "list", false, "list the published releases, newest first, without installing")
	c.MarkFlagsMutuallyExclusive("list", "version")
	c.MarkFlagsMutuallyExclusive("list", "check")
	src.register(c, "GitHub repository that publishes flai releases")
	return c
}

// releaseSource is where published releases are read: a GitHub repository
// and the API that answers for it. self-upgrade and the dashboard's versions
// and upgrade --published read the same one.
type releaseSource struct{ repo, api string }

// register adds --repo and the hidden --api to c.
func (s *releaseSource) register(c *cobra.Command, repoUsage string) {
	c.Flags().StringVar(&s.repo, "repo", "bytepunx/system-flow", repoUsage)
	// FLAI_RELEASES_API points flai at another source of releases that
	// answers as GitHub's API does: a mirror, or a test's stand-in. flai host
	// upgrade and check run self-upgrade, so they follow it too (S-0107).
	defaultAPI := "https://api.github.com"
	if v := strings.TrimSpace(os.Getenv("FLAI_RELEASES_API")); v != "" {
		defaultAPI = v
	}
	c.Flags().StringVar(&s.api, "api", defaultAPI, "GitHub API base URL (default FLAI_RELEASES_API, else api.github.com)")
	_ = c.Flags().MarkHidden("api")
}

// options are the selfupgrade options for this source, with a token.
func (s releaseSource) options(a *app, version string) selfupgrade.Options {
	return selfupgrade.Options{Repo: s.repo, APIBase: s.api, Token: a.githubToken(), Version: version}
}

// installTarget is where self-upgrade installs and the version it compares
// against.
type installTarget struct {
	current, running, dest string
	moved                  bool // dest is not the running flai, which is inside a project
}

// installTargetFor is the running flai and its version, or dir; or, for a
// flai inside a project, the install directory and the flai there (S-0111).
func installTargetFor(dir string) (installTarget, error) {
	running, err := executablePath()
	if err != nil {
		return installTarget{}, err
	}
	t := installTarget{current: buildinfo.Version, running: running, dest: running}
	switch {
	case dir != "":
		t.dest = filepath.Join(dir, flaiBinaryName())
	case insideProject(running):
		// a checkout's build (S-0111): the release goes where install.sh
		// puts it, and the checkout's bin is left to the build
		t.dest, t.moved = filepath.Join(installDir(), flaiBinaryName()), true
		t.current = installedVersion(t.dest)
	}
	return t, nil
}

// listedFlai is one published flai release as self-upgrade --list prints it.
type listedFlai struct {
	Version      string           `json:"version"`
	Tag          string           `json:"tag"`
	Published    time.Time        `json:"published,omitzero"`
	Installed    bool             `json:"installed"`
	Latest       bool             `json:"latest"`
	BelowMinimum []projectMinimum `json:"below_minimum,omitempty"`
}

// listFlaiReleases prints the published flai releases, newest first, marking
// the installed one, the newest, and those below a project's minimum.
func (a *app) listFlaiReleases(cmd *cobra.Command, opt selfupgrade.Options, dir string) error {
	published, err := selfupgrade.List(cmd.Context(), opt, selfupgrade.TagPrefix)
	if err != nil {
		return err
	}
	t, err := installTargetFor(dir)
	if err != nil {
		return err
	}
	mins := a.flaiMinimums()
	listed := make([]listedFlai, len(published))
	found := false
	for i, r := range published {
		listed[i] = listedFlai{Version: r.Version, Tag: r.Tag, Published: r.Published, Installed: r.Version == t.current, Latest: i == 0, BelowMinimum: belowMinimum(r.Version, mins)}
		found = found || listed[i].Installed
	}
	if a.jsonOut {
		return a.printJSON(listed)
	}
	if len(listed) == 0 {
		fmt.Fprintf(a.out, "no published flai release in %s\n", opt.Repo)
		return nil
	}
	fmt.Fprintf(a.out, "published flai releases in %s, newest first:\n", opt.Repo)
	if err := writeFlaiReleases(a.out, listed); err != nil {
		return err
	}
	if !found {
		installed := t.current
		if installed == "" {
			installed = "none"
		}
		fmt.Fprintf(a.out, "installed at %s: %s, which is not a published release\n", t.dest, installed)
	}
	return nil
}

// writeFlaiReleases writes one line per release: its version, its publish
// date, and its marks.
func writeFlaiReleases(out io.Writer, listed []listedFlai) error {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	for _, r := range listed {
		date := "-"
		if !r.Published.IsZero() {
			date = r.Published.UTC().Format(time.DateOnly)
		}
		var marks []string
		if r.Installed {
			marks = append(marks, "installed")
		}
		if r.Latest {
			marks = append(marks, "latest")
		}
		for _, m := range r.BelowMinimum {
			marks = append(marks, fmt.Sprintf("below minimum %s of %s", m.Minimum, m.Project))
		}
		fmt.Fprintf(w, "  %s\t%s\t%s\n", r.Version, date, strings.Join(marks, ", "))
	}
	return w.Flush()
}

// projectMinimum is a project's flai.minimum and the name it goes by.
type projectMinimum struct {
	Project string `json:"project"`
	Minimum string `json:"minimum"`
}

// flaiMinimums are the flai.minimum of the project in the working directory
// and of each project flai serve serves (its registry, ADR-0117 §7), once per
// project and minimum. A manifest that cannot be read is warned about and
// left out.
func (a *app) flaiMinimums() []projectMinimum {
	type candidate struct{ path, name string }
	var candidates []candidate
	if wd, err := a.workingDir(); err == nil {
		if p, err := manifest.Find(wd); err == nil {
			candidates = append(candidates, candidate{p, filepath.Base(filepath.Dir(p))})
		}
	}
	served, err := a.serveDir().Projects()
	if err != nil {
		a.logger().Warn("served projects not read for their flai.minimum", "component", "selfupgrade", "err", err.Error())
	}
	for _, e := range served {
		candidates = append(candidates, candidate{filepath.Join(e.Root, manifest.File), orDefault(e.Name, e.Key)})
	}
	var out []projectMinimum
	for _, c := range candidates {
		m, err := readMinimum(c.path, c.name)
		if err != nil {
			a.logger().Warn("project manifest not read for its flai.minimum", "component", "selfupgrade", "path", c.path, "err", err.Error())
			continue
		}
		if m.Minimum != "" && !slices.Contains(out, m) {
			out = append(out, m)
		}
	}
	return out
}

// readMinimum is the flai.minimum of the manifest at path, under its name or
// fallback; a minimum above the running flai is read too.
func readMinimum(path, fallback string) (projectMinimum, error) {
	m, err := manifest.Load(path)
	var tooOld *manifest.TooOldError
	switch {
	case errors.As(err, &tooOld):
		return projectMinimum{Project: fallback, Minimum: tooOld.Minimum}, nil
	case err != nil:
		return projectMinimum{}, err
	}
	return projectMinimum{Project: orDefault(m.Name, fallback), Minimum: m.Flai.Minimum}, nil
}

// belowMinimum are the minimums version is below.
func belowMinimum(version string, mins []projectMinimum) []projectMinimum {
	var out []projectMinimum
	for _, m := range mins {
		if buildinfo.Below(version, m.Minimum) {
			out = append(out, m)
		}
	}
	return out
}

// githubToken is GITHUB_TOKEN, GH_TOKEN, or the gh CLI's session token.
func (a *app) githubToken() string { return selfupgrade.Token(a.runner) }

// executablePath is the running flai; tests put one where they need it.
var executablePath = func() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	if real, err := filepath.EvalSymlinks(exe); err == nil {
		exe = real
	}
	return exe, nil
}

func flaiBinaryName() string {
	if runtime.GOOS == "windows" {
		return "flai.exe"
	}
	return "flai"
}

// installDir is where a release is installed when it is not over the running
// flai: FLAI_INSTALL_DIR, else ~/.flai/bin, as install.sh has it.
func installDir() string {
	if d := strings.TrimSpace(os.Getenv("FLAI_INSTALL_DIR")); d != "" {
		return d
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".flai", "bin")
	}
	return filepath.Join(home, ".flai", "bin")
}

// insideProject says whether path lies in a system-flow project: a folder
// above it holds system-flow.yaml. flai is built into bin/ of its own
// repository, and a release installed there would be overwritten by the next
// build, and exists on no other machine (S-0111).
func insideProject(path string) bool {
	for dir := filepath.Dir(path); ; {
		if _, err := os.Stat(filepath.Join(dir, "system-flow.yaml")); err == nil {
			return true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return false
		}
		dir = parent
	}
}

// installedVersion is the version of the flai at path, "" when there is none
// or it does not answer.
func installedVersion(path string) string {
	out, err := exec.Command(path, "version", "--json").Output()
	if err != nil {
		return ""
	}
	var v struct {
		Version string `json:"version"`
	}
	if json.Unmarshal(out, &v) != nil {
		return ""
	}
	return v.Version
}
