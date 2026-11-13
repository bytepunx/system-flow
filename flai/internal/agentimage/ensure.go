package agentimage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/images"
	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/execx"
)

// Mode is what Ensure does with the images: make the harness images, pulling
// or building them, then the project's layer over them, or build the layer alone.
type Mode string

// The modes of Ensure.
const (
	// ModeEnsure makes each missing or stale harness image: it pulls it, or
	// builds it when the pull cannot run; then it builds each layer that is
	// missing or stale or whose harness image it made.
	ModeEnsure Mode = "ensure"
	// ModeBuild builds each harness image, current or not, without the pull,
	// then each layer over it.
	ModeBuild Mode = "build"
	// ModeLayer builds each layer, current or not, over the harness image
	// already on this host, and leaves every harness image as it is. With
	// the operator's build command named, it runs that command, asking it
	// for the layer alone, since the operator's machinery owns the layer as
	// well as the harness images (ADR-0295 §7).
	ModeLayer Mode = "layer"
)

// Way is how Ensure left an image.
type Way string

// The ways Ensure leaves an image.
const (
	// WayPulled is a harness image pulled by the digest its release's signed list names.
	WayPulled Way = "pulled"
	// WayBuilt is a harness image built on this host; Outcome.Built.From says from what.
	WayBuilt Way = "built"
	// WayLayered is the project's layer built over its harness image; Outcome.Layered says what it made.
	WayLayered Way = "layered"
	// WayCurrent is an image already current, left as it is.
	WayCurrent Way = "current"
	// WayLeft is a harness image ModeLayer leaves as it is, current or not.
	WayLeft Way = "left"
	// WayOwn is the operator's own image, which flai never makes.
	WayOwn Way = "own"
	// WayCommand is an image the operator's build command made, in place of
	// flai's pull and build (ADR-0295 §7).
	WayCommand Way = "command"
)

// Ensurer makes the harness images through Pull and Build (Puller.Pull and
// Builder.Build, or fakes), and the project's layer over them through Layer
// (Builder.Layer, or a fake) in the project at Root, the layer's build
// context, saying what it does on Out (discarded when nil).
//
// With Command, the operator's build command from the host's configuration
// (ADR-0295 §7), nil for none, it runs that command through Run in Root in
// place of Pull, Build, and Layer, and judges the images again through
// Inspect (Docker through execx.System with Run's Program when nil).
type Ensurer struct {
	Pull    func(ctx context.Context, image, version string, tags []string) (Pulled, error)
	Build   func(ctx context.Context, inner string, src BuildSource, tags []string) (Built, error)
	Layer   func(ctx context.Context, root string, img Image, base string) (Layered, error)
	Root    string
	Out     io.Writer
	Command []string
	Run     Runner
	Inspect Inspector
}

// Outcome is what Ensure did with one image: the image, the way it took,
// what the pull made when it pulled, what the build made when it built, what
// the layer's build made when it built the layer, and why the pull failed
// when the build stood in for it ("" when it did not).
type Outcome struct {
	Image    Image
	Way      Way
	Pulled   Pulled
	Built    Built
	Layered  Layered
	FellBack string
}

// Ensure makes the images of statuses as mode asks, for flai src.Version,
// from the development tree src.Tree when it names one, which is never
// pulled for (ADR-0295 §1, §3): each harness image first, then each layer
// over the harness image of its inner harness among statuses (§4, §7). It
// leaves a current harness image in ModeEnsure and the operator's own image
// always, and builds a layer when it is missing or stale or its harness
// image was pulled or built in this run, and in ModeBuild and ModeLayer
// always. A pull refused because its list or image does not verify stops it
// with no build in its place; any failure stops it before the layers, and a
// layer whose harness image is not on this host stops it, naming that image.
// With e.Command it runs the command in their place, as byCommand says, for
// a development flai too, whose images come from the operator's command (§3).
func (e Ensurer) Ensure(ctx context.Context, mode Mode, statuses []Status, src BuildSource) ([]Outcome, error) {
	if mode != ModeEnsure && mode != ModeBuild && mode != ModeLayer {
		return nil, fmt.Errorf("make the agent images in the mode %q: there is no such mode; name %s, %s, or %s", mode, ModeEnsure, ModeBuild, ModeLayer)
	}
	out := e.Out
	if out == nil {
		out = io.Discard
	}
	if len(e.Command) > 0 {
		outcomes, err := e.byCommand(ctx, mode, statuses, src, out)
		if err != nil {
			return outcomes, err
		}
		summarize(out, outcomes, src, e.Command)
		return outcomes, nil
	}
	if src.Tree == "" && buildinfo.Development(src.Version) {
		refs, files := []string{}, []string{images.Path(images.Bedrock)}
		for _, s := range statuses {
			if makes(s) && acts(mode, s) {
				refs = append(refs, s.Ref)
				files = append(files, images.Path(s.Inner))
			}
		}
		if len(refs) > 0 {
			return nil, fmt.Errorf("make the agent images %s for flai %s: a development flai has no release to pull its images from or take its Linux flai and pilot from; run it in a tree that holds %s, or run a released flai", strings.Join(refs, ", "), shown(src.Version), strings.Join(files, " and "))
		}
	}
	outcomes := make([]Outcome, 0, len(statuses))
	bases := map[string]baseRun{}
	for _, s := range statuses {
		if s.Kind == KindLayer {
			continue
		}
		o, err := e.one(ctx, mode, s, src, out)
		if err != nil {
			return outcomes, err
		}
		outcomes = append(outcomes, o)
		bases[s.Inner] = baseRun{status: s, way: o.Way}
	}
	for _, s := range statuses {
		if s.Kind != KindLayer {
			continue
		}
		b, ok := bases[s.Inner]
		if !ok {
			return outcomes, fmt.Errorf("build the layer image %s: no harness image of %s is among the agent images, so there is none to build it over; name the images Needed gives, which hold the harness image under each layer", s.Ref, s.Inner)
		}
		o, err := e.layer(ctx, mode, s, b, out)
		if err != nil {
			return outcomes, err
		}
		outcomes = append(outcomes, o)
	}
	summarize(out, outcomes, src, nil)
	return outcomes, nil
}

// byCommand makes the images of statuses as mode asks with the operator's
// build command e.Command, in place of the pull, the build, and the layer's
// build (ADR-0295 §7): once for each inner harness among statuses, in their
// order, that has an image to make, which in ModeEnsure is one missing or
// stale, in ModeBuild any, and in ModeLayer its layer alone, with
// FLAI_IMAGE_* naming them (commandImages). The operator's own image (--image)
// is never asked made. After each run it judges that inner harness's images
// again with Check, and fails, naming the command, each image, and why, when
// one it makes is still missing or stale; a non-zero exit or a stall fails it
// as RunCommand says.
func (e Ensurer) byCommand(ctx context.Context, mode Mode, statuses []Status, src BuildSource, out io.Writer) ([]Outcome, error) {
	run := e.Run
	if run.Out == nil {
		run.Out = out
	}
	var outcomes []Outcome
	for _, inner := range innersOf(statuses) {
		var group []Status
		var h, l *Status
		for _, s := range statuses {
			if s.Inner == inner {
				group = append(group, s)
			}
		}
		for i := range group {
			if group[i].Kind == KindHarness {
				h = &group[i]
			} else {
				l = &group[i]
			}
		}
		if h == nil {
			return outcomes, fmt.Errorf("make the layer image %s with the operator's build command: no harness image of %s is among the agent images, so there is none to build it over; name the images Needed gives, which hold the harness image under each layer", l.Ref, inner)
		}
		c, err := commandImages(e.Root, mode, *h, l, src)
		if err != nil {
			return outcomes, err
		}
		made := make([]Outcome, 0, len(group))
		for _, s := range group {
			made = append(made, e.commandOutcome(mode, s, c, src, out))
		}
		if len(c.Stale) > 0 {
			if err := RunCommand(ctx, run, e.Root, e.Command, c); err != nil {
				return outcomes, err
			}
			if err := e.judgeCommand(mode, group, src); err != nil {
				return outcomes, err
			}
		}
		outcomes = append(outcomes, made...)
	}
	return outcomes, nil
}

// commandOutcome is the way byCommand leaves s, given c, what the command
// is asked to make, saying on out what it does before the command runs.
func (e Ensurer) commandOutcome(mode Mode, s Status, c CommandImages, src BuildSource, out io.Writer) Outcome {
	o := Outcome{Image: s.Image}
	part := PartHarness
	if s.Kind == KindLayer {
		part = PartLayer
	}
	switch {
	case s.Source == SourceHost:
		fmt.Fprintf(out, "agent image: %s is the operator's own (flai serve agent harness container --image %s=%s), which flai never pulls, builds, or judges; left as it is\n", s.Ref, s.Inner, s.Ref)
		o.Way = WayOwn
	case slices.Contains(c.Stale, part):
		fmt.Fprintf(out, "agent image: %s %s; the operator's build command %s makes it, in place of flai's pull and build\n", s.Ref, stateWords(s), CommandWords(e.Command))
		o.Way = WayCommand
	case mode == ModeLayer:
		fmt.Fprintf(out, "agent image: %s %s; left as it is, as only the project's layer is built\n", s.Ref, stateWords(s))
		o.Way = WayLeft
	default:
		fmt.Fprintf(out, "agent image: %s is current for flai %s; nothing to do\n", s.Ref, shown(src.Version))
		o.Way = WayCurrent
	}
	return o
}

// judgeCommand judges group, one inner harness's images, again after the
// operator's build command ran, and fails naming the command and each image
// mode makes (all but the operator's own, and in ModeLayer the layer alone)
// that is still missing or stale, and why.
func (e Ensurer) judgeCommand(mode Mode, group []Status, src BuildSource) error {
	imgs := make([]Image, 0, len(group))
	for _, s := range group {
		imgs = append(imgs, s.Image)
	}
	after, err := Check(e.Root, imgs, e.inspector(), func() (string, error) { return src.Version, nil })
	if err != nil {
		return fmt.Errorf("judge the agent images again after the operator's build command %s: %w", CommandWords(e.Command), err)
	}
	var left []string
	for _, s := range after {
		if s.Source == SourceHost || (mode == ModeLayer && s.Kind != KindLayer) {
			continue
		}
		switch s.State {
		case StateMissing:
			left = append(left, s.Ref+" is missing")
		case StateStale:
			left = append(left, s.Ref+" "+stateWords(s))
		}
	}
	if len(left) == 0 {
		return nil
	}
	return fmt.Errorf("the operator's build command %s exited 0 but left the agent images wrong: %s; make it build each with the tag and the labels its FLAI_IMAGE_* variables name, then run it again", CommandWords(e.Command), strings.Join(left, "; and "))
}

// innersOf are the inner harnesses of statuses, each once, in their order.
func innersOf(statuses []Status) []string {
	var inners []string
	for _, s := range statuses {
		if !slices.Contains(inners, s.Inner) {
			inners = append(inners, s.Inner)
		}
	}
	return inners
}

// inspector is e's Inspect, or Docker through execx.System with its
// Runner's Program (docker when empty) when it has none.
func (e Ensurer) inspector() Inspector {
	if e.Inspect != nil {
		return e.Inspect
	}
	program := e.Run.Program
	if program == "" {
		program = "docker"
	}
	return Docker(execx.System{}, program)
}

// baseRun is the harness image a layer builds over: its status before Ensure
// ran and the way Ensure left it.
type baseRun struct {
	status Status
	way    Way
}

// renewed reports whether Ensure pulled or built b in this run.
func (b baseRun) renewed() bool {
	return b.way == WayPulled || b.way == WayBuilt
}

// makes reports whether s is an image flai makes: a harness image of flai's default.
func makes(s Status) bool {
	return s.Kind == KindHarness && s.Source == SourceDefault
}

// acts reports whether mode makes the harness image s: any in ModeBuild,
// one missing or stale in ModeEnsure, none in ModeLayer.
func acts(mode Mode, s Status) bool {
	switch mode {
	case ModeBuild:
		return true
	case ModeEnsure:
		return s.State != StateCurrent
	}
	return false
}

// one makes the harness image s as mode asks, saying on out what it does before it does it.
func (e Ensurer) one(ctx context.Context, mode Mode, s Status, src BuildSource, out io.Writer) (Outcome, error) {
	o := Outcome{Image: s.Image}
	switch {
	case s.Source == SourceHost:
		fmt.Fprintf(out, "agent image: %s is the operator's own (flai serve agent harness container --image %s=%s), which flai never pulls, builds, or judges; left as it is\n", s.Ref, s.Inner, s.Ref)
		o.Way = WayOwn
		return o, nil
	case mode == ModeLayer:
		fmt.Fprintf(out, "agent image: %s %s; left as it is, as only the project's layer is built\n", s.Ref, stateWords(s))
		o.Way = WayLeft
		return o, nil
	case !acts(mode, s):
		fmt.Fprintf(out, "agent image: %s is current for flai %s; nothing to do\n", s.Ref, shown(src.Version))
		o.Way = WayCurrent
		return o, nil
	}
	tags := HarnessTags(s.Ref, s.Inner, src.Version)
	state := stateWords(s)
	switch {
	case src.Tree != "":
		fmt.Fprintf(out, "agent image: %s %s; building it from the development tree %s, as a development flai is never pulled for\n", s.Ref, state, src.Tree)
		return e.build(ctx, o, s.Inner, src, tags)
	case mode == ModeBuild:
		fmt.Fprintf(out, "agent image: %s %s; building it on this host, without the pull, from the Dockerfiles flai carries and flai %s's release binaries\n", s.Ref, state, src.Version)
		return e.build(ctx, o, s.Inner, src, tags)
	}
	fmt.Fprintf(out, "agent image: %s %s; pulling %s of flai %s by the digest its release's signed list names\n", s.Ref, state, Published(s.Inner), src.Version)
	pulled, err := e.Pull(ctx, s.Inner, src.Version, tags)
	if err == nil {
		o.Way, o.Pulled = WayPulled, pulled
		return o, nil
	}
	if !FallBack(err) {
		return o, err
	}
	fmt.Fprintf(out, "agent image: the pull of %s failed (%v); building it on this host instead, from the Dockerfiles flai carries and flai %s's release binaries\n", s.Ref, err, src.Version)
	o.FellBack = err.Error()
	o, berr := e.build(ctx, o, s.Inner, src, tags)
	if berr != nil {
		return o, fmt.Errorf("make the agent image %s: its pull failed (%w), and the build on this host that stands in for it failed too: %w", s.Ref, err, berr)
	}
	return o, nil
}

// layer builds the layer image s as mode asks over b, the harness image of
// its inner harness, saying on out what it does before it does it: it
// refuses when b is not on this host, and leaves a current layer over a
// harness image this run did not make in ModeEnsure.
func (e Ensurer) layer(ctx context.Context, mode Mode, s Status, b baseRun, out io.Writer) (Outcome, error) {
	o := Outcome{Image: s.Image}
	h := b.status
	if !h.Present && !b.renewed() {
		if h.Source == SourceHost {
			return o, fmt.Errorf("build the layer image %s: the harness image %s it builds over, the operator's own (flai serve agent harness container --image %s=%s), is not on this host, and flai never makes it; pull or build %s on this host, or name another image, then run it again", s.Ref, h.Ref, h.Inner, h.Ref, h.Ref)
		}
		return o, fmt.Errorf("build the layer image %s: the harness image %s it builds over is not on this host; run flai image ensure first, then build the layer again", s.Ref, h.Ref)
	}
	switch {
	case mode == ModeEnsure && s.State == StateCurrent && !b.renewed():
		fmt.Fprintf(out, "agent image: %s, the project's layer from %s, is current over %s; nothing to do\n", s.Ref, s.Dockerfile, h.Ref)
		o.Way = WayCurrent
		return o, nil
	case b.renewed():
		fmt.Fprintf(out, "agent image: %s, the project's layer from %s, %s; building it over %s, which was just %s\n", s.Ref, s.Dockerfile, stateWords(s), h.Ref, b.way)
	default:
		fmt.Fprintf(out, "agent image: %s, the project's layer from %s, %s; building it over %s\n", s.Ref, s.Dockerfile, stateWords(s), h.Ref)
	}
	layered, err := e.Layer(ctx, e.Root, s.Image, h.Ref)
	if err != nil {
		return o, err
	}
	o.Way, o.Layered = WayLayered, layered
	return o, nil
}

// build builds inner's harness image from src with tags into o.
func (e Ensurer) build(ctx context.Context, o Outcome, inner string, src BuildSource, tags []string) (Outcome, error) {
	built, err := e.Build(ctx, inner, src, tags)
	if err != nil {
		return o, err
	}
	o.Way, o.Built = WayBuilt, built
	return o, nil
}

// stateWords is what s's state is, in the words a line says it: "is
// missing", "is stale (<why>)", or "is current".
func stateWords(s Status) string {
	if s.State != StateStale {
		return "is " + string(s.State)
	}
	reasons := make([]string, 0, len(s.Reasons))
	for _, r := range s.Reasons {
		reasons = append(reasons, r.Message)
	}
	return "is stale (" + strings.Join(reasons, "; ") + ")"
}

// summarize says on out which way each of outcomes took, src being what a
// build was made from and command the operator's build command, nil for none.
func summarize(out io.Writer, outcomes []Outcome, src BuildSource, command []string) {
	for _, o := range outcomes {
		var said string
		switch o.Way {
		case WayPulled:
			said = "pulled by digest as " + o.Pulled.Reference + ", tagged " + strings.Join(o.Pulled.Tags, ", ")
			if o.Pulled.Cached {
				said += ", by the verified list cached on this host, as the release could not be reached"
			}
		case WayBuilt:
			if o.Built.From == FromTree {
				said = "built from the development tree " + src.Tree
			} else {
				said = "built on this host from flai " + o.Built.Version + "'s release binaries and the Dockerfiles flai carries"
			}
			if o.FellBack != "" {
				said += ", as its pull failed"
			}
		case WayLayered:
			said = "the project's layer, built over " + o.Layered.Base + " (" + o.Layered.BaseID + ") from " + o.Image.Dockerfile + ", labelled with its recipe " + o.Layered.Recipe
		case WayCurrent:
			said = "current, left as it is"
		case WayLeft:
			said = "left as it is, as only the project's layer was built"
		case WayOwn:
			said = "the operator's own, left as it is"
		case WayCommand:
			said = "made by the operator's build command " + CommandWords(command)
		}
		fmt.Fprintf(out, "agent image: %s: %s\n", o.Image.Ref, said)
	}
}

// HarnessTags are the tags a harness image of inner gets for flai version:
// ref, the name the start runs it by, when it is not "", and
// flai-agent-<inner>:<version>, without a leading v, when that is another.
func HarnessTags(ref, inner, version string) []string {
	var tags []string
	if ref != "" {
		tags = append(tags, ref)
	}
	if own := "flai-agent-" + inner + ":" + strings.TrimPrefix(version, "v"); own != ref {
		tags = append(tags, own)
	}
	return tags
}

// LogPath is the log file a run of the agent images started at now writes
// under root's .flai-cache: .flai-cache/agent-images/<UTC stamp>.log, the
// stamp as 20060102T150405Z.
func LogPath(root string, now time.Time) string {
	return filepath.Join(root, ".flai-cache", "agent-images", now.UTC().Format("20060102T150405Z")+".log")
}
