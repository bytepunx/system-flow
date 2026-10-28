---
id: T-2201
type: task
nature: remediation
title: The curate scheduler holds its triggers while the last run's usage limit has not reset, and starts the curator again once it has
status: backlog
parent: S-0514
owner: alex
created: 2026-10-28T12:53:51Z
updated: 2026-10-28T12:53:51Z
transitions: []
stream: S-0514
tags: [flai, serve, curator]
touches: [flai/internal/serve/curate_schedule.go, flai/internal/serve/curate_schedule_test.go, flai/internal/serve/curate.go]
after: [T-2200]
---
# T-2201 The curate scheduler holds its triggers while the last run's usage limit has not reset, and starts the curator again once it has

## Work

Waits for T-2200, which puts the reset time and the retry mark on the curator's run.

In `flai/internal/serve/curate_schedule.go`, have `look` read the project's curator run from `s.o.Dir.AgentStates()[s.e.Root].Curator` before it chooses a trigger:

- While that run ended on the limit (`failed`, with a reset time) and `now` is before the reset time, start nothing: keep `pending` as it is (acceptances seen meanwhile are added to it as today), and let a due schedule time mark the cadence as seen but start nothing, so it is folded into the retry rather than coming round again.
- At the first look at or after the reset time, when no curator runs and the run is not yet retried (no later run exists, and the run is not itself a retry that ended on the limit again; T-2202 raises that), start one run: its trigger is the retried run's trigger with the held acceptances joined in ID order (`accepted S-0498, S-0501, S-0503`; a retried `asked` with no acceptances stays `asked`; a retried `schedule daily` stays so), through `curateCheck` as every start, and `pending` is cleared. Mark the run as a retry of the limited one through `launcher.curate` (`curate.go`), which takes the mark and writes it on the record and in its journal entry's detail (`started again after the usage limit reset at ...`), so `flai serve journal` shows it.
- Because the state is the run in `serve/agents.json`, a `flai serve` started after the reset retries at its first look, even though the scheduler's own memory is new; while the `curate` action is off, nothing is retried, as nothing is started.

Update the file's and `look`'s comments, which say a refused start is not tried again, to say what is tried again.

Tests in `curate_schedule_test.go`, with `newCurationLab`, `lookAt`, `accept`, and `ended`: a run that ended on the limit at 16:03 with a reset at 08:00 the next day, looks before the reset start nothing although a story is accepted and the schedule comes round; the first look after the reset starts one run whose trigger names the retried run's stories and the accepted one, marked a retry; a run that ended failed without a limit is not retried, as today; a scheduler made after the reset, from the run on disk, retries at its first look; with `curate` off nothing is retried.

## Done when

- A curator run that ended on the usage limit is started again once at the first look after its reset time, with the acceptances held meanwhile in its trigger, marked a retry, and journalled.
- Nothing starts the curator between the limit and its reset, and a failure without a limit is not retried.
- `flai test` on the task's paths passes.

## Notes

Layer 3 of S-0514's plan. `curate.go` is also T-2200's, in the layer before, and T-2202's, in the layer after.
