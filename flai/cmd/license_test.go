package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/license"
)

// S-0231: flai license prints the license built into the binary, whole.
func TestLicensePrintsTheEmbeddedLicense(t *testing.T) {
	out, errOut, code := runCLI(t, "license")
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errOut)
	}
	if out != license.Text() {
		t.Fatalf("output is not the license text:\n%s", out)
	}
	if !strings.HasPrefix(out, "# system-flow Shield License 1.0\n") || !strings.Contains(out, "## Prohibited Uses") {
		t.Fatalf("output lacks the license's heading or its prohibited uses:\n%s", out)
	}

	out, _, code = runCLI(t, "license", "--json")
	if code != 0 {
		t.Fatalf("json exit %d", code)
	}
	var v map[string]string
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("not json: %v\n%s", err, out)
	}
	if v["name"] != "system-flow Shield License 1.0" || v["text"] != license.Text() {
		t.Fatalf("json name=%q, text matches=%v", v["name"], v["text"] == license.Text())
	}
}
