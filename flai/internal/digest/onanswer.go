package digest

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Actions OnAnswer reports in AnswerEntry.Action.
const (
	ActionWritten = "written"
	ActionUpdated = "updated"
)

// AnswerEntry is what OnAnswer did with an operator's answer: the entry it
// wrote or updated, or, when it wrote nothing, a warning saying why.
type AnswerEntry struct {
	// ID and Path name the entry, its path from the main checkout's root,
	// slash separated; empty with a warning.
	ID   string `json:"id,omitempty"`
	Path string `json:"path,omitempty"`
	// Action is ActionWritten or ActionUpdated; empty with a warning.
	Action string `json:"action,omitempty"`
	// Warning says what failed, on what, and what takes the answer in
	// instead; empty when an entry was written or updated.
	Warning string `json:"warning,omitempty"`
}

// OnAnswer writes the digest entry for the thread's newest entry when it is
// the operator's answer by Answers' rule, in the main checkout's digest
// folder, and is called once the thread is saved. With no entry citing the
// thread, or only absorbed or compacted ones, it writes one as Seed would;
// with an active one, the highest numbered, it appends the answer to that
// entry's body. It is nil when the newest entry is no answer, or when an
// entry citing the thread already holds the answer's time. It never fails
// the answer: what goes wrong is the Warning, for the caller to log.
func OnAnswer(repo *workitem.Repo, r execx.Runner, th *threads.Thread) *AnswerEntry {
	if repo == nil || th == nil {
		return nil
	}
	es := th.Entries()
	if len(es) == 0 {
		return nil
	}
	main := mainRepo(repo)
	owner := main.Manifest.Owner
	answers := Answers(th, owner, Lookup(main))
	if len(answers) == 0 || answers[len(answers)-1].Ref != es[len(es)-1].At {
		return nil
	}
	c := answers[len(answers)-1]
	warn := func(what string) *AnswerEntry {
		return &AnswerEntry{Warning: fmt.Sprintf("no digest entry was written for the operator's answer on %s at %s: %s", c.Source, c.Ref, what)}
	}
	failed := func(err error) *AnswerEntry {
		var refused *RefusedError
		if errors.As(err, &refused) {
			return warn(fmt.Sprintf("the check refused it under %s; the answer is saved on the thread, and the next flai digest seed or curator run takes it in", brokenRules(refused.Problems)))
		}
		return warn(err.Error() + "; the answer is saved on the thread, and the next flai digest seed or curator run takes it in")
	}
	secret := func() *AnswerEntry {
		return warn("it has a secret's shape (rule " + RuleSecret + "), so it is left on the thread alone; an entry says where a secret lives, never its value")
	}

	entries, err := List(main)
	if err != nil {
		return failed(err)
	}
	if skipReason(c, entries) != "" {
		return nil
	}
	env := EnvFor(main)
	if e := answeredIn(entries, c.Source); e != nil {
		block := answerBlock(c, owner)
		if len(Secrets(block)) > 0 {
			return secret()
		}
		body := withBlock(e.Body, block)
		res, err := Edit(main, e.ID, EditOptions{Body: &body, Env: env})
		if err != nil {
			return failed(err)
		}
		return &AnswerEntry{ID: res.ID, Path: res.Path, Action: ActionUpdated}
	}
	title, body := seedTitle(c.Title), seedBody(c, owner)
	if len(Secrets(title+"\n"+body)) > 0 {
		return secret()
	}
	res, err := New(main, r, NewOptions{Title: title, Kind: c.Kind, Sources: []string{c.Source}, Topics: c.Topics, Paths: c.Paths, Body: body, Now: learned(c), Env: env})
	if err != nil {
		return failed(err)
	}
	return &AnswerEntry{ID: res.ID, Path: res.Path, Action: ActionWritten}
}

// mainRepo is repo with its root at the main checkout, where threads and
// the entries their answers become are written, when repo is a worktree.
func mainRepo(repo *workitem.Repo) *workitem.Repo {
	if !repo.IsWorktree() {
		return repo
	}
	main := *repo
	main.Root = repo.MainRoot
	return &main
}

// answeredIn is the entry a later answer on the thread source is added to:
// the highest numbered active entry whose sources name it, the newest, which
// a curator's merge of the entries citing it leaves; nil when there is none.
func answeredIn(entries []Entry, source string) *Entry {
	source = canonicalSource(source)
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		if e.Status == "active" && slices.Contains(clean(e.Sources, canonicalSource), source) {
			return &entries[i]
		}
	}
	return nil
}

// answerBlock is a later answer as an entry already citing its thread takes
// it in: a sentence naming the operator and the answer's time, and the
// answer, with the recommendation it overruled, each fenced.
func answerBlock(c Candidate, owner string) string {
	blocks := []string{"The operator, " + owner + ", answered " + c.Source + " again at " + c.Ref + "."}
	if c.Kind == "recommendation" {
		blocks = append(blocks, "The recommendation it overruled:", fenced(c.Overruled), "The operator's answer:", fenced(c.Text))
	} else {
		blocks = append(blocks, "The answer:", fenced(c.Text))
	}
	return strings.Join(blocks, "\n\n")
}

// withBlock is body with block added as its last paragraphs, before the
// source line a body ends with, "Source: ..." or "Sources: ...", when it has
// one, so the line still closes it.
func withBlock(body, block string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return block
	}
	if i := strings.LastIndex(body, "\n\n"); i >= 0 {
		last := body[i+2:]
		if !strings.Contains(last, "\n") && (strings.HasPrefix(last, "Source: ") || strings.HasPrefix(last, "Sources: ")) {
			return body[:i] + "\n\n" + block + "\n\n" + last
		}
	}
	return body + "\n\n" + block
}
