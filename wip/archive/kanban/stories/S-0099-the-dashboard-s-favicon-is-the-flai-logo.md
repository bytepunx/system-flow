---
id: S-0099
type: story
nature: improvement
title: The dashboard's favicon is the flai logo
status: done
owner: alex
created: 2026-09-23T04:57:07Z
updated: 2026-09-23T05:04:08Z
transitions:
  - to: ready
    at: 2026-09-23T04:57:15Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T04:57:15Z
    by: system-flow
  - to: review
    at: 2026-09-23T05:00:11Z
    by: system-flow
  - to: done
    at: 2026-09-23T05:04:08Z
    by: alex
tags: []
touches: [flaiover/src]
---

# S-0099 The dashboard's favicon is the flai logo

## Goal

The dashboard's browser tab shows the flai logo the operator supplied (a 128×128 PNG: "flai" in blue on black) instead of the SvelteKit default.

## Acceptance criteria
- [x] The dashboard's favicon is the operator's flai logo
- [x] The old favicon is removed

## Tasks
- T-0354 Replace the favicon with the flai logo

## Notes
Asked for by the operator in conversation on 2026-09-23, with the image attached.
