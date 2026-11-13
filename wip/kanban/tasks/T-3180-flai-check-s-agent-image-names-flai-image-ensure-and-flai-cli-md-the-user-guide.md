---
id: T-3180
type: task
nature: feature
title: flai check's agent.image names flai image ensure, and flai-cli.md, the user guide, and the runbook describe the image job
status: in-progress
parent: S-0751
owner: alex
created: 2026-11-12T16:02:48Z
updated: 2026-11-13T15:43:10Z
transitions:
  - to: ready
    at: 2026-11-13T15:43:10Z
    by: agent-S-0751
  - to: in-progress
    at: 2026-11-13T15:43:10Z
    by: agent-S-0751
stream: S-0751
tags: [cli]
touches: [flai/internal/check/image.go, flai/internal/check/image_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/operators/runbooks/agent-container.md]
after: [T-3178, T-3179]
---
# T-3180 flai check's agent.image names flai image ensure, and flai-cli.md, the user guide, and the runbook describe the image job

## Work

- `flai/internal/check/image.go`: the "not built" warning (lines 86 to 88) names `flai image ensure` and says `flai serve` builds the image at the next start. Update `TestAgentImage`.
- `design/system/flai-cli.md` § The container harness (S-0508), "The start in `flai serve`": the job, its one-at-a-time rule, the held run, its record in `agents.json`, the journal entries, and a failure's words. Mark the job built in the section S-0733 wrote. Bump `updated`.
- `docs/users/flai.md` § In a container, the error table (lines 539 to 540): the waiting run and the failure in place of the refusal.
- `docs/operators/runbooks/agent-container.md` § Build: images come back by themselves at the next start, where the job's log is, and when to run `flai image ensure` by hand.

It waits for T-3178 and T-3179, whose behaviour it describes.

## Done when

- `go test ./internal/check/` passes.
- The three documents describe the job as built, and none tells the operator to run `make agent-images` for a missing image.
- `flai check --strict` and the markdown lint pass.

## Notes
