---
id: T-0483
type: task
nature: feature
title: The work-management convention says to sync and re-run the tests on an overlap notice
status: done
parent: S-0132
owner: alex
created: 2026-09-26T18:17:27Z
updated: 2026-09-26T18:23:06Z
transitions:
  - to: ready
    at: 2026-09-26T18:17:41Z
    by: agent-S-0132
  - to: in-progress
    at: 2026-09-26T18:21:59Z
    by: agent-S-0132
  - to: done
    at: 2026-09-26T18:23:06Z
    by: agent-S-0132
stream: S-0132
tags: []
touches: [template/root/design/conventions, design/conventions]
---
# T-0483 The work-management convention says to sync and re-run the tests on an overlap notice

## Work

Add the rule to the baseline in `template/root/design/conventions/work-management.md`, copy it above the marker in `design/conventions/work-management.md`, and bump the template's version and changelog as template changes do.

## Done when

- Both files carry the rule; `make smoke` (template render and `flai check --strict`) passes.

## Notes
