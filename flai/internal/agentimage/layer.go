package agentimage

import (
	"context"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// Layered is what Builder.Layer made: the layer image, the tags given it,
// the harness image it was built over and that image's ID, and the recipe
// it is labelled with.
type Layered struct {
	Image  Image
	Tags   []string
	Base   string
	BaseID string
	Recipe string
}

// Layer builds the project's layer image img, in the project at root, over
// the harness image base, which must be on this host (ADR-0295 §2, §4): from
// img's Dockerfile with the project's root as the context, BASE set to base
// and each of img's build arguments passed, labelled with its recipe and
// base's ID. It says what it builds on b.Run.Out before it runs the build.
func (b Builder) Layer(ctx context.Context, root string, img Image, base string) (Layered, error) {
	switch {
	case img.Kind != KindLayer:
		return Layered{}, fmt.Errorf("build the layer image %s: it is a %s image, not the project's layer from agent.image.layer; name the layer image For gives under a manifest that names agent.image.layer", shown(img.Ref), img.Kind)
	case img.Ref == "" || img.Dockerfile == "":
		return Layered{}, fmt.Errorf("build the layer image %s: it names no tag or no Dockerfile (%s); name the layer image For gives under a manifest that names agent.image.layer", shown(img.Ref), shown(img.Dockerfile))
	}
	if base == "" {
		return Layered{}, fmt.Errorf("build the layer image %s: no harness image was named to build it over; name the harness image of %s", img.Ref, img.Inner)
	}
	data, err := layerDockerfile(root, img, "to build its image "+img.Ref)
	if err != nil {
		return Layered{}, err
	}
	info, present, err := b.inspector()(base)
	if err != nil {
		return Layered{}, fmt.Errorf("read the ID of the harness image %s, which the layer image %s builds over: %w", base, img.Ref, err)
	}
	if !present {
		return Layered{}, fmt.Errorf("build the layer image %s: the harness image %s it builds over is not on this host; run flai image ensure first, then build the layer again", img.Ref, base)
	}
	recipe := Recipe(data, "", img.Args)
	fmt.Fprintf(b.out(), "agent image: building %s over %s (%s) from %s with %s\n", img.Ref, base, info.ID, img.Dockerfile, argWords(img.Args))
	args := []string{"build", "-f", filepath.FromSlash(img.Dockerfile), "--build-arg", manifest.BaseArg + "=" + base}
	for _, name := range slices.Sorted(maps.Keys(img.Args)) {
		args = append(args, "--build-arg", name+"="+img.Args[name])
	}
	args = append(args, "--label", LabelRecipe+"="+recipe, "--label", LabelBase+"="+info.ID, "-t", img.Ref, ".")
	if err := b.Run.Run(ctx, Step{Name: "build " + img.Ref, Dir: root, Args: args}); err != nil {
		return Layered{}, err
	}
	return Layered{Image: img, Tags: []string{img.Ref}, Base: base, BaseID: info.ID, Recipe: recipe}, nil
}

// layerDockerfile is the bytes of the Dockerfile the layer image img is
// built from, under root, read for the purpose why.
func layerDockerfile(root string, img Image, why string) ([]byte, error) {
	path := filepath.Join(root, filepath.FromSlash(img.Dockerfile))
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s, the Dockerfile agent.image.layer names, %s: %w; fix agent.image.layer in %s or add the Dockerfile", path, why, err, manifest.File)
	}
	return data, nil
}

// argWords is the build arguments args in the words a line says them: each
// NAME=value by name, or the Dockerfile's own ARG defaults when there are none.
func argWords(args map[string]string) string {
	if len(args) == 0 {
		return "its Dockerfile's own ARG defaults"
	}
	words := make([]string, 0, len(args))
	for _, name := range slices.Sorted(maps.Keys(args)) {
		words = append(words, name+"="+args[name])
	}
	return "the build arguments " + strings.Join(words, ", ")
}
