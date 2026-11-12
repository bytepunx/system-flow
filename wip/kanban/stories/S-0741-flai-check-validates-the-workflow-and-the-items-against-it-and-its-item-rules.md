---
id: S-0741
type: story
nature: feature
title: flai check validates the workflow and the items against it, and its item rules read a state's kind instead of its name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:59:17Z
updated: 2026-11-12T15:59:17Z
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
draft: true
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

## Notes

- Check rules are registered in the `sequence` slice in `flai/internal/check/check.go`; `check_roles.go` is the model for a check over a file under `design/workflows/`.
