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
		return "", fmt.Errorf("%s is an experiment; an experiment stays on its branch and is not accepted onto main (ADR-0025)", it.ID)
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

// PendingItem is one accepted item counted toward a PendingPlan.
type PendingItem struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Level string `json:"level"`
}

// PendingPlan is one component's batch: everything accepted and unreleased
// for it since its last tag, at the highest delivery level among it
// (S-0087). Publishing merges several stories' releases into one, rather
// than each cutting its own: a feature and two remediations against the
// same component since its last tag produce one minor release, not three.
type PendingPlan struct {
	Component manifest.Project `json:"component"`
	Level     string           `json:"level"`
	From      Version          `json:"from"`
	To        Version          `json:"to"`
	Tag       string           `json:"tag,omitempty"`
	Version   string           `json:"version,omitempty"`
	Items     []PendingItem    `json:"items"`
	Files     []string         `json:"files"`
}

var rank = map[string]int{None: 0, Patch: 1, Minor: 2, Major: 3}

// higher is the more urgent of two bump levels; "" counts as None.
func higher(a, b string) string {
	if rank[b] > rank[a] {
		return b
	}
	return a
}

var acceptedSubject = regexp.MustCompile(`^chore: \[([EST]-\d+)\] accept and archive`)

// acceptedIDsIn lists the items accepted by a commit in revRange (a git
// revision range, or a single ref meaning everything reachable from it),
// each once, oldest first.
func acceptedIDsIn(r execx.Runner, root, revRange string) ([]string, error) {
	out, err := r.Run(root, "git", "log", "--reverse", "--format=%s", revRange)
	if err != nil {
		return nil, err
	}
	var ids []string
	seen := map[string]bool{}
	for _, l := range strings.Split(out, "\n") {
		if m := acceptedSubject.FindStringSubmatch(l); m != nil && !seen[m[1]] {
			seen[m[1]] = true
			ids = append(ids, m[1])
		}
	}
	return ids, nil
}

// componentBoundary is the ref marking the last publish for a component, or
// "" when it has never been released: everything accepted since is pending.
// A code component's tag is that ref directly. A template has no local tag
// (Tag creates none for it; its version lives in template.yaml, published to
// its own remote); the commit that set its current version stands in.
func componentBoundary(r execx.Runner, root string, p manifest.Project, cur Version) (string, error) {
	if cur == (Version{}) {
		return "", nil
	}
	if p.Kind != "template" {
		return p.Name + "/v" + cur.String(), nil
	}
	out, err := r.Run(root, "git", "log", "-n1", "--format=%H", "-S", "version: "+cur.String(), "--", filepath.Join(p.Path, "template.yaml"))
	if err != nil || strings.TrimSpace(out) == "" {
		return "", nil
	}
	return strings.TrimSpace(out), nil
}

// Pending computes the batch: one PendingPlan per component with something
// accepted and unreleased for it, from every item a "chore: [ID] accept and
// archive" commit names since that component's last publish (S-0087). An
// item is looked up once and its Compute()-equivalent reused for every
// component it touches or delivers to.
func Pending(r execx.Runner, root string, m manifest.Manifest, repo *workitem.Repo) ([]*PendingPlan, error) {
	plans := map[string]*Plan{} // item ID -> its own plan, computed once
	itemPlan := func(id string) (*Plan, error) {
		if p, ok := plans[id]; ok {
			return p, nil
		}
		it, err := repo.Get(id)
		if err != nil {
			return nil, err
		}
		var parent *workitem.Item
		if it.Parent != "" {
			parent, _ = repo.Get(it.Parent)
		}
		p, err := Compute(r, root, m, it, parent, "")
		if err != nil {
			return nil, err
		}
		plans[id] = p
		return p, nil
	}

	var out []*PendingPlan
	for _, proj := range m.Projects {
		cur, err := CurrentVersion(r, root, proj)
		if err != nil {
			return nil, err
		}
		boundary, err := componentBoundary(r, root, proj, cur)
		if err != nil {
			return nil, err
		}
		revRange := "HEAD"
		if boundary != "" {
			revRange = boundary + "..HEAD"
		}
		ids, err := acceptedIDsIn(r, root, revRange)
		if err != nil {
			return nil, err
		}
		pp := &PendingPlan{Component: proj, From: cur}
		seenFile := map[string]bool{}
		for _, id := range ids {
			ip, err := itemPlan(id)
			if err != nil {
				continue // an item that cannot be recomputed (removed, malformed) is skipped, not fatal to the batch
			}
			for _, st := range ip.Steps {
				if st.Component.Name != proj.Name {
					continue
				}
				pp.Level = higher(pp.Level, st.Level)
				pp.Items = append(pp.Items, PendingItem{ID: id, Title: ip.Title, Level: st.Level})
				for _, f := range st.Files {
					if !seenFile[f] {
						seenFile[f] = true
						pp.Files = append(pp.Files, f)
					}
				}
			}
		}
		if len(pp.Items) == 0 || pp.Level == "" || pp.Level == None {
			continue
		}
		sort.Strings(pp.Files)
		pp.To = cur.Bumped(pp.Level)
		if proj.Kind == "template" {
			pp.Version = filepath.ToSlash(filepath.Join(proj.Path, "template.yaml"))
		} else {
			pp.Tag = proj.Name + "/v" + pp.To.String()
		}
		out = append(out, pp)
	}
	return out, nil
}
