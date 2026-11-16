---
id: T-3393
type: task
nature: remediation
title: Show the held-start container tests pass under load in the race run, and close I-0369
status: done
parent: S-0795
owner: alex
created: 2026-11-16T00:41:11Z
updated: 2026-11-16T01:23:49Z
transitions:
  - to: ready
    at: 2026-11-16T00:52:37Z
    by: agent-S-0795
  - to: in-progress
    at: 2026-11-16T00:52:39Z
    by: agent-S-0795
  - to: done
    at: 2026-11-16T01:23:49Z
    by: agent-S-0795
stream: S-0795
tags: [flai, serve, test]
touches: [design/issues/I-0369-testcontainerstrategicstartunderalayerwaitsfortheharnessimagealone-fails-under.md]
after: [T-3392]
usage:
  source: log
  seconds: 1870
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 7925
      cache_read: 4393864
      cache_write: 84705
      cost: 1.7151
---
# T-3393 Show the held-start container tests pass under load in the race run, and close I-0369

## Work

Waits for T-3392: it checks that fix and closes the issue on it.

Run the image-held container tests in `flai/internal/serve/container_test.go` under `-race` and under load: `go test -race -count=50 -run 'TestContainer.*(Held|Image|Layer)' ./internal/serve/` with the package's other tests running beside it, or with `-cpu` raised. Then run `flai test race` once. I-0369's two instances failed only under load, so a quiet run proves nothing: record the commands, counts, and results in this task's Notes. If the window still shows, or the check at line 1616 should now read the state after the planner ends, change the test only so that it checks the same thing without racing, and say why in the narrative.

Then close the issue: `flai issue close I-0369 --reason "<what fixed it>"`, naming T-3392's change and the runs that show it.

## Done when

- The image-held container tests pass in every run of the loaded `-race` loop, and `flai test race` passes, each recorded in Notes.
- I-0369 is closed with a reason that names the fix.
- Both of S-0795's criteria can be ticked.

## Notes

Runs of 2026-11-16, on the 24-CPU host, with test binaries built by `go test -race -c` in `flai/`: `serve.fix.test` from the branch at 587b9d4, and `serve.main.test` with `-overlay` putting back main's `agents_state.go` and `agents_start.go`. Each ran from `flai/internal/serve` with `FLAI_STORY` and `FLAI_STARTED_BY` unset, `-test.skip '^TestGateway'`, pinned with `taskset` beside busy loops pinned to the same CPUs (D-0492).

| Binary | Tests | Load | Count | Result |
|--------|-------|------|-------|--------|
| main's code | the held-start group | CPUs 20-23, `GOMAXPROCS=2`, 8 loops | 20 | passed, 64 s: this load does not reproduce it |
| main's code | `TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone` | CPU 23, `GOMAXPROCS=1`, 40 loops | 100 | 22 failed at `container_test.go:1618`, "the held run is kept as a past run", 740 s: I-0369 reproduced |
| the fix | `TestContainerStrategicStartUnderALayerWaitsForTheHarnessImageAlone` | CPU 22, `GOMAXPROCS=1`, 40 loops | 100 | all passed, 747 s |
| the fix | the held-start group, 16 tests (below) | CPUs 20-23, `GOMAXPROCS=2`, 8 loops, with `flai test race` running beside it | 50 | all 800 passed, 261 s |

The held-start group is `-test.run 'TestContainer.*(Held|Image|Layer)'`.

`scripts/flai.sh test race` in the worktree, while the group loop ran: `passed race (4m54.3s)`, exit 0.

The test was not changed: its check at line 1616 reads what it should, and passes once the replacement is one write.
