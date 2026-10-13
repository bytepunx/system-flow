package selfupgrade

import (
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
)

// verdictsFormat is the format of verdicts.json; a file of another format is
// read as empty.
const verdictsFormat = 1

// verdictsPath is the verdict cache's file under flai's cache_dir.
func verdictsPath(cacheDir string) string {
	return filepath.Join(cacheDir, "releases", "verdicts.json")
}

// verdictKey identifies one signature check. A replaced asset gets a new API
// URL, and a build that trusts other keys a new fingerprint, so either misses.
type verdictKey struct {
	Tag          string `json:"tag"`
	ChecksumsURL string `json:"checksums_url"`
	SignatureURL string `json:"signature_url"`
	PublicKeyURL string `json:"public_key_url,omitempty"` // "" when the release publishes no cosign.pub
	Keys         string `json:"keys"`                     // KeysFingerprint of the trusted keys
}

// verdict is a check's outcome: verified when kind is "", else the refusal.
type verdict struct {
	kind   SignatureKind
	reason string
}

// err answers the verdict as VerifyRelease answered it for the release tag.
func (v verdict) err(tag string) error {
	if v.kind == "" {
		return nil
	}
	return &SignatureError{Tag: tag, Kind: v.kind, Reason: v.reason}
}

// verdictEntry is one verdict as verdicts.json holds it.
type verdictEntry struct {
	verdictKey
	Verified bool          `json:"verified"`
	Kind     SignatureKind `json:"kind,omitempty"`
	Reason   string        `json:"reason,omitempty"`
}

// verdictsFile is the whole of verdicts.json.
type verdictsFile struct {
	Format   int            `json:"format"`
	Verdicts []verdictEntry `json:"verdicts"`
}

// Verdicts keeps each release's signature verdict, so that a release checked
// once is not downloaded again until its assets or the trusted keys change.
// It is safe for concurrent use.
type Verdicts struct {
	path  string // "" stores nothing on disk
	mu    sync.Mutex
	known map[verdictKey]verdict
	dirty bool // a verdict was stored since the file was read or written
}

// OpenVerdicts loads the verdicts kept under cacheDir, flai's cache_dir; a
// missing, unreadable, or corrupt file loads as empty. With cacheDir "" the
// verdicts are kept in memory only.
func OpenVerdicts(cacheDir string) *Verdicts {
	v := &Verdicts{known: map[verdictKey]verdict{}}
	if cacheDir == "" {
		return v
	}
	v.path = verdictsPath(cacheDir)
	data, err := os.ReadFile(v.path)
	if err != nil {
		return v
	}
	var f verdictsFile
	if json.Unmarshal(data, &f) != nil || f.Format != verdictsFormat {
		return v
	}
	for _, e := range f.Verdicts {
		switch {
		case e.Verified && e.Kind == "":
			v.known[e.verdictKey] = verdict{}
		case !e.Verified && (e.Kind == SignatureBad || e.Kind == SignatureUnknownKey):
			v.known[e.verdictKey] = verdict{kind: e.Kind, reason: e.Reason}
		}
	}
	return v
}

// Verify answers what VerifyRelease answers for rel, from the kept verdict
// when there is one for its assets and the trusted keys. It keeps a verdict
// of verified, bad, or unknown key, never a failed download or a missing
// asset, which costs no download to find again.
func (v *Verdicts) Verify(ctx context.Context, opt Options, rel *Release) error {
	key, ok := keyOf(opt, rel)
	if !ok {
		return VerifyRelease(ctx, opt, rel)
	}
	v.mu.Lock()
	known, hit := v.known[key]
	v.mu.Unlock()
	if hit {
		return known.err(rel.Tag)
	}
	err := VerifyRelease(ctx, opt, rel)
	var refused *SignatureError
	switch {
	case err == nil:
		v.store(key, verdict{})
	case errors.As(err, &refused) && (refused.Kind == SignatureBad || refused.Kind == SignatureUnknownKey):
		v.store(key, verdict{kind: refused.Kind, reason: refused.Reason})
	}
	return err
}

// keyOf is the key of rel's check with opt's trusted keys; false when the
// check is not kept: rel lacks checksums.txt or its signature, which
// VerifyRelease refuses without a download, or the trusted keys cannot be
// read or fingerprinted, which VerifyRelease answers or checks uncached.
func keyOf(opt Options, rel *Release) (verdictKey, bool) {
	sums, ok := rel.Asset(checksumsName)
	if !ok {
		return verdictKey{}, false
	}
	sig, ok := rel.Asset(signatureName)
	if !ok {
		return verdictKey{}, false
	}
	keys, err := TrustedKeys(opt)
	if err != nil {
		return verdictKey{}, false
	}
	fingerprint, err := KeysFingerprint(keys)
	if err != nil {
		return verdictKey{}, false
	}
	pub, _ := rel.Asset(publicKeyName)
	return verdictKey{Tag: rel.Tag, ChecksumsURL: sums.URL, SignatureURL: sig.URL, PublicKeyURL: pub.URL, Keys: fingerprint}, true
}

func (v *Verdicts) store(key verdictKey, got verdict) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.known[key] = got
	v.dirty = true
}

// Save writes the verdicts to the file under cache_dir when one was stored
// since it was read, replacing it in one step; another process's verdicts
// saved in between are replaced, and found again on its next check.
func (v *Verdicts) Save() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.path == "" || !v.dirty {
		return nil
	}
	f := verdictsFile{Format: verdictsFormat, Verdicts: make([]verdictEntry, 0, len(v.known))}
	for key, got := range v.known {
		f.Verdicts = append(f.Verdicts, verdictEntry{verdictKey: key, Verified: got.kind == "", Kind: got.kind, Reason: got.reason})
	}
	slices.SortFunc(f.Verdicts, func(a, b verdictEntry) int {
		return cmp.Or(strings.Compare(a.Tag, b.Tag), strings.Compare(a.ChecksumsURL, b.ChecksumsURL),
			strings.Compare(a.SignatureURL, b.SignatureURL), strings.Compare(a.PublicKeyURL, b.PublicKeyURL), strings.Compare(a.Keys, b.Keys))
	})
	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("save the release signature verdicts: %w", err)
	}
	dir := filepath.Dir(v.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("save the release signature verdicts: %w; check that %s can be created and written, or remove it", err, dir)
	}
	if err := atomicfile.WriteFile(v.path, append(data, '\n'), 0o600); err != nil {
		return fmt.Errorf("save the release signature verdicts to %s: %w; check that %s is writable", v.path, err, dir)
	}
	v.dirty = false
	return nil
}
