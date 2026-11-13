package agentimage

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/images"
	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/selfupgrade"
)

// The targets of bedrock's Dockerfile a build names (ADR-0295 §1).
const (
	targetPrebuilt = "prebuilt"
	targetTree     = "tree"
)

// The ways Build makes the images, as Built.From names them.
const (
	// FromRelease is a build from a flai release's Linux binaries and the Dockerfiles flai carries.
	FromRelease = "release"
	// FromTree is a build from a development tree's own Dockerfiles and source.
	FromTree = "tree"
)

// BuildSource is what a build makes the harness images from: a flai
// release's Linux binaries with the Dockerfiles this flai carries, or, with
// Tree set, a development tree's own Dockerfiles and source (ADR-0295 §1, §3).
type BuildSource struct {
	Version string // the flai version the images must hold, labelled dev.system-flow.flai
	Tree    string // a development tree's root, which holds flai/images/<image>/Dockerfile; "" for a release
	Commit  string // the tree's commit, labelled dev.system-flow.commit; "" for a release
	Date    string // the tree commit's date, passed as FLAI_DATE
}

// Builder builds the harness images on the host.
type Builder struct {
	Releases selfupgrade.Options    // flai's releases, for a release's binaries and checksums
	Keys     []*ecdsa.PublicKey     // the keys checksums.txt must verify against; selfupgrade.TrustedKeys(Releases) when empty
	Run      Runner                 // runs each docker build
	Inspect  Inspector              // reads the bedrock image's ID for the harness image's base label; Docker through execx.System when nil
	Arch     func() (string, error) // the engine's architecture as buildx names it, amd64 or arm64; EngineArch through execx.System when nil
	TempDir  string                 // where a release build's context is made; os.TempDir() when ""
}

// Built is what Build made: the inner harness, the version, the tags given
// its image, the bedrock image's tag it was built over, and From,
// FromRelease or FromTree.
type Built struct {
	Inner   string
	Version string
	Tags    []string
	Bedrock string
	From    string
}

// BedrockTag is the tag a build gives the bedrock image of version: flai-agent-bedrock:<version>.
func BedrockTag(version string) string {
	return "flai-agent-bedrock:" + version
}

// EngineArch is the Arch that asks program, the docker client, through r:
// docker info --format '{{.Architecture}}', x86_64 as amd64 and aarch64 as
// arm64, any other refused naming it.
func EngineArch(r execx.Runner, program string) func() (string, error) {
	return func() (string, error) {
		if _, err := r.LookPath(program); err != nil {
			return "", fmt.Errorf("ask the engine its architecture: %s is not on PATH (%w): install Docker, or name its client with flai serve agent harness container --program <path>", program, ErrNoDocker)
		}
		out, err := r.Run("", program, "info", "--format", "{{.Architecture}}")
		if err != nil {
			said := lastLine(out)
			if said == "" {
				said = firstLine(err.Error())
			}
			return "", fmt.Errorf("%s info could not say the engine's architecture (%s): start Docker, then build again", program, said)
		}
		// docker info may print warnings beside the answer.
		for _, line := range strings.Split(out, "\n") {
			switch strings.TrimSpace(line) {
			case "x86_64", "amd64":
				return "amd64", nil
			case "aarch64", "arm64":
				return "arm64", nil
			}
		}
		return "", fmt.Errorf("the engine's architecture is %q, as %s info says it, and flai's releases carry Linux binaries for amd64 and arm64 only: build on an amd64 or arm64 engine, or name an image of your own with flai serve agent harness container --image <inner>=<image>", lastLine(out), program)
	}
}

// Build builds bedrock, then inner's harness image over it, tagging the
// harness image with each of tags: from src.Tree's own Dockerfiles and
// source when it names a tree, else from the Dockerfiles this flai carries
// and the Linux flai and pilot of the release src.Version, checked against
// the release's signed checksums.txt (ADR-0295 §1, §3). It says what each
// phase does on b.Run.Out before it runs it.
func (b Builder) Build(ctx context.Context, inner string, src BuildSource, tags []string) (Built, error) {
	if inner == images.Bedrock {
		return Built{}, fmt.Errorf("build the agent image of %s: bedrock is the base every harness image builds over, not an inner harness; name an inner harness, such as claude-code", inner)
	}
	if _, err := images.Dockerfile(inner); err != nil {
		return Built{}, err
	}
	if len(tags) == 0 {
		return Built{}, fmt.Errorf("build the agent image of %s: no tag was given for it; name at least one, such as flai-agent-%s:<version>", inner, inner)
	}
	version := strings.TrimPrefix(src.Version, "v")
	if version == "" {
		return Built{}, fmt.Errorf("build the agent image of %s: no flai version was given for the images to hold; name the version, such as the one flai version prints", inner)
	}
	built := Built{Inner: inner, Version: version, Tags: slices.Clone(tags), Bedrock: BedrockTag(version)}
	if src.Tree != "" {
		built.From = FromTree
		return built, b.fromTree(ctx, inner, src, version, tags)
	}
	if buildinfo.Development(version) {
		return Built{}, fmt.Errorf("build the agent image of %s for flai %s: a development flai has no release to take binaries from; build in its tree", inner, version)
	}
	built.From = FromRelease
	return built, b.fromRelease(ctx, inner, version, tags)
}

// fromRelease builds bedrock's prebuilt target from the release version's
// Linux flai and pilot for the engine's architecture, in a context folder of
// its own that holds them and the Dockerfiles this flai carries, then inner's
// harness image over it, and removes the folder.
func (b Builder) fromRelease(ctx context.Context, inner, version string, tags []string) error {
	out := b.out()
	arch, err := b.arch()()
	if err != nil {
		return err
	}
	if arch != "amd64" && arch != "arm64" {
		return fmt.Errorf("build the agent images of flai %s: the engine's architecture is %q, and flai's releases carry Linux binaries for amd64 and arm64 only; build on an amd64 or arm64 engine", version, arch)
	}
	rel, err := selfupgrade.Find(ctx, b.Releases, selfupgrade.TagPrefix, version)
	if err != nil {
		return fmt.Errorf("find flai %s's release, for its Linux flai and pilot: %w", version, err)
	}
	fmt.Fprintf(out, "agent image: fetching flai and pilot for linux/%s from the release %s, checked against its signed checksums.txt\n", arch, rel.Tag)
	bins, err := b.releaseBinaries(ctx, rel, version, arch)
	if err != nil {
		return err
	}
	bedrockFile, err := images.Dockerfile(images.Bedrock)
	if err != nil {
		return err
	}
	innerFile, err := images.Dockerfile(inner)
	if err != nil {
		return err
	}
	dir, err := os.MkdirTemp(b.TempDir, "flai-agent-image-")
	if err != nil {
		return fmt.Errorf("make a context folder for the agent images of flai %s: %w; free space in the temporary folder, or name another", version, err)
	}
	defer func() { _ = os.RemoveAll(dir) }()
	files := map[string][]byte{
		images.Bedrock + ".Dockerfile": bedrockFile,
		inner + ".Dockerfile":          innerFile,
	}
	for name, data := range bins {
		files[filepath.Join("linux-"+arch, name)] = data
	}
	if err := writeContext(dir, files); err != nil {
		return fmt.Errorf("write the context of the agent images of flai %s: %w; free space in the temporary folder, or name another", version, err)
	}
	bedrock := BedrockTag(version)
	fmt.Fprintf(out, "agent image: building %s from flai %s's release binaries for linux/%s, checked against its signed checksums.txt\n", bedrock, version, arch)
	args := []string{"build", "-f", filepath.Join(dir, images.Bedrock+".Dockerfile"),
		"--target", targetPrebuilt,
		"--build-arg", "FLAI_VERSION=" + version,
		"--label", LabelRecipe + "=" + Recipe(bedrockFile, targetPrebuilt, nil),
		"-t", bedrock, dir}
	if err := b.Run.Run(ctx, Step{Name: "build " + bedrock, Dir: dir, Args: args}); err != nil {
		return err
	}
	return b.harness(ctx, out, harnessBuild{
		inner: inner, bedrock: bedrock, tags: tags, dir: dir, context: dir,
		file: filepath.Join(dir, inner+".Dockerfile"), data: innerFile, from: "the Dockerfile this flai carries",
	})
}

// releaseBinaries are the flai and pilot for linux/arch of rel, version's
// release, by name, each from its archive once the archive matches
// rel's checksums.txt and that file's signature verifies against b's keys.
func (b Builder) releaseBinaries(ctx context.Context, rel selfupgrade.Release, version, arch string) (map[string][]byte, error) {
	keys := b.Keys
	if len(keys) == 0 {
		var err error
		if keys, err = selfupgrade.TrustedKeys(b.Releases); err != nil {
			return nil, fmt.Errorf("verify the signature of release %s, for its Linux flai and pilot: %w", rel.Tag, err)
		}
	}
	sums, _, err := selfupgrade.VerifySigned(ctx, b.Releases, rel, "checksums.txt", keys)
	var refused *selfupgrade.SignatureError
	if errors.As(err, &refused) {
		return nil, &refusal{err: refused, msg: fmt.Sprintf("refuse to build the agent images from release %s: %s (%s), so its binaries cannot be trusted; build nothing from it by other means, and install a release this flai can verify, or report it to the repository's maintainers", rel.Tag, refused.Summary(), refused.Reason)}
	}
	if err != nil {
		return nil, fmt.Errorf("fetch the signed checksums.txt of release %s, for its Linux flai and pilot: %w", rel.Tag, err)
	}
	archives := []struct{ binary, archive string }{
		{"flai", selfupgrade.ArchiveName(version, "linux", arch)},
		{"pilot", fmt.Sprintf("pilot_%s_linux_%s.tar.gz", version, arch)},
	}
	bins := map[string][]byte{}
	for _, a := range archives {
		data, err := selfupgrade.FetchAsset(ctx, b.Releases, rel, a.archive)
		if err != nil {
			return nil, fmt.Errorf("fetch %s, the %s for linux/%s, from release %s: %w", a.archive, a.binary, arch, rel.Tag, err)
		}
		if err := selfupgrade.Verify(data, a.archive, sums); err != nil {
			return nil, fmt.Errorf("refuse to build the agent images from release %s: %s does not match its signed checksums.txt (%w); the download may be cut short or altered: build again, and report it to the repository's maintainers if it persists", rel.Tag, a.archive, err)
		}
		bin, err := untar(data, a.binary)
		if err != nil {
			return nil, fmt.Errorf("take %s out of %s from release %s: %w; report the release to the repository's maintainers", a.binary, a.archive, rel.Tag, err)
		}
		bins[a.binary] = bin
	}
	return bins, nil
}

// fromTree builds bedrock's tree target, then inner's harness image over it,
// from src.Tree's own Dockerfiles with the tree's root as the context, as
// scripts/agent-image.sh builds them.
func (b Builder) fromTree(ctx context.Context, inner string, src BuildSource, version string, tags []string) error {
	out := b.out()
	bedrockPath, bedrockFile, err := treeDockerfile(src.Tree, images.Bedrock)
	if err != nil {
		return err
	}
	innerPath, innerFile, err := treeDockerfile(src.Tree, inner)
	if err != nil {
		return err
	}
	bedrock := BedrockTag(version)
	at := ""
	if src.Commit != "" {
		at = " at commit " + src.Commit
	}
	fmt.Fprintf(out, "agent image: building %s from the development tree %s%s, flai and pilot compiled from its source\n", bedrock, src.Tree, at)
	args := []string{"build", "-f", bedrockPath, "--target", targetTree, "--build-arg", "FLAI_VERSION=" + version}
	// FLAI_COMMIT and FLAI_DATE are passed but kept out of the recipe: a
	// development image is not judged by its commit (ADR-0295 §3), which has
	// a label of its own, so a commit to the tree rebuilds nothing.
	for _, arg := range []struct{ name, value string }{{"FLAI_COMMIT", src.Commit}, {"FLAI_DATE", src.Date}} {
		if arg.value != "" {
			args = append(args, "--build-arg", arg.name+"="+arg.value)
		}
	}
	args = append(args, "--label", LabelRecipe+"="+Recipe(bedrockFile, targetTree, nil), "-t", bedrock, ".")
	if err := b.Run.Run(ctx, Step{Name: "build " + bedrock, Dir: src.Tree, Args: args}); err != nil {
		return err
	}
	return b.harness(ctx, out, harnessBuild{
		inner: inner, bedrock: bedrock, tags: tags, dir: src.Tree, context: ".",
		file: innerPath, data: innerFile, from: "the tree's " + images.Path(inner),
	})
}

// harnessBuild is one harness image's build: the inner harness, the bedrock
// tag it builds over, its tags, the folder the build runs in and its
// context there, its Dockerfile's path and bytes, and where that Dockerfile
// is from, for the line that says what the build does.
type harnessBuild struct {
	inner, bedrock string
	tags           []string
	dir, context   string
	file           string
	data           []byte
	from           string
}

// harness builds h's harness image over its bedrock image, labelled with
// its recipe and the bedrock image's ID.
func (b Builder) harness(ctx context.Context, out io.Writer, h harnessBuild) error {
	info, present, err := b.inspector()(h.bedrock)
	if err != nil {
		return fmt.Errorf("read the ID of %s, which the harness image of %s builds over: %w", h.bedrock, h.inner, err)
	}
	if !present {
		return fmt.Errorf("%s was built, but docker says it is not on this host, so the harness image of %s cannot build over it: check that the build and docker image inspect use the same engine (docker context show), then build again", h.bedrock, h.inner)
	}
	fmt.Fprintf(out, "agent image: building %s over %s (%s) from %s\n", strings.Join(h.tags, ", "), h.bedrock, info.ID, h.from)
	args := []string{"build", "-f", h.file,
		"--build-arg", "BASE=" + h.bedrock,
		"--label", LabelRecipe + "=" + Recipe(h.data, "", nil),
		"--label", LabelBase + "=" + info.ID}
	for _, tag := range h.tags {
		args = append(args, "-t", tag)
	}
	args = append(args, h.context)
	return b.Run.Run(ctx, Step{Name: "build " + h.tags[0], Dir: h.dir, Args: args})
}

// treeDockerfile is the path and bytes of image's Dockerfile in the tree at root.
func treeDockerfile(root, image string) (string, []byte, error) {
	path := filepath.Join(root, filepath.FromSlash(images.Path(image)))
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("read %s, the Dockerfile of the agent image %s in the development tree %s: %w; build from the root of a tree that holds %s", path, image, root, err, images.Path(image))
	}
	return path, data, nil
}

// writeContext writes each of files, by its path under dir, into dir: the
// binaries executable, the Dockerfiles not.
func writeContext(dir string, files map[string][]byte) error {
	for name, data := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		mode := os.FileMode(0o644)
		if !strings.HasSuffix(name, ".Dockerfile") {
			mode = 0o755
		}
		if err := os.WriteFile(path, data, mode); err != nil {
			return err
		}
	}
	return nil
}

// untar is the regular file named name, in any folder, in the tar.gz data.
func untar(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("the archive holds no %s", name)
		}
		if err != nil {
			return nil, err
		}
		if h.Typeflag == tar.TypeReg && filepath.Base(h.Name) == name {
			return io.ReadAll(tr)
		}
	}
}

// refusal is a build refused for a reason its message says in full, which
// errors.As still finds under it.
type refusal struct {
	msg string
	err error
}

// Error is the refusal's message.
func (e *refusal) Error() string { return e.msg }

// Unwrap is the reason the build was refused.
func (e *refusal) Unwrap() error { return e.err }

// out is where b says what each phase does: its Runner's Out, discarded when nil.
func (b Builder) out() io.Writer {
	if b.Run.Out == nil {
		return io.Discard
	}
	return b.Run.Out
}

// program is the docker client b's Runner runs.
func (b Builder) program() string {
	if b.Run.Program == "" {
		return "docker"
	}
	return b.Run.Program
}

// arch is b's Arch, or EngineArch through execx.System when it has none.
func (b Builder) arch() func() (string, error) {
	if b.Arch != nil {
		return b.Arch
	}
	return EngineArch(execx.System{}, b.program())
}

// inspector is b's Inspect, or Docker through execx.System when it has none.
func (b Builder) inspector() Inspector {
	if b.Inspect != nil {
		return b.Inspect
	}
	return Docker(execx.System{}, b.program())
}
