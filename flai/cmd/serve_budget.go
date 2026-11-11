package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/serve"
)

// flai serve budget: the host's budget states (S-0534, ADR-0209 §5), which
// hold every start on a provider, or on a harness's program and login, once
// a run ended on its usage limit, and their clear by hand (S-0720).

// budgetClearWhy is why a clear by hand cleared a state, as its journal entry
// and its notice say.
const budgetClearWhy = "cleared by hand with flai serve budget clear"

func newServeBudgetCmd(a *app) *cobra.Command {
	c := &cobra.Command{
		Use:   "budget",
		Short: "The budget states that hold starts on this host after a usage limit, and how long each holds",
		Long: `When a run ends on its harness's usage limit, flai serve sets a budget state
for the account the run went through, and holds every start on it, on every
project it serves on this host, until the limit resets. The account is the
provider the run named, or else its harness, the program the harness runs,
and the login the host gives it, such as claude-code meclod, so that a limit
on one Claude login holds no start on another. A state set by an older flai
is kept under the harness alone, such as claude-code, and holds every start
on that harness until it clears.

A state holds until the reset the API reported, or else a backoff that
starts at 5m and doubles up to 1h. Then one start goes through as the probe
and every other waits while it runs: its end on anything but the limit
clears the state, and its end on the limit holds it again.

This command lists each state: its key, the limit and status the harness
named, when it was set, until when it holds, whether it holds a start now,
and the probe that runs, with what it works and where. flai serve budget
clear <key> lifts one now, and with no key every one, so that the next start
on it goes; each clear is journalled (flai serve journal) and sent through
the channel flai serve notify sets. A change of a harness's program, its
arguments, or its login with flai serve agent harness clears the state of the
old one on its own.`,
		Example: `  flai serve budget
  flai serve budget --json
  flai serve budget clear 'claude-code meclod'
  flai serve budget clear`,
		Args: cobra.NoArgs,
		RunE: func(*cobra.Command, []string) error { return a.showBudgets() },
	}
	c.AddCommand(&cobra.Command{
		Use:   "clear [<key>]",
		Short: "Clear one budget state by its key, or every one, so that the next start on it goes",
		Long: `Removes the budget state kept under key, as flai serve budget lists it, or
every state when no key is given, so that flai serve holds no start on it
and the next start goes. A key with more than one word is one argument:
quote it, as flai serve budget clear 'claude-code meclod'. Each state cleared
is journalled, with method serve.budget, and its notice is sent through the
channel flai serve notify sets, when one is set. A budget state's file that
does not parse is replaced with none by a clear with no key.`,
		Example: `  flai serve budget clear 'claude-code meclod'
  flai serve budget clear anthropic --json
  flai serve budget clear`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error { return a.clearBudgets(cmd.Context(), args) },
	})
	return c
}

// budgetView is a budget state as flai serve budget --json prints it: the
// state, and whether it holds a start at the time it was read.
type budgetView struct {
	serve.Budget
	Held bool `json:"held"`
}

// showBudgets prints the host's budget states, each with whether it holds a
// start now.
func (a *app) showBudgets() error {
	states, err := a.serveDir().Budgets()
	if err != nil {
		return fmt.Errorf("read the host's budget state: %w", err)
	}
	now := a.now()
	if a.jsonOut {
		out := make([]budgetView, 0, len(states))
		for _, b := range states {
			out = append(out, budgetView{Budget: b, Held: b.Held(now)})
		}
		return a.printJSON(out)
	}
	if len(states) == 0 {
		fmt.Fprintln(a.out, "no budget state is set; flai serve holds no start for a usage limit")
		return nil
	}
	for _, b := range states {
		fmt.Fprintf(a.out, "%s: %s\n", b.Provider, budgetHolds(b, now))
		fmt.Fprintf(a.out, "  limit: %s, set %s\n", budgetLimit(b), b.Since.UTC().Format(time.RFC3339))
		if p := b.Probe; p != nil {
			fmt.Fprintf(a.out, "  probe: %s in %s, started %s\n", p.Work, p.Root, p.Started)
		}
		fmt.Fprintf(a.out, "  clear it with: %s\n", serve.BudgetClearCommand(b.Provider))
	}
	return nil
}

// budgetHolds says whether b holds a start at now and until when: until its
// until, with the reset the API reported or the backoff it took; after it,
// while its probe runs; else that the next start goes as the probe.
func budgetHolds(b serve.Budget, now time.Time) string {
	until := b.Until.UTC().Format(time.RFC3339)
	switch {
	case now.Before(b.Until) && !b.Reset.IsZero() && b.Reset.Equal(b.Until):
		return "holds every start on it until " + until + ", the reset the API reported"
	case now.Before(b.Until):
		return "holds every start on it until " + until + ", a backoff of " + b.Backoff + ", the API naming no reset"
	case b.Held(now):
		return "holds every start on it while its probe runs, its hold having ended at " + until
	}
	return "holds no start now: its hold ended at " + until + ", and the next start on it goes as the probe"
}

// budgetLimit names the limit b is on: the error the harness named, with its
// HTTP status, or the status alone, or the harness's usage limit.
func budgetLimit(b serve.Budget) string {
	limit := b.Limit
	switch {
	case limit != "" && b.Status != 0:
		limit += " (" + strconv.Itoa(b.Status) + ")"
	case b.Status != 0:
		limit = "HTTP " + strconv.Itoa(b.Status)
	case limit == "":
		limit = "the harness's usage limit"
	}
	return limit
}

// budgetCleared is what flai serve budget clear --json prints: the states it
// removed, and, when the file did not parse and was replaced, why.
type budgetCleared struct {
	Removed []serve.Budget `json:"removed"`
	Unread  string         `json:"unread,omitempty"`
}

// clearBudgets clears the budget state of the key in args, or every one with
// none, journals each state it removed and sends its notice, and says what it
// cleared. A key that names no state is an error naming those there are.
func (a *app) clearBudgets(ctx context.Context, args []string) error {
	dir := a.serveDir()
	res, err := dir.ClearBudgets(args)
	if err != nil {
		return err
	}
	if len(args) == 1 && len(res.Removed) == 0 {
		return noBudgetState(dir, args[0])
	}
	cfg, _, err := a.loadConfig()
	if err != nil {
		a.logger().Warn("budget notices not sent: the host configuration is not read", "component", "serve", "err", err.Error())
	}
	a.announceBudgetClears(ctx, res.Removed, budgetClearing{why: budgetClearWhy, by: "flai serve budget clear"}, cfg.Notify)
	out := budgetCleared{Removed: res.Removed}
	if out.Removed == nil {
		out.Removed = []serve.Budget{}
	}
	if res.Unread != nil {
		out.Unread = res.Unread.Error()
	}
	if a.jsonOut {
		return a.printJSON(out)
	}
	if res.Unread != nil {
		fmt.Fprintf(a.out, "replaced %s, which did not parse (%v), with no budget state; flai serve holds no start until a run ends on a usage limit again\n",
			filepath.Join(string(dir), "budget.json"), res.Unread)
	}
	if len(res.Removed) == 0 && res.Unread == nil {
		fmt.Fprintln(a.out, "no budget state is set; nothing to clear")
	}
	for _, b := range res.Removed {
		fmt.Fprintf(a.out, "cleared the budget state of %s, on %s until %s; starts on it go again\n", b.Provider, budgetLimit(b), b.Until.UTC().Format(time.RFC3339))
	}
	return nil
}

// noBudgetState is the error of a clear of key, which names no state: the
// keys there are, or that there is none.
func noBudgetState(dir serve.Dir, key string) error {
	states, err := dir.Budgets()
	if err != nil {
		return fmt.Errorf("there is no budget state of %q, and the others are not read: %w", key, err)
	}
	if len(states) == 0 {
		return fmt.Errorf("there is no budget state of %q: none is set, so flai serve holds no start for a usage limit", key)
	}
	keys := make([]string, 0, len(states))
	for _, b := range states {
		keys = append(keys, strconv.Quote(b.Provider))
	}
	return fmt.Errorf("there is no budget state of %q; there %s: %s (flai serve budget lists them; quote a key of more than one word)",
		key, oneOrMany(len(keys), "is", "are"), strings.Join(keys, ", "))
}

// budgetClearing is a clear as it is announced: why, by whom, and the
// project and root its journal entries name, empty for a clear of the host's.
type budgetClearing struct{ why, by, project, root string }

// announceBudgetClears journals each state removed under serve.budget with
// c's why, as flai serve journals its own clears, and sends its notice
// through notify when a channel is set; a failure of either is warned of,
// since the state is cleared already.
func (a *app) announceBudgetClears(ctx context.Context, removed []serve.Budget, c budgetClearing, notify config.Notify) {
	if len(removed) == 0 {
		return
	}
	now := a.now()
	if by := a.writer(); by != "agent" {
		c.by = by
	}
	for _, b := range removed {
		e := serve.BudgetClearedEntry(b, c.why, c.by, now)
		e.Project, e.Root = c.project, c.root
		if err := a.journal(e); err != nil {
			a.logger().Warn("budget clear not journalled", "component", "serve", "provider", b.Provider, "err", err.Error())
		}
	}
	if notify.URL == "" && strings.TrimSpace(notify.Command) == "" {
		return
	}
	all, err := a.serveDir().Projects()
	if err != nil {
		a.logger().Warn("budget notice names no project: the registry is not read", "component", "serve", "err", err.Error())
	}
	projects := serve.ProjectNames(all)
	for _, b := range removed {
		if err := serve.SendBudget(ctx, notify, serve.BudgetClearedNotice(b, c.why, now, projects)); err != nil {
			a.logger().Warn("budget notice not sent", "component", "serve", "provider", b.Provider, "err", err.Error())
		}
	}
}
