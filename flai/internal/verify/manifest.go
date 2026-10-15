package verify

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/release"
)

// FromManifest is the manifest's test tiers as a run takes them, in order.
func FromManifest(tiers []manifest.TestTier) []Tier {
	out := make([]Tier, 0, len(tiers))
	for _, t := range tiers {
		out = append(out, Tier{
			Name:       t.Name,
			Command:    slices.Clone(t.Command),
			Dir:        t.Dir,
			Paths:      slices.Clone(t.Paths),
			Format:     t.FormatOrDefault(),
			AllOnly:    t.AllOnly,
			AllCommand: slices.Clone(t.AllCommand),
			FixCommand: slices.Clone(t.FixCommand),
			Covers:     slices.Clone(t.Covers),
		})
	}
	return out
}

// CheckoutTiers are the test tiers in the manifest of the checkout at root,
// a story's worktree as its branch has them or the main checkout, or the
// default tier when the manifest declares none. It refuses a manifest it
// cannot read and tiers that are not valid, naming what to fix.
func CheckoutTiers(root string) ([]Tier, error) { return CheckoutTiersFor(root, nil, "") }

// CheckoutTiersFor is CheckoutTiers for the flai at version running, asking
// git through r. When the tiers are not valid and running is older than the
// newest flai release reachable from the checkout's HEAD, the refusal says
// so and names the upgrade instead of the fix, since the manifest's tests
// may use what this flai does not know (ADR-0158, I-0141). With no runner,
// a running version that is not a release, or no newer release, it refuses
// as CheckoutTiers does.
func CheckoutTiersFor(root string, r execx.Runner, running string) ([]Tier, error) {
	m, err := manifest.Load(filepath.Join(root, "system-flow.yaml"))
	if err != nil {
		return nil, fmt.Errorf("read the test tiers: %w", err)
	}
	tiers, err := m.TestTiers(os.DirFS(root))
	if err != nil {
		if out := release.FlaiOutdated(r, root, running); out != nil {
			return nil, fmt.Errorf("the manifest's tests are not valid for this flai, so no tier runs: %w; this flai is %s, older than flai %s in the project's history, so the manifest's tests may use what this flai does not know: upgrade it with %s",
				err, out.Running, out.Newest, out.Upgrade)
		}
		return nil, fmt.Errorf("the manifest's tests are not valid, so no tier runs: %w; fix them with flai manifest set tests='<JSON list>'", err)
	}
	return FromManifest(tiers), nil
}
