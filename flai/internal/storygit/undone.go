package storygit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/bytepunx/system-flow/flai/internal/atomicfile"
	"github.com/bytepunx/system-flow/flai/internal/execx"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// What a rebased story branch undoes of the main branch (ADR-0155, I-0124):
// a branch that changed a path and reverted it, rebased onto a main that made
// the same change since, loses its first commit, now empty, and replays the
// revert cleanly over main's change. No conflict shows, so after a clean
// rebase the sync compares, for each path main changed since the branch's old
// base and the branch still changes against main, the lines main added with
// the lines the branch's net change removes, by their position in main's file
// (ADR-0309): a line of main's is undone only when the branch removes the line
// of the number main added it at, not a line of the same text elsewhere
// (I-0385), and a markdown line a line the branch adds keeps is not (ADR-0297).
// Each path where the branch removes any of them is recorded per story in
// .flai-cache/undone/S-nnnn.json in the main checkout, beside flai verify's
// last result, because a later rebase starts past main's change and can no
// longer find it. The record is host-local.

// Undone is a path where a rebased story branch removes lines the main branch
// added since the branch's old base.
type Undone struct {
	Path    string   `json:"path"`     // as git names it
	Lines   []string `json:"lines"`    // main's added lines the branch removes, as text
	OldBase string   `json:"old_base"` // the branch's merge base with main before the rebase
}

// UndoneEntry is an Undone in a story's record, with when a sync found it.
type UndoneEntry struct {
	Undone
	FoundAt time.Time `json:"found_at"` // UTC, to the second
}

// UndoneRecord is a story's record of the paths its rebased branch undoes.
type UndoneRecord struct {
	Story string        `json:"story"`
	Paths []UndoneEntry `json:"paths"` // sorted by path
}

// UndoneRecordPath is where story's record of undone paths is in the cache
// folder cacheDir.
func UndoneRecordPath(cacheDir, story string) string {
	return filepath.Join(cacheDir, "undone", story+".json")
}

// FindUndone is the paths where the branch checked out in the worktree at dir,
// rebased onto onto, the main branch's commit the rebase took, removes lines
// main added between oldBase, the branch's merge base with main before the
// rebase, and onto, each at the line main added it in onto's file
// (ADR-0309), in the order of their numbers. A commit main gains after onto
// is not counted, so a line only it adds is none the branch removes (I-0338).
// It finds none when oldBase is onto, since main then added nothing the
// branch did not already have.
func FindUndone(r execx.Runner, dir, onto, oldBase string) ([]Undone, error) {
	found := []Undone{}
	if onto == oldBase {
		return found, nil
	}
	mains, err := nameOnly(r, dir, nameOnlyArgs(oldBase, onto)...)
	if err != nil {
		return nil, fmt.Errorf("list the paths the main branch changed between %s and %s: %w", oldBase, onto, err)
	}
	branch, err := nameOnly(r, dir, nameOnlyArgs(onto, "HEAD")...)
	if err != nil {
		return nil, fmt.Errorf("list the paths the branch in %s changes against %s: %w", dir, onto, err)
	}
	for _, p := range mains {
		if !slices.Contains(branch, p) {
			continue
		}
		added, mainRemoved, err := readHunks(r, dir, oldBase, onto, p)
		if err != nil {
			return nil, fmt.Errorf("read what the main branch added to %s between %s and %s: %w", p, oldBase, onto, err)
		}
		branchAdded, removed, err := readHunks(r, dir, onto, "HEAD", p)
		if err != nil {
			return nil, fmt.Errorf("read what the branch in %s removes from %s against %s: %w", dir, p, onto, err)
		}
		if lines := notKept(p, atPositions(added, removed), texts(branchAdded), texts(mainRemoved)); len(lines) > 0 {
			found = append(found, Undone{Path: p, Lines: lines, OldBase: oldBase})
		}
	}
	return found, nil
}

// RecordedUndone is the paths story's record holds that its branch still
// undoes, each checked against the main branch now and given with the
// recorded lines the branch still removes at the line main added them in
// main's file now (ADR-0309), counted as a multiset of the record's text, so
// a line the sync did not record is not counted (I-0338). A path whose branch
// removes none of its recorded lines there, as after it is restored or when
// the text compare before ADR-0309 recorded it, drops out, and the record is
// written again without it. A story with no record has none. It runs one git
// diff of the branch for all the record's paths, and one per distinct old
// base for the paths the branch still removes lines of, not one or two per
// path (I-0349).
func RecordedUndone(r execx.Runner, repo *workitem.Repo, story string) ([]Undone, error) {
	cacheDir := repo.CacheDir()
	rec, err := readUndone(cacheDir, story)
	if err != nil {
		return nil, err
	}
	out := []Undone{}
	if len(rec.Paths) == 0 {
		return out, nil
	}
	main, err := MainBranch(r, repo.MainRoot)
	if err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(rec.Paths))
	for _, e := range rec.Paths {
		paths = append(paths, e.Path)
	}
	branch, err := readHunksOf(r, repo.MainRoot, main, Branch(story), paths)
	if err != nil {
		return nil, fmt.Errorf("check whether %s still undoes %s's change to %s: %w; a story with no branch has nothing to undo, so delete %s", Branch(story), main, strings.Join(paths, ", "), err, UndoneRecordPath(cacheDir, story))
	}
	// the paths the branch still removes lines of, by the old base each was
	// found at, each base read once
	byBase := map[string][]string{}
	for _, e := range rec.Paths {
		if len(branch[e.Path].removed) > 0 {
			byBase[e.OldBase] = append(byBase[e.OldBase], e.Path)
		}
	}
	mains := make(map[string]map[string]hunks, len(byBase))
	for base, ps := range byBase {
		if mains[base], err = readHunksOf(r, repo.MainRoot, base, main, ps); err != nil {
			return nil, fmt.Errorf("read what %s changed in %s since %s: %w; delete %s and run flai stream sync %s again", main, strings.Join(ps, ", "), base, err, UndoneRecordPath(cacheDir, story), story)
		}
	}
	kept := []UndoneEntry{}
	for _, e := range rec.Paths {
		b := branch[e.Path]
		if len(b.removed) == 0 {
			continue
		}
		m := mains[e.OldBase][e.Path]
		lines := undoneLines(e.Lines, atPositions(m.added, b.removed))
		if e.Lines = notKept(e.Path, lines, texts(b.added), texts(m.removed)); len(e.Lines) == 0 {
			continue
		}
		kept = append(kept, e)
		out = append(out, e.Undone)
	}
	if len(kept) != len(rec.Paths) {
		rec.Paths = kept
		if err := writeUndone(cacheDir, rec); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// withRecordedUndone sets res.UndoesMain, after a clean sync, to every path
// the story's branch still undoes of the main branch's change: what this
// rebase found, and what the story's record holds besides, one entry per
// path, sorted by path. A later rebase starts past main's change and finds
// nothing itself, so the record carries the warning on; Rebase alone, as
// flai accept runs it, gives only its own findings. A record that cannot be
// read is logged and the sync stands.
func withRecordedUndone(o SyncOptions, res *SyncResult) {
	recorded, err := RecordedUndone(o.Runner, o.Repo, o.Story.ID)
	if err != nil {
		o.Log.Warn("the paths the story's record says its branch undoes of the main branch not read; the sync stands", "component", "git", "story", o.Story.ID, "err", err)
		return
	}
	found := len(res.UndoesMain)
	for _, u := range recorded {
		if !slices.ContainsFunc(res.UndoesMain[:found], func(f Undone) bool { return f.Path == u.Path }) {
			res.UndoesMain = append(res.UndoesMain, u)
		}
	}
	slices.SortFunc(res.UndoesMain, func(a, b Undone) int { return strings.Compare(a.Path, b.Path) })
}

// checkUndone sets res.UndoesMain to what the branch rebased in the worktree
// at dir onto onto, the main branch's commit the rebase took, undoes of what
// main changed between oldBase and onto, and adds it to the story's record,
// found at o.Now, leaving out the files flai writes itself. A commit main
// gains during the sync is not counted (I-0338). A check or a record that
// fails is logged and the sync stands; an empty oldBase, which could not be
// read, checks nothing.
func checkUndone(o SyncOptions, dir, onto, oldBase string, res *SyncResult) {
	if oldBase == "" {
		return
	}
	id := o.Story.ID
	found, err := FindUndone(o.Runner, dir, onto, oldBase)
	if err != nil {
		o.Log.Warn("what the rebased branch undoes of the main branch not checked; the sync stands", "component", "git", "story", id, "err", err)
		return
	}
	// flai writes these itself from both sides, so main's lines are not lost
	// there: the generated files, and the files the rebase merged (ADR-0098,
	// ADR-0126)
	own := slices.Concat(o.Files.generatedPaths(), res.Regenerated, res.Merged)
	found = slices.DeleteFunc(found, func(u Undone) bool { return slices.Contains(own, u.Path) })
	if len(found) == 0 {
		return
	}
	res.UndoesMain = found
	for _, u := range found {
		o.Log.Warn("the rebased branch undoes the main branch's change to a path", "component", "git", "story", id, "path", u.Path, "lines", len(u.Lines), "old_base", u.OldBase)
	}
	if err := recordUndone(o.Repo.CacheDir(), id, found, o.Now); err != nil {
		o.Log.Warn("undone paths not recorded; the sync stands", "component", "git", "story", id, "err", err)
	}
}

// recordUndone adds found, found at now, to story's record in cacheDir,
// replacing what the record held for the same path.
func recordUndone(cacheDir, story string, found []Undone, now time.Time) error {
	rec, err := readUndone(cacheDir, story)
	if err != nil {
		return err
	}
	at := now.UTC().Truncate(time.Second)
	for _, u := range found {
		rec.Paths = slices.DeleteFunc(rec.Paths, func(e UndoneEntry) bool { return e.Path == u.Path })
		rec.Paths = append(rec.Paths, UndoneEntry{Undone: u, FoundAt: at})
	}
	slices.SortFunc(rec.Paths, func(a, b UndoneEntry) int { return strings.Compare(a.Path, b.Path) })
	return writeUndone(cacheDir, rec)
}

// readUndone is story's record in cacheDir, empty when it has none.
func readUndone(cacheDir, story string) (UndoneRecord, error) {
	rec := UndoneRecord{Story: story, Paths: []UndoneEntry{}}
	path := UndoneRecordPath(cacheDir, story)
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return rec, nil
	}
	if err != nil {
		return rec, fmt.Errorf("read the undone paths of %s at %s: %w", story, path, err)
	}
	if err := json.Unmarshal(data, &rec); err != nil {
		return rec, fmt.Errorf("read the undone paths of %s at %s: %w; delete it and run flai stream sync %s again", story, path, err, story)
	}
	return rec, nil
}

// writeUndone stores rec in cacheDir in one step, or removes it when it holds
// no path.
func writeUndone(cacheDir string, rec UndoneRecord) error {
	path := UndoneRecordPath(cacheDir, rec.Story)
	if len(rec.Paths) == 0 {
		if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("remove the undone paths of %s at %s: %w; check that the project's .flai-cache is writable", rec.Story, path, err)
		}
		return nil
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("store the undone paths of %s: %w", rec.Story, err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("store the undone paths of %s: make %s: %w; check that the project's .flai-cache is writable", rec.Story, filepath.Dir(path), err)
	}
	if err := atomicfile.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("store the undone paths of %s at %s: %w; check that the project's .flai-cache is writable", rec.Story, path, err)
	}
	return nil
}

// nameOnlyArgs are git's arguments to list the paths that differ between
// from and to, unquoted and with no renames, so that each is the path itself.
func nameOnlyArgs(from, to string) []string {
	return []string{"-c", "core.quotePath=false", "diff", "--name-only", "--no-renames", from, to}
}

// diffArgs are git's arguments for the lines of path that differ between
// from and to, with no context.
func diffArgs(from, to, path string) []string {
	return []string{"diff", "--no-ext-diff", "--no-color", "--no-renames", "-U0", from, to, "--", path}
}

// numbered is a line of a file a diff adds or removes, with its number in
// the file it is in.
type numbered struct {
	n    int
	text string
}

// readHunks is the lines of path that git, run in dir, says differ between
// from and to: what to adds, numbered in to's file, and what it removes,
// numbered in from's file.
func readHunks(r execx.Runner, dir, from, to, path string) (added, removed []numbered, err error) {
	out, err := r.Run(dir, "git", diffArgs(from, to, path)...)
	if err != nil {
		return nil, nil, err
	}
	return diffHunks(out)
}

// hunks is what a diff adds to a file, numbered in the new file, and
// removes, numbered in the old one.
type hunks struct {
	added, removed []numbered
}

// diffAllArgs are git's arguments for the lines of every one of paths that
// differ between from and to, with no context, each file under its own
// header diff --git a/<path> b/<path>, whatever the host's diff.noprefix or
// diff.mnemonicPrefix, and with no path quoted that git need not quote.
func diffAllArgs(from, to string, paths []string) []string {
	args := []string{"-c", "core.quotePath=false", "diff", "--no-ext-diff", "--no-color", "--no-renames", "--src-prefix=a/", "--dst-prefix=b/", "-U0", from, to, "--"}
	return append(args, paths...)
}

// readHunksOf is, for each of paths, the lines git, run in dir, says differ
// between from and to, as readHunks gives them, from one git diff of them
// all, split by each file's header (I-0349). A path git prints no diff for
// has none. A path git still quotes, one holding a quote, a backslash, or a
// control character, matches no header, so it is read with a diff of its
// own.
func readHunksOf(r execx.Runner, dir, from, to string, paths []string) (map[string]hunks, error) {
	out, err := r.Run(dir, "git", diffAllArgs(from, to, paths)...)
	if err != nil {
		return nil, err
	}
	header := make(map[string]string, len(paths))
	for _, p := range paths {
		header["diff --git a/"+p+" b/"+p] = p
	}
	files := map[string][]string{}
	current, unknown := "", false
	for _, l := range strings.Split(out, "\n") {
		if strings.HasPrefix(l, "diff ") {
			p, ok := header[l]
			current, unknown = p, unknown || !ok
		}
		if current != "" {
			files[current] = append(files[current], l)
		}
	}
	res := make(map[string]hunks, len(paths))
	for _, p := range paths {
		lines, ok := files[p]
		if !ok && unknown {
			added, removed, err := readHunks(r, dir, from, to, p)
			if err != nil {
				return nil, err
			}
			res[p] = hunks{added, removed}
			continue
		}
		added, removed, err := diffHunks(strings.Join(lines, "\n"))
		if err != nil {
			return nil, err
		}
		res[p] = hunks{added, removed}
	}
	return res, nil
}

// diffHunks is the lines a unified diff adds, numbered in the new file, and
// removes, numbered in the old one, from each hunk's header, @@ -a,b +c,d @@,
// in the diff's order; a context line counts on both sides. Each line's text
// has its trailing white space left out, since the runner trims the end of
// what git prints. File headers and git's "\ No newline at end of file" are
// not lines.
func diffHunks(diff string) (added, removed []numbered, err error) {
	inHunk := false
	oldN, newN := 0, 0
	for _, l := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(l, "diff "):
			inHunk = false
		case strings.HasPrefix(l, "@@"):
			if oldN, newN, err = hunkStarts(l); err != nil {
				return nil, nil, err
			}
			inHunk = true
		case !inHunk, strings.HasPrefix(l, `\`):
		case strings.HasPrefix(l, "+"):
			added = append(added, numbered{newN, strings.TrimRight(l[1:], " \t\r")})
			newN++
		case strings.HasPrefix(l, "-"):
			removed = append(removed, numbered{oldN, strings.TrimRight(l[1:], " \t\r")})
			oldN++
		default:
			oldN++
			newN++
		}
	}
	return added, removed, nil
}

// hunkStarts is the first line numbers of the old and the new side that the
// hunk header @@ -a[,b] +c[,d] @@ names, a and c.
func hunkStarts(header string) (oldN, newN int, err error) {
	f := strings.Fields(header)
	if len(f) < 3 || !strings.HasPrefix(f[1], "-") || !strings.HasPrefix(f[2], "+") {
		return 0, 0, fmt.Errorf("read the git diff hunk header %q: it is not @@ -a,b +c,d @@", header)
	}
	if oldN, err = rangeStart(f[1][1:]); err != nil {
		return 0, 0, fmt.Errorf("read the git diff hunk header %q: %w", header, err)
	}
	if newN, err = rangeStart(f[2][1:]); err != nil {
		return 0, 0, fmt.Errorf("read the git diff hunk header %q: %w", header, err)
	}
	return oldN, newN, nil
}

// rangeStart is the first line number of a hunk header's range, a or a,b.
func rangeStart(rng string) (int, error) {
	start, _, _ := strings.Cut(rng, ",")
	return strconv.Atoi(start)
}

// atPositions is the text of the lines of added that removed removes at the
// same number, in added's order: main's lines the branch removes where main
// added them, both numbered in the same file (ADR-0309).
func atPositions(added, removed []numbered) []string {
	at := make(map[int]bool, len(removed))
	for _, l := range removed {
		at[l.n] = true
	}
	out := []string{}
	for _, l := range added {
		if at[l.n] {
			out = append(out, l.text)
		}
	}
	return out
}

// texts is the text of lines, in their order.
func texts(lines []numbered) []string {
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		out = append(out, l.text)
	}
	return out
}

// undoneLines is the lines of added that removed holds, counted as a
// multiset of line text, in added's order: a line added twice and removed
// once is given once.
func undoneLines(added, removed []string) []string {
	left := map[string]int{}
	for _, l := range removed {
		left[l]++
	}
	out := []string{}
	for _, l := range added {
		if left[l] > 0 {
			left[l]--
			out = append(out, l)
		}
	}
	return out
}

// notKept is the lines of undone, main's lines the branch removes from the
// file at path where main added them (ADR-0309), that no line the branch
// adds keeps (I-0360, ADR-0297). In a markdown
// file, a line of main's is kept, not undone, when a line the branch adds
// holds every word of it, as a multiset of the words split on white space
// with the punctuation at each word's ends left out, in any order: so a
// sentence the branch adds to main's one-line paragraph or table row, or a
// word it moves or a bracket it drops there, keeps main's line. Each added
// line keeps one of main's lines at most, the first it holds in undone's
// order. An added line that is one of mainRemoved, the lines main removed
// since the branch's old base, keeps none: the branch restores main's old
// line, which is the revert ADR-0155 catches. A line with no words, and every
// line of a file that is not markdown, where a line that keeps every word,
// such as one commented out, can still undo main's change, is undone as the
// position compare found it (ADR-0155, ADR-0309).
func notKept(path string, undone, branchAdded, mainRemoved []string) []string {
	if !mayKeep(path, undone, branchAdded) {
		return undone
	}
	bags := make([]map[string]int, 0, len(branchAdded))
	for _, a := range branchAdded {
		if !slices.Contains(mainRemoved, a) {
			bags = append(bags, wordBag(a))
		}
	}
	out := []string{}
	for _, l := range undone {
		if !keepOnce(wordBag(l), bags) {
			out = append(out, l)
		}
	}
	return out
}

// mayKeep reports whether a line the branch adds to the file at path may
// keep one of undone: the file is markdown and both are non-empty.
func mayKeep(path string, undone, branchAdded []string) bool {
	return len(undone) > 0 && len(branchAdded) > 0 && strings.EqualFold(filepath.Ext(path), ".md")
}

// keepOnce takes out of bags the first bag that holds every word of words,
// and reports whether it found one; a line with no words is never kept.
func keepOnce(words map[string]int, bags []map[string]int) bool {
	if len(words) == 0 {
		return false
	}
	for i, b := range bags {
		if b != nil && holdsAll(b, words) {
			bags[i] = nil
			return true
		}
	}
	return false
}

// holdsAll reports whether bag holds each word of words at least as often.
func holdsAll(bag, words map[string]int) bool {
	for w, n := range words {
		if bag[w] < n {
			return false
		}
	}
	return true
}

// wordBag is the words of line, split on white space with the punctuation
// at each word's ends left out, counted; a word of punctuation alone is none.
func wordBag(line string) map[string]int {
	bag := map[string]int{}
	for _, w := range strings.Fields(line) {
		if w = strings.TrimFunc(w, unicode.IsPunct); w != "" {
			bag[w]++
		}
	}
	return bag
}
