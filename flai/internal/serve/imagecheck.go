package serve

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/dashboard"
	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// dashboardContainer is the one dashboard container flai dashboard starts on
// this host, for every project it serves.
const dashboardContainer = "flaiover"

// ImageKind is what flai serve concluded about the image the dashboard
// container runs (design/system/release-signing.md § flai measures the
// container before it dials).
type ImageKind string

// The verdicts of an image check.
const (
	// ImageListed: a verified digest list names one of the image's RepoDigests; dial.
	ImageListed ImageKind = "listed"
	// ImageUnlisted: no verified digest list names the image, or it has no RepoDigests; do not dial.
	ImageUnlisted ImageKind = "unlisted"
	// ImageNotMeasured: Docker cannot say what the dashboard runs here; fall back to the stamp check.
	ImageNotMeasured ImageKind = "not-measured"
)

// ImageVerdict is what an image check found about the dashboard container.
type ImageVerdict struct {
	Kind        ImageKind `json:"kind"`
	Container   string    `json:"container,omitempty"`    // the container's ID; "" when not measured before Docker named one
	Image       string    `json:"image,omitempty"`        // the image ID the container runs, sha256:…; "" as Container
	Version     string    `json:"version,omitempty"`      // the container's FLAIOVER_VERSION, X.Y.Z; "" when it carries none
	RepoDigests []string  `json:"repo_digests,omitempty"` // the image's RepoDigests, <repository>@sha256:…; none for a local build
	Digest      string    `json:"digest,omitempty"`       // listed: the digest, sha256:…, a verified list names
	Release     string    `json:"release,omitempty"`      // listed: the flaiover release whose verified list names Digest
	Reason      string    `json:"reason,omitempty"`       // unlisted and not measured: why, in words for a log and a status line
	// Allowed is set on an unlisted image dialled because dashboard.allow_unsigned
	// is in force for the project (S-0239).
	Allowed bool `json:"allowed,omitempty"`
}

// ListFetcher fetches version's digest list from its flaiover release,
// verifies it, and caches it; it answers an error for a list that does not
// verify or cannot be fetched.
type ListFetcher func(ctx context.Context, version string) (dashboard.List, error)

// ResolverFetcher is a ListFetcher through r, which fetches, verifies, and
// caches a release's list, refusing one that names another repository than
// repository, the image repository the signed lists name.
func ResolverFetcher(r dashboard.Resolver, repository string) ListFetcher {
	return func(ctx context.Context, version string) (dashboard.List, error) {
		res, err := r.Resolve(ctx, repository, version)
		if err != nil {
			return dashboard.List{}, err
		}
		return res.List, nil
	}
}

// ImageCheck is what an image check needs: Docker, the dashboard's URL, the
// verified lists cached beside flai serve's state, and a way to fetch one.
type ImageCheck struct {
	Runner execx.Runner    // runs docker
	URL    string          // the dashboard URL flai serve dials
	Cache  dashboard.Cache // the verified digest lists, with the keys they verify against
	Fetch  ListFetcher     // fetches the running version's list when none is cached; nil fetches none
}

// containerInspect is what CheckImage reads of docker inspect's container.
type containerInspect struct {
	ID     string `json:"Id"`
	Image  string `json:"Image"`
	Config struct {
		Env []string `json:"Env"`
	} `json:"Config"`
}

// imageInspect is what CheckImage reads of docker image inspect's image.
type imageInspect struct {
	RepoDigests []string `json:"RepoDigests"`
}

// CheckImage asks Docker which image the container flaiover runs and judges
// its RepoDigests against the verified digest lists: listed, unlisted, or
// not measured when Docker cannot say what the dashboard at c.URL runs.
//
// Which digest RepoDigests holds depends on how the image was pulled; S-0234
// checked it (release-signing.md § flai measures the container before it
// dials): a pull by tag or by the index digest records the index digest, a
// pull by a platform's manifest digest records that manifest digest, and
// Docker's classic store can hold both on one image. A list names the index
// and every platform's digest, so each entry is compared with all of them.
func CheckImage(ctx context.Context, c ImageCheck) ImageVerdict {
	if why := notThisHost(c.URL); why != "" {
		return notMeasured(ImageVerdict{}, why)
	}
	if _, err := c.Runner.LookPath("docker"); err != nil {
		return notMeasured(ImageVerdict{}, "docker is not on PATH where flai serve runs")
	}
	v, why := inspectContainer(c.Runner)
	if why != "" {
		return notMeasured(v, why)
	}
	out, err := c.Runner.Run("", "docker", "image", "inspect", v.Image)
	var images []imageInspect
	if err == nil {
		err = decodeOne(out, &images)
	}
	if err != nil {
		return notMeasured(v, fmt.Sprintf("docker image inspect of %s failed: %v", v.Image, err))
	}
	v.RepoDigests = images[0].RepoDigests
	return judge(ctx, c, v)
}

// inspectContainer is the container's ID, image ID, and FLAIOVER_VERSION, or
// why Docker could not say.
func inspectContainer(r execx.Runner) (ImageVerdict, string) {
	out, err := r.Run("", "docker", "inspect", "--type", "container", dashboardContainer)
	if err != nil {
		if strings.Contains(out, "No such container") || strings.Contains(out, "No such object") {
			return ImageVerdict{}, "no container " + dashboardContainer + " is on this host"
		}
		return ImageVerdict{}, fmt.Sprintf("docker inspect of the container %s failed: %v", dashboardContainer, err)
	}
	var containers []containerInspect
	if err := decodeOne(out, &containers); err != nil {
		return ImageVerdict{}, fmt.Sprintf("docker inspect of the container %s failed: %v", dashboardContainer, err)
	}
	ct := containers[0]
	v := ImageVerdict{Container: ct.ID, Image: ct.Image}
	for _, e := range ct.Config.Env {
		if ver, ok := strings.CutPrefix(e, "FLAIOVER_VERSION="); ok && buildinfo.Bare(ver) {
			v.Version = ver
		}
	}
	return v, ""
}

// judge matches v's RepoDigests against the cached verified lists, then, when
// none names them and no list of v.Version is cached that verifies, against
// v.Version's list fetched and verified.
func judge(ctx context.Context, c ImageCheck, v ImageVerdict) ImageVerdict {
	if len(v.RepoDigests) == 0 {
		return unlisted(v, "its image has no RepoDigests, as a local build has none, so no signed digest list names it")
	}
	lists, skipped := c.Cache.All()
	if listed, ok := match(v, lists); ok {
		return listed
	}
	why := "no verified digest list names any of its RepoDigests"
	if len(skipped) > 0 {
		why += fmt.Sprintf(" (cached lists skipped as not verified: %v)", errors.Join(skipped...))
	}
	if v.Version == "" {
		return unlisted(v, why+"; the container carries no FLAIOVER_VERSION, so no list was fetched")
	}
	if _, found, err := c.Cache.Get(v.Version); found && err == nil {
		return unlisted(v, why)
	}
	if c.Fetch == nil {
		return unlisted(v, fmt.Sprintf("%s; no list of flaiover %s is cached", why, v.Version))
	}
	l, err := c.Fetch(ctx, v.Version)
	if err != nil {
		return unlisted(v, fmt.Sprintf("%s; the flaiover %s digest list was not fetched and verified: %v", why, v.Version, err))
	}
	if listed, ok := match(v, []dashboard.List{l}); ok {
		return listed
	}
	return unlisted(v, fmt.Sprintf("%s, nor flaiover %s's, fetched and verified", why, v.Version))
}

// match is v listed by the first of lists naming one of its RepoDigests.
func match(v ImageVerdict, lists []dashboard.List) (ImageVerdict, bool) {
	for _, l := range lists {
		for _, rd := range v.RepoDigests {
			if l.Matches(rd) {
				_, v.Digest, _ = strings.Cut(rd, "@")
				v.Kind, v.Release, v.Reason = ImageListed, l.Version, ""
				return v, true
			}
		}
	}
	return v, false
}

// notThisHost is why the dashboard at raw is not measured here, its host
// being neither localhost nor a loopback address; "" when it is.
func notThisHost(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Hostname() == "" {
		return fmt.Sprintf("the dashboard URL %q names no host", raw)
	}
	host := u.Hostname()
	if strings.EqualFold(host, "localhost") {
		return ""
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return ""
	}
	return fmt.Sprintf("the dashboard URL's host %s is not this host, so Docker here cannot say what it runs", host)
}

// decodeOne decodes docker inspect's JSON array into out, refusing one that
// does not hold exactly one object.
func decodeOne[T any](out string, into *[]T) error {
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), into); err != nil {
		return fmt.Errorf("its output is not the JSON docker inspect writes: %w", err)
	}
	if len(*into) != 1 {
		return fmt.Errorf("it described %d objects, not 1", len(*into))
	}
	return nil
}

// notMeasured is v not measured, for why.
func notMeasured(v ImageVerdict, why string) ImageVerdict {
	v.Kind, v.Reason = ImageNotMeasured, why
	return v
}

// unlisted is v unlisted, for why.
func unlisted(v ImageVerdict, why string) ImageVerdict {
	v.Kind, v.Reason = ImageUnlisted, why
	return v
}
