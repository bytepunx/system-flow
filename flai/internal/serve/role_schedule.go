package serve

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

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
// starts the role on no item through startRole, with the schedule as its
// trigger ("schedule daily"). A start the checks refuse, as while a run of
// the role goes, is logged and not tried again until the schedule next
// comes round; a start the budget state holds is tried again at each look
// until it starts (S-0534). A role whose file takes items is not started:
// a schedule gives it no item to run on. Times missed between two looks come
// round once, and a schedule set or changed, or a role or trigger added,
// comes round first at its next time from then.
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
// cadence, and heldDue, the trigger of the start the budget state held,
// which each look tries again until it starts; "" when none waits.
type roleFollow struct {
	cadence cadence
	heldDue string
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
		s.lookAt(ctx, rs.file, rs.setting, f, repo.Manifest, cfg.RoleActions[rs.file.Action], now)
	}
	s.follows = follows
}

// lookAt starts the role file names when setting has come round for f by
// now, or when f holds a start the budget state held, while on, its host
// action, is. While it is off the schedule keeps up and a held start is
// dropped. A setting that is not one the role may name, or whose value does
// not parse, is warned of once and starts nothing.
func (s *roleScheduler) lookAt(ctx context.Context, file roles.Role, setting string, f *roleFollow, m manifest.Manifest, on bool, now time.Time) {
	if setting != "roles."+file.Name+".schedule" && !slices.Contains(roles.Schedules, setting) {
		s.warnOnce("role schedule not valid", fmt.Errorf("the role %s's file %s names the schedule %s, which is neither roles.%s.schedule nor one of %s, as flai check reports",
			file.Name, file.Path, setting, file.Name, strings.Join(roles.Schedules, ", ")), "role", file.Name)
		return
	}
	_, sched, err := m.ScheduleSetting(setting)
	if err != nil {
		s.warnOnce("role schedule not valid", err, "role", file.Name)
		return
	}
	if !on {
		f.cadence.mark, f.heldDue = now, ""
		return
	}
	trigger := f.heldDue
	if f.cadence.due(sched, now) {
		trigger = scheduleTrigger(sched)
	}
	if trigger == "" {
		return
	}
	f.heldDue = ""
	if len(file.Item) > 0 {
		s.starter.log("scheduled role not started", "role", file.Name, "trigger", trigger,
			"reason", "the role runs on an item, as its file "+file.Path+" says, and a schedule gives it none")
		return
	}
	s.starter.mu.Lock()
	defer s.starter.mu.Unlock()
	_, held, err := startRole(ctx, s.o, s.e, s.starter, file.Name, "", trigger)
	var no *Refused
	switch {
	case errors.As(err, &no):
		s.starter.log("scheduled role not started", "role", file.Name, "trigger", trigger, "reason", no.Why)
	case err != nil:
		s.starter.warn("scheduled role not started", "role", file.Name, "trigger", trigger, "err", err)
	case held != nil:
		f.heldDue = trigger
		s.starter.heldOnce(s.said, "scheduled role held", *held, "role", file.Name, "trigger", trigger)
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
// once, for each role that is not built-in: the built-in roles keep their
// own schedulers.
func scheduledRoles(files []roles.Role) []scheduledRole {
	var out []scheduledRole
	seen := map[roleSetting]bool{}
	for _, f := range files {
		if _, builtin := builtinRoles[f.Name]; builtin {
			continue
		}
		for _, t := range f.Trigger {
			r := scheduledRole{file: f, setting: t.Schedule}
			if t.Kind != roles.TriggerSchedule || seen[r.key()] {
				continue
			}
			seen[r.key()] = true
			out = append(out, r)
		}
	}
	return out
}
