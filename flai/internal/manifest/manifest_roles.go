package manifest

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/cron"
)

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
	// SetCostOfDelay lets it add or change the cost of delay inputs of an
	// epic or a story neither done, cancelled, nor archived, whether or not
	// it or its epic has inputs already, and nothing else (S-0523,
	// ADR-0224).
	SetCostOfDelay bool `yaml:"set_cost_of_delay,omitempty" json:"set_cost_of_delay,omitempty"`
	// FinalizeDrafts lets it finalize a draft story.
	FinalizeDrafts bool `yaml:"finalize_drafts,omitempty" json:"finalize_drafts,omitempty"`
	// PromoteToReady lets it move a story to ready.
	PromoteToReady bool `yaml:"promote_to_ready,omitempty" json:"promote_to_ready,omitempty"`
	// OrderReady lets it write the order of the ready column.
	OrderReady bool `yaml:"order_ready,omitempty" json:"order_ready,omitempty"`
	// AnswerThreads is how it answers threads: one of AnswerModes; empty
	// means AnswerOff.
	AnswerThreads string `yaml:"answer_threads,omitempty" json:"answer_threads,omitempty"`
	// ConfirmAfter is how long a recommendation of the orchestrator's may
	// await the operator before flai serve confirms it (ADR-0206), a Go
	// duration such as 48h; empty, or 0, means it never does. It is the one
	// permission that is a duration.
	ConfirmAfter string `yaml:"confirm_after,omitempty" json:"confirm_after,omitempty"`
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
	PermitSetCostOfDelay     = "set_cost_of_delay"
	PermitFinalizeDrafts     = "finalize_drafts"
	PermitPromoteToReady     = "promote_to_ready"
	PermitOrderReady         = "order_ready"
	PermitAnswerThreads      = "answer_threads"
	PermitConfirmAfter       = "confirm_after"
	PermitAcceptReviews      = "accept_reviews"
	PermitPublish            = "publish"
	PermitEditSubprojects    = "edit_subprojects"
)

// PermissionNames are the keys of orchestration.permissions, in the order
// they are listed to the operator.
var PermissionNames = []string{PermitPlanBacklogEpics, PermitPlanBacklogStories, PermitSetCostOfDelay, PermitFinalizeDrafts, PermitPromoteToReady, PermitOrderReady, PermitAnswerThreads, PermitConfirmAfter, PermitAcceptReviews, PermitPublish, PermitEditSubprojects}

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

// ConfirmAfterDuration is orchestration.permissions.confirm_after as a
// duration: zero when it is empty or 0, and zero means flai serve confirms no
// recommendation.
func (p Permissions) ConfirmAfterDuration() (time.Duration, error) {
	s := strings.TrimSpace(p.ConfirmAfter)
	if s == "" {
		return 0, nil
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return 0, fmt.Errorf("orchestration.permissions.confirm_after %q is not a duration; write one such as 48h or 90m, or unset it for none confirmed", p.ConfirmAfter)
	}
	if d < 0 {
		return 0, fmt.Errorf("orchestration.permissions.confirm_after %q is negative; write a duration such as 48h or 90m, or unset it for none confirmed", p.ConfirmAfter)
	}
	return d, nil
}

// Allows reports whether the permission name, one of PermissionNames, is on:
// a boolean set true, answer_threads set to recommend or autonomous
// (AnswerMode tells them apart), or confirm_after set to a duration of more
// than zero (ConfirmAfterDuration says how long). A name that is no
// permission is off.
func (p Permissions) Allows(name string) bool {
	switch name {
	case PermitPlanBacklogEpics:
		return p.PlanBacklogEpics
	case PermitPlanBacklogStories:
		return p.PlanBacklogStories
	case PermitSetCostOfDelay:
		return p.SetCostOfDelay
	case PermitFinalizeDrafts:
		return p.FinalizeDrafts
	case PermitPromoteToReady:
		return p.PromoteToReady
	case PermitOrderReady:
		return p.OrderReady
	case PermitAnswerThreads:
		m := p.AnswerMode()
		return m == AnswerRecommend || m == AnswerAutonomous
	case PermitConfirmAfter:
		d, err := p.ConfirmAfterDuration()
		return err == nil && d > 0
	case PermitAcceptReviews:
		return p.AcceptReviews
	case PermitPublish:
		return p.Publish
	case PermitEditSubprojects:
		return p.EditSubprojects
	}
	return false
}

// errors are what is wrong with the permissions under the manifest key key,
// such as orchestration.permissions, one sentence each.
func (p Permissions) errors(key string) []string {
	var errs []string
	for _, k := range p.unknown {
		errs = append(errs, fmt.Sprintf("%s has no permission %q; write one of %s, or remove it", key, k, strings.Join(PermissionNames, ", ")))
	}
	if m := p.AnswerMode(); !slices.Contains(AnswerModes, m) {
		errs = append(errs, fmt.Sprintf("%s.answer_threads %q is not a way of answering threads; write off (leave them to the operator), recommend (reply with a recommendation), or autonomous (answer them), or remove it for off", key, p.AnswerThreads))
	}
	return errs
}

// RoleSettings are the operator's settings for one role a role file defines,
// roles.<name> in the manifest (ADR-0207 §5).
type RoleSettings struct {
	// Permissions are what the role may do without the operator, each off
	// when unset, as its role file's permissions entries name them under
	// roles.<name>.permissions.<p>.
	Permissions Permissions `yaml:"permissions,omitempty" json:"permissions,omitzero"`
	// Agent is the role's agent over the project's, roles.<name>.agent: what
	// it sets wins, and what it leaves out is the project's agent's (S-0417).
	Agent *Agent `yaml:"agent,omitempty" json:"agent,omitempty"`
	// Schedule is when flai serve starts the role, roles.<name>.schedule,
	// for a role file whose trigger names it: a five-field cron expression
	// in UTC or daily; empty means no schedule (S-0418).
	Schedule string `yaml:"schedule,omitempty" json:"schedule,omitempty"`
}

// builtinAgent is where a built-in role's agent and schedule are set
// instead of roles.<name>.agent and roles.<name>.schedule: the manifest
// keys, and the agent the first gives.
type builtinAgent struct {
	key      string
	agent    func(Manifest) *Agent
	schedule string
}

// builtinAgents are the built-in roles' agents and schedules, by role name.
var builtinAgents = map[string]builtinAgent{
	"planner":      {"planning.agent", Manifest.PlanningAgent, "planning.schedule"},
	"orchestrator": {"orchestration.agent", Manifest.OrchestrationAgent, "orchestration.schedule"},
	"analyzer":     {"analysis.agent", Manifest.AnalysisAgent, "analysis.schedule"},
	"curator":      {"digest.agent", Manifest.DigestAgent, "digest.schedule"},
}

// roleName is a role's name as its role file writes it: kebab-case.
var roleName = regexp.MustCompile(`^[a-z][a-z0-9]*(?:-[a-z0-9]+)*$`)

// RoleErrors are what is wrong with the roles' settings, one sentence each,
// by role name; none when they are valid.
func (m Manifest) RoleErrors() []string {
	var errs []string
	for _, name := range slices.Sorted(maps.Keys(m.Roles)) {
		key := "roles." + name
		if !roleName.MatchString(name) {
			errs = append(errs, fmt.Sprintf("%s is not a role's name such as reviewer (lower case letters and digits, joined by single hyphens); write the name its role file gives", key))
			continue
		}
		r := m.Roles[name]
		errs = append(errs, r.Permissions.errors(key+".permissions")...)
		b, builtin := builtinAgents[name]
		if builtin && r.Agent != nil {
			errs = append(errs, fmt.Sprintf("%s.agent is not read: the %s's agent is %s; set that instead, and remove this", key, name, b.key))
		} else {
			errs = append(errs, r.Agent.problems(key+".agent")...)
		}
		if builtin && strings.TrimSpace(r.Schedule) != "" {
			errs = append(errs, fmt.Sprintf("%s.schedule is not read: the %s's schedule is %s; set that instead, and remove this", key, name, b.schedule))
		} else if _, err := m.RoleSchedule(name); err != nil {
			errs = append(errs, err.Error())
		}
	}
	return errs
}

// RoleSchedule is roles.<name>.schedule parsed: nil when it is empty. It is
// read for a role a role file defines, not for a built-in role, whose
// schedule is its own block's (RoleErrors).
func (m Manifest) RoleSchedule(name string) (*cron.Schedule, error) {
	return parseSchedule("roles."+name, m.Roles[name].Schedule)
}

// ScheduleSetting is the schedule the manifest setting named holds, as a
// role file's {schedule: <setting>} trigger names it: the value as written
// and parsed, nil when it is empty. The setting is roles.<name>.schedule or
// one of planning.schedule, analysis.schedule, digest.schedule, and
// orchestration.schedule; any other is an error.
func (m Manifest) ScheduleSetting(setting string) (string, *cron.Schedule, error) {
	switch setting {
	case "planning.schedule":
		s, err := m.Planning.PlanSchedule()
		return m.Planning.Schedule, s, err
	case "analysis.schedule":
		s, err := m.Analysis.AnalysisSchedule()
		return m.Analysis.Schedule, s, err
	case "digest.schedule":
		s, err := m.Digest.DigestSchedule()
		return m.Digest.Schedule, s, err
	case "orchestration.schedule":
		s, err := m.Orchestration.OrchestrationSchedule()
		return m.Orchestration.Schedule, s, err
	}
	rest, ok := strings.CutPrefix(setting, "roles.")
	name, isRole := strings.CutSuffix(rest, ".schedule")
	if !ok || !isRole || !roleName.MatchString(name) {
		return "", nil, fmt.Errorf("%s is no manifest setting that holds a schedule; name roles.<name>.schedule, planning.schedule, analysis.schedule, digest.schedule, or orchestration.schedule", setting)
	}
	s, err := m.RoleSchedule(name)
	return m.Roles[name].Schedule, s, err
}

// RoleAgent is the agent the role name is started with: roles.<name>.agent
// merged over the project's agent, as PlanningAgent is; for a built-in role,
// the agent its own block gives, such as PlanningAgent for the planner. Nil
// when neither sets anything.
func (m Manifest) RoleAgent(name string) *Agent {
	if b, builtin := builtinAgents[name]; builtin {
		return b.agent(m)
	}
	return m.Agent.With(m.Roles[name].Agent)
}

// UndeclaredRoleSettings are the roles' settings that nothing reads, one
// sentence each, by role name: a role no role file defines, a permission of
// the orchestrator's turned on here rather than under
// orchestration.permissions, and a permission turned on that its role file
// declares no permissions entry for. declared
// holds, by the name of each role a role file defines, the permissions its
// file's entries name as roles.<name>.permissions.<p>.
func (m Manifest) UndeclaredRoleSettings(declared map[string][]string) []string {
	var out []string
	for _, name := range slices.Sorted(maps.Keys(m.Roles)) {
		key := "roles." + name
		perms, defined := declared[name]
		switch {
		case !roleName.MatchString(name):
			continue // RoleErrors reports it
		case !defined:
			out = append(out, fmt.Sprintf("%s names no role a role file defines, so nothing reads it; give the role a file under workflows/roles in the design folder, or remove it", key))
			continue
		}
		p := m.Roles[name].Permissions
		for _, perm := range PermissionNames {
			switch {
			case !p.Allows(perm):
			case name == "orchestrator":
				out = append(out, fmt.Sprintf("%s.permissions.%s is not read: the orchestrator's permissions are orchestration.permissions; set orchestration.permissions.%s instead, and remove this", key, perm, perm))
			case !slices.Contains(perms, perm):
				out = append(out, fmt.Sprintf("%s.permissions.%s is on, but the role's file has no permissions entry for it, so it allows nothing; add one with setting %s.permissions.%s to the file, or remove it", key, perm, key, perm))
			}
		}
	}
	return out
}

// OrchestrationAgent is the agent the orchestrator is started with:
// orchestration.agent merged over the project's agent, as PlanningAgent is.
// Nil when neither sets anything.
func (m Manifest) OrchestrationAgent() *Agent {
	return m.Agent.With(m.Orchestration.Agent)
}
