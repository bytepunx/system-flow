// Package protected holds the paths flai protects, its own policy kept per
// harness (ADR-0130): a person is asked before a file write to one, headless
// through permission_prompt, and no allow rule pre-approves it; the
// acceptance preview reads the list too (ADR-0106). Claude Code's list
// follows "Protected paths" in Claude Code's documentation as of Claude Code
// 2.1.290 (ADR-0106). A harness's own files, such as .codex/ or
// .agents/plugins/, join the list with that harness's adapter, not before:
// OpenCode's list is Claude Code's plus .opencode/, its config file, and
// AGENTS.md, which it reads as instructions (S-0406). The container
// harness's list holds this repository's files that change what a
// containerised agent may do, the images and the scripts
// that build and run them (ADR-0203), by their path from the repository root.
// Path and Changed answer for every harness's list at once. .git
// stays apart, in Git, because permission_prompt never approves it whichever
// the harness.
package protected

import (
	"maps"
	"path/filepath"
	"slices"
	"strings"
)

// ClaudeCode is the harness name of Claude Code, the one guard.ClaudeCode
// names; protected cannot import guard, which imports it.
const ClaudeCode = "claude-code"

// OpenCode is the harness name of OpenCode, the one the opencode adapter
// registers.
const OpenCode = "opencode"

// Container is the harness name of the container adapter, which runs another
// harness in a container (ADR-0203).
const Container = "container"

// List is one harness's protected paths.
type List struct {
	// Folders are the protected folders: every path below one is protected.
	// ".config/git" is the two folders in order.
	Folders []string
	// Files are the protected files, by base name, wherever they lie.
	Files []string
	// RootFolders are protected folders by their path from the repository
	// root: every path below one is protected, and a folder of the same
	// path below another folder is not.
	RootFolders []string
	// RootFiles are protected files by their path from the repository root.
	RootFiles []string
}

// claudeCode is Claude Code's list, which OpenCode's extends.
var claudeCode = List{
	Folders: []string{".git", ".config/git", ".vscode", ".idea", ".husky", ".cargo", ".devcontainer", ".yarn", ".mvn", ".claude"},
	Files: []string{
		".gitconfig", ".gitmodules",
		".bashrc", ".bash_profile", ".bash_login", ".bash_aliases", ".bash_logout", ".zshrc", ".zprofile", ".zshenv", ".zlogin", ".zlogout", ".profile", ".envrc",
		".npmrc", ".yarnrc", ".yarnrc.yml", ".pnp.cjs", ".pnp.loader.mjs", ".pnpmfile.cjs", "bunfig.toml", ".bunfig.toml",
		".bazelrc", ".bazelversion", ".bazeliskrc",
		".pre-commit-config.yaml", "lefthook.yml", "lefthook.yaml", ".lefthook.yml", ".lefthook.yaml",
		"gradle-wrapper.properties", "maven-wrapper.properties",
		".devcontainer.json",
		".ripgreprc", "pyrightconfig.json",
		".mcp.json", ".claude.json",
	},
}

// lists are the protected paths of each harness, by harness name.
var lists = map[string]List{
	ClaudeCode: claudeCode,
	OpenCode: {
		Folders: append(slices.Clone(claudeCode.Folders), ".opencode"),
		Files:   append(slices.Clone(claudeCode.Files), "opencode.json", "opencode.jsonc", "AGENTS.md"),
	},
	Container: {
		RootFolders: []string{"flai/images"},
		RootFiles:   []string{"scripts/agent-image.sh", "scripts/agent-container.sh"},
	},
}

// For returns a copy of harness's protected paths, empty for a harness with
// no list.
func For(harness string) List {
	l := lists[harness]
	return List{Folders: slices.Clone(l.Folders), Files: slices.Clone(l.Files), RootFolders: slices.Clone(l.RootFolders), RootFiles: slices.Clone(l.RootFiles)}
}

// Harnesses are the names of the harnesses with a list, sorted.
func Harnesses() []string { return slices.Sorted(maps.Keys(lists)) }

// parts splits a path relative to a worktree or repository root, with either
// separator, into its names, dropping empty and "." ones.
func parts(rel string) []string {
	var out []string
	for _, p := range strings.Split(filepath.ToSlash(rel), "/") {
		if p != "" && p != "." {
			out = append(out, p)
		}
	}
	return out
}

// has reports whether names, a path's names from the root, is a path of l:
// its base name is one of l's files, one of its folders is one of l's
// folders, it is one of l's root files, or it lies below one of l's root
// folders.
func (l List) has(names []string) bool {
	if slices.Contains(l.Files, names[len(names)-1]) {
		return true
	}
	path := strings.Join(names, "/")
	if slices.Contains(l.RootFiles, path) {
		return true
	}
	for _, f := range l.RootFolders {
		if strings.HasPrefix(path, f+"/") {
			return true
		}
	}
	dirs := names[:len(names)-1]
	for _, f := range l.Folders {
		seq := strings.Split(f, "/")
		for i := 0; i+len(seq) <= len(dirs); i++ {
			if slices.Equal(dirs[i:i+len(seq)], seq) {
				return true
			}
		}
	}
	return false
}

// Path reports whether rel, a path relative to a worktree or repository root,
// is protected for some harness: one of its folders is a protected folder,
// its base name is a protected file, it is a protected root file or lies
// below a protected root folder, or it is or lies in .git. An absolute path
// matches the folders and files but not the root ones.
func Path(rel string) bool {
	names := parts(rel)
	if len(names) == 0 {
		return false
	}
	if Git(rel) {
		return true
	}
	for _, l := range lists {
		if l.has(names) {
			return true
		}
	}
	return false
}

// Git reports whether rel, a path relative to a worktree or repository root,
// has a .git folder or file in it: a worktree's .git file is git's link to
// the repository, as protected as the repository's own .git folder.
func Git(rel string) bool {
	return slices.Contains(parts(rel), ".git")
}

// Changed returns the paths protected for some harness, in the list's order:
// for a story's branch, the changed files only the operator accepts.
func Changed(paths []string) []string {
	var out []string
	for _, p := range paths {
		if Path(p) {
			out = append(out, p)
		}
	}
	return out
}
