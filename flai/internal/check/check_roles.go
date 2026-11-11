package check

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/formulas"
	"github.com/bytepunx/system-flow/flai/internal/guard"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

func (c *checker) conventions() {
	set, errs, err := conventions.Load(c.repo)
	if err != nil {
		c.add(Error, "conventions.read", conventions.Dir(c.repo), 1, "%v", err)
		return
	}
	for _, f := range set.Validate(errs, c.definedRoles()...) {
		path := filepath.Join(c.repo.Root, f.Path)
		line := 1
		switch f.Rule {
		case "conventions.front-matter", "conventions.order":
			line = keyLine(path, "order")
			if strings.Contains(f.Message, "audience") {
				line = keyLine(path, "audience")
			}
		}
		c.add(f.Level, f.Rule, path, line, "%s", f.Message)
	}
	for _, f := range set.Files {
		if _, bad := errs[f.Path]; bad {
			continue
		}
		path := filepath.Join(c.repo.Root, f.Path)
		doc, err := topics.Parse(f.Raw, []string{topics.All})
		if err != nil {
			c.add(Warning, "doc.topics", path, keyLine(path, "topics"), "%v", err)
			continue
		}
		c.topicWords(path, doc)
	}
}

// definedRoles are the roles the project's role files define (ADR-0207),
// which a convention's roles may list. A file that does not parse defines
// none; c.roles reports it.
func (c *checker) definedRoles() []string {
	list, _, err := roles.Load(c.repo.Root, c.repo.Manifest.Layout["design"])
	if err != nil {
		return nil
	}
	var out []string
	for _, r := range list {
		if r.Role != "" {
			out = append(out, r.Role)
		}
	}
	return out
}

// RolesVocabulary is what flai check holds role files' names to (ADR-0207):
// the reads every agent has, from the guard. Tools and Commands are nil, so
// the check of an entry's MCP tool or command name is skipped: the MCP
// server and the command tree import this package, and a check that read
// them in one caller and not another would differ where it runs. A test in
// cmd holds the shipped role files' names to both.
func RolesVocabulary() roles.Vocabulary {
	return roles.Vocabulary{Reads: guard.MCPReads}
}

// roles validates the design folder's role files (ADR-0207): each problem
// roles.Load or roles.Validate finds is an error under its code, on its line,
// and each built-in role file missing is an advisory roles.baseline warning.
func (c *checker) roles() {
	design := c.repo.Manifest.Layout["design"]
	list, problems, err := roles.Load(c.repo.Root, design)
	if err != nil {
		c.add(Error, roles.CodeFile, filepath.Join(c.repo.Root, filepath.FromSlash(roles.Dir(design))), 1, "%v", err)
		return
	}
	c.missingBaselineRoles(design)
	c.agentCopies(design)
	for _, p := range append(problems, roles.Validate(list, design, RolesVocabulary())...) {
		c.add(Error, p.Code, filepath.Join(c.repo.Root, filepath.FromSlash(p.Path)), p.Line, "%s", p.Message)
	}
	if len(problems) > 0 {
		return // a file that does not parse defines no role, so the manifest's roles cannot be held to them
	}
	mf := filepath.Join(c.repo.Root, manifest.File)
	for _, msg := range c.repo.Manifest.UndeclaredRoleSettings(declaredPermissions(list)) {
		c.add(Warning, "manifest.roles", mf, keyLine(mf, "roles"), "%s", msg)
	}
}

// formulas validates the design folder's formula files (ADR-0238 §6): each
// problem formulas.Load or formulas.Validate finds is a formulas.file error on
// its line. A role step is held to the names of the roles the role files
// define, and a script step to the scripts on disk. A project with no formula
// files' folder has none to report.
func (c *checker) formulas() {
	design := c.repo.Manifest.Layout["design"]
	list, problems, err := formulas.Load(c.repo.Root, design)
	if err != nil {
		c.add(Error, formulas.CodeFile, filepath.Join(c.repo.Root, filepath.FromSlash(formulas.Dir(design))), 1, "%v", err)
		return
	}
	names := formulas.Names{Roles: c.roleNames(), Script: formulas.ScriptOnDisk(c.repo.Root)}
	for _, p := range append(problems, formulas.Validate(list, names)...) {
		c.add(Error, p.Code, filepath.Join(c.repo.Root, filepath.FromSlash(p.Path)), p.Line, "%s", p.Message)
	}
	c.analysisFormula(list, problems)
}

// analysisFormula holds analysis.formula to the formula files (S-0424): each
// way formulas.AnalysisFormulaErrors finds it falls short, the formula no file
// defines or an input of it with no default, is a manifest.analysis error. A
// formula whose file does not parse is left to the formulas.file error on it.
func (c *checker) analysisFormula(list []formulas.Formula, problems []formulas.Problem) {
	mf := filepath.Join(c.repo.Root, manifest.File)
	for _, msg := range formulas.AnalysisFormulaErrors(list, problems, c.repo.Manifest.Layout["design"], c.repo.Manifest.Analysis.Formula) {
		c.add(Error, "manifest.analysis", mf, keyLine(mf, "analysis"), "%s", msg)
	}
}

// roleNames are the names of the roles the project's role files define
// (ADR-0207), which a formula's role step names (ADR-0238 §3), as definedRoles
// are their verbs. A file that does not parse defines none; c.roles reports it.
func (c *checker) roleNames() []string {
	list, _, err := roles.Load(c.repo.Root, c.repo.Manifest.Layout["design"])
	if err != nil {
		return nil
	}
	var out []string
	for _, r := range list {
		if r.Name != "" {
			out = append(out, r.Name)
		}
	}
	return out
}

// missingBaselineRoles advises roles.baseline for each role flai carries a
// baseline for whose file the project lacks: flai runs the baseline in its
// place, and flai upgrade adds the file (S-0419).
func (c *checker) missingBaselineRoles(design string) {
	for _, r := range roles.Baseline(design) {
		path := filepath.Join(c.repo.Root, filepath.FromSlash(r.Path))
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			continue
		}
		c.advise(roles.CodeBaseline, path, 1, "%s is missing: flai runs the %s from the baseline it carries; run flai upgrade to add the file from the template", r.Path, r.Name)
	}
}

// agentCopies warns roles.definition for each built-in role whose
// .claude/agents/<name>.md exists and is neither the link ADR-0236 names,
// ../../.flai-cache/agents/<name>.md, whether or not its target exists, nor
// a file whose description, tools, model, and prompt equal the definition
// flai generates from the role's file, naming the first that differs. A
// missing file is no finding, and neither is a role whose own definition
// cannot be generated: the role file's findings cover it. A warning, since
// only a session started by hand reads the file, so --strict fails on it.
func (c *checker) agentCopies(design string) {
	for _, base := range roles.Baseline(design) {
		name := base.Name
		rel := filepath.Join(".claude", "agents", name+".md")
		path := filepath.Join(c.repo.Root, rel)
		link := "../../.flai-cache/agents/" + name + ".md"
		fix := fmt.Sprintf("run flai role definition --write, then ln -sf %s %s (ADR-0236)", link, filepath.ToSlash(rel))
		info, err := os.Lstat(path)
		switch {
		case errors.Is(err, os.ErrNotExist):
			continue
		case err != nil:
			c.add(Warning, roles.CodeDefinition, path, 1, "%s cannot be read: %v; it should be the link %s: %s", filepath.ToSlash(rel), err, link, fix)
			continue
		case info.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(path)
			if err == nil && target == link {
				continue
			}
			if err != nil {
				target = fmt.Sprintf("a target that cannot be read (%v)", err)
			}
			c.add(Warning, roles.CodeDefinition, path, 1, "%s is a symbolic link to %s; it should be the link %s, to the definition flai generates for the %s: %s", filepath.ToSlash(rel), target, link, name, fix)
			continue
		case !info.Mode().IsRegular():
			c.add(Warning, roles.CodeDefinition, path, 1, "%s is neither a file nor a symbolic link; it should be the link %s: %s", filepath.ToSlash(rel), link, fix)
			continue
		}
		r, fromBaseline, err := roles.Builtin(c.repo.Root, design, name)
		if err != nil {
			continue
		}
		if fromBaseline {
			if _, err := os.Stat(filepath.Join(c.repo.Root, filepath.FromSlash(r.Path))); err == nil {
				continue // the role file does not parse, which its roles.file error reports
			}
		}
		want, err := roles.Generate(r)
		if err != nil {
			continue // the role file's own roles.definition errors report it
		}
		source := "its role file " + r.Path
		if fromBaseline {
			source = "the baseline flai carries for " + r.Path
		}
		data, err := os.ReadFile(path)
		if err != nil {
			c.add(Warning, roles.CodeDefinition, path, 1, "%s cannot be read: %v; replace it with the link to the definition flai generates for the %s: %s", filepath.ToSlash(rel), err, name, fix)
			continue
		}
		got, err := roles.ParseAgentFile(data)
		if err != nil {
			c.add(Warning, roles.CodeDefinition, path, 1, "%s is not a Claude Code agent file: %v; replace it with the link to the definition flai generates from %s: %s", filepath.ToSlash(rel), err, source, fix)
			continue
		}
		field, line := agentDiff(got, want, path)
		if field == "" {
			continue
		}
		c.add(Warning, roles.CodeDefinition, path, line, "%s is a copy whose %s differs from the definition flai generates for the %s from %s; replace it with the link to that definition: %s", filepath.ToSlash(rel), field, name, source, fix)
	}
}

// agentDiff is the first of description, tools, model, and prompt in which
// the agent file at path, read as got, differs from the definition want, and
// its line there: the key's in the front matter, or 1 for the prompt. An
// empty field means none differs.
func agentDiff(got, want roles.Generated, path string) (string, int) {
	model := func(m string) string {
		if m == "" {
			return roles.Inherit
		}
		return m
	}
	switch {
	case strings.TrimSpace(got.Description) != strings.TrimSpace(want.Description):
		return "description", keyLine(path, "description")
	case !slices.Equal(got.Tools, want.Tools):
		return "tools", keyLine(path, "tools")
	case model(got.Model) != model(want.Model):
		return "model", keyLine(path, "model")
	case strings.TrimSpace(got.Prompt) != strings.TrimSpace(want.Prompt):
		return "prompt", 1
	}
	return "", 0
}

// declaredPermissions are, by the name of each role the files define, the
// permissions its entries name as roles.<name>.permissions.<p>.
func declaredPermissions(list []roles.Role) map[string][]string {
	out := map[string][]string{}
	for _, r := range list {
		perms := out[r.Name]
		for _, p := range r.Permissions {
			if perm, owner, ok := roles.ParseSetting(p.Setting); ok && owner == r.Name && strings.HasPrefix(p.Setting, "roles.") {
				perms = append(perms, perm)
			}
		}
		out[r.Name] = perms
	}
	return out
}

// unknownFields reports each front-matter field this flai does not know: the
// listing paths read past them (S-0181), and check stays strict.
func (c *checker) unknownFields(rule, path string, fields []workitem.Field) {
	for i, msg := range workitem.UnknownFieldErrors(fields) {
		c.add(Error, rule, path, keyLine(path, fields[i].Name), "%s", msg)
	}
}
