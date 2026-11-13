package agentimage

import (
	"bytes"
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/dashboard"
)

// ensureFake records the pulls and builds an Ensurer asks of it, answering
// each pull with pullErr and each build with buildErr.
type ensureFake struct {
	pulls    [][]string // image, version, then the tags
	builds   []BuildSource
	tags     [][]string
	pullErr  error
	buildErr error
}

func (f *ensureFake) pull(_ context.Context, image, version string, tags []string) (Pulled, error) {
	f.pulls = append(f.pulls, append([]string{image, version}, tags...))
	if f.pullErr != nil {
		return Pulled{}, f.pullErr
	}
	return Pulled{Image: image, Version: version, Reference: Published(image) + "@" + indexD, Tags: tags}, nil
}

func (f *ensureFake) build(_ context.Context, inner string, src BuildSource, tags []string) (Built, error) {
	f.builds = append(f.builds, src)
	f.tags = append(f.tags, tags)
	if f.buildErr != nil {
		return Built{}, f.buildErr
	}
	from := FromRelease
	if src.Tree != "" {
		from = FromTree
	}
	return Built{Inner: inner, Version: src.Version, Tags: tags, Bedrock: BedrockTag(src.Version), From: from}, nil
}

// ensurer is an Ensurer over f, writing to out.
func (f *ensureFake) ensurer(out *bytes.Buffer) Ensurer {
	return Ensurer{Pull: f.pull, Build: f.build, Out: out}
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

func TestEnsureLeavesTheOperatorsOwnImageAndTheLayer(t *testing.T) {
	t.Parallel()
	own := Status{Image: Image{Kind: KindHarness, Inner: "claude-code", Ref: "mine:1", Source: SourceHost}, Present: true, State: StateOwn}
	layer := Status{Image: Image{Kind: KindLayer, Inner: "claude-code", Ref: "flai-agent-t-claude-code:local", Source: SourceLayer, Dockerfile: "images/agent/Dockerfile"}, State: StateStale}
	for _, mode := range []Mode{ModeEnsure, ModeBuild} {
		f := &ensureFake{}
		var out bytes.Buffer
		got, err := f.ensurer(&out).Ensure(context.Background(), mode, []Status{own, layer}, BuildSource{Version: "dev"})
		if err != nil {
			t.Fatalf("%s: %v", mode, err)
		}
		if len(f.pulls)+len(f.builds) != 0 || len(got) != 2 || got[0].Way != WayOwn || got[1].Way != WayNotBuilt {
			t.Errorf("%s: outcomes %+v; pulls %v, builds %v", mode, got, f.pulls, f.builds)
		}
		text := out.String()
		if !strings.Contains(text, "mine:1 is the operator's own") || !strings.Contains(text, "the project's own build makes it") {
			t.Errorf("%s: output\n%s", mode, text)
		}
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
