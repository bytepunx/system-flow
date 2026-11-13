package agentimage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/images"
	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
)

// The variables the operator's build command gets beside the inherited
// environment (ADR-0295 §7).
const (
	// EnvVersion is the flai version the images must hold, without a leading v.
	EnvVersion = "FLAI_IMAGE_VERSION"
	// EnvInner is the inner harness the images are for, such as claude-code.
	EnvInner = "FLAI_IMAGE_INNER"
	// EnvBedrock is the tag the bedrock image must have.
	EnvBedrock = "FLAI_IMAGE_BEDROCK"
	// EnvHarness is the tag the harness image must have, the name a start runs it by.
	EnvHarness = "FLAI_IMAGE_HARNESS"
	// EnvLayer is the tag the project's layer must have, empty when the manifest names none.
	EnvLayer = "FLAI_IMAGE_LAYER"
	// EnvRecipeBedrock is the recipe label the bedrock image must carry.
	EnvRecipeBedrock = "FLAI_IMAGE_RECIPE_BEDROCK"
	// EnvRecipeHarness is the recipe label the harness image must carry.
	EnvRecipeHarness = "FLAI_IMAGE_RECIPE_HARNESS"
	// EnvRecipeLayer is the recipe label the project's layer must carry.
	EnvRecipeLayer = "FLAI_IMAGE_RECIPE_LAYER"
	// EnvStale is the images the run asks made, of bedrock, harness, and
	// layer, in that order, separated by spaces.
	EnvStale = "FLAI_IMAGE_STALE"
)

// The images of one inner harness, as FLAI_IMAGE_STALE names them.
const (
	PartBedrock = "bedrock"
	PartHarness = "harness"
	PartLayer   = "layer"
)

// CommandImages is what the operator's build command is asked to make for
// one inner harness, as its FLAI_IMAGE_* variables say it: the flai version
// the images must hold, the inner harness, the tags and the recipe labels
// the bedrock, harness, and layer images must have ("" where flai has none
// to give), and the images the run asks made, of PartBedrock, PartHarness,
// and PartLayer, in that order.
type CommandImages struct {
	Version       string
	Inner         string
	Bedrock       string
	Harness       string
	Layer         string
	RecipeBedrock string
	RecipeHarness string
	RecipeLayer   string
	Stale         []string
}

// Env is c as NAME=value pairs, every FLAI_IMAGE_* variable named, an empty
// one with an empty value, in the order ADR-0295 §7 lists them.
func (c CommandImages) Env() []string {
	return []string{
		EnvVersion + "=" + c.Version,
		EnvInner + "=" + c.Inner,
		EnvBedrock + "=" + c.Bedrock,
		EnvHarness + "=" + c.Harness,
		EnvLayer + "=" + c.Layer,
		EnvRecipeBedrock + "=" + c.RecipeBedrock,
		EnvRecipeHarness + "=" + c.RecipeHarness,
		EnvRecipeLayer + "=" + c.RecipeLayer,
		EnvStale + "=" + strings.Join(c.Stale, " "),
	}
}

// commandImages is what the build command is asked to make, as mode asks,
// for the inner harness whose harness image is h and whose layer is l (nil
// when the manifest names none), of the project at root, for the flai src
// names. FLAI_IMAGE_STALE holds, in ModeEnsure, the harness image of flai's
// default and bedrock under it when that image is missing or stale, and the
// layer when it is missing or stale or its harness image is asked made; in
// ModeBuild every image flai would build, current or not; in ModeLayer the
// layer alone. A harness image the operator set with --image is never asked
// made: its tag is the one the layer builds over, and bedrock and the
// harness recipe are left empty. The recipes are those flai's own build
// would label (ADR-0295 §2): bedrock's prebuilt target and the Dockerfiles
// flai carries for a release, the tree's own Dockerfiles and bedrock's tree
// target for a development tree, none for a development flai with no tree.
func commandImages(root string, mode Mode, h Status, l *Status, src BuildSource) (CommandImages, error) {
	version := strings.TrimPrefix(src.Version, "v")
	c := CommandImages{Version: version, Inner: h.Inner, Harness: h.Ref}
	own := h.Source == SourceHost
	if !own {
		c.Bedrock = BedrockTag(version)
		bedrock, inner, err := harnessRecipes(h.Inner, src)
		if err != nil {
			return CommandImages{}, err
		}
		c.RecipeBedrock, c.RecipeHarness = bedrock, inner
	}
	harnessAsked := !own && (mode == ModeBuild || (mode == ModeEnsure && h.State != StateCurrent))
	if harnessAsked {
		c.Stale = append(c.Stale, PartBedrock, PartHarness)
	}
	if l == nil {
		return c, nil
	}
	c.Layer = l.Ref
	data, err := layerDockerfile(root, l.Image, "for the recipe its image "+l.Ref+" must carry")
	if err != nil {
		return CommandImages{}, err
	}
	c.RecipeLayer = Recipe(data, "", l.Args)
	if mode != ModeEnsure || l.State != StateCurrent || harnessAsked {
		c.Stale = append(c.Stale, PartLayer)
	}
	return c, nil
}

// harnessRecipes are the recipe labels flai's own build gives the bedrock
// image and inner's harness image for src: from the tree's own Dockerfiles,
// bedrock at its tree target, when src names a tree; none for a development
// flai with no tree; else from the Dockerfiles flai carries, bedrock at its
// prebuilt target, as the release's published images are built.
func harnessRecipes(inner string, src BuildSource) (string, string, error) {
	var bedrock, harness []byte
	target := targetPrebuilt
	switch {
	case src.Tree != "":
		var err error
		if _, bedrock, err = treeDockerfile(src.Tree, images.Bedrock); err != nil {
			return "", "", err
		}
		if _, harness, err = treeDockerfile(src.Tree, inner); err != nil {
			return "", "", err
		}
		target = targetTree
	case buildinfo.Development(src.Version):
		return "", "", nil
	default:
		var err error
		if bedrock, err = images.Dockerfile(images.Bedrock); err != nil {
			return "", "", err
		}
		if harness, err = images.Dockerfile(inner); err != nil {
			return "", "", err
		}
	}
	return Recipe(bedrock, target, nil), Recipe(harness, "", nil), nil
}

// RunCommand runs argv, the operator's build command (ADR-0295 §7), in
// root, the project's main checkout, with no shell, argv[0] being the
// program, through r, so that its output streams to r.Out and a command
// silent for r.Stall is stopped with a *StallError, as flai's own build
// steps are. Its environment is the inherited one and c's FLAI_IMAGE_*
// variables. It is run once and never again on output that names a network
// failure, as flai cannot tell what failed inside another's command: a
// non-zero exit is a *CommandError.
func RunCommand(ctx context.Context, r Runner, root string, argv []string, c CommandImages) error {
	if len(argv) == 0 || strings.TrimSpace(argv[0]) == "" {
		return fmt.Errorf("run the operator's build command for the agent images of %s: it names no program; name one with flai serve agent image command -- <program> [args...]", c.Inner)
	}
	words := CommandWords(argv)
	if r.Out != nil {
		fmt.Fprintf(r.Out, "agent image: running the operator's build command %s in %s for the images of %s, in place of flai's pull and build, with %s\n", words, root, c.Inner, envWords(c.Env()))
	}
	r.Program = argv[0]
	r.Delays = []time.Duration{} // never tried again
	err := r.Run(ctx, Step{Name: "the build command " + words, Dir: root, Env: c.Env(), Args: argv[1:]})
	var step *StepError
	if errors.As(err, &step) {
		return &CommandError{Command: argv, Inner: c.Inner, Said: step.Said, Err: step.Err}
	}
	if err != nil {
		return fmt.Errorf("make the agent images of %s with the operator's build command: %w", c.Inner, err)
	}
	return nil
}

// CommandError is the operator's build command that exited non-zero: its
// words, the inner harness it ran for, the last line of its output that says
// why (its error's when it wrote none), and its error.
type CommandError struct {
	Command []string
	Inner   string
	Said    string
	Err     error
}

// Error says which command failed, why, and what to do.
func (e *CommandError) Error() string {
	why := e.Said
	if exit := firstLine(e.Err.Error()); exit != why {
		why += "; " + exit
	}
	return fmt.Sprintf("the operator's build command %s failed for the agent images of %s (%s): its output is above and in the log; mend it, or clear it with flai serve agent image command --clear for flai's own pull and build, then run it again", CommandWords(e.Command), e.Inner, why)
}

// Unwrap is the command's own error.
func (e *CommandError) Unwrap() error { return e.Err }

// CommandWords is argv as a line says it: each word as it is, or quoted as
// quoted quotes it.
func CommandWords(argv []string) string {
	words := make([]string, 0, len(argv))
	for _, w := range argv {
		words = append(words, quoted(w))
	}
	return strings.Join(words, " ")
}

// envWords is env, NAME=value pairs, as a line says them, each value quoted
// as CommandWords quotes a word.
func envWords(env []string) string {
	words := make([]string, 0, len(env))
	for _, kv := range env {
		name, value, _ := strings.Cut(kv, "=")
		words = append(words, name+"="+quoted(value))
	}
	return strings.Join(words, " ")
}

// quoted is w, or w quoted as Go quotes a string when it is empty or holds
// a space, a quote, a backslash, or a character Go escapes.
func quoted(w string) string {
	if w == "" || strings.ContainsAny(w, " '") || strconv.Quote(w) != `"`+w+`"` {
		return strconv.Quote(w)
	}
	return w
}
