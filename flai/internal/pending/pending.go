// Package pending answers one question offline and with no credential: has
// the main checkout accepted something that has not been pushed (S-0063,
// ADR-0026)? It compares the checked-out branch with its remote-tracking
// branch, which git updates on every fetch and push from this clone. It
// cannot see a push made from another clone until someone fetches here, and
// callers say so where they show the answer.
package pending

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// Unpushed is what the branch holds that its remote-tracking branch does not.
type Unpushed struct {
	Branch   string `json:"branch"`
	Upstream string `json:"upstream"` // the remote-tracking branch, e.g. origin/main
	Remote   string `json:"remote"`   // its remote, e.g. origin
	Commits  int    `json:"commits"`  // commits ahead of the remote-tracking branch
	// Acceptances are the items accepted by commits that are ahead.
	Acceptances []string `json:"acceptances"`
	// Tags are the tags reachable from the branch and not from the
	// remote-tracking branch: the release tags of unpushed acceptances.
	Tags []string `json:"tags"`
	// Behind is how many commits the remote-tracking branch has that the
	// branch lacks; above zero the two have diverged and a push would fail.
	Behind int `json:"behind,omitempty"`
	// Command is what pushes it from a shell on the host.
	Command string `json:"command"`
}

// acceptance is the subject flai accept commits with.
var acceptance = regexp.MustCompile(`^chore: \[([EST]-\d+)\] accept and archive`)

func lines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// Detect returns nil when there is nothing to say: no git, a detached HEAD,
// no upstream, or nothing ahead.
func Detect(r execx.Runner, root string) *Unpushed {
	u := tracking(r, root)
	if u == nil {
		return nil
	}
	branch, upstream := u.Branch, u.Upstream
	count, err := r.Run(root, "git", "rev-list", "--count", upstream+".."+branch)
	ahead, _ := strconv.Atoi(strings.TrimSpace(count))
	if err != nil || ahead == 0 {
		return nil
	}
	u.Commits = ahead
	if subjects, err := r.Run(root, "git", "log", "--format=%s", upstream+".."+branch); err == nil {
		for _, s := range lines(subjects) {
			if m := acceptance.FindStringSubmatch(s); m != nil {
				u.Acceptances = append(u.Acceptances, m[1])
			}
		}
	}
	if tags, err := r.Run(root, "git", "tag", "--merged", branch, "--no-merged", upstream); err == nil {
		u.Tags = append(u.Tags, lines(tags)...)
	}
	if behind, err := r.Run(root, "git", "rev-list", "--count", branch+".."+upstream); err == nil {
		u.Behind, _ = strconv.Atoi(strings.TrimSpace(behind))
	}
	return u
}

// tracking is the checked-out branch with its remote-tracking branch and
// remote, and nothing yet counted; nil for no git, a detached HEAD, or no
// upstream.
func tracking(r execx.Runner, root string) *Unpushed {
	if r == nil {
		return nil
	}
	if _, err := r.LookPath("git"); err != nil {
		return nil
	}
	branch, err := r.Run(root, "git", "rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.TrimSpace(branch)
	if err != nil || branch == "" || branch == "HEAD" {
		return nil
	}
	upstream, err := r.Run(root, "git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", branch+"@{upstream}")
	upstream = strings.TrimSpace(upstream)
	if err != nil || upstream == "" || !strings.Contains(upstream, "/") {
		// A branch pushed without -u has no upstream configured, but origin's
		// remote-tracking branch of the same name says the same thing.
		upstream = "origin/" + branch
		if _, err := r.Run(root, "git", "rev-parse", "--verify", "--quiet", "refs/remotes/"+upstream); err != nil {
			return nil
		}
	}
	remote, _, _ := strings.Cut(upstream, "/")
	return &Unpushed{Branch: branch, Upstream: upstream, Remote: remote, Acceptances: []string{}, Tags: []string{}, Command: "flai push --pending"}
}

// TagsOnly is a push of tags just made on commits the remote already has
// (S-0144): a release published after its acceptance was pushed tags a
// commit that is not ahead, so Detect, which finds tags through the commits
// ahead, has nothing to say. Nil when there are no tags or no upstream.
func TagsOnly(r execx.Runner, root string, tags []string) *Unpushed {
	if len(tags) == 0 {
		return nil
	}
	u := tracking(r, root)
	if u == nil {
		return nil
	}
	u.Tags = append(u.Tags, tags...)
	return u
}

// Pushes are the pushes that send what is unpushed: Batches of the branch
// and its tags, or with no commits ahead the tags alone, TagsPerPush at a
// time, since a branch that is not ahead has nothing to send and one that is
// behind would be refused.
func (u *Unpushed) Pushes() [][]string {
	if u.Commits > 0 {
		return Batches(u.Branch, u.Tags)
	}
	var out [][]string
	for tags := u.Tags; len(tags) > 0; {
		n := min(len(tags), TagsPerPush)
		out = append(out, append([]string{}, tags[:n]...))
		tags = tags[n:]
	}
	return out
}

// Pending reports whether what is ahead includes an acceptance: only then is
// it flai's business to say so or to push it.
func (u *Unpushed) Pending() bool { return u != nil && len(u.Acceptances) > 0 }

// Refs are what a push sends: the branch and the tags on unpushed commits.
func (u *Unpushed) Refs() []string { return append([]string{u.Branch}, u.Tags...) }

// TagsPerPush is the most tags one push may carry: GitHub creates no events
// for tags when more than three arrive at once, so a workflow that releases on
// a tag never runs (I-0026).
const TagsPerPush = 3

// Batches splits a push into the pushes it is made of: the tags TagsPerPush
// at a time and the branch with the last of them. The branch goes last
// because Detect finds unpushed tags through it; were it pushed first, a push
// interrupted halfway would leave tags nothing reports.
func Batches(branch string, tags []string) [][]string {
	var out [][]string
	for len(tags) > TagsPerPush {
		out = append(out, append([]string{}, tags[:TagsPerPush]...))
		tags = tags[TagsPerPush:]
	}
	return append(out, append([]string{branch}, tags...))
}
