package agentimage

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/dashboard"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

// Registry is where each flai release publishes its agent images.
const Registry = "ghcr.io/bytepunx"

// digestFolder is the folder under flai serve's state folder that keeps, in
// one folder per image, each verified digest list of the agent images.
const digestFolder = "agent-image-digests"

// imageRe is an image's name as Published and DigestListName take it: one
// path segment of lower-case letters, digits, and inner hyphens.
var imageRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Published is the repository image ("bedrock" or an inner harness such as
// "claude-code") is published to: ghcr.io/bytepunx/flai-agent-<image>.
func Published(image string) string {
	return Registry + "/flai-agent-" + image
}

// DigestListName is the release asset of version's digest list of image:
// flai-agent-<image>_<version>.digests.
func DigestListName(image, version string) string {
	return product(image).ListName(version)
}

// product is the digest lists of image, as a flai release publishes them.
func product(image string) dashboard.Product {
	return dashboard.Product{Name: "flai-agent-" + image, Release: "flai"}
}

// Puller pulls a flai release's agent image by the index digest its signed
// list names (ADR-0295 §1).
type Puller struct {
	Releases selfupgrade.Options // flai's releases: repository, API, token, HTTP client
	Keys     []*ecdsa.PublicKey  // the keys a list must verify against; selfupgrade.TrustedKeys(Releases) when empty
	CacheDir string              // flai serve's state folder, whose agent-image-digests/<image>/ keeps each verified list; "" keeps none
	Run      Runner              // runs docker pull and docker tag
	Inspect  Inspector           // reads the pulled image's RepoDigests; nil asks Run's docker client through execx.System
}

// Pulled is what Pull made: the image, the release version, the reference
// pulled (<repository>@<index digest>), the tags given it, and whether the
// list came from the cache because the release could not be reached.
type Pulled struct {
	Image     string
	Version   string
	Reference string
	Tags      []string
	Cached    bool
}

// Pull pulls image of flai version by the index digest the release's
// verified list names and tags it with each of tags. It finds the release
// flai/v<version>, fetches the list and its signature, and refuses, as a
// *dashboard.Refusal, a list the release does not publish or that does not
// verify, is for another version, or names another repository than
// Published(image), and an image pulled whose RepoDigests the list does not
// name. When the releases cannot be reached it pulls by the list cached for
// version, when one verifies. FallBack says which failures a build may stand
// in for.
func (p Puller) Pull(ctx context.Context, image, version string, tags []string) (Pulled, error) {
	if !imageRe.MatchString(image) {
		return Pulled{}, fmt.Errorf("pulling the agent image %q refused: it is not an image name such as bedrock or claude-code (lower-case letters, digits, and inner hyphens)", image)
	}
	if buildinfo.Development(version) {
		return Pulled{}, fmt.Errorf("pulling the agent image %s for flai %q refused: a development flai has no release whose images can be pulled; build its images on this host instead (ADR-0295 §3)", Published(image), version)
	}
	version = strings.TrimPrefix(version, "v")
	keys := p.Keys
	if len(keys) == 0 {
		var err error
		if keys, err = selfupgrade.TrustedKeys(p.Releases); err != nil {
			return Pulled{}, fmt.Errorf("reading the release keys to verify the digest list of the agent image %s failed: %w; this flai build is broken, so install a release of flai", Published(image), err)
		}
	}
	l, cached, err := p.list(ctx, image, version, keys)
	if err != nil {
		return Pulled{}, err
	}
	return p.pull(ctx, image, l, cached, tags)
}

// list is image's verified list of flai version, fetched from its release
// and cached, or read from the cache, cached true, when the release cannot
// be reached.
func (p Puller) list(ctx context.Context, image, version string, keys []*ecdsa.PublicKey) (dashboard.List, bool, error) {
	repo := Published(image)
	prod := product(image)
	cache := dashboard.Cache{Keys: keys, Product: prod}
	if p.CacheDir != "" {
		cache.Dir = filepath.Join(p.CacheDir, digestFolder, image)
	}
	published, err := selfupgrade.List(ctx, p.Releases, selfupgrade.TagPrefix)
	if err != nil {
		return p.offline(ctx, cache, image, version, err)
	}
	rel, err := find(published, version, p.Releases.Repo, prod)
	if err != nil {
		return dashboard.List{}, false, err
	}
	l, data, sig, err := prod.FetchList(ctx, p.Releases, rel, keys)
	var refused *dashboard.Refusal
	switch {
	case errors.As(err, &refused):
		return dashboard.List{}, false, err
	case err != nil:
		return p.offline(ctx, cache, image, version, err)
	}
	if l.Repository != repo {
		return dashboard.List{}, false, &dashboard.Refusal{Version: version, Kind: dashboard.RefusalRepository, Product: prod, Reason: fmt.Sprintf("the signed list %s names %s, not %s", prod.ListName(version), l.Repository, repo)}
	}
	if cache.Dir != "" {
		if err := cache.Put(version, data, sig); err != nil {
			return dashboard.List{}, false, fmt.Errorf("pulling the agent image %s of flai %s failed: %w", repo, version, err)
		}
	}
	return l, false, nil
}

// find is the release of published at version, or a no-release refusal of
// prod naming repo's published versions.
func find(published []selfupgrade.Release, version, repo string, prod dashboard.Product) (selfupgrade.Release, error) {
	for _, rel := range published {
		if rel.Version == version {
			return rel, nil
		}
	}
	if repo == "" {
		repo = "the release repository"
	}
	reason := fmt.Sprintf("%s has no published release tagged %s%s (drafts and prereleases are not published)", repo, selfupgrade.TagPrefix, version)
	if len(published) > 0 {
		reason += "; published are " + selfupgrade.VersionList(published)
	}
	return selfupgrade.Release{}, &dashboard.Refusal{Version: version, Kind: dashboard.RefusalNoRelease, Product: prod, Reason: reason}
}

// offline is, when version's release of image could not be reached for
// unreached, its list cached in cache, cached true, or why there is none.
func (p Puller) offline(ctx context.Context, cache dashboard.Cache, image, version string, unreached error) (dashboard.List, bool, error) {
	if ctx.Err() != nil {
		return dashboard.List{}, false, stopped(image, version, ctx.Err())
	}
	repo := Published(image)
	failed := func(why string) error {
		return &unreachableError{repo: repo, version: version, err: unreached, why: why}
	}
	if cache.Dir == "" {
		return dashboard.List{}, false, failed("no cache of verified lists was given")
	}
	l, found, err := cache.Get(version)
	switch {
	case err != nil:
		return dashboard.List{}, false, failed(fmt.Sprintf("the cached list is not used: %v", err))
	case !found:
		return dashboard.List{}, false, failed("no list of it is cached")
	case l.Repository != repo:
		return dashboard.List{}, false, failed(fmt.Sprintf("its cached list names %s, not %s", l.Repository, repo))
	}
	return l, true, nil
}

// pull runs docker pull of l's reference, checks the pulled image's
// RepoDigests against l, and tags it with each of tags.
func (p Puller) pull(ctx context.Context, image string, l dashboard.List, cached bool, tags []string) (Pulled, error) {
	ref := l.Reference()
	if err := p.Run.Run(ctx, Step{Name: "pull " + ref, Args: []string{"pull", ref}}); err != nil {
		if ctx.Err() != nil {
			return Pulled{}, stopped(image, l.Version, ctx.Err())
		}
		var failed *StepError
		var stalled *StallError
		if errors.As(err, &failed) || errors.As(err, &stalled) {
			return Pulled{}, &pullError{err: err}
		}
		return Pulled{}, err
	}
	inspect := p.Inspect
	if inspect == nil {
		program := p.Run.Program
		if program == "" {
			program = "docker"
		}
		inspect = Docker(execx.System{}, program)
	}
	info, present, err := inspect(ref)
	switch {
	case err != nil:
		return Pulled{}, fmt.Errorf("checking the agent image %s after its pull failed: %w", ref, err)
	case !present:
		return Pulled{}, fmt.Errorf("checking the agent image %s after its pull failed: docker pull succeeded, but docker image inspect finds no such image; check that the engine answers (docker info), then run it again", ref)
	}
	if !listed(l, info.RepoDigests) {
		have := "none"
		if len(info.RepoDigests) > 0 {
			have = strings.Join(info.RepoDigests, ", ")
		}
		prod := product(image)
		return Pulled{}, &dashboard.Refusal{Version: l.Version, Kind: dashboard.RefusalUnlisted, Product: prod, Reason: fmt.Sprintf("the image pulled as %s has the RepoDigests %s, none of which the signed list %s names", ref, have, prod.ListName(l.Version))}
	}
	for _, tag := range tags {
		if err := p.Run.Run(ctx, Step{Name: "tag " + ref + " as " + tag, Args: []string{"tag", ref, tag}}); err != nil {
			if ctx.Err() != nil {
				return Pulled{}, stopped(image, l.Version, ctx.Err())
			}
			return Pulled{}, err
		}
	}
	return Pulled{Image: image, Version: l.Version, Reference: ref, Tags: tags, Cached: cached}, nil
}

// listed reports whether one of repoDigests is l's repository at a digest l names.
func listed(l dashboard.List, repoDigests []string) bool {
	for _, rd := range repoDigests {
		if l.Matches(rd) {
			return true
		}
	}
	return false
}

// FallBack reports whether err, from Pull, lets the build stand in for the
// pull: the release could not be reached and no list is cached, the release
// or its list does not exist, or docker pull failed or stalled. A list or an
// image that does not verify is never stood in for, since the build would
// hide tampering; nor a failed tag, a development version, or a stopped pull.
func FallBack(err error) bool {
	var stop *stoppedError
	var unreached *unreachableError
	var pulling *pullError
	var refused *dashboard.Refusal
	switch {
	case err == nil, errors.As(err, &stop):
		return false
	case errors.As(err, &unreached), errors.As(err, &pulling):
		return true
	case errors.As(err, &refused):
		return refused.Kind == dashboard.RefusalMissing || refused.Kind == dashboard.RefusalNoRelease
	}
	return false
}

// unreachableError is a release that could not be reached, with no list of
// it cached that verifies.
type unreachableError struct {
	repo, version string
	err           error  // why the release could not be reached
	why           string // why no cached list stands in
}

// Error says which image's list could not be had, why, and what to do.
func (e *unreachableError) Error() string {
	return fmt.Sprintf("pulling the agent image %s of flai %s failed: its release could not be reached (%v), and %s; check the connection and the token, then run it again", e.repo, e.version, e.err, e.why)
}

// Unwrap is why the release could not be reached.
func (e *unreachableError) Unwrap() error { return e.err }

// pullError is a docker pull that failed or stalled.
type pullError struct{ err error }

// Error is the step's own.
func (e *pullError) Error() string { return e.err.Error() }

// Unwrap is the *StepError or *StallError.
func (e *pullError) Unwrap() error { return e.err }

// stoppedError is a pull stopped because its context ended.
type stoppedError struct {
	repo, version string
	err           error
}

// stopped is the pull of image of flai version, stopped for err, its context's.
func stopped(image, version string, err error) error {
	return &stoppedError{repo: Published(image), version: version, err: err}
}

// Error says which pull was stopped and why.
func (e *stoppedError) Error() string {
	return fmt.Sprintf("pulling the agent image %s of flai %s was stopped before it finished: %v", e.repo, e.version, e.err)
}

// Unwrap is the context's error.
func (e *stoppedError) Unwrap() error { return e.err }
