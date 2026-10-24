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
	"strconv"
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
	// Providers override, by name, the api, base_url, and models of the
	// host's providers an agent may name (ADR-0129); key_env is the host's
	// alone, and ProviderErrors refuses it here.
	Providers map[string]Provider `yaml:"providers,omitempty" json:"providers,omitempty"`
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
	// Digest is the curator's agent (ADR-0196), schedule, and compaction
	// window (ADR-0199).
	Digest Digest `yaml:"digest,omitempty" json:"digest,omitzero"`
	// Acceptance is how a story in review is accepted: on this machine, or
	// through a pull request on the project's forge (S-0392). Load refuses a
	// value Acceptance.Errors finds wrong.
	Acceptance Acceptance `yaml:"acceptance,omitempty" json:"acceptance,omitzero"`
	// MergeQueue is the merge queue's say about its batches (S-0399,
	// ADR-0173). Load refuses a value MergeQueue.Errors finds wrong.
	MergeQueue MergeQueue `yaml:"merge_queue,omitempty" json:"merge_queue,omitzero"`
	// Claims is how stories' touches claim paths: the shared paths whose
	// overlaps hold no story (S-0295, ADR-0096). A pattern that is not valid
	// does not stop the load; Claims.Errors names it, and it frees nothing.
	Claims Claims `yaml:"claims,omitempty" json:"claims,omitzero"`
	// Design is where the project's design system is kept (S-0381,
	// ADR-0146). Load refuses a path in it that is not inside the repository.
	Design Design `yaml:"design,omitempty" json:"design,omitzero"`
	// Flai is what the project asks of the flai that reads it (S-0181).
	Flai Requirement `yaml:"flai,omitempty" json:"flai,omitzero"`
}

// Requirement is the oldest flai that may read the project.
type Requirement struct {
	// Minimum is a flai release, X.Y.Z: one that knows every front-matter
	// field the project's items carry. A publish raises it to the newest
	// flai release whose front-matter fields changed, no newer than the flai
	// doing the publish, so it names only a release that can be installed
	// (release.RaiseMinimum, I-0107).
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

// Design is the project's say about where its design system is kept: one per
// project (ADR-0146).
type Design struct {
	// Tokens is the repository path of the project's design token file, a
	// DTCG file; empty means the default DesignTokens gives.
	Tokens string `yaml:"tokens,omitempty" json:"tokens,omitempty"`
	// Stylesheet is the repository path of the generated Tailwind v4
	// stylesheet flai design css --write writes the design tokens to, a file
	// under a sub-project (S-0385, ADR-0152, ADR-0159); empty means none, and
	// nothing is written or checked.
	Stylesheet string `yaml:"stylesheet,omitempty" json:"stylesheet,omitempty"`
}

// exampleDesignStylesheet is the stylesheet a refusal of design.stylesheet
// gives as one to write.
const exampleDesignStylesheet = "web/src/theme.css"

// DefaultDesignTokens is the design token file when design.tokens is not set
// and the design layout folder is design.
const DefaultDesignTokens = "design/" + designTokensUnder

// designTokensUnder is the default design token file under the design layout
// folder.
const designTokensUnder = "system/design-system.tokens.json"

// DesignTokens is the repository path of the design token file:
// design.tokens, or, when it is empty, design-system.tokens.json in the
// system folder of the design layout folder, DefaultDesignTokens under the
// default layout.
func (m Manifest) DesignTokens() string {
	if s := strings.TrimSpace(m.Design.Tokens); s != "" {
		return s
	}
	if d := m.Layout["design"]; d != "" {
		return path.Join(d, designTokensUnder)
	}
	return DefaultDesignTokens
}

// DesignStylesheet is the repository path of the generated design stylesheet,
// design.stylesheet, or "" when it is not set.
func (m Manifest) DesignStylesheet() string {
	return strings.TrimSpace(m.Design.Stylesheet)
}

// DesignStylesheetProject is the sub-project the design stylesheet lies under,
// the entry of projects whose path holds it; false when design.stylesheet is
// not set or lies outside every sub-project.
func (m Manifest) DesignStylesheetProject() (Project, bool) {
	s := m.DesignStylesheet()
	if s == "" {
		return Project{}, false
	}
	s = path.Clean(s)
	for _, p := range m.Projects {
		if d := path.Clean(strings.TrimSpace(p.Path)); strings.HasPrefix(s, d+"/") {
			return p, true
		}
	}
	return Project{}, false
}

// DesignGuidance is the repository path of the design system's guidance,
// design-system.md in the system folder of the design layout folder,
// design/system/design-system.md under the default layout (ADR-0146).
func (m Manifest) DesignGuidance() string {
	d := m.Layout["design"]
	if d == "" {
		d = "design"
	}
	return path.Join(d, "system/design-system.md")
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

// Digest is the project's say about the curator, the agent that takes in
// what the archive and the threads gained as digest entries: the agent it
// runs as (ADR-0196), when flai serve runs it, and how long an entry may go
// uncited before it is compacted (ADR-0199).
type Digest struct {
	// Agent is the curator's agent over the project's: what it sets wins,
	// and what it leaves out is the project's agent's.
	Agent *Agent `yaml:"agent,omitempty" json:"agent,omitempty"`
	// Schedule is when flai serve runs the curator: a five-field cron
	// expression in UTC or daily; empty means no schedule.
	Schedule string `yaml:"schedule,omitempty" json:"schedule,omitempty"`
	// CompactAfter is how long an active entry may go with no citation
	// before the curator compacts it, a Go duration such as 720h; empty
	// means DefaultCompactAfter, and 0 turns compaction off.
	CompactAfter string `yaml:"compact_after,omitempty" json:"compact_after,omitempty"`
}

// DefaultCompactAfter is how long an active digest entry may go with no
// citation before it is compacted when digest.compact_after is not set: 30
// days.
const DefaultCompactAfter = 720 * time.Hour

// DigestSchedule is digest.schedule parsed: nil when it is empty.
func (d Digest) DigestSchedule() (*cron.Schedule, error) {
	return parseSchedule("digest", d.Schedule)
}

// CompactAfterDuration is digest.compact_after as a duration:
// DefaultCompactAfter when it is empty, and zero when compaction is off.
func (d Digest) CompactAfterDuration() (time.Duration, error) {
	s := strings.TrimSpace(d.CompactAfter)
	if s == "" {
		return DefaultCompactAfter, nil
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("digest.compact_after %q is not a duration; write one such as 720h or 336h, or 0 to turn compaction off", d.CompactAfter)
	}
	if dur < 0 {
		return 0, fmt.Errorf("digest.compact_after %q is negative; write a duration such as 720h or 336h, or 0 to turn compaction off", d.CompactAfter)
	}
	return dur, nil
}

// Errors are what is wrong with the digest settings, one sentence each;
// none when they are valid.
func (d Digest) Errors() []string {
	var errs []string
	if _, err := d.DigestSchedule(); err != nil {
		errs = append(errs, err.Error())
	}
	if _, err := d.CompactAfterDuration(); err != nil {
		errs = append(errs, err.Error())
	}
	return append(errs, d.Agent.problems("digest.agent")...)
}

// DigestAgent is the agent the curator is started with: digest.agent merged
// over the project's agent, as AnalysisAgent is. Nil when neither sets
// anything.
func (m Manifest) DigestAgent() *Agent {
	return m.Agent.With(m.Digest.Agent)
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
	// PlanBacklogStories lets it ask the planner to plan a story in the
	// backlog that lacks a plan, settle what that planner asks on its
	// threads, and give a backlog story without any its cost of delay
	// inputs (S-0328, ADR-0119).
	PlanBacklogStories bool `yaml:"plan_backlog_stories,omitempty" json:"plan_backlog_stories,omitempty"`
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
	// EditSubprojects lets it add, change, and remove the manifest's
	// sub-projects, as flai subproject does (S-0491).
	EditSubprojects bool `yaml:"edit_subprojects,omitempty" json:"edit_subprojects,omitempty"`

	// unknown are the keys under orchestration.permissions that name no
	// permission, for Errors.
	unknown []string
}

// The permissions, as orchestration.permissions keys them.
const (
	PermitPlanBacklogEpics   = "plan_backlog_epics"
	PermitPlanBacklogStories = "plan_backlog_stories"
	PermitFinalizeDrafts     = "finalize_drafts"
	PermitPromoteToReady     = "promote_to_ready"
	PermitOrderReady         = "order_ready"
	PermitAnswerThreads      = "answer_threads"
	PermitAcceptReviews      = "accept_reviews"
	PermitPublish            = "publish"
	PermitEditSubprojects    = "edit_subprojects"
)

// PermissionNames are the keys of orchestration.permissions, in the order
// they are listed to the operator.
var PermissionNames = []string{PermitPlanBacklogEpics, PermitPlanBacklogStories, PermitFinalizeDrafts, PermitPromoteToReady, PermitOrderReady, PermitAnswerThreads, PermitAcceptReviews, PermitPublish, PermitEditSubprojects}

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
	case PermitPlanBacklogStories:
		return p.PlanBacklogStories
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
	case PermitEditSubprojects:
		return p.EditSubprojects
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

// Acceptance is the project's say about how a story in review is accepted
// (S-0392). Each key empty means its default, so a manifest without the
// block accepts as flai always has.
type Acceptance struct {
	// Mode is one of AcceptanceModes; empty means AcceptLocal.
	Mode string `yaml:"mode,omitempty" json:"mode,omitempty"`
	// Forge is the forge, one of Forges, that pull-request acceptance goes
	// through; empty means ForgeGitHub. Under AcceptLocal it is unused.
	Forge string `yaml:"forge,omitempty" json:"forge,omitempty"`
	// MergeMethod is how the forge merges an accepted pull request, one of
	// MergeMethods; empty means MergeMerge.
	MergeMethod string `yaml:"merge_method,omitempty" json:"merge_method,omitempty"`
}

// The values of acceptance.mode.
const (
	// AcceptLocal merges a story's branch on this machine, as flai accept
	// does. It is the default.
	AcceptLocal = "local"
	// AcceptPullRequest accepts a story through a pull request on the
	// project's forge.
	AcceptPullRequest = "pull-request"
)

// AcceptanceModes are the values of acceptance.mode, in the order they are
// listed to the operator.
var AcceptanceModes = []string{AcceptLocal, AcceptPullRequest}

// ForgeGitHub is GitHub, the default forge.
const ForgeGitHub = "github"

// Forges are the values of acceptance.forge.
var Forges = []string{ForgeGitHub}

// The values of acceptance.merge_method, as the forge merges a pull request.
const (
	// MergeMerge merges with a merge commit. It is the default.
	MergeMerge = "merge"
	// MergeSquash squashes the branch into one commit.
	MergeSquash = "squash"
	// MergeRebase rebases the branch's commits onto the base.
	MergeRebase = "rebase"
)

// MergeMethods are the values of acceptance.merge_method, in the order they
// are listed to the operator.
var MergeMethods = []string{MergeMerge, MergeSquash, MergeRebase}

// ModeOrDefault is acceptance.mode, or AcceptLocal when it is empty. It does
// not say whether the mode is one of AcceptanceModes; Errors does.
func (a Acceptance) ModeOrDefault() string {
	if s := strings.TrimSpace(a.Mode); s != "" {
		return s
	}
	return AcceptLocal
}

// PullRequest reports whether stories are accepted through pull requests.
func (a Acceptance) PullRequest() bool {
	return a.ModeOrDefault() == AcceptPullRequest
}

// ForgeOrDefault is acceptance.forge, or ForgeGitHub when it is empty, under
// either mode. It does not say whether the forge is one of Forges; Errors
// does.
func (a Acceptance) ForgeOrDefault() string {
	if s := strings.TrimSpace(a.Forge); s != "" {
		return s
	}
	return ForgeGitHub
}

// MergeMethodOrDefault is acceptance.merge_method, or MergeMerge when it is
// empty. It does not say whether the method is one of MergeMethods; Errors
// does.
func (a Acceptance) MergeMethodOrDefault() string {
	if s := strings.TrimSpace(a.MergeMethod); s != "" {
		return s
	}
	return MergeMerge
}

// Errors are what is wrong with the acceptance settings, one sentence each;
// none when they are valid. A forge or a merge method set under AcceptLocal
// is valid, and unused.
func (a Acceptance) Errors() []string {
	var errs []string
	if m := a.ModeOrDefault(); !slices.Contains(AcceptanceModes, m) {
		errs = append(errs, fmt.Sprintf("acceptance.mode %q is not an acceptance mode; write local (merge on this machine) or pull-request (merge through a pull request on the forge), or remove it for local", a.Mode))
	}
	if f := a.ForgeOrDefault(); !slices.Contains(Forges, f) {
		errs = append(errs, fmt.Sprintf("acceptance.forge %q is not a forge flai reaches; write %s, or remove it for github", a.Forge, strings.Join(Forges, ", ")))
	}
	if m := a.MergeMethodOrDefault(); !slices.Contains(MergeMethods, m) {
		errs = append(errs, fmt.Sprintf("acceptance.merge_method %q is not a merge method; write merge (a merge commit), squash (one commit), or rebase (the commits onto the base), or remove it for merge", a.MergeMethod))
	}
	return errs
}

// MergeQueue is the project's say about the merge queue flai serve runs
// behind the merge-queue host action (S-0399, ADR-0173). A manifest without
// the block batches DefaultBatchMax stories at most.
type MergeQueue struct {
	// BatchMax is the most enqueued stories one batch holds, a whole number
	// of one or more; unset means DefaultBatchMax.
	BatchMax *int `yaml:"batch_max,omitempty" json:"batch_max,omitempty"`
}

// DefaultBatchMax is the most stories a batch of the merge queue holds when
// merge_queue.batch_max is not set (ADR-0173).
const DefaultBatchMax = 4

// BatchMax is merge_queue.batch_max, or DefaultBatchMax when it is not set.
// It does not say whether the value is one or more; MergeQueue.Errors does.
func (m Manifest) BatchMax() int {
	if b := m.MergeQueue.BatchMax; b != nil {
		return *b
	}
	return DefaultBatchMax
}

// batchMaxProblem says why merge_queue.batch_max, as data writes it, is not
// a whole number, which a load into MergeQueue would otherwise fail on with
// the decoder's words or, for a fraction such as 2.5, cut to a whole one; ""
// when it is one, is not set, or data does not parse, which the load says.
func batchMaxProblem(data []byte) string {
	var raw struct {
		MergeQueue struct {
			BatchMax any `yaml:"batch_max"`
		} `yaml:"merge_queue"`
	}
	if yaml.Unmarshal(data, &raw) != nil {
		return ""
	}
	written := ""
	switch v := raw.MergeQueue.BatchMax.(type) {
	case nil, int, int64, uint64:
		return ""
	case float64:
		if v == math.Trunc(v) && !math.IsInf(v, 0) {
			return ""
		}
		written = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		written = strconv.Quote(v)
	default:
		written = fmt.Sprint(v)
	}
	return fmt.Sprintf("merge_queue.batch_max %s is not a whole number; write the most enqueued stories one batch of the merge queue holds, such as %d, or remove it for %d", written, DefaultBatchMax, DefaultBatchMax)
}

// Errors are what is wrong with the merge queue's settings, one sentence
// each; none when they are valid.
func (q MergeQueue) Errors() []string {
	if b := q.BatchMax; b != nil && *b < 1 {
		return []string{fmt.Sprintf("merge_queue.batch_max %d is not a number of one or more; write the most enqueued stories one batch of the merge queue holds, such as %d, or remove it for %d", *b, DefaultBatchMax, DefaultBatchMax)}
	}
	return nil
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
	// An argument that is PlaceholderPackages, PlaceholderFiles, or
	// PlaceholderImporters alone stands for what the tier's paths selected;
	// nothing else is interpreted.
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
	// FixCommand, when set, is an argument list run before flai commits the
	// files the tier's paths select, with placeholders as Command takes
	// them. It rewrites those files in place; its exit status is ignored.
	FixCommand []string `yaml:"fix_command,omitempty" json:"fix_command,omitempty"`
	// Covers are the names of the tiers this one runs everything of, so that
	// when it is selected with them only this one runs (S-0342).
	Covers []string `yaml:"covers,omitempty" json:"covers,omitempty"`
	// Queue leaves the tier to the merge queue: while the queue runs, flai
	// verify leaves it out and the queue's batch runs it (S-0405).
	Queue bool `yaml:"queue,omitempty" json:"queue,omitempty"`
}

// The placeholders a tier's command may hold, each as an argument of its own.
const (
	// PlaceholderPackages stands for the packages of the paths the tier
	// selected.
	PlaceholderPackages = "{packages}"
	// PlaceholderFiles stands for the files the tier selected.
	PlaceholderFiles = "{files}"
	// PlaceholderImporters stands for the packages PlaceholderPackages names
	// and every package of the same module that imports one of them,
	// directly or through other packages, in its own files or its test
	// files.
	PlaceholderImporters = "{importers}"
)

// tierPlaceholders are the placeholders a tier's command may hold, in the
// order they are listed to the operator.
var tierPlaceholders = []string{PlaceholderPackages, PlaceholderFiles, PlaceholderImporters}

// tierPlaceholderList names tierPlaceholders in a sentence.
const tierPlaceholderList = PlaceholderPackages + ", " + PlaceholderFiles + ", or " + PlaceholderImporters

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
	tiers := map[string]int{}
	for i, t := range m.Tests {
		if name := strings.TrimSpace(t.Name); name != "" {
			if _, ok := tiers[name]; !ok {
				tiers[name] = i
			}
		}
	}
	seen := map[string]int{}
	for i, t := range m.Tests {
		at := func(field string) string { return testField(i, t, field) }
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
		if len(t.FixCommand) > 0 {
			errs = append(errs, commandProblems(at("fix_command"), t.FixCommand, true)...)
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
		named := map[string]bool{}
		for _, c := range t.Covers {
			covered := strings.TrimSpace(c)
			switch _, ok := tiers[covered]; {
			case named[covered]:
				errs = append(errs, fmt.Sprintf("%s names %q twice; name each tier it covers once", at("covers"), c))
			case !ok:
				errs = append(errs, fmt.Sprintf("%s %q names no tier; write the name of another tier, or remove it", at("covers"), c))
			case covered == name:
				errs = append(errs, fmt.Sprintf("%s names the tier itself; name only other tiers, or remove %q", at("covers"), c))
			}
			named[covered] = true
		}
	}
	for _, cycle := range coverCycles(m.Tests, tiers) {
		names := make([]string, len(cycle))
		for k, i := range cycle {
			names[k] = strconv.Quote(strings.TrimSpace(m.Tests[i].Name))
		}
		first := cycle[0]
		errs = append(errs, fmt.Sprintf("%s makes a cycle through the tiers %s, which cover one another; remove a name from one of their covers, so that no tier covers one that covers it", testField(first, m.Tests[first], "covers"), strings.Join(names, ", ")))
	}
	return errs
}

// testField is a tier's field as TestErrors names it, such as tests[1].format
// (tier "unit"), with the tier's name when it has one.
func testField(i int, t TestTier, field string) string {
	s := fmt.Sprintf("tests[%d].%s", i, field)
	if name := strings.TrimSpace(t.Name); name != "" && field != "name" {
		s += fmt.Sprintf(" (tier %q)", name)
	}
	return s
}

// coverCycles are the groups of tiers that cover one another through covers,
// each the tiers' indexes in list order and each group once; tiers maps a
// name to the index of the first tier with it. A tier naming itself is not a
// cycle here, as TestErrors refuses it on its own.
func coverCycles(list []TestTier, tiers map[string]int) [][]int {
	covers := make([][]int, len(list))
	for i, t := range list {
		for _, c := range t.Covers {
			if j, ok := tiers[strings.TrimSpace(c)]; ok && j != i {
				covers[i] = append(covers[i], j)
			}
		}
	}
	reaches := make([][]bool, len(list))
	for i := range list {
		reaches[i] = make([]bool, len(list))
		next := slices.Clone(covers[i])
		for len(next) > 0 {
			j := next[len(next)-1]
			next = next[:len(next)-1]
			if !reaches[i][j] {
				reaches[i][j] = true
				next = append(next, covers[j]...)
			}
		}
	}
	var cycles [][]int
	reported := make([]bool, len(list))
	for i := range list {
		if reported[i] || !reaches[i][i] {
			continue
		}
		var cycle []int
		for j := range list {
			if reaches[i][j] && reaches[j][i] {
				cycle = append(cycle, j)
				reported[j] = true
			}
		}
		cycles = append(cycles, cycle)
	}
	return cycles
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
		case slices.Contains(tierPlaceholders, a) && !placeholders:
			errs = append(errs, fmt.Sprintf("%s argument %s is a placeholder, and all_command runs for no paths; write it without %s", field, a, tierPlaceholderList))
		case slices.Contains(tierPlaceholders, a):
		case placeholder.MatchString(a):
			errs = append(errs, fmt.Sprintf("%s argument %s is not a placeholder flai fills; write %s, or remove it", field, a, tierPlaceholderList))
		case slices.ContainsFunc(tierPlaceholders, func(p string) bool { return strings.Contains(a, p) }):
			errs = append(errs, fmt.Sprintf("%s argument %q holds a placeholder inside it, which flai does not fill; write %s as an argument of its own", field, a, tierPlaceholderList))
		}
	}
	return errs
}

// dirProblem says what is wrong with a tier's dir, or "" when it is valid: a
// clean path inside the repository, relative to its root, or empty.
func dirProblem(d string) string {
	return repoPathProblem(d, "a folder", "flai", "flai/web")
}

// repoPathProblem says what is wrong with p, a path the manifest names, or ""
// when it is valid: a clean path inside the repository, relative to its root,
// or empty. what is what p names, such as a folder, and example and clean
// are paths the reason gives as ones to write.
func repoPathProblem(p, what, example, clean string) string {
	return pathProblem(p, "the repository", "the repository root", what, example, clean)
}

// pathProblem says what is wrong with p, a path the manifest names inside
// within, such as the repository, whose root is root, such as the repository
// root, or "" when it is valid: a clean path inside it, relative to its root,
// or empty. what, example, and clean are as repoPathProblem's.
func pathProblem(p, within, root, what, example, clean string) string {
	switch {
	case p == "":
		return ""
	case isAbsolute(p):
		return "is an absolute path; write it relative to " + root + ", such as " + example
	case path.Clean(p) == ".." || strings.HasPrefix(path.Clean(p), "../"):
		return "leaves " + within + "; write " + what + " inside it, relative to its root, such as " + example
	case path.Clean(p) != p:
		return "is not a clean path; write it with single slashes, no . or .. segments, and none at the end, such as " + clean
	}
	return ""
}

// isAbsolute reports whether p is an absolute path, on any system: one
// beginning with a slash or a backslash, or with a drive letter.
func isAbsolute(p string) bool {
	return strings.HasPrefix(p, "/") || strings.HasPrefix(p, `\`) || len(p) > 1 && p[1] == ':'
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
	// Components is where a Svelte sub-project's components are found and
	// what the harness imports before them (S-0437); nil means the defaults
	// ComponentSettings gives. Only a kind SupportsComponents takes it.
	Components *Components `yaml:"components,omitempty" json:"components,omitempty"`
	// Install is how the sub-project installs its dependencies, for flai
	// accept to run in the main checkout after an acceptance that changed
	// them (S-0483, ADR-0185); nil means flai installs nothing for it.
	Install *Install `yaml:"install,omitempty" json:"install,omitempty"`
}

// Install is a sub-project's install command and the paths whose change
// needs it (S-0483, ADR-0185). Its globs have the syntax of components'
// globs, relative to the sub-project's folder.
type Install struct {
	// Command is the program and its arguments, run with no shell in the
	// sub-project's folder.
	Command []string `yaml:"command" json:"command"`
	// Paths are the globs of the files whose change needs Command run.
	Paths []string `yaml:"paths" json:"paths"`
}

// exampleInstallPath is the glob a refusal of install.paths gives as one to
// write.
const exampleInstallPath = "package.json"

// InstallErrors are what is wrong with the sub-projects' install keys, one
// sentence each, each beginning with the sub-project's index and field, such
// as projects[1].install.command: a command that is empty or has no program
// first, and paths that are empty or hold a glob that is not valid. None when
// every key is valid.
func (m Manifest) InstallErrors() []string {
	var errs []string
	for i, p := range m.Projects {
		in := p.Install
		if in == nil {
			continue
		}
		at := func(field string) string { return projectField(i, p, "install", field) }
		switch {
		case len(in.Command) == 0:
			errs = append(errs, at("command")+" is empty; write the program and its arguments as a list, run in the sub-project's folder, such as [pnpm, install]")
		case strings.TrimSpace(in.Command[0]) == "":
			errs = append(errs, at("command")+" has no program first; write the program, then its arguments, such as [pnpm, install]")
		}
		if len(in.Paths) == 0 {
			errs = append(errs, at("paths")+" is empty; write the globs of the files whose change needs the install, relative to the sub-project's folder, such as "+exampleInstallPath)
		}
		for _, g := range in.Paths {
			if r := projectGlobProblem(g, exampleInstallPath); r != "" {
				errs = append(errs, fmt.Sprintf("%s pattern %q %s", at("paths"), g, r))
			}
		}
	}
	return errs
}

// InstallMatches reports whether the sub-project's install command needs to
// run for changed, paths relative to the repository root: whether one of
// them lies under the sub-project's folder and a valid glob of install.paths
// matches it there. False when the sub-project has no install key.
func (p Project) InstallMatches(changed []string) bool {
	if p.Install == nil {
		return false
	}
	dir := path.Clean(strings.TrimSpace(p.Path)) + "/"
	for _, c := range changed {
		rel, ok := strings.CutPrefix(cleanEntry(c), dir)
		if !ok || rel == "" {
			continue
		}
		segs := strings.Split(rel, "/")
		for _, g := range p.Install.Paths {
			if projectGlobProblem(g, exampleInstallPath) == "" && patternMatches(g, segs, false) {
				return true
			}
		}
	}
	return false
}

// The kinds of sub-project whose components flai reads (S-0437).
const (
	// ProjectKindSvelteKit is a SvelteKit app.
	ProjectKindSvelteKit = "sveltekit"
	// ProjectKindSvelte is a Svelte project that is not a SvelteKit app, such
	// as a component package built with Vite.
	ProjectKindSvelte = "svelte"
)

// ComponentKinds are the kinds of sub-project whose components flai reads,
// in the order they are listed to the operator.
var ComponentKinds = []string{ProjectKindSvelteKit, ProjectKindSvelte}

// DefaultComponentInclude is the glob, relative to the sub-project's folder,
// that finds its components when components.include is not set: every
// .svelte file under src/lib, so that routes are left out.
const DefaultComponentInclude = "src/lib/**/*.svelte"

// exampleComponentSetup is the setup a refusal of components.setup gives as
// one to write.
const exampleComponentSetup = "src/app.css"

// Components is a sub-project's say about its components (S-0437). Its globs
// have the syntax of claims.shared's, relative to the sub-project's folder.
type Components struct {
	// Include are the globs of the sub-project's component files; empty means
	// DefaultComponentInclude.
	Include []string `yaml:"include,omitempty" json:"include,omitempty"`
	// Exclude are the globs of files Include finds that are not components to
	// list; empty means none.
	Exclude []string `yaml:"exclude,omitempty" json:"exclude,omitempty"`
	// Setup is a module of the sub-project, relative to its folder, that the
	// harness imports before any component, such as the app's stylesheet;
	// empty means none.
	Setup string `yaml:"setup,omitempty" json:"setup,omitempty"`
}

// SupportsComponents reports whether flai reads the components of a
// sub-project of kind: one of ComponentKinds.
func SupportsComponents(kind string) bool {
	return slices.Contains(ComponentKinds, kind)
}

// ComponentSettings are the sub-project's components settings in effect:
// its components key over the defaults, include DefaultComponentInclude,
// nothing excluded, and no setup. False, with none, when its kind is not one
// SupportsComponents takes. It does not say whether the settings are valid;
// ComponentErrors does.
func (p Project) ComponentSettings() (Components, bool) {
	if !SupportsComponents(p.Kind) {
		return Components{}, false
	}
	c := Components{Include: []string{DefaultComponentInclude}}
	if p.Components == nil {
		return c, true
	}
	if len(p.Components.Include) > 0 {
		c.Include = slices.Clone(p.Components.Include)
	}
	c.Exclude = slices.Clone(p.Components.Exclude)
	c.Setup = strings.TrimSpace(p.Components.Setup)
	return c, true
}

// ComponentErrors are what is wrong with the sub-projects' components keys,
// one sentence each, each beginning with the sub-project's index and field,
// such as projects[1].components.include: a key on a sub-project of a kind
// SupportsComponents does not take, a glob that is not valid, and a setup
// that is not a clean path inside the sub-project or, when fsys, the
// repository root, is not nil, names no file under the sub-project's folder.
// None when every key is valid.
func (m Manifest) ComponentErrors(fsys fs.FS) []string {
	var errs []string
	for i, p := range m.Projects {
		c := p.Components
		if c == nil {
			continue
		}
		at := func(field string) string { return projectField(i, p, "components", field) }
		if !SupportsComponents(p.Kind) {
			errs = append(errs, fmt.Sprintf("%s is set on a sub-project of kind %q, whose components flai does not read; remove it, or set kind to %s if the sub-project is one", at(""), p.Kind, strings.Join(ComponentKinds, " or ")))
			continue
		}
		for _, list := range []struct {
			field string
			globs []string
		}{{"include", c.Include}, {"exclude", c.Exclude}} {
			for _, g := range list.globs {
				if r := projectGlobProblem(g, DefaultComponentInclude); r != "" {
					errs = append(errs, fmt.Sprintf("%s pattern %q %s", at(list.field), g, r))
				}
			}
		}
		setup := strings.TrimSpace(c.Setup)
		if r := pathProblem(setup, "the sub-project", "the sub-project's folder", "a module", exampleComponentSetup, exampleComponentSetup); r != "" {
			errs = append(errs, fmt.Sprintf("%s %q %s", at("setup"), c.Setup, r))
		} else if setup != "" && fsys != nil && !isFile(fsys, path.Join(path.Clean(strings.TrimSpace(p.Path)), setup)) {
			errs = append(errs, fmt.Sprintf("%s %q names no file under %s; write the path of a module of the sub-project, relative to its folder, such as %s, or remove it", at("setup"), c.Setup, path.Clean(strings.TrimSpace(p.Path)), exampleComponentSetup))
		}
	}
	return errs
}

// projectField is a field of a sub-project's key as ComponentErrors and
// InstallErrors name it, such as projects[1].components.include (sub-project
// "web"), or the key itself when field is empty.
func projectField(i int, p Project, key, field string) string {
	s := fmt.Sprintf("projects[%d].%s", i, key)
	if field != "" {
		s += "." + field
	}
	if name := strings.TrimSpace(p.Name); name != "" {
		s += fmt.Sprintf(" (sub-project %q)", name)
	}
	return s
}

// projectGlobProblem says what is wrong with a glob of a sub-project's key,
// such as components.include or install.paths, or "" when it is valid: a path
// relative to the sub-project's folder and inside it, with the syntax of
// claims.shared's. example is a glob the reason gives as one to write.
func projectGlobProblem(g, example string) string {
	switch {
	case strings.TrimSpace(g) == "":
		return "is empty; write a glob relative to the sub-project's folder, such as " + example + ", or remove it"
	case isAbsolute(g):
		return "is an absolute path; write it relative to the sub-project's folder, such as " + example
	}
	for _, seg := range strings.Split(g, "/") {
		switch seg {
		case "..":
			return "has a .. segment, which leaves the sub-project; write a glob inside its folder, relative to it, such as " + example
		case ".":
			return "has a . segment; write the glob without it, such as " + example + " for ./" + example
		}
	}
	return patternProblem(g)
}

// Dashboard is how the flaiover image runs for this project.
type Dashboard struct {
	Image string `yaml:"image"`
	Tag   string `yaml:"tag"`
	Port  int    `yaml:"port"`
	Bind  string `yaml:"bind"` // host address to publish on; default 0.0.0.0
	// Autocommit commits documents saved from the dashboard; nil means true (ADR-0023).
	Autocommit *bool `yaml:"autocommit,omitempty"`
	// AllowUnsigned, when set, is this project's dashboard.allow_unsigned over
	// the host configuration's (ADR-0070); nil leaves the host's.
	AllowUnsigned *bool `yaml:"allow_unsigned,omitempty"`
	// NotifyURL, when set, is where the dashboard's server posts new inbox
	// entries (S-0042). flai itself never calls it.
	NotifyURL string `yaml:"notify_url,omitempty"`
}

// Autocommit reports whether documents saved from the dashboard are committed.
func (m Manifest) Autocommit() bool {
	return m.Dashboard.Autocommit == nil || *m.Dashboard.Autocommit
}

// AllowUnsignedOr reports whether unsigned dashboards and flais are allowed
// for this project, given the host configuration's dashboard.allow_unsigned.
func (d Dashboard) AllowUnsignedOr(host bool) bool {
	if d.AllowUnsigned != nil {
		return *d.AllowUnsigned
	}
	return host
}

// Load parses the manifest at path.
func Load(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}
	if p := batchMaxProblem(data); p != "" {
		return Manifest{}, fmt.Errorf("%s: %s", path, p)
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
	if r := repoPathProblem(strings.TrimSpace(m.Design.Tokens), "a file", DefaultDesignTokens, DefaultDesignTokens); r != "" {
		return Manifest{}, fmt.Errorf("%s: design.tokens %q %s", path, m.Design.Tokens, r)
	}
	if r := repoPathProblem(m.DesignStylesheet(), "a file", exampleDesignStylesheet, exampleDesignStylesheet); r != "" {
		return Manifest{}, fmt.Errorf("%s: design.stylesheet %q %s", path, m.Design.Stylesheet, r)
	}
	if errs := m.Acceptance.Errors(); len(errs) > 0 {
		return Manifest{}, fmt.Errorf("%s: %s", path, strings.Join(errs, "; "))
	}
	if errs := m.MergeQueue.Errors(); len(errs) > 0 {
		return Manifest{}, fmt.Errorf("%s: %s", path, strings.Join(errs, "; "))
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
