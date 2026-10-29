package usage

import (
	"reflect"
	"testing"
	"time"
)

// The lines below are from flai serve's Claude Code logs, with long texts,
// signatures, and tool lists cut short.

const (
	ccSession = "02012299-7375-4739-9589-a2df3d786616"
	ccInit    = `{"type":"system","subtype":"init","cwd":"/home/alex/git/bytepunx/system-flow","session_id":"02012299-7375-4739-9589-a2df3d786616",` +
		`"tools":["Task","Bash","Edit","mcp__flai__inbox","mcp__flai__wait_for_events"],"mcp_servers":[{"name":"flai","status":"connected"}],` +
		`"model":"claude-opus-5-5","permissionMode":"default","uuid":"19434b5a-e8cb-4d3a-bd4a-e342e01d0222"}`
	// a message streamed as two lines, its thinking and then its tool_use
	ccThinking = `{"type":"assistant","message":{"model":"claude-opus-5-5","id":"msg_011CfoPY3H91E7G44isCcNkx","type":"message","role":"assistant",` +
		`"content":[{"type":"thinking","thinking":"","signature":"CAQS+wQKEAgSGAI4AUIIdGhpbmtpbmcS"}],"container":null,"stop_reason":null,` +
		`"usage":{"input_tokens":2,"cache_creation_input_tokens":10026,"cache_read_input_tokens":15445,"cache_creation":{"ephemeral_5m_input_tokens":0,` +
		`"ephemeral_1h_input_tokens":10026},"output_tokens":8,"service_tier":"standard","inference_geo":"not_available"}},"parent_tool_use_id":null,` +
		`"session_id":"02012299-7375-4739-9589-a2df3d786616","uuid":"0c4f6a3e-6f53-4a8e-9a43-2d43b1a2f001","timestamp":"2026-10-07T18:58:19.928Z",` +
		`"request_id":"req_011CfoPY2xo47kepzBTe4JEx","thinking_duration_ms":5724}`
	ccPrime = `{"type":"assistant","message":{"model":"claude-opus-5-5","id":"msg_011CfoPY3H91E7G44isCcNkx","type":"message","role":"assistant",` +
		`"content":[{"type":"tool_use","id":"toolu_01TrqWjm35rFmsKrVmgkGXCd","name":"mcp__flai__prime","input":{"role":"orchestrate"},"caller":{"type":"direct"}}],` +
		`"container":null,"stop_reason":null,"stop_sequence":null,"stop_details":null,"usage":{"input_tokens":2,"cache_creation_input_tokens":10026,` +
		`"cache_read_input_tokens":15445,"cache_creation":{"ephemeral_5m_input_tokens":0,"ephemeral_1h_input_tokens":10026},"output_tokens":8,` +
		`"service_tier":"standard","inference_geo":"not_available"},"input_transformations":[],"diagnostics":null,"context_management":null},` +
		`"parent_tool_use_id":null,"session_id":"02012299-7375-4739-9589-a2df3d786616","uuid":"ee705e55-6e6d-49ef-b273-b59f7636b283",` +
		`"timestamp":"2026-10-07T18:58:19.930Z","request_id":"req_011CfoPY2xo47kepzBTe4JEx","wire_tool_inputs":{"toolu_01TrqWjm35rFmsKrVmgkGXCd":` +
		`{"role":"orchestrate"}},"tool_use_meta":[{"id":"toolu_01TrqWjm35rFmsKrVmgkGXCd","display_name":"Prime","server_display_name":"flai"}]}`
	// the story's agent starts a sub-agent, which calls a tool and is answered
	ccAgent = `{"type":"assistant","message":{"model":"claude-opus-5-5","id":"msg_011CfoPaaywXTD6s7PmhtRWf","type":"message","role":"assistant",` +
		`"content":[{"type":"tool_use","id":"toolu_012KZUQK62GXpHqoSDii8udM","name":"Agent","input":{"subagent_type":"verifier",` +
		`"description":"Verify S-0298 diff vs criteria","run_in_background":false,"prompt":"Review story S-0298 for acceptance by the orchestrator."},` +
		`"caller":{"type":"direct"}}],"usage":{"input_tokens":2,"cache_creation_input_tokens":550,"cache_read_input_tokens":81119,` +
		`"cache_creation":{"ephemeral_5m_input_tokens":0,"ephemeral_1h_input_tokens":550},"output_tokens":8}},` +
		`"parent_tool_use_id":null,"session_id":"02012299-7375-4739-9589-a2df3d786616","uuid":"5619f16d-0719-4d0e-ba1e-e974f21d2b67",` +
		`"timestamp":"2026-10-07T18:58:55.926Z","request_id":"req_011CfoPaagq4ekr7a8hcAC5s"}`
	ccSubPrompt = `{"type":"user","message":{"role":"user","content":[{"type":"text","text":"Review story S-0298 for acceptance by the orchestrator."}]},` +
		`"parent_tool_use_id":"toolu_012KZUQK62GXpHqoSDii8udM","session_id":"02012299-7375-4739-9589-a2df3d786616",` +
		`"uuid":"e27aaaab-d3f7-4702-9b7f-c694a99625d1","timestamp":"2026-10-07T18:59:00.264Z","subagent_type":"verifier",` +
		`"task_description":"Verify S-0298 diff vs criteria"}`
	ccSubCall = `{"type":"assistant","message":{"model":"claude-sonnet-5-5","id":"msg_011CfoPbGqoZEF5TnqszRSGR","type":"message","role":"assistant",` +
		`"content":[{"type":"tool_use","id":"toolu_01QRAZiSJC7AjTHrgkmaD3ug","name":"mcp__flai__prime","input":{"story":"S-0298","role":"verify"},` +
		`"caller":{"type":"direct"}}],"container":null,"stop_reason":null,"stop_sequence":null,"stop_details":null,"usage":{"input_tokens":2,` +
		`"cache_creation_input_tokens":11274,"cache_read_input_tokens":0,"cache_creation":{"ephemeral_5m_input_tokens":11274,"ephemeral_1h_input_tokens":0},` +
		`"output_tokens":24,"service_tier":"standard","inference_geo":"not_available"},"input_transformations":[],"diagnostics":null,"context_management":null},` +
		`"parent_tool_use_id":"toolu_012KZUQK62GXpHqoSDii8udM","session_id":"02012299-7375-4739-9589-a2df3d786616","uuid":"aeabcba0-30fa-4ec2-ad98-3d97da4e2b4e",` +
		`"timestamp":"2026-10-07T18:59:03.202Z","request_id":"req_011CfoPbGTF1QXKp1YHhTPHt","subagent_type":"verifier","task_description":"Verify S-0298 diff vs criteria",` +
		`"tool_use_meta":[{"id":"toolu_01QRAZiSJC7AjTHrgkmaD3ug","display_name":"Prime","server_display_name":"flai"}]}`
	ccSubResult = `{"type":"user","message":{"role":"user","content":[{"tool_use_id":"toolu_01QRAZiSJC7AjTHrgkmaD3ug","type":"tool_result",` +
		`"content":"{\"briefs_left_out\":131,\"budget\":40960}"}]},"parent_tool_use_id":"toolu_012KZUQK62GXpHqoSDii8udM",` +
		`"session_id":"02012299-7375-4739-9589-a2df3d786616","uuid":"6810eb37-78bc-4097-b3b9-a4783082cbe9","timestamp":"2026-10-07T18:59:03.385Z",` +
		`"tool_use_result":{"_meta":{"io.modelcontextprotocol/serverInfo":{"name":"flai","title":"system-flow repository","version":"1.35.1"}}},` +
		`"subagent_type":"verifier","task_description":"Verify S-0298 diff vs criteria"}`
	// the story's agent waits, and the wait reports a story moved
	ccWait = `{"type":"assistant","message":{"model":"claude-opus-5-5","id":"msg_011CfoPhQijdyTEkMx579emh","type":"message","role":"assistant",` +
		`"content":[{"type":"tool_use","id":"toolu_012JT6msan4fbTsxqYMzFqnr","name":"mcp__flai__wait_for_events","input":{"timeout_seconds":1800},` +
		`"caller":{"type":"direct"}}],"usage":{"input_tokens":2,"cache_creation_input_tokens":810,"cache_read_input_tokens":99837,` +
		`"cache_creation":{"ephemeral_5m_input_tokens":0,"ephemeral_1h_input_tokens":810},"output_tokens":24}},` +
		`"parent_tool_use_id":null,"session_id":"02012299-7375-4739-9589-a2df3d786616","uuid":"55e53e97-062e-4ea3-a168-3e9f2cec7b5e",` +
		`"timestamp":"2026-10-07T19:00:28.283Z","request_id":"req_011CfoPhQSsSreL1Eu6AUaWv"}`
	ccWoken    = `{"changed":[],"end":false,"events":[{"at":"2026-10-07T18:59:43Z","by":"alex","id":"S-0298","kind":"moved","to":"done","type":"story"}],"events_omitted":0,"timed_out":false}`
	ccWokeLine = `{"type":"user","message":{"role":"user","content":[{"tool_use_id":"toolu_012JT6msan4fbTsxqYMzFqnr","type":"tool_result",` +
		`"content":"{\"changed\":[],\"end\":false,\"events\":[{\"at\":\"2026-10-07T18:59:43Z\",\"by\":\"alex\",\"id\":\"S-0298\",\"kind\":\"moved\",` +
		`\"to\":\"done\",\"type\":\"story\"}],\"events_omitted\":0,\"timed_out\":false}"}]},"parent_tool_use_id":null,` +
		`"session_id":"02012299-7375-4739-9589-a2df3d786616","uuid":"3207b801-09fd-470b-b153-e93b78bddb41","timestamp":"2026-10-07T19:00:28.636Z",` +
		`"tool_use_result":{"content":"{\"changed\":[]}"}}`
	// the run's end: what Claude Code reports it cost, per model
	ccResult = `{"type":"result","subtype":"success","is_error":false,"duration_api_ms":2254675,"num_turns":59,"stop_reason":"end_turn",` +
		`"session_id":"02012299-7375-4739-9589-a2df3d786616","total_cost_usd":17.6853053,"usage":{"input_tokens":130,"cache_creation_input_tokens":14829,` +
		`"cache_read_input_tokens":1038913,"output_tokens":1929,"output_tokens_details":{"thinking_tokens":372},"service_tier":"standard"},` +
		`"modelUsage":{"claude-opus-5-5":{"inputTokens":1631,"outputTokens":98106,"thinkingTokens":13611,"cacheReadInputTokens":6449243,` +
		`"cacheCreationInputTokens":605181,"webSearchRequests":0,"costUSD":15.62246575,"contextWindow":1000000,"maxOutputTokens":64000,` +
		`"canonicalModel":"claude-opus-5-5","provider":"firstParty","costBasis":"list"},"claude-haiku-4-5-20251001":{"inputTokens":853006,` +
		`"outputTokens":55817,"cacheReadInputTokens":3654023,"cacheCreationInputTokens":260277,"webSearchRequests":24,"costUSD":2.06283955}},` +
		`"permission_denials":[],"terminal_reason":"completed","result":"S-0298 is accepted.\n\nNext: S-0280."}`
	// a run that ends on the usage limit (I-0212): the rate_limit_event that
	// rejects it, and its result
	ccLimitSession = "1441b203-a736-430d-a0da-597f87cc4f63"
	ccLimitText    = "You've hit your weekly limit · resets Oct 28, 2am (America/Denver)"
	ccLimitEvent   = `{"type":"rate_limit_event","rate_limit_info":{"status":"rejected","resetsAt":1793174400,"rateLimitType":"seven_day",` +
		`"overageStatus":"rejected","overageDisabledReason":"org_level_disabled","isUsingOverage":false,"unifiedWindows":{"five_hour":` +
		`{"utilization":0.21,"resetsAt":1793033400,"status":"allowed"},"seven_day":{"utilization":1,"resetsAt":1793174400,"status":"rejected",` +
		`"surpassedThreshold":1}}},"uuid":"2132f9f6-9747-4763-94b0-eb399bbc28fd","session_id":"1441b203-a736-430d-a0da-597f87cc4f63"}`
	ccLimitResult = `{"duration_api_ms":0,"stop_reason":"stop_sequence","session_id":"1441b203-a736-430d-a0da-597f87cc4f63","total_cost_usd":0,` +
		`"usage":{"output_tokens_details":{"thinking_tokens":0},"input_tokens":0,"cache_creation_input_tokens":0,"cache_read_input_tokens":0,` +
		`"output_tokens":0,"service_tier":"standard","inference_geo":"","iterations":[],"speed":"standard","fallback_credit":null},"modelUsage":{},` +
		`"permission_denials":[],"terminal_reason":"api_error","fast_mode_state":"off","safety_stops":0,"is_error":true,"num_turns":1,` +
		`"subtype":"success","api_error_status":429,"api_error_origin":"first_party","api_error":"usage_limit_reached",` +
		`"result":"You've hit your weekly limit · resets Oct 28, 2am (America/Denver)","type":"result","duration_ms":412,` +
		`"uuid":"08775fa3-fee4-4409-9d86-9783aa52ef89","queued_turn_count":0,"result_index":0,"model_usage_since_last_result":{}}`
	// a rate_limit_event that lets the run go on
	ccAllowedEvent = `{"type":"rate_limit_event","rate_limit_info":{"status":"allowed_warning","resetsAt":1793174400,"rateLimitType":"seven_day",` +
		`"utilization":0.9},"session_id":"1441b203-a736-430d-a0da-597f87cc4f63"}`
)

func ccAt(t *testing.T, s string) time.Time {
	t.Helper()
	at, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatal(err)
	}
	return at
}

// readLines reads lines with one reader, in order, and is the events they
// hold; it fails when a line is not read as a record.
func readLines(t *testing.T, r Reader, lines ...string) []Event {
	t.Helper()
	var out []Event
	for _, l := range lines {
		events, ok := r.Read([]byte(l + "\n"))
		if !ok {
			t.Fatalf("line not read as a record: %s", l)
		}
		out = append(out, events...)
	}
	return out
}

func checkEvents(t *testing.T, got, want []Event) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d events, want %d:\n%+v", len(got), len(want), got)
	}
	for i := range want {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Errorf("event %d =\n%+v\nwant\n%+v", i, got[i], want[i])
		}
	}
}

func TestClaudeCodeReadsASessionStart(t *testing.T) {
	checkEvents(t, readLines(t, ClaudeCode(), ccInit), []Event{{Kind: SessionStart, Session: ccSession, Model: "claude-opus-5-5"}})
}

// Claude Code streams a message once per content block, each line with the
// message's ID and a usage that need not be its final count: each line is a
// Call as it is, the repeats left to the measurement.
func TestClaudeCodeReadsEachLineOfARepeatedMessageAsItIs(t *testing.T) {
	tokens := &Tokens{Input: 2, Output: 8, CacheRead: 15445, CacheWrite: 10026}
	call := Event{Kind: Call, Session: ccSession, Model: "claude-opus-5-5", MessageID: "msg_011CfoPY3H91E7G44isCcNkx", Tokens: tokens}
	thinking, prime := call, call
	thinking.At, prime.At = ccAt(t, "2026-10-07T18:58:19.928Z"), ccAt(t, "2026-10-07T18:58:19.930Z")
	use := prime
	use.Kind, use.Tokens, use.Tool, use.ToolID = ToolCall, nil, "mcp__flai__prime", "toolu_01TrqWjm35rFmsKrVmgkGXCd"
	checkEvents(t, readLines(t, ClaudeCode(), ccThinking, ccPrime), []Event{thinking, prime, use})

	// through OpenRouter the IDs are the upstream's, and the count grows
	// over the repeats
	bedrock := func(at string, out int, block string) string {
		return `{"type":"assistant","message":{"model":"claude-opus-5-5","id":"msg_bdrk_01Xq7cM2Jt3o4Fv9GQkS1aLp","type":"message","role":"assistant",` +
			`"content":[` + block + `],"usage":{"input_tokens":3,"cache_creation_input_tokens":548,"cache_read_input_tokens":260197,"output_tokens":` +
			itoa(out) + `}},"parent_tool_use_id":null,"session_id":"88da66f1-3cc4-4952-abab-8db192d5e4bb","timestamp":"` + at + `"}`
	}
	got := readLines(t, ClaudeCode(),
		bedrock("2026-10-12T20:09:01.100Z", 1, `{"type":"text","text":"Reading the story."}`),
		bedrock("2026-10-12T20:09:01.900Z", 97, `{"type":"tool_use","id":"toolu_bdrk_01Ff2","name":"mcp__flai__inbox","input":{"story":"S-0351"}}`))
	if len(got) != 3 || got[0].MessageID != "msg_bdrk_01Xq7cM2Jt3o4Fv9GQkS1aLp" || got[1].MessageID != got[0].MessageID ||
		got[0].Tokens.Output != 1 || got[1].Tokens.Output != 97 || got[2].Kind != ToolCall || got[2].ToolID != "toolu_bdrk_01Ff2" {
		t.Errorf("OpenRouter repeats = %+v, want two calls of one ID with 1 and 97 output tokens, then the inbox call", got)
	}
}

// A sub-agent's start is the Agent call's description and prompt; its lines
// carry the ID of that call and its description.
func TestClaudeCodeReadsASubAgentsStartAndItsCalls(t *testing.T) {
	const agentID, described = "toolu_012KZUQK62GXpHqoSDii8udM", "Verify S-0298 diff vs criteria"
	start := Event{Kind: Call, At: ccAt(t, "2026-10-07T18:58:55.926Z"), Session: ccSession, Model: "claude-opus-5-5", MessageID: "msg_011CfoPaaywXTD6s7PmhtRWf",
		Tokens: &Tokens{Input: 2, Output: 8, CacheRead: 81119, CacheWrite: 550}}
	agent := start
	agent.Kind, agent.Tokens, agent.Tool, agent.ToolID = ToolCall, nil, "Agent", agentID
	agent.Description, agent.Prompt = described, "Review story S-0298 for acceptance by the orchestrator."
	sub := Event{Kind: Call, At: ccAt(t, "2026-10-07T18:59:03.202Z"), Session: ccSession, Model: "claude-sonnet-5-5", MessageID: "msg_011CfoPbGqoZEF5TnqszRSGR",
		Tokens: &Tokens{Input: 2, Output: 24, CacheWrite: 11274}, ParentCall: agentID, ParentDescription: described}
	subUse := sub
	subUse.Kind, subUse.Tokens, subUse.Tool, subUse.ToolID = ToolCall, nil, "mcp__flai__prime", "toolu_01QRAZiSJC7AjTHrgkmaD3ug"
	answer := Event{Kind: ToolResult, At: ccAt(t, "2026-10-07T18:59:03.385Z"), Session: ccSession, ParentCall: agentID, ParentDescription: described,
		ToolID: "toolu_01QRAZiSJC7AjTHrgkmaD3ug"}
	// the sub-agent's prompt is a record that holds no event
	checkEvents(t, readLines(t, ClaudeCode(), ccAgent, ccSubPrompt, ccSubCall, ccSubResult), []Event{start, agent, sub, subUse, answer})
}

// A Bash call's command and the file an Edit or Write call changes are read,
// each tool_use of a line in its order.
func TestClaudeCodeReadsEachToolCallsCommandAndPath(t *testing.T) {
	line := `{"type":"assistant","timestamp":"2026-09-24T08:30:12.511Z","session_id":"s","parent_tool_use_id":null,"message":{"id":"msg_01","model":"claude-opus-5-5",` +
		`"content":[{"type":"tool_use","id":"toolu_0144q7YY5NtL2R9RDmBXaXzw","name":"Bash","input":{"command":"cat wip/agents/index.md && git status --short",` +
		`"description":"Read the index"}},{"type":"text","text":"Then the narrative."},{"type":"tool_use","id":"toolu_01Hk6fb6DiQe7qgfTijrZcuw","name":"Edit",` +
		`"input":{"replace_all":false,"file_path":"/home/alex/git/bytepunx/system-flow/wip/agents/S-0016.md","old_string":"a","new_string":"b"}}]}}`
	got := readLines(t, ClaudeCode(), line)
	if len(got) != 3 || got[0].Kind != Call || got[0].Tokens != nil {
		t.Fatalf("events = %+v, want a call with no usage reported and two tool calls", got)
	}
	bash, edit := got[1], got[2]
	if bash.Tool != "Bash" || bash.Command != "cat wip/agents/index.md && git status --short" || bash.Description != "Read the index" || bash.MessageID != "msg_01" {
		t.Errorf("bash = %+v", bash)
	}
	if edit.Tool != "Edit" || edit.ToolID != "toolu_01Hk6fb6DiQe7qgfTijrZcuw" || edit.Path != "/home/alex/git/bytepunx/system-flow/wip/agents/S-0016.md" || edit.Command != "" {
		t.Errorf("edit = %+v", edit)
	}
}

// A wait_for_events call's result is read as text, whether its content is a
// string or text blocks, failed or not, and whoever called; another tool's
// result is not.
func TestClaudeCodeReadsTheResultOfAWait(t *testing.T) {
	r := ClaudeCode()
	readLines(t, r, ccWait,
		waitCall("m1", "2026-10-07T19:01:00Z", "w1", ""), waitCall("m2", "2026-10-07T19:02:00Z", "w2", ""),
		waitCall("m3", "2026-10-07T19:03:00Z", "w3", "toolu_A"), waitCall("m4", "2026-10-07T19:04:00Z", "w4", ""),
		opusCall("m5", "2026-10-07T19:05:00Z", 100, "", "", toolUse("o1", "mcp__flai__wait_for_work", `{}`)))
	both := `{"type":"user","timestamp":"2026-10-07T19:06:00Z","session_id":"s","parent_tool_use_id":null,"message":{"role":"user","content":[` +
		`{"tool_use_id":"o1","type":"tool_result","content":` + quoted(timedOut) + `},` +
		`{"tool_use_id":"w4","type":"tool_result","content":[{"type":"text","text":"first"},{"type":"image","source":{}},{"type":"text","text":"second"}]}]}}`
	cases := []struct {
		line, id, result string
		isError          bool
	}{
		{ccWokeLine, "toolu_012JT6msan4fbTsxqYMzFqnr", ccWoken, false},
		{toolResult("2026-10-07T19:01:30Z", "w1", quoted(timedOut), false), "w1", timedOut, false},
		{toolResult("2026-10-07T19:02:30Z", "w2", `[{"type":"text","text":`+quoted(timedOut)+`}]`, false), "w2", timedOut, false},
		{toolResult("2026-10-07T19:03:30Z", "w3", quoted(timedOut), false), "w3", timedOut, false},
		{toolResult("2026-10-07T19:04:30Z", "w2", quoted("Connection closed"), true), "w2", "Connection closed", true},
	}
	for _, c := range cases {
		got := readLines(t, r, c.line)
		if len(got) != 1 || got[0].Kind != ToolResult || got[0].ToolID != c.id || got[0].Result != c.result || got[0].IsError != c.isError {
			t.Errorf("result of %s = %+v, want %q with is_error %v", c.id, got, c.result, c.isError)
		}
	}
	got := readLines(t, r, both)
	if len(got) != 2 || got[0].ToolID != "o1" || got[0].Result != "" || got[1].ToolID != "w4" || got[1].Result != "first\nsecond" {
		t.Errorf("two results in a line = %+v, want wait_for_work's unread and the wait's text blocks joined", got)
	}
	// a reader of its own has seen no wait
	if fresh := readLines(t, ClaudeCode(), ccWokeLine); len(fresh) != 1 || fresh[0].Result != "" {
		t.Errorf("a result whose call was not read = %+v, want its text unread", fresh)
	}
}

// A result is an End: the run's final text, total cost, usage, and each
// model's totals and cost; a result that reports none of them is an End
// still, with nothing.
func TestClaudeCodeReadsARunsEnd(t *testing.T) {
	cost := 17.6853053
	want := Event{Kind: End, Session: ccSession, CostUSD: &cost, Tokens: &Tokens{Input: 130, Output: 1929, CacheRead: 1038913, CacheWrite: 14829},
		Models: []Model{
			{Model: "claude-haiku-4-5-20251001", Input: 853006, Output: 55817, CacheRead: 3654023, CacheWrite: 260277, Cost: 2.06283955},
			{Model: "claude-opus-5-5", Input: 1631, Output: 98106, CacheRead: 6449243, CacheWrite: 605181, Cost: 15.62246575},
		},
		Result: "S-0298 is accepted.\n\nNext: S-0280."}
	checkEvents(t, readLines(t, ClaudeCode(), ccResult), []Event{want})
	checkEvents(t, readLines(t, ClaudeCode(), `{"type":"result","subtype":"error_max_turns","session_id":"s"}`), []Event{{Kind: End, Session: "s"}})
}

// A line that is not a record, or a record that holds no event, is passed
// over unread; a record of a kind that holds events may hold none.
func TestClaudeCodePassesOverLinesThatHoldNoEvent(t *testing.T) {
	r := ClaudeCode()
	for _, l := range []string{
		"",
		`not json: a line the harness wrote to stderr`,
		`{"type":"assistant","message":`,
		`{"type":"system","subtype":"thinking_tokens","estimated_tokens":50,"estimated_tokens_delta":50,"session_id":"s"}`,
		`{"type":"system","subtype":"session_title_changed","title":"sf S-0339","session_id":"s"}`,
	} {
		if events, ok := r.Read([]byte(l)); ok || events != nil {
			t.Errorf("Read(%q) = %+v, %v; want nothing, not a record read", l, events, ok)
		}
	}
	for _, l := range []string{ccSubPrompt, ccLimitEvent, ccAllowedEvent} {
		if events, ok := r.Read([]byte(l)); !ok || len(events) != 0 {
			t.Errorf("Read(%.60q) = %+v, %v; want read, no event", l, events, ok)
		}
	}
}

// limitEnd is the End of I-0212's result, on the usage limit with the reset
// time given.
func limitEnd(resets time.Time) Event {
	cost := 0.0
	return Event{Kind: End, Session: ccLimitSession, CostUSD: &cost, Tokens: &Tokens{}, Result: ccLimitText,
		Limit: &Limit{Error: "usage_limit_reached", Status: 429, Resets: resets}}
}

// A run that ends on the usage limit, as I-0212's log ends: the
// rate_limit_event that rejects it gives the End's reset time, and the result
// its error and status.
func TestClaudeCodeReadsARunsEndOnTheUsageLimit(t *testing.T) {
	resets := time.Date(2026, 10, 28, 8, 0, 0, 0, time.UTC)
	checkEvents(t, readLines(t, ClaudeCode(), ccLimitEvent, ccLimitResult), []Event{limitEnd(resets)})
}

// A result that is no error ends on no limit, even after a rate_limit_event
// rejected the run.
func TestClaudeCodeReadsNoLimitOnAResultThatIsNoError(t *testing.T) {
	got := readLines(t, ClaudeCode(), ccLimitEvent, ccResult)
	if len(got) != 1 || got[0].Limit != nil {
		t.Errorf("events = %+v, want one End on no limit", got)
	}
}

// A result on the limit with no rejecting rate_limit_event before it, in the
// run or since the run's last End, carries its error with no reset time: as
// does one after an event that allows the run.
func TestClaudeCodeReadsNoResetTimeWithoutARejectingRateLimitEvent(t *testing.T) {
	for name, lines := range map[string][]string{
		"no event":                {ccLimitResult},
		"an allowing event":       {ccAllowedEvent, ccLimitResult},
		"allowed after rejected":  {ccLimitEvent, ccAllowedEvent, ccLimitResult},
		"rejected in a past run":  {ccLimitEvent, ccInit, ccLimitResult},
		"rejected before the End": {ccLimitEvent, ccLimitResult, ccLimitResult},
	} {
		got := readLines(t, ClaudeCode(), lines...)
		if want := limitEnd(time.Time{}); !reflect.DeepEqual(got[len(got)-1], want) {
			t.Errorf("%s: End = %+v %+v, want %+v %+v", name, got[len(got)-1], got[len(got)-1].Limit, want, want.Limit)
		}
	}
}

// A result that is an error on no limit, one on another status that no
// rate_limit_event rejected, ends on none.
func TestClaudeCodeReadsNoLimitOnAnotherError(t *testing.T) {
	for _, l := range []string{
		`{"type":"result","subtype":"success","is_error":true,"api_error_status":529,"api_error":"overloaded_error","session_id":"s"}`,
		`{"type":"result","subtype":"error_max_turns","is_error":true,"session_id":"s"}`,
	} {
		if got := readLines(t, ClaudeCode(), ccAllowedEvent, l); len(got) != 1 || got[0].Limit != nil {
			t.Errorf("%s: events = %+v, want one End on no limit", l, got)
		}
	}
}
