package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/config"
)

// The container's writes that git would later run on the host (S-0064,
// ADR-0027). The clone is mounted read-write because acceptance needs it,
// .git included. A process in the container could therefore write a hook, or
// a setting that makes git run a command, and the operator's next git command
// on the host would run it as them. None of that shows in git status. These
// paths are mounted again, read-only, over the read-write clone: the hooks,
// the repository's git config, .git/info (exclude can hide a planted file
// from git status, attributes can name a filter), and two files flai itself
// reads on the host. Acceptance, saves, and moves write none of them; git
// branch -d rewrites the config with the same content and tolerates failing.

// guardedPaths are relative to the clone. Directories are created on the
// host when missing, so that there is something to mount and nothing for the
// container to create in their place.
var guardedDirs = []string{".git/hooks", ".git/info"}
var guardedFiles = []string{".git/config", ".flai-cache/dashboard.token", ".flai-cache/" + agentKeyFileName}

// gitGuardArgs returns the read-only mounts and what to tell the operator.
// hostRoot is the clone on the host, mount where the container sees it.
func (a *app) gitGuardArgs(hostRoot, mount string) (args []string, guarded []string, note string) {
	gitDir := filepath.Join(hostRoot, ".git")
	if st, err := os.Stat(gitDir); err != nil || !st.IsDir() {
		return nil, nil, "no .git directory at " + hostRoot + " (not a git repository, or a linked worktree): nothing to protect from the container"
	}
	add := func(rel string) {
		src := filepath.Join(hostRoot, rel)
		if strings.ContainsAny(src, ":") && !strings.HasPrefix(src, "/") {
			return
		}
		args = append(args, "--volume", src+":"+mount+"/"+filepath.ToSlash(rel)+":ro")
		guarded = append(guarded, filepath.ToSlash(rel))
	}
	dirs := append([]string{}, guardedDirs...)
	// flai serve's registry, when flai's config (and so the registry) is kept
	// inside the clone: it names the projects flai serves, the dashboards it
	// dials, and the credential files it reads, none of which the container
	// may choose (ADR-0029).
	if rel, err := filepath.Rel(hostRoot, string(a.serveDir())); err == nil && !strings.HasPrefix(rel, "..") {
		dirs = append(dirs, rel)
	}
	for _, rel := range dirs {
		if err := os.MkdirAll(filepath.Join(hostRoot, rel), 0o700); err != nil {
			continue
		}
		add(rel)
	}
	files := append([]string{}, guardedFiles...)
	// A host flai config kept inside the clone (FLAI_CONFIG, as this
	// repository's scripts do) names the image and the push key flai
	// dashboard will use next time: the container must not rewrite it.
	if cfg, err := filepath.Abs(config.ResolvePath(a.configPath)); err == nil {
		if rel, err := filepath.Rel(hostRoot, cfg); err == nil && !strings.HasPrefix(rel, "..") {
			files = append(files, rel)
		}
	}
	for _, rel := range files {
		if st, err := os.Stat(filepath.Join(hostRoot, rel)); err == nil && st.Mode().IsRegular() {
			add(rel)
		}
	}
	return args, guarded, ""
}

// executingKeys are git settings that make git run a command, read another
// file as config, or talk to a different remote than the one shown.
var executingKeys = []string{
	"core.hookspath", "core.fsmonitor", "core.sshcommand", "core.editor", "core.pager", "core.askpass",
	"credential.helper", "include.path", "gpg.program", "gpg.ssh.program", "sequence.editor",
	"diff.external", "uploadpack.packobjectshook",
}

func executingKey(key, value string) bool {
	k := strings.ToLower(key)
	for _, e := range executingKeys {
		if k == e {
			return true
		}
	}
	switch {
	case strings.HasPrefix(k, "alias."):
		return strings.HasPrefix(strings.TrimSpace(value), "!")
	case strings.HasPrefix(k, "includeif."), strings.HasPrefix(k, "credential.") && strings.HasSuffix(k, ".helper"):
		return true
	case strings.HasPrefix(k, "filter.") && (strings.HasSuffix(k, ".clean") || strings.HasSuffix(k, ".smudge") || strings.HasSuffix(k, ".process")):
		return true
	case strings.HasPrefix(k, "diff.") && (strings.HasSuffix(k, ".textconv") || strings.HasSuffix(k, ".command")):
		return true
	case strings.HasPrefix(k, "merge.") && strings.HasSuffix(k, ".driver"):
		return true
	case strings.HasPrefix(k, "url.") && (strings.HasSuffix(k, ".insteadof") || strings.HasSuffix(k, ".pushinsteadof")):
		return true
	case strings.HasPrefix(k, "remote.") && strings.HasSuffix(k, ".pushurl"):
		return true
	}
	return false
}

// auditClone lists what the clone already holds of the kinds the read-only
// mounts prevent from now on: an enabled hook, and git settings that run a
// command or redirect a remote. They may be the operator's own; flai cannot
// tell, so it names them and lets the operator look.
func (a *app) auditClone(hostRoot string) []string {
	var found []string
	hooks := filepath.Join(hostRoot, ".git", "hooks")
	if entries, err := os.ReadDir(hooks); err == nil {
		for _, e := range entries {
			if e.IsDir() || strings.HasSuffix(e.Name(), ".sample") {
				continue
			}
			if info, err := e.Info(); err == nil && info.Mode()&0o111 != 0 {
				found = append(found, "hook .git/hooks/"+e.Name()+" runs on every matching git command on this host")
			}
		}
	}
	if _, err := a.runner.LookPath("git"); err == nil {
		if out, err := a.runner.Run(hostRoot, "git", "config", "--local", "--list"); err == nil {
			for _, line := range strings.Split(out, "\n") {
				key, value, ok := strings.Cut(strings.TrimSpace(line), "=")
				if ok && executingKey(key, value) {
					found = append(found, fmt.Sprintf("git setting %s = %s in .git/config", key, value))
				}
			}
		}
	}
	sort.Strings(found)
	return found
}

func guardMessage(guarded []string, note string, found []string) string {
	var b strings.Builder
	switch {
	case note != "":
		fmt.Fprintf(&b, "  %s\n", note)
	case len(guarded) > 0:
		fmt.Fprintf(&b, "  read-only in the container: %s\n    so it cannot leave a hook or a git setting that runs on this host (ADR-0027)\n", strings.Join(guarded, ", "))
	}
	if len(found) > 0 {
		b.WriteString("  already in this clone, and run by git on this host; check that they are yours:\n")
		for _, f := range found {
			fmt.Fprintf(&b, "    %s\n", f)
		}
	}
	return b.String()
}

// runningGuard reports whether a running dashboard container has the git
// hooks mounted read-only: one started by an older flai does not.
func (a *app) runningGuard(container string) bool {
	out, err := a.runner.Run("", "docker", "inspect", "--format", "{{range .Mounts}}{{if not .RW}}{{.Destination}}{{\"\\n\"}}{{end}}{{end}}", container)
	if err != nil {
		return false
	}
	for _, l := range strings.Split(out, "\n") {
		if strings.HasSuffix(strings.TrimSpace(l), "/.git/hooks") {
			return true
		}
	}
	return false
}
