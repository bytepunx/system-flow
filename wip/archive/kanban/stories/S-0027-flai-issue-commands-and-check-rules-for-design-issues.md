---
id: S-0027
type: story
nature: feature
title: flai issue commands and check rules for design/issues
status: done
parent: E-0005
owner: alex
created: 2026-09-15T22:37:54Z
updated: 2026-09-17T00:04:08Z
transitions:
  - to: ready
    at: 2026-09-16T23:54:25Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:54:25Z
    by: agent
  - to: review
    at: 2026-09-16T23:59:09Z
    by: agent
  - to: done
    at: 2026-09-17T00:04:08Z
    by: alex
tags: [conventions, cli]
---

# S-0027 flai issue commands and check rules for design/issues

## Goal
The continuous-improvement convention is cheap to follow: flai records and increments issues, regenerates summary.md, and check validates the folder, so the table never drifts.

## Acceptance criteria
- [x] `flai issue new "<title>" --class defect|blocker|efficiency|impression [--cost 20m]` creates `design/issues/I-nnn-slug.md` from the schema with count 1 and both timestamps
- [x] `flai issue bump I-nnn [--cost] [--note "..."]` increments count, sets last_reported, averages cost, appends the instance to the body
- [x] `flai issue summary` regenerates summary.md with count, average cost, and total cost (count times average), most expensive first; new and bump call it
- [x] `flai issue close I-nnn --reason` sets status closed and removes it from the summary
- [x] check rules: issue front matter schema, summary matches open issues, IDs unique
- [x] flai prime (S-0025) includes the summary table when it is non-empty

## Tasks
- T-0060 issues package: parse, write, list, IDs, new, bump, close, summary
- T-0061 flai issue commands and prime integration
- T-0062 check rules for design/issues, fixtures, docs

## Notes
