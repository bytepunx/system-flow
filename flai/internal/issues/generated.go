package issues

import (
	"time"

	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Generated is every committed file flai writes from others, so that a
// rebase of story's branch stopped on them alone is resolved by writing them
// again in its worktree (ADR-0098): the summary alone today, written from
// the issue files at now's time.
func Generated(r *workitem.Repo, story string, now func() time.Time) []storygit.GeneratedFile {
	wt := RepoFor(r, story)
	return []storygit.GeneratedFile{{
		Path: SummaryPath(r),
		Regenerate: func() error {
			_, err := WriteSummary(wt, now())
			return err
		},
	}}
}
