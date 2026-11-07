package kluster

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// initRoot is a repository root holding manifest as its system-flow.yaml,
// and the manifest it loads as.
func initRoot(t *testing.T, text string) (string, manifest.Manifest) {
	t.Helper()
	root := t.TempDir()
	path := filepath.Join(root, manifest.File)
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		t.Fatal(err)
	}
	m, err := manifest.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	return root, m
}

const initManifest = "version: 1\nname: Auth Star\nkey: authstar\n# the layout\nlayout:\n  design: design   # kept\n  docs: docs\n  wip: wip\n"

// S-0586: Init appends a block naming only config: kluster.yaml, every
// other line kept, writes kluster.yaml naming the cluster, k3d, and spire,
// and the result has no manifest.kluster finding.
func TestInitWritesTheBlockAndTheConfig(t *testing.T) {
	root, m := initRoot(t, initManifest)
	s, err := Init(root, m)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(s.Wrote, []string{"system-flow.yaml", "kluster.yaml"}) || s.Config != "kluster.yaml" || s.Kept {
		t.Errorf("scaffold = %+v", s)
	}
	got, _ := os.ReadFile(filepath.Join(root, manifest.File))
	if want := initManifest + "kluster:\n  config: kluster.yaml\n"; string(got) != want {
		t.Errorf("manifest:\n%s\nwant:\n%s", got, want)
	}
	cfg, err := manifest.ReadKlusterConfig(os.DirFS(root), "kluster.yaml")
	if err != nil || cfg.Name != "auth-star-local" || cfg.Provider != "k3d" || cfg.Profile != ScaffoldProfile {
		t.Errorf("config = %+v, %v", cfg, err)
	}
	after, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil {
		t.Fatal(err)
	}
	if after.Kluster == nil || after.Kluster.Config != "kluster.yaml" || len(after.Kluster.Deploy) != 0 {
		t.Errorf("block = %+v", after.Kluster)
	}
	if errs := after.KlusterErrors(os.DirFS(root)); len(errs) != 0 {
		t.Errorf("the scaffold is faulty: %q", errs)
	}
}

// S-0586: Init keeps a kluster.yaml that is there, and writes nothing for a
// manifest with a block already.
func TestInitKeepsAConfigAndABlock(t *testing.T) {
	root, m := initRoot(t, initManifest)
	mine := "name: mine\nprofile: authstar\n"
	if err := os.WriteFile(filepath.Join(root, "kluster.yaml"), []byte(mine), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := Init(root, m)
	if err != nil || !s.Kept || !slices.Equal(s.Wrote, []string{"system-flow.yaml"}) {
		t.Fatalf("scaffold = %+v, %v", s, err)
	}
	if got, _ := os.ReadFile(filepath.Join(root, "kluster.yaml")); string(got) != mine {
		t.Errorf("kluster.yaml changed:\n%s", got)
	}

	block := initManifest + "kluster:\n  config: integration/kluster.yaml\n"
	root, m = initRoot(t, block)
	s, err = Init(root, m)
	if err != nil || len(s.Wrote) != 0 || s.Config != "integration/kluster.yaml" {
		t.Errorf("with a block: %+v, %v", s, err)
	}
	if got, _ := os.ReadFile(filepath.Join(root, manifest.File)); string(got) != block {
		t.Errorf("with a block, the manifest changed:\n%s", got)
	}
	if _, err := os.Stat(filepath.Join(root, "kluster.yaml")); err == nil {
		t.Error("with a block, kluster.yaml was written")
	}
}

// S-0586: a block the manifest's own validation would fault, as for a name
// that makes no DNS label, is refused naming the field, and nothing is
// written.
func TestInitRefusesABlockTheCheckWouldFault(t *testing.T) {
	text := "version: 1\nname: \"!!!\"\nkey: x\nlayout:\n  design: design\n  docs: docs\n  wip: wip\n"
	root, m := initRoot(t, text)
	_, err := Init(root, m)
	if err == nil || !strings.Contains(err.Error(), "kluster.namespace") {
		t.Fatalf("Init = %v, want a refusal naming kluster.namespace", err)
	}
	if got, _ := os.ReadFile(filepath.Join(root, manifest.File)); string(got) != text {
		t.Errorf("the refused Init changed the manifest:\n%s", got)
	}
	if _, err := os.Stat(filepath.Join(root, "kluster.yaml")); err == nil {
		t.Error("the refused Init wrote kluster.yaml")
	}
}

// S-0586: the scaffold's cluster name is the project's name as a DNS label
// with -local after it, within 63 characters.
func TestScaffoldName(t *testing.T) {
	for in, want := range map[string]string{"Auth Star": "auth-star-local", "!!!": "local", strings.Repeat("a", 70): strings.Repeat("a", 57) + "-local"} {
		if got := ScaffoldName(in); got != want {
			t.Errorf("ScaffoldName(%q) = %q, want %q", in, got, want)
		}
	}
}
