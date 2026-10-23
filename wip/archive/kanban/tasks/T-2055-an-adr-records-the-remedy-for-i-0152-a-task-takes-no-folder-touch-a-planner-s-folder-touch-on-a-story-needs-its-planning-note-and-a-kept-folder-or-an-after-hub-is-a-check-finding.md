---
id: T-2055
type: task
nature: improvement
title: "An ADR records the remedy for I-0152: a task takes no folder touch, a planner's folder touch on a story needs its Planning note, and a kept folder or an after hub is a check finding"
status: done
parent: S-0468
owner: alex
created: 2026-10-22T21:04:24Z
updated: 2026-10-23T03:39:07Z
transitions:
  - to: ready
    at: 2026-10-23T03:34:08Z
    by: agent-S-0468
  - to: in-progress
    at: 2026-10-23T03:34:09Z
    by: agent-S-0468
  - to: done
    at: 2026-10-23T03:39:07Z
    by: agent-S-0468
stream: S-0468
tags: [cli]
touches: [design/adrs, design/adrs/README.md, design/system/workflow.md, design/system/strategic-agents.md, design/system/flai-cli.md]
usage:
  source: log
  seconds: 298
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 46
      output: 15705
      cache_read: 3072856
      cache_write: 100215
      cost: 1.5115
---
# T-2055 An ADR records the remedy for I-0152: a task takes no folder touch, a planner's folder touch on a story needs its Planning note, and a kept folder or an after hub is a check finding

## Work

Waits for nothing: the first layer. Write the ADR with `flai adr new` (or the MCP tool `adr_new`), topics `planning` and `workflow`, and accept it. Its context is I-0152's instance and `design/analysis/2026-10-15-bottlenecks.md` § Finding 3: E-0016's planned stories declared `flai/cmd`, `flai/internal/{harness,hostapi,mcpserver,serve,guard,manifest}`, `flaiover/src/routes`, and `docs/users`, and S-0218 and S-0223 were `after` hubs, so one story in progress or one hub awaiting acceptance held up to ten ready stories while three of five lanes stood empty; `exact_share` was 0.49 to 0.60. ADR-0046 (the hold), ADR-0096 § 3 (tasks narrow a folder claim, which needs tasks that name files), ADR-0134 (shares), and I-0087 (closed) are what stands already; the convention and `planPrompt` ask for files but nothing enforces it.

The decision, in these terms so that the code tasks build one thing:

1. A folder touch is a touches entry that, read as a path (a component's name or tag means its sub-project's path), names a folder of the repository tree holding at least one file. A path absent from the tree is a file to be made. A folder lying wholly inside a `claims.shared` pattern does not count: it holds nothing.
2. A task takes no folder touch. `flai task new --touches`, `flai edit --touches`, `flai touches`, and MCP `item_new` and `item_edit` refuse one, writing nothing, naming the folder and how many files it holds, and saying to name the files the task changes, or to leave the folder to the story's touches with a `### Planning` note.
3. A story's folder touch written by the planner (`FLAI_ROLE=plan`, or an author named `planner-…`) is refused unless the story's Notes name the folder's path under `### Planning`. The operator's and the story's agent's story touches are not refused.
4. `flai check` reports `touches.folder`, a warning, on a story in progress whose claim (ADR-0096's) keeps a folder touch, naming the folder, its file count, and the ready stories it holds. A check scoped to a story keeps the finding on that story and drops another story's, recording no issue, as it does `wip.overlap` (ADR-0115).
5. `flai check` reports `wip.after-hub`, a warning, on a story not done that three or more stories in backlog, ready, or in progress name in `after`, naming them.
6. `after` on a story names only a story whose accepted change the story needs on main, never an order of preference: the pull order is for that.

Alternatives to record: warning instead of refusing a task's folder touch (the convention already asks and was not followed); refusing every writer's story folder touch (the operator's words are theirs); releasing an `after` hold when the hub reaches review (the hub's change is not on main, so the dependent branch cannot build on it); a manifest threshold for the hub (three is stated, a setting adds a settings row and a test for nothing yet asked for).

Then describe it where the design lives: `design/system/workflow.md` § Branches and collisions (ADR-0019), beside ADR-0096's claim rules; `design/system/strategic-agents.md` § Enriching a story (S-0210) › What the planner writes, including that `planner.md` lacks the folder rule until the convention task adds it; and `design/system/flai-cli.md` where check rule names are listed. Link the ADR from each and bump `updated`.

## Done when

- The ADR is accepted, indexed in `design/adrs/README.md`, and states the six points and the alternatives.
- The three design sections describe the refusals and the two findings and link the ADR.
- `flai check --strict` and the markdown lint pass on the files changed.

## Notes

Layer 1, alone. `design/adrs` is in `claims.shared`, so this task's folder touch holds nothing.
