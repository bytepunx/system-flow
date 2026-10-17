package serve

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/hostapi"
	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Reading a story's pull request back (S-0394, ADR-0167). After the forge
// pass pushes a story in review, it reads the pull request's head commit,
// its required checks, and its reviews, and keeps them, with an event for
// each change, in .flai-cache/forge/status/S-nnnn.json. A required check
// that failed at the head, or a reviewer whose latest deciding review
// requests changes, sends the story back: one thread on the story for
// everything new at that read, a SendBack kept for each failure so that it
// never sends the story back twice, and a move from review to in-progress
// by flai. The thread is resolved once the story is back in review at
// another commit, or is done or cancelled.

const (
	// sendBackBy moves a story back from review to in-progress.
	sendBackBy = "flai"
	// quoteLimit is how many characters of a review's comment, or a check's
	// summary, a send-back thread quotes.
	quoteLimit = 500
)

// readLook is what one look's reads of pull requests did, logged once when
// the look ends.
type readLook struct {
	failed   []string
	errs     []string
	sentBack []string
}

func (r *readLook) log(logger *slog.Logger, root string) {
	if len(r.failed) > 0 {
		logger.Warn("pull requests in review not read", "component", "serve", "root", root,
			"failed", strings.Join(r.failed, ","), "err", strings.Join(r.errs, "; "))
	}
	if len(r.sentBack) > 0 {
		logger.Info("stories sent back to in-progress", "component", "serve", "root", root, "stories", strings.Join(r.sentBack, ","))
	}
}

// readBack reads the story's pull request into the kept status, resolves the
// send-back threads of an earlier commit, and sends the story back on a
// failure not sent back for before; a read that fails touches nothing.
func (p *forgePass) readBack(ctx context.Context, repo *workitem.Repo, it *workitem.Item, token ForgeToken, look *readLook) {
	st, err := p.readPullRequest(ctx, repo, it.ID, token)
	if err != nil {
		look.failed = append(look.failed, it.ID)
		look.errs = append(look.errs, it.ID+": "+err.Error())
		return
	}
	why := fmt.Sprintf("%s is back in review at %s, so flai serve reads its checks and reviews again.", it.ID, short(st.HeadSHA))
	p.resolveSendBacks(repo, st, func(b forge.SendBack) bool { return b.Commit != st.HeadSHA }, why)
	if p.sendBack(repo, st) {
		look.sentBack = append(look.sentBack, it.ID)
	}
}

// readPullRequest reads story's pull request, its required checks on the
// head commit, and its reviews from the forge, and keeps them.
func (p *forgePass) readPullRequest(ctx context.Context, repo *workitem.Repo, story string, token ForgeToken) (forge.Status, error) {
	push := ForgePush{Repo: repo, Story: story, Git: gitOf(p.o), Token: token, NewForge: p.o.NewForge}
	f, err := push.adapter(push.git(), mainCheckout(repo), story)
	if err != nil {
		return forge.Status{}, err
	}
	was, _, err := forge.ReadStatus(repo.CacheDir(), story)
	if err != nil {
		return forge.Status{}, err
	}
	fctx, cancel := context.WithTimeout(ctx, forgeTimeout)
	defer cancel()
	pr, found, err := f.FindPullRequest(fctx, storygit.Branch(story))
	if err != nil {
		return forge.Status{}, readReason(f, story, err)
	}
	if !found {
		return forge.Status{}, fmt.Errorf("%s has no open pull request on %s %s to read; flai serve opens it at its next look", story, f.Name(), f.Repository())
	}
	checks, err := f.Checks(fctx, pr)
	if err != nil {
		return forge.Status{}, readReason(f, story, err)
	}
	reviews, err := f.Reviews(fctx, pr)
	if err != nil {
		return forge.Status{}, readReason(f, story, err)
	}
	was.Story = story
	next := was.Record(pr, checks, reviews, p.o.Now())
	return next, forge.WriteStatus(repo.CacheDir(), next)
}

// readReason says that the forge did not answer for story's pull request;
// a token or access failure has the push's failure thread already.
func readReason(f forge.Forge, story string, err error) error {
	return fmt.Errorf("could not read the checks and reviews of %s's pull request on %s %s (%w); flai serve reads them again at its next look", story, f.Name(), f.Repository(), err)
}

// failure is one failure at a pull request's head that sends its story back.
type failure struct {
	key  string
	kind string // forge.EventCheck or forge.EventReview
	name string // the check's name or the reviewer's login
	url  string // the check's log or the review's page
	text string // the check's summary or the review's comment
}

// failures are the failures at st's head commit that flai has not sent the
// story back for: each required check that failed, and each reviewer whose
// latest review that approves, requests changes, or was dismissed requests
// changes, as the forge counts a reviewer's decision.
func failures(st forge.Status) []failure {
	var out []failure
	add := func(f failure) {
		if _, sent := st.Sent(f.key); !sent {
			out = append(out, f)
		}
	}
	for _, c := range st.Checks {
		if c.State == forge.CheckFailed {
			add(failure{key: forge.CheckKey(c.Name, st.HeadSHA), kind: forge.EventCheck, name: c.Name, url: c.URL, text: c.Summary})
		}
	}
	latest := map[string]string{}
	for _, r := range st.Reviews {
		switch r.State {
		case forge.ReviewApproved, forge.ReviewChangesRequested, forge.ReviewDismissed:
			latest[r.Author] = r.URL
		}
	}
	for _, r := range st.Reviews {
		if r.State == forge.ReviewChangesRequested && latest[r.Author] == r.URL {
			add(failure{key: forge.ReviewKey(r.URL), kind: forge.EventReview, name: r.Author, url: r.URL, text: r.Body})
		}
	}
	return out
}

// sendBack opens one thread for the failures at st's head not sent back
// for before, records them, and moves the story to in-progress; with none,
// it moves again a story whose move after an earlier send back at this head
// failed. It answers whether the story was moved.
func (p *forgePass) sendBack(repo *workitem.Repo, st forge.Status) bool {
	it, err := repo.Get(st.Story)
	if err != nil || it.Status != workitem.Review {
		return false
	}
	now := p.o.Now()
	backs := unmoved(repo, it, st)
	if fresh := failures(st); len(fresh) > 0 {
		th, err := p.openSendBack(repo, st, fresh, now)
		if err != nil {
			p.o.Logger.Warn("pull request failure thread not opened", "component", "serve", "root", p.e.Root, "story", st.Story, "err", err.Error())
			return false
		}
		at := now.UTC().Format(time.RFC3339)
		for _, f := range fresh {
			b := forge.SendBack{Key: f.key, Kind: f.kind, Name: f.name, Commit: st.HeadSHA, URL: f.url, Thread: th.ID, At: at}
			st.SentBack = append(st.SentBack, b)
			backs = append(backs, b)
		}
		if err := forge.WriteStatus(repo.CacheDir(), st); err != nil {
			p.o.Logger.Warn("pull request failure not recorded", "component", "serve", "root", p.e.Root, "story", st.Story, "thread", th.ID, "err", err.Error())
			return false
		}
	}
	if len(backs) == 0 {
		return false
	}
	reason := sendBackReason(st.PullRequest, backs)
	if _, err := repo.Transition(it, workitem.InProgress, sendBackBy, reason, now); err != nil {
		p.o.Logger.Warn("story not sent back to in-progress", "component", "serve", "root", p.e.Root, "story", st.Story, "err", err.Error())
		return false
	}
	if p.o.Host.Record != nil {
		p.o.Host.Record(hostapi.Entry{At: now.UTC().Format(time.RFC3339), Action: hostapi.ActionForge, Method: "serve.forge", Project: p.e.Key, Root: p.e.Root, By: forgeAuthor,
			Outcome: "done", Detail: fmt.Sprintf("sent %s back to in-progress: %s", st.Story, reason)})
	}
	return true
}

// unmoved are the send backs at st's head whose thread is open and after
// which it never moved to in-progress: its move failed after the thread
// was opened.
func unmoved(repo *workitem.Repo, it *workitem.Item, st forge.Status) []forge.SendBack {
	var out []forge.SendBack
	for _, b := range st.SentBack {
		if b.Commit == st.HeadSHA && openThread(repo, b.Thread) && !movedSince(it, b.At) {
			out = append(out, b)
		}
	}
	return out
}

// movedSince is whether it moved to in-progress at or after at.
func movedSince(it *workitem.Item, at string) bool {
	since, err := time.Parse(time.RFC3339, at)
	if err != nil {
		return true
	}
	for _, tr := range it.Transitions {
		if tr.To != workitem.InProgress {
			continue
		}
		if t, err := time.Parse(time.RFC3339, tr.At); err == nil && !t.Before(since) {
			return true
		}
	}
	return false
}

// openThread is whether id names a thread flai serve opened that is open.
func openThread(repo *workitem.Repo, id string) bool {
	if id == "" {
		return false
	}
	th, err := threads.Get(repo, id)
	return err == nil && th.Open() && th.Opener() == forgeAuthor
}

// openSendBack opens the one thread on st's story for its fresh failures.
// A review's comment or a check's summary is the forge's text, not flai's:
// when the project's lint refuses it quoted, it is opened with each in a
// code block instead.
func (p *forgePass) openSendBack(repo *workitem.Repo, st forge.Status, fresh []failure, now time.Time) (*threads.Thread, error) {
	opt := threads.NewOptions{Title: sendBackTitle(st, fresh), On: st.Story, Author: forgeAuthor, Text: sendBackText(st, fresh, false), Now: now}
	th, err := threads.New(repo, opt)
	var lint *mdlint.Error
	if errors.As(err, &lint) {
		opt.Text = sendBackText(st, fresh, true)
		return threads.New(repo, opt)
	}
	return th, err
}

// resolveSendBacks resolves each open thread flai serve opened for a send
// back of st that pick picks, saying why; one someone resolved is left so.
func (p *forgePass) resolveSendBacks(repo *workitem.Repo, st forge.Status, pick func(forge.SendBack) bool, why string) {
	done := map[string]bool{}
	for _, b := range st.SentBack {
		if done[b.Thread] || !pick(b) || !openThread(repo, b.Thread) {
			continue
		}
		done[b.Thread] = true
		if _, err := threads.Resolve(repo, b.Thread, forgeAuthor, why, p.o.Now()); err != nil {
			p.o.Logger.Warn("pull request failure thread not resolved", "component", "serve", "root", p.e.Root, "story", st.Story, "thread", b.Thread, "err", err.Error())
		}
	}
}

// closedSendBacks resolves the open send-back threads of each story that is
// done or cancelled, archived or not: flai serve reads its pull request no
// more. A story in progress keeps its thread, which its agent is working.
func (p *forgePass) closedSendBacks(repo *workitem.Repo, items []*workitem.Item) {
	status := map[string]string{}
	for _, it := range items {
		if it.Type == workitem.Story {
			status[it.ID] = it.Status
		}
	}
	files, _ := filepath.Glob(forge.StatusPath(repo.CacheDir(), "*"))
	for _, f := range files {
		id := strings.TrimSuffix(filepath.Base(f), ".json")
		st, ok, err := forge.ReadStatus(repo.CacheDir(), id)
		if err != nil || !ok || !slices.ContainsFunc(st.SentBack, func(b forge.SendBack) bool { return openThread(repo, b.Thread) }) {
			continue
		}
		s, listed := status[id]
		if !listed {
			it, err := repo.Get(id)
			if err != nil {
				continue
			}
			s = it.Status
		}
		if s == workitem.Done || s == workitem.Cancelled {
			why := fmt.Sprintf("%s is %s, so flai serve reads its pull request no more.", id, s)
			p.resolveSendBacks(repo, st, func(forge.SendBack) bool { return true }, why)
		}
	}
}

// OpenSendBack is the latest failure flai sent story back to in-progress
// for, when the thread flai serve opened for it is still open; false when
// there is none, or its thread is resolved.
func OpenSendBack(repo *workitem.Repo, story string) (forge.SendBack, bool, error) {
	st, ok, err := forge.ReadStatus(repo.CacheDir(), story)
	if err != nil || !ok || len(st.SentBack) == 0 {
		return forge.SendBack{}, false, err
	}
	latest := st.SentBack[len(st.SentBack)-1]
	if !openThread(repo, latest.Thread) {
		return forge.SendBack{}, false, nil
	}
	return latest, true, nil
}

// sendBackTitle names the pull request and what failed on it.
func sendBackTitle(st forge.Status, fresh []failure) string {
	var checks, reviews []failure
	for _, f := range fresh {
		if f.kind == forge.EventCheck {
			checks = append(checks, f)
		} else {
			reviews = append(reviews, f)
		}
	}
	var what string
	switch {
	case len(checks) == 1 && len(reviews) == 0:
		what = fmt.Sprintf("the required check %s failed", checks[0].name)
	case len(checks) == 0 && len(reviews) == 1:
		what = fmt.Sprintf("%s requested changes", reviews[0].name)
	case len(reviews) == 0:
		what = fmt.Sprintf("%d required checks failed", len(checks))
	case len(checks) == 0:
		what = fmt.Sprintf("%d reviews requested changes", len(reviews))
	default:
		what = fmt.Sprintf("%s failed and %s requested changes", count(len(checks), "required check"), count(len(reviews), "review"))
	}
	return fmt.Sprintf("Pull request #%d of %s: %s", st.PullRequest, st.Story, what)
}

// count is n things, "s" added for more than one.
func count(n int, thing string) string {
	if n == 1 {
		return "1 " + thing
	}
	return fmt.Sprintf("%d %ss", n, thing)
}

// sendBackText names each failure, with a link to its log or its review and
// a brief quote of what the forge said, and what happens now; fenced puts
// each quote in a code block.
func sendBackText(st forge.Status, fresh []failure, fenced bool) string {
	var b strings.Builder
	pr := fmt.Sprintf("pull request #%d", st.PullRequest)
	fmt.Fprintf(&b, "flai serve read %s of %s at %s, and sends the story back.\n\n", link(pr, st.URL), st.Story, codeSpan(short(st.HeadSHA)))
	for _, f := range fresh {
		var head string
		if f.kind == forge.EventCheck {
			head = fmt.Sprintf("The required check %s failed at %s", codeSpan(f.name), codeSpan(short(st.HeadSHA)))
			if f.url != "" {
				head += fmt.Sprintf(" (%s)", link("its log", f.url))
			}
		} else {
			head = f.name + " requested changes"
			if f.url != "" {
				head += fmt.Sprintf(" (%s)", link("the review", f.url))
			}
		}
		switch {
		case strings.TrimSpace(f.text) != "":
			fmt.Fprintf(&b, "%s:\n\n%s\n\n", head, quote(f.text, fenced))
		case f.kind == forge.EventReview:
			fmt.Fprintf(&b, "%s, with no comment.\n\n", head)
		default:
			fmt.Fprintf(&b, "%s.\n\n", head)
		}
	}
	fmt.Fprintf(&b, "flai moved %[1]s back to in-progress. Fix the branch, close out, and move the story to review: flai serve pushes the fix, and resolves this thread once the story is back in review at a new commit.", st.Story)
	return b.String()
}

// sendBackReason names each failure the story is sent back for, for its move.
func sendBackReason(pr int, backs []forge.SendBack) string {
	parts := make([]string, 0, len(backs))
	for _, b := range backs {
		if b.Kind == forge.EventCheck {
			parts = append(parts, fmt.Sprintf("the required check %s failed on %s of pull request #%d", b.Name, short(b.Commit), pr))
		} else {
			parts = append(parts, fmt.Sprintf("%s requested changes on pull request #%d", b.Name, pr))
		}
	}
	return strings.Join(parts, "; ")
}

// quote is the forge's text s, trimmed to quoteLimit characters: a
// blockquote of its paragraphs, each on one line with its bare URLs in code
// spans, or, fenced, a code block that keeps it as it is.
func quote(s string, fenced bool) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\r\n", "\n"))
	if r := []rune(s); len(r) > quoteLimit {
		s = strings.TrimSpace(string(r[:quoteLimit])) + "…"
	}
	if fenced {
		fence := "```"
		for strings.Contains(s, fence) {
			fence += "`"
		}
		return fence + "text\n" + s + "\n" + fence
	}
	var paras, lines []string
	flush := func() {
		if len(lines) > 0 {
			para := mdlint.QuoteBareURLs(strings.Join(strings.Fields(strings.Join(lines, " ")), " "))
			if strings.HasPrefix(para, "#") {
				para = `\` + para
			}
			paras = append(paras, "> "+para)
			lines = nil
		}
	}
	for l := range strings.SplitSeq(s, "\n") {
		if strings.TrimSpace(l) == "" {
			flush()
			continue
		}
		lines = append(lines, l)
	}
	flush()
	return strings.Join(paras, "\n>\n")
}

// link is text linked to url, or text alone when there is no url.
func link(text, url string) string {
	if url == "" {
		return text
	}
	return "[" + text + "](" + url + ")"
}

// codeSpan is s in a code span whose fence is longer than any run of
// backticks in s.
func codeSpan(s string) string {
	fence := "`"
	for strings.Contains(s, fence) {
		fence += "`"
	}
	if strings.HasPrefix(s, "`") || strings.HasSuffix(s, "`") {
		s = " " + s + " "
	}
	return fence + s + fence
}
