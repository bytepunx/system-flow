---
id: I-0319
title: A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration
class: defect
status: open
count: 1
cost: 15m
first_reported: 2026-11-09T10:34:28Z
last_reported: 2026-11-09T10:34:28Z
updated: 2026-11-09T16:55:57Z
---

# I-0319 A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration

## Description
A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration

## Instances

### 2026-11-09T10:34:28Z
Story: S-0610.
S-0610's timing runs (flai test --all, CLOSE_OUT_STORY=S-0610, 2026-11-09) failed four times on one vitest test, `flaiover/src/.../components.svelte.test.ts` › the props panel and the event log › "lists the calls of the component's function props, and clears them on another pick", which timed out at vitest's 5,000 ms. Two failures came in serial runs with the installed flai at loads of 91 to 102, two in parallel runs at 86 to 125. It touches no flai and posts about 60 messages to the page, so it fails under host load whatever runs beside it. Each failed run cost a 7 to 15 minute rerun.

## Remediation

Story S-0640 remediates this issue, created from it at 2026-11-09T16:55:57Z.
