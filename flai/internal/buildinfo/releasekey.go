package buildinfo

import (
	"crypto/ecdsa"
	"crypto/x509"
	_ "embed"
	"encoding/base64"
	"fmt"
	"strings"
)

// The release key's public half, as PEM (ADR-0070, S-0232): cosign.pub, which
// `cosign generate-key-pair` wrote, ECDSA P-256 in a PKIX "PUBLIC KEY" block.
// GoReleaser publishes the same file as a release asset, and flaiover's
// src/lib/server/release.ts carries the same PEM; its test keeps them equal.
// docs/operators/runbooks/release-key.md records its fingerprint.
//
//go:embed cosign.pub
var releaseKeyPEM string

// ReleaseKeyPEM is the public key that verifies a release's signature, as PEM.
func ReleaseKeyPEM() string { return releaseKeyPEM }

// testReleaseKey is one more public key a test build of flai trusts to verify
// a release's signature, beside the release key: the base64 of its PKIX DER
// form, which is a PEM "PUBLIC KEY" block's body on one line, since a -X value
// cannot hold a newline. Only the linker sets it:
//
//	go build -ldflags "-X github.com/bytepunx/system-flow/flai/internal/buildinfo.testReleaseKey=<base64>"
//
// No release build sets it: GoReleaser's ldflags in flai/.goreleaser.yaml do
// not, and a test keeps it so. Only scripts/install-test.sh does, for the flai
// that upgrades from the local releases it signs with a key of its own
// (S-0233). It verifies releases only, never a release stamp.
var testReleaseKey string

// testReleaseKeyFlag is the linker flag that sets testReleaseKey.
const testReleaseKeyFlag = "-X github.com/bytepunx/system-flow/flai/internal/buildinfo.testReleaseKey"

// LinkedTestReleaseKeys returns the test key this build was linked with, to
// trust beside the release key, or none, as in every release build. A key that
// does not parse is an error naming the flag that set it.
func LinkedTestReleaseKeys() ([]*ecdsa.PublicKey, error) {
	if testReleaseKey == "" {
		return nil, nil
	}
	der, err := base64.StdEncoding.DecodeString(strings.TrimSpace(testReleaseKey))
	if err != nil {
		return nil, fmt.Errorf("the test release key this flai was built with (%s) is not base64: %w; give the body of a PEM PUBLIC KEY on one line", testReleaseKeyFlag, err)
	}
	key, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("parse the test release key this flai was built with (%s): %w; give the body of a PEM PUBLIC KEY on one line", testReleaseKeyFlag, err)
	}
	ec, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("the test release key this flai was built with (%s) is %T, not ECDSA", testReleaseKeyFlag, key)
	}
	return []*ecdsa.PublicKey{ec}, nil
}
