package workitem

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/filelock"
)

// itemsLockWait is how long Save waits for another writer of a work item to
// finish with the items lock. A holder holds it for one read, merge, and
// write of one file, so the wait is far longer than any holder needs; the
// lock is the operating system's, freed when its holder dies, and never
// broken for its age (I-0140).
const itemsLockWait = 30 * time.Second

// ConflictError is a Save refused because another writer saved the item
// since it was read, changing parts of it this copy changes too (I-0175).
type ConflictError struct {
	ID string
	// Keys are the parts both changed: front-matter keys, "body" for what
	// is above the body's first section, and "body section ## Heading".
	Keys []string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("save %s: another writer changed it since it was read, and both changed %s, so nothing was written: read %s again and make the change again",
		e.ID, strings.Join(e.Keys, ", "), e.ID)
}

// lockItems takes the one lock every work item write is made under, an
// operating system lock on a file in .flai-cache, so that two writers do not
// each read a file, then write it, the second over the first (I-0175).
func (r *Repo) lockItems() (unlock func(), err error) {
	return filelock.Lock(filepath.Join(r.CacheDir(), "items.lock"), itemsLockWait)
}

// mergeSaved is the document Save writes for it, under the items lock: it as
// it is when it was made in memory or its file is missing, does not parse,
// or is still what it was read as; otherwise it merged, part by part, with
// what another writer saved since, and it is then made the merged item. A
// part both changed differently, save updated, is a *ConflictError.
func mergeSaved(it *Item) (string, error) {
	mine := it.Marshal()
	if it.readAs == "" {
		return mine, nil
	}
	data, err := os.ReadFile(it.Path)
	if os.IsNotExist(err) {
		return mine, nil
	}
	if err != nil {
		return "", fmt.Errorf("save %s: read %s for what another writer saved since: %w", it.ID, it.Path, err)
	}
	// a file that does not parse holds no change to keep, and is written
	// over as before
	theirs, ok := canonical(string(data))
	if !ok || theirs == it.readAs || theirs == mine {
		return mine, nil
	}
	doc, both, err := merge3(it.readAs, mine, theirs)
	if err != nil {
		return "", fmt.Errorf("save %s: merge with what another writer saved: %w", it.ID, err)
	}
	if len(both) > 0 {
		return "", &ConflictError{ID: it.ID, Keys: both}
	}
	merged, err := ParseItem(doc)
	if err != nil {
		return "", fmt.Errorf("save %s: the merge with what another writer saved does not parse, so nothing was written (%w): read %s again and make the change again", it.ID, err, it.ID)
	}
	merged.Path, merged.Archived = it.Path, it.Archived
	*it = *merged
	return it.Marshal(), nil
}

// Restore takes a write out of the work item file at path, under the items
// lock: wrote is what the write left there and before what it replaced. The
// file is given before when no other writer changed it since, or it is
// missing or one of the three does not parse; otherwise only the write's own
// change is taken out, part by part, keeping the other writers', and a part
// they changed too is a *ConflictError, with nothing written (S-0479, I-0175).
func (r *Repo) Restore(path string, wrote, before []byte) error {
	unlock, err := r.lockItems()
	if err != nil {
		return fmt.Errorf("restore %s: %w", path, err)
	}
	defer unlock()
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("restore %s: read it for what another writer saved since: %w", path, err)
	}
	base, okBase := canonical(string(wrote))
	mine, okMine := canonical(string(before))
	cur, perr := ParseItem(string(data))
	// a file no other writer changed since, or one of the three that does
	// not parse, holds no change to keep
	if err != nil || perr != nil || !okBase || !okMine || bytes.Equal(data, wrote) || cur.Marshal() == base {
		return atomicfile.WriteFile(path, before, 0o644)
	}
	doc, both, err := merge3(base, mine, cur.Marshal())
	if err != nil {
		return fmt.Errorf("restore %s: merge with what another writer saved: %w", cur.ID, err)
	}
	if len(both) > 0 {
		return &ConflictError{ID: cur.ID, Keys: both}
	}
	merged, err := ParseItem(doc)
	if err != nil {
		return fmt.Errorf("restore %s: the merge with what another writer saved does not parse, so nothing was written (%w): read %s and take the change out by hand", cur.ID, err, cur.ID)
	}
	return atomicfile.WriteFile(path, []byte(merged.Marshal()), 0o644)
}

// canonical is doc as Marshal writes it, and whether it parses.
func canonical(doc string) (string, bool) {
	it, err := ParseItem(doc)
	if err != nil {
		return "", false
	}
	return it.Marshal(), true
}

// merge3 merges mine and theirs, two documents changed from base, part by
// part: each top-level front-matter key, and each section of the body, so
// that one writer's note and another's new child line are both kept. The
// parts both changed differently are returned by name, and the document
// then means nothing.
func merge3(base, mine, theirs string) (string, []string, error) {
	var fms, bodies [3]parts
	for i, d := range []string{base, mine, theirs} {
		fm, body, err := SplitFrontMatter(d)
		if err != nil {
			return "", nil, err
		}
		fms[i], bodies[i] = frontParts(fm), bodyParts(body)
	}
	fm, both := mergeParts(fms[0], fms[1], fms[2], func(k string) string { return k })
	body, bodyBoth := mergeParts(bodies[0], bodies[1], bodies[2], func(k string) string {
		if k == "" {
			return "body"
		}
		heading, _, _ := strings.Cut(k, "\x00")
		return "body section " + heading
	})
	return "---\n" + fm + "---\n" + body, append(both, bodyBoth...), nil
}

// mergeParts joins base's parts as mine and theirs changed them, in mine's
// order with theirs' new parts after the part they follow there, and returns
// the names of the parts both changed differently. A part one left as in
// base is the other's, and a part both changed alike is that; updated, both
// changed, is the later.
func mergeParts(base, mine, theirs parts, name func(key string) string) (string, []string) {
	order := slices.Clone(mine.order)
	for i, k := range theirs.order {
		if slices.Contains(order, k) {
			continue
		}
		at := 0
		if i > 0 {
			at = slices.Index(order, theirs.order[i-1]) + 1
		}
		order = slices.Insert(order, at, k)
	}
	var b strings.Builder
	var both []string
	for _, k := range order {
		bp, mp, tp := base.text[k], mine.text[k], theirs.text[k]
		switch {
		case mp == bp:
			b.WriteString(tp)
		case tp == bp, mp == tp:
			b.WriteString(mp)
		case k == "updated":
			// "updated: " and a TimeFormat stamp, which sorts as it reads
			b.WriteString(max(mp, tp))
		default:
			both = append(both, name(k))
		}
	}
	return b.String(), both
}

// parts is a document cut into named parts, each with its lines, in the
// order written; a part a document lacks is empty.
type parts struct {
	order []string
	text  map[string]string
}

func (p *parts) add(key, line string) {
	if p.text == nil {
		p.text = map[string]string{}
	}
	if _, seen := p.text[key]; !seen {
		p.order = append(p.order, key)
	}
	p.text[key] += line
}

// frontParts cuts front matter as UnknownFields does: a line that starts a
// top-level key starts its part, named by the key, and every other line,
// indented or not, belongs to the key above.
func frontParts(fm string) parts {
	var p parts
	key := ""
	for _, line := range strings.SplitAfter(fm, "\n") {
		if line == "" {
			continue
		}
		if m := topKey.FindStringSubmatch(line); m != nil {
			key = m[1]
		}
		p.add(key, line)
	}
	return p
}

// bodyParts cuts a body at its "## " section headings: what is above the
// first is the part named "", and each section is a part named by its
// heading and how many sections above it have that heading.
func bodyParts(body string) parts {
	var p parts
	key := ""
	seen := map[string]int{}
	for _, line := range strings.SplitAfter(body, "\n") {
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "## ") {
			heading := strings.TrimRight(line, "\n")
			key = fmt.Sprintf("%s\x00%d", heading, seen[heading])
			seen[heading]++
		}
		p.add(key, line)
	}
	return p
}
