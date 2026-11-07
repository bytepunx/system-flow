package check

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// S-0577: an invalid kluster: block is the finding manifest.kluster, an
// error on the block's line naming the key and the reason; a valid block,
// and a manifest with none, give no finding.
func TestKlusterBlockIsChecked(t *testing.T) {
	const valid = "kluster:\n  config: integration/kluster.yaml\n  deploy: [scripts/kluster-deploy.sh]\n"
	withConfig := func(t *testing.T, extra, config string) []Finding {
		t.Helper()
		repo := planningProject(t, extra)
		if config != "" {
			p := filepath.Join(repo.Root, "integration", "kluster.yaml")
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte(config), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return findings(t, repo, "manifest.kluster")
	}
	if got := withConfig(t, "", ""); len(got) != 0 {
		t.Errorf("no block: %+v", got)
	}
	if got := withConfig(t, valid, "name: t-integration\nprofile: t\n"); len(got) != 0 {
		t.Errorf("valid block: %+v", got)
	}
	for _, c := range []struct {
		name, extra, config, want string
	}{
		{"config does not exist", valid, "", `kluster.config "integration/kluster.yaml" does not exist`},
		{"unknown key", valid + "  namespce: t\n", "name: t\n", `kluster has no key "namespce"`},
		{"deploy missing", "kluster:\n  config: integration/kluster.yaml\n", "name: t\n", "kluster.deploy is required"},
		{"provider not k3d", valid, "name: t\nprovider: kind\n", `names the provider "kind"`},
		{"namespace story form", valid + "  namespace: s-0042\n", "name: t\n", `kluster.namespace "s-0042" has a story namespace's form`},
		{"timeout", valid + "  timeout: forever\n", "name: t\n", `kluster.timeout "forever" is not a duration`},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := withConfig(t, c.extra, c.config)
			if len(got) != 1 {
				t.Fatalf("want one finding: %+v", got)
			}
			// the block opens on line 8, after the planning project's seven
			if f := got[0]; f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, c.want) {
				t.Errorf("finding: %+v, want an error on line 8 with %q", f, c.want)
			}
		})
	}
}
