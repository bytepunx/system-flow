package hostapi

import (
	"encoding/json"
	"slices"
	"strconv"
	"strings"

	"github.com/bytepunx/system-flow/flai/internal/channel"
	ctxpack "github.com/bytepunx/system-flow/flai/internal/context"
	"github.com/bytepunx/system-flow/flai/internal/issues"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// build turns a write's params into its command line and standard input.
type build = func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error)

// one is a write that runs its command line.
func one(b build) spec { return spec{build: b} }

// read is a method that only reads, run as a write is.
func read(b build) spec { return spec{build: b, reads: true} }

// moveArgs is item.move's command line.
func moveArgs(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
	in, e := decode[struct {
		ID                 string `json:"id"`
		To                 string `json:"to"`
		Reason             string `json:"reason"`
		IncludeUncommitted bool   `json:"include_uncommitted"`
		DryRun             bool   `json:"dry_run"`
		Finalize           bool   `json:"finalize"`
	}](raw)
	if e != nil {
		return nil, "", e
	}
	if e := needID(in.ID); e != nil {
		return nil, "", e
	}
	if !isState(in.To) {
		return nil, "", bad("%q is not a state", in.To)
	}
	args := []string{"move", in.ID, in.To, "--by=" + owner(p)}
	preview := in.DryRun && in.To == workitem.Cancelled
	switch {
	case text(in.Reason) != "":
		args = append(args, "--reason="+text(in.Reason))
	case preview:
		args = append(args, "--reason=preview") // flai validates the move before it lists anything
	}
	if preview {
		args = append(args, "--dry-run")
	}
	// --yes is the designer's choice, in the acceptance confirmation, to include uncommitted files.
	if in.IncludeUncommitted && in.To == workitem.Done {
		args = append(args, "--yes")
	}
	// and, on a move to ready, the designer's say that a draft story is
	// finished (S-0199); without it flai refuses a draft as a rule
	if in.Finalize && in.To == workitem.Ready {
		args = append(args, "--yes")
	}
	return args, "", nil
}

// itemWriteSpecs are the writes to work items, the board, issues, and acceptance.
func itemWriteSpecs() map[string]spec {
	return map[string]spec{
		// item.move to done from review is acceptance: flai move refuses it as
		// flai accept does, so its refusal (exit 4) is Refused with the
		// refused object, findings or undoes_main, as data, as accept.run's
		// is (I-0145).
		"item.move": {exits: map[int]int{4: Refused}, build: moveArgs},

		// story.start: a ready story started in one call, as flai story start
		// starts it (S-0274): moved to in-progress, its stream opened, its
		// pack built, and the inbox read, as the agent the channel runs flai
		// as, FLAI_AGENT. The answer is flai story start --json's result. A
		// story not ready or held (exit 4) is refused with flai's reason and
		// the refused object as data, and nothing changed. A step that failed
		// after the move (exit 3) is flai's error, naming the step and the
		// command that finishes it, with the result as data, so the dashboard
		// sees the story moved and what the steps before it did.
		"story.start": {exits: map[int]int{3: channel.CodeInternal, 4: Refused}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Story  string `json:"story"`
				Budget string `json:"budget"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needStory(in.Story); e != nil {
				return nil, "", e
			}
			args := []string{"story", "start", in.Story}
			if budget := strings.TrimSpace(in.Budget); budget != "" {
				if _, err := ctxpack.ParseSize(budget); err != nil {
					return nil, "", bad("%s", err)
				}
				args = append(args, "--budget="+budget)
			}
			return args, "", nil
		}},

		"item.order": one(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID     string `json:"id"`
				Before string `json:"before"`
				After  string `json:"after"`
				Top    bool   `json:"top"`
				Bottom bool   `json:"bottom"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			given := 0
			for _, on := range []bool{in.Before != "", in.After != "", in.Top, in.Bottom} {
				if on {
					given++
				}
			}
			if given != 1 {
				return nil, "", bad("say where it goes with exactly one of before, after, top, bottom")
			}
			for _, id := range []string{in.ID, in.Before + in.After} {
				if id != "" {
					if e := needID(id); e != nil {
						return nil, "", e
					}
				}
			}
			switch {
			case in.Before != "":
				return []string{"order", in.ID, "--before=" + in.Before}, "", nil
			case in.After != "":
				return []string{"order", in.ID, "--after=" + in.After}, "", nil
			case in.Top:
				return []string{"order", in.ID, "--top"}, "", nil
			}
			return []string{"order", in.ID, "--bottom"}, "", nil
		}),

		// board.limit: a column's WIP limit, from the board's lane menu (S-0167).
		"board.limit": one(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Column string `json:"column"`
				Limit  *int   `json:"limit"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if !slices.Contains(workitem.LimitedColumns, in.Column) {
				return nil, "", bad("%q has no WIP limit; only %s do", in.Column, strings.Join(workitem.LimitedColumns, ", "))
			}
			if in.Limit == nil || *in.Limit < 0 || *in.Limit > maxWIPLimit {
				return nil, "", bad("limit is a whole number from 0, for none, to %d", maxWIPLimit)
			}
			return []string{"board", "limit", "--", in.Column, strconv.Itoa(*in.Limit)}, "", nil
		}),

		// board.archive: the done or cancelled items named, archived with
		// flai archive, from the board's Archive and Archive All menu entries
		// (S-0343). An empty list is refused: flai archive without IDs
		// archives every closed item.
		"board.archive": one(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				IDs []string `json:"ids"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if len(in.IDs) == 0 {
				return nil, "", bad("name the items to archive in ids; flai archive without IDs would archive every closed item")
			}
			for _, id := range in.IDs {
				if e := needID(id); e != nil {
					return nil, "", e
				}
			}
			return append([]string{"archive", "--"}, in.IDs...), "", nil
		}),

		"item.block": one(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID     string `json:"id"`
				Reason string `json:"reason"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			if text(in.Reason) == "" {
				return nil, "", bad("a reason is required")
			}
			return []string{"block", in.ID, "--reason=" + text(in.Reason)}, "", nil
		}),

		"item.unblock": one(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"unblock", in.ID}, "", nil
		}),

		"item.new": {exits: map[int]int{4: Refused}, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Type    string          `json:"type"`
				Title   string          `json:"title"`
				Nature  string          `json:"nature"`
				Parent  string          `json:"parent"`
				Tags    []string        `json:"tags"`
				Touches []string        `json:"touches"`
				Topics  []string        `json:"topics"` // what a story or epic is about (S-0135)
				Agent   *manifest.Agent `json:"agent"`
				Body    string          `json:"body"`
				Draft   bool            `json:"draft"` // a story not yet finalized (S-0199)
				// the cost of delay inputs (S-0204): an absent or empty key is
				// left out; the value is the planner's, not the form's
				CostOfDelay *struct {
					RevenuePerWeek   *string `json:"revenue_per_week"`
					PenaltyPerWeek   *string `json:"penalty_per_week"`
					TimeLostPerCycle *string `json:"time_lost_per_cycle"`
					Value            *string `json:"value"`
				} `json:"cost_of_delay"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Type != workitem.Epic && in.Type != workitem.Story {
				return nil, "", bad("type must be epic or story; tasks are written by the agent that pulls a story")
			}
			title := text(in.Title)
			if title == "" {
				return nil, "", bad("a title is required")
			}
			if len([]rune(title)) > 200 {
				return nil, "", bad("the title is longer than 200 characters")
			}
			nature := in.Nature
			if nature == "" {
				nature = "feature"
			}
			ok := false
			for _, n := range workitem.Natures {
				ok = ok || n == nature
			}
			if !ok {
				return nil, "", bad("nature must be one of %s", strings.Join(workitem.Natures, ", "))
			}
			if strings.TrimSpace(in.Body) == "" {
				return nil, "", bad("write what the item is for: the body is empty")
			}
			args := []string{in.Type, "new", "--nature=" + nature, "--owner=" + owner(p)}
			switch {
			case in.Type == workitem.Story && in.Parent == "":
				// no epic (S-0092): not every story belongs to one
			case in.Type == workitem.Story:
				if !anyItemID.MatchString(in.Parent) || !strings.HasPrefix(in.Parent, "E-") {
					return nil, "", bad("a story's parent must be an epic")
				}
				args = append(args, "--epic="+in.Parent)
			case in.Parent != "":
				return nil, "", bad("an epic has no parent")
			}
			if e := topicWords(in.Topics); e != nil {
				return nil, "", e
			}
			for flag, values := range map[string][]string{"tag": in.Tags, "touches": in.Touches} {
				for _, v := range values {
					if e := listEntry(flag, v); e != nil {
						return nil, "", e
					}
				}
			}
			for _, v := range in.Tags {
				args = append(args, "--tag="+strings.TrimSpace(v))
			}
			for _, v := range in.Touches {
				args = append(args, "--touches="+strings.TrimSpace(v))
			}
			for _, v := range in.Topics {
				args = append(args, "--topics="+strings.TrimSpace(v))
			}
			if !in.Agent.IsZero() {
				if in.Type != workitem.Story {
					return nil, "", bad("only a story carries an agent")
				}
				agent, e := agentArgs(in.Agent)
				if e != nil {
					return nil, "", e
				}
				args = append(args, agent...)
			}
			if in.Draft {
				if in.Type != workitem.Story {
					return nil, "", bad("only a story is a draft")
				}
				args = append(args, "--draft")
			}
			if c := in.CostOfDelay; c != nil {
				if c.Value != nil {
					return nil, "", bad("a new item's cost of delay takes its inputs only: the value is the planner's")
				}
				// flai refuses an amount that is not a number, and checks the duration
				for _, k := range []struct {
					flag  string
					value *string
				}{{"revenue-per-week", c.RevenuePerWeek}, {"penalty-per-week", c.PenaltyPerWeek}, {"time-lost-per-cycle", c.TimeLostPerCycle}} {
					if k.value != nil && text(*k.value) != "" {
						args = append(args, "--"+k.flag+"="+text(*k.value))
					}
				}
			}
			return append(args, "--body-stdin", "--autocommit", "--trailer="+Trailer, "--", title), in.Body, nil
		}},

		// item.edit: an item's own words, changed after it was made (S-0085),
		// from what item.show read. flai edit does the work: the retitle kept in
		// step everywhere, the parents' lists, the check with the change in
		// place, one commit.
		"item.edit": {exits: map[int]int{3: Conflict, 4: Refused}, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID      string    `json:"id"`
				Hash    string    `json:"hash"`
				Title   *string   `json:"title"`
				Nature  *string   `json:"nature"`
				Tags    *[]string `json:"tags"`
				Touches *[]string `json:"touches"`
				Topics  *[]string `json:"topics"` // what a story or epic is about (S-0135)
				After   *[]string `json:"after"`  // stories a story waits for (S-0130)
				Parent  *string   `json:"parent"`
				Body    *string   `json:"body"`
				// Agent replaces the story's agent: absent leaves it, null (or an
				// empty object) removes it (S-0103)
				Agent json.RawMessage `json:"agent"`
				// the planning data (S-0199): draft makes a story a draft or
				// finalizes one; in cost_of_delay and forecast an absent key is
				// left, an empty one removed; a clear removes the block
				Draft       *bool `json:"draft"`
				CostOfDelay *struct {
					RevenuePerWeek   *string `json:"revenue_per_week"`
					PenaltyPerWeek   *string `json:"penalty_per_week"`
					TimeLostPerCycle *string `json:"time_lost_per_cycle"`
					Value            *string `json:"value"`
				} `json:"cost_of_delay"`
				ClearCostOfDelay bool `json:"clear_cost_of_delay"`
				Forecast         *struct {
					Duration *string `json:"duration"`
					Delivery *string `json:"delivery"`
					Basis    *string `json:"basis"`
				} `json:"forecast"`
				ClearForecast bool `json:"clear_forecast"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			// the dashboard always edits what it read: without the hash a change an
			// agent made meanwhile would be overwritten unseen
			if !contentHash.MatchString(in.Hash) {
				return nil, "", bad("hash is required: the one item.show gave for what was edited")
			}
			args := []string{"edit", in.ID, "--hash=" + in.Hash, "--by=" + owner(p), "--autocommit", "--trailer=" + Trailer}
			if in.Title != nil {
				if text(*in.Title) == "" {
					return nil, "", bad("a title is required")
				}
				args = append(args, "--title="+text(*in.Title))
			}
			if in.Nature != nil {
				if !isNature(*in.Nature) {
					return nil, "", bad("%q is not a nature", *in.Nature)
				}
				args = append(args, "--nature="+*in.Nature)
			}
			if in.Parent != nil {
				if e := needID(*in.Parent); e != nil {
					return nil, "", e
				}
				args = append(args, "--parent="+*in.Parent)
			}
			if in.Topics != nil {
				if e := topicWords(*in.Topics); e != nil {
					return nil, "", e
				}
			}
			for _, l := range []struct {
				flag, clear string
				values      *[]string
			}{{"tag", "--clear-tags", in.Tags}, {"touches", "--clear-touches", in.Touches}, {"topics", "--clear-topics", in.Topics}, {"after", "--clear-after", in.After}} {
				if l.values == nil {
					continue
				}
				if len(*l.values) == 0 {
					args = append(args, l.clear)
					continue
				}
				for _, v := range *l.values {
					if e := listEntry(l.flag, v); e != nil {
						return nil, "", e
					}
					args = append(args, "--"+l.flag+"="+strings.TrimSpace(v))
				}
			}
			if len(in.Agent) > 0 {
				var agent *manifest.Agent
				if err := json.Unmarshal(in.Agent, &agent); err != nil {
					return nil, "", bad("agent is an object with harness, model, provider, config, and roles, or null: %s", err)
				}
				more, e := agentArgs(agent)
				if e != nil {
					return nil, "", e
				}
				args = append(append(args, "--clear-agent"), more...)
			}
			switch {
			case in.Draft == nil:
			case *in.Draft:
				args = append(args, "--draft")
			default:
				args = append(args, "--no-draft")
			}
			type key struct {
				flag  string
				value *string
			}
			var keys []key
			if c := in.CostOfDelay; c != nil {
				if in.ClearCostOfDelay {
					return nil, "", bad("clear_cost_of_delay and cost_of_delay contradict each other: give one or the other")
				}
				keys = append(keys, key{"revenue-per-week", c.RevenuePerWeek}, key{"penalty-per-week", c.PenaltyPerWeek}, key{"time-lost-per-cycle", c.TimeLostPerCycle}, key{"cost-of-delay-value", c.Value})
			}
			if in.ClearCostOfDelay {
				args = append(args, "--clear-cost-of-delay")
			}
			if f := in.Forecast; f != nil {
				if in.ClearForecast {
					return nil, "", bad("clear_forecast and forecast contradict each other: give one or the other")
				}
				keys = append(keys, key{"forecast-duration", f.Duration}, key{"forecast-delivery", f.Delivery}, key{"forecast-basis", f.Basis})
			}
			if in.ClearForecast {
				args = append(args, "--clear-forecast")
			}
			// an empty value removes the key; flai checks what is left
			for _, k := range keys {
				if k.value != nil {
					args = append(args, "--"+k.flag+"="+text(*k.value))
				}
			}
			stdin := ""
			if in.Body != nil {
				if strings.TrimSpace(*in.Body) == "" {
					return nil, "", bad("write what the item is for: the body is empty")
				}
				args, stdin = append(args, "--body-stdin"), *in.Body
			}
			if len(args) == 6 {
				return nil, "", bad("nothing to change")
			}
			return args, stdin, nil
		}},

		// item.finalize: the designer's say, from the story page, that a draft
		// story is finished (S-0201). flai edit clears the draft flag and
		// records the owner as who finalized it, and when; a story that is not
		// a draft is refused, since there is nothing to finalize.
		"item.finalize": {exits: map[int]int{4: Refused}, judge: notADraft, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID string `json:"id"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needStory(in.ID); e != nil {
				return nil, "", e
			}
			return []string{"edit", in.ID, "--no-draft", "--by=" + owner(p), "--autocommit", "--trailer=" + Trailer}, "", nil
		}},

		// item.criteria: the designer ticks or unticks an item's acceptance
		// criteria by number from its page (S-0282). flai criteria tick or
		// untick changes those boxes and no other byte, as the owner, committed;
		// the hash is what item.show gave, so a change made meanwhile is a
		// conflict rather than overwritten. One call ticks or unticks, since
		// they are two commands of flai's.
		"item.criteria": {exits: map[int]int{3: Conflict, 4: Refused}, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID     string `json:"id"`
				Hash   string `json:"hash"`
				Tick   []int  `json:"tick"`
				Untick []int  `json:"untick"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			if !contentHash.MatchString(in.Hash) {
				return nil, "", bad("hash is required: the one item.show gave for the criteria ticked")
			}
			verb, ns := "tick", in.Tick
			switch {
			case len(in.Tick) > 0 && len(in.Untick) > 0:
				return nil, "", bad("tick or untick in one call, not both")
			case len(in.Untick) > 0:
				verb, ns = "untick", in.Untick
			case len(in.Tick) == 0:
				return nil, "", bad("name the criteria to tick or untick by their numbers")
			}
			words := make([]string, len(ns))
			for i, n := range ns {
				if n < 1 {
					return nil, "", bad("criteria are numbered from 1: %d is not one", n)
				}
				words[i] = strconv.Itoa(n)
			}
			return []string{"criteria", verb, in.ID, strings.Join(words, ","), "--hash=" + in.Hash, "--by=" + owner(p), "--autocommit", "--trailer=" + Trailer}, "", nil
		}},

		// task.done: a task closed in one call, as flai task done closes it
		// (ADR-0107): commit, sync, move, log, widen touches, check, and inbox,
		// stopping at the first step that fails, with a message needed only
		// when there is something to commit (ADR-0128). The answer is flai
		// task done --json's result. A stop at the sync (exit 3) is a conflict and one at
		// the check (exit 4) a refusal, each with the result as data, as the
		// ADR maps them; a stop at another step (exit 1 with the result) is
		// taskStopped's error, with the result as data too, so the dashboard
		// sees what the steps before it did. Exit 1 without a result is a run
		// that could not start, read as any other failure of flai's.
		"task.done": {exits: map[int]int{3: Conflict, 4: Refused}, answers: []int{1}, judge: taskStopped, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID      string `json:"id"`
				Message string `json:"message"`
				Log     string `json:"log"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needTask(in.ID); e != nil {
				return nil, "", e
			}
			// the message is the commit's, body and all, so its lines are kept;
			// it is needed only when there is something to commit, which flai
			// task done judges, so a blank one is not passed (ADR-0128)
			args := []string{"task", "done", in.ID}
			if strings.TrimSpace(in.Message) != "" {
				args = append(args, "--message="+in.Message)
			}
			if entry := text(in.Log); entry != "" {
				args = append(args, "--log="+entry)
			}
			return args, "", nil
		}},

		"item.template": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Type string `json:"type"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if in.Type != workitem.Epic && in.Type != workitem.Story {
				return nil, "", bad("type must be epic or story")
			}
			return []string{in.Type, "new", "--print-body"}, "", nil
		}),

		// issue.list: every issue, with the stories it names and the open
		// story that links it, or only those a story recorded, read from its
		// worktree while it has one (S-0198).
		"issue.list": read(func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Story string `json:"story"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			args := []string{"issue", "list"}
			if in.Story != "" {
				if e := needStory(in.Story); e != nil {
					return nil, "", e
				}
				args = append(args, "--story="+in.Story)
			}
			return args, "", nil
		}),

		// issue.story: a backlog story made from an open issue, as item.new
		// makes one (S-0198). It reads the issue from the main checkout.
		"issue.story": {exits: map[int]int{4: Refused}, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID   string `json:"id"`
				Epic string `json:"epic"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if !issueID.MatchString(in.ID) {
				return nil, "", bad("%q is not an issue ID", in.ID)
			}
			args := []string{"issue", "story", in.ID, "--owner=" + owner(p), "--autocommit", "--trailer=" + Trailer}
			if in.Epic != "" {
				if !anyItemID.MatchString(in.Epic) || !strings.HasPrefix(in.Epic, "E-") {
					return nil, "", bad("a story's parent must be an epic")
				}
				args = append(args, "--epic="+in.Epic)
			}
			return args, "", nil
		}},

		// issue.new, issue.bump, and issue.close: an issue filed, another
		// occurrence of one recorded, and one closed, from the dashboard
		// (S-0275), as flai issue new, bump, and close do them, in the main
		// checkout and committed there on their own with the dashboard's
		// trailer (--autocommit). Each commits the main checkout and widens
		// no story's touches: --commit, which commits on a story's branch, is
		// the story's agent's, in its worktree. story names the story the
		// occurrence belongs to; close takes none. The answer is the issue
		// as --json gives it, with committed, commit, and commit_error.
		"issue.new": {exits: map[int]int{4: Refused}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				Title string `json:"title"`
				Class string `json:"class"`
				Story string `json:"story"`
				Cost  string `json:"cost"`
				Note  string `json:"note"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			title := text(in.Title)
			if title == "" {
				return nil, "", bad("a title is required: the problem, as a sentence")
			}
			if len([]rune(title)) > 200 {
				return nil, "", bad("the title is longer than 200 characters")
			}
			if !slices.Contains(issues.Classes, in.Class) {
				return nil, "", bad("class must be one of %s", strings.Join(issues.Classes, ", "))
			}
			occurred, e := occurrenceArgs(in.Story, in.Cost, in.Note)
			if e != nil {
				return nil, "", e
			}
			args := append([]string{"issue", "new", "--class=" + in.Class}, occurred...)
			return append(args, "--autocommit", "--trailer="+Trailer, "--", title), "", nil
		}},

		"issue.bump": {exits: map[int]int{4: Refused}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID    string `json:"id"`
				Story string `json:"story"`
				Cost  string `json:"cost"`
				Note  string `json:"note"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if !issueID.MatchString(in.ID) {
				return nil, "", bad("%q is not an issue ID", in.ID)
			}
			occurred, e := occurrenceArgs(in.Story, in.Cost, in.Note)
			if e != nil {
				return nil, "", e
			}
			args := append([]string{"issue", "bump", in.ID}, occurred...)
			return append(args, "--autocommit", "--trailer="+Trailer), "", nil
		}},

		"issue.close": {exits: map[int]int{4: Refused}, build: func(_ channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID     string `json:"id"`
				Reason string `json:"reason"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if !issueID.MatchString(in.ID) {
				return nil, "", bad("%q is not an issue ID", in.ID)
			}
			reason := text(in.Reason)
			if reason == "" {
				return nil, "", bad("a reason is required: what closed it, a story ID, a fix, or why it no longer applies")
			}
			return []string{"issue", "close", in.ID, "--reason=" + reason, "--autocommit", "--trailer=" + Trailer}, "", nil
		}},

		// accept.run is the operator's acceptance, or with by: orchestrator
		// the orchestrator's, with the commit its verifier passed and its
		// evidence (ADR-0093). flai's refusal for lint the fix left (exit 4,
		// ADR-0140) is Refused, with each finding, its file, line, rule, and
		// checkout, as data, so the dashboard opens each file (S-0378).
		// keep_removals are the operator's only (ADR-0155): each a path the
		// branch undoes of main's change and the reason its removal is kept,
		// passed as --keep-removal and --reason pairs in order. flai's refusal
		// for a path still undone (exit 4) is Refused with undoes_main as data.
		// Under pull-request acceptance it merges the story's pull request on
		// the forge (ADR-0170): flai accept runs here, in flai serve's
		// environment, with the forge token's variable, so it merges directly
		// and never hands off. It is refused while the forge host action is
		// off, before flai runs, after orchestratorAccepts.
		// include_uncommitted, which older dashboards send, passes --yes:
		// since ADR-0248 every acceptance commits the main checkout's
		// uncommitted files under wip, design, and the shared paths and
		// leaves the rest as they are, and --yes includes the rest in the
		// acceptance commit too; the dashboard offers no choice.
		"accept.run": {progress: true, exits: map[int]int{4: Refused}, needs: forgeAcceptance, build: func(p channel.Project, raw json.RawMessage) ([]string, string, *channel.Error) {
			in, e := decode[struct {
				ID                 string `json:"id"`
				IncludeUncommitted bool   `json:"include_uncommitted"`
				By                 string `json:"by"`
				Verified           string `json:"verified"`
				Evidence           string `json:"evidence"`
				KeepRemovals       []struct {
					Path   string `json:"path"`
					Reason string `json:"reason"`
				} `json:"keep_removals"`
			}](raw)
			if e != nil {
				return nil, "", e
			}
			if e := needID(in.ID); e != nil {
				return nil, "", e
			}
			var keeps []string
			for _, k := range in.KeepRemovals {
				path, e := undonePath("keep_removals", k.Path)
				if e != nil {
					return nil, "", e
				}
				reason := text(k.Reason)
				if reason == "" {
					return nil, "", bad("the reason for keeping the removal at %s is blank: say why the branch's removal of the main branch's lines is kept", path)
				}
				keeps = append(keeps, "--keep-removal="+path, "--reason="+reason)
			}
			if len(keeps) > 0 && workitem.IsOrchestrator(in.By) {
				return nil, "", bad("keep_removals are the operator's acceptance only: the orchestrator keeps no removal of the main branch's lines, so leave by out (ADR-0155)")
			}
			// Acceptance merges and archives only; it releases nothing (S-0087).
			// Publishing what has accumulated is a step of its own (publish.run).
			args, stdin := append([]string{"accept", in.ID, "--by=" + owner(p)}, keeps...), ""
			switch {
			case workitem.IsOrchestrator(in.By):
				if e := orchestratorAccepts(p, in.ID); e != nil {
					return nil, "", e
				}
				// flai refuses a missing commit or evidence, naming what to give
				args = []string{"accept", in.ID, "--by=" + workitem.ActivityOrchestrator}
				if in.Verified != "" {
					if !commitName.MatchString(in.Verified) {
						return nil, "", bad("verified is the commit the orchestrator's verifier passed, in hex: %q is not one", in.Verified)
					}
					args = append(args, "--verified="+in.Verified)
				}
				if strings.TrimSpace(in.Evidence) != "" {
					args, stdin = append(args, "--evidence=-"), in.Evidence
				}
			case in.By != "":
				return nil, "", bad("by is %s for the orchestrator's acceptance, or left out for the operator's", workitem.ActivityOrchestrator)
			case in.Verified != "" || in.Evidence != "":
				return nil, "", bad("verified and evidence are the orchestrator's acceptance: give them with by %s", workitem.ActivityOrchestrator)
			}
			if in.IncludeUncommitted {
				args = append(args, "--yes")
			}
			return args, stdin, nil
		}},
	}
}
