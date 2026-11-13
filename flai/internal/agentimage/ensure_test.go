package agentimage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/dashboard"
)

// ensureFake records the pulls, builds, and layer builds an Ensurer asks
// of it, and each in the order asked as calls, answering each pull with
// pullErr, each build with buildErr, and each layer build with layerErr.
type ensureFake struct {
	pulls    [][]string // image, version, then the tags
	builds   []BuildSource
	tags     [][]string
	layers   [][]string // root, the layer's ref, then the base
	calls    []string
	pullErr  error
	buildErr error
	layerErr error
}

func (f *ensureFake) pull(_ context.Context, image, version string, tags []string) (Pulled, error) {
	f.pulls = append(f.pulls, append([]string{image, version}, tags...))
	f.calls = append(f.calls, "pull "+image)
	if f.pullErr != nil {
		return Pulled{}, f.pullErr
	}
	return Pulled{Image: image, Version: version, Reference: Published(image) + "@" + indexD, Tags: tags}, nil
}

func (f *ensureFake) build(_ context.Context, inner string, src BuildSource, tags []string) (Built, error) {
	f.builds = append(f.builds, src)
	f.tags = append(f.tags, tags)
	f.calls = append(f.calls, "build "+inner)
	if f.buildErr != nil {
		return Built{}, f.buildErr
	}
	from := FromRelease
	if src.Tree != "" {
		from = FromTree
	}
	return Built{Inner: inner, Version: src.Version, Tags: tags, Bedrock: BedrockTag(src.Version), From: from}, nil
}

func (f *ensureFake) layer(_ context.Context, root string, img Image, base string) (Layered, error) {
	f.layers = append(f.layers, []string{root, img.Ref, base})
	f.calls = append(f.calls, "layer "+img.Ref)
	if f.layerErr != nil {
		return Layered{}, f.layerErr
	}
	return Layered{Image: img, Tags: []string{img.Ref}, Base: base, BaseID: "sha256:base", Recipe: "sha256:recipe"}, nil
}

// ensurer is an Ensurer over f in the project ensureRoot, writing to out.
func (f *ensureFake) ensurer(out *bytes.Buffer) Ensurer {
	return Ensurer{Pull: f.pull, Build: f.build, Layer: f.layer, Root: ensureRoot, Out: out}
}

const (
	ensureRoot     = "/project"
	ensureLayerRef = "flai-agent-t-claude-code:local"
	ensureLayerDF  = "images/agent/Dockerfile"
)

// layerStatus is the project's claude-code layer image in state.
func layerStatus(state State) Status {
	s := Status{Image: Image{Kind: KindLayer, Inner: "claude-code", Ref: ensureLayerRef, Source: SourceLayer, Dockerfile: ensureLayerDF, Args: map[string]string{"GO": "1.25"}}, State: state, Present: state != StateMissing}
	if state == StateStale {
		s.Reasons = []Reason{{Kind: ReasonRecipe, Message: "built from another Dockerfile"}}
	}
	return s
}

const ensureRef = "flai-agent-claude-code:local"

// harnessStatus is flai's default claude-code harness image in state.
func harnessStatus(state State) Status {
	s := Status{Image: Image{Kind: KindHarness, Inner: "claude-code", Ref: ensureRef, Source: SourceDefault}, State: state, Present: state != StateMissing}
	if state == StateStale {
		s.Reasons = []Reason{{Kind: ReasonFlai, Message: "built for flai 1.0.0, not " + pullV}}
	}
	return s
}

func TestEnsurePullsAMissingOrStaleImageWithItsTags(t *testing.T) {
	t.Parallel()
	for _, state := range []State{StateMissing, StateStale} {
		f := &ensureFake{}
		var out bytes.Buffer
		got, err := f.ensurer(&out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(state)}, BuildSource{Version: pullV})
		if err != nil {
			t.Fatalf("%s: %v", state, err)
		}
		want := [][]string{{"claude-code", pullV, ensureRef, "flai-agent-claude-code:" + pullV}}
		if !reflect.DeepEqual(f.pulls, want) || len(f.builds) != 0 {
			t.Errorf("%s: pulls %v, builds %v", state, f.pulls, f.builds)
		}
		if len(got) != 1 || got[0].Way != WayPulled || got[0].FellBack != "" || got[0].Pulled.Reference != claudeRepo+"@"+indexD {
			t.Errorf("%s: outcomes %+v", state, got)
		}
		text := out.String()
		if !strings.Contains(text, "agent image: "+ensureRef+" is "+string(state)) || !strings.Contains(text, "agent image: "+ensureRef+": pulled by digest as "+claudeRepo+"@"+indexD) {
			t.Errorf("%s: output\n%s", state, text)
		}
	}
}

func TestEnsureLeavesACurrentImage(t *testing.T) {
	t.Parallel()
	f := &ensureFake{}
	var out bytes.Buffer
	got, err := f.ensurer(&out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateCurrent)}, BuildSource{Version: pullV})
	if err != nil || len(got) != 1 || got[0].Way != WayCurrent || len(f.pulls)+len(f.builds) != 0 {
		t.Fatalf("outcomes %+v, %v; pulls %v, builds %v", got, err, f.pulls, f.builds)
	}
	if !strings.Contains(out.String(), "is current for flai "+pullV+"; nothing to do") {
		t.Errorf("output\n%s", out.String())
	}
}

func TestEnsureBuildsFromTheReleaseWhenThePullFallsBack(t *testing.T) {
	t.Parallel()
	failed := &pullError{err: &StepError{Step: "pull x", Tries: 3, Said: "connection reset by peer", Network: true, Err: errors.New("exit status 1")}}
	f := &ensureFake{pullErr: failed}
	var out bytes.Buffer
	got, err := f.ensurer(&out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateMissing)}, BuildSource{Version: pullV})
	if err != nil {
		t.Fatal(err)
	}
	if len(f.pulls) != 1 || len(f.builds) != 1 || f.builds[0] != (BuildSource{Version: pullV}) {
		t.Fatalf("pulls %v, builds %v", f.pulls, f.builds)
	}
	if !reflect.DeepEqual(f.tags[0], []string{ensureRef, "flai-agent-claude-code:" + pullV}) {
		t.Errorf("build tags %v", f.tags[0])
	}
	if len(got) != 1 || got[0].Way != WayBuilt || got[0].Built.From != FromRelease || got[0].FellBack != failed.Error() {
		t.Errorf("outcomes %+v", got)
	}
	text := out.String()
	for _, want := range []string{
		"the pull of " + ensureRef + " failed (" + failed.Error() + "); building it on this host instead, from the Dockerfiles flai carries",
		ensureRef + ": built on this host from flai " + pullV + "'s release binaries and the Dockerfiles flai carries, as its pull failed",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("output lacks %q:\n%s", want, text)
		}
	}
}

func TestEnsureStopsOnARefusalWithNoBuild(t *testing.T) {
	t.Parallel()
	refused := &dashboard.Refusal{Version: pullV, Kind: dashboard.RefusalBad, Product: product("claude-code"), Reason: "the signature does not verify"}
	f := &ensureFake{pullErr: refused}
	_, err := f.ensurer(&bytes.Buffer{}).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateMissing)}, BuildSource{Version: pullV})
	var got *dashboard.Refusal
	if !errors.As(err, &got) || got.Kind != dashboard.RefusalBad {
		t.Fatalf("error %v, want the refusal", err)
	}
	if len(f.builds) != 0 {
		t.Errorf("built %v after a refusal", f.builds)
	}
}

func TestEnsureSaysWhenTheFallbackBuildFails(t *testing.T) {
	t.Parallel()
	f := &ensureFake{pullErr: &pullError{err: errors.New("pull failed")}, buildErr: errors.New("build failed")}
	_, err := f.ensurer(&bytes.Buffer{}).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateMissing)}, BuildSource{Version: pullV})
	if err == nil || !strings.Contains(err.Error(), "its pull failed (pull failed)") || !strings.Contains(err.Error(), "build failed") {
		t.Fatalf("error %v", err)
	}
}

func TestEnsureBuildsFromADevelopmentTreeAndNeverPulls(t *testing.T) {
	t.Parallel()
	src := BuildSource{Version: "dev", Tree: "/tree", Commit: "abc123"}
	for _, mode := range []Mode{ModeEnsure, ModeBuild} {
		f := &ensureFake{}
		var out bytes.Buffer
		got, err := f.ensurer(&out).Ensure(context.Background(), mode, []Status{harnessStatus(StateStale)}, src)
		if err != nil {
			t.Fatalf("%s: %v", mode, err)
		}
		if len(f.pulls) != 0 || len(f.builds) != 1 || f.builds[0] != src {
			t.Errorf("%s: pulls %v, builds %v", mode, f.pulls, f.builds)
		}
		if !reflect.DeepEqual(f.tags[0], []string{ensureRef, "flai-agent-claude-code:dev"}) {
			t.Errorf("%s: build tags %v", mode, f.tags[0])
		}
		if len(got) != 1 || got[0].Way != WayBuilt || got[0].Built.From != FromTree {
			t.Errorf("%s: outcomes %+v", mode, got)
		}
		if !strings.Contains(out.String(), ensureRef+": built from the development tree /tree") {
			t.Errorf("%s: output\n%s", mode, out.String())
		}
	}
}

func TestEnsureBuildModeBuildsACurrentImageWithoutThePull(t *testing.T) {
	t.Parallel()
	f := &ensureFake{}
	var out bytes.Buffer
	got, err := f.ensurer(&out).Ensure(context.Background(), ModeBuild, []Status{harnessStatus(StateCurrent)}, BuildSource{Version: pullV})
	if err != nil {
		t.Fatal(err)
	}
	if len(f.pulls) != 0 || len(f.builds) != 1 || f.builds[0] != (BuildSource{Version: pullV}) {
		t.Fatalf("pulls %v, builds %v", f.pulls, f.builds)
	}
	if len(got) != 1 || got[0].Way != WayBuilt || got[0].Built.From != FromRelease || got[0].FellBack != "" {
		t.Errorf("outcomes %+v", got)
	}
	if !strings.Contains(out.String(), "is current; building it on this host, without the pull") {
		t.Errorf("output\n%s", out.String())
	}
}

func TestEnsureLeavesTheOperatorsOwnImage(t *testing.T) {
	t.Parallel()
	own := Status{Image: Image{Kind: KindHarness, Inner: "claude-code", Ref: "mine:1", Source: SourceHost}, Present: true, State: StateOwn}
	for _, mode := range []Mode{ModeEnsure, ModeBuild, ModeLayer} {
		f := &ensureFake{}
		var out bytes.Buffer
		got, err := f.ensurer(&out).Ensure(context.Background(), mode, []Status{own}, BuildSource{Version: "dev"})
		if err != nil {
			t.Fatalf("%s: %v", mode, err)
		}
		if len(f.calls) != 0 || len(got) != 1 || got[0].Way != WayOwn {
			t.Errorf("%s: outcomes %+v; calls %v", mode, got, f.calls)
		}
		if !strings.Contains(out.String(), "mine:1 is the operator's own") {
			t.Errorf("%s: output\n%s", mode, out.String())
		}
	}
}

// S-0745: a harness image pulled in this run is a new base, so the layer
// over it is built after it, current or not, in the project's root over the
// harness image's ref.
func TestEnsureBuildsTheLayerAfterAPull(t *testing.T) {
	t.Parallel()
	for _, state := range []State{StateMissing, StateStale, StateCurrent} {
		f := &ensureFake{}
		var out bytes.Buffer
		got, err := f.ensurer(&out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateStale), layerStatus(state)}, BuildSource{Version: pullV})
		if err != nil {
			t.Fatalf("%s: %v", state, err)
		}
		if want := []string{"pull claude-code", "layer " + ensureLayerRef}; !reflect.DeepEqual(f.calls, want) {
			t.Errorf("%s: calls %v, want %v", state, f.calls, want)
		}
		if want := [][]string{{ensureRoot, ensureLayerRef, ensureRef}}; !reflect.DeepEqual(f.layers, want) {
			t.Errorf("%s: layers %v, want %v", state, f.layers, want)
		}
		if len(got) != 2 || got[0].Way != WayPulled || got[1].Way != WayLayered || got[1].Layered.Base != ensureRef || got[1].Layered.BaseID != "sha256:base" {
			t.Errorf("%s: outcomes %+v", state, got)
		}
		text := out.String()
		for _, want := range []string{
			ensureLayerRef + ", the project's layer from " + ensureLayerDF + ", is " + string(state),
			"; building it over " + ensureRef + ", which was just pulled\n",
			"agent image: " + ensureLayerRef + ": the project's layer, built over " + ensureRef + " (sha256:base) from " + ensureLayerDF + ", labelled with its recipe sha256:recipe\n",
		} {
			if !strings.Contains(text, want) {
				t.Errorf("%s: output lacks %q:\n%s", state, want, text)
			}
		}
	}
}

// S-0745: flai image build builds the harness image, then the layer over
// it, both current.
func TestEnsureBuildModeBuildsTheLayerAfterTheHarnessImage(t *testing.T) {
	t.Parallel()
	f := &ensureFake{}
	var out bytes.Buffer
	got, err := f.ensurer(&out).Ensure(context.Background(), ModeBuild, []Status{harnessStatus(StateCurrent), layerStatus(StateCurrent)}, BuildSource{Version: pullV})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"build claude-code", "layer " + ensureLayerRef}; !reflect.DeepEqual(f.calls, want) {
		t.Errorf("calls %v, want %v", f.calls, want)
	}
	if len(got) != 2 || got[0].Way != WayBuilt || got[1].Way != WayLayered {
		t.Errorf("outcomes %+v", got)
	}
	if !strings.Contains(out.String(), "building it over "+ensureRef+", which was just built\n") {
		t.Errorf("output\n%s", out.String())
	}
}

// S-0745: a missing or stale layer over a current harness image is built
// over it, and the harness image is left as it is.
func TestEnsureBuildsAMissingOrStaleLayerOverACurrentHarnessImage(t *testing.T) {
	t.Parallel()
	for _, state := range []State{StateMissing, StateStale} {
		f := &ensureFake{}
		var out bytes.Buffer
		got, err := f.ensurer(&out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateCurrent), layerStatus(state)}, BuildSource{Version: pullV})
		if err != nil {
			t.Fatalf("%s: %v", state, err)
		}
		if want := []string{"layer " + ensureLayerRef}; !reflect.DeepEqual(f.calls, want) {
			t.Errorf("%s: calls %v, want %v", state, f.calls, want)
		}
		if len(got) != 2 || got[0].Way != WayCurrent || got[1].Way != WayLayered {
			t.Errorf("%s: outcomes %+v", state, got)
		}
		if state == StateStale && !strings.Contains(out.String(), "is stale (built from another Dockerfile); building it over "+ensureRef+"\n") {
			t.Errorf("%s: output\n%s", state, out.String())
		}
	}
}

// S-0745: a current layer over a current harness image is left as it is.
func TestEnsureLeavesACurrentLayer(t *testing.T) {
	t.Parallel()
	f := &ensureFake{}
	var out bytes.Buffer
	got, err := f.ensurer(&out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateCurrent), layerStatus(StateCurrent)}, BuildSource{Version: pullV})
	if err != nil || len(f.calls) != 0 || len(got) != 2 || got[1].Way != WayCurrent {
		t.Fatalf("outcomes %+v, %v; calls %v", got, err, f.calls)
	}
	if !strings.Contains(out.String(), ensureLayerRef+", the project's layer from "+ensureLayerDF+", is current over "+ensureRef+"; nothing to do") {
		t.Errorf("output\n%s", out.String())
	}
}

// S-0745: the harness images are all made before any layer, whatever
// order the statuses come in.
func TestEnsureMakesTheHarnessImagesBeforeTheLayers(t *testing.T) {
	t.Parallel()
	other := Status{Image: Image{Kind: KindHarness, Inner: "opencode", Ref: "flai-agent-opencode:local", Source: SourceDefault}, State: StateMissing}
	f := &ensureFake{}
	if _, err := f.ensurer(&bytes.Buffer{}).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateMissing), layerStatus(StateMissing), other}, BuildSource{Version: pullV}); err != nil {
		t.Fatal(err)
	}
	if want := []string{"pull claude-code", "pull opencode", "layer " + ensureLayerRef}; !reflect.DeepEqual(f.calls, want) {
		t.Errorf("calls %v, want %v", f.calls, want)
	}
}

// S-0745: a harness image that fails stops Ensure with its error, and the
// layer over it is not built.
func TestEnsureBuildsNoLayerAfterAFailedHarnessImage(t *testing.T) {
	t.Parallel()
	refused := &dashboard.Refusal{Version: pullV, Kind: dashboard.RefusalBad, Product: product("claude-code"), Reason: "the signature does not verify"}
	for name, f := range map[string]*ensureFake{
		"refused":      {pullErr: refused},
		"build failed": {pullErr: &pullError{err: errors.New("pull failed")}, buildErr: errors.New("build failed")},
	} {
		got, err := f.ensurer(&bytes.Buffer{}).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateMissing), layerStatus(StateStale)}, BuildSource{Version: pullV})
		if err == nil {
			t.Fatalf("%s: no error", name)
		}
		if name == "refused" && !errors.As(err, new(*dashboard.Refusal)) {
			t.Errorf("%s: error %v, want the refusal", name, err)
		}
		if len(f.layers) != 0 || len(got) != 0 {
			t.Errorf("%s: layers %v, outcomes %+v", name, f.layers, got)
		}
	}
}

// S-0745: a layer whose build fails stops Ensure with its error.
func TestEnsureSaysWhenTheLayerFails(t *testing.T) {
	t.Parallel()
	f := &ensureFake{layerErr: errors.New("layer failed")}
	got, err := f.ensurer(&bytes.Buffer{}).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateCurrent), layerStatus(StateMissing)}, BuildSource{Version: pullV})
	if err == nil || err.Error() != "layer failed" || len(got) != 1 {
		t.Fatalf("outcomes %+v, error %v", got, err)
	}
}

// S-0745: ModeLayer builds the layer, current or not, over the harness
// image on this host, current or stale, and neither pulls nor builds it.
func TestEnsureLayerModeBuildsTheLayerAlone(t *testing.T) {
	t.Parallel()
	for _, state := range []State{StateCurrent, StateStale} {
		f := &ensureFake{}
		var out bytes.Buffer
		got, err := f.ensurer(&out).Ensure(context.Background(), ModeLayer, []Status{harnessStatus(state), layerStatus(StateCurrent)}, BuildSource{Version: "dev"})
		if err != nil {
			t.Fatalf("%s: %v", state, err)
		}
		if want := []string{"layer " + ensureLayerRef}; !reflect.DeepEqual(f.calls, want) {
			t.Errorf("%s: calls %v, want %v", state, f.calls, want)
		}
		if len(got) != 2 || got[0].Way != WayLeft || got[1].Way != WayLayered {
			t.Errorf("%s: outcomes %+v", state, got)
		}
		if !strings.Contains(out.String(), "agent image: "+ensureRef+" is "+string(state)) || !strings.Contains(out.String(), "left as it is, as only the project's layer is built") {
			t.Errorf("%s: output\n%s", state, out.String())
		}
	}
}

// S-0745: a layer whose harness image is not on this host, and that this
// run did not make, is not built: Ensure fails naming the harness image and
// what makes it, flai image ensure for flai's default and the operator for
// their own.
func TestEnsureRefusesALayerWhoseHarnessImageIsMissing(t *testing.T) {
	t.Parallel()
	ownMissing := Status{Image: Image{Kind: KindHarness, Inner: "claude-code", Ref: "mine:1", Source: SourceHost}, State: StateMissing}
	for _, c := range []struct {
		name     string
		mode     Mode
		harness  Status
		wantSaid []string
	}{
		{"default, layer mode", ModeLayer, harnessStatus(StateMissing), []string{"the harness image " + ensureRef + " it builds over is not on this host", "run flai image ensure first"}},
		{"own, ensure", ModeEnsure, ownMissing, []string{"the harness image mine:1 it builds over, the operator's own", "flai never makes it"}},
		{"own, build", ModeBuild, ownMissing, []string{"the harness image mine:1", "flai never makes it"}},
	} {
		f := &ensureFake{}
		_, err := f.ensurer(&bytes.Buffer{}).Ensure(context.Background(), c.mode, []Status{c.harness, layerStatus(StateMissing)}, BuildSource{Version: pullV})
		if err == nil {
			t.Fatalf("%s: no error", c.name)
		}
		for _, want := range append(c.wantSaid, "build the layer image "+ensureLayerRef) {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("%s: error %q lacks %q", c.name, err, want)
			}
		}
		if len(f.calls) != 0 {
			t.Errorf("%s: calls %v", c.name, f.calls)
		}
	}
}

// S-0745: the layer over the operator's own image, on this host, is built
// over it, which is never touched.
func TestEnsureBuildsTheLayerOverTheOperatorsOwnImage(t *testing.T) {
	t.Parallel()
	own := Status{Image: Image{Kind: KindHarness, Inner: "claude-code", Ref: "mine:1", Source: SourceHost}, Present: true, State: StateOwn}
	f := &ensureFake{}
	got, err := f.ensurer(&bytes.Buffer{}).Ensure(context.Background(), ModeEnsure, []Status{own, layerStatus(StateStale)}, BuildSource{Version: pullV})
	if err != nil {
		t.Fatal(err)
	}
	if want := [][]string{{ensureRoot, ensureLayerRef, "mine:1"}}; !reflect.DeepEqual(f.layers, want) || len(f.calls) != 1 {
		t.Errorf("layers %v, calls %v", f.layers, f.calls)
	}
	if len(got) != 2 || got[0].Way != WayOwn || got[1].Way != WayLayered {
		t.Errorf("outcomes %+v", got)
	}
}

// A layer with no harness image of its inner harness among the statuses
// has nothing to build over.
func TestEnsureRefusesALayerWithNoHarnessImage(t *testing.T) {
	t.Parallel()
	f := &ensureFake{}
	_, err := f.ensurer(&bytes.Buffer{}).Ensure(context.Background(), ModeEnsure, []Status{layerStatus(StateMissing)}, BuildSource{Version: pullV})
	if err == nil || !strings.Contains(err.Error(), "no harness image of claude-code is among the agent images") || len(f.calls) != 0 {
		t.Fatalf("error %v, calls %v", err, f.calls)
	}
}

func TestEnsureRefusesADevelopmentFlaiWithNoTree(t *testing.T) {
	t.Parallel()
	f := &ensureFake{}
	var out bytes.Buffer
	_, err := f.ensurer(&out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateMissing)}, BuildSource{Version: "dev"})
	if err == nil || !strings.Contains(err.Error(), "a development flai has no release") || !strings.Contains(err.Error(), "flai/images/bedrock/Dockerfile and flai/images/claude-code/Dockerfile") {
		t.Fatalf("error %v", err)
	}
	if len(f.pulls)+len(f.builds) != 0 || out.Len() != 0 {
		t.Errorf("pulls %v, builds %v, output %q", f.pulls, f.builds, out.String())
	}
	// a current image needs nothing a release would give
	if _, err := f.ensurer(&out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateCurrent)}, BuildSource{Version: "dev"}); err != nil {
		t.Errorf("current image: %v", err)
	}
}

func TestEnsureRefusesAnUnknownMode(t *testing.T) {
	t.Parallel()
	if _, err := (&ensureFake{}).ensurer(&bytes.Buffer{}).Ensure(context.Background(), "push", nil, BuildSource{Version: pullV}); err == nil {
		t.Fatal("no error for an unknown mode")
	}
}

func TestHarnessTags(t *testing.T) {
	t.Parallel()
	for _, c := range []struct {
		ref, version string
		want         []string
	}{
		{ensureRef, "v1.2.3", []string{ensureRef, "flai-agent-claude-code:1.2.3"}},
		{"flai-agent-claude-code:1.2.3", "1.2.3", []string{"flai-agent-claude-code:1.2.3"}},
		{"", "1.2.3", []string{"flai-agent-claude-code:1.2.3"}},
	} {
		if got := HarnessTags(c.ref, "claude-code", c.version); !reflect.DeepEqual(got, c.want) {
			t.Errorf("HarnessTags(%q, %q) = %v, want %v", c.ref, c.version, got, c.want)
		}
	}
}

func TestLogPath(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 11, 13, 9, 8, 7, 0, time.FixedZone("x", 3600))
	want := filepath.Join("/root", ".flai-cache", "agent-images", "20261113T080807Z.log")
	if got := LogPath("/root", now); got != want {
		t.Errorf("LogPath = %q, want %q", got, want)
	}
}

// commandFake is the Exec of an Ensurer's build command: it records each
// run, writes a line, then blocks until its context ends when block is set,
// else lets make put what the command builds into images, as docker would,
// from the variables it got, and returns err.
type commandFake struct {
	mu     sync.Mutex
	runs   []commandCall
	images *fakeImages
	make   func(env map[string]string, images map[string]Info)
	err    error
	block  bool
}

// commandCall is one run of the build command: its program, folder,
// variables, and arguments.
type commandCall struct {
	program, dir string
	env          map[string]string
	args         []string
}

func (f *commandFake) exec(ctx context.Context, dir string, env []string, out io.Writer, program string, args ...string) error {
	vars := map[string]string{}
	for _, kv := range env {
		k, v, _ := strings.Cut(kv, "=")
		vars[k] = v
	}
	f.mu.Lock()
	f.runs = append(f.runs, commandCall{program: program, dir: dir, env: vars, args: args})
	f.mu.Unlock()
	if _, err := fmt.Fprintln(out, "the operator's machinery at work"); err != nil {
		return err
	}
	if f.block {
		<-ctx.Done()
		return ctx.Err()
	}
	if f.make != nil {
		f.make(vars, f.images.images)
	}
	return f.err
}

// makesAll builds what FLAI_IMAGE_STALE names as the variables say: the
// harness image labelled with the version, and the layer over it labelled
// with its recipe and the harness image's ID.
func makesAll(env map[string]string, images map[string]Info) {
	stale := strings.Fields(env["FLAI_IMAGE_STALE"])
	if slices.Contains(stale, PartHarness) {
		images[env["FLAI_IMAGE_HARNESS"]] = Info{ID: "sha256:made-harness", Labels: labels(env["FLAI_IMAGE_VERSION"], env["FLAI_IMAGE_RECIPE_HARNESS"], "")}
	}
	if slices.Contains(stale, PartLayer) {
		base := images[env["FLAI_IMAGE_HARNESS"]].ID
		images[env["FLAI_IMAGE_LAYER"]] = Info{ID: "sha256:made-layer", Labels: labels(env["FLAI_IMAGE_VERSION"], env["FLAI_IMAGE_RECIPE_LAYER"], base)}
	}
}

var ensureCommand = []string{"/opt/make-images", "--fast"}

// commandEnsurer is an Ensurer whose build command is ensureCommand, run by
// c, in a project at a new root with the layer's Dockerfile, on a host whose
// images are c.images, with f's pull, build, and layer, which it must not call.
func commandEnsurer(t *testing.T, f *ensureFake, c *commandFake, out *bytes.Buffer) Ensurer {
	t.Helper()
	if c.images == nil {
		c.images = &fakeImages{images: map[string]Info{}}
	}
	e := f.ensurer(out)
	e.Root = layerRoot(t, commandLayerData)
	e.Command = ensureCommand
	e.Run = Runner{Exec: c.exec, Sleep: (&sleeps{}).sleep}
	e.Inspect = c.images.inspect
	return e
}

// S-0759, ADR-0295 §7: with a build command named, a missing harness image
// and layer are made by the command alone, run once in the root with the
// variables, and judged again: the images it made are current.
func TestEnsureRunsTheBuildCommandInPlaceOfThePullAndTheBuild(t *testing.T) {
	t.Parallel()
	f := &ensureFake{}
	c := &commandFake{make: makesAll}
	var out bytes.Buffer
	e := commandEnsurer(t, f, c, &out)
	got, err := e.Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateMissing), layerStatus(StateMissing)}, BuildSource{Version: pullV})
	if err != nil {
		t.Fatalf("Ensure: %v\n%s", err, out.String())
	}
	if len(f.calls) != 0 {
		t.Errorf("flai's own pull and build ran: %v", f.calls)
	}
	if len(c.runs) != 1 {
		t.Fatalf("the command ran %d times, want once", len(c.runs))
	}
	run := c.runs[0]
	if run.program != ensureCommand[0] || !reflect.DeepEqual(run.args, ensureCommand[1:]) || run.dir != e.Root {
		t.Errorf("ran %s %v in %s", run.program, run.args, run.dir)
	}
	if run.env["FLAI_IMAGE_STALE"] != "bedrock harness layer" || run.env["FLAI_IMAGE_LAYER"] != ensureLayerRef || run.env["FLAI_IMAGE_VERSION"] != pullV {
		t.Errorf("variables %v", run.env)
	}
	if len(got) != 2 || got[0].Way != WayCommand || got[1].Way != WayCommand {
		t.Errorf("outcomes %+v", got)
	}
	after, err := Check(e.Root, []Image{harnessStatus(StateMissing).Image, layerStatus(StateMissing).Image}, c.images.inspect, func() (string, error) { return pullV, nil })
	if err != nil || after[0].State != StateCurrent || after[1].State != StateCurrent {
		t.Errorf("after: %+v, %v", after, err)
	}
	text := out.String()
	for _, want := range []string{
		"agent image: " + ensureRef + " is missing; the operator's build command /opt/make-images --fast makes it, in place of flai's pull and build\n",
		"agent image: running the operator's build command /opt/make-images --fast in " + e.Root + " for the images of claude-code",
		"the operator's machinery at work\n",
		"agent image: " + ensureLayerRef + ": made by the operator's build command /opt/make-images --fast\n",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("output lacks %q:\n%s", want, text)
		}
	}
}

// S-0759: with every image current, ensure mode runs no command.
func TestEnsureRunsNoBuildCommandWhenAllIsCurrent(t *testing.T) {
	t.Parallel()
	f := &ensureFake{}
	c := &commandFake{make: makesAll}
	var out bytes.Buffer
	got, err := commandEnsurer(t, f, c, &out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateCurrent), layerStatus(StateCurrent)}, BuildSource{Version: pullV})
	if err != nil || len(c.runs) != 0 || len(f.calls) != 0 || len(got) != 2 || got[0].Way != WayCurrent || got[1].Way != WayCurrent {
		t.Fatalf("outcomes %+v, %v; runs %v, calls %v", got, err, c.runs, f.calls)
	}
}

// S-0759: build mode runs the command, current or not, asking for every
// image; layer mode asks for the layer alone and leaves, and does not
// judge, a stale harness image under it.
func TestEnsureBuildAndLayerModesRunTheBuildCommand(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		mode    Mode
		harness State
		stale   string
		ways    []Way
	}{
		{ModeBuild, StateCurrent, "bedrock harness layer", []Way{WayCommand, WayCommand}},
		{ModeLayer, StateStale, "layer", []Way{WayLeft, WayCommand}},
	} {
		f := &ensureFake{}
		c := &commandFake{make: makesAll, images: &fakeImages{images: map[string]Info{ensureRef: {ID: "sha256:old", Labels: labels("1.0.0", "", "")}}}}
		var out bytes.Buffer
		if tc.mode == ModeLayer {
			c.make = func(env map[string]string, images map[string]Info) {
				images[env["FLAI_IMAGE_LAYER"]] = Info{ID: "sha256:made-layer", Labels: labels(env["FLAI_IMAGE_VERSION"], env["FLAI_IMAGE_RECIPE_LAYER"], "sha256:old")}
			}
		}
		got, err := commandEnsurer(t, f, c, &out).Ensure(context.Background(), tc.mode, []Status{harnessStatus(tc.harness), layerStatus(StateCurrent)}, BuildSource{Version: pullV})
		if err != nil {
			t.Fatalf("%s: %v\n%s", tc.mode, err, out.String())
		}
		if len(c.runs) != 1 || c.runs[0].env["FLAI_IMAGE_STALE"] != tc.stale || len(f.calls) != 0 {
			t.Errorf("%s: runs %+v, calls %v", tc.mode, c.runs, f.calls)
		}
		if len(got) != 2 || got[0].Way != tc.ways[0] || got[1].Way != tc.ways[1] {
			t.Errorf("%s: outcomes %+v", tc.mode, got)
		}
	}
}

// S-0759: a command that exits non-zero fails the run, naming it, and
// nothing of flai's own runs in its place.
func TestEnsureFailsWhenTheBuildCommandExitsNonZero(t *testing.T) {
	t.Parallel()
	f := &ensureFake{}
	c := &commandFake{err: errExit}
	var out bytes.Buffer
	got, err := commandEnsurer(t, f, c, &out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateMissing)}, BuildSource{Version: pullV})
	var failed *CommandError
	if !errors.As(err, &failed) || !strings.Contains(err.Error(), "/opt/make-images --fast failed for the agent images of claude-code (the operator's machinery at work; exit status 1)") {
		t.Fatalf("Ensure = %v, want a *CommandError naming the command", err)
	}
	if len(got) != 0 || len(f.calls) != 0 || len(c.runs) != 1 {
		t.Errorf("outcomes %+v, calls %v, runs %d", got, f.calls, len(c.runs))
	}
}

// S-0759, ADR-0295 §7: a command that exits 0 but leaves an image missing,
// or carrying other labels, fails the run, naming the command, the image,
// and why.
func TestEnsureFailsWhenTheBuildCommandLeavesAnImageWrong(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		make func(env map[string]string, images map[string]Info)
		want string
	}{
		{"nothing made", nil, ensureRef + " is missing; and " + ensureLayerRef + " is missing"},
		{"other labels", func(env map[string]string, images map[string]Info) {
			images[env["FLAI_IMAGE_HARNESS"]] = Info{ID: "sha256:h", Labels: labels("1.0.0", "", "")}
			images[env["FLAI_IMAGE_LAYER"]] = Info{ID: "sha256:l", Labels: labels(pullV, "sha256:other", "sha256:h")}
		}, ensureRef + " is stale (built for flai 1.0.0, not " + pullV + "); and " + ensureLayerRef + " is stale (built from another Dockerfile"},
	} {
		f := &ensureFake{}
		c := &commandFake{make: tc.make}
		var out bytes.Buffer
		_, err := commandEnsurer(t, f, c, &out).Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateMissing), layerStatus(StateMissing)}, BuildSource{Version: pullV})
		if err == nil || !strings.Contains(err.Error(), "the operator's build command /opt/make-images --fast exited 0 but left the agent images wrong: "+tc.want) || !strings.Contains(err.Error(), "FLAI_IMAGE_* variables") {
			t.Errorf("%s: %v", tc.name, err)
		}
	}
}

// S-0759, ADR-0295 §5: a command silent past the stall limit is stopped,
// and the run fails with a *StallError naming it.
func TestEnsureStopsASilentBuildCommand(t *testing.T) {
	t.Parallel()
	f := &ensureFake{}
	c := &commandFake{block: true}
	var out bytes.Buffer
	e := commandEnsurer(t, f, c, &out)
	e.Run.Stall = 50 * time.Millisecond
	_, err := e.Ensure(context.Background(), ModeEnsure, []Status{harnessStatus(StateMissing)}, BuildSource{Version: pullV})
	var stall *StallError
	if !errors.As(err, &stall) || stall.Step != "the build command /opt/make-images --fast" {
		t.Fatalf("Ensure = %v, want a *StallError naming the command", err)
	}
}

// S-0759, ADR-0295 §3, §7: the operator's own image is never asked made nor
// judged, the layer over it is, and a development flai with no tree runs
// the command where its own build would refuse.
func TestEnsureBuildCommandLeavesTheOperatorsOwnImage(t *testing.T) {
	t.Parallel()
	own := Status{Image: Image{Kind: KindHarness, Inner: "claude-code", Ref: "mine:1", Source: SourceHost}, Present: true, State: StateOwn}
	f := &ensureFake{}
	c := &commandFake{images: &fakeImages{images: map[string]Info{"mine:1": {ID: "sha256:mine"}}}}
	c.make = func(env map[string]string, images map[string]Info) {
		images[env["FLAI_IMAGE_LAYER"]] = Info{ID: "sha256:l", Labels: labels("dev", env["FLAI_IMAGE_RECIPE_LAYER"], "sha256:mine")}
	}
	var out bytes.Buffer
	got, err := commandEnsurer(t, f, c, &out).Ensure(context.Background(), ModeEnsure, []Status{own, layerStatus(StateMissing)}, BuildSource{Version: "dev"})
	if err != nil {
		t.Fatalf("Ensure: %v\n%s", err, out.String())
	}
	if len(c.runs) != 1 || c.runs[0].env["FLAI_IMAGE_STALE"] != "layer" || c.runs[0].env["FLAI_IMAGE_HARNESS"] != "mine:1" || c.runs[0].env["FLAI_IMAGE_BEDROCK"] != "" {
		t.Errorf("runs %+v", c.runs)
	}
	if len(got) != 2 || got[0].Way != WayOwn || got[1].Way != WayCommand {
		t.Errorf("outcomes %+v", got)
	}
}
