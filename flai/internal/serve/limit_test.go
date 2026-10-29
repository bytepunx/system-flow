package serve

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bytepunx/system-flow/flai/internal/usage"
)

// limitStream is I-0212's log, a curator run that ends on the harness's
// usage limit, with its init and its synthetic message cut short: the
// rate_limit_event that rejects the run, resetting at 1793174400, and the
// result on usage_limit_reached with status 429.
var limitStream = []string{
	`{"type":"system","subtype":"session_title_changed","title":"sf curate","uuid":"c64ec637-da4a-43ca-9a5f-c5c32cb2fe12",` +
		`"session_id":"1441b203-a736-430d-a0da-597f87cc4f63"}`,
	`{"type":"system","subtype":"init","cwd":"/home/alex/git/bytepunx/system-flow","session_id":"1441b203-a736-430d-a0da-597f87cc4f63",` +
		`"tools":["Read","Bash","Task","mcp__flai__prime"],"model":"claude-opus-5-5","permissionMode":"acceptEdits","claude_code_version":"2.1.310"}`,
	`{"type":"rate_limit_event","rate_limit_info":{"status":"rejected","resetsAt":1793174400,"rateLimitType":"seven_day",` +
		`"overageStatus":"rejected","overageDisabledReason":"org_level_disabled","isUsingOverage":false,"unifiedWindows":{"five_hour":` +
		`{"utilization":0.21,"resetsAt":1793033400,"status":"allowed"},"seven_day":{"utilization":1,"resetsAt":1793174400,"status":"rejected",` +
		`"surpassedThreshold":1}}},"uuid":"2132f9f6-9747-4763-94b0-eb399bbc28fd","session_id":"1441b203-a736-430d-a0da-597f87cc4f63"}`,
	`{"type":"assistant","message":{"id":"0524a12d-830e-40ad-a6b5-d047f3d66543","model":"<synthetic>","role":"assistant",` +
		`"stop_reason":"stop_sequence","type":"message","usage":{"input_tokens":0,"output_tokens":0,"cache_creation_input_tokens":0,` +
		`"cache_read_input_tokens":0},"content":[{"type":"text","text":"You've hit your weekly limit · resets Oct 28, 2am (America/Denver)"}]},` +
		`"parent_tool_use_id":null,"session_id":"1441b203-a736-430d-a0da-597f87cc4f63","uuid":"131b6d79-7075-4ad8-a033-867a47b53203",` +
		`"timestamp":"2026-10-26T16:03:54.269Z","error":"rate_limit","is_api_error_message":true,"api_error_status":429,` +
		`"api_error":"usage_limit_reached","api_error_origin":"first_party"}`,
	`{"duration_api_ms":0,"stop_reason":"stop_sequence","session_id":"1441b203-a736-430d-a0da-597f87cc4f63","total_cost_usd":0,` +
		`"usage":{"output_tokens_details":{"thinking_tokens":0},"input_tokens":0,"cache_creation_input_tokens":0,"cache_read_input_tokens":0,` +
		`"output_tokens":0,"service_tier":"standard","inference_geo":"","iterations":[],"speed":"standard","fallback_credit":null},"modelUsage":{},` +
		`"permission_denials":[],"terminal_reason":"api_error","fast_mode_state":"off","safety_stops":0,"is_error":true,"num_turns":1,` +
		`"subtype":"success","api_error_status":429,"api_error_origin":"first_party","api_error":"usage_limit_reached",` +
		`"result":"You've hit your weekly limit · resets Oct 28, 2am (America/Denver)","type":"result","duration_ms":412,` +
		`"uuid":"08775fa3-fee4-4409-9d86-9783aa52ef89","queued_turn_count":0,"result_index":0,"model_usage_since_last_result":{}}`,
}

// limitResetsAt is when I-0212's limit resets, limitWhyText the why of a run
// that ended on it, and limitSaid the harness's text at its end.
const (
	limitResetsAt = "2026-10-28T08:00:00Z"
	limitWhyText  = "ended on the harness's usage limit (usage_limit_reached, 429); it resets at " + limitResetsAt
	limitSaid     = "You've hit your weekly limit · resets Oct 28, 2am (America/Denver)"
)

// limitLog writes lines as a run's log, one to a line, and returns its path.
func limitLog(t *testing.T, lines ...string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "t-curator-20261026T160352Z.log")
	if err := os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

// warned is a warn that keeps what it was told.
type warned []string

func (w *warned) warn(msg string, args ...any) {
	*w = append(*w, fmt.Sprint(append([]any{msg}, args...)...))
}

// A run whose log ends on the usage limit has it, with the error, the
// status, and the reset time the harness named, and the harness's text.
func TestARunsLimitIsReadFromItsLog(t *testing.T) {
	var w warned
	lim, said := runLimit(&AgentRun{Log: limitLog(t, limitStream...)}, w.warn)
	want := usage.Limit{Error: "usage_limit_reached", Status: 429, Resets: time.Date(2026, 10, 28, 8, 0, 0, 0, time.UTC)}
	if lim == nil || *lim != want {
		t.Fatalf("limit = %+v, want %+v", lim, want)
	}
	if got := limitWhy(lim); got != limitWhyText {
		t.Errorf("why = %q, want %q", got, limitWhyText)
	}
	if got := limitResets(lim); got != limitResetsAt {
		t.Errorf("resets = %q, want %q", got, limitResetsAt)
	}
	if said != limitSaid {
		t.Errorf("said = %q, want %q", said, limitSaid)
	}
	if len(w) != 0 {
		t.Errorf("warned: %v", w)
	}
}

// A run whose log ends otherwise, whose log is missing, or that names none,
// has no limit, and nothing is warned of; a log that cannot be read is
// warned of, and read as ending on none.
func TestARunWithoutALimitHasNone(t *testing.T) {
	worked := limitLog(t, streamCall("s", "m1", runStart, 999), streamFinal("Curated.", 2000, 0.5))
	for name, log := range map[string]string{"worked": worked, "missing": filepath.Join(t.TempDir(), "gone.log"), "none named": ""} {
		var w warned
		if lim, said := runLimit(&AgentRun{Log: log}, w.warn); lim != nil || said != "" || len(w) != 0 {
			t.Errorf("%s: limit = %+v, warned %v; want none, nothing warned", name, lim, w)
		}
	}
	var w warned
	if lim, _ := runLimit(&AgentRun{Log: t.TempDir()}, w.warn); lim != nil || len(w) != 1 || !strings.Contains(w[0], "run's log not read for a usage limit") {
		t.Errorf("an unreadable log: limit = %+v, warned %v; want none, warned once", lim, w)
	}
}

// The why leaves out what the harness did not name, and says when it named
// no reset time.
func TestALimitsWhyNamesWhatTheHarnessNamed(t *testing.T) {
	for _, c := range []struct {
		lim  usage.Limit
		want string
	}{
		{usage.Limit{Error: "usage_limit_reached", Status: 429}, "ended on the harness's usage limit (usage_limit_reached, 429); it named no reset time"},
		{usage.Limit{Status: 429, Resets: time.Date(2026, 10, 28, 2, 0, 0, 0, time.FixedZone("MDT", -6*3600))}, "ended on the harness's usage limit (429); it resets at " + limitResetsAt},
		{usage.Limit{}, "ended on the harness's usage limit; it named no reset time"},
	} {
		if got := limitWhy(&c.lim); got != c.want {
			t.Errorf("why of %+v = %q, want %q", c.lim, got, c.want)
		}
	}
}
