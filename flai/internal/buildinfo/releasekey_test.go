package buildinfo

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
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
