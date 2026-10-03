package release

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/pending"
)

// RemoteTags is how this clone's release tags and branch stand against its
// remote's. Pending works from local tags only, and flai never fetches: a
// clone whose branch moved without the tags published from another clone
// sees every item since its last local tag as unpublished, and would tag
// versions the remote already has (S-0174); and a clone whose remote branch
// has commits it lacks would tag and commit what its push cannot send
// (ADR-0067).
type RemoteTags struct {
	Remote string `json:"remote"`
	// Behind are the code components whose highest release tag on the
	// remote is newer than the highest here.
	Behind []Lag `json:"behind,omitempty"`
	// Branch is set when the remote's head of the branch this clone's
	// branch tracks is not in this clone's history.
	Branch *BranchLag `json:"branch,omitempty"`
	// Unchecked says why the remote could not be asked; Behind and Branch
	// are then unknown, not empty.
	Unchecked string `json:"unchecked,omitempty"`
	// Fix is the command that brings this clone in step with the remote.
	Fix     string `json:"fix"`
	Message string `json:"message"`
}

// BranchLag is the remote branch this clone's branch tracks, when it has
// commits this clone lacks.
type BranchLag struct {
	Upstream string `json:"upstream"` // the remote-tracking branch, such as origin/main
	Head     string `json:"head"`     // the remote's commit at the head of that branch
	// Fetched says this clone has Head, fetched but not merged; false means
	// it has to be fetched first.
	Fetched bool `json:"fetched"`
}

// Lag is one component whose release tags here lag the remote's.
type Lag struct {
	Component string `json:"component"`
	Local     string `json:"local,omitempty"` // the highest tag here; "" for none
	Remote    string `json:"remote"`          // the highest tag on the remote
}

// Lagging reports whether the remote has a release tag newer than this
// clone's for some component, or commits on its branch this clone lacks.
func (t *RemoteTags) Lagging() bool { return t != nil && (len(t.Behind) > 0 || t.Branch != nil) }

// Refusal is why a publish is refused, "" when nothing stands in its way: a
// lagging clone would tag versions already published, and one that cannot
// ask its remote does not know whether it would, and could not push them.
func (t *RemoteTags) Refusal() string {
	switch {
	case t == nil:
		return ""
	case t.Lagging():
		return "conflict: refusing to publish: " + t.Message
	default:
		return fmt.Sprintf("conflict: refusing to publish: could not ask %s for its release tags and branch (%s), so whether it already has these versions is unknown, and the push needs it anyway. Run this again when %s can be reached", t.Remote, t.Unchecked, t.Remote)
	}
}

// RemoteTTL is how long a remote's tags and branch head are kept per root:
// the dashboard asks for the pending plan on every item change, and a kept
// answer can only miss a tag or commit published in the meantime, never
// invent one. Zero asks the
// remote every time, as a test that publishes from another clone needs.
var RemoteTTL = time.Minute

// remoteTimeout bounds how long the remote is waited for. A git that has not
// answered by then is left to finish on its own; nothing waits on it.
var remoteTimeout = 15 * time.Second

var remoteKept = struct {
	sync.Mutex
	answers map[string]remoteAnswer
}{answers: map[string]remoteAnswer{}}

type remoteAnswer struct {
	at   time.Time
	tags []string
	head string // the remote's head of the branch asked for; "" when it has none
}

// CheckRemote asks the remote the checked-out branch tracks, else origin,
// for its release tags and its head of that branch, in one ls-remote. It
// compares each code component's highest tag with the highest here, and
// whether this clone's history holds the remote's head. Nil when there is
// nothing to say: no remote, or in step with it. A template is not compared
// by tag: its version is in its template.yaml, not a tag here. The remote
// lacking the branch is not behind.
func CheckRemote(r execx.Runner, root string, m manifest.Manifest) *RemoteTags {
	remote := remoteOf(r, root)
	if remote == "" {
		return nil
	}
	upstream, branch := upstreamOf(r, root, remote)
	t := &RemoteTags{Remote: remote}
	theirs, err := askRemote(r, root, remote, branch)
	if err == nil {
		t.Behind, err = tagsBehind(r, root, m, theirs.tags)
	}
	if err != nil {
		t.Unchecked = reason(err)
		t.Fix = "git fetch --tags " + remote
		t.Message = fmt.Sprintf("could not ask %s for its release tags and branch (%s): what is pending is worked out from this clone's tags alone and may include what is already published", remote, t.Unchecked)
		return t
	}
	t.Branch = branchBehind(r, root, upstream, theirs.head)
	if !t.Lagging() {
		return nil
	}
	t.Fix, t.Message = advice(t)
	return t
}

// CheckRemoteAfresh is CheckRemote asking the remote again, whatever answer
// is kept for root: after a push the remote refused, the answer kept from the
// check before tagging is the one that let the publish through (S-0242).
func CheckRemoteAfresh(r execx.Runner, root string, m manifest.Manifest) *RemoteTags {
	remoteKept.Lock()
	for key := range remoteKept.answers {
		if strings.HasPrefix(key, root+"\x00") {
			delete(remoteKept.answers, key)
		}
	}
	remoteKept.Unlock()
	return CheckRemote(r, root, m)
}

// IsReleaseTag reports whether tag is one flai release makes, <name>/vX.Y.Z
// for a code component in m, rather than a tag of the user's own.
func IsReleaseTag(m manifest.Manifest, tag string) bool {
	for _, p := range m.Projects {
		if p.Kind == "template" {
			continue
		}
		if v, ok := strings.CutPrefix(tag, p.Name+"/v"); ok {
			if _, ok := ParseVersion(v); ok {
				return true
			}
		}
	}
	return false
}

// tagsBehind are the code components whose highest release tag among
// theirs is newer than the highest here.
func tagsBehind(r execx.Runner, root string, m manifest.Manifest, theirs []string) ([]Lag, error) {
	var code []manifest.Project
	for _, p := range m.Projects {
		if p.Kind != "template" {
			code = append(code, p)
		}
	}
	if len(code) == 0 {
		return nil, nil
	}
	ours, err := r.Run(root, "git", "tag", "--list")
	if err != nil {
		return nil, err
	}
	var behind []Lag
	for _, p := range code {
		have, want := highestTag(tagLines(ours), p.Name), highestTag(theirs, p.Name)
		if !less(have, want) {
			continue
		}
		lag := Lag{Component: p.Name, Remote: p.Name + "/v" + want.String()}
		if have != (Version{}) {
			lag.Local = p.Name + "/v" + have.String()
		}
		behind = append(behind, lag)
	}
	return behind, nil
}

// branchBehind is the remote's head of upstream when this clone's history
// lacks it, nil when the remote has no such branch or HEAD already holds it.
func branchBehind(r execx.Runner, root, upstream, head string) *BranchLag {
	if head == "" {
		return nil
	}
	if _, err := r.Run(root, "git", "merge-base", "--is-ancestor", head, "HEAD"); err == nil {
		return nil
	}
	_, missing := r.Run(root, "git", "cat-file", "-e", head+"^{commit}")
	return &BranchLag{Upstream: upstream, Head: head, Fetched: missing == nil}
}

// advice is the command that brings this clone in step with the remote, and
// the message saying what lags and what to run.
func advice(t *RemoteTags) (fix, message string) {
	var lags, steps []string
	fetch := ""
	if len(t.Behind) > 0 {
		parts := make([]string, len(t.Behind))
		for i, l := range t.Behind {
			here := l.Local
			if here == "" {
				here = "none"
			}
			parts[i] = fmt.Sprintf("%s (here %s)", l.Remote, here)
		}
		lags = append(lags, fmt.Sprintf("this clone is missing release tags %s has: %s. What looks unpublished may already be published", t.Remote, strings.Join(parts, ", ")))
		fetch = "git fetch --tags " + t.Remote
	}
	if b := t.Branch; b != nil {
		state := "not fetched here"
		if b.Fetched {
			state = "fetched here but not merged"
		}
		lags = append(lags, fmt.Sprintf("%s has commits this clone lacks (at %s, %s), so a publish would tag and commit what its push cannot send", b.Upstream, short(b.Head), state))
		if !b.Fetched && fetch == "" {
			fetch = "git fetch " + t.Remote
		}
	}
	if fetch != "" {
		steps = append(steps, fetch)
	}
	if b := t.Branch; b != nil {
		steps = append(steps, "git merge "+b.Upstream)
	}
	fix = strings.Join(steps, " && ")
	if t.Branch == nil {
		return fix, fmt.Sprintf("%s; fetch the tags with %s, then look again", strings.Join(lags, "; and "), fix)
	}
	return fix, fmt.Sprintf("%s. Run %s (or rebase onto %s instead of merging), then flai release --pending again", strings.Join(lags, "; and "), fix, t.Branch.Upstream)
}

// short is a commit's abbreviated name.
func short(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}

// remoteOf is the remote to ask: the one the branch tracks, else origin if
// there is one, else "".
func remoteOf(r execx.Runner, root string) string {
	if remote := pending.Remote(r, root); remote != "" {
		return remote
	}
	if _, err := r.Run(root, "git", "remote", "get-url", "origin"); err == nil {
		return "origin"
	}
	return ""
}

// upstreamOf is the remote-tracking branch the checked-out branch tracks on
// remote, such as origin/main, and that branch's name on the remote: the
// configured upstream, else remote's branch of the same name. Both "" for a
// detached HEAD.
func upstreamOf(r execx.Runner, root, remote string) (upstream, branch string) {
	head, err := r.Run(root, "git", "rev-parse", "--abbrev-ref", "HEAD")
	head = strings.TrimSpace(head)
	if err != nil || head == "" || head == "HEAD" {
		return "", ""
	}
	if up, err := r.Run(root, "git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", head+"@{upstream}"); err == nil {
		if name, ok := strings.CutPrefix(strings.TrimSpace(up), remote+"/"); ok && name != "" {
			return remote + "/" + name, name
		}
	}
	return remote + "/" + head, head
}

// askRemote is every tag name on the remote and its head of branch, kept
// for RemoteTTL. With no branch, only the tags are asked for.
func askRemote(r execx.Runner, root, remote, branch string) (remoteAnswer, error) {
	key := root + "\x00" + remote + "\x00" + branch
	remoteKept.Lock()
	a, ok := remoteKept.answers[key]
	remoteKept.Unlock()
	if ok && time.Since(a.at) < RemoteTTL {
		return a, nil
	}
	args := []string{"ls-remote", "--tags", "--refs", remote}
	if branch != "" {
		args = []string{"ls-remote", "--refs", remote, "refs/tags/*", "refs/heads/" + branch}
	}
	type result struct {
		out string
		err error
	}
	ch := make(chan result, 1)
	go func() {
		out, err := r.Run(root, "git", args...)
		ch <- result{out, err}
	}()
	var res result
	select {
	case res = <-ch:
	case <-time.After(remoteTimeout):
		return remoteAnswer{}, fmt.Errorf("%s did not answer within %s", remote, remoteTimeout)
	}
	if res.err != nil {
		return remoteAnswer{}, res.err
	}
	a = remoteAnswer{at: time.Now()}
	for _, l := range tagLines(res.out) {
		sha, ref, ok := strings.Cut(l, "\t")
		if !ok {
			continue
		}
		if tag, ok := strings.CutPrefix(ref, "refs/tags/"); ok {
			a.tags = append(a.tags, tag)
		} else if branch != "" && ref == "refs/heads/"+branch {
			a.head = sha
		}
	}
	remoteKept.Lock()
	remoteKept.answers[key] = a
	remoteKept.Unlock()
	return a, nil
}

func tagLines(s string) []string {
	var out []string
	for _, l := range strings.Split(s, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

// reason is what a failed git says went wrong: its fatal: line when it
// printed one, else the error's first line.
func reason(err error) string {
	lines := tagLines(err.Error())
	for _, l := range lines {
		if strings.HasPrefix(l, "fatal: ") {
			return strings.TrimPrefix(l, "fatal: ")
		}
	}
	if len(lines) == 0 {
		return err.Error()
	}
	return lines[0]
}
