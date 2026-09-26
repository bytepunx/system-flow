package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

func TestWebURL(t *testing.T) {
	for remote, want := range map[string]string{
		"git@github.com:arobson/astro-blog.git":           "https://github.com/arobson/astro-blog",
		"github.com:owner/repo":                           "https://github.com/owner/repo",
		"ssh://git@github.com/owner/repo.git":             "https://github.com/owner/repo",
		"ssh://git@gitlab.example.com:2222/grp/sub/r.git": "https://gitlab.example.com/grp/sub/r",
		"git://example.org/owner/repo":                    "https://example.org/owner/repo",
		"https://github.com/owner/repo.git\n":             "https://github.com/owner/repo",
		"https://user:token@github.com/owner/repo.git":    "https://github.com/owner/repo",
		"http://git.local/owner/repo/":                    "https://git.local/owner/repo",
		"/srv/git/repo.git":                               "",
		"../repo":                                         "",
		"file:///srv/git/repo.git":                        "",
		"":                                                "",
	} {
		if got := webURL(remote); got != want {
			t.Errorf("webURL(%q) = %q, want %q", remote, got, want)
		}
	}
}

func TestCheckRepoURL(t *testing.T) {
	for _, ok := range []string{"", "https://github.com/owner/repo", "http://git.local:3000/owner/repo"} {
		if err := checkRepoURL(ok); err != nil {
			t.Errorf("%q refused: %v", ok, err)
		}
	}
	for bad, why := range map[string]string{
		"https://github.com:owner/repo":       "is not a URL",
		"git@github.com:owner/repo.git":       "is not a URL",
		"github.com/owner/repo":               "is not an http or https URL",
		"ssh://git@github.com/owner/repo.git": "is not an http or https URL",
		"https:///owner/repo":                 "has no host",
		"https://github.com":                  "names no repository",
		"https://github.com/":                 "names no repository",
	} {
		err := checkRepoURL(bad)
		if err == nil || !strings.Contains(err.Error(), why) || !strings.Contains(err.Error(), "such as https://github.com/owner/repo") {
			t.Errorf("%q: %v, want %q", bad, err, why)
		}
	}
}

// S-0120: an import offers the origin remote as an https URL for repo_url,
// and refuses a value that is not a URL before writing anything.
func TestImportTakesRepoURLFromOrigin(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	if _, err := os.Stat(filepath.Join("..", "..", "template", "template.yaml")); err != nil {
		t.Skip("prototype template not present")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := legacyRepo(t)
	if _, err := (execx.System{}).Run(root, "git", "remote", "add", "origin", "git@github.com:arobson/astro-blog.git"); err != nil {
		t.Fatal(err)
	}

	_, errOut, code := runIn(t, ".", "import", root, "--template", "../../template", "--yes", "--var", "repo_url=https://github.com:arobson/astro-blog")
	if code == 0 || !strings.Contains(errOut, `https://github.com:arobson/astro-blog\" is not a URL (invalid port`) {
		t.Errorf("a value that is not a URL: %d %s", code, errOut)
	}
	if _, err := os.Stat(filepath.Join(root, manifest.File)); err == nil {
		t.Error("the manifest was written before the refusal")
	}

	if _, errOut, code := runIn(t, ".", "import", root, "--template", "../../template", "--yes"); code != 0 {
		t.Fatalf("import: %d %s", code, errOut)
	}
	m, err := manifest.Load(filepath.Join(root, manifest.File))
	if err != nil || m.Repo != "https://github.com/arobson/astro-blog" {
		t.Errorf("repo: %q %v", m.Repo, err)
	}
}
