package issues

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A rebase that stops on an issue file both sides changed merges it by its
// instances (ADR-0126): the instances both sides added are kept, the count and
// cost add up, and every other part takes the side that changed it.

// ErrNotMerged is returned, wrapped with the reason, when an issue file's
// versions do not merge without a person's judgment.
var ErrNotMerged = errors.New("issue file not merged")

// Merge merges three versions of one issue file's text, base and the two sides
// that changed it, ours and theirs, by their instances, and returns the merged
// text; it returns an error wrapping ErrNotMerged when they do not merge.
func Merge(base, ours, theirs string) (string, error) {
	if strings.TrimSpace(base) == "" {
		return "", fmt.Errorf("%w: there is no base version to merge the two sides from; merge the file by hand", ErrNotMerged)
	}
	var v [3]*Issue
	var p [3]instances
	for i, side := range [3]struct{ name, text string }{{"base", base}, {"ours", ours}, {"theirs", theirs}} {
		is, err := Parse(side.text)
		if err == nil {
			err = is.Validate()
		}
		if err != nil {
			return "", fmt.Errorf("%w: the %s version does not read as an issue (%w); merge the file by hand", ErrNotMerged, side.name, err)
		}
		parts, ok := splitInstances(is.Body)
		if !ok {
			return "", fmt.Errorf("%w: the %s version of %s has no ## Instances section; merge the file by hand", ErrNotMerged, side.name, is.ID)
		}
		v[i], p[i] = is, parts
	}
	b, o, t := v[0], v[1], v[2]

	var conflicts []string
	str := func(name, inBase, inOurs, inTheirs string) string {
		m, ok := pick(inBase, inOurs, inTheirs)
		if !ok {
			conflicts = append(conflicts, name)
		}
		return m
	}
	m := &Issue{
		ID:            str("id", b.ID, o.ID, t.ID),
		Title:         str("title", b.Title, o.Title, t.Title),
		Class:         str("class", b.Class, o.Class, t.Class),
		Status:        str("status", b.Status, o.Status, t.Status),
		FirstReported: str("first_reported", b.FirstReported, o.FirstReported, t.FirstReported),
		LastReported:  later(o.LastReported, t.LastReported),
		Updated:       later(o.Updated, t.Updated),
		Usage:         o.Usage,
	}
	if str("usage", usageBlock(b.Usage), usageBlock(o.Usage), usageBlock(t.Usage)) != usageBlock(o.Usage) {
		m.Usage = t.Usage
	}
	m.Unknown = mergeUnknown(b.Unknown, o.Unknown, t.Unknown, func(name string) { conflicts = append(conflicts, name) })
	parts := instances{
		head: str("the body before its instances", p[0].head, p[1].head, p[2].head),
		tail: str("the body after its Instances section", p[0].tail, p[1].tail, p[2].tail),
	}
	oursAdded, oursKept := addedTo(p[0].blocks, p[1].blocks)
	theirsAdded, theirsKept := addedTo(p[0].blocks, p[2].blocks)
	if !oursKept || !theirsKept {
		conflicts = append(conflicts, "an instance the base recorded")
	}
	if o.Count < b.Count || t.Count < b.Count {
		conflicts = append(conflicts, "the count, lower on a side than the base's")
	}
	if len(conflicts) > 0 {
		return "", fmt.Errorf("%w: %s: the two sides changed %s differently; merge the file by hand", ErrNotMerged, b.ID, joinAnd(conflicts))
	}

	blocks, _ := addInstances(p[0].blocks, oursAdded)
	blocks, held := addInstances(blocks, theirsAdded)
	parts.blocks = inOrder(blocks)
	m.Body = parts.body()
	n := mergedTally(tallyOf(b), tallyOf(o), tallyOf(t), held)
	m.Count, m.Cost = n.count, n.cost()
	if err := m.Validate(); err != nil {
		return "", fmt.Errorf("%w: %s: the merge is not a valid issue (%w); merge the file by hand", ErrNotMerged, b.ID, err)
	}
	return m.Marshal(), nil
}

// addOccurrences adds from's occurrences to into, as a merge adds a side's:
// from's instances that into does not hold already, in timestamp order, its
// count less those it held already, and its cost, weighted by count (see
// tally.plus); first_reported takes the earlier of the two, and last_reported
// and updated the later. Every other field and the rest of into's body stay
// as they are. Both bodies must have an Instances section.
func addOccurrences(into, from *Issue) error {
	ip, ok := splitInstances(into.Body)
	if !ok {
		return fmt.Errorf("%s has no ## Instances section to add %s's occurrences to; add the section to %s by hand", into.ID, from.ID, into.ID)
	}
	fp, ok := splitInstances(from.Body)
	if !ok {
		return fmt.Errorf("%s has no ## Instances section to take its occurrences from; add the section to %s by hand", from.ID, from.ID)
	}
	blocks, held := addInstances(ip.blocks, fp.blocks)
	ip.blocks = inOrder(blocks)
	n := tallyOf(into).plus(tallyOf(from).less(held))
	into.Body = ip.body()
	into.Count, into.Cost = n.count, n.cost()
	into.FirstReported = earlier(into.FirstReported, from.FirstReported)
	into.LastReported = later(into.LastReported, from.LastReported)
	into.Updated = later(into.Updated, from.Updated)
	return nil
}

// instance is one occurrence's block of an Instances section: its "### <ts>"
// heading and the text under it, trimmed.
type instance struct {
	heading, text string
}

// instances is an issue body split at its Instances section: the body up to
// the first instance, the instances, and the body from the section after it.
type instances struct {
	head   string
	blocks []instance
	tail   string
}

// splitInstances splits body at its Instances section, or reports false when
// it has none. Text in the section before its first instance stays in head.
func splitInstances(body string) (instances, bool) {
	start, end, ok := section(body, "## Instances")
	if !ok {
		return instances{}, false
	}
	p := instances{head: body[:start], tail: body[end:]}
	pieces := strings.Split("\n"+body[start:end], "\n### ")
	if pre := strings.TrimSpace(pieces[0]); pre != "" {
		p.head += "\n" + pre + "\n"
	}
	for _, piece := range pieces[1:] {
		heading, text, _ := strings.Cut(piece, "\n")
		p.blocks = append(p.blocks, instance{heading: "### " + strings.TrimSpace(heading), text: strings.TrimSpace(text)})
	}
	return p, true
}

// body joins the parts back into a body, laid out as insertInstance writes
// it: a blank line before, between, and after the instances.
func (p instances) body() string {
	if len(p.blocks) == 0 {
		return p.head + "\n" + p.tail
	}
	var b strings.Builder
	b.WriteString(p.head)
	for _, in := range p.blocks {
		b.WriteString("\n" + in.heading + "\n")
		if in.text != "" {
			b.WriteString(in.text + "\n")
		}
	}
	if p.tail != "" {
		b.WriteString("\n" + p.tail)
	}
	return b.String()
}

// addedTo returns the instances of side that base does not hold, compared by
// heading and text, and whether side holds every instance base does.
func addedTo(base, side []instance) ([]instance, bool) {
	inSide := map[instance]bool{}
	for _, in := range side {
		inSide[in] = true
	}
	kept := true
	inBase := map[instance]bool{}
	for _, in := range base {
		inBase[in] = true
		kept = kept && inSide[in]
	}
	var added []instance
	for _, in := range side {
		if !inBase[in] {
			added = append(added, in)
		}
	}
	return added, kept
}

// addInstances appends to into each of more that it does not hold already,
// compared by heading and text, and returns the result with how many of more
// it held already.
func addInstances(into, more []instance) ([]instance, int) {
	out := append([]instance(nil), into...)
	have := map[instance]bool{}
	for _, in := range out {
		have[in] = true
	}
	held := 0
	for _, in := range more {
		if have[in] {
			held++
			continue
		}
		have[in] = true
		out = append(out, in)
	}
	return out, held
}

// inOrder sorts the instances by the timestamps of their headings, keeping
// the order of those of one timestamp, and joins those under one heading, as
// insertInstance does, since sibling headings must differ: the joined text
// leaves out a story line the instance names already.
func inOrder(blocks []instance) []instance {
	sorted := append([]instance(nil), blocks...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].heading < sorted[j].heading })
	var out []instance
	for _, in := range sorted {
		if len(out) == 0 || out[len(out)-1].heading != in.heading {
			out = append(out, in)
			continue
		}
		last := &out[len(out)-1]
		text := in.text
		if first, rest, _ := strings.Cut(text, "\n"); storyLineRe.MatchString(first) && contains(strings.Split(last.text, "\n"), first) {
			text = rest
		}
		last.text = strings.TrimSpace(last.text + "\n\n" + strings.TrimSpace(text))
	}
	return out
}

// tally is a count of occurrences and their total cost, which is not known
// when no occurrence's cost was recorded.
type tally struct {
	count int
	total time.Duration
	known bool
}

// tallyOf is the issue's occurrences: its count, at its average cost.
func tallyOf(is *Issue) tally {
	d, err := time.ParseDuration(is.Cost)
	if is.Cost == "" || err != nil {
		return tally{count: is.Count}
	}
	return tally{count: is.Count, total: d * time.Duration(is.Count), known: true}
}

// average is the cost per occurrence, zero for none.
func (t tally) average() time.Duration {
	if t.count < 1 {
		return 0
	}
	return t.total / time.Duration(t.count)
}

// at is t's occurrences, each at cost avg.
func (t tally) at(avg time.Duration) tally {
	return tally{count: t.count, total: avg * time.Duration(t.count), known: true}
}

// plus is t's occurrences and u's. When only one of them knows its cost, the
// other's occurrences are taken at its average, as BumpWith takes an issue's
// earlier occurrences at the first cost recorded and an occurrence with no
// cost at the issue's average.
func (t tally) plus(u tally) tally {
	n := tally{count: t.count + u.count}
	switch {
	case u.count == 0:
		return t
	case t.count == 0:
		return u
	case t.known && u.known:
		return tally{count: n.count, total: t.total + u.total, known: true}
	case t.known:
		return n.at(t.average())
	case u.known:
		return n.at(u.average())
	}
	return n
}

// less is t less k of its occurrences, at most all of them, taken at its
// average: those recorded twice.
func (t tally) less(k int) tally {
	k = min(k, t.count)
	if !t.known {
		return tally{count: t.count - k}
	}
	return tally{count: t.count - k, total: t.total - t.average()*time.Duration(k), known: true}
}

// over is the occurrences t has over base: the count and the cost added.
func (t tally) over(base tally) tally {
	return tally{count: t.count - base.count, total: t.total - base.total, known: t.known && base.known}
}

// cost is the average cost as an issue records it, rounded to the minute and
// normalised as BumpWith writes it, or "" when it is not known.
func (t tally) cost() string {
	if !t.known || t.count < 1 {
		return ""
	}
	return normalise(t.average().Round(time.Minute).String())
}

// mergedTally is the occurrences of a merge: the base's, and those each side
// added over it, the cost a side added being its total less the base's, less
// the held instances theirs added that ours had added too, taken at the
// average of what theirs added. A side whose cost is not known is taken at
// the base's average, as an occurrence with no cost leaves an issue's average
// as it is; a base whose cost is not known is taken at the average of the
// sides that know theirs, weighted by count, as plus does; when none of the
// three knows, the cost is not known. Neither side's count may be below the
// base's.
func mergedTally(base, ours, theirs tally, held int) tally {
	if !base.known {
		sides := tally{}
		for _, s := range []tally{ours, theirs} {
			if s.known {
				sides = sides.plus(s)
			}
		}
		if sides.known {
			base = base.at(sides.average())
		}
	}
	if base.known {
		if !ours.known {
			ours = ours.at(base.average())
		}
		if !theirs.known {
			theirs = theirs.at(base.average())
		}
	}
	fromOurs, fromTheirs := ours.over(base), theirs.over(base).less(held)
	return tally{count: base.count + fromOurs.count + fromTheirs.count, total: max(base.total+fromOurs.total+fromTheirs.total, 0), known: base.known}
}

// mergeUnknown merges the front matter fields this flai does not know, each
// by name as pick does, in the order ours, theirs, then the base wrote them,
// calling conflict with the name of each both sides changed differently. A
// field a side removed and the other left as it was is removed.
func mergeUnknown(base, ours, theirs []workitem.Field, conflict func(name string)) []workitem.Field {
	raw := func(fields []workitem.Field) map[string]string {
		m := map[string]string{}
		for _, f := range fields {
			m[f.Name] = f.Raw
		}
		return m
	}
	b, o, t := raw(base), raw(ours), raw(theirs)
	var out []workitem.Field
	seen := map[string]bool{}
	for _, fields := range [][]workitem.Field{ours, theirs, base} {
		for _, f := range fields {
			if seen[f.Name] {
				continue
			}
			seen[f.Name] = true
			m, ok := pick(b[f.Name], o[f.Name], t[f.Name])
			if !ok {
				conflict(f.Name)
			}
			if m != "" {
				out = append(out, workitem.Field{Name: f.Name, Raw: m})
			}
		}
	}
	return out
}

// pick merges one part three ways: the sides' value when they agree, the side
// that changed it when only one did, and false when both changed it
// differently.
func pick[T comparable](base, ours, theirs T) (T, bool) {
	switch {
	case ours == theirs, theirs == base:
		return ours, true
	case ours == base:
		return theirs, true
	}
	return ours, false
}

// earlier is the earlier of two UTC timestamps.
func earlier(a, b string) string {
	if b < a {
		return b
	}
	return a
}

// later is the later of two UTC timestamps.
func later(a, b string) string {
	if b > a {
		return b
	}
	return a
}
