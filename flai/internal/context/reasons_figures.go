package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/release"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The paths of the figures the release and ready reasons load, which are no
// files: each is the MCP tool that gives the same figures afresh
// (ADR-0217 §3).
const (
	FigureRelease    = "release_evaluate"
	FigureBoard      = "board"
	FigureCandidates = "promote_candidates"
)

// loadReview loads a review reason: the story whole and its last flai verify
// result, the report flai verify stores in the project's cache at
// .flai-cache/verify/<story>.json, as flai verify --last --json prints it,
// or, when it has none, an item that says so and how to run it. The report
// is read as the file it is stored in: the verify package imports this one,
// through check, so the reading of verify.LastReport is not called here.
func loadReview(l *reasonLoad, r Reason) error {
	story, err := l.item(r.ID, r)
	if err != nil {
		return err
	}
	path := filepath.Join(l.repo.CacheDir(), "verify", story.ID+".json")
	title := "last flai verify result of " + story.ID
	it, err := l.file(path, story.ID, title, r)
	if errors.Is(err, fs.ErrNotExist) {
		root := l.repo.MainRoot
		if root == "" {
			root = l.repo.Root
		}
		text := fmt.Sprintf("%s has no stored flai verify result: it has not been verified, or its result was cleared. flai verify %s, or the MCP tool verify, runs it in the story's worktree and stores the result, which flai verify %s --last prints.\n", story.ID, story.ID, story.ID)
		it, err = Item{Path: filepath.ToSlash(rel(root, path)), ID: story.ID, Title: title, Text: text}, nil
	}
	if err != nil {
		return err
	}
	l.add(it, r)
	return nil
}

// loadRelease loads a release reason: the release policy's figures as the
// MCP tool release_evaluate gives them, the release package's EvaluateRepo
// in JSON. What is released is read from git tags; when they cannot be
// read, the item says why and that release_evaluate on the host gives the
// figures, so that the rest of the pack still loads.
func loadRelease(l *reasonLoad, r Reason) error {
	title := "the release policy's figures, as release_evaluate and flai release --evaluate give them"
	ev, err := release.EvaluateRepo(execx.System{}, l.repo.Root, l.repo.Manifest, l.repo)
	if err != nil {
		l.add(Item{Path: FigureRelease, Title: title, Text: fmt.Sprintf("flai could not work out the release figures: %v. Call release_evaluate, or run flai release --evaluate on the host, which say why.\n", err)}, r)
		return nil
	}
	data, err := json.MarshalIndent(ev, "", "  ")
	if err != nil {
		return fmt.Errorf("%s --reason %s: writing the release figures: %w", l.cmd, r, err)
	}
	l.add(Item{Path: FigureRelease, Title: title, Text: string(data) + "\n"}, r)
	return nil
}

// loadReady loads a ready reason: the board's stories by column, with the
// WIP limits, the pull order, and each ready story's hold, as the MCP tool
// board and flai board give them, and the backlog stories that could go to
// ready, ranked, with why each other one cannot, as promote_candidates and
// flai promote --candidates give them. Both are text in the form flai's
// commands print, not their JSON, which is tens of kilobytes on a large
// board (ADR-0209); the backlog column is a count, its stories being the
// candidates and the others. The board leaves out the merge queue, which
// the board tool adds and the ready column does not need.
func loadReady(l *reasonLoad, r Reason) error {
	items, err := l.repo.List(true)
	if err != nil {
		return fmt.Errorf("%s --reason %s: listing the items for the board: %w", l.cmd, r, err)
	}
	board, err := l.repo.LoadBoard()
	if err != nil {
		return fmt.Errorf("%s --reason %s: reading the board: %w; fix board.md, which flai check reports", l.cmd, r, err)
	}
	pending := release.PendingIDs(execx.System{}, l.repo.Root, l.repo.Manifest, l.repo)
	view := workitem.NewBoardView(items, board, time.Now(), false, pending, l.repo.Holds(items))
	l.add(Item{Path: FigureBoard, Title: "the board's stories, WIP limits, pull order, and holds, as board and flai board give them", Text: boardText(view)}, r)
	got, err := l.repo.PromotionCandidates(0)
	if err != nil {
		return fmt.Errorf("%s --reason %s: ranking the promotion candidates: %w", l.cmd, r, err)
	}
	l.add(Item{Path: FigureCandidates, Title: "the backlog stories that could go to ready, as promote_candidates and flai promote --candidates give them", Text: candidatesText(got)}, r)
	return nil
}

// boardText is the board as flai board prints it, the backlog as a count:
// each column with its stories in it of its WIP limit, each story with its
// nature, its age in the column, whether it is blocked or held, its touches,
// and why it is held; then whether a story may be pulled, the pull order,
// the limits broken, and what is accepted and not yet published.
func boardText(v workitem.BoardView) string {
	var b strings.Builder
	for _, st := range workitem.States {
		limit := ""
		if l, ok := v.WIPLimits[st]; ok && l > 0 {
			limit = fmt.Sprintf(" (%d/%d)", v.Counts[st], l)
		}
		cards := v.Columns[st]
		if st == workitem.Backlog {
			fmt.Fprintf(&b, "%s%s: stories %d, each under promote_candidates\n", st, limit, len(cards))
			continue
		}
		fmt.Fprintf(&b, "%s%s\n", st, limit)
		if len(cards) == 0 {
			b.WriteString("  -\n")
		}
		for _, cd := range cards {
			flag := ""
			if cd.Blocked {
				flag = " BLOCKED"
			}
			if cd.Held != nil {
				flag += " HELD"
			}
			fmt.Fprintf(&b, "  %s  %s  %s  %s%s\n", cd.ID, cd.Title, cd.Nature, cd.Age, flag)
			if len(cd.Touches) > 0 {
				fmt.Fprintf(&b, "      touches %s\n", strings.Join(cd.Touches, ", "))
			}
			if cd.Held != nil {
				fmt.Fprintf(&b, "      %s\n", cd.Held.Reason)
			}
		}
	}
	if hold := v.PullHold(); hold != "" {
		fmt.Fprintf(&b, "\nno story may be pulled: %s\n", hold)
	}
	if len(v.Order) > 0 {
		fmt.Fprintf(&b, "\npull order: %s\n", strings.Join(v.Order, ", "))
	}
	for _, br := range v.Breaches {
		fmt.Fprintf(&b, "\nover the WIP limit: %s\n", br)
	}
	if len(v.Unpublished) > 0 {
		fmt.Fprintf(&b, "\naccepted, not yet published: %s\n", strings.Join(v.Unpublished, ", "))
	}
	return b.String()
}

// candidatesText is the promotion candidates as flai promote --candidates
// prints them: by the policy, each candidate's rank, ID, figure, and title,
// then each other backlog story with every reason it is not a candidate.
func candidatesText(c workitem.Candidates) string {
	var b strings.Builder
	if len(c.Candidates) == 0 && len(c.Others) == 0 {
		return "no backlog stories\n"
	}
	if len(c.Candidates) == 0 {
		fmt.Fprintf(&b, "no candidates by %s\n", c.Policy)
	} else {
		fmt.Fprintf(&b, "candidates by %s:\n", c.Policy)
	}
	for _, r := range c.Candidates {
		fmt.Fprintf(&b, "  %d  %s  %s  %s\n", r.Position, r.ID, r.Text, r.Title)
	}
	if len(c.Others) > 0 {
		b.WriteString("not candidates:\n")
	}
	for _, o := range c.Others {
		fmt.Fprintf(&b, "  %s  %s\n", o.ID, o.Title)
		for _, reason := range o.Reasons {
			fmt.Fprintf(&b, "    - %s\n", reason)
		}
	}
	return b.String()
}
