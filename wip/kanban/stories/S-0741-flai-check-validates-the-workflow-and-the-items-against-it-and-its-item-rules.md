---
id: S-0741
type: story
nature: feature
title: flai check validates the workflow and the items against it, and its item rules read a state's kind instead of its name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:59:17Z
updated: 2026-11-13T07:53:06Z
transitions: []
tags: [cli]
topics: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/check_workflow.go, flai/internal/check/workflow_test.go, flai/internal/check/check_items.go, flai/internal/check/items_test.go, flai/internal/check/scope.go, flai/internal/check/sequence_test.go, flai/internal/check/vocabulary_test.go, design/system/flai-cli.md, docs/users/flai.md]
after: [S-0738]
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
    - kind: orchestrator
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 14
          cache_read: 108314
          cache_write: 18020
          cost: 0.0361
draft: true
forecast:
  duration: 29m
  delivery: 2026-11-13T15:41:00Z
  basis: "Its own forecast of 29m; 18th in the pull order with an in-progress limit of 3, behind S-0660, S-0760, S-0735, S-0745, S-0759, S-0734, S-0751, S-0757, S-0740, S-0661, S-0663, S-0662, S-0241, S-0730, S-0736, S-0737, S-0738 and S-0739."
  by: flai
  at: 2026-11-13T07:53:06Z
---
# S-0741 flai check validates the workflow and the items against it, and its item rules read a state's kind instead of its name

## Goal

`flai check` tells the operator what is wrong with a workflow before a move trips on it, and its rules on items follow the workflow rather than the six names. It reports a state no transition reaches, a state with no way to an accepted or a cancelled one, a workflow with no pull or no acceptance, a timer that names a state the file lacks, a rule the vocabulary lacks, and an item in a state the file lacks, each with what to do.

## Acceptance criteria

- [ ] `flai check` reports each of these as a finding of its own code, under a `workflow.` prefix, naming the state, transition, timer, or item and what to do about it: a state no transition reaches, a state with no way to an accepted or a cancelled state, no transition with the pull effect, no transition with the acceptance effect, a timer naming a state the file lacks, a rule the vocabulary lacks, and an item whose status the workflow lacks
- [ ] The baseline workflow passes every one of them
- [ ] The item rules in `check_items.go` and `scope.go` name no state: each reads a state's kind, so a story in a wait state is checked as one in `review` is today
- [ ] Tests give each finding on a workflow written to fail it, and run the item rules against a workflow with a state renamed and one added between `in-progress` and `review`
- [ ] With the baseline, every existing `check` test passes unchanged
- [ ] `design/system/flai-cli.md` and `docs/users/flai.md` list the new findings with what clears each

## Tasks

- T-3204 flai check reports the workflow findings: an unreachable state, a dead end, no pull, no acceptance, a bad timer, an unknown rule, and an item in a missing state
- T-3205 flai check's item rules and its scope read a state's kind instead of its name
- T-3211 flai-cli.md and the user guide list the workflow findings and what clears each

## Notes

- Check rules are registered in the `sequence` slice in `flai/internal/check/check.go`; `check_roles.go` is the model for a check over a file under `design/workflows/`.

### Planning

- Touches, all files:
  - `flai/internal/check/check.go`, `check_workflow.go`, `workflow_test.go`, `sequence_test.go`, `vocabulary_test.go`: layout. A new check beside `check_roles.go`, registered in `sequence`, whose tests list every code.
  - `flai/internal/check/check_items.go`, `items_test.go`, `scope.go`: layout. 21 references to a state's name in `check_items.go`, one in `scope.go`.
  - `design/system/flai-cli.md`, `docs/users/flai.md`: design (criterion 6). Both shared paths.
  - Left out from `flai touches suggest`: the operator guide, `flaiover-dashboard.md`, and `workflow.md` (14% to 26%). The runbook for changing the board is S-0755's.
- `after`: S-0738. The rule finding checks names against the vocabulary S-0738 adds, and the findings read the parser and the workflow S-0737 gives, which S-0738 builds on.
- Tasks: T-3204 (the findings) and T-3205 (the item rules) together, no path in common; then T-3211 (the documents), which lists T-3204's codes.
- Forecast 29m, as `flai forecast` gives it for size 16. `check_roles.go` is a close model for the findings.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
