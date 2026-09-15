---
id: S-027
type: story
nature: feature
title: flai issue commands and check rules for design/issues
status: backlog
parent: E-005
owner: alex
created: 2026-09-15T22:37:54Z
updated: 2026-09-15T22:37:54Z
transitions: []
tags: [conventions, cli]
---

# S-027 flai issue commands and check rules for design/issues

## Goal
The continuous-improvement convention is cheap to follow: flai records and increments issues, regenerates summary.md, and check validates the folder, so the table never drifts.

## Acceptance criteria
- [ ] `flai issue new "<title>" --class defect|blocker|efficiency|impression [--cost 20m]` creates `design/issues/I-nnn-slug.md` from the schema with count 1 and both timestamps
- [ ] `flai issue bump I-nnn [--cost] [--note "..."]` increments count, sets last_reported, averages cost, appends the instance to the body
- [ ] `flai issue summary` regenerates summary.md; new and bump call it
- [ ] `flai issue close I-nnn --reason` sets status closed and removes it from the summary
- [ ] check rules: issue front matter schema, summary matches open issues, IDs unique
- [ ] flai prime (S-025) includes the summary table when it is non-empty

## Tasks

## Notes
