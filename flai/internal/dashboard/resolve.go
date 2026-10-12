package dashboard

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

// Resolved is the image a dashboard tag resolves to: a flaiover release, its
// verified digest list, and the image to pull and run by digest.
type Resolved struct {
	Version   string // the flaiover release, X.Y.Z without its v
	List      List   // the release's verified digest list
	Reference string // <repository>@sha256:<index digest>, the image to pull and run
	// Cached is true when the list came from the cache because the release
	// could not be reached; Unreachable then says why.
	Cached      bool
	Unreachable error
}

// RefusalKind says why a release's image is not started.
type RefusalKind string

// The reasons a release's image is not started.
const (
	// RefusalNotAVersion: the tag is neither latest nor a release version X.Y.Z.
	RefusalNotAVersion RefusalKind = "not-a-version"
	// RefusalNoRelease: no published flaiover release has the version.
	RefusalNoRelease RefusalKind = "no-release"
	// RefusalMissing: the release publishes no digest list or no signature for it.
	RefusalMissing RefusalKind = "missing"
	// RefusalBad: no key this flai trusts verifies the list's signature, nor the release's own key.
	RefusalBad RefusalKind = "bad"
	// RefusalUnknownKey: the list is signed by the release's own key, which this flai does not trust.
	RefusalUnknownKey RefusalKind = "unknown-key"
	// RefusalMalformed: the list verifies but is not in the digest list format.
	RefusalMalformed RefusalKind = "malformed"
	// RefusalVersion: the list verifies but is for another release.
	RefusalVersion RefusalKind = "version"
	// RefusalRepository: the list verifies but names another image repository.
	RefusalRepository RefusalKind = "repository"
)

// Refusal is why a flaiover release's image is not started; nothing is
// cached for it.
type Refusal struct {
	Version string      // the release, X.Y.Z, or the tag as configured when it names none
	Kind    RefusalKind // why it is refused
	Reason  string      // what was found
}

// Error names the release, why it is not started, and what to do.
func (e *Refusal) Error() string {
	head := fmt.Sprintf("flaiover %s is not started: %s", e.Version, e.Reason)
	switch e.Kind {
	case RefusalNotAVersion:
		return head + "; only latest or a release version X.Y.Z is run, from its signed digest list: set dashboard.tag to one of them, or run an image unsigned with --build, or with --image naming another repository"
	case RefusalNoRelease:
		return head + "; set dashboard.tag to a published version, or to latest for the newest"
	case RefusalMissing:
		return head + ", so its image cannot be verified (a release from before signing has none); set dashboard.tag to a signed release, such as latest"
	case RefusalUnknownKey:
		return head + ", which this flai does not know (the release key may have changed since this flai was built); upgrade flai first (flai self-upgrade), then start the dashboard again"
	case RefusalMalformed:
		return head
	case RefusalVersion:
		return head + "; the release is mislabelled: run no image it names, and report it to the repository's maintainers"
	case RefusalRepository:
		return head + "; set dashboard.image to the repository the release lists, or run the configured one unsigned with --image"
	}
	return head + "; the list may have been altered since it was signed: run no image it names, and report the release to the repository's maintainers"
}

// Resolver resolves dashboard.tag to a flaiover release's verified image,
// fetching releases through Releases and keeping each verified list in Cache.
// A list must verify against Cache.Keys or, when it has none, the keys
// selfupgrade.TrustedKeys answers for Releases, which the cache then uses too.
type Resolver struct {
	Releases selfupgrade.Options // the release client's options: repository, API, token, HTTP client
	Cache    Cache               // where verified lists are kept, and the keys they verify against
}

// Resolve answers the image repository runs at tag: "" or "latest" for the
// newest published flaiover release, "X.Y.Z" or "vX.Y.Z" for that release. A
// release refused, for the tag, its list, its signature, or its repository,
// is a *Refusal. When the releases cannot be reached it answers the cached
// list of the version, or for latest the newest cached list of repository,
// with Cached set; with none cached, the failure to reach them.
func (r Resolver) Resolve(ctx context.Context, repository, tag string) (Resolved, error) {
	want, err := wantedVersion(tag)
	if err != nil {
		return Resolved{}, err
	}
	cache, err := r.trustingCache()
	if err != nil {
		return Resolved{}, err
	}
	published, err := selfupgrade.List(ctx, r.Releases, selfupgrade.DashboardTagPrefix)
	if err != nil {
		return cache.offline(repository, want, err)
	}
	rel, err := pick(published, want, r.Releases.Repo)
	if err != nil {
		return Resolved{}, err
	}
	l, data, sig, err := FetchList(ctx, r.Releases, rel, cache.Keys)
	var refused *Refusal
	switch {
	case errors.As(err, &refused):
		return Resolved{}, err
	case err != nil:
		return cache.offline(repository, rel.Version, err)
	}
	if l.Repository != repository {
		return Resolved{}, &Refusal{Version: rel.Version, Kind: RefusalRepository, Reason: fmt.Sprintf("the signed list names %s, not the configured %s", l.Repository, repository)}
	}
	if err := cache.Put(rel.Version, data, sig); err != nil {
		return Resolved{}, fmt.Errorf("resolving flaiover %s failed: %w", rel.Version, err)
	}
	return resolved(l), nil
}

// wantedVersion is tag's release version, X.Y.Z, or "" for the newest,
// refusing a tag that is neither.
func wantedVersion(tag string) (string, error) {
	if tag == "" || tag == "latest" {
		return "", nil
	}
	if v := strings.TrimPrefix(tag, "v"); buildinfo.Bare(v) {
		return v, nil
	}
	return "", &Refusal{Version: tag, Kind: RefusalNotAVersion, Reason: fmt.Sprintf("dashboard.tag %q is not a flaiover release version", tag)}
}

// trustingCache is r.Cache with the keys a list must verify against.
func (r Resolver) trustingCache() (Cache, error) {
	cache := r.Cache
	if len(cache.Keys) > 0 {
		return cache, nil
	}
	keys, err := selfupgrade.TrustedKeys(r.Releases)
	if err != nil {
		return Cache{}, fmt.Errorf("reading the release keys to verify a flaiover digest list failed: %w; this flai build is broken, so install a release of flai", err)
	}
	cache.Keys = keys
	return cache, nil
}

// pick is the release of published at want, or the newest when want is "".
func pick(published []selfupgrade.Release, want, repo string) (selfupgrade.Release, error) {
	if repo == "" {
		repo = "the release repository"
	}
	if len(published) == 0 {
		version := want
		if version == "" {
			version = "latest"
		}
		return selfupgrade.Release{}, &Refusal{Version: version, Kind: RefusalNoRelease, Reason: fmt.Sprintf("%s has no published flaiover release (no GitHub release tagged %sX.Y.Z that is neither a draft nor a prerelease)", repo, selfupgrade.DashboardTagPrefix)}
	}
	if want == "" {
		return published[0], nil
	}
	for _, rel := range published {
		if rel.Version == want {
			return rel, nil
		}
	}
	return selfupgrade.Release{}, &Refusal{Version: want, Kind: RefusalNoRelease, Reason: fmt.Sprintf("%s has no published release tagged %s%s (drafts and prereleases are not published); published are %s", repo, selfupgrade.DashboardTagPrefix, want, selfupgrade.VersionList(published))}
}

// FetchList answers rel's digest list, parsed, with the list and its
// signature as published, once the signature verifies against one of keys.
// A list rel does not publish, whose signature does not verify, that does
// not parse, or that is for another release is a *Refusal; a failed
// download is answered as one, not as a refusal.
func FetchList(ctx context.Context, opt selfupgrade.Options, rel selfupgrade.Release, keys []*ecdsa.PublicKey) (l List, data, sig []byte, err error) {
	name := ListName(rel.Version)
	data, sig, err = selfupgrade.VerifySigned(ctx, opt, rel, name, keys)
	var unverified *selfupgrade.SignatureError
	if errors.As(err, &unverified) {
		return List{}, nil, nil, &Refusal{Version: rel.Version, Kind: RefusalKind(unverified.Kind), Reason: fmt.Sprintf("the digest list %s of release %s cannot be verified: %s", name, rel.Tag, unverified.Reason)}
	}
	if err != nil {
		return List{}, nil, nil, err
	}
	l, err = Parse(data)
	if err != nil {
		return List{}, nil, nil, &Refusal{Version: rel.Version, Kind: RefusalMalformed, Reason: fmt.Sprintf("%s of release %s verifies, but %v", name, rel.Tag, err)}
	}
	if l.Version != rel.Version {
		return List{}, nil, nil, &Refusal{Version: rel.Version, Kind: RefusalVersion, Reason: fmt.Sprintf("%s of release %s verifies, but it lists flaiover %s", name, rel.Tag, l.Version)}
	}
	return l, data, sig, nil
}

// offline answers, when the releases could not be reached for unreached, the
// cached list of version, or for "" the newest cached list of repository.
func (c Cache) offline(repository, version string, unreached error) (Resolved, error) {
	failed := func(why string) error {
		what := "flaiover " + version
		if version == "" {
			what = "the newest flaiover release"
		}
		return fmt.Errorf("resolving %s failed: its release could not be reached (%w), and %s; check the connection and the token, then run it again", what, unreached, why)
	}
	if version == "" {
		lists, skipped := c.All()
		for _, l := range lists {
			if l.Repository == repository {
				return offlineResolved(l, unreached), nil
			}
		}
		return Resolved{}, failed(withSkipped("no list of "+repository+" is cached", skipped))
	}
	l, found, err := c.Get(version)
	switch {
	case err != nil:
		return Resolved{}, failed(fmt.Sprintf("the cached list is not used: %v", err))
	case !found:
		return Resolved{}, failed("no list of it is cached")
	case l.Repository != repository:
		return Resolved{}, failed(fmt.Sprintf("its cached list names %s, not the configured %s", l.Repository, repository))
	}
	return offlineResolved(l, unreached), nil
}

// withSkipped adds to why the cached lists skipped because they do not verify.
func withSkipped(why string, skipped []error) string {
	if len(skipped) == 0 {
		return why
	}
	return fmt.Sprintf("%s that verifies (skipped: %v)", why, errors.Join(skipped...))
}

// resolved is the answer for a verified list l.
func resolved(l List) Resolved {
	return Resolved{Version: l.Version, List: l, Reference: l.Reference()}
}

// offlineResolved is the answer for l, read from the cache because the
// release could not be reached for unreached.
func offlineResolved(l List, unreached error) Resolved {
	out := resolved(l)
	out.Cached, out.Unreachable = true, unreached
	return out
}
