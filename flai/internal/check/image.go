package check

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/harness"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// ErrNoDocker is an ImageProbe's answer when docker is not on PATH.
var ErrNoDocker = errors.New("docker is not on PATH")

// ImageProbe reports whether image is built on this host: ErrNoDocker when
// docker is not on PATH, another error when docker could not say.
type ImageProbe func(image string) (bool, error)

// DockerImage is the ImageProbe that asks docker image inspect through r.
func DockerImage(r execx.Runner) ImageProbe {
	return func(image string) (bool, error) {
		if _, err := r.LookPath("docker"); err != nil {
			return false, ErrNoDocker
		}
		out, err := r.Run("", "docker", "image", "inspect", "--format", "{{.Id}}", image)
		if err == nil {
			return true, nil
		}
		// "Error: No such image" or "Error response from daemon: No such
		// image", by the client's version; anything else, such as a daemon
		// not running, leaves the image unknown.
		if strings.Contains(strings.ToLower(out+" "+err.Error()), "no such image") {
			return false, nil
		}
		return false, err
	}
}

// AgentImage adds agent.image when system-flow.yaml names a layer under
// agent.image.layer whose Dockerfile does not exist, a warning, or whose
// image, built over the claude-code harness image, is not on this host, a
// warning --strict passes over, since only building it clears it, which no
// story branch does. When docker is not on PATH or cannot say, it adds a note
// --strict passes over that the image went unchecked. Those three are marked
// unchecked, with the manifest and the layer as the paths whose change brings
// them inside a story, so a check scoped to a story that changes neither
// leaves them out (ADR-0255, I-0264); the warning for a layer that is not a
// file is not, since a story can fix it. A manifest naming no layer gets
// nothing. The caller supplies the probe so this package runs no commands
// (S-0510).
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
	if st, err := os.Stat(filepath.Join(repo.Root, filepath.FromSlash(layer))); err != nil || st.IsDir() {
		add(false, "agent.image.layer names "+layer+", which is not a file: fix the path in "+manifest.File+" or add the Dockerfile, then build its image with `make agent-images`")
		return
	}
	image := repo.Manifest.LayerImage(harness.ClaudeCode)
	built, err := probe(image)
	switch {
	case errors.Is(err, ErrNoDocker):
		add(true, "the image "+image+" of agent.image.layer "+layer+" is not checked: docker is not on PATH; install Docker to check it, and build it with `make agent-images`")
	case err != nil:
		add(true, "the image "+image+" of agent.image.layer "+layer+" is not checked: docker could not say whether it is built ("+firstLine(err.Error())+"); start Docker to check it")
	case !built:
		add(true, "the image "+image+" of agent.image.layer "+layer+" is not built on this host, so flai serve cannot run the project's agents in it: build it with `make agent-images`, or `scripts/agent-image.sh --layer` over the harness image already built")
	}
}

// layerLine returns the 1-based line of layer: under agent: and image: in a
// manifest as flai writes it, or 1.
func layerLine(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 1
	}
	defer func() { _ = f.Close() }()
	sc := bufio.NewScanner(f)
	agent, image, n := false, false, 0
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
			agent, image = strings.HasPrefix(l, "agent:"), false
		case agent && indent == 2:
			image = strings.HasPrefix(trimmed, "image:")
		case agent && image && strings.HasPrefix(trimmed, "layer:"):
			return n
		}
	}
	return 1
}

// firstLine is s up to its first newline.
func firstLine(s string) string {
	line, _, _ := strings.Cut(s, "\n")
	return line
}
