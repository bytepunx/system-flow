package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
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

// githubOrigin is the system's git, but for origin's URL, which it answers
// as acme/widgets on GitHub: flai pushes to the scratch bare remote origin
// is, and calls the stand-in GitHub about the repository.
type githubOrigin struct{ execx.System }

func (g githubOrigin) Run(dir, name string, args ...string) (string, error) {
	if name == "git" && slices.Equal(args, []string{"remote", "get-url", "origin"}) {
		return "git@github.com:acme/widgets.git", nil
	}
	return g.System.Run(dir, name, args...)
}

// runForge runs the CLI in root as runIn does, with git through
// githubOrigin, and fails t if the sentinel token appears in what it printed.
func runForge(t *testing.T, root string, args ...string) (string, string, int) {
	t.Helper()
	t.Setenv("FLAI_CACHE_DIR", filepath.Join(t.TempDir(), "cache"))
	var out, errOut bytes.Buffer
	a := &app{out: &out, errOut: &errOut, cwd: root, runner: githubOrigin{},
		clock: func() time.Time { return time.Date(2026, 9, 15, 21, 0, 0, 0, time.UTC) }}
	cmd := newRootCmdWith(a)
	cmd.SetArgs(args)
	code := 0
	if err := cmd.Execute(); err != nil {
		var ee *exitError
		code = 1
		if errors.As(err, &ee) {
			code = ee.code
		}
		if !errors.As(err, &ee) || ee.msg != "" {
			a.fail(err)
		}
	}
	if strings.Contains(out.String()+errOut.String(), forgeSentinel) {
		t.Fatalf("the token was printed:\nstdout: %s\nstderr: %s", out.String(), errOut.String())
	}
	return out.String(), errOut.String(), code
}

// standInPull is a pull request as the stand-in GitHub keeps and answers it.
type standInPull struct {
	Number  int    `json:"number"`
	HTMLURL string `json:"html_url"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	State   string `json:"state"`
	Head    struct {
		Ref string `json:"ref"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

// pullsStandIn is a stand-in GitHub keeping the pull requests of
// acme/widgets, sent the sentinel token, and counting what it is asked.
type pullsStandIn struct {
	mu                    sync.Mutex
	pulls                 []standInPull
	calls, posts, patches int
}

// githubPullsStandIn points flai at a new pullsStandIn.
func githubPullsStandIn(t *testing.T) *pullsStandIn {
	t.Helper()
	gh := &pullsStandIn{}
	srv := httptest.NewServer(http.HandlerFunc(gh.serve))
	t.Cleanup(srv.Close)
	old := githubAPIBase
	githubAPIBase = srv.URL
	t.Cleanup(func() { githubAPIBase = old })
	return gh
}

func (gh *pullsStandIn) serve(w http.ResponseWriter, r *http.Request) {
	gh.mu.Lock()
	defer gh.mu.Unlock()
	gh.calls++
	const pulls = "/repos/acme/widgets/pulls"
	if r.Header.Get("Authorization") != "Bearer "+forgeSentinel {
		http.Error(w, `{"message":"Bad credentials"}`, http.StatusUnauthorized)
		return
	}
	var in struct{ Title, Head, Base, Body string }
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&in)
	}
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.Method == http.MethodGet && r.URL.Path == pulls:
		open := []standInPull{}
		for _, p := range gh.pulls {
			if "acme:"+p.Head.Ref == r.URL.Query().Get("head") {
				open = append(open, p)
			}
		}
		_ = json.NewEncoder(w).Encode(open)
	case r.Method == http.MethodPost && r.URL.Path == pulls:
		gh.posts++
		p := standInPull{Number: len(gh.pulls) + 1, Title: in.Title, Body: in.Body, State: "open"}
		p.HTMLURL = fmt.Sprintf("https://github.com/acme/widgets/pull/%d", p.Number)
		p.Head.Ref, p.Base.Ref = in.Head, in.Base
		gh.pulls = append(gh.pulls, p)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(p)
	case r.Method == http.MethodPatch && strings.HasPrefix(r.URL.Path, pulls+"/"):
		n, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, pulls+"/"))
		if n < 1 || n > len(gh.pulls) {
			http.Error(w, `{"message":"Not Found"}`, http.StatusNotFound)
			return
		}
		gh.patches++
		p := &gh.pulls[n-1]
		p.Title, p.Body, p.Base.Ref = in.Title, in.Body, in.Base
		_ = json.NewEncoder(w).Encode(p)
	default:
		http.Error(w, `{"message":"unexpected request"}`, http.StatusTeapot)
	}
}

// counts is how many calls, creations, and updates the stand-in has had.
func (gh *pullsStandIn) counts() (calls, posts, patches int) {
	gh.mu.Lock()
	defer gh.mu.Unlock()
	return gh.calls, gh.posts, gh.patches
}

// pullRequestProject is a git project under acceptance.mode mode, with a
// host configuration of its own, GITHUB_TOKEN empty, a scratch bare remote
// as origin, and story S-0001, "Ship it", with two criteria, moved to status.
func pullRequestProject(t *testing.T, mode, status string) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("GITHUB_TOKEN", "")
	root := t.TempDir()
	m := "version: 1\nname: t\nkey: t\nlayout:\n  design: design\n  docs: docs\n  wip: wip\nacceptance:\n  mode: " + mode + "\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".gitignore"), []byte(".flai-cache/\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"wip/kanban/epics", "wip/kanban/stories", "wip/kanban/tasks", "wip/agents"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	bare := filepath.Join(t.TempDir(), "origin.git")
	gitIn(t, root, "init", "-q", "-b", "main")
	gitIn(t, filepath.Dir(bare), "init", "-q", "--bare", "-b", "main", bare)
	gitIn(t, root, "remote", "add", "origin", bare)

	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 9, 15, 20, 0, 0, 0, time.UTC)
	epic, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Forge", Owner: "alex", Now: now})
	if err != nil {
		t.Fatal(err)
	}
	st, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: "Ship it", Parent: epic.ID, Owner: "alex", Touches: []string{"docs/ship.md"}, Now: now})
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(st.Path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(st.Path, []byte(strings.Replace(string(data), "## Acceptance criteria\n", "## Acceptance criteria\n- [ ] works\n- [ ] pushes\n", 1)), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, to := range []string{workitem.Ready, workitem.InProgress, workitem.Review} {
		if to == workitem.Review {
			if status == workitem.InProgress {
				break
			}
			if _, err := repo.Create(workitem.NewOptions{Type: workitem.Task, Title: "Work", Parent: st.ID, Owner: "builder", Now: now}); err != nil {
				t.Fatal(err)
			}
		}
		it, err := repo.Get(st.ID)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := repo.Transition(it, to, "alex", "the test moves it", now); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// enableForge turns the forge host action on for root.
func enableForge(t *testing.T, root string) {
	t.Helper()
	if _, errOut, code := runIn(t, root, "serve", "enable", "forge"); code != 0 {
		t.Fatalf("serve enable forge: %d %s", code, errOut)
	}
}

// S-0393 criteria 1, 2, and 6: on the operator's word, with the forge host
// action off, flai forge push pushes the branch to a scratch bare remote and
// opens the pull request; again with nothing new, it is up to date and
// leaves the request as it is; flai forge status then shows the request and
// origin's branch matching, without calling the forge.
func TestForgePushOpensThePullRequestThenIsUpToDate(t *testing.T) {
	root := pullRequestProject(t, "pull-request", workitem.Review)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "project")
	gitIn(t, root, "branch", "story/S-0001")
	head := gitIn(t, root, "rev-parse", "refs/heads/story/S-0001")
	gh := githubPullsStandIn(t)
	t.Setenv("GITHUB_TOKEN", forgeSentinel)
	const url = "https://github.com/acme/widgets/pull/1"

	out, errOut, code := runForge(t, root, "forge", "push", "S-1")
	want := "story: S-0001\noutcome: pushed\nbranch: pushed story/S-0001 at " + head[:12] + " to origin\npull request: #1 opened, " + url + "\n"
	if code != 0 || out != want {
		t.Fatalf("the first push: %d\nstdout: %q\nwant:   %q\nstderr: %s", code, out, want, errOut)
	}
	if onOrigin := gitIn(t, root, "ls-remote", "origin", "refs/heads/story/S-0001"); !strings.HasPrefix(onOrigin, head+"\t") {
		t.Errorf("origin's story/S-0001 = %q, want %s", onOrigin, head)
	}
	if _, posts, _ := gh.counts(); posts != 1 || gh.pulls[0].Title != "[S-0001] Ship it" || gh.pulls[0].Base.Ref != "main" || gh.pulls[0].Head.Ref != "story/S-0001" || !strings.Contains(gh.pulls[0].Body, "- [ ] works") {
		t.Errorf("the request opened: %d posts, %+v", posts, gh.pulls)
	}

	out, errOut, code = runForge(t, root, "forge", "push", "S-0001")
	want = "story: S-0001\noutcome: up-to-date\nbranch: origin had story/S-0001 at " + head[:12] + " already\npull request: #1 unchanged, " + url + "\n"
	if code != 0 || out != want {
		t.Fatalf("a push with nothing new: %d\nstdout: %q\nwant:   %q\nstderr: %s", code, out, want, errOut)
	}
	out, errOut, code = runForge(t, root, "forge", "push", "S-0001", "--json")
	var st map[string]any
	if err := json.Unmarshal([]byte(out), &st); err != nil || code != 0 {
		t.Fatalf("push --json: %d %v\nstdout: %s\nstderr: %s", code, err, out, errOut)
	}
	if st["story"] != "S-0001" || st["outcome"] != "up-to-date" || st["request"] != "unchanged" || st["pull_request"] != 1.0 || st["url"] != url || st["local_head"] != head || st["remote_head"] != head || st["last_pushed"] != head {
		t.Errorf("push --json = %v", st)
	}
	if _, posts, patches := gh.counts(); posts != 1 || patches != 0 {
		t.Errorf("the request was opened %d times and updated %d; want once and never", posts, patches)
	}

	before, _, _ := gh.counts()
	out, errOut, code = runForge(t, root, "forge", "status", "S-0001")
	want = "story: S-0001 (review)\nmode: pull-request (forge github)\nhost action: forge is off\ntoken: read from GITHUB_TOKEN, set\n" +
		"pull request: #1, " + url + "\nbranch: origin's story/S-0001 matches the local head at " + head[:12] + "\nlast look: up-to-date, at 2026-09-15T21:00:00Z\n"
	if code != 0 || out != want {
		t.Fatalf("status after a push: %d\nstdout: %q\nwant:   %q\nstderr: %s", code, out, want, errOut)
	}
	out, errOut, code = runForge(t, root, "forge", "status", "S-0001", "--json")
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil || code != 0 {
		t.Fatalf("status --json: %d %v\nstdout: %s\nstderr: %s", code, err, out, errOut)
	}
	wantJSON := map[string]any{"story": "S-0001", "status": "review", "mode": "pull-request", "forge": "github", "host_action": false, "token_env": "GITHUB_TOKEN", "token_set": true,
		"pull_request": 1.0, "url": url, "branch": "story/S-0001", "local_head": head, "remote_head": head, "matches": true, "outcome": "up-to-date", "at": "2026-09-15T21:00:00Z"}
	if !reflect.DeepEqual(got, wantJSON) {
		t.Errorf("status --json = %v\nwant %v", got, wantJSON)
	}
	if after, _, _ := gh.counts(); after != before {
		t.Errorf("flai forge status called the forge %d times; want none", after-before)
	}

	// S-0394 criterion 2: once flai serve has read the pull request back,
	// flai forge status prints its checks and reviews from the kept state,
	// and --json gives it, without the events, still calling no forge
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	kept := keptPullRequestFixture("S-0001")
	if err := forge.WriteStatus(repo.CacheDir(), kept); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runForge(t, root, "forge", "status", "S-0001")
	want = "pull request: #1, " + url + "\n" +
		"  read at abc123def456, 2026-10-17T16:00:00Z\n" +
		"  check build: failed https://github.com/o/r/runs/1 (2 tests failed)\n" +
		"  check lint: pending\n" +
		"  review alice: changes-requested https://github.com/o/r/pull/7#pullrequestreview-1\n" +
		"  review bob: approved\n" +
		"branch: "
	if code != 0 || !strings.Contains(out, want) {
		t.Errorf("status with the pull request read: %d\nstdout: %s\nwant it to hold: %s\nstderr: %s", code, out, want, errOut)
	}
	out, errOut, code = runForge(t, root, "forge", "status", "S-0001", "--json")
	var read struct {
		PullRequest     int           `json:"pull_request"`
		PullRequestRead *forge.Status `json:"pull_request_read"`
	}
	if err := json.Unmarshal([]byte(out), &read); err != nil || code != 0 {
		t.Fatalf("status --json with the pull request read: %d %v\nstdout: %s\nstderr: %s", code, err, out, errOut)
	}
	kept.Events = nil
	if read.PullRequest != 1 || !reflect.DeepEqual(read.PullRequestRead, &kept) {
		t.Errorf("status --json: pull_request %d, pull_request_read %+v\nwant %+v", read.PullRequest, read.PullRequestRead, kept)
	}
	if after, _, _ := gh.counts(); after != before {
		t.Errorf("flai forge status called the forge %d times; want none", after-before)
	}

	// a kept state that cannot be read fails the command with what clears it
	if err := os.WriteFile(forge.StatusPath(repo.CacheDir(), "S-0001"), []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, errOut, code := runForge(t, root, "forge", "status", "S-0001"); code != 1 || !strings.Contains(errOut, "delete it, and flai serve reads the pull request afresh") {
		t.Errorf("status with a corrupt kept state: %d %s", code, errOut)
	}
}

// S-0393 criterion 6: flai forge push does nothing, and says why with exit
// status 1, under local acceptance; and fails naming the variable when the
// token is unset, before any git or forge call.
func TestForgePushIsSkippedOrFailsWithItsReason(t *testing.T) {
	for _, tc := range []struct {
		name, mode, out string
		token           bool
		errHas          []string
	}{
		{name: "local acceptance", mode: "local", token: true, out: "story: S-0001\noutcome: skipped\n",
			errHas: []string{"acceptance.mode is local", "set acceptance.mode: pull-request in system-flow.yaml"}},
		{name: "token unset", mode: "pull-request", out: "story: S-0001\noutcome: failed\n",
			errHas: []string{"S-0001's pull request cannot be opened without the github token", "GITHUB_TOKEN, which is not set", "--token-env"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := pullRequestProject(t, tc.mode, workitem.Review)
			gh := githubPullsStandIn(t)
			if tc.token {
				t.Setenv("GITHUB_TOKEN", forgeSentinel)
			}
			out, errOut, code := runForge(t, root, "forge", "push", "S-0001")
			if code != 1 || out != tc.out {
				t.Fatalf("exit %d, want 1\nstdout: %q, want %q\nstderr: %s", code, out, tc.out, errOut)
			}
			for _, s := range tc.errHas {
				if !strings.Contains(errOut, s) {
					t.Errorf("stderr lacks %q: %s", s, errOut)
				}
			}
			if calls, _, _ := gh.counts(); calls != 0 {
				t.Errorf("the forge was called %d times; want none", calls)
			}
		})
	}
}

// S-0393 criterion 2: before anything is pushed, flai forge status says why,
// the first reason that applies, with exit status 0 and no forge call.
func TestForgeStatusSaysWhyNothingWasPushed(t *testing.T) {
	for _, tc := range []struct {
		name, mode, status string
		hostAction, token  bool
		why                []string
	}{
		{name: "mode local", mode: "local", status: workitem.Review, hostAction: true, token: true,
			why: []string{"acceptance.mode is local, so S-0001 is accepted on this machine"}},
		{name: "host action off", mode: "pull-request", status: workitem.Review, token: true,
			why: []string{"the forge host action is off", "flai serve enable forge", "flai forge push S-0001 pushes the story on your word"}},
		{name: "token unset", mode: "pull-request", status: workitem.Review, hostAction: true,
			why: []string{"GITHUB_TOKEN, which is not set", "flai serve reads the environment it started with"}},
		{name: "not in review", mode: "pull-request", status: workitem.InProgress, hostAction: true, token: true,
			why: []string{"S-0001 is in-progress, not in review"}},
		{name: "nothing stops it", mode: "pull-request", status: workitem.Review, hostAction: true, token: true,
			why: []string{"nothing stops flai serve from pushing S-0001", "flai forge push S-0001"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := pullRequestProject(t, tc.mode, tc.status)
			if tc.hostAction {
				enableForge(t, root)
			}
			gh := githubPullsStandIn(t)
			if tc.token {
				t.Setenv("GITHUB_TOKEN", forgeSentinel)
			}
			out, errOut, code := runForge(t, root, "forge", "status", "S-0001")
			if code != 0 {
				t.Fatalf("exit %d\nstdout: %s\nstderr: %s", code, out, errOut)
			}
			for _, s := range []string{"pull request: none yet\n", "branch: neither this repository nor origin has story/S-0001\n", "last look: none yet\n", "\nwhy: "} {
				if !strings.Contains(out, s) {
					t.Errorf("stdout lacks %q: %s", s, out)
				}
			}
			why := out[strings.Index(out, "\nwhy: ")+1:]
			for _, s := range tc.why {
				if !strings.Contains(why, s) {
					t.Errorf("why lacks %q: %s", s, why)
				}
			}
			if calls, _, _ := gh.counts(); calls != 0 {
				t.Errorf("the forge was called %d times; want none", calls)
			}
		})
	}
}

// After flai forge push fails, flai forge status names the token while it
// is unset, and the failure, when and why, once nothing before it applies.
func TestForgeStatusNamesTheLastFailure(t *testing.T) {
	root := pullRequestProject(t, "pull-request", workitem.Review)
	enableForge(t, root)
	githubPullsStandIn(t)
	if _, errOut, code := runForge(t, root, "forge", "push", "S-0001"); code != 1 {
		t.Fatalf("push without the token: %d %s", code, errOut)
	}
	out, errOut, code := runForge(t, root, "forge", "status", "S-0001")
	if code != 0 || !strings.Contains(out, "last look: failed, at 2026-09-15T21:00:00Z\n") || !strings.Contains(out, "why: the github token is read from the variable GITHUB_TOKEN, which is not set") {
		t.Fatalf("status with the token unset: %d\nstdout: %s\nstderr: %s", code, out, errOut)
	}
	t.Setenv("GITHUB_TOKEN", forgeSentinel)
	out, errOut, code = runForge(t, root, "forge", "status", "S-0001")
	if code != 0 || !strings.Contains(out, "token: read from GITHUB_TOKEN, set\n") || !strings.Contains(out, "why: the last push, at 2026-09-15T21:00:00Z, failed: S-0001's pull request cannot be opened without the github token") {
		t.Fatalf("status with the token set: %d\nstdout: %s\nstderr: %s", code, out, errOut)
	}
}

// flai forge status --json gives the same findings, and says why origin
// could not be read when it cannot.
func TestForgeStatusJSON(t *testing.T) {
	root := pullRequestProject(t, "pull-request", workitem.Review)
	gitIn(t, root, "remote", "remove", "origin")
	out, errOut, code := runForge(t, root, "forge", "status", "S-0001", "--json")
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil || code != 0 {
		t.Fatalf("status --json: %d %v\nstdout: %s\nstderr: %s", code, err, out, errOut)
	}
	remoteErr, _ := got["remote_error"].(string)
	why, _ := got["why"].(string)
	delete(got, "remote_error")
	delete(got, "why")
	want := map[string]any{"story": "S-0001", "status": "review", "mode": "pull-request", "forge": "github", "host_action": false, "token_env": "GITHUB_TOKEN", "token_set": false, "branch": "story/S-0001", "matches": false}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("status --json = %v\nwant %v", got, want)
	}
	if !strings.Contains(remoteErr, "could not read story/S-0001 on origin") || !strings.Contains(why, "the forge host action is off") {
		t.Errorf("remote_error %q, why %q", remoteErr, why)
	}
	out, _, _ = runForge(t, root, "forge", "status", "S-0001")
	if !strings.Contains(out, "branch: could not read story/S-0001 on origin") {
		t.Errorf("the branch line does not say origin could not be read: %s", out)
	}
}
