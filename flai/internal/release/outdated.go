package release

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// Outdated says that the running flai is older than the newest flai release
// tagged in the history of the project it serves (S-0181): it may not have
// rules or fields the project already uses.
type Outdated struct {
	Running string `json:"running"`
	Newest  string `json:"newest"`
	// NewestRules is the newest release's task-done rules number, the
	// buildinfo.TaskDoneRulesPath file at the commit its tag tags, or 0
	// when that commit has none (ADR-0180).
	NewestRules int    `json:"newest_rules,omitempty"`
	Upgrade     string `json:"upgrade"`
	Message     string `json:"message"`
}

// outdatedFor is how long an answer of FlaiOutdated is kept per project: the
// MCP inbox and the dashboard ask often, and a release lands rarely.
const outdatedFor = time.Minute

type outdatedAnswer struct {
	at  time.Time
	out *Outdated
}

var outdatedSeen sync.Map // root and running version → outdatedAnswer

// FlaiOutdated compares the running flai's version with the newest flai/vX.Y.Z
// tag reachable from the project's HEAD, and says so when it is older, with
// that tag's task-done rules number. nil when it is not, when running is not
// a release (a dev build), or when git cannot say. An answer git gave is kept
// for a minute; a look whose git failed is not kept, so the next call asks
// git again rather than serving "not outdated" for a minute (I-0252).
func FlaiOutdated(r execx.Runner, root, running string) *Outdated {
	if _, ok := buildinfo.Semver(running); !ok || r == nil {
		return nil
	}
	key := root + "\x00" + running
	if a, ok := outdatedSeen.Load(key); ok && time.Since(a.(outdatedAnswer).at) < outdatedFor {
		return a.(outdatedAnswer).out
	}
	out, answered := flaiOutdated(r, root, running)
	if answered {
		outdatedSeen.Store(key, outdatedAnswer{at: time.Now(), out: out})
	}
	return out
}

// flaiOutdated is FlaiOutdated's look, uncached; answered is false when git
// could not list the project's flai tags.
func flaiOutdated(r execx.Runner, root, running string) (out *Outdated, answered bool) {
	tags, err := r.Run(root, "git", "tag", "--merged", "HEAD", "--list", "flai/v*")
	if err != nil {
		return nil, false
	}
	newest := highestTag(strings.Split(strings.TrimSpace(tags), "\n"), "flai")
	if newest == (Version{}) || !buildinfo.Below(running, newest.String()) {
		return nil, true
	}
	return &Outdated{
		Running:     strings.TrimPrefix(running, "v"),
		Newest:      newest.String(),
		NewestRules: tagRules(r, root, "flai/v"+newest.String()),
		Upgrade:     buildinfo.UpgradeCommand,
		Message: fmt.Sprintf("this flai is %s, older than flai %s in the project's history, so it may lack rules and fields the project uses: upgrade it with %s",
			strings.TrimPrefix(running, "v"), newest, buildinfo.UpgradeCommand),
	}, true
}

// tagRules is the task-done rules number at the commit tag tags, read from
// buildinfo.TaskDoneRulesPath there (ADR-0180), or 0 when that commit has no
// such file or it holds no positive integer: a release from before the
// number, which keeps ADR-0142's minor-version rule.
func tagRules(r execx.Runner, root, tag string) int {
	text, err := r.Run(root, "git", "show", tag+":"+buildinfo.TaskDoneRulesPath)
	if err != nil {
		return 0
	}
	n, _ := buildinfo.ParseTaskDoneRules(text)
	return n
}

// Behind says that main has commits under the flai sub-project's path that
// the running flai was not built from (ADR-0311).
type Behind struct {
	// Running is the running flai's version, without a leading v.
	Running string `json:"running"`
	// Commit is its build commit as buildinfo gives it, "" when unknown.
	Commit string `json:"commit,omitempty"`
	// By names the point compared from: "commit" (its build commit), "tag"
	// (its own release tag), or "newest" (the newest flai/v* tag merged in main).
	By string `json:"by"`
	// From is the commit or tag compared from.
	From string `json:"from"`
	// Main is the full hash of the main ref's commit compared to.
	Main string `json:"main"`
	// Path is the sub-project path counted under, "" for the whole tree.
	Path string `json:"path,omitempty"`
	// Commits counts the commits reachable from Main and not from From that
	// touch Path.
	Commits int `json:"commits"`
	// Newest is the newest flai/vX.Y.Z tag merged in main, as X.Y.Z, "" when none.
	Newest string `json:"newest,omitempty"`
	// Message is one sentence naming the running flai, the count, From, and the remedy.
	Message string `json:"message"`
}

// behindRemedy is what brings the running flai up to main's flai sub-project.
const behindRemedy = "publish the release with flai release --pending and install it with flai host upgrade"

// FlaiBehind says whether main, in the repository at root, has commits under
// path ("" for the whole tree) that the running flai was not built from,
// comparing from its build commit, else its own release tag, else the newest
// flai/v* tag merged in main (ADR-0311). nil when there are none, or when git
// cannot place the running flai at all; an error when git cannot resolve
// main, list its flai tags, or count its commits.
func FlaiBehind(r execx.Runner, root, main, path string, running buildinfo.Info) (*Behind, error) {
	out, err := r.Run(root, "git", "rev-parse", "--verify", main+"^{commit}")
	if err != nil {
		return nil, fmt.Errorf("resolve %q to a commit in %s to see whether the running flai is behind it: %w; check that %s is the main checkout and holds %q", main, root, err, root, main)
	}
	mainCommit := strings.TrimSpace(out)
	tags, err := r.Run(root, "git", "tag", "--merged", mainCommit, "--list", "flai/v*")
	if err != nil {
		return nil, fmt.Errorf("list the flai/v* tags merged in %s in %s: %w; check that git can read the repository's tags there", main, root, err)
	}
	newest := highestTag(strings.Split(strings.TrimSpace(tags), "\n"), "flai")
	by, from := behindFrom(r, root, running, newest)
	if by == "" {
		return nil, nil
	}
	args := []string{"rev-list", "--count", from + ".." + mainCommit}
	if path != "" {
		args = append(args, "--", path)
	}
	out, err = r.Run(root, "git", args...)
	if err != nil {
		return nil, fmt.Errorf("count the commits in %s from %s to %s in %s: %w; check that both are commits git can read there", main, from, shortHash(mainCommit), root, err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return nil, fmt.Errorf("count the commits in %s from %s in %s: git rev-list answered %q, not a number; check the git on the host", main, from, root, out)
	}
	if n == 0 {
		return nil, nil
	}
	b := &Behind{
		Running: strings.TrimPrefix(running.Version, "v"),
		Commit:  knownCommit(running.Commit),
		By:      by,
		From:    from,
		Main:    mainCommit,
		Path:    path,
		Commits: n,
	}
	if newest != (Version{}) {
		b.Newest = newest.String()
	}
	b.Message = behindMessage(b)
	return b, nil
}

// behindFrom picks the point FlaiBehind compares from and how it was found,
// or "" when git cannot place the running flai.
func behindFrom(r execx.Runner, root string, running buildinfo.Info, newest Version) (by, from string) {
	if c := knownCommit(running.Commit); c != "" {
		c = strings.TrimSuffix(strings.TrimSuffix(c, "-dirty"), "+dirty")
		if _, err := r.Run(root, "git", "cat-file", "-e", c+"^{commit}"); err == nil {
			return "commit", c
		}
	}
	if v, ok := buildinfo.Semver(running.Version); ok {
		tag := fmt.Sprintf("flai/v%d.%d.%d", v[0], v[1], v[2])
		if _, err := r.Run(root, "git", "rev-parse", "--verify", "--quiet", "refs/tags/"+tag+"^{commit}"); err == nil {
			return "tag", tag
		}
	}
	if newest != (Version{}) {
		return "newest", "flai/v" + newest.String()
	}
	return "", ""
}

// knownCommit is a build commit as buildinfo gives it, or "" when it gives
// none or "unknown".
func knownCommit(c string) string {
	if c == "unknown" {
		return ""
	}
	return c
}

// shortHash is a commit hash's first seven characters.
func shortHash(c string) string {
	if len(c) > 7 {
		return c[:7]
	}
	return c
}

// behindMessage is the one sentence that says how far b's flai is behind and
// what to do.
func behindMessage(b *Behind) string {
	where := "in the tree"
	if b.Path != "" {
		where = "under " + strings.TrimSuffix(b.Path, "/") + "/"
	}
	count := fmt.Sprintf("%d commits", b.Commits)
	if b.Commits == 1 {
		count = "1 commit"
	}
	var since string
	switch b.By {
	case "commit":
		since = fmt.Sprintf("flai %s, built from commit %s, lacks %s main has %s since that commit", b.Running, shortHash(b.From), count, where)
	case "tag":
		since = fmt.Sprintf("flai %s, whose build commit git cannot place, lacks %s main has %s since its release tag %s", b.Running, count, where, b.From)
	default:
		since = fmt.Sprintf("flai %s, which git cannot place by build commit or release tag, lacks %s main has %s since the newest release tag %s", b.Running, count, where, b.From)
	}
	return since + ": " + behindRemedy
}
