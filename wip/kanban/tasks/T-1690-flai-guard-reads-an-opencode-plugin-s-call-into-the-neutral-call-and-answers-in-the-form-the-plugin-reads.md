---
id: T-1690
type: task
nature: feature
title: flai guard reads an OpenCode plugin's call into the neutral Call and answers in the form the plugin reads
status: backlog
parent: S-0407
owner: alex
created: 2026-10-16T05:15:23Z
updated: 2026-10-16T05:25:08Z
transitions: []
stream: S-0407
tags: [cli]
touches: [flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go, flai/internal/guard/guard.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md]
---
# T-1690 flai guard reads an OpenCode plugin's call into the neutral Call and answers in the form the plugin reads

## Work

- Add an `opencode` reader in `flai/internal/guard/opencode.go`, registered in `readers` in `guard.go`. Define the input the plugin sends: the tool, its arguments, the session, and the sub-agent sign S-0406 found.
- Map `bash` to `Shell`, `edit`, `write`, and `patch` to `FileWrite`, and flai's MCP tools under their OpenCode names to `FlaiTool` without the prefix. Set `Subagent` from the sign, and `ReadOnly` for a sub-agent whose definition grants no edit.
- Record a sub-agent session's start and end where the plugin sees them, as `subagents.go` does for Claude Code, so ADR-0092's refusal holds.
- `Answer` exits 2 with the reason on stderr on a refusal and 0 otherwise.
- Test each mapping, a sub-agent's `flai` write refused, the agent's own let through, and `flai guard --harness opencode` end to end in `flai/cmd/guard_test.go`.
- Waits for nothing in this story: S-0406's finding defines the sign.

## Done when

`flai guard --harness opencode` refuses and allows the same calls `--harness claude-code` does for the same caller, and its tests pass.

## Notes
