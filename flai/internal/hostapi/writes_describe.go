package hostapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	"github.com/bytepunx/system-flow/flai/internal/verify"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// orchestrateNow is the write that runs flai serve orchestrate <verb>, gated
// by the orchestrate action, and journalled as what it did.
func orchestrateNow(verb string) spec {
	return spec{action: ActionOrchestrate, describe: describeOrchestrate, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
		if _, e := decode[struct{}](raw); e != nil {
			return nil, "", e
		}
		return []string{"serve", "orchestrate", verb}, "", nil
	}}
}

// describeOrchestrate reads flai serve orchestrate stop's and start's --json
// shape for the journal.
func describeOrchestrate(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Held         bool `json:"held"`
		Orchestrator struct {
			Agent   string `json:"agent"`
			Command string `json:"command"`
			PID     int    `json:"pid"`
		} `json:"orchestrator"`
	}
	_ = json.Unmarshal(w.Data, &said)
	run := said.Orchestrator
	switch {
	case said.Held && run.PID == 0: // a run that could not be started, held as it is
		return "done", fmt.Sprintf("held the orchestrator %s stopped", run.Agent)
	case said.Held:
		return "done", fmt.Sprintf("stopped the orchestrator %s (pid %d) and held it stopped", run.Agent, run.PID)
	}
	return "done", fmt.Sprintf("lifted the hold and started %s as %s (pid %d)", run.Command, run.Agent, run.PID)
}

// describeAnalyze reads flai analyze's --json shape for the journal.
func describeAnalyze(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Focus   string `json:"focus"`
		Agent   string `json:"agent"`
		Command string `json:"command"`
		PID     int    `json:"pid"`
	}
	_ = json.Unmarshal(w.Data, &said)
	return "done", fmt.Sprintf("started %s to analyze, focus %s, as %s (pid %d)", said.Command, said.Focus, said.Agent, said.PID)
}

// curateSinceReads says whether since is a curator's since, in a form flai
// curate takes: a YYYY-MM-DD day or an RFC 3339 time.
func curateSinceReads(since string) bool {
	for _, form := range []string{time.DateOnly, time.RFC3339} {
		if _, err := time.Parse(form, since); err == nil {
			return true
		}
	}
	return false
}

// describeCurate reads flai curate's --json shape for the journal.
func describeCurate(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Since   string `json:"since"`
		Agent   string `json:"agent"`
		Command string `json:"command"`
		PID     int    `json:"pid"`
	}
	_ = json.Unmarshal(w.Data, &said)
	from := "everything, its first run"
	if said.Since != "" {
		from = "since " + said.Since
	}
	return "done", fmt.Sprintf("started %s to curate, %s, as %s (pid %d)", said.Command, from, said.Agent, said.PID)
}

// describePlan reads flai plan's --json shape for the journal.
func describePlan(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Item    string `json:"item"`
		Agent   string `json:"agent"`
		Command string `json:"command"`
		PID     int    `json:"pid"`
	}
	_ = json.Unmarshal(w.Data, &said)
	return "done", fmt.Sprintf("started %s to plan %s as %s (pid %d)", said.Command, said.Item, said.Agent, said.PID)
}

// roleName is a role's name as a role file gives it, in kebab case.
var roleName = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)

// roleRunArgs is role.run's command line: flai role run with the role's
// name, its item when one is given, and each input as --input=name=value in
// the order given. An input's name is one word and its value one line.
func roleRunArgs(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
	in, e := decode[struct {
		Name   string `json:"name"`
		Item   string `json:"item"`
		Inputs []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"inputs"`
	}](raw)
	if e != nil {
		return nil, "", e
	}
	if !roleName.MatchString(in.Name) {
		return nil, "", bad("%q is not a role's name: a role file names it in kebab case, such as reviewer", in.Name)
	}
	args := []string{"role", "run", in.Name}
	if in.Item != "" {
		if e := needID(in.Item); e != nil {
			return nil, "", e
		}
		args = append(args, in.Item)
	}
	for _, input := range in.Inputs {
		if !listValue.MatchString(input.Name) {
			return nil, "", bad("input %q is not an input's name: one word, such as focus, without = or a leading dash", input.Name)
		}
		if strings.ContainsAny(input.Value, "\x00\n\r") {
			return nil, "", bad("the value of input %s holds a line break or NUL: give it on one line", input.Name)
		}
		args = append(args, "--input="+input.Name+"="+input.Value)
	}
	return args, "", nil
}

// describeRole reads flai role run's --json shape for the journal.
func describeRole(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Role    string `json:"role"`
		Item    string `json:"item"`
		Agent   string `json:"agent"`
		Command string `json:"command"`
		PID     int    `json:"pid"`
	}
	_ = json.Unmarshal(w.Data, &said)
	on := "on no item"
	if said.Item != "" {
		on = "on " + said.Item
	}
	return "done", fmt.Sprintf("started %s to run the role %s %s as %s (pid %d)", said.Command, said.Role, on, said.Agent, said.PID)
}

// agentNow is the write that runs flai serve agent <verb> for a story, gated
// by the agent action, and journalled as what it did and why.
func agentNow(verb, did, why string) spec {
	return spec{action: ActionAgent, describe: describeAgentNow(did, why), build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
		id, e := storyOf(raw)
		if e != nil {
			return nil, "", e
		}
		return []string{"serve", "agent", verb, id}, "", nil
	}}
}

// storyOf is the story an agent method names, refused when it names none or
// an item that is not a story.
func storyOf(raw json.RawMessage) (string, *channel.Error) {
	in, e := decode[struct {
		ID string `json:"id"`
	}](raw)
	if e != nil {
		return "", e
	}
	if e := needID(in.ID); e != nil {
		return "", e
	}
	if !strings.HasPrefix(in.ID, "S-") {
		return "", bad("%s is not a story; only a story has an agent", in.ID)
	}
	return in.ID, nil
}

// describeAgentStop reads flai serve agent stop's --json shape for the
// journal.
func describeAgentStop(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Story string `json:"story"`
		Agent string `json:"agent"`
		PID   int    `json:"pid"`
	}
	_ = json.Unmarshal(w.Data, &said)
	return "done", fmt.Sprintf("stopped %s's agent %s (pid %d)", said.Story, said.Agent, said.PID)
}

// describeAgentNow reads flai serve agent start's, restart's, and commit's --json shape
// for the journal.
func describeAgentNow(did, why string) func(res any, err *channel.Error) (outcome, detail string) {
	return func(res any, err *channel.Error) (outcome, detail string) {
		if err != nil {
			return "failed", err.Message
		}
		w, _ := res.(Written)
		var said struct {
			Story   string `json:"story"`
			Agent   string `json:"agent"`
			Command string `json:"command"`
			PID     int    `json:"pid"`
			Queued  string `json:"queued"`
		}
		_ = json.Unmarshal(w.Data, &said)
		if said.Queued != "" {
			return "done", fmt.Sprintf("queued another agent for %s until there is room", said.Story)
		}
		return "done", fmt.Sprintf("%s %s for %s as %s%s (pid %d)", did, said.Command, said.Story, said.Agent, why, said.PID)
	}
}

// describeChecksRun reads flai checks run/cancel's own --json shape
// (ChecksRun) for the journal, rather than describe's publish shape.
func describeChecksRun(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Story   string `json:"story"`
		Outcome string `json:"outcome"`
	}
	_ = json.Unmarshal(w.Data, &said)
	if said.Outcome == "" {
		return "done", said.Story
	}
	return "done", said.Story + ": " + said.Outcome
}

// testRun is what test.run is asked: the item whose worktree the tiers run
// in, none for the main checkout, and flai test's own arguments.
type testRun struct {
	ID    string   `json:"id"`
	Paths []string `json:"paths"`
	All   bool     `json:"all"`
	Max   int      `json:"max"`
}

// testArgs is test.run's command line, before --json: flai test with --all
// and --max when given, and the paths after --. A path is relative to the
// checkout's root and stays inside it; flai test refuses one that is not there.
func testArgs(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
	in, e := decode[testRun](raw)
	if e != nil {
		return nil, "", e
	}
	if in.ID != "" {
		if e := needID(in.ID); e != nil {
			return nil, "", e
		}
	}
	if in.Max < 0 {
		return nil, "", bad("max %d is not a number of findings; give 1 or more, or leave it out for %d", in.Max, verify.DefaultMax)
	}
	args := []string{"test"}
	if in.All {
		args = append(args, "--all")
	}
	if in.Max > 0 {
		args = append(args, "--max="+strconv.Itoa(in.Max))
	}
	if len(in.Paths) > 0 {
		args = append(args, "--")
	}
	for _, p := range in.Paths {
		p = strings.TrimSpace(p)
		switch {
		case strings.ContainsAny(p, "\x00\n\r"):
			return nil, "", bad("%q in paths holds a line break or NUL", p)
		case !filepath.IsLocal(p):
			return nil, "", bad("%q in paths is not a file or folder inside the checkout; give it relative to the repository's root, such as flai/internal/verify", p)
		}
		args = append(args, filepath.ToSlash(p))
	}
	return args, "", nil
}

// testDir is the checkout test.run runs in: the worktree of the item named,
// which must exist, or the project's main checkout, as the MCP tool test's is.
func testDir(p channel.Project, raw json.RawMessage) (string, *channel.Error) {
	in, e := decode[testRun](raw)
	if e != nil {
		return "", e
	}
	return worktreeDir(p, in.ID, "leave id out to run in the main checkout")
}

// docDir is the checkout doc.show and doc.save run in: the worktree of the
// story named, which must exist, or the project's main checkout (S-0378).
func docDir(p channel.Project, raw json.RawMessage) (string, *channel.Error) {
	in, e := decode[struct {
		Story string `json:"story"`
	}](raw)
	if e != nil {
		return "", e
	}
	return worktreeDir(p, in.Story, "leave story out to edit the main checkout")
}

// worktreeDir is the worktree of the item id names, which must exist, or
// the project's main checkout when id is empty; instead says how to do
// without one.
func worktreeDir(p channel.Project, id, instead string) (string, *channel.Error) {
	if id == "" {
		return p.Root, nil
	}
	repo, err := workitem.Open(p.Root)
	if err != nil {
		return "", failed(err)
	}
	it, err := repo.Get(id)
	if errors.Is(err, workitem.ErrNotFound) {
		return "", &channel.Error{Code: NotFound, Message: fmt.Sprintf("%s is not a work item here; name a story with a worktree, or %s", id, instead)}
	}
	if err != nil {
		return "", failed(err)
	}
	wt := repo.WorktreePath(it.ID)
	if st, err := os.Stat(wt); err != nil || !st.IsDir() {
		return "", bad("%s has no worktree; open it with flai stream open %s on the host, or %s", it.ID, it.ID, instead)
	}
	return wt, nil
}

// testResult refuses an answer of flai test's that is not its result, an
// object: flai test --json prints one whenever it answers.
func testResult(w Written) *channel.Error {
	if !isObject(w.Data) {
		return &channel.Error{Code: channel.CodeInternal, Message: "flai test answered without its result on standard output; run flai test --json in the checkout on the host to see why"}
	}
	return nil
}

// verifyRun is what verify.run is asked: the story, and the most findings
// to answer.
type verifyRun struct {
	ID  string `json:"id"`
	Max int    `json:"max"`
}

// verifyArgs is verify.run's command line, before --json: flai verify with
// the story and --max when given.
func verifyArgs(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
	in, e := decode[verifyRun](raw)
	if e != nil {
		return nil, "", e
	}
	if e := needStory(in.ID); e != nil {
		return nil, "", e
	}
	if in.Max < 0 {
		return nil, "", bad("max %d is not a number of findings; give 1 or more, or leave it out for %d", in.Max, verify.DefaultMax)
	}
	args := []string{"verify", in.ID}
	if in.Max > 0 {
		args = append(args, "--max="+strconv.Itoa(in.Max))
	}
	return args, "", nil
}

// verifyResult refuses an answer of flai verify's that is not its report,
// an object: flai verify --json prints one whenever it answers.
func verifyResult(w Written) *channel.Error {
	if !isObject(w.Data) {
		return &channel.Error{Code: channel.CodeInternal, Message: "flai verify answered without its report on standard output; run flai verify <story> --json in the project on the host to see why"}
	}
	return nil
}

// verifyStored refuses an answer of flai verify --last's that is neither
// the stored report, an object, nor null, which says there is none.
func verifyStored(w Written) *channel.Error {
	if !isObject(w.Data) && string(w.Data) != "null" {
		return &channel.Error{Code: channel.CodeInternal, Message: "flai verify --last answered neither the last report nor null; run flai verify <story> --last --json in the project on the host to see why"}
	}
	return nil
}

// describeVerifyRun reads flai verify --json's report for the journal: the
// story, and passed or the step it stopped at.
func describeVerifyRun(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said verify.Report
	_ = json.Unmarshal(w.Data, &said)
	switch {
	case said.Passed:
		return "done", said.Story + ": passed"
	case said.StoppedAt != "":
		return "done", said.Story + ": stopped at " + said.StoppedAt
	}
	return "done", said.Story + ": failed"
}

// taskStopped reads flai task done's result and refuses one that stopped:
// the run reached exit 1 at a step besides the sync and the check (ADR-0107),
// which answers counts as an answer only so that the result reaches the
// dashboard. It is an error in flai's words with the result as data, a rule
// when a workflow rule refused the step (a move flai move does not allow),
// and a failure of flai's otherwise. A result that ran to the end is the
// answer.
func taskStopped(w Written) *channel.Error {
	var res struct {
		Task    string `json:"task"`
		Stopped string `json:"stopped"`
		Error   string `json:"error"`
	}
	_ = json.Unmarshal(w.Data, &res) // an answer that is not the result did not stop
	if res.Stopped == "" {
		return nil
	}
	var data map[string]any
	_ = json.Unmarshal(w.Data, &data)
	code := channel.CodeInternal
	if strings.HasPrefix(res.Error, "rule: ") || strings.Contains(res.Error, ": rule: ") {
		code = Rule
	}
	return &channel.Error{Code: code, Message: fmt.Sprintf("%s stopped at the %s step: %s", res.Task, res.Stopped, res.Error), Data: data}
}

// isObject reports whether out is a JSON object.
func isObject(out []byte) bool {
	out = bytes.TrimSpace(out)
	return len(out) > 0 && out[0] == '{' && json.Valid(out)
}

// answered reads an exit the spec counts as an answer as success, when its
// standard output holds the answer.
func answered(ran Ran, exits []int) Ran {
	if ran.Exit != 0 && slices.Contains(exits, ran.Exit) && isObject(ran.Stdout) {
		ran.Exit = 0
	}
	return ran
}

// describeTestRun reads flai test --json's result for the journal: passed,
// or the tier that failed.
func describeTestRun(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said verify.Result
	_ = json.Unmarshal(w.Data, &said)
	if said.Passed {
		return "done", "passed"
	}
	for _, t := range said.Tiers {
		if t.State == verify.Failed {
			return "done", "the tier " + t.Name + " failed"
		}
	}
	return "done", "failed"
}

// describeHost reads flai host upgrade's --json shape for the journal,
// naming the release asked for, version, or the newest when it is ""
// (ADR-0117); a process action says what was asked (its say).
func describeHost(version string) func(res any, err *channel.Error) (outcome, detail string) {
	return func(res any, err *channel.Error) (outcome, detail string) {
		if err != nil {
			return "failed", err.Message
		}
		w, _ := res.(Written)
		var said struct {
			Restarting bool `json:"restarting"`
			Upgrade    struct {
				Installed string `json:"installed"`
				Previous  string `json:"previous"`
				Current   string `json:"current"`
			} `json:"upgrade"`
		}
		_ = json.Unmarshal(w.Data, &said)
		switch {
		case said.Restarting && version != "":
			return "done", fmt.Sprintf("installed flai %s, as asked, over %s; the host restarts on it", said.Upgrade.Installed, said.Upgrade.Previous)
		case said.Restarting:
			return "done", fmt.Sprintf("installed the newest flai, %s, over %s; the host restarts on it", said.Upgrade.Installed, said.Upgrade.Previous)
		case said.Upgrade.Current != "":
			return "done", "flai " + said.Upgrade.Current + " is the latest"
		}
		return "done", ""
	}
}

// describeDashboardRestart, describeDashboardUpgrade, and describeDashboardStop
// read cmd/dashboard.go and cmd/dashboard_upgrade.go's own --json shapes for
// the journal, rather than describe's publish shape above. A failed
// call is described the same way regardless: "failed", err.Message.

func describeDashboardRestart(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Container  string `json:"container"`
		WasRunning bool   `json:"was_running"`
		Ref        string `json:"ref"`
	}
	_ = json.Unmarshal(w.Data, &said)
	verb := "started"
	if said.WasRunning {
		verb = "restarted"
	}
	return "done", fmt.Sprintf("%s %s (%s)", verb, said.Container, said.Ref)
}

// describeDashboardUpgrade names the release asked for, tag, or the
// configured one when it is "" (ADR-0117, ADR-0118).
func describeDashboardUpgrade(tag string) func(res any, err *channel.Error) (outcome, detail string) {
	asked := "the configured release"
	if tag != "" {
		asked = "release " + tag + ", as asked"
	}
	return func(res any, err *channel.Error) (outcome, detail string) {
		if err != nil {
			return "failed", err.Message
		}
		w, _ := res.(Written)
		var said struct {
			Container string `json:"container"`
			Outcome   string `json:"outcome"`
			From      string `json:"from"`
			To        string `json:"to"`
		}
		_ = json.Unmarshal(w.Data, &said)
		switch said.Outcome {
		case "up-to-date":
			return "done", fmt.Sprintf("%s already running %s, %s", said.Container, said.To, asked)
		case "started":
			return "done", fmt.Sprintf("started %s (%s), %s", said.Container, said.To, asked)
		case "upgraded":
			return "done", fmt.Sprintf("upgraded %s from %s to %s, %s", said.Container, said.From, said.To, asked)
		default:
			return "done", said.Container
		}
	}
}

func describeDashboardStop(res any, err *channel.Error) (outcome, detail string) {
	if err != nil {
		return "failed", err.Message
	}
	w, _ := res.(Written)
	var said struct {
		Container string   `json:"container"`
		State     string   `json:"state"`
		Serves    []string `json:"serves"`
	}
	_ = json.Unmarshal(w.Data, &said)
	switch said.State {
	case "stopped":
		return "done", said.Container + " stopped"
	case "still-running":
		return "done", said.Container + " unregistered; still serving " + strings.Join(said.Serves, ", ")
	case "not-running":
		return "done", said.Container + " was not running"
	default:
		return "done", said.Container
	}
}

// withJSON asks for JSON before a -- that ends the flags, else at the end.
func withJSON(args []string) []string {
	for i, a := range args {
		if a == "--" {
			out := append([]string{}, args[:i]...)
			return append(append(out, "--json"), args[i:]...)
		}
	}
	return append(append([]string{}, args...), "--json")
}
