// Package sendback is what every send back of a story to its agent shares
// (ADR-0248, ADR-0302): the state the workflow sends a story back to, and the
// thread by flai, titled with storygit.SentBackPrefix, that tells its agent
// why and that flai serve starts the agent on. flai accept's send back, flai
// serve's Rebase and resolve job (S-0804), and the send back by hand, a
// story's move along the send-back transition by anyone but flai and its own
// agent (I-0397), all use it, so none imports another.
package sendback

import (
	"fmt"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workflow"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Target is the state w's transition with the send-back effect takes it to
// out of its state, else w's first state of the active kind: on the
// baseline, in-progress.
func Target(w workflow.Workflow, it *workitem.Item) string {
	for _, t := range w.From(it.Type, it.Status) {
		if t.HasEffect(workflow.EffectSendBack) {
			return t.To
		}
	}
	return workitem.FirstOfKind(w, workflow.KindActive)
}

// Thread opens the send back's thread on story, by storygit.SentBackAuthor,
// titled title and saying text at now, or adds text to the one still open
// from an earlier send back, which flai serve starts story's agent on
// (ADR-0248). title begins with storygit.SentBackPrefix.
func Thread(repo *workitem.Repo, story, title, text string, now time.Time) (*threads.Thread, error) {
	all, err := threads.For(repo, story)
	if err != nil {
		return nil, err
	}
	for _, th := range all {
		if th.Open() && th.Opener() == storygit.SentBackAuthor && strings.HasPrefix(th.Title, storygit.SentBackPrefix) {
			return threads.Reply(repo, th.ID, storygit.SentBackAuthor, text, now)
		}
	}
	return threads.New(repo, threads.NewOptions{Title: title, On: story, Author: storygit.SentBackAuthor, Text: text, Now: now})
}

// ByHandSuffix ends the title of a send back by hand's thread, which
// ByHandTitle words and IsByHand tells apart from the acceptance's and the
// rebase and resolve job's.
const ByHandSuffix = " by hand"

// ByHandTitle is the title of the thread of by's send back by hand of story
// from the state from to the state to: "Acceptance sent back: alex sent
// S-0823 back from review to in-progress by hand".
func ByHandTitle(by, story, from, to string) string {
	return fmt.Sprintf("%s%s sent %s back from %s to %s%s", storygit.SentBackPrefix, by, story, from, to, ByHandSuffix)
}

// IsByHand reports whether title is a send back by hand's, as ByHandTitle
// words it.
func IsByHand(title string) bool {
	_, ok := ByHandMover(title)
	return ok
}

// ByHandMover is who made the send back by hand that ByHandTitle titled
// title; false when title is not one.
func ByHandMover(title string) (string, bool) {
	rest, ok := strings.CutPrefix(title, storygit.SentBackPrefix)
	if !ok || !strings.HasSuffix(rest, ByHandSuffix) || !strings.Contains(rest, " back from ") {
		return "", false
	}
	i := strings.LastIndex(rest, " sent ")
	if i <= 0 {
		return "", false
	}
	return rest[:i], true
}

// OwnAgent reports whether by is the name flai serve gives story's agent,
// <name>-<story>, agent-S-nnnn unless the agent's configuration names it
// otherwise.
func OwnAgent(by, story string) bool {
	return strings.HasSuffix(by, "-"+story) && len(by) > len(story)+1
}

// ByHand opens the thread of a send back by hand, or adds to the one still
// open from an earlier send back, as Thread does, and reports whether it did:
// when story it moved from the state from to the state to along w's
// transition with the send-back effect, by anyone but storygit.SentBackAuthor,
// whose own send backs open their thread, and the story's own agent, which
// is on it already (I-0397). Its words name by, the move, and reason, and
// tell the story's agent to do what it asks in the story's worktree and move
// the story to from again.
func ByHand(repo *workitem.Repo, w workflow.Workflow, it *workitem.Item, from, to, by, reason string, now time.Time) (*threads.Thread, bool, error) {
	if it.Type != workitem.Story || by == storygit.SentBackAuthor || OwnAgent(by, it.ID) {
		return nil, false, nil
	}
	if t, ok := w.Find(it.Type, from, to); !ok || !t.HasEffect(workflow.EffectSendBack) {
		return nil, false, nil
	}
	th, err := Thread(repo, it.ID, ByHandTitle(by, it.ID, from, to), ByHandText(by, it.ID, from, to, reason), now)
	if err != nil {
		return nil, false, err
	}
	return th, true, nil
}

// ByHandText is what the thread of by's send back by hand of story from the
// state from to the state to says, with reason, or with none given.
func ByHandText(by, story, from, to, reason string) string {
	sent := fmt.Sprintf("%s sent %s back from %s to %s by hand", by, story, from, to)
	do := "read this, do what it asks"
	if reason = strings.TrimSpace(reason); reason == "" {
		sent, do = sent+" and gave no reason.", "ask "+by+" on this thread what to do, do it"
	} else {
		sent += ": " + reason
	}
	return fmt.Sprintf("%s\n\n%s's agent: %s in the story's worktree, and move %s to %s again; flai serve resolves this thread once it is back in %s.",
		sent, story, do, story, from, from)
}
