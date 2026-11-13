package agentimage

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/bytepunx/system-flow/flai/images"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

const (
	buildVersion = "1.2.3"
	bedrockID    = "sha256:bedrock"
	harnessTag   = "flai-agent-claude-code:1.2.3"
)

// signingKey is a fresh ECDSA P-256 key, as the release key is.
func signingKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

// targz is a tar.gz holding a license and name with content, as GoReleaser packs a binary.
func targz(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, f := range []struct {
		name string
		data []byte
	}{{"LICENSE.md", []byte("license")}, {name, content}} {
		if err := tw.WriteHeader(&tar.Header{Name: f.name, Mode: 0o755, Size: int64(len(f.data)), Typeflag: tar.TypeReg}); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(f.data); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// sha is the hex SHA-256 of data.
func sha(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// release is a fake GitHub serving one flai release, flai/v1.2.3, with
// checksums.txt and its signature and the flai and pilot archives for
// linux/amd64, counting the requests for each asset.
type release struct {
	mu    sync.Mutex
	asked map[string]int
	srv   *httptest.Server
}

// releaseFiles are the assets a fake release serves, by name.
type releaseFiles map[string][]byte

// goodRelease is a release's assets: flai and pilot archives for
// linux/amd64 whose binaries hold flaiBin and pilotBin, and checksums.txt
// over them signed by key.
func goodRelease(t *testing.T, key *ecdsa.PrivateKey, flaiBin, pilotBin []byte) releaseFiles {
	t.Helper()
	flaiTar := targz(t, "flai", flaiBin)
	pilotTar := targz(t, "pilot", pilotBin)
	flaiName := selfupgrade.ArchiveName(buildVersion, "linux", "amd64")
	pilotName := "pilot_" + buildVersion + "_linux_amd64.tar.gz"
	sums := []byte(fmt.Sprintf("%s  %s\n%s  %s\n", sha(flaiTar), flaiName, sha(pilotTar), pilotName))
	return releaseFiles{
		"checksums.txt":     sums,
		"checksums.txt.sig": signBlob(t, key, sums),
		flaiName:            flaiTar,
		pilotName:           pilotTar,
	}
}

// serveRelease serves files as the assets of flai/v1.2.3.
func serveRelease(t *testing.T, files releaseFiles) *release {
	t.Helper()
	f := &release{asked: map[string]int{}}
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/o/r/releases" {
			var assets []string
			for name := range files {
				assets = append(assets, fmt.Sprintf(`{"name":%q,"url":"%s/assets/%s"}`, name, f.srv.URL, name))
			}
			fmt.Fprintf(w, `[{"tag_name":"flai/v%s","published_at":"2026-11-13T10:00:00Z","assets":[%s]}]`, buildVersion, strings.Join(assets, ","))
			return
		}
		name, ok := strings.CutPrefix(r.URL.Path, "/assets/")
		data, found := files[name]
		if !ok || !found {
			http.NotFound(w, r)
			return
		}
		f.mu.Lock()
		f.asked[name]++
		f.mu.Unlock()
		_, _ = w.Write(data)
	}))
	t.Cleanup(f.srv.Close)
	return f
}

// count is how many times the asset name was asked for.
func (f *release) count(name string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.asked[name]
}

// build is one docker invocation the fake docker saw: its folder and arguments.
type build struct {
	dir  string
	args []string
}

// buildDocker is an Exec that records each docker invocation and runs during
// on each before it returns.
type buildDocker struct {
	builds []build
	during func(b build)
}

// exec is the fake's Exec.
func (d *buildDocker) exec(_ context.Context, dir string, _ []string, _ io.Writer, _ string, args ...string) error {
	b := build{dir: dir, args: args}
	d.builds = append(d.builds, b)
	if d.during != nil {
		d.during(b)
	}
	return nil
}

// releaseBuilder is a Builder over the fake release f, trusting key, with
// the fake docker d, an engine on amd64, and bedrock built as bedrockID.
func releaseBuilder(t *testing.T, f *release, key *ecdsa.PrivateKey, d *buildDocker, out io.Writer) Builder {
	t.Helper()
	inspect := &fakeImages{images: map[string]Info{BedrockTag(buildVersion): {ID: bedrockID}, BedrockTag("dev"): {ID: bedrockID}}}
	return Builder{
		Releases: selfupgrade.Options{Repo: "o/r", APIBase: f.srv.URL, Attempts: 1},
		Keys:     []*ecdsa.PublicKey{&key.PublicKey},
		Run:      Runner{Exec: d.exec, Out: out},
		Inspect:  inspect.inspect,
		Arch:     func() (string, error) { return "amd64", nil },
		TempDir:  t.TempDir(),
	}
}

// embedded is the Dockerfile this flai carries for image.
func embedded(t *testing.T, image string) []byte {
	t.Helper()
	data, err := images.Dockerfile(image)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// S-0740, ADR-0295 §1: a release build fetches the flai and pilot archives
// for the engine's architecture, checks each against the signed
// checksums.txt, puts their binaries in the context bedrock's prebuilt target
// is built from, then builds the harness image over bedrock with its base,
// its recipe and base labels, and every tag, and removes the context.
func TestBuildFromAReleaseBuildsBedrockFromItsBinaries(t *testing.T) {
	key := signingKey(t)
	f := serveRelease(t, goodRelease(t, key, []byte("flai binary"), []byte("pilot binary")))
	var out bytes.Buffer
	d := &buildDocker{}
	d.during = func(b build) {
		if len(d.builds) != 1 {
			return
		}
		ctxDir := b.args[len(b.args)-1]
		for name, want := range map[string]string{"flai": "flai binary", "pilot": "pilot binary"} {
			path := filepath.Join(ctxDir, "linux-amd64", name)
			got, err := os.ReadFile(path)
			if err != nil || string(got) != want {
				t.Errorf("the context holds %s = %q, %v; want %q", path, got, err, want)
			}
			if st, err := os.Stat(path); err != nil || st.Mode().Perm()&0o111 == 0 {
				t.Errorf("%s is not executable: %v, %v", path, st, err)
			}
		}
		for name, image := range map[string]string{"bedrock.Dockerfile": images.Bedrock, "claude-code.Dockerfile": "claude-code"} {
			if got, err := os.ReadFile(filepath.Join(ctxDir, name)); err != nil || !bytes.Equal(got, embedded(t, image)) {
				t.Errorf("the context's %s is not the embedded Dockerfile of %s: %v", name, image, err)
			}
		}
	}
	b := releaseBuilder(t, f, key, d, &out)
	tags := []string{harnessTag, "flai-agent-claude-code:local"}
	got, err := b.Build(context.Background(), "claude-code", BuildSource{Version: "v" + buildVersion}, tags)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	want := Built{Inner: "claude-code", Version: buildVersion, Tags: tags, Bedrock: "flai-agent-bedrock:1.2.3", From: FromRelease}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Build = %+v, want %+v", got, want)
	}
	for _, name := range []string{"checksums.txt", "checksums.txt.sig", "flai_1.2.3_linux_amd64.tar.gz", "pilot_1.2.3_linux_amd64.tar.gz"} {
		if n := f.count(name); n != 1 {
			t.Errorf("%s was fetched %d times, want once", name, n)
		}
	}
	if len(d.builds) != 2 {
		t.Fatalf("docker ran %d builds, want 2: %v", len(d.builds), d.builds)
	}
	ctxDir := d.builds[0].args[len(d.builds[0].args)-1]
	wantBedrock := []string{"build", "-f", filepath.Join(ctxDir, "bedrock.Dockerfile"), "--target", "prebuilt",
		"--build-arg", "FLAI_VERSION=1.2.3",
		"--label", LabelRecipe + "=" + Recipe(embedded(t, images.Bedrock), "prebuilt", nil),
		"-t", "flai-agent-bedrock:1.2.3", ctxDir}
	if !reflect.DeepEqual(d.builds[0].args, wantBedrock) || d.builds[0].dir != ctxDir {
		t.Errorf("bedrock was built in %s with\n%q\nwant in %s with\n%q", d.builds[0].dir, d.builds[0].args, ctxDir, wantBedrock)
	}
	wantHarness := []string{"build", "-f", filepath.Join(ctxDir, "claude-code.Dockerfile"),
		"--build-arg", "BASE=flai-agent-bedrock:1.2.3",
		"--label", LabelRecipe + "=" + Recipe(embedded(t, "claude-code"), "", nil),
		"--label", LabelBase + "=" + bedrockID,
		"-t", harnessTag, "-t", "flai-agent-claude-code:local", ctxDir}
	if !reflect.DeepEqual(d.builds[1].args, wantHarness) {
		t.Errorf("the harness image was built with\n%q\nwant\n%q", d.builds[1].args, wantHarness)
	}
	if left, err := os.ReadDir(b.TempDir); err != nil || len(left) != 0 {
		t.Errorf("the context was left behind: %v, %v", left, err)
	}
	for _, said := range []string{
		"building flai-agent-bedrock:1.2.3 from flai 1.2.3's release binaries for linux/amd64, checked against its signed checksums.txt",
		"building " + harnessTag + ", flai-agent-claude-code:local over flai-agent-bedrock:1.2.3 (" + bedrockID + ")",
	} {
		if !strings.Contains(out.String(), said) {
			t.Errorf("the build did not say %q:\n%s", said, out.String())
		}
	}
}

// S-0740, ADR-0295 §1: a checksums.txt whose signature does not verify
// against the trusted keys refuses the build, saying so, before any archive
// is fetched or any docker build runs.
func TestBuildFromAReleaseRefusesAnUnverifiedChecksums(t *testing.T) {
	key := signingKey(t)
	files := goodRelease(t, signingKey(t), []byte("flai"), []byte("pilot"))
	f := serveRelease(t, files)
	d := &buildDocker{}
	_, err := releaseBuilder(t, f, key, d, nil).Build(context.Background(), "claude-code", BuildSource{Version: buildVersion}, []string{harnessTag})
	var refused *selfupgrade.SignatureError
	if !errors.As(err, &refused) || refused.Kind != selfupgrade.SignatureBad {
		t.Fatalf("Build = %v, want a refusal for a bad signature", err)
	}
	if !strings.Contains(err.Error(), "refuse to build the agent images from release flai/v1.2.3") || !strings.Contains(err.Error(), "signature does not verify") {
		t.Errorf("the refusal does not say what it refused and why: %v", err)
	}
	if len(d.builds) != 0 {
		t.Errorf("docker ran %d builds after the refusal", len(d.builds))
	}
	if n := f.count("flai_1.2.3_linux_amd64.tar.gz") + f.count("pilot_1.2.3_linux_amd64.tar.gz"); n != 0 {
		t.Errorf("the archives were fetched %d times after the refusal", n)
	}
}

// S-0740, ADR-0295 §1: an archive whose SHA-256 is not the one its signed
// checksums.txt names refuses the build before any docker build runs.
func TestBuildFromAReleaseRefusesAnArchiveThatDoesNotMatch(t *testing.T) {
	key := signingKey(t)
	files := goodRelease(t, key, []byte("flai"), []byte("pilot"))
	files["pilot_1.2.3_linux_amd64.tar.gz"] = targz(t, "pilot", []byte("altered"))
	d := &buildDocker{}
	_, err := releaseBuilder(t, serveRelease(t, files), key, d, nil).Build(context.Background(), "claude-code", BuildSource{Version: buildVersion}, []string{harnessTag})
	if err == nil || !strings.Contains(err.Error(), "pilot_1.2.3_linux_amd64.tar.gz does not match its signed checksums.txt") {
		t.Fatalf("Build = %v, want a refusal naming the pilot archive", err)
	}
	if len(d.builds) != 0 {
		t.Errorf("docker ran %d builds after the refusal", len(d.builds))
	}
}

// writeTree writes a development tree under a new folder whose bedrock and
// claude-code Dockerfiles hold bedrock and inner, and is its root.
func writeTree(t *testing.T, bedrock, inner string) string {
	t.Helper()
	root := t.TempDir()
	for image, data := range map[string]string{images.Bedrock: bedrock, "claude-code": inner} {
		path := filepath.Join(root, filepath.FromSlash(images.Path(image)))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// treeBuilder is a Builder for a tree build: no releases, an Arch that must
// not be asked, and bedrock built as bedrockID.
func treeBuilder(t *testing.T, d *buildDocker) Builder {
	t.Helper()
	inspect := &fakeImages{images: map[string]Info{BedrockTag("dev"): {ID: bedrockID}}}
	return Builder{
		Run:     Runner{Exec: d.exec},
		Inspect: inspect.inspect,
		Arch: func() (string, error) {
			t.Error("a tree build asked the engine's architecture")
			return "", errors.New("not asked")
		},
	}
}

// S-0740, ADR-0295 §3: a tree build builds bedrock's tree target and the
// harness image from the tree's own Dockerfiles, in the tree's root, with
// FLAI_COMMIT and FLAI_DATE passed and its recipe their bytes', so a change
// of commit alone leaves the recipe as it was.
func TestBuildFromATreeUsesTheTreesDockerfiles(t *testing.T) {
	const bedrock, inner = "FROM scratch AS tree\n# the tree's bedrock\n", "ARG BASE\nFROM ${BASE}\n# the tree's claude-code\n"
	root := writeTree(t, bedrock, inner)
	d := &buildDocker{}
	src := BuildSource{Version: "dev", Tree: root, Commit: "abc1234", Date: "2026-11-13T10:00:00Z"}
	got, err := treeBuilder(t, d).Build(context.Background(), "claude-code", src, []string{"flai-agent-claude-code:local"})
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if got.From != FromTree || got.Bedrock != "flai-agent-bedrock:dev" {
		t.Errorf("Build = %+v, want from the tree over flai-agent-bedrock:dev", got)
	}
	if len(d.builds) != 2 {
		t.Fatalf("docker ran %d builds, want 2: %v", len(d.builds), d.builds)
	}
	recipe := Recipe([]byte(bedrock), "tree", nil)
	if recipe == Recipe(embedded(t, images.Bedrock), "tree", nil) {
		t.Fatal("the fixture's bedrock Dockerfile is the embedded one")
	}
	wantBedrock := []string{"build", "-f", filepath.Join(root, "flai", "images", "bedrock", "Dockerfile"), "--target", "tree",
		"--build-arg", "FLAI_VERSION=dev", "--build-arg", "FLAI_COMMIT=abc1234", "--build-arg", "FLAI_DATE=2026-11-13T10:00:00Z",
		"--label", LabelRecipe + "=" + recipe, "-t", "flai-agent-bedrock:dev", "."}
	if !reflect.DeepEqual(d.builds[0].args, wantBedrock) || d.builds[0].dir != root {
		t.Errorf("bedrock was built in %s with\n%q\nwant in %s with\n%q", d.builds[0].dir, d.builds[0].args, root, wantBedrock)
	}
	wantHarness := []string{"build", "-f", filepath.Join(root, "flai", "images", "claude-code", "Dockerfile"),
		"--build-arg", "BASE=flai-agent-bedrock:dev",
		"--label", LabelRecipe + "=" + Recipe([]byte(inner), "", nil),
		"--label", LabelBase + "=" + bedrockID,
		"-t", "flai-agent-claude-code:local", "."}
	if !reflect.DeepEqual(d.builds[1].args, wantHarness) || d.builds[1].dir != root {
		t.Errorf("the harness image was built in %s with\n%q\nwant in %s with\n%q", d.builds[1].dir, d.builds[1].args, root, wantHarness)
	}

	again := &buildDocker{}
	src.Commit, src.Date = "", ""
	if _, err := treeBuilder(t, again).Build(context.Background(), "claude-code", src, []string{"flai-agent-claude-code:local"}); err != nil {
		t.Fatalf("Build with no commit: %v", err)
	}
	wantBare := []string{"build", "-f", filepath.Join(root, "flai", "images", "bedrock", "Dockerfile"), "--target", "tree",
		"--build-arg", "FLAI_VERSION=dev", "--label", LabelRecipe + "=" + recipe, "-t", "flai-agent-bedrock:dev", "."}
	if !reflect.DeepEqual(again.builds[0].args, wantBare) {
		t.Errorf("with no commit or date, bedrock was built with\n%q\nwant\n%q", again.builds[0].args, wantBare)
	}
}

// S-0740: a tree build that misses a Dockerfile fails naming it, and builds nothing.
func TestBuildFromATreeNamesAMissingDockerfile(t *testing.T) {
	root := writeTree(t, "FROM scratch AS tree\n", "FROM scratch\n")
	missing := filepath.Join(root, "flai", "images", "claude-code", "Dockerfile")
	if err := os.Remove(missing); err != nil {
		t.Fatal(err)
	}
	d := &buildDocker{}
	_, err := treeBuilder(t, d).Build(context.Background(), "claude-code", BuildSource{Version: "dev", Tree: root}, []string{"flai-agent-claude-code:local"})
	if err == nil || !strings.Contains(err.Error(), missing) {
		t.Fatalf("Build = %v, want an error naming %s", err, missing)
	}
	if len(d.builds) != 0 {
		t.Errorf("docker ran %d builds", len(d.builds))
	}
}

// S-0740, ADR-0295 §3: a development flai with no tree has no release to
// take binaries from, and an inner harness flai carries no Dockerfile for,
// or bedrock itself, is refused; none of them builds anything.
func TestBuildRefusesWhatItCannotBuild(t *testing.T) {
	for _, c := range []struct {
		name, inner, version, want string
	}{
		{"development flai", "claude-code", "dev", "a development flai has no release to take binaries from; build in its tree"},
		{"snapshot flai", "claude-code", "0.0.0-abc1234", "a development flai has no release to take binaries from"},
		{"unknown inner", "other", buildVersion, "flai carries no Dockerfile for the agent image \"other\""},
		{"bedrock", images.Bedrock, buildVersion, "not an inner harness"},
	} {
		t.Run(c.name, func(t *testing.T) {
			d := &buildDocker{}
			b := Builder{Run: Runner{Exec: d.exec}, Arch: func() (string, error) { return "amd64", nil }}
			if _, err := b.Build(context.Background(), c.inner, BuildSource{Version: c.version}, []string{harnessTag}); err == nil || !strings.Contains(err.Error(), c.want) {
				t.Errorf("Build = %v, want %q", err, c.want)
			}
			if len(d.builds) != 0 {
				t.Errorf("docker ran %d builds", len(d.builds))
			}
		})
	}
}

// S-0740, ADR-0295 §1: EngineArch asks docker info for the engine's
// architecture, names x86_64 amd64 and aarch64 arm64, and refuses another,
// naming it, or a docker that is not on PATH.
func TestEngineArchNamesTheEnginesArchitectureAsBuildxDoes(t *testing.T) {
	for _, c := range []struct {
		out, want string
	}{
		{"x86_64", "amd64"},
		{"aarch64\n", "arm64"},
		{"WARNING: No swap limit support\nx86_64", "amd64"},
	} {
		r := &fakeRunner{path: []string{"docker"}, out: c.out}
		got, err := EngineArch(r, "docker")()
		if err != nil || got != c.want {
			t.Errorf("EngineArch with %q = %q, %v; want %q", c.out, got, err, c.want)
		}
		want := []string{"", "docker", "info", "--format", "{{.Architecture}}"}
		if len(r.runs) != 1 || !reflect.DeepEqual(r.runs[0], want) {
			t.Errorf("EngineArch ran %q, want %q", r.runs, want)
		}
	}
	if _, err := EngineArch(&fakeRunner{path: []string{"docker"}, out: "s390x"}, "docker")(); err == nil || !strings.Contains(err.Error(), `"s390x"`) {
		t.Errorf("EngineArch on s390x = %v, want a refusal naming it", err)
	}
	if _, err := EngineArch(&fakeRunner{}, "docker")(); !errors.Is(err, ErrNoDocker) {
		t.Errorf("EngineArch with no docker = %v, want ErrNoDocker", err)
	}
	if _, err := EngineArch(&fakeRunner{path: []string{"docker"}, out: "Cannot connect to the Docker daemon", err: errExit}, "docker")(); err == nil || !strings.Contains(err.Error(), "Cannot connect to the Docker daemon") {
		t.Errorf("EngineArch with no engine = %v, want what docker said", err)
	}
}

// BedrockTag names the bedrock image of a version without a registry.
func TestBedrockTag(t *testing.T) {
	if got := BedrockTag(buildVersion); got != "flai-agent-bedrock:1.2.3" {
		t.Errorf("BedrockTag = %q", got)
	}
}
