package messages

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// Asking about a hold, and sharing paths (ADR-0134). A ready story held on
// overlap alone has no agent until it starts, so flai asks for it: a message
// by flai, from the held story to the story in progress that holds it, about
// the overlapping paths. That story's agent answers by narrowing its touches,
// by sharing the paths with a split of the work, which the conversation's
// front matter keeps under shares, or by replying why the hold stands.

// AskAuthor is the author of the message that asks the holding story's agent
// about a hold, flai itself.
const AskAuthor = "flai"

// AskOptions describe an ask about a hold: the ready story held, the story in
// progress that holds it, and the paths the two overlap on.
type AskOptions struct {
	Held   string   // the ready story held on overlap
	Holder string   // the story in progress that holds it
	Paths  []string // the paths of the hold, taken as given, cleaned
	Now    time.Time
}

// AskHold asks the agent of opt.Holder about the hold on opt.Held, in the two
// stories' open conversation, whichever started it, or in a new one from the
// held story: the held story's title and goal, the paths, and the three
// answers. It asks once per pair while the held story stays in ready: when
// the open conversation already holds a message by flai for the held story
// since it last entered ready, nothing is written, and it returns that
// conversation and false. The bool says whether it wrote. The held story must
// be ready and the holder in progress; the file is refused, and nothing
// written, when the project's lint rejects it.
func AskHold(r *workitem.Repo, opt AskOptions) (*Conversation, bool, error) {
	held, err := storyIn(r, opt.Held, "held story", workitem.Ready)
	if err != nil {
		return nil, false, err
	}
	holder, err := storyIn(r, opt.Holder, "holding story", workitem.InProgress)
	if err != nil {
		return nil, false, err
	}
	if held.ID == holder.ID {
		return nil, false, fmt.Errorf("%s cannot hold itself: name the story in progress that holds it", held.ID)
	}
	paths, err := cleanPaths("paths", opt.Paths)
	if err != nil {
		return nil, false, err
	}
	if len(paths) == 0 {
		return nil, false, fmt.Errorf("an ask about the hold on %s needs the paths %s holds it on", held.ID, holder.ID)
	}
	now := opt.Now.UTC().Format(workitem.TimeFormat)
	c, err := openBetween(r, held.ID, holder.ID)
	if err != nil {
		return nil, false, err
	}
	if c == nil {
		id := NextID(r)
		c, err = begin(r, message{id: id, text: askText(id, held, holder, paths), author: AskAuthor, from: held, to: holder, about: paths, now: now})
		return c, err == nil, err
	}
	if askedSince(c, held) {
		return c, false, nil
	}
	was := c.Marshal()
	before := betweenLine(c.From, c.To, c.About)
	for _, p := range paths {
		if !contains(c.About, p) {
			c.About = append(c.About, p)
		}
	}
	c.Body = strings.Replace(c.Body, "\n"+before+"\n", "\n"+betweenLine(c.From, c.To, c.About)+"\n", 1)
	c.add(now, AskAuthor, storyOf(c, held.ID), askText(c.ID, held, holder, paths))
	if err := save(r, c, was); err != nil {
		return nil, false, err
	}
	return c, true, nil
}

// askedSince says whether c holds a message by flai written for held at or
// after held last entered ready.
func askedSince(c *Conversation, held *workitem.Item) bool {
	since := readySince(held)
	for _, e := range c.Entries() {
		if e.Author != AskAuthor || workitem.CanonicalID(e.Story) != held.ID {
			continue
		}
		if at, err := time.Parse(workitem.TimeFormat, e.At); err == nil && !at.Before(since) {
			return true
		}
	}
	return false
}

// readySince is when the story last entered ready, else when it was created.
func readySince(it *workitem.Item) time.Time {
	for i := len(it.Transitions) - 1; i >= 0; i-- {
		if it.Transitions[i].To == workitem.Ready {
			t, _ := time.Parse(workitem.TimeFormat, it.Transitions[i].At)
			return t
		}
	}
	t, _ := time.Parse(workitem.TimeFormat, it.Created)
	return t
}

// askText is the ask about the hold on held by holder, in conversation id.
func askText(id string, held, holder *workitem.Item, paths []string) string {
	goal := held.ID + " has no goal written."
	if s, ok := workitem.NarrativeSection(held.Body, "Goal"); ok && strings.TrimSpace(s.Text) != "" {
		goal = "Its goal:\n\n" + quote(strings.TrimSpace(s.Text))
	}
	return fmt.Sprintf(`%[2]s holds %[1]s on overlap alone: may %[1]s start beside it?

%[1]s, %[3]q, is ready, and its claim overlaps %[2]s's on %[4]s. It has no agent until it starts, so flai asks for it (ADR-0134).

%[5]s

Answer in one of three ways:

1. Narrow your touches with `+"`flai touches %[2]s --remove <path>`"+`, or on the task that names the path, if %[2]s will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `+"`flai message share %[6]s --paths <path> \"<split>\"`"+`, or the MCP tool `+"`message_share`"+`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and %[1]s can start.
3. Reply with `+"`message_reply`"+` saying why the hold stands: %[1]s waits until %[2]s moves to review, is cancelled, or is sent back.`,
		held.ID, holder.ID, held.Title, quotedPaths(paths), goal, id)
}

// quotedPaths is the paths in backticks, joined with commas.
func quotedPaths(paths []string) string {
	quoted := make([]string, len(paths))
	for i, p := range paths {
		quoted[i] = "`" + p + "`"
	}
	return strings.Join(quoted, ", ")
}

// ShareOptions describe a share: the conversation, who shares and for which
// story, the paths, and the split of the work.
type ShareOptions struct {
	ID     string   // the conversation between the holding and the held story
	Story  string   // the story the sharer writes for: the holding story's agent's
	Author string   // who shares: that agent, or the operator, the holding story's owner
	Paths  []string // what is shared, each inside the two stories' overlap
	Split  string   // who changes what
	Now    time.Time
}

// Share records on conversation opt.ID that its story in progress, the
// holder, shares opt.Paths with its ready story, the held, split as
// opt.Split says (ADR-0134): a share in its front matter, and an entry by
// opt.Author for the holder giving the split. While the share is in force an
// overlap between the two on those paths does not hold. Only the holding
// story's agent, writing for it, or the operator, its owner, may share; each
// path must lie inside the two stories' claims' overlap. A conversation that
// reads as closed, an empty split, and no path are refused, and nothing is
// written; so is a file the project's lint rejects.
func Share(r *workitem.Repo, opt ShareOptions) (*Conversation, error) {
	split := strings.TrimSpace(opt.Split)
	if split == "" {
		return nil, fmt.Errorf("a share needs the split: say who changes what in the paths shared")
	}
	author := cleanAuthor(opt.Author)
	if author == "" {
		return nil, fmt.Errorf("an author is required (--by, FLAI_AGENT, or config author)")
	}
	paths, err := cleanPaths("paths", opt.Paths)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("a share needs at least one path: give --paths with each path the two stories share")
	}
	c, err := Get(r, opt.ID)
	if err != nil {
		return nil, err
	}
	if closed, why := c.Closed(r); closed {
		return nil, fmt.Errorf("%s is closed (%s) and takes no share: a share is made in the open conversation between the two stories", c.ID, why)
	}
	held, holder, err := sides(r, c)
	if err != nil {
		return nil, err
	}
	if workitem.CanonicalID(strings.TrimSpace(opt.Story)) != holder.ID && (holder.Owner == "" || author != cleanAuthor(holder.Owner)) {
		return nil, fmt.Errorf("%s may not share %s's paths with %s: only %s's agent, writing for %s, or the operator, %s, may (ADR-0134); nothing was written", sharer(opt.Story, author), holder.ID, held.ID, holder.ID, holder.ID, ownerWords(holder))
	}
	if err := insideOverlap(r, held, holder, paths); err != nil {
		return nil, err
	}
	now := opt.Now.UTC().Format(workitem.TimeFormat)
	was := c.Marshal()
	c.Shares = append(c.Shares, workitem.Share{Holder: holder.ID, Held: held.ID, Paths: paths, Split: split, By: author, At: now})
	c.add(now, author, storyOf(c, holder.ID), fmt.Sprintf("%s shares %s with %s, split so:\n\n%s\n\nThe overlap between the two on them no longer holds (ADR-0134).", holder.ID, quotedPaths(paths), held.ID, quote(split)))
	if err := save(r, c, was); err != nil {
		return nil, err
	}
	return c, nil
}

// sides finds the conversation's ready story, the held, and the other, the
// holder, which must be in progress.
func sides(r *workitem.Repo, c *Conversation) (held, holder *workitem.Item, err error) {
	from, err := r.Get(c.From)
	if err != nil {
		return nil, nil, fmt.Errorf("read %s, of %s: %w", c.From, c.ID, err)
	}
	to, err := r.Get(c.To)
	if err != nil {
		return nil, nil, fmt.Errorf("read %s, of %s: %w", c.To, c.ID, err)
	}
	switch {
	case from.Status == workitem.Ready && to.Status == workitem.Ready:
		return nil, nil, fmt.Errorf("%s and %s, the two stories of %s, are both ready: a share is made by the story in progress that holds a ready one; nothing was written", from.ID, to.ID, c.ID)
	case from.Status == workitem.Ready:
		held, holder = from, to
	case to.Status == workitem.Ready:
		held, holder = to, from
	default:
		return nil, nil, fmt.Errorf("neither %s (%s) nor %s (%s), the two stories of %s, is ready: a share lifts the hold on a ready story, so there is nothing to share; nothing was written", from.ID, stateWords(from.Status), to.ID, stateWords(to.Status), c.ID)
	}
	if holder.Status != workitem.InProgress {
		return nil, nil, fmt.Errorf("%s is %s, and only a story in progress holds a ready one (ADR-0096), so it has nothing to share with %s; nothing was written", holder.ID, stateWords(holder.Status), held.ID)
	}
	return held, holder, nil
}

// sharer names who tried to share, in a refusal.
func sharer(story, author string) string {
	if s := workitem.CanonicalID(strings.TrimSpace(story)); s != "" {
		return author + ", writing for " + s + ","
	}
	return author
}

// ownerWords names the holding story's owner, in a refusal.
func ownerWords(holder *workitem.Item) string {
	if holder.Owner == "" {
		return "its owner, of whom it names none"
	}
	return "its owner " + holder.Owner
}

// insideOverlap refuses a path that does not overlap both an entry of the
// holder's claim and one of the held story's, naming the overlap.
func insideOverlap(r *workitem.Repo, held, holder *workitem.Item, paths []string) error {
	items, err := r.List(false)
	if err != nil {
		return fmt.Errorf("read the work items for the claims of %s and %s: %w; run flai check to see what is wrong", holder.ID, held.ID, err)
	}
	h := workitem.NewHolds(items, r.Manifest.Projects)
	heldClaim, holderClaim := h.Claim(held), h.Claim(holder)
	var overlap []string
	for _, a := range heldClaim {
		for _, b := range holderClaim {
			if p := deeper(a, b); workitem.PathsOverlap(a, b) && !contains(overlap, p) {
				overlap = append(overlap, p)
			}
		}
	}
	if len(overlap) == 0 {
		return fmt.Errorf("the claims of %s and %s do not overlap, so %s does not hold %s and there is nothing to share; nothing was written", held.ID, holder.ID, holder.ID, held.ID)
	}
	for _, p := range paths {
		if !overlapsAny(p, heldClaim) || !overlapsAny(p, holderClaim) {
			return fmt.Errorf("--paths %s lies outside the overlap of %s's and %s's claims, which is %s: share only paths both stories claim; nothing was written", p, held.ID, holder.ID, quotedPaths(overlap))
		}
	}
	return nil
}

// overlapsAny says whether p overlaps an entry of claim.
func overlapsAny(p string, claim []string) bool {
	for _, e := range claim {
		if workitem.PathsOverlap(p, e) {
			return true
		}
	}
	return false
}

// deeper is the longer of two paths, the second when they are as long: of
// two that overlap, the one both stories change.
func deeper(a, b string) string {
	if len(strings.TrimSuffix(a, "/")) > len(strings.TrimSuffix(b, "/")) {
		return a
	}
	return b
}

// storyIn finds the story id, which must exist, not be archived, and be in
// state; role names it in the error.
func storyIn(r *workitem.Repo, id, role, state string) (*workitem.Item, error) {
	canon := workitem.CanonicalID(strings.TrimSpace(id))
	if !storyPattern.MatchString(canon) {
		return nil, fmt.Errorf("%q, the %s, is not a story ID such as S-0001", id, role)
	}
	it, err := r.Get(canon)
	if errors.Is(err, workitem.ErrNotFound) {
		return nil, fmt.Errorf("%s, the %s, does not exist", canon, role)
	}
	if err != nil {
		return nil, fmt.Errorf("read %s, the %s: %w", canon, role, err)
	}
	if it.Archived {
		return nil, fmt.Errorf("%s, the %s, is archived; it must be %s", it.ID, role, stateWords(state))
	}
	if it.Status != state {
		return nil, fmt.Errorf("%s, the %s, is %s; it must be %s", it.ID, role, stateWords(it.Status), stateWords(state))
	}
	return it, nil
}

// cleanPaths returns the paths cleaned as cleanAbout does, the error naming
// flag.
func cleanPaths(flag string, paths []string) ([]string, error) {
	var out []string
	for _, p := range paths {
		rel, err := cleanPath(p)
		if err != nil {
			return nil, fmt.Errorf("--%s: %w", flag, err)
		}
		if !contains(out, rel) {
			out = append(out, rel)
		}
	}
	return out, nil
}

// validateShares checks the schema of a conversation's shares.
func validateShares(shares []workitem.Share) []string {
	var errs []string
	for i, sh := range shares {
		where := fmt.Sprintf("shares[%d]", i)
		for _, side := range [][2]string{{"holder", sh.Holder}, {"held", sh.Held}} {
			if !storyPattern.MatchString(side[1]) {
				errs = append(errs, fmt.Sprintf("%s: %s %q must be a story ID such as S-0001", where, side[0], side[1]))
			}
		}
		if len(sh.Paths) == 0 {
			errs = append(errs, where+": paths must name at least one path")
		}
		for _, p := range sh.Paths {
			if _, err := cleanPath(p); err != nil {
				errs = append(errs, where+": paths: "+err.Error())
			}
		}
		if strings.TrimSpace(sh.Split) == "" {
			errs = append(errs, where+": split is required")
		}
		if strings.TrimSpace(sh.By) == "" {
			errs = append(errs, where+": by is required")
		}
		if _, err := time.Parse(workitem.TimeFormat, sh.At); err != nil {
			errs = append(errs, fmt.Sprintf("%s: at %q is not a UTC timestamp", where, sh.At))
		}
	}
	return errs
}
