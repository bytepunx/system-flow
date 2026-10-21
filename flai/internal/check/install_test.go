package check

import (
	"strings"
	"testing"
)

// S-0483: a sub-project's install key that is not valid is an error on the
// projects key, naming the sub-project's index and field; a valid key, and
// none, are not reported.
func TestInstallKeysAreChecked(t *testing.T) {
	const web = "projects:\n  - name: web\n    path: web\n    kind: sveltekit\n"
	for _, extra := range []string{
		web,
		web + "    install:\n      command: [pnpm, install]\n      paths: [package.json, pnpm-lock.yaml]\n",
	} {
		if got := findings(t, planningProject(t, extra), "manifest.install"); len(got) != 0 {
			t.Errorf("valid install %q: %+v", extra, got)
		}
	}
	got := findings(t, planningProject(t, web+"    install:\n      command: [pnpm, install]\n"), "manifest.install")
	if len(got) != 1 {
		t.Fatalf("want one error, the paths: %+v", got)
	}
	if f := got[0]; f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, `projects[0].install.paths (sub-project "web") is empty`) {
		t.Errorf("finding: %+v, want an error on line 8 naming projects[0].install.paths", f)
	}
}
