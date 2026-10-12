package buildinfo

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"
)

// A release stamp is the statement "flai <version> <commit>" and its
// signature by the release key, which the release workflow puts in the
// binary through the linker (design/system/release-signing.md, S-0235).

// StampState says whether a build carries a valid release stamp.
type StampState string

// The states of a release stamp.
const (
	StampSigned  StampState = "signed"
	StampMissing StampState = "missing"
	StampInvalid StampState = "invalid"
)

// stampComponent is the component a flai release stamp names.
const stampComponent = "flai"

// Stamp is a build's release stamp and what verifying it found.
type Stamp struct {
	Statement string     `json:"statement"`
	Signature string     `json:"signature"`
	State     StampState `json:"state"`
	// Reason says why an invalid stamp is not valid.
	Reason string `json:"reason,omitempty"`
	// Development is true for a build that claims no release version,
	// which carries no stamp and is never warned about.
	Development bool `json:"development"`
}

// Development reports whether version is a development build's: dev, empty,
// or 0.0.0 and a suffix, as scripts/flai.sh and goreleaser --snapshot build.
func Development(version string) bool {
	return version == "" || version == "dev" || strings.HasPrefix(version, "0.0.0")
}

// ReleaseKeys returns the public keys a release stamp is verified with.
func ReleaseKeys() ([]*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(ReleaseKeyPEM()))
	if block == nil {
		return nil, fmt.Errorf("the embedded release key holds no PEM block; rebuild flai from a tree whose flai/internal/buildinfo/cosign.pub is the release key")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse the embedded release key: %w; rebuild flai from a tree whose flai/internal/buildinfo/cosign.pub is the release key", err)
	}
	ec, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("the embedded release key is %T, not ECDSA; rebuild flai from a tree whose flai/internal/buildinfo/cosign.pub is the release key", key)
	}
	return []*ecdsa.PublicKey{ec}, nil
}

// OwnStamp verifies this binary's release stamp with the embedded release key.
func OwnStamp() Stamp {
	keys, err := ReleaseKeys()
	s := VerifyStamp(ReleaseStatement, ReleaseSignature, Version, Commit, keys)
	if err != nil && s.State == StampInvalid {
		s.Reason = err.Error()
	}
	return s
}

// VerifyStamp checks that signature signs statement by one of keys, and that
// statement names flai, version, and commit.
func VerifyStamp(statement, signature, version, commit string, keys []*ecdsa.PublicKey) Stamp {
	s := Stamp{Statement: statement, Signature: signature, Development: Development(version)}
	s.State, s.Reason = verify(statement, signature, version, commit, keys)
	return s
}

func verify(statement, signature, version, commit string, keys []*ecdsa.PublicKey) (StampState, string) {
	switch {
	case statement == "" && signature == "":
		return StampMissing, ""
	case signature == "":
		return StampInvalid, "a statement without a signature"
	case statement == "":
		return StampInvalid, "a signature without a statement"
	}
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return StampInvalid, "the signature is not base64"
	}
	if !signedByOneOf(statement, sig, keys) {
		return StampInvalid, "the signature does not verify against the release key"
	}
	fields := strings.Fields(statement)
	if len(fields) != 3 || strings.Join(fields, " ") != statement {
		return StampInvalid, fmt.Sprintf("the statement %q is not %q", statement, "flai <version> <commit>")
	}
	switch {
	case fields[0] != stampComponent:
		return StampInvalid, fmt.Sprintf("the statement names %s, not %s", fields[0], stampComponent)
	case fields[1] != version:
		return StampInvalid, fmt.Sprintf("the statement names version %s, the build is %s", fields[1], version)
	case fields[2] != commit:
		return StampInvalid, fmt.Sprintf("the statement names commit %s, the build is %s", fields[2], commit)
	}
	return StampSigned, ""
}

// PeerFault says why a peer's release stamp is not signed: the words a
// connection refused with close code 4403 gives after "unsigned peer: "
// (S-0237).
type PeerFault string

// The faults of a peer's release stamp.
const (
	PeerNoStamp        PeerFault = "missing stamp"
	PeerBadSignature   PeerFault = "invalid signature"
	PeerWrongComponent PeerFault = "wrong component"
)

// PeerStamp is what checking the release stamp a peer sent found.
type PeerStamp struct {
	State StampState `json:"state"`
	// Fault says why a stamp that is not signed is not; empty when it is.
	Fault PeerFault `json:"fault,omitempty"`
	// Reason says in full why, for a log.
	Reason string `json:"reason,omitempty"`
	// Component, Version, and Commit are what a well-formed statement names,
	// whether or not its signature verifies: trust them only when State is
	// signed.
	Component string `json:"component,omitempty"`
	Version   string `json:"version,omitempty"`
	Commit    string `json:"commit,omitempty"`
}

// VerifyPeerStamp checks a peer's release stamp: that signature signs
// statement by one of keys, and that statement names component. Unlike
// VerifyStamp it does not compare the version or the commit with this
// build's, since the peer is another release. A stamp that is not well
// formed counts as an invalid signature.
func VerifyPeerStamp(statement, signature, component string, keys []*ecdsa.PublicKey) PeerStamp {
	var p PeerStamp
	if fields := strings.Fields(statement); len(fields) == 3 && strings.Join(fields, " ") == statement {
		p.Component, p.Version, p.Commit = fields[0], fields[1], fields[2]
	}
	invalid := func(reason string) PeerStamp {
		p.State, p.Fault, p.Reason = StampInvalid, PeerBadSignature, reason
		return p
	}
	switch {
	case statement == "" && signature == "":
		p.State, p.Fault, p.Reason = StampMissing, PeerNoStamp, "the peer sent no release stamp"
		return p
	case signature == "":
		return invalid("a statement without a signature")
	case statement == "":
		return invalid("a signature without a statement")
	}
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return invalid("the signature is not base64")
	}
	if !signedByOneOf(statement, sig, keys) {
		return invalid("the signature does not verify against the release key")
	}
	if p.Component == "" {
		return invalid(fmt.Sprintf("the statement %q is not %q", statement, component+" <version> <commit>"))
	}
	if p.Component != component {
		p.State, p.Fault, p.Reason = StampInvalid, PeerWrongComponent, fmt.Sprintf("the statement names %s, not %s", p.Component, component)
		return p
	}
	p.State = StampSigned
	return p
}

func signedByOneOf(statement string, sig []byte, keys []*ecdsa.PublicKey) bool {
	sum := sha256.Sum256([]byte(statement))
	for _, k := range keys {
		if ecdsa.VerifyASN1(k, sum[:], sig) {
			return true
		}
	}
	return false
}
