package hostapi

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/perf"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func designMethods() map[string]channel.Method {
	return map[string]channel.Method{
		// design.get: as flai design show --json; the project's design system,
		// its tokens per theme, their problems, renderings, and contrast. No
		// token file is no error but an answer that says so. With a story,
		// as flai design show --story: the design system in that open
		// story's worktree (S-0384).
		"design.get": func(ctx context.Context, p channel.Project, raw json.RawMessage) (any, *channel.Error) {
			in, e := decode[struct {
				Story string `json:"story"`
			}](raw)
			if e != nil {
				return nil, e
			}
			if in.Story != "" {
				return storyDesign(ctx, p, in.Story)
			}
			done := perf.Track(ctx, "manifest.load")
			m, err := manifest.Load(filepath.Join(p.Root, manifest.File))
			done()
			if err != nil {
				return nil, failed(err)
			}
			defer perf.Track(ctx, "design.read")()
			shown, err := designsystem.Show(os.DirFS(p.Root), m.DesignTokens(), m.DesignGuidance())
			if err != nil {
				return nil, failed(err)
			}
			shown.Stylesheet = m.DesignStylesheet()
			return shown, nil
		},
		// design.proposals: as flai design proposals --json; the open
		// stories whose branches, or whose worktrees uncommitted, change the
		// design system, each with its ID, title, state, and the paths it
		// changes; none is [] (S-0384).
		"design.proposals": func(ctx context.Context, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
			done := perf.Track(ctx, "repo.open")
			repo, err := workitem.Open(p.Root)
			done()
			if err != nil {
				return nil, failed(err)
			}
			r := execx.Timed(ctx, execx.System{})
			repo.Git = r
			defer perf.Track(ctx, "design.proposals")()
			out, err := designsystem.Proposals(r, repo)
			if err != nil {
				return nil, failed(err)
			}
			return out, nil
		},
	}
}

// storyDesign is the design system in the worktree of the open story id
// names. An ID no item has is NotFound; an item that is not a story, a
// story done or cancelled, and a story without a worktree are refused; a
// worktree that cannot be read fails.
func storyDesign(ctx context.Context, p channel.Project, id string) (any, *channel.Error) {
	if e := needID(id); e != nil {
		return nil, e
	}
	done := perf.Track(ctx, "repo.open")
	repo, err := workitem.Open(p.Root)
	done()
	if err != nil {
		return nil, failed(err)
	}
	defer perf.Track(ctx, "design.read")()
	shown, _, err := designsystem.ShowStory(repo, id)
	var refused *designsystem.Refused
	switch {
	case errors.Is(err, workitem.ErrNotFound):
		return nil, &channel.Error{Code: NotFound, Message: inRepo(p.Root, err)}
	case errors.As(err, &refused):
		return nil, bad("%s", inRepo(p.Root, err))
	case err != nil:
		return nil, failed(err)
	}
	return shown, nil
}

// inRepo is err's message with the project's root left out of the paths it
// names, which an answer gives relative to the repository.
func inRepo(root string, err error) string {
	return strings.ReplaceAll(err.Error(), filepath.Clean(root)+string(filepath.Separator), "")
}
