---
id: S-0341
type: story
nature: improvement
title: flai verify resumes at the tier that failed when the branch head and its base are unchanged, so a retry re-runs the failure and not the tiers that passed
status: ready
owner: alex
created: 2026-10-08T07:59:12Z
updated: 2026-10-08T07:59:29Z
transitions:
  - to: ready
    at: 2026-10-08T07:59:29Z
    by: system-flow
tags: [cli]
topics: [testing]
touches: [flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/verify/run.go, flai/internal/verify/run_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, scripts/close-out.sh, docs/users/flai.md, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 22
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 103
          cache_read: 925166
          cache_write: 2463
          cost: 0.2286
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: alex
    at: 2026-10-08T07:59:12Z
---
# S-0341 flai verify resumes at the tier that failed when the branch head and its base are unchanged, so a retry re-runs the failure and not the tiers that passed

## Goal

`flai verify` runs every selected tier from the start on every call. When a late tier fails for a reason outside the branch, the agent fixes nothing and runs the close-out again, and pays the whole suite again before reaching the tier that failed. On 2026-10-08 S-0326's agent ran the close-out five times at the same head: each run passed rebase, sync, narrative, check, gofmt, vet, golangci-lint, go-test, markdown, and integration, about four minutes, and then smoke failed on a network fault. Twenty minutes of Go tests were spent to learn five times what the first run had shown.

A verify run at the same branch head, against the same main commit, with the same selected tiers as the last recorded run, can trust that run's passed tiers: nothing they read has changed. Such a run starts at the first tier the record shows as failed or not reached, after the cheap checks that guard the premise, and says which tiers it reused.

## Acceptance criteria

- [ ] `flai verify S-nnnn` reads the story's last record and, when the branch head, the main commit it was verified against, and the selected tiers are the same, re-runs the rebase, sync, narrative, and check steps and then only the tiers from the first one the record shows as failed or not reached; a tier it did not run is reported as `reused` with the time of the run it comes from, in the text and in `--json`.
- [ ] Any change to the head, the base, the selected tiers, the manifest's `tests`, or a tier's command makes the next run a full run, and `--fresh` forces one; the record written by a resumed run holds every tier's state, reused ones included, so `--last` and the review page show a whole result.
- [ ] `scripts/close-out.sh` runs `flai verify` as it does, so a second close-out at an unchanged head after a smoke failure reaches smoke within the cheap checks' time; the close-out's last line says when tiers were reused.
- [ ] Tests cover a resumed run, each condition that forces a full run, and `--fresh`.
- [ ] `docs/users/flai.md` and `design/system/flai-cli.md` describe when a run resumes and how to force a full one.

## Tasks

## Notes

The record is `.flai-cache/verify/S-nnnn.json`, which `flai verify --last` prints; it already carries the head and the main commit a run was verified against. A commit in the close-out between verify and the next run changes the head and so forces a full run, which is right: the tiers ran against other code.
