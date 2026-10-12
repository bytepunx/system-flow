// Package statsread reads from a project what flai stats computes its
// metrics from, for the command and for flai serve's read of it alike, so
// that both pass metrics.Compute the same (S-0205).
package statsread

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/issues"
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
// changed, the conversations between stories, and the conflict log. The
// caller sets the window and grouping. When git cannot be read, Commits stays
// nil, which leaves claims.drift out; a conversation that cannot be read is
// left out; and a conflict log that cannot be read sets ConflictsUnreadable,
// which leaves coordination's conflicts out. Read logs a warning on log for
// each; any other read that fails stops it.
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
	return items, opt, nil
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
