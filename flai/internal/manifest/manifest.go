// Package manifest reads system-flow.yaml, the file that marks a conforming
// project. See design/system/project-manifest.md and ADR-0011.
package manifest

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"math"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-yaml"

	"github.com/bytepunx/system-flow/flai/internal/buildinfo"
	"github.com/bytepunx/system-flow/flai/internal/cron"
)

// File is the manifest file name at a project root.
const File = "system-flow.yaml"

// ErrNotFound means no manifest exists in the directory or any parent.
var ErrNotFound = errors.New("no " + File + " found in this directory or any parent")

// Manifest is the schema of system-flow.yaml.
type Manifest struct {
	Version     int               `yaml:"version"`
	Name        string            `yaml:"name"`
	Key         string            `yaml:"key"`
	Description string            `yaml:"description"`
	Owner       string            `yaml:"owner"`
	Repo        string            `yaml:"repo"`
	Template    Template          `yaml:"template"`
	Layout      map[string]string `yaml:"layout"`
	Projects    []Project         `yaml:"projects"`
	Dashboard   Dashboard         `yaml:"dashboard"`
	// Checks are the commands a story in review is checked with, in the
	// story's worktree, in order (S-0082). Used only when the host's own
	// configuration names none: that is the operator's choice of where to
	// name them, not a merge of both.
	Checks []NamedCommand `yaml:"checks,omitempty" json:"checks,omitempty"`
	// Tests are the project's test tiers, cheapest first, that flai test runs
	// for the paths each selects (S-0273). Nil, with no tests key, means the
	// default TestTiers gives; an empty list means none.
	Tests []TestTier `yaml:"tests,omitempty" json:"tests,omitempty"`
	// Agent is the project's default agent, copied into every story created
	// while it is set (S-0103). flai agent sets it.
	Agent *Agent `yaml:"agent,omitempty" json:"agent,omitempty"`
	// Prime is how flai prime --story builds this project's context packs
	// (ADR-0049).
	Prime Prime `yaml:"prime,omitempty" json:"prime,omitzero"`
	// Issues is how flai check treats the project's open issues (S-0198).
	Issues Issues `yaml:"issues,omitempty" json:"issues,omitzero"`
	// Planning is the units the planner's numbers are in (S-0199), and the
	// planner's agent (S-0208).
	Planning Planning `yaml:"planning,omitempty" json:"planning,omitzero"`
	// Orchestration is the policy that orders the ready column and the
	// policy that says when to release (S-0217).
	Orchestration Orchestration `yaml:"orchestration,omitempty" json:"orchestration,omitzero"`
	// Analysis is the analyzer's agent and schedule (S-0223).
	Analysis Analysis `yaml:"analysis,omitempty" json:"analysis,omitzero"`
	// Claims is how stories' touches claim paths: the shared paths whose
	// overlaps hold no story (S-0295, ADR-0096). A pattern that is not valid
	// does not stop the load; Claims.Errors names it, and it frees nothing.
	Claims Claims `yaml:"claims,omitempty" json:"claims,omitzero"`
	// Flai is what the project asks of the flai that reads it (S-0181).
	Flai Requirement `yaml:"flai,omitempty" json:"flai,omitzero"`
}

// Requirement is the oldest flai that may read the project.
type Requirement struct {
	// Minimum is a flai release, X.Y.Z: one that knows every front-matter
	// field the project's items carry. Publishing a flai release whose
	// front-matter fields changed raises it (release.RaiseMinimum).
	Minimum string `yaml:"minimum,omitempty" json:"minimum,omitempty"`
}

// TooOldError is a manifest whose minimum flai is newer than the running one.
type TooOldError struct {
	Path    string
	Minimum string
	Running string
}

func (e *TooOldError) Error() string {
	return fmt.Sprintf("%s needs flai %s or newer, and this is flai %s, which may not know the fields its items carry: upgrade it with %s, then run this again",
		e.Path, e.Minimum, e.Running, buildinfo.UpgradeCommand)
}

// Prime is the project's say about its context packs.
type Prime struct {
	// Budget is the size a story's context pack fits, such as 80KB or
	// 81920 (bytes); empty means flai's default.
	Budget string `yaml:"budget,omitempty" json:"budget,omitempty"`
}

// Issues is the project's say about its open issues.
type Issues struct {
	// StoryAfter is how long an issue may stay open with no open story
	// linking it before flai check warns, a Go duration such as 168h; empty
	// means DefaultStoryAfter, and 0 turns the warning off.
	StoryAfter string `yaml:"story_after,omitempty" json:"story_after,omitempty"`
}

// DefaultStoryAfter is how long an issue may stay open with no open story
// linking it when issues.story_after is not set: 7 days.
const DefaultStoryAfter = 168 * time.Hour

// StoryAfterDuration is issues.story_after as a duration: DefaultStoryAfter
// when it is empty, and zero when the warning is off.
func (i Issues) StoryAfterDuration() (time.Duration, error) {
	s := strings.TrimSpace(i.StoryAfter)
	if s == "" {
		return DefaultStoryAfter, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("issues.story_after %q is not a duration; write one such as 168h or 24h, or 0 to turn the warning off", i.StoryAfter)
	}
	if d < 0 {
		return 0, fmt.Errorf("issues.story_after %q is negative; write a duration such as 168h or 24h, or 0 to turn the warning off", i.StoryAfter)
	}
	return d, nil
}

// Planning is the project's say about the units of its planning data and
// the agent that plans.
type Planning struct {
	// Currency is the ISO 4217 code of every amount items carry, such as
	// EUR; empty means DefaultCurrency.
	Currency string `yaml:"currency,omitempty" json:"currency,omitempty"`
	// HourRate is what an hour of work costs, in Currency; unset means
	// unknown, not free.
	HourRate *float64 `yaml:"hour_rate,omitempty" json:"hour_rate,omitempty"`
	// Cycle is the period a cost of delay's time lost is counted over, a Go
	// duration such as 168h; empty means DefaultCycle.
	Cycle string `yaml:"cycle,omitempty" json:"cycle,omitempty"`
	// DefaultDuration is the work a story is forecast to take when there is
	// too little history to forecast it from (S-0210), a Go duration such as
	// 2h; empty means DefaultDuration.
	DefaultDuration string `yaml:"default_duration,omitempty" json:"default_duration,omitempty"`
	// Agent is the planner's agent over the project's (S-0208): what it sets
	// wins, and what it leaves out is the project's agent's.
	Agent *Agent `yaml:"agent,omitempty" json:"agent,omitempty"`
	// Replan is what flai serve does when a story is accepted or cancelled
	// or the pull order changes (S-0211, ADR-0084): ReplanNever,
	// ReplanDeterministic, or ReplanAgent; empty means ReplanDeterministic.
	Replan string `yaml:"replan,omitempty" json:"replan,omitempty"`
	// Schedule is when flai serve runs the planner over the ready column
	// (S-0211, ADR-0084): a five-field cron expression in UTC or daily; empty
	// means no schedule.
	Schedule string `yaml:"schedule,omitempty" json:"schedule,omitempty"`
}

// The values of planning.replan.
const (
	// ReplanNever does nothing when work ahead completes or the order
	// changes.
	ReplanNever = "never"
	// ReplanDeterministic plays the board out again and moves each forecast
	// delivery that changed, with no agent. It is the default.
	ReplanDeterministic = "deterministic"
	// ReplanAgent does what ReplanDeterministic does and queues the planner
	// for each story whose delivery moved.
	ReplanAgent = "agent"
)

// DefaultCurrency is the currency of amounts when planning.currency is not
// set.
const DefaultCurrency = "USD"

// DefaultCycle is the planning cycle when planning.cycle is not set: a week.
const DefaultCycle = 168 * time.Hour

// DefaultDuration is a story's forecast duration when there is too little
// history and planning.default_duration is not set: an hour.
const DefaultDuration = time.Hour

var currencyCode = regexp.MustCompile(`^[A-Z]{3}$`)

// CurrencyCode is planning.currency, or DefaultCurrency when it is empty.
func (p Planning) CurrencyCode() string {
	if p.Currency != "" {
		return p.Currency
	}
	return DefaultCurrency
}

// CycleDuration is planning.cycle as a duration: DefaultCycle when it is
// empty.
func (p Planning) CycleDuration() (time.Duration, error) {
	s := strings.TrimSpace(p.Cycle)
	if s == "" {
		return DefaultCycle, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("planning.cycle %q is not a duration longer than zero; write one such as 168h or 336h, or remove it for a week", p.Cycle)
	}
	return d, nil
}

// FallbackDuration is planning.default_duration as a duration:
// DefaultDuration when it is empty.
func (p Planning) FallbackDuration() (time.Duration, error) {
	s := strings.TrimSpace(p.DefaultDuration)
	if s == "" {
		return DefaultDuration, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("planning.default_duration %q is not a duration longer than zero; write one such as 1h or 90m, or remove it for an hour", p.DefaultDuration)
	}
	return d, nil
}

// ReplanPolicy is planning.replan: ReplanDeterministic when it is empty.
func (p Planning) ReplanPolicy() (string, error) {
	switch s := strings.TrimSpace(p.Replan); s {
	case "":
		return ReplanDeterministic, nil
	case ReplanNever, ReplanDeterministic, ReplanAgent:
		return s, nil
	}
	return "", fmt.Errorf("planning.replan %q is not a replan policy; write never (do nothing), deterministic (play the board out again and move forecast deliveries), or agent (do that and queue the planner for each story whose delivery moved), or remove it for deterministic", p.Replan)
}

// PlanSchedule is planning.schedule parsed: nil when it is empty.
func (p Planning) PlanSchedule() (*cron.Schedule, error) {
	return parseSchedule("planning", p.Schedule)
}

// parseSchedule is the schedule key under block parsed: nil when it is empty.
func parseSchedule(block, spec string) (*cron.Schedule, error) {
	s := strings.TrimSpace(spec)
	if s == "" {
		return nil, nil
	}
	sched, err := cron.Parse(s)
	if err != nil {
		// cron's refusals begin with the schedule quoted: name the key.
		return nil, fmt.Errorf("%s.%w", block, err)
	}
	return &sched, nil
}

// Errors are what is wrong with the planning settings, one sentence each;
// none when they are valid.
func (p Planning) Errors() []string {
	var errs []string
	if p.Currency != "" && !currencyCode.MatchString(p.Currency) {
		errs = append(errs, fmt.Sprintf("planning.currency %q is not an ISO 4217 code; write three capital letters such as USD or EUR, or remove it for USD", p.Currency))
	}
	if r := p.HourRate; r != nil && (math.IsNaN(*r) || math.IsInf(*r, 0) || *r < 0) {
		errs = append(errs, fmt.Sprintf("planning.hour_rate %v is not an amount of zero or more; write what an hour of work costs in the project's currency, or remove it", *r))
	}
	if _, err := p.CycleDuration(); err != nil {
		errs = append(errs, err.Error())
	}
	if _, err := p.FallbackDuration(); err != nil {
		errs = append(errs, err.Error())
	}
	if _, err := p.ReplanPolicy(); err != nil {
		errs = append(errs, err.Error())
	}
	if _, err := p.PlanSchedule(); err != nil {
		errs = append(errs, err.Error())
	}
	return append(errs, p.Agent.problems("planning.agent")...)
}

// PlanningAgent is the agent the planner is started with: planning.agent
// merged over the project's agent, as a story's agent is merged over the
// default (ADR-0037, ADR-0065). Nil when neither sets anything.
func (m Manifest) PlanningAgent() *Agent {
	return m.Agent.With(m.Planning.Agent)
}

// Analysis is the project's say about the analyzer: the agent it runs as
// and when flai serve runs it (S-0223).
type Analysis struct {
	// Agent is the analyzer's agent over the project's: what it sets wins,
	// and what it leaves out is the project's agent's.
	Agent *Agent `yaml:"agent,omitempty" json:"agent,omitempty"`
	// Schedule is when flai serve runs the analyzer: a five-field cron
	// expression in UTC or daily; empty means no schedule.
	Schedule string `yaml:"schedule,omitempty" json:"schedule,omitempty"`
}

// AnalysisSchedule is analysis.schedule parsed: nil when it is empty.
func (a Analysis) AnalysisSchedule() (*cron.Schedule, error) {
	return parseSchedule("analysis", a.Schedule)
}

// Errors are what is wrong with the analysis settings, one sentence each;
// none when they are valid.
func (a Analysis) Errors() []string {
	var errs []string
	if _, err := a.AnalysisSchedule(); err != nil {
		errs = append(errs, err.Error())
	}
	return append(errs, a.Agent.problems("analysis.agent")...)
}

// AnalysisAgent is the agent the analyzer is started with: analysis.agent
// merged over the project's agent, as PlanningAgent is. Nil when neither
// sets anything.
func (m Manifest) AnalysisAgent() *Agent {
	return m.Agent.With(m.Analysis.Agent)
}

// Orchestration is the project's say about the order of the ready column and
// when accepted work is released (S-0217), and about what the orchestrator
// may do and the agent it runs as (S-0218).
type Orchestration struct {
	// Policy orders the ready column: one of OrderPolicies, the names flai
	// order --by takes; empty means OrderFIFO, which leaves the operator's
	// order alone.
	Policy string `yaml:"policy,omitempty" json:"policy,omitempty"`
	// Release says when accepted stories not yet released are due a release.
	Release Release `yaml:"release,omitempty" json:"release,omitzero"`
	// Permissions are what the orchestrator may do without the operator.
	Permissions Permissions `yaml:"permissions,omitempty" json:"permissions,omitzero"`
	// Agent is the orchestrator's agent over the project's: what it sets
	// wins, and what it leaves out is the project's agent's.
	Agent *Agent `yaml:"agent,omitempty" json:"agent,omitempty"`
}

// Permissions are what the orchestrator may do without the operator, each
// off when unset (S-0218). flai guard holds an orchestrator session to them.
type Permissions struct {
	// PlanBacklogEpics lets it ask the planner to draft stories for an epic
	// in the backlog.
	PlanBacklogEpics bool `yaml:"plan_backlog_epics,omitempty" json:"plan_backlog_epics,omitempty"`
	// FinalizeDrafts lets it finalize a draft story.
	FinalizeDrafts bool `yaml:"finalize_drafts,omitempty" json:"finalize_drafts,omitempty"`
	// PromoteToReady lets it move a story to ready.
	PromoteToReady bool `yaml:"promote_to_ready,omitempty" json:"promote_to_ready,omitempty"`
	// OrderReady lets it write the order of the ready column.
	OrderReady bool `yaml:"order_ready,omitempty" json:"order_ready,omitempty"`
	// AnswerThreads is how it answers threads: one of AnswerModes; empty
	// means AnswerOff.
	AnswerThreads string `yaml:"answer_threads,omitempty" json:"answer_threads,omitempty"`
	// AcceptReviews lets it accept a story in review.
	AcceptReviews bool `yaml:"accept_reviews,omitempty" json:"accept_reviews,omitempty"`
	// Publish lets it release and push accepted work.
	Publish bool `yaml:"publish,omitempty" json:"publish,omitempty"`

	// unknown are the keys under orchestration.permissions that name no
	// permission, for Errors.
	unknown []string
}

// The permissions, as orchestration.permissions keys them.
const (
	PermitPlanBacklogEpics = "plan_backlog_epics"
	PermitFinalizeDrafts   = "finalize_drafts"
	PermitPromoteToReady   = "promote_to_ready"
	PermitOrderReady       = "order_ready"
	PermitAnswerThreads    = "answer_threads"
	PermitAcceptReviews    = "accept_reviews"
	PermitPublish          = "publish"
)

// PermissionNames are the keys of orchestration.permissions, in the order
// they are listed to the operator.
var PermissionNames = []string{PermitPlanBacklogEpics, PermitFinalizeDrafts, PermitPromoteToReady, PermitOrderReady, PermitAnswerThreads, PermitAcceptReviews, PermitPublish}

// The values of orchestration.permissions.answer_threads.
const (
	// AnswerOff leaves threads to the operator. It is the default.
	AnswerOff = "off"
	// AnswerRecommend replies with a recommendation for the operator.
	AnswerRecommend = "recommend"
	// AnswerAutonomous answers threads itself.
	AnswerAutonomous = "autonomous"
)

// AnswerModes are the values of orchestration.permissions.answer_threads.
var AnswerModes = []string{AnswerOff, AnswerRecommend, AnswerAutonomous}

// UnmarshalYAML reads the permissions and keeps the keys that name none, so
// that Errors can name a misspelt one rather than leave it silently off.
func (p *Permissions) UnmarshalYAML(unmarshal func(any) error) error {
	type plain Permissions
	if err := unmarshal((*plain)(p)); err != nil {
		return err
	}
	var keys map[string]any
	if err := unmarshal(&keys); err != nil {
		return err
	}
	p.unknown = nil
	for _, k := range slices.Sorted(maps.Keys(keys)) {
		if !slices.Contains(PermissionNames, k) {
			p.unknown = append(p.unknown, k)
		}
	}
	return nil
}

// AnswerMode is answer_threads, or AnswerOff when it is empty. It does not
// say whether the value is one of AnswerModes; Errors does.
func (p Permissions) AnswerMode() string {
	if s := strings.TrimSpace(p.AnswerThreads); s != "" {
		return s
	}
	return AnswerOff
}

// Allows reports whether the permission name, one of PermissionNames, is on:
// a boolean set true, or answer_threads set to recommend or autonomous
// (AnswerMode tells them apart). A name that is no permission is off.
func (p Permissions) Allows(name string) bool {
	switch name {
	case PermitPlanBacklogEpics:
		return p.PlanBacklogEpics
	case PermitFinalizeDrafts:
		return p.FinalizeDrafts
	case PermitPromoteToReady:
		return p.PromoteToReady
	case PermitOrderReady:
		return p.OrderReady
	case PermitAnswerThreads:
		m := p.AnswerMode()
		return m == AnswerRecommend || m == AnswerAutonomous
	case PermitAcceptReviews:
		return p.AcceptReviews
	case PermitPublish:
		return p.Publish
	}
	return false
}

func (p Permissions) errors() []string {
	var errs []string
	for _, k := range p.unknown {
		errs = append(errs, fmt.Sprintf("orchestration.permissions has no permission %q; write one of %s, or remove it", k, strings.Join(PermissionNames, ", ")))
	}
	if m := p.AnswerMode(); !slices.Contains(AnswerModes, m) {
		errs = append(errs, fmt.Sprintf("orchestration.permissions.answer_threads %q is not a way of answering threads; write off (leave them to the operator), recommend (reply with a recommendation), or autonomous (answer them), or remove it for off", p.AnswerThreads))
	}
	return errs
}

// OrchestrationAgent is the agent the orchestrator is started with:
// orchestration.agent merged over the project's agent, as PlanningAgent is.
// Nil when neither sets anything.
func (m Manifest) OrchestrationAgent() *Agent {
	return m.Agent.With(m.Orchestration.Agent)
}

// The values of orchestration.policy, the same as flai order --by's.
const (
	// OrderCOD orders by cost of delay, the largest first.
	OrderCOD = "cod"
	// OrderWSJF orders by cost of delay divided by forecast duration, the
	// largest first.
	OrderWSJF = "wsjf"
	// OrderThroughput orders by forecast duration, the shortest first.
	OrderThroughput = "throughput"
	// OrderFIFO leaves the operator's order alone. It is the default.
	OrderFIFO = "fifo"
)

// OrderPolicies are the values of orchestration.policy and flai order --by,
// in the order they are listed to the operator.
var OrderPolicies = []string{OrderCOD, OrderWSJF, OrderThroughput, OrderFIFO}

// Release is the project's say about when accepted stories not yet released
// are due a release.
type Release struct {
	// Policy is one of ReleasePolicies; empty means ReleaseJudgement.
	Policy string `yaml:"policy,omitempty" json:"policy,omitempty"`
	// Value is, for ReleaseThreshold, the unreleased cost of delay per week,
	// in the project's currency, at which a release is due; unset means none.
	Value *float64 `yaml:"value,omitempty" json:"value,omitempty"`
	// Count is, for ReleaseThreshold, the number of accepted stories not yet
	// released at which a release is due; unset means none.
	Count *int `yaml:"count,omitempty" json:"count,omitempty"`
	// Epic is, for ReleaseTheme, the epic whose stories, every one accepted,
	// make a release due.
	Epic string `yaml:"epic,omitempty" json:"epic,omitempty"`
	// Tag is, for ReleaseTheme, the tag whose stories, every one accepted,
	// make a release due.
	Tag string `yaml:"tag,omitempty" json:"tag,omitempty"`
	// WholeEpics, under every policy, holds a release back while a story
	// accepted and not yet released belongs to an epic in neither review nor
	// done (S-0222). Off when unset.
	WholeEpics bool `yaml:"whole_epics,omitempty" json:"whole_epics,omitempty"`
}

// The values of orchestration.release.policy.
const (
	// ReleaseJudgement leaves the release to the operator: it is never met by
	// itself. It is the default.
	ReleaseJudgement = "judgement"
	// ReleaseThreshold is met by the unreleased cost of delay per week against
	// Value, or the accepted stories not yet released against Count.
	ReleaseThreshold = "threshold"
	// ReleaseTheme is met when every story of an epic or a tag is accepted.
	ReleaseTheme = "theme"
)

// ReleasePolicies are the values of orchestration.release.policy.
var ReleasePolicies = []string{ReleaseJudgement, ReleaseThreshold, ReleaseTheme}

var epicID = regexp.MustCompile(`^E-\d{3,}$`)

// PolicyOrDefault is orchestration.policy, or OrderFIFO when it is empty. It
// does not say whether the policy is one of OrderPolicies; Errors does.
func (o Orchestration) PolicyOrDefault() string {
	if s := strings.TrimSpace(o.Policy); s != "" {
		return s
	}
	return OrderFIFO
}

// PolicyOrDefault is orchestration.release.policy, or ReleaseJudgement when
// it is empty. It does not say whether the policy is one of ReleasePolicies;
// Errors does.
func (r Release) PolicyOrDefault() string {
	if s := strings.TrimSpace(r.Policy); s != "" {
		return s
	}
	return ReleaseJudgement
}

// Errors are what is wrong with the orchestration settings, one sentence
// each; none when they are valid.
func (o Orchestration) Errors() []string {
	var errs []string
	if p := o.PolicyOrDefault(); !slices.Contains(OrderPolicies, p) {
		errs = append(errs, fmt.Sprintf("orchestration.policy %q is not an order policy; write cod (cost of delay), wsjf (cost of delay by duration), throughput (shortest first), or fifo (the operator's order), or remove it for fifo", o.Policy))
	}
	errs = append(errs, o.Release.errors()...)
	errs = append(errs, o.Permissions.errors()...)
	return append(errs, o.Agent.problems("orchestration.agent")...)
}

func (r Release) errors() []string {
	var errs []string
	switch r.PolicyOrDefault() {
	case ReleaseJudgement:
	case ReleaseThreshold:
		if r.Value == nil && r.Count == nil {
			errs = append(errs, "orchestration.release is a threshold with neither value nor count; write value (the unreleased cost of delay per week), count (the accepted stories not yet released), or both")
		}
		if v := r.Value; v != nil && (math.IsNaN(*v) || math.IsInf(*v, 0) || *v < 0) {
			errs = append(errs, fmt.Sprintf("orchestration.release.value %v is not an amount of zero or more; write the unreleased cost of delay per week, in the project's currency, at which a release is due", *v))
		}
		if c := r.Count; c != nil && *c < 0 {
			errs = append(errs, fmt.Sprintf("orchestration.release.count %d is not a number of zero or more; write the number of accepted stories not yet released at which a release is due", *c))
		}
	case ReleaseTheme:
		epic, tag := strings.TrimSpace(r.Epic), strings.TrimSpace(r.Tag)
		switch {
		case epic == "" && tag == "":
			errs = append(errs, "orchestration.release is a theme with neither epic nor tag; write epic (an epic ID such as E-0001) or tag (a tag), one of them")
		case epic != "" && tag != "":
			errs = append(errs, "orchestration.release is a theme with both epic and tag; keep one of them")
		case epic != "" && !epicID.MatchString(epic):
			errs = append(errs, fmt.Sprintf("orchestration.release.epic %q is not an epic ID; write one such as E-0001", r.Epic))
		}
	default:
		errs = append(errs, fmt.Sprintf("orchestration.release.policy %q is not a release policy; write judgement (the operator decides), threshold (a value or a count of unreleased work), or theme (an epic or a tag accepted), or remove it for judgement", r.Policy))
	}
	return errs
}

// NamedCommand is one command by name: an argument list, run as it stands,
// never through a shell. In an argument, {story} is replaced by the story's
// ID and {root} by the directory it runs in; nothing else is interpreted.
type NamedCommand struct {
	Name    string   `yaml:"name" json:"name"`
	Command []string `yaml:"command" json:"command"`
}

// TestTier is one tier of the project's tests: a command, the paths that
// select it, and how its output becomes findings (S-0273).
type TestTier struct {
	// Name names the tier, unique among the project's tiers.
	Name string `yaml:"name" json:"name"`
	// Command is an argument list, run as it stands, never through a shell.
	// An argument that is PlaceholderPackages or PlaceholderFiles alone
	// stands for what the tier's paths selected; nothing else is
	// interpreted.
	Command []string `yaml:"command" json:"command"`
	// Dir is the folder the command runs in, relative to the repository
	// root; empty means the root.
	Dir string `yaml:"dir,omitempty" json:"dir,omitempty"`
	// Paths are glob patterns relative to the repository root, as
	// claims.shared's are, that select the tier for a changed or named path;
	// one beginning with ! takes paths out again. Required unless AllOnly.
	Paths []string `yaml:"paths,omitempty" json:"paths,omitempty"`
	// Format is how the command's output becomes findings: one of
	// TestFormats; empty means FormatPlain.
	Format string `yaml:"format,omitempty" json:"format,omitempty"`
	// AllOnly runs the tier only when every tier is asked for, as flai test
	// --all does, never for paths.
	AllOnly bool `yaml:"all_only,omitempty" json:"all_only,omitempty"`
	// AllCommand, when set, runs instead of Command when every tier is asked
	// for. It takes no placeholders.
	AllCommand []string `yaml:"all_command,omitempty" json:"all_command,omitempty"`
}

// The placeholders a tier's command may hold, each as an argument of its own.
const (
	// PlaceholderPackages stands for the packages of the paths the tier
	// selected.
	PlaceholderPackages = "{packages}"
	// PlaceholderFiles stands for the files the tier selected.
	PlaceholderFiles = "{files}"
)

// The values of a tier's format.
const (
	// FormatGoTestJSON is go test -json's events.
	FormatGoTestJSON = "go-test-json"
	// FormatVitestJSON is vitest's JSON reporter.
	FormatVitestJSON = "vitest-json"
	// FormatGolangciJSON is golangci-lint's JSON output.
	FormatGolangciJSON = "golangci-json"
	// FormatGofmtList is gofmt -l's list of files not formatted.
	FormatGofmtList = "gofmt-list"
	// FormatPlain is any output; the exit status alone says whether the tier
	// passed. It is the default.
	FormatPlain = "plain"
)

// TestFormats are the values of a tier's format, in the order they are
// listed to the operator.
var TestFormats = []string{FormatGoTestJSON, FormatVitestJSON, FormatGolangciJSON, FormatGofmtList, FormatPlain}

// DefaultTestScript is the script the default tier runs when the manifest
// has no tests key.
const DefaultTestScript = "scripts/test.sh"

// placeholder is an argument that looks like a placeholder: a word in
// braces.
var placeholder = regexp.MustCompile(`^\{[A-Za-z_][A-Za-z0-9_]*\}$`)

// FormatOrDefault is the tier's format, or FormatPlain when it is empty. It
// does not say whether the format is one of TestFormats; TestErrors does.
func (t TestTier) FormatOrDefault() string {
	if s := strings.TrimSpace(t.Format); s != "" {
		return s
	}
	return FormatPlain
}

// TestTiers are the project's test tiers, cheapest first: the tests key's
// when it is there, and otherwise one plain tier named test that runs
// DefaultTestScript for every path, when fsys, the repository root, holds
// that file; none when it does not. It refuses tiers TestErrors finds wrong,
// with every problem.
func (m Manifest) TestTiers(fsys fs.FS) ([]TestTier, error) {
	if m.Tests != nil {
		if errs := m.TestErrors(); len(errs) > 0 {
			return nil, errors.New(strings.Join(errs, "; "))
		}
		return m.Tests, nil
	}
	if !isFile(fsys, DefaultTestScript) {
		return nil, nil
	}
	return []TestTier{{Name: "test", Command: []string{DefaultTestScript}, Paths: []string{"**"}, Format: FormatPlain}}, nil
}

// isFile reports whether name is a file in fsys; anything it cannot stat is
// not.
func isFile(fsys fs.FS, name string) bool {
	st, err := fs.Stat(fsys, name)
	return err == nil && !st.IsDir()
}

// TestErrors are what is wrong with the tests key, one sentence each, each
// beginning with the tier's index and field, such as tests[1].format; none
// when every tier is valid.
func (m Manifest) TestErrors() []string {
	var errs []string
	seen := map[string]int{}
	for i, t := range m.Tests {
		at := func(field string) string {
			s := fmt.Sprintf("tests[%d].%s", i, field)
			if name := strings.TrimSpace(t.Name); name != "" && field != "name" {
				s += fmt.Sprintf(" (tier %q)", name)
			}
			return s
		}
		name := strings.TrimSpace(t.Name)
		if name == "" {
			errs = append(errs, at("name")+" is empty; give the tier a name, such as test or unit")
		} else if j, ok := seen[name]; ok {
			errs = append(errs, fmt.Sprintf("%s %q is tests[%d]'s name too; give each tier a name of its own", at("name"), t.Name, j))
		} else {
			seen[name] = i
		}
		if len(t.Command) == 0 {
			errs = append(errs, at("command")+" is empty; write the program and its arguments as a list, such as [scripts/test.sh] or [go, test, \"{packages}\"]")
		} else {
			errs = append(errs, commandProblems(at("command"), t.Command, true)...)
		}
		if len(t.AllCommand) > 0 {
			errs = append(errs, commandProblems(at("all_command"), t.AllCommand, false)...)
		}
		if r := dirProblem(t.Dir); r != "" {
			errs = append(errs, fmt.Sprintf("%s %q %s", at("dir"), t.Dir, r))
		}
		selects := 0
		for _, p := range t.Paths {
			glob, out := strings.CutPrefix(p, "!")
			if r := patternProblem(glob); r != "" {
				errs = append(errs, fmt.Sprintf("%s pattern %q %s", at("paths"), p, r))
			} else if !out {
				selects++
			}
		}
		switch {
		case t.AllOnly || selects > 0:
		case len(t.Paths) == 0:
			errs = append(errs, at("paths")+" is empty; write the globs that select the tier, such as flai/**, or set all_only for a tier that runs only when every tier is asked for")
		default:
			errs = append(errs, at("paths")+" has no pattern that selects a path, only ones beginning with ! or ones not valid; add one, such as flai/**")
		}
		if f := t.FormatOrDefault(); !slices.Contains(TestFormats, f) {
			errs = append(errs, fmt.Sprintf("%s %q is not a format flai reads; write %s, or remove it for plain", at("format"), t.Format, strings.Join(TestFormats, ", ")))
		}
	}
	return errs
}

// commandProblems are what is wrong with a tier's argument list, each
// sentence beginning with field. A placeholder is allowed only where
// placeholders is set, and only as an argument of its own.
func commandProblems(field string, args []string, placeholders bool) []string {
	var errs []string
	if strings.TrimSpace(args[0]) == "" {
		errs = append(errs, field+" has no program first; write the program, then its arguments, such as [scripts/test.sh]")
	} else if placeholder.MatchString(args[0]) {
		errs = append(errs, fmt.Sprintf("%s begins with the placeholder %s; write the program first, then its arguments", field, args[0]))
	}
	for _, a := range args[1:] {
		switch {
		case (a == PlaceholderPackages || a == PlaceholderFiles) && !placeholders:
			errs = append(errs, fmt.Sprintf("%s argument %s is a placeholder, and all_command runs for no paths; write it without %s or %s", field, a, PlaceholderPackages, PlaceholderFiles))
		case a == PlaceholderPackages || a == PlaceholderFiles:
		case placeholder.MatchString(a):
			errs = append(errs, fmt.Sprintf("%s argument %s is not a placeholder flai fills; write %s or %s, or remove it", field, a, PlaceholderPackages, PlaceholderFiles))
		case strings.Contains(a, PlaceholderPackages) || strings.Contains(a, PlaceholderFiles):
			errs = append(errs, fmt.Sprintf("%s argument %q holds a placeholder inside it, which flai does not fill; write %s or %s as an argument of its own", field, a, PlaceholderPackages, PlaceholderFiles))
		}
	}
	return errs
}

// dirProblem says what is wrong with a tier's dir, or "" when it is valid: a
// clean path inside the repository, relative to its root, or empty.
func dirProblem(d string) string {
	switch {
	case d == "":
		return ""
	case strings.HasPrefix(d, "/") || strings.HasPrefix(d, `\`) || len(d) > 1 && d[1] == ':':
		return "is an absolute path; write it relative to the repository root, such as flai"
	case d == ".." || strings.HasPrefix(d, "../"):
		return "leaves the repository; write a folder inside it, relative to its root, such as flai"
	case path.Clean(d) != d:
		return "is not a clean path; write it with single slashes, no . or .. segments, and none at the end, such as flai/web"
	}
	return ""
}

// Template records which template produced the project.
type Template struct {
	Repo    string `yaml:"repo"`
	Ref     string `yaml:"ref"`
	Version string `yaml:"version"`
	Applied string `yaml:"applied"`
}

// Project is a releasable component at the repo root: a code sub-project
// or the template. Tags are story tags that mean "this story delivers to
// this component" (for example cli for flai).
type Project struct {
	Name string   `yaml:"name" json:"name"`
	Path string   `yaml:"path" json:"path"`
	Kind string   `yaml:"kind" json:"kind"`
	Tags []string `yaml:"tags,omitempty" json:"tags,omitempty"`
}

// Dashboard is how the flaiover image runs for this project.
type Dashboard struct {
	Image string `yaml:"image"`
	Tag   string `yaml:"tag"`
	Port  int    `yaml:"port"`
	Bind  string `yaml:"bind"` // host address to publish on; default 0.0.0.0
	// Autocommit commits documents saved from the dashboard; nil means true (ADR-0023).
	Autocommit *bool `yaml:"autocommit,omitempty"`
	// NotifyURL, when set, is where the dashboard's server posts new inbox
	// entries (S-0042). flai itself never calls it.
	NotifyURL string `yaml:"notify_url,omitempty"`
}

// Autocommit reports whether documents saved from the dashboard are committed.
func (m Manifest) Autocommit() bool {
	return m.Dashboard.Autocommit == nil || *m.Dashboard.Autocommit
}

// Load parses the manifest at path.
func Load(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return Manifest{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if m.Version == 0 {
		return Manifest{}, fmt.Errorf("%s: version is required", path)
	}
	if m.Name == "" {
		return Manifest{}, fmt.Errorf("%s: name is required", path)
	}
	for _, k := range []string{"design", "docs", "wip"} {
		if m.Layout[k] == "" {
			return Manifest{}, fmt.Errorf("%s: layout.%s is required", path, k)
		}
	}
	if min := m.Flai.Minimum; min != "" {
		if _, ok := buildinfo.Semver(min); !ok {
			return Manifest{}, fmt.Errorf("%s: flai.minimum %q is not a release version like 1.27.0", path, min)
		}
		if buildinfo.Below(buildinfo.Version, min) {
			return Manifest{}, &TooOldError{Path: path, Minimum: min, Running: buildinfo.Version}
		}
	}
	return m, nil
}

// SetMinimum writes flai.minimum into the manifest at path, keeping the rest
// of the file as written: the minimum line under a flai: block is replaced
// or added, and a manifest with no flai: block gets one at its end.
func SetMinimum(path, version string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	line := "  minimum: " + version
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	block := -1
	for i, l := range lines {
		if strings.TrimRight(l, " ") == "flai:" {
			block = i
			break
		}
	}
	if block < 0 {
		lines = append(lines, "flai:", line)
	} else {
		at := -1
		for i := block + 1; i < len(lines) && (strings.HasPrefix(lines[i], " ") || lines[i] == ""); i++ {
			if strings.HasPrefix(strings.TrimSpace(lines[i]), "minimum:") {
				at = i
				break
			}
		}
		if at >= 0 {
			lines[at] = line
		} else {
			lines = append(lines[:block+1], append([]string{line}, lines[block+1:]...)...)
		}
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}

// Find walks up from start looking for the manifest and returns its path.
func Find(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		p := filepath.Join(dir, File)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrNotFound
		}
		dir = parent
	}
}

// Dir returns the folder for a layout key, resolved against the project root.
func (m Manifest) Dir(root, key string) string {
	return filepath.Join(root, m.Layout[key])
}
