---
id: T-3482
type: task
nature: feature
title: "flai serve starts a state job when a story enters a state with runs: it rebases the branch, runs the state's suites, records each step, and sends the story back by the state's send-back transition when one fails"
status: done
parent: S-0810
owner: alex
created: 2026-11-16T23:51:15Z
updated: 2026-11-18T10:37:07Z
transitions:
  - to: ready
    at: 2026-11-18T09:47:08Z
    by: agent-S-0810
  - to: in-progress
    at: 2026-11-18T09:47:08Z
    by: agent-S-0810
  - to: done
    at: 2026-11-18T10:37:07Z
    by: agent-S-0810
stream: S-0810
tags: [flai]
touches: [flai/internal/serve/statejob.go, flai/internal/serve/statejob_test.go, flai/internal/serve/resolvejob.go, flai/internal/serve/verifyjob.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_restart_test.go, flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go, flai/internal/sendback/sendback.go, flai/internal/sendback/sendback_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-3479]
usage:
  source: log
  seconds: 2999
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 222
      output: 103723
      cache_read: 25252817
      cache_write: 808670
      cost: 11.9462
---
# T-3482 flai serve starts a state job when a story enters a state with runs: it rebases the branch, runs the state's suites, records each step, and sends the story back by the state's send-back transition when one fails

## Work

Waits for T-3479, whose accessors say which state runs what and which transition sends a story back from it. It builds on S-0804's resolve job (`flai/internal/serve/resolvejob.go`), its `closeout` record and progress lines, and its `flai/internal/sendback` package, on main once S-0804 is accepted.

1. A new job kind in `flai/internal/closeout/closeout.go`, `state`, with the state it runs for and one record per step (the rebase, each suite, and the `kluster:` steps T-3484 adds), each with its start, its end, and how it ended: passed, failed, or skipped with its reason.
2. In a new `flai/internal/serve/statejob.go`, a pass `flai serve` runs on each look, wired in `flai/internal/serve/serve.go`: for each story whose last transition entered a state that runs something on entry, start one job, keyed by that transition's time so that it runs once per entry and again on a later entry. A job already running for the story is left to run; one `flai serve`'s restart interrupts is swept as `interrupted` and sends nothing back, as a close-out job is.
3. The job's body, sharing the resolve job's rebase rather than copying it: rebase the branch onto the main branch as `flai stream sync` does, aborting a rebase that stops on conflicts; then run the tiers of the state's suites that the story's paths select, through `flai verify S-nnnn` limited to those suites, so that the verify record keeps the result at the rebased head. Add that limit to `flai verify` (`flai/cmd/verify.go`, `flai/internal/verify/story.go`) with S-0805's suite of each tier, unless S-0805 or S-0807 gave it one; regenerate `docs/users/flai-reference.md` when a flag is added. `all` runs every tier the paths select. Output goes to the job's log.
4. At the end: a job that passed leaves the story in the state. A rebase that stopped or a step that failed sends the story back as `flai` by the send-back transition out of that state, with the send-back thread `flai/internal/sendback` opens naming the state, the step, its findings or the conflicting paths, and the job's log. Read the story again first: one that has left the state is left as it is. A state with no send-back transition records the failure and moves nothing.
5. `flai/internal/serve/agents_restart.go` starts the story's agent on that thread, as it does for the acceptance's send-back, and resolves the thread once the story is back in the state the job ran for.
6. Tests with a fake runner, as `verifyjob_test.go` does: a pass, a rebase conflict, a failing suite, a second entry, a story that left the state, the interrupted sweep, a state with no `runs:` that starts nothing, and a state with no send-back transition.

## Done when

- `flai test` passes on the changed paths.
- Each end above is covered by a test that fails without the change.

## Notes
