package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

// dashboardTags are the stand-in's flaiover tags: three releases, out of
// order, and a prerelease, which is not one.
var dashboardTags = []string{"flaiover/v0.3.0", "flaiover/v0.10.0", "flaiover/v0.4.0", "flaiover/v0.5.0-rc.1"}

// flaioverRelease is one flaiover GitHub release dashboardStandIn serves: its
// version, and its assets by name.
type flaioverRelease struct {
	version string
	assets  map[string][]byte
}

// digestListFor is a well-formed digest list naming version.
func digestListFor(version string) []byte {
	ref := "ghcr.io/bytepunx/flaiover:" + version
	return []byte("sha256:" + strings.Repeat("a", 64) + " " + ref + " index\n" +
		"sha256:" + strings.Repeat("b", 64) + " " + ref + " linux/amd64\n")
}

// signedFlaiover is release version publishing list as its digest list,
// signed by key, beside cosign.pub, key's public half.
func signedFlaiover(t *testing.T, version string, list []byte, key *ecdsa.PrivateKey) flaioverRelease {
	t.Helper()
	name := "flaiover_" + version + ".digests"
	return flaioverRelease{version: version, assets: map[string][]byte{
		name:          list,
		name + ".sig": signBlob(t, key, list),
		"cosign.pub":  publicPEM(t, key),
	}}
}

// otherKey is a key this flai does not trust.
func otherKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

// dashboardStandIn answers as GitHub's API does for o/r: the flaiover tags
// given, and rels as its GitHub releases, each asset at
// /assets/<version>/<name>; asked lists the paths asked for so far.
func dashboardStandIn(t *testing.T, tags []string, rels ...flaioverRelease) (url string, asked func() []string) {
	t.Helper()
	var mu sync.Mutex
	var paths []string
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		mu.Unlock()
		switch r.URL.Path {
		case "/repos/o/r/releases":
			list := make([]string, len(rels))
			for i, rel := range rels {
				var assets []string
				for name := range rel.assets {
					assets = append(assets, fmt.Sprintf(`{"name":%q,"url":"%s/assets/%s/%s"}`, name, srv.URL, rel.version, name))
				}
				list[i] = fmt.Sprintf(`{"tag_name":"flaiover/v%s","assets":[%s]}`, rel.version, strings.Join(assets, ","))
			}
			fmt.Fprint(w, "["+strings.Join(list, ",")+"]")
		case "/repos/o/r/git/matching-refs/tags/flaiover/v":
			refs := make([]string, len(tags))
			for i, tag := range tags {
				refs[i] = `{"ref":"refs/tags/` + tag + `"}`
			}
			fmt.Fprint(w, "["+strings.Join(refs, ",")+"]")
		default:
			for _, rel := range rels {
				if name, ok := strings.CutPrefix(r.URL.Path, "/assets/"+rel.version+"/"); ok {
					if data, ok := rel.assets[name]; ok {
						_, _ = w.Write(data)
						return
					}
				}
			}
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv.URL, func() []string {
		mu.Lock()
		defer mu.Unlock()
		return slices.Clone(paths)
	}
}

// signedDashboardStandIn serves dashboardTags with each release's digest
// list signed by a key this flai trusts for the rest of t.
func signedDashboardStandIn(t *testing.T) (url string, asked func() []string) {
	t.Helper()
	key := trustTestKey(t)
	var rels []flaioverRelease
	for _, v := range []string{"0.3.0", "0.10.0", "0.4.0"} {
		rels = append(rels, signedFlaiover(t, v, digestListFor(v), key))
	}
	return dashboardStandIn(t, dashboardTags, rels...)
}

// versionsProject is a project whose configured dashboard tag is tag.
func versionsProject(t *testing.T, tag string) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("GITHUB_TOKEN", "tok")
	root := tempProject(t)
	m := "version: 1\nname: My Proj\nkey: m\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  port: 5555\n  tag: " + tag + "\n"
	if err := os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte(m), 0o644); err != nil {
		t.Fatal(err)
	}
	return root
}

// marksOf reads a versions listing into its versions, in order, and the
// marks on each.
func marksOf(out string) ([]string, map[string]string) {
	var versions []string
	marks := map[string]string{}
	for _, l := range strings.Split(out, "\n") {
		if !strings.HasPrefix(l, "  ") {
			continue
		}
		f := strings.Fields(l)
		versions = append(versions, f[0])
		marks[f[0]] = strings.Join(f[1:], " ")
	}
	return versions, marks
}

// S-0298: dashboard versions lists the published dashboard releases newest
// first, marking the release the running container carries, the configured
// tag, and the newest.
func TestDashboardVersions(t *testing.T) {
	root := versionsProject(t, "0.4.0")
	api, _ := signedDashboardStandIn(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{"flaiover": true},
		env: map[string]string{"flaiover": "PATH=/usr/bin\nFLAIOVER_VERSION=0.3.0\n"}}
	out, errOut, code := runWith(t, root, f, "dashboard", "versions", "--repo", "o/r", "--api", api)
	if code != 0 {
		t.Fatalf("versions: %d %s", code, errOut)
	}
	versions, marks := marksOf(out)
	if !slices.Equal(versions, []string{"0.10.0", "0.4.0", "0.3.0"}) {
		t.Fatalf("newest first, releases only: %v\n%s", versions, out)
	}
	if marks["0.10.0"] != "latest" || marks["0.4.0"] != "configured" || marks["0.3.0"] != "running" {
		t.Errorf("marks: %v\n%s", marks, out)
	}
	if strings.Contains(out, "running:") || strings.Contains(out, "configured:") {
		t.Errorf("both are releases, so no note: %s", out)
	}

	out, errOut, code = runWith(t, root, f, "dashboard", "versions", "--repo", "o/r", "--api", api, "--json")
	if code != 0 {
		t.Fatalf("versions --json: %d %s", code, errOut)
	}
	var listed []map[string]any
	if err := json.Unmarshal([]byte(out), &listed); err != nil || len(listed) != 3 {
		t.Fatalf("json: %v %s", err, out)
	}
	want := []map[string]any{
		{"version": "0.10.0", "tag": "flaiover/v0.10.0", "running": false, "configured": false, "latest": true},
		{"version": "0.4.0", "tag": "flaiover/v0.4.0", "running": false, "configured": true, "latest": false},
		{"version": "0.3.0", "tag": "flaiover/v0.3.0", "running": true, "configured": false, "latest": false},
	}
	for i, w := range want {
		if len(listed[i]) != len(w) {
			t.Errorf("fields of %d: %v", i, listed[i])
		}
		for k, v := range w {
			if listed[i][k] != v {
				t.Errorf("%d %s: %v, want %v", i, k, listed[i][k], v)
			}
		}
	}
}

// S-0298: with no FLAIOVER_VERSION the image tag says what runs; a tag that
// is not a release is said to be none.
func TestDashboardVersionsNotesWhatIsNotARelease(t *testing.T) {
	root := versionsProject(t, "latest")
	api, _ := signedDashboardStandIn(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{"flaiover": true},
		containerRef: map[string]string{"flaiover": "ghcr.io/bytepunx/flaiover:0.4.0"}}
	out, errOut, code := runWith(t, root, f, "dashboard", "versions", "--repo", "o/r", "--api", api)
	if _, marks := marksOf(out); code != 0 || marks["0.4.0"] != "running" {
		t.Fatalf("running from the image tag: %d %s %s", code, out, errOut)
	}
	if !strings.Contains(out, "configured: ghcr.io/bytepunx/flaiover:latest, whose tag names no published release") {
		t.Errorf("latest is no release: %s", out)
	}
	f.containerRef = nil // the fake's container then runs 0.2.0, which is not published
	out, _, _ = runWith(t, root, f, "dashboard", "versions", "--repo", "o/r", "--api", api)
	if !strings.Contains(out, "running: 0.2.0, which is not a published release") {
		t.Errorf("unpublished running: %s", out)
	}
}

// S-0298: the list needs no Docker; without it nothing is marked running.
func TestDashboardVersionsWithoutDocker(t *testing.T) {
	root := versionsProject(t, "0.4.0")
	api, _ := signedDashboardStandIn(t)
	f := &fakeRunner{missing: true, images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard", "versions", "--repo", "o/r", "--api", api)
	if code != 0 {
		t.Fatalf("versions: %d %s", code, errOut)
	}
	if versions, marks := marksOf(out); len(versions) != 3 || strings.Contains(out, "running,") || marks["0.4.0"] != "configured" {
		t.Errorf("listed without running: %v %s", marks, out)
	}
	if !strings.Contains(out, "running: no flaiover container found") {
		t.Errorf("says nothing was found running: %s", out)
	}
	for _, c := range f.calls {
		if strings.HasPrefix(c, "docker ") {
			t.Errorf("no docker call without docker: %s", c)
		}
	}
}

// S-0298, ADR-0117 §3: with --published a tag that is not a published
// dashboard release is refused, naming the published ones, before any
// docker call.
func TestDashboardUpgradePublishedRefusesAnUnpublishedTag(t *testing.T) {
	root := versionsProject(t, "latest")
	api, asked := signedDashboardStandIn(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{"flaiover": true}}
	for _, tag := range []string{"0.9.9", "v0.4.0", "0.5.0-rc.1"} {
		_, errOut, code := runWith(t, root, f, "dashboard", "upgrade", "--published", "--tag", tag, "--repo", "o/r", "--api", api)
		if code == 0 || !strings.Contains(errOut, "published are 0.10.0, 0.4.0, 0.3.0") {
			t.Errorf("%s refused: %d %s", tag, code, errOut)
		}
	}
	_, errOut, code := runWith(t, root, f, "dashboard", "upgrade", "--published", "--repo", "o/r", "--api", api)
	if code == 0 || !strings.Contains(errOut, "--tag") {
		t.Errorf("--published needs --tag: %d %s", code, errOut)
	}
	if n := len(asked()); n != 3 {
		t.Errorf("each tag is checked once, and no tag is not: %d asks", n)
	}
	for _, c := range f.calls {
		if strings.HasPrefix(c, "docker ") {
			t.Errorf("no docker call before the refusal: %s", c)
		}
	}
}

// S-0298, ADR-0117 §4 and §5: a published tag deploys as an upgrade does,
// for that container once: dashboard.tag is not written. S-0236, ADR-0070:
// both the release that ran and the one deployed run by digest.
func TestDashboardUpgradePublishedDeploysAPublishedTag(t *testing.T) {
	root := versionsProject(t, "latest")
	api, _ := signedDashboardStandIn(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, errOut, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	manifestBefore, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	out, errOut, code := runWith(t, root, f, "dashboard", "upgrade", "--published", "--tag", "0.3.0", "--repo", "o/r", "--api", api)
	from, to := flaioverRef("0.10.0"), flaioverRef("0.3.0")
	if code != 0 || !strings.Contains(out, "upgraded flaiover from "+from+" to "+to) {
		t.Fatalf("upgrade: %d %s %s", code, out, errOut)
	}
	if !slices.Contains(f.calls, "docker pull --quiet --platform linux/amd64 "+to) || f.containerRef["flaiover"] != to || f.containerLabel["flaiover"] != to {
		t.Errorf("pulled and swapped to 0.3.0 by digest: %v %v", f.containerRef, f.calls)
	}
	cfg, _ := os.ReadFile(os.Getenv("FLAI_CONFIG"))
	if strings.Contains(string(cfg), "0.3.0") {
		t.Errorf("dashboard.tag is not written: %s", cfg)
	}
	if manifestAfter, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml")); string(manifestAfter) != string(manifestBefore) {
		t.Errorf("the manifest is not changed: %s", manifestAfter)
	}
	// a tag without --published is used as it is, unchecked, for an image of
	// another repository; for the signed one it must name a release still
	if _, errOut, code := runWith(t, root, f, "dashboard", "upgrade", "--image", mirrorImage, "--tag", "mirror-build", "--api", "http://127.0.0.1:1"); code != 0 || f.containerLabel["flaiover"] != mirrorImage+":mirror-build" {
		t.Errorf("unchecked tag: %d %s", code, errOut)
	}
	if _, errOut, code := runWith(t, root, f, "dashboard", "upgrade", "--tag", "mirror-build"); code == 0 || !strings.Contains(errOut, "is not a flaiover release version") {
		t.Errorf("the signed repository resolves the tag: %d %s", code, errOut)
	}
}

// S-0236, ADR-0070: dashboard versions marks each release with no digest
// list this flai can verify, saying why, in its table and its JSON, and
// leaves a signed one unmarked.
func TestDashboardVersionsMarksReleasesWithNoListItCanVerify(t *testing.T) {
	root := versionsProject(t, "0.9.0")
	key, other := trustTestKey(t), otherKey(t)
	noList := signedFlaiover(t, "0.8.0", digestListFor("0.8.0"), key)
	delete(noList.assets, "flaiover_0.8.0.digests")
	altered := signedFlaiover(t, "0.7.0", digestListFor("0.7.0"), key)
	altered.assets["flaiover_0.7.0.digests"] = digestListFor("0.7.1")
	unknown := signedFlaiover(t, "0.6.0", digestListFor("0.6.0"), other)
	forAnother := signedFlaiover(t, "0.5.0", digestListFor("0.4.0"), key)
	malformed := signedFlaiover(t, "0.4.0", []byte("not a digest list\n"), key)
	api, _ := dashboardStandIn(t,
		[]string{"flaiover/v0.10.0", "flaiover/v0.9.0", "flaiover/v0.8.0", "flaiover/v0.7.0", "flaiover/v0.6.0", "flaiover/v0.5.0", "flaiover/v0.4.0"},
		signedFlaiover(t, "0.10.0", digestListFor("0.10.0"), key), noList, altered, unknown, forAnother, malformed)
	f := &fakeRunner{missing: true, images: map[string]bool{}, running: map[string]bool{}}

	out, errOut, code := runWith(t, root, f, "dashboard", "versions", "--repo", "o/r", "--api", api)
	if code != 0 {
		t.Fatalf("versions: %d %s", code, errOut)
	}
	_, marks := marksOf(out)
	wantMarks := map[string]string{
		"0.10.0": "latest",
		"0.9.0":  "configured, not deployable: it has no GitHub release with a signed digest list",
		"0.8.0":  "not deployable: it publishes no signed digest list",
		"0.7.0":  "not deployable: its digest list's signature does not verify, so it may have been altered",
		"0.6.0":  "not deployable: its digest list is signed by a key this flai does not trust",
		"0.5.0":  "not deployable: its digest list is for another version",
		"0.4.0":  "not deployable: its digest list is not in the format release-flaiover.yml writes",
	}
	for v, w := range wantMarks {
		if marks[v] != w {
			t.Errorf("%s: %q, want %q\n%s", v, marks[v], w, out)
		}
	}

	out, errOut, code = runWith(t, root, f, "dashboard", "versions", "--repo", "o/r", "--api", api, "--json")
	if code != 0 {
		t.Fatalf("versions --json: %d %s", code, errOut)
	}
	var listed []map[string]any
	if err := json.Unmarshal([]byte(out), &listed); err != nil || len(listed) != 7 {
		t.Fatalf("json: %v %s", err, out)
	}
	wantSignature := []string{"", "missing", "missing", "bad", "unknown-key", "bad", "bad"}
	for i, w := range wantSignature {
		got := listed[i]
		if w == "" {
			for _, k := range []string{"unverifiable", "signature", "reason"} {
				if _, ok := got[k]; ok {
					t.Errorf("a signed release has no %s: %v", k, got)
				}
			}
			continue
		}
		reason := strings.TrimPrefix(wantMarks[got["version"].(string)], "configured, ")
		if got["unverifiable"] != true || got["signature"] != w || got["reason"] != strings.TrimPrefix(reason, "not deployable: ") {
			t.Errorf("%d: %v, want signature %s", i, got, w)
		}
	}
}

// S-0236, ADR-0070: with --published a release with no digest list this flai
// can verify is refused, naming why and the deployable ones, before any
// docker call; TestDashboardUpgradePublishedDeploysAPublishedTag deploys a
// signed one.
func TestDashboardUpgradePublishedRefusesAnUnverifiableRelease(t *testing.T) {
	root := versionsProject(t, "latest")
	key := trustTestKey(t)
	api, _ := dashboardStandIn(t, []string{"flaiover/v0.10.0", "flaiover/v0.4.0", "flaiover/v0.3.0"},
		signedFlaiover(t, "0.10.0", digestListFor("0.10.0"), otherKey(t)),
		signedFlaiover(t, "0.3.0", digestListFor("0.3.0"), key))
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	for tag, why := range map[string]string{
		"0.10.0": "it is not deployable: its digest list is signed by a key this flai does not trust",
		"0.4.0":  "it is not deployable: it has no GitHub release with a signed digest list",
	} {
		_, errOut, code := runWith(t, root, f, "dashboard", "upgrade", "--published", "--tag", tag, "--repo", "o/r", "--api", api)
		if code == 0 || !strings.Contains(errOut, why) || !strings.Contains(errOut, "deployable are 0.3.0") {
			t.Errorf("%s refused: %d %s", tag, code, errOut)
		}
	}
	for _, c := range f.calls {
		if strings.HasPrefix(c, "docker ") {
			t.Errorf("no docker call before the refusal: %s", c)
		}
	}

	// with none deployable, the refusal says so
	none, _ := dashboardStandIn(t, []string{"flaiover/v0.4.0"})
	_, errOut, code := runWith(t, root, f, "dashboard", "upgrade", "--published", "--tag", "0.4.0", "--repo", "o/r", "--api", none)
	if code == 0 || !strings.Contains(errOut, "no published dashboard release of o/r has a digest list this flai can verify") {
		t.Errorf("none deployable: %d %s", code, errOut)
	}
}
