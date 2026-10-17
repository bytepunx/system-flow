// Package forge is the contract pull-request acceptance reaches a project's
// forge through: find a story's pull request by its head branch, open or
// update it, read the checks and reviews on it, and merge it by a method.
// A Forge is bound to one repository when it is made, and gets its token
// from a TokenFunc at each call rather than holding it. GitHub is the one
// adapter; what is particular to a forge (its API paths, how it names a
// repository, how it reports checks) stays inside its adapter, so another
// forge implements the same contract without a change to its callers.
package forge

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Forge is a forge's view of one repository, as acceptance uses it.
type Forge interface {
	// Name is the forge's name, as acceptance.forge spells it: "github".
	Name() string
	// Repository is the repository the Forge is bound to, for display: "owner/name".
	Repository() string
	// Access is whether the token can read the repository and push to it.
	Access(ctx context.Context) (Access, error)
	// FindPullRequest is the open pull request whose head branch is head;
	// false when there is none.
	FindPullRequest(ctx context.Context, head string) (PullRequest, bool, error)
	// OpenPullRequest opens a pull request from in.Head into in.Base, or,
	// when one is open for in.Head already, updates its title, body, and
	// base to in's, and answers the pull request as it then stands.
	OpenPullRequest(ctx context.Context, in PullRequestInput) (PullRequest, error)
	// Checks are the checks on the pull request's head commit, with each
	// check the base branch requires that has not reported yet as pending.
	Checks(ctx context.Context, pr PullRequest) ([]Check, error)
	// Reviews are the reviews submitted on the pull request, oldest first.
	Reviews(ctx context.Context, pr PullRequest) ([]Review, error)
	// Merge merges the pull request by method, refusing with
	// ErrMergeRefused when the forge will not merge it or its head is no
	// longer pr.HeadSHA.
	Merge(ctx context.Context, pr PullRequest, method MergeMethod) (Merged, error)
}

// TokenFunc answers the token a Forge calls its forge with; an adapter
// calls it at each request and keeps nothing it answers.
type TokenFunc func() (string, error)

// Access is what the token may do in the repository.
type Access struct {
	// Read is whether the token can read the repository.
	Read bool
	// Push is whether the token can push to it.
	Push bool
}

// PullRequestState is where a pull request stands.
type PullRequestState string

// The states a pull request can be in.
const (
	PullRequestOpen   PullRequestState = "open"
	PullRequestClosed PullRequestState = "closed"
	PullRequestMerged PullRequestState = "merged"
)

// PullRequest is a pull request as the forge reports it.
type PullRequest struct {
	// Number is the pull request's number in its repository.
	Number int
	// URL is the pull request's page on the forge.
	URL string
	// Head is the branch it merges from.
	Head string
	// Base is the branch it merges into.
	Base string
	// HeadSHA is the commit at the head of Head when it was read.
	HeadSHA string
	// Title is its title.
	Title string
	// Body is its description.
	Body string
	// State is whether it is open, closed, or merged.
	State PullRequestState
	// Draft is whether it is a draft.
	Draft bool
}

// PullRequestInput is what OpenPullRequest opens or updates a pull request with.
type PullRequestInput struct {
	// Head is the branch to merge from.
	Head string
	// Base is the branch to merge into.
	Base string
	// Title is the pull request's title.
	Title string
	// Body is its description.
	Body string
}

// CheckState is where a check stands.
type CheckState string

// The states a check can be in. A check that finished without failing,
// including one skipped or neutral, has passed, as the forges count it
// toward a required check; one that finished any other way has failed.
const (
	CheckPending CheckState = "pending"
	CheckPassed  CheckState = "passed"
	CheckFailed  CheckState = "failed"
)

// Check is one check on a pull request's head commit.
type Check struct {
	// Name is the check's name, as the base branch's rules require it.
	Name string
	// State is whether it is pending, passed, or failed.
	State CheckState
	// Required is whether the base branch requires it to pass.
	Required bool
	// URL is the page with its details and log; empty when the forge has none.
	URL string
	// Summary is the forge's one line on what it found; may be empty.
	Summary string
}

// ReviewState is what a review says.
type ReviewState string

// The states a review can be in.
const (
	ReviewApproved         ReviewState = "approved"
	ReviewChangesRequested ReviewState = "changes-requested"
	ReviewCommented        ReviewState = "commented"
	ReviewDismissed        ReviewState = "dismissed"
	ReviewPending          ReviewState = "pending"
)

// Review is one review of a pull request.
type Review struct {
	// Author is the reviewer's login.
	Author string
	// State is what the review says.
	State ReviewState
	// Body is the review's comment; may be empty.
	Body string
	// URL is the review's page on the forge.
	URL string
	// SubmittedAt is when it was submitted, RFC 3339; empty while pending.
	SubmittedAt string
}

// MergeMethod is how a pull request is merged into its base.
type MergeMethod string

// The merge methods, as acceptance.merge_method spells them.
const (
	MergeMerge  MergeMethod = "merge"
	MergeSquash MergeMethod = "squash"
	MergeRebase MergeMethod = "rebase"
)

// MergeMethods are the merge methods, in the order they are offered.
func MergeMethods() []MergeMethod {
	return []MergeMethod{MergeMerge, MergeSquash, MergeRebase}
}

// Valid is whether m is one of MergeMethods.
func (m MergeMethod) Valid() bool {
	switch m {
	case MergeMerge, MergeSquash, MergeRebase:
		return true
	}
	return false
}

// ParseMergeMethod is the merge method s names.
func ParseMergeMethod(s string) (MergeMethod, error) {
	m := MergeMethod(s)
	if !m.Valid() {
		names := make([]string, 0, len(MergeMethods()))
		for _, v := range MergeMethods() {
			names = append(names, string(v))
		}
		return "", fmt.Errorf("merge method %q is not one flai knows; use one of %s", s, strings.Join(names, ", "))
	}
	return m, nil
}

// Merged is what a merge made.
type Merged struct {
	// SHA is the commit the merge left at the head of the base branch.
	SHA string
}

// ErrNotFound is the forge's answer for a repository or pull request that
// does not exist or that the token cannot see.
var ErrNotFound = errors.New("not found on the forge")

// ErrUnauthorized is the forge's answer for a token it does not accept.
var ErrUnauthorized = errors.New("the forge did not accept the token")

// ErrForbidden is the forge's answer for a token it accepts that may not
// do what was asked, or that has used up its rate limit.
var ErrForbidden = errors.New("the forge refused the token this request")

// ErrMergeRefused is the forge's answer for a pull request it will not
// merge: a required check or review is outstanding, it conflicts, or its
// head has moved since it was read.
var ErrMergeRefused = errors.New("the forge refused to merge the pull request")
