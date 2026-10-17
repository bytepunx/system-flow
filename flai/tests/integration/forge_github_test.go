package integration

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/forge"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// TestGitHubForgeAgainstARealRepository drives the GitHub adapter against a
// repository on github.com (S-0392). It only reads: it opens, updates, and
// merges nothing, so anyone may run it with their own token.
//
//	FLAI_FORGE_TEST_REPO=owner/name \
//	GITHUB_TOKEN=... \
//	FLAI_FORGE_TEST_HEAD=branch-with-an-open-pull-request \
//	go test ./tests/integration/ -run TestGitHubForge -v
//
// The token is read the way flai reads it, from the variable the host names
// for github with no forges entry, GITHUB_TOKEN. FLAI_FORGE_TEST_HEAD is
// optional: without it the pull request's checks and reviews are not read.
// FLAI_FORGE_TEST_API is optional too, a GitHub Enterprise API address.
// Without the repository or the token it skips, naming what is missing, so
// the ordinary tiers are unaffected.
func TestGitHubForgeAgainstARealRepository(t *testing.T) {
	tokenEnv := config.Config{}.ForgeTokenEnv(manifest.ForgeGitHub)
	token := config.Config{}.ForgeTokenFunc(manifest.ForgeGitHub, os.Getenv)
	repo := strings.TrimSpace(os.Getenv("FLAI_FORGE_TEST_REPO"))
	var missing []string
	if repo == "" {
		missing = append(missing, "FLAI_FORGE_TEST_REPO")
	}
	if _, err := token(); err != nil {
		missing = append(missing, tokenEnv)
	}
	if len(missing) > 0 {
		t.Skipf("set FLAI_FORGE_TEST_REPO (owner/name) and %s to test the GitHub forge against a real repository; missing: %s",
			tokenEnv, strings.Join(missing, ", "))
	}
	owner, name, ok := strings.Cut(repo, "/")
	if !ok || owner == "" || name == "" || strings.Contains(name, "/") {
		t.Fatalf("FLAI_FORGE_TEST_REPO is %q, not owner/name; give it as GitHub names the repository, such as bytepunx/system-flow", repo)
	}
	gh, err := forge.NewGitHub(forge.GitHubOptions{
		BaseURL: os.Getenv("FLAI_FORGE_TEST_API"),
		Owner:   owner,
		Name:    name,
		Token:   token,
	})
	if err != nil {
		t.Fatalf("make the GitHub forge for %s: %v", repo, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	access, err := gh.Access(ctx)
	if err != nil {
		t.Fatalf("Access to %s: %v", gh.Repository(), err)
	}
	if !access.Read {
		t.Errorf("Access to %s: Read is false, want true", gh.Repository())
	}
	t.Logf("access to %s: read %t, push %t", gh.Repository(), access.Read, access.Push)

	absent := fmt.Sprintf("flai-forge-test-absent-%d", time.Now().UnixNano())
	if pr, found, err := gh.FindPullRequest(ctx, absent); err != nil {
		t.Errorf("FindPullRequest(%s): %v", absent, err)
	} else if found {
		t.Errorf("FindPullRequest(%s) found #%d, want none", absent, pr.Number)
	}

	head := strings.TrimSpace(os.Getenv("FLAI_FORGE_TEST_HEAD"))
	if head == "" {
		t.Logf("FLAI_FORGE_TEST_HEAD is not set; set it to a branch of %s with an open pull request to read its checks and reviews", gh.Repository())
		return
	}
	pr, found, err := gh.FindPullRequest(ctx, head)
	if err != nil {
		t.Fatalf("FindPullRequest(%s): %v", head, err)
	}
	if !found {
		t.Fatalf("FindPullRequest(%s) found no open pull request in %s; FLAI_FORGE_TEST_HEAD must name a branch of %s with one open", head, gh.Repository(), owner)
	}
	if pr.Number <= 0 || pr.Head != head || pr.Base == "" || pr.HeadSHA == "" || pr.State != forge.PullRequestOpen {
		t.Errorf("FindPullRequest(%s) = #%d head %q base %q sha %q state %q; want a number, head %q, a base, a sha, and open",
			head, pr.Number, pr.Head, pr.Base, pr.HeadSHA, pr.State, head)
	}
	t.Logf("pull request #%d from %s into %s at %s, draft %t", pr.Number, pr.Head, pr.Base, pr.HeadSHA, pr.Draft)

	checks, err := gh.Checks(ctx, pr)
	if err != nil {
		t.Errorf("Checks on #%d: %v", pr.Number, err)
	}
	for _, c := range checks {
		t.Logf("check %s: %s, required %t", c.Name, c.State, c.Required)
	}
	t.Logf("%d checks", len(checks))

	reviews, err := gh.Reviews(ctx, pr)
	if err != nil {
		t.Errorf("Reviews on #%d: %v", pr.Number, err)
	}
	for _, r := range reviews {
		t.Logf("review by %s: %s", r.Author, r.State)
	}
	t.Logf("%d reviews", len(reviews))
}
