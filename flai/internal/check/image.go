package check

import (
	"bufio"
	"errors"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/agentimage"
	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ErrNoDocker is an ImageProbe's Inspect answer when docker is not on PATH,
// agentimage's own.
var ErrNoDocker = agentimage.ErrNoDocker

// ImageProbe is how AgentImage asks this host about the agent images:
// Inspect says whether each is here, and its labels; Runner runs the
// project's flai.command for the flai version they must hold
// (agentimage.ProjectFlai), and must be set for a manifest that names one;
// HostImages, when set, is the image the host's container harness names for
// each inner harness, for the project whose main checkout is main, which is
// the operator's own and never judged.
type ImageProbe struct {
	Inspect    agentimage.Inspector
	Runner     execx.Runner
	HostImages func(main string) map[string]string
}

// DockerImage is the ImageProbe that asks docker image inspect through r,
// runs the project's flai through r, and reads the host's images from the
// configuration file (HostImages).
func DockerImage(r execx.Runner) ImageProbe {
	return ImageProbe{Inspect: agentimage.Docker(r, "docker"), Runner: r, HostImages: HostImages}
}

// HostImages is the image the host's container harness names for each inner
// harness, with flai serve agent harness container --image, for the project
// whose main checkout is main, the project's own entry over the host-wide
// one (ADR-0261), as flai serve reads them: from the configuration file
// FLAI_CONFIG names, else ~/.flai/config.json. A file that is missing or
// cannot be read or parsed names none, and is never written; flai serve
// refuses such a file itself.
func HostImages(main string) map[string]string {
	data, err := os.ReadFile(config.ResolvePath(""))
	if err != nil {
		return nil
	}
	cfg, err := config.Parse(data)
	if err != nil {
		return nil
	}
	return cfg.Agent.ForProject(main).Harnesses[config.ContainerHarness].Images
}

// AgentImage adds agent.image when system-flow.yaml names a layer under
// agent.image.layer whose Dockerfile does not exist, a warning, and, for
// each image the project's agents need (agentimage.Needed: the harness
// images and the layer built over claude-code's) that is not on this host or
// is stale (ADR-0295), a warning --strict passes over, since only building
// it clears it, which no story branch does; each names flai image status.
// Each name agent.image.args gives that no ARG of the layer declares is a
// warning on its line, which a story can fix.
// The host's own images, named with --image, are never judged. When docker
// is not on PATH or cannot say, or the images cannot be judged, it adds a
// note --strict passes over that the images went unchecked. Those are marked
// unchecked, with the manifest and the layer as the paths whose change
// brings them inside a story, so a check scoped to a story that changes
// neither leaves them out (ADR-0255, I-0264); the warnings for a layer that
// is not a file and for an undeclared arg are not, since a story can fix
// them. A manifest naming no
// layer gets nothing. The caller supplies the probe so this package runs no
// commands (S-0510).
func AgentImage(res *Result, repo *workitem.Repo, probe ImageProbe) {
	layer := repo.Manifest.Layer()
	if layer == "" {
		return
	}
	mf := filepath.Join(repo.Root, manifest.File)
	line := layerLine(mf)
	add := func(advisory bool, msg string) {
		f := Finding{Level: Warning, Rule: "agent.image", Path: manifest.File, Line: line, Message: msg, advisory: advisory}
		if advisory {
			// Only the host clears it (ADR-0255).
			f.unchecked, f.inside = true, []string{manifest.File, layer}
		}
		res.Findings = append(res.Findings, f)
		res.Warnings++
		if advisory {
			res.Advisory++
		}
	}
	// The manifest's own check keeps layer a clean path inside the repository.
	dockerfile := filepath.Join(repo.Root, filepath.FromSlash(layer))
	if st, err := os.Stat(dockerfile); err != nil || st.IsDir() {
		add(false, "agent.image.layer names "+layer+", which is not a file: fix the path in "+manifest.File+" or add the Dockerfile, then build its image with `make agent-images`")
		return
	}
	undeclaredArgs(res, mf, layer, dockerfile, repo.Manifest.LayerArgs())
	// The main checkout, as flai serve and flai image status judge the images
	// from it: the layer's Dockerfile, the project's flai, and the host's
	// images there, so a check in a story's worktree judges them the same.
	main := repo.MainRoot
	if main == "" {
		main = repo.Root
	}
	var host map[string]string
	if probe.HostImages != nil {
		host = probe.HostImages(main)
	}
	images := agentimage.Needed(repo.Manifest, agentimage.Agents(repo.Manifest), host)
	var inspectErr error
	inspect := func(ref string) (agentimage.Info, bool, error) {
		info, present, err := probe.Inspect(ref)
		if err != nil {
			inspectErr = err
		}
		return info, present, err
	}
	flai := func() (string, error) {
		return agentimage.ProjectFlai(main, repo.Manifest, probe.Runner, buildinfo.Get().Version)
	}
	statuses, err := agentimage.Check(main, images, inspect, flai)
	image := repo.Manifest.LayerImage(harness.ClaudeCode)
	switch {
	case errors.Is(inspectErr, ErrNoDocker):
		add(true, "the image "+image+" of agent.image.layer "+layer+" is not checked: docker is not on PATH; install Docker to check it, and build it with `make agent-images`")
		return
	case inspectErr != nil:
		add(true, "the image "+image+" of agent.image.layer "+layer+" is not checked: docker could not say whether it is built ("+firstLine(inspectErr.Error())+"); start Docker to check it")
		return
	case err != nil:
		add(true, "the images of agent.image.layer "+layer+" are not checked: "+firstLine(err.Error())+"; see `flai image status`")
		return
	}
	for _, s := range statuses {
		if msg := imageFinding(s, layer); msg != "" {
			add(true, msg)
		}
	}
}

// imageFinding is what AgentImage says of the agent image s, under the layer
// agent.image.layer names: that it is not built or is stale, and what to do,
// or "" for an image that is current or the operator's own.
func imageFinding(s agentimage.Status, layer string) string {
	what := "the image " + s.Ref + " of agent.image.layer " + layer
	switch {
	case s.Kind == agentimage.KindHarness && s.Source == agentimage.SourceHost:
		what = "the image " + s.Ref + " the host names for the inner harness " + s.Inner
	case s.Kind == agentimage.KindHarness:
		what = "the harness image " + s.Ref + " of the inner harness " + s.Inner
	}
	switch s.State {
	case agentimage.StateMissing:
		switch {
		case s.Kind == agentimage.KindLayer:
			return what + " is not built on this host, so flai serve cannot run the project's agents in it: see `flai image status`, and build it with `make agent-images`, or `scripts/agent-image.sh --layer` over the harness image already built"
		case s.Source == agentimage.SourceHost:
			return what + " is not on this host: see `flai image status`, and build it, or name another with `flai serve agent harness container --image " + s.Inner + "=<image>`"
		}
		return what + " is not built on this host, so flai serve cannot run the project's agents in it, nor build a layer over it: see `flai image status`, and build it with `make agent-images`"
	case agentimage.StateStale:
		why := make([]string, 0, len(s.Reasons))
		for _, r := range s.Reasons {
			why = append(why, r.Message)
		}
		return what + " is stale: " + strings.Join(why, "; ") + "; see `flai image status`, and rebuild it with `make agent-images`"
	}
	return ""
}

// undeclaredArgs adds an agent.image warning, on its line of the manifest at
// mf, for each of args the layer's Dockerfile declares with no ARG, which the
// build would pass to nothing; a story can fix it, so it is not advisory.
func undeclaredArgs(res *Result, mf, layer, dockerfile string, args map[string]string) {
	if len(args) == 0 {
		return
	}
	add := func(line int, msg string) {
		res.Findings = append(res.Findings, Finding{Level: Warning, Rule: "agent.image", Path: manifest.File, Line: line, Message: msg})
		res.Warnings++
	}
	data, err := os.ReadFile(dockerfile)
	if err != nil {
		add(layerLine(mf), "agent.image.args are not checked against the layer "+layer+": "+firstLine(err.Error())+"; make the Dockerfile readable")
		return
	}
	declared := dockerfileArgs(data)
	for _, n := range slices.Sorted(maps.Keys(args)) {
		if !declared[n] {
			add(argLine(mf, n), "agent.image.args names "+n+", which no ARG of the layer "+layer+" declares, so its build would not use it: add `ARG "+n+"` to the Dockerfile, or leave "+n+" out of agent.image.args in "+manifest.File)
		}
	}
}

// dockerfileArgs are the names the ARG instructions of a Dockerfile declare,
// in any stage: ARG NAME and ARG NAME=default, several on one instruction,
// over continued lines, with comment lines left out.
func dockerfileArgs(data []byte) map[string]bool {
	out := map[string]bool{}
	var logical strings.Builder
	flush := func() {
		fields := dockerfileFields(logical.String())
		logical.Reset()
		if len(fields) < 2 || !strings.EqualFold(fields[0], "ARG") {
			return
		}
		for _, f := range fields[1:] {
			name, _, _ := strings.Cut(f, "=")
			out[name] = true
		}
	}
	for _, l := range strings.Split(string(data), "\n") {
		l = strings.TrimRight(l, " \t\r")
		if strings.HasPrefix(strings.TrimSpace(l), "#") {
			continue
		}
		if cont, ok := strings.CutSuffix(l, `\`); ok {
			logical.WriteString(cont + " ")
			continue
		}
		logical.WriteString(l)
		flush()
	}
	flush()
	return out
}

// dockerfileFields splits an instruction at its blanks, keeping a quoted
// value whole.
func dockerfileFields(s string) []string {
	var fields []string
	var cur strings.Builder
	var quote rune
	for _, r := range s {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			}
			cur.WriteRune(r)
		case r == '"' || r == '\'':
			quote = r
			cur.WriteRune(r)
		case r == ' ' || r == '\t':
			if cur.Len() > 0 {
				fields = append(fields, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		fields = append(fields, cur.String())
	}
	return fields
}

// layerLine returns the 1-based line of layer: under agent: and image: in a
// manifest as flai writes it, or 1.
func layerLine(path string) int {
	layer, _ := imageLines(path, "")
	return layer
}

// argLine returns the 1-based line of name's entry under agent.image.args in
// a manifest as flai writes it, else layerLine's.
func argLine(path, name string) int {
	layer, entry := imageLines(path, name)
	if entry == 0 {
		return layer
	}
	return entry
}

// imageLines are the 1-based lines of layer: under agent: and image: in the
// manifest at path, or 1, and of arg's entry under its args:, or 0.
func imageLines(path, arg string) (layer, entry int) {
	layer = 1
	f, err := os.Open(path)
	if err != nil {
		return layer, 0
	}
	defer func() { _ = f.Close() }()
	sc := bufio.NewScanner(f)
	agent, image, args, foundLayer := false, false, false, false
	child, n := 0, 0
	for sc.Scan() {
		n++
		l := sc.Text()
		trimmed := strings.TrimSpace(l)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		indent := len(l) - len(strings.TrimLeft(l, " "))
		switch {
		case indent == 0:
			agent, image, args = strings.HasPrefix(l, "agent:"), false, false
		case agent && indent == 2:
			image, args, child = strings.HasPrefix(trimmed, "image:"), false, 0
		case agent && image && (child == 0 || indent <= child):
			// A key of image itself, at the indent of its first.
			child = indent
			args = strings.HasPrefix(trimmed, "args:")
			if !foundLayer && strings.HasPrefix(trimmed, "layer:") {
				layer, foundLayer = n, true
			}
		case args && arg != "" && entry == 0 && strings.HasPrefix(trimmed, arg+":"):
			entry = n
		}
	}
	return layer, entry
}

// firstLine is s up to its first newline.
func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}
