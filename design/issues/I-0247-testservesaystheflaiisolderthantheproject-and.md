---
id: I-0247
title: TestServeSaysTheFlaiIsOlderThanTheProject and TestServeStartsADefinedRoleOnAnEventAndOnASchedule fail the close-out's integration tier on timing, on a branch that does not touch flai/internal/serve
class: defect
status: open
count: 3
cost: 7m
first_reported: 2026-11-02T01:00:19Z
last_reported: 2026-11-02T04:06:44Z
updated: 2026-11-02T21:41:39Z
---

# I-0247 TestServeSaysTheFlaiIsOlderThanTheProject and TestServeStartsADefinedRoleOnAnEventAndOnASchedule fail the close-out's integration tier on timing, on a branch that does not touch flai/internal/serve

## Description
TestServeSaysTheFlaiIsOlderThanTheProject and TestServeStartsADefinedRoleOnAnEventAndOnASchedule fail the close-out's integration tier on timing, on a branch that does not touch flai/internal/serve

## Instances

### 2026-11-02T01:00:19Z
Story: S-0518.
S-0518's close-out at 2026-11-02T00:48Z ran the integration tier (7m56s), selected because `template/root/.claude/settings.json` matches its `when` of `template/**`, and failed on two `flai/internal/serve` tests the story does not touch: `outdated_test.go:175 TestServeSaysTheFlaiIsOlderThanTheProject: project.info: <nil>` and `serve_test.go:1009 TestServeStartsADefinedRoleOnAnEventAndOnASchedule: never happened: the watcher's run for the move ended`. Neither skips under `-short`, and `flai test` on the two files passed at once after, in `./internal/serve`'s own run (1m38s). The host was running four other stories' agents and their tests at the time, so both waits look too short for a loaded host, like I-0226's two serve tests.

### 2026-11-02T04:02:11Z
Story: S-0518.
S-0518's close-out at 2026-11-02T03:58Z, after the T-2425 changelog fix, failed the integration tier (3m20s) again on `serve_test.go:1009 TestServeStartsADefinedRoleOnAnEventAndOnASchedule: never happened: the watcher's run for the move ended`. The branch changes nothing under `flai/internal/serve`; it was selected because `template/root/.claude/settings.json` matches `template/**`.

### 2026-11-02T04:06:44Z
Story: S-0518.
The next close-out of S-0518, at 2026-11-02T04:03Z, failed the integration tier (4m0s) on a third test outside the story, `flai/internal/mcpserver/folder_test.go:371 TestWaitForEventsAcrossAFolder: changed paths are relative to the folder: []` (I-0186's test, closed). The host's load average was 76.81 (1 min), 59.37 (5 min), 40.60 (15 min) right after. The waits in these tests do not hold on a host this loaded, whichever test meets them.

## Remediation

Story S-0569 remediates this issue, created from it at 2026-11-02T21:41:39Z.
