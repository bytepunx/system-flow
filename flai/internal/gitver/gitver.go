// Package gitver reads the version of the git on PATH, for the few places
// where flai's behaviour depends on it (ADR-0022).
package gitver

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// RelativeWorktrees is the first git that understands relative worktree
// links and the extensions.relativeWorktrees repository extension.
var RelativeWorktrees = Version{Major: 2, Minor: 48}

// MergeTree is the first git whose merge-tree has --write-tree, the trial
// merge flai stream sync runs against other story branches (S-0131).
var MergeTree = Version{Major: 2, Minor: 38}

// Version is a git release, major and minor; patch levels never matter here.
type Version struct{ Major, Minor int }

func (v Version) String() string { return fmt.Sprintf("%d.%d", v.Major, v.Minor) }

// AtLeast reports whether v is the same release as min or a later one.
func (v Version) AtLeast(min Version) bool {
	return v.Major > min.Major || (v.Major == min.Major && v.Minor >= min.Minor)
}

var versionRE = regexp.MustCompile(`(\d+)\.(\d+)`)

// Parse reads the output of "git version", for example
// "git version 2.47.3" or "git version 2.39.3 (Apple Git-146)".
func Parse(out string) (Version, error) {
	m := versionRE.FindStringSubmatch(out)
	if m == nil {
		return Version{}, fmt.Errorf("no version in %q", out)
	}
	major, _ := strconv.Atoi(m[1])
	minor, _ := strconv.Atoi(m[2])
	return Version{Major: major, Minor: minor}, nil
}

// Installed asks the git on PATH for its version.
func Installed(r execx.Runner) (Version, error) {
	out, err := r.Run("", "git", "version")
	if err != nil {
		return Version{}, fmt.Errorf("git version: %w", err)
	}
	return Parse(out)
}
