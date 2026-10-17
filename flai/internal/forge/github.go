package forge

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"
)

// GitHubAPI is the address of github.com's REST API.
const GitHubAPI = "https://api.github.com"

// githubAPIVersion is the REST API version the adapter is written to.
const githubAPIVersion = "2022-11-28"

// githubTimeout bounds one call when the caller gives no client of its own.
const githubTimeout = 30 * time.Second

// maxAnswer bounds how much of an answer is decoded: a page of 100 check
// runs, each with an output summary of up to 64 KiB, fits.
const maxAnswer = 16 << 20

// maxErrorAnswer bounds how much of a refusal is read for its message.
const maxErrorAnswer = 64 << 10

// perPage is how many items a list asks for at once, GitHub's most.
const perPage = 100

// maxPages bounds how many pages one list follows, 5000 items at perPage.
const maxPages = 50

// GitHubOptions is what NewGitHub makes a GitHub from.
type GitHubOptions struct {
	// BaseURL is the REST API's address; empty is GitHubAPI.
	BaseURL string
	// Owner is the account or organization that owns the repository.
	Owner string
	// Name is the repository's name.
	Name string
	// Token answers the token for each call.
	Token TokenFunc
	// Client sends the calls; nil is a client with a 30 s timeout.
	Client *http.Client
}

// GitHub is the Forge for one repository on GitHub, over its REST API.
type GitHub struct {
	base   string
	owner  string
	name   string
	token  TokenFunc
	client *http.Client
}

var _ Forge = (*GitHub)(nil)

// NewGitHub is the GitHub forge for the repository o names.
func NewGitHub(o GitHubOptions) (*GitHub, error) {
	base := strings.TrimRight(cmp.Or(o.BaseURL, GitHubAPI), "/")
	u, err := url.Parse(base)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" {
		return nil, fmt.Errorf("GitHub's API address %q is not an http or https URL; give one such as %s", base, GitHubAPI)
	}
	if o.Owner == "" || o.Name == "" || strings.Contains(o.Owner, "/") || strings.Contains(o.Name, "/") {
		return nil, fmt.Errorf("the GitHub repository %q/%q is not an owner and a name; give both, as the origin remote names them", o.Owner, o.Name)
	}
	if o.Token == nil {
		return nil, fmt.Errorf("the GitHub forge for %s/%s has no way to read its token; give it a TokenFunc", o.Owner, o.Name)
	}
	client := o.Client
	if client == nil {
		client = &http.Client{Timeout: githubTimeout}
	}
	return &GitHub{base: base, owner: o.Owner, name: o.Name, token: o.Token, client: client}, nil
}

// Name is "github".
func (g *GitHub) Name() string { return "github" }

// Repository is the repository as owner/name.
func (g *GitHub) Repository() string { return g.owner + "/" + g.name }

// Access reads the repository and answers whether the token may push to
// it: reading it at all is read access.
func (g *GitHub) Access(ctx context.Context) (Access, error) {
	var repo struct {
		Permissions *struct {
			Push     bool `json:"push"`
			Maintain bool `json:"maintain"`
			Admin    bool `json:"admin"`
		} `json:"permissions"`
	}
	if _, err := g.call(ctx, http.MethodGet, g.repoPath(""), nil, &repo); err != nil {
		return Access{}, fmt.Errorf("read the repository %s on GitHub: %w", g.Repository(), err)
	}
	access := Access{Read: true}
	if p := repo.Permissions; p != nil {
		access.Push = p.Push || p.Maintain || p.Admin
	}
	return access, nil
}

// FindPullRequest is the open pull request from head, the first GitHub
// lists when there are several into different bases.
func (g *GitHub) FindPullRequest(ctx context.Context, head string) (PullRequest, bool, error) {
	if head == "" {
		return PullRequest{}, false, errors.New("find a pull request on GitHub: no head branch given; name the story's branch")
	}
	q := url.Values{"head": {g.owner + ":" + head}, "state": {"open"}, "per_page": {strconv.Itoa(perPage)}}
	var pulls []githubPull
	if _, err := g.call(ctx, http.MethodGet, g.repoPath("/pulls?"+q.Encode()), nil, &pulls); err != nil {
		return PullRequest{}, false, fmt.Errorf("find the open pull request from %s in %s on GitHub: %w", head, g.Repository(), err)
	}
	if len(pulls) == 0 {
		return PullRequest{}, false, nil
	}
	return pulls[0].pullRequest(), true, nil
}

// OpenPullRequest opens a pull request from in.Head into in.Base, or
// updates the title, body, and base of the one open from in.Head.
func (g *GitHub) OpenPullRequest(ctx context.Context, in PullRequestInput) (PullRequest, error) {
	if in.Head == "" || in.Base == "" || in.Title == "" {
		return PullRequest{}, fmt.Errorf("open a pull request in %s on GitHub: head %q, base %q, and title %q must all be given", g.Repository(), in.Head, in.Base, in.Title)
	}
	open, found, err := g.FindPullRequest(ctx, in.Head)
	if err != nil {
		return PullRequest{}, err
	}
	var pull githubPull
	if found {
		edit := map[string]string{"title": in.Title, "body": in.Body, "base": in.Base}
		if _, err := g.call(ctx, http.MethodPatch, g.repoPath("/pulls/"+strconv.Itoa(open.Number)), edit, &pull); err != nil {
			return PullRequest{}, fmt.Errorf("update pull request #%d from %s in %s on GitHub: %w", open.Number, in.Head, g.Repository(), err)
		}
		return pull.pullRequest(), nil
	}
	create := map[string]string{"title": in.Title, "head": in.Head, "base": in.Base, "body": in.Body}
	if _, err := g.call(ctx, http.MethodPost, g.repoPath("/pulls"), create, &pull); err != nil {
		return PullRequest{}, fmt.Errorf("open a pull request from %s into %s in %s on GitHub: %w", in.Head, in.Base, g.Repository(), err)
	}
	return pull.pullRequest(), nil
}

// Checks are the latest check run of each name and the latest commit
// status of each context on the pull request's head commit, with a pending
// check for each the base branch requires that has reported neither, in
// order of name. A check is required when the base branch's rulesets or
// its classic protection name it.
func (g *GitHub) Checks(ctx context.Context, pr PullRequest) ([]Check, error) {
	if pr.HeadSHA == "" || pr.Base == "" {
		return nil, fmt.Errorf("read the checks on pull request #%d in %s on GitHub: it has no head commit or base branch; read it with FindPullRequest first", pr.Number, g.Repository())
	}
	required, err := g.requiredChecks(ctx, pr.Base)
	if err != nil {
		return nil, fmt.Errorf("read the checks pull request #%d's base %s requires in %s on GitHub: %w", pr.Number, pr.Base, g.Repository(), err)
	}
	var checks []Check
	runs := g.repoPath("/commits/" + url.PathEscape(pr.HeadSHA) + "/check-runs?filter=latest&per_page=" + strconv.Itoa(perPage))
	err = pages(ctx, g, runs, func(page githubCheckRuns) error {
		for _, run := range page.CheckRuns {
			checks = append(checks, run.check())
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read the check runs on pull request #%d's head %s in %s on GitHub: %w", pr.Number, pr.HeadSHA, g.Repository(), err)
	}
	statuses := g.repoPath("/commits/" + url.PathEscape(pr.HeadSHA) + "/status?per_page=" + strconv.Itoa(perPage))
	err = pages(ctx, g, statuses, func(page githubCombinedStatus) error {
		for _, status := range page.Statuses {
			checks = append(checks, status.check())
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read the commit statuses on pull request #%d's head %s in %s on GitHub: %w", pr.Number, pr.HeadSHA, g.Repository(), err)
	}
	reported := map[string]bool{}
	for i := range checks {
		checks[i].Required = required[checks[i].Name]
		reported[checks[i].Name] = true
	}
	for name := range required {
		if !reported[name] {
			checks = append(checks, Check{Name: name, State: CheckPending, Required: true, Summary: "has not reported yet"})
		}
	}
	slices.SortStableFunc(checks, func(a, b Check) int { return strings.Compare(a.Name, b.Name) })
	return checks, nil
}

// Reviews are the reviews on the pull request, oldest first.
func (g *GitHub) Reviews(ctx context.Context, pr PullRequest) ([]Review, error) {
	var reviews []Review
	path := g.repoPath("/pulls/" + strconv.Itoa(pr.Number) + "/reviews?per_page=" + strconv.Itoa(perPage))
	err := pages(ctx, g, path, func(page []githubReview) error {
		for _, r := range page {
			review, err := r.review()
			if err != nil {
				return err
			}
			reviews = append(reviews, review)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read the reviews on pull request #%d in %s on GitHub: %w", pr.Number, g.Repository(), err)
	}
	return reviews, nil
}

// RequiredApprovals is the highest count of approving reviews that base's
// rulesets and its classic protection require. Either may be out of reach
// without failing the read, as for requiredChecks.
func (g *GitHub) RequiredApprovals(ctx context.Context, base string) (int, error) {
	if base == "" {
		return 0, fmt.Errorf("read the approvals a base branch requires in %s on GitHub: no base branch given; name the pull request's base", g.Repository())
	}
	most := 0
	rules := g.repoPath("/rules/branches/" + url.PathEscape(base) + "?per_page=" + strconv.Itoa(perPage))
	err := pages(ctx, g, rules, func(page []githubRule) error {
		for _, rule := range page {
			if rule.Type == "pull_request" && rule.Parameters != nil {
				most = max(most, rule.Parameters.RequiredApprovingReviewCount)
			}
		}
		return nil
	})
	if err != nil && !outOfReach(err) {
		return 0, fmt.Errorf("read the rules on %s in %s on GitHub for the approvals it requires: %w", base, g.Repository(), err)
	}
	var reviews struct {
		RequiredApprovingReviewCount int `json:"required_approving_review_count"`
	}
	_, err = g.call(ctx, http.MethodGet, g.repoPath("/branches/"+url.PathEscape(base)+"/protection/required_pull_request_reviews"), nil, &reviews)
	if err != nil && !outOfReach(err) {
		return 0, fmt.Errorf("read the protection of %s in %s on GitHub for the approvals it requires: %w", base, g.Repository(), err)
	}
	return max(most, reviews.RequiredApprovingReviewCount), nil
}

// Merge merges the pull request by method, provided its head is still
// pr.HeadSHA; GitHub's refusal (405, or 409 for a moved head) is
// ErrMergeRefused.
func (g *GitHub) Merge(ctx context.Context, pr PullRequest, method MergeMethod) (Merged, error) {
	if !method.Valid() {
		_, err := ParseMergeMethod(string(method))
		return Merged{}, fmt.Errorf("merge pull request #%d in %s on GitHub: %w", pr.Number, g.Repository(), err)
	}
	if pr.HeadSHA == "" {
		return Merged{}, fmt.Errorf("merge pull request #%d in %s on GitHub: it has no head commit to merge at; read it with FindPullRequest first", pr.Number, g.Repository())
	}
	in := map[string]string{"merge_method": string(method), "sha": pr.HeadSHA}
	var out struct {
		SHA     string `json:"sha"`
		Merged  bool   `json:"merged"`
		Message string `json:"message"`
	}
	_, err := g.call(ctx, http.MethodPut, g.repoPath("/pulls/"+strconv.Itoa(pr.Number)+"/merge"), in, &out)
	var refused *apiError
	if errors.As(err, &refused) && (refused.status == http.StatusMethodNotAllowed || refused.status == http.StatusConflict) {
		return Merged{}, fmt.Errorf("merge pull request #%d by %s in %s on GitHub: %w: %w", pr.Number, method, g.Repository(), ErrMergeRefused, err)
	}
	if err != nil {
		return Merged{}, fmt.Errorf("merge pull request #%d by %s in %s on GitHub: %w", pr.Number, method, g.Repository(), err)
	}
	if !out.Merged {
		return Merged{}, fmt.Errorf("merge pull request #%d by %s in %s on GitHub: it answered %q without merging: %w", pr.Number, method, g.Repository(), out.Message, ErrMergeRefused)
	}
	return Merged{SHA: out.SHA}, nil
}

// requiredChecks are the names of the checks branch requires, from its
// rulesets and its classic protection. Either may be out of reach without
// failing the read: GitHub answers 403 or 404 for protection to a token
// that is not an admin, and for rulesets on a plan without them.
func (g *GitHub) requiredChecks(ctx context.Context, branch string) (map[string]bool, error) {
	required := map[string]bool{}
	rules := g.repoPath("/rules/branches/" + url.PathEscape(branch) + "?per_page=" + strconv.Itoa(perPage))
	err := pages(ctx, g, rules, func(page []githubRule) error {
		for _, rule := range page {
			if rule.Type != "required_status_checks" || rule.Parameters == nil {
				continue
			}
			for _, c := range rule.Parameters.RequiredStatusChecks {
				required[c.Context] = true
			}
		}
		return nil
	})
	if err != nil && !outOfReach(err) {
		return nil, err
	}
	var protection struct {
		Contexts []string `json:"contexts"`
		Checks   []struct {
			Context string `json:"context"`
		} `json:"checks"`
	}
	_, err = g.call(ctx, http.MethodGet, g.repoPath("/branches/"+url.PathEscape(branch)+"/protection/required_status_checks"), nil, &protection)
	if err != nil && !outOfReach(err) {
		return nil, err
	}
	for _, c := range protection.Contexts {
		required[c] = true
	}
	for _, c := range protection.Checks {
		required[c.Context] = true
	}
	return required, nil
}

// outOfReach is whether err is GitHub's 403 or 404.
func outOfReach(err error) bool {
	return errors.Is(err, ErrForbidden) || errors.Is(err, ErrNotFound)
}

// repoPath is rest under the repository's path.
func (g *GitHub) repoPath(rest string) string {
	return "/repos/" + url.PathEscape(g.owner) + "/" + url.PathEscape(g.name) + rest
}

// pages GETs path and each page after it that GitHub's Link header names,
// handing each page to add, up to maxPages; it stops at add's first error.
func pages[T any](ctx context.Context, g *GitHub, path string, add func(T) error) error {
	for page := 0; path != ""; page++ {
		if page == maxPages {
			return fmt.Errorf("GET %s: the list runs past %d pages of %d; flai reads no further", path, maxPages, perPage)
		}
		var v T
		next, err := g.call(ctx, http.MethodGet, path, nil, &v)
		if err != nil {
			return err
		}
		if err := add(v); err != nil {
			return err
		}
		path = next
	}
	return nil
}

// call sends method to path, the API path with its query, with in as a
// JSON body when it is not nil, and decodes a 2xx answer into out. It
// answers the path of the next page the Link header names, or "". The
// token goes in the Authorization header only; no error carries it.
func (g *GitHub) call(ctx context.Context, method, path string, in, out any) (string, error) {
	token, err := g.token()
	if err != nil {
		return "", fmt.Errorf("read the token for GitHub: %w", err)
	}
	if token == "" {
		return "", errors.New("the token for GitHub is empty; set the variable the host names for it")
	}
	var body io.Reader
	if in != nil {
		b, err := json.Marshal(in)
		if err != nil {
			return "", fmt.Errorf("encode the body of %s %s: %w", method, path, err)
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, g.base+path, body)
	if err != nil {
		return "", fmt.Errorf("make the request %s %s: %w; check GitHub's API address", method, path, err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", githubAPIVersion)
	req.Header.Set("User-Agent", "flai")
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s %s: %w; check the network and GitHub's API address", method, path, err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", refusal(method, path, resp)
	}
	if out != nil {
		if err := json.NewDecoder(io.LimitReader(resp.Body, maxAnswer)).Decode(out); err != nil {
			return "", fmt.Errorf("decode the answer to %s %s, read to at most %d MiB: %w", method, path, maxAnswer>>20, err)
		}
	}
	return g.nextPage(resp.Header.Get("Link"))
}

// nextPage is the API path of the rel="next" link in link, or "" when it
// has none. It refuses a link outside the API, which would take the token
// with it.
func (g *GitHub) nextPage(link string) (string, error) {
	for part := range strings.SplitSeq(link, ",") {
		target, params, ok := strings.Cut(part, ";")
		if !ok || !strings.Contains(params, `rel="next"`) {
			continue
		}
		target = strings.TrimSpace(target)
		target = strings.TrimSuffix(strings.TrimPrefix(target, "<"), ">")
		path, under := strings.CutPrefix(target, g.base)
		if !under || !strings.HasPrefix(path, "/") {
			return "", fmt.Errorf("GitHub's next page %s is not under its API address %s; flai sends its token nowhere else", target, g.base)
		}
		return path, nil
	}
	return "", nil
}

// apiError is GitHub's answer outside 2xx to one call.
type apiError struct {
	method  string
	path    string
	status  int
	message string
}

// refusal is the apiError for resp, with GitHub's message and the messages
// of its errors read from its body.
func refusal(method, path string, resp *http.Response) *apiError {
	var body struct {
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Field   string `json:"field"`
			Code    string `json:"code"`
		} `json:"errors"`
	}
	_ = json.NewDecoder(io.LimitReader(resp.Body, maxErrorAnswer)).Decode(&body)
	message := body.Message
	for _, e := range body.Errors {
		detail := cmp.Or(e.Message, strings.TrimSpace(e.Field+" "+e.Code))
		if detail != "" {
			message += "; " + detail
		}
	}
	return &apiError{method: method, path: path, status: resp.StatusCode, message: strings.TrimPrefix(message, "; ")}
}

// Error names the call, the status, GitHub's message, and what to check.
func (e *apiError) Error() string {
	s := fmt.Sprintf("%s %s answered %d %s", e.method, e.path, e.status, http.StatusText(e.status))
	if e.message != "" {
		s += ": " + e.message
	}
	switch e.status {
	case http.StatusUnauthorized:
		s += "; check the token in the variable the host names for GitHub"
	case http.StatusForbidden:
		s += "; check that the token may do this in the repository, or wait for GitHub's rate limit to reset"
	case http.StatusNotFound:
		s += "; check the repository's owner and name, and that the token can read it"
	}
	return s
}

// Unwrap is the contract's error for the status, if it has one.
func (e *apiError) Unwrap() error {
	switch e.status {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	}
	return nil
}

// githubPull is a pull request as GitHub's API answers it.
type githubPull struct {
	Number   int    `json:"number"`
	HTMLURL  string `json:"html_url"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	State    string `json:"state"`
	Draft    bool   `json:"draft"`
	MergedAt string `json:"merged_at"`
	Head     struct {
		Ref string `json:"ref"`
		SHA string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`
}

// pullRequest is p in the contract's terms.
func (p githubPull) pullRequest() PullRequest {
	state := PullRequestState(p.State)
	if p.MergedAt != "" {
		state = PullRequestMerged
	}
	return PullRequest{
		Number: p.Number, URL: p.HTMLURL, Head: p.Head.Ref, Base: p.Base.Ref, HeadSHA: p.Head.SHA,
		Title: p.Title, Body: p.Body, State: state, Draft: p.Draft,
	}
}

// githubCheckRuns is a page of check runs.
type githubCheckRuns struct {
	CheckRuns []githubCheckRun `json:"check_runs"`
}

// githubCheckRun is one check run.
type githubCheckRun struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	HTMLURL    string `json:"html_url"`
	DetailsURL string `json:"details_url"`
	Output     struct {
		Title string `json:"title"`
	} `json:"output"`
}

// check is r in the contract's terms: completed as success, neutral, or
// skipped is passed, completed any other way is failed, and anything not
// completed is pending.
func (r githubCheckRun) check() Check {
	state := CheckPending
	if r.Status == "completed" {
		switch r.Conclusion {
		case "success", "neutral", "skipped":
			state = CheckPassed
		default:
			state = CheckFailed
		}
	}
	summary := r.Output.Title
	if summary == "" {
		summary = strings.ReplaceAll(cmp.Or(r.Conclusion, r.Status), "_", " ")
	}
	return Check{Name: r.Name, State: state, URL: cmp.Or(r.HTMLURL, r.DetailsURL), Summary: summary}
}

// githubCombinedStatus is a page of a commit's statuses, the latest of each context.
type githubCombinedStatus struct {
	Statuses []githubStatus `json:"statuses"`
}

// githubStatus is one commit status.
type githubStatus struct {
	Context     string `json:"context"`
	State       string `json:"state"`
	Description string `json:"description"`
	TargetURL   string `json:"target_url"`
}

// check is s in the contract's terms: success is passed, pending is
// pending, and failure or error is failed.
func (s githubStatus) check() Check {
	state := CheckFailed
	switch s.State {
	case "success":
		state = CheckPassed
	case "pending":
		state = CheckPending
	}
	return Check{Name: s.Context, State: state, URL: s.TargetURL, Summary: s.Description}
}

// githubReview is one review.
type githubReview struct {
	User *struct {
		Login string `json:"login"`
	} `json:"user"`
	State       string `json:"state"`
	Body        string `json:"body"`
	HTMLURL     string `json:"html_url"`
	SubmittedAt string `json:"submitted_at"`
}

// githubReviewStates maps GitHub's review states to the contract's.
var githubReviewStates = map[string]ReviewState{
	"APPROVED":          ReviewApproved,
	"CHANGES_REQUESTED": ReviewChangesRequested,
	"COMMENTED":         ReviewCommented,
	"DISMISSED":         ReviewDismissed,
	"PENDING":           ReviewPending,
}

// review is r in the contract's terms; a state GitHub has added since this
// adapter was written is an error rather than a guess.
func (r githubReview) review() (Review, error) {
	state, ok := githubReviewStates[r.State]
	if !ok {
		return Review{}, fmt.Errorf("a review at %s is in state %q, which flai does not know; upgrade flai", r.HTMLURL, r.State)
	}
	review := Review{State: state, Body: r.Body, URL: r.HTMLURL, SubmittedAt: r.SubmittedAt}
	if r.User != nil {
		review.Author = r.User.Login
	}
	return review, nil
}

// githubRule is one rule that applies to a branch.
type githubRule struct {
	Type       string `json:"type"`
	Parameters *struct {
		RequiredStatusChecks []struct {
			Context string `json:"context"`
		} `json:"required_status_checks"`
		RequiredApprovingReviewCount int `json:"required_approving_review_count"`
	} `json:"parameters"`
}

// ParseGitHubRemote is the owner and name of the GitHub repository a git
// remote URL names, in the https, scp-like (git@github.com:o/r), or ssh://
// form, with or without .git. It refuses a remote on any host but
// github.com, and never repeats credentials a URL carries.
func ParseGitHubRemote(remote string) (owner, name string, err error) {
	var host, path, shown string
	if strings.Contains(remote, "://") {
		u, err := url.Parse(remote)
		if err != nil {
			return "", "", errors.New("the remote is not a URL flai can read; give its fetch URL as git remote get-url origin prints it")
		}
		u.User = nil
		shown = u.String()
		switch u.Scheme {
		case "https", "http", "ssh", "git":
		default:
			return "", "", fmt.Errorf("the remote %s is not an https, ssh, or git URL; GitHub's remotes are one of these", shown)
		}
		host, path = u.Hostname(), u.Path
	} else {
		at, rest, ok := strings.Cut(remote, ":")
		if !ok || strings.Contains(at, "/") {
			return "", "", errors.New("the remote is neither a URL nor host:path; give its fetch URL as git remote get-url origin prints it")
		}
		host = at[strings.LastIndex(at, "@")+1:]
		path = rest
		shown = host + ":" + path
	}
	if !strings.EqualFold(host, "github.com") {
		return "", "", fmt.Errorf("the remote %s is on %s, not github.com; flai's GitHub forge reaches github.com only", shown, host)
	}
	parts := strings.Split(strings.TrimSuffix(strings.Trim(path, "/"), ".git"), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("the remote %s does not name a repository as owner/name", shown)
	}
	return parts[0], parts[1], nil
}
