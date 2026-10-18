package forge

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// testToken is the token the fake GitHub expects; no test prints it.
const testToken = "fake-github-token-0123456789"

// fakeRequest is one request the fake GitHub received.
type fakeRequest struct {
	Method string
	Path   string
	Query  string
	Body   map[string]any
}

// fakeGitHub is a stand-in for GitHub's REST API: tests register answers on
// it by pattern, and it records every request. Each request must carry
// testToken as a bearer and GitHub's headers; a request to a pattern no
// test registered answers 404, as GitHub does.
type fakeGitHub struct {
	t      *testing.T
	mux    *http.ServeMux
	server *httptest.Server
	mu     sync.Mutex
	got    []fakeRequest
}

// newFakeGitHub starts a fake GitHub that stops when the test ends.
func newFakeGitHub(t *testing.T) *fakeGitHub {
	t.Helper()
	f := &fakeGitHub{t: t, mux: http.NewServeMux()}
	f.server = httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(f.server.Close)
	return f
}

// serve records r, checks its headers, and hands it to the registered answer.
func (f *fakeGitHub) serve(w http.ResponseWriter, r *http.Request) {
	req := fakeRequest{Method: r.Method, Path: r.URL.Path, Query: r.URL.RawQuery}
	if b, _ := io.ReadAll(r.Body); len(b) > 0 {
		if err := json.Unmarshal(b, &req.Body); err != nil {
			f.t.Errorf("%s %s: the body is not a JSON object: %v", r.Method, r.URL.Path, err)
		}
	}
	f.mu.Lock()
	f.got = append(f.got, req)
	f.mu.Unlock()
	if r.Header.Get("Authorization") != "Bearer "+testToken {
		f.t.Errorf("%s %s: the Authorization header does not carry the token as a bearer", r.Method, r.URL.Path)
	}
	for name, want := range map[string]string{"Accept": "application/vnd.github+json", "X-GitHub-Api-Version": githubAPIVersion, "User-Agent": "flai"} {
		if got := r.Header.Get(name); got != want {
			f.t.Errorf("%s %s: header %s = %q, want %q", r.Method, r.URL.Path, name, got, want)
		}
	}
	if _, pattern := f.mux.Handler(r); pattern == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "Not Found"})
		return
	}
	f.mux.ServeHTTP(w, r)
}

// handle answers requests matching pattern with h.
func (f *fakeGitHub) handle(pattern string, h http.HandlerFunc) { f.mux.HandleFunc(pattern, h) }

// reply answers requests matching pattern with status and body as JSON.
func (f *fakeGitHub) reply(pattern string, status int, body any) {
	f.handle(pattern, func(w http.ResponseWriter, _ *http.Request) { writeJSON(w, status, body) })
}

// requests are the requests received so far.
func (f *fakeGitHub) requests() []fakeRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return slices.Clone(f.got)
}

// calls are the requests received so far as "METHOD path".
func (f *fakeGitHub) calls() []string {
	var calls []string
	for _, r := range f.requests() {
		calls = append(calls, r.Method+" "+r.Path)
	}
	return calls
}

// forge is a GitHub forge for octo/repo on the fake.
func (f *fakeGitHub) forge() *GitHub {
	f.t.Helper()
	g, err := NewGitHub(GitHubOptions{
		BaseURL: f.server.URL, Owner: "octo", Name: "repo",
		Token:  func() (string, error) { return testToken, nil },
		Client: f.server.Client(),
	})
	if err != nil {
		f.t.Fatal(err)
	}
	return g
}

// writeJSON writes v as JSON with status.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// githubPullJSON is pull request number from head at sha into main, as GitHub answers it.
func githubPullJSON(number int, head, sha, title string) map[string]any {
	return map[string]any{
		"number": number, "html_url": "https://github.com/octo/repo/pull/" + strconv.Itoa(number),
		"title": title, "body": "the body", "state": "open", "draft": false, "merged_at": nil,
		"head": map[string]any{"ref": head, "sha": sha},
		"base": map[string]any{"ref": "main"},
	}
}

// noToken fails the test when err's text carries the token, without printing either.
func noToken(t *testing.T, err error) {
	t.Helper()
	if err != nil && strings.Contains(err.Error(), testToken) {
		t.Error("an error carries the token")
	}
}

func TestGitHubNameAndRepository(t *testing.T) {
	g := newFakeGitHub(t).forge()
	if g.Name() != "github" || g.Repository() != "octo/repo" {
		t.Errorf("Name, Repository = %q, %q; want github, octo/repo", g.Name(), g.Repository())
	}
}

func TestGitHubAccess(t *testing.T) {
	for _, tc := range []struct {
		name        string
		permissions any
		want        Access
	}{
		{"push", map[string]bool{"pull": true, "push": true}, Access{Read: true, Push: true}},
		{"admin", map[string]bool{"pull": true, "admin": true}, Access{Read: true, Push: true}},
		{"read only", map[string]bool{"pull": true}, Access{Read: true}},
		{"no permissions reported", nil, Access{Read: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeGitHub(t)
			f.reply("GET /repos/octo/repo", http.StatusOK, map[string]any{"full_name": "octo/repo", "permissions": tc.permissions})
			got, err := f.forge().Access(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("Access = %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestGitHubFindPullRequest(t *testing.T) {
	f := newFakeGitHub(t)
	f.handle("GET /repos/octo/repo/pulls", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("head") == "octo:story/S-0001" && r.URL.Query().Get("state") == "open" {
			pull := githubPullJSON(7, "story/S-0001", "abc123", "S-0001 A story")
			pull["draft"] = true
			writeJSON(w, http.StatusOK, []any{pull})
			return
		}
		writeJSON(w, http.StatusOK, []any{})
	})
	g := f.forge()
	pr, found, err := g.FindPullRequest(context.Background(), "story/S-0001")
	if err != nil || !found {
		t.Fatalf("FindPullRequest = found %v, %v; want found", found, err)
	}
	want := PullRequest{
		Number: 7, URL: "https://github.com/octo/repo/pull/7", Head: "story/S-0001", Base: "main",
		HeadSHA: "abc123", Title: "S-0001 A story", Body: "the body", State: PullRequestOpen, Draft: true,
	}
	if pr != want {
		t.Errorf("FindPullRequest = %+v, want %+v", pr, want)
	}
	if _, found, err := g.FindPullRequest(context.Background(), "story/S-0002"); err != nil || found {
		t.Errorf("FindPullRequest for a branch with none open = found %v, %v; want not found", found, err)
	}
	if _, _, err := g.FindPullRequest(context.Background(), ""); err == nil {
		t.Error("FindPullRequest with no head answered no error")
	}
}

// S-0397: a pull request is read by its number in any state, so that one
// merged on the forge is found though FindPullRequest sees open ones only.
func TestGitHubPullRequestByNumber(t *testing.T) {
	f := newFakeGitHub(t)
	merged := githubPullJSON(7, "story/S-0001", "abc123", "S-0001 A story")
	merged["state"], merged["merged_at"] = "closed", "2026-10-17T10:00:00Z"
	f.reply("GET /repos/octo/repo/pulls/7", http.StatusOK, merged)
	g := f.forge()
	pr, err := g.PullRequest(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	if pr.Number != 7 || pr.Head != "story/S-0001" || pr.HeadSHA != "abc123" || pr.State != PullRequestMerged {
		t.Errorf("PullRequest(7) = %+v, want #7 from story/S-0001 at abc123, merged", pr)
	}
	_, err = g.PullRequest(context.Background(), 8)
	if !errors.Is(err, ErrNotFound) || !strings.Contains(err.Error(), "#8") {
		t.Errorf("PullRequest(8), which GitHub does not have = %v, want ErrNotFound naming #8", err)
	}
	noToken(t, err)
	if _, err := g.PullRequest(context.Background(), 0); err == nil {
		t.Error("PullRequest(0) answered no error")
	}
	if calls := f.calls(); len(calls) != 2 {
		t.Errorf("calls = %v, want the two reads of a number", calls)
	}
}

func TestGitHubPullRequestStates(t *testing.T) {
	merged := githubPullJSON(1, "a", "s", "t")
	merged["state"], merged["merged_at"] = "closed", "2026-10-17T10:00:00Z"
	closed := githubPullJSON(2, "b", "s", "t")
	closed["state"] = "closed"
	for _, tc := range []struct {
		pull map[string]any
		want PullRequestState
	}{{githubPullJSON(3, "c", "s", "t"), PullRequestOpen}, {closed, PullRequestClosed}, {merged, PullRequestMerged}} {
		b, _ := json.Marshal(tc.pull)
		var p githubPull
		if err := json.Unmarshal(b, &p); err != nil {
			t.Fatal(err)
		}
		if got := p.pullRequest().State; got != tc.want {
			t.Errorf("pull request #%v: State = %q, want %q", tc.pull["number"], got, tc.want)
		}
	}
}

func TestGitHubOpenPullRequestOpensOneWhenNoneIsOpen(t *testing.T) {
	f := newFakeGitHub(t)
	f.reply("GET /repos/octo/repo/pulls", http.StatusOK, []any{})
	f.reply("POST /repos/octo/repo/pulls", http.StatusCreated, githubPullJSON(8, "story/S-0001", "def456", "S-0001 A story"))
	pr, err := f.forge().OpenPullRequest(context.Background(), PullRequestInput{Head: "story/S-0001", Base: "main", Title: "S-0001 A story", Body: "Closes S-0001"})
	if err != nil {
		t.Fatal(err)
	}
	if pr.Number != 8 || pr.HeadSHA != "def456" {
		t.Errorf("OpenPullRequest = %+v, want #8 at def456", pr)
	}
	if got, want := f.calls(), []string{"GET /repos/octo/repo/pulls", "POST /repos/octo/repo/pulls"}; !slices.Equal(got, want) {
		t.Fatalf("calls = %v, want %v", got, want)
	}
	body := f.requests()[1].Body
	for k, want := range map[string]string{"head": "story/S-0001", "base": "main", "title": "S-0001 A story", "body": "Closes S-0001"} {
		if body[k] != want {
			t.Errorf("POST body %s = %v, want %q", k, body[k], want)
		}
	}
}

func TestGitHubOpenPullRequestUpdatesTheOpenOne(t *testing.T) {
	f := newFakeGitHub(t)
	f.reply("GET /repos/octo/repo/pulls", http.StatusOK, []any{githubPullJSON(7, "story/S-0001", "abc123", "Old title")})
	f.reply("PATCH /repos/octo/repo/pulls/7", http.StatusOK, githubPullJSON(7, "story/S-0001", "abc123", "New title"))
	pr, err := f.forge().OpenPullRequest(context.Background(), PullRequestInput{Head: "story/S-0001", Base: "main", Title: "New title", Body: "New body"})
	if err != nil {
		t.Fatal(err)
	}
	if pr.Number != 7 || pr.Title != "New title" {
		t.Errorf("OpenPullRequest = %+v, want #7 titled New title", pr)
	}
	if got, want := f.calls(), []string{"GET /repos/octo/repo/pulls", "PATCH /repos/octo/repo/pulls/7"}; !slices.Equal(got, want) {
		t.Fatalf("calls = %v, want %v", got, want)
	}
	body := f.requests()[1].Body
	for k, want := range map[string]string{"base": "main", "title": "New title", "body": "New body"} {
		if body[k] != want {
			t.Errorf("PATCH body %s = %v, want %q", k, body[k], want)
		}
	}
}

func TestGitHubOpenPullRequestRefusesAMissingField(t *testing.T) {
	f := newFakeGitHub(t)
	if _, err := f.forge().OpenPullRequest(context.Background(), PullRequestInput{Head: "story/S-0001", Base: "main"}); err == nil {
		t.Error("OpenPullRequest without a title answered no error")
	}
	if calls := f.calls(); len(calls) != 0 {
		t.Errorf("calls = %v, want none", calls)
	}
}

func TestGitHubOpenPullRequestCarriesGitHubsValidationErrors(t *testing.T) {
	f := newFakeGitHub(t)
	f.reply("GET /repos/octo/repo/pulls", http.StatusOK, []any{})
	f.reply("POST /repos/octo/repo/pulls", http.StatusUnprocessableEntity, map[string]any{
		"message": "Validation Failed",
		"errors":  []any{map[string]string{"resource": "PullRequest", "code": "custom", "message": "No commits between main and story/S-0001"}},
	})
	_, err := f.forge().OpenPullRequest(context.Background(), PullRequestInput{Head: "story/S-0001", Base: "main", Title: "t"})
	if err == nil {
		t.Fatal("OpenPullRequest answered no error to a 422")
	}
	for _, want := range []string{"POST /repos/octo/repo/pulls", "422", "Validation Failed", "No commits between main and story/S-0001"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not name %q", err, want)
		}
	}
	noToken(t, err)
}

// checksPR is the pull request the checks tests read.
var checksPR = PullRequest{Number: 7, Head: "story/S-0001", Base: "main", HeadSHA: "abc123"}

func TestGitHubChecksMergesRunsStatusesAndRequiredChecks(t *testing.T) {
	f := newFakeGitHub(t)
	f.reply("GET /repos/octo/repo/rules/branches/main", http.StatusOK, []any{
		map[string]any{"type": "pull_request", "parameters": map[string]any{"required_approving_review_count": 1}},
		map[string]any{"type": "required_status_checks", "parameters": map[string]any{
			"required_status_checks": []any{map[string]any{"context": "build"}, map[string]any{"context": "never-ran"}},
		}},
	})
	f.reply("GET /repos/octo/repo/branches/main/protection/required_status_checks", http.StatusOK, map[string]any{
		"contexts": []string{"ci/legacy"}, "checks": []any{map[string]any{"context": "lint"}},
	})
	f.handle("GET /repos/octo/repo/commits/abc123/check-runs", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("filter") != "latest" {
			t.Errorf("check runs read with filter %q, want latest", r.URL.Query().Get("filter"))
		}
		writeJSON(w, http.StatusOK, map[string]any{"total_count": 4, "check_runs": []any{
			map[string]any{"name": "build", "status": "completed", "conclusion": "success", "html_url": "https://github.com/octo/repo/runs/1", "output": map[string]any{"title": "Build passed"}},
			map[string]any{"name": "lint", "status": "completed", "conclusion": "failure", "html_url": "https://github.com/octo/repo/runs/2", "output": map[string]any{"title": "2 errors"}},
			map[string]any{"name": "docs", "status": "completed", "conclusion": "skipped", "details_url": "https://ci.example/docs"},
			map[string]any{"name": "e2e", "status": "in_progress", "conclusion": nil},
			map[string]any{"name": "flaky", "status": "completed", "conclusion": "timed_out"},
		}})
	})
	f.reply("GET /repos/octo/repo/commits/abc123/status", http.StatusOK, map[string]any{"state": "failure", "statuses": []any{
		map[string]any{"context": "ci/legacy", "state": "success", "description": "All good", "target_url": "https://ci.example/1"},
		map[string]any{"context": "coverage", "state": "error", "description": "Coverage fell"},
		map[string]any{"context": "deploy", "state": "pending"},
	}})
	got, err := f.forge().Checks(context.Background(), checksPR)
	if err != nil {
		t.Fatal(err)
	}
	want := []Check{
		{Name: "build", State: CheckPassed, Required: true, URL: "https://github.com/octo/repo/runs/1", Summary: "Build passed"},
		{Name: "ci/legacy", State: CheckPassed, Required: true, URL: "https://ci.example/1", Summary: "All good"},
		{Name: "coverage", State: CheckFailed, Summary: "Coverage fell"},
		{Name: "deploy", State: CheckPending},
		{Name: "docs", State: CheckPassed, URL: "https://ci.example/docs", Summary: "skipped"},
		{Name: "e2e", State: CheckPending, Summary: "in progress"},
		{Name: "flaky", State: CheckFailed, Summary: "timed out"},
		{Name: "lint", State: CheckFailed, Required: true, URL: "https://github.com/octo/repo/runs/2", Summary: "2 errors"},
		{Name: "never-ran", State: CheckPending, Required: true, Summary: "has not reported yet"},
	}
	if !slices.Equal(got, want) {
		t.Errorf("Checks =\n%+v\nwant\n%+v", got, want)
	}
}

func TestGitHubChecksTreatProtectionOutOfReachAsNone(t *testing.T) {
	for _, status := range []int{http.StatusForbidden, http.StatusNotFound} {
		f := newFakeGitHub(t)
		f.reply("GET /repos/octo/repo/rules/branches/main", http.StatusOK, []any{})
		f.reply("GET /repos/octo/repo/branches/main/protection/required_status_checks", status, map[string]string{"message": "Must have admin rights to Repository."})
		f.reply("GET /repos/octo/repo/commits/abc123/check-runs", http.StatusOK, map[string]any{"check_runs": []any{
			map[string]any{"name": "build", "status": "completed", "conclusion": "success"},
		}})
		f.reply("GET /repos/octo/repo/commits/abc123/status", http.StatusOK, map[string]any{"statuses": []any{}})
		got, err := f.forge().Checks(context.Background(), checksPR)
		if err != nil {
			t.Fatalf("protection answering %d: %v", status, err)
		}
		if want := []Check{{Name: "build", State: CheckPassed, Summary: "success"}}; !slices.Equal(got, want) {
			t.Errorf("protection answering %d: Checks = %+v, want %+v", status, got, want)
		}
	}
}

func TestGitHubChecksFailWhenRulesCannotBeRead(t *testing.T) {
	f := newFakeGitHub(t)
	f.reply("GET /repos/octo/repo/rules/branches/main", http.StatusInternalServerError, map[string]string{"message": "Server Error"})
	_, err := f.forge().Checks(context.Background(), checksPR)
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Errorf("Checks with rules answering 500 = %v, want an error naming 500", err)
	}
}

func TestGitHubChecksNeedTheHeadCommit(t *testing.T) {
	f := newFakeGitHub(t)
	if _, err := f.forge().Checks(context.Background(), PullRequest{Number: 7, Base: "main"}); err == nil {
		t.Error("Checks without a head commit answered no error")
	}
	if calls := f.calls(); len(calls) != 0 {
		t.Errorf("calls = %v, want none", calls)
	}
}

func TestGitHubReviews(t *testing.T) {
	f := newFakeGitHub(t)
	f.reply("GET /repos/octo/repo/pulls/7/reviews", http.StatusOK, []any{
		map[string]any{"user": map[string]any{"login": "ann"}, "state": "APPROVED", "body": "", "html_url": "https://github.com/octo/repo/pull/7#r1", "submitted_at": "2026-10-17T10:00:00Z"},
		map[string]any{"user": map[string]any{"login": "bob"}, "state": "CHANGES_REQUESTED", "body": "Rename it", "html_url": "u2", "submitted_at": "2026-10-17T11:00:00Z"},
		map[string]any{"user": map[string]any{"login": "cat"}, "state": "COMMENTED", "body": "Nice", "html_url": "u3", "submitted_at": "2026-10-17T12:00:00Z"},
		map[string]any{"user": nil, "state": "DISMISSED", "html_url": "u4", "submitted_at": "2026-10-17T13:00:00Z"},
		map[string]any{"user": map[string]any{"login": "dan"}, "state": "PENDING", "html_url": "u5"},
	})
	got, err := f.forge().Reviews(context.Background(), checksPR)
	if err != nil {
		t.Fatal(err)
	}
	want := []Review{
		{Author: "ann", State: ReviewApproved, URL: "https://github.com/octo/repo/pull/7#r1", SubmittedAt: "2026-10-17T10:00:00Z"},
		{Author: "bob", State: ReviewChangesRequested, Body: "Rename it", URL: "u2", SubmittedAt: "2026-10-17T11:00:00Z"},
		{Author: "cat", State: ReviewCommented, Body: "Nice", URL: "u3", SubmittedAt: "2026-10-17T12:00:00Z"},
		{State: ReviewDismissed, URL: "u4", SubmittedAt: "2026-10-17T13:00:00Z"},
		{Author: "dan", State: ReviewPending, URL: "u5"},
	}
	if !slices.Equal(got, want) {
		t.Errorf("Reviews =\n%+v\nwant\n%+v", got, want)
	}
}

func TestGitHubReviewsRefuseAStateFlaiDoesNotKnow(t *testing.T) {
	f := newFakeGitHub(t)
	f.reply("GET /repos/octo/repo/pulls/7/reviews", http.StatusOK, []any{map[string]any{"state": "BLESSED", "html_url": "u1"}})
	_, err := f.forge().Reviews(context.Background(), checksPR)
	if err == nil || !strings.Contains(err.Error(), "BLESSED") {
		t.Errorf("Reviews with an unknown state = %v, want an error naming it", err)
	}
}

// The calls RequiredApprovals makes for main: its rules, then its classic protection.
const (
	rulesPath   = "GET /repos/octo/repo/rules/branches/main"
	reviewsPath = "GET /repos/octo/repo/branches/main/protection/required_pull_request_reviews"
)

func TestGitHubRequiredApprovals(t *testing.T) {
	adminOnly := map[string]string{"message": "Must have admin rights to Repository."}
	rules := func(counts ...int) []any {
		out := []any{map[string]any{"type": "required_status_checks", "parameters": map[string]any{"required_status_checks": []any{}}}}
		for _, n := range counts {
			out = append(out, map[string]any{"type": "pull_request", "parameters": map[string]any{"required_approving_review_count": n}})
		}
		return out
	}
	for _, tc := range []struct {
		name          string
		rulesStatus   int
		rules         any
		reviewsStatus int
		reviews       any
		want          int
	}{
		{"rulesets alone", http.StatusOK, rules(1, 2), http.StatusNotFound, map[string]string{"message": "Branch not protected"}, 2},
		{"classic protection alone", http.StatusOK, []any{}, http.StatusOK, map[string]any{"required_approving_review_count": 3}, 3},
		{"both, the higher of them", http.StatusOK, rules(2), http.StatusOK, map[string]any{"required_approving_review_count": 1}, 2},
		{"both out of reach, 403", http.StatusForbidden, adminOnly, http.StatusForbidden, adminOnly, 0},
		{"both out of reach, 404", http.StatusNotFound, map[string]string{"message": "Not Found"}, http.StatusNotFound, map[string]string{"message": "Not Found"}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newFakeGitHub(t)
			f.reply(rulesPath, tc.rulesStatus, tc.rules)
			f.reply(reviewsPath, tc.reviewsStatus, tc.reviews)
			got, err := f.forge().RequiredApprovals(context.Background(), "main")
			if err != nil {
				t.Fatal(err)
			}
			if got != tc.want {
				t.Errorf("RequiredApprovals = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestGitHubRequiredApprovalsFailWhenGitHubFails(t *testing.T) {
	serverError := map[string]string{"message": "Server Error"}
	for _, path := range []string{rulesPath, reviewsPath} {
		f := newFakeGitHub(t)
		if path == rulesPath {
			f.reply(rulesPath, http.StatusInternalServerError, serverError)
			f.reply(reviewsPath, http.StatusOK, map[string]any{"required_approving_review_count": 1})
		} else {
			f.reply(rulesPath, http.StatusOK, []any{})
			f.reply(reviewsPath, http.StatusInternalServerError, serverError)
		}
		_, err := f.forge().RequiredApprovals(context.Background(), "main")
		if err == nil {
			t.Errorf("RequiredApprovals with %s answering 500 answered no error", path)
			continue
		}
		for _, want := range []string{"500", "main", "octo/repo"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("RequiredApprovals with %s answering 500: error %q does not name %q", path, err, want)
			}
		}
		noToken(t, err)
	}
}

func TestGitHubRequiredApprovalsNeedABase(t *testing.T) {
	f := newFakeGitHub(t)
	if _, err := f.forge().RequiredApprovals(context.Background(), ""); err == nil {
		t.Error("RequiredApprovals without a base answered no error")
	}
	if calls := f.calls(); len(calls) != 0 {
		t.Errorf("calls = %v, want none", calls)
	}
}

func TestGitHubFollowsTheNextPage(t *testing.T) {
	f := newFakeGitHub(t)
	f.handle("GET /repos/octo/repo/pulls/7/reviews", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("per_page") != "100" {
			t.Errorf("reviews read with per_page %q, want 100", r.URL.Query().Get("per_page"))
		}
		w.Header().Set("Link", `<`+f.server.URL+`/repositories/1/pulls/7/reviews?per_page=100&page=2>; rel="next", <`+f.server.URL+`/repositories/1/pulls/7/reviews?per_page=100&page=2>; rel="last"`)
		writeJSON(w, http.StatusOK, []any{map[string]any{"user": map[string]any{"login": "ann"}, "state": "APPROVED"}})
	})
	f.handle("GET /repositories/1/pulls/7/reviews", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Errorf("second page read with page %q, want 2", r.URL.Query().Get("page"))
		}
		writeJSON(w, http.StatusOK, []any{map[string]any{"user": map[string]any{"login": "bob"}, "state": "COMMENTED"}})
	})
	got, err := f.forge().Reviews(context.Background(), checksPR)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0].Author != "ann" || got[1].Author != "bob" {
		t.Errorf("Reviews = %+v, want ann's then bob's", got)
	}
}

func TestGitHubRefusesANextPageOffItsAPI(t *testing.T) {
	f := newFakeGitHub(t)
	f.handle("GET /repos/octo/repo/pulls/7/reviews", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Link", `<https://elsewhere.example/steal?page=2>; rel="next"`)
		writeJSON(w, http.StatusOK, []any{})
	})
	_, err := f.forge().Reviews(context.Background(), checksPR)
	if err == nil || !strings.Contains(err.Error(), "elsewhere.example") {
		t.Errorf("Reviews with a next page off the API = %v, want an error naming it", err)
	}
}

func TestGitHubStopsAfterMaxPages(t *testing.T) {
	f := newFakeGitHub(t)
	f.handle("GET /repos/octo/repo/pulls/7/reviews", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Link", `<`+f.server.URL+r.URL.Path+`>; rel="next"`)
		writeJSON(w, http.StatusOK, []any{})
	})
	_, err := f.forge().Reviews(context.Background(), checksPR)
	if err == nil || !strings.Contains(err.Error(), "pages") {
		t.Errorf("Reviews on an endless list = %v, want an error", err)
	}
	if n := len(f.calls()); n != maxPages {
		t.Errorf("calls = %d, want %d", n, maxPages)
	}
}

func TestGitHubMerge(t *testing.T) {
	for _, method := range MergeMethods() {
		t.Run(string(method), func(t *testing.T) {
			f := newFakeGitHub(t)
			f.reply("PUT /repos/octo/repo/pulls/7/merge", http.StatusOK, map[string]any{"sha": "merged789", "merged": true, "message": "Pull Request successfully merged"})
			got, err := f.forge().Merge(context.Background(), checksPR, method)
			if err != nil {
				t.Fatal(err)
			}
			if got.SHA != "merged789" {
				t.Errorf("Merge SHA = %q, want merged789", got.SHA)
			}
			body := f.requests()[0].Body
			if body["merge_method"] != string(method) || body["sha"] != "abc123" {
				t.Errorf("PUT body = %v, want merge_method %s and sha abc123", body, method)
			}
		})
	}
}

func TestGitHubMergeRefused(t *testing.T) {
	for _, tc := range []struct {
		status  int
		message string
	}{
		{http.StatusMethodNotAllowed, "Required status check \"build\" is failing."},
		{http.StatusConflict, "Head branch was modified. Review and try the merge again."},
	} {
		f := newFakeGitHub(t)
		f.reply("PUT /repos/octo/repo/pulls/7/merge", tc.status, map[string]string{"message": tc.message})
		_, err := f.forge().Merge(context.Background(), checksPR, MergeSquash)
		if !errors.Is(err, ErrMergeRefused) {
			t.Errorf("Merge answered %d: %v, want ErrMergeRefused", tc.status, err)
			continue
		}
		if !strings.Contains(err.Error(), tc.message) {
			t.Errorf("Merge error %q does not carry GitHub's message %q", err, tc.message)
		}
		noToken(t, err)
	}
}

func TestGitHubMergeRefusesAnUnknownMethodWithoutACall(t *testing.T) {
	f := newFakeGitHub(t)
	if _, err := f.forge().Merge(context.Background(), checksPR, MergeMethod("octopus")); err == nil {
		t.Error("Merge by octopus answered no error")
	}
	if _, err := f.forge().Merge(context.Background(), PullRequest{Number: 7}, MergeMerge); err == nil {
		t.Error("Merge without a head commit answered no error")
	}
	if calls := f.calls(); len(calls) != 0 {
		t.Errorf("calls = %v, want none", calls)
	}
}

func TestGitHubErrorsCarryTheContractsErrors(t *testing.T) {
	for _, tc := range []struct {
		status int
		want   error
		hint   string
	}{
		{http.StatusUnauthorized, ErrUnauthorized, "check the token"},
		{http.StatusForbidden, ErrForbidden, "rate limit"},
		{http.StatusNotFound, ErrNotFound, "owner and name"},
	} {
		f := newFakeGitHub(t)
		f.reply("GET /repos/octo/repo", tc.status, map[string]string{"message": "GitHub says no"})
		_, err := f.forge().Access(context.Background())
		if !errors.Is(err, tc.want) {
			t.Errorf("Access answered %d: %v, want %v", tc.status, err, tc.want)
			continue
		}
		for _, want := range []string{"GET /repos/octo/repo", http.StatusText(tc.status), "GitHub says no", tc.hint, "octo/repo"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("Access answered %d: error %q does not name %q", tc.status, err, want)
			}
		}
		noToken(t, err)
	}
}

func TestGitHubTokenErrorStopsTheCall(t *testing.T) {
	f := newFakeGitHub(t)
	missing := errors.New("GITHUB_TOKEN is not set")
	for _, token := range []TokenFunc{
		func() (string, error) { return "", missing },
		func() (string, error) { return "", nil },
	} {
		g, err := NewGitHub(GitHubOptions{BaseURL: f.server.URL, Owner: "octo", Name: "repo", Token: token, Client: f.server.Client()})
		if err != nil {
			t.Fatal(err)
		}
		if _, _, err := g.FindPullRequest(context.Background(), "story/S-0001"); err == nil {
			t.Error("FindPullRequest with no token answered no error")
		}
	}
	g, _ := NewGitHub(GitHubOptions{BaseURL: f.server.URL, Owner: "octo", Name: "repo", Token: func() (string, error) { return "", missing }})
	if _, err := g.Access(context.Background()); !errors.Is(err, missing) {
		t.Errorf("Access = %v, want the token func's error", err)
	}
	if calls := f.calls(); len(calls) != 0 {
		t.Errorf("calls = %v, want none", calls)
	}
}

func TestNewGitHubRefusesAnIncompleteRepository(t *testing.T) {
	token := func() (string, error) { return testToken, nil }
	for _, o := range []GitHubOptions{
		{Name: "repo", Token: token},
		{Owner: "octo", Token: token},
		{Owner: "octo/x", Name: "repo", Token: token},
		{Owner: "octo", Name: "repo"},
		{BaseURL: "ftp://example.com", Owner: "octo", Name: "repo", Token: token},
		{BaseURL: "not a url", Owner: "octo", Name: "repo", Token: token},
	} {
		if _, err := NewGitHub(o); err == nil {
			t.Errorf("NewGitHub(%+v) answered no error", GitHubOptions{BaseURL: o.BaseURL, Owner: o.Owner, Name: o.Name})
		}
	}
	g, err := NewGitHub(GitHubOptions{Owner: "octo", Name: "repo", Token: token})
	if err != nil {
		t.Fatal(err)
	}
	if g.base != GitHubAPI {
		t.Errorf("NewGitHub with no BaseURL: base %q, want %s", g.base, GitHubAPI)
	}
}

func TestParseGitHubRemote(t *testing.T) {
	for _, remote := range []string{
		"https://github.com/octo/repo",
		"https://github.com/octo/repo.git",
		"https://github.com/octo/repo/",
		"https://x-access-token:" + testToken + "@github.com/octo/repo.git",
		"http://github.com/octo/repo",
		"git@github.com:octo/repo.git",
		"git@github.com:octo/repo",
		"git@GitHub.com:octo/repo.git",
		"ssh://git@github.com/octo/repo.git",
		"ssh://git@github.com:22/octo/repo",
		"git://github.com/octo/repo.git",
	} {
		owner, name, err := ParseGitHubRemote(remote)
		if err != nil || owner != "octo" || name != "repo" {
			t.Errorf("ParseGitHubRemote(%s) = %q, %q, %v; want octo, repo", strings.ReplaceAll(remote, testToken, "<token>"), owner, name, err)
		}
	}
	for _, tc := range []struct {
		remote, want string
	}{
		{"https://gitlab.com/octo/repo.git", "gitlab.com"},
		{"git@bitbucket.org:octo/repo.git", "bitbucket.org"},
		{"https://github.com/octo", "owner/name"},
		{"https://github.com/octo/repo/extra", "owner/name"},
		{"https://github.com/", "owner/name"},
		{"file:///srv/repo.git", "https, ssh, or git"},
		{"/srv/repo.git", "host:path"},
		{"https://" + testToken + "@gitlab.com/octo/repo.git", "gitlab.com"},
		{"https://x:" + testToken + "@github.com/octo", "owner/name"},
	} {
		_, _, err := ParseGitHubRemote(tc.remote)
		shown := strings.ReplaceAll(tc.remote, testToken, "<token>")
		if err == nil {
			t.Errorf("ParseGitHubRemote(%s) answered no error", shown)
			continue
		}
		if !strings.Contains(err.Error(), tc.want) {
			t.Errorf("ParseGitHubRemote(%s) error %q does not name %q", shown, err, tc.want)
		}
		noToken(t, err)
	}
}
