---
id: I-0092
title: Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-06T19:46:45Z
last_reported: 2026-10-06T19:46:45Z
updated: 2026-10-06T19:46:46Z
---

# I-0092 Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite

## Description
Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite

## Instances

### 2026-10-06T19:46:45Z
Story: S-0278.
S-0292's close-out and S-0221's both bumped I-0078 to 9, so S-0221's sync before review stopped on I-0078's front matter and its agent merged it by hand to count 10 with both instances (I-0074's fifth instance); S-0257 and S-0262 met the same in I-0073 (I-0074's third). The operator asked on TH-0173 for a story for it: S-0278 regenerates only design/issues/summary.md, which is derived, while an issue's count and instances are data both branches add to.

## Remediation

Story S-0297 remediates this issue, created from it at 2026-10-06T19:46:46Z.
