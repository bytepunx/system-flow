package verify

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
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
		})
	}
	return out
}

// CheckoutTiers are the test tiers in the manifest of the checkout at root,
// a story's worktree as its branch has them or the main checkout, or the
// default tier when the manifest declares none. It refuses a manifest it
// cannot read and tiers that are not valid, naming what to fix.
func CheckoutTiers(root string) ([]Tier, error) {
	m, err := manifest.Load(filepath.Join(root, "system-flow.yaml"))
	if err != nil {
		return nil, fmt.Errorf("read the test tiers: %w", err)
	}
	tiers, err := m.TestTiers(os.DirFS(root))
	if err != nil {
		return nil, fmt.Errorf("the manifest's tests are not valid, so no tier runs: %w; fix them with flai manifest set tests='<JSON list>'", err)
	}
	return FromManifest(tiers), nil
}
