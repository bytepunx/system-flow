package buildinfo

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"testing"
)

func testKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// sign signs statement as cosign sign-blob does: ECDSA over its SHA-256, in
// ASN.1 DER, base64.
func sign(t *testing.T, k *ecdsa.PrivateKey, statement string) string {
	t.Helper()
	sum := sha256.Sum256([]byte(statement))
	sig, err := ecdsa.SignASN1(rand.Reader, k, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(sig)
}

func TestVerifyStamp(t *testing.T) {
	release, other := testKey(t), testKey(t)
	keys := []*ecdsa.PublicKey{&release.PublicKey}
	const version, commit = "1.42.0", "0ef6d26"
	good := "flai 1.42.0 0ef6d26"
	for _, c := range []struct {
		name, statement, signature string
		state                      StampState
		reason                     string
	}{
		{"signed by the release key", good, sign(t, release, good), StampSigned, ""},
		{"missing", "", "", StampMissing, ""},
		{"signed by another key", good, sign(t, other, good), StampInvalid, "does not verify"},
		{"another statement's signature", good, sign(t, release, "flai 1.42.0 1234567"), StampInvalid, "does not verify"},
		{"statement only", good, "", StampInvalid, "without a signature"},
		{"signature only", "", sign(t, release, good), StampInvalid, "without a statement"},
		{"signature not base64", good, "not base64!", StampInvalid, "not base64"},
		{"wrong component", "flaiover 1.42.0 0ef6d26", sign(t, release, "flaiover 1.42.0 0ef6d26"), StampInvalid, "names flaiover"},
		{"another version", "flai 1.41.0 0ef6d26", sign(t, release, "flai 1.41.0 0ef6d26"), StampInvalid, "version 1.41.0"},
		{"another commit", "flai 1.42.0 1234567", sign(t, release, "flai 1.42.0 1234567"), StampInvalid, "commit 1234567"},
		{"malformed statement", "flai 1.42.0 0ef6d26\n", sign(t, release, "flai 1.42.0 0ef6d26\n"), StampInvalid, "is not"},
		{"too many fields", "flai 1.42.0 0ef6d26 x", sign(t, release, "flai 1.42.0 0ef6d26 x"), StampInvalid, "is not"},
	} {
		t.Run(c.name, func(t *testing.T) {
			s := VerifyStamp(c.statement, c.signature, version, commit, keys)
			if s.State != c.state || !strings.Contains(s.Reason, c.reason) || (c.reason == "") != (s.Reason == "") {
				t.Errorf("state %s reason %q, want %s with %q", s.State, s.Reason, c.state, c.reason)
			}
			if s.Statement != c.statement || s.Signature != c.signature || s.Development {
				t.Errorf("stamp %+v does not carry what it verified", s)
			}
		})
	}
}

// A signature is tried against every key the binary knows.
func TestVerifyStampTriesEveryKey(t *testing.T) {
	first, second := testKey(t), testKey(t)
	statement := "flai 1.42.0 0ef6d26"
	s := VerifyStamp(statement, sign(t, second, statement), "1.42.0", "0ef6d26", []*ecdsa.PublicKey{&first.PublicKey, &second.PublicKey})
	if s.State != StampSigned {
		t.Errorf("signed by the second key: %+v", s)
	}
}

func TestDevelopment(t *testing.T) {
	for v, want := range map[string]bool{
		"dev": true, "": true, "0.0.0-0ef6d26": true, "0.0.0": true,
		"1.42.0": false, "0.39.0": false,
	} {
		if got := Development(v); got != want {
			t.Errorf("Development(%q) = %v, want %v", v, got, want)
		}
		if got := VerifyStamp("", "", v, "", nil); got.Development != want || got.State != StampMissing {
			t.Errorf("VerifyStamp for version %q: %+v", v, got)
		}
	}
}

// A test build carries no stamp, and the embedded release key parses.
func TestOwnStampOfATestBuildIsMissing(t *testing.T) {
	keys, err := ReleaseKeys()
	if err != nil || len(keys) != 1 {
		t.Fatalf("release keys: %v %v", keys, err)
	}
	if s := OwnStamp(); s.State != StampMissing || !s.Development {
		t.Errorf("own stamp of a test build: %+v", s)
	}
}
