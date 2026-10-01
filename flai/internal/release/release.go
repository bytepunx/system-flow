// Package release computes and applies semver releases on acceptance, per
// design/conventions/git.md: the component an item delivers to gets the
// delivery-type bump, other touched components a patch.
package release

import (
	"encoding/json"
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

// MarshalJSON writes a Version as its plain "X.Y.Z" string, not the struct
// fields: flaiover reads From/To as strings (found in S-0087's T-0318 review,
// where the done column's Publish banner showed "[object Object]" for them).
func (v Version) MarshalJSON() ([]byte, error) { return json.Marshal(v.String()) }

// UnmarshalJSON reads a Version from its plain "X.Y.Z" string, the mirror of MarshalJSON.
func (v *Version) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	parsed, ok := ParseVersion(s)
	if !ok {
		return fmt.Errorf("%q is not a version", s)
	}
	*v = parsed
	return nil
}

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
		return templateVersion(root, p)
	}
	out, err := r.Run(root, "git", "tag", "--list", p.Name+"/v*")
	if err != nil {
		return Version{}, err
	}
	return highestTag(strings.Split(strings.TrimSpace(out), "\n"), p.Name), nil
}

// templateVersion reads the version: line of a template's template.yaml.
func templateVersion(root string, p manifest.Project) (Version, error) {
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

// highestTag is the highest version among the tags named <name>/vX.Y.Z.
func highestTag(tags []string, name string) Version {
	var best Version
	for _, t := range tags {
		if !strings.HasPrefix(t, name+"/v") {
			continue
		}
		if v, ok := ParseVersion(strings.TrimPrefix(t, name+"/")); ok && less(best, v) {
			best = v
		}
	}
	return best
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
	commits, err := Commits(r, root, it.ID)
	if err != nil {
		return nil, err
	}
	files, err := TouchedFiles(r, root, commits)
	if err != nil {
		return nil, err
	}
	return computePlan(m, it, parent, deliver, level, commits, files, func(p manifest.Project) (Version, error) { return CurrentVersion(r, root, p) })
}

// computePlan is Compute once the item's commits, the files they touched, and a way
// to find a component's current version are known.
func computePlan(m manifest.Manifest, it, parent *workitem.Item, deliver, level string, commits, files []string, current func(manifest.Project) (Version, error)) (*Plan, error) {
	plan := &Plan{Item: it.ID, Title: it.Title, Level: level, Commits: commits}
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
		cur, err := current(p)
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
	bullet := fmt.Sprintf("- %s %s (%s).\n", plan.Item, plan.Title, st.Level)
	return bumpVersionAndChangelog(root, st.Version, st.Component.Path, st.To, now, bullet)
}

// bumpVersionAndChangelog writes the version file's version: line and adds
// one changelog entry above the previous one, its body the given bullets.
func bumpVersionAndChangelog(root, versionFile, componentPath string, to Version, now time.Time, bullets string) error {
	path := filepath.Join(root, versionFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for i, l := range lines {
		if strings.HasPrefix(l, "version:") {
			lines[i] = "version: " + to.String()
		}
	}
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644); err != nil {
		return err
	}
	cl := filepath.Join(root, componentPath, "CHANGELOG.md")
	entry := fmt.Sprintf("## %s - %s\n\n%s", to, now.UTC().Format("2006-01-02"), bullets)
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

// ApplyPending performs a batch plan: the version file and one changelog
// entry naming every item it bundles (files only; the caller commits).
func ApplyPending(plan *PendingPlan, root string, now time.Time) error {
	if plan.Version == "" {
		return nil
	}
	var bullets strings.Builder
	for _, it := range plan.Items {
		fmt.Fprintf(&bullets, "- %s %s (%s).\n", it.ID, it.Title, it.Level)
	}
	return bumpVersionAndChangelog(root, plan.Version, plan.Component.Path, plan.To, now, bullets.String())
}

// TagPending creates the plan's tag on HEAD, unless it already exists: a
// prior run of flai release --pending may have created it and failed before
// pushing, and running it again must not fail on an existing tag (S-0087).
// "" is returned for a template plan (Tag() never tags one) or one already
// tagged.
func TagPending(r execx.Runner, root string, plan *PendingPlan) (string, error) {
	if plan.Tag == "" {
		return "", nil
	}
	if out, err := r.Run(root, "git", "tag", "--list", plan.Tag); err == nil && strings.TrimSpace(out) != "" {
		return "", nil
	}
	items := make([]string, len(plan.Items))
	for i, it := range plan.Items {
		items[i] = it.ID
	}
	msg := fmt.Sprintf("%s v%s: %s (%s)", plan.Component.Name, plan.To, strings.Join(items, ", "), plan.Level)
	if _, err := r.Run(root, "git", "tag", "-a", plan.Tag, "-m", msg); err != nil {
		return "", err
	}
	return plan.Tag, nil
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

// Unplanned is an item accepted since a component's last publish whose own
// plan could not be computed (I-0024): removed, with no release rule, or
// touching two components with no tag saying which it delivers to. It is in
// no PendingPlan, and its work misses its bump if published as it is.
type Unplanned struct {
	ID     string `json:"id"`
	Title  string `json:"title,omitempty"`
	Reason string `json:"reason"`
}

// Batch is everything Pending finds: the plans, and the accepted items it
// could not plan.
type Batch struct {
	Plans     []*PendingPlan `json:"plans"`
	Unplanned []Unplanned    `json:"unplanned,omitempty"`
}

// Pending computes the batch's plans: one PendingPlan per component with
// something accepted and unreleased for it, from every item a "chore: [ID]
// accept and archive" commit names since that component's last publish
// (S-0087). An item that cannot be planned is left out; PendingBatch names
// it.
func Pending(r execx.Runner, root string, m manifest.Manifest, repo *workitem.Repo) ([]*PendingPlan, error) {
	b, err := PendingBatch(r, root, m, repo)
	if err != nil {
		return nil, err
	}
	return b.Plans, nil
}

// PendingBatch is Pending with the items it could not plan, each once. An
// item is looked up once and its Compute()-equivalent reused for every
// component it touches or delivers to. What it needs of git comes from the
// repository's kept history (S-0157): one git process while HEAD and the
// tags are unchanged.
func PendingBatch(r execx.Runner, root string, m manifest.Manifest, repo *workitem.Repo) (*Batch, error) {
	batch := &Batch{}
	if len(m.Projects) == 0 {
		return batch, nil // no component, nothing to release, and no git asked
	}
	kept.Lock()
	defer kept.Unlock()
	h, err := readHistory(r, root)
	if err != nil {
		return nil, err
	}
	type since struct {
		cur Version
		ids []string
	}
	ranges := make([]since, len(m.Projects))
	var commits []string
	for i, proj := range m.Projects {
		cur, err := h.currentVersion(root, proj)
		if err != nil {
			return nil, err
		}
		boundary, err := h.boundary(r, root, proj, cur)
		if err != nil {
			return nil, err
		}
		ids, err := h.accepted(r, root, boundary)
		if err != nil {
			return nil, err
		}
		ranges[i] = since{cur, ids}
		for _, id := range ids {
			commits = append(commits, h.itemCommits(id)...)
		}
	}
	if err := h.list(r, root, commits); err != nil {
		return nil, err
	}

	plans := map[string]*Plan{} // item ID -> its own plan, computed once
	unplanned := map[string]bool{}
	itemPlan := func(id string) (*Plan, error) {
		if p, ok := plans[id]; ok {
			return p, nil
		}
		title := ""
		fail := func(err error) (*Plan, error) {
			if !unplanned[id] {
				unplanned[id] = true
				batch.Unplanned = append(batch.Unplanned, Unplanned{ID: id, Title: title, Reason: err.Error()})
			}
			return nil, err
		}
		it, err := repo.Get(id)
		if err != nil {
			return fail(err)
		}
		title = it.Title
		var parent *workitem.Item
		if it.Parent != "" {
			parent, _ = repo.Get(it.Parent)
		}
		level, err := LevelFor(it)
		if err != nil {
			return fail(err)
		}
		commits := h.itemCommits(id)
		p, err := computePlan(m, it, parent, "", level, commits, h.files(commits), func(p manifest.Project) (Version, error) { return h.currentVersion(root, p) })
		if err != nil {
			return fail(err)
		}
		plans[id] = p
		return p, nil
	}

	for i, proj := range m.Projects {
		cur, ids := ranges[i].cur, ranges[i].ids
		pp := &PendingPlan{Component: proj, From: cur}
		seenFile := map[string]bool{}
		for _, id := range ids {
			ip, err := itemPlan(id)
			if err != nil {
				continue // not fatal to the batch: named in Unplanned
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
		batch.Plans = append(batch.Plans, pp)
	}
	return batch, nil
}

// PendingIDs is every story and epic ID Pending finds still unpublished, for
// a caller that only needs to know which items those are — the board,
// keeping a done, archived item visible until it is published (S-0087) —
// not the full per-component plan. A history that cannot be read (no git,
// for instance) is nothing pending rather than an error the caller must
// handle: whether to show it is a lesser concern than whether to show the
// board at all.
func PendingIDs(r execx.Runner, root string, m manifest.Manifest, repo *workitem.Repo) map[string]bool {
	b, err := PendingBatch(r, root, m, repo)
	if err != nil {
		return nil
	}
	ids := map[string]bool{}
	for _, p := range b.Plans {
		for _, it := range p.Items {
			ids[it.ID] = true
		}
	}
	// an item that could not be planned is unpublished too, and stays in
	// view until it is put right (I-0024)
	for _, u := range b.Unplanned {
		ids[u.ID] = true
	}
	return ids
}
