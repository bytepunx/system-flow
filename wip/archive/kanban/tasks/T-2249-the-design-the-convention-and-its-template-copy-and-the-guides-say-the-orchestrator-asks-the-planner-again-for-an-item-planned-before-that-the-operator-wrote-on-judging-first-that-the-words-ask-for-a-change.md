---
id: T-2249
type: task
nature: improvement
title: The design, the convention and its template copy, and the guides say the orchestrator asks the planner again for an item planned before that the operator wrote on, judging first that the words ask for a change
status: done
parent: S-0524
owner: alex
created: 2026-10-28T14:42:36Z
updated: 2026-10-30T14:25:32Z
transitions:
  - to: ready
    at: 2026-10-30T14:20:38Z
    by: agent-S-0524
  - to: in-progress
    at: 2026-10-30T14:20:39Z
    by: agent-S-0524
  - to: done
    at: 2026-10-30T14:25:32Z
    by: agent-S-0524
stream: S-0524
tags: [flai, orchestrator, planner, docs]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-2246, T-2248]
usage:
  source: log
  seconds: 293
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 49
      output: 18282
      cache_read: 3822682
      cache_write: 165020
      cost: 2.0929
---
# T-2249 The design, the convention and its template copy, and the guides say the orchestrator asks the planner again for an item planned before that the operator wrote on, judging first that the words ask for a change

## Work

Third layer: waits for T-2246, whose ADR it links, and T-2248, which settles the trigger's form and the prompt the design describes.

Record the behaviour T-2247 and T-2248 build, linking the ADR T-2246 wrote:

- `design/system/strategic-agents.md`: under The planner › Starting it, the orchestrator's exception names the third kind; under The orchestrator › What it does with each permission (S-0219), the `plan_backlog_epics` and `plan_backlog_stories` rows say the candidates include an epic or story planned before that the operator has written on since its newest planner run, with the thread named, and that the orchestrator plans it only when the operator's words ask for a change to the plan; under The planner › Planning again › What a run records, the trigger table gains `orchestrator; TH-nnnn`, said as the orchestrator asked after the operator wrote on that thread. Bump `updated`.
- `design/system/flai-cli.md`: the `flai plan --candidates` row names the third kind.
- `design/conventions/strategic-agents.md`, under As the orchestrator, the `plan_backlog_epics` and `plan_backlog_stories` rules: also an epic or story the candidates list as planned before, when the thread it names carries an operator's entry that asks for a change to the plan, a story or task to add, split, merge, or drop, a changed outcome, or a changed figure, and never for a confirmation that changes nothing. Make the same change in `template/root/design/conventions/strategic-agents.md`, keeping the two copies equal above the baseline marker, and add the convention change to `template/CHANGELOG.md`.
- `docs/users/flai.md`: the `flai plan --candidates` description and the orchestrator's planning section (the `plan_backlog_epics` row and the `plan_backlog_stories` paragraph) say what the list now holds and what the orchestrator does with it. `docs/users/flai-reference.md`: regenerate with `make flai-reference` when the command's help text changed in T-2247, or edit the `--candidates` prose where it does not come from the help.
- Where `flai/internal/harness/harness.go`'s orchestrator prompt names the candidates (`each story flai plan --candidates lists`), it stays true as written; change it only if T-2248 left it saying less than the convention. That file is T-2248's.

## Done when

- The sections above describe the third kind of candidate, the orchestrator's judgement, and the `orchestrator; TH-nnnn` trigger, and link the ADR from T-2246.
- The two copies of the convention are equal above `<!-- system-flow:end-of-baseline -->`, and the template changelog has the entry.
- `flai check --strict` and the markdown lint pass on the files changed, and `flai test` on the changed paths passes.

## Notes

`docs/users/flai.md`, `docs/users/flai-reference.md`, and `design/system/flai-cli.md` are in the manifest's shared paths (`claims.shared`), so touching them holds no ready story (ADR-0096).
