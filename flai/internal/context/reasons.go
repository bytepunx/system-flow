package context

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/conventions"
	"github.com/bytepunx/system-flow/flai/internal/digest"
	"github.com/bytepunx/system-flow/flai/internal/threads"
	"github.com/bytepunx/system-flow/flai/internal/topics"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// The kinds of reason the orchestrator is primed for (ADR-0209 §3,
// ADR-0217 §1).
const (
	ReasonThread  = "thread"
	ReasonReview  = "review"
	ReasonPlanned = "planned"
	ReasonDraft   = "draft"
	ReasonRelease = "release"
	ReasonReady   = "ready"
)

// ReasonKinds is every kind of reason, in the order ADR-0217 lists them.
var ReasonKinds = []string{ReasonThread, ReasonReview, ReasonPlanned, ReasonDraft, ReasonRelease, ReasonReady}

// reasonIDs is the ID prefixes each kind of reason takes, none for a kind
// given bare (ADR-0217 §1).
var reasonIDs = map[string][]string{
	ReasonThread:  {"TH"},
	ReasonReview:  {"S"},
	ReasonPlanned: {"E", "S"},
	ReasonDraft:   {"S"},
	ReasonRelease: nil,
	ReasonReady:   nil,
}

// reasonID is an ID as CanonicalID pads it: its prefix, a dash, and its
// digits.
var reasonID = regexp.MustCompile(`^([A-Z]+)-\d+$`)

// StepReason is the step of what a reason loads: a thread, a work item, or,
// from T-2309, a figure.
const StepReason = "reason"

// Reason is one reason the orchestrator is primed for: its kind and, for a
// kind that takes one, the ID it names in the padding flai allocates.
type Reason struct {
	Kind string `json:"kind"`
	ID   string `json:"id,omitempty"`
}

// String is the reason as flai prime --reason takes it: kind:ID, or the
// kind alone.
func (r Reason) String() string {
	if r.ID == "" {
		return r.Kind
	}
	return r.Kind + ":" + r.ID
}

// ReasonForms is the forms of every kind of reason, as a refusal names
// them: "thread:TH-nnnn, review:S-nnnn, planned:E-nnnn or planned:S-nnnn,
// draft:S-nnnn, release, or ready".
func ReasonForms() string {
	var each []string
	for _, k := range ReasonKinds {
		var forms []string
		for _, p := range reasonIDs[k] {
			forms = append(forms, k+":"+p+"-nnnn")
		}
		if len(forms) == 0 {
			forms = []string{k}
		}
		each = append(each, strings.Join(forms, " or "))
	}
	return strings.Join(each[:len(each)-1], ", ") + ", or " + each[len(each)-1]
}

// ParseReasons reads the reasons flai prime --reason and MCP prime's
// reasons give (ADR-0217 §1), each <kind>:<id> or a bare <kind>, with the
// ID in any zero padding, and drops a reason given twice. An unknown kind,
// a kind without the ID it needs or with one it does not take, and an ID
// of the wrong type are refused, naming the kinds and their forms.
func ParseReasons(list []string) ([]Reason, error) {
	var out []Reason
	for _, s := range list {
		r, err := parseReason(s)
		if err != nil {
			return nil, fmt.Errorf("flai prime --reason %s: %w; a reason is %s", s, err, ReasonForms())
		}
		if !slices.Contains(out, r) {
			out = append(out, r)
		}
	}
	return out, nil
}

// parseReason reads one reason, saying what is wrong with it when it does
// not parse.
func parseReason(s string) (Reason, error) {
	kind, id, withID := strings.Cut(strings.TrimSpace(s), ":")
	kind, id = strings.ToLower(strings.TrimSpace(kind)), strings.TrimSpace(id)
	prefixes, known := reasonIDs[kind]
	switch {
	case !known:
		return Reason{}, fmt.Errorf("no such kind of reason %q", kind)
	case len(prefixes) == 0 && withID:
		return Reason{}, fmt.Errorf("%s takes no ID", kind)
	case len(prefixes) == 0:
		return Reason{Kind: kind}, nil
	case id == "":
		return Reason{}, fmt.Errorf("%s needs the ID of %s", kind, idNames(prefixes))
	}
	canon := workitem.CanonicalID(id)
	if prefixes[0] == "TH" {
		canon = threads.CanonicalID(id)
	}
	m := reasonID.FindStringSubmatch(canon)
	if m == nil || !slices.Contains(prefixes, m[1]) {
		return Reason{}, fmt.Errorf("%s is not the ID of %s", id, idNames(prefixes))
	}
	return Reason{Kind: kind, ID: canon}, nil
}

// idNames names what IDs with these prefixes are, as "a thread (TH-nnnn)"
// or "an epic (E-nnnn) or a story (S-nnnn)".
func idNames(prefixes []string) string {
	name := map[string]string{"TH": "a thread (TH-nnnn)", "E": "an epic (E-nnnn)", "S": "a story (S-nnnn)"}
	var out []string
	for _, p := range prefixes {
		out = append(out, name[p])
	}
	return strings.Join(out, " or ")
}

// SectionWords is the words a convention heading may carry beside all for
// a reason pack to read (ADR-0217 §2): the kinds of reason, then the
// built-in strategic roles' topics, each once, which flai check accepts on
// a heading as it accepts a story's topics.
func SectionWords() []string {
	out := slices.Clone(ReasonKinds)
	for _, role := range conventions.StrategicRoles {
		for _, t := range roleTopics[role] {
			if !slices.Contains(out, t) {
				out = append(out, t)
			}
		}
	}
	return out
}

// MatchesReasons says whether a reason pack keeps a convention section with
// these topics (ADR-0217 §2): one whose topics are all alone, as every pack
// does; one whose topics are all and other words only when one of those
// words is in words, the kinds of the run's reasons and the role's topics;
// and one without all as Matches has it for the pack's topics.
func MatchesReasons(section, words, pack []string) bool {
	if !slices.Contains(section, topics.All) {
		return Matches(section, pack)
	}
	others := false
	for _, t := range section {
		if t == topics.All {
			continue
		}
		if slices.Contains(words, t) {
			return true
		}
		others = true
	}
	return !others
}

// ForReasons builds the orchestrator's pack primed for the reasons of its
// run (ADR-0209 §3, ADR-0217), each as ParseReasons reads it: the
// conventions the orchestrate role reads, once, with the sections
// MatchesReasons leaves out taken out; then, per reason, what it needs, as
// reasonLoaders loads it, each item once however many reasons name it; and
// the digest's answers and recommendations ranked against the reasons'
// threads, as answersFor takes them for the general pack (ADR-0194), within
// a sixteenth of the budget rather than an eighth (ADR-0218 §2). It holds no design or tech briefs, no ranked
// sections, no open issues, and no catalog: doc_search and doc_get find the
// rest. What the reasons load is kept when it takes the pack over the
// budget, as what a story names is. budget is a size as ParseSize reads
// it; empty means the project's prime.budget, else DefaultBudget.
func ForReasons(repo *workitem.Repo, reasons []string, budget string) (*Pack, error) {
	role := conventions.RoleOrchestrate
	cmd := "flai prime --role " + role
	rs, err := ParseReasons(reasons)
	if err != nil {
		return nil, err
	}
	if len(rs) == 0 {
		return nil, fmt.Errorf("%s: a pack primed for reasons needs one at least; give --reason %s, or none of them for the general pack", cmd, ReasonForms())
	}
	r, err := resolveRole(repo, cmd, role)
	if err != nil {
		return nil, err
	}
	size, err := Budget(budget, repo.Manifest.Prime.Budget)
	if err != nil {
		return nil, err
	}
	pr, err := loadProject(repo, cmd)
	if err != nil {
		return nil, err
	}
	packTopics := withTopics(nil, role, r.topics)
	names := topics.Names(packTopics)
	words := slices.Clone(r.topics)
	for _, rr := range rs {
		words = append(words, rr.Kind)
	}
	keep := func(section []string) bool { return MatchesReasons(section, words, names) }
	pack, err := build(repo.Root, "", "", packTopics, pr.readBy(role, false, names), "", size, keep)
	if err != nil {
		return nil, err
	}
	pack.Role, pack.agent = role, r.agent
	l := &reasonLoad{repo: repo, cmd: cmd}
	for _, rr := range rs {
		pack.Reasons = append(pack.Reasons, rr.String())
		load, ok := reasonLoaders[rr.Kind]
		if !ok {
			continue
		}
		if err := load(l, rr); err != nil {
			return nil, err
		}
	}
	if pack.Size.Bytes > pack.Budget {
		pack.Exceeded = ExceededConventions
	}
	pack.Items = append(pack.Items, l.items...)
	pack.size()
	if pack.Exceeded == "" && pack.Size.Bytes > pack.Budget {
		pack.Exceeded = ExceededNamed
	}
	pack.addDigestWithin(l.answers(), pack.Exceeded != "", pack.Budget/ReasonDigestShare)
	pack.size()
	return pack, nil
}

// reasonLoader loads into l what one reason needs (ADR-0217 §3).
type reasonLoader func(l *reasonLoad, r Reason) error

// reasonLoaders is the loader of each kind of reason: those that load
// documents here, and those that load the figures flai works out, the last
// verify result, release_evaluate's figures, and the board and the
// promotion candidates, in reasons_figures.go.
var reasonLoaders = map[string]reasonLoader{
	ReasonThread:  loadThread,
	ReasonReview:  loadReview,
	ReasonPlanned: loadPlanned,
	ReasonDraft:   loadDraft,
	ReasonRelease: loadRelease,
	ReasonReady:   loadReady,
}

// reasonLoad is what the reasons of one pack have loaded so far: the items,
// each once, in the order the reasons loaded them, and the threads the
// digest is ranked against.
type reasonLoad struct {
	repo    *workitem.Repo
	cmd     string
	items   []Item
	threads []*threads.Thread
	all     []*workitem.Item // the items outside the archive, read once when a loader first needs them
	trimmed map[string]bool  // the paths of the items loaded trimmed, as a thread's item is (ADR-0218 §1)
}

// add puts it in the pack for reason, or, when an earlier reason loaded the
// same file, adds reason to that item's.
func (l *reasonLoad) add(it Item, reason Reason) {
	l.put(it, reason, false)
}

// put is add for an item loaded whole or, with trimmed, cut as a thread's
// item is (ADR-0218 §1). The pack holds a file once, whole when any of its
// reasons loads it whole: a whole item an earlier reason loaded trimmed
// takes the trimmed one's text, in its place, and a trimmed one an earlier
// reason loaded whole adds its reason alone.
func (l *reasonLoad) put(it Item, reason Reason, trimmed bool) {
	if i := slices.IndexFunc(l.items, func(o Item) bool { return o.Path == it.Path }); i >= 0 {
		if l.trimmed[it.Path] && !trimmed {
			l.items[i].Text, l.items[i].Size = it.Text, len(it.Text)
			delete(l.trimmed, it.Path)
		}
		if l.items[i].Reason != reason.String() {
			l.items[i].Also = appendNew(l.items[i].Also, reason.String())
		}
		return
	}
	if trimmed {
		if l.trimmed == nil {
			l.trimmed = map[string]bool{}
		}
		l.trimmed[it.Path] = true
	}
	it.Step, it.Reason, it.Size = StepReason, reason.String(), len(it.Text)
	l.items = append(l.items, it)
}

// file is the file at path whole, as an item with its path relative to the
// main checkout, which holds the work items and the threads.
func (l *reasonLoad) file(path, id, title string, reason Reason) (Item, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Item{}, fmt.Errorf("%s --reason %s: reading %s: %w", l.cmd, reason, path, err)
	}
	return Item{Path: l.rel(path), ID: id, Title: title, Text: string(data)}, nil
}

// rel is path relative to the main checkout, which holds the work items and
// the threads, with forward slashes, as a reason's item names it.
func (l *reasonLoad) rel(path string) string {
	root := l.repo.MainRoot
	if root == "" {
		root = l.repo.Root
	}
	return filepath.ToSlash(rel(root, path))
}

// threadItemSections is the sections of a work item a thread on it loads,
// in this order (ADR-0218 §1).
var threadItemSections = []string{"Goal", "Acceptance criteria"}

// threadItem loads the work item id as a thread on it needs it (ADR-0218
// §1): trimmed, as trimItem cuts it, under the item's own path, so that a
// reason that loads it whole takes its place.
func (l *reasonLoad) threadItem(id string, reason Reason) error {
	it, err := l.repo.Get(id)
	if err != nil {
		return fmt.Errorf("%s --reason %s: %w; give a reason for an item that exists", l.cmd, reason, err)
	}
	l.put(Item{Path: l.rel(it.Path), ID: it.ID, Title: it.Title, Text: trimItem(it)}, reason, true)
	return nil
}

// trimItem is it as a thread on it loads it (ADR-0218 §1): a front matter
// of its ID, title, type, and status; its Goal and Acceptance criteria
// sections, those of them it has, each under its heading with what is under
// it verbatim; and a line saying what is left out and that item_get gives
// it whole.
func trimItem(it *workitem.Item) string {
	var b strings.Builder
	fmt.Fprintf(&b, "---\nid: %s\ntitle: %s\ntype: %s\nstatus: %s\n---\n", it.ID, strconv.Quote(it.Title), it.Type, it.Status)
	for _, name := range threadItemSections {
		s, ok := workitem.NarrativeSection(it.Body, name)
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "\n## %s\n", name)
		if s.Text != "" {
			fmt.Fprintf(&b, "\n%s\n", s.Text)
		}
	}
	fmt.Fprintf(&b, "\n%s's notes, tasks, and planning are left out of this pack (ADR-0218); item_get gives it whole.\n", it.ID)
	return b.String()
}

// item loads the work item id whole for reason.
func (l *reasonLoad) item(id string, reason Reason) (*workitem.Item, error) {
	it, err := l.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("%s --reason %s: %w; give a reason for an item that exists", l.cmd, reason, err)
	}
	if err := l.workItem(it, reason); err != nil {
		return nil, err
	}
	return it, nil
}

// workItem loads it whole for reason.
func (l *reasonLoad) workItem(it *workitem.Item, reason Reason) error {
	loaded, err := l.file(it.Path, it.ID, it.Title, reason)
	if err != nil {
		return err
	}
	l.add(loaded, reason)
	return nil
}

// thread loads th whole, its entries with it, for reason, and ranks the
// digest against it.
func (l *reasonLoad) thread(th *threads.Thread, reason Reason) error {
	loaded, err := l.file(th.Path, th.ID, th.Title, reason)
	if err != nil {
		return err
	}
	l.add(loaded, reason)
	if !slices.ContainsFunc(l.threads, func(o *threads.Thread) bool { return o.ID == th.ID }) {
		l.threads = append(l.threads, th)
	}
	return nil
}

// loadThread loads a thread reason: the thread with its entries and, when
// it is on a work item, that item's ID, title, type, status, goal, and
// acceptance criteria, as threadItem loads it (ADR-0218 §1). A thread on a
// document names its path in its anchor and its first line, which is all
// the pack gives of it; doc_get reads the document.
func loadThread(l *reasonLoad, r Reason) error {
	th, err := threads.Get(l.repo, r.ID)
	if err != nil {
		return fmt.Errorf("%s --reason %s: %w; give a thread that exists", l.cmd, r, err)
	}
	if err := l.thread(th, r); err != nil {
		return err
	}
	if th.Anchor.Item == "" {
		return nil
	}
	return l.threadItem(th.Anchor.Item, r)
}

// loadDraft loads a draft reason: the story whole.
func loadDraft(l *reasonLoad, r Reason) error {
	_, err := l.item(r.ID, r)
	return err
}

// loadPlanned loads a planned reason: the item planned whole; the items
// under it that are not cancelled, an epic's stories or a story's tasks,
// each whole, as the planner wrote them; and each open thread on it that
// the planner opened, with its entries, which the digest is ranked against
// too.
func loadPlanned(l *reasonLoad, r Reason) error {
	it, err := l.item(r.ID, r)
	if err != nil {
		return err
	}
	if l.all == nil {
		if l.all, err = l.repo.List(false); err != nil {
			return fmt.Errorf("%s --reason %s: listing the items under %s: %w", l.cmd, r, it.ID, err)
		}
	}
	for _, child := range workitem.Children(l.all, it.ID) {
		if child.Status == workitem.Cancelled {
			continue
		}
		if err := l.workItem(child, r); err != nil {
			return err
		}
	}
	on, err := threads.For(l.repo, it.ID)
	if err != nil {
		return fmt.Errorf("%s --reason %s: listing the threads on %s: %w", l.cmd, r, it.ID, err)
	}
	for _, th := range on {
		if th.Open() && workitem.IsPlanner("", th.Opener()) {
			if err := l.thread(th, r); err != nil {
				return err
			}
		}
	}
	return nil
}

// answers is the digest items the pack tries, as answersFor ranks the
// digest against the reasons' threads. A digest entry that does not parse
// is flai check's to report; the pack does without the digest rather than
// fail, as the general pack does.
func (l *reasonLoad) answers() []Item {
	if len(l.threads) == 0 {
		return nil
	}
	entries, err := digest.List(l.repo)
	if err != nil {
		return nil
	}
	return answersFor(entries, l.threads)
}

// reasonsHead is groupHead for a pack primed for reasons, whose digest is
// ranked against the reasons' threads rather than those awaiting the
// operator.
var reasonsHead = func() map[string]string {
	m := maps.Clone(strategicHead)
	m[groupDigest] = "\ndigest\n======\n\nThe answers and recommendations of the digest that rank highest against each thread the reasons load (reason: thread TH-nnnn; rank n), within a sixteenth of the budget (ADR-0194, ADR-0218): each with its body, or as one line when it is compacted. They are what the operator answered and confirmed before, not the design. When one answers a thread, cite its D-nnnn and give its path as a source on thread_reply; ask the operator only what no entry answers. Read one whole with the MCP doc_get and its path, or flai doc show <path>.\n"
	return m
}()

// reasonLines is the header's line naming a pack primed for reasons and
// what it holds.
func (p *Pack) reasonLines() string {
	return fmt.Sprintf("role: %s, %s, primed for %s (ADR-0217): the conventions the role reads with the sections those reasons do not need left out, what each reason loads, and the answers and recommendations of the digest that rank highest against the reasons' threads; doc_search and doc_get find the rest.\n", p.Role, p.agent, strings.Join(p.Reasons, ", "))
}
