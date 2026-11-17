package guard

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// toolResultMark is in every line of Claude Code's session transcript that
// carries a tool call's result, so that only those lines are parsed.
var toolResultMark = []byte(`"tool_result"`)

// transcriptLine is the part of a line of Claude Code's session transcript
// the guard reads: the message's content, which a tool call's result gives as
// a list of items, and toolUseResult, which for an Agent call's result is an
// object or, when the result is only an error, a string.
type transcriptLine struct {
	Message struct {
		Content json.RawMessage `json:"content"`
	} `json:"message"`
	ToolUseResult json.RawMessage `json:"toolUseResult"`
}

// contentItem is an item of a message's content: a tool_result names the
// tool call it is the result of.
type contentItem struct {
	Type      string `json:"type"`
	ToolUseID string `json:"tool_use_id"`
}

// agentResult is the part of an Agent call's toolUseResult the guard reads:
// the sub-agent it names, and whether it is a background launch's notice,
// written while the sub-agent goes on running.
type agentResult struct {
	AgentID string `json:"agentId"`
	IsAsync bool   `json:"isAsync"`
	Status  string `json:"status"`
}

// asyncLaunched is the status of a background launch's notice.
const asyncLaunched = "async_launched"

// subAgentMeta is the part of the file Claude Code keeps beside its session
// transcript for each sub-agent the guard reads: the tool call that launched
// it, and whether that call ran it in the foreground.
type subAgentMeta struct {
	ToolUseID    string `json:"toolUseId"`
	RequestShape string `json:"requestShape"`
}

// foreground is the requestShape of a sub-agent launched in the foreground.
const foreground = "foreground"

// results is what a session transcript holds of its tool calls' results:
// the sub-agents a result names that is not a background launch's notice,
// and the tool calls that have a result.
type results struct {
	agents map[string]bool
	uses   map[string]bool
}

// EndedSubAgents is those of ids, in their order, that Claude Code's session
// transcript at transcript shows have ended (ADR-0315): a tool call's result
// there names the sub-agent and is not a background launch's notice, or the
// sub-agent's meta file beside the transcript names a foreground launch whose
// call has its result there. A line that does not parse is passed over, and
// a meta file that cannot be read leaves only the first form. A transcript
// that cannot be read is an error, and then none has ended.
func EndedSubAgents(transcript string, ids []string) ([]string, error) {
	f, err := os.Open(transcript)
	if err != nil {
		return nil, fmt.Errorf("read the session transcript to find the sub-agents that ended: %w", err)
	}
	defer func() { _ = f.Close() }()
	res, err := scanResults(f)
	if err != nil {
		return nil, fmt.Errorf("read the session transcript %s to find the sub-agents that ended: %w", transcript, err)
	}
	dir := filepath.Join(strings.TrimSuffix(transcript, ".jsonl"), "subagents")
	var ended []string
	for _, id := range ids {
		if res.agents[id] {
			ended = append(ended, id)
			continue
		}
		if m, ok := readSubAgentMeta(dir, id); ok && m.RequestShape == foreground && m.ToolUseID != "" && res.uses[m.ToolUseID] {
			ended = append(ended, id)
		}
	}
	return ended, nil
}

// scanResults reads a session transcript line by line, parsing only the
// lines that carry a tool call's result. A line can run to megabytes, so it
// is read whole rather than with a scanner's bounded buffer.
func scanResults(r io.Reader) (results, error) {
	res := results{agents: map[string]bool{}, uses: map[string]bool{}}
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadBytes('\n')
		if bytes.Contains(line, toolResultMark) {
			res.add(line)
		}
		if errors.Is(err, io.EOF) {
			return res, nil
		}
		if err != nil {
			return res, err
		}
	}
}

// add notes the results line carries, if it parses and holds a tool_result.
func (res results) add(line []byte) {
	var l transcriptLine
	if json.Unmarshal(line, &l) != nil {
		return
	}
	var items []contentItem
	if json.Unmarshal(l.Message.Content, &items) != nil {
		return
	}
	found := false
	for _, it := range items {
		if it.Type == "tool_result" && it.ToolUseID != "" {
			res.uses[it.ToolUseID] = true
			found = true
		}
	}
	var a agentResult
	if !found || json.Unmarshal(l.ToolUseResult, &a) != nil {
		return
	}
	if a.AgentID != "" && !a.IsAsync && a.Status != asyncLaunched {
		res.agents[a.AgentID] = true
	}
}

// readSubAgentMeta reads the meta file of the sub-agent id in dir, and
// whether it could; an ID that is not a safe file name names none.
func readSubAgentMeta(dir, id string) (subAgentMeta, bool) {
	var m subAgentMeta
	if !sessionName.MatchString(id) {
		return m, false
	}
	data, err := os.ReadFile(filepath.Join(dir, "agent-"+id+".meta.json"))
	if err != nil || json.Unmarshal(data, &m) != nil {
		return m, false
	}
	return m, true
}
