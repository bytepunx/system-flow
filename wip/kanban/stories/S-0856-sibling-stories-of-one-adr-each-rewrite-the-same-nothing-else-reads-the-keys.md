---
id: S-0856
type: story
nature: improvement
title: Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict
status: backlog
owner: alex
created: 2026-11-27T22:55:28Z
updated: 2026-11-27T23:57:10Z
transitions: []
tags: []
topics: [documentation, check]
touches: [design/system/workflow.md, docs/operators/runbooks/change-the-board.md, design/conventions/documentation.md, design/conventions/strategic-agents.md, design/system/documentation-standard.md, design/system/agent-context.md, flai/internal/check/check_docs.go, flai/internal/check/docs_stories.go, flai/internal/check/docs_stories_test.go, design/system/flai-cli.md, docs/users/flai.md, design/adrs, design/adrs/README.md, design/issues/I-0414-sibling-stories-of-one-adr-each-rewrite-the-same-nothing-else-reads-the-keys.md, design/issues/summary.md]
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
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: flai
    at: 2026-11-27T22:55:28Z
---
# S-0856 Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict

## Goal

This story remediates [I-0414](../../../design/issues/I-0414-sibling-stories-of-one-adr-each-rewrite-the-same-nothing-else-reads-the-keys.md), "Sibling stories of one ADR each rewrite the same 'nothing else reads the keys yet' sentence in a hub design doc and runbook, so each acceptance stops the next on a rebase conflict", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0414 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0414 is closed with `flai issue close I-0414 --reason` saying what fixed it

## Tasks

## Notes

The analyzer found I-0414: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0414. time_lost_per_cycle 2h carried over from I-0414's Impact section. Its Impact time_lost_per_cycle was taken rather than the 1h57m derived from its cost and count.

Strategic usage carried over from I-0414: analyzer 1m25s, 1.4522 USD.

### Planning

Proposed remedy, from I-0414's three instances and I-0387 before it (the same shape on `design/system/server-performance.md`):

1. Rewrite the hub lines that exist now. `design/system/workflow.md` § Suites a step needs and runs keeps one sentence that lists which steps read the suite keys and which stories are "still to come", and `docs/operators/runbooks/change-the-board.md` keeps its twin, "Nothing else reads the keys yet but ...". Each step gets its own sentence or paragraph, in the present tense, with no running list of who landed. The same section's "until S-0813 retires them" is a forward reference of the same kind.
2. Make the rule explicit where agents read it: a project addition to `design/conventions/documentation.md` (describe each step on its own line, keep no list of which stories have landed or are still to come), and one to `design/conventions/strategic-agents.md` for the planner (sibling stories that change one document each get their own section or paragraph to write, named in their criteria).
3. A mechanical guard: an advisory `flai check` finding on a line of a design or docs document, outside ADRs and code spans, that names a story not yet started (backlog or ready). A "still to come (S-nnnn)" list is what makes every later sibling rewrite the line, and the first sibling's close-out would have noted it. It is advisory, as `wip.after-reason` is, so `--strict` passes over it; an ADR records it.

Touches, and where each came from:

- `design/system/workflow.md`, `docs/operators/runbooks/change-the-board.md`: design, the two files I-0414 names.
- `design/conventions/documentation.md`, `design/conventions/strategic-agents.md`: design, the project additions of remedy 2.
- `design/system/documentation-standard.md`: design, where the content rules every project's documents follow are stated; the new check is one.
- `flai/internal/check/check_docs.go`: layout, where `documentation()` walks the design and docs files.
- `flai/internal/check/docs_stories.go`, `flai/internal/check/docs_stories_test.go`: layout, new files for the rule and its test.
- `design/system/agent-context.md`: layout, its line naming S-0241, a backlog story, is the one other hit the rule finds today.
- `design/system/flai-cli.md`, `docs/users/flai.md`: co-change, 69% of the commits that changed the starting paths; they list `flai check`'s rules and its advisory warnings.
- `design/adrs` (folder) and `design/adrs/README.md`: the ADR for the rule. The folder stays because the ADR's file name is not known until its title is written; the index is generated.
- `design/issues/I-0414-...md`, `design/issues/summary.md`: declared by the second criterion.

After: S-0808, the last of ADR-0312's sibling stories, is in review and rewrites the very sentence of `design/system/workflow.md` this story restructures, and adds § What the acceptance's gates read (S-0808). This story reads that text and the behaviour it describes, the acceptance's gates reading `verify`, to write each step's own sentence. Written before S-0808 lands, the rewrite would describe the gates as they were and conflict with S-0808 on the line, which is I-0414 again.

Forecast: flai's 2m rests on size 2, two criteria and no touches. It is adjusted to 3h for four tasks: two document rewrites, a check rule with its test and ADR, and the close of the issue; the median improvement story of this size that adds a check rule took two to three hours of agent time.

Cost of delay: 300 USD a week, as `flai cod` gives it from the 2h `time_lost_per_cycle` I-0414 carried over, at 150 USD an hour and one cycle a week. It stands: S-0813, S-0814, and S-0815 under E-0034 still change the same section, so the cycle recurs.
