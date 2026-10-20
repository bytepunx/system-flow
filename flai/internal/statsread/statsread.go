// Package statsread reads from a project what flai stats computes its
// metrics from, for the command and for flai serve's read of it alike, so
// that both pass metrics.Compute the same (S-0205).
package statsread

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

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
// project: the strategic agents' activity documents, the threads, the
// issues, the manifest's projects, the board's in-progress limit and pull
// order, planning.default_duration, the files each story's commits
// changed, the conversations between stories, the conflict log, the merge
// queue's event log, and, under pull-request acceptance only, each story's
// kept forge events. The caller sets the window and grouping. When git
// cannot be read, Commits stays nil, which leaves claims.drift out; a
// conversation that cannot be read is left out; a conflict log that cannot
// be read sets ConflictsUnreadable, which leaves coordination's conflicts
// out; a merge queue event log that cannot be read leaves QueueLogged false,
// as a missing one does, which leaves the queue section out; and a story's
// kept forge state that cannot be read is left out. Read logs a warning on
// log for each; any other read that fails stops it.
func Read(r execx.Runner, repo *workitem.Repo, log *slog.Logger) ([]*workitem.Item, metrics.Options, error) {
	var opt metrics.Options
	items, err := repo.List(true)
	if err != nil {
		return nil, opt, err
	}
	if opt.Activities, err = repo.Activities(); err != nil {
		return nil, opt, fmt.Errorf("cannot read the strategic agents' activity documents: %w; flai writes them, so restore the file from git or run flai check to see what is wrong", err)
	}
	if opt.Threads, err = threads.List(repo); err != nil {
		return nil, opt, fmt.Errorf("cannot read the threads: %w; restore the thread from git or run flai check to see what is wrong", err)
	}
	if opt.Issues, err = issues.List(repo); err != nil {
		return nil, opt, fmt.Errorf("cannot read the issues: %w; restore the issue from git or run flai check to see what is wrong", err)
	}
	board, err := repo.LoadBoard()
	if err != nil {
		return nil, opt, fmt.Errorf("cannot read the board's in-progress limit and pull order: %w; restore board.md from git or run flai check to see what is wrong", err)
	}
	opt.WIPLimit, opt.Order = board.WIPLimits[workitem.InProgress], board.Order
	if opt.Fallback, err = repo.Manifest.Planning.FallbackDuration(); err != nil {
		return nil, opt, fmt.Errorf("cannot project the ready column's cost of delay: %w", err)
	}
	opt.Projects = repo.Manifest.Projects
	opt.Shared = repo.SharedClaims()
	if opt.Commits, err = storygit.Committed(r, repo); err != nil {
		log.Warn("cannot read the stories' commits for touches drift", "component", "stats", "err", err.Error(),
			"detail", "claims.drift is left out; run flai stats in the project's git checkout, with git installed, to compare touches with commits")
	}
	opt.Conversations = conversations(repo, log)
	conflicts, err := storygit.ReadConflicts(repo.CacheDir())
	if err != nil {
		log.Warn("cannot read the conflict log", "component", "stats", "err", err.Error(),
			"detail", "coordination's conflicts are left out and coordination.logged is false; correct or remove the line the error names in "+storygit.ConflictsPath(repo.CacheDir()))
	}
	opt.Conflicts, opt.ConflictsUnreadable = conflicts, err != nil
	opt.QueueEvents, opt.QueueLogged = queueEvents(repo.CacheDir(), log)
	if repo.Manifest.Acceptance.PullRequest() {
		opt.Forge = forgeEvents(repo.CacheDir(), log)
	}
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
