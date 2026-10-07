---
id: S-0330
type: story
nature: feature
title: flai message sends a message from one open story's agent to another's, kept apart from the operator's threads
status: in-progress
parent: E-0018
owner: alex
created: 2026-10-07T20:10:22Z
updated: 2026-10-07T20:24:33Z
transitions:
  - to: ready
    at: 2026-10-07T20:24:28Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T20:24:33Z
    by: agent-S-0330
tags: [flai]
topics: [cli]
touches: [design/adrs, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/cmd/message.go, flai/cmd/message_test.go, flai/cmd/root.go, flai/cmd/accept.go, flai/cmd/accept_threads_test.go, flai/cmd/archive.go, flai/cmd/archive_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, wip/messages/README.md, design/system/agent-coordination.md, design/system/repository-layout.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 93
  estimated: true
  turns:
    - day: 2026-10-07
      work: 19
  models:
    - model: claude-opus-5-5
      input: 38
      output: 329
      cache_read: 1553026
      cache_write: 107187
      cost: 0.7421
  strategic:
    - kind: orchestrator
      seconds: 301
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 46
          output: 819
          cache_read: 3661515
          cache_write: 10356
          cost: 0.9048
cost_of_delay:
  value: 124.03
  by: planner-E-0018
  at: 2026-10-07T20:21:58Z
forecast:
  duration: 48m
  delivery: 2026-10-08T07:30:00Z
  basis: "flai forecast: median 114 s per unit of size over 29 done large-band feature stories on claude-opus-5-5, times size 25; 39th in the pull order with an in-progress limit of 3."
  by: planner-E-0018
  at: 2026-10-07T20:19:18Z
finalized:
  by: orchestrator
  at: 2026-10-07T20:24:14Z
---
# S-0330 flai message sends a message from one open story's agent to another's, kept apart from the operator's threads

## Goal

Give the agents of two open stories a channel of their own. Today they coordinate on threads, which are the operator's: every one an agent opens for another lands in the operator's inbox as awaiting them, and nothing addresses it to the other story's agent. A message is addressed from one story to another, so that the agents can agree on who changes what before their branches conflict, without asking the operator.

This story builds the store and the CLI: `flai message send`, `reply`, `list`, and `show`, with conversations closed when either story leaves the open columns. The MCP tools and the inbox come in S-0331.

## Acceptance criteria

- [ ] An ADR records how a message is addressed (from the sender's story to an open story), where messages are kept, how a conversation ends, and why messages are apart from threads (ADR-0020, ADR-0109).
- [ ] `flai message send <S-nnnn> "<text>" --from <S-nnnn> [--about <path>…]` starts a conversation between two stories in progress or in review; one to or from a story in any other state is refused with the reason.
- [ ] `flai message reply`, `flai message list [--story S-nnnn] [--all]`, and `flai message show` reply to a conversation, list the open ones a story is part of with which side it awaits, and print one with every entry; each has `--json`.
- [ ] No message appears in `flai thread list`, among the threads awaiting the operator, or in the dashboard's designer inbox.
- [ ] A conversation reads as closed once either of its stories is accepted, cancelled, or archived, and accepting or archiving a story writes an entry saying why in each conversation it closes.
- [ ] `flai check --strict` validates the message files and the markdown flai writes for them.
- [ ] `design/system/agent-coordination.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, and the generated `docs/users/flai-reference.md` describe messages.

## Tasks

- T-1185 An ADR records how messages between stories are addressed, kept, and closed, and why they are apart from threads
- T-1186 A messages package writes, reads, and lists conversations between two open stories
- T-1187 flai message send, reply, list, and show work the conversations, with --json
- T-1188 Acceptance and flai archive close a story's conversations, and flai check validates message files
- T-1189 The design and the users' guide describe messages between stories, and the command reference is regenerated

## Notes

### Planning

Planned by planner-E-0018 on 2026-10-07, as the first story of E-0018. The plan thread on E-0018 lists every story, its tasks, and the assumptions.

Assumption: messages are kept as one markdown file per conversation under `wip/messages/`, as threads are under `wip/threads/`. T-1185's ADR decides it and may choose otherwise; the touches follow the assumption.

Layers:

1. T-1185, the ADR.
2. T-1186, the package.
3. T-1187, the command; T-1188, closing and checking. They share no path and run together.
4. T-1189, the docs.

Touches:

- **Declared:** none before planning.
- **Layout:**
  - `flai/internal/messages/messages.go` and its test: a new package beside `flai/internal/threads`.
  - `flai/cmd/message.go` and its test, registered in `flai/cmd/root.go` as `newThreadCmd` is.
  - `flai/cmd/accept.go`, `flai/cmd/archive.go`, and their tests: where `threads.ResolveOnItems` closes threads at acceptance and archive (ADR-0109).
  - `flai/internal/check/check.go` and its test: where thread files are validated.
  - `wip/messages/README.md`: every folder a reader lands in has one.
- **Co-change:** `flai touches suggest` from the thread files gave `design/system/flai-cli.md` (38%), `flai/cmd/root.go` (31%), `flai/internal/check/check.go` (31%), and `docs/users/flai.md` (23%).
- **Design:** `design/system/agent-coordination.md` and `design/system/repository-layout.md`, which describe coordination and the `wip/` folders; `docs/users/flai-reference.md`, generated from the help.
- **Folder touch kept:** `design/adrs`. T-1185's ADR number and slug are picked by `flai adr new`. It lies inside `claims.shared` and holds no story.
- **Not taken:** `flai/internal/threads/threads.go` and `flai/internal/workitem/*`, which `suggest` listed at 15 to 23%; messages get their own package, and threads do not change.

Forecast 48m, delivery 2026-10-08T07:30Z: `flai forecast` gave it, 114 s per unit over 29 done large-band feature stories on claude-opus-5-5, times size 25 (7 criteria, 18 touches). It stands: the new package is about the size of `threads`, and its layers are the usual ADR, code, docs.

Cost of delay 124.03 USD a week: `flai cod` gave S-0330's share of E-0018's 1000 USD a week of penalty, 48m of the epic's 6h27m. It stands; the share by duration undervalues this story, which every other story of E-0018 waits on, but no input says by how much.
