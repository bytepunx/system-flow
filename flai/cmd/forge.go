package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/serve"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// forgeCheckTimeout bounds flai forge check's call to the forge.
const forgeCheckTimeout = 30 * time.Second

// githubAPIBase is the GitHub API flai's forge commands call; tests point it at a stand-in.
var githubAPIBase = forge.GitHubAPI

// newForgeCmd groups the commands about the project's forge (ADR-0165).
func newForgeCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "forge",
		Short: "The project's forge, through which pull-request acceptance merges a story",
		Long: `A project whose system-flow.yaml sets acceptance.mode: pull-request accepts a
story through a pull request on its forge, acceptance.forge (github, the
default), rather than by merging on this machine (ADR-0165). flai reaches the
forge's repository named by the origin remote, with a token it reads from the
environment variable this host names for the forge: GITHUB_TOKEN unless
flai serve forge <forge> --token-env <NAME> names another.

With the forge host action on (flai serve enable forge), flai serve pushes
each story in review to origin and opens or updates its one pull request
(ADR-0166). flai forge check says whether the token can read and push to the
repository; flai forge push does one story's push and pull request on your
word; flai forge status says what flai last did for a story, and why nothing
was pushed when nothing was.`,
		Example: `  flai forge check
  flai forge push S-0042
  flai forge status S-0042
  flai forge status S-0042 --json`,
	}
	c.AddCommand(newForgeCheckCmd(a), newForgePushCmd(a), newForgeStatusCmd(a))
	return c
}

func newForgePushCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "push S-nnnn",
		Short: "Push a story in review to origin and open or update its pull request, on the operator's word",
		Long: `Do for one story what flai serve does at each look under the forge host
action (ADR-0166), whether that action is on or not: push story/S-nnnn to
origin with this host's git credentials when origin's branch is missing or not
the local head, plainly, or with --force-with-lease on the head flai last
pushed when a sync rebased it, then open the story's one pull request against
the main branch, found by its head branch, or update it when its title, body,
or base is not what flai would write. The token is read from the variable flai
forge check names, in this shell's environment. What was done is recorded
where flai forge status reads it.

It needs acceptance.mode: pull-request in system-flow.yaml and the story in
review. flai guard refuses it to every agent: it is the operator's.

It prints the outcome (pushed, or up-to-date when origin had the local head
already), the branch it pushed and whether with the lease, and the pull
request's number, what was done to it (opened, updated, or unchanged), and its
URL. The exit status is 1, with the reason and what clears it, when nothing
was done (the mode is local, or the story is not in review) or the push or the
request failed; what was done before a failure is printed still. With --json
the result is the state flai records: {"story", "outcome", "reason", "at",
"local_head", "remote_head", "last_pushed", "pushed_at", "forced",
"pull_request", "url", "request", "thread", "thread_reason"}, empty fields
left out; "forced" is about the push at "pushed_at".`,
		Example: `  flai forge push S-0042
  flai forge push S-0042 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			st, err := serve.PushStory(cmd.Context(), serve.ForgePush{
				Repo: repo, Story: args[0], Git: a.runner, Token: a.serveForgeToken(repo.Manifest.Acceptance.ForgeOrDefault()),
				NewForge: newForge, Now: a.now,
			})
			if st.Story == "" {
				return err
			}
			if perr := a.printForgePush(st); perr != nil {
				return errors.Join(err, perr)
			}
			if err != nil {
				return err
			}
			if st.Outcome == serve.ForgeSkipped {
				return errors.New(st.Reason)
			}
			return nil
		},
	}
}

// printForgePush prints what flai forge push did, one finding a line, or the
// state as JSON; a reason goes out through the error path.
func (a *app) printForgePush(st serve.ForgeState) error {
	if a.jsonOut {
		return a.printJSON(st)
	}
	fmt.Fprintf(a.out, "story: %s\noutcome: %s\n", st.Story, st.Outcome)
	branch := storyBranch(st.Story)
	switch {
	case st.Outcome == serve.ForgeUpToDate:
		fmt.Fprintf(a.out, "branch: origin had %s at %s already\n", branch, short(st.RemoteHead))
	case st.Outcome == serve.ForgePushed || (st.Outcome == serve.ForgeFailed && st.PushedAt != "" && st.PushedAt == st.At):
		lease := ""
		if st.Forced {
			lease = ", with --force-with-lease on the head flai last pushed, after a sync rebased it"
		}
		fmt.Fprintf(a.out, "branch: pushed %s at %s to origin%s\n", branch, short(st.LastPushed), lease)
	}
	if st.Request != "" {
		fmt.Fprintf(a.out, "pull request: #%d %s, %s\n", st.PullRequest, st.Request, st.URL)
	}
	return nil
}

func newForgeStatusCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "status S-nnnn",
		Short: "Say what flai last did for a story's pull request, whether origin has its branch, and why nothing was pushed",
		Long: `Print, for one story under pull-request acceptance (ADR-0166), its pull
request's number and URL as flai last opened or found it ("none yet" before
that), and whether origin's story/S-nnnn matches the local branch, read now
with git ls-remote origin from this host; when origin cannot be read, it says
so and why. It prints what flai's last push or look did and when, and, once flai serve
has read the pull request back (ADR-0167), the head commit and time of that
read, each required check's name, state, log link, and summary, and each
review's author, state, and link.

When nothing was pushed, or the last push failed or no longer matches, it
says why, the first of these that applies: acceptance.mode is local; the
forge host action is off for this project on this host (flai forge push still
pushes on your word); the token's variable is unset, as this shell sees it
(flai serve reads the environment it started with); the last push failed, with
its reason, when, and its thread; the story is not in review; or nothing stops
flai serve, which pushes it at its next look. The token is never printed.

It does not call the forge, so it works offline and costs no rate limit. The
exit status is 0 whenever it can report. With --json the result is {"story",
"status", "mode", "forge", "host_action", "token_env", "token_set",
"pull_request", "url", "branch", "local_head", "remote_head", "remote_error",
"matches", "outcome", "reason", "at", "thread", "why", "pull_request_read"},
empty fields left out but for the booleans; "pull_request_read" is the pull
request as flai serve last read it, as flai show --json gives it under
"pull_request": {"story", "pull_request", "url", "head_sha", "read_at",
"checks", "reviews", "sent_back"}.`,
		Example: `  flai forge status S-0042
  flai forge status S-0042 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			cfg, _, err := a.loadConfig()
			if err != nil {
				return err
			}
			s, err := a.forgeStatus(repo, args[0], cfg)
			if err != nil {
				return err
			}
			return a.printForgeStatus(s)
		},
	}
}

// forgeStatus is what flai forge status found; it never holds the token.
type forgeStatus struct {
	Story       string             `json:"story"`
	Status      string             `json:"status"`
	Mode        string             `json:"mode"`
	Forge       string             `json:"forge"`
	HostAction  bool               `json:"host_action"`
	TokenEnv    string             `json:"token_env"`
	TokenSet    bool               `json:"token_set"`
	PullRequest int                `json:"pull_request,omitempty"`
	URL         string             `json:"url,omitempty"`
	Branch      string             `json:"branch"`
	LocalHead   string             `json:"local_head,omitempty"`
	RemoteHead  string             `json:"remote_head,omitempty"`
	RemoteError string             `json:"remote_error,omitempty"`
	Matches     bool               `json:"matches"`
	Outcome     serve.ForgeOutcome `json:"outcome,omitempty"`
	Reason      string             `json:"reason,omitempty"`
	At          string             `json:"at,omitempty"`
	Thread      string             `json:"thread,omitempty"`
	Why         string             `json:"why,omitempty"`
	// PullRequestRead is the pull request as flai serve last read it, nil
	// before it has (ADR-0167); pull_request is taken by the number.
	PullRequestRead *forge.Status `json:"pull_request_read,omitempty"`
}

// forgeStatus reads what flai last did for the story id from its state
// file, and origin's and the local branch from git, without calling the forge.
func (a *app) forgeStatus(repo *workitem.Repo, id string, cfg config.Config) (forgeStatus, error) {
	it, err := repo.Get(id)
	if err != nil {
		return forgeStatus{}, err
	}
	if it.Type != workitem.Story {
		return forgeStatus{}, fmt.Errorf("%s is a %s: flai pushes a story's branch, so name its story", it.ID, it.Type)
	}
	st, ok, err := serve.ReadForgeState(repo, it.ID)
	if err != nil {
		return forgeStatus{}, err
	}
	acc := repo.Manifest.Acceptance
	root := mainRootOf(repo)
	s := forgeStatus{
		Story: it.ID, Status: it.Status, Mode: acc.ModeOrDefault(), Forge: acc.ForgeOrDefault(),
		HostAction: cfg.ActionEnabled(hostapi.ActionForge, root), Branch: storyBranch(it.ID),
		PullRequest: st.PullRequest, URL: st.URL, Outcome: st.Outcome, Reason: st.Reason, At: st.At, Thread: st.Thread,
	}
	if s.PullRequestRead, err = keptPullRequest(repo, it.ID); err != nil {
		return forgeStatus{}, err
	}
	s.TokenEnv = cfg.ForgeTokenEnv(s.Forge)
	_, tokenErr := cfg.ForgeTokenFunc(s.Forge, os.Getenv)()
	s.TokenSet = tokenErr == nil
	if local, err := a.runner.Run(root, "git", "rev-parse", "--verify", "--quiet", "refs/heads/"+s.Branch+"^{commit}"); err == nil {
		s.LocalHead = strings.TrimSpace(local)
	}
	s.RemoteHead, s.RemoteError = a.originHead(root, s.Branch)
	s.Matches = s.LocalHead != "" && s.RemoteHead == s.LocalHead
	pushed := ok && (st.Outcome == serve.ForgePushed || st.Outcome == serve.ForgeUpToDate)
	if !pushed || (!s.Matches && s.RemoteError == "") {
		s.Why = whyNotPushed(s, acc.PullRequest(), tokenErr)
	}
	return s, nil
}

// keptPullRequest is story's pull request as flai serve last read it, without
// its events, or nil when it has not read one; it never calls the forge.
func keptPullRequest(repo *workitem.Repo, story string) (*forge.Status, error) {
	st, ok, err := forge.ReadStatus(repo.CacheDir(), story)
	if err != nil || !ok {
		return nil, err
	}
	st.Events = nil
	return &st, nil
}

// pullRequestLines are a kept pull request's required checks, then its
// reviews, one a line: its name, its state, and its link when it has one.
func pullRequestLines(pr *forge.Status) []string {
	lines := make([]string, 0, len(pr.Checks)+len(pr.Reviews))
	for _, c := range pr.Checks {
		l := fmt.Sprintf("check %s: %s", c.Name, c.State)
		if c.URL != "" {
			l += " " + c.URL
		}
		if c.Summary != "" {
			l += " (" + c.Summary + ")"
		}
		lines = append(lines, l)
	}
	for _, r := range pr.Reviews {
		l := fmt.Sprintf("review %s: %s", r.Author, r.State)
		if r.URL != "" {
			l += " " + r.URL
		}
		lines = append(lines, l)
	}
	return lines
}

// originHead is the head of branch on origin, "" when origin has none, or
// why origin could not be read.
func (a *app) originHead(root, branch string) (head, why string) {
	ref := "refs/heads/" + branch
	out, err := a.runner.Run(root, "git", "ls-remote", "origin", ref)
	if err != nil {
		return "", fmt.Sprintf("could not read %s on origin (%s): check that origin is reachable from this host and that its git credentials can read it", branch, gitSaid(out, err))
	}
	for line := range strings.SplitSeq(out, "\n") {
		if sha, name, ok := strings.Cut(strings.TrimSpace(line), "\t"); ok && name == ref {
			return sha, ""
		}
	}
	return "", ""
}

// whyNotPushed is the first reason that applies for s's branch not being on
// origin as the local head, each saying what clears it.
func whyNotPushed(s forgeStatus, pullRequest bool, tokenErr error) string {
	switch {
	case !pullRequest:
		return fmt.Sprintf("acceptance.mode is %s, so %s is accepted on this machine and nothing is pushed: set acceptance.mode: pull-request in system-flow.yaml to have stories in review pushed and their pull requests opened", s.Mode, s.Story)
	case !s.HostAction:
		return fmt.Sprintf("the forge host action is off for this project on this host, so flai serve pushes nothing: turn it on with flai serve enable forge; flai forge push %s pushes the story on your word without it", s.Story)
	case tokenErr != nil:
		return fmt.Sprintf("%v (as this shell sees it; flai serve reads the environment it started with)", tokenErr)
	case s.Outcome == serve.ForgeFailed:
		why := fmt.Sprintf("the last push, at %s, failed: %s", s.At, s.Reason)
		if s.Thread != "" {
			why += fmt.Sprintf(" (thread %s)", s.Thread)
		}
		return why
	case s.Status != workitem.Review:
		return fmt.Sprintf("%s is %s, not in review: flai pushes a story's branch and opens its pull request once the story is in review", s.Story, s.Status)
	}
	return fmt.Sprintf("nothing stops flai serve from pushing %s: it pushes it at its next look, once a minute while it runs, or push it now with flai forge push %s", s.Story, s.Story)
}

// printForgeStatus prints s, one finding a line, or as JSON.
func (a *app) printForgeStatus(s forgeStatus) error {
	if a.jsonOut {
		return a.printJSON(s)
	}
	action := "off"
	if s.HostAction {
		action = "on"
	}
	fmt.Fprintf(a.out, "story: %s (%s)\nmode: %s (forge %s)\nhost action: forge is %s\n", s.Story, s.Status, s.Mode, s.Forge, action)
	set := "not set"
	if s.TokenSet {
		set = "set"
	}
	fmt.Fprintf(a.out, "token: read from %s, %s\n", s.TokenEnv, set)
	if s.PullRequest > 0 {
		fmt.Fprintf(a.out, "pull request: #%d, %s\n", s.PullRequest, s.URL)
	} else {
		fmt.Fprintln(a.out, "pull request: none yet")
	}
	if pr := s.PullRequestRead; pr != nil {
		fmt.Fprintf(a.out, "  read at %s, %s\n", short(pr.HeadSHA), pr.ReadAt)
		for _, l := range pullRequestLines(pr) {
			fmt.Fprintf(a.out, "  %s\n", l)
		}
	}
	fmt.Fprintf(a.out, "branch: %s\n", branchFinding(s))
	if s.Outcome != "" {
		fmt.Fprintf(a.out, "last look: %s, at %s\n", s.Outcome, s.At)
	} else {
		fmt.Fprintln(a.out, "last look: none yet")
	}
	if s.Why != "" {
		fmt.Fprintf(a.out, "why: %s\n", s.Why)
	}
	return nil
}

// branchFinding says whether origin's branch matches the local one.
func branchFinding(s forgeStatus) string {
	switch {
	case s.RemoteError != "":
		return s.RemoteError
	case s.LocalHead == "" && s.RemoteHead == "":
		return fmt.Sprintf("neither this repository nor origin has %s", s.Branch)
	case s.LocalHead == "":
		return fmt.Sprintf("this repository has no %s; origin's is at %s", s.Branch, short(s.RemoteHead))
	case s.RemoteHead == "":
		return fmt.Sprintf("origin has no %s; the local head is at %s", s.Branch, short(s.LocalHead))
	case s.Matches:
		return fmt.Sprintf("origin's %s matches the local head at %s", s.Branch, short(s.LocalHead))
	}
	return fmt.Sprintf("origin's %s is at %s, the local head at %s", s.Branch, short(s.RemoteHead), short(s.LocalHead))
}

func newForgeCheckCmd(a *app) *cobra.Command {
	return &cobra.Command{
		Use:   "check",
		Short: "Say whether the forge token this host names can read and push to the repository behind origin",
		Long: `Print the forge (acceptance.forge, github by default, under either
acceptance mode, so the token can be checked before switching), the
repository the origin remote names, the environment variable the token is
read from, and whether the token can read the repository and push to it. The
variable is named; the token never is.

The exit status is 1, with the reason, when the project has no origin remote,
origin is not a repository on the forge, the variable is unset or empty, the
forge does not accept the token, the repository is not found or not visible to
the token, the forge refuses the token or has rate-limited it, or the token can
read but not push. What was learned before the failure is printed still. With
--json the result is {"forge", "repository", "token_env", "read", "push",
"ok", "reason"}.`,
		Example: `  flai forge check
  flai forge check --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.project()
			if err != nil {
				return err
			}
			cfg, _, err := a.loadConfig()
			if err != nil {
				return err
			}
			r := a.checkForge(repo.Root, repo.Manifest.Acceptance.ForgeOrDefault(), cfg)
			if err := a.printForgeCheck(r); err != nil {
				return err
			}
			if !r.OK {
				return errors.New(r.Reason)
			}
			return nil
		},
	}
}

// forgeCheck is what flai forge check found; it never holds the token.
type forgeCheck struct {
	Forge      string `json:"forge"`
	Repository string `json:"repository"`
	TokenEnv   string `json:"token_env"`
	Read       bool   `json:"read"`
	Push       bool   `json:"push"`
	OK         bool   `json:"ok"`
	Reason     string `json:"reason,omitempty"`
	asked      bool   // the forge answered, so Read and Push are its word
}

// checkForge asks name whether the token cfg names for it can read and push
// to the repository origin names in root, stopping at the first failure with
// its reason.
func (a *app) checkForge(root, name string, cfg config.Config) forgeCheck {
	r := forgeCheck{Forge: name, TokenEnv: cfg.ForgeTokenEnv(name)}
	remote, err := a.runner.Run(root, "git", "remote", "get-url", "origin")
	if err != nil {
		r.Reason = fmt.Sprintf("the project has no origin remote to read its repository from (%s): add the forge's repository with git remote add origin <url>", gitSaid(remote, err))
		return r
	}
	token := forge.TokenFunc(cfg.ForgeTokenFunc(name, os.Getenv))
	f, err := newForge(name, strings.TrimSpace(remote), token)
	if err != nil {
		r.Reason = err.Error()
		return r
	}
	r.Repository = f.Repository()
	if _, err := token(); err != nil {
		r.Reason = err.Error()
		return r
	}
	ctx, cancel := context.WithTimeout(context.Background(), forgeCheckTimeout)
	defer cancel()
	access, err := f.Access(ctx)
	if err != nil {
		r.asked = errors.Is(err, forge.ErrUnauthorized) || errors.Is(err, forge.ErrNotFound) || errors.Is(err, forge.ErrForbidden)
		r.Reason = forgeAccessReason(r, err)
		return r
	}
	r.asked, r.Read, r.Push = true, access.Read, access.Push
	if !r.Push {
		r.Reason = fmt.Sprintf("the token in %s can read %s but not push to it: give it push access to the repository (contents: write for a fine-grained token, the repo scope for a classic one)", r.TokenEnv, r.Repository)
		return r
	}
	r.OK = true
	return r
}

// newForge is the adapter for the forge name, bound to the repository remote
// names, calling the GitHub API at githubAPIBase.
func newForge(name, remote string, token forge.TokenFunc) (forge.Forge, error) {
	return forge.New(name, remote, token, githubAPIBase)
}

// forgeAccessReason says why the forge would not answer whether the token
// can read the repository, naming the variable and never the token.
func forgeAccessReason(r forgeCheck, err error) string {
	switch {
	case errors.Is(err, forge.ErrUnauthorized):
		return fmt.Sprintf("%s did not accept the token in %s (%v): put a valid token in %s, or name the variable that holds one with flai serve forge %s --token-env <NAME>", r.Forge, r.TokenEnv, err, r.TokenEnv, r.Forge)
	case errors.Is(err, forge.ErrNotFound):
		return fmt.Sprintf("%s has no repository %s that the token in %s can see (%v): check that origin names it and that the token is granted access to it", r.Forge, r.Repository, r.TokenEnv, err)
	case errors.Is(err, forge.ErrForbidden):
		return fmt.Sprintf("%s refused the token in %s (%v): check that the token may read %s, or wait for the rate limit to reset", r.Forge, r.TokenEnv, err, r.Repository)
	}
	return fmt.Sprintf("could not ask %s about %s: %v", r.Forge, r.Repository, err)
}

// gitSaid is git's own words for a failed command, or the error when it said nothing.
func gitSaid(out string, err error) string {
	if s := strings.TrimSpace(out); s != "" {
		return s
	}
	return err.Error()
}

// printForgeCheck prints r, one finding a line, or as JSON; the reason goes
// out through the error path.
func (a *app) printForgeCheck(r forgeCheck) error {
	if a.jsonOut {
		return a.printJSON(r)
	}
	fmt.Fprintf(a.out, "forge: %s\n", r.Forge)
	if r.Repository != "" {
		fmt.Fprintf(a.out, "repository: %s (from origin)\n", r.Repository)
	}
	fmt.Fprintf(a.out, "token: read from %s\n", r.TokenEnv)
	if r.asked {
		fmt.Fprintf(a.out, "read: %s\npush: %s\n", yesNo(r.Read), yesNo(r.Push))
	}
	return nil
}

// yesNo is b as yes or no.
func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}
