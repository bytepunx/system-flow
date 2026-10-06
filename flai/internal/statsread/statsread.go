// Package statsread reads from a project what flai stats computes its
// metrics from, for the command and for flai serve's read of it alike, so
// that both pass metrics.Compute the same (S-0205).
package statsread

import (
	"fmt"
	"log/slog"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/metrics"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Read returns every item, archived ones too, and the options read from the
// project: the strategic agents' activity documents, the threads, the
// issues, the manifest's projects, the board's in-progress limit, and the files each
// story's commits changed. The caller sets the window and grouping. When git
// cannot be read, Commits stays nil, which leaves claims.drift out, and Read
// logs a warning saying so on log; any other read that fails stops it.
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
		return nil, opt, fmt.Errorf("cannot read the board's in-progress limit: %w; restore board.md from git or run flai check to see what is wrong", err)
	}
	opt.WIPLimit = board.WIPLimits[workitem.InProgress]
	opt.Projects = repo.Manifest.Projects
	if opt.Commits, err = storygit.Committed(r, repo); err != nil {
		log.Warn("cannot read the stories' commits for touches drift", "component", "stats", "err", err.Error(),
			"detail", "claims.drift is left out; run flai stats in the project's git checkout, with git installed, to compare touches with commits")
	}
	return items, opt, nil
}
