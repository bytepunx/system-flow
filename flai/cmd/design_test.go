package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/designsystem"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
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
		len(got.Problems) != 0 || len(got.Themes) != 2 || len(got.Contrast) != 2 || !strings.HasPrefix(got.CSS, ":root {") || !strings.HasPrefix(got.Tailwind, "@custom-variant dark ") {
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
	if code != 0 || out != designsystem.Tailwind(sys) || !strings.HasPrefix(out, "@custom-variant dark ") || !strings.Contains(out, "\n@theme {\n") {
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

// designStylesheetFile is the design stylesheet the tests name in the
// manifest's design.stylesheet: under folders that are not there yet.
const designStylesheetFile = "web/src/theme.css"

// nameDesignStylesheet sets the manifest's design.stylesheet to
// designStylesheetFile.
func nameDesignStylesheet(t *testing.T, root string) {
	t.Helper()
	data := readProjectFile(t, root, manifest.File) + "design:\n  stylesheet: " + designStylesheetFile + "\n"
	writeProjectFile(t, root, manifest.File, data)
}

// S-0385: css --write writes designsystem.Stylesheet to the file
// design.stylesheet names, creating its folders, prints the path, and
// replaces a stylesheet that has fallen behind; --json says what was written
// from what; --write and --tailwind are mutually exclusive.
func TestDesignCSSWrite(t *testing.T) {
	root := designProject(t)
	if _, errOut, code := runIn(t, root, "design", "init"); code != 0 {
		t.Fatal(errOut)
	}
	nameDesignStylesheet(t, root)
	sys, err := designsystem.Load(os.DirFS(root), designTokensFile)
	if err != nil {
		t.Fatal(err)
	}
	want := designsystem.Stylesheet(sys, designTokensFile)
	out, errOut, code := runIn(t, root, "design", "css", "--write")
	if code != 0 || out != "wrote "+designStylesheetFile+"\n" {
		t.Fatalf("css --write: %d %q %s", code, out, errOut)
	}
	got := readProjectFile(t, root, designStylesheetFile)
	if got != want || !strings.HasPrefix(got, "/* Generated by flai design css --write from "+designTokensFile+".") || !strings.Contains(got, "\n@theme {\n") {
		t.Errorf("the stylesheet written:\n%s", got)
	}

	writeProjectFile(t, root, designStylesheetFile, "@theme {\n  --color-primary: red;\n}\n")
	out, errOut, code = runIn(t, root, "design", "css", "--write", "--json")
	var written struct {
		Stylesheet string `json:"stylesheet"`
		Tokens     string `json:"tokens"`
	}
	if code != 0 || json.Unmarshal([]byte(out), &written) != nil || written.Stylesheet != designStylesheetFile || written.Tokens != designTokensFile {
		t.Errorf("css --write --json: %d %q %s", code, out, errOut)
	}
	if got := readProjectFile(t, root, designStylesheetFile); got != want {
		t.Errorf("a stale stylesheet was not replaced:\n%s", got)
	}

	writeProjectFile(t, root, designStylesheetFile, "mine\n")
	if out, errOut, code := runIn(t, root, "design", "css", "--write", "--tailwind"); code == 0 || out != "" || !strings.Contains(errOut, "tailwind") {
		t.Errorf("css --write --tailwind: %d %q %s", code, out, errOut)
	}
	if got := readProjectFile(t, root, designStylesheetFile); got != "mine\n" {
		t.Errorf("css --write --tailwind wrote the stylesheet:\n%s", got)
	}
}

// S-0385: css --write refuses with exit status 1 and writes nothing when the
// manifest names no stylesheet, when there is no token file, and when the
// token file has problems.
func TestDesignCSSWriteRefuses(t *testing.T) {
	t.Run("no stylesheet named", func(t *testing.T) {
		root := designProject(t)
		if _, errOut, code := runIn(t, root, "design", "init"); code != 0 {
			t.Fatal(errOut)
		}
		out, errOut, code := runIn(t, root, "design", "css", "--write")
		if code != 1 || out != "" || !strings.Contains(errOut, "names no design stylesheet") || !strings.Contains(errOut, "design.stylesheet in "+manifest.File) {
			t.Errorf("css --write: %d %q %s", code, out, errOut)
		}
		if projectFileExists(root, designStylesheetFile) || projectFileExists(root, "web") {
			t.Error("css --write wrote something")
		}
	})
	t.Run("no token file", func(t *testing.T) {
		root := designProject(t)
		nameDesignStylesheet(t, root)
		out, errOut, code := runIn(t, root, "design", "css", "--write")
		if code != 1 || out != "" || !strings.Contains(errOut, "no token file at "+designTokensFile) {
			t.Errorf("css --write: %d %q %s", code, out, errOut)
		}
		if projectFileExists(root, "web") {
			t.Error("css --write wrote something")
		}
	})
	t.Run("a token file with problems", func(t *testing.T) {
		root := designProject(t)
		nameDesignStylesheet(t, root)
		writeProjectFile(t, root, designTokensFile, "{ \"color\": {\n")
		out, errOut, code := runIn(t, root, "design", "css", "--write")
		if code != 1 || out != "" || !strings.Contains(errOut, "renders nothing from the token file "+designTokensFile+" while it has problems") {
			t.Errorf("css --write: %d %q %s", code, out, errOut)
		}
		if projectFileExists(root, "web") {
			t.Error("css --write wrote something")
		}
	})
}

// S-0385: init writes the design stylesheet from the starter tokens when the
// manifest names one, as a third path written and the stylesheet key in
// --json, and no stylesheet when it names none.
func TestDesignInitWritesTheStylesheetWhenNamed(t *testing.T) {
	root := designProject(t)
	nameDesignStylesheet(t, root)
	out, errOut, code := runIn(t, root, "design", "init")
	if want := "wrote " + designTokensFile + "\nwrote " + designGuidanceFile + "\nwrote " + designStylesheetFile + "\n"; code != 0 || out != want {
		t.Fatalf("init: %d %q %s", code, out, errOut)
	}
	sys, err := designsystem.Load(fstest.MapFS{"tokens.json": {Data: designsystem.StarterTokens}}, "tokens.json")
	if err != nil {
		t.Fatal(err)
	}
	if got := readProjectFile(t, root, designStylesheetFile); got != designsystem.Stylesheet(sys, designTokensFile) {
		t.Errorf("the stylesheet init wrote is not the starter tokens' stylesheet:\n%s", got)
	}

	root = designProject(t)
	nameDesignStylesheet(t, root)
	out, errOut, code = runIn(t, root, "design", "init", "--json")
	var got map[string]any
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil || got["stylesheet"] != designStylesheetFile {
		t.Errorf("init --json: %d %q %s", code, out, errOut)
	}

	root = designProject(t)
	out, errOut, code = runIn(t, root, "design", "init", "--json")
	got = nil
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("init --json: %d %q %s", code, out, errOut)
	}
	if _, ok := got["stylesheet"]; ok || projectFileExists(root, "web") {
		t.Errorf("init with no stylesheet named wrote one: %q", out)
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

// designStory is a story a scratch project holds for flai design show
// --story and flai design proposals: its title, its state, and whether it
// has a worktree.
type designStory struct {
	title, status string
	worktree      bool
}

// plainDesignTokens is a token file of one colour, with no theme modifier.
const plainDesignTokens = `{"color":{"primary":{"$type":"color","$value":{"colorSpace":"srgb","components":[0,0,1],"hex":"#0000ff"}}}}`

// designStories adds E-0001 and the stories given, S-0001 onwards, to the
// scratch project at root, each story with a worktree linked to the project
// as git links one, holding a copy of the project's manifest.
func designStories(t *testing.T, root string, stories ...designStory) *workitem.Repo {
	t.Helper()
	repo, err := workitem.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 14, 9, 0, 0, 0, time.UTC)
	if _, err := repo.Create(workitem.NewOptions{Type: workitem.Epic, Title: "Epic", Now: now}); err != nil {
		t.Fatal(err)
	}
	for _, s := range stories {
		it, err := repo.Create(workitem.NewOptions{Type: workitem.Story, Title: s.title, Parent: "E-0001", Now: now})
		if err != nil {
			t.Fatal(err)
		}
		it.Status = s.status
		if err := repo.Save(it); err != nil {
			t.Fatal(err)
		}
		if s.worktree {
			wt := repo.WorktreePath(it.ID)
			writeProjectFile(t, wt, ".git", "gitdir: "+filepath.Join(root, ".git", "worktrees", it.ID)+"\n")
			writeProjectFile(t, wt, manifest.File, readProjectFile(t, root, manifest.File))
		}
	}
	return repo
}

// S-0384: show --story refuses, with exit status 1 and the reason, an ID no
// story has, a done story, and a story with no worktree.
func TestDesignShowStoryRefuses(t *testing.T) {
	root := designProject(t)
	repo := designStories(t, root,
		designStory{"Open", workitem.InProgress, true},
		designStory{"Done", workitem.Done, true},
		designStory{"No worktree", workitem.Ready, false},
	)
	cases := []struct {
		name, story string
		want        []string
	}{
		{"an ID no item has", "S-0099", []string{"read the design system of S-0099", "name a story of this project"}},
		{"a done story", "S-0002", []string{"S-0002: it is done", "name an open story"}},
		{"a story without a worktree", "S-0003", []string{"S-0003: it has no worktree at " + repo.WorktreePath("S-0003"), "flai stream open S-0003"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			for _, args := range [][]string{{"design", "show", "--story", c.story}, {"design", "show", "--story", c.story, "--json"}} {
				out, errOut, code := runIn(t, root, args...)
				if code != 1 || out != "" {
					t.Errorf("%v: %d %q %s", args, code, out, errOut)
				}
				for _, w := range c.want {
					if !strings.Contains(errOut, w) {
						t.Errorf("%v: the error %q does not say %q", args, errOut, w)
					}
				}
			}
		})
	}
}

// S-0384: show --story reads the design system from the story's worktree,
// not the main checkout, names the story and its worktree, and reads the
// same from inside the worktree; --json is the worktree's design system.
func TestDesignShowStoryReadsItsWorktree(t *testing.T) {
	root := designProject(t)
	repo := designStories(t, root, designStory{"Proposal", workitem.InProgress, true})
	wt := repo.WorktreePath("S-0001")
	writeProjectFile(t, wt, designTokensFile, plainDesignTokens)
	writeProjectFile(t, wt, designGuidanceFile, "# Design system\n\n## Colour\n")
	writeProjectFile(t, root, designTokensFile, "{")

	for _, dir := range []string{root, wt} {
		out, errOut, code := runIn(t, dir, "design", "show", "--story", "S-0001")
		if code != 0 {
			t.Fatalf("show --story in %s: %d %s", dir, code, errOut)
		}
		for _, want := range []string{
			"story       S-0001, read from its worktree " + wt + "\n",
			"token file  " + designTokensFile + ", a plain token file\n",
			"guidance    " + designGuidanceFile + ", 2 headings\n",
			"tokens      light 1, dark 1\n",
		} {
			if !strings.Contains(out, want) {
				t.Errorf("show --story in %s does not say %q:\n%s", dir, want, out)
			}
		}
	}

	out, errOut, code := runIn(t, root, "design", "show", "--story", "S-0001", "--json")
	var got designsystem.Shown
	if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
		t.Fatalf("show --story --json: %d %q %s", code, out, errOut)
	}
	if !got.Exists || got.Resolver || got.Guidance != designGuidanceFile || len(got.Themes) != 2 || len(got.Themes[0].Tokens) != 1 || got.Themes[0].Tokens[0].Path != "color.primary" {
		t.Errorf("show --story --json is not the worktree's design system: %+v", got)
	}

	if err := os.Remove(filepath.Join(wt, filepath.FromSlash(designTokensFile))); err != nil {
		t.Fatal(err)
	}
	out, errOut, code = runIn(t, root, "design", "show", "--story", "S-0001")
	if want := "story S-0001, read from its worktree " + wt + "\nno token file at " + designTokensFile + ": run flai design init to write one and its guidance\n"; code != 0 || out != want {
		t.Errorf("show --story of a worktree with no token file: %d %q %s", code, out, errOut)
	}
}

// S-0384: with no open story's branch changing the design system, proposals
// prints [] with --json, and says so naming the token file and the guidance;
// a project that is not a git repository has no branches, so no proposals.
func TestDesignProposalsNone(t *testing.T) {
	root := designProject(t)
	designStories(t, root, designStory{"Open", workitem.InProgress, true})
	out, errOut, code := runIn(t, root, "design", "proposals", "--json")
	if code != 0 || out != "[]\n" {
		t.Errorf("proposals --json: %d %q %s", code, out, errOut)
	}
	out, errOut, code = runIn(t, root, "design", "proposals")
	if want := "no open story's branch changes the design system, the token file " + designTokensFile + " or the guidance " + designGuidanceFile + "\n"; code != 0 || out != want {
		t.Errorf("proposals: %d %q %s", code, out, errOut)
	}
}

// S-0384: proposals lists the open stories whose branch, or whose worktree
// uncommitted, changes the token file or the guidance, with their IDs,
// titles, states, and paths, the same from the main checkout and from a
// worktree; show --story reads what a proposal commits. Integration: it
// runs real git.
func TestDesignProposalsListsTheStoriesChangingTheDesignSystem(t *testing.T) {
	root := syncProject(t)
	// git status lists a new file in a folder git does not track as the folder
	writeProjectFile(t, root, "design/system/.gitkeep", "")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "design/system")
	one := openSyncStory(t, root, 1, "design")
	two := openSyncStory(t, root, 2, "design")
	three := openSyncStory(t, root, 3, "docs")
	commitIn(t, one, designTokensFile, plainDesignTokens)
	writeProjectFile(t, two, designGuidanceFile, "# Design system\n")
	commitIn(t, three, "docs/guide.md", "three's line\n")

	want := []designsystem.Proposal{
		{ID: "S-0001", Title: "Story 1", Status: workitem.InProgress, Paths: []string{designTokensFile}},
		{ID: "S-0002", Title: "Story 2", Status: workitem.InProgress, Paths: []string{designGuidanceFile}},
	}
	for _, dir := range []string{root, three} {
		out, errOut, code := runIn(t, dir, "design", "proposals", "--json")
		var got []designsystem.Proposal
		if code != 0 || json.Unmarshal([]byte(out), &got) != nil {
			t.Fatalf("proposals --json in %s: %d %q %s", dir, code, out, errOut)
		}
		if len(got) != len(want) {
			t.Fatalf("proposals --json in %s:\n got %+v\nwant %+v", dir, got, want)
		}
		for i := range want {
			if got[i].ID != want[i].ID || got[i].Title != want[i].Title || got[i].Status != want[i].Status || strings.Join(got[i].Paths, ",") != strings.Join(want[i].Paths, ",") {
				t.Errorf("proposal %d in %s:\n got %+v\nwant %+v", i, dir, got[i], want[i])
			}
		}
	}

	out, errOut, code := runIn(t, root, "design", "proposals")
	if code != 0 {
		t.Fatalf("proposals: %d %s", code, errOut)
	}
	for _, line := range []string{
		"ID      STATE        TITLE    CHANGES\n",
		"S-0001  in-progress  Story 1  " + designTokensFile + "\n",
		"S-0002  in-progress  Story 2  " + designGuidanceFile + "\n",
	} {
		if !strings.Contains(out, line) {
			t.Errorf("proposals does not say %q:\n%s", line, out)
		}
	}

	out, errOut, code = runIn(t, root, "design", "show", "--story", "S-0001")
	if code != 0 || !strings.Contains(out, "token file  "+designTokensFile+", a plain token file\n") {
		t.Errorf("show --story S-0001: %d %s\n%s", code, errOut, out)
	}
}
