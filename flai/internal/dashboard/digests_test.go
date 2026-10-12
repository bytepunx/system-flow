package dashboard

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"testing"
)

const (
	repo   = "ghcr.io/bytepunx/flaiover"
	index  = "sha256:0dce07ffc38f1fecc91cbb659524525db687eb12030bea40209f2ecf3636f15d"
	amd64  = "sha256:b3c17aafc5f6026830b885843021a783cc6e129937cc4f1f1dc44e546ff38ec0"
	arm64  = "sha256:b8ec2bdf93dfbd0cdbe8b756bdd0ce30d66242a93bd99d5e869f0668a475eb0c"
	other  = "sha256:1111111111111111111111111111111111111111111111111111111111111111"
	goodV  = "0.42.3"
	goodRf = repo + ":" + goodV
)

// goodList is the example design/system/release-signing.md gives.
const goodList = index + " " + goodRf + " index\n" +
	amd64 + " " + goodRf + " linux/amd64\n" +
	arm64 + " " + goodRf + " linux/arm64\n"

// listFor is a well-formed list for version with index digest idx and the
// example's linux/amd64 and linux/arm64 digests.
func listFor(version, idx string) []byte {
	ref := repo + ":" + version
	return []byte(idx + " " + ref + " index\n" + amd64 + " " + ref + " linux/amd64\n" + arm64 + " " + ref + " linux/arm64\n")
}

// newKey is a fresh P-256 key, as cosign generates.
func newKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

// sign is data's signature by key as cosign sign-blob writes it: base64 of
// the ASN.1 signature over its SHA-256, with a final newline.
func sign(t *testing.T, key *ecdsa.PrivateKey, data []byte) []byte {
	t.Helper()
	sum := sha256.Sum256(data)
	der, err := ecdsa.SignASN1(rand.Reader, key, sum[:])
	if err != nil {
		t.Fatal(err)
	}
	return []byte(base64.StdEncoding.EncodeToString(der) + "\n")
}

func TestAGoodListParses(t *testing.T) {
	l, err := Parse([]byte(goodList))
	if err != nil {
		t.Fatal(err)
	}
	if l.Version != goodV || l.Repository != repo {
		t.Fatalf("version %q repository %q, want %q %q", l.Version, l.Repository, goodV, repo)
	}
	if l.Index() != index {
		t.Fatalf("index %q, want %q", l.Index(), index)
	}
	if got := l.Reference(); got != repo+"@"+index {
		t.Fatalf("reference %q", got)
	}
	if got := strings.Join(l.Digests(), ","); got != index+","+amd64+","+arm64 {
		t.Fatalf("digests %s", got)
	}
	if want := (Entry{Digest: amd64, Ref: goodRf, Kind: "linux/amd64"}); l.Entries[1] != want {
		t.Fatalf("entry 2 %+v, want %+v", l.Entries[1], want)
	}
	if !l.Has(arm64) || l.Has(other) {
		t.Fatal("Has answers wrongly")
	}
	if got := ListName(goodV); got != "flaiover_0.42.3.digests" {
		t.Fatalf("ListName %q", got)
	}
}

func TestAListWithAVariantAndARegistryPortParses(t *testing.T) {
	ref := "localhost:5000/flaiover:1.0.0"
	data := index + " " + ref + " index\n" + amd64 + " " + ref + " linux/arm/v7\n"
	l, err := Parse([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	if l.Repository != "localhost:5000/flaiover" || l.Version != "1.0.0" {
		t.Fatalf("repository %q version %q", l.Repository, l.Version)
	}
}

func TestEachMalformedShapeIsRefusedNamingTheLine(t *testing.T) {
	line1 := index + " " + goodRf + " index"
	line2 := amd64 + " " + goodRf + " linux/amd64"
	line3 := arm64 + " " + goodRf + " linux/arm64"
	cases := []struct {
		name, data, line, why string
	}{
		{"empty", "", "at it", "empty"},
		{"crlf", line1 + "\r\n" + line2 + "\r\n", "line 1", "carriage return"},
		{"no final newline", line1 + "\n" + line2, "line 2", "newline"},
		{"blank line", line1 + "\n\n" + line2 + "\n", "line 2", "blank"},
		{"trailing blank line", line1 + "\n" + line2 + "\n\n", "line 3", "blank"},
		{"two fields", line1 + "\n" + amd64 + " " + goodRf + "\n", "line 2", "2 fields"},
		{"four fields", line1 + "\n" + line2 + " extra\n", "line 2", "4 fields"},
		{"two spaces", line1 + "\n" + amd64 + "  " + goodRf + " linux/amd64\n", "line 2", "4 fields"},
		{"tab separated", line1 + "\n" + strings.ReplaceAll(line2, " ", "\t") + "\n", "line 2", "1 fields"},
		{"bad digest", line1 + "\nsha256:abc " + goodRf + " linux/amd64\n", "line 2", "digest"},
		{"upper-case hex", line1 + "\n" + strings.ToUpper(amd64[:9]) + amd64[9:] + " " + goodRf + " linux/amd64\n", "line 2", "digest"},
		{"upper-case hex digits", line1 + "\nsha256:" + strings.ToUpper(amd64[7:]) + " " + goodRf + " linux/amd64\n", "line 2", "digest"},
		{"other algorithm", line1 + "\nsha512:" + amd64[7:] + " " + goodRf + " linux/amd64\n", "line 2", "digest"},
		{"ref without version", line1 + "\n" + amd64 + " " + repo + " linux/amd64\n", "line 2", "reference"},
		{"ref with v", index + " " + repo + ":v0.42.3 index\n", "line 1", "version"},
		{"refs disagree on version", line1 + "\n" + amd64 + " " + repo + ":0.42.4 linux/amd64\n", "line 2", "line 1's"},
		{"refs disagree on repository", line1 + "\n" + amd64 + " ghcr.io/else/flaiover:0.42.3 linux/amd64\n", "line 2", "line 1's"},
		{"first line not index", line2 + "\n" + line3 + "\n", "line 1", "first line is the index"},
		{"index twice", line1 + "\n" + line1 + "\n", "line 2", "only the first line"},
		{"unsorted platforms", line1 + "\n" + line3 + "\n" + line2 + "\n", "line 3", "sorted"},
		{"duplicate platform", line1 + "\n" + line2 + "\n" + arm64 + " " + goodRf + " linux/amd64\n", "line 3", "line 2's too"},
		{"bad kind", line1 + "\n" + amd64 + " " + goodRf + " amd64\n", "line 2", "kind"},
		{"not utf-8", line1 + "\n" + amd64 + " " + goodRf + " linux/\xff\n", "line 2", "UTF-8"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Parse([]byte(tc.data))
			if err == nil {
				t.Fatal("parsed, want refused")
			}
			if !strings.Contains(err.Error(), "failed "+tc.line+":") && !strings.Contains(err.Error(), "failed at "+tc.line+":") {
				t.Fatalf("error does not name %s: %v", tc.line, err)
			}
			if !strings.Contains(err.Error(), tc.why) {
				t.Fatalf("error does not say %q: %v", tc.why, err)
			}
		})
	}
}

func TestARepoDigestMatchesByRepositoryAndDigest(t *testing.T) {
	l, err := Parse([]byte(goodList))
	if err != nil {
		t.Fatal(err)
	}
	for _, rd := range []string{repo + "@" + index, repo + "@" + amd64, repo + "@" + arm64} {
		if !l.Matches(rd) {
			t.Errorf("%s does not match, want it to", rd)
		}
	}
	for _, rd := range []string{
		repo + "@" + other,
		"ghcr.io/else/flaiover@" + index,
		"ghcr.io/bytepunx/flaiover:0.42.3@" + index,
		repo + index,
		repo + "@" + index + "@" + index,
		"",
	} {
		if l.Matches(rd) {
			t.Errorf("%s matches, want it not to", rd)
		}
	}
}

func TestVerifyTrustsOnlyAListSignedByATrustedKey(t *testing.T) {
	trusted, stranger := newKey(t), newKey(t)
	keys := []*ecdsa.PublicKey{&trusted.PublicKey}
	data := []byte(goodList)

	l, err := Verify(data, sign(t, trusted, data), keys)
	if err != nil {
		t.Fatalf("good signature refused: %v", err)
	}
	if l.Index() != index {
		t.Fatalf("index %q", l.Index())
	}

	cases := map[string][]byte{
		"unknown key":   sign(t, stranger, data),
		"other content": sign(t, trusted, []byte(goodList+"x")),
		"not base64":    []byte("not a signature"),
		"empty":         nil,
	}
	for name, sig := range cases {
		if _, err := Verify(data, sig, keys); err == nil || !strings.Contains(err.Error(), "does not verify") {
			t.Errorf("%s: err %v, want a refusal saying it does not verify", name, err)
		}
	}
}

func TestVerifyRefusesASignedListThatDoesNotParse(t *testing.T) {
	key := newKey(t)
	data := []byte(strings.TrimSuffix(goodList, "\n"))
	if _, err := Verify(data, sign(t, key, data), []*ecdsa.PublicKey{&key.PublicKey}); err == nil || !strings.Contains(err.Error(), "line 3") {
		t.Fatalf("err %v, want a refusal naming line 3", err)
	}
}
