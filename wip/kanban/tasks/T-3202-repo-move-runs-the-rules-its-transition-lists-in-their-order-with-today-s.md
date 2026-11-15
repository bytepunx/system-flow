---
id: T-3202
type: task
nature: feature
title: Repo.Move runs the rules its transition lists, in their order, with today's refusals and warnings word for word, and each rule is tested
status: done
parent: S-0738
owner: alex
created: 2026-11-12T16:06:14Z
updated: 2026-11-15T00:29:53Z
transitions:
  - to: ready
    at: 2026-11-15T00:10:24Z
    by: agent-S-0738
  - to: in-progress
    at: 2026-11-15T00:10:24Z
    by: agent-S-0738
  - to: done
    at: 2026-11-15T00:29:53Z
    by: agent-S-0738
stream: S-0738
tags: [cli]
touches: [flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go]
after: [T-3197]
usage:
  source: log
  seconds: 1169
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 55
      output: 18874
      cache_read: 4582862
      cache_write: 143403
      cost: 2.226
---
# T-3202 Repo.Move runs the rules its transition lists, in their order, with today's refusals and warnings word for word, and each rule is tested

## Work

Waits for T-3197: it runs the rules the baseline lists after that task. Map each rule name to its check in `workitem`, the code now at `rules.go` lines 58 to 118, and make `Repo.Move` run the transition's rules in order, refusing or warning as each rule's severity says, with today's text. Add tests for each rule, the uncommitted-worktree rule and the every-child-closed rule at done included, and one that puts a rule on a transition of a workflow other than the baseline.

## Done when

Every rule refuses or warns in its own test with today's words, and the existing move tests pass unchanged.

## Notes
