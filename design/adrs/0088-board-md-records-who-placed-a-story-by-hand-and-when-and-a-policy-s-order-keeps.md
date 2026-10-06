---
id: ADR-0088
title: "board.md records who placed a story by hand and when, and a policy's order keeps a placement of the last day"
status: accepted
date: 2026-10-06
supersedes: []
superseded_by: []
topics: [orchestration, board, cli]
---

# ADR-0088 board.md records who placed a story by hand and when, and a policy's order keeps a placement of the last day

## Context

The pull order is the `order` list in `wip/kanban/board.md`'s front matter (S-0057). `flai order <story> --before|--after|--top|--bottom` rewrites it, and the dashboard's drag runs that command. The list says where each story is and nothing about who put it there or when; `Reordered` in `changes.go` can only compare two lists. S-0217 gave the orchestrator `flai order --by <policy> --apply`, which rewrites the whole ready column by cost of delay, WSJF, throughput, or FIFO, and ADR-0087 lets it run that behind `order_ready`. An operator who drags a story to the top because they know something the figures do not would see the orchestrator undo it at its next decision. S-0219 asks that the orchestrator not reorder a story the operator placed by hand in the last day, and that needs a record of each placement that every host and the metrics can read, and a rule for how a policy order treats it.

## Decision

`board.md` records each placement by hand beside the order, and an order computed by a policy keeps a recent one in place.

1. **Where it lives.** A `placed` map in `board.md`'s front matter, after `order`, keyed by story ID, each entry with `by` and `at` (UTC, `2006-01-02T15:04:05Z`). It is written only when it has an entry, sorted by ID, so a board nobody has placed on is unchanged.
2. **What records it.** Every placement `flai order <story>` makes replaces that story's entry. `by` is `--placed-by`, else `FLAI_AGENT`, else the config author, as `flai move` stamps a transition. The dashboard's drag runs `flai order` with `FLAI_AGENT=flaiover`, so it is recorded as `flaiover`. When `flai serve` runs the orchestrator (`FLAI_ROLE=orchestrate`), the entry is `orchestrator` whatever `FLAI_AGENT` or `--placed-by` say, so that its placements never pass for one by hand. `flai order --by --apply` records nothing: it is a policy's order, not a placement.
3. **When it goes.** `flai move` drops a story's entry when it takes the story out of its column: to in-progress or cancelled, and also between ready and backlog, since a story that changes column is no longer where it was placed.
4. **Who counts as by hand.** Anyone but the orchestrator, recognised by the name `flai serve` runs it under, `orchestrator` (ADR-0087).
5. **The window.** `flai order --by <policy>` takes `--keep-placed <duration>`, a day by default, `0` for none. A ready story placed by hand within the window keeps its position in the column; the other stories fill the remaining positions in the policy's order. The printed order marks each kept story with who placed it and when; `--json` gives each a `kept` object and the order a `keep_placed` duration. The MCP tool `order_by_policy` and the host API's `order.by` read the same order with the default window, so the orchestrator sees what `--apply` will write.

## Consequences

- The operator's drag holds against the orchestrator for a day without the operator turning `order_ready` off. A story kept in place still shows its figure, so the operator sees what the policy would have done.
- The record is plain front matter in a tracked file: every host, `flai stats`, and the dashboard read it from git with no cache or service, and its history is in the board's commits.
- `board.md` grows by three lines per placed ready or backlog story, bounded by those two columns, since a story leaving them takes its entry with it.
- A placement through the dashboard records `flaiover`, not the operator's name; telling it apart from the operator's own terminal needs the host API to pass `--placed-by`.
- Recognising the orchestrator by name means a person or agent who sets `FLAI_AGENT=orchestrator` is treated as it. `FLAI_ROLE` covers the orchestrator whose harness settings override `FLAI_AGENT` in its shell.
- An ordering run with `--keep-placed 0` restores the policy's order over everything, for an operator who wants it.

## Alternatives considered

- **A `placed` block in each story's front matter.** A placement is a fact about the board, not the story; it would make every drag an edit of a story file and a change in its history, and the order already lives in `board.md`.
- **Reading placements from the git history of `board.md`.** It does not say who made a change through the dashboard, or which of several placements in one commit, and it is slow on every policy order.
- **A record in `.flai-cache` or in `flai serve`'s state.** It is per host and untracked, so another host and the metrics could not read it, and a placement made with no `flai serve` running would be lost.
- **Recording the placements `--apply` makes as the orchestrator's.** They would carry nothing the window uses, and an operator who ran `--apply` by hand would freeze the whole column for a day.
- **Pinning a hand-placed story to its neighbours rather than its position.** Neighbours the policy moves would carry it with them, which is not where the operator put it.
- **Keeping placements with no window.** An old drag would hold the policy off for good; a day is long enough for the operator's decision to be acted on and short enough that the policy takes over again.
