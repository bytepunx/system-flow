package selfupgrade

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

// assetRequests counts the requests f answered for release assets.
func (f *releases) assetRequests() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for path, count := range f.asked {
		if strings.HasPrefix(path, "/assets/") {
			n += count
		}
	}
	return n
}

// unknownKeys are the keys of a build that trusts only the unknown test key.
func unknownKeys(t *testing.T) []*ecdsa.PublicKey {
	t.Helper()
	key, err := parsePublicKey(testdataFile(t, "unknown.pub"))
	if err != nil {
		t.Fatal(err)
	}
	return []*ecdsa.PublicKey{key}
}

// withURL is rel with asset name's URL replaced, as a replaced asset's is.
func withURL(rel *Release, name, url string) *Release {
	out := *rel
	out.Assets = append([]Asset(nil), rel.Assets...)
	for i := range out.Assets {
		if out.Assets[i].Name == name {
			out.Assets[i].URL = url
		}
	}
	return &out
}

// saved saves v and opens what it saved afresh.
func saved(t *testing.T, v *Verdicts, dir string) *Verdicts {
	t.Helper()
	if err := v.Save(); err != nil {
		t.Fatal(err)
	}
	return OpenVerdicts(dir)
}

// sameAnswer fails unless got is the answer want is: both nil, or the same
// refusal.
func sameAnswer(t *testing.T, what string, got, want error) {
	t.Helper()
	if (got == nil) != (want == nil) {
		t.Fatalf("%s: got %v, want %v", what, got, want)
	}
	if want == nil {
		return
	}
	var g, w *SignatureError
	if !errors.As(got, &g) || !errors.As(want, &w) || g.Tag != w.Tag || g.Kind != w.Kind || g.Reason != w.Reason || got.Error() != want.Error() {
		t.Fatalf("%s: got %v, want %v", what, got, want)
	}
}

func TestVerdictsAnswerAKeptVerdictWithoutADownload(t *testing.T) {
	altered := signed(t, "1.1.0", "trusted", []byte("bin"))
	altered.assets[checksumsName] = append(altered.assets[checksumsName], []byte("0000  other\n")...)
	f := serveReleases(t,
		signed(t, "2.0.0", "trusted", []byte("bin")),
		altered,
		signed(t, "1.0.0", "unknown", []byte("bin")),
	)
	opt := f.options(trustedKeys(t))
	ctx := context.Background()
	for version, kind := range map[string]SignatureKind{"2.0.0": "", "1.1.0": SignatureBad, "1.0.0": SignatureUnknownKey} {
		t.Run(version, func(t *testing.T) {
			rel := releaseOf(t, f, version)
			want := VerifyRelease(ctx, opt, rel)
			if kind != "" {
				_ = refusal(t, want, kind)
			}
			dir := t.TempDir()
			v := OpenVerdicts(dir)
			sameAnswer(t, "the first check", v.Verify(ctx, opt, rel), want)

			before := f.assetRequests()
			v = saved(t, v, dir)
			sameAnswer(t, "the kept verdict", v.Verify(ctx, opt, rel), want)
			if got := f.assetRequests(); got != before {
				t.Errorf("a kept verdict downloads nothing: %d asset requests, want %d", got, before)
			}
		})
	}
}

func TestVerdictsMissWhenAnAssetOrTheKeysChange(t *testing.T) {
	f := serveReleases(t, signed(t, "2.0.0", "trusted", []byte("bin")))
	ctx := context.Background()
	rel := releaseOf(t, f, "2.0.0")
	trusted := f.options(trustedKeys(t))
	for name, tc := range map[string]struct {
		opt   Options
		rel   *Release
		asked string // the asset the changed check downloads again
		kind  SignatureKind
	}{
		"checksums.txt replaced": {trusted, withURL(rel, checksumsName, f.srv.URL+"/assets/2.0.0/checksums.txt?id=2"), "/assets/2.0.0/checksums.txt", ""},
		"signature replaced":     {trusted, withURL(rel, signatureName, f.srv.URL+"/assets/2.0.0/checksums.txt.sig?id=2"), "/assets/2.0.0/checksums.txt.sig", ""},
		"other trusted keys":     {f.options(unknownKeys(t)), rel, "/assets/2.0.0/checksums.txt", SignatureUnknownKey},
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			v := OpenVerdicts(dir)
			if err := v.Verify(ctx, trusted, rel); err != nil {
				t.Fatal(err)
			}
			v = saved(t, v, dir)
			before := f.count(tc.asked)
			err := v.Verify(ctx, tc.opt, tc.rel)
			if f.count(tc.asked) == before {
				t.Errorf("the changed check downloads %s again", tc.asked)
			}
			sameAnswer(t, "the changed check", err, VerifyRelease(ctx, tc.opt, tc.rel))
			if tc.kind != "" {
				_ = refusal(t, err, tc.kind)
			}
		})
	}
}

func TestVerdictsNeverKeepADownloadFailure(t *testing.T) {
	pauseFor(t, 0)
	var failed int
	var mu sync.Mutex
	broken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		failed++
		mu.Unlock()
		http.Error(w, "unavailable", http.StatusInternalServerError)
	}))
	t.Cleanup(broken.Close)
	f := serveReleases(t, signed(t, "2.0.0", "trusted", []byte("bin")))
	rel := withURL(releaseOf(t, f, "2.0.0"), checksumsName, broken.URL+"/checksums.txt")
	opt := f.options(trustedKeys(t))
	ctx := context.Background()
	dir := t.TempDir()
	v := OpenVerdicts(dir)
	for attempt := 1; attempt <= 2; attempt++ {
		err := v.Verify(ctx, opt, rel)
		var se *SignatureError
		if err == nil || errors.As(err, &se) || !strings.Contains(err.Error(), "500") {
			t.Fatalf("check %d: the failed download is answered as one: %v", attempt, err)
		}
		mu.Lock()
		got := failed
		mu.Unlock()
		if got != attempt {
			t.Fatalf("check %d downloads checksums.txt again: %d requests, want %d", attempt, got, attempt)
		}
	}
	if err := v.Save(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(verdictsPath(dir)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("nothing was kept, so nothing is written: %v", err)
	}
}

func TestVerdictsNeverKeepAMissingSignature(t *testing.T) {
	f := serveReleases(t, without(signed(t, "2.0.0", "trusted", []byte("bin")), signatureName))
	dir := t.TempDir()
	v := OpenVerdicts(dir)
	_ = refusal(t, v.Verify(context.Background(), f.options(trustedKeys(t)), releaseOf(t, f, "2.0.0")), SignatureMissing)
	if err := v.Save(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(verdictsPath(dir)); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("a missing signature is not kept, so nothing is written: %v", err)
	}
}

func TestCorruptVerdictsLoadAsEmpty(t *testing.T) {
	f := serveReleases(t, signed(t, "2.0.0", "trusted", []byte("bin")))
	opt := f.options(trustedKeys(t))
	rel := releaseOf(t, f, "2.0.0")
	ctx := context.Background()
	for name, content := range map[string]string{
		"not JSON":       "{\"format\": 1, \"verdicts\": [",
		"another format": `{"format": 99, "verdicts": []}`,
		"a missing verdict kept": `{"format": 1, "verdicts": [{"tag": "flai/v2.0.0", "checksums_url": "` + f.srv.URL + `/assets/2.0.0/checksums.txt", ` +
			`"signature_url": "` + f.srv.URL + `/assets/2.0.0/checksums.txt.sig", "public_key_url": "` + f.srv.URL + `/assets/2.0.0/cosign.pub", ` +
			`"keys": "` + fingerprint(t, trustedKeys(t)) + `", "verified": false, "kind": "missing", "reason": "it publishes no checksums.txt.sig"}]}`,
	} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			if err := os.MkdirAll(filepath.Dir(verdictsPath(dir)), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(verdictsPath(dir), []byte(content), 0o600); err != nil {
				t.Fatal(err)
			}
			v := OpenVerdicts(dir)
			before := f.assetRequests()
			if err := v.Verify(ctx, opt, rel); err != nil {
				t.Fatalf("the check runs as without a file: %v", err)
			}
			if f.assetRequests() == before {
				t.Error("nothing was read from the file, so the check downloads")
			}
			before = f.assetRequests()
			if err := saved(t, v, dir).Verify(ctx, opt, rel); err != nil || f.assetRequests() != before {
				t.Errorf("the file is replaced by one holding the new verdict: %v, %d asset requests", err, f.assetRequests()-before)
			}
		})
	}
}

func TestVerdictsWithoutACacheDirKeepNothingOnDisk(t *testing.T) {
	f := serveReleases(t, signed(t, "2.0.0", "trusted", []byte("bin")))
	opt := f.options(trustedKeys(t))
	rel := releaseOf(t, f, "2.0.0")
	v := OpenVerdicts("")
	if err := v.Verify(context.Background(), opt, rel); err != nil {
		t.Fatal(err)
	}
	if err := v.Save(); err != nil {
		t.Errorf("saving without a cache_dir does nothing: %v", err)
	}
}

func TestVerdictsAreSafeForConcurrentChecks(t *testing.T) {
	f := serveReleases(t,
		signed(t, "2.0.0", "trusted", []byte("bin")),
		signed(t, "1.0.0", "unknown", []byte("bin")),
		without(signed(t, "0.9.0", "trusted", []byte("bin")), signatureName),
	)
	opt := f.options(trustedKeys(t))
	ctx := context.Background()
	rels := []*Release{releaseOf(t, f, "2.0.0"), releaseOf(t, f, "1.0.0"), releaseOf(t, f, "0.9.0")}
	dir := t.TempDir()
	v := OpenVerdicts(dir)
	var wg sync.WaitGroup
	for range 4 {
		for _, rel := range rels {
			wg.Go(func() { _ = v.Verify(ctx, opt, rel) })
		}
	}
	wg.Wait()
	v = saved(t, v, dir)
	before := f.assetRequests()
	for _, rel := range rels {
		sameAnswer(t, rel.Version, v.Verify(ctx, opt, rel), VerifyRelease(ctx, opt, rel))
	}
	// Only the checks of 2.0.0 and 1.0.0 were kept; VerifyRelease above
	// downloads their assets again, 2 for 2.0.0 and 3 for 1.0.0.
	if got := f.assetRequests() - before; got != 5 {
		t.Errorf("the kept verdicts download nothing: %d asset requests, want the 5 VerifyRelease made", got)
	}
}

func TestKeysFingerprintIgnoresOrder(t *testing.T) {
	trusted, unknown := trustedKeys(t)[0], unknownKeys(t)[0]
	ab := fingerprint(t, []*ecdsa.PublicKey{trusted, unknown})
	if ba := fingerprint(t, []*ecdsa.PublicKey{unknown, trusted}); ab != ba {
		t.Errorf("the same keys in another order: %s, want %s", ba, ab)
	}
	if a := fingerprint(t, []*ecdsa.PublicKey{trusted}); a == ab || a == fingerprint(t, []*ecdsa.PublicKey{unknown}) {
		t.Errorf("other keys have another fingerprint: %s", a)
	}
}

// fingerprint is KeysFingerprint of keys, failing the test on an error.
func fingerprint(t *testing.T, keys []*ecdsa.PublicKey) string {
	t.Helper()
	fp, err := KeysFingerprint(keys)
	if err != nil {
		t.Fatal(err)
	}
	return fp
}
