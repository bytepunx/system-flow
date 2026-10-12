---
id: I-0129
title: An empty narrative file under wip/agents makes every flai stream log in the project fail
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-08T20:55:33Z
last_reported: 2026-10-08T20:55:33Z
updated: 2026-10-12T16:52:52Z
---

# I-0129 An empty narrative file under wip/agents makes every flai stream log in the project fail

## Description
An empty narrative file under wip/agents makes every flai stream log in the project fail

## Instances

### 2026-10-08T20:55:33Z
Story: S-0232.
At 2026-10-08T20:55Z, `flai stream log S-0232` failed twice with `wip/agents/S-0337.md: no front matter`. That file is empty (0 bytes, untracked, dated 05:04) and is not S-0232's. Writing one story's narrative log should not depend on another story's narrative parsing: flai could skip or warn on a narrative it cannot read when it rewrites wip/agents/index.md, rather than refusing the write. S-0232 left the file alone, as it is not its to delete.

## Remediation

Story S-0373 remediates this issue, created from it at 2026-10-12T16:52:52Z.
