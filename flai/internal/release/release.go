// Package release computes and applies semver releases on acceptance, per
// design/conventions/git.md: the component an item delivers to gets the
// delivery-type bump, other touched components a patch.
package release

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Bump levels.
const (
	Major = "major"
	Minor = "minor"
	Patch = "patch"
	// None is the level of an item that is accepted but cuts no release:
	// a research story (ADR-0025).
	None = "none"
)

// Version is a parsed semver.
type Version struct{ Major, Minor, Patch int }

func (v Version) String() string { return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch) }

// Bumped returns the next version for a level.
func (v Version) Bumped(level string) Version {
	switch level {
	case Major:
		return Version{v.Major + 1, 0, 0}
	case Minor:
		return Version{v.Major, v.Minor + 1, 0}
	default:
		return Version{v.Major, v.Minor, v.Patch + 1}
	}
}

var semverRe = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)$`)

// ParseVersion reads X.Y.Z with an optional v.
func ParseVersion(s string) (Version, bool) {
	m := semverRe.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return Version{}, false
	}
	a, _ := strconv.Atoi(m[1])
	b, _ := strconv.Atoi(m[2])
	c, _ := strconv.Atoi(m[3])
	return Version{a, b, c}, true
}

// Step is one component's release.
type Step struct {
	Component manifest.Project `json:"component"`
	Delivered bool             `json:"delivered"` // the item delivers to this component
	Level     string           `json:"level"`
	From      Version          `json:"from"`
	To        Version          `json:"to"`
	Tag       string           `json:"tag,omitempty"`     // code components: <name>/vX.Y.Z
	Version   string           `json:"version,omitempty"` // template kind: file to bump
	Files     []string         `json:"files"`             // touched files that mapped here
}

// Plan is the computed release for an item.
type Plan struct {
	Item      string   `json:"item"`
	Title     string   `json:"title"`
	Level     string   `json:"level"`
	Commits   []string `json:"commits"`
	Untouched []string `json:"untouched_files"` // touched files mapping to no component
	Steps     []Step   `json:"steps"`
	Skipped   string   `json:"skipped,omitempty"` // reason nothing releases
	// Unreleased names the components whose files the item's commits touched
	// although it cuts no release: code landing on main without a version.
	Unreleased []Unreleased `json:"unreleased,omitempty"`
}

// Unreleased is a component a no-release item touched.
type Unreleased struct {
	Component string   `json:"component"`
	Files     []string `json:"files"`
}

// LevelFor maps an item's type and nature to a bump level. A research story
// is accepted without a release (None); an experiment stays on its branch and
// is refused (ADR-0025).
func LevelFor(it *workitem.Item) (string, error) {
	if it.Type == workitem.Epic {
		return Major, nil
	}
	switch it.Nature {
	case "feature":
		return Minor, nil
	case "remediation", "improvement":
		return Patch, nil
	case "research":
		return None, nil
	case "experiment":
		return "", fmt.Errorf("%s is an experiment; an experiment stays on its branch and is not accepted onto main (ADR-0025); flai accept --no-release lands one deliberately", it.ID)
	}
	return "", fmt.Errorf("%s has nature %q, no release rule", it.ID, it.Nature)
}

// Commits returns the hashes of commits whose message names the item.
func Commits(r execx.Runner, root, id string) ([]string, error) {
	out, err := r.Run(root, "git", "log", "--format=%H", "--fixed-strings", "--grep=["+id+"]")
	if err != nil {
		return nil, err
	}
	var hashes []string
	for _, l := range strings.Split(strings.TrimSpace(out), "\n") {
		if l != "" {
			hashes = append(hashes, l)
		}
	}
	return hashes, nil
}

// TouchedFiles lists files changed by the commits.
func TouchedFiles(r execx.Runner, root string, commits []string) ([]string, error) {
	seen := map[string]bool{}
	var files []string
	for _, c := range commits {
		out, err := r.Run(root, "git", "diff-tree", "--no-commit-id", "--name-only", "-r", "--root", c)
		if err != nil {
			return nil, err
		}
		for _, f := range strings.Split(strings.TrimSpace(out), "\n") {
			if f != "" && !seen[f] {
				seen[f] = true
				files = append(files, f)
			}
		}
	}
	sort.Strings(files)
	return files, nil
}

// CurrentVersion finds a component's latest version: the highest
// <name>/vX.Y.Z tag for code components, template.yaml for the template.
func CurrentVersion(r execx.Runner, root string, p manifest.Project) (Version, error) {
	if p.Kind == "template" {
		data, err := os.ReadFile(filepath.Join(root, p.Path, "template.yaml"))
		if err != nil {
			return Version{}, err
		}
		for _, l := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(l, "version:") {
				if v, ok := ParseVersion(strings.TrimSpace(strings.TrimPrefix(l, "version:"))); ok {
					return v, nil
				}
			}
		}
		return Version{}, fmt.Errorf("%s/template.yaml has no semver version", p.Path)
	}
	out, err := r.Run(root, "git", "tag", "--list", p.Name+"/v*")
	if err != nil {
		return Version{}, err
	}
	var best Version
	for _, t := range strings.Split(strings.TrimSpace(out), "\n") {
		if v, ok := ParseVersion(strings.TrimPrefix(t, p.Name+"/")); ok && less(best, v) {
			best = v
		}
	}
	return best, nil
}

func less(a, b Version) bool {
	if a.Major != b.Major {
		return a.Major < b.Major
	}
	if a.Minor != b.Minor {
		return a.Minor < b.Minor
	}
	return a.Patch < b.Patch
}

// Compute builds the plan for an accepted item. deliver names the
// component the item delivers to when tags do not say.
func Compute(r execx.Runner, root string, m manifest.Manifest, it *workitem.Item, parent *workitem.Item, deliver string) (*Plan, error) {
	level, err := LevelFor(it)
	if err != nil {
		return nil, err
	}
	plan := &Plan{Item: it.ID, Title: it.Title, Level: level}
	plan.Commits, err = Commits(r, root, it.ID)
	if err != nil {
		return nil, err
	}
	files, err := TouchedFiles(r, root, plan.Commits)
	if err != nil {
		return nil, err
	}
	touched := map[string][]string{}
	for _, f := range files {
		matched := false
		for _, p := range m.Projects {
			if f == p.Path || strings.HasPrefix(f, p.Path+"/") {
				touched[p.Name] = append(touched[p.Name], f)
				matched = true
				break
			}
		}
		if !matched {
			plan.Untouched = append(plan.Untouched, f)
		}
	}
	if level == None {
		plan.Skipped = fmt.Sprintf("%s is research: its findings land on main and are pushed, and research cuts no release whatever it touched (ADR-0025)", it.ID)
		for _, p := range m.Projects {
			if files, was := touched[p.Name]; was {
				plan.Unreleased = append(plan.Unreleased, Unreleased{Component: p.Name, Files: files})
			}
		}
		return plan, nil
	}
	delivered := deliverTarget(m, it, parent, deliver, touched)
	if len(touched) == 0 && delivered == "" {
		plan.Skipped = "no component touched by the item's commits; design, docs, and wip changes release nothing"
		return plan, nil
	}
	if delivered == "" {
		names := make([]string, 0, len(touched))
		for n := range touched {
			names = append(names, n)
		}
		sort.Strings(names)
		return nil, fmt.Errorf("%s touches %s but no tag says which it delivers to; add a tag matching a project or pass --deliver", it.ID, strings.Join(names, ", "))
	}
	for _, p := range m.Projects {
		files, was := touched[p.Name]
		isDelivered := p.Name == delivered
		if !was && !isDelivered {
			continue
		}
		cur, err := CurrentVersion(r, root, p)
		if err != nil {
			return nil, err
		}
		lvl := Patch
		if isDelivered {
			lvl = level
		}
		st := Step{Component: p, Delivered: isDelivered, Level: lvl, From: cur, To: cur.Bumped(lvl), Files: files}
		if p.Kind == "template" {
			st.Version = filepath.ToSlash(filepath.Join(p.Path, "template.yaml"))
		} else {
			st.Tag = p.Name + "/v" + st.To.String()
		}
		plan.Steps = append(plan.Steps, st)
	}
	sort.SliceStable(plan.Steps, func(i, j int) bool { return plan.Steps[i].Delivered && !plan.Steps[j].Delivered })
	return plan, nil
}

// deliverTarget picks the component the item delivers to. --deliver is the
// operator's explicit word and wins. Otherwise a tag decides, the item's own
// before its parent's, and only among components the item's commits touched:
// a tag naming a component no commit touched never delivers, so that
// component gets no release at all (I-0016: a story of pure CLI work tagged
// [dashboard, cli] gave the dashboard a minor release with no files in it).
// When the tags name several touched components, the one with the most
// touched files wins, the earlier tag breaking a tie. With no tag deciding,
// the only touched component delivers.
func deliverTarget(m manifest.Manifest, it, parent *workitem.Item, deliver string, touched map[string][]string) string {
	if deliver != "" {
		return deliver
	}
	match := func(tags []string) string {
		best, files := "", 0
		for _, t := range tags {
			for _, p := range m.Projects {
				named := t == p.Name
				for _, alias := range p.Tags {
					named = named || t == alias
				}
				if n := len(touched[p.Name]); named && n > files {
					best, files = p.Name, n
				}
			}
		}
		return best
	}
	if n := match(it.Tags); n != "" {
		return n
	}
	if parent != nil {
		if n := match(parent.Tags); n != "" {
			return n
		}
	}
	if len(touched) == 1 {
		for n := range touched {
			return n
		}
	}
	return ""
}

// Apply performs the plan: template version and changelog bumps (files
// only; the caller commits), then annotated tags on HEAD.
func Apply(r execx.Runner, root string, plan *Plan, now time.Time) error {
	for _, st := range plan.Steps {
		if st.Version != "" {
			if err := bumpTemplate(root, st, plan, now); err != nil {
				return err
			}
		}
	}
	return nil
}

// Tag creates the annotated tags for the plan on HEAD.
func Tag(r execx.Runner, root string, plan *Plan) ([]string, error) {
	var tags []string
	for _, st := range plan.Steps {
		if st.Tag == "" {
			continue
		}
		msg := fmt.Sprintf("%s v%s: %s %s (%s %s)", st.Component.Name, st.To, plan.Item, plan.Title, st.Level, map[bool]string{true: "delivered", false: "incidental"}[st.Delivered])
		if _, err := r.Run(root, "git", "tag", "-a", st.Tag, "-m", msg); err != nil {
			return tags, err
		}
		tags = append(tags, st.Tag)
	}
	return tags, nil
}

func bumpTemplate(root string, st Step, plan *Plan, now time.Time) error {
	path := filepath.Join(root, st.Version)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, "version:") {
			lines[i] = "version: " + st.To.String()
		}
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		return err
	}
	cl := filepath.Join(root, st.Component.Path, "CHANGELOG.md")
	entry := fmt.Sprintf("## %s - %s\n\n- %s %s (%s).\n", st.To, now.UTC().Format("2006-01-02"), plan.Item, plan.Title, st.Level)
	existing, err := os.ReadFile(cl)
	if err != nil {
		return os.WriteFile(cl, []byte("# Changelog\n\n"+entry), 0o644)
	}
	s := string(existing)
	if i := strings.Index(s, "\n## "); i >= 0 {
		s = s[:i+1] + entry + "\n" + s[i+1:]
	} else {
		s = strings.TrimRight(s, "\n") + "\n\n" + entry
	}
	return os.WriteFile(cl, []byte(s), 0o644)
}
