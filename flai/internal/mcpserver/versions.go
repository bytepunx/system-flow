package mcpserver

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

const versionsDescription = "The published releases of flai and of the dashboard (S-0298, ADR-0117), each list newest first; it installs nothing. flai lists the GitHub releases tagged flai/vX.Y.Z that are neither drafts nor prereleases, each with its version, tag, and published date, marking running (the flai this server runs), latest (the newest), and below_minimum (older than the project's flai.minimum, which it cannot open). dashboard lists the repository's flaiover/vX.Y.Z tags, whose image tag is the bare X.Y.Z, marking latest. running and minimum say the version this server runs and the project's flai.minimum, empty when it sets none. The releases come from FLAI_RELEASES_API when set, else api.github.com, as flai self-upgrade reads them. No MCP tool installs or deploys a release: that is the operator's, with flai host upgrade --version X.Y.Z or flai dashboard upgrade --published --tag X.Y.Z on the host, or the dashboard's Updates page."

// VersionsIn asks for the published releases of flai and the dashboard.
type VersionsIn struct {
	Project string `json:"project,omitempty" jsonschema:"the project, by key or folder, whose flai.minimum marks the flai releases below it: needed only when the server serves more than one"`
}

func (in VersionsIn) project() string { return in.Project }

// FlaiVersion is one published flai release, as flai self-upgrade --list
// --json gives it, marked against this server rather than the installed flai.
type FlaiVersion struct {
	Version      string    `json:"version"`
	Tag          string    `json:"tag"`
	Published    time.Time `json:"published,omitzero"`
	Running      bool      `json:"running" jsonschema:"the flai this server runs"`
	Latest       bool      `json:"latest" jsonschema:"the newest published flai"`
	BelowMinimum bool      `json:"below_minimum" jsonschema:"older than the project's flai.minimum: it cannot open the project"`
}

// DashboardVersion is one published dashboard release, as flai dashboard
// versions --json gives it, without what needs Docker or the host's config.
type DashboardVersion struct {
	Version string `json:"version" jsonschema:"the image tag, X.Y.Z"`
	Tag     string `json:"tag"`
	Latest  bool   `json:"latest" jsonschema:"the newest published dashboard"`
}

// VersionsOut is both lists of published releases, newest first.
type VersionsOut struct {
	Flai      []FlaiVersion      `json:"flai"`
	Dashboard []DashboardVersion `json:"dashboard"`
	Running   string             `json:"running" jsonschema:"the version of the flai this server runs"`
	Minimum   string             `json:"minimum" jsonschema:"the project's flai.minimum; empty when it sets none"`
}

// versions answers the published flai and dashboard releases, marking the
// flai this server runs, the newest of each, and any flai below the
// project's minimum (ADR-0117 §6).
func (s *server) versions(ctx context.Context, _ *mcp.CallToolRequest, _ VersionsIn) (*mcp.CallToolResult, VersionsOut, error) {
	opt := selfupgrade.Options{APIBase: strings.TrimSpace(os.Getenv("FLAI_RELEASES_API")), Token: selfupgrade.Token(s.runner)}
	flai, err := selfupgrade.List(ctx, opt, selfupgrade.TagPrefix)
	if err != nil {
		return nil, VersionsOut{}, err
	}
	dashboard, err := selfupgrade.ListTags(ctx, opt, selfupgrade.DashboardTagPrefix)
	if err != nil {
		return nil, VersionsOut{}, err
	}
	minimum := s.repo.Manifest.Flai.Minimum
	out := VersionsOut{Flai: make([]FlaiVersion, len(flai)), Dashboard: make([]DashboardVersion, len(dashboard)), Running: s.version, Minimum: minimum}
	for i, r := range flai {
		out.Flai[i] = FlaiVersion{Version: r.Version, Tag: r.Tag, Published: r.Published, Running: r.Version == s.version, Latest: i == 0, BelowMinimum: buildinfo.Below(r.Version, minimum)}
	}
	for i, r := range dashboard {
		out.Dashboard[i] = DashboardVersion{Version: r.Version, Tag: r.Tag, Latest: i == 0}
	}
	return nil, out, nil
}
