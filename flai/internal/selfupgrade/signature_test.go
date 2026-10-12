package selfupgrade

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// testdataFile reads a fixture under testdata.
func testdataFile(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// privateKey is the test key pair name's private half, testdata/<name>.key.
func privateKey(t *testing.T, name string) *ecdsa.PrivateKey {
	t.Helper()
	block, _ := pem.Decode(testdataFile(t, name+".key"))
	if block == nil {
		t.Fatalf("testdata/%s.key holds no PEM block", name)
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

// trustedKeys are the keys a test's binary trusts: the trusted test key's
// public half, read from testdata/trusted.pub.
func trustedKeys(t *testing.T) []*ecdsa.PublicKey {
	t.Helper()
	key, err := parsePublicKey(testdataFile(t, "trusted.pub"))
	if err != nil {
		t.Fatal(err)
	}
	return []*ecdsa.PublicKey{key}
}

// sign is the signature cosign sign-blob writes for data with key: base64 of
// an ASN.1 ECDSA signature over its SHA-256, and a newline.
func sign(t *testing.T, key *ecdsa.PrivateKey, data []byte) []byte {
	t.Helper()
	sum := sha256.Sum256(data)
	sig, err := ecdsa.SignASN1(rand.Reader, key, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	return []byte(base64.StdEncoding.EncodeToString(sig) + "\n")
}

// fakeRelease is one release a releases fake serves: its version and its
// assets by name.
type fakeRelease struct {
	version string
	assets  map[string][]byte
}

// releases is a fake of the releases API serving rels, newest first as
// given, with each asset at /assets/<version>/<name>; it counts the requests
// for each path.
type releases struct {
	srv   *httptest.Server
	mu    sync.Mutex
	asked map[string]int
}

func serveReleases(t *testing.T, rels ...fakeRelease) *releases {
	t.Helper()
	f := &releases{asked: map[string]int{}}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /repos/o/r/releases", func(w http.ResponseWriter, _ *http.Request) {
		list := make([]ghRelease, 0, len(rels))
		for i, rel := range rels {
			gh := ghRelease{TagName: TagPrefix + rel.version, PublishedAt: time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC).Add(-time.Duration(i) * time.Hour)}
			for name := range rel.assets {
				gh.Assets = append(gh.Assets, ghAsset{Name: name, URL: f.srv.URL + "/assets/" + rel.version + "/" + name})
			}
			list = append(list, gh)
		}
		_ = json.NewEncoder(w).Encode(list)
	})
	mux.HandleFunc("GET /assets/{version}/{name}", func(w http.ResponseWriter, r *http.Request) {
		for _, rel := range rels {
			if data, ok := rel.assets[r.PathValue("name")]; ok && rel.version == r.PathValue("version") {
				_, _ = w.Write(data)
				return
			}
		}
		http.NotFound(w, r)
	})
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.asked[r.URL.Path]++
		f.mu.Unlock()
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *releases) count(path string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.asked[path]
}

func (f *releases) options(keys []*ecdsa.PublicKey) Options {
	return Options{Repo: "o/r", APIBase: f.srv.URL, OS: "linux", Arch: "amd64", Keys: keys}
}

// signed is a release of version whose archive holds content, its
// checksums.txt signed with the test key signer, and that key's public half
// published as cosign.pub.
func signed(t *testing.T, version, signer string, content []byte) fakeRelease {
	t.Helper()
	archive := tarball(t, "flai", content)
	sum := sha256.Sum256(archive)
	sums := []byte(hex.EncodeToString(sum[:]) + "  " + ArchiveName(version, "linux", "amd64") + "\n")
	return fakeRelease{version: version, assets: map[string][]byte{
		ArchiveName(version, "linux", "amd64"): archive,
		checksumsName:                          sums,
		signatureName:                          sign(t, privateKey(t, signer), sums),
		publicKeyName:                          testdataFile(t, signer+".pub"),
	}}
}

// without is rel with the named assets taken out.
func without(rel fakeRelease, names ...string) fakeRelease {
	for _, n := range names {
		delete(rel.assets, n)
	}
	return rel
}

// download resolves version from f and downloads it with keys.
func download(t *testing.T, f *releases, version string, keys []*ecdsa.PublicKey) ([]byte, error) {
	t.Helper()
	opt := f.options(keys)
	opt.Version = version
	rel, err := Resolve(context.Background(), opt)
	if err != nil {
		t.Fatal(err)
	}
	return Download(context.Background(), opt, rel)
}

// refusal is err as a *SignatureError of kind, failing the test otherwise.
func refusal(t *testing.T, err error, kind SignatureKind) *SignatureError {
	t.Helper()
	var se *SignatureError
	if !errors.As(err, &se) || se.Kind != kind {
		t.Fatalf("want a %s signature refusal, got %v", kind, err)
	}
	return se
}

func TestGoodSignatureInstalls(t *testing.T) {
	content := []byte("#!/bin/sh\necho signed\n")
	f := serveReleases(t, signed(t, "2.0.0", "trusted", content))
	bin, err := download(t, f, "2.0.0", trustedKeys(t))
	if err != nil || !bytes.Equal(bin, content) {
		t.Fatalf("a release signed with a trusted key installs: %q %v", bin, err)
	}
	if f.count("/assets/2.0.0/cosign.pub") != 0 {
		t.Error("a signature a trusted key verifies needs no cosign.pub")
	}
}

func TestBadSignatureIsRefused(t *testing.T) {
	altered := signed(t, "2.0.0", "trusted", []byte("bin"))
	altered.assets[checksumsName] = append(altered.assets[checksumsName], []byte("0000  flai_2.0.0_darwin_arm64.tar.gz\n")...)
	notBase64 := signed(t, "2.0.0", "trusted", []byte("bin"))
	notBase64.assets[signatureName] = []byte("not a signature!\n")
	noKey := without(signed(t, "2.0.0", "unknown", []byte("bin")), publicKeyName)
	notAKey := signed(t, "2.0.0", "unknown", []byte("bin"))
	notAKey.assets[publicKeyName] = []byte("not a key\n")
	for name, tc := range map[string]struct {
		rel    fakeRelease
		reason string
	}{
		"checksums.txt altered after signing": {altered, "neither a key this flai trusts nor the release's cosign.pub verifies it"},
		"a signature that is not base64":      {notBase64, "checksums.txt.sig is not base64"},
		"no cosign.pub to check against":      {noKey, "the release publishes no cosign.pub"},
		"a cosign.pub that is not a key":      {notAKey, "the release's cosign.pub cannot check it: it holds no PEM block"},
	} {
		t.Run(name, func(t *testing.T) {
			f := serveReleases(t, tc.rel)
			_, err := download(t, f, "2.0.0", trustedKeys(t))
			se := refusal(t, err, SignatureBad)
			msg := err.Error()
			if se.Tag != "flai/v2.0.0" || !strings.Contains(msg, "release flai/v2.0.0 is not installed") || !strings.Contains(msg, tc.reason) || !strings.Contains(msg, "may have been altered") {
				t.Errorf("names the release, the reason, and what it means: %v", err)
			}
			if f.count("/assets/2.0.0/flai_2.0.0_linux_amd64.tar.gz") != 0 {
				t.Error("the archive is not downloaded before the signature verifies")
			}
		})
	}
}

func TestMissingSignatureIsRefused(t *testing.T) {
	for _, asset := range []string{signatureName, checksumsName} {
		f := serveReleases(t, without(signed(t, "2.0.0", "trusted", []byte("bin")), asset))
		_, err := download(t, f, "2.0.0", trustedKeys(t))
		se := refusal(t, err, SignatureMissing)
		if se.Tag != "flai/v2.0.0" || !strings.Contains(err.Error(), "release flai/v2.0.0 is not installed: it publishes no "+asset+", so its signature cannot be verified") {
			t.Errorf("without %s: names the release and the reason: %v", asset, err)
		}
	}
}

func TestUnknownKeyNamesTheNewestVerifiableRelease(t *testing.T) {
	f := serveReleases(t,
		signed(t, "2.1.0", "unknown", []byte("bin")),
		signed(t, "2.0.0", "unknown", []byte("bin")),
		without(signed(t, "1.9.0", "trusted", []byte("bin")), signatureName),
		signed(t, "1.8.0", "trusted", []byte("bin")),
		signed(t, "1.7.0", "trusted", []byte("bin")),
	)
	_, err := download(t, f, "2.1.0", trustedKeys(t))
	se := refusal(t, err, SignatureUnknownKey)
	if se.Tag != "flai/v2.1.0" || se.Bridge != "1.8.0" {
		t.Errorf("the refused release and the newest one a trusted key signed: %+v", se)
	}
	if !strings.Contains(err.Error(), "release flai/v2.1.0 is not installed") || !strings.Contains(err.Error(), "install flai 1.8.0 first") || !strings.Contains(err.Error(), "flai self-upgrade --version 1.8.0") {
		t.Errorf("names the release to install first and how: %v", err)
	}
	if f.count("/assets/2.0.0/cosign.pub") != 0 || f.count("/assets/1.7.0/checksums.txt") != 0 {
		t.Error("the search reads only each candidate's checksums.txt and signature, and stops at the first that verifies")
	}

	err = VerifyRelease(context.Background(), f.options(trustedKeys(t)), releaseOf(t, f, "2.1.0"))
	if se := refusal(t, err, SignatureUnknownKey); se.Bridge != "" || !strings.Contains(err.Error(), "install the newest release this flai can verify first") {
		t.Errorf("the cheap check does not search for a release to install first: %v", err)
	}
}

func TestUnknownKeyWithNoVerifiableRelease(t *testing.T) {
	f := serveReleases(t,
		signed(t, "2.0.0", "unknown", []byte("bin")),
		without(signed(t, "1.0.0", "trusted", []byte("bin")), signatureName),
	)
	_, err := download(t, f, "2.0.0", trustedKeys(t))
	se := refusal(t, err, SignatureUnknownKey)
	if se.Bridge != "" || !strings.Contains(err.Error(), "no published release is signed by a key this flai trusts") {
		t.Errorf("says that no release can be verified: %v", err)
	}
}

func TestDefaultKeysAreTheEmbeddedReleaseKey(t *testing.T) {
	f := serveReleases(t, signed(t, "2.0.0", "trusted", []byte("bin")))
	_, err := download(t, f, "2.0.0", nil)
	if se := refusal(t, err, SignatureUnknownKey); se.Tag != "flai/v2.0.0" {
		t.Errorf("without Options.Keys only the embedded release key is trusted, not the test key: %+v", se)
	}
}

func TestVerifyRelease(t *testing.T) {
	f := serveReleases(t,
		signed(t, "2.0.0", "trusted", []byte("bin")),
		without(signed(t, "1.0.0", "trusted", []byte("bin")), signatureName),
	)
	opt := f.options(trustedKeys(t))
	if err := VerifyRelease(context.Background(), opt, releaseOf(t, f, "2.0.0")); err != nil {
		t.Errorf("a release a trusted key signed verifies: %v", err)
	}
	if f.count("/assets/2.0.0/flai_2.0.0_linux_amd64.tar.gz") != 0 {
		t.Error("the check downloads no archive")
	}
	if se := refusal(t, VerifyRelease(context.Background(), opt, releaseOf(t, f, "1.0.0")), SignatureMissing); se.Tag != "flai/v1.0.0" {
		t.Errorf("a release without a signature is refused by name: %+v", se)
	}
}

// releaseOf is version as f lists it.
func releaseOf(t *testing.T, f *releases, version string) *Release {
	t.Helper()
	opt := f.options(nil)
	opt.Version = version
	rel, err := Resolve(context.Background(), opt)
	if err != nil {
		t.Fatal(err)
	}
	return rel
}

func TestSignatureDownloadFailureIsNotARefusal(t *testing.T) {
	pauseFor(t, 0)
	f := serveReleases(t, signed(t, "2.0.0", "trusted", []byte("bin")))
	rel := releaseOf(t, f, "2.0.0")
	for i := range rel.Assets {
		if rel.Assets[i].Name == signatureName {
			rel.Assets[i].URL = f.srv.URL + "/assets/2.0.0/gone"
		}
	}
	err := VerifyRelease(context.Background(), f.options(trustedKeys(t)), rel)
	var se *SignatureError
	if err == nil || errors.As(err, &se) || !strings.Contains(err.Error(), "download checksums.txt.sig of release flai/v2.0.0") {
		t.Errorf("a failed download is named as one, not as a refusal: %v", err)
	}
}
