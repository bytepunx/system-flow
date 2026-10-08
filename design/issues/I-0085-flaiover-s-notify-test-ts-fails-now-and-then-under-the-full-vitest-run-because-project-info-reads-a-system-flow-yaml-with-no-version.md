---
id: I-0085
title: flaiover's notify.test.ts fails now and then under the full vitest run because project.info reads a system-flow.yaml with no version
class: defect
status: closed
count: 3
cost: 4m
first_reported: 2026-10-06T07:08:02Z
last_reported: 2026-10-08T06:20:59Z
updated: 2026-10-08T09:06:54Z
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

### 2026-10-08T06:20:59Z
Story: S-0336.
T-1212's `flai test` on flaiover's repo and the new /api/messages route failed once in `flaiover/src/lib/server/notify.test.ts:78` ("inbox webhook posts an entry that appears after it started…: expected null not to be null"), a file the task did not touch, then passed on three runs straight after.

## Remediation

Story S-0288 remediates this issue, created from it at 2026-10-06T09:56:51Z.
Closed 2026-10-08T09:06:54Z: S-0288 fixed it in the test's setup, not in Repo. notify.test.ts's setManifest truncated system-flow.yaml with writeFile and never reported the change, so a manifest read that an earlier repo.changed had started could see the file empty ("version is required") or cache it without notify_url (startNotifier returned null). setManifest now writes a temporary file, renames it over the manifest, and calls repo.changed('system-flow.yaml') as flai on the host does. A new test, "reads the notify_url set after a change was reported and the manifest read", reproduces the stale cached read and failed against the old setManifest; the whole flaiover vitest suite then passed five runs in a row.
