---
id: I-0085
title: flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version
class: defect
status: open
count: 2
cost: 5m
first_reported: 2026-10-06T07:08:02Z
last_reported: 2026-10-06T19:28:59Z
updated: 2026-10-06T19:28:59Z
---

# I-0085 flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version

## Description
flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version

## Instances

### 2026-10-06T07:08:02Z
Story: S-0220.
S-0220's second close-out stopped at vitest on one test, flaiover/src/lib/server/notify.test.ts 'inbox webhook > posts an entry that appears after it started, once, and never what was already there': bin/flai hostapi project.info said /tmp/flaiover-notify-*/system-flow.yaml: version is required. The first close-out, minutes earlier on the same tests, passed it (915), and the file alone passes 3 of 3. Likely cause, not proven: setManifest rewrites system-flow.yaml with writeFile, which truncates it first, so a project.info asked at that moment by a notifier still polling reads a manifest with no version. S-0220 changes nothing there.

### 2026-10-06T19:28:59Z
Story: S-0295.
S-0295's close-out stopped at the flaiover step on notify.test.ts:78 (expected null not to be null); it passed alone and the whole flaiover suite passed on the next run

## Remediation

Story S-0288 remediates this issue, created from it at 2026-10-06T09:56:51Z.
