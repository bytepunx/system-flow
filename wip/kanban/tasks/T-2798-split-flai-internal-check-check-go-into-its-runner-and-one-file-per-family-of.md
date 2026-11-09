---
id: T-2798
type: task
nature: improvement
title: Split flai/internal/check/check.go into its runner and one file per family of rules
status: backlog
parent: S-0631
owner: alex
created: 2026-11-09T22:04:54Z
updated: 2026-11-09T22:04:54Z
transitions: []
stream: S-0631
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/check_layout.go, flai/internal/check/check_items.go, flai/internal/check/check_docs.go, flai/internal/check/check_roles.go, flai/internal/check/check_records.go]
---
# T-2798 Split flai/internal/check/check.go into its runner and one file per family of rules

## Work

`check.go` is 2,235 lines. Keep the runner in `check.go` and move each family of rules into its own file of package `check`. What any rule does stays the same, and so do the order the rules run in and the findings they give. The line ranges below are approximate, from the plan's read on 2026-11-09:

| Lines | Declarations | File |
|-------|--------------|------|
| 46–310 | `Finding`, `Result`, `checker`, `rule`, `sequence`, `Run`, `runRule`, `add`, `advise`, `keyLine`, `headingLine` | `check.go`, which keeps them |
| 311–456 | `layout`, `sharedPatterns`, `cacheIgnored`, `requiredHeadings` | `check_layout.go` |
| 457–1148 | `workItems`, `oneItem`, `epicLags`, `overlap`, `componentTag`, `after`, `afterHub`, `afterCycle`, `narratives`, `board` | `check_items.go` |
| 1149–1364 | `adrIndex`, `documentation`, `Vocabulary`, `docTopics`, `criteriaSection` | `check_docs.go` |
| 1365–1646 | `conventions`, `roles`, `formulas`, `agentCopies`, `unknownFields` | `check_roles.go` |
| 1647–2235 | `threads`, `messages`, `issues`, `experiments`, `analysis`, `digest`, `markdown`, `componentExamples`, `conflictMarkers` | `check_records.go` |

`items.go` (64 lines) is a fit for the work-item rules if that reads better. Keep every file under 1,200.

Change no test file. `check_test.go` stays as it is: the size test goes in a file of its own, in a later task of this story.

This task waits for nothing and shares no path with the other splits, so all of them run together.

## Done when

- No non-test file in `flai/internal/check` is over 2,000 lines, and none of the new files is over 1,200.
- `scripts/flai.sh check --strict` on the main checkout gives the same findings before and after.
- `git diff --stat` shows no `_test.go` file changed.
- `flai test flai/internal/check` passes, gofmt, goimports, and golangci-lint included.

## Notes
