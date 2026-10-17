package forge

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
)

// Status is a story's pull request as flai serve last read it from the
// forge, kept in .flai-cache/forge/status/S-nnnn.json; it never holds the token.
type Status struct {
	// Story is the story's ID.
	Story string `json:"story"`
	// PullRequest is the pull request's number in its repository.
	PullRequest int `json:"pull_request"`
	// URL is the pull request's page on the forge.
	URL string `json:"url"`
	// HeadSHA is the commit at the head of the pull request when it was read.
	HeadSHA string `json:"head_sha"`
	// ReadAt is when it was last read, RFC 3339.
	ReadAt string `json:"read_at"`
	// Checks are the required checks only, in the forge's order.
	Checks []CheckStatus `json:"checks"`
	// Reviews are every review, oldest first.
	Reviews []ReviewStatus `json:"reviews"`
	// RequiredApprovals is the approving reviews the base branch requires, as last read.
	RequiredApprovals int `json:"required_approvals"`
	// Events are each change, oldest first, appended and never rewritten.
	Events []StatusEvent `json:"events,omitempty"`
	// SentBack are the failures flai sent the story back to in-progress for.
	SentBack []SendBack `json:"sent_back,omitempty"`
}

// CheckStatus is one required check on the pull request's head commit, as last read.
type CheckStatus struct {
	// Name is the check's name, as the base branch's rules require it.
	Name string `json:"name"`
	// State is whether it is pending, passed, or failed.
	State CheckState `json:"state"`
	// URL is the page with its details and log; empty when the forge has none.
	URL string `json:"url,omitempty"`
	// Summary is the forge's one line on what it found; may be empty.
	Summary string `json:"summary,omitempty"`
	// Since is when it entered State, RFC 3339.
	Since string `json:"since"`
}

// ReviewStatus is one review of the pull request, as last read.
type ReviewStatus struct {
	// Author is the reviewer's login.
	Author string `json:"author"`
	// State is what the review says.
	State ReviewState `json:"state"`
	// URL is the review's page on the forge.
	URL string `json:"url,omitempty"`
	// Body is the review's comment; may be empty.
	Body string `json:"body,omitempty"`
	// SubmittedAt is when it was submitted, RFC 3339; empty while pending.
	SubmittedAt string `json:"submitted_at,omitempty"`
}

// The kinds of change a StatusEvent records.
const (
	EventHead   = "head"
	EventCheck  = "check"
	EventReview = "review"
)

// StatusEvent is one change to a story's pull request, as flai serve saw it.
type StatusEvent struct {
	// At is when it was seen, RFC 3339.
	At string `json:"at"`
	// Kind is EventHead, EventCheck, or EventReview.
	Kind string `json:"kind"`
	// Name is the check's name or the reviewer's login; empty for a head.
	Name string `json:"name,omitempty"`
	// From is the state, or for a head the commit, before; empty when first seen.
	From string `json:"from,omitempty"`
	// To is the state, or for a head the commit, after.
	To string `json:"to"`
	// Commit is the head commit it was seen at.
	Commit string `json:"commit"`
}

// SendBack is a failure flai sent the story back to in-progress for.
type SendBack struct {
	// Key is the failure's CheckKey or ReviewKey.
	Key string `json:"key"`
	// Kind is EventCheck or EventReview.
	Kind string `json:"kind"`
	// Name is the check's name or the reviewer's login.
	Name string `json:"name"`
	// Commit is the head commit it failed at.
	Commit string `json:"commit"`
	// URL is the check's log or the review's page.
	URL string `json:"url,omitempty"`
	// Thread is the thread flai opened on the story for it.
	Thread string `json:"thread,omitempty"`
	// At is when flai sent the story back, RFC 3339.
	At string `json:"at"`
}

// StatusPath is where flai keeps story's pull request status under cacheDir.
func StatusPath(cacheDir, story string) string {
	return filepath.Join(cacheDir, "forge", "status", story+".json")
}

// ReadStatus is story's pull request as flai serve last read it, and false
// with no error when it has not read it.
func ReadStatus(cacheDir, story string) (Status, bool, error) {
	path := StatusPath(cacheDir, story)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Status{}, false, nil
	}
	if err != nil {
		return Status{}, false, fmt.Errorf("read the pull request status of %s at %s: %w; check that the project's .flai-cache is readable", story, path, err)
	}
	var st Status
	if err := json.Unmarshal(data, &st); err != nil {
		return Status{}, false, fmt.Errorf("read the pull request status of %s at %s: %w; delete it, and flai serve reads the pull request afresh", story, path, err)
	}
	return st, true, nil
}

// WriteStatus keeps st as its story's pull request status under cacheDir,
// replacing the file whole.
func WriteStatus(cacheDir string, st Status) error {
	path := StatusPath(cacheDir, st.Story)
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return fmt.Errorf("record the pull request status of %s: %w", st.Story, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("record the pull request status of %s: %w; check that the project's .flai-cache is writable", st.Story, err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("record the pull request status of %s at %s: %w; check that the project's .flai-cache is writable", st.Story, path, err)
	}
	return nil
}

// Record is s with a fresh read of its pull request from the forge, at
// now, and an event appended for each change since s.
func (s Status) Record(pr PullRequest, checks []Check, reviews []Review, now time.Time) Status {
	at := now.UTC().Format(time.RFC3339)
	next := Status{
		Story:       s.Story,
		PullRequest: pr.Number,
		URL:         pr.URL,
		HeadSHA:     pr.HeadSHA,
		ReadAt:      at,
		Checks:      make([]CheckStatus, 0, len(checks)),
		Reviews:     make([]ReviewStatus, 0, len(reviews)),
		Events:      slices.Clone(s.Events),
		SentBack:    slices.Clone(s.SentBack),
	}
	event := func(kind, name, from, to string) {
		next.Events = append(next.Events, StatusEvent{At: at, Kind: kind, Name: name, From: from, To: to, Commit: pr.HeadSHA})
	}
	if pr.HeadSHA != s.HeadSHA {
		event(EventHead, "", s.HeadSHA, pr.HeadSHA)
	}
	wasCheck := make(map[string]CheckStatus, len(s.Checks))
	for _, c := range s.Checks {
		wasCheck[c.Name] = c
	}
	for _, c := range checks {
		if !c.Required {
			continue
		}
		cs := CheckStatus{Name: c.Name, State: c.State, URL: c.URL, Summary: c.Summary, Since: at}
		was, seen := wasCheck[c.Name]
		if seen && was.State == c.State {
			cs.Since = was.Since
		} else {
			event(EventCheck, c.Name, string(was.State), string(c.State))
		}
		next.Checks = append(next.Checks, cs)
	}
	wasReview := make(map[string]ReviewStatus, len(s.Reviews))
	for _, r := range s.Reviews {
		wasReview[r.URL] = r
	}
	for _, r := range reviews {
		was, seen := wasReview[r.URL]
		if !seen || was.State != r.State {
			event(EventReview, r.Author, string(was.State), string(r.State))
		}
		next.Reviews = append(next.Reviews, ReviewStatus{Author: r.Author, State: r.State, URL: r.URL, Body: r.Body, SubmittedAt: r.SubmittedAt})
	}
	return next
}

// Decisions are each reviewer's latest review that approves, requests
// changes, or was dismissed, as the forge counts a reviewer's decision: a
// later review that only comments does not change it. One per reviewer, in
// the order of each reviewer's first review; a reviewer who only commented
// has none.
func (s Status) Decisions() []ReviewStatus {
	var order []string
	latest := map[string]ReviewStatus{}
	seen := map[string]bool{}
	for _, r := range s.Reviews {
		if !seen[r.Author] {
			seen[r.Author] = true
			order = append(order, r.Author)
		}
		switch r.State {
		case ReviewApproved, ReviewChangesRequested, ReviewDismissed:
			latest[r.Author] = r
		}
	}
	var out []ReviewStatus
	for _, author := range order {
		if r, ok := latest[author]; ok {
			out = append(out, r)
		}
	}
	return out
}

// CheckKey is the key a failed check is sent back under: the check at one commit.
func CheckKey(name, commit string) string {
	return "check:" + name + "@" + commit
}

// ReviewKey is the key a review requesting changes is sent back under: the review itself.
func ReviewKey(url string) string {
	return "review:" + url
}

// Sent is the send back recorded for key, and whether there is one.
func (s Status) Sent(key string) (SendBack, bool) {
	for _, b := range s.SentBack {
		if b.Key == key {
			return b, true
		}
	}
	return SendBack{}, false
}
