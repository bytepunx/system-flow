package images

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDockerfileIsTheTreesFile checks that each embedded Dockerfile is the
// file beside this package, byte for byte.
func TestDockerfileIsTheTreesFile(t *testing.T) {
	for _, image := range []string{Bedrock, "claude-code"} {
		got, err := Dockerfile(image)
		if err != nil {
			t.Fatalf("Dockerfile(%q): %v", image, err)
		}
		want, err := os.ReadFile(filepath.Join(image, "Dockerfile"))
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("Dockerfile(%q) is not %s/Dockerfile", image, image)
		}
		if p := Path(image); p != "flai/images/"+image+"/Dockerfile" {
			t.Errorf("Path(%q) = %q", image, p)
		}
	}
}

// TestDockerfileRefusesAnImageItCarriesNone checks the refusal names the
// images it carries.
func TestDockerfileRefusesAnImageItCarriesNone(t *testing.T) {
	_, err := Dockerfile("system-flow")
	if err == nil || !strings.Contains(err.Error(), "bedrock's and claude-code's") {
		t.Fatalf("Dockerfile(system-flow) error = %v, want one naming the Dockerfiles carried", err)
	}
}
