package hostapi

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/docedit"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// spec is one write: how its params become a command line.
type spec struct {
	// action names the host action this method is: asked for while the
	// operator has not enabled it, it is refused and says what enables it.
	action string
	// actionOf, when set, is the host action a call is, read from its
	// params, where it is not the same for every call: role.run's is the
	// one the file of the role it names names (S-0417). The call is gated
	// and journalled under it as under action; "" gates nothing and
	// journals nothing, leaving the call to its command's refusal. Asked
	// once build has accepted the params.
	actionOf func(p channel.Project, raw json.RawMessage) string
	// needs, when set, names the host action a call is in the project it is
	// asked of, or "" for none, where action would gate every call:
	// accept.run under pull-request acceptance merges on the forge
	// (ADR-0170). The call is refused and journalled as action's are while
	// the operator has not enabled it, and is not journalled otherwise.
	// Asked once build has accepted the params.
	needs func(p channel.Project) string
	// uses names a host action this method performs as part of its work when
	// it is enabled, and does without when it is not; it is journalled then.
	uses string
	// hostwide marks a setting kept for every project (S-0105): it needs its
	// action enabled for every project, not this one only.
	hostwide bool
	// hostwideOf, when set, says from its params whether a call changes a
	// setting for every project, and so is gated as hostwide is, where the
	// same setting may also be kept for this project alone: settings.agent
	// and settings.harness with everywhere true (ADR-0261). Asked once build
	// has accepted the params.
	hostwideOf func(raw json.RawMessage) bool
	// record journals every call under this name without gating it: an act
	// on the host the operator consented to some other way (S-0098: import,
	// by naming the folder the repository is in).
	record string
	// about names the project the write changes when it is not the one the
	// request arrives through (S-0122: serving or removing a project from
	// the settings page, which may have no connection of its own): the
	// action is asked of that project, and the journal names it. Called
	// only once build has accepted the params.
	about func(raw json.RawMessage) (key, root string)
	// build validates and returns the arguments (without --json) and what goes on standard input.
	build func(p channel.Project, raw json.RawMessage) (args []string, stdin string, err *channel.Error)
	// do, when set, answers in flai serve's own process instead of running
	// flai with build's arguments, once build has accepted the params and
	// the host action allows it: components.start and components.stop, whose
	// dev servers flai serve supervises (S-0438). Its answer is journalled
	// and remembered by request ID as a command's is.
	do func(ctx context.Context, host Host, p channel.Project, raw json.RawMessage) (any, *channel.Error)
	// dir, when set, is the folder the command runs in instead of the
	// project's: test.run's story worktree (S-0273). It is asked once build
	// has accepted the params, before the host action is.
	dir func(p channel.Project, raw json.RawMessage) (string, *channel.Error)
	// exits maps exit codes that carry a payload on standard output to error codes.
	exits map[int]int
	// answers are exit codes besides 0 whose JSON object on standard output
	// is the answer, not an error: flai test's 1, a tier that failed
	// (S-0273), and flai task done's 1, which its judge reads (ADR-0107).
	// Without one on standard output the exit is the error it would
	// otherwise be.
	answers []int
	// judge, when set, reads a successful answer and may refuse it after all:
	// item.finalize refuses a story flai edit left unchanged, which was not a
	// draft (S-0201).
	judge func(w Written) *channel.Error
	// progress sends each log event to the dashboard while the command runs.
	progress bool
	// reads marks a command that changes nothing: no request ID is asked for.
	reads bool
	// describe says in a line what a successful call did, for the journal.
	// Nil means the publish shape (describe, below); a host action with
	// a different answer shape (dashboard.restart, dashboard.stop) sets its
	// own, or asked (dashboard.upgrade, host.upgrade). A failed call is
	// always "failed", err.Message, whichever this is.
	describe func(res any, err *channel.Error) (outcome, detail string)
	// asked, when set, is describe made from what was asked, for an answer
	// that does not say it: an upgrade names the release asked for, or the
	// newest (ADR-0117). Called only once build has accepted the params.
	asked func(raw json.RawMessage) func(res any, err *channel.Error) (outcome, detail string)
	// say, when set, is the journal's line for a call that succeeded, made
	// from what was asked (S-0105: which setting became what).
	say func(raw json.RawMessage) string
	// detachTimeout, when nonzero, runs the command with its own background
	// context instead of the request's: a write whose own success can close
	// the WebSocket connection the request arrived on (dashboard.restart
	// stops the container answering it; dashboard.upgrade does on a
	// successful swap; dashboard.stop does when it is the last project)
	// must not be killed by that connection's own context cancelling out
	// from under it (S-0081, found live: exec.CommandContext SIGKILLed the
	// subprocess mid-restart when the container's own disconnect cancelled
	// the request context that spawned it, leaving no container running at
	// all). Progress still tries the original context; sending on a closed
	// connection just fails silently (channel.Progress → send → conn.Write).
	detachTimeout time.Duration
}

func text(v string) string { return strings.Join(strings.Fields(v), " ") }

// owner is who the dashboard acts for: the manifest's owner, decided here.
// The dashboard does not get to name them.
func owner(p channel.Project) string {
	if m, err := manifest.Load(filepath.Join(p.Root, manifest.File)); err == nil && m.Owner != "" {
		return m.Owner
	}
	return "designer"
}

// orchestratorAccepts refuses the orchestrator's acceptance of id in p while
// the project's orchestration.permissions.accept_reviews is off, as flai
// accept does, so that the dashboard is told why before anything runs
// (ADR-0093).
func orchestratorAccepts(p channel.Project, id string) *channel.Error {
	m, err := manifest.Load(filepath.Join(p.Root, manifest.File))
	if err != nil {
		return &channel.Error{Code: Rule, Message: fmt.Sprintf("cannot read the project's manifest to check orchestration.permissions.%s, so the orchestrator does not accept %s: %v", manifest.PermitAcceptReviews, id, err)}
	}
	if err := (&workitem.Repo{Manifest: m}).OrchestratorPermits(manifest.PermitAcceptReviews, "accepts a story", id); err != nil {
		return &channel.Error{Code: Rule, Message: err.Error()}
	}
	return nil
}

// forgeAcceptance is the host action an acceptance in p is: ActionForge
// under pull-request acceptance, whose merge is on the forge (ADR-0170),
// and none under local acceptance. A manifest that cannot be read asks for
// none: flai accept refuses it itself.
func forgeAcceptance(p channel.Project) string {
	if m, err := manifest.Load(filepath.Join(p.Root, manifest.File)); err == nil && m.Acceptance.PullRequest() {
		return ActionForge
	}
	return ""
}

var (
	commitName = regexp.MustCompile(`^[0-9a-f]{4,64}$`)
	anyItemID  = regexp.MustCompile(`^[EST]-\d{1,6}$`)
	issueID    = regexp.MustCompile(`^I-\d{1,6}$`)
	threadID   = regexp.MustCompile(`^(?i:TH-)?\d{1,6}$|^TH-\d{1,6}$`)
	adrNumber  = regexp.MustCompile(`^(?:ADR-)?(\d{1,4})$`)
	listValue  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/@+-]*$`)
	sinceValue = regexp.MustCompile(`^\d+[dwh]$`)
	docHash    = regexp.MustCompile(`^[A-Fa-f0-9]{8,128}$`)
)

// upgradeRelease reads the release an upgrade's params name in field,
// host.upgrade's version or dashboard.upgrade's tag: "" for none, or a bare
// X.Y.Z, with no v, suffix, or leading zeros (ADR-0117 §3), refused before
// any command runs.
func upgradeRelease(raw json.RawMessage, field string) (string, *channel.Error) {
	in, e := decode[struct {
		Version string `json:"version"`
		Tag     string `json:"tag"`
	}](raw)
	if e != nil {
		return "", e
	}
	v, other, otherField := in.Version, in.Tag, "tag"
	if field == "tag" {
		v, other, otherField = in.Tag, in.Version, "version"
	}
	if other != "" {
		return "", bad("%s is not a parameter of this upgrade: give the release as %s", otherField, field)
	}
	if v == "" {
		return "", nil
	}
	if !buildinfo.Bare(v) {
		return "", bad("%s %q is not a release: give a bare X.Y.Z, such as 1.2.3, or leave %s out", field, v, field)
	}
	return v, nil
}

func needID(id string) *channel.Error {
	if !anyItemID.MatchString(id) {
		return bad("%q is not an item ID", id)
	}
	return nil
}

// needStory checks id is a story's ID.
func needStory(id string) *channel.Error {
	if e := needID(id); e != nil {
		return e
	}
	if !strings.HasPrefix(id, "S-") {
		return bad("%s is not a story", id)
	}
	return nil
}

// needTask checks id is a task's ID: flai task done closes a task, and a
// story or an epic moves on by other means (ADR-0107).
func needTask(id string) *channel.Error {
	if e := needID(id); e != nil {
		return e
	}
	if !strings.HasPrefix(id, "T-") {
		return bad("%s is not a task: task.done closes a task (T-nnnn); a story goes to review with item.move after its close-out, and an epic follows its stories", id)
	}
	return nil
}

// topicWords refuses a topic that is not one word (ADR-0047) before flai is
// asked, so that the dashboard says why.
func topicWords(list []string) *channel.Error {
	for _, v := range list {
		if v = strings.TrimSpace(v); v != "" && !workitem.ValidTopic(v) {
			return bad("topic %q is not one word: letters, digits, dot, dash, or underscore", v)
		}
	}
	return nil
}

// listEntry refuses v, a value of the list flag, where flai would: a touch
// by workitem.CheckTouch (I-0071), so one starting with a dot passes; the rest
// are single words or paths.
func listEntry(flag, v string) *channel.Error {
	v = strings.TrimSpace(v)
	if flag == "touches" {
		if err := workitem.CheckTouch(strings.TrimSuffix(v, "/")); err != nil {
			return bad("%s", err)
		}
		return nil
	}
	if !listValue.MatchString(v) {
		return bad("%s values are single words or paths without commas", flag)
	}
	return nil
}

// occurrenceArgs checks what issue.new and issue.bump say of an occurrence,
// the story it belongs to, its cost, and a note, each optional, and gives
// the flags for those given.
func occurrenceArgs(story, cost, note string) ([]string, *channel.Error) {
	var args []string
	if story != "" {
		if e := needStory(story); e != nil {
			return nil, e
		}
		args = append(args, "--story="+story)
	}
	if cost = strings.TrimSpace(cost); cost != "" {
		if d, err := time.ParseDuration(cost); err != nil || d < 0 {
			return nil, bad("cost %q is not a duration like 20m", cost)
		}
		args = append(args, "--cost="+cost)
	}
	if note = text(note); note != "" {
		args = append(args, "--note="+note)
	}
	return args, nil
}

func adrNo(v string) (string, *channel.Error) {
	m := adrNumber.FindStringSubmatch(strings.TrimSpace(v))
	if m == nil {
		return "", bad("%s is not an ADR", v)
	}
	n, _ := strconv.Atoi(m[1])
	if n == 0 {
		return "", bad("%s is not an ADR", v)
	}
	return strconv.Itoa(n), nil
}

// docRel checks a document path without requiring the file to exist (a save
// may create it): Markdown, inside the manifest's folders, no climbing.
func docRel(p channel.Project, rel string) (string, *channel.Error) {
	clean := filepath.ToSlash(filepath.Clean(rel))
	if rel == "" || filepath.IsAbs(rel) || strings.HasPrefix(clean, "..") || strings.HasPrefix(clean, "-") || strings.ContainsRune(rel, 0) {
		return "", bad("%q is not a path inside the repository", rel)
	}
	if !strings.HasSuffix(clean, ".md") {
		return "", bad("only Markdown files are edited")
	}
	dirs, err := layoutDirs(p.Root)
	if err != nil {
		return "", failed(err)
	}
	for _, d := range dirs {
		if strings.HasPrefix(clean, d+"/") {
			return clean, nil
		}
	}
	return "", bad("%s is not under the project's design, docs, or wip folder", clean)
}

// fileRel checks a path file.get and file.save name, as docRel does a
// document's: inside the repository, never a flag, not markdown, which
// doc.show and doc.save edit, and a file flai file reaches under the
// project's manifest (ADR-0147), refused with flai's reason before it runs.
func fileRel(p channel.Project, rel string) (string, *channel.Error) {
	clean := filepath.ToSlash(filepath.Clean(rel))
	if rel == "" || filepath.IsAbs(rel) || strings.HasPrefix(clean, "..") || strings.HasPrefix(clean, "-") || strings.ContainsRune(rel, 0) {
		return "", bad("%q is not a path inside the repository", rel)
	}
	if strings.EqualFold(filepath.Ext(clean), ".md") {
		return "", bad("%s is a markdown file; read and save it with doc.show and doc.save", clean)
	}
	m, err := manifest.Load(filepath.Join(p.Root, manifest.File))
	if err != nil {
		return "", failed(err)
	}
	reached, err := docedit.FileReach(m, clean)
	if err != nil {
		return "", bad("%s", err.Error())
	}
	return reached, nil
}

// undonePath checks a path a story branch undoes of main's change (ADR-0155),
// as stream.restore and accept.run's keep_removals name it: relative to the
// repository's root, inside it with no .. at all, and never a flag. Whether
// the story's record holds it is flai's to say.
func undonePath(field, given string) (string, *channel.Error) {
	p := strings.TrimSpace(given)
	clean := filepath.ToSlash(filepath.Clean(p))
	switch {
	case strings.ContainsAny(p, "\x00\n\r"):
		return "", bad("%q in %s holds a line break or NUL", given, field)
	case strings.HasPrefix(p, "-"):
		return "", bad("%q in %s is not a path: give it as git names it from the repository's root, such as docs/users/flai.md", given, field)
	case !filepath.IsLocal(p) || clean == "." || slices.Contains(strings.Split(filepath.ToSlash(p), "/"), ".."):
		return "", bad("%q in %s is not a file inside the checkout: give it relative to the repository's root, with no .., such as docs/users/flai.md", given, field)
	}
	return clean, nil
}

// decode reads params into v, refusing what is not an object.
func decode[T any](raw json.RawMessage) (T, *channel.Error) {
	var v T
	return v, params(raw, &v)
}

// specs are every write, the settings among them (S-0105).
func specs() map[string]spec {
	table := itemSpecs()
	for name, sp := range settingsSpecs() {
		table[name] = sp
	}
	// The host's agent command, name, and harnesses are kept per project
	// over a host-wide fallback (ADR-0261): a write of this project's needs
	// the settings action on here, and one of every project's, everywhere
	// true, needs it on for every project.
	for _, name := range []string{"settings.agent", "settings.harness"} {
		sp := table[name]
		sp.hostwide, sp.hostwideOf = false, asksEverywhere
		table[name] = sp
	}
	return table
}

// asksEverywhere is whether a call's params say everywhere true: it changes
// the setting for every project, not this one only (ADR-0261).
func asksEverywhere(raw json.RawMessage) bool {
	var in struct {
		Everywhere bool `json:"everywhere"`
	}
	_ = json.Unmarshal(raw, &in)
	return in.Everywhere
}

// contentHash is what docedit.Hash gives: a SHA-256 in hex.
var contentHash = regexp.MustCompile(`^[0-9a-f]{64}$`)

func isNature(v string) bool {
	for _, n := range workitem.Natures {
		if n == v {
			return true
		}
	}
	return false
}
