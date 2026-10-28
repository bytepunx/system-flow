---
id: I-0221
title: The close-out's go-test tier runs only the changed packages, so a change that breaks a caller package's tests passes before review
class: defect
status: open
count: 1
cost: 20m
first_reported: 2026-10-28T19:01:11Z
last_reported: 2026-10-28T19:01:11Z
updated: 2026-10-28T19:01:11Z
---

# I-0221 The close-out's go-test tier runs only the changed packages, so a change that breaks a caller package's tests passes before review

## Description
The close-out's go-test tier runs only the changed packages, so a change that breaks a caller package's tests passes before review

## Instances

### 2026-10-28T19:01:11Z
Story: S-0531.
S-0531 changed flai/internal/storygit/committed.go so that Committed wrote .flai-cache/stories-commits.json, making the folder when it was missing. The close-out passed, since its go-test tier ran only ./internal/storygit and ./internal/statsread. But flai check calls Committed for claims drift, and five tests in flai/cmd (TestCheckCommand, TestCheckStoryPassesOverFindingsOutsideIt, and others) failed on layout.gitignore, because the fixture flai/internal/metrics/testdata/good gained a .flai-cache with no .gitignore. Only a flai test run on flai/cmd/check.go, which took 7 minutes, found it. The tiers left to the merge queue would have found it after acceptance. That fixture sits in the source tree, so the failing run also left a git-ignored .flai-cache in it that failed the next run until it was removed. A tier that also ran the tests of the packages importing a changed package would have found it before review.

## Remediation
