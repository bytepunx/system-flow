package cmd

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/host"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// fakeRunner records docker calls and answers from a script.
type fakeRunner struct {
	calls    []string
	missing  bool            // docker not on PATH
	images   map[string]bool // images present
	running  map[string]bool // containers running
	private  map[string]bool // images that need a login to pull
	loggedIn bool
	noToken  bool   // gh has no token
	identity bool   // git config has a user
	mounts   string // mount destinations of a running container, one per line

	imageIDs       map[string]string // ref -> the image ID a pull or inspect of it reports now; default sha256:<ref>
	containerRef   map[string]string // container name -> the ref it was last started with: docker run's image argument, its {{.Config.Image}}
	containerLabel map[string]string // container name -> its dashboardRefLabel; none for a container an older flai started
	containerImage map[string]string // container name -> the image ID it was started from (a snapshot, not re-derived, so a later change to imageIDs is a real difference a check can find)
	probeAddr      string            // host:port `docker port` reports for the upgrade probe container; default 127.0.0.1:19999
	health         map[string]string // container name -> docker's HEALTHCHECK verdict; default none
	env            map[string]string // container name -> its environment, one VAR=value per line

	// repoDigests are each image ID's RepoDigests; by default, the ref it was
	// pulled by when that is by digest, as a pull by digest leaves it.
	repoDigests map[string][]string

	// removeDelay stands in for a --rm container Docker is still removing
	// after docker stop: the name stays listed by that many `docker ps --all`
	// queries, and a docker run under it meanwhile fails with exit 125 and a
	// name conflict (I-0094). lingering counts what is left of it per name.
	removeDelay int
	lingering   map[string]int
	runFails    map[string]bool // "name ref" -> docker run of that container from that ref fails

	daemonPlatform string            // what `docker version` reports the server runs; default linux/amd64
	imagePlatforms map[string]string // ref -> the os/arch of the local image; default the daemon's
	published      []string          // platforms the registry has an image for; default any
}

// platform is the daemon's os/arch.
func (f *fakeRunner) platform() string {
	if f.daemonPlatform != "" {
		return f.daemonPlatform
	}
	return "linux/amd64"
}

// idFor is the image ID a pull or inspect of ref reports right now.
func (f *fakeRunner) idFor(ref string) string {
	if id, ok := f.imageIDs[ref]; ok {
		return id
	}
	return "sha256:" + ref
}

// repoDigestsOf is the RepoDigests docker image inspect reports for image ID
// id, as JSON: f.repoDigests' entry, else the by-digest reference whose image
// it is, as a pull by digest leaves it, else none.
func (f *fakeRunner) repoDigestsOf(id string) string {
	digests, ok := f.repoDigests[id]
	if !ok {
		refs := []string{strings.TrimPrefix(id, "sha256:")}
		for ref, refID := range f.imageIDs {
			if refID == id {
				refs = append(refs, ref)
			}
		}
		for _, ref := range refs {
			if strings.Contains(ref, "@sha256:") && f.idFor(ref) == id {
				digests = []string{ref}
			}
		}
	}
	data, _ := json.Marshal(append([]string{}, digests...))
	return string(data)
}

func (f *fakeRunner) RunInput(dir, name, input string, args ...string) (string, error) {
	if name == "docker" && args[0] == "login" {
		f.calls = append(f.calls, "docker "+strings.Join(args, " ")+" <"+strings.TrimSpace(input)+">")
		if strings.TrimSpace(input) == "" {
			return "", fmt.Errorf("docker login: empty password")
		}
		f.loggedIn = true
		return "Login Succeeded", nil
	}
	return f.Run(dir, name, args...)
}

func (f *fakeRunner) LookPath(name string) (string, error) {
	if name == "docker" && f.missing {
		return "", errors.New("not found")
	}
	return "/usr/bin/" + name, nil
}

func (f *fakeRunner) Run(dir, name string, args ...string) (string, error) {
	f.calls = append(f.calls, name+" "+strings.Join(args, " "))
	if name == "gh" {
		switch strings.Join(args, " ") {
		case "auth token":
			if f.noToken {
				return "", fmt.Errorf("gh auth token: not logged in")
			}
			return "ghp_fake", nil
		case "api user --jq .login":
			return "tester", nil
		}
		return "", nil
	}
	if name == "git" {
		switch {
		case args[0] == "rev-parse":
			return "abc1234", nil
		case args[0] == "config" && f.identity && args[1] == "user.name":
			return "Test User", nil
		case args[0] == "config" && f.identity && args[1] == "user.email":
			return "test@example.com", nil
		}
		return "", fmt.Errorf("git: no tags")
	}
	if name != "docker" {
		return "", nil
	}
	switch args[0] {
	case "version":
		return f.platform(), nil
	case "image":
		ref := args[2]
		if len(args) >= 5 && args[3] == "--format" && args[4] == "{{json .RepoDigests}}" {
			return f.repoDigestsOf(ref), nil
		}
		if !f.images[ref] {
			return "", fmt.Errorf("docker image inspect %s: exit status 1\nError: No such image", ref)
		}
		if len(args) >= 5 && args[3] == "--format" && args[4] == "{{.Id}}" {
			return f.idFor(ref), nil
		}
		if len(args) >= 5 && args[3] == "--format" && args[4] == "{{.Os}}/{{.Architecture}}" {
			if p, ok := f.imagePlatforms[ref]; ok {
				return p, nil
			}
			return f.platform(), nil
		}
		return "[]", nil
	case "pull":
		ref, platform := args[len(args)-1], ""
		if len(args) >= 5 && args[2] == "--platform" {
			platform = args[3]
		}
		if f.private[ref] && !f.loggedIn {
			return "", fmt.Errorf("docker pull --quiet %s: exit status 1\nError response from daemon: error from registry: unauthorized", ref)
		}
		if f.published != nil && !slices.Contains(f.published, platform) {
			return "", fmt.Errorf("docker pull --quiet %s: exit status 1\nError response from daemon: no matching manifest for %s in the manifest list entries", ref, platform)
		}
		f.images[ref] = true
		if f.imagePlatforms == nil {
			f.imagePlatforms = map[string]string{}
		}
		f.imagePlatforms[ref] = orDefault(platform, f.platform())
		return ref, nil
	case "build":
		f.images[args[4]] = true
		return "sha256:built", nil
	case "run":
		name, ref := args[4], args[len(args)-1]
		if f.lingering[name] > 0 {
			return "", fmt.Errorf("docker %s: exit status 125\ndocker: Error response from daemon: Conflict. The container name \"/%s\" is already in use by container \"0123456789abcdef\". You have to remove (or rename) that container to be able to reuse that name.", strings.Join(args, " "), name)
		}
		if f.runFails[name+" "+ref] {
			return "", fmt.Errorf("docker %s: exit status 125\ndocker: Error response from daemon: driver failed programming external connectivity on endpoint %s", strings.Join(args, " "), name)
		}
		f.running[name] = true
		if f.containerRef == nil {
			f.containerRef = map[string]string{}
		}
		if f.containerImage == nil {
			f.containerImage = map[string]string{}
		}
		if f.containerLabel == nil {
			f.containerLabel = map[string]string{}
		}
		f.containerRef[name] = ref
		delete(f.containerLabel, name)
		for i, arg := range args[:len(args)-1] {
			if v, ok := strings.CutPrefix(args[i+1], dashboardRefLabel+"="); arg == "--label" && ok {
				f.containerLabel[name] = v
			}
		}
		f.containerImage[name] = f.idFor(ref)
		if strings.HasPrefix(ref, "sha256:") {
			f.containerImage[name] = ref // started from an image ID, not a tag
		}
		return "0123456789abcdef", nil
	case "ps":
		name := strings.TrimSuffix(strings.TrimPrefix(args[len(args)-1], "name=^/"), "$")
		if f.running[name] {
			return "0123456789ab", nil
		}
		if slices.Contains(args, "--all") && f.lingering[name] > 0 {
			f.lingering[name]--
			return "0123456789ab", nil
		}
		return "", nil
	case "inspect":
		if len(args) >= 3 && args[1] == "--format" && args[2] == "{{.Image}}" {
			return f.containerImage[args[len(args)-1]], nil
		}
		if strings.Contains(strings.Join(args, " "), ".Mounts") {
			return f.mounts, nil
		}
		if strings.Contains(strings.Join(args, " "), ".State.Health") {
			return f.health[args[len(args)-1]], nil
		}
		if strings.Contains(strings.Join(args, " "), ".Config.Env") {
			return f.env[args[len(args)-1]], nil
		}
		if len(args) >= 3 && args[1] == "--format" && args[2] == containerInfoFormat {
			name := args[len(args)-1]
			if ref, ok := f.containerRef[name]; ok {
				return ref + "|5555|" + f.containerLabel[name], nil
			}
			return "ghcr.io/bytepunx/flaiover:0.2.0|5555|", nil
		}
		return "", nil
	case "stop":
		delete(f.running, args[1])
		if f.removeDelay > 0 {
			if f.lingering == nil {
				f.lingering = map[string]int{}
			}
			f.lingering[args[1]] = f.removeDelay
		}
		return args[1], nil
	case "rm":
		delete(f.running, args[len(args)-1])
		delete(f.lingering, args[len(args)-1])
		return "", nil
	case "port":
		if f.probeAddr != "" {
			return f.probeAddr, nil
		}
		return "127.0.0.1:19999", nil
	case "logs":
		return "hello from container", nil
	}
	return "", nil
}

func runWith(t *testing.T, dir string, r *fakeRunner, args ...string) (string, string, int) {
	t.Helper()
	return runWithApp(t, &app{cwd: dir, runner: r}, args...)
}

// runWithApp is runWith for a test that needs to set more of app than the
// runner: an upgrade test fakes healthProbe and sleep so a health check
// that never succeeds does not make the test wait for real.
func runWithApp(t *testing.T, a *app, args ...string) (string, string, int) {
	t.Helper()
	t.Setenv("FLAI_CACHE_DIR", filepath.Join(t.TempDir(), "cache"))
	var out, errOut bytes.Buffer
	a.out, a.errOut = &out, &errOut
	if a.clock == nil {
		a.clock = func() time.Time { return time.Date(2026, 9, 17, 1, 0, 0, 0, time.UTC) }
	}
	if a.hostStarter == nil {
		// No test starts a real flai host: the test binary is not flai.
		a.hostStarter = func() (host.Status, bool, error) { return host.Status{PID: 4242}, true, nil }
	}
	if a.sleep == nil {
		a.sleep = func(time.Duration) {} // tests never wait for a real interval
	}
	if a.healthProbe == nil {
		a.healthProbe = func(string) bool { return true } // tests never reach the network; a running container answers
	}
	if dashboardReleasesOf != t {
		lazySignedDashboardReleases(t) // nor does flai dashboard, for a test that names no releases of its own
	}
	root := newRootCmdWith(a)
	root.SetArgs(args)
	code := 0
	if err := root.Execute(); err != nil {
		var ee *exitError
		if errors.As(err, &ee) {
			code = ee.code
		} else {
			a.fail(err)
			code = 1
		}
	}
	return out.String(), errOut.String(), code
}

// testReleaseVersions are the flaiover releases the dashboard tests' stand-in
// publishes, newest first.
var testReleaseVersions = []string{"0.10.0", "0.4.0", "0.3.0", "0.2.0"}

// flaioverDigest is the index digest release version's test list names.
func flaioverDigest(version string) string {
	sum := sha256.Sum256([]byte("index " + version))
	return "sha256:" + hex.EncodeToString(sum[:])
}

// flaioverRef is what flai dashboard pulls and runs for release version.
func flaioverRef(version string) string {
	return signedDashboardImage + "@" + flaioverDigest(version)
}

// signedListOf is release version's digest list, its index at
// flaioverDigest(version), beside a platform manifest.
func signedListOf(version string) []byte {
	ref := signedDashboardImage + ":" + version
	platform := sha256.Sum256([]byte("linux/amd64 " + version))
	return []byte(flaioverDigest(version) + " " + ref + " index\n" +
		"sha256:" + hex.EncodeToString(platform[:]) + " " + ref + " linux/amd64\n")
}

// dashboardReleasesOf is the test dashboardReleases answers for, nil when
// none has set it.
var dashboardReleasesOf *testing.T

// setDashboardReleases makes flai dashboard resolve its releases through opt
// for the rest of t.
func setDashboardReleases(t *testing.T, opt func() selfupgrade.Options) {
	t.Helper()
	old, oldOf := dashboardReleases, dashboardReleasesOf
	t.Cleanup(func() { dashboardReleases, dashboardReleasesOf = old, oldOf })
	dashboardReleases, dashboardReleasesOf = opt, t
}

// releaseKey is a fresh key for a test's releases to be signed with.
func releaseKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

// signedReleases are testReleaseVersions, each publishing signedListOf it
// signed by key.
func signedReleases(t *testing.T, key *ecdsa.PrivateKey) []flaioverRelease {
	t.Helper()
	rels := make([]flaioverRelease, len(testReleaseVersions))
	for i, v := range testReleaseVersions {
		rels[i] = signedFlaiover(t, v, signedListOf(v), key)
	}
	return rels
}

// dashboardReleasesAt are the options for o/r's releases at api, verified
// with trusted.
func dashboardReleasesAt(api string, trusted *ecdsa.PrivateKey) selfupgrade.Options {
	return selfupgrade.Options{Repo: "o/r", APIBase: api, Keys: []*ecdsa.PublicKey{&trusted.PublicKey}}
}

// serveDashboardReleases serves rels as o/r's flaiover releases and makes flai
// dashboard resolve them for the rest of t, trusting trusted; asked lists the
// paths asked for so far.
func serveDashboardReleases(t *testing.T, trusted *ecdsa.PrivateKey, rels ...flaioverRelease) (asked func() []string) {
	t.Helper()
	tags := make([]string, len(rels))
	for i, r := range rels {
		tags[i] = selfupgrade.DashboardTagPrefix + r.version
	}
	api, asked := dashboardStandIn(t, tags, rels...)
	setDashboardReleases(t, func() selfupgrade.Options { return dashboardReleasesAt(api, trusted) })
	return asked
}

// useSignedDashboardReleases serves testReleaseVersions, each with its list
// signed by a key flai dashboard trusts, for the rest of t.
func useSignedDashboardReleases(t *testing.T) (asked func() []string) {
	t.Helper()
	key := releaseKey(t)
	return serveDashboardReleases(t, key, signedReleases(t, key)...)
}

// lazySignedDashboardReleases is useSignedDashboardReleases, started the first
// time flai dashboard asks for its releases in t: every dashboard test gets it
// from runWithApp unless it sets its own.
func lazySignedDashboardReleases(t *testing.T) {
	t.Helper()
	var opt *selfupgrade.Options
	setDashboardReleases(t, func() selfupgrade.Options {
		if opt == nil {
			key := releaseKey(t)
			tags := make([]string, len(testReleaseVersions))
			for i, v := range testReleaseVersions {
				tags[i] = selfupgrade.DashboardTagPrefix + v
			}
			api, _ := dashboardStandIn(t, tags, signedReleases(t, key)...)
			o := dashboardReleasesAt(api, key)
			opt = &o
		}
		return *opt
	})
}

// lastRunOf is the last docker run f was asked for, "" when none.
func lastRunOf(f *fakeRunner) string { return lastRun(f.calls) }

// called reports whether any of f's calls starts with prefix.
func called(f *fakeRunner, prefix string) bool {
	return slices.ContainsFunc(f.calls, func(c string) bool { return strings.HasPrefix(c, prefix) })
}

// S-0236, ADR-0070: flai dashboard resolves dashboard.tag to a flaiover
// release, verifies its digest list, keeps it in the cache beside flai
// serve's state, and pulls and runs the image by the list's index digest.
func TestDashboardRunsTheReleaseByDigest(t *testing.T) {
	root := dashboardProject(t)
	asked := useSignedDashboardReleases(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard")
	if code != 0 {
		t.Fatalf("start: %d %s", code, errOut)
	}
	ref := flaioverRef("0.10.0")
	if !slices.Contains(f.calls, "docker pull --quiet --platform linux/amd64 "+ref) {
		t.Errorf("latest is pulled by the newest release's index digest:\n%s", strings.Join(f.calls, "\n"))
	}
	if run := lastRunOf(f); !strings.HasSuffix(run, " --label "+dashboardRefLabel+"="+ref+" "+ref) {
		t.Errorf("run by digest, labelled with it: %s", run)
	}
	if !strings.Contains(out, "image "+ref) || !strings.Contains(out, "flaiover 0.10.0, run by digest from its verified digest list") || strings.Contains(out, "unsigned") {
		t.Errorf("output: %s", out)
	}
	dir := string((&app{}).serveDir())
	cached, err := os.ReadFile(filepath.Join(dir, "dashboard-digests", "0.10.0.digests"))
	if err != nil || string(cached) != string(signedListOf("0.10.0")) {
		t.Errorf("the verified list is cached beside flai serve's state: %v %q", err, cached)
	}
	if _, err := os.Stat(filepath.Join(dir, "dashboard-digests", "0.10.0.digests.sig")); err != nil {
		t.Errorf("with its signature: %v", err)
	}
	if r, _ := readDashboardRecord(dir); r.Ref != ref {
		t.Errorf("the record keeps the by-digest reference: %+v", r)
	}
	if !slices.Contains(asked(), "/repos/o/r/releases") {
		t.Errorf("the releases are asked: %v", asked())
	}

	// a version is that release; --json says what runs and that it is signed
	g := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code = runWith(t, root, g, "dashboard", "--tag", "0.3.0", "--json")
	var js map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &js) != nil {
		t.Fatalf("start 0.3.0: %d %s %s", code, out, errOut)
	}
	if js["image"] != flaioverRef("0.3.0") || js["signed"] != true || js["release"] != "0.3.0" || js["cached_list"] != nil {
		t.Errorf("--json: %v", js)
	}
	if !slices.Contains(g.calls, "docker pull --quiet --platform linux/amd64 "+flaioverRef("0.3.0")) {
		t.Errorf("0.3.0 by digest:\n%s", strings.Join(g.calls, "\n"))
	}
}

// S-0236, ADR-0070: a release with no list this flai can verify, or a tag that
// names no release, starts no container: nothing is pulled or run, nothing is
// cached, and the error says why.
func TestDashboardRefusesAReleaseWithNoListItCanVerify(t *testing.T) {
	root := dashboardProject(t)
	key := releaseKey(t)
	noList := signedFlaiover(t, "0.10.0", signedListOf("0.10.0"), key)
	delete(noList.assets, "flaiover_0.10.0.digests")
	altered := signedFlaiover(t, "0.4.0", signedListOf("0.4.0"), key)
	altered.assets["flaiover_0.4.0.digests"] = signedListOf("0.4.1")
	unknown := signedFlaiover(t, "0.3.0", signedListOf("0.3.0"), releaseKey(t))
	serveDashboardReleases(t, key, noList, altered, unknown)
	for tag, why := range map[string]string{
		"latest": "flaiover 0.10.0 is not started",
		"0.4.0":  "the list may have been altered since it was signed",
		"0.3.0":  "which this flai does not know",
		"0.9.9":  "has no published release tagged flaiover/v0.9.9",
		"edge":   "is not a flaiover release version",
	} {
		f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
		out, errOut, code := runWith(t, root, f, "dashboard", "--tag", tag)
		if code == 0 || !strings.Contains(errOut, why) {
			t.Errorf("%s refused: %d %s %s", tag, code, out, errOut)
		}
		if called(f, "docker pull") || called(f, "docker run") {
			t.Errorf("%s: nothing is pulled or run:\n%s", tag, strings.Join(f.calls, "\n"))
		}
	}
	if entries, _ := os.ReadDir(filepath.Join(string((&app{}).serveDir()), "dashboard-digests")); len(entries) != 0 {
		t.Errorf("nothing refused is cached: %v", entries)
	}
}

// S-0236: when the releases cannot be reached, the cached verified list of the
// release runs, and flai dashboard says so; with none cached nothing starts.
func TestDashboardRunsTheCachedListWhenTheReleasesCannotBeReached(t *testing.T) {
	root := dashboardProject(t)
	key := releaseKey(t)
	serveDashboardReleases(t, key, signedReleases(t, key)...)
	if _, errOut, code := runWith(t, root, &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}, "dashboard", "--tag", "0.4.0"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	setDashboardReleases(t, func() selfupgrade.Options { return dashboardReleasesAt("http://127.0.0.1:1", key) })
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard", "--tag", "0.4.0")
	if code != 0 || !strings.Contains(out, "flaiover 0.4.0, run by digest from its verified digest list cached beside flai serve's state: the release could not be reached") {
		t.Fatalf("cached: %d %s %s", code, out, errOut)
	}
	if run := lastRunOf(f); !strings.HasSuffix(run, " "+flaioverRef("0.4.0")) {
		t.Errorf("the cached release by digest: %s", run)
	}
	g := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	_, errOut, code = runWith(t, root, g, "dashboard", "--tag", "0.3.0")
	if code == 0 || !strings.Contains(errOut, "no list of it is cached") || !strings.Contains(errOut, "no dashboard container was started") || called(g, "docker run") {
		t.Errorf("nothing cached: %d %s", code, errOut)
	}
}

// S-0236: --image naming another repository, a configured image of another
// repository, and --build run what they name, as before, asking no release,
// and say the container is unsigned.
func TestDashboardRunsAnotherImageUnsigned(t *testing.T) {
	root := dashboardProject(t)
	asked := useSignedDashboardReleases(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard", "--image", "registry.example/team/flaiover", "--tag", "1.2")
	if code != 0 || !strings.Contains(out, "unsigned: no signed digest list names registry.example/team/flaiover:1.2") {
		t.Fatalf("--image: %d %s %s", code, out, errOut)
	}
	if run := lastRunOf(f); !strings.HasSuffix(run, " registry.example/team/flaiover:1.2") || !slices.Contains(f.calls, "docker pull --quiet --platform linux/amd64 registry.example/team/flaiover:1.2") {
		t.Errorf("runs what --image names: %s", run)
	}
	js, _, _ := runWith(t, mirrorProject(t), &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}, "dashboard", "--json")
	if !strings.Contains(js, `"signed": false`) || !strings.Contains(js, `"image": "`+mirrorLatest+`"`) || strings.Contains(js, `"release"`) {
		t.Errorf("a configured image of another repository: %s", js)
	}
	if n := len(asked()); n != 0 {
		t.Errorf("an unsigned image asks no release: %v", asked())
	}
}

// S-0236: status shows the running image's repository digest, the release
// whose cached verified list names it, and whether it is on a signed list.
func TestDashboardStatusShowsTheDigestAndItsRelease(t *testing.T) {
	root := dashboardProject(t)
	useSignedDashboardReleases(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, errOut, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	type status struct {
		Image   string `json:"image"`
		Digest  string `json:"digest"`
		Release string `json:"release"`
		Signed  bool   `json:"signed"`
	}
	read := func() (status, string) {
		t.Helper()
		js, errOut, code := runWith(t, root, f, "dashboard", "status", "--json")
		var st status
		if code != 0 || json.Unmarshal([]byte(js), &st) != nil {
			t.Fatalf("status --json: %d %s %s", code, js, errOut)
		}
		out, _, _ := runWith(t, root, f, "dashboard", "status")
		return st, out
	}
	st, out := read()
	if st != (status{Image: flaioverRef("0.10.0"), Digest: flaioverDigest("0.10.0"), Release: "0.10.0", Signed: true}) {
		t.Errorf("signed: %+v", st)
	}
	if !strings.Contains(out, "digest: "+flaioverDigest("0.10.0")+" (flaiover 0.10.0, on a signed list)") {
		t.Errorf("signed text: %s", out)
	}

	// the entry a cached list names is the one shown, whichever comes first;
	// 0.4.0's list is cached once a start resolves it
	if _, errOut, code := runWith(t, root, &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}, "dashboard", "--tag", "0.4.0"); code != 0 {
		t.Fatalf("start 0.4.0: %s", errOut)
	}
	id := f.containerImage["flaiover"]
	f.repoDigests = map[string][]string{id: {"registry.example/flaiover@sha256:" + strings.Repeat("c", 64), flaioverRef("0.4.0")}}
	if st, _ = read(); st.Digest != flaioverDigest("0.4.0") || st.Release != "0.4.0" || !st.Signed {
		t.Errorf("the listed entry: %+v", st)
	}

	// a digest no list names, or a listed digest under another repository, is unsigned
	f.repoDigests = map[string][]string{id: {"registry.example/flaiover@" + flaioverDigest("0.10.0")}}
	st, out = read()
	if st.Digest != flaioverDigest("0.10.0") || st.Release != "" || st.Signed {
		t.Errorf("another repository: %+v", st)
	}
	if !strings.Contains(out, "digest: "+flaioverDigest("0.10.0")+" (unsigned: no signed digest list names it)") {
		t.Errorf("unsigned text: %s", out)
	}
	f.repoDigests = map[string][]string{id: {}}
	st, out = read()
	if st.Digest != "" || st.Signed || !strings.Contains(out, "unsigned: its image has no repository digest") {
		t.Errorf("no digest: %+v %s", st, out)
	}
	js, _, _ := runWith(t, root, f, "dashboard", "status", "--json")
	if strings.Contains(js, `"digest"`) || strings.Contains(js, `"release"`) {
		t.Errorf("digest and release are left out when there is none: %s", js)
	}
}

// S-0236, ADR-0118: flai dashboard restart starts the release it ran again by
// digest, asking no release and pulling nothing.
func TestDashboardRestartKeepsTheDigest(t *testing.T) {
	root := dashboardProject(t)
	asked := useSignedDashboardReleases(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, errOut, code := runWith(t, root, f, "dashboard", "--tag", "0.4.0"); code != 0 {
		t.Fatalf("start: %s", errOut)
	}
	ref, ranImage, before := flaioverRef("0.4.0"), f.containerImage["flaiover"], len(asked())
	f.calls = nil
	out, errOut, code := runWith(t, root, f, "dashboard", "restart", "--json")
	if code != 0 || !strings.Contains(out, `"ref": "`+ref+`"`) {
		t.Fatalf("restart: %d %s %s", code, out, errOut)
	}
	if f.containerLabel["flaiover"] != ref || f.containerImage["flaiover"] != ranImage || called(f, "docker pull") || len(asked()) != before {
		t.Errorf("restart keeps the digest: label %s image %s\n%s", f.containerLabel["flaiover"], f.containerImage["flaiover"], strings.Join(f.calls, "\n"))
	}
	if r, _ := readDashboardRecord(string((&app{}).serveDir())); r.Ref != ref {
		t.Errorf("record after restart: %+v", r)
	}
}

func TestDashboardLifecycle(t *testing.T) {
	cfg := filepath.Join(t.TempDir(), "cfg.json")
	t.Setenv("FLAI_CONFIG", cfg)
	serveDir := string(serve.DirFor(cfg))
	root := tempProject(t)
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: My Proj\nkey: m\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  port: 5555\n"), 0o644)

	missing := &fakeRunner{missing: true, images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, missing, "dashboard")
	if code == 0 || !strings.Contains(errOut, "docker is required") || !strings.Contains(errOut, "Install Docker") {
		t.Fatalf("missing docker: %d %s", code, errOut)
	}

	_ = out
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code = runWith(t, root, f, "dashboard", "--tag", "0.2.0")
	if code != 0 {
		t.Fatalf("run: %s", errOut)
	}
	joined := strings.Join(f.calls, "\n")
	for _, want := range []string{
		"docker image inspect " + flaioverRef("0.2.0"),
		"docker pull --quiet --platform linux/amd64 " + flaioverRef("0.2.0"),
		"--name flaiover",
		"--publish 0.0.0.0:5555:3000",
		"--mount type=bind,source=" + filepath.Join(serveDir, "dashboard.token") + ",target=/run/secrets/flaiover_token,readonly",
		"--env FLAIOVER_TOKEN_FILE=/run/secrets/flaiover_token",
		"--user " + fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"--label " + dashboardRefLabel + "=" + flaioverRef("0.2.0") + " " + flaioverRef("0.2.0"),
	} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing %q in calls:\n%s", want, joined)
		}
	}
	if !strings.Contains(out, "http://localhost:5555") || !strings.Contains(out, "flai dashboard stop") || !strings.Contains(out, "http://localhost:5555/login#token=") || !strings.Contains(out, "no file of any project") {
		t.Errorf("output: %s", out)
	}
	// ADR-0031: a port and two secrets, and nothing of the project or of the
	// person who started it.
	var run string
	for _, c := range f.calls {
		if strings.HasPrefix(c, "docker run ") {
			run = c
		}
	}
	for _, gone := range []string{"--volume", root + ":", "PROJECT_DIR", "GIT_", ".git/", "push_key", "known_hosts", "/etc/passwd", "gitignore"} {
		if strings.Contains(run, gone) {
			t.Errorf("%q has no place in the container's arguments:\n%s", gone, run)
		}
	}
	if n := strings.Count(run, "--mount "); n != 2 || strings.Count(run, "target=/run/secrets/") != 2 || strings.Count(run, ",readonly") != 2 {
		t.Errorf("the only mounts are the two secrets, read-only:\n%s", run)
	}
	tokenFile := filepath.Join(serveDir, "dashboard.token")
	tokenData, err := os.ReadFile(tokenFile)
	if err != nil || len(strings.TrimSpace(string(tokenData))) < 40 {
		t.Fatalf("token file: %v %q", err, tokenData)
	}
	if st, _ := os.Stat(tokenFile); st.Mode().Perm() != 0o600 {
		t.Errorf("token file mode %v", st.Mode().Perm())
	}
	if strings.Contains(errOut, strings.TrimSpace(string(tokenData))) {
		t.Error("token must not be logged")
	}
	// token command: stable, then rotated with a restart of the running container
	out, _, code = runWith(t, root, f, "dashboard", "token")
	if code != 0 || !strings.Contains(out, "token: "+strings.TrimSpace(string(tokenData))) || !strings.Contains(out, "login: http://localhost:5555/login#token=") {
		t.Errorf("token: %d %s", code, out)
	}
	out, _, code = runWith(t, root, f, "dashboard", "token", "--rotate")
	rotated, _ := os.ReadFile(tokenFile)
	if code != 0 || string(rotated) == string(tokenData) || !strings.Contains(out, "token: "+strings.TrimSpace(string(rotated))) || !strings.Contains(strings.Join(f.calls, "\n"), "docker restart flaiover") || !strings.Contains(out, "restarted flaiover") {
		t.Errorf("rotate: %d %s", code, out)
	}
	// S-0105: rotated with the dashboard left running, which takes the new token from the answer
	before := len(f.calls)
	out, _, code = runWith(t, root, f, "dashboard", "token", "--rotate", "--no-restart", "--json")
	again, _ := os.ReadFile(tokenFile)
	if code != 0 || string(again) == string(rotated) || !strings.Contains(out, `"restarted": false`) || !strings.Contains(out, strings.TrimSpace(string(again))) || strings.Contains(strings.Join(f.calls[before:], "\n"), "docker restart") {
		t.Errorf("rotate --no-restart: %d %s", code, out)
	}
	// precedence: flag port beats manifest port; manifest port beat config's 4242
	out, _, _ = runWith(t, root, &fakeRunner{images: map[string]bool{flaioverRef("0.10.0"): true}, running: map[string]bool{}}, "dashboard", "--port", "6000", "--json")
	if !strings.Contains(out, `"url": "http://localhost:6000"`) {
		t.Errorf("flag precedence: %s", out)
	}
	// the host's git identity stays on the host: flai serve commits there (S-0075)
	idr := &fakeRunner{images: map[string]bool{flaioverRef("0.10.0"): true}, running: map[string]bool{}, identity: true}
	_, _, _ = runWith(t, root, idr, "dashboard")
	if j := strings.Join(idr.calls, "\n"); strings.Contains(j, "Test User") || strings.Contains(j, "git config") {
		t.Errorf("the identity is not read, let alone passed:\n%s", j)
	}
	// --bind restricts the published address; the manifest can set it too
	bnd := &fakeRunner{images: map[string]bool{flaioverRef("0.10.0"): true}, running: map[string]bool{}}
	out, _, _ = runWith(t, root, bnd, "dashboard", "--bind", "127.0.0.1")
	if !strings.Contains(strings.Join(bnd.calls, "\n"), "--publish 127.0.0.1:5555:3000") || !strings.Contains(out, "reachable from this host only") {
		t.Errorf("--bind: %s\n%s", out, strings.Join(bnd.calls, "\n"))
	}
	// present image is not pulled unless --pull
	g := &fakeRunner{images: map[string]bool{flaioverRef("0.10.0"): true}, running: map[string]bool{}}
	_, _, _ = runWith(t, root, g, "dashboard")
	if strings.Contains(strings.Join(g.calls, "\n"), "docker pull") {
		t.Error("pulled although present")
	}
	h := &fakeRunner{images: map[string]bool{flaioverRef("0.10.0"): true}, running: map[string]bool{}}
	_, _, _ = runWith(t, root, h, "dashboard", "--pull")
	if !strings.Contains(strings.Join(h.calls, "\n"), "docker pull --quiet --platform linux/amd64 "+flaioverRef("0.10.0")) {
		t.Error("--pull should pull")
	}
	// already running for this same project: still registers, says so, starts no second container
	out, _, _ = runWith(t, root, f, "dashboard")
	if !strings.Contains(out, "already runs") || !strings.Contains(out, "now also serves") || strings.Count(strings.Join(f.calls, "\n"), "docker run ") != 1 {
		t.Errorf("second start: %s\n%s", out, strings.Join(f.calls, "\n"))
	}
	out, _, _ = runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "running at http://localhost:5555 ("+flaioverRef("0.2.0")+")") {
		t.Errorf("status must report the running container: %s", out)
	}
	if strings.Contains(out, "older flai") || strings.Contains(out, "push") || strings.Contains(out, "read-only") {
		t.Errorf("status says nothing of keys or git paths any more: %s", out)
	}
	// a container an older flai started still has the project mounted: say so
	f.mounts = "/run/secrets/flaiover_token\n" + root + "\n" + root + "/.git/hooks\n"
	out, _, _ = runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "started by an older flai") || !strings.Contains(out, root) || !strings.Contains(out, "flai dashboard stop") {
		t.Errorf("a stale container: %s", out)
	}
	js, _, _ := runWith(t, root, f, "dashboard", "status", "--json")
	if !strings.Contains(js, `"stale_mounts"`) || strings.Contains(js, "git_read_only") || strings.Contains(js, "push_key") {
		t.Errorf("status --json: %s", js)
	}
	f.mounts = ""
	out, _, _ = runWith(t, root, f, "dashboard", "logs")
	if !strings.Contains(out, "hello from container") {
		t.Errorf("logs: %s", out)
	}
	out, _, code = runWith(t, root, f, "dashboard", "stop")
	if code != 0 || !strings.Contains(out, "stopped flaiover") || f.running["flaiover"] {
		t.Errorf("stop: %d %s", code, out)
	}
	out, _, _ = runWith(t, root, f, "dashboard", "stop")
	if !strings.Contains(out, "not running") {
		t.Errorf("stop twice: %s", out)
	}
}

// S-0119: the pull asks for the daemon's platform, a present image built for
// another platform is pulled again, and a registry with nothing for the
// platform gets an error that names it.
func TestDashboardPullsTheDaemonsPlatform(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	img := flaioverRef("0.10.0")

	f := &fakeRunner{daemonPlatform: "linux/arm64", images: map[string]bool{}, running: map[string]bool{}}
	_, errOut, code := runWith(t, root, f, "dashboard", "--pull")
	if joined := strings.Join(f.calls, "\n"); code != 0 || !strings.Contains(joined, "docker pull --quiet --platform linux/arm64 "+img) {
		t.Fatalf("the pull should name the daemon's platform: %d %s\n%s", code, errOut, joined)
	}

	// An amd64 image pulled before an arm64 one was published is replaced.
	g := &fakeRunner{daemonPlatform: "linux/arm64", images: map[string]bool{img: true}, imagePlatforms: map[string]string{img: "linux/amd64"}, running: map[string]bool{}}
	_, errOut, code = runWith(t, root, g, "dashboard")
	if joined := strings.Join(g.calls, "\n"); code != 0 || !strings.Contains(joined, "docker pull --quiet --platform linux/arm64 "+img) || g.imagePlatforms[img] != "linux/arm64" {
		t.Fatalf("an image for another platform should be pulled again: %d %s\n%s", code, errOut, joined)
	}

	// A present image for the daemon's platform is reused.
	h := &fakeRunner{daemonPlatform: "linux/arm64", images: map[string]bool{img: true}, imagePlatforms: map[string]string{img: "linux/arm64"}, running: map[string]bool{}}
	_, _, _ = runWith(t, root, h, "dashboard")
	if strings.Contains(strings.Join(h.calls, "\n"), "docker pull") {
		t.Errorf("pulled although present for this platform:\n%s", strings.Join(h.calls, "\n"))
	}

	// Nothing published for the platform: say which, and what to do.
	m := &fakeRunner{daemonPlatform: "linux/arm64", published: []string{"linux/amd64"}, images: map[string]bool{}, running: map[string]bool{}}
	_, errOut, code = runWith(t, root, m, "dashboard")
	if code == 0 || !strings.Contains(errOut, "has no image for linux/arm64") || !strings.Contains(errOut, "--build") || strings.Contains(strings.Join(m.calls, "\n"), "docker run") {
		t.Errorf("no image for the platform: %d %s", code, errOut)
	}
}

func TestDashboardPrivateRegistryAndBuild(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GH_TOKEN", "")
	root := tempProject(t)
	img := flaioverRef("0.10.0")

	// Unauthorized pull, token from gh: login and retry.
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}, private: map[string]bool{img: true}}
	out, errOut, code := runWith(t, root, f, "dashboard")
	joined := strings.Join(f.calls, "\n")
	if code != 0 || !strings.Contains(joined, "docker login ghcr.io --username tester --password-stdin <ghp_fake>") || strings.Count(joined, "docker pull --quiet --platform linux/amd64 "+img) != 2 || !strings.Contains(out, "http://localhost:4242") {
		t.Fatalf("login retry: %d %s %s\n%s", code, out, errOut, joined)
	}
	if strings.Contains(errOut, "ghp_fake") {
		t.Error("token must not be logged")
	}

	// No token anywhere: the error names the fix.
	g := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}, private: map[string]bool{img: true}, noToken: true}
	_, errOut, code = runWith(t, root, g, "dashboard")
	if code == 0 || !strings.Contains(errOut, "read:packages") || strings.Contains(strings.Join(g.calls, "\n"), "docker login") {
		t.Errorf("no token: %d %s", code, errOut)
	}

	// Docker Hub images get no login attempt.
	h := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}, private: map[string]bool{"nginx:latest": true}}
	_, errOut, code = runWith(t, root, h, "dashboard", "--image", "nginx")
	if code == 0 || strings.Contains(strings.Join(h.calls, "\n"), "docker login") {
		t.Errorf("hub image: %d %s", code, errOut)
	}

	// --build outside the monorepo refuses; inside it builds and runs flaiover:local.
	b := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	_, errOut, code = runWith(t, root, b, "dashboard", "--build")
	if code == 0 || !strings.Contains(errOut, "flaiover/Dockerfile") {
		t.Errorf("build without dockerfile: %d %s", code, errOut)
	}
	_ = os.MkdirAll(filepath.Join(root, "flaiover"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "flaiover", "Dockerfile"), []byte("FROM scratch\n"), 0o644)
	out, errOut, code = runWith(t, root, b, "dashboard", "--build")
	joined = strings.Join(b.calls, "\n")
	if code != 0 || !strings.Contains(joined, "docker build -f "+filepath.Join(root, "flaiover", "Dockerfile")+" -t flaiover:local --build-arg FLAI_VERSION=dev --build-arg FLAI_COMMIT=abc1234") || !strings.Contains(joined, "flaiover:local") || strings.Contains(joined, "docker pull") || !strings.Contains(out, "image flaiover:local") {
		t.Errorf("build: %d %s %s\n%s", code, out, errOut, joined)
	}
}

// S-0235: only release-flaiover.yml stamps an image, on a release tag; a
// dashboard built here passes no stamp, so its image has none.
func TestDashboardBuildPassesNoReleaseStamp(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	_ = os.MkdirAll(filepath.Join(root, "flaiover"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "flaiover", "Dockerfile"), []byte("FROM scratch\n"), 0o644)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	_, errOut, code := runWith(t, root, f, "dashboard", "--build")
	if code != 0 {
		t.Fatalf("build: %d %s", code, errOut)
	}
	builds := 0
	for _, call := range f.calls {
		if !strings.HasPrefix(call, "docker build ") {
			continue
		}
		builds++
		if !strings.Contains(call, "--build-arg FLAIOVER_VERSION=") {
			t.Errorf("the build passes no build arguments at all, so this test checks nothing: %s", call)
		}
		if strings.Contains(call, "FLAIOVER_STAMP_") {
			t.Errorf("flai dashboard --build passed a release stamp: %s", call)
		}
	}
	if builds != 1 {
		t.Errorf("want one docker build, got %d:\n%s", builds, strings.Join(f.calls, "\n"))
	}
}

// S-0077: the push key is retired. Asked for, it is named, explained, and
// ignored, and the dashboard starts.
func TestDashboardPushKeyIsRetired(t *testing.T) {
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	img := map[string]bool{flaioverRef("0.10.0"): true}
	f := &fakeRunner{images: img, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard", "--push-key", "/home/me/.ssh/deploy")
	if code != 0 || !strings.Contains(out, "--push-key: retired and ignored") || !strings.Contains(out, "flai release --pending") || strings.Contains(out, "flai push") || !strings.Contains(out, "ADR-0031") {
		t.Fatalf("flag: %d %s %s", code, out, errOut)
	}
	if j := strings.Join(f.calls, "\n"); strings.Contains(j, "deploy") || strings.Contains(j, "ssh-keygen") {
		t.Errorf("the key is not looked at or passed:\n%s", j)
	}
	if _, _, code := runWith(t, root, f, "config", "set", "dashboard.push_key", "/home/me/.ssh/deploy"); code != 0 {
		t.Fatal("an old config still loads, and the key can still be set, so it can be cleared")
	}
	out, _, _ = runWith(t, root, &fakeRunner{images: img, running: map[string]bool{}}, "dashboard")
	if !strings.Contains(out, "dashboard.push_key: retired and ignored") {
		t.Errorf("config: %s", out)
	}
	js, _, _ := runWith(t, root, &fakeRunner{images: img, running: map[string]bool{}}, "dashboard", "--json")
	if !strings.Contains(js, `"retired"`) || strings.Contains(js, `"push_key"`) || strings.Contains(js, `"mount"`) {
		t.Errorf("--json: %s", js)
	}
	if help, _, _ := runWith(t, root, f, "dashboard", "--help"); strings.Contains(help, "push-key") {
		t.Error("the retired flags are hidden")
	}
}

// T-0322: flai dashboard restart, check, upgrade.

func dashboardProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: My Proj\nkey: m\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  port: 5555\n  tag: latest\n"), 0o644)
	return root
}

// mirrorImage is an image of a repository other than the signed one, which
// flai dashboard runs by its tag, unsigned (S-0236); mirrorLatest is its tag
// the mirror project configures.
const (
	mirrorImage  = "registry.example/flaiover"
	mirrorLatest = mirrorImage + ":latest"
)

// mirrorProject is dashboardProject configured with mirrorImage, for the
// tests of a tag that can name another image tomorrow (I-0116), which only an
// unsigned image has, and of check and upgrade, which compare by tag.
func mirrorProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	root := tempProject(t)
	_ = os.WriteFile(filepath.Join(root, "system-flow.yaml"), []byte("version: 1\nname: My Proj\nkey: m\nlayout:\n  design: design\n  docs: docs\n  wip: wip\ndashboard:\n  image: "+mirrorImage+"\n  port: 5555\n  tag: latest\n"), 0o644)
	return root
}

func TestDashboardRestartStartsWhenNotRunning(t *testing.T) {
	root := dashboardProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard", "restart")
	if code != 0 {
		t.Fatalf("restart: %d %s", code, errOut)
	}
	if !f.running["flaiover"] {
		t.Error("not started")
	}
	if !strings.Contains(out, "started flaiover") {
		t.Errorf("output: %s", out)
	}
}

func TestDashboardRestartCyclesTheRunningContainerWithoutChangingItsImage(t *testing.T) {
	root := mirrorProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	ranImage := f.containerImage["flaiover"]
	f.imageIDs = map[string]string{mirrorLatest: "sha256:newer"} // a newer build under the same tag
	f.calls = nil                                                // only the restart's own calls matter below
	out, errOut, code := runWith(t, root, f, "dashboard", "restart")
	if code != 0 {
		t.Fatalf("restart: %d %s", code, errOut)
	}
	if !strings.Contains(out, "restarted flaiover") {
		t.Errorf("output: %s", out)
	}
	if f.containerImage["flaiover"] != ranImage {
		t.Errorf("restart runs %s, want the image it was already running, %s (a plain restart must never silently upgrade)", f.containerImage["flaiover"], ranImage)
	}
	joined := strings.Join(f.calls, "\n")
	if strings.Count(joined, "docker pull") != 0 {
		t.Errorf("restart pulled; it should reuse the local image:\n%s", joined)
	}
}

// I-0116: flai dashboard check pulls the tag the dashboard runs, so after a
// newer latest is pulled the tag names another image than the one that runs.
// A restart starts the image that runs, by its ID, and still names the tag.
func TestDashboardRestartRunsTheImageIDNotWhatTheTagNamesNow(t *testing.T) {
	root := mirrorProject(t)
	tag := mirrorLatest
	f := &fakeRunner{images: map[string]bool{tag: true}, running: map[string]bool{}, imageIDs: map[string]string{tag: "sha256:old"}}
	if _, errOut, code := runWith(t, root, f, "dashboard"); code != 0 || f.containerImage["flaiover"] != "sha256:old" {
		t.Fatalf("start: %d %s", code, errOut)
	}
	f.imageIDs[tag] = "sha256:new" // what the pull of a newer latest leaves
	if out, errOut, code := runWith(t, root, f, "dashboard", "check", "--json"); code != 0 || !strings.Contains(out, `"upgrade_available": true`) {
		t.Fatalf("check: %d %s %s", code, out, errOut)
	}
	for _, again := range []string{"first", "second"} {
		f.calls = nil
		out, errOut, code := runWith(t, root, f, "dashboard", "restart", "--json")
		if code != 0 {
			t.Fatalf("%s restart: %d %s", again, code, errOut)
		}
		var runs []string
		for _, c := range f.calls {
			if strings.HasPrefix(c, "docker run ") {
				runs = append(runs, c)
			}
		}
		if len(runs) != 1 || !strings.HasSuffix(runs[0], " sha256:old") || !strings.Contains(runs[0], " --label "+dashboardRefLabel+"="+tag+" ") {
			t.Errorf("%s restart runs the image that ran, by its ID, labelled with its tag:\n%s", again, strings.Join(f.calls, "\n"))
		}
		if joined := strings.Join(f.calls, "\n"); strings.Contains(joined, "docker pull") {
			t.Errorf("%s restart pulled:\n%s", again, joined)
		}
		if !strings.Contains(out, `"ref": "`+tag+`"`) || !strings.Contains(out, `"image": "sha256:old"`) || f.containerImage["flaiover"] != "sha256:old" {
			t.Errorf("%s restart: runs %s, says %s", again, f.containerImage["flaiover"], out)
		}
	}
	out, _, _ := runWith(t, root, f, "dashboard", "restart")
	if !strings.Contains(out, "restarted flaiover ("+tag+")") {
		t.Errorf("restart names the tag, not the image ID: %s", out)
	}
	out, _, _ = runWith(t, root, f, "dashboard", "status")
	if !strings.Contains(out, "("+tag+")") || strings.Contains(out, "sha256:") {
		t.Errorf("status names the tag, not the image ID: %s", out)
	}
}

// A running dashboard whose image ID docker does not say is started again
// from the reference it shows, as before I-0116.
func TestDashboardRestartFallsBackToTheRefWhenTheImageIDIsUnknown(t *testing.T) {
	root := dashboardProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{"flaiover": true}}
	out, errOut, code := runWith(t, root, f, "dashboard", "restart")
	if code != 0 || !strings.Contains(out, "restarted flaiover (ghcr.io/bytepunx/flaiover:0.2.0)") || f.containerRef["flaiover"] != "ghcr.io/bytepunx/flaiover:0.2.0" || f.containerLabel["flaiover"] != "ghcr.io/bytepunx/flaiover:0.2.0" {
		t.Errorf("restart: %d %s %s\n%s", code, out, errOut, strings.Join(f.calls, "\n"))
	}
}

// containerInfo names the reference the label records, which a container
// started from an image ID needs, and a container an older flai started
// without the label by its {{.Config.Image}}.
func TestContainerInfoPrefersTheRefLabel(t *testing.T) {
	f := &fakeRunner{
		containerRef:   map[string]string{"labelled": "sha256:abc", "older": "ghcr.io/bytepunx/flaiover:0.1.0"},
		containerLabel: map[string]string{"labelled": "ghcr.io/bytepunx/flaiover:latest"},
	}
	a := &app{runner: f}
	if image, url := a.containerInfo("labelled"); image != "ghcr.io/bytepunx/flaiover:latest" || url != "http://localhost:5555" {
		t.Errorf("labelled: %q %q", image, url)
	}
	if image, url := a.containerInfo("older"); image != "ghcr.io/bytepunx/flaiover:0.1.0" || url != "http://localhost:5555" {
		t.Errorf("no label: %q %q", image, url)
	}
}

// containerInfo reads its fields by position, so an empty one, such as a
// container that publishes no port, does not shift the others.
func TestContainerInfoReadsEmptyFields(t *testing.T) {
	a := &app{runner: &answerRunner{out: "ghcr.io/bytepunx/flaiover:0.1.0||ghcr.io/bytepunx/flaiover:latest\n"}}
	if image, url := a.containerInfo("flaiover"); image != "ghcr.io/bytepunx/flaiover:latest" || url != "" {
		t.Errorf("no port: %q %q", image, url)
	}
	a = &app{runner: &answerRunner{out: "ghcr.io/bytepunx/flaiover:0.1.0\n"}}
	if image, url := a.containerInfo("flaiover"); image != "ghcr.io/bytepunx/flaiover:0.1.0" || url != "" {
		t.Errorf("one field: %q %q", image, url)
	}
}

// answerRunner answers every command with out.
type answerRunner struct{ out string }

func (r *answerRunner) Run(string, string, ...string) (string, error) { return r.out, nil }

func (r *answerRunner) RunInput(string, string, string, ...string) (string, error) {
	return r.out, nil
}

func (r *answerRunner) LookPath(name string) (string, error) { return "/usr/bin/" + name, nil }

func TestDashboardCheckReportsNotRunningPlainly(t *testing.T) {
	root := dashboardProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard", "check")
	if code != 0 {
		t.Fatalf("check: %d %s", code, errOut)
	}
	if !strings.Contains(out, "not running") {
		t.Errorf("output: %s", out)
	}
}

func TestDashboardCheckFindsNothingWhenAlreadyCurrent(t *testing.T) {
	root := mirrorProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	out, errOut, code := runWith(t, root, f, "dashboard", "check", "--json")
	if code != 0 {
		t.Fatalf("check: %d %s", code, errOut)
	}
	if !strings.Contains(out, `"upgrade_available": false`) {
		t.Errorf("output: %s", out)
	}
}

func TestDashboardCheckFindsAnUpdate(t *testing.T) {
	root := mirrorProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	f.imageIDs = map[string]string{mirrorLatest: "sha256:newer"}
	out, errOut, code := runWith(t, root, f, "dashboard", "check", "--json")
	if code != 0 {
		t.Fatalf("check: %d %s", code, errOut)
	}
	if !strings.Contains(out, `"upgrade_available": true`) {
		t.Errorf("output: %s", out)
	}
}

func TestDashboardUpgradeStartsWhenNotRunning(t *testing.T) {
	root := dashboardProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	out, errOut, code := runWith(t, root, f, "dashboard", "upgrade")
	if code != 0 {
		t.Fatalf("upgrade: %d %s", code, errOut)
	}
	if !f.running["flaiover"] || !strings.Contains(out, "started flaiover") {
		t.Errorf("not started: %s %s", out, errOut)
	}
}

func TestDashboardUpgradeReportsAlreadyUpToDate(t *testing.T) {
	root := mirrorProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	a := &app{cwd: root, runner: f, healthProbe: func(string) bool { t.Fatal("no probe needed when already current"); return false }}
	out, errOut, code := runWithApp(t, a, "dashboard", "upgrade")
	if code != 0 {
		t.Fatalf("upgrade: %d %s", code, errOut)
	}
	if !strings.Contains(out, "already running") {
		t.Errorf("output: %s", out)
	}
	if f.running[dashboardUpgradeProbeName] {
		t.Error("no temporary container should have been started")
	}
}

// dashboardUpgradeProbeName is the temporary container's name for the
// fixture project's fixed container name.
const dashboardUpgradeProbeName = "flaiover" + upgradeProbeSuffix

func TestDashboardUpgradeSwapsToAHealthyNewImage(t *testing.T) {
	root := mirrorProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	f.imageIDs = map[string]string{mirrorLatest: "sha256:newer"}
	var probed []string
	a := &app{cwd: root, runner: f, healthProbe: func(url string) bool { probed = append(probed, url); return true }}
	out, errOut, code := runWithApp(t, a, "dashboard", "upgrade")
	if code != 0 {
		t.Fatalf("upgrade: %d %s", code, errOut)
	}
	if !strings.Contains(out, "upgraded flaiover") {
		t.Errorf("output: %s", out)
	}
	if len(probed) != 1 || !strings.HasSuffix(probed[0], "/_health") || !strings.HasPrefix(probed[0], "http://127.0.0.1:") {
		t.Errorf("health probe: %v", probed)
	}
	if f.running[dashboardUpgradeProbeName] {
		t.Error("the temporary container should have been stopped")
	}
	if !f.running["flaiover"] {
		t.Error("the real container should be running again")
	}
	if f.containerImage["flaiover"] != "sha256:newer" {
		t.Errorf("running image = %s, want the new one", f.containerImage["flaiover"])
	}
}

func TestDashboardUpgradeLeavesThePreviousContainerRunningWhenTheNewOneNeverAnswersHealthy(t *testing.T) {
	root := mirrorProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	f.imageIDs = map[string]string{mirrorLatest: "sha256:broken"}
	a := &app{cwd: root, runner: f, healthProbe: func(string) bool { return false }}
	out, errOut, code := runWithApp(t, a, "dashboard", "upgrade")
	if code == 0 {
		t.Fatalf("upgrade should have failed: %s", out)
	}
	if !strings.Contains(errOut, "keeps running unchanged") {
		t.Errorf("errOut: %s", errOut)
	}
	if !f.running["flaiover"] {
		t.Fatal("the previous container must still be running")
	}
	if f.containerImage["flaiover"] == "sha256:broken" {
		t.Error("the real container must not have been touched")
	}
	if f.running[dashboardUpgradeProbeName] {
		t.Error("the temporary container should have been cleaned up")
	}
}

// I-0094: the dashboard runs with --rm, so after docker stop Docker still
// holds its name until it has removed it, and a docker run under that name
// at once fails with a conflict. An upgrade waits for the name first.
func TestDashboardUpgradeSwapsWhileDockerStillRemovesTheStoppedContainer(t *testing.T) {
	root := mirrorProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	f.imageIDs = map[string]string{mirrorLatest: "sha256:newer"}
	f.removeDelay = 3
	f.calls = nil
	out, errOut, code := runWith(t, root, f, "dashboard", "upgrade")
	if code != 0 {
		t.Fatalf("upgrade: %d %s", code, errOut)
	}
	if !strings.Contains(out, "upgraded flaiover") || !f.running["flaiover"] || f.containerImage["flaiover"] != "sha256:newer" {
		t.Errorf("not upgraded: %s", out)
	}
	if slices.Contains(f.calls, "docker rm -f flaiover") {
		t.Errorf("a name Docker releases in time is not removed by force:\n%s", strings.Join(f.calls, "\n"))
	}
}

func TestDashboardRestartWaitsWhileDockerStillRemovesTheStoppedContainer(t *testing.T) {
	root := dashboardProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	f.removeDelay = 3
	out, errOut, code := runWith(t, root, f, "dashboard", "restart")
	if code != 0 {
		t.Fatalf("restart: %d %s", code, errOut)
	}
	if !strings.Contains(out, "restarted flaiover") || !f.running["flaiover"] {
		t.Errorf("not restarted: %s", out)
	}
}

func TestDashboardUpgradeRemovesANameDockerNeverReleases(t *testing.T) {
	root := mirrorProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	f.imageIDs = map[string]string{mirrorLatest: "sha256:newer"}
	f.removeDelay = nameReleaseAttempts + 5 // still listed when the wait gives up
	f.calls = nil
	slept := 0
	a := &app{cwd: root, runner: f, sleep: func(time.Duration) { slept++ }}
	out, errOut, code := runWithApp(t, a, "dashboard", "upgrade")
	if code != 0 {
		t.Fatalf("upgrade: %d %s", code, errOut)
	}
	if !strings.Contains(out, "upgraded flaiover") || f.containerImage["flaiover"] != "sha256:newer" {
		t.Errorf("not upgraded: %s", out)
	}
	asked := 0
	for _, c := range f.calls {
		if c == "docker ps --all --quiet --filter name=^/flaiover$" {
			asked++
		}
	}
	if asked != nameReleaseAttempts || slept != nameReleaseAttempts-1 {
		t.Errorf("asked docker %d times and slept %d, want %d and %d", asked, slept, nameReleaseAttempts, nameReleaseAttempts-1)
	}
	// the start is followed only by reading the image ID it recorded (I-0116)
	if n := len(f.calls); f.calls[n-3] != "docker rm -f flaiover" || !strings.HasPrefix(f.calls[n-2], "docker run ") || f.calls[n-1] != "docker inspect --format {{.Image}} flaiover" {
		t.Errorf("the name is removed by force, then the new image started:\n%s", strings.Join(f.calls, "\n"))
	}
}

func TestDashboardUpgradeStartsThePreviousImageAgainWhenTheNewOneFailsToStart(t *testing.T) {
	root := mirrorProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	previous := f.containerImage["flaiover"]
	f.imageIDs = map[string]string{mirrorLatest: "sha256:newer"}
	f.runFails = map[string]bool{"flaiover " + mirrorLatest: true}
	out, errOut, code := runWith(t, root, f, "dashboard", "upgrade", "--json")
	if code == 0 {
		t.Fatalf("the upgrade failed, so it exits non-zero, --json or not: %s", out)
	}
	for _, want := range []string{"upgrade to " + mirrorLatest + " failed", "did not start after flaiover stopped", "driver failed programming", "flaiover runs the previous image " + mirrorLatest + " again"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("missing %q in: %s", want, errOut)
		}
	}
	if !f.running["flaiover"] || f.containerImage["flaiover"] != previous {
		t.Errorf("flaiover should run the previous image %s again (by its ID, as the pull moved the tag), runs %s", previous, f.containerImage["flaiover"])
	}
	if f.containerLabel["flaiover"] != mirrorLatest {
		t.Errorf("the previous image is labelled with the tag it ran as, so status names that: %q", f.containerLabel["flaiover"])
	}
}

func TestDashboardUpgradeSaysTheDashboardIsDownWhenThePreviousImageFailsToStartToo(t *testing.T) {
	root := mirrorProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	previous := f.containerImage["flaiover"]
	f.imageIDs = map[string]string{mirrorLatest: "sha256:newer"}
	f.runFails = map[string]bool{"flaiover " + mirrorLatest: true, "flaiover " + previous: true}
	_, errOut, code := runWith(t, root, f, "dashboard", "upgrade")
	if code == 0 {
		t.Fatal("upgrade should have failed")
	}
	for _, want := range []string{"upgrade to " + mirrorLatest + " failed", "starting the previous image " + mirrorLatest + " again failed too", "the dashboard is down: run flai dashboard to start it"} {
		if !strings.Contains(errOut, want) {
			t.Errorf("missing %q in: %s", want, errOut)
		}
	}
	if f.running["flaiover"] {
		t.Error("nothing started, so nothing runs")
	}
}

func TestDashboardRestartSaysTheDashboardIsDownWhenItDoesNotStartAgain(t *testing.T) {
	root := dashboardProject(t)
	f := &fakeRunner{images: map[string]bool{}, running: map[string]bool{}}
	if _, _, code := runWith(t, root, f, "dashboard"); code != 0 {
		t.Fatal("start")
	}
	ran := f.containerImage["flaiover"]
	f.runFails = map[string]bool{"flaiover " + ran: true}
	_, errOut, code := runWith(t, root, f, "dashboard", "restart")
	if code == 0 || !strings.Contains(errOut, "restart flaiover: it stopped but did not start again from "+flaioverRef("0.10.0")+" (image "+ran+")") || !strings.Contains(errOut, "the dashboard is down: run flai dashboard to start it") {
		t.Errorf("restart: %d %s", code, errOut)
	}
}
