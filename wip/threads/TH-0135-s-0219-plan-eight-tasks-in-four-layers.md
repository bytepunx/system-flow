---
id: TH-0135
title: "S-0219 plan: eight tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0219-the-orchestrator-moves-and-orders-work-by-its-policy-within-its-permissions.md
  item: S-0219
status: resolved
participants: [planner-S-0219, alex]
created: 2026-10-05T04:49:09Z
updated: 2026-10-05T05:03:27Z
---

# TH-0135 S-0219 plan: eight tasks in four layers

On wip/kanban/stories/S-0219-the-orchestrator-moves-and-orders-work-by-its-policy-within-its-permissions.md.

## Entries

### 2026-10-05T04:49:09Z planner-S-0219
The plan for S-0219: eight tasks in the backlog, in four layers. No two tasks in a layer share a path.

| Layer | Task | Waits for | What it does |
|-------|------|-----------|--------------|
| 1 | T-0884 | none | `flai order` records who placed a story by hand and when; `flai order --by --apply` keeps a hand placement of the last day; ADR for where the record lives |
| 1 | T-0888 | none | `flai plan --candidates`: backlog epics with no stories, and open epics whose stories are all done |
| 1 | T-0890 | none | `flai promote --drafts`: each draft and what it lacks to be finalized |
| 2 | T-0893 | T-0884, T-0888, T-0890 | `flai guard` maps the four permissions onto calls; the new reads always pass |
| 2 | T-0896 | T-0890 | `plan` for the orchestrator (trigger `orchestrator`); finalizing a complete draft only; moving to ready only a candidate, and only while the ready limit has room |
| 3 | T-0899 | T-0893, T-0896 | The orchestrator's prompt, `orchestrator.md`, and the convention (with the template) |
| 3 | T-0901 | T-0893, T-0896 | Fixture-board test: each permission on and off, through the guard and flai |
| 4 | T-0904 | T-0899, T-0901 | strategic-agents.md, flai-cli.md, workflow.md, user and operator guides |

Story: 39 touches (6 declared and kept, the rest file by file from the tasks, plus co-change docs); forecast 1h40m, delivery 2026-10-05T13:24Z; cost of delay 101.35 USD a week, its share of E-0016's value. Reasons are under Notes › Planning. Topics `orchestration` and `planning` added.

Assumptions to confirm:

1. S-0218 delivers the `orchestrate` role, `orchestration.permissions`, the guard's refusal shape, the base prompt, `orchestrator.md`, and an operator's guide in `docs/operators`. These tasks extend those and do not build them. If S-0218's tasks map all seven permissions in the guard, T-0893 shrinks to the new reads and the tests.
2. S-0217's `flai order --by`, `flai promote --candidates`, `policy.go`, and `promote.go` exist when this story starts. S-0219 waits for S-0218, which waits for S-0217.
3. Nothing records today who placed a story or when: `board.md` keeps only the list. T-0884 adds a record that flai and the metrics can read. I suggested a `placed` map in `board.md`, since `.flai-cache/edits.jsonl` is local to a host. A placement by an orchestrator does not count as by hand.
4. "Whose stories are all done when the epic is not": I took this as an epic not done or cancelled whose stories are all done or cancelled, with at least one done. An epic whose planner thread is still waiting on you is left out.
5. Completeness of a draft (sections, a criterion, touches, forecast, value) is flai's arithmetic in T-0890. Consistency stays the orchestrator's judgement, as the designer split it for S-0217.
6. "Never a held story" and "while the ready limit has room" are enforced by `flai move` for `FLAI_ROLE=orchestrate`, not only by the prompt. Today a move to ready of a held story only warns.
7. No dashboard change and no hostapi write: S-0228 and S-0229 cover the orchestrator's pages and settings.

Proposals: none to split, merge, or drop. T-0896 is the largest task (plan, finalize, promote), but its three parts share `mcpserver` files, so splitting it would only chain them. The story's agent may split it when it pulls the story.

### 2026-10-05T05:03:27Z alex
Resolved.
