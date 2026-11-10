package manifest

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/cron"
)

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
	// Formula is the formula flai serve makes an instance of, in place of
	// starting the analyzer, each time Schedule comes round: a formula's
	// kebab-case name; empty for none (S-0424).
	Formula string `yaml:"formula,omitempty" json:"formula,omitempty"`
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
	if f := a.Formula; f != "" && !roleName.MatchString(f) {
		errs = append(errs, fmt.Sprintf("analysis.formula %q is not a formula's name such as analysis-cadence (lower case letters and digits, joined by single hyphens); write the name its file under design/workflows/formulas gives", f))
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

// Threads is the project's say about reminding the operator of a thread
// that awaits them (ADR-0205): how long it may await them before flai serve
// reminds them, and how often it reminds them again.
type Threads struct {
	// RemindAfter is how long a thread may await the operator before they
	// are reminded of it, a Go duration such as 4h; empty, or 0, means no
	// reminders.
	RemindAfter string `yaml:"remind_after,omitempty" json:"remind_after,omitempty"`
	// RemindEvery is how often the operator is reminded again of a thread
	// that still awaits them, a Go duration such as 12h; empty means
	// DefaultRemindEvery.
	RemindEvery string `yaml:"remind_every,omitempty" json:"remind_every,omitempty"`
}

// DefaultRemindEvery is how often the operator is reminded again of a
// thread that still awaits them when threads.remind_every is not set: a
// day.
const DefaultRemindEvery = 24 * time.Hour

// RemindAfterDuration is threads.remind_after as a duration: zero when it is
// empty or 0, and zero means no reminders.
func (t Threads) RemindAfterDuration() (time.Duration, error) {
	s := strings.TrimSpace(t.RemindAfter)
	if s == "" {
		return 0, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("threads.remind_after %q is not a duration; write one such as 4h or 90m, or unset it for no reminders", t.RemindAfter)
	}
	if d < 0 {
		return 0, fmt.Errorf("threads.remind_after %q is negative; write a duration such as 4h or 90m, or unset it for no reminders", t.RemindAfter)
	}
	return d, nil
}

// RemindEveryDuration is threads.remind_every as a duration:
// DefaultRemindEvery when it is empty.
func (t Threads) RemindEveryDuration() (time.Duration, error) {
	s := strings.TrimSpace(t.RemindEvery)
	if s == "" {
		return DefaultRemindEvery, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("threads.remind_every %q is not a duration; write one such as 24h or 12h, or unset it for 24h", t.RemindEvery)
	}
	if d <= 0 {
		return 0, fmt.Errorf("threads.remind_every %q is not more than zero; write a duration such as 24h or 12h, or unset it for 24h", t.RemindEvery)
	}
	return d, nil
}

// Errors are what is wrong with the threads settings, one sentence each;
// none when they are valid.
func (t Threads) Errors() []string {
	var errs []string
	if _, err := t.RemindAfterDuration(); err != nil {
		errs = append(errs, err.Error())
	}
	if _, err := t.RemindEveryDuration(); err != nil {
		errs = append(errs, err.Error())
	}
	return errs
}

// Orchestration is the project's say about the order of the ready column and
// when accepted work is released (S-0217), about what the orchestrator may do
// and the agent it runs as (S-0218), and about when flai serve starts and
// stops its runs (ADR-0209).
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
	// Coalesce is how long after the first queued reason flai serve starts
	// an orchestrator run, with every reason seen meanwhile, a Go duration
	// such as 30s; empty means DefaultCoalesce.
	Coalesce string `yaml:"coalesce,omitempty" json:"coalesce,omitempty"`
	// RunTimeout is how long an orchestrator run may hold before flai serve
	// stops it, a Go duration such as 20m; empty means DefaultRunTimeout.
	RunTimeout string `yaml:"run_timeout,omitempty" json:"run_timeout,omitempty"`
	// Schedule is when flai serve starts an orchestrator run for a periodic
	// look: a five-field cron expression in UTC or daily; empty means no
	// schedule.
	Schedule string `yaml:"schedule,omitempty" json:"schedule,omitempty"`
}

// DefaultCoalesce is how long after the first queued reason flai serve starts
// an orchestrator run when orchestration.coalesce is not set: 30 seconds.
const DefaultCoalesce = 30 * time.Second

// DefaultRunTimeout is how long an orchestrator run may hold before flai
// serve stops it when orchestration.run_timeout is not set: 20 minutes.
const DefaultRunTimeout = 20 * time.Minute

// CoalesceDuration is orchestration.coalesce as a duration: DefaultCoalesce
// when it is empty.
func (o Orchestration) CoalesceDuration() (time.Duration, error) {
	s := strings.TrimSpace(o.Coalesce)
	if s == "" {
		return DefaultCoalesce, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("orchestration.coalesce %q is not a duration; write one such as 30s or 2m, or unset it for 30s", o.Coalesce)
	}
	if d <= 0 {
		return 0, fmt.Errorf("orchestration.coalesce %q is not more than zero; write a duration such as 30s or 2m, or unset it for 30s", o.Coalesce)
	}
	return d, nil
}

// RunTimeoutDuration is orchestration.run_timeout as a duration:
// DefaultRunTimeout when it is empty.
func (o Orchestration) RunTimeoutDuration() (time.Duration, error) {
	s := strings.TrimSpace(o.RunTimeout)
	if s == "" {
		return DefaultRunTimeout, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("orchestration.run_timeout %q is not a duration; write one such as 20m or 1h, or unset it for 20m", o.RunTimeout)
	}
	if d <= 0 {
		return 0, fmt.Errorf("orchestration.run_timeout %q is not more than zero; write a duration such as 20m or 1h, or unset it for 20m", o.RunTimeout)
	}
	return d, nil
}

// OrchestrationSchedule is orchestration.schedule parsed: nil when it is
// empty.
func (o Orchestration) OrchestrationSchedule() (*cron.Schedule, error) {
	return parseSchedule("orchestration", o.Schedule)
}
