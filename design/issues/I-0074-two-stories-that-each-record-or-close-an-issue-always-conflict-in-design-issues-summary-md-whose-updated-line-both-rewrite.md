---
id: I-0074
title: Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite
class: efficiency
status: open
count: 1
cost: 3m
first_reported: 2026-10-05T03:22:09Z
last_reported: 2026-10-05T03:22:09Z
updated: 2026-10-05T03:22:09Z
---

# I-0074 Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite

## Description
Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite

## Instances

### 2026-10-05T03:22:09Z
Story: S-0260.
S-0260 closed I-0071 and S-0253 closed I-0066, each with flai issue close on its own branch. flai stream sync's trial merge then opened TH-0126: both changed summary.md's updated timestamp, and the two removed rows sit two lines apart, so the hunks meet. The resolution is mechanical (keep both removals, the later timestamp, or regenerate with flai issue summary), but whichever story is accepted second stops on it, and each story's agent spends a turn reading the thread. summary.md is generated: a rebase or merge could regenerate it rather than merge it line by line, or the trial merge could leave it out.

## Remediation
