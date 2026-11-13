package agentimage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/images"
	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

// Mode is what Ensure does with each harness image: pull it first, or build it without the pull.
type Mode string

// The modes of Ensure.
const (
	// ModeEnsure makes each missing or stale harness image: it pulls it, or builds it when the pull cannot run.
	ModeEnsure Mode = "ensure"
	// ModeBuild builds each harness image, current or not, without the pull.
	ModeBuild Mode = "build"
)

// Way is how Ensure left an image.
type Way string

// The ways Ensure leaves an image.
const (
	// WayPulled is a harness image pulled by the digest its release's signed list names.
	WayPulled Way = "pulled"
	// WayBuilt is a harness image built on this host; Outcome.Built.From says from what.
	WayBuilt Way = "built"
	// WayCurrent is a harness image already current, left as it is.
	WayCurrent Way = "current"
	// WayOwn is the operator's own image, which flai never makes.
	WayOwn Way = "own"
	// WayNotBuilt is the project's layer, which flai does not build yet.
	WayNotBuilt Way = "not built"
)

// Ensurer makes the harness images through Pull and Build (Puller.Pull and
// Builder.Build, or fakes), saying what it does on Out (discarded when nil).
type Ensurer struct {
	Pull  func(ctx context.Context, image, version string, tags []string) (Pulled, error)
	Build func(ctx context.Context, inner string, src BuildSource, tags []string) (Built, error)
	Out   io.Writer
}

// Outcome is what Ensure did with one image: the image, the way it took,
// what the pull made when it pulled, what the build made when it built, and
// why the pull failed when the build stood in for it ("" when it did not).
type Outcome struct {
	Image    Image
	Way      Way
	Pulled   Pulled
	Built    Built
	FellBack string
}

// Ensure makes each harness image of statuses as mode asks, for flai
// src.Version, from the development tree src.Tree when it names one, which
// is never pulled for (ADR-0295 §1, §3). It leaves a current image in
// ModeEnsure, the operator's own image always, and the project's layer,
// which it does not build. A pull refused because its list or image does
// not verify stops it with no build in its place.
func (e Ensurer) Ensure(ctx context.Context, mode Mode, statuses []Status, src BuildSource) ([]Outcome, error) {
	if mode != ModeEnsure && mode != ModeBuild {
		return nil, fmt.Errorf("make the agent images in the mode %q: there is no such mode; name %s or %s", mode, ModeEnsure, ModeBuild)
	}
	out := e.Out
	if out == nil {
		out = io.Discard
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
	for _, s := range statuses {
		o, err := e.one(ctx, mode, s, src, out)
		if err != nil {
			return outcomes, err
		}
		outcomes = append(outcomes, o)
	}
	summarize(out, outcomes, src)
	return outcomes, nil
}

// makes reports whether s is an image flai makes: a harness image of flai's default.
func makes(s Status) bool {
	return s.Kind == KindHarness && s.Source == SourceDefault
}

// acts reports whether mode makes the image s: any in ModeBuild, else one missing or stale.
func acts(mode Mode, s Status) bool {
	return mode == ModeBuild || s.State != StateCurrent
}

// one makes the image s as mode asks, saying on out what it does before it does it.
func (e Ensurer) one(ctx context.Context, mode Mode, s Status, src BuildSource, out io.Writer) (Outcome, error) {
	o := Outcome{Image: s.Image}
	switch {
	case s.Kind == KindLayer:
		fmt.Fprintf(out, "agent image: %s, the project's layer from %s, is %s; flai does not build a project's layer yet, so the project's own build makes it\n", s.Ref, s.Dockerfile, s.State)
		o.Way = WayNotBuilt
		return o, nil
	case s.Source == SourceHost:
		fmt.Fprintf(out, "agent image: %s is the operator's own (flai serve agent harness container --image %s=%s), which flai never pulls, builds, or judges; left as it is\n", s.Ref, s.Inner, s.Ref)
		o.Way = WayOwn
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

// summarize says on out which way each of outcomes took, src being what a build was made from.
func summarize(out io.Writer, outcomes []Outcome, src BuildSource) {
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
		case WayCurrent:
			said = "current, left as it is"
		case WayOwn:
			said = "the operator's own, left as it is"
		case WayNotBuilt:
			said = "the project's layer, not built by flai"
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
