package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// dashboardTags are the stand-in's flaiover tags: three releases, out of
// order, and a prerelease, which is not one.
var dashboardTags = []string{"flaiover/v0.3.0", "flaiover/v0.10.0", "flaiover/v0.4.0", "flaiover/v0.5.0-rc.1"}

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
	api, _ := releasesStandIn(t, "[]", dashboardTags...)
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
	api, _ := releasesStandIn(t, "[]", dashboardTags...)
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
	api, _ := releasesStandIn(t, "[]", dashboardTags...)
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
	api, asked := releasesStandIn(t, "[]", dashboardTags...)
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
// for that container once: dashboard.tag is not written.
func TestDashboardUpgradePublishedDeploysAPublishedTag(t *testing.T) {
	root := versionsProject(t, "latest")
	api, _ := releasesStandIn(t, "[]", dashboardTags...)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, errOut, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	manifestBefore, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml"))
	out, errOut, code := runWith(t, root, f, "dashboard", "upgrade", "--published", "--tag", "0.3.0", "--repo", "o/r", "--api", api)
	if code != 0 || !strings.Contains(out, "upgraded flaiover from ghcr.io/bytepunx/flaiover:latest to ghcr.io/bytepunx/flaiover:0.3.0") {
		t.Fatalf("upgrade: %d %s %s", code, out, errOut)
	}
	if !slices.Contains(f.calls, "docker pull --quiet --platform linux/amd64 ghcr.io/bytepunx/flaiover:0.3.0") || f.containerRef["flaiover"] != "ghcr.io/bytepunx/flaiover:0.3.0" {
		t.Errorf("pulled and swapped to 0.3.0: %v %v", f.containerRef, f.calls)
	}
	cfg, _ := os.ReadFile(os.Getenv("FLAI_CONFIG"))
	if strings.Contains(string(cfg), "0.3.0") {
		t.Errorf("dashboard.tag is not written: %s", cfg)
	}
	if manifestAfter, _ := os.ReadFile(filepath.Join(root, "system-flow.yaml")); string(manifestAfter) != string(manifestBefore) {
		t.Errorf("the manifest is not changed: %s", manifestAfter)
	}
	// a tag without --published is used as it is, unchecked
	if _, errOut, code := runWith(t, root, f, "dashboard", "upgrade", "--tag", "mirror-build", "--api", "http://127.0.0.1:1"); code != 0 {
		t.Errorf("unchecked tag: %d %s", code, errOut)
	}
}
