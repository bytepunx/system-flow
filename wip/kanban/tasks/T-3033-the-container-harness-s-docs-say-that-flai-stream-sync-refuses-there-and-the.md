---
id: T-3033
type: task
nature: remediation
title: The container harness's docs say that flai stream sync refuses there and the MCP tools task_done and verify sync on the host, and I-0272 is closed
status: backlog
parent: S-0689
owner: alex
created: 2026-11-10T10:33:41Z
updated: 2026-11-10T10:33:41Z
transitions: []
stream: S-0689
tags: [docs, container]
touches: [docs/users/flai.md, design/system/flai-cli.md, docs/operators/runbooks/agent-container.md, design/issues/I-0272-flai-stream-sync-run-in-a-story-s-agent-container-cannot-rebase-since-the-main.md, design/issues/summary.md]
after: [T-3026, T-3027, T-3030]
---
# T-3033 The container harness's docs say that flai stream sync refuses there and the MCP tools task_done and verify sync on the host, and I-0272 is closed

## Work

Write down what T-3026, T-3027, and T-3030 built, in the words they settled, which is why it waits for all three.

- `docs/users/flai.md`, § In a container (near lines 496-567): say that `flai stream sync` and the shell's `flai task done` refuse in a story's container, and that the MCP tools `task_done` and `verify` sync on the host. Remove the `stream_sync` MCP tool it names under "What it does not do yet": there is no such tool. Say in the `flai stream sync` section what a stop with no conflict and an unrun trial merge now report.
- `design/system/flai-cli.md`, the container harness section (near lines 713-754): the same refusal, beside the note that the main checkout's manifest is mounted so `flai verify` runs.
- `docs/operators/runbooks/agent-container.md` (near lines 77-88): beside "commits are the host's through `task_done` and the close-out", say that the shell's sync refuses there.

Then close the issue: `flai issue close I-0272 --reason "<what fixed it>"`, naming the refusal, the container prompt, and the stop and trial-merge reports, which regenerates `design/issues/summary.md`.

## Done when

- The three documents say what the container harness refuses and what syncs on the host, and none names a `stream_sync` MCP tool.
- I-0272 is closed with a reason naming what fixed it, and `design/issues/summary.md` no longer lists it.
- The markdown lint and `flai check --strict` pass.

## Notes
