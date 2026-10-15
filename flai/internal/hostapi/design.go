package hostapi

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/perf"
)

func designMethods() map[string]channel.Method {
	return map[string]channel.Method{
		// design.get: as flai design show --json; the project's design system,
		// its tokens per theme, their problems, renderings, and contrast. No
		// token file is no error but an answer that says so.
		"design.get": func(ctx context.Context, p channel.Project, _ json.RawMessage) (any, *channel.Error) {
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
			return shown, nil
		},
	}
}
