---
id: I-0085
title: flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-06T07:08:02Z
last_reported: 2026-10-06T07:08:02Z
updated: 2026-10-06T07:08:02Z
---

# I-0085 flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version

## Description
flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version

## Instances

### 2026-10-06T07:08:02Z
Story: S-0220.
S-0220's second close-out stopped at vitest on one test, flaiover/src/lib/server/notify.test.ts 'inbox webhook > posts an entry that appears after it started, once, and never what was already there': bin/flai hostapi project.info said /tmp/flaiover-notify-*/system-flow.yaml: version is required. The first close-out, minutes earlier on the same tests, passed it (915), and the file alone passes 3 of 3. Likely cause, not proven: setManifest rewrites system-flow.yaml with writeFile, which truncates it first, so a project.info asked at that moment by a notifier still polling reads a manifest with no version. S-0220 changes nothing there.

## Remediation
