// Package sendback is what every send back of a story to its agent shares
// (ADR-0248, ADR-0302): the state the workflow sends a story back to, and the
// thread by flai, titled with storygit.SentBackPrefix, that tells its agent
// why and that flai serve starts the agent on. flai accept's send back and
// flai serve's Rebase and resolve job (S-0804) both use it, so neither
// imports the other.
package sendback

import (
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
