package dashboard

import (
	"crypto/ecdsa"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newCache is a cache in a fresh state directory trusting key.
func newCache(t *testing.T, key *ecdsa.PrivateKey) Cache {
	t.Helper()
	return NewCache(t.TempDir(), []*ecdsa.PublicKey{&key.PublicKey})
}

// put caches version's list, signed by key, failing the test if refused.
func put(t *testing.T, c Cache, key *ecdsa.PrivateKey, version string, data []byte) {
	t.Helper()
	if err := c.Put(version, data, sign(t, key, data)); err != nil {
		t.Fatal(err)
	}
}

func TestNewCacheIsTheDashboardDigestsFolderUnderServesState(t *testing.T) {
	c := NewCache("/state/serve", nil)
	if c.Dir != filepath.Join("/state/serve", "dashboard-digests") {
		t.Fatalf("dir %q", c.Dir)
	}
}

func TestAVerifiedListIsCachedAndReadBackByVersionAndDigest(t *testing.T) {
	key := newKey(t)
	c := newCache(t, key)
	data := []byte(goodList)
	put(t, c, key, goodV, data)

	for _, name := range []string{"0.42.3.digests", "0.42.3.digests.sig"} {
		info, err := os.Stat(filepath.Join(c.Dir, name))
		if err != nil {
			t.Fatalf("%s not written: %v", name, err)
		}
		if info.Mode().Perm() != 0o600 {
			t.Errorf("%s mode %v, want 0600", name, info.Mode().Perm())
		}
	}

	l, ok, err := c.Get(goodV)
	if err != nil || !ok {
		t.Fatalf("Get: ok %v err %v", ok, err)
	}
	if l.Version != goodV || l.Index() != index {
		t.Fatalf("Get answered %+v", l)
	}

	for _, d := range []string{index, amd64, arm64} {
		l, ok, err := c.Lookup(d)
		if err != nil || !ok || l.Version != goodV {
			t.Fatalf("Lookup %s: version %q ok %v err %v", d, l.Version, ok, err)
		}
	}
	if _, ok, err := c.Lookup(other); ok || err != nil {
		t.Fatalf("Lookup of an unlisted digest: ok %v err %v, want false and nil", ok, err)
	}
}

func TestNothingCachedIsNotFound(t *testing.T) {
	c := newCache(t, newKey(t))
	if _, ok, err := c.Get(goodV); ok || err != nil {
		t.Fatalf("Get: ok %v err %v, want false and nil", ok, err)
	}
	if _, ok, err := c.Lookup(index); ok || err != nil {
		t.Fatalf("Lookup: ok %v err %v, want false and nil", ok, err)
	}
	if lists, skipped := c.All(); len(lists) != 0 || len(skipped) != 0 {
		t.Fatalf("All: %v %v, want nothing", lists, skipped)
	}
}

func TestLookupAnswersTheNewestVersionThatListsTheDigest(t *testing.T) {
	key := newKey(t)
	c := newCache(t, key)
	// An image rebuilt to the same digest is listed by both releases.
	put(t, c, key, "0.9.0", listFor("0.9.0", other))
	put(t, c, key, "0.10.0", listFor("0.10.0", index))
	put(t, c, key, "0.2.0", listFor("0.2.0", index))

	l, ok, err := c.Lookup(amd64)
	if err != nil || !ok || l.Version != "0.10.0" {
		t.Fatalf("Lookup: version %q ok %v err %v, want 0.10.0", l.Version, ok, err)
	}
	l, ok, err = c.Lookup(other)
	if err != nil || !ok || l.Version != "0.9.0" {
		t.Fatalf("Lookup: version %q ok %v err %v, want 0.9.0", l.Version, ok, err)
	}
	lists, skipped := c.All()
	if len(skipped) != 0 {
		t.Fatal(skipped)
	}
	var got []string
	for _, l := range lists {
		got = append(got, l.Version)
	}
	if strings.Join(got, " ") != "0.10.0 0.9.0 0.2.0" {
		t.Fatalf("All in order %v, want newest first", got)
	}
}

func TestACachedListThatNoLongerVerifiesIsNotTrusted(t *testing.T) {
	tamper := map[string]func(t *testing.T, c Cache){
		"list altered": func(t *testing.T, c Cache) {
			altered := strings.Replace(goodList, amd64, other, 1)
			write(t, filepath.Join(c.Dir, "0.42.3.digests"), altered)
		},
		"signature replaced": func(t *testing.T, c Cache) {
			write(t, filepath.Join(c.Dir, "0.42.3.digests.sig"), string(sign(t, newKey(t), []byte(goodList))))
		},
		"signature removed": func(t *testing.T, c Cache) {
			if err := os.Remove(filepath.Join(c.Dir, "0.42.3.digests.sig")); err != nil {
				t.Fatal(err)
			}
		},
		"renamed to another version": func(t *testing.T, c Cache) {
			for _, ext := range []string{".digests", ".digests.sig"} {
				if err := os.Rename(filepath.Join(c.Dir, "0.42.3"+ext), filepath.Join(c.Dir, "0.43.0"+ext)); err != nil {
					t.Fatal(err)
				}
			}
		},
	}
	for name, change := range tamper {
		t.Run(name, func(t *testing.T) {
			key := newKey(t)
			c := newCache(t, key)
			put(t, c, key, goodV, []byte(goodList))
			change(t, c)

			for _, v := range []string{goodV, "0.43.0"} {
				if _, ok, err := c.Get(v); ok {
					t.Fatalf("Get %s trusts a tampered list", v)
				} else if err != nil && !strings.Contains(err.Error(), c.Dir) {
					t.Fatalf("Get %s error does not name the file: %v", v, err)
				}
			}
			_, ok, err := c.Lookup(index)
			if ok {
				t.Fatal("Lookup trusts a tampered list")
			}
			if err == nil || !strings.Contains(err.Error(), filepath.Join(c.Dir, "0.4")) {
				t.Fatalf("Lookup error does not name the skipped file: %v", err)
			}
			lists, skipped := c.All()
			if len(lists) != 0 || len(skipped) != 1 {
				t.Fatalf("All: %d lists %d skipped, want 0 and 1", len(lists), len(skipped))
			}
		})
	}
}

func TestAVerifiedListIsStillFoundBesideOneThatIsNot(t *testing.T) {
	key := newKey(t)
	c := newCache(t, key)
	put(t, c, key, goodV, []byte(goodList))
	put(t, c, key, "0.1.0", listFor("0.1.0", other))
	write(t, filepath.Join(c.Dir, "0.42.3.digests"), strings.Replace(goodList, "0.42.3", "0.42.9", 1))

	l, ok, err := c.Lookup(amd64)
	if err != nil || !ok || l.Version != "0.1.0" {
		t.Fatalf("Lookup: version %q ok %v err %v, want 0.1.0 only", l.Version, ok, err)
	}
	if _, ok, _ := c.Lookup(index); ok {
		t.Fatal("Lookup answered the tampered list's index")
	}
}

func TestPutRefusesAListItCannotTrust(t *testing.T) {
	key := newKey(t)
	c := newCache(t, key)
	data := []byte(goodList)
	cases := map[string]struct {
		version string
		sig     []byte
		why     string
	}{
		"unknown key":     {goodV, sign(t, newKey(t), data), "does not verify"},
		"no signature":    {goodV, nil, "does not verify"},
		"another version": {"0.42.4", sign(t, key, data), "lists version 0.42.3"},
		"not a version":   {"../x", sign(t, key, data), "lists version 0.42.3"},
	}
	for name, tc := range cases {
		if err := c.Put(tc.version, data, tc.sig); err == nil || !strings.Contains(err.Error(), tc.why) {
			t.Errorf("%s: err %v, want a refusal saying %q", name, err, tc.why)
		}
	}
	if _, err := os.Stat(c.Dir); !os.IsNotExist(err) {
		t.Fatalf("a refused Put wrote %s: %v", c.Dir, err)
	}
}

func TestGetRefusesAVersionThatIsNotARelease(t *testing.T) {
	c := newCache(t, newKey(t))
	if _, ok, err := c.Get("../serve"); ok || err == nil {
		t.Fatalf("ok %v err %v, want a refusal", ok, err)
	}
}

func TestAllNamesACachedFileThatIsNotARelease(t *testing.T) {
	key := newKey(t)
	c := newCache(t, key)
	put(t, c, key, goodV, []byte(goodList))
	write(t, filepath.Join(c.Dir, "latest.digests"), goodList)
	write(t, filepath.Join(c.Dir, ".0.42.3.digests.123.tmp"), "partial")

	lists, skipped := c.All()
	if len(lists) != 1 || len(skipped) != 1 || !strings.Contains(skipped[0].Error(), "latest") {
		t.Fatalf("All: %d lists, skipped %v; want 1 list and latest.digests named", len(lists), skipped)
	}
}

// write replaces path's content, failing the test if it cannot.
func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
