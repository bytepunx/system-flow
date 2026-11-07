---
id: I-0247
title: TestServeSaysTheFlaiIsOlderThanTheProject and TestServeStartsADefinedRoleOnAnEventAndOnASchedule fail the close-out's integration tier on timing, on a branch that does not touch flai/internal/serve
class: defect
status: closed
count: 5
cost: 8m
first_reported: 2026-11-02T01:00:19Z
last_reported: 2026-11-07T02:26:34Z
updated: 2026-11-07T05:27:19Z
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

### 2026-11-07T02:10:46Z
Story: S-0578.
S-0578, which touches flai/internal/kluster and flai/cmd/kluster.go only, failed two close-outs at the integration tier on timing tests outside it: at 9ae5a92, TestHostRunsServeAndTheMCPServersAndTakesThemWithIt (flai/cmd/host_integration_test.go:112, "never the MCP server's own state"); at 0401e58, TestServeStartsADefinedRoleOnAnEventAndOnASchedule (flai/internal/serve/serve_test.go:1009, "never happened: the watcher's run for the move ended"). S-0580's agent ran its own work on the host at the same time.

### 2026-11-07T02:26:34Z
Story: S-0578.
S-0578's third and fourth close-outs at 0264cef failed the integration tier on timing tests outside the story again: TestServeStartsADefinedRoleOnAnEventAndOnASchedule (serve_test.go:1009) at a load average of 64, then, once the load fell to about 23 at the start, flai/internal/mcpserver/folder_test.go:362 TestWaitForEventsAcrossAFolder ("timed_out:true"), with the load back near 37 by its end. Four close-outs, four different or repeated timing failures, none in a package the story changes.

## Remediation

Story S-0569 remediates this issue, created from it at 2026-11-02T21:41:39Z.
Closed 2026-11-07T05:27:19Z: fixed by S-0569. Each test raced its own setup; the load only widened the window. Each fix waits for the event it depends on.
- TestWaitForEventsAcrossAFolder and its two wait_for_work siblings slept 150 ms, then changed the project, before the wait had fingerprinted it. They now wait for the wait's first progress heartbeat and hold up to 30 s (T-2522).
- TestServeSaysTheFlaiIsOlderThanTheProject read the watcher's change notification for serve's upgrade thread as project.info's reply, because channeltest's Ask returns the next message. It now asks until an answer names the project (T-2525).
- TestServeStartsADefinedRoleOnAnEventAndOnASchedule moved the story before serve's watcher took its first snapshot, so only the next minute's look saw the move. It now waits for the watcher to report a probe file (T-2528; serve's window is I-0273).
- TestHostRunsServeAndTheMCPServersAndTakesThemWithIt lost the first free MCP port to concurrent runs, so the state came with another PID. It now seeds its own kernel-picked port (T-2529; the host's same-port restart is I-0271).
Each fix reproduced its failure before it changed anything, under GOMAXPROCS=2 -race beside busy loops pinned to the test's CPUs. All six tests then passed 60 of 60 together, pinned to CPUs 20-23 beside 6 busy loops.
