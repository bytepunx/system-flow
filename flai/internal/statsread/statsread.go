// Package statsread reads from a project what flai stats computes its
// metrics from, for the command and for flai serve's read of it alike, so
// that both pass metrics.Compute the same (S-0205).
package statsread

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/mergequeue"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Read returns every item, archived ones too, and the options read from the
// project: the strategic agents' activity documents and the roles the
// project's role files define, the threads, the
// issues, the manifest's projects, the board's in-progress limit and pull
// order, planning.default_duration, the files each story's commits
// changed, the conversations between stories, the conflict log, the merge
// queue's event log, the manifest's owner and the digest's entries, and,
// under pull-request acceptance only, each story's kept forge events. The
// caller sets the window and grouping. When git cannot be read, Commits
// stays nil, which leaves claims.drift out; a conversation that cannot be
// read is left out; a conflict log that cannot be read sets
// ConflictsUnreadable, which leaves coordination's conflicts out; a merge
// queue event log that cannot be read leaves QueueLogged false, as a missing
// one does, which leaves the queue section out; a digest that cannot be read
// sets DigestUnreadable, which leaves the digest section out, while a
// missing digest folder is no entries; and a story's kept forge state that
// cannot be read is left out. Read logs a warning on log for each; any other
// read that fails stops it.
//
// Root, ItemsGeneration, and CommitsKey are set so that Compute keeps the hold
// replay and the touches drift between calls while the items and the commits
// stay the same (S-0608).
//
// The reads run at once, each in its own goroutine, so Read takes about as
// long as its slowest read (S-0607). When several fail, the error is the one
// of the first in the order above, and each warning is logged with its
// detail as before; the warnings may come in another order, and those of the
// reads that warn are logged even when another read fails.
func Read(r execx.Runner, repo *workitem.Repo, log *slog.Logger) ([]*workitem.Item, metrics.Options, error) {
	var (
		wg                                            sync.WaitGroup
		items                                         []*workitem.Item
		activities                                    []*workitem.Activity
		roles                                         []string
		threadList                                    []*threads.Thread
		issueList                                     []*issues.Issue
		board                                         *workitem.Board
		commits                                       map[string][]string
		commitsKey                                    string
		generation                                    uint64
		convs                                         []*messages.Conversation
		conflicts                                     []storygit.Conflict
		queue                                         []mergequeue.Event
		queueLogged, digestUnreadable                 bool
		entries                                       []digest.Entry
		forged                                        map[string][]forge.StatusEvent
		errItems, errActivities, errRoles, errThreads error
		errIssues, errBoard, errConflicts             error
	)
	cacheDir := repo.CacheDir()
	wg.Go(func() {
		// The generation read before and after the list says the items listed
		// are those of that generation; when it moved meanwhile, none is given,
		// so nothing is kept from them (S-0608).
		before := repo.Generation()
		items, errItems = repo.List(true)
		if after := repo.Generation(); after == before {
			generation = after
		}
	})
	wg.Go(func() { activities, errActivities = repo.Activities() })
	wg.Go(func() { roles, errRoles = repo.Roles() })
	wg.Go(func() { threadList, errThreads = threads.List(repo) })
	wg.Go(func() { issueList, errIssues = issues.List(repo) })
	wg.Go(func() { board, errBoard = repo.LoadBoard() })
	wg.Go(func() {
		var err error
		if commits, commitsKey, err = storygit.CommittedKeyed(r, repo); err != nil {
			log.Warn("cannot read the stories' commits for touches drift", "component", "stats", "err", err.Error(),
				"detail", "claims.drift is left out; run flai stats in the project's git checkout, with git installed, to compare touches with commits")
		}
	})
	wg.Go(func() { convs = conversations(repo, log) })
	wg.Go(func() {
		if conflicts, errConflicts = storygit.ReadConflicts(cacheDir); errConflicts != nil {
			log.Warn("cannot read the conflict log", "component", "stats", "err", errConflicts.Error(),
				"detail", "coordination's conflicts are left out and coordination.logged is false; correct or remove the line the error names in "+storygit.ConflictsPath(cacheDir))
		}
	})
	wg.Go(func() { queue, queueLogged = queueEvents(cacheDir, log) })
	wg.Go(func() {
		var err error
		if entries, err = digest.List(repo); err != nil {
			log.Warn("cannot read the digest", "component", "stats", "err", err.Error(),
				"detail", "the digest section is left out; run flai check to see what is wrong with the entry the error names in "+digest.Dir(repo.Manifest))
			entries, digestUnreadable = nil, true
		}
	})
	if repo.Manifest.Acceptance.PullRequest() {
		wg.Go(func() { forged = forgeEvents(cacheDir, log) })
	}
	wg.Wait()

	var opt metrics.Options
	if errItems != nil {
		return nil, opt, errItems
	}
	if errActivities != nil {
		return nil, opt, fmt.Errorf("cannot read the strategic agents' activity documents: %w; flai writes them, so restore the file from git or run flai check to see what is wrong", errActivities)
	}
	if errRoles != nil {
		return nil, opt, errRoles
	}
	if errThreads != nil {
		return nil, opt, fmt.Errorf("cannot read the threads: %w; restore the thread from git or run flai check to see what is wrong", errThreads)
	}
	if errIssues != nil {
		return nil, opt, fmt.Errorf("cannot read the issues: %w; restore the issue from git or run flai check to see what is wrong", errIssues)
	}
	if errBoard != nil {
		return nil, opt, fmt.Errorf("cannot read the board's in-progress limit and pull order: %w; restore board.md from git or run flai check to see what is wrong", errBoard)
	}
	var err error
	if opt.Fallback, err = repo.Manifest.Planning.FallbackDuration(); err != nil {
		return nil, opt, fmt.Errorf("cannot project the ready column's cost of delay: %w", err)
	}
	opt.Activities, opt.Roles, opt.Threads, opt.Issues = activities, roles, threadList, issueList
	opt.WIPLimit, opt.Order = board.WIPLimits[workitem.InProgress], board.Order
	opt.Projects = repo.Manifest.Projects
	opt.Shared = repo.SharedClaims()
	opt.Commits = commits
	opt.Root, opt.ItemsGeneration, opt.CommitsKey = repo.Root, generation, commitsKey
	opt.Conversations = convs
	opt.Conflicts, opt.ConflictsUnreadable = conflicts, errConflicts != nil
	opt.QueueEvents, opt.QueueLogged = queue, queueLogged
	opt.Owner = repo.Manifest.Owner
	opt.Entries, opt.DigestUnreadable = entries, digestUnreadable
	opt.Forge = forged
	return items, opt, nil
}

// queueEvents reads the merge queue's event log under cacheDir, and whether
// it is there and could be read: a missing log is no warning, and an empty
// one is a log with no events. A log that cannot be read is warned of on log,
// so that flai stats never fails on it (S-0404).
func queueEvents(cacheDir string, log *slog.Logger) ([]mergequeue.Event, bool) {
	path := mergequeue.EventsPath(cacheDir)
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		return nil, false
	}
	events, err := mergequeue.ReadEvents(cacheDir)
	if err != nil {
		log.Warn("cannot read the merge queue's event log", "component", "stats", "err", err.Error(),
			"detail", "the queue section is left out; correct or remove the line the error names in "+path)
		return nil, false
	}
	return events, true
}

// CloseOuts reads flai serve's host journal at path, the journal.jsonl of the
// serve folder beside the configuration flai reads, for the close-outs of
// repo's project, whose main checkout it names (S-0546, ADR-0229). It is
// apart from Read because only flai stats on the host reads it. A missing
// journal is nil with no warning; one that cannot be opened or read is nil
// and warned of on log, so that flai stats never fails on it. Either leaves
// the close-outs section out. A line that is not an entry, once its NUL bytes
// are removed, is skipped, and one warning names how many were and the first
// one's number, keeping the section (ADR-0268).
func CloseOuts(path string, repo *workitem.Repo, log *slog.Logger) *metrics.CloseOutJournal {
	root := repo.Root
	if repo.MainRoot != "" {
		root = repo.MainRoot
	}
	entries, skipped, err := journalEntries(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		log.Warn("cannot read flai serve's journal", "component", "stats", "err", err.Error(),
			"detail", "the close-outs section is left out; correct or remove the line the error names in "+path)
		return nil
	}
	if skipped.count > 0 {
		log.Warn("skipped lines of flai serve's journal that are not entries", "component", "stats",
			"path", path, "skipped", skipped.count, "first_line", skipped.first, "err", skipped.err.Error(),
			"detail", "the close-outs those lines held are not counted; correct or remove them in "+path)
	}
	return &metrics.CloseOutJournal{Root: root, Entries: entries}
}

// skippedLines counts the journal lines that are not entries, with the first
// one's number and decode error.
type skippedLines struct {
	count, first int
	err          error
}

// journalEntries decodes every line of the host journal at path, each with
// its NUL bytes removed, as a write a host's stop cut short leaves them
// (I-0290); an empty line is skipped, and a line that is still not an entry
// is skipped and counted in skipped. Only a journal that cannot be opened or
// read is an error.
func journalEntries(path string) ([]metrics.JournalEntry, skippedLines, error) {
	var skipped skippedLines
	f, err := os.Open(path)
	if err != nil {
		return nil, skipped, err
	}
	defer func() { _ = f.Close() }()
	entries := []metrics.JournalEntry{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for n := 1; sc.Scan(); n++ {
		line := bytes.TrimSpace(bytes.ReplaceAll(sc.Bytes(), []byte{0}, nil))
		if len(line) == 0 {
			continue
		}
		var e metrics.JournalEntry
		if err := json.Unmarshal(line, &e); err != nil {
			if skipped.count == 0 {
				skipped.first, skipped.err = n, err
			}
			skipped.count++
			continue
		}
		entries = append(entries, e)
	}
	if err := sc.Err(); err != nil {
		return nil, skipped, fmt.Errorf("read flai serve's journal %s: %w", path, err)
	}
	return entries, skipped, nil
}

// forgeEvents reads the events of every story's kept forge state under
// cacheDir, by the story's canonical ID, an empty map when there is none. A
// state that cannot be read is left out with a warning on log, so that flai
// stats never fails on one (ADR-0168).
func forgeEvents(cacheDir string, log *slog.Logger) map[string][]forge.StatusEvent {
	out := map[string][]forge.StatusEvent{}
	files, err := filepath.Glob(forge.StatusPath(cacheDir, "*"))
	if err != nil {
		log.Warn("cannot list the kept pull request statuses", "component", "stats", "err", err.Error(),
			"detail", "the checks and reviewers waits are left out; check that the project's .flai-cache is readable")
		return out
	}
	sort.Strings(files)
	var skipped []string
	var first error
	for _, f := range files {
		story := strings.TrimSuffix(filepath.Base(f), ".json")
		st, ok, err := forge.ReadStatus(cacheDir, story)
		if err != nil {
			skipped = append(skipped, f)
			if first == nil {
				first = err
			}
			continue
		}
		if ok {
			out[workitem.CanonicalID(story)] = st.Events
		}
	}
	if len(skipped) > 0 {
		log.Warn("cannot read some kept pull request statuses", "component", "stats", "count", len(skipped), "paths", skipped, "err", first.Error(),
			"detail", "their stories' checks and reviewers waits are left out; delete each, and flai serve reads its pull request afresh while the story is in review")
	}
	return out
}

// conversations reads every conversation under wip/messages in the main
// checkout, in ID order, whatever its status. One that cannot be read is
// left out with a warning on log, so that flai stats never fails on one.
func conversations(repo *workitem.Repo, log *slog.Logger) []*messages.Conversation {
	files, err := filepath.Glob(filepath.Join(messages.Dir(repo), "MS-*.md"))
	if err != nil {
		log.Warn("cannot list the conversations", "component", "stats", "err", err.Error(),
			"detail", "coordination's conversations and shares are left out")
		return nil
	}
	sort.Strings(files)
	var out []*messages.Conversation
	var skipped []string
	var first error
	for _, f := range files {
		c, err := messages.Read(f)
		if err != nil {
			skipped = append(skipped, f)
			if first == nil {
				first = err
			}
			continue
		}
		out = append(out, c)
	}
	if len(skipped) > 0 {
		log.Warn("cannot read some conversations", "component", "stats", "count", len(skipped), "paths", skipped, "err", first.Error(),
			"detail", "coordination leaves them out; run flai check to see what is wrong with each")
	}
	return out
}
