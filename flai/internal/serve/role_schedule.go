package serve

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/cron"
	"github.com/bytepunx/system-flow/flai/internal/manifest"
	"github.com/bytepunx/system-flow/flai/internal/roles"
	"github.com/bytepunx/system-flow/flai/internal/workitem"
)

// A role a project defines, started on a schedule (S-0418, ADR-0207 §3).
//
// A role scheduler per served project looks whenever the project's launcher
// does, as the analysis scheduler does (S-0223). At each look it reads the
// project's role files, and for each role that is not built-in and whose
// trigger has a {schedule: <setting>} entry it follows that setting: the
// role's own roles.<name>.schedule, or one of planning.schedule,
// analysis.schedule, digest.schedule, and orchestration.schedule, which the
// role then shares with the built-in role that setting is for. Each time the
// setting comes round with the host action the role's file names on, it
// starts the role through startRole, with the schedule as its trigger
// ("schedule daily"): a role that takes no item once, on no item, under an
// entry with no condition; a role whose file takes stories once for each
// ready story, in the board's pull order (roles.EachReadyStory), under an
// entry with when: each-ready-story (S-0571, ADR-0250 §5). A role that takes
// items under an entry with no condition, and one that takes no story under
// each-ready-story, is not started, which is logged and which flai check
// reports. A start the checks refuse, as while a run of the role goes on the
// same item, is logged, naming the item, and not tried again until the
// schedule next comes round; a start the budget state holds is kept, per
// item, and tried again at each look until it starts (S-0534), or until its
// story is no longer ready. Times missed between two looks come round once,
// and a schedule set or changed, or a role or trigger added, comes round
// first at its next time from then. The built-in roles' own launchers follow
// the settings their files name, by the same functions (scheduleSettings,
// scheduleSettingValid, matching), and start their roles themselves
// (ADR-0250 §4).
//
// What it has seen is kept in memory only, as the analysis scheduler's is.
// While a role's host action is off its schedules start nothing and keep up,
// so that turning it on does not act on the past.

// roleScheduler starts the roles one project defines on the schedules their
// files name (S-0418).
type roleScheduler struct {
	o       Options
	e       Entry
	starter *launcher

	mu sync.Mutex
	// follows are the schedules followed, by role and setting.
	follows map[roleSetting]*roleFollow
	said    map[string]bool // the problems warned of, and the holds logged (heldOnce)
}

// roleSetting is a role's name and a manifest setting its trigger names a
// schedule by.
type roleSetting struct {
	role, setting string
}

// roleFollow is how a role scheduler follows one role's schedule: its
// cadence, and held, the triggers of the starts the budget state held, by
// the item each is on ("" for none), which each look tries again until they
// start.
type roleFollow struct {
	cadence cadence
	held    map[string]string
}

// newRoleScheduler is the role scheduler for project e, starting its roles
// with starter, the project's launcher. Times before now do not come round.
func newRoleScheduler(o Options, e Entry, starter *launcher) *roleScheduler {
	s := &roleScheduler{o: o, e: e, starter: starter, follows: map[roleSetting]*roleFollow{}, said: map[string]bool{}}
	now := o.Now().UTC()
	if repo, err := workitem.Open(e.Root); err == nil {
		if files, _, err := roles.Load(repo.Root, repo.Manifest.Layout["design"]); err == nil {
			for _, rs := range scheduledRoles(files) {
				s.follow(rs.key(), repo.Manifest, now)
			}
		}
	}
	return s
}

// look starts each role whose schedule has come round since the last look,
// or whose start was held by the budget state, while its host action is on,
// unless the checks refuse it. Schedules no role file names any longer are
// no longer followed.
func (s *roleScheduler) look(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.o.Now().UTC()
	repo, err := workitem.Open(s.e.Root)
	if err != nil {
		s.warnOnce("role scheduler could not read the project", err)
		return
	}
	files, _, err := roles.Load(repo.Root, repo.Manifest.Layout["design"])
	if err != nil {
		s.warnOnce("role scheduler could not read the role files", err)
		return
	}
	var cfg AgentConfig
	if s.o.Agent != nil {
		cfg = s.o.Agent(s.e.Root)
	}
	follows := map[roleSetting]*roleFollow{}
	for _, rs := range scheduledRoles(files) {
		k := rs.key()
		f := s.follow(k, repo.Manifest, now)
		follows[k] = f
		s.lookAt(ctx, rs.file, rs.setting, f, repo, cfg.RoleActions[rs.file.Action], now)
	}
	s.follows = follows
}

// lookAt starts the role file names when setting has come round for f by
// now, on what scheduleTargets gives, and tries again each start f holds
// that the budget state held, while on, its host action, is. While it is off
// the schedule keeps up and the held starts are dropped. A setting that is
// not one the role may name, or whose value does not parse, is warned of
// once and starts nothing.
func (s *roleScheduler) lookAt(ctx context.Context, file roles.Role, setting string, f *roleFollow, repo *workitem.Repo, on bool, now time.Time) {
	if err := scheduleSettingValid(file, setting); err != nil {
		s.warnOnce("role schedule not valid", err, "role", file.Name)
		return
	}
	_, sched, err := repo.Manifest.ScheduleSetting(setting)
	if err != nil {
		s.warnOnce("role schedule not valid", err, "role", file.Name)
		return
	}
	if !on {
		f.cadence.mark, f.held = now, nil
		return
	}
	trigger := ""
	if f.cadence.due(sched, now) {
		trigger = scheduleTrigger(sched)
	}
	if trigger == "" && len(f.held) == 0 {
		return
	}
	targets, ok := s.scheduleTargets(file, setting, repo, trigger)
	if !ok {
		f.held = nil
		return
	}
	held := f.held
	f.held = nil
	s.starter.mu.Lock()
	defer s.starter.mu.Unlock()
	for _, item := range targets {
		t := trigger
		if t == "" {
			t = held[item]
		}
		if t == "" {
			continue
		}
		s.startScheduled(ctx, file, item, t, f)
	}
}

// scheduleTargets are the items the role file names is started on when
// setting comes round, in order: "", for no item, for a role that takes
// none under an entry with no condition; and the ready stories in the
// board's pull order for one that takes stories under an entry with when:
// each-ready-story (ADR-0250 §5). It is false, with why logged when the
// setting came round (trigger set), when the file's entries start the role
// on nothing: a schedule gives a role that takes items no item, and
// each-ready-story gives one that takes no story none, as flai check
// reports.
func (s *roleScheduler) scheduleTargets(file roles.Role, setting string, repo *workitem.Repo, trigger string) ([]string, bool) {
	plain, each := false, false
	for _, t := range matching(file.Trigger, roles.Occurrence{Schedule: setting}) {
		switch t.When {
		case "":
			plain = true
		case roles.CondEachReadyStory:
			each = true
		}
	}
	notStarted := func(reason string) ([]string, bool) {
		if trigger != "" {
			s.starter.log("scheduled role not started", "role", file.Name, "trigger", trigger, "reason", reason)
		}
		return nil, false
	}
	switch {
	case each && slices.Contains(file.Item, string(workitem.Story)):
		items, err := repo.List(true)
		if err != nil {
			s.warnOnce("role scheduler could not read the project", err)
			return nil, false
		}
		board, err := repo.LoadBoard()
		if err != nil {
			s.warnOnce("role scheduler could not read the project", err)
			return nil, false
		}
		return roles.EachReadyStory(board.Order, items), true
	case each && !plain:
		return notStarted("the role takes no story, as its file " + file.Path + " says, and " + roles.CondEachReadyStory + " starts it only on ready stories")
	case len(file.Item) > 0:
		return notStarted("the role runs on an item, as its file " + file.Path + " says, and a schedule gives it none")
	case plain:
		return []string{""}, true
	}
	return nil, false
}

// startScheduled starts the role file names on item, "" for none, with
// trigger, logging a start the checks refuse with the item it was for, and
// keeping in f a start the budget state holds, to try again at the next
// look. The caller holds the launcher's lock.
func (s *roleScheduler) startScheduled(ctx context.Context, file roles.Role, item, trigger string, f *roleFollow) {
	args := []any{"role", file.Name}
	if item != "" {
		args = append(args, "item", item)
	}
	args = append(args, "trigger", trigger)
	_, held, err := startRole(ctx, s.o, s.e, s.starter, file.Name, item, trigger)
	var no *Refused
	switch {
	case errors.As(err, &no):
		s.starter.log("scheduled role not started", append(args, "reason", no.Why)...)
	case err != nil:
		s.starter.warn("scheduled role not started", append(args, "err", err)...)
	case held != nil:
		if f.held == nil {
			f.held = map[string]string{}
		}
		f.held[item] = trigger
		s.starter.heldOnce(s.said, "scheduled role held", *held, args...)
	}
}

// follow is how k is followed: as before, or from now when it was not, and
// with the setting's value in m seen at now, so that a schedule set or
// changed comes round first at its next time from now.
func (s *roleScheduler) follow(k roleSetting, m manifest.Manifest, now time.Time) *roleFollow {
	f := s.follows[k]
	if f == nil {
		f = &roleFollow{cadence: cadence{mark: now}}
		s.follows[k] = f
	}
	spec, _, _ := m.ScheduleSetting(k.setting)
	f.cadence.see(spec, now)
	return f
}

// warnOnce warns of err with msg, and args, the first time it is seen.
func (s *roleScheduler) warnOnce(msg string, err error, args ...any) {
	if key := msg + "\x00" + err.Error(); !s.said[key] {
		s.said[key] = true
		s.starter.warn(msg, append(args, "err", err)...)
	}
}

// scheduledRole is a role file and a setting one of its trigger entries
// names a schedule by.
type scheduledRole struct {
	file    roles.Role
	setting string
}

// key is the role and setting r follows.
func (r scheduledRole) key() roleSetting { return roleSetting{role: r.file.Name, setting: r.setting} }

// scheduledRoles are the settings files' schedule triggers name, each role's
// once, for each role that is not built-in. The built-in roles' own
// launchers start them, judging the same entries by scheduleSettings and
// matching (ADR-0250 §4): the planner's replanner here, and the analyzer's
// and the curator's schedulers.
func scheduledRoles(files []roles.Role) []scheduledRole {
	var out []scheduledRole
	for _, f := range files {
		if _, builtin := builtinRoles[f.Name]; builtin {
			continue
		}
		for _, setting := range scheduleSettings(f) {
			out = append(out, scheduledRole{file: f, setting: setting})
		}
	}
	return out
}

// scheduleSettings are the settings f's schedule trigger entries name, each
// once, in the file's order: what a launcher of f follows a cadence for.
// Whether a setting's coming round starts f, and how, is matching's to say
// of a schedule occurrence on it, by the entries' conditions.
func scheduleSettings(f roles.Role) []string {
	var out []string
	for _, t := range f.Trigger {
		if t.Kind == roles.TriggerSchedule && !slices.Contains(out, t.Schedule) {
			out = append(out, t.Schedule)
		}
	}
	return out
}

// scheduleSettingValid is nil when setting is one f's schedule entry may
// name, roles.<name>.schedule or one of roles.Schedules, and otherwise what
// flai check reports of it.
func scheduleSettingValid(f roles.Role, setting string) error {
	if setting == "roles."+f.Name+".schedule" || slices.Contains(roles.Schedules, setting) {
		return nil
	}
	return fmt.Errorf("the role %s's file %s names the schedule %s, which is neither roles.%s.schedule nor one of %s, as flai check reports",
		f.Name, f.Path, setting, f.Name, strings.Join(roles.Schedules, ", "))
}

// The planner's, the analyzer's, and the curator's launchers follow the
// settings their role files' schedule entries name by these (ADR-0250 §4).

// startCadences are the cadences of the settings f's schedule entries name,
// by setting, each following the setting's value in m from now, so that
// times before now do not come round.
func startCadences(f roles.Role, m manifest.Manifest, now time.Time) map[string]*cadence {
	out := map[string]*cadence{}
	for _, setting := range scheduleSettings(f) {
		spec, _, _ := m.ScheduleSetting(setting)
		out[setting] = &cadence{spec: strings.TrimSpace(spec), mark: now}
	}
	return out
}

// followSchedules are the settings f's schedule entries name, by setting,
// parsed from m, and the cadences that follow them: each setting's from
// cadences, or a new one from now for a setting newly named, seen at now, so
// that one set or changed comes round first at its next time from now.
// Cadences of settings f no longer names are let go. A setting that is
// unset, that f may not name, or that cannot be parsed is left out of the
// schedules, and the latter two are passed to warn.
func followSchedules(f roles.Role, cadences map[string]*cadence, m manifest.Manifest, now time.Time, warn func(error)) (map[string]*cron.Schedule, map[string]*cadence) {
	scheds, follows := map[string]*cron.Schedule{}, map[string]*cadence{}
	for _, setting := range scheduleSettings(f) {
		c := cadences[setting]
		if c == nil {
			c = &cadence{mark: now}
		}
		follows[setting] = c
		if err := scheduleSettingValid(f, setting); err != nil {
			warn(err)
			continue
		}
		spec, sched, err := m.ScheduleSetting(setting)
		c.see(spec, now)
		if err != nil {
			warn(err)
			continue
		}
		if sched != nil {
			scheds[setting] = sched
		}
	}
	return scheds, follows
}

// markAll marks each of cadences at now, so that times before now do not
// come round: what a scheduler does while its host action is off.
func markAll(cadences map[string]*cadence, now time.Time) {
	for _, c := range cadences {
		c.mark = now
	}
}

// servedPlainly reports whether o meets an entry among triggers that names
// no condition: the only entries the analyzer's and the curator's launchers
// serve, as flai check reports of the others (ADR-0250 §4).
func servedPlainly(triggers []roles.Trigger, o roles.Occurrence) bool {
	return slices.ContainsFunc(matching(triggers, o), func(t roles.Trigger) bool { return t.When == "" })
}
