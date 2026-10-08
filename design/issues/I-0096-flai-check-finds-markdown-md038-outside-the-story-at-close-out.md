---
id: I-0096
title: "flai check finds `markdown.MD038` outside the story at close-out"
class: efficiency
status: closed
count: 1
first_reported: 2026-10-06T22:10:04Z
last_reported: 2026-10-06T22:10:04Z
updated: 2026-10-08T04:45:26Z
---

# I-0096 flai check finds `markdown.MD038` outside the story at close-out

## Description
flai check finds `markdown.MD038` outside the story at close-out

## Instances

### 2026-10-06T22:10:04Z
Story: S-0227.
flai check found outside the story:
`wip/agents/S-0229.md`: MD038/no-space-in-code Spaces inside code span elements [Context: "'rule: '"]

## Remediation

Story S-0318 remediates this issue, created from it at 2026-10-07T18:59:49Z.
Closed 2026-10-08T04:45:26Z: Fixed by S-0318. The cause: S-0229's agent hand-edited its narrative's ## Decisions with a code span ending in a space, which broke MD038, and S-0227's close-out recorded it, though only S-0229 could fix it. ADR-0123 decides that a check scoped to a story leaves out a markdown finding on the narrative of another story that is neither done nor cancelled, so a close-out records none; that story's own close-out still finds it, and the unscoped check still warns. TestScopedCheckLeavesOutAMarkdownFindingOnAnotherOpenStorysNarrative and TestScopedCheckNotesEveryOtherMarkdownFindingOutsideTheStory in flai/internal/check/scope_test.go, and TestCheckRecordIssuesLeavesOutAMarkdownFindingOnAnotherOpenStorysNarrative in flai/cmd/check_test.go, reproduce it.
