---
id: T-3202
type: task
nature: feature
title: Repo.Move runs the rules its transition lists, in their order, with today's refusals and warnings word for word, and each rule is tested
status: backlog
parent: S-0738
owner: alex
created: 2026-11-12T16:06:14Z
updated: 2026-11-12T16:06:14Z
transitions: []
stream: S-0738
tags: [cli]
touches: [flai/internal/workitem/rules.go, flai/internal/workitem/rules_test.go, flai/internal/workitem/workitem_test.go, flai/internal/workitem/worktree_test.go]
after: [T-3197]
---
# T-3202 Repo.Move runs the rules its transition lists, in their order, with today's refusals and warnings word for word, and each rule is tested

## Work

Waits for T-3197: it runs the rules the baseline lists after that task. Map each rule name to its check in `workitem`, the code now at `rules.go` lines 58 to 118, and make `Repo.Move` run the transition's rules in order, refusing or warning as each rule's severity says, with today's text. Add tests for each rule, the uncommitted-worktree rule and the every-child-closed rule at done included, and one that puts a rule on a transition of a workflow other than the baseline.

## Done when

Every rule refuses or warns in its own test with today's words, and the existing move tests pass unchanged.

## Notes
