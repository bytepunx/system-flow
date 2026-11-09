---
id: T-2724
type: task
nature: remediation
title: Show the outdated serve tests hold under load since S-0569, and fix what still fails
status: done
parent: S-0617
owner: alex
created: 2026-11-09T10:34:49Z
updated: 2026-11-09T10:47:10Z
transitions:
  - to: ready
    at: 2026-11-09T10:40:23Z
    by: agent-S-0617
  - to: in-progress
    at: 2026-11-09T10:40:25Z
    by: agent-S-0617
  - to: done
    at: 2026-11-09T10:47:10Z
    by: agent-S-0617
stream: S-0617
tags: [flai, testing]
touches: [flai/internal/serve/outdated_test.go]
usage:
  source: log
  seconds: 405
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 3733
      cache_read: 3044154
      cache_write: 48913
      cost: 1.0677
---
# T-2724 Show the outdated serve tests hold under load since S-0569, and fix what still fails

## Work

I-0252's one instance, S-0506 on 2026-11-02, failed at the old `outdated_test.go:175`, `t.Errorf("project.info: %v %s", rerr, res)`, with `project.info: <nil>`: no error and no result. That is I-0247's symptom: `channeltest`'s `Ask` returned the watcher's change notification, which came before the reply on a loaded host. S-0569 (commit `ea5280843`, 2026-11-06) fixed it after the instance: `askProjectInfo` asks until an answer names the project. `project.info` works out `flai_outdated` itself, at each call (`flai/internal/hostapi/hostapi.go`), so it does not wait on serve's own look, and the issue's reading, "asked before flai serve has worked out the outdated flai", does not hold.

- Run the file's tests many times on a loaded host, for example `go test ./internal/serve -run 'TestServe(SaysTheFlaiIsOlderThanTheProject|AsksItsHostToUpgradeWithAutoUpgradeOn|AsksTheOperatorWithAutoUpgradeOff)' -count=50 -race` from `flai/` while another `go test ./...` runs beside it. Record the command, the load, and the result in this task's Notes.
- If a run fails, read its output, find its cause, and fix it in `flai/internal/serve/outdated_test.go`, or say on a thread on the story which file the fix needs when it is outside this task's touches.
- If none fails, change nothing here: the Notes are the evidence for I-0252's close.

Waits for nothing: it runs the serve tests, which do not depend on the release package's change.

## Done when

- The Notes record the stress run's command, its load, and its result.
- Every run passes, with any fix the runs called for made and its cause named in the Notes.

## Notes

Stress run on 2026-11-09 at story/S-0617's base `12446dd49`, which holds S-0569's fix `ea5280843`, in a worktree of its own (`task/T-2724`) so that T-2723's edit could not reach the build.

- Tests: `TestServeSaysTheFlaiIsOlderThanTheProject`, `TestServeAsksItsHostToUpgradeWithAutoUpgradeOn`, `TestServeAsksTheOperatorWithAutoUpgradeOff`.
- Build, from `flai/`: `GOMAXPROCS=2 go test -race -c -o /tmp/serve-T2724.test ./internal/serve`.
- Run, from `flai/internal/serve`: `GOMAXPROCS=2 taskset -c 20,21 /tmp/serve-T2724.test -test.count=50 -test.timeout=40m -test.run '^TestServe(SaysTheFlaiIsOlderThanTheProject|AsksItsHostToUpgradeWithAutoUpgradeOn|AsksTheOperatorWithAutoUpgradeOff)$' -test.skip '^TestGateway'`.
- Run 1: 4 busy loops pinned to CPUs 20-21 beside it; host of 24 CPUs at load average 85.6 to 80.0. Passed, 51 s.
- Run 2: the same with `-test.v` and 6 busy loops on CPUs 20-21; load average 80.0 to 74.7. Each test passed 50 times, 61 s.

100 runs of each test on 2 loaded CPUs, none failed, so `outdated_test.go` is unchanged. I-0252's symptom, `project.info: <nil>`, did not recur after S-0569's fix.
