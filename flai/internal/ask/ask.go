// Package ask is the guard's hold-and-ask of ADR-0130, which ADR-0086,
// ADR-0097, and ADR-0124 describe for permission_prompt: it holds a write
// to a protected path, asks the story's owner or the project's owner on a
// thread on the story, and answers allow or deny within the caller's bound.
// Any harness's permission handling, or the guard's own hold, calls it.
package ask

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// allowWords are the first words of an answer that allow the write.
var allowWords = []string{"allow", "yes", "approve", "approved", "ok"}

// Request is a write held for an owner's answer.
type Request struct {
	// Story is the in-progress story whose worktree the write is in, and
	// Owner its owner.
	Story string
	Owner string
	// Agent is who asks: its own entries on the thread are never an answer.
	Agent string
	// Tool is the harness's name for the tool that writes, and Path the path
	// it writes, relative to the story's worktree.
	Tool string
	Path string
	// Why says why a person must answer, as the harness has it.
	Why string
	// Change shows what the write would change, under the question.
	Change string
}

// Text is the thread's first entry for r when the project's owner is
// projectOwner: who asks, why a person must answer, who may and how, and the
// change itself.
func (r Request) Text(projectOwner string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s asks to %s `%s` in %s's worktree. %s\n\n", r.Agent, r.Tool, r.Path, r.Story, r.Why)
	if who := answeredBy(r.Owner, projectOwner); who != "" {
		fmt.Fprintf(&b, "Reply `allow`, as %s, to let it write. Anything else refuses it, and your words go back to the agent as the reason; a reply by anyone else is not an answer.\n\n", who)
	} else {
		b.WriteString("Reply `allow` to let it write. Anything else refuses it, and your words go back to the agent as the reason.\n\n")
	}
	b.WriteString(r.Change)
	return b.String()
}

// answeredBy names who may answer, as the request says it: the story's owner,
// and the project's owner too when the two differ; empty when neither is named.
func answeredBy(owner, projectOwner string) string {
	switch {
	case owner != "" && projectOwner != "" && owner != projectOwner:
		return fmt.Sprintf("%s, the story's owner, or %s, the project's owner", owner, projectOwner)
	case owner != "":
		return owner + ", the story's owner"
	case projectOwner != "":
		return projectOwner + ", the project's owner"
	}
	return ""
}

// Answerers are who may answer a request: the story's owner and the
// project's owner, those named, each once (I-0081). None means anyone.
func Answerers(owner, projectOwner string) []string {
	var who []string
	for _, name := range []string{owner, projectOwner} {
		if name != "" && !slices.Contains(who, name) {
			who = append(who, name)
		}
	}
	return who
}

// Allowed reports whether an answer's first word lets the write through.
func Allowed(text string) bool {
	words := strings.Fields(strings.ToLower(text))
	if len(words) == 0 {
		return false
	}
	return slices.Contains(allowWords, strings.TrimRight(words[0], ".,;:!?"))
}

// Asker holds requests on the threads of one repository.
type Asker struct {
	Repo *workitem.Repo
	// Now stamps the thread's entries; nil is time.Now.
	Now func() time.Time
	// Poll is how often a held request reads its thread; zero is a second.
	Poll time.Duration
	// Closing, when it closes, ends every wait as the bound does; nil never
	// does.
	Closing <-chan struct{}
	// Logger records a thread it could not resolve; nil logs nothing.
	Logger *slog.Logger
}

// Hold asks on a thread on the story for an answer from the story's owner or
// the project's owner (ADR-0097; anyone but the agent when neither is named)
// and resolves the thread with what was decided; it answers whether the
// write is allowed and, when not, why. The agent's own entry and another
// agent's are not answers. It asks on the open thread in which the agent
// asked the same request before, and opens one only when there is none, so
// an answer given after an earlier call stopped waiting is taken (ADR-0124).
// Unanswered within bound, or when ctx or Closing ends the wait, it refuses
// the write and leaves the thread open.
func (a Asker) Hold(ctx context.Context, req Request, bound time.Duration) (allowed bool, reason string) {
	text := req.Text(a.Repo.Manifest.Owner)
	th, asked, err := a.askedBefore(req.Story, req.Agent, text)
	if err != nil {
		return false, fmt.Sprintf("could not read %s's threads for an earlier ask of this write: %v", req.Story, err)
	}
	if th == nil {
		th, err = threads.New(a.Repo, threads.NewOptions{
			Title:  fmt.Sprintf("Allow %s %s?", req.Tool, req.Path),
			On:     req.Story,
			Author: req.Agent,
			Text:   text,
			Now:    a.now(),
		})
		if err != nil {
			return false, fmt.Sprintf("could not ask the operator on a thread: %v", err)
		}
		_ = a.mirror(th)
		// Only entries after the opening one can answer: the content shown may
		// itself hold lines that read as entries.
		asked = len(th.Entries())
	}
	answer, ok := a.await(ctx, th.ID, asked, req.Agent, Answerers(req.Owner, a.Repo.Manifest.Owner), bound)
	if !ok {
		return false, fmt.Sprintf("No answer on %s yet, so %s was not written. The thread stays open, and the owner's answer is taken when you make the same write again. Go on with work that does not need this write. When nothing else is left, write your narrative's Current state and Next steps naming %s, and end: flai serve starts you again when it is answered.", th.ID, req.Path, th.ID)
	}
	if Allowed(answer.Text) {
		a.settle(th.ID, req.Agent, fmt.Sprintf("allowed by %s: %s may %s %s", answer.Author, req.Agent, req.Tool, req.Path))
		return true, ""
	}
	a.settle(th.ID, req.Agent, fmt.Sprintf("refused by %s: %s was not written", answer.Author, req.Path))
	return false, "the operator refused: " + answer.Text
}

func (a Asker) now() time.Time {
	if a.Now == nil {
		return time.Now()
	}
	return a.Now()
}

// mirror copies the thread into its story's narrative.
func (a Asker) mirror(th *threads.Thread) error {
	if story := threads.StoryOf(a.Repo, th); story != "" {
		return threads.MirrorNarrative(a.Repo, story)
	}
	return nil
}

// askedBefore is the open thread on the story whose opening entries are the
// agent's request text, and how many entries that text reads as; nil when the
// agent has not asked it.
func (a Asker) askedBefore(story, agent, text string) (*threads.Thread, int, error) {
	all, err := threads.For(a.Repo, story)
	if err != nil {
		return nil, 0, err
	}
	for _, th := range all {
		if !th.Open() {
			continue
		}
		opening := openingEntries(th.Created, agent, text)
		if entries := th.Entries(); len(entries) >= len(opening) && sameEntries(entries[:len(opening)], opening) {
			return th, len(opening), nil
		}
	}
	return nil, 0, nil
}

// openingEntries are the entries a thread that author opened at created with
// text reads as, its body opening as threads.New writes it: one, unless the
// content shown holds lines that read as entries.
func openingEntries(created, author, text string) []threads.Entry {
	body := fmt.Sprintf("### %s %s\n%s\n", created, author, strings.TrimSpace(text))
	return (&threads.Thread{Body: body}).Entries()
}

func sameEntries(x, y []threads.Entry) bool {
	return slices.EqualFunc(x, y, func(p, q threads.Entry) bool {
		return p.At == q.At && p.Author == q.Author && p.Text == q.Text && p.Recommendation == q.Recommendation
	})
}

// await reads the thread, at once and then every Poll, until an entry after
// the first asked ones is by one of who, never the agent (by anyone but the
// agent when who is empty), ctx or Closing ends, or bound passes.
func (a Asker) await(ctx context.Context, id string, asked int, agent string, who []string, bound time.Duration) (threads.Entry, bool) {
	poll := a.Poll
	if poll <= 0 {
		poll = time.Second
	}
	limit := time.NewTimer(bound)
	defer limit.Stop()
	tick := time.NewTicker(poll)
	defer tick.Stop()
	for {
		// An error is the thread being written, or gone for a moment: it is
		// read again on the next tick.
		if th, err := threads.Get(a.Repo, id); err == nil {
			if answer, ok := answerOn(th, asked, agent, who); ok {
				return answer, true
			}
		}
		select {
		case <-ctx.Done():
			return threads.Entry{}, false
		case <-a.Closing: // nil, and so never ready, unless the Asker was given one
			return threads.Entry{}, false
		case <-limit.C:
			return threads.Entry{}, false
		case <-tick.C:
		}
	}
}

// answerOn is the first entry after the first asked ones by one of who, never
// the agent (by anyone but the agent when who is empty).
func answerOn(th *threads.Thread, asked int, agent string, who []string) (threads.Entry, bool) {
	entries := th.Entries()
	for i := asked; i < len(entries); i++ {
		if by := entries[i].Author; by != agent && (len(who) == 0 || slices.Contains(who, by)) {
			return entries[i], true
		}
	}
	return threads.Entry{}, false
}

// settle resolves the thread as the agent, unless the operator resolved it.
func (a Asker) settle(id, agent, reason string) {
	th, err := threads.Get(a.Repo, id)
	if err != nil || !th.Open() {
		return
	}
	if th, err = threads.Resolve(a.Repo, id, agent, reason, a.now()); err == nil {
		_ = a.mirror(th)
	} else if a.Logger != nil {
		a.Logger.Warn("permission thread not resolved", "component", "ask", "agent", agent, "thread", id, "err", err.Error())
	}
}
