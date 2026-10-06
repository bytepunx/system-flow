---
id: I-0074
title: Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite
class: efficiency
status: closed
count: 5
cost: 3m
first_reported: 2026-10-05T03:22:09Z
last_reported: 2026-10-06T11:45:29Z
updated: 2026-10-06T20:03:27Z
---

# I-0074 Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite

## Description
Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite

## Instances

### 2026-10-05T03:22:09Z
Story: S-0260.
S-0260 closed I-0071 and S-0253 closed I-0066, each with flai issue close on its own branch. flai stream sync's trial merge then opened TH-0126: both changed summary.md's updated timestamp, and the two removed rows sit two lines apart, so the hunks meet. The resolution is mechanical (keep both removals, the later timestamp, or regenerate with flai issue summary), but whichever story is accepted second stops on it, and each story's agent spends a turn reading the thread. summary.md is generated: a rebase or merge could regenerate it rather than merge it line by line, or the trial merge could leave it out.

### 2026-10-05T03:48:36Z
Story: S-0253.
Syncing S-0253 onto main after S-0260's acceptance stopped on design/issues/summary.md: both rewrote its updated line, S-0253 closing I-0066 and recording the close-out's findings.

### 2026-10-05T04:38:59Z
Story: S-0257.
S-0257 closed I-0069 and S-0262 closed I-0072, so flai stream sync's trial merge conflicted in design/issues/summary.md (TH-0130). Each close-out then recorded a threads.archived instance in I-0073, which added a second conflict, in I-0073 itself.

### 2026-10-06T11:36:35Z
Story: S-0292.
S-0292 closed I-0090 and S-0221 bumped I-0058; both rewrote design/issues/summary.md's updated line, so flai stream sync S-0292's trial merge reported a conflict with story/S-0221 and opened TH-0175.

### 2026-10-06T11:45:29Z
Story: S-0221.
The sync before review stopped on design/issues/summary.md (I-0058's row and the updated line) and then on I-0078, which S-0292's close-out and S-0221's had both bumped to 9; merged by hand to 10 with both instances.

## Remediation

Story S-0278 remediates this issue, created from it at 2026-10-05T04:40:47Z.
Closed 2026-10-06T20:03:27Z: S-0278 (ADR-0098): flai stream sync and flai accept regenerate design/issues/summary.md and continue the rebase when it is the only conflict, and the trial merge between open story branches leaves it out, so no thread is opened over it. Two branches bumping the same issue is I-0092.
