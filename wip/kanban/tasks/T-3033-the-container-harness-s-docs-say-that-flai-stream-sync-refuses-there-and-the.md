---
id: T-3033
type: task
nature: remediation
title: The container harness's docs say that flai stream sync refuses there and the MCP tools task_done and verify sync on the host, and I-0272 is closed
status: done
parent: S-0689
owner: alex
created: 2026-11-10T10:33:41Z
updated: 2026-11-11T16:04:38Z
transitions:
  - to: ready
    at: 2026-11-11T16:01:05Z
    by: agent-S-0689
  - to: in-progress
    at: 2026-11-11T16:01:05Z
    by: agent-S-0689
  - to: done
    at: 2026-11-11T16:04:38Z
    by: agent-S-0689
stream: S-0689
tags: [docs, container]
touches: [docs/users/flai.md, design/system/flai-cli.md, design/system/agent-adapters.md, docs/operators/runbooks/agent-container.md, design/issues/I-0272-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md, design/issues/summary.md]
after: [T-3026, T-3027, T-3030]
usage:
  source: log
  seconds: 213
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 17043
      cache_read: 3009657
      cache_write: 155104
      cost: 1.8391
---
# T-3033 The container harness's docs say that flai stream sync refuses there and the MCP tools task_done and verify sync on the host, and I-0272 is closed

## Work

Write down what T-3026, T-3027, and T-3030 built, in the words they settled, which is why it waits for all three.

- `docs/users/flai.md`, § In a container (near lines 496-567): say that `flai stream sync` and the shell's `flai task done` refuse in a story's container, and that the MCP tools `task_done` and `verify` sync on the host. No `stream_sync` MCP tool exists: drop it where this document names one, under "What it does not do yet" (near line 564) and in the worktree paragraph (near line 863). Say in the `flai stream sync` section what a stop with no conflict and an unrun trial merge now report.
- `design/system/flai-cli.md`, the container harness section (near lines 713-754): the same refusal, beside the note that the main checkout's manifest is mounted so `flai verify` runs. Drop `stream_sync` from the tools the `flai stream` row names (near line 72); the `flai mcp` row (near line 142) already says neither server has one.
- `design/system/agent-adapters.md`: drop `stream_sync` from the MCP server row (near line 477) and the open question (near line 514), which name it among the tools that run on the host.
- `docs/operators/runbooks/agent-container.md` (near lines 77-88): beside "commits are the host's through `task_done` and the close-out", say that the shell's sync refuses there.

Then close the issue: `flai issue close I-0272 --reason "<what fixed it>"`, naming the refusal, the container prompt, and the stop and trial-merge reports, which regenerates `design/issues/summary.md`.

## Done when

- The four documents say what the container harness refuses and what syncs on the host, and none names a `stream_sync` MCP tool.
- I-0272 is closed with a reason naming what fixed it, and `design/issues/summary.md` no longer lists it.
- The markdown lint and `flai check --strict` pass.

## Notes
