---
id: T-3189
type: task
nature: feature
title: The workflow package names the rule vocabulary, each rule refusing or warning, and refuses a transition naming a rule it lacks
status: done
parent: S-0738
owner: alex
created: 2026-11-12T16:05:23Z
updated: 2026-11-15T00:03:38Z
transitions:
  - to: ready
    at: 2026-11-14T23:56:31Z
    by: agent-S-0738
  - to: in-progress
    at: 2026-11-14T23:56:31Z
    by: agent-S-0738
  - to: done
    at: 2026-11-15T00:03:38Z
    by: agent-S-0738
stream: S-0738
tags: [cli]
touches: [flai/internal/workflow/rules.go, flai/internal/workflow/rules_test.go, flai/internal/workflow/workflow.go]
usage:
  source: log
  seconds: 427
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 659
      cache_read: 245572
      cache_write: 7504
      cost: 0.1155
---
# T-3189 The workflow package names the rule vocabulary, each rule refusing or warning, and refuses a transition naming a rule it lacks

## Work

Add the rule names the ADR fixes to `flai/internal/workflow/rules.go`, each with whether it refuses or warns: criteria present, not a draft, at least one task, no open question, nothing uncommitted, every criterion ticked, every child closed, and a reason. Make `Parse` refuse a transition whose rule the vocabulary lacks, naming the rule and the transition. The rules' checks stay in `workitem`, which knows the item, its narrative, and its worktree; this package holds only the names and their severity.

## Done when

`rules_test.go` lists every rule with its severity and refuses a workflow naming an unknown one.

## Notes
