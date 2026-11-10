---
id: TH-1002
title: "Plan for S-0668: stream_state and flai stream state write Context, Decisions, and Open questions; five tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0668-a-story-s-agent-in-a-container-cannot-write-its-narrative-s-decisions-context.md
  item: S-0668
status: open
participants: [system-flow]
created: 2026-11-10T08:51:46Z
updated: 2026-11-10T08:51:46Z
---

# TH-1002 Plan for S-0668: stream_state and flai stream state write Context, Decisions, and Open questions; five tasks in three layers

On wip/kanban/stories/S-0668-a-story-s-agent-in-a-container-cannot-write-its-narrative-s-decisions-context.md.

## Entries

### 2026-11-10T08:51:46Z system-flow
No thread named "asked" exists on S-0668, and the story had no earlier plan, tasks, or Planning notes. So this is a first plan, not a replan.

## Proposed fix

The MCP server runs on the host, and a container's agent reaches it over HTTP. So `stream_state` already writes the main checkout's `wip/` from inside a container. I propose widening what it writes, and not mounting `wip/` read-write, which ADR-0253 refused.

- `context` replaces `## Context`.
- `decisions` appends one bullet per text to `## Decisions`.
- `questions` appends one dated bullet per text to `## Open questions`, above the threads block, which stays as it is.

The same three go on `flai stream state` (`--context`, `--decision`, `--question`) and on the host channel's `stream.state`.

## Tasks

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-2944 `SetStreamState` writes the three sections | none | `flai/internal/workitem/narrative.go`, its test |
| 2 | T-2947 CLI flags, reference, and flag index | T-2944 | `flai/cmd/stream.go`, `flai/cmd/stream_state_test.go`, `docs/users/flai-reference.md`, `docs/operators/settings.md` |
| 2 | T-2949 MCP fields, and the test that reproduces I-0277 | T-2944 | `flai/internal/mcpserver/stream.go`, its test |
| 3 | T-2950 host channel `stream.state` | T-2947 | `flai/internal/hostapi/writes_streams.go`, `writes_test.go` |
| 3 | T-2951 design, user guide, tooling.md in both copies, template changelog, close I-0277 | T-2947, T-2949 | `design/system/agent-narrative.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, the two `tooling.md`, `template/CHANGELOG.md`, the I-0277 file, `design/issues/summary.md` |

No two tasks in the same layer share a path.

## Figures

- Forecast: 45m, up from flai's 34m. S-0271 built `stream state` across the same three surfaces and took about 41m.
- Cost of delay: 12.50 USD a week, from the inputs flai set from I-0277. It stands as computed.
- `after`: none.

## Assumptions

- Decisions and Open questions are appended to, never replaced, so an earlier session's entries cannot be lost. Fixing a past decision is still done by hand where `wip/` is writable.
- A question written this way holds the move to review, as a hand-written one does. Questions for the designer should stay threads.
- No ADR is needed: this extends a command within ADR-0253's mounts. The design goes in `design/system/agent-narrative.md`.

## Proposal

- T-2950, the host channel, is optional for I-0277: the dashboard does not write these sections today. I kept it so that the three surfaces stay one call, as the design says. Drop it if you would rather keep the story to the CLI and MCP.
