package selfupgrade

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

// The release assets a signature check reads: the checksums file, cosign's
// signature over it, and the release key's public half (S-0232).
const (
	checksumsName = "checksums.txt"
	signatureName = "checksums.txt.sig"
	publicKeyName = "cosign.pub"
)

// SignatureKind says why a release's signature was refused.
type SignatureKind string

// The reasons a release's signature is refused.
const (
	// SignatureMissing: the release publishes no checksums.txt.sig, or no
	// checksums.txt, such as a release from before releases were signed.
	SignatureMissing SignatureKind = "missing"
	// SignatureBad: the signature verifies against no key this binary trusts
	// nor the release's own cosign.pub; the release may have been altered.
	SignatureBad SignatureKind = "bad"
	// SignatureUnknownKey: the signature verifies against the release's own
	// cosign.pub, a key this binary does not trust.
	SignatureUnknownKey SignatureKind = "unknown-key"
)

// SignatureError refuses a release whose checksums.txt this binary cannot
// verify.
type SignatureError struct {
	Tag    string        // the refused release
	Kind   SignatureKind // why it is refused
	Reason string        // what the check found
	// Bridge is, for SignatureUnknownKey, the newest published version this
	// binary can verify, to install first; "" when none or not looked for.
	Bridge   string
	searched bool // the published releases were searched for a bridge
}

// Error names the release, why it is refused, and what to install instead.
func (e *SignatureError) Error() string {
	switch e.Kind {
	case SignatureMissing:
		return fmt.Sprintf("release %s is not installed: %s, so its signature cannot be verified (a release from before signing has none); install a signed release, such as the latest", e.Tag, e.Reason)
	case SignatureUnknownKey:
		head := fmt.Sprintf("release %s is not installed: %s, which this flai does not trust (the release key may have changed since this flai was built)", e.Tag, e.Reason)
		switch {
		case e.Bridge != "":
			return fmt.Sprintf("%s; install flai %s first, the newest release this flai can verify (flai self-upgrade --version %s), then upgrade again", head, e.Bridge, e.Bridge)
		case e.searched:
			return head + "; no published release is signed by a key this flai trusts: install with install.sh, and compare the release key's fingerprint with docs/operators/runbooks/release-key.md"
		}
		return head + "; install the newest release this flai can verify first, then upgrade again"
	}
	return fmt.Sprintf("release %s is not installed: the signature over its checksums.txt does not verify (%s); the release may have been altered since it was signed: do not install it by other means, and report it to the repository's maintainers", e.Tag, e.Reason)
}

// trustedKeys are the keys a release's signature must verify against:
// Options.Keys, else the release key embedded in the binary.
func (o Options) trustedKeys() ([]*ecdsa.PublicKey, error) {
	if len(o.Keys) > 0 {
		return o.Keys, nil
	}
	return buildinfo.ReleaseKeys()
}

// VerifyRelease checks that rel's checksums.txt is signed by a key this binary
// trusts, downloading checksums.txt, its signature, and, when no trusted key
// verifies it, the release's cosign.pub, but never an archive. It answers nil,
// a *SignatureError, or the failure to download or to read the trusted keys.
func VerifyRelease(ctx context.Context, opt Options, rel *Release) error {
	opt.defaults()
	keys, err := opt.trustedKeys()
	if err != nil {
		return fmt.Errorf("verify the signature of release %s: %w", rel.Tag, err)
	}
	_, err = opt.verifiedSums(ctx, rel, keys)
	return err
}

// verifiedSums answers rel's checksums.txt once its signature verifies against
// one of keys, else a *SignatureError saying why not.
func (o Options) verifiedSums(ctx context.Context, rel *Release, keys []*ecdsa.PublicKey) ([]byte, error) {
	sums, sig, err := o.signedSums(ctx, rel)
	if err != nil {
		return nil, err
	}
	if signedByOneOf(sums, sig, keys) {
		return sums, nil
	}
	pubAsset, ok := rel.Asset(publicKeyName)
	if !ok {
		return nil, &SignatureError{Tag: rel.Tag, Kind: SignatureBad, Reason: "no key this flai trusts verifies it, and the release publishes no " + publicKeyName}
	}
	pubPEM, err := o.getBytes(ctx, pubAsset.URL)
	if err != nil {
		return nil, fmt.Errorf("download %s of release %s: %w", publicKeyName, rel.Tag, err)
	}
	own, err := parsePublicKey(pubPEM)
	if err != nil {
		return nil, &SignatureError{Tag: rel.Tag, Kind: SignatureBad, Reason: fmt.Sprintf("no key this flai trusts verifies it, and the release's %s cannot check it: %v", publicKeyName, err)}
	}
	if signedByOneOf(sums, sig, []*ecdsa.PublicKey{own}) {
		return nil, &SignatureError{Tag: rel.Tag, Kind: SignatureUnknownKey, Reason: "it is signed by the key it publishes as " + publicKeyName}
	}
	return nil, &SignatureError{Tag: rel.Tag, Kind: SignatureBad, Reason: "neither a key this flai trusts nor the release's " + publicKeyName + " verifies it"}
}

// signedSums downloads rel's checksums.txt and the signature over it, decoded
// from the base64 cosign writes.
func (o Options) signedSums(ctx context.Context, rel *Release) (sums, sig []byte, err error) {
	sumsAsset, ok := rel.Asset(checksumsName)
	if !ok {
		return nil, nil, &SignatureError{Tag: rel.Tag, Kind: SignatureMissing, Reason: "it publishes no " + checksumsName}
	}
	sigAsset, ok := rel.Asset(signatureName)
	if !ok {
		return nil, nil, &SignatureError{Tag: rel.Tag, Kind: SignatureMissing, Reason: "it publishes no " + signatureName}
	}
	sums, err = o.getBytes(ctx, sumsAsset.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("download %s of release %s: %w", checksumsName, rel.Tag, err)
	}
	encoded, err := o.getBytes(ctx, sigAsset.URL)
	if err != nil {
		return nil, nil, fmt.Errorf("download %s of release %s: %w", signatureName, rel.Tag, err)
	}
	sig, err = base64.StdEncoding.DecodeString(strings.TrimSpace(string(encoded)))
	if err != nil {
		return nil, nil, &SignatureError{Tag: rel.Tag, Kind: SignatureBad, Reason: signatureName + " is not base64"}
	}
	return sums, sig, nil
}

// withBridge adds to refused, an unknown-key refusal, the newest published
// release whose signature verifies against one of keys, downloading only the
// checksums.txt and the signature of each candidate.
func (o Options) withBridge(ctx context.Context, refused *SignatureError, keys []*ecdsa.PublicKey) error {
	published, err := List(ctx, o, TagPrefix)
	if err != nil {
		return fmt.Errorf("%w; looking for a release this flai can verify failed: %w", refused, err)
	}
	refused.searched = true
	for i := range published {
		rel := &published[i]
		if rel.Tag == refused.Tag {
			continue
		}
		sums, sig, err := o.signedSums(ctx, rel)
		var unverifiable *SignatureError
		switch {
		case errors.As(err, &unverifiable):
			continue
		case err != nil:
			return fmt.Errorf("%w; looking for a release this flai can verify failed: %w", refused, err)
		}
		if signedByOneOf(sums, sig, keys) {
			refused.Bridge = rel.Version
			return refused
		}
	}
	return refused
}

// signedByOneOf reports whether sig, an ASN.1 ECDSA signature over the
// SHA-256 of data, verifies against one of keys.
func signedByOneOf(data, sig []byte, keys []*ecdsa.PublicKey) bool {
	sum := sha256.Sum256(data)
	for _, k := range keys {
		if ecdsa.VerifyASN1(k, sum[:], sig) {
			return true
		}
	}
	return false
}

// parsePublicKey reads an ECDSA public key from a PKIX PEM block, as
// cosign.pub holds it.
func parsePublicKey(data []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("it holds no PEM block")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ec, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("it is %T, not ECDSA", key)
	}
	return ec, nil
}
