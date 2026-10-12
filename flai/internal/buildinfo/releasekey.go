package buildinfo

import _ "embed"

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
