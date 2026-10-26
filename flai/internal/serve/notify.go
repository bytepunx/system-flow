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
// where things are.

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

// Send sends r through the channel ch: a POST of its JSON to ch.URL, one
// attempt within NotifyURLTimeout, an answer other than 2xx being an error,
// and a run of ch.Command through the host's shell, sh -c, or cmd /C on
// Windows, with its JSON on standard input, within NotifyCommandTimeout. It
// does each that is set, the second whether or not the first failed, and
// returns both errors joined. With neither set it does nothing and returns
// nil. An error names the webhook's host alone, never its URL, which may
// carry a secret, nor the command, which may too.
func Send(ctx context.Context, ch config.Notify, r Reminder) error {
	return notifier{client: http.DefaultClient, shell: runShell, urlTimeout: NotifyURLTimeout, commandTimeout: NotifyCommandTimeout}.send(ctx, ch, r)
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

func (n notifier) send(ctx context.Context, ch config.Notify, r Reminder) error {
	if ch.URL == "" && strings.TrimSpace(ch.Command) == "" {
		return nil
	}
	body, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("encode the reminder of %s: %w", r.Thread, err)
	}
	var errs []error
	if ch.URL != "" {
		errs = append(errs, n.post(ctx, ch.URL, r.Thread, body))
	}
	if strings.TrimSpace(ch.Command) != "" {
		errs = append(errs, n.run(ctx, ch.Command, r.Thread, append(body, '\n')))
	}
	return errors.Join(errs...)
}

// post POSTs body to target once, within n.urlTimeout.
func (n notifier) post(ctx context.Context, target, thread string, body []byte) error {
	host := "an unreadable URL"
	if u, err := url.Parse(target); err == nil {
		host = u.Host
	}
	fail := func(what string) error {
		return fmt.Errorf("post the reminder of %s to the webhook at %s: %s; check that it answers, or change it with flai serve notify set --url <url>", thread, host, what)
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

// run runs command through the shell with stdin, within n.commandTimeout.
func (n notifier) run(ctx context.Context, command, thread string, stdin []byte) error {
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
	return fmt.Errorf("run the notify command for the reminder of %s: %s; try it in a shell with the reminder's JSON on its standard input, or change it with flai serve notify set --command <command>", thread, what)
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
