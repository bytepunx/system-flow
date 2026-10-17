package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

// forgeSentinel is the token the tests put in the environment; no output may carry it.
const forgeSentinel = "ghp_forge-check-sentinel-0123456789"

// forgeProject is a git project whose origin is remote, none when it is "",
// with a host configuration of its own and GITHUB_TOKEN empty.
func forgeProject(t *testing.T, remote string) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("GITHUB_TOKEN", "")
	root := tempProject(t)
	git(t, root, "init", "-q", "-b", "main")
	if remote != "" {
		git(t, root, "remote", "add", "origin", remote)
	}
	return root
}

// githubStandIn points flai at a stand-in GitHub that answers GET
// /repos/acme/widgets, sent the sentinel token, with status and body, and
// answers how many calls it has had.
func githubStandIn(t *testing.T, status int, body string) (calls func() int) {
	t.Helper()
	var mu sync.Mutex
	n := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		n++
		mu.Unlock()
		if r.Method != http.MethodGet || r.URL.Path != "/repos/acme/widgets" || r.Header.Get("Authorization") != "Bearer "+forgeSentinel {
			http.Error(w, `{"message":"unexpected request"}`, http.StatusTeapot)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)
	old := githubAPIBase
	githubAPIBase = srv.URL
	t.Cleanup(func() { githubAPIBase = old })
	return func() int {
		mu.Lock()
		defer mu.Unlock()
		return n
	}
}

// runForgeCheck runs flai forge check in root and fails t if the sentinel
// token appears in what it printed.
func runForgeCheck(t *testing.T, root string, args ...string) (string, string, int) {
	t.Helper()
	out, errOut, code := runIn(t, root, append([]string{"forge", "check"}, args...)...)
	if strings.Contains(out+errOut, forgeSentinel) {
		t.Fatalf("the token was printed:\nstdout: %s\nstderr: %s", out, errOut)
	}
	return out, errOut, code
}

func TestForgeCheck(t *testing.T) {
	const pushable = `{"full_name":"acme/widgets","permissions":{"pull":true,"push":true}}`
	const readOnly = `{"full_name":"acme/widgets","permissions":{"pull":true,"push":false}}`
	const passed = "forge: github\nrepository: acme/widgets (from origin)\ntoken: read from GITHUB_TOKEN\nread: yes\npush: yes\n"
	for _, tc := range []struct {
		name    string
		remote  string
		token   bool
		status  int
		body    string
		code    int
		out     string // the whole of standard output when given
		has     []string
		errHas  []string
		noCalls bool
	}{
		{name: "scp remote, push", remote: "git@github.com:acme/widgets.git", token: true, status: 200, body: pushable, out: passed},
		{name: "https remote, push", remote: "https://github.com/acme/widgets", token: true, status: 200, body: pushable, out: passed},
		{name: "https remote with credentials", remote: "https://x-access-token:" + forgeSentinel + "@github.com/acme/widgets.git", token: true, status: 200, body: pushable, out: passed},
		{
			name: "read but not push", remote: "git@github.com:acme/widgets.git", token: true, status: 200, body: readOnly, code: 1,
			has:    []string{"read: yes\n", "push: no\n"},
			errHas: []string{"the token in GITHUB_TOKEN can read acme/widgets but not push to it", "contents: write"},
		},
		{
			name: "token rejected", remote: "git@github.com:acme/widgets.git", token: true, status: 401, body: `{"message":"Bad credentials"}`, code: 1,
			has:    []string{"token: read from GITHUB_TOKEN\n", "read: no\n"},
			errHas: []string{"github did not accept the token in GITHUB_TOKEN", "Bad credentials", "flai serve forge github --token-env <NAME>"},
		},
		{
			name: "repository not visible", remote: "git@github.com:acme/widgets.git", token: true, status: 404, body: `{"message":"Not Found"}`, code: 1,
			has:    []string{"repository: acme/widgets (from origin)\n", "read: no\n"},
			errHas: []string{"github has no repository acme/widgets that the token in GITHUB_TOKEN can see"},
		},
		{
			name: "forbidden", remote: "git@github.com:acme/widgets.git", token: true, status: 403, body: `{"message":"API rate limit exceeded"}`, code: 1,
			errHas: []string{"github refused the token in GITHUB_TOKEN", "rate limit"},
		},
		{
			name: "variable unset", remote: "git@github.com:acme/widgets.git", code: 1, noCalls: true,
			out:    "forge: github\nrepository: acme/widgets (from origin)\ntoken: read from GITHUB_TOKEN\n",
			errHas: []string{"GITHUB_TOKEN, which is not set", "flai serve forge github --token-env <NAME>"},
		},
		{
			name: "no origin", token: true, code: 1, noCalls: true,
			out:    "forge: github\ntoken: read from GITHUB_TOKEN\n",
			errHas: []string{"no origin remote", "git remote add origin"},
		},
		{
			name: "origin on another forge", remote: "https://oauth2:" + forgeSentinel + "@gitlab.com/acme/widgets.git", token: true, code: 1, noCalls: true,
			out:    "forge: github\ntoken: read from GITHUB_TOKEN\n",
			errHas: []string{"origin is not a repository on github", "gitlab.com, not github.com"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := forgeProject(t, tc.remote)
			calls := githubStandIn(t, tc.status, tc.body)
			if tc.token {
				t.Setenv("GITHUB_TOKEN", forgeSentinel)
			}
			out, errOut, code := runForgeCheck(t, root)
			if code != tc.code {
				t.Fatalf("exit %d, want %d\nstdout: %s\nstderr: %s", code, tc.code, out, errOut)
			}
			if tc.out != "" && out != tc.out {
				t.Errorf("stdout = %q, want %q", out, tc.out)
			}
			for _, s := range tc.has {
				if !strings.Contains(out, s) {
					t.Errorf("stdout lacks %q: %s", s, out)
				}
			}
			for _, s := range tc.errHas {
				if !strings.Contains(errOut, s) {
					t.Errorf("stderr lacks %q: %s", s, errOut)
				}
			}
			if tc.noCalls && calls() != 0 {
				t.Errorf("the forge was called %d times; want none", calls())
			}
		})
	}
}

// The variable the host names for the forge is the one read and named.
func TestForgeCheckReadsTheVariableTheHostNames(t *testing.T) {
	root := forgeProject(t, "git@github.com:acme/widgets.git")
	githubStandIn(t, 200, `{"permissions":{"push":true}}`)
	if _, errOut, code := runIn(t, root, "serve", "forge", "github", "--token-env", "FLAI_GH_TOKEN"); code != 0 {
		t.Fatalf("serve forge: %d %s", code, errOut)
	}
	t.Setenv("FLAI_GH_TOKEN", "")
	t.Setenv("GITHUB_TOKEN", forgeSentinel)
	out, errOut, code := runForgeCheck(t, root)
	if code != 1 || !strings.Contains(out, "token: read from FLAI_GH_TOKEN\n") || !strings.Contains(errOut, "FLAI_GH_TOKEN, which is not set") {
		t.Fatalf("unset FLAI_GH_TOKEN: %d\nstdout: %s\nstderr: %s", code, out, errOut)
	}
	t.Setenv("FLAI_GH_TOKEN", forgeSentinel)
	t.Setenv("GITHUB_TOKEN", "")
	out, errOut, code = runForgeCheck(t, root)
	if code != 0 || out != "forge: github\nrepository: acme/widgets (from origin)\ntoken: read from FLAI_GH_TOKEN\nread: yes\npush: yes\n" {
		t.Fatalf("FLAI_GH_TOKEN set: %d\nstdout: %s\nstderr: %s", code, out, errOut)
	}
}

func TestForgeCheckJSON(t *testing.T) {
	root := forgeProject(t, "git@github.com:acme/widgets.git")
	t.Setenv("GITHUB_TOKEN", forgeSentinel)
	githubStandIn(t, 200, `{"permissions":{"push":true}}`)
	out, errOut, code := runForgeCheck(t, root, "--json")
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil || code != 0 {
		t.Fatalf("passing --json: %d %v\nstdout: %s\nstderr: %s", code, err, out, errOut)
	}
	want := map[string]any{"forge": "github", "repository": "acme/widgets", "token_env": "GITHUB_TOKEN", "read": true, "push": true, "ok": true}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("passing --json = %v, want %v", got, want)
	}

	githubStandIn(t, 200, `{"permissions":{"push":false}}`)
	out, errOut, code = runForgeCheck(t, root, "--json")
	got = nil
	if err := json.Unmarshal([]byte(out), &got); err != nil || code != 1 {
		t.Fatalf("failing --json: %d %v\nstdout: %s\nstderr: %s", code, err, out, errOut)
	}
	reason, _ := got["reason"].(string)
	delete(got, "reason")
	want = map[string]any{"forge": "github", "repository": "acme/widgets", "token_env": "GITHUB_TOKEN", "read": true, "push": false, "ok": false}
	if !reflect.DeepEqual(got, want) || !strings.Contains(reason, "can read acme/widgets but not push") || !strings.Contains(errOut, reason) {
		t.Errorf("failing --json = %v, reason %q\nstderr: %s", got, reason, errOut)
	}
}
