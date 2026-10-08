---
id: I-0092
title: Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite
class: efficiency
status: closed
count: 3
cost: 7m
first_reported: 2026-10-06T19:46:45Z
last_reported: 2026-10-07T09:16:07Z
updated: 2026-10-08T09:49:58Z
---

# I-0092 Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite

## Description
Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite

## Instances

### 2026-10-06T19:46:45Z
Story: S-0278.
S-0292's close-out and S-0221's both bumped I-0078 to 9, so S-0221's sync before review stopped on I-0078's front matter and its agent merged it by hand to count 10 with both instances (I-0074's fifth instance); S-0257 and S-0262 met the same in I-0073 (I-0074's third). The operator asked on TH-0173 for a story for it: S-0278 regenerates only design/issues/summary.md, which is derived, while an issue's count and instances are data both branches add to.

### 2026-10-07T00:54:07Z
Story: S-0254.
S-0254's close-out bumped I-0073 and I-0078 while S-0228's acceptance bumped them on main; the final sync stopped on both, and on summary.md, until the two instances were merged by hand and the summary regenerated

### 2026-10-07T09:16:07Z
Story: S-0275.
S-0275 bumped I-0079, I-0058, and I-0111 while S-0213 and S-0246 bumped or closed the same issues; each trial merge conflicted, and S-0275 dropped its bumps, keeping the occurrences in its narrative.

## Remediation

Story S-0297 remediates this issue, created from it at 2026-10-06T19:46:46Z.
Closed 2026-10-08T09:49:58Z: S-0326 (ADR-0126): the rebase in flai stream sync, flai task done, and flai accept merges an issue file both sides changed by its instances, adding up the count and taking the later last_reported and updated, and the trial merge leaves issue files out. S-0297 confirmed it covers every instance; tests in flai/cmd/stream_sync_test.go reproduce each: TestSyncMergesAnIssueBothSidesBumped and TestSyncTrialMergeLeavesOutIssueFiles (S-0278, S-0221), TestAcceptMergesAnIssueBothSidesBumped (S-0254, S-0228), and TestSyncMergesABumpIntoTheIssueMainClosed (S-0275).
