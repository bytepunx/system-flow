---
id: S-0856
type: story
nature: improvement
title: Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict
status: backlog
owner: alex
created: 2026-11-27T22:55:28Z
updated: 2026-11-28T00:03:04Z
transitions: []
tags: []
topics: [documentation, check]
touches: [design/system/workflow.md, docs/operators/runbooks/change-the-board.md, design/conventions/documentation.md, design/conventions/strategic-agents.md, design/system/documentation-standard.md, design/system/agent-context.md, flai/internal/check/check_docs.go, flai/internal/check/docs_stories.go, flai/internal/check/docs_stories_test.go, design/system/flai-cli.md, docs/users/flai.md, design/adrs, design/adrs/README.md, design/issues/I-0414-sibling-stories-of-one-adr-each-rewrite-the-same-nothing-else-reads-the-keys.md, design/issues/summary.md, flai/internal/check/scope.go, flai/internal/check/scope_test.go]
after: [S-0808]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 215
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 210
          output: 38942
          cache_read: 12135375
          cache_write: 549141
          cost: 5.33
    - kind: orchestrator
      seconds: 83
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 21
          output: 96
          cache_read: 1309564
          cache_write: 158628
          cost: 0.4675
    - kind: analyzer
      seconds: 85
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 55
          output: 5532
          cache_read: 391176
          cache_write: 29514
          cost: 0.9652
        - model: claude-opus-5-5
          input: 29
          output: 4664
          cache_read: 1264948
          cache_write: 28114
          cost: 0.487
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: flai
    at: 2026-11-27T22:55:28Z
  value: 300
  by: planner-S-0856
  at: 2026-11-27T23:57:39Z
forecast:
  duration: 1h30m
  delivery: 2026-11-28T09:44:00Z
  basis: "Its own forecast of 1h30m; 16th in the pull order with an in-progress limit of 3, behind S-0837, S-0844, S-0846, S-0835, S-0814, S-0852, S-0813, S-0822, S-0848, S-0853, S-0854, S-0824, S-0845, S-0241, S-0815, S-0849, S-0851 and S-0855."
  by: flai
  at: 2026-11-28T00:03:04Z
finalized:
  by: alex
  at: 2026-11-28T00:01:56Z
---
# S-0856 Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict

## Goal

This story remediates [I-0414](../../../design/issues/I-0414-sibling-stories-of-one-adr-each-rewrite-the-same-nothing-else-reads-the-keys.md), "Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0414 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0414 is closed with `flai issue close I-0414 --reason` saying what fixed it

## Tasks
- T-3627 Rewrite the suite keys' hub sentence in workflow.md and the board runbook so each step has its own sentence, with no running list of stories
- T-3632 flai check advises docs.unstarted-story on a design or docs line that names a story in backlog or ready, with an ADR and a test that reproduces I-0414's list
- T-3635 Project additions: documentation keeps no running list of sibling stories, and the planner gives each sibling of one document its own section
- T-3637 Close I-0414 with the reason naming the rewritten hub lines, the docs.unstarted-story check, and the project additions

## Notes

The analyzer found I-0414: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0414. time_lost_per_cycle 2h carried over from I-0414's Impact section. Its Impact time_lost_per_cycle was taken rather than the 1h57m derived from its cost and count.

Strategic usage carried over from I-0414: analyzer 1m25s, 1.4522 USD.

### Planning

Proposed remedy, from I-0414's three instances and I-0387 before it (the same shape on `design/system/server-performance.md`):

1. Rewrite the hub lines that exist now. `design/system/workflow.md` § Suites a step needs and runs keeps one paragraph that names each step that reads the suite keys "since S-nnnn" and then lists them again by story, and `docs/operators/runbooks/change-the-board.md` keeps its twin, "Nothing else reads the keys yet but ...", which still says the acceptance's gates run as they did although S-0808 changed them. Each step gets its own sentence or paragraph, in the present tense, with no running list of who landed. The same section's "until S-0813 retires them" is a forward reference of the same kind.
2. Make the rule explicit where agents read it: a project addition to `design/conventions/documentation.md` (describe each step on its own line, keep no list of which stories have landed or are still to come), and one to `design/conventions/strategic-agents.md` for the planner (sibling stories that change one document each get their own section or paragraph to write, named in their criteria).
3. A mechanical guard: an advisory `flai check` finding on a line of a design or docs document, outside ADRs and code spans, that names a story not yet started (backlog or ready). A "still to come (S-nnnn)" list is what makes every later sibling rewrite the line, and the first sibling's close-out would have noted it. It is advisory, as `wip.after-reason` is, so `--strict` passes over it; an ADR records it.

Touches, and where each came from:

- `design/system/workflow.md`, `docs/operators/runbooks/change-the-board.md`: design, the two files I-0414 names.
- `design/conventions/documentation.md`, `design/conventions/strategic-agents.md`: design, the project additions of remedy 2.
- `design/system/documentation-standard.md`: design, where the content rules every project's documents follow are stated; the new check is one.
- `flai/internal/check/check_docs.go`: layout, where `documentation()` walks the design and docs files.
- `flai/internal/check/docs_stories.go`, `flai/internal/check/docs_stories_test.go`: layout, new files for the rule and its test.
- `flai/internal/check/scope.go`, `flai/internal/check/scope_test.go`: layout, where a check scoped to a story drops the advisory findings outside it (its filter lists `wip.after-reason` and the others); the new rule joins them, so a close-out records no issue for another story's document. flai added them to the story when T-3632 named them, which is the touch edit by system-flow at 2026-11-27T23:57:57Z.
- `design/system/agent-context.md`: layout, its line naming S-0241, a backlog story, is the one other hit the rule finds today.
- `design/system/flai-cli.md`, `docs/users/flai.md`: co-change, 69% of the commits that changed the starting paths; they list `flai check`'s rules and its advisory warnings.
- `design/adrs` (folder) and `design/adrs/README.md`: the ADR for the rule. The folder stays because the ADR's number is not known until it is written: T-3618 of S-0852 and T-3641 of S-0854 plan ADRs at the same time, and all three had guessed 0322. No task names the ADR's file for that reason; the index is generated.
- `design/issues/I-0414-...md`, `design/issues/summary.md`: declared by the second criterion.
- Left out of `flai touches suggest`'s list: `docs/operators/index.md`, `docs/users/flai-reference.md`, `template/CHANGELOG.md`, and the rest it ranks at 9% or below. The rule is a check finding, not a command or flag, so the reference and the operators' index do not change, and the project additions are not baseline, so the template does not.

After: S-0808, the last of ADR-0312's sibling stories to change the hub sentence, rewrote the very sentence of `design/system/workflow.md` this story restructures, and added § What the acceptance's gates read (S-0808). This story reads that text and the behaviour it describes, the acceptance's gates reading `verify`, to write each step's own sentence. S-0808 was accepted on 2026-11-27, so the `after` is met and holds nothing.

Forecast: `flai forecast` gives 35m, size 19 (2 criteria, 17 touches) at the 109 s median of 187 done improvement stories in the large band. It is raised to 1h30m: the size counts touches, not the four tasks, one of which is a new Go check rule with its test and an ADR, and the close-out alone ran 602 s at W47's p50. Delivery is flai's, from that duration and the story's place in the pull order.

Cost of delay: 300 USD a week, as `flai cod` gives it from the 2h `time_lost_per_cycle` I-0414 carried over, at 150 USD an hour and one cycle a week. It stands: S-0813 and S-0815 under E-0034 still change the same section of `design/system/workflow.md`, so the cycle can recur.
