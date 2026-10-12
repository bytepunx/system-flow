package buildinfo

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"os"
	"strings"
	"testing"
)

// The fingerprint docs/operators/runbooks/release-key.md records: the SHA-256
// of the key's DER form, as openssl pkey -pubin -outform DER | openssl dgst
// -sha256 prints it.
const releaseKeyFingerprint = "ee6df72313c9dbac1cf2270ed7f01cc523786d37bfce2dec12210215d6c5da40"

func TestReleaseKeyIsTheOperatorsP256PublicKey(t *testing.T) {
	block, rest := pem.Decode([]byte(ReleaseKeyPEM()))
	if block == nil {
		t.Fatalf("cosign.pub holds no PEM block: %q", ReleaseKeyPEM())
	}
	if block.Type != "PUBLIC KEY" {
		t.Errorf("PEM block type = %q, want PUBLIC KEY", block.Type)
	}
	if len(rest) != 0 {
		t.Errorf("cosign.pub has %d bytes after its PEM block", len(rest))
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		t.Fatalf("parse the release key: %v", err)
	}
	ec, ok := key.(*ecdsa.PublicKey)
	if !ok {
		t.Fatalf("release key is %T, want *ecdsa.PublicKey", key)
	}
	if ec.Curve != elliptic.P256() {
		t.Errorf("release key curve = %s, want P-256", ec.Curve.Params().Name)
	}
	sum := sha256.Sum256(block.Bytes)
	if got := hex.EncodeToString(sum[:]); got != releaseKeyFingerprint {
		t.Errorf("release key fingerprint = %s, want %s, as the release-key runbook records", got, releaseKeyFingerprint)
	}
}

// A build that no linker flag gave a test key, as every release build is,
// trusts no key beside the release key.
func TestABuildTrustsNoTestKeyUnlessLinkedWithOne(t *testing.T) {
	keys, err := LinkedTestReleaseKeys()
	if err != nil || keys != nil {
		t.Errorf("LinkedTestReleaseKeys() = %v, %v; want no key and no error", keys, err)
	}
}

// No release build links a test key: GoReleaser builds every release, and
// its ldflags do not set testReleaseKey.
func TestReleaseBuildsLinkNoTestKey(t *testing.T) {
	b, err := os.ReadFile("../../.goreleaser.yaml")
	if err != nil {
		t.Fatalf("read flai/.goreleaser.yaml: %v", err)
	}
	config := string(b)
	if !strings.Contains(config, "-X github.com/bytepunx/system-flow/flai/internal/buildinfo.Version=") {
		t.Fatal("flai/.goreleaser.yaml sets no buildinfo.Version through its ldflags; is it still the file that builds releases?")
	}
	if strings.Contains(config, "testReleaseKey") {
		t.Error("flai/.goreleaser.yaml mentions testReleaseKey; a release build must trust the release key only")
	}
}

func TestALinkedTestKeyParsesOrSaysWhichFlag(t *testing.T) {
	t.Cleanup(func() { testReleaseKey = "" })
	block, _ := pem.Decode([]byte(ReleaseKeyPEM()))
	testReleaseKey = base64.StdEncoding.EncodeToString(block.Bytes)
	keys, err := LinkedTestReleaseKeys()
	if err != nil || len(keys) != 1 {
		t.Fatalf("LinkedTestReleaseKeys() with a key = %v, %v; want the key", keys, err)
	}
	want, err := ReleaseKeys()
	if err != nil {
		t.Fatal(err)
	}
	if !keys[0].Equal(want[0]) {
		t.Error("the linked key parsed to another key")
	}
	for _, bad := range []string{"not base64!", base64.StdEncoding.EncodeToString([]byte("not a key"))} {
		testReleaseKey = bad
		if _, err := LinkedTestReleaseKeys(); err == nil || !strings.Contains(err.Error(), testReleaseKeyFlag) {
			t.Errorf("LinkedTestReleaseKeys() with %q answered %v, want an error naming %s", bad, err, testReleaseKeyFlag)
		}
	}
}
