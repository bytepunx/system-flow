---
id: S-0267
type: story
nature: improvement
title: "The close-out's flai tier runs each Go test once: the short behaviour run is dropped where the full run follows it"
status: done
owner: alex
created: 2026-10-05T00:06:35Z
updated: 2026-10-05T00:25:39Z
transitions:
  - to: ready
    at: 2026-10-05T00:14:13Z
    by: alex
  - to: in-progress
    at: 2026-10-05T00:18:02Z
    by: agent-S-0267
  - to: review
    at: 2026-10-05T00:25:22Z
    by: agent-S-0267
  - to: done
    at: 2026-10-05T00:25:39Z
    by: alex
tags: []
topics: [code]
touches: [scripts/flai-test.sh, scripts/test.sh, scripts/flaiover-unit.sh, scripts/close-out.sh, scripts/README.md, Makefile, design/system/devex.md, docs/operators/index.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 464
  models:
    - model: claude-opus-5-5
      input: 98
      output: 22263
      cache_read: 3965318
      cache_write: 106060
      cost: 2.0872
    - model: claude-sonnet-5
      input: 36
      output: 10235
      cache_read: 603516
      cache_write: 56984
      cost: 0.3656
  strategic:
    - kind: planner
      seconds: 461
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 287
          output: 36124
          cache_read: 9704604
          cache_write: 244567
          cost: 5.59
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: planner-S-0267
    at: 2026-10-05T00:23:54Z
  value: 150
  by: planner-S-0267
  at: 2026-10-05T00:24:09Z
forecast:
  duration: 25m
  delivery: 2026-10-05T00:43:00Z
  basis: "flai forecast's 16m (87 s per unit over 6 done medium improvement stories, size 11) raised to 25m: verifying takes about 10m of test runs alone, flai-test.sh before and after plus make test and make integration; started 00:18Z"
  by: planner-S-0267
  at: 2026-10-05T00:19:15Z
---
# S-0267 The close-out's flai tier runs each Go test once: the short behaviour run is dropped where the full run follows it

## Goal

`scripts/flai-test.sh` runs `scripts/test.sh` (`go test -race -short -count=1 ./...` plus flaiover's vitest) and then `scripts/integration.sh` (`go test -race -count=1 ./...`) on the same packages. The full run is a superset of the short one, so every close-out, and every `make flai-test`, runs the short Go tests only to run them again. Timed in the main checkout on 2026-10-04: the short tier took 50 seconds and the full tier 67, so about 50 seconds of each close-out run is duplicate testing; the S-0248 verifier ran close-out three times and paid it three times. The tiers themselves stay as they are for day-to-day work, where `make test` is the quick run between tasks: it is only the sequence that runs both that should run the Go tests once.

## Acceptance criteria
- [x] `scripts/flai-test.sh`, and so `scripts/close-out.sh` and `make flai-test`, run the Go tests once with the race detector, the full run, and still run flaiover's unit tests when `flaiover/node_modules` is present
- [x] `make test` and `make integration` keep their current behaviour for a developer running one tier
- [x] `design/system/devex.md`, `scripts/README.md`, and `docs/operators/index.md` describe the tiers as they then run

## Tasks
- T-0842 scripts/flai-test.sh runs the Go tests once, the full run with the race detector, and still runs flaiover's unit tests
- T-0843 The design, scripts README, operator docs, Makefile help, and close-out header describe the tiers as flai-test.sh now runs them
- T-0846 flai-test.sh runs the Go tests once, the full run, and flaiover's unit tests through a script of their own
- T-0848 The design, the scripts index, and the operator guide describe the tiers as flai-test.sh runs them

## Notes

Found by the operator's review of the S-0248 agent log on 2026-10-04.

### Planning

Touches:

- `scripts/flai-test.sh`, `scripts/test.sh`: from the goal and the first criterion. `flai-test.sh` stops calling `test.sh`'s Go run, and `test.sh` gives up its vitest block so that `flai-test.sh` can run that block alone.
- `scripts/flaiover-unit.sh`: from the code layout. It is the predicted new home of the vitest block that `test.sh` and `flai-test.sh` would share. The story agent's T-0846 chose the same shape.
- `design/system/devex.md`, `scripts/README.md`, `docs/operators/index.md`: from the design, named by the third criterion.
- `Makefile`: from co-change, 9% with the named paths. Its `flai-test` help says "all three test tiers in order".
- `scripts/close-out.sh`: from the code layout. Its header comment says `flai-test.sh` runs "every tier". The comment changes and the behaviour does not.
- `flai touches suggest` listed nothing else that bears on the story. Its top co-changes, `flai-cli.md` and the dashboard design, change alongside `docs/operators/index.md` for unrelated reasons.
- The story declared no touches, so none was kept.
- No `.github` workflow calls these scripts.
- `template/root/scripts` has its own generic tiers and no `flai-test.sh`, so the template is not touched.

Topics: `code` is added, because `devex.md` carries it.

Forecast: 25m, delivery 2026-10-05T00:43Z.

- `flai forecast` gave 16m (87 s per unit over 6 done medium improvement stories, size 11).
- I raised it because checking the criteria needs about 10 minutes of test runs: `flai-test.sh` before and after the change, about 3 minutes each, plus `make test` and `make integration`, about 2 minutes together.
- The delivery counts 25 minutes from 00:18Z, when the story agent pulled it.

Cost of delay: 150 USD/week, the value `flai cod` gives, and it stands.

- The operator set `time_lost_per_cycle: 1h` in TH-0114 ("use a 1h CoD"), and I recorded it on that instruction.
- `flai cod` works it out as 1h lost per 168h cycle at an `hour_rate` of 150.
- The hour estimates the duplicate short Go runs a week: about 55 stories a week change `flai/`, about 1.5 close-out runs each, and about 50 seconds each.

Tasks: the story agent pulled the story at 00:18:02Z, while I was creating T-0842 and T-0843, and wrote T-0846 and T-0848 with the same split. It cancelled mine as duplicates, as TH-0114 proposed, and T-0848 took the `scripts/close-out.sh` header comment.
