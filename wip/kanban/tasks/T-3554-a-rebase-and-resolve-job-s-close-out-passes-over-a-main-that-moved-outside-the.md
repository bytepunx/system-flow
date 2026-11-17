---
id: T-3554
type: task
nature: improvement
title: A rebase and resolve job's close-out passes over a main that moved outside the story's change set, in both scripts and the container close-out
status: backlog
parent: S-0824
owner: alex
created: 2026-11-17T08:33:52Z
updated: 2026-11-17T08:33:52Z
transitions: []
stream: S-0824
tags: [flai, serve, template]
touches: [scripts/close-out.sh, template/root/scripts/close-out.sh, template/CHANGELOG.md, flai/internal/closeout/script_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_container.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/resolvejob.go, flai/internal/serve/resolvejob_test.go]
after: [T-3552, T-3553]
---
# T-3554 A rebase and resolve job's close-out passes over a main that moved outside the story's change set, in both scripts and the container close-out

## Work

Fix the cause I-0398 describes. This task waits for two others:

- T-3552, for the `flai verify --sync-only` flag the last check passes;
- T-3553, for the round line's words, which `closeout.Status` reads.

The resolve job asks for the close-out the ADR names. `startJob` and `command` in `flai/internal/serve/verifyjob.go` build the close-out a resolve job runs. For a job of kind `resolve` only, they set the ADR's variable, for example `CLOSE_OUT_RESOLVE=1`. `closeOutEnv` drops it from the host's environment, as it drops `CLOSE_OUT_FRESH`, so a close-out job never inherits it.

The scripts:

- `scripts/close-out.sh` and `template/root/scripts/close-out.sh` read the variable, unset it before the tiers run, and pass the flag to the last check. When the last check passes over a moved main, the close-out ends in that round. Its last line says main moved and that nothing it brought in sends back a tier.
- When the check stops on a tier sent back, `sync_round` prints the round line with a why that names the tiers and the first overlapping paths, as T-3553 reads it.
- Note the change in `template/CHANGELOG.md`, with the flai version the scripts' flag needs.

The container close-out: `verifyjob_container.go` runs its own three-round loop (`closeOutRounds`). Give it the same last check and the same round line for a resolve job.

Tests that reproduce I-0398:

- in `flai/internal/closeout/script_test.go`, against both scripts and the stub flai:
  - with the variable, a last check that passes over a moved main ends in round 1;
  - without it, the same stub's stop at sync runs round 2, as today;
  - a stop on a tier sent back runs round 2, and its round line names the tier;
- in `flai/internal/serve/resolvejob_test.go`, a resolve job's close-out gets the variable and a close-out job's does not;
- in `flai/internal/serve/verifyjob_test.go`, the container loop ends in round 1 when its last check passes over a moved main.

## Done when

- A resolve job's close-out whose main moved without sending back a tier ends in the round it verified, on the host and in a container.
- A close-out job, and the agent's close-out before review, still sync and verify again as ADR-0156 says.
- Both scripts behave alike, and the tests pass with `flai test` on the changed paths.

## Notes
