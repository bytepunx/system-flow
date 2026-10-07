package mcpserver

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// releasesStandIn answers as GitHub's API does for bytepunx/system-flow:
// flai releases out of order, with a draft, a prerelease, and another
// component's release among them, and dashboard tags with one that is not a
// plain version. auth records the Authorization header of each request.
func releasesStandIn(t *testing.T) (auth func() []string) {
	t.Helper()
	var mu sync.Mutex
	var seen []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen = append(seen, r.Header.Get("Authorization"))
		mu.Unlock()
		switch r.URL.Path {
		case "/repos/bytepunx/system-flow/releases":
			_, _ = w.Write([]byte(`[
				{"tag_name":"flai/v1.9.0","published_at":"2026-08-01T10:00:00Z"},
				{"tag_name":"flai/v1.10.0","published_at":"2026-09-01T10:00:00Z"},
				{"tag_name":"flai/v1.11.0","draft":true},
				{"tag_name":"flai/v1.12.0-rc1","prerelease":true},
				{"tag_name":"template/v1.0.70","published_at":"2026-09-02T10:00:00Z"},
				{"tag_name":"flai/v1.2.0","published_at":"2026-01-01T10:00:00Z"}
			]`))
		case "/repos/bytepunx/system-flow/git/matching-refs/tags/flaiover/v":
			_, _ = w.Write([]byte(`[{"ref":"refs/tags/flaiover/v0.9.0"},{"ref":"refs/tags/flaiover/v0.40.0"},{"ref":"refs/tags/flaiover/v0.41.0-rc1"},{"ref":"refs/tags/flaiover/v0.10.0"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	t.Setenv("FLAI_RELEASES_API", srv.URL)
	return func() []string {
		mu.Lock()
		defer mu.Unlock()
		return append([]string{}, seen...)
	}
}

// versionsOf is a list in the versions answer as its versions and marks.
func versionsOf(t *testing.T, list any, marks ...string) []string {
	t.Helper()
	entries, ok := list.([]any)
	if !ok {
		t.Fatalf("not a list: %v", list)
	}
	var out []string
	for _, e := range entries {
		m := e.(map[string]any)
		s := m["version"].(string)
		for _, mark := range marks {
			if m[mark] == true {
				s += " " + mark
			}
		}
		out = append(out, s)
	}
	return out
}

func TestVersionsListsTheFlaiAndDashboardReleasesNewestFirstWithTheirMarks(t *testing.T) {
	auth := releasesStandIn(t)
	t.Setenv("GITHUB_TOKEN", "tok")
	t.Setenv("GH_TOKEN", "")
	f := setupWith(t, func(o *Options) { o.Version = "1.9.0" })
	f.repo.Manifest.Flai.Minimum = "1.9.0"

	out, failed := f.call(t, "versions", map[string]any{})
	if failed != "" {
		t.Fatalf("versions failed: %s", failed)
	}
	if got, want := strings.Join(versionsOf(t, out["flai"], "running", "latest", "below_minimum"), ", "), "1.10.0 latest, 1.9.0 running, 1.2.0 below_minimum"; got != want {
		t.Errorf("flai: got %s, want %s", got, want)
	}
	if got, want := strings.Join(versionsOf(t, out["dashboard"], "latest"), ", "), "0.40.0 latest, 0.10.0, 0.9.0"; got != want {
		t.Errorf("dashboard: got %s, want %s", got, want)
	}
	if out["running"] != "1.9.0" || out["minimum"] != "1.9.0" {
		t.Errorf("running and minimum: %v, %v", out["running"], out["minimum"])
	}
	newest := out["flai"].([]any)[0].(map[string]any)
	if newest["tag"] != "flai/v1.10.0" || newest["published"] != "2026-09-01T10:00:00Z" {
		t.Errorf("the newest flai's tag and date: %v", newest)
	}
	dash := out["dashboard"].([]any)[0].(map[string]any)
	if dash["tag"] != "flaiover/v0.40.0" {
		t.Errorf("the newest dashboard's tag: %v", dash)
	}
	if _, ok := dash["published"]; ok {
		t.Errorf("a dashboard tag has no publish date: %v", dash)
	}
	for _, h := range auth() {
		if h != "Bearer tok" {
			t.Errorf("a request went without GITHUB_TOKEN: %q", h)
		}
	}
}

func TestVersionsMarksNoFlaiBelowAMinimumTheProjectDoesNotSet(t *testing.T) {
	releasesStandIn(t)
	f := setupWith(t, func(o *Options) { o.Version = "dev" })
	out, failed := f.call(t, "versions", map[string]any{})
	if failed != "" {
		t.Fatalf("versions failed: %s", failed)
	}
	if got, want := strings.Join(versionsOf(t, out["flai"], "running", "below_minimum"), ", "), "1.10.0, 1.9.0, 1.2.0"; got != want {
		t.Errorf("flai: got %s, want %s", got, want)
	}
	if out["minimum"] != "" {
		t.Errorf("minimum: %v", out["minimum"])
	}
}

// ghRunner answers gh auth token as a signed-in gh CLI does.
type ghRunner struct{ noGitRunner }

func (ghRunner) LookPath(name string) (string, error) {
	if name == "gh" {
		return "/usr/bin/gh", nil
	}
	return "", errors.New(name + " not on PATH")
}

func (ghRunner) Run(_, name string, args ...string) (string, error) {
	if name == "gh" && strings.Join(args, " ") == "auth token" {
		return "gh-session", nil
	}
	return "", errors.New("not gh auth token")
}

func TestVersionsReadsWithTheGhSessionWithoutATokenInTheEnvironment(t *testing.T) {
	auth := releasesStandIn(t)
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GH_TOKEN", "")
	f := setupWith(t, func(o *Options) { o.Runner = ghRunner{} })
	if _, failed := f.call(t, "versions", map[string]any{}); failed != "" {
		t.Fatalf("versions failed: %s", failed)
	}
	seen := auth()
	if len(seen) == 0 {
		t.Fatal("the stand-in was not asked")
	}
	for _, h := range seen {
		if h != "Bearer gh-session" {
			t.Errorf("a request went without gh's session token: %q", h)
		}
	}
}

func TestVersionsSaysWhatItCouldNotList(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	t.Cleanup(srv.Close)
	t.Setenv("FLAI_RELEASES_API", srv.URL)
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GH_TOKEN", "")
	f := setup(t)
	_, failed := f.call(t, "versions", map[string]any{})
	for _, want := range []string{"list the flai/v releases of bytepunx/system-flow", "404", "set GITHUB_TOKEN or run gh auth login"} {
		if !strings.Contains(failed, want) {
			t.Errorf("the refusal does not say %q: %s", want, failed)
		}
	}
}
