package serve

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytepunx/system-flow/flai/internal/dashboard"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

const (
	imgRepo      = "ghcr.io/bytepunx/flaiover"
	imgVersion   = "0.42.3"
	imgIndex     = "sha256:0dce07ffc38f1fecc91cbb659524525db687eb12030bea40209f2ecf3636f15d"
	imgAMD64     = "sha256:b3c17aafc5f6026830b885843021a783cc6e129937cc4f1f1dc44e546ff38ec0"
	imgARM64     = "sha256:b8ec2bdf93dfbd0cdbe8b756bdd0ce30d66242a93bd99d5e869f0668a475eb0c"
	imgOther     = "sha256:1111111111111111111111111111111111111111111111111111111111111111"
	imgContainer = "4c1f0d6e9a2b"
	imgID        = "sha256:9f2a7c11d3e4b5a6978877665544332211009988776655443322110099887766"
	imgURL       = "http://localhost:4242"
)

// imgList is the signed digest list of flaiover imgVersion.
var imgList = []byte(imgIndex + " " + imgRepo + ":" + imgVersion + " index\n" +
	imgAMD64 + " " + imgRepo + ":" + imgVersion + " linux/amd64\n" +
	imgARM64 + " " + imgRepo + ":" + imgVersion + " linux/arm64\n")

// fakeDocker answers docker's commands from canned output, keyed by their
// arguments; absent is docker not on PATH.
type fakeDocker struct {
	absent bool
	out    map[string]string
	fail   map[string]bool
	calls  []string
}

func (f *fakeDocker) Run(_, name string, args ...string) (string, error) {
	key := strings.Join(args, " ")
	f.calls = append(f.calls, name+" "+key)
	out, ok := f.out[key]
	if !ok || f.fail[key] {
		return out, fmt.Errorf("%s %s: exit status 1\n%s", name, key, out)
	}
	return out, nil
}

func (f *fakeDocker) RunInput(dir, name, _ string, args ...string) (string, error) {
	return f.Run(dir, name, args...)
}

func (f *fakeDocker) LookPath(name string) (string, error) {
	if f.absent {
		return "", errors.New("exec: \"" + name + "\": executable file not found in $PATH")
	}
	return "/usr/bin/" + name, nil
}

// dockerWith is Docker running the container flaiover from an image with
// repoDigests, carrying FLAIOVER_VERSION version.
func dockerWith(t *testing.T, version string, repoDigests []string) *fakeDocker {
	t.Helper()
	container := []map[string]any{{
		"Id":     imgContainer,
		"Image":  imgID,
		"Config": map[string]any{"Env": []string{"PATH=/usr/bin", "FLAIOVER_VERSION=" + version}},
	}}
	image := []map[string]any{{"Id": imgID, "RepoDigests": repoDigests}}
	return &fakeDocker{out: map[string]string{
		"inspect --type container flaiover": mustJSON(t, container),
		"image inspect " + imgID:            mustJSON(t, image),
	}}
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// imgKey is a fresh P-256 key, as cosign generates.
func imgKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// imgSign is data's signature by key as cosign sign-blob writes it.
func imgSign(t *testing.T, key *ecdsa.PrivateKey, data []byte) []byte {
	t.Helper()
	sum := sha256.Sum256(data)
	der, err := ecdsa.SignASN1(rand.Reader, key, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	return []byte(base64.StdEncoding.EncodeToString(der) + "\n")
}

// imgCache is a cache in a fresh state directory trusting key.
func imgCache(t *testing.T, key *ecdsa.PrivateKey) dashboard.Cache {
	t.Helper()
	return dashboard.NewCache(t.TempDir(), []*ecdsa.PublicKey{&key.PublicKey})
}

// releaseServer is a stand-in for GitHub's releases of o/r publishing the
// flaiover imgVersion release with list signed by key; down answers 503.
func releaseServer(t *testing.T, key *ecdsa.PrivateKey, list []byte, down bool) *httptest.Server {
	t.Helper()
	name := dashboard.ListName(imgVersion)
	assets := map[string][]byte{name: list, name + ".sig": imgSign(t, key, list)}
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if down {
			http.Error(w, "unavailable", http.StatusServiceUnavailable)
			return
		}
		if r.URL.Path == "/repos/o/r/releases" {
			type asset struct{ Name, URL string }
			rel := struct {
				TagName string  `json:"tag_name"`
				Assets  []asset `json:"assets"`
			}{TagName: selfupgrade.DashboardTagPrefix + imgVersion}
			for n := range assets {
				rel.Assets = append(rel.Assets, asset{Name: n, URL: srv.URL + "/assets/" + n})
			}
			_ = json.NewEncoder(w).Encode([]any{rel})
			return
		}
		data, ok := assets[strings.TrimPrefix(r.URL.Path, "/assets/")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write(data)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// fetcherFrom fetches through a resolver reading srv's releases into cache.
func fetcherFrom(srv *httptest.Server, cache dashboard.Cache) ListFetcher {
	r := dashboard.Resolver{
		Releases: selfupgrade.Options{Repo: "o/r", APIBase: srv.URL, Client: srv.Client()},
		Cache:    cache,
	}
	return ResolverFetcher(r, imgRepo)
}

// noFetch fails the test if a list is fetched.
func noFetch(t *testing.T) ListFetcher {
	return func(context.Context, string) (dashboard.List, error) {
		t.Errorf("a list was fetched, though none needed to be")
		return dashboard.List{}, errors.New("not fetched")
	}
}

func cachePut(t *testing.T, c dashboard.Cache, key *ecdsa.PrivateKey) {
	t.Helper()
	if err := c.Put(imgVersion, imgList, imgSign(t, key, imgList)); err != nil {
		t.Fatal(err)
	}
}

func TestAnImageACachedVerifiedListNamesIsListed(t *testing.T) {
	key := imgKey(t)
	cache := imgCache(t, key)
	cachePut(t, cache, key)
	docker := dockerWith(t, imgVersion, []string{imgRepo + "@" + imgAMD64})

	v := CheckImage(context.Background(), ImageCheck{Runner: docker, URL: imgURL, Cache: cache, Fetch: noFetch(t)})

	if v.Kind != ImageListed || v.Digest != imgAMD64 || v.Release != imgVersion || v.Reason != "" {
		t.Fatalf("verdict %+v, want listed at %s by flaiover %s", v, imgAMD64, imgVersion)
	}
	if v.Container != imgContainer || v.Image != imgID || v.Version != imgVersion {
		t.Fatalf("container %q image %q version %q, want %q %q %q", v.Container, v.Image, v.Version, imgContainer, imgID, imgVersion)
	}
}

func TestAnImageNoVerifiedListNamesIsUnlisted(t *testing.T) {
	key := imgKey(t)
	cache := imgCache(t, key)
	cachePut(t, cache, key)
	rds := []string{imgRepo + "@" + imgOther}
	docker := dockerWith(t, imgVersion, rds)

	v := CheckImage(context.Background(), ImageCheck{Runner: docker, URL: imgURL, Cache: cache, Fetch: noFetch(t)})

	if v.Kind != ImageUnlisted || v.Digest != "" || v.Release != "" {
		t.Fatalf("verdict %+v, want unlisted", v)
	}
	if len(v.RepoDigests) != 1 || v.RepoDigests[0] != rds[0] || v.Image != imgID || v.Container != imgContainer {
		t.Fatalf("verdict %+v, want the image %s and its digests %v", v, imgID, rds)
	}
	if !strings.Contains(v.Reason, "no verified digest list names") {
		t.Fatalf("reason %q does not say no list names it", v.Reason)
	}
}

func TestAListedDigestOfAnotherRepositoryIsUnlisted(t *testing.T) {
	key := imgKey(t)
	cache := imgCache(t, key)
	cachePut(t, cache, key)
	docker := dockerWith(t, imgVersion, []string{"example.com/someone/flaiover@" + imgAMD64})

	v := CheckImage(context.Background(), ImageCheck{Runner: docker, URL: imgURL, Cache: cache, Fetch: noFetch(t)})

	if v.Kind != ImageUnlisted {
		t.Fatalf("verdict %+v, want unlisted", v)
	}
}

func TestALocalImageWithNoRepoDigestsIsUnlisted(t *testing.T) {
	key := imgKey(t)
	cache := imgCache(t, key)
	cachePut(t, cache, key)
	docker := dockerWith(t, "0.0.0", nil)

	v := CheckImage(context.Background(), ImageCheck{Runner: docker, URL: imgURL, Cache: cache, Fetch: noFetch(t)})

	if v.Kind != ImageUnlisted || len(v.RepoDigests) != 0 || v.Image != imgID {
		t.Fatalf("verdict %+v, want unlisted with no digests", v)
	}
	if !strings.Contains(v.Reason, "no RepoDigests") {
		t.Fatalf("reason %q does not say the image has no RepoDigests", v.Reason)
	}
}

func TestAMissingListIsFetchedVerifiedAndCachedThenJudged(t *testing.T) {
	key := imgKey(t)
	cache := imgCache(t, key)
	srv := releaseServer(t, key, imgList, false)
	docker := dockerWith(t, imgVersion, []string{imgRepo + "@" + imgIndex})

	v := CheckImage(context.Background(), ImageCheck{Runner: docker, URL: imgURL, Cache: cache, Fetch: fetcherFrom(srv, cache)})

	if v.Kind != ImageListed || v.Digest != imgIndex || v.Release != imgVersion {
		t.Fatalf("verdict %+v, want listed at the index digest by flaiover %s", v, imgVersion)
	}
	if _, found, err := cache.Get(imgVersion); !found || err != nil {
		t.Fatalf("the fetched list is not cached: found %v err %v", found, err)
	}
}

func TestAFetchedListThatDoesNotVerifyAdmitsNothing(t *testing.T) {
	cache := imgCache(t, imgKey(t))
	srv := releaseServer(t, imgKey(t), imgList, false)
	docker := dockerWith(t, imgVersion, []string{imgRepo + "@" + imgAMD64})

	v := CheckImage(context.Background(), ImageCheck{Runner: docker, URL: imgURL, Cache: cache, Fetch: fetcherFrom(srv, cache)})

	if v.Kind != ImageUnlisted {
		t.Fatalf("verdict %+v, want unlisted", v)
	}
	if !strings.Contains(v.Reason, "not fetched and verified") {
		t.Fatalf("reason %q does not say the list did not verify", v.Reason)
	}
	if _, found, _ := cache.Get(imgVersion); found {
		t.Fatal("a list that does not verify was cached")
	}
}

func TestACachedListThatDoesNotVerifyAdmitsNothing(t *testing.T) {
	key := imgKey(t)
	cache := imgCache(t, key)
	cachePut(t, cache, key)
	tampered := []byte(imgOther + " " + imgRepo + ":" + imgVersion + " index\n")
	if err := os.WriteFile(filepath.Join(cache.Dir, imgVersion+".digests"), tampered, 0o600); err != nil {
		t.Fatal(err)
	}
	docker := dockerWith(t, imgVersion, []string{imgRepo + "@" + imgOther})

	for _, down := range []bool{true, false} {
		srv := releaseServer(t, key, imgList, down)
		v := CheckImage(context.Background(), ImageCheck{Runner: docker, URL: imgURL, Cache: cache, Fetch: fetcherFrom(srv, cache)})
		if v.Kind != ImageUnlisted {
			t.Fatalf("release down %v: verdict %+v, want unlisted", down, v)
		}
	}
	if l, found, err := cache.Get(imgVersion); !found || err != nil || l.Has(imgOther) {
		t.Fatalf("the tampered list was not replaced by the release's: found %v err %v", found, err)
	}
}

func TestADashboardDockerCannotInspectHereIsNotMeasured(t *testing.T) {
	key := imgKey(t)
	cache := imgCache(t, key)
	cachePut(t, cache, key)
	rds := []string{imgRepo + "@" + imgAMD64}
	noContainer := dockerWith(t, imgVersion, rds)
	noContainer.out["inspect --type container flaiover"] = "Error: No such container: flaiover"
	noContainer.fail = map[string]bool{"inspect --type container flaiover": true}

	cases := []struct {
		name   string
		docker *fakeDocker
		url    string
		reason string
	}{
		{"no docker", &fakeDocker{absent: true}, imgURL, "not on PATH"},
		{"no container", noContainer, imgURL, "no container flaiover"},
		{"another host by address", dockerWith(t, imgVersion, rds), "http://192.168.1.20:4242", "not this host"},
		{"another host by name", dockerWith(t, imgVersion, rds), "https://dash.example.com", "not this host"},
		{"no host", dockerWith(t, imgVersion, rds), "", "names no host"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := CheckImage(context.Background(), ImageCheck{Runner: tc.docker, URL: tc.url, Cache: cache, Fetch: noFetch(t)})
			if v.Kind != ImageNotMeasured || !strings.Contains(v.Reason, tc.reason) {
				t.Fatalf("verdict %+v, want not measured saying %q", v, tc.reason)
			}
			if strings.Contains(tc.reason, "host") && len(tc.docker.calls) != 0 {
				t.Fatalf("docker was asked %v about a dashboard on another host", tc.docker.calls)
			}
		})
	}
}

func TestADashboardOnThisHostIsMeasured(t *testing.T) {
	key := imgKey(t)
	cache := imgCache(t, key)
	cachePut(t, cache, key)
	for _, u := range []string{"http://localhost:4242", "http://LOCALHOST:4242", "http://127.0.0.1:4242", "http://[::1]:4242"} {
		docker := dockerWith(t, imgVersion, []string{imgRepo + "@" + imgARM64})
		v := CheckImage(context.Background(), ImageCheck{Runner: docker, URL: u, Cache: cache, Fetch: noFetch(t)})
		if v.Kind != ImageListed {
			t.Fatalf("%s: verdict %+v, want listed", u, v)
		}
	}
}
