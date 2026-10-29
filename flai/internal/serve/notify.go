package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/config"
	"github.com/bytepunx/system-flow/flai/internal/threads"
)

// Reminding the operator off the dashboard (S-0498, ADR-0205). flai serve
// sends a Reminder of each thread that has awaited the operator past the
// manifest's threads.remind_after through the channel the host configuration
// keeps, config.Notify: a POST of its JSON to a webhook, a shell command with
// its JSON on standard input, or both. The body is the one flaiover posts to
// dashboard.notify_url, project and entry, extended with the thread's own
// fields, so that one receiver can take both. It holds no token and no file
// content: a thread's title, the first line of a pending recommendation, and
// where things are. A change of the host's budget state goes through the
// same channel as a BudgetNotice (S-0534).

const (
	// NotifyURLTimeout bounds the one POST of a reminder to the webhook.
	NotifyURLTimeout = 5 * time.Second
	// NotifyCommandTimeout bounds one run of the reminder's shell command.
	NotifyCommandTimeout = 30 * time.Second
	// notifySaid is how much of a failing command's output its error quotes.
	notifySaid = 300
)

// Reminder is what flai serve sends of one thread that awaits the operator.
type Reminder struct {
	// Project is the served project's name, as its manifest says it.
	Project string `json:"project"`
	// Thread is the thread's ID, TH-nnnn.
	Thread string `json:"thread"`
	// Title is the thread's title.
	Title string `json:"title"`
	// Item is the ID of the work item the thread is on, or empty for a
	// thread on a document.
	Item string `json:"item"`
	// AwaitingSeconds is how long the thread has awaited the operator: from
	// Since to when the reminder was made.
	AwaitingSeconds int64 `json:"awaiting_seconds"`
	// Since is when the thread's last entry, the one it awaits the operator
	// on, was written.
	Since time.Time `json:"since"`
	// Recommendation is the recommendation awaiting the operator's
	// confirmation, or nil when none is pending.
	Recommendation *ReminderRecommendation `json:"recommendation,omitempty"`
	// Href is the thread's link in the dashboard, a path such as
	// /items/S-0001?thread=TH-0001; ReminderHref makes it.
	Href string `json:"href"`
}

// ReminderRecommendation is a pending recommendation, as a reminder carries
// it: who wrote it, its first line, and what it cites.
type ReminderRecommendation struct {
	// Author is who wrote the recommendation.
	Author string `json:"author"`
	// FirstLine is the first line of its text.
	FirstLine string `json:"first_line"`
	// Source is what it cites, or nil when it cites nothing.
	Source *threads.Source `json:"source,omitempty"`
}

// reminderEntry is the entry flaiover posts to dashboard.notify_url, which a
// reminder carries too, for a receiver of both.
type reminderEntry struct {
	Key   string `json:"key"`
	Kind  string `json:"kind"`
	Title string `json:"title"`
	Href  string `json:"href"`
	At    string `json:"at"`
}

// MarshalJSON writes the reminder with Since in UTC to the second and with
// flaiover's entry, kind thread and keyed thread:TH-nnnn as the dashboard's
// inbox keys it.
func (r Reminder) MarshalJSON() ([]byte, error) {
	type plain Reminder
	since := ""
	if !r.Since.IsZero() {
		since = r.Since.UTC().Format(time.RFC3339)
	}
	return json.Marshal(struct {
		plain
		Since string        `json:"since"`
		Entry reminderEntry `json:"entry"`
	}{
		plain: plain(r),
		Since: since,
		Entry: reminderEntry{Key: "thread:" + r.Thread, Kind: "thread", Title: r.Title, Href: r.Href, At: since},
	})
}

// ReminderHref is the dashboard's link to thread, as flaiover's inbox links
// one: the page of item, opened on the thread, or, for a thread on a
// document, the threads page opened on it.
func ReminderHref(item, thread string) string {
	if item != "" {
		return "/items/" + url.PathEscape(item) + "?thread=" + url.QueryEscape(thread)
	}
	return "/threads?thread=" + url.QueryEscape(thread)
}

// BudgetHref is the dashboard's link a budget notice carries: its root, as
// the budget state's banner is on every page.
const BudgetHref = "/"

// BudgetNotice is what flai serve sends of a change of one provider's budget
// state (S-0534, ADR-0209 §5) through the reminder's channel: once per
// change, for the host, not per project. Its JSON has kind budget, which a
// thread's reminder has not, and flaiover's entry, of kind budget, as a
// reminder's has of kind thread, so that one receiver can take both. It holds
// no token and no file content.
type BudgetNotice struct {
	// State is the change: set, extended, or cleared.
	State string `json:"state"`
	// Provider is the provider the state is kept under, or the runs'
	// harness when they named none.
	Provider string `json:"provider"`
	// Limit is the error the harness named, such as usage_limit_reached, and
	// Status its HTTP status; each empty when it named none.
	Limit  string `json:"limit"`
	Status int    `json:"status,omitempty"`
	// Since is when the state was set, and Until when its hold ends: for a
	// clear, as they were.
	Since time.Time `json:"since"`
	Until time.Time `json:"until"`
	// At is when the state changed.
	At time.Time `json:"at"`
	// Projects are the names of the projects the host serves, each of which
	// the state holds.
	Projects []string `json:"projects"`
	// Href is the dashboard's root, BudgetHref.
	Href string `json:"href"`
}

// budgetNoticeTitles are what each change says, by state, of provider, in
// flaiover's entry.
var budgetNoticeTitles = map[string]string{
	budgetSet:      "%s is out of budget: every start on it waits until %s",
	budgetExtended: "%s is still out of budget: every start on it waits until %s",
	budgetCleared:  "%s has budget again: starts on it go again, the hold that was until %s over",
}

// MarshalJSON writes the notice with kind budget, its times in UTC to the
// second, and flaiover's entry, kind budget and keyed budget:<provider>.
func (n BudgetNotice) MarshalJSON() ([]byte, error) {
	type plain BudgetNotice
	at := func(t time.Time) string {
		if t.IsZero() {
			return ""
		}
		return t.UTC().Format(time.RFC3339)
	}
	projects := n.Projects
	if projects == nil {
		projects = []string{}
	}
	title := n.Provider + " budget " + n.State
	if f, ok := budgetNoticeTitles[n.State]; ok {
		title = fmt.Sprintf(f, n.Provider, at(n.Until))
	}
	return json.Marshal(struct {
		Kind string `json:"kind"`
		plain
		Since    string        `json:"since"`
		Until    string        `json:"until"`
		At       string        `json:"at"`
		Projects []string      `json:"projects"`
		Entry    reminderEntry `json:"entry"`
	}{
		Kind:     "budget",
		plain:    plain(n),
		Since:    at(n.Since),
		Until:    at(n.Until),
		At:       at(n.At),
		Projects: projects,
		Entry:    reminderEntry{Key: "budget:" + n.Provider, Kind: "budget", Title: title, Href: n.Href, At: at(n.At)},
	})
}

// Send sends r through the channel ch: a POST of its JSON to ch.URL, one
// attempt within NotifyURLTimeout, an answer other than 2xx being an error,
// and a run of ch.Command through the host's shell, sh -c, or cmd /C on
// Windows, with its JSON on standard input, within NotifyCommandTimeout. It
// does each that is set, the second whether or not the first failed, and
// returns both errors joined. With neither set it does nothing and returns
// nil. An error names the webhook's host alone, never its URL, which may
// carry a secret, nor the command, which may too.
func Send(ctx context.Context, ch config.Notify, r Reminder) error {
	return defaultNotifier().send(ctx, ch, "the reminder of "+r.Thread, r)
}

// SendBudget sends n through the channel ch, as Send sends a reminder.
func SendBudget(ctx context.Context, ch config.Notify, n BudgetNotice) error {
	return defaultNotifier().send(ctx, ch, "the budget notice of "+n.Provider, n)
}

// defaultNotifier is the notifier of Send: the network, the host's shell,
// and the reminder's timeouts.
func defaultNotifier() notifier {
	return notifier{client: http.DefaultClient, shell: runShell, urlTimeout: NotifyURLTimeout, commandTimeout: NotifyCommandTimeout}
}

// shellRunner runs command through the host's shell with stdin on its
// standard input until it ends or ctx does, and answers what it wrote on
// standard output and standard error, with an error when it could not run
// or exited non-zero.
type shellRunner func(ctx context.Context, command string, stdin []byte) ([]byte, error)

// notifier is Send with its I/O and its bounds apart, for tests.
type notifier struct {
	client         *http.Client
	shell          shellRunner
	urlTimeout     time.Duration
	commandTimeout time.Duration
}

// send sends v's JSON through ch, as Send does; about names what v is, such
// as "the reminder of TH-0001", for its errors.
func (n notifier) send(ctx context.Context, ch config.Notify, about string, v any) error {
	if ch.URL == "" && strings.TrimSpace(ch.Command) == "" {
		return nil
	}
	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("encode %s: %w", about, err)
	}
	var errs []error
	if ch.URL != "" {
		errs = append(errs, n.post(ctx, ch.URL, about, body))
	}
	if strings.TrimSpace(ch.Command) != "" {
		errs = append(errs, n.run(ctx, ch.Command, about, append(body, '\n')))
	}
	return errors.Join(errs...)
}

// post POSTs body, which is about, to target once, within n.urlTimeout.
func (n notifier) post(ctx context.Context, target, about string, body []byte) error {
	host := "an unreadable URL"
	if u, err := url.Parse(target); err == nil {
		host = u.Host
	}
	fail := func(what string) error {
		return fmt.Errorf("post %s to the webhook at %s: %s; check that it answers, or change it with flai serve notify set --url <url>", about, host, what)
	}
	ctx, cancel := context.WithTimeout(ctx, n.urlTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(body))
	if err != nil {
		return fail("the URL cannot be requested")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "flai")
	resp, err := n.client.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fail(fmt.Sprintf("it gave no answer within %s", n.urlTimeout))
		}
		// *url.Error repeats the URL, which may carry a secret
		var ue *url.Error
		if errors.As(err, &ue) {
			err = ue.Err
		}
		return fail(err.Error())
	}
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
	_ = resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fail("it answered " + resp.Status)
	}
	return nil
}

// run runs command through the shell with stdin, which is about, within
// n.commandTimeout.
func (n notifier) run(ctx context.Context, command, about string, stdin []byte) error {
	ctx, cancel := context.WithTimeout(ctx, n.commandTimeout)
	defer cancel()
	out, err := n.shell(ctx, command, stdin)
	if err == nil {
		return nil
	}
	what := err.Error()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		what = fmt.Sprintf("it did not end within %s", n.commandTimeout)
	}
	if said := strings.TrimSpace(string(out)); said != "" {
		if len(said) > notifySaid {
			said = "..." + said[len(said)-notifySaid:]
		}
		what += ": " + said
	}
	return fmt.Errorf("run the notify command for %s: %s; try it in a shell with its JSON on its standard input, or change it with flai serve notify set --command <command>", about, what)
}

// runShell is the shellRunner that runs command through sh -c, or cmd /C on
// Windows, in flai serve's own environment and directory.
func runShell(ctx context.Context, command string, stdin []byte) ([]byte, error) {
	shell, flag := "sh", "-c"
	if runtime.GOOS == "windows" {
		shell, flag = "cmd", "/C"
	}
	cmd := exec.CommandContext(ctx, shell, flag, command)
	cmd.Stdin = bytes.NewReader(stdin)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	// a child the shell started may hold the output open past the shell
	cmd.WaitDelay = 2 * time.Second
	err := cmd.Run()
	return out.Bytes(), err
}
