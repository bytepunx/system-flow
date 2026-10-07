// Package protected holds the paths Claude Code protects: it asks a person
// before an Edit, Write, MultiEdit, or NotebookEdit of one, headless through
// the permission prompt tool, and no allow rule pre-approves it. The list
// follows "Protected paths" in Claude Code's documentation as of Claude Code
// 2.1.290 (ADR-0106); permission_prompt and the acceptance preview read it.
package protected

import (
	"path/filepath"
	"slices"
	"strings"
)

// folders are the protected folders: every path below one is protected.
// ".config/git" is the two folders in order.
var folders = []string{".git", ".config/git", ".vscode", ".idea", ".husky", ".cargo", ".devcontainer", ".yarn", ".mvn", ".claude"}

// files are the protected files, by base name, wherever they lie.
var files = []string{
	".gitconfig", ".gitmodules",
	".bashrc", ".bash_profile", ".bash_login", ".bash_aliases", ".bash_logout", ".zshrc", ".zprofile", ".zshenv", ".zlogin", ".zlogout", ".profile", ".envrc",
	".npmrc", ".yarnrc", ".yarnrc.yml", ".pnp.cjs", ".pnp.loader.mjs", ".pnpmfile.cjs", "bunfig.toml", ".bunfig.toml",
	".bazelrc", ".bazelversion", ".bazeliskrc",
	".pre-commit-config.yaml", "lefthook.yml", "lefthook.yaml", ".lefthook.yml", ".lefthook.yaml",
	"gradle-wrapper.properties", "maven-wrapper.properties",
	".devcontainer.json",
	".ripgreprc", "pyrightconfig.json",
	".mcp.json", ".claude.json",
}

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

// Path reports whether rel, a path relative to a worktree or repository root,
// is one Claude Code protects: one of its folders is a protected folder, its
// base name is a protected file, or it is or lies in .git.
func Path(rel string) bool {
	names := parts(rel)
	if len(names) == 0 {
		return false
	}
	if Git(rel) || slices.Contains(files, names[len(names)-1]) {
		return true
	}
	dirs := names[:len(names)-1]
	for _, f := range folders {
		seq := strings.Split(f, "/")
		for i := 0; i+len(seq) <= len(dirs); i++ {
			if slices.Equal(dirs[i:i+len(seq)], seq) {
				return true
			}
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

// Changed returns the paths of the list Claude Code protects, in the list's
// order: for a story's branch, the changed files only the operator accepts.
func Changed(paths []string) []string {
	var out []string
	for _, p := range paths {
		if Path(p) {
			out = append(out, p)
		}
	}
	return out
}
