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
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Merging a story's pull request through the remote (S-0397, ADR-0170).
// Under pull-request acceptance, flai accept merges the story's pull request
// on the forge by acceptance.merge_method, then deletes story/S-nnnn on
// origin. MergeStory does that for one story, with the token of the process
// that calls it. An agent has no token, so flai accept hands the merge to
// flai serve on the host: it writes a request under
// .flai-cache/forge/merge, which flai serve's merge pass takes by renaming
// it, merges behind the forge host action with its own token, and answers in
// a file beside it. No listener and no secret is added: the hand-off is the
// project's own cache folder, which only the host's user writes.

const (
	// mergeEvery is how often the merge pass looks for a request.
	mergeEvery = time.Second
	// The suffixes of a request, a request flai serve took, and its answer.
	mergeRequestSuffix = ".request.json"
	mergeTakenSuffix   = ".taken.json"
	mergeAnswerSuffix  = ".answer.json"
)

// MergeOutcome is what a merge of a story's pull request came to.
type MergeOutcome string

// The outcomes of a merge.
const (
	// MergeMerged is a pull request the forge merged at this request.
	MergeMerged MergeOutcome = "merged"
	// MergeAlready is a pull request the forge had merged before, by a
	// reviewer, say, or by an earlier acceptance that stopped after it.
	MergeAlready MergeOutcome = "already-merged"
	// MergeRefused is a pull request not merged, for the reason given.
	MergeRefused MergeOutcome = "refused"
)

// MergeAnswer is what a merge of a story's pull request did; it never holds
// the token.
type MergeAnswer struct {
	// Story is the story's ID.
	Story string `json:"story"`
	// PullRequest and URL are its pull request, once found.
	PullRequest int    `json:"pull_request,omitempty"`
	URL         string `json:"url,omitempty"`
	// Method is the merge method asked for.
	Method forge.MergeMethod `json:"method,omitempty"`
	// Outcome is whether it merged, was merged already, or was refused, and
	// Reason why when it was refused, saying what clears it.
	Outcome MergeOutcome `json:"outcome"`
	Reason  string       `json:"reason,omitempty"`
	// SHA is the commit the merge left at the head of the base branch;
	// empty for a pull request merged before.
	SHA string `json:"sha,omitempty"`
	// BranchDeleted is whether story/S-nnnn was deleted on origin after the
	// merge, and BranchReason why not when it was not; a branch not deleted
	// leaves the merge standing.
	BranchDeleted bool   `json:"branch_deleted,omitempty"`
	BranchReason  string `json:"branch_reason,omitempty"`
	// At is when the merge was asked of the forge, RFC 3339.
	At string `json:"at"`
}

// Merged is whether the pull request is merged on the forge, at this request or before.
func (a MergeAnswer) Merged() bool {
	return a.Outcome == MergeMerged || a.Outcome == MergeAlready
}

// ForgeMerge is what MergeStory merges one story's pull request with.
type ForgeMerge struct {
	// Repo is the project; Story the story whose pull request is merged.
	Repo  *workitem.Repo
	Story string
	// Head is the head of the story's branch the acceptance gate passed at:
	// a pull request whose head is another is not merged. Empty skips the check.
	Head string
	// Method is the merge method; empty is the manifest's
	// acceptance.merge_method, or its default.
	Method forge.MergeMethod
	// Git runs git; the system's when nil.
	Git execx.Runner
	// Token reads the token of the forge acceptance.forge names.
	Token ForgeToken
	// NewForge makes the forge's adapter; forge.New with GitHub's API when nil.
	NewForge NewForge
	// Now is the clock; time.Now when nil.
	Now func() time.Time
}

// MergeStory merges the story's pull request on the forge by the merge
// method, provided its head is still m.Head, and then deletes story/S-nnnn on
// origin. A pull request merged before is answered as such, and its branch
// deleted, without a second merge. Every failure is an answer refused with
// its reason, saying what clears it, and changes nothing.
func MergeStory(ctx context.Context, m ForgeMerge) MergeAnswer {
	now := time.Now
	if m.Now != nil {
		now = m.Now
	}
	a := MergeAnswer{Story: m.Story, Method: m.Method, At: now().UTC().Format(time.RFC3339)}
	if err := m.merge(ctx, &a, now()); err != nil {
		a.Outcome, a.Reason = MergeRefused, err.Error()
	}
	return a
}

// merge does MergeStory's work, filling a as it goes; its error is the
// reason the merge was refused.
func (m ForgeMerge) merge(ctx context.Context, a *MergeAnswer, now time.Time) error {
	it, err := m.Repo.Get(m.Story)
	if err != nil {
		return err
	}
	if it.Type != workitem.Story {
		return fmt.Errorf("%s is a %s: flai merges a story's pull request, so name its story", it.ID, it.Type)
	}
	acc := m.Repo.Manifest.Acceptance
	if !acc.PullRequest() {
		return fmt.Errorf("acceptance.mode is %s, so %s is merged on this machine by flai accept and nothing is merged on the forge: set acceptance.mode: pull-request in system-flow.yaml to merge stories through their pull requests",
			acc.ModeOrDefault(), it.ID)
	}
	if a.Method == "" {
		if a.Method, err = forge.ParseMergeMethod(acc.MergeMethodOrDefault()); err != nil {
			return fmt.Errorf("acceptance.merge_method in system-flow.yaml: %w", err)
		}
	} else if !a.Method.Valid() {
		_, err := forge.ParseMergeMethod(string(a.Method))
		return err
	}
	git, root := m.git(), mainCheckout(m.Repo)
	f, err := m.adapter(git, root, it.ID)
	if err != nil {
		return err
	}
	fctx, cancel := context.WithTimeout(ctx, forgeTimeout)
	defer cancel()
	pr, err := m.find(fctx, f, it.ID)
	if err != nil {
		return err
	}
	a.PullRequest, a.URL = pr.Number, pr.URL
	switch {
	case pr.State == forge.PullRequestMerged:
		a.Outcome = MergeAlready
	case pr.State == forge.PullRequestClosed:
		return fmt.Errorf("%s's pull request #%d on %s %s was closed without merging, so flai merges nothing: reopen it on the forge and accept again, or send %s back to in-progress",
			it.ID, pr.Number, f.Name(), f.Repository(), it.ID)
	case m.Head != "" && pr.HeadSHA != m.Head:
		return fmt.Errorf("%s's pull request #%d on %s %s is at %s, not at %s where its checks and reviews passed, so flai merges nothing: flai serve reads its checks and reviews again at its next look; accept once they pass at the new head",
			it.ID, pr.Number, f.Name(), f.Repository(), short(pr.HeadSHA), short(m.Head))
	default:
		merged, err := f.Merge(fctx, pr, a.Method)
		if err != nil {
			return mergeReason(f, m.Token.Env, it.ID, pr.Number, a.Method, err)
		}
		a.Outcome, a.SHA = MergeMerged, merged.SHA
		pr.State = forge.PullRequestMerged
	}
	if err := keepMerged(m.Repo, it.ID, pr, now); err != nil {
		return fmt.Errorf("%s's pull request #%d is merged on %s %s, but flai could not record it: %w; accept again, and flai finds the merge rather than making another", it.ID, pr.Number, f.Name(), f.Repository(), err)
	}
	a.BranchDeleted, a.BranchReason = deleteStoryBranch(git, root, it.ID)
	return nil
}

func (m ForgeMerge) git() execx.Runner {
	if m.Git == nil {
		return execx.System{}
	}
	return m.Git
}

// adapter is the project's forge, bound to origin's repository, once the
// token can be read.
func (m ForgeMerge) adapter(git execx.Runner, root, story string) (forge.Forge, error) {
	name := m.Repo.Manifest.Acceptance.ForgeOrDefault()
	if m.Token.Read == nil {
		return nil, fmt.Errorf("flai was given no way to read the %s token, so it merges nothing: accept with flai serve running on the host (flai serve start), which merges with its token", name)
	}
	if _, err := m.Token.Read(); err != nil {
		return nil, fmt.Errorf("%s's pull request cannot be merged without the %s token, so flai merges nothing: %w (flai serve, and flai host that starts it, read it from the environment they started with)", story, name, err)
	}
	url, err := git.Run(root, "git", "remote", "get-url", "origin")
	if err != nil {
		return nil, fmt.Errorf("the project has no origin remote whose pull request for %s flai could merge (%s): add the forge's repository with git remote add origin, then its URL", storygit.Branch(story), gitSaid(url, err))
	}
	if m.NewForge == nil {
		return forge.New(name, strings.TrimSpace(url), m.Token.Read, forge.GitHubAPI)
	}
	return m.NewForge(name, strings.TrimSpace(url), m.Token.Read)
}

// find is the story's pull request: the one numbered as flai's push record
// or its kept status names it, or else the one open from its branch. A
// numbered one that is closed gives way to one open from the branch, which
// flai opened after it.
func (m ForgeMerge) find(ctx context.Context, f forge.Forge, story string) (forge.PullRequest, error) {
	branch := storygit.Branch(story)
	pushed, _, err := ReadForgeState(m.Repo, story)
	if err != nil {
		return forge.PullRequest{}, err
	}
	kept, _, err := forge.ReadStatus(m.Repo.CacheDir(), story)
	if err != nil {
		return forge.PullRequest{}, err
	}
	var pr forge.PullRequest
	if n := max(pushed.PullRequest, kept.PullRequest); n > 0 {
		if pr, err = f.PullRequest(ctx, n); err != nil {
			return forge.PullRequest{}, mergeReason(f, m.Token.Env, story, n, "", err)
		}
		if pr.Head == branch && pr.State != forge.PullRequestClosed {
			return pr, nil
		}
	}
	open, found, err := f.FindPullRequest(ctx, branch)
	if err != nil {
		return forge.PullRequest{}, mergeReason(f, m.Token.Env, story, 0, "", err)
	}
	switch {
	case found:
		return open, nil
	case pr.Head == branch:
		return pr, nil // closed, with none opened after it
	}
	return forge.PullRequest{}, fmt.Errorf("%s has no pull request on %s %s from %s, so flai merges nothing: flai serve pushes the branch and opens it while the story is in review, behind the forge host action (flai serve enable forge); accept once it has",
		story, f.Name(), f.Repository(), branch)
}

// mergeReason says why the forge did not merge, or could not be asked to
// merge, story's pull request number (0 before it is found), and what clears
// it, naming the token's variable and never the token. A refusal carries the
// forge's own reason: a conflict, a branch out of date, or a rule.
func mergeReason(f forge.Forge, env, story string, number int, method forge.MergeMethod, err error) error {
	name, repo := f.Name(), f.Repository()
	request := story + "'s pull request"
	if number > 0 {
		request = fmt.Sprintf("%s's pull request #%d", story, number)
	}
	switch {
	case errors.Is(err, forge.ErrMergeRefused):
		return fmt.Errorf("%s refused to merge %s by %s, and nothing changed (%w): clear what it names on the forge (bring main into the branch for a conflict or a branch out of date, or meet the rule), then accept again", name, request, method, err)
	case errors.Is(err, forge.ErrUnauthorized):
		return fmt.Errorf("%s did not accept the token in %s (%w): put a valid token in %s where flai serve starts, or name the variable that holds one with flai serve forge %s --token-env NAME", name, env, err, env, name)
	case errors.Is(err, forge.ErrForbidden):
		return fmt.Errorf("%s refused the token in %s this request (%w): give it contents and pull requests write access to %s (contents: write and pull requests: write for a fine-grained token, the repo scope for a classic one), or wait for the rate limit to reset", name, env, err, repo)
	case errors.Is(err, forge.ErrNotFound):
		return fmt.Errorf("%s has no %s in %s that the token in %s can see (%w): check that origin names the repository and that the token is granted access to it", name, request, repo, env, err)
	}
	return fmt.Errorf("could not merge %s on %s %s (%w): check that the forge is reachable from this host, then accept again; a merge that went through is found, not repeated", request, name, repo, err)
}

// deleteStoryBranch deletes story/S-nnnn on origin with the host's git
// credentials, answering whether it did and, when it did not, why: origin
// has no such branch, the forge having deleted it on the merge, or git
// refused. A branch not deleted leaves the merge standing.
func deleteStoryBranch(git execx.Runner, root, story string) (bool, string) {
	branch := storygit.Branch(story)
	remote, err := remoteHead(git, root, branch)
	if err != nil {
		return false, fmt.Sprintf("the merge stands, but %s; delete the branch with git push origin --delete %s", err, branch)
	}
	if remote == "" {
		return false, fmt.Sprintf("origin has no %s to delete: the forge deleted it on the merge, or it was deleted before", branch)
	}
	if out, err := git.Run(root, "git", "push", "--quiet", "origin", "--delete", "refs/heads/"+branch); err != nil {
		return false, fmt.Sprintf("the merge stands, but git could not delete %s on origin (%s): give the host's git credentials push access to origin, or delete it with git push origin --delete %s", branch, gitSaid(out, err), branch)
	}
	return true, ""
}

// MergeRequest is flai accept's request that flai serve merge a story's pull
// request, written where MergeRequestPath says.
type MergeRequest struct {
	// Story is the story whose pull request is merged.
	Story string `json:"story"`
	// Head is the head the acceptance gate passed at; empty skips the check.
	Head string `json:"head,omitempty"`
	// Method is the merge method; empty is the manifest's.
	Method forge.MergeMethod `json:"method,omitempty"`
	// By is who accepts the story, and At when they asked, RFC 3339.
	By string `json:"by,omitempty"`
	At string `json:"at,omitempty"`
}

// mergeDir is the folder of merge requests and their answers under the
// project's cache at cacheDir.
func mergeDir(cacheDir string) string {
	return filepath.Join(cacheDir, "forge", "merge")
}

// MergeRequestPath is where flai accept writes its request that flai serve
// merge story's pull request.
func MergeRequestPath(repo *workitem.Repo, story string) string {
	return filepath.Join(mergeDir(repo.CacheDir()), story+mergeRequestSuffix)
}

// MergeTakenPath is where the request is once flai serve has taken it.
func MergeTakenPath(repo *workitem.Repo, story string) string {
	return filepath.Join(mergeDir(repo.CacheDir()), story+mergeTakenSuffix)
}

// MergeAnswerPath is where flai serve answers the request.
func MergeAnswerPath(repo *workitem.Repo, story string) string {
	return filepath.Join(mergeDir(repo.CacheDir()), story+mergeAnswerSuffix)
}

// MergeWait is how long RequestMerge waits for flai serve; a zero field is
// its default.
type MergeWait struct {
	// Take is how long flai serve has to take the request: 30s.
	Take time.Duration
	// Answer is how long it has to answer once it took it: 5m.
	Answer time.Duration
	// Poll is how often the files are looked at: 250ms.
	Poll time.Duration
}

func (w MergeWait) orDefaults() MergeWait {
	if w.Take <= 0 {
		w.Take = 30 * time.Second
	}
	if w.Answer <= 0 {
		w.Answer = 5 * time.Minute
	}
	if w.Poll <= 0 {
		w.Poll = 250 * time.Millisecond
	}
	return w
}

// RequestMerge asks flai serve on the host to merge req.Story's pull request
// and waits for its answer. A request flai serve does not take within w.Take
// is withdrawn, and one it took but did not answer within w.Answer is an
// error; either says what to do. A request withdrawn because ctx ended is
// never merged after.
func RequestMerge(ctx context.Context, repo *workitem.Repo, req MergeRequest, w MergeWait) (MergeAnswer, error) {
	w = w.orDefaults()
	request, taken, answer := MergeRequestPath(repo, req.Story), MergeTakenPath(repo, req.Story), MergeAnswerPath(repo, req.Story)
	for _, stale := range []string{answer, taken} {
		if err := os.Remove(stale); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return MergeAnswer{}, fmt.Errorf("remove the stale %s before asking flai serve to merge %s's pull request: %w; check that the project's .flai-cache is writable", stale, req.Story, err)
		}
	}
	data, err := json.MarshalIndent(req, "", "  ")
	if err != nil {
		return MergeAnswer{}, fmt.Errorf("write the request that flai serve merge %s's pull request: %w", req.Story, err)
	}
	if err := os.MkdirAll(filepath.Dir(request), 0o755); err != nil {
		return MergeAnswer{}, fmt.Errorf("write the request that flai serve merge %s's pull request: %w; check that the project's .flai-cache is writable", req.Story, err)
	}
	if err := atomicfile.WriteFile(request, append(data, '\n'), 0o644); err != nil {
		return MergeAnswer{}, fmt.Errorf("write the request that flai serve merge %s's pull request at %s: %w; check that the project's .flai-cache is writable", req.Story, request, err)
	}
	asked := time.Now()
	var tookAt time.Time
	tick := time.NewTicker(w.Poll)
	defer tick.Stop()
	for {
		if a, ok, err := readMergeAnswer(answer, req.Story); err != nil || ok {
			_ = os.Remove(answer)
			_ = os.Remove(taken)
			return a, err
		}
		switch {
		case tookAt.IsZero() && !exists(request):
			tookAt = time.Now()
		case tookAt.IsZero() && time.Since(asked) >= w.Take:
			if err := os.Remove(request); err == nil {
				return MergeAnswer{}, fmt.Errorf("flai serve did not take the request to merge %s's pull request within %s, so nothing was merged: start flai serve on this host (flai serve start) with the forge host action on (flai serve enable forge), then accept again", req.Story, w.Take)
			}
			tookAt = time.Now() // taken as the wait ran out
		case !tookAt.IsZero() && time.Since(tookAt) >= w.Answer:
			return MergeAnswer{}, fmt.Errorf("flai serve took the request to merge %s's pull request but did not answer within %s; its log says why. A merge it made is found, not repeated, when you accept again", req.Story, w.Answer)
		}
		select {
		case <-ctx.Done():
			if tookAt.IsZero() && os.Remove(request) == nil {
				return MergeAnswer{}, fmt.Errorf("stopped waiting for flai serve to take the request to merge %s's pull request, and withdrew it: %w", req.Story, ctx.Err())
			}
			return MergeAnswer{}, fmt.Errorf("stopped waiting for flai serve's answer to the request to merge %s's pull request, which it took: %w; a merge it made is found, not repeated, when you accept again", req.Story, ctx.Err())
		case <-tick.C:
		}
	}
}

// readMergeAnswer is the answer at path, and false with no error while there is none.
func readMergeAnswer(path, story string) (MergeAnswer, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return MergeAnswer{}, false, nil
	}
	if err != nil {
		return MergeAnswer{}, false, fmt.Errorf("read flai serve's answer to the request to merge %s's pull request at %s: %w; check that the project's .flai-cache is readable", story, path, err)
	}
	var a MergeAnswer
	if err := json.Unmarshal(data, &a); err != nil {
		return MergeAnswer{}, false, fmt.Errorf("read flai serve's answer to the request to merge %s's pull request at %s: %w; accept again, and a merge it made is found, not repeated", story, path, err)
	}
	return a, true, nil
}

// mergePass takes the merge requests of one served project every second and
// answers each: with the forge host action off it refuses, and otherwise it
// merges with flai serve's token. Its git and its forge come from the
// options, so that it is tested without a remote or a forge.
type mergePass struct {
	o     Options
	e     Entry
	wg    *sync.WaitGroup
	every time.Duration
}

func newMergePass(o Options, e Entry, wg *sync.WaitGroup) *mergePass {
	return &mergePass{o: o, e: e, wg: wg, every: mergeEvery}
}

// start looks now and then every second in the background until ctx ends.
func (p *mergePass) start(ctx context.Context) {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			p.look(ctx)
			select {
			case <-ctx.Done():
				return
			case <-time.After(p.every):
			}
		}
	}()
}

// look takes each request waiting and answers it; a request another look
// took first is left to it.
func (p *mergePass) look(ctx context.Context) {
	// the project is opened only when a request waits: the served root is
	// the main checkout, whose .flai-cache is the project's cache
	requests, _ := filepath.Glob(filepath.Join(mergeDir(filepath.Join(p.e.Root, ".flai-cache")), "*"+mergeRequestSuffix))
	if len(requests) == 0 {
		return
	}
	repo, err := workitem.Open(p.e.Root)
	if err != nil {
		p.o.Logger.Warn("merge requests not taken", "component", "serve", "root", p.e.Root, "err", err.Error())
		return
	}
	for _, request := range requests {
		if ctx.Err() != nil {
			return
		}
		story := strings.TrimSuffix(filepath.Base(request), mergeRequestSuffix)
		taken := MergeTakenPath(repo, story)
		if err := os.Rename(request, taken); err != nil {
			continue
		}
		a, off := p.answer(ctx, repo, story, taken)
		if err := writeMergeAnswer(repo, a); err != nil {
			p.o.Logger.Warn("merge answer not written", "component", "serve", "root", p.e.Root, "story", story, "err", err.Error())
		}
		p.journal(a, off)
		if a.Merged() {
			p.o.Logger.Info("pull request merged", "component", "serve", "root", p.e.Root, "story", story, "outcome", string(a.Outcome),
				"pull_request", a.PullRequest, "method", string(a.Method), "branch_deleted", a.BranchDeleted)
		} else {
			p.o.Logger.Warn("pull request not merged", "component", "serve", "root", p.e.Root, "story", story, "outcome", string(a.Outcome), "reason", a.Reason)
		}
	}
}

// answer merges the story's pull request as the request at taken asks, or
// refuses, saying whether that was for the forge host action being off.
func (p *mergePass) answer(ctx context.Context, repo *workitem.Repo, story, taken string) (MergeAnswer, bool) {
	refused := func(why string) MergeAnswer {
		return MergeAnswer{Story: story, Outcome: MergeRefused, Reason: why, At: p.o.Now().UTC().Format(time.RFC3339)}
	}
	var req MergeRequest
	data, err := os.ReadFile(taken)
	if err == nil {
		err = json.Unmarshal(data, &req)
	}
	if err != nil {
		return refused(fmt.Sprintf("flai serve could not read the request to merge %s's pull request at %s (%v), so it merged nothing: accept again", story, taken, err)), false
	}
	if req.Story != story {
		return refused(fmt.Sprintf("the request at %s names %s, not %s, so flai serve merged nothing: accept again", taken, req.Story, story)), false
	}
	if p.o.Host.Enabled == nil || !p.o.Host.Enabled(hostapi.ActionForge, p.e.Root) {
		return refused("the forge host action is off for this project, so flai serve merges nothing: turn it on with flai serve enable forge"), true
	}
	var token ForgeToken
	if p.o.ForgeToken != nil {
		token = p.o.ForgeToken(repo.Manifest.Acceptance.ForgeOrDefault())
	}
	return MergeStory(ctx, ForgeMerge{Repo: repo, Story: story, Head: req.Head, Method: req.Method, Git: gitOf(p.o), Token: token, NewForge: p.o.NewForge, Now: p.o.Now}), false
}

// journal records a merge, or a refusal, in the host journal; off is a
// refusal for the forge host action being off.
func (p *mergePass) journal(a MergeAnswer, off bool) {
	if p.o.Host.Record == nil {
		return
	}
	entry := hostapi.Entry{At: a.At, Action: hostapi.ActionForge, Method: "serve.merge", Project: p.e.Key, Root: p.e.Root, By: forgeAuthor, Outcome: "done"}
	switch {
	case a.Outcome == MergeMerged:
		entry.Detail = fmt.Sprintf("merged %s's pull request #%d by %s at %s", a.Story, a.PullRequest, a.Method, short(a.SHA))
	case a.Outcome == MergeAlready:
		entry.Detail = fmt.Sprintf("found %s's pull request #%d merged already", a.Story, a.PullRequest)
	case off:
		entry.Outcome, entry.Detail = "disabled", a.Reason
	default:
		entry.Outcome, entry.Detail = "failed", fmt.Sprintf("did not merge %s's pull request: %s", a.Story, a.Reason)
	}
	p.o.Host.Record(entry)
}

// writeMergeAnswer writes a where flai accept waits for it.
func writeMergeAnswer(repo *workitem.Repo, a MergeAnswer) error {
	path := MergeAnswerPath(repo, a.Story)
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return err
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("answer the request to merge %s's pull request at %s: %w; check that the project's .flai-cache is writable", a.Story, path, err)
	}
	return nil
}
