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

// RemoteTags is how this clone's release tags stand against its remote's
// (S-0174). Pending works from local tags only, and flai never fetches: a
// clone whose branch moved without the tags published from another clone
// sees every item since its last local tag as unpublished, and would tag
// versions the remote already has.
type RemoteTags struct {
	Remote string `json:"remote"`
	// Behind are the code components whose highest release tag on the
	// remote is newer than the highest here.
	Behind []Lag `json:"behind,omitempty"`
	// Unchecked says why the remote could not be asked; Behind is then
	// unknown, not empty.
	Unchecked string `json:"unchecked,omitempty"`
	// Fix is the command that brings the remote's tags here.
	Fix     string `json:"fix"`
	Message string `json:"message"`
}

// Lag is one component whose release tags here lag the remote's.
type Lag struct {
	Component string `json:"component"`
	Local     string `json:"local,omitempty"` // the highest tag here; "" for none
	Remote    string `json:"remote"`          // the highest tag on the remote
}

// Lagging reports whether the remote has a release tag newer than this
// clone's for some component.
func (t *RemoteTags) Lagging() bool { return t != nil && len(t.Behind) > 0 }

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
		return fmt.Sprintf("conflict: refusing to publish: could not ask %s for its release tags (%s), so whether it already has these versions is unknown, and the push needs it anyway. Run this again when %s can be reached", t.Remote, t.Unchecked, t.Remote)
	}
}

// RemoteTTL is how long a remote's tags are kept per root: the dashboard
// asks for the pending plan on every item change, and a kept answer can only
// miss a tag published in the meantime, never invent one. Zero asks the
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
}

// CheckRemote asks the remote the checked-out branch tracks, else origin,
// for its release tags, and compares each code component's highest with the
// highest here. Nil when there is nothing to say: no code component, no
// remote, or every component in step. A template is not compared: its
// version is in its template.yaml, not a tag here.
func CheckRemote(r execx.Runner, root string, m manifest.Manifest) *RemoteTags {
	var code []manifest.Project
	for _, p := range m.Projects {
		if p.Kind != "template" {
			code = append(code, p)
		}
	}
	if len(code) == 0 {
		return nil
	}
	remote := remoteOf(r, root)
	if remote == "" {
		return nil
	}
	t := &RemoteTags{Remote: remote, Fix: "git fetch --tags " + remote}
	theirs, err := remoteTags(r, root, remote)
	if err == nil {
		var ours string
		if ours, err = r.Run(root, "git", "tag", "--list"); err == nil {
			for _, p := range code {
				have, want := highestTag(tagLines(ours), p.Name), highestTag(theirs, p.Name)
				if !less(have, want) {
					continue
				}
				lag := Lag{Component: p.Name, Remote: p.Name + "/v" + want.String()}
				if have != (Version{}) {
					lag.Local = p.Name + "/v" + have.String()
				}
				t.Behind = append(t.Behind, lag)
			}
		}
	}
	switch {
	case err != nil:
		t.Unchecked = reason(err)
		t.Message = fmt.Sprintf("could not ask %s for its release tags (%s): what is pending is worked out from this clone's tags alone and may include what is already published", remote, t.Unchecked)
	case len(t.Behind) == 0:
		return nil
	default:
		parts := make([]string, len(t.Behind))
		for i, l := range t.Behind {
			here := l.Local
			if here == "" {
				here = "none"
			}
			parts[i] = fmt.Sprintf("%s (here %s)", l.Remote, here)
		}
		t.Message = fmt.Sprintf("this clone is missing release tags %s has: %s. What looks unpublished may already be published; fetch the tags with %s, then look again", remote, strings.Join(parts, ", "), t.Fix)
	}
	return t
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

// remoteTags is every tag name on the remote, kept for RemoteTTL.
func remoteTags(r execx.Runner, root, remote string) ([]string, error) {
	key := root + "\x00" + remote
	remoteKept.Lock()
	a, ok := remoteKept.answers[key]
	remoteKept.Unlock()
	if ok && time.Since(a.at) < RemoteTTL {
		return a.tags, nil
	}
	type result struct {
		out string
		err error
	}
	ch := make(chan result, 1)
	go func() {
		out, err := r.Run(root, "git", "ls-remote", "--tags", "--refs", remote)
		ch <- result{out, err}
	}()
	var res result
	select {
	case res = <-ch:
	case <-time.After(remoteTimeout):
		return nil, fmt.Errorf("%s did not answer within %s", remote, remoteTimeout)
	}
	if res.err != nil {
		return nil, res.err
	}
	var tags []string
	for _, l := range tagLines(res.out) {
		if _, ref, ok := strings.Cut(l, "\t"); ok {
			tags = append(tags, strings.TrimPrefix(ref, "refs/tags/"))
		}
	}
	remoteKept.Lock()
	remoteKept.answers[key] = remoteAnswer{at: time.Now(), tags: tags}
	remoteKept.Unlock()
	return tags, nil
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
