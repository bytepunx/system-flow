package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/itemedit"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Pushing a story in review under pull-request acceptance (S-0393,
// ADR-0166). Once a minute, for each served project whose manifest sets
// acceptance.mode: pull-request and whose forge host action is on, flai
// serve looks at each story in review. It pushes story/S-nnnn to origin
// with the host's git credentials when origin's branch is missing or not
// the local head: plainly when origin's head is in the local branch, and,
// when a sync rebased the branch, with --force-with-lease naming the head
// flai last pushed, which it keeps per story under .flai-cache/forge. It
// then opens the story's one pull request, found by its head branch, or
// updates it when its title, body, or base is not what flai would write.
// A failure opens one thread on the story, adds an entry to it only when
// the reason changes, and is resolved by the first look that succeeds.
// flai forge push runs the same PushStory on the operator's word.

const (
	// forgeEvery is how often the stories in review are looked at.
	forgeEvery = time.Minute
	// forgeAuthor opens, writes on, and resolves the failure thread.
	forgeAuthor = "flai serve"
	// forgeTimeout bounds the forge calls for one story.
	forgeTimeout = 2 * time.Minute
	// pullRequestTemplate is the project's pull request template, from its root.
	pullRequestTemplate = ".github/pull_request_template.md"
)

// ForgeOutcome is what flai did with a story's branch at its last push.
type ForgeOutcome string

// The outcomes of a push.
const (
	// ForgePushed is a push of the branch, and its pull request opened or updated.
	ForgePushed ForgeOutcome = "pushed"
	// ForgeUpToDate is a branch origin already had, and its pull request opened or updated.
	ForgeUpToDate ForgeOutcome = "up-to-date"
	// ForgeFailed is a push or a request that failed, for the reason given.
	ForgeFailed ForgeOutcome = "failed"
	// ForgeSkipped is nothing done, for the reason given.
	ForgeSkipped ForgeOutcome = "skipped"
)

// The changes a push makes to a story's pull request.
const (
	RequestOpened    = "opened"
	RequestUpdated   = "updated"
	RequestUnchanged = "unchanged"
)

// ForgeState is what flai last did for one story under pull-request
// acceptance, kept in .flai-cache/forge/S-nnnn.json; it never holds the token.
type ForgeState struct {
	// Story is the story's ID.
	Story string `json:"story"`
	// Outcome is what the last push did, and Reason why when it failed or was skipped.
	Outcome ForgeOutcome `json:"outcome"`
	Reason  string       `json:"reason,omitempty"`
	// At is when the last push ran, RFC 3339.
	At string `json:"at"`
	// LocalHead and RemoteHead are story/S-nnnn here and on origin as the
	// last push left them; RemoteHead is empty when origin has no such branch.
	LocalHead  string `json:"local_head,omitempty"`
	RemoteHead string `json:"remote_head,omitempty"`
	// LastPushed is the head flai last pushed, or last found on origin as
	// the local head, which a push after a rebase names as its lease;
	// PushedAt is when it last pushed, and Forced whether that push was
	// with the lease.
	LastPushed string `json:"last_pushed,omitempty"`
	PushedAt   string `json:"pushed_at,omitempty"`
	Forced     bool   `json:"forced,omitempty"`
	// PullRequest and URL are the story's pull request as last opened or
	// found; Request is what the last push did to it: RequestOpened,
	// RequestUpdated, or RequestUnchanged, and empty when it did not reach it.
	PullRequest int    `json:"pull_request,omitempty"`
	URL         string `json:"url,omitempty"`
	Request     string `json:"request,omitempty"`
	// Thread is the thread flai serve opened on the story for a failure,
	// and ThreadReason the reason it last wrote there; empty once a push
	// succeeds.
	Thread       string `json:"thread,omitempty"`
	ThreadReason string `json:"thread_reason,omitempty"`
}

// ForgeStatePath is where flai keeps what it last did for story.
func ForgeStatePath(repo *workitem.Repo, story string) string {
	return filepath.Join(repo.CacheDir(), "forge", story+".json")
}

// ReadForgeState is what flai last did for story, and false with no error
// when it has done nothing for it.
func ReadForgeState(repo *workitem.Repo, story string) (ForgeState, bool, error) {
	path := ForgeStatePath(repo, story)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return ForgeState{}, false, nil
	}
	if err != nil {
		return ForgeState{}, false, fmt.Errorf("read what flai last pushed for %s at %s: %w", story, path, err)
	}
	var st ForgeState
	if err := json.Unmarshal(data, &st); err != nil {
		return ForgeState{}, false, fmt.Errorf("read what flai last pushed for %s at %s: %w; delete it, and flai records the next push afresh", story, path, err)
	}
	return st, true, nil
}

func writeForgeState(repo *workitem.Repo, st ForgeState) error {
	path := ForgeStatePath(repo, st.Story)
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("record what flai pushed for %s: %w; check that the project's .flai-cache is writable", st.Story, err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("record what flai pushed for %s at %s: %w; check that the project's .flai-cache is writable", st.Story, path, err)
	}
	return nil
}

// ForgeToken is how flai reads a forge's token: the variable it is read
// from, named in reasons, and the function that reads it at each call.
type ForgeToken struct {
	Env  string
	Read forge.TokenFunc
}

// NewForge makes the adapter for the forge name, bound to the repository
// the origin URL remote names.
type NewForge func(name, remote string, token forge.TokenFunc) (forge.Forge, error)

// ForgePush is what PushStory pushes one story with.
type ForgePush struct {
	// Repo is the project; Story the story to push.
	Repo  *workitem.Repo
	Story string
	// Git runs git; the system's when nil.
	Git execx.Runner
	// Token reads the token of the forge acceptance.forge names.
	Token ForgeToken
	// NewForge makes the forge's adapter; forge.New with GitHub's API when nil.
	NewForge NewForge
	// HostAction, when set, answers whether the forge host action is on for
	// the project, and the push is skipped when it is not; nil, for the
	// operator's flai forge push, pushes without it.
	HostAction func() bool
	// Now is the clock; time.Now when nil.
	Now func() time.Time
}

// PushStory pushes the story's branch to origin when origin's is missing or
// not the local head, and opens or updates its pull request, under
// pull-request acceptance, and records what it did. A push it skips, for
// the mode, the host action, or a story not in review, answers why in a
// ForgeSkipped state and no error; a push that fails answers its reason as
// the error, and in a ForgeFailed state.
func PushStory(ctx context.Context, p ForgePush) (ForgeState, error) {
	it, err := p.Repo.Get(p.Story)
	if err != nil {
		return ForgeState{}, err
	}
	if it.Type != workitem.Story {
		return ForgeState{}, fmt.Errorf("%s is a %s: flai pushes a story's branch, so name its story", it.ID, it.Type)
	}
	st, _, err := ReadForgeState(p.Repo, it.ID)
	if err != nil {
		return ForgeState{}, err
	}
	st.Story, st.At, st.Reason, st.Request = it.ID, p.now().UTC().Format(time.RFC3339), "", ""
	if why := p.skip(it); why != "" {
		st.Outcome, st.Reason = ForgeSkipped, why
		return st, writeForgeState(p.Repo, st)
	}
	err = p.push(ctx, it, &st)
	if err != nil {
		st.Outcome, st.Reason = ForgeFailed, err.Error()
	}
	if werr := writeForgeState(p.Repo, st); werr != nil {
		return st, errors.Join(err, werr)
	}
	return st, err
}

// skip is why the story is not pushed, or "" when it is.
func (p ForgePush) skip(it *workitem.Item) string {
	switch {
	case !p.Repo.Manifest.Acceptance.PullRequest():
		return fmt.Sprintf("acceptance.mode is %s, so %s is accepted on this machine and flai pushes nothing: set acceptance.mode: pull-request in system-flow.yaml to have stories in review pushed and their pull requests opened",
			p.Repo.Manifest.Acceptance.ModeOrDefault(), it.ID)
	case p.HostAction != nil && !p.HostAction():
		return fmt.Sprintf("the forge host action is off for this project, so flai serve pushes nothing: turn it on with flai serve enable forge, or push one story with flai forge push %s", it.ID)
	case it.Status != workitem.Review:
		return fmt.Sprintf("%s is %s, not in review: flai pushes a story's branch and opens its pull request once the story is in review", it.ID, it.Status)
	}
	return ""
}

func (p ForgePush) now() time.Time {
	if p.Now == nil {
		return time.Now()
	}
	return p.Now()
}

func (p ForgePush) git() execx.Runner {
	if p.Git == nil {
		return execx.System{}
	}
	return p.Git
}

func (p ForgePush) newForge(name, remote string) (forge.Forge, error) {
	if p.NewForge == nil {
		return forge.New(name, remote, p.Token.Read, forge.GitHubAPI)
	}
	return p.NewForge(name, remote, p.Token.Read)
}

// push pushes it's branch and opens or updates its pull request, recording
// in st what it did as it goes, so that a push that succeeded is kept when
// the request after it fails.
func (p ForgePush) push(ctx context.Context, it *workitem.Item, st *ForgeState) error {
	git, root := p.git(), mainCheckout(p.Repo)
	name := p.Repo.Manifest.Acceptance.ForgeOrDefault()
	branch := storygit.Branch(it.ID)
	if p.Token.Read == nil {
		return fmt.Errorf("flai was given no way to read the %s token, so it pushes nothing: start flai serve with flai serve start", name)
	}
	if _, err := p.Token.Read(); err != nil {
		return fmt.Errorf("%s's pull request cannot be opened without the %s token, so flai pushes nothing: %w (flai serve, and flai host that starts it, read it from the environment they started with)", it.ID, name, err)
	}
	url, err := git.Run(root, "git", "remote", "get-url", "origin")
	if err != nil {
		return fmt.Errorf("the project has no origin remote to push %s to (%s): add the forge's repository with git remote add origin, then its URL", branch, gitSaid(url, err))
	}
	f, err := p.newForge(name, strings.TrimSpace(url))
	if err != nil {
		return err
	}
	local, err := git.Run(root, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+branch+"^{commit}")
	if err != nil || local == "" {
		return fmt.Errorf("this repository has no branch %s to push: open it with flai stream open %s, or move %s back to in-progress", branch, it.ID, it.ID)
	}
	st.LocalHead = local
	remote, err := remoteHead(git, root, branch)
	if err != nil {
		return err
	}
	st.RemoteHead = remote
	plan, err := planPush(it.ID, local, remote, st.LastPushed, func() bool {
		_, err := git.Run(root, "git", "merge-base", "--is-ancestor", remote, local)
		return err == nil
	})
	if err != nil {
		return err
	}
	if plan == pushNone {
		st.Outcome, st.LastPushed = ForgeUpToDate, local
	} else {
		if err := pushBranch(git, root, it.ID, local, remote, st.LastPushed, plan); err != nil {
			return err
		}
		st.Outcome, st.LastPushed, st.RemoteHead = ForgePushed, local, local
		st.PushedAt, st.Forced = st.At, plan == pushLease
	}
	return p.request(ctx, f, git, root, it, st)
}

// request opens the story's pull request, or updates the one its branch has
// when its title, body, or base is not what flai would write.
func (p ForgePush) request(ctx context.Context, f forge.Forge, git execx.Runner, root string, it *workitem.Item, st *ForgeState) error {
	base, err := storygit.MainBranch(git, root)
	if err != nil {
		return fmt.Errorf("could not tell the main branch %s's pull request merges into: %w", it.ID, err)
	}
	spr, err := storyPullRequest(p.Repo, root, it)
	if err != nil {
		return err
	}
	in := forge.PullRequestInput{Head: storygit.Branch(it.ID), Base: base, Title: spr.PullRequestTitle()}
	if in.Body, err = spr.PullRequestBody(); err != nil {
		return fmt.Errorf("%w; flai tries again at its next look", err)
	}
	fctx, cancel := context.WithTimeout(ctx, forgeTimeout)
	defer cancel()
	pr, found, err := f.FindPullRequest(fctx, in.Head)
	if err != nil {
		return forgeReason(f, p.Token.Env, it.ID, err)
	}
	change := RequestOpened
	if found {
		change = RequestUpdated
		if sameRequest(pr, in) {
			change = RequestUnchanged
		}
	}
	if change != RequestUnchanged {
		if pr, err = f.OpenPullRequest(fctx, in); err != nil {
			return forgeReason(f, p.Token.Env, it.ID, err)
		}
	}
	st.PullRequest, st.URL, st.Request = pr.Number, pr.URL, change
	return nil
}

// sameRequest is whether pr already has what in would write: the forge may
// end its lines with CRLF and drop the trailing blank lines of a body.
func sameRequest(pr forge.PullRequest, in forge.PullRequestInput) bool {
	norm := func(s string) string { return strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n")) }
	return pr.Title == in.Title && pr.Base == in.Base && norm(pr.Body) == norm(in.Body)
}

// storyPullRequest is what the story's pull request is built from: its
// criteria, its last flai verify result, and the project's template.
func storyPullRequest(repo *workitem.Repo, root string, it *workitem.Item) (forge.StoryPullRequest, error) {
	spr := forge.StoryPullRequest{ID: it.ID, Title: it.Title}
	for _, c := range itemedit.Criteria(it.Body) {
		spr.Criteria = append(spr.Criteria, forge.Criterion{Text: c.Text, Ticked: c.Ticked})
	}
	rep, ok, err := verify.LastReport(repo, it.ID)
	if err != nil {
		return spr, err
	}
	if ok {
		spr.Verify = &forge.VerifyResult{Passed: rep.Passed, StoppedAt: rep.StoppedAt, Commit: rep.Commit, RanAt: rep.RanAt}
	}
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(pullRequestTemplate)))
	switch {
	case errors.Is(err, fs.ErrNotExist):
	case err != nil:
		return spr, fmt.Errorf("read the project's pull request template %s: %w; check that it is readable", pullRequestTemplate, err)
	default:
		spr.Template = string(data)
	}
	return spr, nil
}

// pushKind is how a story's branch is pushed.
type pushKind int

const (
	// pushNone pushes nothing: origin has the local head.
	pushNone pushKind = iota
	// pushPlain pushes as a fast-forward, or creates origin's branch.
	pushPlain
	// pushLease pushes over origin's branch, which a sync rebased, with a
	// lease on the head flai last pushed.
	pushLease
)

// planPush decides how story's branch, at local here and remote on origin
// (empty when origin has none), is pushed; isAncestor says whether remote
// is in the local branch, asked only when the heads differ. A branch whose
// origin head is not in it and that flai never pushed is refused: flai
// cannot tell the commits it would force over are its own.
func planPush(story, local, remote, lastPushed string, isAncestor func() bool) (pushKind, error) {
	branch := storygit.Branch(story)
	switch {
	case remote == "":
		return pushPlain, nil
	case remote == local:
		return pushNone, nil
	case isAncestor():
		return pushPlain, nil
	case lastPushed == "":
		return pushNone, fmt.Errorf("origin's %[1]s is at %[2]s, which is not in the local branch at %[3]s, and flai has not pushed %[4]s before, so it does not know those commits are its own and will not force over them: "+
			"bring them into the story's branch (git fetch origin %[1]s, then rebase it onto them or merge them), or, when they are not wanted, delete origin's branch with git push origin --delete %[1]s; flai pushes it at its next look",
			branch, short(remote), short(local), story)
	}
	return pushLease, nil
}

// pushBranch pushes local to origin's branch of story, plainly or with a
// lease on lastPushed; it never forces without the lease, nor any other branch.
func pushBranch(git execx.Runner, root, story, local, remote, lastPushed string, kind pushKind) error {
	branch := storygit.Branch(story)
	ref := "refs/heads/" + branch
	args := []string{"push", "--quiet"}
	if kind == pushLease {
		args = append(args, "--force-with-lease="+ref+":"+lastPushed)
	}
	args = append(args, "origin", local+":"+ref)
	out, err := git.Run(root, "git", args...)
	if err == nil {
		return nil
	}
	if kind == pushLease {
		return fmt.Errorf("git refused to push %[1]s, which a sync rebased, over origin's at %[2]s with a lease on %[3]s, the head flai last pushed (%[4]s): someone pushed to origin's %[1]s since. "+
			"Bring their commits into the story's branch (git fetch origin %[1]s, then rebase it onto them or merge them), or, when they are not wanted, put origin's branch back with git push --force-with-lease=%[5]s:%[2]s origin %[3]s:%[5]s; flai pushes it at its next look",
			branch, short(remote), lastPushed, gitSaid(out, err), ref)
	}
	return fmt.Errorf("git refused to push %s to origin (%s): give the host's git credentials push access to origin's repository; flai tries again at its next look", branch, gitSaid(out, err))
}

// remoteHead is the head of branch on origin, or "" when origin has none.
func remoteHead(git execx.Runner, root, branch string) (string, error) {
	ref := "refs/heads/" + branch
	out, err := git.Run(root, "git", "ls-remote", "origin", ref)
	if err != nil {
		return "", fmt.Errorf("could not read %s on origin (%s): check that origin is reachable from this host and that its git credentials can read it; flai tries again at its next look", branch, gitSaid(out, err))
	}
	for line := range strings.SplitSeq(out, "\n") {
		if sha, name, ok := strings.Cut(strings.TrimSpace(line), "\t"); ok && name == ref {
			return sha, nil
		}
	}
	return "", nil
}

// forgeReason says why the forge refused the story's pull request and what
// clears it, naming the token's variable and never the token.
func forgeReason(f forge.Forge, env, story string, err error) error {
	name, repo := f.Name(), f.Repository()
	switch {
	case errors.Is(err, forge.ErrUnauthorized):
		return fmt.Errorf("%s did not accept the token in %s (%w): put a valid token in %s where flai serve starts, or name the variable that holds one with flai serve forge %s --token-env NAME", name, env, err, env, name)
	case errors.Is(err, forge.ErrForbidden):
		return fmt.Errorf("%s refused the token in %s this request (%w): give it pull requests write access to %s (pull requests: write for a fine-grained token, the repo scope for a classic one), or wait for the rate limit to reset", name, env, err, repo)
	case errors.Is(err, forge.ErrNotFound):
		return fmt.Errorf("%s has no repository %s that the token in %s can see (%w): check that origin names it and that the token is granted access to it", name, repo, env, err)
	}
	return fmt.Errorf("could not open or update %s's pull request on %s %s (%w); flai tries again at its next look", story, name, repo, err)
}

// gitSaid is git's own words for a failed command, or the error when it said nothing.
func gitSaid(out string, err error) string {
	if s := strings.TrimSpace(out); s != "" {
		return s
	}
	return err.Error()
}

// short is the first twelve characters of a commit's name.
func short(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}

// mainCheckout is the project's main checkout, which holds origin's
// configuration, the story branches, and the pull request template.
func mainCheckout(repo *workitem.Repo) string {
	if repo.MainRoot != "" {
		return repo.MainRoot
	}
	return repo.Root
}

// forgePass pushes the stories in review of one served project once a
// minute; its git and its forge come from the options, so that it is
// tested without a remote or a forge.
type forgePass struct {
	o  Options
	e  Entry
	wg *sync.WaitGroup
}

func newForgePass(o Options, e Entry, wg *sync.WaitGroup) *forgePass {
	return &forgePass{o: o, e: e, wg: wg}
}

// start looks now and then every minute in the background until ctx ends:
// a push or a forge call may take long, and the project's other looks must
// not wait for it.
func (p *forgePass) start(ctx context.Context) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			p.look(ctx)
			select {
			case <-ctx.Done():
				return
			case <-time.After(forgeEvery):
			}
		}
	}()
}

func (p *forgePass) enabled() bool {
	return p.o.Host.Enabled != nil && p.o.Host.Enabled(hostapi.ActionForge, p.e.Root)
}

func (p *forgePass) token(name string) ForgeToken {
	if p.o.ForgeToken == nil {
		return ForgeToken{}
	}
	return p.o.ForgeToken(name)
}

// look pushes each story in review, when the project accepts through pull
// requests and the forge host action is on, and settles each one's failure
// thread; it resolves the failure thread of a story that left review.
func (p *forgePass) look(ctx context.Context) {
	if !p.enabled() {
		return
	}
	repo, err := workitem.Open(p.e.Root)
	if err != nil || !repo.Manifest.Acceptance.PullRequest() {
		return
	}
	items, err := repo.List(false)
	if err != nil {
		p.o.Logger.Warn("stories in review not pushed", "component", "serve", "root", p.e.Root, "err", err.Error())
		return
	}
	token := p.token(repo.Manifest.Acceptance.ForgeOrDefault())
	inReview := map[string]bool{}
	var pushed, opened, updated, failed []string
	for _, it := range items {
		if it.Type != workitem.Story || it.Status != workitem.Review {
			continue
		}
		if ctx.Err() != nil {
			return
		}
		inReview[it.ID] = true
		was, _, _ := ReadForgeState(repo, it.ID)
		st, err := PushStory(ctx, ForgePush{Repo: repo, Story: it.ID, Git: gitOf(p.o), Token: token, NewForge: p.o.NewForge, HostAction: p.enabled, Now: p.o.Now})
		if st.Story == "" {
			p.o.Logger.Warn("story in review not pushed", "component", "serve", "root", p.e.Root, "story", it.ID, "err", err.Error())
			continue
		}
		p.settle(repo, st)
		p.journal(was, st)
		switch st.Outcome {
		case ForgeFailed:
			failed = append(failed, st.Story)
		case ForgePushed:
			pushed = append(pushed, st.Story)
		}
		switch st.Request {
		case RequestOpened:
			opened = append(opened, st.Story)
		case RequestUpdated:
			updated = append(updated, st.Story)
		}
	}
	p.leftReview(repo, inReview)
	if len(pushed)+len(opened)+len(updated) > 0 {
		p.o.Logger.Info("stories in review pushed", "component", "serve", "root", p.e.Root,
			"pushed", strings.Join(pushed, ","), "opened", strings.Join(opened, ","), "updated", strings.Join(updated, ","))
	}
	if len(failed) > 0 {
		p.o.Logger.Warn("stories in review not pushed", "component", "serve", "root", p.e.Root, "failed", strings.Join(failed, ","))
	}
}

// settle keeps the story's one failure thread: it opens it on a failure,
// adds an entry when the reason changes, and resolves it once a push
// succeeds. A thread someone resolved while the same failure stands is not
// opened again until the reason changes.
func (p *forgePass) settle(repo *workitem.Repo, st ForgeState) {
	if st.Outcome == ForgeSkipped {
		return
	}
	th := failureThread(repo, st.Story)
	now := p.o.Now()
	next := st
	switch {
	case st.Outcome != ForgeFailed:
		if th != nil {
			if _, err := threads.Resolve(repo, th.ID, forgeAuthor, pushedNote(st), now); err != nil {
				p.o.Logger.Warn("forge failure thread not resolved", "component", "serve", "root", p.e.Root, "story", st.Story, "thread", th.ID, "err", err.Error())
				return
			}
		}
		next.Thread, next.ThreadReason = "", ""
	case th != nil:
		next.Thread = th.ID
		if st.ThreadReason != st.Reason {
			if _, err := threads.Reply(repo, th.ID, forgeAuthor, failureText(st), now); err != nil {
				p.o.Logger.Warn("forge failure not written on its thread", "component", "serve", "root", p.e.Root, "story", st.Story, "thread", th.ID, "err", err.Error())
				return
			}
			next.ThreadReason = st.Reason
		}
	case st.Thread != "" && st.ThreadReason == st.Reason:
		// resolved by someone while the same failure stands
	default:
		opened, err := threads.New(repo, threads.NewOptions{Title: failureTitle(st.Story), On: st.Story, Author: forgeAuthor, Text: failureText(st), Now: now})
		if err != nil {
			p.o.Logger.Warn("forge failure thread not opened", "component", "serve", "root", p.e.Root, "story", st.Story, "err", err.Error())
			return
		}
		next.Thread, next.ThreadReason = opened.ID, st.Reason
	}
	if next != st {
		if err := writeForgeState(repo, next); err != nil {
			p.o.Logger.Warn("forge failure thread not recorded", "component", "serve", "root", p.e.Root, "story", st.Story, "err", err.Error())
		}
	}
}

// leftReview resolves the failure thread of each story flai recorded one
// for that is no longer in review: flai pushes it no more.
func (p *forgePass) leftReview(repo *workitem.Repo, inReview map[string]bool) {
	files, _ := filepath.Glob(filepath.Join(repo.CacheDir(), "forge", "*.json"))
	for _, f := range files {
		id := strings.TrimSuffix(filepath.Base(f), ".json")
		if inReview[id] {
			continue
		}
		st, ok, err := ReadForgeState(repo, id)
		if err != nil || !ok || st.Thread == "" {
			continue
		}
		if th := failureThread(repo, id); th != nil {
			why := fmt.Sprintf("%s is no longer in review, so flai serve pushes it no more.", id)
			if _, err := threads.Resolve(repo, th.ID, forgeAuthor, why, p.o.Now()); err != nil {
				p.o.Logger.Warn("forge failure thread not resolved", "component", "serve", "root", p.e.Root, "story", id, "thread", th.ID, "err", err.Error())
				continue
			}
		}
		st.Thread, st.ThreadReason = "", ""
		if err := writeForgeState(repo, st); err != nil {
			p.o.Logger.Warn("forge failure thread not recorded", "component", "serve", "root", p.e.Root, "story", id, "err", err.Error())
		}
	}
}

// journal records a push or a request that changed something, and a
// failure the first time its reason is seen.
func (p *forgePass) journal(was, st ForgeState) {
	if p.o.Host.Record == nil {
		return
	}
	entry := hostapi.Entry{At: st.At, Action: hostapi.ActionForge, Method: "serve.forge", Project: p.e.Key, Root: p.e.Root, By: forgeAuthor}
	switch {
	case st.Outcome == ForgeFailed:
		if was.Outcome == ForgeFailed && was.Reason == st.Reason {
			return
		}
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("could not push %s or open its pull request: %s", st.Story, st.Reason)
	case st.Outcome == ForgePushed || st.Request == RequestOpened || st.Request == RequestUpdated:
		entry.Outcome, entry.Detail = "done", pushedNote(st)
	default:
		return
	}
	p.o.Host.Record(entry)
}

// failureThread is the open thread flai serve opened on story for a failure, or nil.
func failureThread(repo *workitem.Repo, story string) *threads.Thread {
	all, err := threads.For(repo, story)
	if err != nil {
		return nil
	}
	title := workitem.CleanTitle(failureTitle(story))
	for _, th := range all {
		if th.Open() && th.Opener() == forgeAuthor && th.Title == title {
			return th
		}
	}
	return nil
}

func failureTitle(story string) string {
	return fmt.Sprintf("flai serve could not push %s's branch or open its pull request", story)
}

// failureText names the reason, which says what clears it, in a block that
// keeps its variable names and commands as they are.
func failureText(st ForgeState) string {
	return fmt.Sprintf("flai serve could not push %s to origin or open or update its pull request:\n\n```text\n%s\n```\n\n"+
		"It tries again every minute, and resolves this thread once the push and the pull request succeed. flai forge status %s says what it last did.",
		storygit.Branch(st.Story), st.Reason, st.Story)
}

// pushedNote says what a push that succeeded did.
func pushedNote(st ForgeState) string {
	branch := fmt.Sprintf("origin's %s is at %s", storygit.Branch(st.Story), short(st.RemoteHead))
	if st.Outcome == ForgePushed {
		branch = fmt.Sprintf("pushed %s at %s", storygit.Branch(st.Story), short(st.LastPushed))
		if st.Forced {
			branch += " with a lease on the head it last pushed, after a sync rebased it"
		}
	}
	request := "its pull request"
	if st.PullRequest > 0 {
		request = fmt.Sprintf("pull request #%d", st.PullRequest)
	}
	switch st.Request {
	case RequestOpened:
		request = "opened " + request
	case RequestUpdated:
		request = "updated " + request
	default:
		request += " is up to date"
	}
	if st.URL != "" {
		request += ", " + st.URL
	}
	return fmt.Sprintf("flai serve %s, and %s.", branch, request)
}
