package check

import (
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A project with no stories and no sub-projects still accepts on a heading
// every word a pack primed for the orchestrator's reasons reads beside all:
// the kinds of reason and the built-in strategic roles' topics (ADR-0217),
// so the template's baseline checks clean in a project made from it.
func TestVocabularyHoldsTheReasonPacksWords(t *testing.T) {
	v := Vocabulary(&workitem.Repo{Root: t.TempDir(), Manifest: manifest.Manifest{}}, nil)
	for _, w := range []string{"all", "code", "thread", "review", "planned", "draft", "release", "ready", "planning", "orchestration", "analysis", "digest"} {
		if !v[w] {
			t.Errorf("Vocabulary lacks %q; want every reason kind and every built-in role's topic", w)
		}
	}
	if v["accepted"] {
		t.Errorf("Vocabulary holds %q, which nothing declares", "accepted")
	}
}
