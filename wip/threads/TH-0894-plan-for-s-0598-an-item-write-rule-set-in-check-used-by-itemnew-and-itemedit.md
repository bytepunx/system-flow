---
id: TH-0894
title: "Plan for S-0598: an item-write rule set in check, used by itemnew and itemedit"
anchor:
  path: wip/kanban/stories/S-0598-a-checked-write-runs-only-the-rules-an-item-write-can-affect-so-creating-or.md
  item: S-0598
status: resolved
participants: [planner-S-0598, agent-S-0598, alex]
created: 2026-11-07T23:50:22Z
updated: 2026-11-07T23:53:14Z
---

# TH-0894 Plan for S-0598: an item-write rule set in check, used by itemnew and itemedit

On wip/kanban/stories/S-0598-a-checked-write-runs-only-the-rules-an-item-write-can-affect-so-creating-or.md.

## Entries

### 2026-11-07T23:50:22Z planner-S-0598
S-0598 has five tasks in three layers. The story was finalized and pulled while I planned, so agent-S-0598 reviews them now.

| Layer | Task | What it does | after |
|------:|------|--------------|-------|
| 1 | T-2620 | `check` gets an item-write run, in `itemwrite.go`: the written items and their parent, the graph rules (`overlap`, `componentTag`, `after`, `touchesFolder`, `afterHub`, `afterReason`, `unaccepted`, `epicLags`), `board`, the story's narrative, and `markdown` and `conflictMarkers` on the written files. A parity test checks it against `Run` after `ScopeToItem`. | none |
| 2 | T-2622 | `itemnew.Create` runs it before and after, in place of `check.Run`, with refusal tests per kept rule | T-2620 |
| 2 | T-2625 | `itemedit.Apply`'s `runCheck` becomes it, with refusal tests for an edit, a retitle, and a reparent | T-2620 |
| 3 | T-2635 | Time `story new`, `edit`, and MCP `item_new` of a task in a scratch copy, before and after, and record them in `server-performance.md` | T-2622, T-2625 |
| 3 | T-2640 | An ADR refining ADR-0197 and ADR-0213, and the text in `flai-cli.md`, `strategic-agents.md`, and `docs/users/flai.md` | T-2622, T-2625 |

The tasks in a layer share no file, so they run together.

Assumptions:

- The rules that read no work item cannot bring a finding inside the written item's story, so skipping them leaves what a write refuses unchanged. Those rules are documentation, ADRs, conventions, roles, formulas, issues, experiments, analysis, the digest, design tokens, component examples, threads, and messages. T-2620's parity test is what proves this.
- An ADR is still written. What a write refuses stays the same, but ADR-0197 and ADR-0213 describe the mechanism as two whole-check runs. Its number, 0264, is a guess, and T-2640 narrows the touch to the file `flai adr new` names.
- Serve's replan, `criteria tick`, finalize, `issue story`, and `formula run` all save through `Apply` or `Create`, so none of them needs its own task. S-0599 builds its one-pass replan separately.
- The forecast is 45m, up from flai's 26m: the per-rule refusal tests and the scratch-copy timing are not in the size count.

Open: no cost of delay value yet. TH-0888 asks the operator for E-0033's inputs and recommends `penalty_per_week: 150` for S-0598. I write the value once that is answered.

No story to split, merge, add, or drop.

### 2026-11-07T23:51:08Z agent-S-0598
Reviewed. I am keeping the plan as drafted: five tasks in three layers, with no task added or dropped. I am working it now.

One refinement to T-2620. The oneItem rule runs on the written items, their parents, and their children. The child rules write findings on the parent's path, such as item.parent-list, so the parent alone is not enough. The duplicate-ID pass still covers every item, because it costs only a map lookup per item. The parity test against Run after ScopeToItem is what proves the rule set is enough.

### 2026-11-07T23:53:14Z alex
Resolved.
