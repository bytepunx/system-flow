package dashboard

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

// fakeRelease is a GitHub release the stand-in serves, with its assets' contents.
type fakeRelease struct {
	tag    string
	assets map[string][]byte
}

// gitHub is a stand-in for GitHub's releases API of o/r; down answers every
// request with 503, failing serves the assets with 500.
type gitHub struct {
	srv      *httptest.Server
	releases []fakeRelease
	down     bool
	failing  bool
}

// newGitHub serves releases at /repos/o/r/releases and their assets at /assets/<release>/<name>.
func newGitHub(t *testing.T, releases ...fakeRelease) *gitHub {
	t.Helper()
	g := &gitHub{releases: releases}
	g.srv = httptest.NewServer(http.HandlerFunc(g.serve))
	t.Cleanup(g.srv.Close)
	return g
}

func (g *gitHub) serve(w http.ResponseWriter, r *http.Request) {
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
	if g.failing {
		http.Error(w, "broken", http.StatusInternalServerError)
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

// listOf is a well-formed list of repository at version with index digest idx.
func listOf(repository, version, idx string) []byte {
	ref := repository + ":" + version
	return []byte(idx + " " + ref + " index\n" + amd64 + " " + ref + " linux/amd64\n")
}

// pubPEM is key's public half as a release's cosign.pub holds it.
func pubPEM(t *testing.T, key *ecdsa.PrivateKey) []byte {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
}

// signedRelease is the flaiover release of version publishing list, its
// signature by key, and key's public half as cosign.pub.
func signedRelease(t *testing.T, key *ecdsa.PrivateKey, version string, list []byte) fakeRelease {
	t.Helper()
	name := ListName(version)
	return fakeRelease{tag: selfupgrade.DashboardTagPrefix + version, assets: map[string][]byte{
		name:          list,
		name + ".sig": sign(t, key, list),
		"cosign.pub":  pubPEM(t, key),
	}}
}

// resolverFor is a resolver reading g's releases and trusting key, with its
// cache in a fresh state directory.
func resolverFor(t *testing.T, g *gitHub, key *ecdsa.PrivateKey) Resolver {
	t.Helper()
	return Resolver{
		Releases: selfupgrade.Options{Repo: "o/r", APIBase: g.srv.URL, Client: g.srv.Client()},
		Cache:    newCache(t, key),
	}
}

// cachedNothing fails the test when c holds any list or any file.
func cachedNothing(t *testing.T, c Cache) {
	t.Helper()
	if entries, err := os.ReadDir(c.Dir); err == nil && len(entries) > 0 {
		t.Errorf("cache holds %d files, want none", len(entries))
	}
}

func TestLatestAndAVersionResolveToTheReleaseAndItsIndexDigest(t *testing.T) {
	key := newKey(t)
	g := newGitHub(t,
		signedRelease(t, key, goodV, listFor(goodV, index)),
		signedRelease(t, key, "0.43.0", listFor("0.43.0", other)),
	)
	r := resolverFor(t, g, key)

	for _, tc := range []struct{ tag, version, idx string }{
		{"latest", "0.43.0", other},
		{"", "0.43.0", other},
		{goodV, goodV, index},
		{"v" + goodV, goodV, index},
	} {
		got, err := r.Resolve(context.Background(), repo, tc.tag)
		if err != nil {
			t.Fatalf("tag %q: %v", tc.tag, err)
		}
		if got.Version != tc.version || got.Reference != repo+"@"+tc.idx || got.List.Index() != tc.idx {
			t.Errorf("tag %q resolved to %+v, want version %s and %s@%s", tc.tag, got, tc.version, repo, tc.idx)
		}
		if got.Cached || got.Unreachable != nil {
			t.Errorf("tag %q: cached %v (%v), want a fresh list", tc.tag, got.Cached, got.Unreachable)
		}
	}
}

func TestAVerifiedListIsCachedWithItsVersion(t *testing.T) {
	key := newKey(t)
	g := newGitHub(t, signedRelease(t, key, goodV, []byte(goodList)))
	r := resolverFor(t, g, key)
	if _, err := r.Resolve(context.Background(), repo, goodV); err != nil {
		t.Fatal(err)
	}
	l, ok, err := r.Cache.Get(goodV)
	if err != nil || !ok || l.Index() != index || l.Repository != repo {
		t.Fatalf("cached %+v ok %v err %v, want the release's list", l, ok, err)
	}
}

func TestWithoutCacheKeysTheReleaseClientsKeysAreTrusted(t *testing.T) {
	key := newKey(t)
	g := newGitHub(t, signedRelease(t, key, goodV, []byte(goodList)))
	r := resolverFor(t, g, key)
	r.Cache.Keys = nil
	r.Releases.Keys = []*ecdsa.PublicKey{&key.PublicKey}
	if _, err := r.Resolve(context.Background(), repo, goodV); err != nil {
		t.Fatal(err)
	}
	r.Cache.Keys = r.Releases.Keys
	if _, ok, err := r.Cache.Get(goodV); err != nil || !ok {
		t.Fatalf("not cached: ok %v err %v", ok, err)
	}
}

func TestEachRefusalNamesTheReleaseAndTheReasonAndCachesNothing(t *testing.T) {
	trusted, unknown := newKey(t), newKey(t)
	good := []byte(goodList)
	without := func(rel fakeRelease, name string) fakeRelease {
		delete(rel.assets, name)
		return rel
	}
	with := func(rel fakeRelease, name string, data []byte) fakeRelease {
		rel.assets[name] = data
		return rel
	}
	name := ListName(goodV)

	for _, tc := range []struct {
		name    string
		release fakeRelease
		tag     string
		kind    RefusalKind
		version string
		says    string
	}{
		{"no such release", signedRelease(t, trusted, goodV, good), "0.1.0", RefusalNoRelease, "0.1.0", "published are 0.42.3"},
		{"a branch tag", signedRelease(t, trusted, goodV, good), "main", RefusalNotAVersion, "main", "unsigned with --build, or with --image naming another repository"},
		{"a sha tag", signedRelease(t, trusted, goodV, good), "sha-abc1234", RefusalNotAVersion, "sha-abc1234", "not a flaiover release version"},
		{"a major.minor tag", signedRelease(t, trusted, goodV, good), "1.2", RefusalNotAVersion, "1.2", "only latest or a release version X.Y.Z"},
		{"no list", without(signedRelease(t, trusted, goodV, good), name), goodV, RefusalMissing, goodV, "publishes no " + name},
		{"no signature", without(signedRelease(t, trusted, goodV, good), name+".sig"), goodV, RefusalMissing, goodV, "publishes no " + name + ".sig"},
		{"a signature over other data", with(signedRelease(t, trusted, goodV, good), name+".sig", sign(t, trusted, []byte("other"))), goodV, RefusalBad, goodV, "may have been altered"},
		{"a signature that is not base64", with(signedRelease(t, trusted, goodV, good), name+".sig", []byte("!!!\n")), goodV, RefusalBad, goodV, "not base64"},
		{"an unknown key with no cosign.pub", without(signedRelease(t, unknown, goodV, good), "cosign.pub"), goodV, RefusalBad, goodV, "publishes no cosign.pub"},
		{"an unknown key", signedRelease(t, unknown, goodV, good), goodV, RefusalUnknownKey, goodV, "upgrade flai first"},
		{"a malformed list", signedRelease(t, trusted, goodV, []byte("garbage\n")), goodV, RefusalMalformed, goodV, "line 1"},
		{"another version's list", signedRelease(t, trusted, goodV, listFor("0.40.0", index)), goodV, RefusalVersion, goodV, "lists flaiover 0.40.0"},
		{"another repository's list", signedRelease(t, trusted, goodV, listOf("ghcr.io/else/flaiover", goodV, index)), goodV, RefusalRepository, goodV, "the signed list names ghcr.io/else/flaiover, not the configured " + repo},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := newGitHub(t, tc.release)
			r := resolverFor(t, g, trusted)
			_, err := r.Resolve(context.Background(), repo, tc.tag)
			var refused *Refusal
			if !errors.As(err, &refused) {
				t.Fatalf("got %v, want a *Refusal", err)
			}
			if refused.Kind != tc.kind || refused.Version != tc.version {
				t.Errorf("kind %s version %s, want %s and %s", refused.Kind, refused.Version, tc.kind, tc.version)
			}
			msg := err.Error()
			if !strings.HasPrefix(msg, "flaiover "+tc.version+" is not started: ") || !strings.Contains(msg, tc.says) {
				t.Errorf("message %q does not name the release and say %q", msg, tc.says)
			}
			cachedNothing(t, r.Cache)
		})
	}
}

func TestLatestWithNoPublishedReleaseIsRefused(t *testing.T) {
	g := newGitHub(t)
	r := resolverFor(t, g, newKey(t))
	_, err := r.Resolve(context.Background(), repo, "latest")
	var refused *Refusal
	if !errors.As(err, &refused) || refused.Kind != RefusalNoRelease || refused.Version != "latest" {
		t.Fatalf("got %v, want a no-release refusal of latest", err)
	}
	if !strings.Contains(err.Error(), "o/r has no published flaiover release") {
		t.Errorf("message %q does not name the repository", err)
	}
}

func TestACachedListIsUsedWhenTheReleasesCannotBeReached(t *testing.T) {
	key := newKey(t)
	g := newGitHub(t)
	g.down = true
	r := resolverFor(t, g, key)
	put(t, r.Cache, key, goodV, []byte(goodList))
	put(t, r.Cache, key, "0.43.0", listFor("0.43.0", other))
	put(t, r.Cache, key, "0.44.0", listOf("ghcr.io/else/flaiover", "0.44.0", other))

	for _, tc := range []struct{ tag, version, idx string }{
		{goodV, goodV, index},
		{"v0.43.0", "0.43.0", other},
		{"latest", "0.43.0", other},
	} {
		got, err := r.Resolve(context.Background(), repo, tc.tag)
		if err != nil {
			t.Fatalf("tag %q: %v", tc.tag, err)
		}
		if !got.Cached || got.Unreachable == nil || got.Version != tc.version || got.Reference != repo+"@"+tc.idx {
			t.Errorf("tag %q resolved to %+v, want cached %s at %s", tc.tag, got, tc.version, tc.idx)
		}
	}
}

func TestACachedListIsUsedWhenTheListCannotBeDownloaded(t *testing.T) {
	key := newKey(t)
	g := newGitHub(t, signedRelease(t, key, goodV, []byte(goodList)))
	g.failing = true
	r := resolverFor(t, g, key)
	put(t, r.Cache, key, goodV, []byte(goodList))
	got, err := r.Resolve(context.Background(), repo, goodV)
	if err != nil || !got.Cached || got.Reference != repo+"@"+index {
		t.Fatalf("got %+v err %v, want the cached list", got, err)
	}
}

func TestACachedListIsUsedWhenTheServerIsGone(t *testing.T) {
	key := newKey(t)
	g := newGitHub(t)
	r := resolverFor(t, g, key)
	put(t, r.Cache, key, goodV, []byte(goodList))
	g.srv.Close()
	got, err := r.Resolve(context.Background(), repo, goodV)
	if err != nil || !got.Cached || got.Version != goodV {
		t.Fatalf("got %+v err %v, want the cached list", got, err)
	}
}

func TestUnreachableReleasesWithNothingUsableCachedFail(t *testing.T) {
	key := newKey(t)
	for _, tc := range []struct {
		tag, cached, says string
	}{
		{goodV, "", "no list of it is cached"},
		{"latest", "", "no list of " + repo + " is cached"},
		{goodV, "ghcr.io/else/flaiover", "its cached list names ghcr.io/else/flaiover"},
		{"latest", "ghcr.io/else/flaiover", "no list of " + repo + " is cached"},
	} {
		g := newGitHub(t)
		g.down = true
		r := resolverFor(t, g, key)
		if tc.cached != "" {
			put(t, r.Cache, key, goodV, listOf(tc.cached, goodV, index))
		}
		_, err := r.Resolve(context.Background(), repo, tc.tag)
		var refused *Refusal
		if err == nil || errors.As(err, &refused) {
			t.Fatalf("tag %q with %q cached: got %v, want the failure to reach the releases", tc.tag, tc.cached, err)
		}
		if msg := err.Error(); !strings.Contains(msg, "could not be reached") || !strings.Contains(msg, "503") || !strings.Contains(msg, tc.says) {
			t.Errorf("tag %q: message %q does not say the releases could not be reached and %q", tc.tag, msg, tc.says)
		}
	}
}
