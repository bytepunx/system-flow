package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Limits on what a diff carries, so a large branch cannot swamp a browser.
const (
	diffFileLimit  = 64 * 1024
	diffTotalLimit = 768 * 1024
)

// diffFile is one file a story branch changed.
type diffFile struct {
	Path      string `json:"path"`
	OldPath   string `json:"old_path,omitempty"` // for a rename
	Status    string `json:"status"`             // added, modified, deleted, renamed
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Binary    bool   `json:"binary"`
	Truncated bool   `json:"truncated"` // the patch was cut at the limit
	Patch     string `json:"patch"`     // unified hunks, empty for a binary file
}

// streamDiff is a story branch against the main branch.
type streamDiff struct {
	Story     string     `json:"story"`
	Branch    string     `json:"branch"`
	Base      string     `json:"base"` // the merge base, abbreviated
	Commits   int        `json:"commits"`
	Files     []diffFile `json:"files"`
	Additions int        `json:"additions"`
	Deletions int        `json:"deletions"`
	Truncated bool       `json:"truncated"` // some patches were cut or left out
}

func newStreamDiffCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "diff <story-id>",
		Short: "What a story branch changes: files and hunks against the main branch",
		Long: `The story's branch (ADR-0019) against its merge base with the main
branch, read with git in the main checkout, so it does not depend on the
worktree. The dashboard's review page shows this (S-0041).`,
		Example: `  flai stream diff S-0041
  flai stream diff S-0041 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			it, err := repo.Get(args[0])
			if err != nil {
				return err
			}
			d, err := a.storyDiff(repo, it.ID)
			if err != nil {
				return err
			}
			if a.jsonOut {
				return a.printJSON(d)
			}
			fmt.Fprintf(a.out, "%s against %s: %d commit(s), %d file(s), +%d -%d\n", d.Branch, d.Base, d.Commits, len(d.Files), d.Additions, d.Deletions)
			for _, f := range d.Files {
				name := f.Path
				if f.OldPath != "" {
					name = f.OldPath + " -> " + f.Path
				}
				fmt.Fprintf(a.out, "  %-9s %-60s +%d -%d\n", f.Status, name, f.Additions, f.Deletions)
			}
			return nil
		},
	}
}

func (a *app) storyDiff(repo *workitem.Repo, id string) (*streamDiff, error) {
	root := repo.MainRoot
	if root == "" || !a.inGitWorkTree(root) {
		return nil, fmt.Errorf("%s has no branch to compare: this project is not a git repository", id)
	}
	branch := storyBranch(id)
	if !a.branchExists(root, branch) {
		return nil, fmt.Errorf("%s has no branch %s: it was never opened with flai stream open, or it has been merged and removed", id, branch)
	}
	main, err := a.mainBranch(root)
	if err != nil {
		return nil, err
	}
	base, err := a.runner.Run(root, "git", "merge-base", main, branch)
	if err != nil {
		return nil, err
	}
	d := &streamDiff{Story: id, Branch: branch, Files: []diffFile{}}
	if short, err := a.runner.Run(root, "git", "rev-parse", "--short", base); err == nil {
		d.Base = short
	}
	if n, err := a.runner.Run(root, "git", "rev-list", "--count", base+".."+branch); err == nil {
		d.Commits, _ = strconv.Atoi(n)
	}
	status, err := a.runner.Run(root, "git", "diff", "--name-status", "-M", "-z", base, branch)
	if err != nil {
		return nil, err
	}
	byPath := map[string]*diffFile{}
	fields := strings.Split(strings.TrimRight(status, "\x00"), "\x00")
	for i := 0; i < len(fields) && fields[i] != ""; i++ {
		f := diffFile{}
		switch code := fields[i]; {
		case strings.HasPrefix(code, "R") && i+2 < len(fields):
			f.Status, f.OldPath, f.Path = "renamed", fields[i+1], fields[i+2]
			i += 2
		case i+1 < len(fields):
			f.Status = map[byte]string{'A': "added", 'D': "deleted", 'M': "modified"}[code[0]]
			if f.Status == "" {
				f.Status = "modified"
			}
			f.Path = fields[i+1]
			i++
		default:
			continue
		}
		d.Files = append(d.Files, f)
	}
	for i := range d.Files {
		byPath[d.Files[i].Path] = &d.Files[i]
	}
	numstat, err := a.runner.Run(root, "git", "diff", "--numstat", "-M", "-z", base, branch)
	if err != nil {
		return nil, err
	}
	// records are "add\tdel\tpath\0", or for a rename "add\tdel\t\0old\0new\0"
	recs := strings.Split(strings.TrimRight(numstat, "\x00"), "\x00")
	for i := 0; i < len(recs); i++ {
		parts := strings.SplitN(recs[i], "\t", 3)
		if len(parts) != 3 {
			continue
		}
		path := parts[2]
		if path == "" && i+2 < len(recs) {
			path = recs[i+2]
			i += 2
		}
		f := byPath[path]
		if f == nil {
			continue
		}
		if parts[0] == "-" {
			f.Binary = true
			continue
		}
		f.Additions, _ = strconv.Atoi(parts[0])
		f.Deletions, _ = strconv.Atoi(parts[1])
		d.Additions += f.Additions
		d.Deletions += f.Deletions
	}
	total := 0
	for i := range d.Files {
		f := &d.Files[i]
		if f.Binary {
			continue
		}
		if total >= diffTotalLimit {
			f.Truncated, d.Truncated = true, true
			continue
		}
		paths := []string{f.Path}
		if f.OldPath != "" {
			paths = append(paths, f.OldPath)
		}
		patch, err := a.runner.Run(root, "git", append([]string{"diff", "--no-color", "-M", base, branch, "--"}, paths...)...)
		if err != nil {
			return nil, err
		}
		// keep the hunks; the header lines repeat what the fields say
		if at := strings.Index(patch, "\n@@"); at >= 0 {
			patch = patch[at+1:]
		} else if !strings.HasPrefix(patch, "@@") {
			patch = ""
		}
		if len(patch) > diffFileLimit {
			patch = patch[:strings.LastIndex(patch[:diffFileLimit], "\n")+1]
			f.Truncated, d.Truncated = true, true
		}
		f.Patch = patch
		total += len(patch)
	}
	return d, nil
}
