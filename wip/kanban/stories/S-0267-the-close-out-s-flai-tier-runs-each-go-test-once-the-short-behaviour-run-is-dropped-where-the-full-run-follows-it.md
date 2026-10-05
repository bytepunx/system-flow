---
id: S-0267
type: story
nature: improvement
title: "The close-out's flai tier runs each Go test once: the short behaviour run is dropped where the full run follows it"
status: backlog
owner: alex
created: 2026-10-05T00:06:35Z
updated: 2026-10-05T00:06:35Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0267 The close-out's flai tier runs each Go test once: the short behaviour run is dropped where the full run follows it

## Goal

`scripts/flai-test.sh` runs `scripts/test.sh` (`go test -race -short -count=1 ./...` plus flaiover's vitest) and then `scripts/integration.sh` (`go test -race -count=1 ./...`) on the same packages. The full run is a superset of the short one, so every close-out, and every `make flai-test`, runs the short Go tests only to run them again. Timed in the main checkout on 2026-10-04: the short tier took 50 seconds and the full tier 67, so about 50 seconds of each close-out run is duplicate testing; the S-0248 verifier ran close-out three times and paid it three times. The tiers themselves stay as they are for day-to-day work, where `make test` is the quick run between tasks: it is only the sequence that runs both that should run the Go tests once.

## Acceptance criteria
- [ ] `scripts/flai-test.sh`, and so `scripts/close-out.sh` and `make flai-test`, run the Go tests once with the race detector, the full run, and still run flaiover's unit tests when `flaiover/node_modules` is present
- [ ] `make test` and `make integration` keep their current behaviour for a developer running one tier
- [ ] `design/system/devex.md`, `scripts/README.md`, and `docs/operators/index.md` describe the tiers as they then run

## Tasks

## Notes

Found by the operator's review of the S-0248 agent log on 2026-10-04.
