package check

import (
	"fmt"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/gitver"
	"github.com/bytepunx/system-flow/flai/internal/storygit"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// RunScoped runs the check on repo as flai check does: with the findings of
// the installed git's version, and, given a story, each finding outside what
// the story owns or changes marked a note (ADR-0085).
func RunScoped(repo *workitem.Repo, r execx.Runner, now time.Time, story string) (*Result, error) {
	res, err := Run(repo, now)
	if err != nil {
		return nil, err
	}
	// The clone may need a newer git than the one installed (ADR-0022).
	// No git on PATH means nothing here depends on its version.
	if v, err := gitver.Installed(r); err == nil {
		GitCompat(res, repo, v)
	}
	if story == "" {
		return res, nil
	}
	changed, err := storygit.StoryChanges(r, repo, story)
	if err != nil {
		return nil, fmt.Errorf("scope the check to %s: %w", story, err)
	}
	if err := ScopeToStory(res, repo, story, changed); err != nil {
		return nil, err
	}
	return res, nil
}
