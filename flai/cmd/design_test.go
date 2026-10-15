package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
)

// The paths flai design writes in a project with the default layout.
const (
	designTokensFile   = "design/system/design-system.tokens.json"
	designGuidanceFile = "design/system/design-system.md"
)

// designProject is a scratch project with no design system.
func designProject(t *testing.T) string {
	t.Helper()
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	return tempProject(t)
}

// writeProjectFile writes data at the repository path name under root.
func writeProjectFile(t *testing.T, root, name, data string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readProjectFile(t *testing.T, root, name string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func projectFileExists(root, name string) bool {
	_, err := os.Stat(filepath.Join(root, filepath.FromSlash(name)))
	return err == nil
}

// S-0381: flai design init writes the starter token file and the guidance,
// its updated today, its topics all, and the token file's path in its body;
// flai check finds nothing in either.
func TestDesignInitWritesTheTokenFileAndTheGuidance(t *testing.T) {
	root := designProject(t)
	out, errOut, code := runIn(t, root, "design", "init")
	if code != 0 {
		t.Fatalf("init: %d %s %s", code, out, errOut)
	}
	if want := "wrote " + designTokensFile + "\nwrote " + designGuidanceFile + "\n"; out != want {
		t.Errorf("output %q, want %q", out, want)
	}
	if got := readProjectFile(t, root, designTokensFile); got != string(designsystem.StarterTokens) {
		t.Error("the token file is not the starter")
	}
	guidance := readProjectFile(t, root, designGuidanceFile)
	for _, want := range []string{"---\ntitle: Design system\nupdated: 2026-09-15\nstatus: active\ntopics: [all]\n---\n\n# Design system\n", "`" + designTokensFile + "`", "\n## Accessibility\n"} {
		if !strings.Contains(guidance, want) {
			t.Errorf("the guidance does not hold %q:\n%s", want, guidance)
		}
	}
	if strings.Contains(guidance, designsystem.StarterTokensMark) {
		t.Errorf("the token file's mark is left in the guidance:\n%s", guidance)
	}
	// the bare test project has findings of its own; none may be about the design system
	if out, _, _ := runIn(t, root, "check", "--strict"); strings.Contains(out, "design-system") || strings.Contains(out, "design.tokens") {
		t.Errorf("flai check has something to say about the design system:\n%s", out)
	}
}

// S-0381: --topics sets the guidance's topics, --json says what was written,
// and a word that is no topic is refused before anything is written.
func TestDesignInitTopics(t *testing.T) {
	root := designProject(t)
	out, errOut, code := runIn(t, root, "design", "init", "--topics", "flaiover,ui", "--json")
	var got struct {
		Tokens   string   `json:"tokens"`
		Guidance string   `json:"guidance"`
		Topics   []string `json:"topics"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("init --json: %d %q %s", code, out, errOut)
	}
	if got.Tokens != designTokensFile || got.Guidance != designGuidanceFile || strings.Join(got.Topics, ",") != "flaiover,ui" {
		t.Errorf("init --json: %+v", got)
	}
	if g := readProjectFile(t, root, designGuidanceFile); !strings.Contains(g, "\ntopics: [flaiover, ui]\n") {
		t.Errorf("the guidance's topics:\n%s", g)
	}

	root = designProject(t)
	if _, errOut, code := runIn(t, root, "design", "init", "--topics", "ui,not a topic"); code != 1 || !strings.Contains(errOut, "is not a topic") {
		t.Errorf("a word that is no topic: %d %s", code, errOut)
	}
	if projectFileExists(root, designTokensFile) || projectFileExists(root, designGuidanceFile) {
		t.Error("a refused init wrote a file")
	}
}

// S-0381: init refuses when either file is there, naming it, and writes
// neither.
func TestDesignInitRefusesWhenEitherFileIsThere(t *testing.T) {
	for _, there := range []string{designTokensFile, designGuidanceFile} {
		t.Run(filepath.Base(there), func(t *testing.T) {
			root := designProject(t)
			writeProjectFile(t, root, there, "mine\n")
			out, errOut, code := runIn(t, root, "design", "init")
			if code != 1 || !strings.Contains(errOut, there+" is there already") || out != "" {
				t.Errorf("init over %s: %d %q %s", there, code, out, errOut)
			}
			if got := readProjectFile(t, root, there); got != "mine\n" {
				t.Errorf("%s was changed: %q", there, got)
			}
			for _, other := range []string{designTokensFile, designGuidanceFile} {
				if other != there && projectFileExists(root, other) {
					t.Errorf("%s was written", other)
				}
			}
		})
	}
}

// S-0381: init writes the token file where the manifest's design.tokens
// says, and the guidance names that path.
func TestDesignInitFollowsDesignTokens(t *testing.T) {
	root := designProject(t)
	data := readProjectFile(t, root, manifest.File) + "design:\n  tokens: ui/brand.tokens.json\n"
	writeProjectFile(t, root, manifest.File, data)
	out, errOut, code := runIn(t, root, "design", "init")
	if code != 0 || !strings.Contains(out, "wrote ui/brand.tokens.json\n") {
		t.Fatalf("init: %d %q %s", code, out, errOut)
	}
	if g := readProjectFile(t, root, designGuidanceFile); !strings.Contains(g, "`ui/brand.tokens.json`") {
		t.Errorf("the guidance does not name the token file:\n%s", g)
	}
}

// S-0381: show without a token file says so and how to write one, and is no
// error, in text and in JSON.
func TestDesignShowWithoutATokenFile(t *testing.T) {
	root := designProject(t)
	out, errOut, code := runIn(t, root, "design", "show")
	if code != 0 || out != "no token file at "+designTokensFile+": run flai design init to write one and its guidance\n" {
		t.Errorf("show: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "design", "show", "--json")
	var got designsystem.Shown
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("show --json: %d %q %s", code, out, errOut)
	}
	if got.Exists || got.Path != designTokensFile || !strings.Contains(got.NoTokenFile, "flai design init") {
		t.Errorf("show --json: %+v", got)
	}
}

// S-0381: show after init summarises the token file, the guidance, the
// tokens of each theme, and each role pair's contrast; --json is
// designsystem.Show's answer.
func TestDesignShowAfterInit(t *testing.T) {
	root := designProject(t)
	if _, errOut, code := runIn(t, root, "design", "init"); code != 0 {
		t.Fatal(errOut)
	}
	out, errOut, code := runIn(t, root, "design", "show")
	if code != 0 {
		t.Fatalf("show: %d %s", code, errOut)
	}
	for _, want := range []string{
		"token file  " + designTokensFile + ", a resolver document\n",
		"guidance    " + designGuidanceFile + ", 9 headings\n",
		"tokens      light 35, dark 35\n",
		"problems    none\n",
		"\ncontrast in light\n  PAIR ",
		"\ncontrast in dark\n",
		"  color.text on color.background ",
		"  color.border on color.background ",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("show does not say %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "FAIL") || strings.Contains(out, " yes\n") {
		t.Errorf("the starter has a pair that fails or is clipped:\n%s", out)
	}

	out, errOut, code = runIn(t, root, "design", "show", "--json")
	var got designsystem.Shown
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("show --json: %d %q %s", code, out, errOut)
	}
	if !got.Exists || !got.Resolver || got.Guidance != designGuidanceFile || len(got.Headings) != 9 || got.Headings[0] != "Design system" ||
		len(got.Problems) != 0 || len(got.Themes) != 2 || len(got.Contrast) != 2 || !strings.HasPrefix(got.CSS, ":root {") || !strings.HasPrefix(got.Tailwind, "@theme {") {
		t.Errorf("show --json: exists %v resolver %v guidance %q headings %q problems %v themes %d contrast %d",
			got.Exists, got.Resolver, got.Guidance, got.Headings, got.Problems, len(got.Themes), len(got.Contrast))
	}
}

// S-0381: show prints each problem as <file>:<line>: <token path>: <message>
// and exits 0, and says when the guidance is not there.
func TestDesignShowPrintsProblems(t *testing.T) {
	root := designProject(t)
	writeProjectFile(t, root, designTokensFile, "{\n  \"color\": {\n    \"background\": { \"$type\": \"color\", \"$value\": { \"colorSpace\": \"hsl\", \"components\": [0, 0, 100] } }\n  }\n}\n")
	out, errOut, code := runIn(t, root, "design", "show")
	if code != 0 {
		t.Fatalf("show: %d %s", code, errOut)
	}
	for _, want := range []string{
		"guidance    " + designGuidanceFile + ", not there\n",
		"\n" + designTokensFile + ":3: color.background: flai does not read the colour space \"hsl\"",
		"\n" + designTokensFile + ":",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("show does not say %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "problems    none") {
		t.Errorf("show says there are no problems:\n%s", out)
	}
}

// S-0381: css prints the custom properties, one block per theme, or with
// --tailwind the @theme block; no token file, or one with problems, is an
// error and prints nothing.
func TestDesignCSS(t *testing.T) {
	root := designProject(t)
	out, errOut, code := runIn(t, root, "design", "css")
	if code != 1 || out != "" || !strings.Contains(errOut, "no token file at "+designTokensFile) || !strings.Contains(errOut, "flai design init") {
		t.Errorf("css without a token file: %d %q %s", code, out, errOut)
	}
	if _, errOut, code := runIn(t, root, "design", "init"); code != 0 {
		t.Fatal(errOut)
	}
	sys, err := designsystem.Load(os.DirFS(root), designTokensFile)
	if err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runIn(t, root, "design", "css")
	if code != 0 || out != designsystem.CSS(sys) || !strings.Contains(out, "\n[data-theme='dark'] {\n") || !strings.Contains(out, "  --color-primary: #2563eb;\n") {
		t.Errorf("css: %d %s\n%s", code, errOut, out)
	}
	out, errOut, code = runIn(t, root, "design", "css", "--tailwind")
	if code != 0 || out != designsystem.Tailwind(sys) || !strings.HasPrefix(out, "@theme {\n") {
		t.Errorf("css --tailwind: %d %s\n%s", code, errOut, out)
	}

	writeProjectFile(t, root, designTokensFile, "{ \"color\": {\n")
	for _, args := range [][]string{{"design", "css"}, {"design", "css", "--tailwind"}} {
		out, errOut, code := runIn(t, root, args...)
		if code != 1 || out != "" || !strings.Contains(errOut, "renders nothing from the token file "+designTokensFile+" while it has problems") || !strings.Contains(errOut, designTokensFile+":") {
			t.Errorf("%v of a file that does not parse: %d %q %s", args, code, out, errOut)
		}
	}
}

// S-0381: the starter token file has no problem, holds every required
// token in both themes, and every role pair passes its contrast threshold
// in light and dark with no colour clipped.
func TestDesignStarterTokensPassContrast(t *testing.T) {
	fsys := fstest.MapFS{"tokens.json": {Data: designsystem.StarterTokens}}
	sys, err := designsystem.Load(fsys, "tokens.json")
	if err != nil {
		t.Fatal(err)
	}
	if !sys.Resolver || len(sys.Problems) != 0 || len(sys.Themes) != 2 {
		t.Fatalf("resolver %v, problems %v, themes %d", sys.Resolver, sys.Problems, len(sys.Themes))
	}
	for _, th := range sys.Themes {
		for _, p := range designsystem.RequiredPaths() {
			if _, ok := th.Tokens.Get(p); !ok {
				t.Errorf("%s: no %s", th.Name, p)
			}
		}
		pairs := designsystem.Contrast(th.Name, th.Tokens)
		if len(pairs) != 20 {
			t.Errorf("%s: %d pairs, want 20", th.Name, len(pairs))
		}
		for _, p := range pairs {
			if !p.Passes || p.Clipped {
				t.Errorf("%s: %s on %s is %.2f against %.1f, clipped %v", th.Name, p.Foreground, p.Background, p.Ratio, p.Threshold, p.Clipped)
			}
		}
	}
	light, dark := sys.Themes[0].Tokens, sys.Themes[1].Tokens
	for _, role := range designsystem.ColorRoles {
		l, _ := light.Get("color." + role)
		d, _ := dark.Get("color." + role)
		if bytes.Equal(mustJSON(t, l.Value.JSON()), mustJSON(t, d.Value.JSON())) {
			t.Errorf("color.%s has no dark value of its own", role)
		}
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// S-0381: in a project made from the template, the token file and the
// guidance init writes pass flai check --strict.
func TestDesignInitInTheTemplatePassesCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("integration: renders the monorepo template")
	}
	proto := filepath.Join("..", "..", "template")
	if _, err := os.Stat(filepath.Join(proto, "template.yaml")); err != nil {
		t.Skip("prototype template not present")
	}
	t.Setenv("FLAI_CONFIG", filepath.Join(t.TempDir(), "cfg.json"))
	dest := filepath.Join(t.TempDir(), "sample")
	if _, errOut, code := runCLI(t, "new", dest, "--template", proto, "--defaults", "--no-git"); code != 0 {
		t.Fatalf("render the template: %s", errOut)
	}
	if out, errOut, code := runIn(t, dest, "design", "init"); code != 0 {
		t.Fatalf("init: %d %s %s", code, out, errOut)
	}
	if out, errOut, code := runIn(t, dest, "check", "--strict"); code != 0 {
		t.Errorf("check --strict: %d\n%s\n%s", code, out, errOut)
	}
}
