package guard

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// The lines of a Claude Code 2.1.325 session transcript the guard reads, in
// the shapes S-0819 measured: a foreground Agent result that is only an
// error, one cut off on an error with the sub-agent's ID, a background
// launch's notice, and lines that are no tool call's result.
const (
	erroredResult = `{"type":"user","message":{"role":"user","content":[{"tool_use_id":"toolu_errored","type":"tool_result","content":"Error: Agent terminated early due to an API error: SSL certificate hostname mismatch","is_error":true}]},"toolUseResult":"Error: Agent terminated early due to an API error: SSL certificate hostname mismatch"}`
	cutOffResult  = `{"type":"user","message":{"role":"user","content":[{"tool_use_id":"toolu_cutoff","type":"tool_result","content":[{"type":"text","text":"terminated early"}],"is_error":null}]},"toolUseResult":{"status":"completed","agentId":"ae3487c9a0041b35d","agentType":"general-purpose"}}`
	asyncNotice   = `{"type":"user","message":{"role":"user","content":[{"tool_use_id":"toolu_async","type":"tool_result","content":"launched"}]},"toolUseResult":{"isAsync":true,"status":"async_launched","agentId":"acf8eff006f041f85"}}`
	aPrompt       = `{"type":"user","message":{"role":"user","content":"a prompt that mentions \"tool_result\" in passing"}}`
	aToolCall     = `{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","id":"toolu_pending","name":"Agent","input":{}}]}}`
	cutLine       = `{"type":"user","message":{"content":[{"tool_use_id":"toolu_pending","type":"tool_result"`
)

// writeTranscript writes lines as a session transcript in a temporary
// folder, and each of metas as the meta file of the sub-agent it is keyed
// by, beside it as Claude Code keeps them.
func writeTranscript(t *testing.T, lines []string, metas map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	transcript := filepath.Join(dir, "s1.jsonl")
	if err := os.WriteFile(transcript, []byte(strings.Join(lines, "\n")), 0o600); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "s1", "subagents")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	for id, meta := range metas {
		if err := os.WriteFile(filepath.Join(sub, "agent-"+id+".meta.json"), []byte(meta), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return transcript
}

// S-0819 (I-0376): a foreground sub-agent has ended once its launch has its
// result in the session transcript, an error string linked to it by its meta
// file or a result naming its ID; a background launch's notice, a launch
// with no result, and a sub-agent with no meta file and no result naming it
// have not, and a line that does not parse is passed over.
func TestEndedSubAgentsAreThoseWhoseForegroundLaunchHasItsResult(t *testing.T) {
	transcript := writeTranscript(t, []string{aPrompt, erroredResult, aToolCall, "not json \"tool_result\"", cutOffResult, asyncNotice, cutLine}, map[string]string{
		"ac8d9be303813c862": `{"agentType":"general-purpose","toolUseId":"toolu_errored","requestShape":"foreground","finished":true}`,
		"acf8eff006f041f85": `{"agentType":"general-purpose","toolUseId":"toolu_async","requestShape":"background"}`,
		"a0pending":         `{"agentType":"general-purpose","toolUseId":"toolu_pending","requestShape":"foreground"}`,
		"a0noshape":         `{"agentType":"general-purpose","toolUseId":"toolu_errored"}`,
		"a0unreadable":      `{"toolUseId":`,
	})
	ids := []string{"acf8eff006f041f85", "ae3487c9a0041b35d", "a0pending", "a0nometa", "ac8d9be303813c862", "a0noshape", "a0unreadable", "../a0"}
	got, err := EndedSubAgents(transcript, ids)
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"ae3487c9a0041b35d", "ac8d9be303813c862"}; !slices.Equal(got, want) {
		t.Errorf("ended %v, want %v", got, want)
	}
}

// S-0819: a line of a transcript runs to megabytes, past a scanner's
// default buffer, and its result is read all the same.
func TestEndedSubAgentsReadsALongLine(t *testing.T) {
	long := strings.Replace(cutOffResult, `"terminated early"`, `"`+strings.Repeat("x", 4<<20)+`"`, 1)
	got, err := EndedSubAgents(writeTranscript(t, []string{long}, nil), []string{"ae3487c9a0041b35d"})
	if err != nil || !slices.Equal(got, []string{"ae3487c9a0041b35d"}) {
		t.Errorf("ended %v, %v", got, err)
	}
}

// S-0819: a transcript that cannot be read is an error naming what was
// read for, and no sub-agent has ended.
func TestEndedSubAgentsFailsOnATranscriptItCannotRead(t *testing.T) {
	for what, path := range map[string]string{
		"missing":  filepath.Join(t.TempDir(), "s1.jsonl"),
		"a folder": t.TempDir(),
	} {
		got, err := EndedSubAgents(path, []string{"ae3487c9a0041b35d"})
		if err == nil || got != nil {
			t.Errorf("%s: ended %v, %v", what, got, err)
			continue
		}
		if !strings.Contains(err.Error(), "session transcript") {
			t.Errorf("%s: %v does not say what it read", what, err)
		}
	}
}
