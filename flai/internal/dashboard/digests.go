// Package dashboard reads the signed digest list a flaiover release publishes
// and keeps each verified list beside flai serve's state, so that flai
// dashboard runs a release's image by digest and a dial trusts only an image
// a verified list names (ADR-0070, design/system/release-signing.md § The
// digest list).
package dashboard

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

// IndexKind is the kind of the digest list's first line, the image index.
const IndexKind = "index"

var (
	digestRe = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	refRe    = regexp.MustCompile(`^([A-Za-z0-9._/:-]+):([^:]+)$`)
	kindRe   = regexp.MustCompile(`^[a-z0-9]+(/[a-z0-9]+){1,2}$`)
)

// Entry is one line of a digest list: a digest, the reference it was pushed
// under, and what it is.
type Entry struct {
	Digest string // sha256: and 64 lower-case hex characters
	Ref    string // <repository>:<version>
	Kind   string // index, or <os>/<arch>[/<variant>]
}

// List is a parsed digest list: the release it is for and its entries, the
// index first.
type List struct {
	Version    string // the release version, X.Y.Z without its v
	Repository string // the image repository, the reference without :<version>
	Entries    []Entry
}

// ListName is the release asset name of version's digest list; its
// signature is ListName(version)+".sig".
func ListName(version string) string {
	return "flaiover_" + version + ".digests"
}

// Index is the image index's digest, the one flai dashboard pulls by.
func (l List) Index() string {
	if len(l.Entries) == 0 {
		return ""
	}
	return l.Entries[0].Digest
}

// Reference is the image to pull: the repository at the index digest.
func (l List) Reference() string {
	return l.Repository + "@" + l.Index()
}

// Digests is every digest the list names, in its order.
func (l List) Digests() []string {
	out := make([]string, len(l.Entries))
	for i, e := range l.Entries {
		out[i] = e.Digest
	}
	return out
}

// Has reports whether the list names digest, whatever its kind.
func (l List) Has(digest string) bool {
	for _, e := range l.Entries {
		if e.Digest == digest {
			return true
		}
	}
	return false
}

// Matches reports whether repoDigest, an entry of an image's RepoDigests
// (<repository>@sha256:…), is the list's repository at a digest it names.
func (l List) Matches(repoDigest string) bool {
	repo, digest, ok := strings.Cut(repoDigest, "@")
	return ok && repo == l.Repository && l.Has(digest)
}

// Parse reads a digest list in the format release-flaiover.yml writes,
// refusing any other shape with an error naming the line.
func Parse(data []byte) (List, error) {
	if len(data) == 0 {
		return List{}, malformed(0, "it is empty")
	}
	body, ok := strings.CutSuffix(string(data), "\n")
	if !ok {
		return List{}, malformed(strings.Count(string(data), "\n")+1, "it does not end in a newline")
	}
	var l List
	kindLine := map[string]int{}
	for i, line := range strings.Split(body, "\n") {
		n := i + 1
		e, err := parseLine(line)
		if err != nil {
			return List{}, malformed(n, err.Error())
		}
		repo, version := splitRef(e.Ref)
		if n == 1 {
			l.Repository, l.Version = repo, version
		} else if repo != l.Repository || version != l.Version {
			return List{}, malformed(n, fmt.Sprintf("its reference %s is not line 1's %s:%s", e.Ref, l.Repository, l.Version))
		}
		if err := checkKind(n, e.Kind, l.Entries, kindLine); err != nil {
			return List{}, malformed(n, err.Error())
		}
		kindLine[e.Kind] = n
		l.Entries = append(l.Entries, e)
	}
	return l, nil
}

// Verify checks sig, the list's signature as its release publishes it,
// against keys, then parses the list; a list it does not verify is refused.
func Verify(data, sig []byte, keys []*ecdsa.PublicKey) (List, error) {
	if !selfupgrade.VerifyBlob(data, sig, keys) {
		return List{}, errors.New("the flaiover digest list's signature does not verify against a key this flai trusts; it may have been altered, so fetch it again from its flaiover release and run no image it names until it verifies")
	}
	return Parse(data)
}

// parseLine splits a line into its three fields and checks each one's shape.
func parseLine(line string) (Entry, error) {
	switch {
	case line == "":
		return Entry{}, errors.New("it is blank")
	case strings.Contains(line, "\r"):
		return Entry{}, errors.New("it holds a carriage return; the list ends lines with LF only")
	case !utf8.ValidString(line):
		return Entry{}, errors.New("it is not UTF-8")
	}
	fields := strings.Split(line, " ")
	if len(fields) != 3 {
		return Entry{}, fmt.Errorf("it has %d fields separated by spaces, not 3 (digest, reference, kind) separated by one space each", len(fields))
	}
	e := Entry{Digest: fields[0], Ref: fields[1], Kind: fields[2]}
	if !digestRe.MatchString(e.Digest) {
		return Entry{}, fmt.Errorf("its digest %q is not sha256: and 64 lower-case hex characters", e.Digest)
	}
	m := refRe.FindStringSubmatch(e.Ref)
	if m == nil {
		return Entry{}, fmt.Errorf("its reference %q is not <repository>:<version>", e.Ref)
	}
	if !buildinfo.Bare(m[2]) {
		return Entry{}, fmt.Errorf("its reference %q has version %q, not a release X.Y.Z without its v", e.Ref, m[2])
	}
	return e, nil
}

// splitRef is a checked reference's repository and version.
func splitRef(ref string) (repo, version string) {
	m := refRe.FindStringSubmatch(ref)
	return m[1], m[2]
}

// checkKind checks that line n's kind follows the entries before it: the
// index first, then platforms sorted in byte order, none twice.
func checkKind(n int, kind string, before []Entry, kindLine map[string]int) error {
	switch {
	case n == 1 && kind != IndexKind:
		return fmt.Errorf("its kind is %q, but the first line is the index", kind)
	case n == 1:
		return nil
	case kind == IndexKind:
		return errors.New("its kind is index, which only the first line is")
	case !kindRe.MatchString(kind):
		return fmt.Errorf("its kind %q is not index or <os>/<arch>[/<variant>]", kind)
	}
	if prev, ok := kindLine[kind]; ok {
		return fmt.Errorf("its kind %s is line %d's too", kind, prev)
	}
	if last := before[len(before)-1].Kind; n > 2 && kind < last {
		return fmt.Errorf("its kind %s sorts before line %d's %s; platform lines are sorted by kind in byte order", kind, n-1, last)
	}
	return nil
}

// malformed refuses a digest list at line n, or the whole list when n is 0.
func malformed(n int, reason string) error {
	at := "it"
	if n > 0 {
		at = fmt.Sprintf("line %d", n)
	}
	return fmt.Errorf("reading the flaiover digest list failed at %s: %s; it is not in the format release-flaiover.yml writes, so fetch it again from its flaiover release and run no image it names", at, reason)
}
