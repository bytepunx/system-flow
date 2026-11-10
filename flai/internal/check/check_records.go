package check

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/analysis"
	"github.com/bytepunx/system-flow/flai/internal/components"
	"github.com/bytepunx/system-flow/flai/internal/conflictmark"
	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/experiment"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/mdlint"
	"github.com/bytepunx/system-flow/flai/internal/messages"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// threads validates wip/threads (ADR-0020): schema, file names, anchors
// that exist, headings that are present, and open threads on archived items.
func (c *checker) threads() {
	dir := threads.Dir(c.repo)
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return // optional until the first thread
	}
	list, err := threads.List(c.repo)
	if err != nil {
		c.add(Error, "threads.front-matter", dir, 1, "%v", err)
		return
	}
	byID := map[string]*workitem.Item{}
	for _, it := range c.items {
		byID[it.ID] = it
	}
	seen := map[string]string{}
	for _, th := range list {
		if err := th.Validate(); err != nil {
			c.add(Error, "threads.front-matter", th.Path, keyLine(th.Path, "id"), "%v", err)
			continue
		}
		if prev, dup := seen[th.ID]; dup {
			c.add(Error, "threads.duplicate-id", th.Path, keyLine(th.Path, "id"), "%s is also defined in %s", th.ID, prev)
		}
		seen[th.ID] = th.Path
		c.unknownFields("threads.unknown-field", th.Path, th.Unknown)
		if base := filepath.Base(th.Path); !strings.HasPrefix(base, th.ID+"-") {
			c.add(Error, "threads.filename", th.Path, 1, "file name should start with %s-", th.ID)
		}
		// An item anchor follows the item: archiving moves the file, the
		// thread stays attached. A plain path anchor must exist as written.
		anchorPath := th.Anchor.Path
		if th.Anchor.Item != "" {
			it, ok := byID[th.Anchor.Item]
			switch {
			case !ok:
				c.add(Error, "threads.anchor", th.Path, keyLine(th.Path, "anchor"), "item %s does not exist", th.Anchor.Item)
				continue
			case it.Archived && th.Open():
				c.add(Warning, "threads.archived", th.Path, keyLine(th.Path, "status"), "%s is %s but %s is archived; resolve it with flai thread resolve %s", th.ID, th.Status, th.Anchor.Item, th.ID)
			}
			anchorPath = it.Path
		}
		abs := anchorPath
		if !filepath.IsAbs(abs) {
			abs = filepath.Join(c.repo.Root, filepath.FromSlash(anchorPath))
		}
		data, err := os.ReadFile(abs)
		if err != nil && !filepath.IsAbs(anchorPath) {
			if d2, err2 := os.ReadFile(filepath.Join(c.repo.MainRoot, filepath.FromSlash(anchorPath))); err2 == nil {
				data, err = d2, nil
			}
		}
		if err != nil {
			c.add(Error, "threads.anchor", th.Path, keyLine(th.Path, "anchor"), "anchor %s does not exist", th.Anchor.Path)
		} else if th.Anchor.Heading != "" && !threads.HasHeading(string(data), th.Anchor.Heading) {
			c.add(Warning, "threads.heading", th.Path, keyLine(th.Path, "anchor"), "heading %q is no longer in %s", th.Anchor.Heading, th.Anchor.Path)
		}
		if len(th.Entries()) == 0 {
			c.add(Warning, "threads.entries", th.Path, 1, "%s has no dated entries", th.ID)
		}
	}
}

// messages validates wip/messages (ADR-0120): each conversation's front
// matter, unique IDs, file names, that both its stories exist, its dated
// entries, and an open conversation of an archived story. Its markdown is
// the wip walk's (markdown). Every markdown file but README.md is read on its
// own, so that one that does not parse is reported and the rest still are.
func (c *checker) messages() {
	dir := messages.Dir(c.repo)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // optional until the first message
	}
	seen := map[string]string{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || name == "README.md" {
			continue
		}
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			c.add(Error, "messages.front-matter", path, 1, "%v", err)
			continue
		}
		conv, err := messages.Parse(string(data))
		if err == nil {
			err = conv.Validate()
		}
		if err != nil {
			c.add(Error, "messages.front-matter", path, keyLine(path, "id"), "%v", err)
			continue
		}
		if prev, dup := seen[conv.ID]; dup {
			c.add(Error, "messages.duplicate-id", path, keyLine(path, "id"), "%s is also defined in %s", conv.ID, prev)
		}
		seen[conv.ID] = path
		if rel, err := filepath.Rel(c.repo.Root, path); err == nil {
			seen[conv.ID] = rel
		}
		c.unknownFields("messages.unknown-field", path, conv.Unknown)
		if !strings.HasPrefix(name, conv.ID+"-") {
			c.add(Error, "messages.filename", path, 1, "file name should start with %s-", conv.ID)
		}
		archived := ""
		for _, side := range []struct{ field, story string }{{"from", conv.From}, {"to", conv.To}} {
			it := c.named(side.story)
			switch {
			case it == nil:
				c.add(Error, "messages.story", path, keyLine(path, side.field), "%s names story %s, which does not exist", side.field, side.story)
			case it.Archived && archived == "":
				archived = it.ID
			}
		}
		// flai accept and flai archive close it; one an older flai left open
		// is closed by hand, as they would have closed it
		if archived != "" && conv.Status == messages.StatusOpen {
			c.add(Warning, "messages.closed", path, keyLine(path, "status"), "%s is open but %s is archived, so it reads as closed; close it as flai archive would: set its status to closed and add a dated entry \"Closed: %s was archived\" under ## Entries", conv.ID, archived, archived)
		}
		if len(conv.Entries()) == 0 {
			c.add(Warning, "messages.entries", path, 1, "%s has no dated entries", conv.ID)
		}
	}
}

func (c *checker) issues() {
	dir := issues.Dir(c.repo)
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return // optional until the first issue is recorded
	}
	list, err := issues.List(c.repo)
	if err != nil {
		c.add(Error, "issues.front-matter", dir, 1, "%v", err)
		return
	}
	seen := map[string]string{}
	for _, is := range list {
		if err := is.Validate(); err != nil {
			c.add(Error, "issues.front-matter", is.Path, keyLine(is.Path, "id"), "%v", err)
		}
		if prev, dup := seen[is.ID]; dup {
			c.add(Error, "issues.duplicate-id", is.Path, keyLine(is.Path, "id"), "%s is also defined in %s", is.ID, prev)
		}
		seen[is.ID] = is.Path
		c.unknownFields("issues.unknown-field", is.Path, is.Unknown)
		if base := filepath.Base(is.Path); !strings.HasPrefix(base, is.ID+"-") {
			c.add(Error, "issues.filename", is.Path, 1, "file name should start with %s-", is.ID)
		}
	}
	c.issuesNoStory(list)
	if len(list) == 0 {
		return // a project that has never had an issue keeps the template's summary
	}
	summaryPath := filepath.Join(dir, issues.SummaryFile)
	data, err := os.ReadFile(summaryPath)
	if err != nil {
		c.add(Warning, "issues.summary", dir, 1, "summary.md is missing; run flai check --fix or flai issue summary, which write it")
		return
	}
	if summaryUpdated.ReplaceAllString(string(data), "") != summaryUpdated.ReplaceAllString(issues.Summary(list, c.now), "") {
		c.add(Warning, "issues.summary", summaryPath, 1, "summary.md differs from what its issue files render; run flai check --fix or flai issue summary, which rewrite it")
	}
}

// summaryUpdated is summary.md's updated line, which says when it was
// written rather than what it holds, so the issues.summary comparison leaves
// it out.
var summaryUpdated = regexp.MustCompile(`(?m)^updated: .*$`)

// issuesNoStory warns about each open issue older than issues.story_after
// that no open story links (S-0198). An invalid story_after is layout's
// finding, and turns this rule off, as 0 does.
func (c *checker) issuesNoStory(list []*issues.Issue) {
	after, err := c.repo.Manifest.Issues.StoryAfterDuration()
	if err != nil || after == 0 {
		return
	}
	setting := "issues.story_after is " + strings.TrimSpace(c.repo.Manifest.Issues.StoryAfter)
	if strings.TrimSpace(c.repo.Manifest.Issues.StoryAfter) == "" {
		setting = "issues.story_after is unset, so " + strings.TrimSuffix(manifest.DefaultStoryAfter.String(), "0m0s")
	}
	for _, is := range issues.NoStory(list, c.items) {
		first, err := time.Parse(workitem.TimeFormat, is.FirstReported)
		if err != nil || c.now.Sub(first) <= after {
			continue // an unreadable first_reported is issues.front-matter's finding
		}
		c.add(Warning, "issues.no-story", is.Path, keyLine(is.Path, "first_reported"),
			"%s has been open since %s with no open story linking it (%s): make one with flai issue story %s",
			is.ID, first.Format(time.DateOnly), setting, is.ID)
	}
}

// experiments validates design/experiments (ADR-0066): each results document's
// front matter, its file name, its sections, and that it names an experiment
// story that exists.
func (c *checker) experiments() {
	dir := filepath.Join(c.repo.Root, filepath.FromSlash(experiment.Dir(c.repo.Manifest)))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // optional until the first experiment records its results
	}
	seen := map[string]string{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || experiment.Skipped(name) {
			continue
		}
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			c.add(Error, "experiments.front-matter", path, 1, "%v", err)
			continue
		}
		f, problems := experiment.Validate(name, string(data))
		for _, p := range problems {
			line := 1
			switch {
			case p.Field != "":
				line = keyLine(path, p.Field)
			case p.Heading != "":
				line = headingLine(path, p.Heading)
			}
			c.add(Error, "experiments.document", path, line, "%s", p.Message)
		}
		if f.Story == "" {
			continue
		}
		if prev, dup := seen[f.Story]; dup {
			c.add(Error, "experiments.duplicate", path, keyLine(path, "story"), "%s already has its results in %s", f.Story, prev)
		}
		seen[f.Story] = path
		switch it := c.named(f.Story); {
		case it == nil:
			c.add(Error, "experiments.story", path, keyLine(path, "story"), "story %s does not exist", f.Story)
		case it.Type != workitem.Story || it.Nature != experiment.Nature:
			c.add(Warning, "experiments.story", path, keyLine(path, "story"), "%s is a %s %s, not an experiment story; results of anything else belong elsewhere", f.Story, it.Nature, it.Type)
		}
	}
}

// analysis validates design/analysis (S-0223): each report's file name and
// front matter, and that the folder's README.md lists it, as adrIndex keeps
// the ADRs' index, with no row linking to a report that is not there.
func (c *checker) analysis() {
	dir := filepath.Join(c.repo.Root, filepath.FromSlash(analysis.Dir(c.repo.Manifest)))
	entries, err := os.ReadDir(dir)
	if err != nil {
		return // optional until the analyzer writes its first report
	}
	index := filepath.Join(dir, analysis.Index)
	listed, _ := os.ReadFile(index) // without an index, every report is unlisted, and said to be
	present := map[string]bool{}
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || analysis.Skipped(name) {
			continue
		}
		present[name] = true
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			c.add(Error, "analysis.front-matter", path, 1, "%v", err)
			continue
		}
		_, problems := analysis.Validate(name, string(data))
		for _, p := range problems {
			line := 1
			if p.Field != "" {
				line = keyLine(path, p.Field)
			}
			c.add(Error, "analysis.report", path, line, "%s", p.Message)
		}
		if !analysis.Lists(string(listed), name) {
			c.add(Warning, "analysis.index", path, 1, "%s/%s does not list it; add a row linking to %s", analysis.Dir(c.repo.Manifest), analysis.Index, name)
		}
	}
	for i, l := range strings.Split(string(listed), "\n") {
		for _, m := range analysisLink.FindAllStringSubmatch(l, -1) {
			if !present[m[1]] {
				c.add(Warning, "analysis.index", index, i+1, "the row links to %s, which is not there", m[1])
			}
		}
	}
}

var analysisLink = regexp.MustCompile(`\]\((?:\./)?(\d{4}-\d{2}-\d{2}-[a-z]+\.md)\)`)

// digest validates the design folder's digest (ADR-0191): each entry's name
// and front matter, what its sources and absorbed_by name, as digest.entry,
// and text with the shape of a secret, as digest.secret, each on its line.
// What each entry's file alone decides is kept between checks in the cache
// folder by the file's size and modification time (S-0677); what its sources
// and absorbed_by name is looked up every check.
func (c *checker) digest() {
	dir := filepath.Join(c.repo.Root, filepath.FromSlash(digest.Dir(c.repo.Manifest)))
	files, err := digest.CheckFolder(dir, c.repo.CacheDir(), digest.EnvFor(c.repo))
	if err != nil {
		return // optional until the first entry is written
	}
	for _, f := range files {
		if f.Err != nil {
			c.add(Error, digest.RuleEntry, f.Path, 1, "%v", f.Err)
			continue
		}
		for _, p := range f.Problems {
			c.add(Error, p.Rule, f.Path, p.Line, "%s", p.Message)
		}
	}
}

// LintFindings are a refused write's markdown lint findings as flai check
// reports them, for the commands that refuse with check findings.
func LintFindings(e *mdlint.Error) []Finding {
	out := make([]Finding, len(e.Findings))
	for i, f := range e.Findings {
		out[i] = Finding{Level: Warning, Rule: "markdown." + f.Rule, Path: e.Path, Line: f.Line, Message: f.String()}
	}
	return out
}

// markdown lints every markdown file under the wip folder with the project's
// markdownlint configuration (S-0179). flai writes them in the main
// checkout, where no story's lint runs (ADR-0019), so they reach main unlinted
// unless flai check says so first. Findings are warnings: --strict fails on
// them, and the commands that write refuse a change that brings one. The
// configuration is the main checkout's, or the worktree's when the main
// checkout's root holds none, as in an agent container that mounts only its
// .git, wip, and system-flow.yaml (ADR-0253); in the main checkout the two
// are the same. A project without a markdownlint configuration is not
// linted.
func (c *checker) markdown() {
	root := c.repo.MainRoot
	if root == "" {
		root = c.repo.Root
	}
	cfg, err := mdlint.Load(root)
	if err == nil && cfg == nil && root != c.repo.Root {
		root = c.repo.Root
		cfg, err = mdlint.Load(root)
	}
	if err != nil {
		c.add(Warning, "markdown.config", root, 1, "the markdownlint configuration cannot be read, so wip is not linted: %v", err)
		return
	}
	if cfg == nil {
		return
	}
	if c.write != nil {
		// a write's run lints only the files it wrote under wip (S-0598)
		wip := c.repo.WipDir()
		for _, path := range c.write.files {
			if !strings.HasSuffix(path, ".md") || !within(path, wip) {
				continue
			}
			data, err := os.ReadFile(path)
			if err != nil {
				continue // unreadable files are skipped, as the other rules skip what they cannot read
			}
			for _, f := range cfg.Lint(string(data)) {
				c.add(Warning, "markdown."+f.Rule, path, f.Line, "%s", f)
			}
		}
		return
	}
	// Each file's findings are kept in the cache folder between runs, so a
	// check lints only the files that changed (S-0600).
	for _, ff := range cfg.LintFolderKept(c.repo.WipDir(), c.repo.CacheDir()) {
		for _, f := range ff.Findings {
			c.add(Warning, "markdown."+f.Rule, ff.Path, f.Line, "%s", f)
		}
	}
}

// nodeRunner runs node for the components.examples check; tests replace it.
var nodeRunner components.Runner = components.System{}

// componentExamples reports components.examples (ADR-0164): each examples
// file under the design folder's components folder that does not fit the
// format, names a component that does not exist, or has examples that do not
// fit the component's props, read with node once for each sub-project with
// examples. Warnings, at the problem's line, or line 1 for a file without its
// component or whose component's props could not be read. A project without
// examples files is not checked. One with no sub-project whose components
// flai reads, or a host without node, gets one note, which --strict passes
// over, and so does a sub-project whose svelte is not installed, on its
// examples folder (ADR-0171); without node or svelte the format and the
// components are still checked. The notes for node and svelte are marked
// unchecked.
func (c *checker) componentExamples() {
	folder := filepath.Join(c.repo.Root, filepath.FromSlash(components.ExamplesFolder(c.repo.Manifest)))
	var files []string
	_ = filepath.WalkDir(folder, func(p string, e fs.DirEntry, err error) error {
		if err == nil && e.Type().IsRegular() && strings.HasSuffix(e.Name(), components.ExamplesSuffix) {
			files = append(files, p)
		}
		return nil //nolint:nilerr // an absent folder has no examples, and an unreadable entry is skipped
	})
	if len(files) == 0 {
		return
	}
	shown := components.ExamplesFolder(c.repo.Manifest)
	projects := map[string]manifest.Project{}
	for _, p := range c.repo.Manifest.Projects {
		if manifest.SupportsComponents(p.Kind) {
			projects[p.Name] = p
		}
	}
	if len(projects) == 0 {
		c.advise("components.examples", folder, 1, "the %d examples files under %s are not checked: system-flow.yaml names no sub-project of kind %s, whose components they are of; add the sub-project under projects with its kind, or remove the files", len(files), shown, strings.Join(manifest.ComponentKinds, " or "))
		return
	}
	_, nodeErr := nodeRunner.LookPath("node")
	// Each sub-project's components, found once, and the files of those with
	// examples, in the order they were found.
	found := map[string]map[string]components.Component{}
	var withExamples []string
	byProject := map[string][]exampleFile{}
	for _, path := range files {
		f, ok := c.exampleFile(path, folder, projects, found)
		if !ok {
			continue
		}
		if _, ok := byProject[f.project]; !ok {
			withExamples = append(withExamples, f.project)
		}
		byProject[f.project] = append(byProject[f.project], f)
	}
	if nodeErr != nil {
		c.adviseUnchecked(folder, "the props of the examples under %s are not checked against their components': flai reads a component's props with node, which is not on PATH; install Node.js (https://nodejs.org) to check them", shown)
		return
	}
	for _, name := range withExamples {
		c.exampleProps(projects[name], byProject[name], folder, shown)
	}
}

// exampleFile is one examples file whose format was checked and whose
// component exists.
type exampleFile struct {
	path     string
	project  string
	comp     components.Component
	examples []components.Example
}

// exampleFile checks the examples file at path, under folder: its format,
// then that its component is one of a sub-project of projects, found once
// into found. False when its component does not exist or its sub-project's
// components could not be found, each reported.
func (c *checker) exampleFile(path, folder string, projects map[string]manifest.Project, found map[string]map[string]components.Component) (exampleFile, bool) {
	const rule = "components.examples"
	data, err := os.ReadFile(path)
	if err != nil {
		c.add(Warning, rule, path, 1, "cannot read the examples file: %v; check that it is a readable file", err)
		return exampleFile{}, false
	}
	examples, problems := components.ParseExamples(data)
	for _, p := range problems {
		c.add(Warning, rule, path, p.Line, "%s", p)
	}
	rel, err := filepath.Rel(folder, path)
	if err != nil {
		return exampleFile{}, false
	}
	name, compPath, ok := components.ExamplesComponent(filepath.ToSlash(rel))
	if !ok {
		c.add(Warning, rule, path, 1, "the file is in no sub-project's folder, so it is no component's; move it to the folder of the sub-project's name, at the component's path with .svelte made %s", components.ExamplesSuffix)
		return exampleFile{}, false
	}
	p, ok := projects[name]
	if !ok {
		c.add(Warning, rule, path, 1, "system-flow.yaml names no sub-project %q of kind %s, so the file's component %s does not exist; move the file to the folder of its sub-project's name, or remove it", name, strings.Join(manifest.ComponentKinds, " or "), compPath)
		return exampleFile{}, false
	}
	comps, ok := found[name]
	if !ok {
		s, _ := p.ComponentSettings()
		list, err := components.Find(c.repo.Root, p.Path, s.Include, s.Exclude)
		if err != nil {
			c.add(Warning, rule, path, 1, "its examples are not checked: %v", err)
			return exampleFile{}, false
		}
		comps = make(map[string]components.Component, len(list))
		for _, comp := range list {
			comps[comp.Path] = comp
		}
		found[name] = comps
	}
	comp, ok := comps[compPath]
	if !ok {
		c.add(Warning, rule, path, 1, "sub-project %q has no component %s, so the file is no component's; move it with its component, or remove it", name, compPath)
		return exampleFile{}, false
	}
	return exampleFile{path: path, project: name, comp: comp, examples: examples}, true
}

// exampleProps reads the props of files' components, all of project, with
// node once, and reports each example that does not fit its component's;
// a file whose component's props could not be read is reported once, with
// the reader's error. When svelte is not installed in project, the reader
// stops for all its files, and the project's folder under folder, the
// examples folder, shown as shown, gets one note instead, marked unchecked
// (ADR-0171).
func (c *checker) exampleProps(project manifest.Project, files []exampleFile, folder, shown string) {
	const rule = "components.examples"
	comps := make([]components.Component, len(files))
	for i, f := range files {
		comps[i] = f.comp
	}
	read, err := components.ReadProps(context.Background(), nodeRunner, c.repo.Root, comps)
	if errors.Is(err, components.ErrSvelteMissing) {
		c.adviseUnchecked(filepath.Join(folder, project.Name), "the props of the examples under %s/%s are not checked against their components': %v", shown, project.Name, err)
		return
	}
	for i, f := range files {
		switch {
		case err != nil:
			c.add(Warning, rule, f.path, 1, "its examples are not checked, since the props of %s in %s could not be read: %v", f.comp.Path, project.Name, err)
		case read[i].Error != "":
			c.add(Warning, rule, f.path, 1, "its examples are not checked, since the props of %s in %s could not be read: %s", f.comp.Path, project.Name, read[i].Error)
		default:
			for _, p := range components.Validate(f.examples, read[i].Props) {
				c.add(Warning, rule, f.path, p.Line, "%s", p)
			}
		}
	}
}

// conflictMarkers reports each line that opens or closes a merge conflict in
// the markdown under the design, docs, and wip folders and at the root: a
// conflict git-added unresolved reached main in a design document and nothing
// caught it (S-0253, I-0066). Errors, so CI's check fails without --strict.
// A whole run keeps each file's marker lines between checks in the cache
// folder by the file's size and modification time, so an unchanged file is
// not read again (S-0677); a write's run reads the files it wrote.
func (c *checker) conflictMarkers() {
	seen := map[string]bool{}
	// lines reads a file's marker lines; a whole run takes them from what
	// the last run kept, by the file's size and modification time (S-0677).
	lines := func(path string) ([]int, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return conflictmark.Lines(string(data)), nil
	}
	scan := func(path string) {
		if seen[path] {
			return
		}
		seen[path] = true
		marks, err := lines(path)
		if err != nil {
			return // unreadable files are skipped, as the other rules skip them
		}
		for _, line := range marks {
			c.add(Error, "markdown.conflict-marker", path, line, "a merge conflict marker: resolve the conflict, keeping what both sides meant, and remove the markers")
		}
	}
	root := filepath.Clean(c.repo.Root)
	var dirs []string
	for _, dir := range []string{c.repo.Manifest.Dir(c.repo.Root, "design"), c.repo.Manifest.Dir(c.repo.Root, "docs"), c.repo.WipDir()} {
		if filepath.Clean(dir) == root {
			continue // a folder laid out at the root: its markdown is the root's, below
		}
		dirs = append(dirs, dir)
	}
	if c.write != nil {
		// a write's run scans only the files it wrote that the whole run scans
		for _, path := range c.write.files {
			if strings.HasSuffix(path, ".md") && (filepath.Dir(path) == root || slices.ContainsFunc(dirs, func(dir string) bool { return within(path, dir) })) {
				scan(path)
			}
		}
		return
	}
	kept := conflictmark.OpenKept(c.repo.CacheDir(), root)
	defer kept.Save()
	lines = kept.Lines
	for _, dir := range dirs {
		_ = filepath.WalkDir(dir, func(path string, e fs.DirEntry, err error) error {
			if err != nil || e.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil //nolint:nilerr // an unreadable entry is skipped, the walk continues
			}
			scan(path)
			return nil
		})
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			scan(filepath.Join(root, e.Name()))
		}
	}
}
