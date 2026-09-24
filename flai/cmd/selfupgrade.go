package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

func newSelfUpgradeCmd(a *app) *cobra.Command {
	var version, dir, repo, apiBase string
	var check bool
	c := &cobra.Command{
		Use:   "self-upgrade",
		Short: "Install the latest flai release over this binary",
		Long: `Resolves the newest flai release on GitHub (or --version), downloads the
archive for this platform and checksums.txt, verifies the SHA-256, and
replaces the running executable. While the repository is private the API
needs a token: GITHUB_TOKEN, GH_TOKEN, or a gh auth login session. The same
operation from a shell is install.sh at the repository root.`,
		Example: `  flai self-upgrade --check
  flai self-upgrade
  flai self-upgrade --version 1.0.3
  flai self-upgrade --dir /usr/local/bin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opt := selfupgrade.Options{Repo: repo, APIBase: apiBase, Token: a.githubToken(), Version: version}
			rel, err := selfupgrade.Resolve(cmd.Context(), opt)
			if err != nil {
				return err
			}
			current := buildinfo.Version
			running, err := executablePath()
			if err != nil {
				return err
			}
			dest, moved := running, false
			switch {
			case dir != "":
				dest = filepath.Join(dir, flaiBinaryName())
			case insideProject(running):
				// a checkout's build (S-0111): the release goes where install.sh
				// puts it, and the checkout's bin is left to the build
				dest, moved = filepath.Join(installDir(), flaiBinaryName()), true
				current = installedVersion(dest)
			}
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
	c.Flags().StringVar(&version, "version", "", "install this release instead of the latest, e.g. 1.0.3")
	c.Flags().StringVar(&dir, "dir", "", "install into this directory instead of over the running binary")
	c.Flags().BoolVar(&check, "check", false, "report the installed and latest versions without installing")
	c.Flags().StringVar(&repo, "repo", "bytepunx/system-flow", "GitHub repository that publishes flai releases")
	// FLAI_RELEASES_API points flai at another source of releases that
	// answers as GitHub's API does: a mirror, or a test's stand-in. flai host
	// upgrade and check run self-upgrade, so they follow it too (S-0107).
	defaultAPI := "https://api.github.com"
	if v := strings.TrimSpace(os.Getenv("FLAI_RELEASES_API")); v != "" {
		defaultAPI = v
	}
	c.Flags().StringVar(&apiBase, "api", defaultAPI, "GitHub API base URL (default FLAI_RELEASES_API, else api.github.com)")
	_ = c.Flags().MarkHidden("api")
	return c
}

// githubToken is GITHUB_TOKEN, GH_TOKEN, or the gh CLI's session token.
func (a *app) githubToken() string {
	for _, k := range []string{"GITHUB_TOKEN", "GH_TOKEN"} {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	if _, err := a.runner.LookPath("gh"); err != nil {
		return ""
	}
	out, err := a.runner.Run("", "gh", "auth", "token")
	if err != nil {
		return ""
	}
	return out
}

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
