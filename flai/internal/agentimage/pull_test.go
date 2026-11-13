package agentimage

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/dashboard"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

const (
	pullV      = "1.2.3"
	indexD     = "sha256:0dce07ffc38f1fecc91cbb659524525db687eb12030bea40209f2ecf3636f15d"
	amd64D     = "sha256:b3c17aafc5f6026830b885843021a783cc6e129937cc4f1f1dc44e546ff38ec0"
	strangerD  = "sha256:1111111111111111111111111111111111111111111111111111111111111111"
	claudeRepo = "ghcr.io/bytepunx/flai-agent-claude-code"
)

// releaseAPI is a stand-in for GitHub's releases API of o/r, serving each
// release's assets at /assets/<release>/<name>; down answers every request with 503.
type releaseAPI struct {
	srv      *httptest.Server
	mu       sync.Mutex
	releases []fakeFlaiRelease
	down     bool
	requests int
}

// fakeFlaiRelease is a GitHub release the stand-in serves, with its assets' contents.
type fakeFlaiRelease struct {
	tag    string
	assets map[string][]byte
}

// newReleaseAPI serves releases until the test ends.
func newReleaseAPI(t *testing.T, releases ...fakeFlaiRelease) *releaseAPI {
	t.Helper()
	g := &releaseAPI{releases: releases}
	g.srv = httptest.NewServer(http.HandlerFunc(g.serve))
	t.Cleanup(g.srv.Close)
	return g
}

func (g *releaseAPI) serve(w http.ResponseWriter, r *http.Request) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.requests++
	if g.down {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
		return
	}
	if r.URL.Path == "/repos/o/r/releases" {
		type asset struct{ Name, URL string }
		type release struct {
			TagName string  `json:"tag_name"`
			Assets  []asset `json:"assets"`
		}
		out := []release{}
		for i, rel := range g.releases {
			gr := release{TagName: rel.tag}
			for name := range rel.assets {
				gr.Assets = append(gr.Assets, asset{Name: name, URL: fmt.Sprintf("%s/assets/%d/%s", g.srv.URL, i, name)})
			}
			out = append(out, gr)
		}
		_ = json.NewEncoder(w).Encode(out)
		return
	}
	rest, ok := strings.CutPrefix(r.URL.Path, "/assets/")
	i, name, _ := strings.Cut(rest, "/")
	n, err := strconv.Atoi(i)
	if !ok || err != nil || n >= len(g.releases) {
		http.NotFound(w, r)
		return
	}
	data, ok := g.releases[n].assets[name]
	if !ok {
		http.NotFound(w, r)
		return
	}
	_, _ = w.Write(data)
}

// setDown makes the stand-in answer 503 to everything, or serve again.
func (g *releaseAPI) setDown(down bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.down = down
}

// releaseKey is a fresh P-256 key, as cosign generates.
func releaseKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// signBlob is data's signature by key as cosign sign-blob writes it.
func signBlob(t *testing.T, key *ecdsa.PrivateKey, data []byte) []byte {
	t.Helper()
	sum := sha256.Sum256(data)
	der, err := ecdsa.SignASN1(rand.Reader, key, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	return []byte(base64.StdEncoding.EncodeToString(der) + "\n")
}

// publicPEM is key's public half as a release's cosign.pub holds it.
func publicPEM(t *testing.T, key *ecdsa.PrivateKey) []byte {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
}

// digestList is a well-formed list of repository at version with index digest indexD.
func digestList(repository, version string) []byte {
	ref := repository + ":" + version
	return []byte(indexD + " " + ref + " index\n" + amd64D + " " + ref + " linux/amd64\n")
}

// flaiRelease is the flai release of version publishing the claude-code
// list, its signature by key, and key's public half as cosign.pub.
func flaiRelease(t *testing.T, key *ecdsa.PrivateKey, version string, list []byte) fakeFlaiRelease {
	t.Helper()
	name := DigestListName("claude-code", version)
	return fakeFlaiRelease{tag: selfupgrade.TagPrefix + version, assets: map[string][]byte{
		name:          list,
		name + ".sig": signBlob(t, key, list),
		"cosign.pub":  publicPEM(t, key),
	}}
}

// dockerCalls is an Exec that records the arguments of each docker run,
// fails, writing said, each run whose first argument is in fail, and writes
// nothing until its context ends for each whose first argument is block.
type dockerCalls struct {
	mu    sync.Mutex
	runs  [][]string
	fail  map[string]string
	block string
}

// exec is the fake's Exec.
func (d *dockerCalls) exec(ctx context.Context, _ string, _ []string, out io.Writer, _ string, args ...string) error {
	d.mu.Lock()
	d.runs = append(d.runs, args)
	said, failing := d.fail[args[0]]
	blocking := d.block == args[0]
	d.mu.Unlock()
	switch {
	case blocking:
		<-ctx.Done()
		return ctx.Err()
	case failing:
		_, _ = io.WriteString(out, said+"\n")
		return errExit
	}
	return nil
}

// ran is each docker run's arguments, in order.
func (d *dockerCalls) ran() [][]string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.runs
}

// pulledAs is an Inspector answering repoDigests for every image.
func pulledAs(repoDigests ...string) Inspector {
	return func(string) (Info, bool, error) {
		return Info{ID: "sha256:img", RepoDigests: repoDigests}, true, nil
	}
}

// pullerFor is a Puller reading g's releases, trusting key, caching under a
// fresh state folder, running docker through d with no retries, and
// inspecting with inspect.
func pullerFor(t *testing.T, g *releaseAPI, key *ecdsa.PrivateKey, d *dockerCalls, inspect Inspector) Puller {
	t.Helper()
	return Puller{
		Releases: selfupgrade.Options{Repo: "o/r", APIBase: g.srv.URL, Client: g.srv.Client()},
		Keys:     []*ecdsa.PublicKey{&key.PublicKey},
		CacheDir: t.TempDir(),
		Run:      Runner{Exec: d.exec, Delays: []time.Duration{}},
		Inspect:  inspect,
	}
}

func TestPublishedAndDigestListNameNameTheReleasesImages(t *testing.T) {
	if got := Published("claude-code"); got != claudeRepo {
		t.Errorf("Published = %q", got)
	}
	if got := DigestListName("bedrock", pullV); got != "flai-agent-bedrock_1.2.3.digests" {
		t.Errorf("DigestListName = %q", got)
	}
}

// S-0740, ADR-0295 §1: a verified list pulls the image by its index digest,
// tags it with each tag given, and is cached under agent-image-digests/<image>.
func TestPullPullsByTheVerifiedListsIndexDigestAndTagsIt(t *testing.T) {
	key := releaseKey(t)
	g := newReleaseAPI(t, flaiRelease(t, key, pullV, digestList(claudeRepo, pullV)))
	d := &dockerCalls{}
	p := pullerFor(t, g, key, d, pulledAs(claudeRepo+"@"+indexD))
	tags := []string{"flai-agent-claude-code:1.2.3", "flai-agent-claude-code:local"}

	got, err := p.Pull(context.Background(), "claude-code", "v"+pullV, tags)
	if err != nil {
		t.Fatal(err)
	}
	ref := claudeRepo + "@" + indexD
	want := Pulled{Image: "claude-code", Version: pullV, Reference: ref, Tags: tags}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Pulled = %+v, want %+v", got, want)
	}
	wantRuns := [][]string{{"pull", ref}, {"tag", ref, tags[0]}, {"tag", ref, tags[1]}}
	if !reflect.DeepEqual(d.ran(), wantRuns) {
		t.Errorf("docker ran %v, want %v", d.ran(), wantRuns)
	}
	for _, name := range []string{"1.2.3.digests", "1.2.3.digests.sig"} {
		if _, err := os.Stat(filepath.Join(p.CacheDir, "agent-image-digests", "claude-code", name)); err != nil {
			t.Errorf("%s not cached: %v", name, err)
		}
	}
}

// S-0740: a list or an image that does not verify is refused with
// dashboard's refusal kinds, never stood in for by a build, and a list
// refused runs no docker pull.
func TestPullRefusesWhatDoesNotVerify(t *testing.T) {
	trusted, stranger := releaseKey(t), releaseKey(t)
	good := digestList(claudeRepo, pullV)
	name := DigestListName("claude-code", pullV)
	for _, c := range []struct {
		name    string
		release fakeFlaiRelease
		inspect Inspector
		kind    dashboard.RefusalKind
		says    string
		pulls   bool
	}{
		{"signed by another key", func() fakeFlaiRelease {
			rel := flaiRelease(t, stranger, pullV, good)
			delete(rel.assets, "cosign.pub")
			return rel
		}(), nil, dashboard.RefusalBad, "may have been altered", false},
		{"signed by the release's own key", flaiRelease(t, stranger, pullV, good), nil, dashboard.RefusalUnknownKey, "upgrade flai first", false},
		{"another repository", flaiRelease(t, trusted, pullV, digestList("ghcr.io/else/flai-agent-claude-code", pullV)), nil, dashboard.RefusalRepository, "names ghcr.io/else/flai-agent-claude-code, not " + claudeRepo, false},
		{"another version", flaiRelease(t, trusted, pullV, digestList(claudeRepo, "1.2.2")), nil, dashboard.RefusalVersion, "lists flai-agent-claude-code 1.2.2", false},
		{"malformed", flaiRelease(t, trusted, pullV, []byte("garbage\n")), nil, dashboard.RefusalMalformed, "reading the flai-agent-claude-code digest list failed at line 1", false},
		{"an unlisted image", flaiRelease(t, trusted, pullV, good), pulledAs(claudeRepo + "@" + strangerD), dashboard.RefusalUnlisted, "RepoDigests " + claudeRepo + "@" + strangerD + ", none of which the signed list " + name + " names", true},
		{"an image with no RepoDigests", flaiRelease(t, trusted, pullV, good), pulledAs(), dashboard.RefusalUnlisted, "RepoDigests none", true},
	} {
		t.Run(c.name, func(t *testing.T) {
			g := newReleaseAPI(t, c.release)
			d := &dockerCalls{}
			p := pullerFor(t, g, trusted, d, c.inspect)
			_, err := p.Pull(context.Background(), "claude-code", pullV, []string{"x:1"})
			var refused *dashboard.Refusal
			if !errors.As(err, &refused) || refused.Kind != c.kind {
				t.Fatalf("Pull = %v, want a %s refusal", err, c.kind)
			}
			msg := err.Error()
			if !strings.HasPrefix(msg, "flai-agent-claude-code 1.2.3 is not pulled: ") || !strings.Contains(msg, c.says) {
				t.Errorf("message %q does not name the image and say %q", msg, c.says)
			}
			if FallBack(err) {
				t.Error("FallBack is true, but a build must never stand in for what does not verify")
			}
			wantRuns := [][]string(nil)
			if c.pulls {
				wantRuns = [][]string{{"pull", claudeRepo + "@" + indexD}}
			}
			if !reflect.DeepEqual(d.ran(), wantRuns) {
				t.Errorf("docker ran %v, want %v (no tag)", d.ran(), wantRuns)
			}
			if c.pulls {
				return
			}
			if entries, err := os.ReadDir(filepath.Join(p.CacheDir, "agent-image-digests", "claude-code")); err == nil && len(entries) > 0 {
				t.Errorf("a refused list was cached: %d files", len(entries))
			}
		})
	}
}

// S-0740: a release that publishes no list, or no such release, lets the
// build stand in, with no docker pull.
func TestPullFallsBackWhenTheReleaseHasNoList(t *testing.T) {
	key := releaseKey(t)
	without := flaiRelease(t, key, pullV, digestList(claudeRepo, pullV))
	delete(without.assets, DigestListName("claude-code", pullV))
	for _, c := range []struct {
		name    string
		release fakeFlaiRelease
		kind    dashboard.RefusalKind
		says    string
	}{
		{"no list", without, dashboard.RefusalMissing, "publishes no flai-agent-claude-code_1.2.3.digests"},
		{"no such release", flaiRelease(t, key, "1.2.2", digestList(claudeRepo, "1.2.2")), dashboard.RefusalNoRelease, "o/r has no published release tagged flai/v1.2.3 (drafts and prereleases are not published); published are 1.2.2"},
	} {
		t.Run(c.name, func(t *testing.T) {
			g := newReleaseAPI(t, c.release)
			d := &dockerCalls{}
			_, err := pullerFor(t, g, key, d, nil).Pull(context.Background(), "claude-code", pullV, nil)
			var refused *dashboard.Refusal
			if !errors.As(err, &refused) || refused.Kind != c.kind || !strings.Contains(err.Error(), c.says) {
				t.Fatalf("Pull = %v, want a %s refusal saying %q", err, c.kind, c.says)
			}
			if !FallBack(err) {
				t.Error("FallBack is false, want the build to stand in")
			}
			if len(d.ran()) != 0 {
				t.Errorf("docker ran %v", d.ran())
			}
		})
	}
}

// S-0740: when the releases cannot be reached, the list cached for the
// version is pulled by, Cached set; with none cached the build stands in.
func TestPullUsesTheCachedListWhenTheReleasesCannotBeReached(t *testing.T) {
	key := releaseKey(t)
	g := newReleaseAPI(t, flaiRelease(t, key, pullV, digestList(claudeRepo, pullV)))
	d := &dockerCalls{}
	p := pullerFor(t, g, key, d, pulledAs(claudeRepo+"@"+indexD))
	if _, err := p.Pull(context.Background(), "claude-code", pullV, nil); err != nil {
		t.Fatal(err)
	}
	g.setDown(true)

	got, err := p.Pull(context.Background(), "claude-code", pullV, []string{"x:1"})
	if err != nil {
		t.Fatal(err)
	}
	if !got.Cached || got.Reference != claudeRepo+"@"+indexD || got.Version != pullV {
		t.Errorf("Pulled = %+v, want the cached list's image", got)
	}

	noCache := pullerFor(t, g, key, d, nil)
	noCache.CacheDir = ""
	for _, c := range []struct {
		name    string
		p       Puller
		version string
	}{
		{"none cached", pullerFor(t, g, key, d, nil), pullV},
		{"no cache", noCache, pullV},
		{"another version cached", p, "1.2.4"},
	} {
		t.Run(c.name, func(t *testing.T) {
			before := len(d.ran())
			_, err := c.p.Pull(context.Background(), "claude-code", c.version, nil)
			var refused *dashboard.Refusal
			if err == nil || errors.As(err, &refused) || !strings.Contains(err.Error(), "could not be reached") || !strings.Contains(err.Error(), "503") {
				t.Fatalf("Pull = %v, want the failure to reach the releases", err)
			}
			if !FallBack(err) {
				t.Error("FallBack is false, want the build to stand in")
			}
			if len(d.ran()) != before {
				t.Errorf("docker ran %v", d.ran()[before:])
			}
		})
	}
}

// S-0740: a docker pull that fails or stalls lets the build stand in; a
// failed tag does not.
func TestPullFallsBackOnAFailedPullButNotAFailedTag(t *testing.T) {
	key := releaseKey(t)
	g := newReleaseAPI(t, flaiRelease(t, key, pullV, digestList(claudeRepo, pullV)))

	d := &dockerCalls{fail: map[string]string{"pull": "Error response from daemon: Get \"https://ghcr.io/v2/\": dial tcp: connection reset by peer"}}
	_, err := pullerFor(t, g, key, d, nil).Pull(context.Background(), "claude-code", pullV, []string{"x:1"})
	var step *StepError
	if !errors.As(err, &step) || !FallBack(err) {
		t.Fatalf("Pull = %v, FallBack %v; want the step's failure and true", err, FallBack(err))
	}
	if len(d.ran()) != 1 {
		t.Errorf("docker ran %v, want the pull alone", d.ran())
	}
	d = &dockerCalls{block: "pull"}
	p := pullerFor(t, g, key, d, nil)
	p.Run.Stall = 20 * time.Millisecond
	_, err = p.Pull(context.Background(), "claude-code", pullV, []string{"x:1"})
	var stall *StallError
	if !errors.As(err, &stall) || !FallBack(err) {
		t.Fatalf("a stalled pull: Pull = %v, FallBack %v; want the stall and true", err, FallBack(err))
	}

	d = &dockerCalls{fail: map[string]string{"tag": "Error: no such image"}}
	_, err = pullerFor(t, g, key, d, pulledAs(claudeRepo+"@"+indexD)).Pull(context.Background(), "claude-code", pullV, []string{"x:1"})
	if !errors.As(err, &step) || FallBack(err) {
		t.Fatalf("Pull = %v, FallBack %v; want the tag's failure and false", err, FallBack(err))
	}
}

// S-0740, ADR-0295 §3: a development flai pulls nothing and asks no release;
// a stopped pull is not stood in for.
func TestPullRefusesADevelopmentFlaiAndAStoppedPull(t *testing.T) {
	key := releaseKey(t)
	g := newReleaseAPI(t, flaiRelease(t, key, pullV, digestList(claudeRepo, pullV)))
	d := &dockerCalls{}
	p := pullerFor(t, g, key, d, nil)
	for _, v := range []string{"dev", "", "0.0.0-abc1234"} {
		_, err := p.Pull(context.Background(), "claude-code", v, nil)
		var refused *dashboard.Refusal
		if err == nil || errors.As(err, &refused) || FallBack(err) || !strings.Contains(err.Error(), "development flai") {
			t.Errorf("version %q: Pull = %v, FallBack %v; want a development refusal, not stood in for", v, err, FallBack(err))
		}
	}
	if _, err := p.Pull(context.Background(), "../x", pullV, nil); err == nil || FallBack(err) {
		t.Errorf("an image name that is not one: Pull = %v", err)
	}
	g.mu.Lock()
	requests := g.requests
	g.mu.Unlock()
	if requests != 0 || len(d.ran()) != 0 {
		t.Errorf("asked the releases %d times and ran docker %v", requests, d.ran())
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := p.Pull(ctx, "claude-code", pullV, nil)
	if err == nil || FallBack(err) || !errors.Is(err, context.Canceled) {
		t.Errorf("a cancelled pull: Pull = %v, FallBack %v; want the cancellation, not stood in for", err, FallBack(err))
	}
	if FallBack(nil) {
		t.Error("FallBack(nil) is true")
	}
}
