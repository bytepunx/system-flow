---
id: I-0319
title: A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration
class: defect
status: closed
count: 2
cost: 11m
first_reported: 2026-11-09T10:34:28Z
last_reported: 2026-11-09T17:32:35Z
updated: 2026-11-10T05:19:50Z
---

# I-0319 A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration

## Description
A flaiover props-panel test times out at 5,000 ms under load and fails the flaiover tier, serial or beside integration

## Instances

### 2026-11-09T10:34:28Z
Story: S-0610.
S-0610's timing runs (flai test --all, CLOSE_OUT_STORY=S-0610, 2026-11-09) failed four times on one vitest test, `flaiover/src/.../components.svelte.test.ts` › the props panel and the event log › "lists the calls of the component's function props, and clears them on another pick", which timed out at vitest's 5,000 ms. Two failures came in serial runs with the installed flai at loads of 91 to 102, two in parallel runs at 86 to 125. It touches no flai and posts about 60 messages to the page, so it fails under host load whatever runs beside it. Each failed run cost a 7 to 15 minute rerun.

### 2026-11-09T17:32:35Z
Story: S-0611.
S-0611's first close-out, job S-0611-20261109T172632Z at 6e33bd405, stopped at the flaiover tier after 5m44s. It failed only on components.svelte.test.ts › the props panel and the event log › "lists the calls of the component's function props, and clears them on another pick", which timed out at 5,000 ms; 2,422 of 2,423 tests passed. S-0611 changes no file under the Components page. The one-minute load average was about 50 to 62. The close-out was asked for again.

## Remediation

Story S-0640 remediates this issue, created from it at 2026-11-09T16:55:57Z.
Closed 2026-11-10T05:19:50Z: S-0640 (T-2805, 81c24831f): the test settled after each of the 55 harness messages it posted, about 1,240 setTimeout(0) turns in one test, so a loaded host pushed it past 5,000 ms. It now posts the burst without settling and settles once, and checks all 50 entries land whole. It takes 672 to 682 ms, down from 2,174 to 2,447 ms, at a load of 12 to 19.
