---
id: TH-0871
title: "S-0590 plan: import writes the lock, three tasks in three layers, 45m, 25 USD a week"
anchor:
  path: wip/kanban/stories/S-0590-flai-import-writes-no-system-flow-lock-yaml-so-the-import-formula-s-flai-check.md
  item: S-0590
status: resolved
participants: [planner-S-0590, orchestrator]
created: 2026-11-07T10:04:42Z
updated: 2026-11-07T10:06:23Z
---

# TH-0871 S-0590 plan: import writes the lock, three tasks in three layers, 45m, 25 USD a week

On wip/kanban/stories/S-0590-flai-import-writes-no-system-flow-lock-yaml-so-the-import-formula-s-flai-check.md.

## Entries

### 2026-11-07T10:04:42Z planner-S-0590
The plan is written: three tasks in three layers. The story was finalized and started while I planned, so agent-S-0590 reviews these tasks as it works them.

## Proposed fix

`runImport` in `flai/cmd/import.go` renders the template but never saves `system-flow.lock.yaml`, while `flai new` does. The import saves the lock from the render's hashes, as `flai new` does, and `commitImportPaths` commits it with `--commit` and `--formula`.

## Tasks

| Layer | Task | Waits for | Why it waits |
|-------|------|-----------|--------------|
| 1 | T-2570: import writes and commits the lock; a reproducing test on `--formula`; help and `docs/users/flai-reference.md` | none | |
| 2 | T-2571: `docs/users/flai.md`, `design/system/template.md`, `design/system/flai-cli.md` say import writes the lock | T-2570 | it documents what T-2570 built |
| 3 | T-2572: run the import tests and the integration tier, close I-0302 | T-2570, T-2571 | the issue closes once the fix and its docs are in |

## Figures

- Forecast 45m, up from flai's 24m. The change under `flai/cmd` makes the close-out run the integration and smoke tiers, and the size counts neither those tiers nor the hash choice.
- Cost of delay 25.00 USD a week, from the 10m input as `flai cod` gives it.

## Assumptions

- The planner run's prompt said the operator had written on a thread named `asked`. No such thread exists, and `flai plan --candidates` listed S-0590 as never planned. So this is a first plan. If you asked for something specific, say so here.
- The lock records the rendered hash of a file the import kept, as `flai new`'s does. `flai upgrade` then treats that file as the project's own and never overwrites it. Hashing the current files, as `--relock` does, would let an upgrade replace them.
- No ADR records "import writes no lock", so the story changes the design documents and writes no ADR. ADR-0015 names only `flai new` and `flai upgrade` as writers of the lock. Say so if you want an ADR that adds import.
- Projects imported before the fix keep the `layout.lock` hint to run `flai upgrade --relock`. This story does not backfill them.

### 2026-11-07T10:06:04Z orchestrator
Approved as written. agent-S-0590 works T-2570 to T-2572 in their layers.

**Why the plan fits:**
- T-2570 fixes the cause and adds the reproducing test, which meets criterion 1.
- T-2572 closes I-0302 with `flai issue close`, which meets criterion 2.
- T-2571 updates the living design in the same story, as `documentation.md` asks.
- The story's ten touches cover every task's touches.
- The figures stand: a 45m forecast with its basis stated, and 25 USD a week from flai's input of 10m. That input was set by flai from I-0302, so I changed nothing.

**Your assumptions:**
- No `asked` thread: this is a first plan. Nothing more is asked for.
- Rendered hash, as `flai new` records it: yes. It keeps a kept file as the project's own.
- No ADR: right. Import writing the lock as `flai new` does is a living-design edit plus a narrative entry. An ADR can follow if the operator wants one.
- No backfill of projects imported before the fix: yes. `flai upgrade --relock` stays their path.

Source: design/conventions/decisions.md § When in doubt

### 2026-11-07T10:06:23Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: the tasks meet both criteria, the touches cover every task, and the figures stand as flai and the planner gave them.
