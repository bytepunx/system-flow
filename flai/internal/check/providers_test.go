package check

import (
	"strings"
	"testing"
)

// S-0349: the manifest's providers may override a host entry's api, base_url,
// and models; a key_env there is an error on the providers key naming the
// host setting that holds it, and a valid override is not reported.
func TestManifestProvidersAreChecked(t *testing.T) {
	for _, extra := range []string{
		"",
		"providers:\n  litellm:\n    base_url: https://litellm.example.com\n    models:\n      opus: anthropic/claude-opus-5-5\n",
	} {
		if got := findings(t, planningProject(t, extra), "manifest.providers"); len(got) != 0 {
			t.Errorf("valid providers %q: %+v", extra, got)
		}
	}
	repo := planningProject(t, "providers:\n  litellm:\n    key_env: LITELLM_API_KEY\n")
	got := findings(t, repo, "manifest.providers")
	if len(got) != 1 {
		t.Fatalf("want one error, the key_env: %+v", got)
	}
	if f := got[0]; f.Level != Error || f.Path != "system-flow.yaml" || f.Line != 8 || !strings.Contains(f.Message, "flai serve agent provider litellm --key-env") {
		t.Errorf("finding: %+v, want an error on line 8 naming flai serve agent provider litellm --key-env", f)
	}
}
