package usage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

// spendLogDate is the layout of /spend/logs' start_date and end_date: a day,
// read as midnight UTC, end_date's included only at its midnight.
const spendLogDate = "2006-01-02"

// maxSpendLogBody bounds how much of the spend log is read: its rows carry
// the request's metadata, and the proxy answers up to 10,000 of them.
const maxSpendLogBody = 64 << 20

// spendLogTruncated is the header the proxy sets when it answered only its
// cap of the most recent rows.
const spendLogTruncated = "x-litellm-spend-logs-truncated"

// LiteLLM reads what a LiteLLM proxy charged for a run from its spend log,
// /spend/logs, which only the proxy's master key reads (S-0356).
type LiteLLM struct {
	root   string
	key    string
	client *http.Client
}

// NewLiteLLM is LiteLLM's spend reader for the proxy at baseURL, its root or
// its /v1, called with key, the proxy's master key.
func NewLiteLLM(baseURL, key string, client *http.Client) *LiteLLM {
	root := strings.TrimSuffix(strings.TrimRight(baseURL, "/"), "/v1")
	return &LiteLLM{root: root, key: key, client: spendClient(client)}
}

// spendLogRow is one row of the proxy's spend log, as far as a join reads it.
type spendLogRow struct {
	APIKey    string  `json:"api_key"`
	Spend     float64 `json:"spend"`
	StartTime string  `json:"startTime"`
	SessionID string  `json:"session_id"`
	EndUser   string  `json:"end_user"`
}

// Run is what the proxy charged for the run: its spend log's rows tagged
// with one of the run's sessions, else the rows of the key whose hash is
// KeyHash within the run's window.
func (l *LiteLLM) Run(ctx context.Context, run RunRef) (*Spent, error) {
	if len(run.Sessions) == 0 && run.KeyHash == "" {
		return nil, fmt.Errorf("LiteLLM: the run has no session ID and no hash of its key, and the window alone could take another story's rows: %w", ErrNoSpend)
	}
	rows, err := l.spendLog(ctx, run)
	if err != nil {
		return nil, err
	}
	spent, err := bySession(rows, run.Sessions, run.From, run.To)
	if err == nil && spent == nil && run.KeyHash != "" {
		spent, err = byWindow(rows, run.KeyHash, run.From, run.To)
	}
	if err != nil {
		return nil, fmt.Errorf("read the answer to GET %s/spend/logs: %w", l.root, err)
	}
	if spent == nil {
		return nil, fmt.Errorf("LiteLLM's spend log holds no row for the run's sessions or for its key from %s to %s; the proxy writes spend in batches, and the next measurement reads it again: %w", run.From.UTC().Format(time.RFC3339), run.To.UTC().Format(time.RFC3339), ErrNoSpend)
	}
	return spent, nil
}

// bySession sums the rows whose session_id or end_user is one of sessions
// and that started within from and to: the unified route sets both to
// Claude Code's session ID, the pass-through only end_user (S-0356), and a
// run that resumes a session keeps its ID, so the window keeps an earlier
// run's rows out. It is nil when no row is.
func bySession(rows []spendLogRow, sessions []string, from, to time.Time) (*Spent, error) {
	var spent *Spent
	for _, r := range rows {
		if !isSession(sessions, r.SessionID) && !isSession(sessions, r.EndUser) {
			continue
		}
		in, err := startedWithin(r.StartTime, from, to)
		if err != nil {
			return nil, err
		}
		if !in {
			continue
		}
		if spent == nil {
			spent = &Spent{By: SpentBySession}
		}
		spent.Cost += r.Spend
		spent.Calls++
	}
	return spent, nil
}

// startedWithin says whether a row that started at startTime started within from
// and to; a zero from or to leaves that side open.
func startedWithin(startTime string, from, to time.Time) (bool, error) {
	at, err := spendLogTime(startTime)
	if err != nil {
		return false, err
	}
	return (from.IsZero() || !at.Before(from)) && (to.IsZero() || !at.After(to)), nil
}

// isSession says whether id is one of sessions; no ID is none.
func isSession(sessions []string, id string) bool {
	return id != "" && slices.Contains(sessions, id)
}

// byWindow sums the rows filed under keyHash that started within from and
// to. It is nil when no row is.
func byWindow(rows []spendLogRow, keyHash string, from, to time.Time) (*Spent, error) {
	var spent *Spent
	for _, r := range rows {
		if r.APIKey != keyHash {
			continue
		}
		in, err := startedWithin(r.StartTime, from, to)
		if err != nil {
			return nil, err
		}
		if !in {
			continue
		}
		if spent == nil {
			spent = &Spent{By: SpentByWindow}
		}
		spent.Cost += r.Spend
		spent.Calls++
	}
	return spent, nil
}

// spendLogTime reads a row's startTime: RFC 3339, the same with a space for
// the T, or with no zone, which is UTC.
func spendLogTime(s string) (time.Time, error) {
	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.999999999", "2006-01-02 15:04:05.999999999Z07:00", "2006-01-02 15:04:05.999999999"} {
		if at, err := time.Parse(layout, s); err == nil {
			return at, nil
		}
	}
	return time.Time{}, fmt.Errorf("a row's startTime %q is in no layout the proxy writes; check the proxy's version against the one S-0356 checked", s)
}

// spendLog GETs the proxy's spend log rows, not its daily summary, for the
// days around the run's window, those of KeyHash's key only when it is set.
// An error names the endpoint and the status, never the key.
func (l *LiteLLM) spendLog(ctx context.Context, run RunRef) ([]spendLogRow, error) {
	endpoint := l.root + "/spend/logs"
	query := url.Values{
		// the proxy takes whole days from midnight UTC to midnight UTC, so
		// the window ends inside the day after To's
		"start_date": {run.From.UTC().Format(spendLogDate)},
		"end_date":   {run.To.UTC().AddDate(0, 0, 1).Format(spendLogDate)},
		"summarize":  {"false"},
	}
	if run.KeyHash != "" {
		query.Set("api_key", run.KeyHash)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+query.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("make the request GET %s: %w; check the provider's base_url", endpoint, err)
	}
	req.Header.Set("Authorization", "Bearer "+l.key)
	resp, err := l.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("read the spend from LiteLLM: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	switch {
	case resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden:
		return nil, fmt.Errorf("GET %s answered %d %s; only the proxy's master key reads the spend log: check the key the provider's spend_key_env names", endpoint, resp.StatusCode, http.StatusText(resp.StatusCode))
	case resp.StatusCode < 200 || resp.StatusCode > 299:
		return nil, fmt.Errorf("GET %s answered %d %s", endpoint, resp.StatusCode, http.StatusText(resp.StatusCode))
	case resp.Header.Get(spendLogTruncated) == "true":
		return nil, fmt.Errorf("GET %s answered only its most recent rows, so the run's may be missing; read it when fewer calls fall in the days around the run", endpoint)
	}
	body := &io.LimitedReader{R: resp.Body, N: maxSpendLogBody}
	rows, err := decodeRows(body)
	if err != nil && body.N == 0 {
		return nil, fmt.Errorf("GET %s answered more than %d MB of rows; read it when fewer calls fall in the days around the run", endpoint, maxSpendLogBody>>20)
	}
	if err != nil {
		return nil, fmt.Errorf("decode the answer to GET %s: %w; check that the provider's base_url is the proxy's root", endpoint, err)
	}
	return rows, nil
}

// decodeRows reads a JSON array of spend log rows one row at a time, so that
// the fields a join does not read are never held.
func decodeRows(r io.Reader) ([]spendLogRow, error) {
	dec := json.NewDecoder(r)
	tok, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if tok != json.Delim('[') {
		return nil, fmt.Errorf("want an array of rows, got %v", tok)
	}
	var rows []spendLogRow
	for dec.More() {
		var row spendLogRow
		if err := dec.Decode(&row); err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	if _, err := dec.Token(); err != nil {
		return nil, err
	}
	return rows, nil
}
