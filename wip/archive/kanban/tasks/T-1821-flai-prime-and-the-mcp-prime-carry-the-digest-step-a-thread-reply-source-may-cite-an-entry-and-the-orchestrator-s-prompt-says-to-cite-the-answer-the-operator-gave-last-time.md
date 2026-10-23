---
id: T-1821
type: task
nature: feature
title: flai prime and the MCP prime carry the digest step, a thread_reply source may cite an entry, and the orchestrator's prompt says to cite the answer the operator gave last time
status: done
parent: S-0430
owner: alex
created: 2026-10-16T06:05:01Z
updated: 2026-10-23T19:43:16Z
transitions:
  - to: ready
    at: 2026-10-23T19:30:15Z
    by: agent-S-0430
  - to: in-progress
    at: 2026-10-23T19:30:16Z
    by: agent-S-0430
  - to: done
    at: 2026-10-23T19:43:15Z
    by: agent-S-0430
stream: S-0430
tags: [cli]
touches: [flai/cmd/prime.go, flai/cmd/prime_test.go, docs/users/flai-reference.md, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/prime_test.go, flai/internal/mcpserver/server_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/mcpserver/server.go]
after: [T-1820]
usage:
  source: log
  seconds: 779
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 63
      output: 24594
      cache_read: 4612137
      cache_write: 124154
      cost: 2.171
---
# T-1821 flai prime and the MCP prime carry the digest step, a thread_reply source may cite an entry, and the orchestrator's prompt says to cite the answer the operator gave last time

## Work

- `flai/cmd/prime_test.go` and `flai/internal/mcpserver/prime_test.go`: a story's pack and the orchestrator's carry digest items, in text, `--json`, and the MCP tool, and a part boundary falls cleanly before an entry.
- `flai/internal/mcpserver/server_test.go`: `thread_reply` with `recommendation` and a `source` naming an entry's path is accepted and written as the entry's last line; fix the source check if it is not.
- `flai/internal/harness/harness.go`: the orchestrator's prompt says that when an entry of its pack, or one `doc_search` finds, answers a thread's question, it recommends that answer and cites the entry as the `source`; `harness_test.go` checks the wording.

Waits for T-1820, so both packs have the step. Runs beside the docs task, which touches other files.

## Done when

The tests pass with `flai test flai/cmd flai/internal/mcpserver flai/internal/harness`.

## Notes

`.claude/agents/orchestrator.md` repeats the prompt in short. Changing it needs the operator's permission; leave it unless the story's agent finds it contradicts the prompt, and say so in the narrative.
